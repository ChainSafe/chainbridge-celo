// Code generated by MockGen. DO NOT EDIT.
// Source: ./chain/writer/writer.go

// Package mock_writer is a generated GoMock package.
package mock_writer

import (
	context "context"
	Bridge "github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
	ethereum "github.com/ethereum/go-ethereum"
	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	common "github.com/ethereum/go-ethereum/common"
	types "github.com/ethereum/go-ethereum/core/types"
	gomock "github.com/golang/mock/gomock"
	big "math/big"
	reflect "reflect"
)

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
func (m *MockBridger) ExecuteProposal(opts *bind.TransactOpts, chainID uint8, depositNonce uint64, data []byte, resourceID [32]byte, signatureHeader, aggregatePublicKey, g1 []byte, hashedMessage, rootHash [32]byte, key, nodes []byte) (*types.Transaction, error) {
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

// MockContractCaller is a mock of ContractCaller interface
type MockContractCaller struct {
	ctrl     *gomock.Controller
	recorder *MockContractCallerMockRecorder
}

// MockContractCallerMockRecorder is the mock recorder for MockContractCaller
type MockContractCallerMockRecorder struct {
	mock *MockContractCaller
}

// NewMockContractCaller creates a new mock instance
func NewMockContractCaller(ctrl *gomock.Controller) *MockContractCaller {
	mock := &MockContractCaller{ctrl: ctrl}
	mock.recorder = &MockContractCallerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockContractCaller) EXPECT() *MockContractCallerMockRecorder {
	return m.recorder
}

// FilterLogs mocks base method
func (m *MockContractCaller) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilterLogs", ctx, q)
	ret0, _ := ret[0].([]types.Log)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FilterLogs indicates an expected call of FilterLogs
func (mr *MockContractCallerMockRecorder) FilterLogs(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilterLogs", reflect.TypeOf((*MockContractCaller)(nil).FilterLogs), ctx, q)
}

// LatestBlock mocks base method
func (m *MockContractCaller) LatestBlock() (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LatestBlock")
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LatestBlock indicates an expected call of LatestBlock
func (mr *MockContractCallerMockRecorder) LatestBlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LatestBlock", reflect.TypeOf((*MockContractCaller)(nil).LatestBlock))
}

// BlockByNumber mocks base method
func (m *MockContractCaller) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockByNumber", ctx, number)
	ret0, _ := ret[0].(*types.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BlockByNumber indicates an expected call of BlockByNumber
func (mr *MockContractCallerMockRecorder) BlockByNumber(ctx, number interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockByNumber", reflect.TypeOf((*MockContractCaller)(nil).BlockByNumber), ctx, number)
}

// CallOpts mocks base method
func (m *MockContractCaller) CallOpts() *bind.CallOpts {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CallOpts")
	ret0, _ := ret[0].(*bind.CallOpts)
	return ret0
}

// CallOpts indicates an expected call of CallOpts
func (mr *MockContractCallerMockRecorder) CallOpts() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CallOpts", reflect.TypeOf((*MockContractCaller)(nil).CallOpts))
}

// Opts mocks base method
func (m *MockContractCaller) Opts() *bind.TransactOpts {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Opts")
	ret0, _ := ret[0].(*bind.TransactOpts)
	return ret0
}

// Opts indicates an expected call of Opts
func (mr *MockContractCallerMockRecorder) Opts() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Opts", reflect.TypeOf((*MockContractCaller)(nil).Opts))
}

// LockAndUpdateOpts mocks base method
func (m *MockContractCaller) LockAndUpdateOpts() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LockAndUpdateOpts")
	ret0, _ := ret[0].(error)
	return ret0
}

// LockAndUpdateOpts indicates an expected call of LockAndUpdateOpts
func (mr *MockContractCallerMockRecorder) LockAndUpdateOpts() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LockAndUpdateOpts", reflect.TypeOf((*MockContractCaller)(nil).LockAndUpdateOpts))
}

// UnlockOpts mocks base method
func (m *MockContractCaller) UnlockOpts() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UnlockOpts")
}

// UnlockOpts indicates an expected call of UnlockOpts
func (mr *MockContractCallerMockRecorder) UnlockOpts() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnlockOpts", reflect.TypeOf((*MockContractCaller)(nil).UnlockOpts))
}

// WaitForBlock mocks base method
func (m *MockContractCaller) WaitForBlock(block *big.Int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitForBlock", block)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitForBlock indicates an expected call of WaitForBlock
func (mr *MockContractCallerMockRecorder) WaitForBlock(block interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitForBlock", reflect.TypeOf((*MockContractCaller)(nil).WaitForBlock), block)
}
