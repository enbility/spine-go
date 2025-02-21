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

func TestUtilsSuite(t *testing.T) {
	suite.Run(t, new(UtilsSuite))
}

type UtilsSuite struct {
	suite.Suite

	localDevice  api.DeviceLocalInterface
	remoteDevice api.DeviceRemoteInterface
}

func (s *UtilsSuite) WriteShipMessageWithPayload([]byte) {}

func (s *UtilsSuite) Test_isMatchingClientOrServerByDeviceAndEntity() {
	result := isMatchingClientOrServerByDeviceAndEntity(nil, nil, nil, nil)
	assert.False(s.T(), result)

	clientAddress := &model.FeatureAddressType{}
	serverAddress := &model.FeatureAddressType{}
	result = isMatchingClientOrServerByDeviceAndEntity(clientAddress, serverAddress, nil, nil)
	assert.False(s.T(), result)

	clientAddress = &model.FeatureAddressType{
		Device:  util.Ptr(model.AddressDeviceType("Device1")),
		Entity:  []model.AddressEntityType{1},
		Feature: util.Ptr(model.AddressFeatureType(100)),
	}
	serverAddress = &model.FeatureAddressType{
		Device:  util.Ptr(model.AddressDeviceType("Device2")),
		Entity:  []model.AddressEntityType{1},
		Feature: util.Ptr(model.AddressFeatureType(100)),
	}
	deviceAddress := util.Ptr(model.AddressDeviceType("Device1"))
	entityAddress := []model.AddressEntityType{2}
	result = isMatchingClientOrServerByDeviceAndEntity(clientAddress, serverAddress, deviceAddress, entityAddress)
	assert.False(s.T(), result)

	entityAddress = []model.AddressEntityType{1}
	result = isMatchingClientOrServerByDeviceAndEntity(clientAddress, serverAddress, deviceAddress, entityAddress)
	assert.True(s.T(), result)

	deviceAddress = util.Ptr(model.AddressDeviceType("Device2"))
	entityAddress = []model.AddressEntityType{2}
	result = isMatchingClientOrServerByDeviceAndEntity(clientAddress, serverAddress, deviceAddress, entityAddress)
	assert.False(s.T(), result)

	entityAddress = []model.AddressEntityType{1}
	result = isMatchingClientOrServerByDeviceAndEntity(clientAddress, serverAddress, deviceAddress, entityAddress)
	assert.True(s.T(), result)
}

func (s *UtilsSuite) Test_addressDetails() {
	s.localDevice = NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)

	remoteSki := "TestRemoteSki"
	sender := NewSender(s)
	remoteDevice := NewDeviceRemote(s.localDevice, remoteSki, sender)
	remoteDevice.address = util.Ptr(model.AddressDeviceType("Address"))

	lF, rF, lR, rR, err := addressDetails(s.localDevice, remoteDevice, nil, nil)
	assert.Nil(s.T(), lF)
	assert.Nil(s.T(), rF)
	assert.Equal(s.T(), "", string(lR))
	assert.Equal(s.T(), "", string(rR))
	assert.NotNil(s.T(), err)

	clientAddress := &model.FeatureAddressType{}
	serverAddress := &model.FeatureAddressType{}
	lF, rF, lR, rR, err = addressDetails(s.localDevice, remoteDevice, clientAddress, serverAddress)
	assert.Nil(s.T(), lF)
	assert.Nil(s.T(), rF)
	assert.Equal(s.T(), "", string(lR))
	assert.Equal(s.T(), "", string(rR))
	assert.NotNil(s.T(), err)

	// setup local device
	entity := NewEntityLocal(s.localDevice, model.EntityTypeTypeCEM, []model.AddressEntityType{1}, time.Second*4)
	localClientFeature := entity.GetOrAddFeature(model.FeatureTypeTypeGeneric, model.RoleTypeClient)
	localServerFeature := entity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	s.localDevice.AddEntity(entity)

	// setup remote device
	remoteDeviceAddress := *remoteDevice.Address()
	remoteEntity := NewEntityRemote(remoteDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{1})

	remoteClientFeature := NewFeatureRemote(remoteEntity.NextFeatureId(), remoteEntity, model.FeatureTypeTypeGeneric, model.RoleTypeClient)
	remoteClientFeature.Address().Device = util.Ptr(remoteDeviceAddress)
	remoteEntity.AddFeature(remoteClientFeature)

	remoteServerFeature := NewFeatureRemote(remoteEntity.NextFeatureId(), remoteEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	remoteServerFeature.Address().Device = util.Ptr(remoteDeviceAddress)
	remoteEntity.AddFeature(remoteServerFeature)

	remoteDevice.AddEntity(remoteEntity)

	clientAddress = &model.FeatureAddressType{
		Device:  remoteClientFeature.Address().Device,
		Entity:  remoteClientFeature.Address().Entity,
		Feature: util.Ptr(model.AddressFeatureType(100)),
	}
	serverAddress = &model.FeatureAddressType{
		Device:  localServerFeature.Address().Device,
		Entity:  localServerFeature.Address().Entity,
		Feature: util.Ptr(model.AddressFeatureType(100)),
	}

	lF, rF, lR, rR, err = addressDetails(s.localDevice, remoteDevice, clientAddress, serverAddress)
	assert.Nil(s.T(), lF)
	assert.Nil(s.T(), rF)
	assert.Equal(s.T(), model.RoleTypeServer, lR)
	assert.Equal(s.T(), model.RoleTypeClient, rR)
	assert.NotNil(s.T(), err)

	clientAddress = remoteClientFeature.Address()
	serverAddress = localServerFeature.Address()

	lF, rF, lR, rR, err = addressDetails(s.localDevice, remoteDevice, clientAddress, serverAddress)
	assert.Equal(s.T(), localServerFeature, lF)
	assert.Equal(s.T(), remoteClientFeature, rF)
	assert.Equal(s.T(), model.RoleTypeServer, lR)
	assert.Equal(s.T(), model.RoleTypeClient, rR)
	assert.Nil(s.T(), err)

	clientAddress = localClientFeature.Address()
	serverAddress = remoteServerFeature.Address()

	lF, rF, lR, rR, err = addressDetails(s.localDevice, remoteDevice, clientAddress, serverAddress)
	assert.Equal(s.T(), localClientFeature, lF)
	assert.Equal(s.T(), remoteServerFeature, rF)
	assert.Equal(s.T(), model.RoleTypeClient, lR)
	assert.Equal(s.T(), model.RoleTypeServer, rR)
	assert.Nil(s.T(), err)
}

func (s *UtilsSuite) Test_DataCopyOfType() {
	s.localDevice = NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)
	localEntity := NewEntityLocal(s.localDevice, model.EntityTypeTypeCEM, NewAddressEntityType([]uint{1}), time.Second*4)
	s.localDevice.AddEntity(localEntity)

	localFeature := NewFeatureLocal(1, localEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	localFeature.AddFunctionType(model.FunctionTypeElectricalConnectionDescriptionListData, true, false)
	localEntity.AddFeature(localFeature)
	assert.Equal(s.T(), 1, len(localEntity.Features()))

	// this function does not exist
	_, err := LocalFeatureDataCopyOfType[*model.NodeManagementUseCaseDataType](localFeature, "dummy")
	assert.NotNil(s.T(), err)

	_, err = LocalFeatureDataCopyOfType[*model.ElectricalConnectionDescriptionListDataType](
		localFeature,
		model.FunctionTypeElectricalConnectionDescriptionListData,
	)
	assert.NotNil(s.T(), err)

	data := &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{},
	}
	updatedData, err1 := localFeature.updateData(false, model.FunctionTypeElectricalConnectionDescriptionListData, data, nil, nil)
	assert.NotNil(s.T(), updatedData)
	assert.Nil(s.T(), err1)

	_, err = LocalFeatureDataCopyOfType[*model.ElectricalConnectionDescriptionListDataType](
		localFeature,
		model.FunctionTypeElectricalConnectionDescriptionListData,
	)
	assert.Nil(s.T(), err)

	_, err = LocalFeatureDataCopyOfType[string](
		localFeature,
		model.FunctionTypeElectricalConnectionDescriptionListData,
	)
	assert.NotNil(s.T(), err)

	ski := "test"
	sender := NewSender(s)

	s.remoteDevice = NewDeviceRemote(s.localDevice, ski, sender)
	remoteEntity := NewEntityRemote(s.remoteDevice, model.EntityTypeTypeCEM, NewAddressEntityType([]uint{1}))
	s.remoteDevice.AddEntity(remoteEntity)

	remoteFeature := NewFeatureRemote(1, remoteEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	remoteEntity.AddFeature(remoteFeature)

	// this function does not exist
	_, err = RemoteFeatureDataCopyOfType[*model.NodeManagementUseCaseDataType](remoteFeature, "dummy")
	assert.NotNil(s.T(), err)
}
