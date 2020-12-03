// Code generated by MockGen. DO NOT EDIT.
// Source: ./chain/listener/listener.go

// Package mock_listener is a generated GoMock package.
package mock_listener

import (
	msg "github.com/ChainSafe/chainbridge-utils/msg"
	gomock "github.com/golang/mock/gomock"
	big "math/big"
	reflect "reflect"
)

// MockBlockSyncer is a mock of BlockSyncer interface
type MockBlockSyncer struct {
	ctrl     *gomock.Controller
	recorder *MockBlockSyncerMockRecorder
}

// MockBlockSyncerMockRecorder is the mock recorder for MockBlockSyncer
type MockBlockSyncerMockRecorder struct {
	mock *MockBlockSyncer
}

// NewMockBlockSyncer creates a new mock instance
func NewMockBlockSyncer(ctrl *gomock.Controller) *MockBlockSyncer {
	mock := &MockBlockSyncer{ctrl: ctrl}
	mock.recorder = &MockBlockSyncerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBlockSyncer) EXPECT() *MockBlockSyncerMockRecorder {
	return m.recorder
}

// Sync mocks base method
func (m *MockBlockSyncer) Sync(latestBlock *big.Int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sync", latestBlock)
	ret0, _ := ret[0].(error)
	return ret0
}

// Sync indicates an expected call of Sync
func (mr *MockBlockSyncerMockRecorder) Sync(latestBlock interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sync", reflect.TypeOf((*MockBlockSyncer)(nil).Sync), latestBlock)
}

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
func (m *MockIRouter) Send(msg msg.Message) error {
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