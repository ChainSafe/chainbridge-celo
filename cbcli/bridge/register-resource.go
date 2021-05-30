package bridge

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/cbcli/cliutils"
	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func registerResource(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Int64("gasLimit")
	gasPrice := cctx.Int64("gasPrice")

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
		return errors.New(fmt.Sprintf("invalid handler address %s", handler))
	}
	handlerAddress := common.HexToAddress(handler)
	targetContract := cctx.String("targetContract")
	if !common.IsHexAddress(targetContract) {
		return errors.New(fmt.Sprintf("invalid targetContract address %s", targetContract))
	}
	targetContractAddress := common.HexToAddress(targetContract)
	resourceId := cctx.String("resourceId")
	resourceIdBytes := common.Hex2Bytes(resourceId)
	resourceIdBytesArr := utils.SliceTo32Bytes(resourceIdBytes)

	fmt.Printf("Registering contract %s with resource ID %s on handler %s", targetContract, resourceId, handler)
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(gasLimit), big.NewInt(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	err = utils.RegisterResource(ethClient, bridgeAddress, handlerAddress, resourceIdBytesArr, targetContractAddress)
	if err != nil {
		return err
	}
	fmt.Println("Resource registered")

	return nil
}

var registerResourceCMD = &cli.Command{
	Name:        "register-resource",
	Description: "Register a resource ID with a contract address for a handler.",
	Action:      registerResource,
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
	},
}
