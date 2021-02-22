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

var setTresholdCMD = &cli.Command{
	Name:        "set-treshold",
	Description: "Sets a new relayer vote threshold.",
	Action:      setTreshold,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "treshold",
			Usage: "New relayer threshold",
		},
		&cli.StringFlag{
			Name:  "bridge",
			Usage: "Bridge contract address",
		},
	},
}

func setTreshold(cctx *cli.Context) error {
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
	treshHold := cctx.Uint64("treshold")
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice))
	if err != nil {
		return err
	}
	err = utils.AdminSetTreshHold(ethClient, bridgeAddress, big.NewInt(0).SetUint64(treshHold))
	if err != nil {
		return err
	}
	log.Info().Msgf("New treshhold set for %v", treshHold)
	return nil
}
