// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	paginations "github.com/bimalabs/framework/v4/paginations"
	mock "github.com/stretchr/testify/mock"

	repositories "github.com/bimalabs/framework/v4/repositories"
)

// Handler is an autogenerated mock type for the Handler type
type Handler struct {
	mock.Mock
}

// All provides a mock function with given fields: v
func (_m *Handler) All(v interface{}) error {
	ret := _m.Called(v)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Bind provides a mock function with given fields: v, id
func (_m *Handler) Bind(v interface{}, id string) error {
	ret := _m.Called(v, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, string) error); ok {
		r0 = rf(v, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Create provides a mock function with given fields: v
func (_m *Handler) Create(v interface{}) error {
	ret := _m.Called(v)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: v, id
func (_m *Handler) Delete(v interface{}, id string) error {
	ret := _m.Called(v, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, string) error); ok {
		r0 = rf(v, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindBy provides a mock function with given fields: v, filters
func (_m *Handler) FindBy(v interface{}, filters ...repositories.Filter) error {
	_va := make([]interface{}, len(filters))
	for _i := range filters {
		_va[_i] = filters[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, v)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, ...repositories.Filter) error); ok {
		r0 = rf(v, filters...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Paginate provides a mock function with given fields: paginator, result
func (_m *Handler) Paginate(paginator *paginations.Pagination, result interface{}) paginations.Metadata {
	ret := _m.Called(paginator, result)

	var r0 paginations.Metadata
	if rf, ok := ret.Get(0).(func(*paginations.Pagination, interface{}) paginations.Metadata); ok {
		r0 = rf(paginator, result)
	} else {
		r0 = ret.Get(0).(paginations.Metadata)
	}

	return r0
}

// Repository provides a mock function with given fields:
func (_m *Handler) Repository() repositories.Repository {
	ret := _m.Called()

	var r0 repositories.Repository
	if rf, ok := ret.Get(0).(func() repositories.Repository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repositories.Repository)
		}
	}

	return r0
}

// Update provides a mock function with given fields: v, id
func (_m *Handler) Update(v interface{}, id string) error {
	ret := _m.Called(v, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, string) error); ok {
		r0 = rf(v, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewHandler creates a new instance of Handler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewHandler(t mockConstructorTestingTNewHandler) *Handler {
	mock := &Handler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
