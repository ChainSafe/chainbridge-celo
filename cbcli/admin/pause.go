package admin

import (
	"github.com/ChainSafe/chainbridge-celo/cbcli/cliutils"
	"github.com/rs/zerolog/log"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/urfave/cli/v2"
)

var pauseCMD = &cli.Command{
	Name:        "pause",
	Description: "Pauses deposits and proposals.",
	Action:      pause,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "bridge",
			Usage: "Bridge contract address",
		},
	},
}

func pause(cctx *cli.Context) error {
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
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice))
	if err != nil {
		return err
	}
	err = utils.AdminPause(ethClient, bridgeAddress)
	if err != nil {
		return err
	}
	log.Info().Msgf("Deposits and proposals are paused")
	return nil
}
