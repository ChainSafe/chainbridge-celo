package cmd

import (
	"github.com/ChainSafe/chainbridge-utils/crypto/secp256k1"
	"github.com/ChainSafe/chainbridge-utils/keystore"

	"github.com/ChainSafe/chainbridge-celo/blockdb"
	"github.com/ChainSafe/chainbridge-celo/chain"
	"github.com/ChainSafe/chainbridge-celo/chain/connection"
	"github.com/ChainSafe/chainbridge-celo/chain/listener"
	"github.com/ChainSafe/chainbridge-celo/chain/writer"
	"github.com/ChainSafe/chainbridge-celo/cmd/cfg"
	"github.com/ChainSafe/chainbridge-celo/flags"
	"github.com/ChainSafe/chainbridge-utils/core"
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
		celoChainConfig, err := chain.ParseChainConfig(chainConfig)
		if err != nil {
			return err
		}

		kpI, err := keystore.KeypairFromAddress(celoChainConfig.From, keystore.EthChain, celoChainConfig.KeystorePath, celoChainConfig.Insecure)
		if err != nil {
			return err
		}
		kp, _ := kpI.(*secp256k1.Keypair)

		conn := connection.NewConnection(celoChainConfig.Endpoint, celoChainConfig.Http, kp, celoChainConfig.GasLimit, celoChainConfig.MaxGasPrice)
		bdb, err := blockdb.NewBlockstoreDB(from, keystorePath, insecure, bloksorePath, chainID, freshStart, startblock)
		if err != nil {
			return err
		}
		l := listener.NewListener(conn, cfg, bs, stop, sysErr, m)
		w := writer.NewWriter(conn, cfg, logger, stop, sysErr, m)

		var newChain core.Chain
		newChain, err = chain.InitializeChain(celoChainConfig, sysErr, conn, l, w, bdb)

		if err != nil {
			return err
		}
		coreApp.AddChain(newChain)

	}

	return nil
}
