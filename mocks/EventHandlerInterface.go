// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	api "github.com/enbility/spine-go/api"
	mock "github.com/stretchr/testify/mock"
)

// EventHandlerInterface is an autogenerated mock type for the EventHandlerInterface type
type EventHandlerInterface struct {
	mock.Mock
}

type EventHandlerInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *EventHandlerInterface) EXPECT() *EventHandlerInterface_Expecter {
	return &EventHandlerInterface_Expecter{mock: &_m.Mock}
}

// HandleEvent provides a mock function with given fields: _a0
func (_m *EventHandlerInterface) HandleEvent(_a0 api.EventPayload) {
	_m.Called(_a0)
}

// EventHandlerInterface_HandleEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HandleEvent'
type EventHandlerInterface_HandleEvent_Call struct {
	*mock.Call
}

// HandleEvent is a helper method to define mock.On call
//   - _a0 api.EventPayload
func (_e *EventHandlerInterface_Expecter) HandleEvent(_a0 interface{}) *EventHandlerInterface_HandleEvent_Call {
	return &EventHandlerInterface_HandleEvent_Call{Call: _e.mock.On("HandleEvent", _a0)}
}

func (_c *EventHandlerInterface_HandleEvent_Call) Run(run func(_a0 api.EventPayload)) *EventHandlerInterface_HandleEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(api.EventPayload))
	})
	return _c
}

func (_c *EventHandlerInterface_HandleEvent_Call) Return() *EventHandlerInterface_HandleEvent_Call {
	_c.Call.Return()
	return _c
}

func (_c *EventHandlerInterface_HandleEvent_Call) RunAndReturn(run func(api.EventPayload)) *EventHandlerInterface_HandleEvent_Call {
	_c.Call.Return(run)
	return _c
}

// NewEventHandlerInterface creates a new instance of EventHandlerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEventHandlerInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *EventHandlerInterface {
	mock := &EventHandlerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
