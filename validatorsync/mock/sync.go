// Code generated by MockGen. DO NOT EDIT.
// Source: ./validatorsync/sync.go

// Package mock_validatorsync is a generated GoMock package.
package mock_validatorsync

import (
	context "context"
	types "github.com/ethereum/go-ethereum/core/types"
	gomock "github.com/golang/mock/gomock"
	big "math/big"
	reflect "reflect"
)

// MockHeaderByNumberGetter is a mock of HeaderByNumberGetter interface
type MockHeaderByNumberGetter struct {
	ctrl     *gomock.Controller
	recorder *MockHeaderByNumberGetterMockRecorder
}

// MockHeaderByNumberGetterMockRecorder is the mock recorder for MockHeaderByNumberGetter
type MockHeaderByNumberGetterMockRecorder struct {
	mock *MockHeaderByNumberGetter
}

// NewMockHeaderByNumberGetter creates a new mock instance
func NewMockHeaderByNumberGetter(ctrl *gomock.Controller) *MockHeaderByNumberGetter {
	mock := &MockHeaderByNumberGetter{ctrl: ctrl}
	mock.recorder = &MockHeaderByNumberGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHeaderByNumberGetter) EXPECT() *MockHeaderByNumberGetterMockRecorder {
	return m.recorder
}

// HeaderByNumber mocks base method
func (m *MockHeaderByNumberGetter) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HeaderByNumber", ctx, number)
	ret0, _ := ret[0].(*types.Header)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HeaderByNumber indicates an expected call of HeaderByNumber
func (mr *MockHeaderByNumberGetterMockRecorder) HeaderByNumber(ctx, number interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HeaderByNumber", reflect.TypeOf((*MockHeaderByNumberGetter)(nil).HeaderByNumber), ctx, number)
}
