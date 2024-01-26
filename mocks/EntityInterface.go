// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	model "github.com/enbility/spine-go/model"
	mock "github.com/stretchr/testify/mock"
)

// EntityInterface is an autogenerated mock type for the EntityInterface type
type EntityInterface struct {
	mock.Mock
}

type EntityInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *EntityInterface) EXPECT() *EntityInterface_Expecter {
	return &EntityInterface_Expecter{mock: &_m.Mock}
}

// Address provides a mock function with given fields:
func (_m *EntityInterface) Address() *model.EntityAddressType {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Address")
	}

	var r0 *model.EntityAddressType
	if rf, ok := ret.Get(0).(func() *model.EntityAddressType); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.EntityAddressType)
		}
	}

	return r0
}

// EntityInterface_Address_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Address'
type EntityInterface_Address_Call struct {
	*mock.Call
}

// Address is a helper method to define mock.On call
func (_e *EntityInterface_Expecter) Address() *EntityInterface_Address_Call {
	return &EntityInterface_Address_Call{Call: _e.mock.On("Address")}
}

func (_c *EntityInterface_Address_Call) Run(run func()) *EntityInterface_Address_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *EntityInterface_Address_Call) Return(_a0 *model.EntityAddressType) *EntityInterface_Address_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EntityInterface_Address_Call) RunAndReturn(run func() *model.EntityAddressType) *EntityInterface_Address_Call {
	_c.Call.Return(run)
	return _c
}

// Description provides a mock function with given fields:
func (_m *EntityInterface) Description() *model.DescriptionType {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Description")
	}

	var r0 *model.DescriptionType
	if rf, ok := ret.Get(0).(func() *model.DescriptionType); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DescriptionType)
		}
	}

	return r0
}

// EntityInterface_Description_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Description'
type EntityInterface_Description_Call struct {
	*mock.Call
}

// Description is a helper method to define mock.On call
func (_e *EntityInterface_Expecter) Description() *EntityInterface_Description_Call {
	return &EntityInterface_Description_Call{Call: _e.mock.On("Description")}
}

func (_c *EntityInterface_Description_Call) Run(run func()) *EntityInterface_Description_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *EntityInterface_Description_Call) Return(_a0 *model.DescriptionType) *EntityInterface_Description_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EntityInterface_Description_Call) RunAndReturn(run func() *model.DescriptionType) *EntityInterface_Description_Call {
	_c.Call.Return(run)
	return _c
}

// EntityType provides a mock function with given fields:
func (_m *EntityInterface) EntityType() model.EntityTypeType {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for EntityType")
	}

	var r0 model.EntityTypeType
	if rf, ok := ret.Get(0).(func() model.EntityTypeType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.EntityTypeType)
	}

	return r0
}

// EntityInterface_EntityType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EntityType'
type EntityInterface_EntityType_Call struct {
	*mock.Call
}

// EntityType is a helper method to define mock.On call
func (_e *EntityInterface_Expecter) EntityType() *EntityInterface_EntityType_Call {
	return &EntityInterface_EntityType_Call{Call: _e.mock.On("EntityType")}
}

func (_c *EntityInterface_EntityType_Call) Run(run func()) *EntityInterface_EntityType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *EntityInterface_EntityType_Call) Return(_a0 model.EntityTypeType) *EntityInterface_EntityType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EntityInterface_EntityType_Call) RunAndReturn(run func() model.EntityTypeType) *EntityInterface_EntityType_Call {
	_c.Call.Return(run)
	return _c
}

// NextFeatureId provides a mock function with given fields:
func (_m *EntityInterface) NextFeatureId() uint {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for NextFeatureId")
	}

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

// EntityInterface_NextFeatureId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NextFeatureId'
type EntityInterface_NextFeatureId_Call struct {
	*mock.Call
}

// NextFeatureId is a helper method to define mock.On call
func (_e *EntityInterface_Expecter) NextFeatureId() *EntityInterface_NextFeatureId_Call {
	return &EntityInterface_NextFeatureId_Call{Call: _e.mock.On("NextFeatureId")}
}

func (_c *EntityInterface_NextFeatureId_Call) Run(run func()) *EntityInterface_NextFeatureId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *EntityInterface_NextFeatureId_Call) Return(_a0 uint) *EntityInterface_NextFeatureId_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EntityInterface_NextFeatureId_Call) RunAndReturn(run func() uint) *EntityInterface_NextFeatureId_Call {
	_c.Call.Return(run)
	return _c
}

// SetDescription provides a mock function with given fields: d
func (_m *EntityInterface) SetDescription(d *model.DescriptionType) {
	_m.Called(d)
}

// EntityInterface_SetDescription_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetDescription'
type EntityInterface_SetDescription_Call struct {
	*mock.Call
}

// SetDescription is a helper method to define mock.On call
//   - d *model.DescriptionType
func (_e *EntityInterface_Expecter) SetDescription(d interface{}) *EntityInterface_SetDescription_Call {
	return &EntityInterface_SetDescription_Call{Call: _e.mock.On("SetDescription", d)}
}

func (_c *EntityInterface_SetDescription_Call) Run(run func(d *model.DescriptionType)) *EntityInterface_SetDescription_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*model.DescriptionType))
	})
	return _c
}

func (_c *EntityInterface_SetDescription_Call) Return() *EntityInterface_SetDescription_Call {
	_c.Call.Return()
	return _c
}

func (_c *EntityInterface_SetDescription_Call) RunAndReturn(run func(*model.DescriptionType)) *EntityInterface_SetDescription_Call {
	_c.Call.Return(run)
	return _c
}

// NewEntityInterface creates a new instance of EntityInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEntityInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *EntityInterface {
	mock := &EntityInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}