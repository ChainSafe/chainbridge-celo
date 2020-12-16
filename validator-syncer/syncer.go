package validator_syncer

import (
	"fmt"
	"github.com/ChainSafe/chainbridge-celo/chain"
	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/cmd/cfg"
	"github.com/ChainSafe/chainbridge-utils/crypto/secp256k1"
	"github.com/ChainSafe/chainbridge-utils/keystore"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"syscall"
)

func Sync(ctx *cli.Context) error {
	startConfig, err := cfg.GetConfig(ctx)
	if err != nil {
		return err
	}
	errChn := make(chan error)
	stopChn := make(chan struct{})
	dataChn := make(chan string)
	celoChainConfig, err := chain.ParseChainConfig(&startConfig.Chains[0], ctx)

	kpI, err := keystore.KeypairFromAddress(celoChainConfig.From, keystore.EthChain, celoChainConfig.KeystorePath, celoChainConfig.Insecure)
	if err != nil {
		return err
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	chainClient, err := client.NewClient(celoChainConfig.Endpoint, celoChainConfig.Http, kp, celoChainConfig.GasLimit, celoChainConfig.MaxGasPrice)

	go validate(stopChn, dataChn, chainClient)

	sysErr := make(chan os.Signal, 1)
	signal.Notify(sysErr,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGQUIT)

	for {
		select {
		case <-stopChn:
			return nil
		case err := <-errChn:
			return err
		case res := <-dataChn:
			fmt.Println(res)
			continue
		case <-sysErr:
			return nil
		}
	}

}

func validate(stopChn <-chan struct{}, dataChn chan string, c *client.Client) {

}
