//Copyright 2020 ChainSafe Systems
//SPDX-License-Identifier: LGPL-3.0-only
package validatorsync

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	blscrypto "github.com/ethereum/go-ethereum/crypto/bls"

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
	s.Equal(computeLastBlockOfEpochForProvidedBlock(big.NewInt(0), 2335), big.NewInt(0))

	s.Equal(computeLastBlockOfEpochForProvidedBlock(big.NewInt(11), 12), big.NewInt(12))

	s.Equal(computeLastBlockOfEpochForProvidedBlock(big.NewInt(251), 12), big.NewInt(252))

}

func (s *WriterTestSuite) TestAggregatePublicKeys() {
	startVals := make([]*istanbul.ValidatorData, 3)
	testKey1, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testKey2, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f292")
	testKey3, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f293")
	blsPK1, _ := blscrypto.ECDSAToBLS(testKey1)
	blsPK2, _ := blscrypto.ECDSAToBLS(testKey2)
	blsPK3, _ := blscrypto.ECDSAToBLS(testKey3)
	pubKey1, _ := blscrypto.PrivateToPublic(blsPK1)
	pubKey2, _ := blscrypto.PrivateToPublic(blsPK2)
	pubKey3, _ := blscrypto.PrivateToPublic(blsPK3)

	startVals[0] = &istanbul.ValidatorData{Address: common.Address{0x0f}, BLSPublicKey: pubKey1}
	startVals[1] = &istanbul.ValidatorData{Address: common.Address{0x1f}, BLSPublicKey: pubKey2}
	startVals[2] = &istanbul.ValidatorData{Address: common.Address{0x2f}, BLSPublicKey: pubKey3}
	apk, err := aggregatePublicKeys(startVals)
	s.Nil(err)
	s.NotNil(apk)
	// checking that function is clear
	apk2, err := aggregatePublicKeys(startVals)
	s.Nil(err)
	s.Equal(apk, apk2)
}

// TestFilterValidatorsWithBitmap tests that filterValidatorsWithBitmap
// returns a slice of validators who signed the current block after applying
// the current block's bitmap.
func (s *WriterTestSuite) TestFilterValidatorsWithBitmap() {
	// init new slice to hold initial validators slice
	startVals := make([]*istanbul.ValidatorData, 5)
	startVals[0] = &istanbul.ValidatorData{Address: common.Address{0x0f}}
	startVals[1] = &istanbul.ValidatorData{Address: common.Address{0x1f}}
	startVals[2] = &istanbul.ValidatorData{Address: common.Address{0x2f}}
	startVals[3] = &istanbul.ValidatorData{Address: common.Address{0x3f}}
	startVals[4] = &istanbul.ValidatorData{Address: common.Address{0x4f}}

	// init sample Istanbul Extra header data
	extra := &types.IstanbulExtra{
		AggregatedSeal: types.IstanbulAggregatedSeal{
			// init bitmap at 0
			Bitmap: big.NewInt(0),
		},
	}

	// loop over validators to set bitmap
	for valIndex := range startVals {
		// validators 4 (index 3) and 5 (index 4) did not sign
		if valIndex == 3 || valIndex == 4 {
			// skip
			continue
		}
		// set bitmap
		extra.AggregatedSeal.Bitmap.SetBit(
			extra.AggregatedSeal.Bitmap, valIndex, 1,
		)
	}

	// init new slice to hold expected validators slice after bitmask applied
	expected := make([]*istanbul.ValidatorData, 3)
	expected[0] = &istanbul.ValidatorData{Address: common.Address{0x0f}}
	expected[1] = &istanbul.ValidatorData{Address: common.Address{0x1f}}
	expected[2] = &istanbul.ValidatorData{Address: common.Address{0x2f}}
	// validator at index 3 => 0x3f did not sign
	// validator at index 4 => 0x4f did not sign

	newValidators := filterValidatorsWithBitmap(startVals, extra.AggregatedSeal.Bitmap)
	s.NotNil(newValidators)
	s.Equal(expected, newValidators)
}
