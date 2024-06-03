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

var approveCMD = &cli.Command{
	Name:        "approve",
	Description: "Approve token in an ERC721 contract for transfer.",
	Action:      approve,
	Flags: []cli.Flag{
		&cli.Int64Flag{
			Name:  "id",
			Usage: "Token id",
		},
		&cli.StringFlag{
			Name:  "erc721Address",
			Usage: "ERC721 contract address",
		},
		&cli.StringFlag{
			Name:  "recipient",
			Usage: "Destination recipient address",
		},
	},
}

func approve(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Uint64("gasLimit")
	gasPrice := cctx.Uint64("gasPrice")
	sender, err := cliutils.DefineSender(cctx)
	if err != nil {
		return err
	}
	erc721 := cctx.String("erc721Address")
	if !common.IsHexAddress(erc721) {
		return errors.New("invalid erc721Address address")
	}
	erc721Address := common.HexToAddress(erc721)

	recipient := cctx.String("recipient")
	if !common.IsHexAddress(recipient) {
		return errors.New("invalid recipient address")
	}
	recipientAddress := common.HexToAddress(recipient)

	id := cctx.Int64("id")
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	err = utils.ERC721Approve(ethClient, erc721Address, recipientAddress, big.NewInt(id))
	if err != nil {
		return err
	}
	log.Info().Msgf("ERC721 token with id %s approved", big.NewInt(id).String())
	return nil
}
