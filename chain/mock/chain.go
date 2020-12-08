// Code generated by MockGen. DO NOT EDIT.
// Source: ./chain/chain.go

// Package mock_chain is a generated GoMock package.
package mock_chain

import (
	context "context"
	Bridge "github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
	ERC20Handler "github.com/ChainSafe/chainbridge-celo/bindings/ERC20Handler"
	ERC721Handler "github.com/ChainSafe/chainbridge-celo/bindings/ERC721Handler"
	GenericHandler "github.com/ChainSafe/chainbridge-celo/bindings/GenericHandler"
	chain "github.com/ChainSafe/chainbridge-celo/chain"
	ethereum "github.com/ethereum/go-ethereum"
	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	common "github.com/ethereum/go-ethereum/common"
	types "github.com/ethereum/go-ethereum/core/types"
	gomock "github.com/golang/mock/gomock"
	big "math/big"
	reflect "reflect"
)

// MockBlockDB is a mock of BlockDB interface
type MockBlockDB struct {
	ctrl     *gomock.Controller
	recorder *MockBlockDBMockRecorder
}

// MockBlockDBMockRecorder is the mock recorder for MockBlockDB
type MockBlockDBMockRecorder struct {
	mock *MockBlockDB
}

// NewMockBlockDB creates a new mock instance
func NewMockBlockDB(ctrl *gomock.Controller) *MockBlockDB {
	mock := &MockBlockDB{ctrl: ctrl}
	mock.recorder = &MockBlockDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBlockDB) EXPECT() *MockBlockDBMockRecorder {
	return m.recorder
}

// StoreBlock mocks base method
func (m *MockBlockDB) StoreBlock(arg0 *big.Int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreBlock", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreBlock indicates an expected call of StoreBlock
func (mr *MockBlockDBMockRecorder) StoreBlock(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreBlock", reflect.TypeOf((*MockBlockDB)(nil).StoreBlock), arg0)
}

// TryLoadLatestBlock mocks base method
func (m *MockBlockDB) TryLoadLatestBlock() (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TryLoadLatestBlock")
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TryLoadLatestBlock indicates an expected call of TryLoadLatestBlock
func (mr *MockBlockDBMockRecorder) TryLoadLatestBlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TryLoadLatestBlock", reflect.TypeOf((*MockBlockDB)(nil).TryLoadLatestBlock))
}

// MockListener is a mock of Listener interface
type MockListener struct {
	ctrl     *gomock.Controller
	recorder *MockListenerMockRecorder
}

// MockListenerMockRecorder is the mock recorder for MockListener
type MockListenerMockRecorder struct {
	mock *MockListener
}

// NewMockListener creates a new mock instance
func NewMockListener(ctrl *gomock.Controller) *MockListener {
	mock := &MockListener{ctrl: ctrl}
	mock.recorder = &MockListenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockListener) EXPECT() *MockListenerMockRecorder {
	return m.recorder
}

// StartPollingBlocks mocks base method
func (m *MockListener) StartPollingBlocks() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartPollingBlocks")
	ret0, _ := ret[0].(error)
	return ret0
}

// StartPollingBlocks indicates an expected call of StartPollingBlocks
func (mr *MockListenerMockRecorder) StartPollingBlocks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartPollingBlocks", reflect.TypeOf((*MockListener)(nil).StartPollingBlocks))
}

// SetContracts mocks base method
func (m *MockListener) SetContracts(bridge *Bridge.Bridge, erc20Handler *ERC20Handler.ERC20Handler, erc721Handler *ERC721Handler.ERC721Handler, genericHandler *GenericHandler.GenericHandler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetContracts", bridge, erc20Handler, erc721Handler, genericHandler)
}

// SetContracts indicates an expected call of SetContracts
func (mr *MockListenerMockRecorder) SetContracts(bridge, erc20Handler, erc721Handler, genericHandler interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetContracts", reflect.TypeOf((*MockListener)(nil).SetContracts), bridge, erc20Handler, erc721Handler, genericHandler)
}

// MockBridger is a mock of Bridger interface
type MockBridger struct {
	ctrl     *gomock.Controller
	recorder *MockBridgerMockRecorder
}

// MockBridgerMockRecorder is the mock recorder for MockBridger
type MockBridgerMockRecorder struct {
	mock *MockBridger
}

// NewMockBridger creates a new mock instance
func NewMockBridger(ctrl *gomock.Controller) *MockBridger {
	mock := &MockBridger{ctrl: ctrl}
	mock.recorder = &MockBridgerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBridger) EXPECT() *MockBridgerMockRecorder {
	return m.recorder
}

// GetProposal mocks base method
func (m *MockBridger) GetProposal(opts *bind.CallOpts, originChainID uint8, depositNonce uint64, dataHash [32]byte) (Bridge.BridgeProposal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProposal", opts, originChainID, depositNonce, dataHash)
	ret0, _ := ret[0].(Bridge.BridgeProposal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProposal indicates an expected call of GetProposal
func (mr *MockBridgerMockRecorder) GetProposal(opts, originChainID, depositNonce, dataHash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProposal", reflect.TypeOf((*MockBridger)(nil).GetProposal), opts, originChainID, depositNonce, dataHash)
}

// HasVotedOnProposal mocks base method
func (m *MockBridger) HasVotedOnProposal(opts *bind.CallOpts, arg0 *big.Int, arg1 [32]byte, arg2 common.Address) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasVotedOnProposal", opts, arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasVotedOnProposal indicates an expected call of HasVotedOnProposal
func (mr *MockBridgerMockRecorder) HasVotedOnProposal(opts, arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasVotedOnProposal", reflect.TypeOf((*MockBridger)(nil).HasVotedOnProposal), opts, arg0, arg1, arg2)
}

// VoteProposal mocks base method
func (m *MockBridger) VoteProposal(opts *bind.TransactOpts, chainID uint8, depositNonce uint64, resourceID, dataHash [32]byte) (*types.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VoteProposal", opts, chainID, depositNonce, resourceID, dataHash)
	ret0, _ := ret[0].(*types.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VoteProposal indicates an expected call of VoteProposal
func (mr *MockBridgerMockRecorder) VoteProposal(opts, chainID, depositNonce, resourceID, dataHash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VoteProposal", reflect.TypeOf((*MockBridger)(nil).VoteProposal), opts, chainID, depositNonce, resourceID, dataHash)
}

// ExecuteProposal mocks base method
func (m *MockBridger) ExecuteProposal(opts *bind.TransactOpts, chainID uint8, depositNonce uint64, data []byte, resourceID [32]byte, signatureHeader, aggregatePublicKey, g1, hashedMessage []byte, rootHash [32]byte, key, nodes []byte) (*types.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecuteProposal", opts, chainID, depositNonce, data, resourceID, signatureHeader, aggregatePublicKey, g1, hashedMessage, rootHash, key, nodes)
	ret0, _ := ret[0].(*types.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecuteProposal indicates an expected call of ExecuteProposal
func (mr *MockBridgerMockRecorder) ExecuteProposal(opts, chainID, depositNonce, data, resourceID, signatureHeader, aggregatePublicKey, g1, hashedMessage, rootHash, key, nodes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecuteProposal", reflect.TypeOf((*MockBridger)(nil).ExecuteProposal), opts, chainID, depositNonce, data, resourceID, signatureHeader, aggregatePublicKey, g1, hashedMessage, rootHash, key, nodes)
}

// MockWriter is a mock of Writer interface
type MockWriter struct {
	ctrl     *gomock.Controller
	recorder *MockWriterMockRecorder
}

// MockWriterMockRecorder is the mock recorder for MockWriter
type MockWriterMockRecorder struct {
	mock *MockWriter
}

// NewMockWriter creates a new mock instance
func NewMockWriter(ctrl *gomock.Controller) *MockWriter {
	mock := &MockWriter{ctrl: ctrl}
	mock.recorder = &MockWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWriter) EXPECT() *MockWriterMockRecorder {
	return m.recorder
}

// SetBridge mocks base method
func (m *MockWriter) SetBridge(bridge chain.Bridger) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetBridge", bridge)
}

// SetBridge indicates an expected call of SetBridge
func (mr *MockWriterMockRecorder) SetBridge(bridge interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBridge", reflect.TypeOf((*MockWriter)(nil).SetBridge), bridge)
}

// MockContractBackendWithBlockFinder is a mock of ContractBackendWithBlockFinder interface
type MockContractBackendWithBlockFinder struct {
	ctrl     *gomock.Controller
	recorder *MockContractBackendWithBlockFinderMockRecorder
}

// MockContractBackendWithBlockFinderMockRecorder is the mock recorder for MockContractBackendWithBlockFinder
type MockContractBackendWithBlockFinderMockRecorder struct {
	mock *MockContractBackendWithBlockFinder
}

// NewMockContractBackendWithBlockFinder creates a new mock instance
func NewMockContractBackendWithBlockFinder(ctrl *gomock.Controller) *MockContractBackendWithBlockFinder {
	mock := &MockContractBackendWithBlockFinder{ctrl: ctrl}
	mock.recorder = &MockContractBackendWithBlockFinderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockContractBackendWithBlockFinder) EXPECT() *MockContractBackendWithBlockFinderMockRecorder {
	return m.recorder
}

// CodeAt mocks base method
func (m *MockContractBackendWithBlockFinder) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CodeAt", ctx, contract, blockNumber)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CodeAt indicates an expected call of CodeAt
func (mr *MockContractBackendWithBlockFinderMockRecorder) CodeAt(ctx, contract, blockNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CodeAt", reflect.TypeOf((*MockContractBackendWithBlockFinder)(nil).CodeAt), ctx, contract, blockNumber)
}

// CallContract mocks base method
func (m *MockContractBackendWithBlockFinder) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CallContract", ctx, call, blockNumber)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CallContract indicates an expected call of CallContract
func (mr *MockContractBackendWithBlockFinderMockRecorder) CallContract(ctx, call, blockNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CallContract", reflect.TypeOf((*MockContractBackendWithBlockFinder)(nil).CallContract), ctx, call, blockNumber)
}

// PendingCodeAt mocks base method
func (m *MockContractBackendWithBlockFinder) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PendingCodeAt", ctx, account)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PendingCodeAt indicates an expected call of PendingCodeAt
func (mr *MockContractBackendWithBlockFinderMockRecorder) PendingCodeAt(ctx, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PendingCodeAt", reflect.TypeOf((*MockContractBackendWithBlockFinder)(nil).PendingCodeAt), ctx, account)
}

// PendingNonceAt mocks base method
func (m *MockContractBackendWithBlockFinder) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PendingNonceAt", ctx, account)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PendingNonceAt indicates an expected call of PendingNonceAt
func (mr *MockContractBackendWithBlockFinderMockRecorder) PendingNonceAt(ctx, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PendingNonceAt", reflect.TypeOf((*MockContractBackendWithBlockFinder)(nil).PendingNonceAt), ctx, account)
}

// SuggestGasPrice mocks base method
func (m *MockContractBackendWithBlockFinder) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SuggestGasPrice", ctx)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SuggestGasPrice indicates an expected call of SuggestGasPrice
func (mr *MockContractBackendWithBlockFinderMockRecorder) SuggestGasPrice(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SuggestGasPrice", reflect.TypeOf((*MockContractBackendWithBlockFinder)(nil).SuggestGasPrice), ctx)
}

// EstimateGas mocks base method
func (m *MockContractBackendWithBlockFinder) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EstimateGas", ctx, call)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EstimateGas indicates an expected call of EstimateGas
func (mr *MockContractBackendWithBlockFinderMockRecorder) EstimateGas(ctx, call interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EstimateGas", reflect.TypeOf((*MockContractBackendWithBlockFinder)(nil).EstimateGas), ctx, call)
}

// SendTransaction mocks base method
func (m *MockContractBackendWithBlockFinder) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendTransaction", ctx, tx)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendTransaction indicates an expected call of SendTransaction
func (mr *MockContractBackendWithBlockFinderMockRecorder) SendTransaction(ctx, tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendTransaction", reflect.TypeOf((*MockContractBackendWithBlockFinder)(nil).SendTransaction), ctx, tx)
}

// FilterLogs mocks base method
func (m *MockContractBackendWithBlockFinder) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilterLogs", ctx, query)
	ret0, _ := ret[0].([]types.Log)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FilterLogs indicates an expected call of FilterLogs
func (mr *MockContractBackendWithBlockFinderMockRecorder) FilterLogs(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilterLogs", reflect.TypeOf((*MockContractBackendWithBlockFinder)(nil).FilterLogs), ctx, query)
}

// SubscribeFilterLogs mocks base method
func (m *MockContractBackendWithBlockFinder) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeFilterLogs", ctx, query, ch)
	ret0, _ := ret[0].(ethereum.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubscribeFilterLogs indicates an expected call of SubscribeFilterLogs
func (mr *MockContractBackendWithBlockFinderMockRecorder) SubscribeFilterLogs(ctx, query, ch interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeFilterLogs", reflect.TypeOf((*MockContractBackendWithBlockFinder)(nil).SubscribeFilterLogs), ctx, query, ch)
}

// LatestBlock mocks base method
func (m *MockContractBackendWithBlockFinder) LatestBlock() (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LatestBlock")
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LatestBlock indicates an expected call of LatestBlock
func (mr *MockContractBackendWithBlockFinderMockRecorder) LatestBlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LatestBlock", reflect.TypeOf((*MockContractBackendWithBlockFinder)(nil).LatestBlock))
}

// MockLogFilterWithLatestBlock is a mock of LogFilterWithLatestBlock interface
type MockLogFilterWithLatestBlock struct {
	ctrl     *gomock.Controller
	recorder *MockLogFilterWithLatestBlockMockRecorder
}

// MockLogFilterWithLatestBlockMockRecorder is the mock recorder for MockLogFilterWithLatestBlock
type MockLogFilterWithLatestBlockMockRecorder struct {
	mock *MockLogFilterWithLatestBlock
}

// NewMockLogFilterWithLatestBlock creates a new mock instance
func NewMockLogFilterWithLatestBlock(ctrl *gomock.Controller) *MockLogFilterWithLatestBlock {
	mock := &MockLogFilterWithLatestBlock{ctrl: ctrl}
	mock.recorder = &MockLogFilterWithLatestBlockMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogFilterWithLatestBlock) EXPECT() *MockLogFilterWithLatestBlockMockRecorder {
	return m.recorder
}

// FilterLogs mocks base method
func (m *MockLogFilterWithLatestBlock) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilterLogs", ctx, q)
	ret0, _ := ret[0].([]types.Log)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FilterLogs indicates an expected call of FilterLogs
func (mr *MockLogFilterWithLatestBlockMockRecorder) FilterLogs(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilterLogs", reflect.TypeOf((*MockLogFilterWithLatestBlock)(nil).FilterLogs), ctx, q)
}

// LatestBlock mocks base method
func (m *MockLogFilterWithLatestBlock) LatestBlock() (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LatestBlock")
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LatestBlock indicates an expected call of LatestBlock
func (mr *MockLogFilterWithLatestBlockMockRecorder) LatestBlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LatestBlock", reflect.TypeOf((*MockLogFilterWithLatestBlock)(nil).LatestBlock))
}
