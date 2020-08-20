// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package chain

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/connection"
	"github.com/ChainSafe/log15"
	"github.com/celo-org/celo-bls-go/bls"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/ethereum/go-ethereum/core/types"
	blscrypto "github.com/ethereum/go-ethereum/crypto/bls"
	"github.com/pkg/errors"
)

const DefaultHeaderNumber = 0

type ValidatorSyncer struct {
	conn       *connection.Connection
	validators []istanbul.ValidatorData
	apk        *bls.PublicKey
	log        log15.Logger
}

func NewValidatorSyncer(conn *connection.Connection, log log15.Logger) *ValidatorSyncer {
	return &ValidatorSyncer{
		conn: conn,
		log:  log,
	}
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

	bitmap := diff.RemovedValidators.Bytes()
	var removedValidators []istanbul.ValidatorData

	v.log.Info("====== removed validators =======")
	for _, i := range bitmap {
		v.log.Info(fmt.Sprintf("address: %s, bls public key: %s", v.validators[i].Address.Hex(), hex.EncodeToString(v.validators[i].BLSPublicKey[:])))
		removedValidators = append(removedValidators, v.validators[i])
		v.validators = append(v.validators[:i], v.validators[i+1:]...)
	}
	v.log.Info("=================================")

	v.log.Info("====== added validators =======")
	var addedValidators []istanbul.ValidatorData
	for i, addr := range diff.AddedValidators {
		v.log.Info(fmt.Sprintf("address: %s, bls public key: %s", addr.Hex(), hex.EncodeToString(diff.AddedValidatorsPublicKeys[i][:])))
		addedValidators = append(addedValidators, istanbul.ValidatorData{Address: addr, BLSPublicKey: diff.AddedValidatorsPublicKeys[i]})
		v.validators = append(v.validators, istanbul.ValidatorData{Address: addr, BLSPublicKey: diff.AddedValidatorsPublicKeys[i]})
	}
	v.log.Info("===============================")

	return addedValidators, removedValidators, nil
}

func (v *ValidatorSyncer) Sync() error {
	err := v.start()
	if err != nil {
		return err
	}

	v.log.Info("Syncing the validators...")
	v.validators, err = v.ExtractValidators(0)
	defer v.close()
	if err != nil {
		return errors.Wrap(err, "failed to extract validators")
	}

	v.log.Info("Extracting validators diff...")
	removedValidators, addedValidators, err := v.ExtractValidatorsDiff(DefaultHeaderNumber)
	if err != nil {
		return errors.Wrap(err, "failed to extract validators diff")
	}

	// if there's a change aggregate a new public key
	if len(removedValidators) > 1 || len(addedValidators) > 1 {
		v.log.Info("Aggregating the public keys...")
		v.apk, err = v.AggregatePublicKeys()
		if err != nil {
			return errors.Wrap(err, "failed to aggregate public keys")
		}
	}

	return nil
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
