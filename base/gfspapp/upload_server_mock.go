// Code generated by MockGen. DO NOT EDIT.
// Source: ./upload_server.go

// Package gfspapp is a generated GoMock package.
package gfspapp

import (
	context "context"
	reflect "reflect"

	gfspserver "github.com/bnb-chain/greenfield-storage-provider/base/types/gfspserver"
	gomock "go.uber.org/mock/gomock"
	metadata "google.golang.org/grpc/metadata"
)

// MockgRPCUploadStream is a mock of gRPCUploadStream interface.
type MockgRPCUploadStream struct {
	ctrl     *gomock.Controller
	recorder *MockgRPCUploadStreamMockRecorder
}

// MockgRPCUploadStreamMockRecorder is the mock recorder for MockgRPCUploadStream.
type MockgRPCUploadStreamMockRecorder struct {
	mock *MockgRPCUploadStream
}

// NewMockgRPCUploadStream creates a new mock instance.
func NewMockgRPCUploadStream(ctrl *gomock.Controller) *MockgRPCUploadStream {
	mock := &MockgRPCUploadStream{ctrl: ctrl}
	mock.recorder = &MockgRPCUploadStreamMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockgRPCUploadStream) EXPECT() *MockgRPCUploadStreamMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockgRPCUploadStream) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockgRPCUploadStreamMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockgRPCUploadStream)(nil).Context))
}

// Recv mocks base method.
func (m *MockgRPCUploadStream) Recv() (*gfspserver.GfSpUploadObjectRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*gfspserver.GfSpUploadObjectRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockgRPCUploadStreamMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockgRPCUploadStream)(nil).Recv))
}

// RecvMsg mocks base method.
func (m_2 *MockgRPCUploadStream) RecvMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RecvMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockgRPCUploadStreamMockRecorder) RecvMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockgRPCUploadStream)(nil).RecvMsg), m)
}

// SendAndClose mocks base method.
func (m *MockgRPCUploadStream) SendAndClose(arg0 *gfspserver.GfSpUploadObjectResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendAndClose", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendAndClose indicates an expected call of SendAndClose.
func (mr *MockgRPCUploadStreamMockRecorder) SendAndClose(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendAndClose", reflect.TypeOf((*MockgRPCUploadStream)(nil).SendAndClose), arg0)
}

// SendHeader mocks base method.
func (m *MockgRPCUploadStream) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockgRPCUploadStreamMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockgRPCUploadStream)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m_2 *MockgRPCUploadStream) SendMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SendMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockgRPCUploadStreamMockRecorder) SendMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockgRPCUploadStream)(nil).SendMsg), m)
}

// SetHeader mocks base method.
func (m *MockgRPCUploadStream) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockgRPCUploadStreamMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockgRPCUploadStream)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockgRPCUploadStream) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockgRPCUploadStreamMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockgRPCUploadStream)(nil).SetTrailer), arg0)
}

// MockgRPCResumableUploadStream is a mock of gRPCResumableUploadStream interface.
type MockgRPCResumableUploadStream struct {
	ctrl     *gomock.Controller
	recorder *MockgRPCResumableUploadStreamMockRecorder
}

// MockgRPCResumableUploadStreamMockRecorder is the mock recorder for MockgRPCResumableUploadStream.
type MockgRPCResumableUploadStreamMockRecorder struct {
	mock *MockgRPCResumableUploadStream
}

// NewMockgRPCResumableUploadStream creates a new mock instance.
func NewMockgRPCResumableUploadStream(ctrl *gomock.Controller) *MockgRPCResumableUploadStream {
	mock := &MockgRPCResumableUploadStream{ctrl: ctrl}
	mock.recorder = &MockgRPCResumableUploadStreamMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockgRPCResumableUploadStream) EXPECT() *MockgRPCResumableUploadStreamMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockgRPCResumableUploadStream) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockgRPCResumableUploadStreamMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockgRPCResumableUploadStream)(nil).Context))
}

// Recv mocks base method.
func (m *MockgRPCResumableUploadStream) Recv() (*gfspserver.GfSpResumableUploadObjectRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*gfspserver.GfSpResumableUploadObjectRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockgRPCResumableUploadStreamMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockgRPCResumableUploadStream)(nil).Recv))
}

// RecvMsg mocks base method.
func (m_2 *MockgRPCResumableUploadStream) RecvMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RecvMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockgRPCResumableUploadStreamMockRecorder) RecvMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockgRPCResumableUploadStream)(nil).RecvMsg), m)
}

// SendAndClose mocks base method.
func (m *MockgRPCResumableUploadStream) SendAndClose(arg0 *gfspserver.GfSpResumableUploadObjectResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendAndClose", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendAndClose indicates an expected call of SendAndClose.
func (mr *MockgRPCResumableUploadStreamMockRecorder) SendAndClose(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendAndClose", reflect.TypeOf((*MockgRPCResumableUploadStream)(nil).SendAndClose), arg0)
}

// SendHeader mocks base method.
func (m *MockgRPCResumableUploadStream) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockgRPCResumableUploadStreamMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockgRPCResumableUploadStream)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m_2 *MockgRPCResumableUploadStream) SendMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SendMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockgRPCResumableUploadStreamMockRecorder) SendMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockgRPCResumableUploadStream)(nil).SendMsg), m)
}

// SetHeader mocks base method.
func (m *MockgRPCResumableUploadStream) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockgRPCResumableUploadStreamMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockgRPCResumableUploadStream)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockgRPCResumableUploadStream) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockgRPCResumableUploadStreamMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockgRPCResumableUploadStream)(nil).SetTrailer), arg0)
}