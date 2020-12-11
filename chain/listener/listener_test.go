// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only
package listener

import (
	"math/big"
	"testing"

	mock_listener "github.com/ChainSafe/chainbridge-celo/chain/listener/mock"
	mock_chain "github.com/ChainSafe/chainbridge-celo/chain/mock"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/ChainSafe/chainbridge-celo/chain"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type ListenerTestSuite struct {
	suite.Suite
	syncerMock       *mock_listener.MockBlockSyncer
	routerMock       *mock_listener.MockIRouter
	clientMock       *mock_chain.MockLogFilterWithLatestBlock
	blockStorerMock  *mock_listener.MockBlockstorer
	gomockController *gomock.Controller
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(ListenerTestSuite))
}

func (s *ListenerTestSuite) SetupSuite()    {}
func (s *ListenerTestSuite) TearDownSuite() {}
func (s *ListenerTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.syncerMock = mock_listener.NewMockBlockSyncer(gomockController)
	s.routerMock = mock_listener.NewMockIRouter(gomockController)
	s.clientMock = mock_chain.NewMockLogFilterWithLatestBlock(gomockController)
	s.blockStorerMock = mock_listener.NewMockBlockstorer(gomockController)
	s.gomockController = gomockController
}
func (s *ListenerTestSuite) TearDownTest() {}

func (s *ListenerTestSuite) TestListenerStartStop() {
	stopChn := make(chan struct{})
	errChn := make(chan error)

	l := NewListener(&chain.CeloChainConfig{StartBlock: big.NewInt(1)}, s.clientMock, s.blockStorerMock, stopChn, errChn, s.syncerMock, s.routerMock)
	close(stopChn)
	s.NotNil(l.pollBlocks())
}

func (s *ListenerTestSuite) TestLatestBlockUpdate() {
	stopChn := make(chan struct{})
	errChn := make(chan error)
	cfg := &chain.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	l := NewListener(cfg, s.clientMock, s.blockStorerMock, stopChn, errChn, s.syncerMock, s.routerMock)

	s.clientMock.EXPECT().LatestBlock().Return(big.NewInt(555), nil)
	s.syncerMock.EXPECT().Sync(gomock.Any()).Return(nil)
	//No event logs found
	s.clientMock.EXPECT().FilterLogs(gomock.Any(), gomock.Any()).Return(make([]types.Log, 0), nil)

	s.blockStorerMock.EXPECT().StoreBlock(big.NewInt(1))

	//ON second call to latest block we stopping goroutine
	s.clientMock.EXPECT().LatestBlock().DoAndReturn(func() (*big.Int, error) { close(stopChn); return nil, errors.New("err") })

	s.NotNil(l.pollBlocks())
	s.Equal(cfg.StartBlock.String(), "2")
}
