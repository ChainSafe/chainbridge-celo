package admin

import (
	"github.com/ChainSafe/chainbridge-celo/cbcli/cliutils"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var setTresholdCMD = &cli.Command{
	Name:        "set-threshold",
	Description: "Sets a new relayer vote threshold.",
	Action:      setThreshold,
	Flags: []cli.Flag{
		&cli.Uint64Flag{
			Name:  "threshold",
			Usage: "New relayer threshold",
		},
		&cli.StringFlag{
			Name:  "bridge",
			Usage: "Bridge contract address",
		},
	},
}

func setThreshold(cctx *cli.Context) error {
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
	threshold := cctx.Uint64("threshold")
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	err = utils.AdminSetThreshHold(ethClient, bridgeAddress, big.NewInt(0).SetUint64(threshold))
	if err != nil {
		return err
	}
	log.Info().Msgf("New threshold set for %v", threshold)
	return nil
}
