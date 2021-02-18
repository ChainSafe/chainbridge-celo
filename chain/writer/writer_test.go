// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only
package writer

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
	"github.com/ChainSafe/chainbridge-celo/chain/config"
	mock_writer "github.com/ChainSafe/chainbridge-celo/chain/writer/mock"
	"github.com/ChainSafe/chainbridge-celo/utils"
	eth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
)

type WriterTestSuite struct {
	suite.Suite
	client           *mock_writer.MockContractCaller
	gomockController *gomock.Controller
	bridgeMock       *mock_writer.MockBridger
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(WriterTestSuite))
}

func (s *WriterTestSuite) SetupSuite()    {}
func (s *WriterTestSuite) TearDownSuite() {}
func (s *WriterTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.client = mock_writer.NewMockContractCaller(gomockController)
	s.bridgeMock = mock_writer.NewMockBridger(gomockController)
	s.gomockController = gomockController
}
func (s *WriterTestSuite) TearDownTest() {}

func (s *WriterTestSuite) TestResolveMessageWrongType() {
	resourceId := [32]byte{1}
	recipient := make([]byte, 32)
	amount := big.NewInt(10)
	stopChn := make(chan struct{})
	errChn := make(chan error)
	m := utils.NewFungibleTransfer(1, 0, utils.Nonce(555), resourceId, nil, nil, amount, recipient)
	m.Type = "123"
	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	s.False(w.ResolveMessage(m))
}

func (s *WriterTestSuite) TestHasVotedError() {
	stopChn := make(chan struct{})
	errChn := make(chan error)

	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	s.client.EXPECT().CallOpts().Return(nil)
	s.client.EXPECT().Opts().Return(&bind.TransactOpts{
		From: common.HexToAddress("0x"),
	})

	hash := crypto.Keccak256Hash([]byte("data"))

	s.bridgeMock.EXPECT().HasVotedOnProposal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, errors.New("error occured"))
	s.False(w.hasVoted(utils.ChainId(3), utils.Nonce(1), hash))
}

func (s *WriterTestSuite) TestShouldVoteProposalIsAlreadyComplete() {
	resourceId := [32]byte{1}
	recipient := make([]byte, 32)
	amount := big.NewInt(10)
	stopChn := make(chan struct{})
	errChn := make(chan error)
	m := utils.NewFungibleTransfer(1, 0, utils.Nonce(555), resourceId, nil, nil, amount, recipient)

	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	// Setting returned proposal to PassedStatus
	prop := Bridge.BridgeProposal{Status: ProposalStatusPassed}
	s.client.EXPECT().CallOpts().Return(nil)
	s.bridgeMock.EXPECT().GetProposal(gomock.Any(), gomock.Any(), uint64(m.DepositNonce), gomock.Any()).Return(prop, nil)
	s.False(w.shouldVote(m, common.Hash{}))
}

func (s *WriterTestSuite) TestShouldVoteProposalIsAlreadyVoted() {
	stopChn := make(chan struct{})
	errChn := make(chan error)
	m := utils.NewFungibleTransfer(1, 0, utils.Nonce(555), [32]byte{1}, nil, nil, big.NewInt(10), make([]byte, 32))

	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	// Setting returned proposal to PassedStatus
	var notPassedStatus uint8 = 0
	prop := Bridge.BridgeProposal{Status: notPassedStatus} // some other status

	s.client.EXPECT().CallOpts().Return(nil)
	s.bridgeMock.EXPECT().GetProposal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(prop, nil)

	s.client.EXPECT().CallOpts().Return(nil)
	s.client.EXPECT().Opts().Return(&bind.TransactOpts{From: common.Address{}})
	s.bridgeMock.EXPECT().HasVotedOnProposal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
	s.False(w.shouldVote(m, common.Hash{}))
}

func (s *WriterTestSuite) TestShouldVoteProposal() {
	stopChn := make(chan struct{})
	errChn := make(chan error)
	m := utils.NewFungibleTransfer(1, 0, utils.Nonce(555), [32]byte{1}, nil, nil, big.NewInt(10), make([]byte, 32))

	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	// Setting returned proposal to PassedStatus
	var notPassedStatus uint8 = 0
	prop := Bridge.BridgeProposal{Status: notPassedStatus} // some other status

	s.client.EXPECT().CallOpts().Return(nil)
	s.bridgeMock.EXPECT().GetProposal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(prop, nil)

	s.client.EXPECT().CallOpts().Return(nil)
	s.client.EXPECT().Opts().Return(&bind.TransactOpts{From: common.Address{}})
	s.bridgeMock.EXPECT().HasVotedOnProposal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
	s.True(w.shouldVote(m, common.Hash{}))
}

func (s *WriterTestSuite) TestVoteProposalAlreadyComplete() {
	stopChn := make(chan struct{})
	errChn := make(chan error)
	m := utils.NewFungibleTransfer(utils.ChainId(1), 0, utils.Nonce(555), [32]byte{1}, nil, nil, big.NewInt(10), make([]byte, 32))

	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	s.client.EXPECT().CallOpts().Return(nil)
	proposal := Bridge.BridgeProposal{
		Status: ProposalStatusPassed,
	}

	s.bridgeMock.EXPECT().GetProposal(gomock.Any(), uint8(m.Source), uint64(m.DepositNonce), gomock.Any()).Return(proposal, nil)

	//Vote proposal should not be called, since proposal already passed
	//s.bridgeMock.EXPECT().VoteProposal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(s.Fail("Vote proposal should not be voted"))

	go func() {
		select {
		case <-errChn:
			s.Fail("err channel has value")
		case <-time.After(time.Second * 10):
			// Closing this goroutine after 10 seconds
			return
		}
	}()

	w.voteProposal(m, common.Hash{})
}

func (s *WriterTestSuite) TestVoteProposalIsNotComplete() {
	stopChn := make(chan struct{})
	errChn := make(chan error)
	m := utils.NewFungibleTransfer(utils.ChainId(1), 0, utils.Nonce(555), [32]byte{1}, nil, nil, big.NewInt(10), make([]byte, 32))

	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	proposal := Bridge.BridgeProposal{
		Status: 0,
	}
	tx := types.NewTransaction(5577006791947779410, common.Address{0x0f}, new(big.Int), 0, new(big.Int), &common.Address{0x0f}, &common.Address{0x0f}, big.NewInt(10), nil)
	s.client.EXPECT().CallOpts().Return(nil)
	s.bridgeMock.EXPECT().GetProposal(gomock.Any(), uint8(m.Source), uint64(m.DepositNonce), gomock.Any()).Return(proposal, nil)
	s.client.EXPECT().LockAndUpdateOpts()
	s.client.EXPECT().Opts()
	s.bridgeMock.EXPECT().VoteProposal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tx, nil)
	s.client.EXPECT().UnlockOpts()

	go func() {
		select {
		case <-errChn:
			s.Fail("err channel has value")
		case <-time.After(time.Second * 10):
			// Closing this goroutine after 10 seconds
			return
		}
	}()

	w.voteProposal(m, common.Hash{})
}

func (s *WriterTestSuite) TestVoteProposalUnexpectedErrorOnVote() {
	stopChn := make(chan struct{})
	errChn := make(chan error)
	m := utils.NewFungibleTransfer(utils.ChainId(1), 0, utils.Nonce(555), [32]byte{1}, nil, nil, big.NewInt(10), make([]byte, 32))

	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	proposal := Bridge.BridgeProposal{
		Status: 0,
	}
	for i := 0; i < TxRetryLimit; i++ {
		s.client.EXPECT().CallOpts().Return(nil)
		s.bridgeMock.EXPECT().GetProposal(gomock.Any(), uint8(m.Source), uint64(m.DepositNonce), gomock.Any()).Return(proposal, nil)
		s.client.EXPECT().LockAndUpdateOpts()
		s.client.EXPECT().Opts()
		s.bridgeMock.EXPECT().VoteProposal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("unexpectedERROR"))
		s.client.EXPECT().UnlockOpts()
	}

	go func() {
		err := <-errChn
		s.NotNil(err)
	}()

	w.voteProposal(m, common.Hash{})
}

func (s *WriterTestSuite) TestVoteProposalLockAndUpdateOptsError() {
	stopChn := make(chan struct{})
	errChn := make(chan error)
	m := utils.NewFungibleTransfer(utils.ChainId(1), 0, utils.Nonce(555), [32]byte{1}, nil, nil, big.NewInt(10), make([]byte, 32))

	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	proposal := Bridge.BridgeProposal{
		Status: 0,
	}
	for i := 0; i < TxRetryLimit; i++ {
		s.client.EXPECT().CallOpts().Return(nil)
		s.bridgeMock.EXPECT().GetProposal(gomock.Any(), uint8(m.Source), uint64(m.DepositNonce), gomock.Any()).Return(proposal, nil)
		s.client.EXPECT().LockAndUpdateOpts().Return(errors.New("error")).Times(1)
		s.client.EXPECT().Opts().Times(0)
		s.bridgeMock.EXPECT().VoteProposal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("unexpectedERROR")).Times(0)
		s.client.EXPECT().UnlockOpts().Times(0)
	}

	go func() {
		err := <-errChn
		s.NotNil(err)
	}()

	w.voteProposal(m, common.Hash{})
}

func (s *WriterTestSuite) TestExecuteProposalLockAndUpdateOptsError() {

	stopChn := make(chan struct{})
	errChn := make(chan error)
	pkg := utils.NewFungibleTransfer(utils.ChainId(1), 0, utils.Nonce(555), [32]byte{1}, nil, nil, big.NewInt(10), make([]byte, 32))

	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	for i := 0; i < TxRetryLimit; i++ {
		s.client.EXPECT().LockAndUpdateOpts().Return(errors.New("error"))
	}

	go func() {
		err := <-errChn
		s.True(err.Error() == ErrFatalTx.Error())
	}()

	w.executeProposal(pkg, []byte{}, common.Hash{})

}

func (s *WriterTestSuite) TestExecuteProposalNonceTooLowError() {

	stopChn := make(chan struct{})
	errChn := make(chan error)

	sig := &utils.SignatureVerification{
		AggregatePublicKey: []byte{},
		BlockHash:          common.Hash{},
		Signature:          []byte{},
	}

	mp := &utils.MerkleProof{
		TxRootHash: common.Hash{},
		Key:        []byte{},
		Nodes:      []byte{},
	}

	message := utils.NewFungibleTransfer(utils.ChainId(1), 0, utils.Nonce(555), [32]byte{1}, mp, sig, big.NewInt(10), make([]byte, 32))

	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	for i := 0; i < TxRetryLimit; i++ {
		s.client.EXPECT().LockAndUpdateOpts().Return(nil)

		s.client.EXPECT().Opts()

		s.client.EXPECT().CallOpts()

		s.bridgeMock.EXPECT().ExecuteProposal(
			gomock.Any(),
			uint8(message.Source),
			uint64(message.DepositNonce),
			[]byte{},
			gomock.Any(),
			[]byte{},
			[]byte{},
			[]byte{},
			gomock.Any(),
			gomock.Any(),
			[]byte{},
			[]byte{},
		).Return(nil, ErrNonceTooLow)

		s.client.EXPECT().UnlockOpts()

		s.bridgeMock.EXPECT().GetProposal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(Bridge.BridgeProposal{}, errors.New("error"))
	}

	go func() {
		err := <-errChn
		s.True(err.Error() == ErrFatalTx.Error())
	}()

	w.executeProposal(message, []byte{}, common.Hash{})

}

func (s *WriterTestSuite) TestExecuteProposalCompleted() {

	stopChn := make(chan struct{})
	errChn := make(chan error)

	sig := &utils.SignatureVerification{
		AggregatePublicKey: []byte{},
		BlockHash:          common.Hash{},
		Signature:          []byte{},
	}

	mp := &utils.MerkleProof{
		TxRootHash: common.Hash{},
		Key:        []byte{},
		Nodes:      []byte{},
	}

	message := utils.NewFungibleTransfer(utils.ChainId(1), 0, utils.Nonce(555), [32]byte{1}, mp, sig, big.NewInt(10), make([]byte, 32))

	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	for i := 0; i < TxRetryLimit; i++ {
		s.client.EXPECT().LockAndUpdateOpts().Return(nil)

		s.client.EXPECT().Opts()

		s.client.EXPECT().CallOpts()

		s.bridgeMock.EXPECT().ExecuteProposal(
			gomock.Any(),
			uint8(message.Source),
			uint64(message.DepositNonce),
			[]byte{},
			gomock.Any(),
			[]byte{},
			[]byte{},
			[]byte{},
			gomock.Any(),
			gomock.Any(),
			[]byte{},
			[]byte{},
		).Return(&types.Transaction{}, nil)

		s.client.EXPECT().UnlockOpts()

		s.bridgeMock.EXPECT().GetProposal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(Bridge.BridgeProposal{}, nil)
	}

	go func() {
		err := <-errChn
		s.True(err.Error() == ErrFatalTx.Error())
	}()

	w.executeProposal(message, []byte{}, common.Hash{})

}

func (s *WriterTestSuite) TestExecuteProposalProposalIsFinalizedError() {

	stopChn := make(chan struct{})
	errChn := make(chan error)

	sig := &utils.SignatureVerification{
		AggregatePublicKey: []byte{},
		BlockHash:          common.Hash{},
		Signature:          []byte{},
	}

	mp := &utils.MerkleProof{
		TxRootHash: common.Hash{},
		Key:        []byte{},
		Nodes:      []byte{},
	}

	message := utils.NewFungibleTransfer(utils.ChainId(1), 0, utils.Nonce(555), [32]byte{1}, mp, sig, big.NewInt(10), make([]byte, 32))

	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	for i := 0; i < TxRetryLimit; i++ {
		s.client.EXPECT().LockAndUpdateOpts().Return(nil)

		s.client.EXPECT().Opts()

		s.client.EXPECT().CallOpts()

		s.bridgeMock.EXPECT().ExecuteProposal(
			gomock.Any(),
			uint8(message.Source),
			uint64(message.DepositNonce),
			[]byte{},
			gomock.Any(),
			[]byte{},
			[]byte{},
			[]byte{},
			gomock.Any(),
			gomock.Any(),
			[]byte{},
			[]byte{},
		).Return(nil, errors.New("fail"))

		s.client.EXPECT().UnlockOpts()

		s.bridgeMock.EXPECT().GetProposal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(Bridge.BridgeProposal{}, errors.New("error"))
	}

	go func() {
		err := <-errChn
		s.True(err.Error() == ErrFatalTx.Error())
	}()

	w.executeProposal(message, []byte{}, common.Hash{})

}

func (s *WriterTestSuite) TestExecuteProposalProposalStatusTransferred() {

	stopChn := make(chan struct{})
	errChn := make(chan error)

	sig := &utils.SignatureVerification{
		AggregatePublicKey: []byte{},
		BlockHash:          common.Hash{},
		Signature:          []byte{},
	}

	mp := &utils.MerkleProof{
		TxRootHash: common.Hash{},
		Key:        []byte{},
		Nodes:      []byte{},
	}

	message := utils.NewFungibleTransfer(utils.ChainId(1), 0, utils.Nonce(555), [32]byte{1}, mp, sig, big.NewInt(10), make([]byte, 32))

	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	for i := 0; i < TxRetryLimit; i++ {
		s.client.EXPECT().LockAndUpdateOpts().Return(nil)

		s.client.EXPECT().Opts()

		s.client.EXPECT().CallOpts()

		s.bridgeMock.EXPECT().ExecuteProposal(
			gomock.Any(),
			uint8(message.Source),
			uint64(message.DepositNonce),
			[]byte{},
			gomock.Any(),
			[]byte{},
			[]byte{},
			[]byte{},
			gomock.Any(),
			gomock.Any(),
			[]byte{},
			[]byte{},
		).Return(nil, errors.New("fail"))

		s.client.EXPECT().UnlockOpts()

		s.bridgeMock.EXPECT().GetProposal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(Bridge.BridgeProposal{Status: ProposalStatusTransferred}, nil)
	}

	go func() {
		err := <-errChn
		s.True(err.Error() == ErrFatalTx.Error())
	}()

	w.executeProposal(message, []byte{}, common.Hash{})

}

func (s *WriterTestSuite) TestExecuteProposalProposalStatusCancelled() {

	stopChn := make(chan struct{})
	errChn := make(chan error)

	sig := &utils.SignatureVerification{
		AggregatePublicKey: []byte{},
		BlockHash:          common.Hash{},
		Signature:          []byte{},
	}

	mp := &utils.MerkleProof{
		TxRootHash: common.Hash{},
		Key:        []byte{},
		Nodes:      []byte{},
	}

	message := utils.NewFungibleTransfer(utils.ChainId(1), 0, utils.Nonce(555), [32]byte{1}, mp, sig, big.NewInt(10), make([]byte, 32))

	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	for i := 0; i < TxRetryLimit; i++ {
		s.client.EXPECT().LockAndUpdateOpts().Return(nil)

		s.client.EXPECT().Opts()

		s.client.EXPECT().CallOpts()

		s.bridgeMock.EXPECT().ExecuteProposal(
			gomock.Any(),
			uint8(message.Source),
			uint64(message.DepositNonce),
			[]byte{},
			gomock.Any(),
			[]byte{},
			[]byte{},
			[]byte{},
			gomock.Any(),
			gomock.Any(),
			[]byte{},
			[]byte{},
		).Return(nil, errors.New("fail"))

		s.client.EXPECT().UnlockOpts()

		s.bridgeMock.EXPECT().GetProposal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(Bridge.BridgeProposal{Status: ProposalStatusCancelled}, nil)
	}

	go func() {
		err := <-errChn
		s.True(err.Error() == ErrFatalTx.Error())
	}()

	w.executeProposal(message, []byte{}, common.Hash{})

}

func (s *WriterTestSuite) TestWatchThenExecuteWaitForBlockError() {
	stopChn := make(chan struct{})
	errChn := make(chan error)
	message := utils.NewFungibleTransfer(utils.ChainId(1), 0, utils.Nonce(555), [32]byte{1}, nil, nil, big.NewInt(10), make([]byte, 32))

	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	latestblock := big.NewInt(3)
	for i := 0; i < ExecuteBlockWatchLimit; i++ {
		for waitRetrys := 0; waitRetrys <= BlockRetryLimit; waitRetrys++ {
			s.client.EXPECT().WaitForBlock(latestblock).Return(errors.New("error"))

		}
	}

	go func() {
		err := <-errChn
		s.True(err.Error() == ErrFatalQuery.Error())
	}()

	w.watchThenExecute(message, []byte{}, common.Hash{}, latestblock)
}

func (s *WriterTestSuite) TestWatchThenExecuteFilterLogsError() {
	stopChn := make(chan struct{})
	errChn := make(chan error)
	message := utils.NewFungibleTransfer(utils.ChainId(1), 0, utils.Nonce(555), [32]byte{1}, nil, nil, big.NewInt(10), make([]byte, 32))

	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	latestblock := big.NewInt(3)
	for i := 0; i < ExecuteBlockWatchLimit; i++ {
		for waitRetrys := 0; waitRetrys <= BlockRetryLimit; waitRetrys++ {
			s.client.EXPECT().WaitForBlock(latestblock).Return(nil)

		}

		query := buildQuery(w.cfg.BridgeContract, utils.ProposalEvent, latestblock, latestblock)

		s.client.EXPECT().FilterLogs(context.Background(), query).Return([]types.Log{}, errors.New("error"))

	}

	go func() {
		err := <-errChn
		s.Nil(err)
	}()

	w.watchThenExecute(message, []byte{}, common.Hash{}, latestblock)
}

func (s *WriterTestSuite) TestWatchThenExecuteFilterLogsError2() {
	stopChn := make(chan struct{})
	errChn := make(chan error)
	message := utils.NewFungibleTransfer(utils.ChainId(1), 0, utils.Nonce(555), [32]byte{1}, nil, nil, big.NewInt(10), make([]byte, 32))

	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	latestblock := big.NewInt(3)

	for i := 0; i < ExecuteBlockWatchLimit; i++ {
		for waitRetrys := 0; waitRetrys <= BlockRetryLimit; waitRetrys++ {
			s.client.EXPECT().WaitForBlock(latestblock).Return(nil)

		}
		query := buildQuery(w.cfg.BridgeContract, utils.ProposalEvent, latestblock, latestblock)
		contractAddress := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")

		logs := []types.Log{
			{
				Address: contractAddress,
				// list of topics provided by the contract.
				Topics: []common.Hash{
					utils.Deposit.GetTopic(),
					crypto.Keccak256Hash(big.NewInt(1).Bytes()),
					contractAddress.Hash(),
					crypto.Keccak256Hash(big.NewInt(1).Bytes()),
				},
				Data: []byte{},
			},
		}

		s.client.EXPECT().FilterLogs(context.Background(), query).Return(logs, nil)

	}

	go func() {
		err := <-errChn
		s.Nil(err)
	}()

	w.watchThenExecute(message, []byte{}, common.Hash{}, latestblock)
}

func (s *WriterTestSuite) TestProposalIsFinalizedError() {
	stopChn := make(chan struct{})
	errChn := make(chan error)
	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	hash := crypto.Keccak256Hash([]byte("data"))

	prop := Bridge.BridgeProposal{Status: ProposalStatusPassed}
	s.client.EXPECT().CallOpts().Return(nil)
	s.bridgeMock.EXPECT().GetProposal(gomock.Any(), gomock.Any(), uint64(1), gomock.Any()).Return(prop, errors.New("error"))

	result := w.proposalIsFinalized(utils.ChainId(3), utils.Nonce(1), hash)

	s.False(result)

}

func (s *WriterTestSuite) TestProposalIsFinalizedSuccess() {

	stopChn := make(chan struct{})
	errChn := make(chan error)
	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	hash := crypto.Keccak256Hash([]byte("data"))

	prop := Bridge.BridgeProposal{Status: ProposalStatusTransferred}
	s.client.EXPECT().CallOpts().Return(nil)
	s.bridgeMock.EXPECT().GetProposal(gomock.Any(), gomock.Any(), uint64(1), gomock.Any()).Return(prop, nil)

	result := w.proposalIsFinalized(utils.ChainId(3), utils.Nonce(1), hash)

	s.True(result)

}

func (s *WriterTestSuite) TestProposalIsCompleteError() {
	stopChn := make(chan struct{})
	errChn := make(chan error)
	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	hash := crypto.Keccak256Hash([]byte("data"))

	prop := Bridge.BridgeProposal{Status: ProposalStatusPassed}
	s.client.EXPECT().CallOpts().Return(nil)
	s.bridgeMock.EXPECT().GetProposal(gomock.Any(), gomock.Any(), uint64(1), gomock.Any()).Return(prop, errors.New("error"))

	result := w.proposalIsComplete(utils.ChainId(3), utils.Nonce(1), hash)

	s.False(result)

}

func (s *WriterTestSuite) TestBuildQuery() {

	startBlock := big.NewInt(112233)

	contractAddress := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")

	expected := eth.FilterQuery{
		FromBlock: startBlock,
		ToBlock:   startBlock,
		Addresses: []common.Address{contractAddress},
		Topics: [][]ethcommon.Hash{
			{utils.Deposit.GetTopic()},
		},
	}

	actual := buildQuery(contractAddress, utils.Deposit, startBlock, startBlock)

	s.Equal(expected, actual)

}

func (s *WriterTestSuite) TestCreateERC20ProposalMalformedPayload() {

	stopChn := make(chan struct{})
	errChn := make(chan error)
	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	message := &utils.Message{
		Source:       utils.ChainId(3),
		Destination:  utils.ChainId(3),
		Type:         utils.FungibleTransfer,
		DepositNonce: utils.Nonce(1),
		ResourceId:   [32]byte{},
		MPParams:     nil,
		SVParams:     nil,
		Payload:      []interface{}{},
	}

	result, err := w.createERC20ProposalData(message)

	s.NotNil(err)
	s.Nil(result)

}

func (s *WriterTestSuite) TestCreateERC20ProposalDataWrongAmountFormat() {

	stopChn := make(chan struct{})
	errChn := make(chan error)
	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	contractAddress := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")

	message := &utils.Message{
		Source:       utils.ChainId(3),
		Destination:  utils.ChainId(3),
		Type:         utils.FungibleTransfer,
		DepositNonce: utils.Nonce(1),
		ResourceId:   [32]byte{},
		MPParams:     nil,
		SVParams:     nil,
		Payload: []interface{}{
			uint64(54),
			contractAddress.Bytes(),
		},
	}

	result, err := w.createERC20ProposalData(message)

	s.NotNil(err)
	s.Nil(result)

}

func (s *WriterTestSuite) TestCreateERC20ProposalDataWrongRecipientFormat() {

	stopChn := make(chan struct{})
	errChn := make(chan error)
	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	contractAddress := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")

	message := &utils.Message{
		Source:       utils.ChainId(3),
		Destination:  utils.ChainId(3),
		Type:         utils.FungibleTransfer,
		DepositNonce: utils.Nonce(1),
		ResourceId:   [32]byte{},
		MPParams:     nil,
		SVParams:     nil,
		Payload: []interface{}{
			big.NewInt(76).Bytes(),
			contractAddress,
		},
	}

	result, err := w.createERC20ProposalData(message)

	s.NotNil(err)
	s.Nil(result)

}

func (s *WriterTestSuite) TestCreateERC20ProposalDataComplete() {

	stopChn := make(chan struct{})
	errChn := make(chan error)
	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	contractAddress := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")

	message := &utils.Message{
		Source:       utils.ChainId(3),
		Destination:  utils.ChainId(3),
		Type:         utils.FungibleTransfer,
		DepositNonce: utils.Nonce(1),
		ResourceId:   [32]byte{},
		MPParams:     nil,
		SVParams:     nil,
		Payload: []interface{}{
			[]byte{},
			contractAddress.Bytes(),
		},
	}

	result, err := w.createERC20ProposalData(message)

	s.NotNil(result)
	s.Nil(err)

}

func (s *WriterTestSuite) TestCreateERC21ProposalMalformedPayload() {

	stopChn := make(chan struct{})
	errChn := make(chan error)
	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	message := &utils.Message{
		Source:       utils.ChainId(3),
		Destination:  utils.ChainId(3),
		Type:         utils.FungibleTransfer,
		DepositNonce: utils.Nonce(1),
		ResourceId:   [32]byte{},
		MPParams:     nil,
		SVParams:     nil,
		Payload:      []interface{}{},
	}

	result, err := w.createErc721ProposalData(message)

	s.NotNil(err)
	s.Nil(result)

}

func (s *WriterTestSuite) TestCreateERC21ProposalDataTokenIDFormat() {

	stopChn := make(chan struct{})
	errChn := make(chan error)
	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	contractAddress := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")

	message := &utils.Message{
		Source:       utils.ChainId(3),
		Destination:  utils.ChainId(3),
		Type:         utils.FungibleTransfer,
		DepositNonce: utils.Nonce(1),
		ResourceId:   [32]byte{},
		MPParams:     nil,
		SVParams:     nil,
		Payload: []interface{}{
			contractAddress,
			big.NewInt(76).Bytes(),
			[]byte{},
		},
	}

	result, err := w.createErc721ProposalData(message)

	s.NotNil(err)
	s.Nil(result)

}

func (s *WriterTestSuite) TestCreateERC21ProposalDataRecipientFormat() {

	stopChn := make(chan struct{})
	errChn := make(chan error)
	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	contractAddress := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")

	message := &utils.Message{
		Source:       utils.ChainId(3),
		Destination:  utils.ChainId(3),
		Type:         utils.FungibleTransfer,
		DepositNonce: utils.Nonce(1),
		ResourceId:   [32]byte{},
		MPParams:     nil,
		SVParams:     nil,
		Payload: []interface{}{
			contractAddress.Bytes(),
			"0x",
			[]byte{},
		},
	}

	result, err := w.createErc721ProposalData(message)

	s.NotNil(err)
	s.Nil(result)

}

func (s *WriterTestSuite) TestCreateERC21ProposalDataMetaDataFormat() {

	stopChn := make(chan struct{})
	errChn := make(chan error)
	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	contractAddress := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")

	message := &utils.Message{
		Source:       utils.ChainId(3),
		Destination:  utils.ChainId(3),
		Type:         utils.FungibleTransfer,
		DepositNonce: utils.Nonce(1),
		ResourceId:   [32]byte{},
		MPParams:     nil,
		SVParams:     nil,
		Payload: []interface{}{
			contractAddress.Bytes(),
			[]byte{},
			uint64(65),
		},
	}

	result, err := w.createErc721ProposalData(message)

	s.NotNil(err)
	s.Nil(result)

}

func (s *WriterTestSuite) TestCreateERC21ProposalDataComplete() {

	stopChn := make(chan struct{})
	errChn := make(chan error)
	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	contractAddress := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")

	message := &utils.Message{
		Source:       utils.ChainId(3),
		Destination:  utils.ChainId(3),
		Type:         utils.FungibleTransfer,
		DepositNonce: utils.Nonce(1),
		ResourceId:   [32]byte{},
		MPParams:     nil,
		SVParams:     nil,
		Payload: []interface{}{
			contractAddress.Bytes(),
			[]byte{},
			[]byte{},
		},
	}

	result, err := w.createErc721ProposalData(message)

	s.NotNil(result)
	s.Nil(err)

}

func (s *WriterTestSuite) TestCreateGenericProposalDataMalformedPayload() {

	stopChn := make(chan struct{})
	errChn := make(chan error)
	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	message := &utils.Message{
		Source:       utils.ChainId(3),
		Destination:  utils.ChainId(3),
		Type:         utils.FungibleTransfer,
		DepositNonce: utils.Nonce(1),
		ResourceId:   [32]byte{},
		MPParams:     nil,
		SVParams:     nil,
		Payload:      []interface{}{},
	}

	result, err := w.createGenericDepositProposalData(message)

	s.NotNil(err)
	s.Nil(result)

}

func (s *WriterTestSuite) TestCreateGenericProposalDataWrongMetadataFormat() {

	stopChn := make(chan struct{})
	errChn := make(chan error)
	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	message := &utils.Message{
		Source:       utils.ChainId(3),
		Destination:  utils.ChainId(3),
		Type:         utils.FungibleTransfer,
		DepositNonce: utils.Nonce(1),
		ResourceId:   [32]byte{},
		MPParams:     nil,
		SVParams:     nil,
		Payload: []interface{}{
			"wrong_metadata",
		},
	}

	result, err := w.createGenericDepositProposalData(message)

	s.NotNil(err)
	s.Nil(result)
}

func (s *WriterTestSuite) TestCreateGenericProposalDataComplete() {

	stopChn := make(chan struct{})
	errChn := make(chan error)
	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	message := &utils.Message{
		Source:       utils.ChainId(3),
		Destination:  utils.ChainId(3),
		Type:         utils.FungibleTransfer,
		DepositNonce: utils.Nonce(1),
		ResourceId:   [32]byte{},
		MPParams:     nil,
		SVParams:     nil,
		Payload: []interface{}{
			[]byte{},
		},
	}

	result, err := w.createGenericDepositProposalData(message)

	s.NotNil(result)
	s.Nil(err)
}
