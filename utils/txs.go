package utils

import (
	"context"
	"fmt"
	"github.com/ChainSafe/chainbridge-utils/msg"
	"github.com/pkg/errors"
	"math/big"
	"time"

	Bridge "github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
	erc20Handler "github.com/ChainSafe/chainbridge-celo/bindings/ERC20Handler"
	erc20 "github.com/ChainSafe/chainbridge-celo/bindings/ERC20PresetMinterPauser"
	erc721Handler "github.com/ChainSafe/chainbridge-celo/bindings/ERC721Handler"
	"github.com/ChainSafe/chainbridge-celo/bindings/ERC721MinterBurnerPauser"
	"github.com/ChainSafe/chainbridge-celo/bindings/GenericHandler"
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

func MintTokens(client *client.Client, erc20Addr common.Address) error {
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
