package bridge

import (
	"github.com/ChainSafe/chainbridge-celo/cbcli/cliutils"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func queryProposal(cctx *cli.Context) error {
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

	chainID := cctx.Uint64("chainId")
	depositNonce := cctx.Uint64("depositNonce")
	dataHash := cctx.String("dataHash")
	dataHashBytes := utils.SliceTo32Bytes(common.Hex2Bytes(dataHash))

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}

	prop, err := utils.QueryProposal(ethClient, bridgeAddress, uint8(chainID), depositNonce, dataHashBytes)
	if err != nil {
		return err
	}
	log.Info().Msgf("proposal with chainID %v and depositNonce %v queried. %+v", chainID, depositNonce, prop)
	return nil
}

var queryProposalCMD = &cli.Command{
	Name:        "query-proposal",
	Description: "Queries an inbound proposal.",
	Action:      queryProposal,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "bridge",
			Usage: "Bridge contract address",
		},
		&cli.Uint64Flag{
			Name:  "chainId",
			Usage: "Source chain ID of proposal",
		},
		&cli.Uint64Flag{
			Name:  "depositNonce",
			Usage: "Deposit nonce of proposal",
		},
		&cli.StringFlag{
			Name:  "dataHash",
			Usage: "Hash of proposal metadata",
		},
	},
}
