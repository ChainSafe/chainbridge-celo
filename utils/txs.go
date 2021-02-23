package utils

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ChainSafe/chainbridge-utils/msg"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/pkg/errors"
	"github.com/status-im/keycard-go/hexutils"
	"math/big"
	"time"

	Bridge "github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
	erc20Handler "github.com/ChainSafe/chainbridge-celo/bindings/ERC20Handler"
	erc20 "github.com/ChainSafe/chainbridge-celo/bindings/ERC20PresetMinterPauser"
	erc721Handler "github.com/ChainSafe/chainbridge-celo/bindings/ERC721Handler"
	"github.com/ChainSafe/chainbridge-celo/bindings/ERC721MinterBurnerPauser"
	"github.com/ChainSafe/chainbridge-celo/bindings/GenericHandler"
	handlerHelper "github.com/ChainSafe/chainbridge-celo/bindings/HandlerHelpers"
	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-utils/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

var AliceKp = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
var BobKp = keystore.TestKeyRing.EthereumKeys[keystore.BobKey]
var EveKp = keystore.TestKeyRing.EthereumKeys[keystore.EveKey]

var (
	DefaultRelayerAddresses = []common.Address{
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.AliceKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.BobKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.CharlieKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.DaveKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.EveKey].Address()),
	}
)

const DefaultGasLimit = 6721975
const DefaultGasPrice = 20000000000

type DeployedContracts struct {
	BridgeAddress         common.Address
	ERC20HandlerAddress   common.Address
	ERC721HandlerAddress  common.Address
	GenericHandlerAddress common.Address
	ERC20TokenAddress     common.Address
}

func Erc20AddMinter(client *client.Client, erc20Contract, handler common.Address) error {
	err := client.LockAndUpdateOpts()
	if err != nil {
		return err
	}

	instance, err := erc20.NewERC20PresetMinterPauser(erc20Contract, client.Client)
	if err != nil {
		return err
	}

	role, err := instance.MINTERROLE(client.CallOpts())
	if err != nil {
		return err
	}

	tx, err := instance.GrantRole(client.Opts(), role, handler)
	if err != nil {
		return err
	}

	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}

	client.UnlockOpts()

	return nil
}

func SetBurnable(client *client.Client, bridge, handler, contract common.Address) error {
	instance, err := Bridge.NewBridge(bridge, client.Client)
	if err != nil {
		return err
	}

	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}

	tx, err := instance.AdminSetBurnable(client.Opts(), handler, contract)
	if err != nil {
		return err
	}

	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}

	client.UnlockOpts()

	return nil
}

func Erc20Approve(client *client.Client, erc20Contract, target common.Address, amount *big.Int) error {
	err := client.LockAndUpdateOpts()
	if err != nil {
		return err
	}

	instance, err := erc20.NewERC20PresetMinterPauser(erc20Contract, client.Client)
	if err != nil {
		return err
	}

	tx, err := instance.Approve(client.Opts(), target, amount)
	if err != nil {
		return err
	}

	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}

	client.UnlockOpts()

	return nil
}

func RegisterResource(client *client.Client, bridge, handler common.Address, rId [32]byte, addr common.Address) error {
	instance, err := Bridge.NewBridge(bridge, client.Client)
	if err != nil {
		return err
	}

	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}

	tx, err := instance.AdminSetResource(client.Opts(), handler, rId, addr)
	if err != nil {
		return err
	}

	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}

	client.UnlockOpts()

	return nil
}

func RegisterGenericResource(client *client.Client, bridge, handler common.Address, rId msg.ResourceId, addr common.Address, depositSig, executeSig [4]byte) error {
	instance, err := Bridge.NewBridge(bridge, client.Client)
	if err != nil {
		return err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}
	tx, err := instance.AdminSetGenericResource(client.Opts(), handler, rId, addr, depositSig, executeSig)
	if err != nil {
		return err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}
	client.UnlockOpts()
	return nil
}

func MintTokens(client *client.Client, erc20Addr common.Address, amount *big.Int) error {
	erc20Contract, err := erc20.NewERC20PresetMinterPauser(erc20Addr, client.Client)
	if err != nil {
		return err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}
	mintTx, err := erc20Contract.Mint(client.Opts(), AliceKp.CommonAddress(), amount)
	if err != nil {
		return err
	}
	err = WaitForTx(client, mintTx)
	if err != nil {
		return err
	}

	client.UnlockOpts()
	return nil
}

// DeployContracts deploys Bridge, Relayer, ERC20Handler, ERC721Handler and CentrifugeAssetHandler and returns the addresses
func DeployContracts(client *client.Client, chainID uint8, initialRelayerThreshold *big.Int, relayerAddresses []common.Address) (*DeployedContracts, error) {
	bridgeAddr, err := DeployBridge(client, chainID, relayerAddresses, initialRelayerThreshold)
	if err != nil {
		return nil, err
	}

	erc20HandlerAddr, err := DeployERC20Handler(client, bridgeAddr)
	if err != nil {
		return nil, err
	}

	erc721HandlerAddr, err := DeployERC721Handler(client, bridgeAddr)
	if err != nil {
		return nil, err
	}

	genericHandlerAddr, err := DeployGenericHandler(client, bridgeAddr)
	if err != nil {
		return nil, err
	}

	erc20Token, err := DeployERC20Token(client)
	if err != nil {
		return nil, err
	}

	dpc := &DeployedContracts{bridgeAddr, erc20HandlerAddr, erc721HandlerAddr, genericHandlerAddr, erc20Token}
	log.Debug().Msgf("Bridge %s \r\nerc20 handler %s \r\nerc721 handler %s \r\ngeneric handler %s \r\nerc20Contract %s", dpc.BridgeAddress.Hex(), dpc.ERC20HandlerAddress.Hex(), dpc.ERC721HandlerAddress.Hex(), dpc.GenericHandlerAddress.Hex(), dpc.ERC20TokenAddress.String())
	return dpc, nil
}

func DeployERC20Token(client *client.Client) (common.Address, error) {
	err := client.LockAndUpdateOpts()
	if err != nil {
		return common.Address{}, err
	}

	bridgeAddr, tx, _, err := erc20.DeployERC20PresetMinterPauser(client.Opts(), client.Client, "test", "TST")
	if err != nil {
		return common.Address{}, err
	}

	err = WaitForTx(client, tx)
	if err != nil {
		return common.Address{}, err
	}

	client.UnlockOpts()

	return bridgeAddr, nil
}

func DeployBridge(client *client.Client, chainID uint8, relayerAddrs []common.Address, initialRelayerThreshold *big.Int) (common.Address, error) {
	err := client.LockAndUpdateOpts()
	if err != nil {
		return common.Address{}, err
	}

	bridgeAddr, tx, _, err := Bridge.DeployBridge(client.Opts(), client.Client, chainID, relayerAddrs, initialRelayerThreshold, big.NewInt(0), big.NewInt(100))
	if err != nil {
		return common.Address{}, err
	}

	err = WaitForTx(client, tx)
	if err != nil {
		return common.Address{}, err
	}

	client.UnlockOpts()

	return bridgeAddr, nil
}

func DeployERC20Handler(client *client.Client, bridgeAddress common.Address) (common.Address, error) {
	err := client.LockAndUpdateOpts()
	if err != nil {
		return common.Address{}, err
	}

	erc20HandlerAddr, tx, _, err := erc20Handler.DeployERC20Handler(client.Opts(), client.Client, bridgeAddress, [][32]byte{}, []common.Address{}, []common.Address{})
	if err != nil {
		return common.Address{}, err
	}

	err = WaitForTx(client, tx)
	if err != nil {
		return common.Address{}, err
	}

	client.UnlockOpts()

	return erc20HandlerAddr, nil
}

func DeployERC721Handler(client *client.Client, bridgeAddress common.Address) (common.Address, error) {
	err := client.LockAndUpdateOpts()
	if err != nil {
		return common.Address{}, err
	}

	erc721HandlerAddr, tx, _, err := erc721Handler.DeployERC721Handler(client.Opts(), client.Client, bridgeAddress, [][32]byte{}, []common.Address{}, []common.Address{})
	if err != nil {
		return common.Address{}, err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return common.Address{}, err
	}

	client.UnlockOpts()

	return erc721HandlerAddr, nil
}

func DeployERC721Token(client *client.Client) (common.Address, error) {
	err := client.LockAndUpdateOpts()
	if err != nil {
		return common.Address{}, err
	}

	addr, tx, _, err := ERC721MinterBurnerPauser.DeployERC721MinterBurnerPauser(client.Opts(), client.Client, "", "", "")
	if err != nil {
		return common.Address{}, err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return common.Address{}, err
	}

	client.UnlockOpts()

	return addr, nil
}

func DeployGenericHandler(client *client.Client, bridgeAddress common.Address) (common.Address, error) {
	err := client.LockAndUpdateOpts()
	if err != nil {
		return common.Address{}, err
	}

	addr, tx, _, err := GenericHandler.DeployGenericHandler(client.Opts(), client.Client, bridgeAddress, [][32]byte{}, []common.Address{}, [][4]byte{}, [][4]byte{})
	if err != nil {
		return common.Address{}, err
	}

	err = WaitForTx(client, tx)
	if err != nil {
		return common.Address{}, err
	}

	client.UnlockOpts()

	return addr, nil
}

func CancelProposal(client *client.Client, bridgeAddress common.Address, chainID uint8, depositNonce uint64, dataHash [32]byte) error {
	bridgeInstance, err := Bridge.NewBridge(bridgeAddress, client.Client)
	if err != nil {
		return err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}
	tx, err := bridgeInstance.CancelProposal(client.Opts(), chainID, depositNonce, dataHash)
	if err != nil {
		return err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}

	client.UnlockOpts()
	return nil
}

func QueryProposal(client *client.Client, bridgeAddress common.Address, chainID uint8, depositNonce uint64, dataHash [32]byte) (*Bridge.BridgeProposal, error) {
	bridgeInstance, err := Bridge.NewBridge(bridgeAddress, client.Client)
	if err != nil {
		return nil, err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return nil, err
	}
	prop, err := bridgeInstance.GetProposal(client.CallOpts(), chainID, depositNonce, dataHash)
	if err != nil {
		return nil, err
	}

	client.UnlockOpts()
	return &prop, nil
}

func QueryResource(client *client.Client, handler common.Address, resourceID [32]byte) (common.Address, error) {
	handlerInstance, err := handlerHelper.NewHandlerHelpersCaller(handler, client.Client)
	if err != nil {
		return common.Address{}, err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return common.Address{}, err
	}
	addr, err := handlerInstance.ResourceIDToTokenContractAddress(client.CallOpts(), resourceID)
	if err != nil {
		return common.Address{}, err
	}
	client.UnlockOpts()
	return addr, nil
}

func AdminIsRelayer(client *client.Client, bridge common.Address, relayer common.Address) (bool, error) {
	bridgeInstance, err := Bridge.NewBridge(bridge, client.Client)
	if err != nil {
		return false, err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return false, err
	}
	prop, err := bridgeInstance.IsRelayer(client.CallOpts(), relayer)
	if err != nil {
		return false, err
	}

	client.UnlockOpts()
	return prop, nil
}

func AdminAddRelyaer(client *client.Client, bridge common.Address, relayer common.Address) error {
	bridgeInstance, err := Bridge.NewBridge(bridge, client.Client)
	if err != nil {
		return err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}
	tx, err := bridgeInstance.AdminAddRelayer(client.Opts(), relayer)
	if err != nil {
		return err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}
	client.UnlockOpts()
	return nil
}

func AdminRemoveRelayer(client *client.Client, bridge common.Address, relayer common.Address) error {
	bridgeInstance, err := Bridge.NewBridge(bridge, client.Client)
	if err != nil {
		return err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}
	tx, err := bridgeInstance.AdminRemoveRelayer(client.Opts(), relayer)
	if err != nil {
		return err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}
	client.UnlockOpts()
	return nil
}

func AdminSetTreshHold(client *client.Client, bridge common.Address, treshHold *big.Int) error {
	bridgeInstance, err := Bridge.NewBridge(bridge, client.Client)
	if err != nil {
		return err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}
	tx, err := bridgeInstance.AdminChangeRelayerThreshold(client.Opts(), treshHold)
	if err != nil {
		return err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}
	client.UnlockOpts()
	return nil
}

func AdminPause(client *client.Client, bridge common.Address) error {
	bridgeInstance, err := Bridge.NewBridge(bridge, client.Client)
	if err != nil {
		return err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}
	tx, err := bridgeInstance.AdminPauseTransfers(client.Opts())
	if err != nil {
		return err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}
	client.UnlockOpts()
	return nil
}

func AdminUnpause(client *client.Client, bridge common.Address) error {
	bridgeInstance, err := Bridge.NewBridge(bridge, client.Client)
	if err != nil {
		return err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}
	tx, err := bridgeInstance.AdminUnpauseTransfers(client.Opts())
	if err != nil {
		return err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}
	client.UnlockOpts()
	return nil
}

func AdminSetFee(client *client.Client, bridge common.Address, newFee *big.Int) error {
	bridgeInstance, err := Bridge.NewBridge(bridge, client.Client)
	if err != nil {
		return err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}
	tx, err := bridgeInstance.AdminChangeFee(client.Opts(), newFee)
	if err != nil {
		return err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}
	client.UnlockOpts()
	return nil
}

func AdminWithdraw(client *client.Client, bridge, handler, token, recipient common.Address, amount *big.Int) error {
	bridgeInstance, err := Bridge.NewBridge(bridge, client.Client)
	if err != nil {
		return err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}
	tx, err := bridgeInstance.AdminWithdraw(client.Opts(), handler, token, recipient, amount)
	if err != nil {
		return err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}
	client.UnlockOpts()
	return nil
}

var AdminRole = "0x0000000000000000000000000000000000000000000000000000000000000000"

func AdminAddAdmin(client *client.Client, bridge common.Address, newAdmin common.Address) error {
	bridgeInstance, err := Bridge.NewBridge(bridge, client.Client)
	if err != nil {
		return err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}
	tx, err := bridgeInstance.GrantRole(client.Opts(), SliceTo32Bytes(hexutils.HexToBytes(AdminRole)), newAdmin)
	if err != nil {
		return err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}
	client.UnlockOpts()
	return nil
}

func AdminRemoveAdmin(client *client.Client, bridge common.Address, addresToRevoke common.Address) error {
	bridgeInstance, err := Bridge.NewBridge(bridge, client.Client)
	if err != nil {
		return err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}
	tx, err := bridgeInstance.RevokeRole(client.Opts(), SliceTo32Bytes(hexutils.HexToBytes(AdminRole)), addresToRevoke)
	if err != nil {
		return err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}
	client.UnlockOpts()
	return nil
}

func ERC20Mint(client *client.Client, amount *big.Int, erc20Address, recipient common.Address) error {
	erc20Instance, err := erc20.NewERC20PresetMinterPauser(erc20Address, client.Client)
	if err != nil {
		return err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}
	tx, err := erc20Instance.Mint(client.Opts(), recipient, amount)
	if err != nil {
		return err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}
	client.UnlockOpts()
	return nil
}

func ERC20MinterRole(client *client.Client, erc20Address common.Address) ([32]byte, error) {
	erc20Instance, err := erc20.NewERC20PresetMinterPauser(erc20Address, client.Client)
	if err != nil {
		return [32]byte{}, err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return [32]byte{}, err
	}
	res, err := erc20Instance.MINTERROLE(client.CallOpts())
	if err != nil {
		return [32]byte{}, err
	}
	client.UnlockOpts()
	return res, nil
}

func ERC20AddMinter(client *client.Client, erc20Address, minter common.Address) error {
	erc20Instance, err := erc20.NewERC20PresetMinterPauser(erc20Address, client.Client)
	if err != nil {
		return err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}
	role, err := ERC20MinterRole(client, erc20Address)
	if err != nil {
		return nil
	}

	tx, err := erc20Instance.GrantRole(client.Opts(), role, minter)
	if err != nil {
		return err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}
	client.UnlockOpts()
	return nil
}

func ERC20Approve(client *client.Client, erc20Address, spender common.Address, amount *big.Int) error {
	erc20Instance, err := erc20.NewERC20PresetMinterPauser(erc20Address, client.Client)
	if err != nil {
		return err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}

	tx, err := erc20Instance.Approve(client.Opts(), spender, amount)
	if err != nil {
		return err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}
	client.UnlockOpts()
	return nil
}

func BridgeDeposit(client *client.Client, bridge common.Address, destChainID uint8, resourceID [32]byte, data []byte) error {
	bridgeInstance, err := Bridge.NewBridge(bridge, client.Client)
	if err != nil {
		return err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}
	tx, err := bridgeInstance.Deposit(client.Opts(), destChainID, resourceID, data)
	if err != nil {
		return err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}
	client.UnlockOpts()
	return nil
}

func ERC20BalanceOf(client *client.Client, erc20Address, dest common.Address) (*big.Int, error) {
	erc20Instance, err := erc20.NewERC20PresetMinterPauser(erc20Address, client.Client)
	if err != nil {
		return nil, err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return nil, err
	}

	balance, err := erc20Instance.BalanceOf(client.CallOpts(), dest)
	if err != nil {
		return nil, err
	}
	client.UnlockOpts()
	return balance, nil
}

func ERC20Allowance(client *client.Client, erc20Address, spender, owner common.Address) (*big.Int, error) {
	erc20Instance, err := erc20.NewERC20PresetMinterPauser(erc20Address, client.Client)
	if err != nil {
		return nil, err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return nil, err
	}

	balance, err := erc20Instance.Allowance(client.CallOpts(), owner, spender)
	if err != nil {
		return nil, err
	}
	client.UnlockOpts()
	return balance, nil
}

func ERC721Mint(client *client.Client, erc721Address, to common.Address, id *big.Int, metadata string) error {
	erc721Instance, err := ERC721MinterBurnerPauser.NewERC721MinterBurnerPauser(erc721Address, client.Client)
	if err != nil {
		return err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}
	tx, err := erc721Instance.Mint(client.Opts(), to, id, metadata)
	if err != nil {
		return err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}
	client.UnlockOpts()
	return nil
}

func ERC721MinterRole(client *client.Client, erc721Address common.Address) ([32]byte, error) {
	erc721Instance, err := ERC721MinterBurnerPauser.NewERC721MinterBurnerPauser(erc721Address, client.Client)
	if err != nil {
		return [32]byte{}, err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return [32]byte{}, err
	}
	res, err := erc721Instance.MINTERROLE(client.CallOpts())
	if err != nil {
		return [32]byte{}, err
	}
	client.UnlockOpts()
	return res, nil
}

func ERC721OwnerOf(client *client.Client, erc721Address common.Address, id *big.Int) (common.Address, error) {
	erc721Instance, err := ERC721MinterBurnerPauser.NewERC721MinterBurnerPauser(erc721Address, client.Client)
	if err != nil {
		return common.Address{}, err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return common.Address{}, err
	}
	res, err := erc721Instance.OwnerOf(client.CallOpts(), id)
	if err != nil {
		return common.Address{}, err
	}
	client.UnlockOpts()
	return res, nil
}

func ERC721AddMinter(client *client.Client, erc721Address, minter common.Address) error {
	erc721Instance, err := ERC721MinterBurnerPauser.NewERC721MinterBurnerPauser(erc721Address, client.Client)
	if err != nil {
		return err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}
	role, err := ERC721MinterRole(client, erc721Address)
	if err != nil {
		return err
	}
	tx, err := erc721Instance.GrantRole(client.Opts(), role, minter)
	if err != nil {
		return err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}
	client.UnlockOpts()
	return nil
}

func ERC721Approve(client *client.Client, erc721Address, recipient common.Address, id *big.Int) error {
	erc721Instance, err := ERC721MinterBurnerPauser.NewERC721MinterBurnerPauser(erc721Address, client.Client)
	if err != nil {
		return err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}
	tx, err := erc721Instance.Approve(client.Opts(), recipient, id)
	if err != nil {
		return err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}
	client.UnlockOpts()
	return nil
}

//nolint
func ERC20Transfer(client *client.Client, erc20 *erc20.ERC20PresetMinterPauser, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	err := client.LockAndUpdateOpts()
	if err != nil {
		return nil, err
	}

	tx, err := erc20.Transfer(client.Opts(), recipient, amount)
	if err != nil {
		return nil, err
	}
	client.UnlockOpts()
	return tx, nil
}

var ExpectedBlockTime = 2 * time.Second

// WaitForTx will query the chain at ExpectedBlockTime intervals, until a receipt is returned.
// Returns an error if the tx failed.
func WaitForTx(client *client.Client, tx *types.Transaction) error {
	retry := 10
	for retry > 0 {
		receipt, err := client.Client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			retry--
			time.Sleep(ExpectedBlockTime)
			continue
		}

		if receipt.Status != 1 {
			return fmt.Errorf("transaction failed on chain")
		}
		return nil
	}
	return nil
}

// WaitForTx will query the chain at ExpectedBlockTime intervals, until a receipt is returned.
// Returns an error if the tx failed.
func WaitAndReturnTxReceipt(client *client.Client, tx *types.Transaction) (*types.Receipt, error) {
	retry := 10
	for retry > 0 {
		receipt, err := client.Client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			retry--
			time.Sleep(ExpectedBlockTime)
			continue
		}
		if receipt.Status != 1 {
			return receipt, fmt.Errorf("transaction failed on chain")
		}
		return receipt, nil
	}
	return nil, errors.New("Tx do not appear")
}

func MakeErc20Deposit(client *client.Client, bridge *Bridge.Bridge, erc20ContractAddr, dest common.Address, amount *big.Int) (*types.Transaction, error) {
	data := ConstructErc20DepositData(dest.Bytes(), amount)
	err := client.LockAndUpdateOpts()
	if err != nil {
		return nil, err
	}

	src := ChainId(5)
	resourceID := SliceTo32Bytes(append(common.LeftPadBytes(erc20ContractAddr.Bytes(), 31), uint8(src)))
	tx, err := bridge.Deposit(client.Opts(), 1, resourceID, data)
	if err != nil {
		return nil, err
	}
	client.UnlockOpts()
	return tx, nil
}

func MakeAndSendERC20Deposit(client *client.Client, bridgeAddress common.Address, recipient common.Address, amount *big.Int, resourceID [32]byte, destChainID uint8) error {
	data := ConstructErc20DepositData(recipient.Bytes(), amount)
	err := client.LockAndUpdateOpts()
	if err != nil {
		return err
	}
	bridgeInstance, err := Bridge.NewBridge(bridgeAddress, client.Client)
	if err != nil {
		return err
	}
	tx, err := bridgeInstance.Deposit(client.Opts(), destChainID, resourceID, data)
	if err != nil {
		return err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}
	client.UnlockOpts()
	return nil
}

func ConstructErc20DepositData(destRecipient []byte, amount *big.Int) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(amount, 32)...)
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(destRecipient))), 32)...)
	data = append(data, destRecipient...)
	return data
}

func MakeAndSendERC721Deposit(client *client.Client, bridgeAddress common.Address, recipient common.Address, id *big.Int, resourceID [32]byte, destChainID uint8) error {
	data := ConstructErc721DepositData(id, recipient.Bytes())
	err := client.LockAndUpdateOpts()
	if err != nil {
		return err
	}
	bridgeInstance, err := Bridge.NewBridge(bridgeAddress, client.Client)
	if err != nil {
		return err
	}
	tx, err := bridgeInstance.Deposit(client.Opts(), destChainID, resourceID, data)
	if err != nil {
		return err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}
	client.UnlockOpts()
	return nil
}

// constructErc20Data constructs the data field to be passed into an erc721 deposit call
func ConstructErc721DepositData(tokenId *big.Int, destRecipient []byte) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(tokenId, 32)...)                               // Resource Id + Token Id
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(destRecipient))), 32)...) // Length of recipient
	data = append(data, destRecipient...)                                                  // Recipient

	return data
}

func ConstructGenericDepositData(metadata []byte) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(metadata))), 32)...)
	data = append(data, metadata...)

	return data
}

//nolint
func Simulate(client *client.Client, block *big.Int, txHash common.Hash, from common.Address) ([]byte, error) {
	tx, _, err := client.Client.TransactionByHash(context.TODO(), txHash)
	if err != nil {
		return nil, err
	}
	msg := ethereum.CallMsg{
		From:                from,
		To:                  tx.To(),
		Gas:                 tx.Gas(),
		FeeCurrency:         tx.FeeCurrency(),
		GatewayFeeRecipient: tx.GatewayFeeRecipient(),
		GatewayFee:          tx.GatewayFee(),
		GasPrice:            tx.GasPrice(),
		Value:               tx.Value(),
		Data:                tx.Data(),
	}
	res, err := client.Client.CallContract(context.TODO(), msg, block)
	if err != nil {
		return nil, err
	}
	bs, err := hex.DecodeString(hexutils.BytesToHex(res))
	if err != nil {
		panic(err)
	}
	log.Debug().Msg(string(bs))
	return nil, nil
}

func BuildQuery(contract common.Address, sig EventSig, startBlock *big.Int, endBlock *big.Int) ethereum.FilterQuery {
	query := ethereum.FilterQuery{
		FromBlock: startBlock,
		ToBlock:   endBlock,
		Addresses: []common.Address{contract},
		Topics: [][]common.Hash{
			{sig.GetTopic()},
		},
	}
	return query
}
