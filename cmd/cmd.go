// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only
package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ChainSafe/chainbridge-celo/router"
	"github.com/ChainSafe/chainbridge-utils/crypto/secp256k1"
	"github.com/ChainSafe/chainbridge-utils/keystore"

	"github.com/ChainSafe/chainbridge-celo/blockdb"
	"github.com/ChainSafe/chainbridge-celo/chain"
	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/chain/listener"
	"github.com/ChainSafe/chainbridge-celo/chain/writer"
	"github.com/ChainSafe/chainbridge-celo/cmd/cfg"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func Run(ctx *cli.Context) error {
	startConfig, err := cfg.GetConfig(ctx)
	if err != nil {
		return err
	}
	errChn := make(chan error)
	stopChn := make(chan struct{})
	r := router.NewRouter()
	for _, c := range startConfig.Chains {
		celoChainConfig, err := chain.ParseChainConfig(&c, ctx)
		if err != nil {
			return err
		}
		kpI, err := keystore.KeypairFromAddress(celoChainConfig.From, keystore.EthChain, celoChainConfig.KeystorePath, celoChainConfig.Insecure)
		if err != nil {
			return err
		}
		kp, _ := kpI.(*secp256k1.Keypair)

		chainClient, err := client.NewClient(celoChainConfig.Endpoint, celoChainConfig.Http, kp, celoChainConfig.GasLimit, celoChainConfig.MaxGasPrice)
		if err != nil {
			return err
		}
		// TODO not to abstract should be moved inside chain initialization
		bdb, err := blockdb.NewBlockStoreDB(kp.Address(), celoChainConfig.BlockstorePath, celoChainConfig.ID, celoChainConfig.FreshStart, celoChainConfig.StartBlock)
		if err != nil {
			return err
		}
		// TODO  ValidatorSyncer
		// TODO ChainMetrics
		w := writer.NewWriter(chainClient, celoChainConfig, stopChn, errChn, nil)
		r.Register(celoChainConfig.ID, w)
		l := listener.NewListener(celoChainConfig, chainClient, bdb, stopChn, errChn, nil, r)
		newChain, err := chain.InitializeChain(celoChainConfig, chainClient, l, w, stopChn)
		if err != nil {
			return err
		}
		err = newChain.Start()
		if err != nil {
			log.Error().Interface("chain", newChain.ID()).Err(err).Msg("failed to start chain")
			return err
		}

	}
	sysErr := make(chan os.Signal, 1)
	signal.Notify(sysErr,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGQUIT)

	select {
	case err := <-errChn:
		log.Error().Err(err).Msg("failed to listen and serve")
		close(stopChn)
		return err
	case sig := <-sysErr:
		log.Info().Msgf("terminating got [%v] signal", sig)
		return nil
	}
}
