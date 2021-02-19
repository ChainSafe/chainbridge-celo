package cbcli

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/status-im/keycard-go/hexutils"
	"github.com/urfave/cli/v2"
)

func RegisterResource(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Int64("gasLimit")
	gasPrice := cctx.Int64("gasPrice")

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
	resourceIdBytes := hexutils.HexToBytes(targetContract)
	resourceIdBytesArr := utils.SliceTo32Bytes(resourceIdBytes)

	fmt.Printf("Registering contract %s with resource ID %s on handler %s", targetContract, resourceId, handler)
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(gasLimit), big.NewInt(gasPrice))
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