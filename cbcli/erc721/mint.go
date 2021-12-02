package erc721

import (
	"github.com/rs/zerolog/log"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/cbcli/cliutils"
	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/celo-org/celo-blockchain/common"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var mintCMD = &cli.Command{
	Name:        "mint",
	Description: "Mint tokens on an ERC721 mintable contract.",
	Action:      mint,
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
			Name:  "metadata",
			Usage: "Metadata (tokenURI) for token",
		},
	},
}

func mint(cctx *cli.Context) error {
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

	id := cctx.Int64("id")

	metadata := cctx.String("metadata")

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	err = utils.ERC721Mint(ethClient, erc721Address, sender.CommonAddress(), big.NewInt(id), metadata)
	if err != nil {
		return err
	}
	log.Info().Msgf("ERC721 token with id %s minted", big.NewInt(id).String())
	return nil
}
