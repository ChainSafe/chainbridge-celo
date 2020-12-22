package validatorsync

import (
	"github.com/pkg/errors"
	"math/big"

	"github.com/celo-org/celo-bls-go/bls"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/bls"
)

var (
	ErrorWrongInitialValidators = errors.New("wrong initial validators")
)

func applyValidatorsDiff(extra *types.IstanbulExtra, validators []*istanbul.ValidatorData) ([]*istanbul.ValidatorData, error) {
	var addedValidators []*istanbul.ValidatorData
	for i, addr := range extra.AddedValidators {
		addedValidators = append(addedValidators, &istanbul.ValidatorData{Address: addr, BLSPublicKey: extra.AddedValidatorsPublicKeys[i]})
	}

	for i := 0; i < extra.RemovedValidators.BitLen(); i++ {
		if len(validators) <= i {
			return nil, ErrorWrongInitialValidators
		}
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
	return newValidators, nil
}

func bitSetToTrue(index int, bits *big.Int) bool {
	andY := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(int64(index)), nil)
	andRes := big.NewInt(0).And(bits, andY)
	return andRes.Cmp(big.NewInt(0)) > 0
}

func aggregatePublicKeys(validators []*istanbul.ValidatorData) (*bls.PublicKey, error) {
	publicKeys := make([]blscrypto.SerializedPublicKey, len(validators))
	for i := range validators {
		publicKeys[i] = validators[i].BLSPublicKey
	}

	publicKeyObjs := make([]*bls.PublicKey, len(publicKeys))
	for i := range publicKeys {
		publicKeyObj, err := bls.DeserializePublicKeyCached(publicKeys[i][:])
		if err != nil {
			return nil, err
		}
		publicKeyObjs[i] = publicKeyObj
		publicKeyObj.Destroy()
	}
	apk, err := bls.AggregatePublicKeys(publicKeyObjs)
	if err != nil {
		return nil, err
	}
	defer apk.Destroy()

	return apk, nil
}
