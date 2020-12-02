// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only
//nolint
//TODO remove nolint when start using this pakage
package validator

import (
	"context"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/celo-org/celo-bls-go/bls"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/ethereum/go-ethereum/core/types"
	blscrypto "github.com/ethereum/go-ethereum/crypto/bls"
	"github.com/pkg/errors"
)

type ValidatorSyncer struct {
	conn       *client.Client
	validators []istanbul.ValidatorData
}

// ExtractValidators pulls the extra data from the block header and extract
// validators and returns an array of validator data
func (v *ValidatorSyncer) ExtractValidators(num uint64) ([]istanbul.ValidatorData, error) {
	header, err := v.conn.Client().HeaderByNumber(context.Background(), new(big.Int).SetUint64(num))
	if err != nil {
		return []istanbul.ValidatorData{}, errors.Wrap(err, "getting the block header by number failed")
	}

	extra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return []istanbul.ValidatorData{}, errors.Wrap(err, "failed to extract istanbul extra from header")
	}
	var validators []istanbul.ValidatorData

	for i := range extra.AddedValidators {
		validator := &istanbul.ValidatorData{
			Address:      extra.AddedValidators[i],
			BLSPublicKey: extra.AddedValidatorsPublicKeys[i],
		}

		validators = append(v.validators, *validator)
	}

	return validators, nil

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
	return nil
}

func (v *ValidatorSyncer) close() {
	v.conn.Close()
}

func (v *ValidatorSyncer) Sync(latestBlock *big.Int) error {
	return nil
}
