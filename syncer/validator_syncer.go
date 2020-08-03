// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package syncer

import (
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

// ExtractValidators pulls the extra data from the block header and extract
// validators and returns an array of validator data
func (v *ValidatorSyncer) ExtractValidators(num uint64) ([]istanbul.ValidatorData, error) {
	header, err := v.conn.Client().HeaderByNumber(context.Background(), new(big.Int).SetUint64(num))
	if err != nil {
		return []istanbul.ValidatorData{}, errors.Wrap(err, "getting the block header by number failed")
	}

	return validator.ExtractValidators(header.Extra), nil
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

	return apk, nil
}

// ExtractValidatorsDiff extracts all values of the IstanbulExtra (aka diff) from the header
func (v *ValidatorSyncer) ExtractValidatorsDiff(num uint64) (*types.IstanbulExtra, error) {
	header, err := v.conn.Client().HeaderByNumber(context.Background(), new(big.Int).SetUint64(num))
	if err != nil {
		return nil, errors.Wrap(err, "getting the block header by number failed")
	}

	diff, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return nil, errors.Wrap(err, "failed to extract validators diff")
	}

	return diff, err
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
