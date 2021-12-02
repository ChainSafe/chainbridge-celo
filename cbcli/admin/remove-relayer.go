package admin

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/cbcli/cliutils"

	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/celo-org/celo-blockchain/common"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var removeRelayerCMD = &cli.Command{
	Name:        "remove-relayer",
	Description: "Removes a relayer.",
	Action:      removeRelayer,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "relayer",
			Usage: "Address to remove",
		},
		&cli.StringFlag{
			Name:  "bridge",
			Usage: "Bridge contract address",
		},
	},
}

func removeRelayer(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Uint64("gasLimit")
	gasPrice := cctx.Uint64("gasPrice")
	sender, err := cliutils.DefineSender(cctx)
	if err != nil {
		return err
	}
	bridgeAddress, err := cliutils.DefineBridgeAddress(cctx)
	if err != nil {
		return err
	}
	relayer := cctx.String("relayer")
	if !common.IsHexAddress(relayer) {
		return fmt.Errorf("invalid bridge address %s", relayer)
	}
	relayerAddress := common.HexToAddress(relayer)
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	err = utils.AdminRemoveRelayer(ethClient, bridgeAddress, relayerAddress)
	if err != nil {
		return err
	}
	log.Info().Msgf("Address %s is relayer now", relayerAddress.String())
	return nil
}
