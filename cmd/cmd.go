package cmd

import (
	"github.com/ChainSafe/ChainBridge/core"
	"github.com/ChainSafe/chainbridge-celo/blockdb"
	"github.com/ChainSafe/chainbridge-celo/chain"
	"github.com/ChainSafe/chainbridge-celo/chain/connection"
	"github.com/ChainSafe/chainbridge-celo/chain/listener"
	"github.com/ChainSafe/chainbridge-celo/config"
	"github.com/urfave/cli/v2"
)

func Run(ctx *cli.Context) error {
	// TODO - Implement run method

	cfg, err := config.GetConfig(ctx)
	if err != nil {
		return err
	}

	sysErr := make(chan error)
	c := core.NewCore(sysErr)

	conn := connection.NewConnection(cfg.endpoint, cfg.http, kp, cfg.gasLimit, cfg.maxGasPrice)

	err = conn.Connect()
	if err != nil {
		return nil, err
	}

	bdb := blockdb.NewBlockstoreDB(from, keystorePath, insecure, bloksorePath, chainID, freshStart, startblock)

	listener := listener.NewListener(conn, cfg, bs, stop, sysErr, m)

	writer := NewWriter(conn, cfg, logger, stop, sysErr, m)

	newChain := chain.InitializeChain()
	return nil
}
