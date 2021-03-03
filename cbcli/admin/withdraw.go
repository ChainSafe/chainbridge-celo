package admin

import (
	"fmt"
	"github.com/ChainSafe/chainbridge-celo/cbcli/cliutils"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var withdrawCMD = &cli.Command{
	Name:        "withdraw",
	Description: "Withdraw tokens from a handler contract.",
	Action:      withdraw,
	Flags: []cli.Flag{
		&cli.Uint64Flag{
			Name:  "amountOrId",
			Usage: "Token ID or amount to withdraw",
		},
		&cli.StringFlag{
			Name:  "bridge",
			Usage: "Bridge contract address",
		},
		&cli.StringFlag{
			Name:  "handler",
			Usage: "Handler contract address",
		},
		&cli.StringFlag{
			Name:  "token",
			Usage: "ERC20 or ERC721 token contract address",
		},
		&cli.StringFlag{
			Name:  "recipient",
			Usage: "Address to withdraw to",
		},
	},
}

func withdraw(cctx *cli.Context) error {
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

	handler := cctx.String("handler")
	if !common.IsHexAddress(handler) {
		return errors.New(fmt.Sprintf("invalid handler address %s", handler))
	}
	handlerAddress := common.HexToAddress(handler)

	token := cctx.String("token")
	if !common.IsHexAddress(token) {
		return errors.New(fmt.Sprintf("invalid token address %s", token))
	}
	tokenAddress := common.HexToAddress(token)

	recipient := cctx.String("recipient")
	if !common.IsHexAddress(recipient) {
		return errors.New(fmt.Sprintf("invalid recipient address %s", recipient))
	}
	recipientAddress := common.HexToAddress(recipient)

	amount := big.NewInt(0).SetUint64(cctx.Uint64("fee"))

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice))
	if err != nil {
		return err
	}
	err = utils.AdminWithdraw(ethClient, bridgeAddress, handlerAddress, tokenAddress, recipientAddress, amount)
	if err != nil {
		return err
	}
	log.Info().Msgf("Withdrawn %s to %s", amount.String(), recipient)
	return nil
}
