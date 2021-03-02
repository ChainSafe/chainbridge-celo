package erc20

import (
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/cbcli/cliutils"
	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/status-im/keycard-go/hexutils"
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
	bridge := cctx.String("bridge")
	if !common.IsHexAddress(bridge) {
		return errors.New("invalid bridge address")
	}
	bridgeAddress := common.HexToAddress(bridge)

	recipient := cctx.String("recipient")
	if !common.IsHexAddress(recipient) {
		return errors.New("invalid minter address")
	}
	recipientAddress := common.HexToAddress(recipient)

	amount := cctx.String("amount")

	realAmount, err := utils.UserAmountToWei(amount, decimals)
	if err != nil {
		return err
	}

	dest := cctx.Uint64("dest")

	resourceId := cctx.String("resourceId")
	resourceIDBytes := utils.SliceTo32Bytes(hexutils.HexToBytes(resourceId))

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice))
	if err != nil {
		return err
	}
	err = utils.MakeAndSendERC20Deposit(ethClient, bridgeAddress, recipientAddress, realAmount, resourceIDBytes, uint8(dest))
	if err != nil {
		return err
	}
	log.Info().Msgf("%s account granted allowance on %v tokens of %s", recipientAddress.String(), amount, sender.CommonAddress().String())
	return nil
}
