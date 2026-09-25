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

func TestSubscriptionManagerSuite(t *testing.T) {
	suite.Run(t, new(SubscriptionManagerSuite))
}

type SubscriptionManagerSuite struct {
	suite.Suite

	localDevice  api.DeviceLocalInterface
	writeHandler *WriteMessageHandler
	remoteDevice,
	remoteDevice2 api.DeviceRemoteInterface

	sut api.SubscriptionManagerInterface
}

func (s *SubscriptionManagerSuite) BeforeTest(suiteName, testName string) {
	s.localDevice = NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)

	s.writeHandler = &WriteMessageHandler{}

	ski := "test"
	sender := NewSender(s.writeHandler)
	s.remoteDevice = NewDeviceRemote(s.localDevice, ski, sender)
	_ = s.localDevice.SetupRemoteDevice(ski, s.writeHandler)

	ski2 := "test2"
	s.remoteDevice2 = NewDeviceRemote(s.localDevice, ski2, sender)
	_ = s.localDevice.SetupRemoteDevice(ski2, s.writeHandler)

	s.sut = NewSubscriptionManager(s.localDevice)
}

// TC_SPINE_SUBS_002: deleting a subscription is idempotent - removing an already-absent subscription still returns success.
func (s *SubscriptionManagerSuite) Test_Subscriptions() {
	entity := NewEntityLocal(s.localDevice, model.EntityTypeTypeCEM, []model.AddressEntityType{1}, time.Second*4)
	s.localDevice.AddEntity(entity)

	localFeature := entity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)

	remoteEntity := NewEntityRemote(s.remoteDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{1})
	s.remoteDevice.AddEntity(remoteEntity)

	remoteDeviceAddress := model.AddressDeviceType("remoteDevice")
	remoteDeviceAddress2 := model.AddressDeviceType("remoteDevice2")
	s.remoteDevice.UpdateDevice(
		&model.NetworkManagementDeviceDescriptionDataType{
			DeviceAddress: &model.DeviceAddressType{Device: &remoteDeviceAddress},
		},
	)
	s.remoteDevice2.UpdateDevice(
		&model.NetworkManagementDeviceDescriptionDataType{
			DeviceAddress: &model.DeviceAddressType{Device: &remoteDeviceAddress2},
		},
	)

	remoteFeature := NewFeatureRemote(remoteEntity.NextFeatureId(), remoteEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
	remoteFeature.Address().Device = &remoteDeviceAddress
	remoteEntity.AddFeature(remoteFeature)
	remoteEntity.Address().Device = &remoteDeviceAddress

	subscrRequest := model.SubscriptionManagementRequestCallType{
		ClientAddress:     remoteFeature.Address(),
		ServerAddress:     localFeature.Address(),
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeDeviceDiagnosis),
	}

	remoteEntity2 := NewEntityRemote(s.remoteDevice2, model.EntityTypeTypeEVSE, []model.AddressEntityType{1})
	s.remoteDevice2.AddEntity(remoteEntity2)

	remoteFeature2 := NewFeatureRemote(remoteEntity2.NextFeatureId(), remoteEntity2, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
	remoteFeature2.Address().Device = &remoteDeviceAddress2
	remoteEntity2.AddFeature(remoteFeature2)
	remoteEntity2.Address().Device = &remoteDeviceAddress2

	subscrRequest2 := model.SubscriptionManagementRequestCallType{
		ClientAddress:     remoteFeature2.Address(),
		ServerAddress:     localFeature.Address(),
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeDeviceDiagnosis),
	}

	subMgr := s.localDevice.SubscriptionManager()

	err := subMgr.AddSubscription(s.remoteDevice, subscrRequest)
	assert.Nil(s.T(), err)

	subs := subMgr.SubscriptionsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 1, len(subs))

	err = subMgr.AddSubscription(s.remoteDevice, subscrRequest)
	assert.Nil(s.T(), err)

	subs = subMgr.SubscriptionsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 1, len(subs))

	err = subMgr.AddSubscription(s.remoteDevice2, subscrRequest2)
	assert.Nil(s.T(), err)

	subs = subMgr.SubscriptionsForRemoteDevice(s.remoteDevice2)
	assert.Equal(s.T(), 1, len(subs))

	var clientFeatureAddress, serverFeatureAddress model.FeatureAddressType
	util.DeepCopy(remoteFeature.Address(), &clientFeatureAddress)
	util.DeepCopy(localFeature.Address(), &serverFeatureAddress)
	subscrDelete := model.SubscriptionManagementDeleteCallType{
		ClientAddress: &clientFeatureAddress,
		ServerAddress: &serverFeatureAddress,
	}

	err = subMgr.RemoveSubscription(s.remoteDevice, subscrDelete)
	assert.Nil(s.T(), err)

	subs = subMgr.SubscriptionsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 0, len(subs))

	err = subMgr.RemoveSubscription(s.remoteDevice, subscrDelete)
	assert.Nil(s.T(), err)

	subMgr = s.localDevice.SubscriptionManager()
	err = subMgr.AddSubscription(s.remoteDevice, subscrRequest)
	assert.Nil(s.T(), err)

	subs = subMgr.SubscriptionsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 1, len(subs))

	subMgr.RemoveSubscriptionsForRemoteEntity(nil)

	subs = subMgr.SubscriptionsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 1, len(subs))

	subMgr.RemoveSubscriptionsForRemoteDevice(nil)

	subs = subMgr.SubscriptionsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 1, len(subs))

	subMgr.RemoveSubscriptionsForRemoteDevice(s.remoteDevice)

	subs = subMgr.SubscriptionsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 0, len(subs))

	subs = subMgr.SubscriptionsForRemoteDevice(s.remoteDevice2)
	assert.Equal(s.T(), 1, len(subs))

	subscrDelete = model.SubscriptionManagementDeleteCallType{
		ClientAddress: &model.FeatureAddressType{
			Device: &remoteDeviceAddress2,
			Entity: remoteFeature2.address.Entity,
		},
		ServerAddress: &model.FeatureAddressType{
			Device: localFeature.Device().Address(),
			Entity: localFeature.Entity().Address().Entity,
		},
	}

	err = subMgr.RemoveSubscription(s.remoteDevice2, subscrDelete)
	assert.Nil(s.T(), err)

	subs = subMgr.SubscriptionsForRemoteDevice(s.remoteDevice2)
	assert.Equal(s.T(), 0, len(subs))

	err = subMgr.AddSubscription(s.remoteDevice, subscrRequest)
	assert.Nil(s.T(), err)

	subs = subMgr.SubscriptionsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 1, len(subs))

	err = subMgr.AddSubscription(s.remoteDevice2, subscrRequest2)
	assert.Nil(s.T(), err)

	subs = subMgr.SubscriptionsForRemoteDevice(s.remoteDevice2)
	assert.Equal(s.T(), 1, len(subs))

	subscrDelete = model.SubscriptionManagementDeleteCallType{
		ClientAddress: &model.FeatureAddressType{
			Device: &remoteDeviceAddress2,
		},
		ServerAddress: &model.FeatureAddressType{
			Device: localFeature.Device().Address(),
		},
	}

	err = subMgr.RemoveSubscription(s.remoteDevice2, subscrDelete)
	assert.Nil(s.T(), err)

	subs = subMgr.SubscriptionsForRemoteDevice(s.remoteDevice2)
	assert.Equal(s.T(), 0, len(subs))

	subs = subMgr.SubscriptionsForRemoteDevice(s.remoteDevice)
	assert.Equal(s.T(), 1, len(subs))
}
