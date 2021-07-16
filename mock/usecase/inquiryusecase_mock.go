// Code generated by MockGen. DO NOT EDIT.
// Source: usecase/inquiryusecase.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/pevin/pevin-golang-training-beginner/model"
)

// MockIInquiryUseCase is a mock of IInquiryUseCase interface.
type MockIInquiryUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockIInquiryUseCaseMockRecorder
}

// MockIInquiryUseCaseMockRecorder is the mock recorder for MockIInquiryUseCase.
type MockIInquiryUseCaseMockRecorder struct {
	mock *MockIInquiryUseCase
}

// NewMockIInquiryUseCase creates a new mock instance.
func NewMockIInquiryUseCase(ctrl *gomock.Controller) *MockIInquiryUseCase {
	mock := &MockIInquiryUseCase{ctrl: ctrl}
	mock.recorder = &MockIInquiryUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIInquiryUseCase) EXPECT() *MockIInquiryUseCaseMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIInquiryUseCase) Create(ctx context.Context, p *model.Inquiry) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, p)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockIInquiryUseCaseMockRecorder) Create(ctx, p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIInquiryUseCase)(nil).Create), ctx, p)
}

// GetByTransactionId mocks base method.
func (m *MockIInquiryUseCase) GetByTransactionId(ctx context.Context, id string) (model.Inquiry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByTransactionId", ctx, id)
	ret0, _ := ret[0].(model.Inquiry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByTransactionId indicates an expected call of GetByTransactionId.
func (mr *MockIInquiryUseCaseMockRecorder) GetByTransactionId(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTransactionId", reflect.TypeOf((*MockIInquiryUseCase)(nil).GetByTransactionId), ctx, id)
}

// InitFromRequest mocks base method.
func (m *MockIInquiryUseCase) InitFromRequest(r *http.Request) (model.Inquiry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitFromRequest", r)
	ret0, _ := ret[0].(model.Inquiry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InitFromRequest indicates an expected call of InitFromRequest.
func (mr *MockIInquiryUseCaseMockRecorder) InitFromRequest(r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitFromRequest", reflect.TypeOf((*MockIInquiryUseCase)(nil).InitFromRequest), r)
}