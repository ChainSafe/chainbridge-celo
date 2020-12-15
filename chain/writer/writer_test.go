// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only
package writer

import (
	"math/big"
	"testing"

	"github.com/ChainSafe/chainbridge-celo/chain"
	mock_writer "github.com/ChainSafe/chainbridge-celo/chain/writer/mock"
	message "github.com/ChainSafe/chainbridge-celo/msg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type WriterTestSuite struct {
	suite.Suite
	client           *mock_writer.MockContractCaller
	gomockController *gomock.Controller
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(WriterTestSuite))
}

func (s *WriterTestSuite) SetupSuite() {
}

func (s *WriterTestSuite) TearDownSuite() {}
func (s *WriterTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.client = mock_writer.NewMockContractCaller(gomockController)
	s.gomockController = gomockController
}
func (s *WriterTestSuite) TearDownTest() {}

func (s *WriterTestSuite) TestResolveMessageWrongType() {
	resourceId := [32]byte{1}
	recipient := make([]byte, 32)
	amount := big.NewInt(10)
	stopChn := make(chan struct{})
	errChn := make(chan error)
	m := message.NewFungibleTransfer(1, 0, 0, amount, resourceId, recipient, nil)
	m.Type = "123"
	cfg := &chain.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	s.False(w.ResolveMessage(m))
}

//TestWriter_start_stop
//TestCreateAndExecuteErc20DepositProposal
//TestCreateAndExecuteErc721Proposal
//TestCreateAndExecuteGenericProposal
//TestDuplicateMessage
