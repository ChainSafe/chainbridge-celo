package validator_syncer

import (
	"github.com/celo-org/celo-bls-go/bls"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/bls"
	"github.com/rs/zerolog/log"
	"math/big"
)

func ExtractValidatorsDiff(headerExtra *types.IstanbulExtra, validators []*istanbul.ValidatorData) ([]*istanbul.ValidatorData, []*istanbul.ValidatorData, error) {
	var addedValidators []*istanbul.ValidatorData
	for i, addr := range headerExtra.AddedValidators {
		addedValidators = append(addedValidators, &istanbul.ValidatorData{Address: addr, BLSPublicKey: headerExtra.AddedValidatorsPublicKeys[i]})
	}

	bitmap := headerExtra.RemovedValidators.Bytes()
	var removedValidators []*istanbul.ValidatorData

	for _, i := range bitmap {
		removedValidators = append(removedValidators, validators[i])
	}

	return addedValidators, removedValidators, nil
}

func ApplyValidatorsDiff(extra *types.IstanbulExtra, validators []*istanbul.ValidatorData) []*istanbul.ValidatorData {
	var addedValidators []*istanbul.ValidatorData
	for i, addr := range extra.AddedValidators {
		addedValidators = append(addedValidators, &istanbul.ValidatorData{Address: addr, BLSPublicKey: extra.AddedValidatorsPublicKeys[i]})
	}

	log.Debug().Msgf("LEN BYTEARR %s", extra.RemovedValidators.BitLen())
	for i := 0; i < extra.RemovedValidators.BitLen(); i++ {
		if bitSetToTrue(i, extra.RemovedValidators) {
			validators[i] = nil
		}
	}
	validators = append(validators, addedValidators...)
	newValidators := make([]*istanbul.ValidatorData, 0)
	for _, v := range validators {
		if v != nil {
			newValidators = append(newValidators, v)
		}
	}
	return newValidators
}

func remove(slice []*istanbul.ValidatorData, s int) []*istanbul.ValidatorData {
	return append(slice[:s], slice[s+1:]...)
}
func bitSetToTrue(index int, bits *big.Int) bool {
	andY := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(int64(index)), nil)
	andRes := big.NewInt(0).And(bits, andY)
	if andRes.Cmp(big.NewInt(0)) > 0 {
		return true
	}
	return false
}

func AggregatePublicKeys(validators []*istanbul.ValidatorData) (*bls.PublicKey, error) {
	var publicKeys []blscrypto.SerializedPublicKey
	for _, validator := range validators {
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
