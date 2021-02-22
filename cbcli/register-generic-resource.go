package cbcli

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/status-im/keycard-go/hexutils"
	"github.com/urfave/cli/v2"
)

func getFunctionBytes(in string) [4]byte {
	res := crypto.Keccak256(bytes.NewBufferString(in).Bytes())
	return utils.SliceTo4Bytes(res)
}

func registerGenericResource(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Int64("gasLimit")
	gasPrice := cctx.Int64("gasPrice")

	depositSig := cctx.String("deposit")
	depositSigBytes := hexutils.HexToBytes(depositSig)
	depositSigBytesArr := utils.SliceTo4Bytes(depositSigBytes)

	executeSig := cctx.String("execute")
	executeSigBytes := hexutils.HexToBytes(executeSig)
	executeSigBytesArr := utils.SliceTo4Bytes(executeSigBytes)

	if cctx.Bool("hash") {
		depositSigBytesArr = getFunctionBytes(depositSig)
		executeSigBytesArr = getFunctionBytes(executeSig)
	}

	sender, err := defineSender(cctx)
	if err != nil {
		return err
	}

	bridge := cctx.String("bridge")
	if !common.IsHexAddress(bridge) {
		return errors.New("invalid bridge address")
	}
	bridgeAddress := common.HexToAddress(bridge)

	handler := cctx.String("handler")
	if !common.IsHexAddress(handler) {
		return errors.New("invalid bridge address")
	}
	handlerAddress := common.HexToAddress(handler)
	targetContract := cctx.String("targetContract")
	if !common.IsHexAddress(targetContract) {
		return errors.New("invalid bridge address")
	}
	targetContractAddress := common.HexToAddress(targetContract)
	resourceId := cctx.String("resourceId")
	resourceIdBytes := hexutils.HexToBytes(resourceId)
	resourceIdBytesArr := utils.SliceTo32Bytes(resourceIdBytes)

	fmt.Printf("Registering contract %s with resource ID %s on handler %s", targetContract, resourceId, handler)
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(gasLimit), big.NewInt(gasPrice))
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
	Description: "Register a resource ID with a contract address for a generic handler..",
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
