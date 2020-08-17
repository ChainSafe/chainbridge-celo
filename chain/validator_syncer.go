// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package chain

import (
	"encoding/hex"
	"context"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/connection"
	"github.com/celo-org/celo-bls-go/bls"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/ethereum/go-ethereum/consensus/istanbul/validator"
	"github.com/ethereum/go-ethereum/core/types"
	blscrypto "github.com/ethereum/go-ethereum/crypto/bls"
	"github.com/pkg/errors"
)

type ValidatorSyncer struct {
	conn       *connection.Connection
	validators []istanbul.ValidatorData
}

var extra = "f4314cb9046bece6aa54bb9533155434d0c76909ec0d01b5adf993cdfee480b43be638b346ca58bc7d63d2d0e8b288de24bb320c02fa254a79fecc14511dc176f4e15c012e7d1b8ea9717c82c07b76ee5d6a5ec4ba710418ae299d3bdce703351f7c465fbaeb7ba814b43d7206546051d90f1b801da2d666868b7b1caff3f51c60f7b5e73ea57232b4107f0ce86e610bdde19c273c5f51776bbc48c0a685741a13d6f70a78e5bb931c2357601fe10bbac66b2a14e6778301863c678d1ece07b7c7bd46550d076625714f93600146393c50577d28924e2881b15a5f7bc585f8cb0ef501fd4503ac00cacf63aff9e92c22b0315ba8a0052b98c2a6b3281549869e84beed58692471c042776e39fcdc32ffd10663d3d3fb6793c9705a272e09311469f92d7765070ec196d4380085ad2fb742ebf5d57d252cf4a1d21c6d5a2ddf4190a3a122bc0167e5307bd90c83109b6609ddf5ce287d56ef43601b00";

// ExtractValidators pulls the extra data from the block header and extract
// validators and returns an array of validator data
func (v *ValidatorSyncer) ExtractValidators(num uint64) ([]istanbul.ValidatorData, error) {
	//header, err := v.conn.Client().HeaderByNumber(context.Background(), new(big.Int).SetUint64(num))
	// if err != nil {
	// 	return []istanbul.ValidatorData{}, errors.Wrap(err, "getting the block header by number failed")
	// }
	extraData, err := hex.DecodeString(extra)
	if err!= nil {
		return nil, err
	}
	return validator.ExtractValidators(extraData), nil
}

// AggregatePublicKeys merges all the validators public keys into one
// and returns it as an aggeragated public key
func (v *ValidatorSyncer) AggregatePublicKeys() (*bls.PublicKey, error) {
	var publicKeys []blscrypto.SerializedPublicKey
	for _, validator := range v.validators {
		publicKeys = append(publicKeys, validator.BLSPublicKey)
	}

	publicKeyObjs := []*bls.PublicKey{}
	for _, publicKey := range publicKeys {
		publicKeyObj, err := bls.DeserializePublicKeyCached(publicKey[:])
		if err != nil {
			return nil, err
		}
		defer publicKeyObj.Destroy()
		publicKeyObjs = append(publicKeyObjs, publicKeyObj)
	}
	apk, err := bls.AggregatePublicKeys(publicKeyObjs)
	if err != nil {
		return nil, err
	}
	defer apk.Destroy()

	return apk, nil
}

// ExtractValidatorsDiff extracts all values of the IstanbulExtra (aka diff) from the header
func (v *ValidatorSyncer) ExtractValidatorsDiff(num uint64) ([]istanbul.ValidatorData, []istanbul.ValidatorData, error) {
	header, err := v.conn.Client().HeaderByNumber(context.Background(), new(big.Int).SetUint64(num))
	if err != nil {
		return []istanbul.ValidatorData{}, []istanbul.ValidatorData{}, errors.Wrap(err, "getting the block header by number failed")
	}

	diff, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return []istanbul.ValidatorData{}, []istanbul.ValidatorData{}, errors.Wrap(err, "failed to extract istanbul extra from header")
	}

	var addedValidators []istanbul.ValidatorData
	for i, addr := range diff.AddedValidators {
		addedValidators = append(addedValidators, istanbul.ValidatorData{Address: addr, BLSPublicKey: diff.AddedValidatorsPublicKeys[i]})
	}

	bitmap := diff.RemovedValidators.Bytes()
	var removedValidators []istanbul.ValidatorData

	for _, i := range bitmap {
		removedValidators = append(removedValidators, v.validators[i])
	}

	return addedValidators, removedValidators, nil
}

func (v *ValidatorSyncer) start() error {
	err := v.conn.Connect()
	if err != nil {
		return err
	}
	return nil
}

func (v *ValidatorSyncer) close() {
	v.conn.Close()
}
