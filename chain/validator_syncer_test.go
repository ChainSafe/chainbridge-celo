// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package chain

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ChainSafe/chainbridge-celo/connection"
	"github.com/ChainSafe/chainbridge-utils/keystore"
	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common"
)

var TestEndpoint = "ws://localhost:8545"
var GasLimit = big.NewInt(connection.DefaultGasLimit)
var GasPrice = big.NewInt(connection.DefaultGasPrice)

var AliceKp = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]

var expectedAddresses = []common.Address{
	// expectedAddresses are packed into our extra data
	// this references `extraData` from our genesis.json
	common.HexToAddress("0xf4314cb9046bece6aa54bb9533155434d0c76909"),
	common.HexToAddress("0x1da2d666868b7b1caff3f51c60f7b5e73ea57232"),
	common.HexToAddress("0xcacf63aff9e92c22b0315ba8a0052b98c2a6b328"),
}

var expectedBlsPublicKeys = []string{
	// expectedBlsPublicKeys are packed into our extra data
	// this references `extraData` from our genesis.json
	"ec0d01b5adf993cdfee480b43be638b346ca58bc7d63d2d0e8b288de24bb320c02fa254a79fecc14511dc176f4e15c012e7d1b8ea9717c82c07b76ee5d6a5ec4ba710418ae299d3bdce703351f7c465fbaeb7ba814b43d7206546051d90f1b80",
	"b4107f0ce86e610bdde19c273c5f51776bbc48c0a685741a13d6f70a78e5bb931c2357601fe10bbac66b2a14e6778301863c678d1ece07b7c7bd46550d076625714f93600146393c50577d28924e2881b15a5f7bc585f8cb0ef501fd4503ac00",
	"1549869e84beed58692471c042776e39fcdc32ffd10663d3d3fb6793c9705a272e09311469f92d7765070ec196d4380085ad2fb742ebf5d57d252cf4a1d21c6d5a2ddf4190a3a122bc0167e5307bd90c83109b6609ddf5ce287d56ef43601b00",
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
		if expectedAddresses[i] != v.Address {
			t.Fatalf("expected %s, got %s", expectedAddresses[i].Hex(), v.Address.Hex())
		}

		blsKeyHex := hex.EncodeToString(v.BLSPublicKey[:])

		if expectedBlsPublicKeys[i] != blsKeyHex {
			t.Fatalf("expected %s, got %s", expectedBlsPublicKeys[i], blsKeyHex)
		}

	}

}

func TestValidatorSyncer_AggregatePublicKeys(t *testing.T) {
	conn := createTestConnection(t)
	vsyncer := ValidatorSyncer{conn: conn}
	err := vsyncer.start()
	if err != nil {
		t.Fatal(err)
	}

	vsyncer.validators, err = vsyncer.ExtractValidators(0)
	defer vsyncer.close()

	_, err = vsyncer.AggregatePublicKeys()
	if err != nil {
		t.Fatalf("failed to aggergate the keys %s", err.Error())
	}
}

func TestValidatorSyncer_ExtractValidatorsDiff(t *testing.T) {
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
	vsyncer.validators = validators

	_, _, err = vsyncer.ExtractValidatorsDiff(0)
	if err != nil {
		t.Fatalf("failed to extract validators diff %s", err.Error())
	}
}
