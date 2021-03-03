package admin

import (
	"fmt"
	"github.com/ChainSafe/chainbridge-celo/cbcli/cliutils"
	"github.com/rs/zerolog/log"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var unpauseCMD = &cli.Command{
	Name:        "unpause",
	Description: "Unpauses deposits and proposals.",
	Action:      unpause,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "bridge",
			Usage: "Bridge contract address",
		},
	},
}

func unpause(cctx *cli.Context) error {
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
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice))
	if err != nil {
		return err
	}
	err = utils.AdminUnpause(ethClient, bridgeAddress)
	if err != nil {
		return err
	}
	log.Info().Msgf("Deposits and proposals are Unpaused")
	return nil
}
