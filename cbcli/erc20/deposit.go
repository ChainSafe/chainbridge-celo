package erc20

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/cbcli/cliutils"
	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var depositCMD = &cli.Command{
	Name:        "deposit",
	Description: "Initiate a transfer of ERC20 tokens.",
	Action:      deposit,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "recipient",
			Usage: "Recipient",
		},
		&cli.StringFlag{
			Name:  "bridge",
			Usage: "bridge contract address",
		},
		&cli.StringFlag{
			Name:  "amount",
			Usage: "Amount to deposit",
		},
		&cli.StringFlag{
			Name:  "value",
			Usage: "Value of ETH that should be sent along with deposit to cover possible fees. In WEI",
			Value: "0",
		},
		&cli.StringFlag{
			Name:  "dest",
			Usage: "Destination chainID",
		},
		&cli.Uint64Flag{
			Name:  "resourceId",
			Usage: "ResourceID for transfer",
		},
		&cli.Uint64Flag{
			Name:     "decimals",
			Usage:    "erc20Token decimals",
			Required: true,
		},
	},
}

func deposit(cctx *cli.Context) error {
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

	recipient := cctx.String("recipient")
	if !common.IsHexAddress(recipient) {
		return errors.New(fmt.Sprintf("invalid recipient address %s", recipient))
	}
	recipientAddress := common.HexToAddress(recipient)

	amount := cctx.String("amount")

	realAmount, err := utils.UserAmountToWei(amount, decimals)
	if err != nil {
		return err
	}

	value, ok := big.NewInt(0).SetString(cctx.String("value"), 10)
	if !ok {
		return errors.New("invalid value format")
	}

	dest := cctx.Uint64("dest")

	resourceId := cctx.String("resourceId")
	resourceIDBytes := utils.SliceTo32Bytes(common.Hex2Bytes(resourceId))

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice))
	if err != nil {
		return err
	}

	ethClient.ClientWithArgs(client.TheClientWithValue(value))

	err = utils.MakeAndSendERC20Deposit(ethClient, bridgeAddress, recipientAddress, realAmount, resourceIDBytes, uint8(dest))
	if err != nil {
		return err
	}
	log.Info().Msgf("%s tokens were transferred to %s from %s", amount, recipientAddress.String(), sender.CommonAddress().String())
	return nil
}
