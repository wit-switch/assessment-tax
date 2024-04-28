// Code generated by MockGen. DO NOT EDIT.
// Source: internal/core/port/service.go
//
// Generated by this command:
//
//	mockgen -source=internal/core/port/service.go -destination=mocks/service.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	csv "encoding/csv"
	reflect "reflect"

	domain "github.com/wit-switch/assessment-tax/internal/core/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockTaxService is a mock of TaxService interface.
type MockTaxService struct {
	ctrl     *gomock.Controller
	recorder *MockTaxServiceMockRecorder
}

// MockTaxServiceMockRecorder is the mock recorder for MockTaxService.
type MockTaxServiceMockRecorder struct {
	mock *MockTaxService
}

// NewMockTaxService creates a new mock instance.
func NewMockTaxService(ctrl *gomock.Controller) *MockTaxService {
	mock := &MockTaxService{ctrl: ctrl}
	mock.recorder = &MockTaxServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaxService) EXPECT() *MockTaxServiceMockRecorder {
	return m.recorder
}

// Calculate mocks base method.
func (m *MockTaxService) Calculate(ctx context.Context, body domain.TaxCalculate, allowKReceipt bool) (*domain.Tax, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Calculate", ctx, body, allowKReceipt)
	ret0, _ := ret[0].(*domain.Tax)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Calculate indicates an expected call of Calculate.
func (mr *MockTaxServiceMockRecorder) Calculate(ctx, body, allowKReceipt any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Calculate", reflect.TypeOf((*MockTaxService)(nil).Calculate), ctx, body, allowKReceipt)
}

// CalculateFromCSV mocks base method.
func (m *MockTaxService) CalculateFromCSV(ctx context.Context, file csv.Reader) ([]domain.TaxCSV, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CalculateFromCSV", ctx, file)
	ret0, _ := ret[0].([]domain.TaxCSV)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CalculateFromCSV indicates an expected call of CalculateFromCSV.
func (mr *MockTaxServiceMockRecorder) CalculateFromCSV(ctx, file any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CalculateFromCSV", reflect.TypeOf((*MockTaxService)(nil).CalculateFromCSV), ctx, file)
}

// UpdateTaxDeduct mocks base method.
func (m *MockTaxService) UpdateTaxDeduct(ctx context.Context, body domain.UpdateTaxDeduct) (*domain.TaxDeduct, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTaxDeduct", ctx, body)
	ret0, _ := ret[0].(*domain.TaxDeduct)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTaxDeduct indicates an expected call of UpdateTaxDeduct.
func (mr *MockTaxServiceMockRecorder) UpdateTaxDeduct(ctx, body any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTaxDeduct", reflect.TypeOf((*MockTaxService)(nil).UpdateTaxDeduct), ctx, body)
}
