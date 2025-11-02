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
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
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
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
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
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
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

func TestNetworkManagementDeviceDescriptionListDataType_ReadPartialData(t *testing.T) {
	testData := &NetworkManagementDeviceDescriptionListDataType{
		NetworkManagementDeviceDescriptionData: []NetworkManagementDeviceDescriptionDataType{
			{
				DeviceAddress: &DeviceAddressType{
					Device: newAddressDeviceType("device1"),
				},
			},
			{
				DeviceAddress: &DeviceAddressType{
					Device: newAddressDeviceType("device2"),
				},
			},
		},
	}

	// Test 1: ReadPartialData with nil filter should return all data
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, testData, result)

	// Test 2: ReadPartialData with filter
	filter := &FilterType{
		NetworkManagementDeviceDescriptionListDataSelectors: &NetworkManagementDeviceDescriptionListDataSelectorsType{
			DeviceAddress: &DeviceAddressType{
				Device: newAddressDeviceType("device1"),
			},
		},
	}
	_, success = testData.ReadPartialData(filter)
	assert.True(t, success)

	// Test 3: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestNetworkManagementEntityDescriptionListDataType_ReadPartialData(t *testing.T) {
	testData := &NetworkManagementEntityDescriptionListDataType{
		NetworkManagementEntityDescriptionData: []NetworkManagementEntityDescriptionDataType{
			{
				EntityType: util.Ptr(EntityTypeTypeBattery),
				Label:      util.Ptr(LabelType("label1")),
			},
			{
				EntityType: util.Ptr(EntityTypeTypeCEM),
				Label:      util.Ptr(LabelType("label2")),
			},
		},
	}

	// Test 1: ReadPartialData with nil filter should return all data
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, testData, result)

	// Test 2: ReadPartialData with simple filter
	filter := &FilterType{
		NetworkManagementEntityDescriptionListDataSelectors: &NetworkManagementEntityDescriptionListDataSelectorsType{},
	}
	_, success = testData.ReadPartialData(filter)
	assert.True(t, success) // Should succeed with empty selector

	// Test 3: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

// Helper functions
func newAddressDeviceType(device string) *AddressDeviceType {
	val := AddressDeviceType(device)
	return &val
}

func newAddressFeatureType(feature uint) *AddressFeatureType {
	val := AddressFeatureType(feature)
	return &val
}

func TestNetworkManagementFeatureDescriptionListDataType_ReadPartialData(t *testing.T) {
	testData := &NetworkManagementFeatureDescriptionListDataType{
		NetworkManagementFeatureDescriptionData: []NetworkManagementFeatureDescriptionDataType{
			{
				FeatureAddress: &FeatureAddressType{
					Device:  newAddressDeviceType("device1"),
					Entity:  []AddressEntityType{1, 1},
					Feature: newAddressFeatureType(1),
				},
				FeatureType: util.Ptr(FeatureTypeType("type1")),
			},
			{
				FeatureAddress: &FeatureAddressType{
					Device:  newAddressDeviceType("device2"),
					Entity:  []AddressEntityType{2, 1},
					Feature: newAddressFeatureType(2),
				},
				FeatureType: util.Ptr(FeatureTypeType("type2")),
			},
		},
	}

	// Test 1: ReadPartialData with nil filter should return all data
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, testData, result)

	// Test 2: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)

	// Test 3: ReadPartialData with filter
	filter := &FilterType{
		NetworkManagementFeatureDescriptionListDataSelectors: &NetworkManagementFeatureDescriptionListDataSelectorsType{
			FeatureAddress: &FeatureAddressType{
				Device: newAddressDeviceType("device1"),
			},
		},
	}
	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData := result.(*NetworkManagementFeatureDescriptionListDataType)
	assert.Len(t, resultData.NetworkManagementFeatureDescriptionData, 1)
	assert.Equal(t, "type1", string(*resultData.NetworkManagementFeatureDescriptionData[0].FeatureType))
}
