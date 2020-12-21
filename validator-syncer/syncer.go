package validator_syncer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/chain/config"
	"github.com/ChainSafe/chainbridge-celo/cmd/cfg"
	"github.com/ChainSafe/chainbridge-utils/crypto/secp256k1"
	"github.com/ChainSafe/chainbridge-utils/keystore"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func Sync(ctx *cli.Context) error {
	startConfig, err := cfg.GetConfig(ctx)
	if err != nil {
		return err
	}
	errChn := make(chan error)
	stopChn := make(chan struct{})
	dataChn := make(chan string)
	celoChainConfig, err := config.ParseChainConfig(&startConfig.Chains[0], ctx)

	kpI, err := keystore.KeypairFromAddress(celoChainConfig.From, keystore.EthChain, celoChainConfig.KeystorePath, celoChainConfig.Insecure)
	if err != nil {
		return err
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	chainClient, err := client.NewClient(celoChainConfig.Endpoint, celoChainConfig.Http, kp, celoChainConfig.GasLimit, celoChainConfig.MaxGasPrice)
	//syncer := validator.NewValidatorSyncer(chainClient)
	db := NewValidatorsDB()

	go SyncValidators(stopChn, errChn, chainClient, db)

	sysErr := make(chan os.Signal, 1)
	signal.Notify(sysErr,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGQUIT)

	for {
		select {
		case <-stopChn:
			return errors.New("stoped")
		case err := <-errChn:
			close(stopChn)
			return err
		case _ = <-dataChn:
			//fmt.Println(res)
			continue
		case <-sysErr:
			close(stopChn)
			return nil
		}
	}
}

type SyncerDB interface {
	GetLatestValidators() []*istanbul.ValidatorData
	GetLatestBlock() *big.Int
	SetValidatorsForBlock(block *big.Int, validators []*istanbul.ValidatorData)
}

func SyncValidators(stopChn <-chan struct{}, errChn chan error, c *client.Client, db SyncerDB) {
	block := db.GetLatestBlock()
	actualValidators := db.GetLatestValidators()
	for {
		select {
		case <-stopChn:
			return
		default:
			header, err := c.HeaderByNumber(context.Background(), block)
			if err != nil {
				if errors.Is(err, ethereum.NotFound) {
					// Block not yet mined, waiting
					time.Sleep(5)
					continue
				}
				errChn <- fmt.Errorf("gettings header by number err: %w", err)
				return
			}
			extra, err := types.ExtractIstanbulExtra(header)
			b := bytes.NewBuffer(extra.RemovedValidators.Bytes())
			if len(extra.AddedValidators) != 0 || b.Len() > 0 {
				actualValidators, err = ApplyValidatorsDiff(extra, db.GetLatestValidators())
				log.Debug().Msgf("New validators %+v", actualValidators)
			}
			db.SetValidatorsForBlock(block, actualValidators)
			block.Add(block, big.NewInt(1))
		}
	}
}
