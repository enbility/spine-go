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

	localDevice api.DeviceLocalInterface
	remoteDevice,
	remoteDevice2 api.DeviceRemoteInterface
	sut api.SubscriptionManagerInterface
}

func (suite *SubscriptionManagerSuite) WriteShipMessageWithPayload([]byte) {}

func (suite *SubscriptionManagerSuite) SetupSuite() {
	suite.localDevice = NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)

	ski := "test"
	sender := NewSender(suite)
	suite.remoteDevice = NewDeviceRemote(suite.localDevice, ski, sender)
	_ = suite.localDevice.SetupRemoteDevice(ski, suite)

	ski2 := "test2"
	suite.remoteDevice2 = NewDeviceRemote(suite.localDevice, ski2, sender)
	_ = suite.localDevice.SetupRemoteDevice(ski2, suite)

	suite.sut = NewSubscriptionManager(suite.localDevice)
}

func (suite *SubscriptionManagerSuite) Test_Subscriptions() {
	entity := NewEntityLocal(suite.localDevice, model.EntityTypeTypeCEM, []model.AddressEntityType{1}, time.Second*4)
	suite.localDevice.AddEntity(entity)

	localFeature := entity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)

	remoteEntity := NewEntityRemote(suite.remoteDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{1})
	suite.remoteDevice.AddEntity(remoteEntity)

	remoteFeature := NewFeatureRemote(remoteEntity.NextFeatureId(), remoteEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
	remoteFeature.Address().Device = util.Ptr(model.AddressDeviceType("remoteDevice"))
	remoteEntity.AddFeature(remoteFeature)
	remoteEntity.Address().Device = util.Ptr(model.AddressDeviceType("remoteDevice"))

	subscrRequest := model.SubscriptionManagementRequestCallType{
		ClientAddress:     remoteFeature.Address(),
		ServerAddress:     localFeature.Address(),
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeDeviceDiagnosis),
	}

	remoteEntity2 := NewEntityRemote(suite.remoteDevice2, model.EntityTypeTypeEVSE, []model.AddressEntityType{1})
	suite.remoteDevice2.AddEntity(remoteEntity2)

	remoteFeature2 := NewFeatureRemote(remoteEntity2.NextFeatureId(), remoteEntity2, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
	remoteFeature2.Address().Device = util.Ptr(model.AddressDeviceType("remoteDevice2"))
	remoteEntity2.AddFeature(remoteFeature2)
	remoteEntity2.Address().Device = util.Ptr(model.AddressDeviceType("remoteDevice2"))

	subscrRequest2 := model.SubscriptionManagementRequestCallType{
		ClientAddress:     remoteFeature2.Address(),
		ServerAddress:     localFeature.Address(),
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeDeviceDiagnosis),
	}

	subMgr := suite.localDevice.SubscriptionManager()
	err := subMgr.AddSubscription(suite.remoteDevice, subscrRequest)
	assert.Nil(suite.T(), err)

	subs := subMgr.Subscriptions(suite.remoteDevice)
	assert.Equal(suite.T(), 1, len(subs))

	err = subMgr.AddSubscription(suite.remoteDevice, subscrRequest)
	assert.NotNil(suite.T(), err)

	subs = subMgr.Subscriptions(suite.remoteDevice)
	assert.Equal(suite.T(), 1, len(subs))

	err = subMgr.AddSubscription(suite.remoteDevice2, subscrRequest2)
	assert.Nil(suite.T(), err)

	subs = subMgr.Subscriptions(suite.remoteDevice2)
	assert.Equal(suite.T(), 1, len(subs))

	subscrDelete := model.SubscriptionManagementDeleteCallType{
		ClientAddress: remoteFeature.Address(),
		ServerAddress: localFeature.Address(),
	}

	err = subMgr.RemoveSubscription(subscrDelete, suite.remoteDevice)
	assert.Nil(suite.T(), err)

	subs = subMgr.Subscriptions(suite.remoteDevice)
	assert.Equal(suite.T(), 0, len(subs))

	err = subMgr.RemoveSubscription(subscrDelete, suite.remoteDevice)
	assert.NotNil(suite.T(), err)

	subMgr = suite.localDevice.SubscriptionManager()
	err = subMgr.AddSubscription(suite.remoteDevice, subscrRequest)
	assert.Nil(suite.T(), err)

	subs = subMgr.Subscriptions(suite.remoteDevice)
	assert.Equal(suite.T(), 1, len(subs))

	subMgr.RemoveSubscriptionsForEntity(nil)

	subs = subMgr.Subscriptions(suite.remoteDevice)
	assert.Equal(suite.T(), 1, len(subs))

	subMgr.RemoveSubscriptionsForDevice(nil)

	subs = subMgr.Subscriptions(suite.remoteDevice)
	assert.Equal(suite.T(), 1, len(subs))

	subMgr.RemoveSubscriptionsForDevice(suite.remoteDevice)

	subs = subMgr.Subscriptions(suite.remoteDevice)
	assert.Equal(suite.T(), 0, len(subs))

	subs = subMgr.Subscriptions(suite.remoteDevice2)
	assert.Equal(suite.T(), 1, len(subs))
}
