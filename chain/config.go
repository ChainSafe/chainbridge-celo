// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package chain

import (
	"math/big"
	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	bridgeContract         common.Address
	erc20HandlerContract   common.Address
	erc721HandlerContract  common.Address
	genericHandlerContract common.Address
	startBlock *big.Int
}
