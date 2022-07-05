// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	repositories "github.com/bimalabs/framework/v4/repositories"
	mock "github.com/stretchr/testify/mock"
)

// Transaction is an autogenerated mock type for the Transaction type
type Transaction struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *Transaction) Execute(_a0 repositories.Repository) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(repositories.Repository) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTransaction interface {
	mock.TestingT
	Cleanup(func())
}

// NewTransaction creates a new instance of Transaction. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTransaction(t mockConstructorTestingTNewTransaction) *Transaction {
	mock := &Transaction{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
