package cbcli

import (
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func CancelProposal(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Uint64("gasLimit")
	gasPrice := cctx.Uint64("gasPrice")
	sender, err := defineSender(cctx)
	if err != nil {
		return err
	}
	bridge := cctx.String("bridge")
	if !common.IsHexAddress(bridge) {
		return errors.New("invalid bridge address")
	}
	bridgeAddress := common.HexToAddress(bridge)

	chainID := cctx.Uint64("chainId")
	depositNonce := cctx.Uint64("depositNonce")
	dataHash := cctx.String("dataHash")
	dataHashBytes := utils.SliceTo32Bytes(common.Hex2Bytes(dataHash))

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice))
	if err != nil {
		return err
	}
	err = utils.CancelProposal(ethClient, bridgeAddress, uint8(chainID), depositNonce, dataHashBytes)
	if err != nil {
		return err
	}
	log.Info().Msgf("Setting proposal with chain ID %v and deposit nonce %v status to 'Cancelled", chainID, depositNonce)
	return nil
}

var cancelProposalCMD = &cli.Command{
	Name:        "cancel-proposal",
	Description: "Cancels an expired proposal.",
	Action:      CancelProposal,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "bridge",
			Usage: "Bridge contract address",
		},
		&cli.Uint64Flag{
			Name:  "chainId",
			Usage: "Chain ID of proposal to cancel",
		},
		&cli.Uint64Flag{
			Name:  "depositNonce",
			Usage: "Deposit nonce of proposal to cancel",
		},
		&cli.StringFlag{
			Name:  "dataHash",
			Usage: "Hash of proposal metadata",
		},
	},
}
