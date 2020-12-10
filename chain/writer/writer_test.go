package writer

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"math/big"
	"testing"

	"github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
	"github.com/ChainSafe/chainbridge-celo/chain"
	"github.com/ChainSafe/chainbridge-celo/chain/writer/mock"
	message "github.com/ChainSafe/chainbridge-celo/msg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
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
	m := message.NewFungibleTransfer(1, 0, message.Nonce(555), amount, resourceId, recipient)
	m.Type = "123"
	cfg := &chain.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	s.False(w.ResolveMessage(&m))
}

func (s *WriterTestSuite) TestShouldVoteProposalIsAlreadyComplete() {
	resourceId := [32]byte{1}
	recipient := make([]byte, 32)
	amount := big.NewInt(10)
	stopChn := make(chan struct{})
	errChn := make(chan error)
	m := message.NewFungibleTransfer(1, 0, message.Nonce(555), amount, resourceId, recipient)

	cfg := &chain.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	w := NewWriter(s.client, cfg, stopChn, errChn, nil)
	w.SetBridge(s.bridgeMock)

	// Setting returned proposal to PassedStatus
	prop := Bridge.BridgeProposal{Status: ProposalStatusPassed}
	s.client.EXPECT().CallOpts().Return(nil)
	s.bridgeMock.EXPECT().GetProposal(gomock.Any(), gomock.Any(), uint64(m.DepositNonce), gomock.Any()).Return(prop, nil)
	s.False(w.shouldVote(&m, common.Hash{}))
}

func (s *WriterTestSuite) TestShouldVoteProposalIsAlreadyVoted() {
	stopChn := make(chan struct{})
	errChn := make(chan error)
	m := message.NewFungibleTransfer(1, 0, message.Nonce(555), big.NewInt(10), [32]byte{1}, make([]byte, 32))

	cfg := &chain.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
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
	s.False(w.shouldVote(&m, common.Hash{}))
}

//TestCreateAndExecuteErc20DepositProposal
//TestCreateAndExecuteErc721Proposal
//TestCreateAndExecuteGenericProposal
//TestDuplicateMessage
