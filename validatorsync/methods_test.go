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

func (s *WriterTestSuite) TestApplyValidatorsDiffOnlyAdd() {
	startVals := make([]*istanbul.ValidatorData, 2)
	startVals[0] = &istanbul.ValidatorData{Address: common.Address{0x0f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals[1] = &istanbul.ValidatorData{Address: common.Address{0x1f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	addedAddresses := []common.Address{{0x3f}}
	extra := &types.IstanbulExtra{
		AddedValidators:           addedAddresses,
		RemovedValidators:         big.NewInt(0),
		AddedValidatorsPublicKeys: []blscrypto.SerializedPublicKey{{0x3f}},
	}
	resVals, err := applyValidatorsDiff(extra, startVals)
	s.Nil(err)
	s.Equal(3, len(resVals))
	s.Equal(resVals[2].BLSPublicKey, blscrypto.SerializedPublicKey{0x3f})
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

func (s *WriterTestSuite) TestDefineBlocksEpochLastBlockNumber() {
	s.Equal(defineBlocksEpochLastBlockNumber(big.NewInt(0), 2335), big.NewInt(0))

	s.Equal(defineBlocksEpochLastBlockNumber(big.NewInt(11), 12), big.NewInt(12))

	s.Equal(defineBlocksEpochLastBlockNumber(big.NewInt(251), 12), big.NewInt(252))

}
