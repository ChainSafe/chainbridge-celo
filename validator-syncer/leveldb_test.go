package validator_syncer

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/ethereum/go-ethereum/crypto/bls"
	"github.com/stretchr/testify/suite"
	"math/big"
	"os"
	"testing"
)

type SyncerDBTestSuite struct {
	suite.Suite
	syncer *SyncerDB
}

func TestRunSyncerDBTestSuite(t *testing.T) {
	suite.Run(t, new(SyncerDBTestSuite))
}
func (s *SyncerDBTestSuite) SetupSuite()    {}
func (s *SyncerDBTestSuite) TearDownSuite() {}
func (s *SyncerDBTestSuite) SetupTest() {
	syncer, err := NewSyncerDB("test/db")
	if err != nil {
		s.Fail(err.Error())
	}
	s.syncer = syncer
}
func (s *SyncerDBTestSuite) TearDownTest() {
	s.syncer.Close()
	os.RemoveAll("./test")
}

func (s *SyncerDBTestSuite) TestPutAndGetLatestKnownBlock() {
	err := s.syncer.setLatestKnownBlock(big.NewInt(420))
	s.Nil(err)
	b, err := s.syncer.GetLatestKnownBlock()
	s.Nil(err)
	s.Equal(0, b.Cmp(big.NewInt(420)))
}

func (s *SyncerDBTestSuite) TestPutAndGetLatestKnownValidators() {
	startVals := make([]*istanbul.ValidatorData, 3)
	startVals[0] = &istanbul.ValidatorData{Address: common.Address{0x0f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals[1] = &istanbul.ValidatorData{Address: common.Address{0x1f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals[2] = &istanbul.ValidatorData{Address: common.Address{0x2f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	err := s.syncer.setLatestKnownValidators(startVals)
	s.Nil(err)

	v, err := s.syncer.GetLatestKnownValidators()
	s.Nil(err)
	s.Equal(3, len(v))
	s.Equal(common.Address{0x0f}, v[0].Address)
}

func (s *SyncerDBTestSuite) TestSetValidatorsForBlock() {
	startVals := make([]*istanbul.ValidatorData, 3)
	startVals[0] = &istanbul.ValidatorData{Address: common.Address{0x0f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals[1] = &istanbul.ValidatorData{Address: common.Address{0x1f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals[2] = &istanbul.ValidatorData{Address: common.Address{0x2f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	err := s.syncer.SetValidatorsForBlock(big.NewInt(420), startVals)
	s.Nil(err)
	v, err := s.syncer.GetLatestKnownValidators()
	s.Nil(err)
	s.Equal(3, len(v))
	s.Equal(common.Address{0x0f}, v[0].Address)
	b, err := s.syncer.GetLatestKnownBlock()
	s.Nil(err)
	s.Equal(0, b.Cmp(big.NewInt(420)))

	validators, err := s.syncer.GetValidatorsForBLock(big.NewInt(420))
	s.Nil(err)
	s.Equal(3, len(validators))
	s.Equal(common.Address{0x0f}, validators[0].Address)

}

func (s *SyncerDBTestSuite) TestGetLatestKnownBlockWithEmptyDB() {
	v, err := s.syncer.GetLatestKnownBlock()
	s.Nil(err)
	s.Equal(0, v.Cmp(big.NewInt(0)))
}

func (s *SyncerDBTestSuite) TestGetLatestKnownValidatorsFromEmptyDB() {
	v, err := s.syncer.GetLatestKnownValidators()
	s.Nil(err)
	s.Equal(0, len(v))
}
