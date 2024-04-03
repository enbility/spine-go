package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestNetworkManagementDeviceDescriptionListDataType(t *testing.T) {
	sut := NetworkManagementDeviceDescriptionListDataType{
		NetworkManagementDeviceDescriptionData: []NetworkManagementDeviceDescriptionDataType{
			{
				DeviceAddress: &DeviceAddressType{
					Device: util.Ptr(AddressDeviceType("test 1")),
				},
				NetworkFeatureSet: util.Ptr(NetworkManagementFeatureSetTypeSimple),
			},
			{
				DeviceAddress: &DeviceAddressType{
					Device: util.Ptr(AddressDeviceType("test 2")),
				},
				NetworkFeatureSet: util.Ptr(NetworkManagementFeatureSetTypeRouter),
			},
		},
	}

	newData := NetworkManagementDeviceDescriptionListDataType{
		NetworkManagementDeviceDescriptionData: []NetworkManagementDeviceDescriptionDataType{
			{
				DeviceAddress: &DeviceAddressType{
					Device: util.Ptr(AddressDeviceType("test 2")),
				},
				NetworkFeatureSet: util.Ptr(NetworkManagementFeatureSetTypeGateway),
			},
		},
	}

	// Act
	success := sut.UpdateList(false, &newData, NewFilterTypePartial(), nil)
	assert.True(t, success)

	data := sut.NetworkManagementDeviceDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, NetworkManagementFeatureSetTypeSimple, *item1.NetworkFeatureSet)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, NetworkManagementFeatureSetTypeGateway, *item2.NetworkFeatureSet)
}

func TestNetworkManagementEntityDescriptionListDataType(t *testing.T) {
	sut := NetworkManagementEntityDescriptionListDataType{
		NetworkManagementEntityDescriptionData: []NetworkManagementEntityDescriptionDataType{
			{
				EntityAddress: &EntityAddressType{
					Device: util.Ptr(AddressDeviceType("test 1")),
					Entity: []AddressEntityType{1},
				},
				EntityType: util.Ptr(EntityTypeTypeBattery),
			},
			{
				EntityAddress: &EntityAddressType{
					Device: util.Ptr(AddressDeviceType("test 1")),
					Entity: []AddressEntityType{2},
				},
				EntityType: util.Ptr(EntityTypeTypeCEM),
			},
		},
	}

	newData := NetworkManagementEntityDescriptionListDataType{
		NetworkManagementEntityDescriptionData: []NetworkManagementEntityDescriptionDataType{
			{
				EntityAddress: &EntityAddressType{
					Device: util.Ptr(AddressDeviceType("test 1")),
					Entity: []AddressEntityType{2},
				},
				EntityType: util.Ptr(EntityTypeTypeEVSE),
			},
		},
	}

	// Act
	success := sut.UpdateList(false, &newData, NewFilterTypePartial(), nil)
	assert.True(t, success)

	data := sut.NetworkManagementEntityDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, EntityTypeTypeBattery, *item1.EntityType)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, EntityTypeTypeEVSE, *item2.EntityType)
}

func TestNetworkManagementFeatureDescriptionListDataType(t *testing.T) {
	sut := NetworkManagementFeatureDescriptionListDataType{
		NetworkManagementFeatureDescriptionData: []NetworkManagementFeatureDescriptionDataType{
			{
				FeatureAddress: &FeatureAddressType{
					Device:  util.Ptr(AddressDeviceType("test 1")),
					Entity:  []AddressEntityType{1},
					Feature: util.Ptr(AddressFeatureType(1)),
				},
				FeatureType: util.Ptr(FeatureTypeTypeActuatorLevel),
			},
			{
				FeatureAddress: &FeatureAddressType{
					Device:  util.Ptr(AddressDeviceType("test 1")),
					Entity:  []AddressEntityType{1},
					Feature: util.Ptr(AddressFeatureType(2)),
				},
				FeatureType: util.Ptr(FeatureTypeTypeAlarm),
			},
		},
	}

	newData := NetworkManagementFeatureDescriptionListDataType{
		NetworkManagementFeatureDescriptionData: []NetworkManagementFeatureDescriptionDataType{
			{
				FeatureAddress: &FeatureAddressType{
					Device:  util.Ptr(AddressDeviceType("test 1")),
					Entity:  []AddressEntityType{1},
					Feature: util.Ptr(AddressFeatureType(2)),
				},
				FeatureType: util.Ptr(FeatureTypeTypeBill),
			},
		},
	}

	// Act
	success := sut.UpdateList(false, &newData, NewFilterTypePartial(), nil)
	assert.True(t, success)

	data := sut.NetworkManagementFeatureDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, FeatureTypeTypeActuatorLevel, *item1.FeatureType)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, FeatureTypeTypeBill, *item2.FeatureType)
}
