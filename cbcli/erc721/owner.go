package erc721

import (
	"github.com/rs/zerolog/log"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/cbcli/cliutils"
	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var ownerCMD = &cli.Command{
	Name:        "owner",
	Description: "Query ownerOf for a token",
	Action:      owner,
	Flags: []cli.Flag{
		&cli.Int64Flag{
			Name:  "id",
			Usage: "Token id",
		},
		&cli.StringFlag{
			Name:  "erc721Address",
			Usage: "ERC721 contract address",
		},
	},
}

func owner(cctx *cli.Context) error {
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

	id := cctx.Int64("id")

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice))
	if err != nil {
		return err
	}
	addr, err := utils.ERC721OwnerOf(ethClient, erc721Address, big.NewInt(id))
	if err != nil {
		return err
	}
	log.Info().Msgf("owner of token %s is ", big.NewInt(id).String(), addr.String())
	return nil
}
