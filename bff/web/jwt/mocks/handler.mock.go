// Code generated by MockGen. DO NOT EDIT.
// Source: ./types.go

// Package jwtmocks is a generated GoMock package.
package jwtmocks

import (
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
)

// MockHandler is a mock of Handler interface.
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler.
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance.
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// CheckSession mocks base method.
func (m *MockHandler) CheckSession(ctx *gin.Context, ssid string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckSession", ctx, ssid)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckSession indicates an expected call of CheckSession.
func (mr *MockHandlerMockRecorder) CheckSession(ctx, ssid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckSession", reflect.TypeOf((*MockHandler)(nil).CheckSession), ctx, ssid)
}

// ClearToken mocks base method.
func (m *MockHandler) ClearToken(ctx *gin.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearToken", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// ClearToken indicates an expected call of ClearToken.
func (mr *MockHandlerMockRecorder) ClearToken(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearToken", reflect.TypeOf((*MockHandler)(nil).ClearToken), ctx)
}

// ExtractTokenString mocks base method.
func (m *MockHandler) ExtractTokenString(ctx *gin.Context) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExtractTokenString", ctx)
	ret0, _ := ret[0].(string)
	return ret0
}

// ExtractTokenString indicates an expected call of ExtractTokenString.
func (mr *MockHandlerMockRecorder) ExtractTokenString(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExtractTokenString", reflect.TypeOf((*MockHandler)(nil).ExtractTokenString), ctx)
}

// SetLoginToken mocks base method.
func (m *MockHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLoginToken", ctx, uid)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetLoginToken indicates an expected call of SetLoginToken.
func (mr *MockHandlerMockRecorder) SetLoginToken(ctx, uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLoginToken", reflect.TypeOf((*MockHandler)(nil).SetLoginToken), ctx, uid)
}

// SwtJWTToken mocks base method.
func (m *MockHandler) SwtJWTToken(ctx *gin.Context, ssid string, uid int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SwtJWTToken", ctx, ssid, uid)
	ret0, _ := ret[0].(error)
	return ret0
}

// SwtJWTToken indicates an expected call of SwtJWTToken.
func (mr *MockHandlerMockRecorder) SwtJWTToken(ctx, ssid, uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SwtJWTToken", reflect.TypeOf((*MockHandler)(nil).SwtJWTToken), ctx, ssid, uid)
}