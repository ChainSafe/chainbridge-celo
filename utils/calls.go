package utils

import (
	"context"
	"encoding/hex"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/bindings/mptp/Bridge"
	erc20 "github.com/ChainSafe/chainbridge-celo/bindings/mptp/ERC20PresetMinterPauser"
	"github.com/ChainSafe/chainbridge-celo/bindings/mptp/ERC721MinterBurnerPauser"
	handlerHelper "github.com/ChainSafe/chainbridge-celo/bindings/mptp/HandlerHelpers"
	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

func QueryProposal(client *client.Client, bridgeAddress common.Address, chainID uint8, depositNonce uint64, dataHash [32]byte) (*Bridge.BridgeProposal, error) {
	bridgeInstance, err := Bridge.NewBridge(bridgeAddress, client.Client)
	if err != nil {
		return nil, err
	}
	prop, err := bridgeInstance.GetProposal(client.CallOpts(), chainID, depositNonce, dataHash)
	if err != nil {
		return nil, err
	}
	return &prop, nil
}

func QueryResource(client *client.Client, handler common.Address, resourceID [32]byte) (common.Address, error) {
	handlerInstance, err := handlerHelper.NewHandlerHelpersCaller(handler, client.Client)
	if err != nil {
		return common.Address{}, err
	}

	addr, err := handlerInstance.ResourceIDToTokenContractAddress(client.CallOpts(), resourceID)
	if err != nil {
		return common.Address{}, err
	}
	return addr, nil
}

func AdminIsRelayer(client *client.Client, bridge common.Address, relayer common.Address) (bool, error) {
	bridgeInstance, err := Bridge.NewBridge(bridge, client.Client)
	if err != nil {
		return false, err
	}

	prop, err := bridgeInstance.IsRelayer(client.CallOpts(), relayer)
	if err != nil {
		return false, err
	}

	return prop, nil
}

func ERC20MinterRole(client *client.Client, erc20Address common.Address) ([32]byte, error) {
	erc20Instance, err := erc20.NewERC20PresetMinterPauser(erc20Address, client.Client)
	if err != nil {
		return [32]byte{}, err
	}
	res, err := erc20Instance.MINTERROLE(client.CallOpts())
	if err != nil {
		return [32]byte{}, err
	}
	return res, nil
}

func ERC20IsMinter(client *client.Client, erc20Address, minter common.Address) (bool, error) {
	erc20Instance, err := erc20.NewERC20PresetMinterPauser(erc20Address, client.Client)
	if err != nil {
		return false, err
	}
	role, err := ERC20MinterRole(client, erc20Address)
	if err != nil {
		return false, nil
	}
	res, err := erc20Instance.HasRole(client.CallOpts(), role, minter)
	if err != nil {
		return false, err
	}
	return res, nil
}

func ERC20BalanceOf(client *client.Client, erc20Address, dest common.Address) (*big.Int, error) {
	erc20Instance, err := erc20.NewERC20PresetMinterPauser(erc20Address, client.Client)
	if err != nil {
		return nil, err
	}
	balance, err := erc20Instance.BalanceOf(client.CallOpts(), dest)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func ERC20Allowance(client *client.Client, erc20Address, spender, owner common.Address) (*big.Int, error) {
	erc20Instance, err := erc20.NewERC20PresetMinterPauser(erc20Address, client.Client)
	if err != nil {
		return nil, err
	}
	balance, err := erc20Instance.Allowance(client.CallOpts(), owner, spender)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func ERC721OwnerOf(client *client.Client, erc721Address common.Address, id *big.Int) (common.Address, error) {
	erc721Instance, err := ERC721MinterBurnerPauser.NewERC721MinterBurnerPauser(erc721Address, client.Client)
	if err != nil {
		return common.Address{}, err
	}
	res, err := erc721Instance.OwnerOf(client.CallOpts(), id)
	if err != nil {
		return common.Address{}, err
	}
	return res, nil
}

// Simulate function gets transaction info by hash and then executes a message call transaction, which is directly executed in the VM
// of the node, but never mined into the blockchain. Execution happens against provided block.
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
	bs, err := hex.DecodeString(common.Bytes2Hex(res))
	if err != nil {
		return nil, err
	}
	log.Debug().Msg(string(bs))
	return bs, nil
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

func ERC721MinterRole(client *client.Client, erc721Address common.Address) ([32]byte, error) {
	erc721Instance, err := ERC721MinterBurnerPauser.NewERC721MinterBurnerPauser(erc721Address, client.Client)
	if err != nil {
		return [32]byte{}, err
	}
	res, err := erc721Instance.MINTERROLE(client.CallOpts())
	if err != nil {
		return [32]byte{}, err
	}
	return res, nil
}

func AdminFeeAmount(client *client.Client, bridge common.Address) (*big.Int, error) {
	bridgeInstance, err := Bridge.NewBridge(bridge, client.Client)
	if err != nil {
		return nil, err
	}

	prop, err := bridgeInstance.Fee(client.CallOpts())
	if err != nil {
		return nil, err
	}

	return prop, nil
}
