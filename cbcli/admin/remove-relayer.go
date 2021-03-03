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
	bridge := cctx.String("bridge")
	if !common.IsHexAddress(bridge) {
		return errors.New(fmt.Sprintf("invalid bridge address %s", bridge))
	}
	bridgeAddress := common.HexToAddress(bridge)
	relayer := cctx.String("relayer")
	if !common.IsHexAddress(bridge) {
		return errors.New(fmt.Sprintf("invalid bridge address %s", relayer))
	}
	relayerAddress := common.HexToAddress(relayer)
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice))
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
