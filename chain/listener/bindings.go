// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only
package listener

import (
	erc20 "github.com/ChainSafe/chainbridge-celo/bindings/ERC20Handler"
	erc721 "github.com/ChainSafe/chainbridge-celo/bindings/ERC721Handler"
	genericHandler "github.com/ChainSafe/chainbridge-celo/bindings/GenericHandler"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type IBridge interface {
	ResourceIDToHandlerAddress(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error)
}

type IERC20Handler interface {
	GetDepositRecord(opts *bind.CallOpts, depositNonce uint64, destId uint8) (erc20.ERC20HandlerDepositRecord, error)
}

type IERC721Handler interface {
	GetDepositRecord(opts *bind.CallOpts, depositNonce uint64, destId uint8) (erc721.ERC721HandlerDepositRecord, error)
}

type IGenericHandler interface {
	GetDepositRecord(opts *bind.CallOpts, depositNonce uint64, destId uint8) (genericHandler.GenericHandlerDepositRecord, error)
}
