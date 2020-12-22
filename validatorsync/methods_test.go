//Copyright 2020 ChainSafe Systems
//SPDX-License-Identifier: LGPL-3.0-only
package validatorsync

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/bls"
	"math/big"
	"testing"

	"github.com/stretchr/testify/suite"
)

type WriterTestSuite struct {
	suite.Suite
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(WriterTestSuite))
}

func (s *WriterTestSuite) SetupSuite()    {}
func (s *WriterTestSuite) TearDownSuite() {}
func (s *WriterTestSuite) SetupTest()     {}
func (s *WriterTestSuite) TearDownTest()  {}

func (s *WriterTestSuite) TestBitSetToTrue() {
	i := big.NewInt(4)
	s.True(bitSetToTrue(2, i))
	s.False(bitSetToTrue(0, i))

	i = big.NewInt(3)
	s.True(bitSetToTrue(0, i))
	s.False(bitSetToTrue(2, i))

	i = big.NewInt(3)
	s.True(bitSetToTrue(1, i))
}

func (s *WriterTestSuite) TestApplyValidatorsDiff() {
	startVals := make([]*istanbul.ValidatorData, 3)
	startVals[0] = &istanbul.ValidatorData{Address: common.Address{0x0f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals[1] = &istanbul.ValidatorData{Address: common.Address{0x1f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals[2] = &istanbul.ValidatorData{Address: common.Address{0x2f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	addedAddresses := []common.Address{{0x3f}}

	extra := &types.IstanbulExtra{
		AddedValidators:           addedAddresses,
		RemovedValidators:         big.NewInt(4),
		AddedValidatorsPublicKeys: []blscrypto.SerializedPublicKey{{0x3f}},
	}
	resVals, err := applyValidatorsDiff(extra, startVals)
	s.Nil(err)
	s.Equal(len(resVals), 3)
	s.Equal(resVals[2].BLSPublicKey, blscrypto.SerializedPublicKey{0x3f})
}

func (s *WriterTestSuite) TestApplyValidatorsDiffEmptyStartVals() {
	startVals := make([]*istanbul.ValidatorData, 0)
	addedAddresses := []common.Address{{0x3f}}

	extra := &types.IstanbulExtra{
		AddedValidators:           addedAddresses,
		RemovedValidators:         big.NewInt(0),
		AddedValidatorsPublicKeys: []blscrypto.SerializedPublicKey{{0x3f}},
	}
	resVals, err := applyValidatorsDiff(extra, startVals)
	s.Nil(err)
	s.Equal(len(resVals), 1)
	s.Equal(resVals[0].BLSPublicKey, blscrypto.SerializedPublicKey{0x3f})
}

func (s *WriterTestSuite) TestApplyValidatorsDiffWithRemovedOnEmptyVals() {
	startVals := make([]*istanbul.ValidatorData, 0)
	addedAddresses := []common.Address{{0x3f}}

	extra := &types.IstanbulExtra{
		AddedValidators:           addedAddresses,
		RemovedValidators:         big.NewInt(1),
		AddedValidatorsPublicKeys: []blscrypto.SerializedPublicKey{{0x3f}},
	}
	resVals, err := applyValidatorsDiff(extra, startVals)
	s.Nil(resVals)
	s.NotNil(err)
	s.Equal(err, ErrorWrongInitialValidators)
}
