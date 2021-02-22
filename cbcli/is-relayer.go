package cbcli

import (
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/rs/zerolog/log"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
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
	sender, err := defineSender(cctx)
	if err != nil {
		return err
	}
	bridge := cctx.String("bridge")
	if !common.IsHexAddress(bridge) {
		return errors.New("invalid bridge address")
	}
	bridgeAddress := common.HexToAddress(bridge)
	relayer := cctx.String("relayer")
	if !common.IsHexAddress(bridge) {
		return errors.New("invalid bridge address")
	}
	relayerAddress := common.HexToAddress(relayer)
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice))
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
		log.Info().Msgf("Requested address %s is not relayer", relayerAddress.String())
	}
	return nil
}
