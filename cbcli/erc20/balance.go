package erc20

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

var balanceCMD = &cli.Command{
	Name:        "balance",
	Description: "Query balance for an account in an ERC20 contract.",
	Action:      balanceOf,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "address",
			Usage: "Address to receive balanceOf",
		},
		&cli.StringFlag{
			Name:  "erc20Address",
			Usage: "erc20 contract address",
		},
	},
}

func balanceOf(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Uint64("gasLimit")
	gasPrice := cctx.Uint64("gasPrice")
	sender, err := cliutils.DefineSender(cctx)
	if err != nil {
		return err
	}
	erc20 := cctx.String("erc20Address")
	if !common.IsHexAddress(erc20) {
		return errors.New("invalid erc20Address address")
	}
	erc20Address := common.HexToAddress(erc20)

	address := cctx.String("address")
	if !common.IsHexAddress(address) {
		return errors.New("invalid target address")
	}
	targetAddress := common.HexToAddress(address)

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice))
	if err != nil {
		return err
	}
	balance, err := utils.ERC20BalanceOf(ethClient, erc20Address, targetAddress)
	if err != nil {
		return err
	}
	log.Info().Msgf("balance of %s is %s", targetAddress.String(), balance.String())
	return nil
}
