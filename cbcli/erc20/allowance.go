package erc20

import (
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/cbcli/cliutils"
	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/celo-org/celo-blockchain/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var allowanceCMD = &cli.Command{
	Name:        "allowance",
	Description: "Get the allowance of a spender for an address",
	Action:      allowance,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "spender",
			Usage: "Address of spender",
		},
		&cli.StringFlag{
			Name:  "owner",
			Usage: "Address to tokens owner",
		},
		&cli.StringFlag{
			Name:  "erc20Address",
			Usage: "erc20 contract address",
		},
	},
}

func allowance(cctx *cli.Context) error {
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

	spender := cctx.String("spender")
	if !common.IsHexAddress(spender) {
		return errors.New("invalid spender address")
	}
	spenderAddress := common.HexToAddress(spender)

	owner := cctx.String("owner")
	if !common.IsHexAddress(owner) {
		return errors.New("invalid owner address")
	}
	ownerAddress := common.HexToAddress(owner)

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	balance, err := utils.ERC20Allowance(ethClient, erc20Address, spenderAddress, ownerAddress)
	if err != nil {
		return err
	}
	log.Info().Msgf("allowance of %s to spend from address %s is %s", spenderAddress.String(), ownerAddress.String(), balance.String())
	return nil
}
