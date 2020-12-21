// Code generated by MockGen. DO NOT EDIT.
// Source: ./chain/chain.go

// Package mock_chain is a generated GoMock package.
package mock_chain

import (
	"github.com/ChainSafe/chainbridge-celo/chain/listener"
	writer "github.com/ChainSafe/chainbridge-celo/chain/writer"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

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
func (m *MockListener) SetContracts(bridge listener.IBridge, erc20Handler listener.IERC20Handler, erc721Handler listener.IERC721Handler, genericHandler listener.IGenericHandler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetContracts", bridge, erc20Handler, erc721Handler, genericHandler)
}

// SetContracts indicates an expected call of SetContracts
func (mr *MockListenerMockRecorder) SetContracts(bridge, erc20Handler, erc721Handler, genericHandler interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetContracts", reflect.TypeOf((*MockListener)(nil).SetContracts), bridge, erc20Handler, erc721Handler, genericHandler)
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
func (m *MockWriter) SetBridge(bridge writer.Bridger) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetBridge", bridge)
}

// SetBridge indicates an expected call of SetBridge
func (mr *MockWriterMockRecorder) SetBridge(bridge interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBridge", reflect.TypeOf((*MockWriter)(nil).SetBridge), bridge)
}
