package erc20

//
//import (
//	"math/big"
//
//	"github.com/ChainSafe/chainbridge-celo/utils"
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/pkg/errors"
//	"github.com/urfave/cli/v2"
//)
//
//var dataHashCMD = &cli.Command{
//	Name:        "data-hash",
//	Description: "Constructs proposal data and returns the hash required for on-chain queries.",
//	Action:      balanceOf,
//	Flags: []cli.Flag{
//		&cli.StringFlag{
//			Name:  "amount",
//			Usage: "Amount to deposit",
//		},
//		&cli.StringFlag{
//			Name:  "recipient",
//			Usage: "Destination recipient address ",
//		},
//		&cli.StringFlag{
//			Name:  "erc20Address",
//			Usage: "erc20 contract address",
//		},
//	},
//}
//
//func dataHash(cctx *cli.Context) error {
//	decimals := big.NewInt(0).SetUint64(cctx.Uint64("decimals"))
//
//	erc20 := cctx.String("erc20Address")
//	if !common.IsHexAddress(erc20) {
//		return errors.New("invalid erc20Address address")
//	}
//	erc20Address := common.HexToAddress(erc20)
//
//	recipient := cctx.String("recipient")
//	if !common.IsHexAddress(recipient) {
//		return errors.New("invalid recipient address")
//	}
//	recipientAddress := common.HexToAddress(recipient)
//	amount := cctx.String("amount")
//	realAmount, err := utils.UserAmountToReal(amount, decimals)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
