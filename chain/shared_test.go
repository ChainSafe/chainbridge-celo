// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package chain

import (
	"math/big"
	"time"
    "log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ChainSafe/chainbridge-utils/keystore"
	"github.com/ChainSafe/chainbridge-utils/msg"
	"github.com/ChainSafe/chainbridge-celo/config"
)

var TestEndpoint string
var GasLimit *big.Int 
var GasPrice *big.Int
var GasLimitUint64 uint64
var ZeroAddress common.Address 
var TestChainID uint8
var TestRelayerThreshold *big.Int
var TestTimeout time.Duration
var DefaultGasLimit *big.Int
var DefaultGasPrice *big.Int
var TestChainId msg.ChainId

var AliceKp = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
var BobKp = keystore.TestKeyRing.EthereumKeys[keystore.BobKey]

func init() {
	configManager, err := config.NewConfigurationManager("../config")
	if err != nil {
		log.Fatal(err)
	}
	TestEndpoint = configManager.GetTestEndPoint()
	DefaultGasLimit = configManager.GetDefaultGasLimit()
	DefaultGasPrice = configManager.GetDefaultGasPrice()
	GasPrice = configManager.GetDefaultGasPrice()
	GasLimit = configManager.GetDefaultGasLimit()
	ZeroAddress = configManager.GetZeroAddress()
	TestChainID =  configManager.GetTestChainID()
	TestChainId =  msg.ChainId(configManager.GetTestChainID())
	TestTimeout = configManager.GetTestTimeout()
	TestRelayerThreshold = configManager.GetTestRelayerThreshold()
	GasLimitUint64 = configManager.GetDefaultGasLimit64()

}
