// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	api "github.com/enbility/spine-go/api"
	mock "github.com/stretchr/testify/mock"
)

// HeartbeatManagerInterface is an autogenerated mock type for the HeartbeatManagerInterface type
type HeartbeatManagerInterface struct {
	mock.Mock
}

type HeartbeatManagerInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *HeartbeatManagerInterface) EXPECT() *HeartbeatManagerInterface_Expecter {
	return &HeartbeatManagerInterface_Expecter{mock: &_m.Mock}
}

// IsHeartbeatRunning provides a mock function with given fields:
func (_m *HeartbeatManagerInterface) IsHeartbeatRunning() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsHeartbeatRunning")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// HeartbeatManagerInterface_IsHeartbeatRunning_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsHeartbeatRunning'
type HeartbeatManagerInterface_IsHeartbeatRunning_Call struct {
	*mock.Call
}

// IsHeartbeatRunning is a helper method to define mock.On call
func (_e *HeartbeatManagerInterface_Expecter) IsHeartbeatRunning() *HeartbeatManagerInterface_IsHeartbeatRunning_Call {
	return &HeartbeatManagerInterface_IsHeartbeatRunning_Call{Call: _e.mock.On("IsHeartbeatRunning")}
}

func (_c *HeartbeatManagerInterface_IsHeartbeatRunning_Call) Run(run func()) *HeartbeatManagerInterface_IsHeartbeatRunning_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *HeartbeatManagerInterface_IsHeartbeatRunning_Call) Return(_a0 bool) *HeartbeatManagerInterface_IsHeartbeatRunning_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *HeartbeatManagerInterface_IsHeartbeatRunning_Call) RunAndReturn(run func() bool) *HeartbeatManagerInterface_IsHeartbeatRunning_Call {
	_c.Call.Return(run)
	return _c
}

// SetLocalFeature provides a mock function with given fields: entity, feature
func (_m *HeartbeatManagerInterface) SetLocalFeature(entity api.EntityLocalInterface, feature api.FeatureLocalInterface) {
	_m.Called(entity, feature)
}

// HeartbeatManagerInterface_SetLocalFeature_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetLocalFeature'
type HeartbeatManagerInterface_SetLocalFeature_Call struct {
	*mock.Call
}

// SetLocalFeature is a helper method to define mock.On call
//   - entity api.EntityLocalInterface
//   - feature api.FeatureLocalInterface
func (_e *HeartbeatManagerInterface_Expecter) SetLocalFeature(entity interface{}, feature interface{}) *HeartbeatManagerInterface_SetLocalFeature_Call {
	return &HeartbeatManagerInterface_SetLocalFeature_Call{Call: _e.mock.On("SetLocalFeature", entity, feature)}
}

func (_c *HeartbeatManagerInterface_SetLocalFeature_Call) Run(run func(entity api.EntityLocalInterface, feature api.FeatureLocalInterface)) *HeartbeatManagerInterface_SetLocalFeature_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(api.EntityLocalInterface), args[1].(api.FeatureLocalInterface))
	})
	return _c
}

func (_c *HeartbeatManagerInterface_SetLocalFeature_Call) Return() *HeartbeatManagerInterface_SetLocalFeature_Call {
	_c.Call.Return()
	return _c
}

func (_c *HeartbeatManagerInterface_SetLocalFeature_Call) RunAndReturn(run func(api.EntityLocalInterface, api.FeatureLocalInterface)) *HeartbeatManagerInterface_SetLocalFeature_Call {
	_c.Call.Return(run)
	return _c
}

// StartHeartbeat provides a mock function with given fields:
func (_m *HeartbeatManagerInterface) StartHeartbeat() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for StartHeartbeat")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HeartbeatManagerInterface_StartHeartbeat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StartHeartbeat'
type HeartbeatManagerInterface_StartHeartbeat_Call struct {
	*mock.Call
}

// StartHeartbeat is a helper method to define mock.On call
func (_e *HeartbeatManagerInterface_Expecter) StartHeartbeat() *HeartbeatManagerInterface_StartHeartbeat_Call {
	return &HeartbeatManagerInterface_StartHeartbeat_Call{Call: _e.mock.On("StartHeartbeat")}
}

func (_c *HeartbeatManagerInterface_StartHeartbeat_Call) Run(run func()) *HeartbeatManagerInterface_StartHeartbeat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *HeartbeatManagerInterface_StartHeartbeat_Call) Return(_a0 error) *HeartbeatManagerInterface_StartHeartbeat_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *HeartbeatManagerInterface_StartHeartbeat_Call) RunAndReturn(run func() error) *HeartbeatManagerInterface_StartHeartbeat_Call {
	_c.Call.Return(run)
	return _c
}

// StopHeartbeat provides a mock function with given fields:
func (_m *HeartbeatManagerInterface) StopHeartbeat() {
	_m.Called()
}

// HeartbeatManagerInterface_StopHeartbeat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StopHeartbeat'
type HeartbeatManagerInterface_StopHeartbeat_Call struct {
	*mock.Call
}

// StopHeartbeat is a helper method to define mock.On call
func (_e *HeartbeatManagerInterface_Expecter) StopHeartbeat() *HeartbeatManagerInterface_StopHeartbeat_Call {
	return &HeartbeatManagerInterface_StopHeartbeat_Call{Call: _e.mock.On("StopHeartbeat")}
}

func (_c *HeartbeatManagerInterface_StopHeartbeat_Call) Run(run func()) *HeartbeatManagerInterface_StopHeartbeat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *HeartbeatManagerInterface_StopHeartbeat_Call) Return() *HeartbeatManagerInterface_StopHeartbeat_Call {
	_c.Call.Return()
	return _c
}

func (_c *HeartbeatManagerInterface_StopHeartbeat_Call) RunAndReturn(run func()) *HeartbeatManagerInterface_StopHeartbeat_Call {
	_c.Call.Return(run)
	return _c
}

// NewHeartbeatManagerInterface creates a new instance of HeartbeatManagerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHeartbeatManagerInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *HeartbeatManagerInterface {
	mock := &HeartbeatManagerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
