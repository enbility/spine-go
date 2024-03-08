package spine

import (
	"testing"
	"time"

	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestDeviceLocalSuite(t *testing.T) {
	suite.Run(t, new(DeviceLocalTestSuite))
}

type DeviceLocalTestSuite struct {
	suite.Suite

	lastMessage string
}

var _ shipapi.ShipConnectionDataWriterInterface = (*DeviceLocalTestSuite)(nil)

func (d *DeviceLocalTestSuite) WriteShipMessageWithPayload(msg []byte) {
	d.lastMessage = string(msg)
}

func (d *DeviceLocalTestSuite) Test_RemoveRemoteDevice() {
	sut := NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)

	ski := "test"
	_ = sut.SetupRemoteDevice(ski, d)
	rDevice := sut.RemoteDeviceForSki(ski)
	assert.NotNil(d.T(), rDevice)

	sut.RemoveRemoteDeviceConnection(ski)

	rDevice = sut.RemoteDeviceForSki(ski)
	assert.Nil(d.T(), rDevice)
}

func (d *DeviceLocalTestSuite) Test_RemoteDevice() {
	sut := NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)
	localEntity := NewEntityLocal(sut, model.EntityTypeTypeCEM, NewAddressEntityType([]uint{1}))
	sut.AddEntity(localEntity)

	f := NewFeatureLocal(1, localEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = NewFeatureLocal(2, localEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	localEntity.AddFeature(f)

	ski := "test"
	remote := sut.RemoteDeviceForSki(ski)
	assert.Nil(d.T(), remote)

	devices := sut.RemoteDevices()
	assert.Equal(d.T(), 0, len(devices))

	_ = sut.SetupRemoteDevice(ski, d)
	remote = sut.RemoteDeviceForSki(ski)
	assert.NotNil(d.T(), remote)

	devices = sut.RemoteDevices()
	assert.Equal(d.T(), 1, len(devices))

	entities := sut.Entities()
	assert.Equal(d.T(), 2, len(entities))

	entity1 := sut.Entity([]model.AddressEntityType{1})
	assert.NotNil(d.T(), entity1)

	entity2 := sut.Entity([]model.AddressEntityType{2})
	assert.Nil(d.T(), entity2)

	featureAddress := &model.FeatureAddressType{
		Entity:  []model.AddressEntityType{1},
		Feature: util.Ptr(model.AddressFeatureType(1)),
	}
	feature1 := sut.FeatureByAddress(featureAddress)
	assert.NotNil(d.T(), feature1)

	feature2 := localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	assert.NotNil(d.T(), feature2)

	featureAddress = &model.FeatureAddressType{
		Device:  remote.Address(),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(0)),
	}

	subscription := model.SubscriptionManagementRequestCallType{
		ClientAddress:     featureAddress,
		ServerAddress:     sut.NodeManagement().Address(),
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeNodeManagement),
	}
	err := sut.SubscriptionManager().AddSubscription(remote, subscription)
	assert.Nil(d.T(), err)

	newSubEntity := NewEntityLocal(sut, model.EntityTypeTypeEV, NewAddressEntityType([]uint{1, 1}))
	f = NewFeatureLocal(1, newSubEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeLoadControlLimitListData, true, true)
	newSubEntity.AddFeature(f)

	sut.AddEntity(newSubEntity)
	// A notification should have been sent
	expectedNotifyMsg := `{"datagram":{"header":{"specificationVersion":"1.3.0","addressSource":{"device":"address","entity":[0],"feature":0},"addressDestination":{"entity":[0],"feature":0},"msgCounter":2,"cmdClassifier":"notify"},"payload":{"cmd":[{"function":"nodeManagementDetailedDiscoveryData","filter":[{"cmdControl":{"partial":{}}}],"nodeManagementDetailedDiscoveryData":{"specificationVersionList":{"specificationVersion":["1.3.0"]},"deviceInformation":{"description":{"deviceAddress":{"device":"address"},"deviceType":"EnergyManagementSystem","networkFeatureSet":"smart"}},"entityInformation":[{"description":{"entityAddress":{"device":"address","entity":[1,1]},"entityType":"EV","lastStateChange":"added"}}],"featureInformation":[{"description":{"featureAddress":{"device":"address","entity":[1,1],"feature":1},"featureType":"LoadControl","role":"server","supportedFunction":[{"function":"loadControlLimitListData","possibleOperations":{"read":{},"write":{"partial":{}}}}]}}]}}]}}}`
	assert.Equal(d.T(), expectedNotifyMsg, d.lastMessage)

	entities = sut.Entities()
	assert.Equal(d.T(), 3, len(entities))

	sut.RemoveEntity(newSubEntity)
	// A notification should have been sent
	expectedNotifyMsg = `{"datagram":{"header":{"specificationVersion":"1.3.0","addressSource":{"device":"address","entity":[0],"feature":0},"addressDestination":{"entity":[0],"feature":0},"msgCounter":3,"cmdClassifier":"notify"},"payload":{"cmd":[{"function":"nodeManagementDetailedDiscoveryData","filter":[{"cmdControl":{"partial":{}}}],"nodeManagementDetailedDiscoveryData":{"specificationVersionList":{"specificationVersion":["1.3.0"]},"deviceInformation":{"description":{"deviceAddress":{"device":"address"},"deviceType":"EnergyManagementSystem","networkFeatureSet":"smart"}},"entityInformation":[{"description":{"entityAddress":{"device":"address","entity":[1,1]},"entityType":"EV","lastStateChange":"removed"}}]}}]}}}`
	assert.Equal(d.T(), expectedNotifyMsg, d.lastMessage)

	entities = sut.Entities()
	assert.Equal(d.T(), 2, len(entities))

	sut.RemoveEntity(entity1)
	// A notification should have been sent
	expectedNotifyMsg = `{"datagram":{"header":{"specificationVersion":"1.3.0","addressSource":{"device":"address","entity":[0],"feature":0},"addressDestination":{"entity":[0],"feature":0},"msgCounter":4,"cmdClassifier":"notify"},"payload":{"cmd":[{"function":"nodeManagementDetailedDiscoveryData","filter":[{"cmdControl":{"partial":{}}}],"nodeManagementDetailedDiscoveryData":{"specificationVersionList":{"specificationVersion":["1.3.0"]},"deviceInformation":{"description":{"deviceAddress":{"device":"address"},"deviceType":"EnergyManagementSystem","networkFeatureSet":"smart"}},"entityInformation":[{"description":{"entityAddress":{"device":"address","entity":[1]},"entityType":"CEM","lastStateChange":"removed"}}]}}]}}}`
	assert.Equal(d.T(), expectedNotifyMsg, d.lastMessage)

	entities = sut.Entities()
	assert.Equal(d.T(), 1, len(entities))

	sut.RemoveRemoteDevice(ski)
	remote = sut.RemoteDeviceForSki(ski)
	assert.Nil(d.T(), remote)
}

func (d *DeviceLocalTestSuite) Test_ProcessCmd_NotifyError() {
	sut := NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)
	localEntity := NewEntityLocal(sut, model.EntityTypeTypeCEM, NewAddressEntityType([]uint{1}))
	localFeature := NewFeatureLocal(50, localEntity, model.FeatureTypeTypeSensing, model.RoleTypeClient)
	localEntity.AddFeature(localFeature)
	sut.AddEntity(localEntity)

	ski := "test"
	_ = sut.SetupRemoteDevice(ski, d)
	remote := sut.RemoteDeviceForSki(ski)
	assert.NotNil(d.T(), remote)
	remoteEntity := NewEntityRemote(remote, model.EntityTypeTypeCEM, []model.AddressEntityType{1})
	remoteFeature := NewFeatureRemote(50, remoteEntity, model.FeatureTypeTypeSensing, model.RoleTypeServer)
	remoteEntity.AddFeature(remoteFeature)
	remote.AddEntity(remoteEntity)

	datagram := model.DatagramType{
		Header: model.HeaderType{
			AddressSource:      remoteFeature.Address(),
			AddressDestination: localFeature.Address(),
			MsgCounter:         util.Ptr(model.MsgCounterType(1)),
			CmdClassifier:      util.Ptr(model.CmdClassifierTypeNotify),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{
				{
					Function: util.Ptr(model.FunctionTypeSensingListData),
					Filter:   filterEmptyPartial(),
					SensingListData: &model.SensingListDataType{
						SensingData: []model.SensingDataType{
							{
								Timestamp: model.NewAbsoluteOrRelativeTimeTypeFromTime(time.Now()),
								Value:     model.NewScaledNumberType(99),
							},
						},
					},
				},
			},
		},
	}

	err := sut.ProcessCmd(datagram, remote)
	assert.NotNil(d.T(), err)

}

func (d *DeviceLocalTestSuite) Test_ProcessCmd_Errors() {
	sut := NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)
	localEntity := NewEntityLocal(sut, model.EntityTypeTypeCEM, NewAddressEntityType([]uint{1}))
	sut.AddEntity(localEntity)

	ski := "test"
	_ = sut.SetupRemoteDevice(ski, d)
	remote := sut.RemoteDeviceForSki(ski)
	assert.NotNil(d.T(), remote)

	datagram := model.DatagramType{
		Header: model.HeaderType{
			AddressSource: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("localdevice")),
			},
			AddressDestination: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("localdevice")),
			},
			MsgCounter:    util.Ptr(model.MsgCounterType(1)),
			CmdClassifier: util.Ptr(model.CmdClassifierTypeRead),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{},
		},
	}

	err := sut.ProcessCmd(datagram, remote)
	assert.NotNil(d.T(), err)

	datagram = model.DatagramType{
		Header: model.HeaderType{
			AddressSource: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("localdevice")),
			},
			AddressDestination: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("localdevice")),
			},
			MsgCounter:    util.Ptr(model.MsgCounterType(1)),
			CmdClassifier: util.Ptr(model.CmdClassifierTypeRead),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{
				{},
			},
		},
	}

	err = sut.ProcessCmd(datagram, remote)
	assert.NotNil(d.T(), err)
}

func (d *DeviceLocalTestSuite) Test_ProcessCmd() {
	sut := NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)
	localEntity := NewEntityLocal(sut, model.EntityTypeTypeCEM, NewAddressEntityType([]uint{1}))
	sut.AddEntity(localEntity)

	f1 := NewFeatureLocal(1, localEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient)
	localEntity.AddFeature(f1)
	f2 := NewFeatureLocal(2, localEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	localEntity.AddFeature(f2)
	f3 := NewFeatureLocal(3, localEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	localEntity.AddFeature(f3)

	ski := "test"
	remoteDeviceName := "remote"
	_ = sut.SetupRemoteDevice(ski, d)
	remote := sut.RemoteDeviceForSki(ski)
	assert.NotNil(d.T(), remote)

	detailedData := &model.NodeManagementDetailedDiscoveryDataType{
		DeviceInformation: &model.NodeManagementDetailedDiscoveryDeviceInformationType{
			Description: &model.NetworkManagementDeviceDescriptionDataType{
				DeviceAddress: &model.DeviceAddressType{
					Device: util.Ptr(model.AddressDeviceType(remoteDeviceName)),
				},
			},
		},
		EntityInformation: []model.NodeManagementDetailedDiscoveryEntityInformationType{
			{
				Description: &model.NetworkManagementEntityDescriptionDataType{
					EntityAddress: &model.EntityAddressType{
						Device: util.Ptr(model.AddressDeviceType(remoteDeviceName)),
						Entity: []model.AddressEntityType{1},
					},
					EntityType: util.Ptr(model.EntityTypeTypeEVSE),
				},
			},
		},
		FeatureInformation: []model.NodeManagementDetailedDiscoveryFeatureInformationType{
			{
				Description: &model.NetworkManagementFeatureDescriptionDataType{
					FeatureAddress: &model.FeatureAddressType{
						Device:  util.Ptr(model.AddressDeviceType(remoteDeviceName)),
						Entity:  []model.AddressEntityType{1},
						Feature: util.Ptr(model.AddressFeatureType(1)),
					},
					FeatureType: util.Ptr(model.FeatureTypeTypeElectricalConnection),
					Role:        util.Ptr(model.RoleTypeServer),
				},
			},
		},
	}
	_, err := remote.AddEntityAndFeatures(true, detailedData)
	assert.Nil(d.T(), err)

	datagram := model.DatagramType{
		Header: model.HeaderType{
			AddressSource: &model.FeatureAddressType{
				Device:  util.Ptr(model.AddressDeviceType(remoteDeviceName)),
				Entity:  []model.AddressEntityType{1},
				Feature: util.Ptr(model.AddressFeatureType(1)),
			},
			AddressDestination: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("localdevice")),
				Entity: []model.AddressEntityType{1},
			},
			MsgCounter:    util.Ptr(model.MsgCounterType(1)),
			CmdClassifier: util.Ptr(model.CmdClassifierTypeRead),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{},
		},
	}

	err = sut.ProcessCmd(datagram, remote)
	assert.NotNil(d.T(), err)

	cmd := model.CmdType{
		ElectricalConnectionParameterDescriptionListData: &model.ElectricalConnectionParameterDescriptionListDataType{},
	}

	datagram.Payload.Cmd = append(datagram.Payload.Cmd, cmd)

	err = sut.ProcessCmd(datagram, remote)
	assert.NotNil(d.T(), err)

	datagram.Header.AddressDestination.Feature = util.Ptr(model.AddressFeatureType(1))

	err = sut.ProcessCmd(datagram, remote)
	assert.NotNil(d.T(), err)

	datagram.Header.AddressDestination.Feature = util.Ptr(model.AddressFeatureType(3))

	err = sut.ProcessCmd(datagram, remote)
	assert.Nil(d.T(), err)

	datagram = model.DatagramType{
		Header: model.HeaderType{
			AddressSource: &model.FeatureAddressType{
				Device:  util.Ptr(model.AddressDeviceType(remoteDeviceName)),
				Entity:  []model.AddressEntityType{1},
				Feature: util.Ptr(model.AddressFeatureType(1)),
			},
			AddressDestination: &model.FeatureAddressType{
				Device:  util.Ptr(model.AddressDeviceType("localdevice")),
				Entity:  []model.AddressEntityType{1},
				Feature: util.Ptr(model.AddressFeatureType(3)),
			},
			MsgCounter:    util.Ptr(model.MsgCounterType(1)),
			CmdClassifier: util.Ptr(model.CmdClassifierTypeWrite),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{
				{ElectricalConnectionParameterDescriptionListData: util.Ptr(model.ElectricalConnectionParameterDescriptionListDataType{})},
			},
		},
	}

	err = sut.ProcessCmd(datagram, remote)
	assert.NotNil(d.T(), err)

	f3.AddFunctionType(model.FunctionTypeElectricalConnectionParameterDescriptionListData, true, true)
	err = sut.ProcessCmd(datagram, remote)
	assert.NotNil(d.T(), err)

}
