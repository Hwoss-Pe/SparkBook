// Code generated by MockGen. DO NOT EDIT.
// Source: types.go

// Package repomocks is a generated GoMock package.
package repomocks

import (
	domain "Webook/payment/domain"
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockPaymentRepository is a mock of PaymentRepository interface.
type MockPaymentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPaymentRepositoryMockRecorder
}

// MockPaymentRepositoryMockRecorder is the mock recorder for MockPaymentRepository.
type MockPaymentRepositoryMockRecorder struct {
	mock *MockPaymentRepository
}

// NewMockPaymentRepository creates a new mock instance.
func NewMockPaymentRepository(ctrl *gomock.Controller) *MockPaymentRepository {
	mock := &MockPaymentRepository{ctrl: ctrl}
	mock.recorder = &MockPaymentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPaymentRepository) EXPECT() *MockPaymentRepositoryMockRecorder {
	return m.recorder
}

// AddPayment mocks base method.
func (m *MockPaymentRepository) AddPayment(ctx context.Context, pmt domain.Payment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPayment", ctx, pmt)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddPayment indicates an expected call of AddPayment.
func (mr *MockPaymentRepositoryMockRecorder) AddPayment(ctx, pmt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPayment", reflect.TypeOf((*MockPaymentRepository)(nil).AddPayment), ctx, pmt)
}

// FindExpiredPayment mocks base method.
func (m *MockPaymentRepository) FindExpiredPayment(ctx context.Context, offset, limit int, t time.Time) ([]domain.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindExpiredPayment", ctx, offset, limit, t)
	ret0, _ := ret[0].([]domain.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindExpiredPayment indicates an expected call of FindExpiredPayment.
func (mr *MockPaymentRepositoryMockRecorder) FindExpiredPayment(ctx, offset, limit, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindExpiredPayment", reflect.TypeOf((*MockPaymentRepository)(nil).FindExpiredPayment), ctx, offset, limit, t)
}

// GetPayment mocks base method.
func (m *MockPaymentRepository) GetPayment(ctx context.Context, bizTradeNO string) (domain.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPayment", ctx, bizTradeNO)
	ret0, _ := ret[0].(domain.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPayment indicates an expected call of GetPayment.
func (mr *MockPaymentRepositoryMockRecorder) GetPayment(ctx, bizTradeNO interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPayment", reflect.TypeOf((*MockPaymentRepository)(nil).GetPayment), ctx, bizTradeNO)
}

// UpdatePayment mocks base method.
func (m *MockPaymentRepository) UpdatePayment(ctx context.Context, pmt domain.Payment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePayment", ctx, pmt)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePayment indicates an expected call of UpdatePayment.
func (mr *MockPaymentRepositoryMockRecorder) UpdatePayment(ctx, pmt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePayment", reflect.TypeOf((*MockPaymentRepository)(nil).UpdatePayment), ctx, pmt)
}