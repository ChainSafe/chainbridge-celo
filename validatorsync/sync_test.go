//Copyright 2020 ChainSafe Systems
//SPDX-License-Identifier: LGPL-3.0-only
package validatorsync

import (
	"github.com/ChainSafe/chainbridge-celo/validatorsync/mock"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/syndtr/goleveldb/leveldb"
)

type SyncTestSuite struct {
	suite.Suite
	syncer *ValidatorsStore
	client *mock_validatorsync.MockHeaderByNumberGetter
}

func TestRunSyncTestSuite(t *testing.T) {
	suite.Run(t, new(SyncTestSuite))
}
func (s *SyncTestSuite) SetupSuite()    {}
func (s *SyncTestSuite) TearDownSuite() {}
func (s *SyncTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	db, err := leveldb.OpenFile("./test/db", nil)
	if err != nil {
		s.Fail(err.Error())
	}
	syncer := NewValidatorsStore(db)
	s.syncer = syncer
	s.client = mock_validatorsync.NewMockHeaderByNumberGetter(gomockController)

}
func (s *SyncTestSuite) TearDownTest() {
	s.syncer.Close()
	os.RemoveAll("./test")
}

func (s *SyncTestSuite) TestStoreBlockValidators() {
	//stopChn := make(chan struct{})
	//errChn := make(chan error)
	//chainID := uint8(1)
	//StoreBlockValidators(stopChn, errChn, s.client, s.syncer, chainID)
}
