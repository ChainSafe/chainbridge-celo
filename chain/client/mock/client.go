// Code generated by MockGen. DO NOT EDIT.
// Source: ./chain/client/client.go

// Package mock_client is a generated GoMock package.
package mock_client

import (
	context "context"
	ethereum "github.com/celo-org/celo-blockchain"
	types "github.com/celo-org/celo-blockchain/core/types"
	gomock "github.com/golang/mock/gomock"
	big "math/big"
	reflect "reflect"
)

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

// BlockByNumber mocks base method
func (m *MockLogFilterWithLatestBlock) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockByNumber", ctx, number)
	ret0, _ := ret[0].(*types.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BlockByNumber indicates an expected call of BlockByNumber
func (mr *MockLogFilterWithLatestBlockMockRecorder) BlockByNumber(ctx, number interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockByNumber", reflect.TypeOf((*MockLogFilterWithLatestBlock)(nil).BlockByNumber), ctx, number)
}
