// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	message "github.com/ThreeDotsLabs/watermill/message"

	mock "github.com/stretchr/testify/mock"
)

// Broker is an autogenerated mock type for the Broker type
type Broker struct {
	mock.Mock
}

// Consume provides a mock function with given fields: queueName
func (_m *Broker) Consume(queueName string) (<-chan *message.Message, error) {
	ret := _m.Called(queueName)

	var r0 <-chan *message.Message
	if rf, ok := ret.Get(0).(func(string) <-chan *message.Message); ok {
		r0 = rf(queueName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan *message.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(queueName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Publish provides a mock function with given fields: queueName, payload
func (_m *Broker) Publish(queueName string, payload message.Payload) error {
	ret := _m.Called(queueName, payload)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, message.Payload) error); ok {
		r0 = rf(queueName, payload)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewBroker interface {
	mock.TestingT
	Cleanup(func())
}

// NewBroker creates a new instance of Broker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBroker(t mockConstructorTestingTNewBroker) *Broker {
	mock := &Broker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
