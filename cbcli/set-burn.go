package cbcli

import (
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func SetBurn(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Int64("gasLimit")
	gasPrice := cctx.Int64("gasPrice")
	bridge := cctx.String("bridge")
	if !common.IsHexAddress(bridge) {
		return errors.New("bridge address is incorrect format")
	}
	handler := cctx.String("handler")
	if !common.IsHexAddress(handler) {
		return errors.New("handler address is incorrect format")
	}
	tokenContract := cctx.String("tokenContract")
	if !common.IsHexAddress(tokenContract) {
		return errors.New("tokenContract address is incorrect format")
	}
	bridgeAddress := common.HexToAddress(bridge)
	handlerAddress := common.HexToAddress(handler)
	tokenContractAddress := common.HexToAddress(tokenContract)

	sender, err := defineSender(cctx)
	if err != nil {
		return err
	}

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(gasLimit), big.NewInt(gasPrice))
	if err != nil {
		return err
	}
	log.Info().Msgf("Setting contract %s as burnable on handler %s", tokenContractAddress, handlerAddress)
	err = utils.SetBurnable(ethClient, bridgeAddress, handlerAddress, tokenContractAddress)
	if err != nil {
		return err
	}
	log.Info().Msg("Burnable set")
	return nil
}
