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
	"github.com/ChainSafe/chainbridge-celo/core"
	"github.com/urfave/cli/v2"
)

func Run(ctx *cli.Context) error {
	startConfig, err := cfg.GetConfig(ctx)
	if err != nil {
		return err
	}
	sysErr := make(chan error)
	coreApp := core.NewCore(sysErr)
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

		conn := connection.NewConnection(celoChainConfig.Endpoint, celoChainConfig.Http, kp, celoChainConfig.GasLimit, celoChainConfig.MaxGasPrice)
		err = celoChainConfig.EnsureContractsHaveBytecode(conn)
		if err != nil {
			return err
		}
		bdb, err := blockdb.NewBlockStoreDB(kp.Address(), celoChainConfig.BlockstorePath, celoChainConfig.ID, celoChainConfig.FreshStart, celoChainConfig.StartBlock)
		if err != nil {
			return err
		}

		stop := make(chan int)
		l := listener.NewListener(conn, celoChainConfig, bdb, stop, sysErr, validatorSyncer)
		// TODO ChainMetrics
		w := writer.NewWriter(conn, celoChainConfig, stop, sysErr, nil)

		newChain, err := chain.InitializeChain(celoChainConfig, sysErr, conn, l, w, bdb)

		if err != nil {
			return err
		}
		coreApp.AddChain(newChain)
	}
	return nil
}
