package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestBindingManagementEntryListDataType_Update(t *testing.T) {
	sut := BindingManagementEntryListDataType{
		BindingManagementEntryData: []BindingManagementEntryDataType{
			{
				BindingId:   util.Ptr(BindingIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				BindingId:   util.Ptr(BindingIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := BindingManagementEntryListDataType{
		BindingManagementEntryData: []BindingManagementEntryDataType{
			{
				BindingId:   util.Ptr(BindingIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.BindingManagementEntryData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.BindingId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.BindingId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestBindingManagementEntryListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	testData := &BindingManagementEntryListDataType{
		BindingManagementEntryData: []BindingManagementEntryDataType{
			{
				ClientAddress: &FeatureAddressType{
					Device:  newAddressDeviceType("client1"),
					Entity:  []AddressEntityType{1, 1},
					Feature: newAddressFeatureType(1),
				},
				ServerAddress: &FeatureAddressType{
					Device:  newAddressDeviceType("server1"),
					Entity:  []AddressEntityType{2, 1},
					Feature: newAddressFeatureType(2),
				},
			},
			{
				ClientAddress: &FeatureAddressType{
					Device:  newAddressDeviceType("client2"),
					Entity:  []AddressEntityType{1, 2},
					Feature: newAddressFeatureType(3),
				},
				ServerAddress: &FeatureAddressType{
					Device:  newAddressDeviceType("server2"),
					Entity:  []AddressEntityType{2, 2},
					Feature: newAddressFeatureType(4),
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
		BindingManagementEntryListDataSelectors: &BindingManagementEntryListDataSelectorsType{
			ClientAddress: &FeatureAddressType{
				Device: newAddressDeviceType("client1"),
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
