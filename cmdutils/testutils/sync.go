//nolint
//Copyright 2020 ChainSafe Systems
//SPDX-License-Identifier: LGPL-3.0-only
package testutils

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/ChainSafe/chainbridge-celo/chain/config"
	"github.com/ChainSafe/chainbridge-celo/chain/sender"
	"github.com/ChainSafe/chainbridge-celo/cmd/cfg"
	"github.com/ChainSafe/chainbridge-celo/validatorsync"
	"github.com/ChainSafe/chainbridge-utils/crypto/secp256k1"
	"github.com/ChainSafe/chainbridge-utils/keystore"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/urfave/cli/v2"
)

// SYnc is the function that only runs ValidatorsSyncer functionality of chainbridge for test purposes
func Sync(ctx *cli.Context) error {
	startConfig, err := cfg.GetConfig(ctx)
	if err != nil {
		return err
	}
	errChn := make(chan error)
	stopChn := make(chan struct{})
	celoChainConfig, err := config.ParseChainConfig(&startConfig.Chains[0], ctx)
	if err != nil {
		return err
	}
	kpI, err := keystore.KeypairFromAddress(celoChainConfig.From, keystore.EthChain, celoChainConfig.KeystorePath, celoChainConfig.Insecure)
	if err != nil {
		return err
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	chainClient, err := sender.NewSender(celoChainConfig.Endpoint, celoChainConfig.Http, kp, celoChainConfig.GasLimit, celoChainConfig.MaxGasPrice)
	if err != nil {
		return err
	}
	db, err := leveldb.OpenFile("./test/db", nil)
	if err != nil {
		return err
	}
	store := validatorsync.NewValidatorsStore(db)

	go validatorsync.SyncBlockValidators(stopChn, errChn, chainClient, store, 1, 12)

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
		case <-sysErr:
			close(stopChn)
			return nil
		}
	}
}

func storeBlockHeaderFile(h *types.Header) {
	f, err := os.Create("file.bin")
	if err != nil {
		log.Fatal().Err(err).Msg("can't open file")
	}
	defer f.Close()

	// TODO
	//err = binary.Write(f, , *h)
	if err != nil {
		log.Fatal().Err(err).Msg("write fail")

	}
}
