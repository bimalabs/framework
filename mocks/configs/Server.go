// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	elastic "github.com/olivere/elastic/v7"
	gorm "gorm.io/gorm"

	grpc "google.golang.org/grpc"

	messengers "github.com/bimalabs/framework/v4/messengers"

	mock "github.com/stretchr/testify/mock"

	runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// Server is an autogenerated mock type for the Server type
type Server struct {
	mock.Mock
}

// Consume provides a mock function with given fields: messenger
func (_m *Server) Consume(messenger *messengers.Messenger) {
	_m.Called(messenger)
}

// Handle provides a mock function with given fields: _a0, server, client
func (_m *Server) Handle(_a0 context.Context, server *runtime.ServeMux, client *grpc.ClientConn) error {
	ret := _m.Called(_a0, server, client)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error); ok {
		r0 = rf(_a0, server, client)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Migrate provides a mock function with given fields: db
func (_m *Server) Migrate(db *gorm.DB) {
	_m.Called(db)
}

// Register provides a mock function with given fields: server
func (_m *Server) Register(server *grpc.Server) {
	_m.Called(server)
}

// Sync provides a mock function with given fields: client
func (_m *Server) Sync(client *elastic.Client) {
	_m.Called(client)
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
