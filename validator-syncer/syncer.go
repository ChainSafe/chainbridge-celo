package validator_syncer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
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

	go validate(stopChn, errChn, dataChn, chainClient, db)

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

func validate(stopChn <-chan struct{}, errChn chan error, dataChn chan string, c *client.Client, db *validatorsDB) {
	block := db.getLatestAddedBlock()
	actualValidators := db.getLatestValidators()
	for {
		select {
		case <-stopChn:
			return
		default:
			header, err := c.HeaderByNumber(context.Background(), block)
			if err != nil {
				if errors.Is(err, ethereum.NotFound) {
					time.Sleep(5)
					continue
				}
				errChn <- fmt.Errorf("gettings header by number err: %w", err)
				return
			}
			extra, err := types.ExtractIstanbulExtra(header)
			b := bytes.NewBuffer(extra.RemovedValidators.Bytes())
			if len(extra.AddedValidators) != 0 || b.Len() > 0 {
				log.Debug().Msgf("EXTRA OF BLOCK %+v %s", extra, block.String())
				actualValidators = ApplyValidatorsDiff(extra, db.latestValidators)
				log.Debug().Msgf("New validators %+v", actualValidators)
			}
			db.setValidatorsForBlock(block, actualValidators)
			block.Add(block, big.NewInt(1))
		}
	}
}
