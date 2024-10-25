// Code generated by MockGen. DO NOT EDIT.
// Source: tokenservice.go
//
// Generated by this command:
//
//	mockgen -source=tokenservice.go -destination=tokenservice_mock.go -package=services
//

// Package services is a generated GoMock package.
package services

import (
	model "booking-event/internal/modules/booking/model"
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockTokenRepository is a mock of TokenRepository interface.
type MockTokenRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTokenRepositoryMockRecorder
}

// MockTokenRepositoryMockRecorder is the mock recorder for MockTokenRepository.
type MockTokenRepositoryMockRecorder struct {
	mock *MockTokenRepository
}

// NewMockTokenRepository creates a new mock instance.
func NewMockTokenRepository(ctrl *gomock.Controller) *MockTokenRepository {
	mock := &MockTokenRepository{ctrl: ctrl}
	mock.recorder = &MockTokenRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenRepository) EXPECT() *MockTokenRepositoryMockRecorder {
	return m.recorder
}

// CreateTokens mocks base method.
func (m *MockTokenRepository) CreateTokens(ctx context.Context, tokens []model.EventToken) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTokens", ctx, tokens)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTokens indicates an expected call of CreateTokens.
func (mr *MockTokenRepositoryMockRecorder) CreateTokens(ctx, tokens any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTokens", reflect.TypeOf((*MockTokenRepository)(nil).CreateTokens), ctx, tokens)
}

// GetByToken mocks base method.
func (m *MockTokenRepository) GetByToken(ctx context.Context, token string) (*model.EventToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByToken", ctx, token)
	ret0, _ := ret[0].(*model.EventToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByToken indicates an expected call of GetByToken.
func (mr *MockTokenRepositoryMockRecorder) GetByToken(ctx, token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByToken", reflect.TypeOf((*MockTokenRepository)(nil).GetByToken), ctx, token)
}

// ReleaseToken mocks base method.
func (m *MockTokenRepository) ReleaseToken(ctx context.Context, eventToken *model.EventToken) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReleaseToken", ctx, eventToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReleaseToken indicates an expected call of ReleaseToken.
func (mr *MockTokenRepositoryMockRecorder) ReleaseToken(ctx, eventToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReleaseToken", reflect.TypeOf((*MockTokenRepository)(nil).ReleaseToken), ctx, eventToken)
}

// SelectAvailableToken mocks base method.
func (m *MockTokenRepository) SelectAvailableToken(ctx context.Context, holderID int32, eventID, quantity int) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectAvailableToken", ctx, holderID, eventID, quantity)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectAvailableToken indicates an expected call of SelectAvailableToken.
func (mr *MockTokenRepositoryMockRecorder) SelectAvailableToken(ctx, holderID, eventID, quantity any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectAvailableToken", reflect.TypeOf((*MockTokenRepository)(nil).SelectAvailableToken), ctx, holderID, eventID, quantity)
}
