// Code generated by MockGen. DO NOT EDIT.
// Source: internal/core/port/repository.go
//
// Generated by this command:
//
//	mockgen -source=internal/core/port/repository.go -destination=mocks/repository.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	domain "github.com/wit-switch/assessment-tax/internal/core/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockTaxRepository is a mock of TaxRepository interface.
type MockTaxRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTaxRepositoryMockRecorder
}

// MockTaxRepositoryMockRecorder is the mock recorder for MockTaxRepository.
type MockTaxRepositoryMockRecorder struct {
	mock *MockTaxRepository
}

// NewMockTaxRepository creates a new mock instance.
func NewMockTaxRepository(ctrl *gomock.Controller) *MockTaxRepository {
	mock := &MockTaxRepository{ctrl: ctrl}
	mock.recorder = &MockTaxRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaxRepository) EXPECT() *MockTaxRepositoryMockRecorder {
	return m.recorder
}

// ListTaxDeduct mocks base method.
func (m *MockTaxRepository) ListTaxDeduct(ctx context.Context, q domain.GetTaxDeduct) ([]domain.TaxDeduct, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTaxDeduct", ctx, q)
	ret0, _ := ret[0].([]domain.TaxDeduct)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTaxDeduct indicates an expected call of ListTaxDeduct.
func (mr *MockTaxRepositoryMockRecorder) ListTaxDeduct(ctx, q any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTaxDeduct", reflect.TypeOf((*MockTaxRepository)(nil).ListTaxDeduct), ctx, q)
}