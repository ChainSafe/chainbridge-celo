package bridge

import (
	"errors"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/cbcli/cliutils"
	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/celo-org/celo-blockchain/common"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func queryResource(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Uint64("gasLimit")
	gasPrice := cctx.Uint64("gasPrice")
	sender, err := cliutils.DefineSender(cctx)
	if err != nil {
		return err
	}
	handlerS := cctx.String("handler")
	if !common.IsHexAddress(handlerS) {
		return errors.New("provided handler address is not valid")
	}
	handlerAddr := common.HexToAddress(handlerS)
	resourceIDs := cctx.String("resourceId")
	resourceID := utils.SliceTo32Bytes(common.Hex2Bytes(resourceIDs))
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	res, err := utils.QueryResource(ethClient, handlerAddr, resourceID)
	if err != nil {
		return err
	}
	log.Info().Msgf("Resource address that associated with ID %s is %s", common.Bytes2Hex(resourceID[:]), res.String())
	return nil
}

var queryResourceCMD = &cli.Command{
	Name:        "query-resource",
	Description: "Queries the contract address associated with the provided resource ID for a specific handler contract.",
	Action:      queryResource,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "handler",
			Usage: "Handler contract address",
		},
		&cli.StringFlag{
			Name:  "resourceId",
			Usage: "ResourceID to query",
		},
	},
}
