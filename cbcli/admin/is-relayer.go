package admin

import (
	"fmt"
	"github.com/ChainSafe/chainbridge-celo/cbcli/cliutils"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var isRelayerCMD = &cli.Command{
	Name:        "is-relayer",
	Description: "Check if an address is registered as a relayer",
	Action:      isRelayer,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "relayer",
			Usage: "Address to check",
		},
		&cli.StringFlag{
			Name:  "bridge",
			Usage: "Bridge contract address",
		},
	},
}

func isRelayer(cctx *cli.Context) error {
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
		return errors.New(fmt.Sprintf("invalid relayer address %s", relayer))
	}
	relayerAddress := common.HexToAddress(relayer)
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	isRelayer, err := utils.AdminIsRelayer(ethClient, bridgeAddress, relayerAddress)
	if err != nil {
		return err
	}
	if isRelayer {
		log.Info().Msgf("Requested address %s is relayer", relayerAddress.String())
	} else {
		log.Info().Msgf("Requested address %s is not a relayer", relayerAddress.String())
	}
	return nil
}
