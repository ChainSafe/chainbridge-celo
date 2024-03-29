package erc721

import (
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
	Description: "Initiates a bridge ERC721 transfer.",
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
		&cli.Int64Flag{
			Name:  "id",
			Usage: "ERC721 token id",
		},
		&cli.StringFlag{
			Name:  "dest",
			Usage: "Destination chainID",
		},
		&cli.StringFlag{
			Name:  "resourceId",
			Usage: "ResourceID for transfer",
		},
		&cli.StringFlag{
			Name:  "value",
			Usage: "Value of ETH that should be sent along with deposit to cover possible fees. In ETH (decimals are allowed)",
			Value: "0",
		},
	},
}

func deposit(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Uint64("gasLimit")
	gasPrice := cctx.Uint64("gasPrice")

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
		return errors.New("invalid minter address")
	}
	recipientAddress := common.HexToAddress(recipient)
	id := cctx.Int64("id")
	dest := cctx.Uint64("dest")
	resourceId := cctx.String("resourceId")
	resourceIDBytes := utils.SliceTo32Bytes(common.Hex2Bytes(resourceId))

	value := cctx.String("value")

	realValue, err := utils.UserAmountToWei(value, big.NewInt(18))
	if err != nil {
		return err
	}

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	ethClient.ClientWithArgs(client.ClientWithValue(realValue))

	err = utils.MakeAndSendERC721Deposit(ethClient, bridgeAddress, recipientAddress, big.NewInt(id), resourceIDBytes, uint8(dest))
	if err != nil {
		return err
	}
	log.Info().Msgf("TokenID %s deposited to recipient address %s", big.NewInt(id).String(), recipientAddress.String())
	return nil
}
