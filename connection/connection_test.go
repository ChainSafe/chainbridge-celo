// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package connection

import (
	"math/big"
	"testing"

	ethutils "github.com/ChainSafe/chainbridge-celo/shared/ethereum"
	ethtest "github.com/ChainSafe/chainbridge-celo/shared/ethereum/testing"
	"github.com/ChainSafe/chainbridge-utils/keystore"
	"github.com/ChainSafe/log15"
	ethcmn "github.com/ethereum/go-ethereum/common"
)

var TestEndpoint = "ws://localhost:8546"
var AliceKp = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
var GasLimit = big.NewInt(DefaultGasLimit)
var GasPrice = big.NewInt(DefaultGasPrice)

func TestConnect(t *testing.T) {
	conn := NewConnection(TestEndpoint, false, AliceKp, log15.Root(), GasLimit)
	err := conn.Connect()
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()
}

// TestContractCode is used to make sure the contracts are deployed correctly.
// This is probably the least intrusive way to check if the contracts exists
func TestContractCode(t *testing.T) {
	client := ethtest.NewClient(t, TestEndpoint, AliceKp)
	contracts, err := ethutils.DeployContracts(client, 0, big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}

	conn := NewConnection(TestEndpoint, false, AliceKp, log15.Root(), GasLimit)
	err = conn.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// The following section checks if the byteCode exists on the chain at the specificed Addresses
	err = conn.EnsureHasBytecode(contracts.BridgeAddress)
	if err != nil {
		t.Fatal(err)
	}

	err = conn.EnsureHasBytecode(ethcmn.HexToAddress("0x0"))
	if err == nil {
		t.Fatal("should detect no bytecode")
	}

}
