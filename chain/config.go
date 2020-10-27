// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package chain

import (
	"math/big"
)

const ZeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")

type Config struct {
	startBlock *big.Int
}
