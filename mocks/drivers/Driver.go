// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	gorm "gorm.io/gorm"
)

// Driver is an autogenerated mock type for the Driver type
type Driver struct {
	mock.Mock
}

// Connect provides a mock function with given fields: host, port, user, password, dbname, debug
func (_m *Driver) Connect(host string, port int, user string, password string, dbname string, debug bool) *gorm.DB {
	ret := _m.Called(host, port, user, password, dbname, debug)

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func(string, int, string, string, string, bool) *gorm.DB); ok {
		r0 = rf(host, port, user, password, dbname, debug)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

// Name provides a mock function with given fields:
func (_m *Driver) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

type mockConstructorTestingTNewDriver interface {
	mock.TestingT
	Cleanup(func())
}

// NewDriver creates a new instance of Driver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDriver(t mockConstructorTestingTNewDriver) *Driver {
	mock := &Driver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
