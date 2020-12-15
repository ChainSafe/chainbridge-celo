// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only
package listener

import (
	"math/big"
	"testing"

	mock_chain "github.com/ChainSafe/chainbridge-celo/chain/mock"
	"github.com/ChainSafe/chainbridge-celo/msg"
	utils "github.com/ChainSafe/chainbridge-celo/shared/ethereum"
	eth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	ERC20Handler "github.com/ChainSafe/chainbridge-celo/bindings/ERC20Handler"
	"github.com/ChainSafe/chainbridge-celo/bindings/ERC721Handler"
	"github.com/ChainSafe/chainbridge-celo/bindings/GenericHandler"
	"github.com/ChainSafe/chainbridge-celo/chain"
	mock_listener "github.com/ChainSafe/chainbridge-celo/chain/listener/mock"
	"github.com/ethereum/go-ethereum/crypto"
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
	bridge           *mock_listener.MockIBridge
	erc20Handler     *mock_listener.MockIERC20Handler
	erc721Handler    *mock_listener.MockIERC721Handler
	genericHandler   *mock_listener.MockIGenericHandler
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(ListenerTestSuite))
}

func (s *ListenerTestSuite) SetupSuite() {
}

func (s *ListenerTestSuite) TearDownSuite() {}
func (s *ListenerTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.syncerMock = mock_listener.NewMockBlockSyncer(gomockController)
	s.routerMock = mock_listener.NewMockIRouter(gomockController)
	s.clientMock = mock_chain.NewMockLogFilterWithLatestBlock(gomockController)
	s.blockStorerMock = mock_listener.NewMockBlockstorer(gomockController)
	s.gomockController = gomockController
	s.bridge = mock_listener.NewMockIBridge(gomockController)
	s.erc20Handler = mock_listener.NewMockIERC20Handler(gomockController)
	s.erc721Handler = mock_listener.NewMockIERC721Handler(gomockController)
	s.genericHandler = mock_listener.NewMockIGenericHandler(gomockController)
}
func (s *ListenerTestSuite) TearDownTest() {}

func (s *ListenerTestSuite) TestListenerStartStop() {
	stopChn := make(chan struct{})
	errChn := make(chan error)

	l := NewListener(&chain.CeloChainConfig{StartBlock: big.NewInt(1)}, s.clientMock, s.blockStorerMock, stopChn, errChn, s.syncerMock, s.routerMock)
	close(stopChn)
	s.NotNil(l.pollBlocks())
}

func (s *ListenerTestSuite) TestLatestBlockUpdateTest() {
	stopChn := make(chan struct{})
	errChn := make(chan error)
	cfg := &chain.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	l := NewListener(cfg, s.clientMock, s.blockStorerMock, stopChn, errChn, s.syncerMock, s.routerMock)
	s.clientMock.EXPECT().LatestBlock().Return(big.NewInt(555), nil)
	s.syncerMock.EXPECT().Sync(gomock.Any()).Return(nil)
	//No event logs found
	s.clientMock.EXPECT().FilterLogs(gomock.Any(), gomock.Any()).Return(make([]types.Log, 0), nil)
	// getDepositEventsAndProofsForBlock todo
	s.blockStorerMock.EXPECT().StoreBlock(big.NewInt(1))

	//ON second call to latest block we stopping goroutine
	s.clientMock.EXPECT().LatestBlock().DoAndReturn(func() (*big.Int, error) { close(stopChn); return nil, errors.New("err") })

	s.NotNil(l.pollBlocks())
	s.Equal(cfg.StartBlock.String(), "2")
}

func (s *ListenerTestSuite) TestHandleErc20DepositedEventSucccess() {

	stopChn := make(chan struct{})
	errChn := make(chan error)
	cfg := &chain.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	listener := NewListener(cfg, s.clientMock, s.blockStorerMock, stopChn, errChn, s.syncerMock, s.routerMock)

	listener.SetContracts(s.bridge, s.erc20Handler, s.erc721Handler, s.genericHandler)

	tokenAddress := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")

	prop := ERC20Handler.ERC20HandlerDepositRecord{
		TokenAddress:                tokenAddress,
		DestinationChainID:          1,
		ResourceID:                  [32]byte{},
		DestinationRecipientAddress: []byte{},
		Depositer:                   tokenAddress,
		Amount:                      big.NewInt(1),
	}

	s.erc20Handler.EXPECT().GetDepositRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return(prop, nil)

	res, err := listener.handleErc20DepositedEvent(3, 0)

	s.NotNil(res)

	s.Nil(err)

}

func (s *ListenerTestSuite) TestHandleErc20DepositedEventFailure() {

	stopChn := make(chan struct{})
	errChn := make(chan error)
	cfg := &chain.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	listener := NewListener(cfg, s.clientMock, s.blockStorerMock, stopChn, errChn, s.syncerMock, s.routerMock)

	listener.SetContracts(s.bridge, s.erc20Handler, s.erc721Handler, s.genericHandler)

	s.erc20Handler.EXPECT().GetDepositRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return(ERC20Handler.ERC20HandlerDepositRecord{}, errors.New("error occured"))

	_, err := listener.handleErc20DepositedEvent(3, 0)

	s.NotNil(err)

}

func (s *ListenerTestSuite) TestHandleErc721DepositedEventSuccess() {

	stopChn := make(chan struct{})
	errChn := make(chan error)

	cfg := &chain.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	listener := NewListener(cfg, s.clientMock, s.blockStorerMock, stopChn, errChn, s.syncerMock, s.routerMock)

	listener.SetContracts(s.bridge, s.erc20Handler, s.erc721Handler, s.genericHandler)

	tokenAddress := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")

	prop := ERC721Handler.ERC721HandlerDepositRecord{
		TokenAddress:                tokenAddress,
		DestinationChainID:          1,
		ResourceID:                  [32]byte{},
		DestinationRecipientAddress: []byte{},
		Depositer:                   tokenAddress,
		TokenID:                     big.NewInt(1),
		MetaData:                    []byte{},
	}

	s.erc721Handler.EXPECT().GetDepositRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return(prop, nil)

	res, err := listener.handleErc721DepositedEvent(3, 0)

	s.NotNil(res)

	s.Nil(err)

}

func (s *ListenerTestSuite) TestHandleErc721DepositedEventFailure() {

	stopChn := make(chan struct{})
	errChn := make(chan error)

	cfg := &chain.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	listener := NewListener(cfg, s.clientMock, s.blockStorerMock, stopChn, errChn, s.syncerMock, s.routerMock)

	listener.SetContracts(s.bridge, s.erc20Handler, s.erc721Handler, s.genericHandler)

	s.erc721Handler.EXPECT().GetDepositRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return(ERC721Handler.ERC721HandlerDepositRecord{}, errors.New("error occured"))

	_, err := listener.handleErc721DepositedEvent(3, 0)

	s.NotNil(err)

}

func (s *ListenerTestSuite) TestHandleGenericDepositedEventSuccess() {

	stopChn := make(chan struct{})
	errChn := make(chan error)

	cfg := &chain.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	listener := NewListener(cfg, s.clientMock, s.blockStorerMock, stopChn, errChn, s.syncerMock, s.routerMock)

	listener.SetContracts(s.bridge, s.erc20Handler, s.erc721Handler, s.genericHandler)

	tokenAddress := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")

	prop := GenericHandler.GenericHandlerDepositRecord{
		DestinationChainID: 1,
		Depositer:          tokenAddress,
		ResourceID:         [32]byte{},
		MetaData:           []byte{},
	}

	s.genericHandler.EXPECT().GetDepositRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return(prop, nil)

	res, err := listener.handleGenericDepositedEvent(3, 0)

	s.NotNil(res)

	s.Nil(err)

}

func (s *ListenerTestSuite) TestHandleGenericDepositedEventFailure() {

	stopChn := make(chan struct{})
	errChn := make(chan error)

	cfg := &chain.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	listener := NewListener(cfg, s.clientMock, s.blockStorerMock, stopChn, errChn, s.syncerMock, s.routerMock)

	listener.SetContracts(s.bridge, s.erc20Handler, s.erc721Handler, s.genericHandler)

	s.genericHandler.EXPECT().GetDepositRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return(GenericHandler.GenericHandlerDepositRecord{}, errors.New("error occured"))

	_, err := listener.handleGenericDepositedEvent(3, 0)

	s.NotNil(err)

}

func (s *ListenerTestSuite) TestGetDepositEventsAndProofsForBlockerERC20() {

	stopChn := make(chan struct{})
	errChn := make(chan error)

	startBlock := big.NewInt(112233)

	address := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")
	bridgeContract := address

	erc20HandlerContractaddress := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")

	cfg := &chain.CeloChainConfig{
		ID:                   3,
		Erc20HandlerContract: erc20HandlerContractaddress,
		StartBlock:           startBlock,
		BridgeContract:       bridgeContract,
	}

	listener := NewListener(cfg, s.clientMock, s.blockStorerMock, stopChn, errChn, s.syncerMock, s.routerMock)

	listener.SetContracts(s.bridge, s.erc20Handler, s.erc721Handler, s.genericHandler)

	query := buildQuery(address, utils.Deposit, startBlock, startBlock)

	logs := []types.Log{
		{
			Address: listener.cfg.Erc20HandlerContract,
			// list of topics provided by the contract.
			Topics: []common.Hash{
				utils.Deposit.GetTopic(),
				crypto.Keccak256Hash(big.NewInt(1).Bytes()),
				address.Hash(),
				crypto.Keccak256Hash(big.NewInt(1).Bytes()),
			},
			Data: []byte{},
		},
	}

	s.clientMock.EXPECT().FilterLogs(context.Background(), query).Return(logs, nil)

	prop := ERC20Handler.ERC20HandlerDepositRecord{
		TokenAddress:                listener.cfg.Erc20HandlerContract,
		DestinationChainID:          1,
		ResourceID:                  [32]byte{},
		DestinationRecipientAddress: []byte{},
		Depositer:                   address,
		Amount:                      big.NewInt(1),
	}

	s.erc20Handler.EXPECT().GetDepositRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return(prop, nil)

	s.bridge.EXPECT().ResourceIDToHandlerAddress(&bind.CallOpts{}, [32]byte(listener.cfg.Erc20HandlerContract.Hash())).Return(listener.cfg.Erc20HandlerContract, nil)

	nonce := msg.Nonce(logs[0].Topics[3].Big().Uint64())

	destID := msg.ChainId(logs[0].Topics[1].Big().Uint64())

	message := msg.NewFungibleTransfer(
		listener.cfg.ID,
		destID,
		nonce,
		prop.Amount,
		prop.ResourceID,
		prop.DestinationRecipientAddress,
		nil,
	)

	s.routerMock.EXPECT().Send(message).Times(1).Return(nil)

	err := listener.getDepositEventsAndProofsForBlock(big.NewInt(112233))

	s.Nil(err)

}

func (s *ListenerTestSuite) TestGetDepositEventsAndProofsForBlockerERC721() {

	stopChn := make(chan struct{})
	errChn := make(chan error)

	startBlock := big.NewInt(112233)

	address := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")
	bridgeContract := address

	erc721HandlerContractaddress := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")

	cfg := &chain.CeloChainConfig{
		ID:                    3,
		Erc721HandlerContract: erc721HandlerContractaddress,
		StartBlock:            startBlock,
		BridgeContract:        bridgeContract,
	}

	listener := NewListener(cfg, s.clientMock, s.blockStorerMock, stopChn, errChn, s.syncerMock, s.routerMock)

	listener.SetContracts(s.bridge, s.erc20Handler, s.erc721Handler, s.genericHandler)

	query := buildQuery(address, utils.Deposit, startBlock, startBlock)

	logs := []types.Log{
		{
			Address: listener.cfg.Erc721HandlerContract,
			// list of topics provided by the contract.
			Topics: []common.Hash{
				utils.Deposit.GetTopic(),
				crypto.Keccak256Hash(big.NewInt(1).Bytes()),
				address.Hash(),
				crypto.Keccak256Hash(big.NewInt(1).Bytes()),
			},
			Data: []byte{},
		},
	}

	s.clientMock.EXPECT().FilterLogs(context.Background(), query).Return(logs, nil)

	prop := ERC721Handler.ERC721HandlerDepositRecord{
		TokenAddress:                listener.cfg.Erc721HandlerContract,
		DestinationChainID:          1,
		ResourceID:                  [32]byte{},
		DestinationRecipientAddress: []byte{},
		Depositer:                   address,
		TokenID:                     big.NewInt(1),
		MetaData:                    []byte{},
	}

	s.erc721Handler.EXPECT().GetDepositRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return(prop, nil)

	s.bridge.EXPECT().ResourceIDToHandlerAddress(&bind.CallOpts{}, [32]byte(listener.cfg.Erc721HandlerContract.Hash())).Return(listener.cfg.Erc721HandlerContract, nil)

	nonce := msg.Nonce(logs[0].Topics[3].Big().Uint64())

	destID := msg.ChainId(logs[0].Topics[1].Big().Uint64())

	message := msg.NewNonFungibleTransfer(
		listener.cfg.ID,
		destID,
		nonce,
		prop.ResourceID,
		prop.TokenID,
		prop.DestinationRecipientAddress,
		prop.MetaData,
		nil,
	)

	s.routerMock.EXPECT().Send(message).Times(1).Return(nil)

	err := listener.getDepositEventsAndProofsForBlock(big.NewInt(112233))

	s.Nil(err)

}

func (s *ListenerTestSuite) TestBuildQuery() {

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
