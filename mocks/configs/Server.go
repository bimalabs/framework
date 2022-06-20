// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// Server is an autogenerated mock type for the Server type
type Server struct {
	mock.Mock
}

// GRpcHandler provides a mock function with given fields: _a0, server, client
func (_m *Server) GRpcHandler(_a0 context.Context, server *runtime.ServeMux, client *grpc.ClientConn) error {
	ret := _m.Called(_a0, server, client)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error); ok {
		r0 = rf(_a0, server, client)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RegisterAutoMigrate provides a mock function with given fields:
func (_m *Server) RegisterAutoMigrate() {
	_m.Called()
}

// RegisterGRpc provides a mock function with given fields: server
func (_m *Server) RegisterGRpc(server *grpc.Server) {
	_m.Called(server)
}

// RegisterQueueConsumer provides a mock function with given fields:
func (_m *Server) RegisterQueueConsumer() {
	_m.Called()
}

// RepopulateData provides a mock function with given fields:
func (_m *Server) RepopulateData() {
	_m.Called()
}

type mockConstructorTestingTNewServer interface {
	mock.TestingT
	Cleanup(func())
}

// NewServer creates a new instance of Server. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewServer(t mockConstructorTestingTNewServer) *Server {
	mock := &Server{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
