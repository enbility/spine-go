package spine

import (
	"testing"
	"time"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestBindingManagerSuite(t *testing.T) {
	suite.Run(t, new(BindingManagerSuite))
}

type BindingManagerSuite struct {
	suite.Suite

	localDevice  api.DeviceLocalInterface
	writeHandler *WriteMessageHandler
	remoteDevice *DeviceRemote

	sut api.BindingManagerInterface
}

func (s *BindingManagerSuite) BeforeTest(suiteName, testName string) {
	s.localDevice = NewDeviceLocal("TestBrandName", "TestDeviceModel", "TestSerialNumber", "TestDeviceCode",
		"TestDeviceAddress", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)
	remoteSki := "TestRemoteSki"

	s.writeHandler = &WriteMessageHandler{}

	sender := NewSender(s.writeHandler)
	s.remoteDevice = NewDeviceRemote(s.localDevice, remoteSki, sender)
	s.remoteDevice.address = util.Ptr(model.AddressDeviceType("Address"))

	s.sut = NewBindingManager(s.localDevice)
}

func (suite *BindingManagerSuite) Test_Bindings() {
	entity := NewEntityLocal(suite.localDevice, model.EntityTypeTypeCEM, []model.AddressEntityType{1}, time.Second*4)
	suite.localDevice.AddEntity(entity)

	localServerFeature := entity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	localServerFeature2 := entity.GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	localServerFeature3 := entity.GetOrAddFeature(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	localClientFeature := entity.GetOrAddFeature(model.FeatureTypeTypeGeneric, model.RoleTypeClient)

	remoteDeviceAddress := model.AddressDeviceType("remoteDevice")
	suite.remoteDevice.UpdateDevice(
		&model.NetworkManagementDeviceDescriptionDataType{
			DeviceAddress: &model.DeviceAddressType{Device: &remoteDeviceAddress},
		},
	)

	remoteEntity := NewEntityRemote(suite.remoteDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{1})

	remoteClientFeature := NewFeatureRemote(remoteEntity.NextFeatureId(), remoteEntity, model.FeatureTypeTypeGeneric, model.RoleTypeClient)
	remoteClientFeature.Address().Device = util.Ptr(remoteDeviceAddress)
	remoteEntity.AddFeature(remoteClientFeature)

	remoteClientFeature2 := NewFeatureRemote(remoteEntity.NextFeatureId(), remoteEntity, model.FeatureTypeTypeGeneric, model.RoleTypeClient)
	remoteClientFeature2.Address().Device = util.Ptr(remoteDeviceAddress)
	remoteEntity.AddFeature(remoteClientFeature2)

	remoteServerFeature := NewFeatureRemote(remoteEntity.NextFeatureId(), remoteEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	remoteServerFeature.Address().Device = util.Ptr(remoteDeviceAddress)
	remoteEntity.AddFeature(remoteServerFeature)

	suite.remoteDevice.AddEntity(remoteEntity)

	bindingMgr := suite.localDevice.BindingManager()

	bindingRequest := model.BindingManagementRequestCallType{
		ClientAddress: util.Ptr(model.FeatureAddressType{
			Device:  util.Ptr(model.AddressDeviceType("dummy")),
			Entity:  []model.AddressEntityType{1000},
			Feature: util.Ptr(model.AddressFeatureType(1000)),
		}),
		ServerAddress: util.Ptr(model.FeatureAddressType{
			Device:  util.Ptr(model.AddressDeviceType("dummy")),
			Entity:  []model.AddressEntityType{1000},
			Feature: util.Ptr(model.AddressFeatureType(1000)),
		}),
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeDeviceDiagnosis),
	}

	err := bindingMgr.AddBinding(suite.remoteDevice, bindingRequest)
	assert.NotNil(suite.T(), err)

	bindingRequest = model.BindingManagementRequestCallType{
		ClientAddress: util.Ptr(model.FeatureAddressType{
			Device:  util.Ptr(model.AddressDeviceType("dummy")),
			Entity:  []model.AddressEntityType{1000},
			Feature: util.Ptr(model.AddressFeatureType(1000)),
		}),
		ServerAddress: localClientFeature.Address(),
	}

	err = bindingMgr.AddBinding(suite.remoteDevice, bindingRequest)
	assert.NotNil(suite.T(), err)

	bindingRequest.ServerFeatureType = util.Ptr(model.FeatureTypeTypeDeviceDiagnosis)

	err = bindingMgr.AddBinding(suite.remoteDevice, bindingRequest)
	assert.NotNil(suite.T(), err)

	bindingRequest.ServerAddress = localServerFeature.Address()

	err = bindingMgr.AddBinding(suite.remoteDevice, bindingRequest)
	assert.NotNil(suite.T(), err)

	bindingRequest.ClientAddress = remoteServerFeature.Address()

	err = bindingMgr.AddBinding(suite.remoteDevice, bindingRequest)
	assert.NotNil(suite.T(), err)

	bindingRequest.ClientAddress = remoteClientFeature.Address()

	err = bindingMgr.AddBinding(suite.remoteDevice, bindingRequest)
	assert.Nil(suite.T(), err)

	subs := bindingMgr.Bindings(suite.remoteDevice)
	assert.Equal(suite.T(), 1, len(subs))

	err = bindingMgr.AddBinding(suite.remoteDevice, bindingRequest)
	assert.NotNil(suite.T(), err)

	subs = bindingMgr.Bindings(suite.remoteDevice)
	assert.Equal(suite.T(), 1, len(subs))

	address := model.FeatureAddressType{
		Device:  entity.Device().Address(),
		Entity:  entity.Address().Entity,
		Feature: util.Ptr(model.AddressFeatureType(10)),
	}
	entries := bindingMgr.BindingsOnFeature(address)
	assert.Equal(suite.T(), 0, len(entries))

	address.Feature = localServerFeature.Address().Feature
	entries = bindingMgr.BindingsOnFeature(address)
	assert.Equal(suite.T(), 1, len(entries))

	bindingRequest2 := model.BindingManagementRequestCallType{
		ClientAddress:     remoteClientFeature.Address(),
		ServerAddress:     localServerFeature2.Address(),
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeMeasurement),
	}

	err = bindingMgr.AddBinding(suite.remoteDevice, bindingRequest2)
	assert.Nil(suite.T(), err)

	address.Feature = localServerFeature2.Address().Feature
	entries = bindingMgr.BindingsOnFeature(address)
	assert.Equal(suite.T(), 1, len(entries))
	entries = bindingMgr.Bindings(suite.remoteDevice)
	assert.Equal(suite.T(), 2, len(entries))

	bindingRequest2 = model.BindingManagementRequestCallType{
		ClientAddress:     remoteClientFeature2.Address(),
		ServerAddress:     localServerFeature3.Address(),
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeLoadControl),
	}

	err = bindingMgr.AddBinding(suite.remoteDevice, bindingRequest2)
	assert.Nil(suite.T(), err)

	address.Feature = localServerFeature3.Address().Feature
	entries = bindingMgr.BindingsOnFeature(address)
	assert.Equal(suite.T(), 1, len(entries))
	entries = bindingMgr.Bindings(suite.remoteDevice)
	assert.Equal(suite.T(), 3, len(entries))

	bindingDelete := model.BindingManagementDeleteCallType{
		ClientAddress: util.Ptr(model.FeatureAddressType{
			Entity:  []model.AddressEntityType{1000},
			Feature: util.Ptr(model.AddressFeatureType(1000)),
		}),
		ServerAddress: util.Ptr(model.FeatureAddressType{
			Entity:  []model.AddressEntityType{1000},
			Feature: util.Ptr(model.AddressFeatureType(1000)),
		}),
	}
	err = bindingMgr.RemoveBinding(bindingDelete, suite.remoteDevice)
	assert.NotNil(suite.T(), err)

	bindingDelete.ClientAddress = remoteServerFeature.Address()

	err = bindingMgr.RemoveBinding(bindingDelete, suite.remoteDevice)
	assert.NotNil(suite.T(), err)

	bindingDelete.ServerAddress = localClientFeature.Address()

	err = bindingMgr.RemoveBinding(bindingDelete, suite.remoteDevice)
	assert.NotNil(suite.T(), err)

	err = bindingMgr.RemoveBinding(bindingDelete, suite.remoteDevice)
	assert.NotNil(suite.T(), err)

	bindingDelete.ServerAddress = localServerFeature.Address()

	err = bindingMgr.RemoveBinding(bindingDelete, suite.remoteDevice)
	assert.NotNil(suite.T(), err)

	bindingDelete.ClientAddress = remoteClientFeature2.Address()

	err = bindingMgr.RemoveBinding(bindingDelete, suite.remoteDevice)
	assert.NotNil(suite.T(), err)

	subs = bindingMgr.Bindings(suite.remoteDevice)
	assert.Equal(suite.T(), 3, len(subs))

	bindingDelete.ClientAddress = remoteClientFeature.Address()

	err = bindingMgr.RemoveBinding(bindingDelete, suite.remoteDevice)
	assert.Nil(suite.T(), err)

	subs = bindingMgr.Bindings(suite.remoteDevice)
	assert.Equal(suite.T(), 2, len(subs))

	err = bindingMgr.RemoveBinding(bindingDelete, suite.remoteDevice)
	assert.NotNil(suite.T(), err)

	err = bindingMgr.AddBinding(suite.remoteDevice, bindingRequest)
	assert.Nil(suite.T(), err)

	subs = bindingMgr.Bindings(suite.remoteDevice)
	assert.Equal(suite.T(), 3, len(subs))

	bindingMgr.RemoveBindingsForDevice(suite.remoteDevice)

	subs = bindingMgr.Bindings(suite.remoteDevice)
	assert.Equal(suite.T(), 0, len(subs))
}
