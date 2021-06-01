package bridge

import (
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"

	"github.com/ChainSafe/chainbridge-celo/cbcli/cliutils"
	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
)

func registerGenericResource(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Int64("gasLimit")
	gasPrice := cctx.Int64("gasPrice")

	depositSig := cctx.String("deposit")
	depositSigBytes := common.Hex2Bytes(depositSig)
	depositSigBytesArr := utils.SliceTo4Bytes(depositSigBytes)

	executeSig := cctx.String("execute")
	executeSigBytes := common.Hex2Bytes(executeSig)
	executeSigBytesArr := utils.SliceTo4Bytes(executeSigBytes)

	if cctx.Bool("hash") {
		depositSigBytesArr = utils.GetSolidityFunctionSig(depositSig)
		executeSigBytesArr = utils.GetSolidityFunctionSig(executeSig)
	}

	sender, err := cliutils.DefineSender(cctx)
	if err != nil {
		return err
	}

	bridgeAddress, err := cliutils.DefineBridgeAddress(cctx)
	if err != nil {
		return err
	}

	handler := cctx.String("handler")
	if !common.IsHexAddress(handler) {
		return fmt.Errorf("invalid handler address %s", handler)
	}
	handlerAddress := common.HexToAddress(handler)
	targetContract := cctx.String("targetContract")
	if !common.IsHexAddress(targetContract) {
		return fmt.Errorf("invalid targetContract address %s", targetContract)
	}
	targetContractAddress := common.HexToAddress(targetContract)
	resourceId := cctx.String("resourceId")
	resourceIdBytes := common.Hex2Bytes(resourceId)
	resourceIdBytesArr := utils.SliceTo32Bytes(resourceIdBytes)

	log.Info().Msgf("Registering contract %s with resource ID %s on handler %s", targetContract, resourceId, handler)
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(gasLimit), big.NewInt(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	err = utils.RegisterGenericResource(ethClient, bridgeAddress, handlerAddress, resourceIdBytesArr, targetContractAddress, depositSigBytesArr, executeSigBytesArr)
	if err != nil {
		return err
	}
	fmt.Println("Resource registered")
	return nil
}

var registerGenericResourceCMD = &cli.Command{
	Name:        "register-generic-resource",
	Description: "Register a resource ID with a contract address for a generic handler.",
	Action:      registerGenericResource,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "bridge",
			Usage: "Bridge contract address",
		},
		&cli.StringFlag{
			Name:  "handler",
			Usage: "Handler address",
		},
		&cli.StringFlag{
			Name:  "targetContract",
			Usage: "Contract address to be registered",
		},
		&cli.StringFlag{
			Name:  "resourceId",
			Usage: "Resource ID to be registered",
		},
		&cli.StringFlag{
			Name:  "deposit",
			Usage: "Deposit function signature",
			Value: "0x00000000",
		},
		&cli.StringFlag{
			Name:  "execute",
			Usage: "Execute proposal function signature",
			Value: "0x00000000",
		},
		&cli.BoolFlag{
			Name:  "hash",
			Usage: "Treat signature inputs as function prototype strings, hash and take the first 4 bytes ",
		},
	},
}
