// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package chain

import (
	"math/big"
<<<<<<< HEAD
=======

>>>>>>> main
	"github.com/ChainSafe/chainbridge-utils/msg"
	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
<<<<<<< HEAD
	id                     msg.ChainId // ChainID
	name                   string      // Human-readable chain name
	endpoint               string      // url for rpc endpoint
	from                   string      // address of key to use
	keystorePath           string      // Location of keyfiles
	blockstorePath         string
	freshStart             bool // Disables loading from blockstore at start
=======
	id                     msg.ChainId
>>>>>>> main
	bridgeContract         common.Address
	erc20HandlerContract   common.Address
	erc721HandlerContract  common.Address
	genericHandlerContract common.Address
<<<<<<< HEAD
	gasLimit               *big.Int
	maxGasPrice            *big.Int
	http                   bool // Config for type of connection
=======
>>>>>>> main
	startBlock             *big.Int
}
