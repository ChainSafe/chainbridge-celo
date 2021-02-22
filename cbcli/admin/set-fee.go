package admin

import (
	"github.com/ChainSafe/chainbridge-celo/cbcli/cliutils"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var setFeeCMD = &cli.Command{
	Name:        "is-relayer",
	Description: "Sets a new relayer vote threshold.",
	Action:      setFee,
	Flags: []cli.Flag{
		&cli.Uint64Flag{
			Name:  "fee",
			Usage: "New fee (in wei)",
		},
		&cli.StringFlag{
			Name:  "bridge",
			Usage: "Bridge contract address",
		},
	},
}

func setFee(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Uint64("gasLimit")
	gasPrice := cctx.Uint64("gasPrice")
	sender, err := cliutils.DefineSender(cctx)
	if err != nil {
		return err
	}
	bridge := cctx.String("bridge")
	if !common.IsHexAddress(bridge) {
		return errors.New("invalid bridge address")
	}
	bridgeAddress := common.HexToAddress(bridge)
	fee := cctx.Uint64("fee")
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice))
	if err != nil {
		return err
	}
	err = utils.AdminSetTreshHold(ethClient, bridgeAddress, big.NewInt(0).SetUint64(fee))
	if err != nil {
		return err
	}
	log.Info().Msgf("Fee set to %v", fee)
	return nil
}
