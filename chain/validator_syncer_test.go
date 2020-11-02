// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package chain

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ChainSafe/chainbridge-celo/connection"
	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common"
)

var GasLimit = big.NewInt(connection.DefaultGasLimit)
var GasPrice = big.NewInt(connection.DefaultGasPrice)


var expectedAddresses = []common.Address{
	// expectedAddresses are packed into our extra data
	// this references `extraData` from our genesis.json
	common.HexToAddress("0xf4314cb9046bece6aa54bb9533155434d0c76909"),
}

var expectedBlsPublicKeys = []string{
	// expectedBlsPublicKeys are packed into our extra data
	// this references `extraData` from our genesis.json
	"ec0d01b5adf993cdfee480b43be638b346ca58bc7d63d2d0e8b288de24bb320c02fa254a79fecc14511dc176f4e15c012e7d1b8ea9717c82c07b76ee5d6a5ec4ba710418ae299d3bdce703351f7c465fbaeb7ba814b43d7206546051d90f1b80",
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
	if err != nil {
		t.Fatal(err)
	}
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
