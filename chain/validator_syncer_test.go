// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package chain

import (
	"math/big"
	"testing"

	"github.com/ChainSafe/chainbridge-celo/connection"
	"github.com/ChainSafe/chainbridge-utils/keystore"
	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common"
)

var TestEndpoint = "ws://localhost:8545"
var AliceKp = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
var GasLimit = big.NewInt(connection.DefaultGasLimit)
var GasPrice = big.NewInt(connection.DefaultGasPrice)

var testAddresses = []common.Address{
	// testAddresses are packed into our extra data
	// this references `extraData` from our genesis.json
	common.HexToAddress("0xecc833a7747eaa8327335e8e0c6b6d8aa3a38d00"),
	common.HexToAddress("0x82c07B76ee5D6a5Ec4bA710418ae299d3bdCE703"),
	common.HexToAddress("0x0000000000000000000000000000000000000000"),
}

func createTestConnection(t *testing.T) *connection.Connection {
	conn := connection.NewConnection(TestEndpoint, false, AliceKp, log15.Root(), GasLimit, GasPrice)
	return conn
}

func TestValidatorSyncer_ExtractValidators(t *testing.T) {
	conn := createTestConnection(t)
	vsyncer := ValidatorSyncer{conn: conn}
	err := vsyncer.start()
	if err != nil {
		t.Fatal(err)
	}

	validators, err := vsyncer.ExtractValidators(0)
	defer vsyncer.close()
	if err != nil {
		t.Fatal(err)
	}

	for i, v := range validators {
		if testAddresses[i] != v.Address {
			t.Fatalf("expected %s, got %s", testAddresses[i].Hex(), v.Address.Hex())
		}

	}

}

func TestValidatorSyncer_ExtractValidatorsDiff(t *testing.T) {
	conn := createTestConnection(t)
	vsyncer := ValidatorSyncer{conn: conn}
	err := vsyncer.start()
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = vsyncer.ExtractValidatorsDiff(0)
	if err != nil {
		t.Fatalf("failed to extract validators diff %s", err.Error())
	}
}
