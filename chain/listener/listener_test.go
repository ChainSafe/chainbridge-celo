// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only
package listener

import (
	"math/big"
	"testing"

	"github.com/ChainSafe/chainbridge-celo/bindings/mptp/ERC20Handler"
	"github.com/ChainSafe/chainbridge-celo/bindings/mptp/ERC721Handler"
	"github.com/ChainSafe/chainbridge-celo/bindings/mptp/GenericHandler"
	"github.com/ChainSafe/chainbridge-celo/chain/client/mock"
	"github.com/ChainSafe/chainbridge-celo/chain/config"
	"github.com/ChainSafe/chainbridge-celo/chain/listener/mock"
	"github.com/ChainSafe/chainbridge-celo/txtrie"
	"github.com/ChainSafe/chainbridge-celo/utils"
	eth "github.com/celo-org/celo-blockchain"
	"github.com/celo-org/celo-blockchain/accounts/abi/bind"
	"github.com/celo-org/celo-blockchain/common"
	ethcommon "github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/core/types"
	"github.com/celo-org/celo-blockchain/crypto"
	"github.com/celo-org/celo-blockchain/rlp"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type ListenerTestSuite struct {
	suite.Suite
	routerMock               *mock_listener.MockIRouter
	clientMock               *mock_client.MockLogFilterWithLatestBlock
	blockStorerMock          *mock_listener.MockBlockstorer
	gomockController         *gomock.Controller
	bridge                   *mock_listener.MockIBridge
	erc20Handler             *mock_listener.MockIERC20Handler
	erc721Handler            *mock_listener.MockIERC721Handler
	genericHandler           *mock_listener.MockIGenericHandler
	validatorsAggregatorMock *mock_listener.MockValidatorsAggregator
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(ListenerTestSuite))
}

func (s *ListenerTestSuite) SetupSuite()    {}
func (s *ListenerTestSuite) TearDownSuite() {}
func (s *ListenerTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.routerMock = mock_listener.NewMockIRouter(gomockController)
	s.clientMock = mock_client.NewMockLogFilterWithLatestBlock(gomockController)
	s.blockStorerMock = mock_listener.NewMockBlockstorer(gomockController)
	s.gomockController = gomockController
	s.bridge = mock_listener.NewMockIBridge(gomockController)
	s.erc20Handler = mock_listener.NewMockIERC20Handler(gomockController)
	s.erc721Handler = mock_listener.NewMockIERC721Handler(gomockController)
	s.genericHandler = mock_listener.NewMockIGenericHandler(gomockController)
	s.validatorsAggregatorMock = mock_listener.NewMockValidatorsAggregator(gomockController)
}
func (s *ListenerTestSuite) TearDownTest() {}

func dummyBlock(number int64) *types.Block {
	header := &types.Header{
		Number:  big.NewInt(number),
		GasUsed: 123213,
		Time:    100,
		Extra:   []byte{01, 02},
	}
	feeCurrencyAddr := common.HexToAddress("02")
	gatewayFeeRecipientAddr := common.HexToAddress("03")
	tx := types.NewTransaction(1, common.HexToAddress("01"), big.NewInt(1), 10000, big.NewInt(10), &feeCurrencyAddr, &gatewayFeeRecipientAddr, big.NewInt(34), []byte{04})
	return types.NewBlock(header, []*types.Transaction{tx}, nil, nil)
}

func (s *ListenerTestSuite) TestListenerStartStop() {
	stopChn := make(chan struct{})
	errChn := make(chan error)

	l := NewListener(&config.CeloChainConfig{StartBlock: big.NewInt(1)}, s.clientMock, s.blockStorerMock, stopChn, errChn, s.routerMock, s.validatorsAggregatorMock)
	close(stopChn)
	s.NotNil(l.pollBlocks())
}

func (s *ListenerTestSuite) TestLatestBlockUpdate() {
	stopChn := make(chan struct{})
	errChn := make(chan error)
	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	l := NewListener(cfg, s.clientMock, s.blockStorerMock, stopChn, errChn, s.routerMock, s.validatorsAggregatorMock)

	s.clientMock.EXPECT().LatestBlock().Return(big.NewInt(555), nil)
	//No event logs found
	s.clientMock.EXPECT().FilterLogs(gomock.Any(), gomock.Any()).Return(make([]types.Log, 0), nil)

	s.blockStorerMock.EXPECT().StoreBlock(big.NewInt(1))

	//ON second call to latest block we stopping goroutine
	s.clientMock.EXPECT().LatestBlock().DoAndReturn(func() (*big.Int, error) { close(stopChn); return nil, errors.New("err") })

	s.NotNil(l.pollBlocks())
	s.Equal(cfg.StartBlock.String(), "2")
}

func (s *ListenerTestSuite) TestHandleErc20DepositedEventSucccess() {

	stopChn := make(chan struct{})
	errChn := make(chan error)
	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	listener := NewListener(cfg, s.clientMock, s.blockStorerMock, stopChn, errChn, s.routerMock, s.validatorsAggregatorMock)

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
	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	listener := NewListener(cfg, s.clientMock, s.blockStorerMock, stopChn, errChn, s.routerMock, s.validatorsAggregatorMock)

	listener.SetContracts(s.bridge, s.erc20Handler, s.erc721Handler, s.genericHandler)

	s.erc20Handler.EXPECT().GetDepositRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return(ERC20Handler.ERC20HandlerDepositRecord{}, errors.New("error occured"))

	_, err := listener.handleErc20DepositedEvent(3, 0)

	s.NotNil(err)

}

func (s *ListenerTestSuite) TestHandleErc721DepositedEventSuccess() {

	stopChn := make(chan struct{})
	errChn := make(chan error)

	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	listener := NewListener(cfg, s.clientMock, s.blockStorerMock, stopChn, errChn, s.routerMock, s.validatorsAggregatorMock)

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

	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	listener := NewListener(cfg, s.clientMock, s.blockStorerMock, stopChn, errChn, s.routerMock, s.validatorsAggregatorMock)

	listener.SetContracts(s.bridge, s.erc20Handler, s.erc721Handler, s.genericHandler)

	s.erc721Handler.EXPECT().GetDepositRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return(ERC721Handler.ERC721HandlerDepositRecord{}, errors.New("error occured"))

	_, err := listener.handleErc721DepositedEvent(3, 0)

	s.NotNil(err)

}

func (s *ListenerTestSuite) TestHandleGenericDepositedEventSuccess() {

	stopChn := make(chan struct{})
	errChn := make(chan error)

	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	listener := NewListener(cfg, s.clientMock, s.blockStorerMock, stopChn, errChn, s.routerMock, s.validatorsAggregatorMock)

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

	cfg := &config.CeloChainConfig{StartBlock: big.NewInt(1), BridgeContract: common.Address{}}
	listener := NewListener(cfg, s.clientMock, s.blockStorerMock, stopChn, errChn, s.routerMock, s.validatorsAggregatorMock)

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

	cfg := &config.CeloChainConfig{
		ID:                   3,
		Erc20HandlerContract: erc20HandlerContractaddress,
		StartBlock:           startBlock,
		BridgeContract:       bridgeContract,
	}

	listener := NewListener(cfg, s.clientMock, s.blockStorerMock, stopChn, errChn, s.routerMock, s.validatorsAggregatorMock)

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
			Data:    []byte{},
			TxIndex: 1,
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

	nonce := utils.Nonce(logs[0].Topics[3].Big().Uint64())

	destID := utils.ChainId(logs[0].Topics[1].Big().Uint64())
	pk := []byte{0x1f}
	s.validatorsAggregatorMock.EXPECT().GetAPKForBlock(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte{0x1f}, nil)
	block := dummyBlock(123)
	s.clientMock.EXPECT().BlockByNumber(gomock.Any(), gomock.Any()).Return(block, nil)

	trie, err := txtrie.CreateNewTrie(block.TxHash(), block.Transactions())
	s.Nil(err)
	//s.Nil(trie.CreateNewTrie(block.TxHash(), block.Transactions()))

	keyRlp, err := rlp.EncodeToBytes(uint(1))
	s.Nil(err)

	proof, key, err := txtrie.RetrieveProof(trie, keyRlp)
	s.Nil(err)

	_ = utils.NewFungibleTransfer(
		listener.cfg.ID,
		destID,
		nonce,
		prop.ResourceID,
		&utils.MerkleProof{
			TxRootHash: block.TxHash(),
			Nodes:      proof,
			Key:        key,
		},
		&utils.SignatureVerification{
			AggregatePublicKey: pk,
			BlockHash:          block.Header().Hash(),
			Signature:          block.EpochSnarkData().Signature,
		},
		prop.Amount,
		prop.DestinationRecipientAddress,
	)

	s.routerMock.EXPECT().Send(gomock.Any()).Times(1).Return(nil)

	transactions := GetTransactions()
	txRoox, _ := GetTxRoot(transactions)

	header := &types.Header{
		TxHash: txRoox,
	}

	block = types.NewBlock(header, transactions, nil, nil)

	s.clientMock.EXPECT().BlockByNumber(context.TODO(), gomock.Any()).Return(block, nil)

	err = listener.getDepositEventsAndProofsForBlock(big.NewInt(112233))

	s.Nil(err)
}

func (s *ListenerTestSuite) TestGetDepositEventsAndProofsForBlockerERC721() {

	stopChn := make(chan struct{})
	errChn := make(chan error)

	startBlock := big.NewInt(112233)

	address := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")
	bridgeContract := address

	erc721HandlerContractaddress := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")

	cfg := &config.CeloChainConfig{
		ID:                    3,
		Erc721HandlerContract: erc721HandlerContractaddress,
		StartBlock:            startBlock,
		BridgeContract:        bridgeContract,
	}

	listener := NewListener(cfg, s.clientMock, s.blockStorerMock, stopChn, errChn, s.routerMock, s.validatorsAggregatorMock)

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

	nonce := utils.Nonce(logs[0].Topics[3].Big().Uint64())

	destID := utils.ChainId(logs[0].Topics[1].Big().Uint64())
	pk := []byte{0x1f}
	s.validatorsAggregatorMock.EXPECT().GetAPKForBlock(gomock.Any(), gomock.Any(), gomock.Any()).Return(pk, nil)
	block := dummyBlock(123)
	s.clientMock.EXPECT().BlockByNumber(gomock.Any(), gomock.Any()).Return(block, nil)
	_ = utils.NewNonFungibleTransfer(
		listener.cfg.ID,
		destID,
		nonce,
		prop.ResourceID,
		&utils.MerkleProof{
			TxRootHash: block.TxHash(),
		},
		&utils.SignatureVerification{
			AggregatePublicKey: pk,
			BlockHash:          block.Header().Hash(),
			Signature:          block.EpochSnarkData().Signature,
		},
		prop.TokenID,
		prop.DestinationRecipientAddress,
		prop.MetaData,
	)

	s.routerMock.EXPECT().Send(gomock.Any()).Times(1).Return(nil)

	err := listener.getDepositEventsAndProofsForBlock(big.NewInt(112233))

	s.Nil(err)

}

func (s *ListenerTestSuite) TestGetDepositEventsAndProofsForBlockerGeneric() {

	stopChn := make(chan struct{})
	errChn := make(chan error)

	startBlock := big.NewInt(112233)

	address := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")
	bridgeContract := address

	genericContractaddress := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")

	cfg := &config.CeloChainConfig{
		ID:                     3,
		GenericHandlerContract: genericContractaddress,
		StartBlock:             startBlock,
		BridgeContract:         bridgeContract,
	}

	listener := NewListener(cfg, s.clientMock, s.blockStorerMock, stopChn, errChn, s.routerMock, s.validatorsAggregatorMock)

	listener.SetContracts(s.bridge, s.erc20Handler, s.erc721Handler, s.genericHandler)

	query := buildQuery(address, utils.Deposit, startBlock, startBlock)

	logs := []types.Log{
		{
			Address: listener.cfg.GenericHandlerContract,
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

	prop := GenericHandler.GenericHandlerDepositRecord{
		DestinationChainID: 1,
		ResourceID:         [32]byte{},
		Depositer:          address,
		MetaData:           []byte{},
	}

	s.genericHandler.EXPECT().GetDepositRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return(prop, nil)

	s.bridge.EXPECT().ResourceIDToHandlerAddress(&bind.CallOpts{}, [32]byte(listener.cfg.GenericHandlerContract.Hash())).Return(listener.cfg.GenericHandlerContract, nil)

	nonce := utils.Nonce(logs[0].Topics[3].Big().Uint64())

	destID := utils.ChainId(logs[0].Topics[1].Big().Uint64())
	pk := []byte{0x1f}
	s.validatorsAggregatorMock.EXPECT().GetAPKForBlock(gomock.Any(), gomock.Any(), gomock.Any()).Return(pk, nil)
	block := dummyBlock(123)
	s.clientMock.EXPECT().BlockByNumber(gomock.Any(), gomock.Any()).Return(block, nil)
	_ = utils.NewGenericTransfer(
		listener.cfg.ID,
		destID,
		nonce,
		prop.ResourceID,
		&utils.MerkleProof{
			TxRootHash: block.TxHash(),
		},
		&utils.SignatureVerification{
			AggregatePublicKey: pk,
			BlockHash:          block.Header().Hash(),
			Signature:          block.EpochSnarkData().Signature,
		},
		prop.MetaData,
	)

	s.routerMock.EXPECT().Send(gomock.Any()).Times(1).Return(nil)

	err := listener.getDepositEventsAndProofsForBlock(big.NewInt(112233))

	s.Nil(err)

}

func (s *ListenerTestSuite) TestGetDepositEventsAndProofsForBlockerFailure() {

	stopChn := make(chan struct{})
	errChn := make(chan error)

	startBlock := big.NewInt(112233)

	contractAddress := common.HexToAddress("0x67C7656EC7ab88b098defB751B7401B5f6d8976F")

	handlerContractaddress := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976C")

	bridgeContract := contractAddress

	cfg := &config.CeloChainConfig{
		ID:                     3,
		Erc20HandlerContract:   handlerContractaddress,
		GenericHandlerContract: handlerContractaddress,
		Erc721HandlerContract:  handlerContractaddress,
		StartBlock:             startBlock,
		BridgeContract:         bridgeContract,
	}

	//cfg := &chain.CeloChainConfig{StartBlock: startBlock, BridgeContract: bridgeContract}

	listener := NewListener(cfg, s.clientMock, s.blockStorerMock, stopChn, errChn, s.routerMock, s.validatorsAggregatorMock)

	listener.SetContracts(s.bridge, s.erc20Handler, s.erc721Handler, s.genericHandler)

	query := buildQuery(contractAddress, utils.Deposit, startBlock, startBlock)

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

	s.clientMock.EXPECT().FilterLogs(context.Background(), query).Return(logs, nil)

	prop := GenericHandler.GenericHandlerDepositRecord{
		DestinationChainID: 1,
		ResourceID:         [32]byte{},
		Depositer:          contractAddress,
		MetaData:           []byte{},
	}

	//should not be called
	s.genericHandler.EXPECT().GetDepositRecord(gomock.Any(), gomock.Any(), gomock.Any()).Times(0).Return(prop, nil)

	s.bridge.EXPECT().ResourceIDToHandlerAddress(&bind.CallOpts{}, [32]byte(contractAddress.Hash())).Return(contractAddress, nil)
	block := dummyBlock(123)
	s.clientMock.EXPECT().BlockByNumber(gomock.Any(), gomock.Any()).Return(block, nil)

	nonce := utils.Nonce(logs[0].Topics[3].Big().Uint64())

	destID := utils.ChainId(logs[0].Topics[1].Big().Uint64())

	_ = utils.NewGenericTransfer(
		listener.cfg.ID,
		destID,
		nonce,
		prop.ResourceID,
		nil,
		nil,
		prop.MetaData,
	)

	//should not be called
	s.routerMock.EXPECT().Send(gomock.Any()).Times(0).Return(nil)

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
