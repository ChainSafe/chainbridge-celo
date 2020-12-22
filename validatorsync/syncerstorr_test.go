//Copyright 2020 ChainSafe Systems
//SPDX-License-Identifier: LGPL-3.0-only
package validatorsync

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/ethereum/go-ethereum/crypto/bls"
	"github.com/stretchr/testify/suite"
	"github.com/syndtr/goleveldb/leveldb"
	"math/big"
	"os"
	"testing"
)

type SyncerDBTestSuite struct {
	suite.Suite
	syncer *SyncerStorr
}

func TestRunSyncerDBTestSuite(t *testing.T) {
	suite.Run(t, new(SyncerDBTestSuite))
}
func (s *SyncerDBTestSuite) SetupSuite()    {}
func (s *SyncerDBTestSuite) TearDownSuite() {}
func (s *SyncerDBTestSuite) SetupTest() {
	db, err := leveldb.OpenFile("./test/db", nil)
	if err != nil {
		s.Fail(err.Error())
	}
	syncer := NewSyncerStorr(db)
	s.syncer = syncer
}
func (s *SyncerDBTestSuite) TearDownTest() {
	s.syncer.Close()
	os.RemoveAll("./test")
}

func (s *SyncerDBTestSuite) TestPutAndGetLatestKnownBlock() {
	chainID := uint8(1)
	err := s.syncer.setLatestKnownBlock(big.NewInt(420), chainID)
	s.Nil(err)
	b, err := s.syncer.GetLatestKnownBlock(chainID)
	s.Nil(err)
	s.Equal(0, b.Cmp(big.NewInt(420)))
}

func (s *SyncerDBTestSuite) TestPutAndGetLatestKnownValidators() {
	chainID := uint8(1)
	startVals := make([]*istanbul.ValidatorData, 3)
	startVals[0] = &istanbul.ValidatorData{Address: common.Address{0x0f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals[1] = &istanbul.ValidatorData{Address: common.Address{0x1f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals[2] = &istanbul.ValidatorData{Address: common.Address{0x2f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	err := s.syncer.setLatestKnownValidators(startVals, chainID)
	s.Nil(err)

	v, err := s.syncer.GetLatestKnownValidators(chainID)
	s.Nil(err)
	s.Equal(3, len(v))
	s.Equal(common.Address{0x0f}, v[0].Address)
}

func (s *SyncerDBTestSuite) TestSetValidatorsForBlock() {
	chainID := uint8(1)
	startVals := make([]*istanbul.ValidatorData, 3)
	startVals[0] = &istanbul.ValidatorData{Address: common.Address{0x0f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals[1] = &istanbul.ValidatorData{Address: common.Address{0x1f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals[2] = &istanbul.ValidatorData{Address: common.Address{0x2f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	err := s.syncer.SetValidatorsForBlock(big.NewInt(420), startVals, chainID)
	s.Nil(err)
	v, err := s.syncer.GetLatestKnownValidators(chainID)
	s.Nil(err)
	s.Equal(3, len(v))
	s.Equal(common.Address{0x0f}, v[0].Address)
	b, err := s.syncer.GetLatestKnownBlock(chainID)
	s.Nil(err)
	s.Equal(0, b.Cmp(big.NewInt(420)))

	validators, err := s.syncer.GetValidatorsForBLock(big.NewInt(420), chainID)
	s.Nil(err)
	s.Equal(3, len(validators))
	s.Equal(common.Address{0x0f}, validators[0].Address)

}

func (s *SyncerDBTestSuite) TestGetLatestKnownBlockWithEmptyDB() {
	chainID := uint8(1)
	v, err := s.syncer.GetLatestKnownBlock(chainID)
	s.Nil(err)
	s.Equal(0, v.Cmp(big.NewInt(0)))
}

func (s *SyncerDBTestSuite) TestGetLatestKnownValidatorsFromEmptyDB() {
	chainID := uint8(1)
	v, err := s.syncer.GetLatestKnownValidators(chainID)
	s.Nil(err)
	s.Equal(0, len(v))
}

func (s *SyncerDBTestSuite) TestTestSetValidatorsForBlockForDifferentChains() {
	chainID1 := uint8(1)
	startVals1 := make([]*istanbul.ValidatorData, 3)
	startVals1[0] = &istanbul.ValidatorData{Address: common.Address{0x0f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals1[1] = &istanbul.ValidatorData{Address: common.Address{0x1f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals1[2] = &istanbul.ValidatorData{Address: common.Address{0x2f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	err := s.syncer.SetValidatorsForBlock(big.NewInt(420), startVals1, chainID1)
	s.Nil(err)

	chainID2 := uint8(2)
	startVals2 := make([]*istanbul.ValidatorData, 2)
	startVals2[0] = &istanbul.ValidatorData{Address: common.Address{0x3f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals2[1] = &istanbul.ValidatorData{Address: common.Address{0x4f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	err = s.syncer.SetValidatorsForBlock(big.NewInt(420), startVals2, chainID2)
	s.Nil(err)

	v, err := s.syncer.GetLatestKnownValidators(chainID1)
	s.Nil(err)
	s.Equal(3, len(v))
	s.Equal(common.Address{0x0f}, v[0].Address)
	b, err := s.syncer.GetLatestKnownBlock(chainID1)
	s.Nil(err)
	s.Equal(0, b.Cmp(big.NewInt(420)))

	validators, err := s.syncer.GetValidatorsForBLock(big.NewInt(420), chainID1)
	s.Nil(err)
	s.Equal(3, len(validators))
	s.Equal(common.Address{0x0f}, validators[0].Address)

	v, err = s.syncer.GetLatestKnownValidators(chainID2)
	s.Nil(err)
	s.Equal(2, len(v))
	s.Equal(common.Address{0x3f}, v[0].Address)
	b, err = s.syncer.GetLatestKnownBlock(chainID2)
	s.Nil(err)
	s.Equal(0, b.Cmp(big.NewInt(420)))

	validators, err = s.syncer.GetValidatorsForBLock(big.NewInt(420), chainID2)
	s.Nil(err)
	s.Equal(2, len(validators))
	s.Equal(common.Address{0x3f}, validators[0].Address)
}
