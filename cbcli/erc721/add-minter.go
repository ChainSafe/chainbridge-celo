package erc721

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

var addMinterCMD = &cli.Command{
	Name:        "add-minter",
	Description: "Add a minter to an ERC721 mintable contact",
	Action:      addMinter,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "minter",
			Usage: "minter address",
		},
		&cli.StringFlag{
			Name:  "erc721Address",
			Usage: "ERC721 contract address",
		},
	},
}

func addMinter(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Uint64("gasLimit")
	gasPrice := cctx.Uint64("gasPrice")
	sender, err := cliutils.DefineSender(cctx)
	if err != nil {
		return err
	}
	erc721 := cctx.String("erc721Address")
	if !common.IsHexAddress(erc721) {
		return errors.New("invalid erc20Address address")
	}
	erc721Address := common.HexToAddress(erc721)

	minter := cctx.String("minter")
	if !common.IsHexAddress(minter) {
		return errors.New("invalid minter address")
	}
	minterAddress := common.HexToAddress(minter)
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	err = utils.ERC721AddMinter(ethClient, erc721Address, minterAddress)
	if err != nil {
		return err
	}
	log.Info().Msgf("Minter with address %s added", minterAddress.String())
	return nil
}
