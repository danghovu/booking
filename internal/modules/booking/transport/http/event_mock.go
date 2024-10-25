// Code generated by MockGen. DO NOT EDIT.
// Source: event.go
//
// Generated by this command:
//
//	mockgen -source=event.go -destination=event_mock.go -package=transporthttp
//

// Package transporthttp is a generated GoMock package.
package transporthttp

import (
	model "booking-event/internal/modules/booking/model"
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockEventHandler is a mock of EventHandler interface.
type MockEventHandler struct {
	ctrl     *gomock.Controller
	recorder *MockEventHandlerMockRecorder
}

// MockEventHandlerMockRecorder is the mock recorder for MockEventHandler.
type MockEventHandlerMockRecorder struct {
	mock *MockEventHandler
}

// NewMockEventHandler creates a new mock instance.
func NewMockEventHandler(ctrl *gomock.Controller) *MockEventHandler {
	mock := &MockEventHandler{ctrl: ctrl}
	mock.recorder = &MockEventHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventHandler) EXPECT() *MockEventHandlerMockRecorder {
	return m.recorder
}

// CreateEvent mocks base method.
func (m *MockEventHandler) CreateEvent(ctx context.Context, params model.CreateEventRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEvent", ctx, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateEvent indicates an expected call of CreateEvent.
func (mr *MockEventHandlerMockRecorder) CreateEvent(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEvent", reflect.TypeOf((*MockEventHandler)(nil).CreateEvent), ctx, params)
}

// QueryEvents mocks base method.
func (m *MockEventHandler) QueryEvents(ctx context.Context, query model.EventQuery) ([]model.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryEvents", ctx, query)
	ret0, _ := ret[0].([]model.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryEvents indicates an expected call of QueryEvents.
func (mr *MockEventHandlerMockRecorder) QueryEvents(ctx, query any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryEvents", reflect.TypeOf((*MockEventHandler)(nil).QueryEvents), ctx, query)
}

// RetrieveEventDetail mocks base method.
func (m *MockEventHandler) RetrieveEventDetail(ctx context.Context, eventID int) (*model.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveEventDetail", ctx, eventID)
	ret0, _ := ret[0].(*model.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveEventDetail indicates an expected call of RetrieveEventDetail.
func (mr *MockEventHandlerMockRecorder) RetrieveEventDetail(ctx, eventID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveEventDetail", reflect.TypeOf((*MockEventHandler)(nil).RetrieveEventDetail), ctx, eventID)
}

// UpdateEvent mocks base method.
func (m *MockEventHandler) UpdateEvent(ctx context.Context, params model.UpdateEventRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEvent", ctx, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEvent indicates an expected call of UpdateEvent.
func (mr *MockEventHandlerMockRecorder) UpdateEvent(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEvent", reflect.TypeOf((*MockEventHandler)(nil).UpdateEvent), ctx, params)
}