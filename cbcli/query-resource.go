package cbcli

import (
	"errors"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func QueryResource(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Uint64("gasLimit")
	gasPrice := cctx.Uint64("gasPrice")
	sender, err := defineSender(cctx)
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
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice))
	if err != nil {
		return err
	}
	res, err := utils.QueryResource(ethClient, handlerAddr, resourceID)
	if err != nil {
		return err
	}
	log.Info().Msgf("Resource address that associated with ID %s is %s", resourceID, res)
	return nil
}
