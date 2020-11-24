package cmd

import (
	"strconv"

	"github.com/ChainSafe/chainbridge-celo/blockdb"
	"github.com/ChainSafe/chainbridge-celo/chain"
	"github.com/ChainSafe/chainbridge-celo/chain/connection"
	"github.com/ChainSafe/chainbridge-celo/chain/listener"
	"github.com/ChainSafe/chainbridge-celo/chain/writer"
	"github.com/ChainSafe/chainbridge-celo/cmd/cfg"
	"github.com/ChainSafe/chainbridge-celo/flags"
	"github.com/ChainSafe/chainbridge-utils/core"
	"github.com/ChainSafe/chainbridge-utils/msg"
	"github.com/urfave/cli/v2"
)

func Run(ctx *cli.Context) error {
	startConfig, err := cfg.GetConfig(ctx)
	if err != nil {
		return err
	}

	// Check for test key flag
	var ks string
	var insecure bool
	if key := ctx.String(flags.TestKeyFlag.Name); key != "" {
		ks = key
		insecure = true
	} else {
		ks = startConfig.KeystorePath
	}

	sysErr := make(chan error)
	coreApp := core.NewCore(sysErr)

	for _, c := range startConfig.Chains {
		chainId, errr := strconv.Atoi(c.Id)
		if errr != nil {
			return errr
		}
		chainConfig := &core.ChainConfig{
			Name:           c.Name,
			Id:             msg.ChainId(chainId),
			Endpoint:       c.Endpoint,
			From:           c.From,
			KeystorePath:   ks,
			Insecure:       insecure,
			BlockstorePath: ctx.String(flags.BlockstorePathFlag.Name),
			FreshStart:     ctx.Bool(flags.FreshStartFlag.Name),
			LatestBlock:    ctx.Bool(flags.LatestBlockFlag.Name),
			Opts:           c.Opts,
		}
		celoChainConfig, err := chain.ParseChainConfig(chainConfig)
		if err != nil {
			return err
		}

		var newChain core.Chain

		conn := connection.NewConnection(celoChainConfig.Endpoint, celoChainConfig.Http, kp, celoChainConfig.GasLimit, celoChainConfig.MaxGasPrice)
		bdb := blockdb.NewBlockstoreDB(from, keystorePath, insecure, bloksorePath, chainID, freshStart, startblock)
		l := listener.NewListener(conn, cfg, bs, stop, sysErr, m)
		w := writer.NewWriter(conn, cfg, logger, stop, sysErr, m)

		newChain, err = chain.InitializeChain(celoChainConfig, sysErr, conn, l, w, bdb)

		if err != nil {
			return err
		}
		coreApp.AddChain(newChain)

	}

	return nil
}
