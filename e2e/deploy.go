// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package e2e

import (
	"github.com/ChainSafe/chainbridge-celo/chain/sender"
	"math/big"
	"time"

	Bridge "github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
	erc20Handler "github.com/ChainSafe/chainbridge-celo/bindings/ERC20Handler"
	erc20 "github.com/ChainSafe/chainbridge-celo/bindings/ERC20PresetMinterPauser"
	erc721Handler "github.com/ChainSafe/chainbridge-celo/bindings/ERC721Handler"
	"github.com/ChainSafe/chainbridge-celo/bindings/GenericHandler"
	"github.com/ChainSafe/chainbridge-celo/pkg"
	"github.com/ChainSafe/chainbridge-utils/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

const TestEndpoint = "ws://localhost:8545"
const TestEndpoint2 = "ws://localhost:8547"
const Chain1ID = 0
const Chain2ID = 1

var AliceKp = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
var BobKp = keystore.TestKeyRing.EthereumKeys[keystore.BobKey]
var EveKp = keystore.TestKeyRing.EthereumKeys[keystore.EveKey]

var (
	RelayerAddresses = []common.Address{
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.AliceKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.BobKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.CharlieKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.DaveKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.EveKey].Address()),
	}
)

type DeployedContracts struct {
	BridgeAddress         common.Address
	ERC20HandlerAddress   common.Address
	ERC721HandlerAddress  common.Address
	GenericHandlerAddress common.Address
	ERC20TokenAddress     common.Address
}

const DefaultGasLimit = 6721975
const DefaultGasPrice = 20000000000

var ExpectedBlockTime = time.Second

func Deploy(ctx *cli.Context) error {
	client, err := sender.NewSender(TestEndpoint, false, AliceKp, big.NewInt(DefaultGasLimit), big.NewInt(DefaultGasPrice))
	if err != nil {
		return err
	}

	dpc, err := DeployContracts(client, 1, big.NewInt(1))
	if err != nil {
		return err
	}
	src := pkg.ChainId(5)
	resourceID := pkg.SliceTo32Bytes(append(common.LeftPadBytes(dpc.ERC20TokenAddress.Bytes(), 31), uint8(src)))
	err = RegisterResource(client, dpc.BridgeAddress, dpc.ERC20HandlerAddress, resourceID, dpc.ERC20TokenAddress)
	if err != nil {
		return err
	}

	err = MintTokens(client, dpc.ERC20TokenAddress)
	if err != nil {
		return err
	}

	tenTokens := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))

	err = Erc20Approve(client, dpc.ERC20TokenAddress, dpc.ERC20HandlerAddress, tenTokens)
	if err != nil {
		return err
	}

	err = Erc20AddMinter(client, dpc.ERC20TokenAddress, dpc.ERC20HandlerAddress)
	if err != nil {
		return err
	}

	err = SetBurnable(client, dpc.BridgeAddress, dpc.ERC20HandlerAddress, dpc.ERC20TokenAddress)
	if err != nil {
		return err
	}

	return nil
}
func Erc20AddMinter(client *sender.Sender, erc20Contract, handler common.Address) error {
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

	client.UnlockNonce()

	return nil
}

func SetBurnable(client *sender.Sender, bridge, handler, contract common.Address) error {
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

	client.UnlockNonce()

	return nil
}

func Erc20Approve(client *sender.Sender, erc20Contract, target common.Address, amount *big.Int) error {
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

	client.UnlockNonce()

	return nil
}

func RegisterResource(client *sender.Sender, bridge, handler common.Address, rId [32]byte, addr common.Address) error {
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

	client.UnlockNonce()

	return nil
}

func MintTokens(client *sender.Sender, erc20Addr common.Address) error {
	erc20Contract, err := erc20.NewERC20PresetMinterPauser(erc20Addr, client.Client)
	if err != nil {
		return err
	}
	err = client.LockAndUpdateOpts()
	if err != nil {
		return err
	}
	tenTokens := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))
	mintTx, err := erc20Contract.Mint(client.Opts(), AliceKp.CommonAddress(), tenTokens)
	if err != nil {
		return err
	}
	err = WaitForTx(client, mintTx)
	if err != nil {
		return err
	}

	client.UnlockNonce()
	return nil
}

// DeployContracts deploys Bridge, Relayer, ERC20Handler, ERC721Handler and CentrifugeAssetHandler and returns the addresses
func DeployContracts(client *sender.Sender, chainID uint8, initialRelayerThreshold *big.Int) (*DeployedContracts, error) {
	bridgeAddr, err := deployBridge(client, chainID, RelayerAddresses, initialRelayerThreshold)
	if err != nil {
		return nil, err
	}

	erc20HandlerAddr, err := deployERC20Handler(client, bridgeAddr)
	if err != nil {
		return nil, err
	}

	erc721HandlerAddr, err := deployERC721Handler(client, bridgeAddr)
	if err != nil {
		return nil, err
	}

	genericHandlerAddr, err := deployGenericHandler(client, bridgeAddr)
	if err != nil {
		return nil, err
	}

	erc20Token, err := deployERC20Token(client)
	if err != nil {
		return nil, err
	}

	dpc := &DeployedContracts{bridgeAddr, erc20HandlerAddr, erc721HandlerAddr, genericHandlerAddr, erc20Token}
	log.Debug().Msgf("Bridge %s \r\nerc20 handler %s \r\nerc721 handler %s \r\ngeneric handler %s \r\nerc20Contract %s", dpc.BridgeAddress.Hex(), dpc.ERC20HandlerAddress.Hex(), dpc.ERC721HandlerAddress.Hex(), dpc.GenericHandlerAddress.Hex(), dpc.ERC20TokenAddress.String())

	return dpc, nil

}

func deployERC20Token(client *sender.Sender) (common.Address, error) {
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

	client.UnlockNonce()

	return bridgeAddr, nil
}

func deployBridge(client *sender.Sender, chainID uint8, relayerAddrs []common.Address, initialRelayerThreshold *big.Int) (common.Address, error) {
	err := client.LockAndUpdateNonce()
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

	client.UnlockNonce()

	return bridgeAddr, nil

}

func deployERC20Handler(client *sender.Sender, bridgeAddress common.Address) (common.Address, error) {
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

	client.UnlockNonce()

	return erc20HandlerAddr, nil
}

func deployERC721Handler(client *sender.Sender, bridgeAddress common.Address) (common.Address, error) {
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

	client.UnlockNonce()

	return erc721HandlerAddr, nil
}

func deployGenericHandler(client *sender.Sender, bridgeAddress common.Address) (common.Address, error) {
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

	client.UnlockNonce()

	return addr, nil
}
