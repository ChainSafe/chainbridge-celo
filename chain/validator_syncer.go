// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package chain

import (
	"context"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/connection"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/ethereum/go-ethereum/consensus/istanbul/validator"
	"github.com/ethereum/go-ethereum/core/types"
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
