package bridge

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

func setBurn(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Int64("gasLimit")
	gasPrice := cctx.Int64("gasPrice")
	bridgeAddress, err := cliutils.DefineBridgeAddress(cctx)
	if err != nil {
		return err
	}
	handler := cctx.String("handler")
	if !common.IsHexAddress(handler) {
		return errors.New("handler address is incorrect format")
	}
	tokenContract := cctx.String("tokenContract")
	if !common.IsHexAddress(tokenContract) {
		return errors.New("tokenContract address is incorrect format")
	}
	handlerAddress := common.HexToAddress(handler)
	tokenContractAddress := common.HexToAddress(tokenContract)

	sender, err := cliutils.DefineSender(cctx)
	if err != nil {
		return err
	}

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(gasLimit), big.NewInt(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	log.Info().Msgf("Setting contract %s as burnable on handler %s", tokenContractAddress.String(), handlerAddress.String())
	err = utils.SetBurnable(ethClient, bridgeAddress, handlerAddress, tokenContractAddress)
	if err != nil {
		return err
	}
	log.Info().Msg("Burnable set")
	return nil
}

var setBurnCMD = &cli.Command{
	Name:        "set-burn",
	Description: "Set a token contract as mintable/burnable in a handler.",
	Action:      setBurn,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "bridge",
			Usage: "Bridge contract address",
		},
		&cli.StringFlag{
			Name:  "handler",
			Usage: "ERC20 handler contract address",
		},
		&cli.StringFlag{
			Name:  "tokenContract",
			Usage: "Token contract to be registered",
		},
	},
}
