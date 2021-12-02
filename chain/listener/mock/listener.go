// Code generated by MockGen. DO NOT EDIT.
// Source: ./chain/listener/listener.go

// Package mock_listener is a generated GoMock package.
package mock_listener

import (
	utils "github.com/ChainSafe/chainbridge-celo/utils"
	types "github.com/celo-org/celo-blockchain/core/types"
	gomock "github.com/golang/mock/gomock"
	big "math/big"
	reflect "reflect"
)

// MockIRouter is a mock of IRouter interface
type MockIRouter struct {
	ctrl     *gomock.Controller
	recorder *MockIRouterMockRecorder
}

// MockIRouterMockRecorder is the mock recorder for MockIRouter
type MockIRouterMockRecorder struct {
	mock *MockIRouter
}

// NewMockIRouter creates a new mock instance
func NewMockIRouter(ctrl *gomock.Controller) *MockIRouter {
	mock := &MockIRouter{ctrl: ctrl}
	mock.recorder = &MockIRouterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIRouter) EXPECT() *MockIRouterMockRecorder {
	return m.recorder
}

// Send mocks base method
func (m *MockIRouter) Send(msg *utils.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send
func (mr *MockIRouterMockRecorder) Send(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockIRouter)(nil).Send), msg)
}

// MockBlockstorer is a mock of Blockstorer interface
type MockBlockstorer struct {
	ctrl     *gomock.Controller
	recorder *MockBlockstorerMockRecorder
}

// MockBlockstorerMockRecorder is the mock recorder for MockBlockstorer
type MockBlockstorerMockRecorder struct {
	mock *MockBlockstorer
}

// NewMockBlockstorer creates a new mock instance
func NewMockBlockstorer(ctrl *gomock.Controller) *MockBlockstorer {
	mock := &MockBlockstorer{ctrl: ctrl}
	mock.recorder = &MockBlockstorerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBlockstorer) EXPECT() *MockBlockstorerMockRecorder {
	return m.recorder
}

// StoreBlock mocks base method
func (m *MockBlockstorer) StoreBlock(arg0 *big.Int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreBlock", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreBlock indicates an expected call of StoreBlock
func (mr *MockBlockstorerMockRecorder) StoreBlock(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreBlock", reflect.TypeOf((*MockBlockstorer)(nil).StoreBlock), arg0)
}

// MockValidatorsAggregator is a mock of ValidatorsAggregator interface
type MockValidatorsAggregator struct {
	ctrl     *gomock.Controller
	recorder *MockValidatorsAggregatorMockRecorder
}

// MockValidatorsAggregatorMockRecorder is the mock recorder for MockValidatorsAggregator
type MockValidatorsAggregatorMockRecorder struct {
	mock *MockValidatorsAggregator
}

// NewMockValidatorsAggregator creates a new mock instance
func NewMockValidatorsAggregator(ctrl *gomock.Controller) *MockValidatorsAggregator {
	mock := &MockValidatorsAggregator{ctrl: ctrl}
	mock.recorder = &MockValidatorsAggregatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockValidatorsAggregator) EXPECT() *MockValidatorsAggregatorMockRecorder {
	return m.recorder
}

// GetAPKForBlock mocks base method
func (m *MockValidatorsAggregator) GetAPKForBlock(block *big.Int, chainID uint8, epochSize uint64, extra *types.IstanbulExtra) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAPKForBlock", block, chainID, epochSize, extra)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAPKForBlock indicates an expected call of GetAPKForBlock
func (mr *MockValidatorsAggregatorMockRecorder) GetAPKForBlock(block, chainID, epochSize, extra interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAPKForBlock", reflect.TypeOf((*MockValidatorsAggregator)(nil).GetAPKForBlock), block, chainID, epochSize, extra)
}
