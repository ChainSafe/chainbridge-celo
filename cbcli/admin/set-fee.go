package admin

import (
	"github.com/ChainSafe/chainbridge-celo/cbcli/cliutils"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var setFeeCMD = &cli.Command{
	Name:        "set-fee",
	Description: "Sets a new fee for deposits.",
	Action:      setFee,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "fee",
			Usage: "New fee (in ether)",
		},
		&cli.StringFlag{
			Name:  "bridge",
			Usage: "Bridge contract address",
		},
		&cli.Uint64Flag{
			Name:     "decimals",
			Usage:    "erc20Token decimals",
			Required: true,
		},
	},
}

func setFee(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Uint64("gasLimit")
	gasPrice := cctx.Uint64("gasPrice")
	decimals := big.NewInt(0).SetUint64(cctx.Uint64("decimals"))
	sender, err := cliutils.DefineSender(cctx)
	if err != nil {
		return err
	}
	bridgeAddress, err := cliutils.DefineBridgeAddress(cctx)
	if err != nil {
		return err
	}
	fee := cctx.String("fee")

	realFeeAmount, err := utils.UserAmountToWei(fee, decimals)
	if err != nil {
		return err
	}
	log.Debug().Msgf("%s", realFeeAmount.String())

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice))
	if err != nil {
		return err
	}
	err = utils.AdminSetFee(ethClient, bridgeAddress, realFeeAmount)
	if err != nil {
		return err
	}
	log.Info().Msgf("Fee set to %v", fee)
	return nil
}
