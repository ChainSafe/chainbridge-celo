package cbcli

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/status-im/keycard-go/hexutils"
	"github.com/urfave/cli/v2"
)

//sig
//ethers.utils.keccak256(ethers.utils.hexlify(ethers.utils.toUtf8Bytes(sig))).substr(0, 10)
func getFunctionBytes(in string) [4]byte {
	res := crypto.Keccak256(bytes.NewBufferString(in).Bytes())
	return utils.SliceTo4Bytes(res)
}

func RegisterGenericResource(cctx *cli.Context) error {
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
	resourceIdBytes := hexutils.HexToBytes(targetContract)
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
