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

func (s *BindingManagerSuite) Test_Bindings() {
	entity := NewEntityLocal(s.localDevice, model.EntityTypeTypeCEM, []model.AddressEntityType{1}, time.Second*4)
	s.localDevice.AddEntity(entity)

	localServerFeature := entity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	localServerFeature2 := entity.GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	localServerFeature3 := entity.GetOrAddFeature(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	localClientFeature := entity.GetOrAddFeature(model.FeatureTypeTypeGeneric, model.RoleTypeClient)

	remoteDeviceAddress := model.AddressDeviceType("remoteDevice")
	s.remoteDevice.UpdateDevice(
		&model.NetworkManagementDeviceDescriptionDataType{
			DeviceAddress: &model.DeviceAddressType{Device: &remoteDeviceAddress},
		},
	)

	remoteEntity := NewEntityRemote(s.remoteDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{1})

	remoteClientFeature := NewFeatureRemote(remoteEntity.NextFeatureId(), remoteEntity, model.FeatureTypeTypeGeneric, model.RoleTypeClient)
	remoteClientFeature.Address().Device = util.Ptr(remoteDeviceAddress)
	remoteEntity.AddFeature(remoteClientFeature)

	remoteClientFeature2 := NewFeatureRemote(remoteEntity.NextFeatureId(), remoteEntity, model.FeatureTypeTypeGeneric, model.RoleTypeClient)
	remoteClientFeature2.Address().Device = util.Ptr(remoteDeviceAddress)
	remoteEntity.AddFeature(remoteClientFeature2)

	remoteServerFeature := NewFeatureRemote(remoteEntity.NextFeatureId(), remoteEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	remoteServerFeature.Address().Device = util.Ptr(remoteDeviceAddress)
	remoteEntity.AddFeature(remoteServerFeature)

	s.remoteDevice.AddEntity(remoteEntity)

	bindingMgr := s.localDevice.BindingManager()

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

	err := bindingMgr.AddBinding(s.remoteDevice, bindingRequest)
	assert.NotNil(s.T(), err)

	bindingRequest = model.BindingManagementRequestCallType{
		ClientAddress: util.Ptr(model.FeatureAddressType{
			Device:  util.Ptr(model.AddressDeviceType("dummy")),
			Entity:  []model.AddressEntityType{1000},
			Feature: util.Ptr(model.AddressFeatureType(1000)),
		}),
		ServerAddress: localClientFeature.Address(),
	}

	err = bindingMgr.AddBinding(s.remoteDevice, bindingRequest)
	assert.NotNil(s.T(), err)

	bindingRequest.ServerFeatureType = util.Ptr(model.FeatureTypeTypeDeviceDiagnosis)

	err = bindingMgr.AddBinding(s.remoteDevice, bindingRequest)
	assert.NotNil(s.T(), err)

	bindingRequest.ServerAddress = localServerFeature.Address()

	err = bindingMgr.AddBinding(s.remoteDevice, bindingRequest)
	assert.NotNil(s.T(), err)

	bindingRequest.ClientAddress = remoteServerFeature.Address()

	err = bindingMgr.AddBinding(s.remoteDevice, bindingRequest)
	assert.NotNil(s.T(), err)

	bindingRequest.ClientAddress = remoteClientFeature.Address()

	err = bindingMgr.AddBinding(s.remoteDevice, bindingRequest)
	assert.Nil(s.T(), err)

	subs := bindingMgr.BindingsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 1, len(subs))

        // adding a binding that already exists isn't an error
	err = bindingMgr.AddBinding(s.remoteDevice, bindingRequest)
	assert.Nil(s.T(), err)

	subs = bindingMgr.BindingsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 1, len(subs))

	address := model.FeatureAddressType{
		Device:  entity.Device().Address(),
		Entity:  entity.Address().Entity,
		Feature: util.Ptr(model.AddressFeatureType(10)),
	}
	entries := bindingMgr.BindingsForFeatureAddress(address)
	assert.Equal(s.T(), 0, len(entries))

	address.Feature = localServerFeature.Address().Feature
	entries = bindingMgr.BindingsForFeatureAddress(address)
	assert.Equal(s.T(), 1, len(entries))

	bindingRequest2 := model.BindingManagementRequestCallType{
		ClientAddress:     remoteClientFeature.Address(),
		ServerAddress:     localServerFeature2.Address(),
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeMeasurement),
	}

	err = bindingMgr.AddBinding(s.remoteDevice, bindingRequest2)
	assert.Nil(s.T(), err)

	address.Feature = localServerFeature2.Address().Feature
	entries = bindingMgr.BindingsForFeatureAddress(address)
	assert.Equal(s.T(), 1, len(entries))
	entries = bindingMgr.BindingsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 2, len(entries))

	bindingRequest2 = model.BindingManagementRequestCallType{
		ClientAddress:     remoteClientFeature2.Address(),
		ServerAddress:     localServerFeature3.Address(),
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeLoadControl),
	}

	err = bindingMgr.AddBinding(s.remoteDevice, bindingRequest2)
	assert.Nil(s.T(), err)

	address.Feature = localServerFeature3.Address().Feature
	entries = bindingMgr.BindingsForFeatureAddress(address)
	assert.Equal(s.T(), 1, len(entries))
	entries = bindingMgr.BindingsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 3, len(entries))

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
        // removing a binding that doesn't exist is considered a success
	err = bindingMgr.RemoveBinding(s.remoteDevice, bindingDelete)
	assert.Nil(s.T(), err)

	bindingDelete.ClientAddress = remoteServerFeature.Address()

	err = bindingMgr.RemoveBinding(s.remoteDevice, bindingDelete)
	assert.Nil(s.T(), err)

	bindingDelete.ServerAddress = localClientFeature.Address()

	err = bindingMgr.RemoveBinding(s.remoteDevice, bindingDelete)
	assert.Nil(s.T(), err)

	bindingDelete.ServerAddress = localServerFeature.Address()

	err = bindingMgr.RemoveBinding(s.remoteDevice, bindingDelete)
	assert.Nil(s.T(), err)

	bindingDelete.ClientAddress = remoteClientFeature2.Address()

	err = bindingMgr.RemoveBinding(s.remoteDevice, bindingDelete)
	assert.Nil(s.T(), err)

	subs = bindingMgr.BindingsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 3, len(subs))

	bindingDelete.ClientAddress = remoteClientFeature.Address()

	err = bindingMgr.RemoveBinding(s.remoteDevice, bindingDelete)
	assert.Nil(s.T(), err)

	subs = bindingMgr.BindingsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 2, len(subs))

	err = bindingMgr.AddBinding(s.remoteDevice, bindingRequest)
	assert.Nil(s.T(), err)

	subs = bindingMgr.BindingsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 3, len(subs))

	bindingMgr.RemoveBindingsForRemoteDevice(s.remoteDevice)

	subs = bindingMgr.BindingsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 0, len(subs))

	err = bindingMgr.AddBinding(s.remoteDevice, bindingRequest)
	assert.Nil(s.T(), err)

	subs = bindingMgr.BindingsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 1, len(subs))

	err = bindingMgr.AddBinding(s.remoteDevice, bindingRequest2)
	assert.Nil(s.T(), err)

	subs = bindingMgr.BindingsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 2, len(subs))

	bindingDelete = model.BindingManagementDeleteCallType{
		ClientAddress: &model.FeatureAddressType{
			Device: &remoteDeviceAddress,
			Entity: remoteClientFeature.address.Entity,
		},
		ServerAddress: &model.FeatureAddressType{
			Device: localServerFeature.Device().Address(),
			Entity: localServerFeature.Entity().Address().Entity,
		},
	}

	err = bindingMgr.RemoveBinding(s.remoteDevice, bindingDelete)
	assert.Nil(s.T(), err)

	subs = bindingMgr.BindingsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 0, len(subs))

	err = bindingMgr.AddBinding(s.remoteDevice, bindingRequest)
	assert.Nil(s.T(), err)

	subs = bindingMgr.BindingsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 1, len(subs))

	err = bindingMgr.AddBinding(s.remoteDevice, bindingRequest2)
	assert.Nil(s.T(), err)

	subs = bindingMgr.BindingsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 2, len(subs))

	bindingDelete = model.BindingManagementDeleteCallType{
		ClientAddress: &model.FeatureAddressType{
			Device: &remoteDeviceAddress,
		},
		ServerAddress: &model.FeatureAddressType{
			Device: localServerFeature.Device().Address(),
		},
	}

	err = bindingMgr.RemoveBinding(s.remoteDevice, bindingDelete)
	assert.Nil(s.T(), err)

	subs = bindingMgr.BindingsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 0, len(subs))
}
