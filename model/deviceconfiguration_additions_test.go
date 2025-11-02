package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestDeviceConfigurationKeyValueListDataType_Update(t *testing.T) {
	sut := DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(DeviceConfigurationKeyIdType(0)),
				Value: &DeviceConfigurationKeyValueValueType{
					Boolean: util.Ptr(true),
				},
			},
			{
				KeyId:             util.Ptr(DeviceConfigurationKeyIdType(1)),
				IsValueChangeable: util.Ptr(true),
				Value: &DeviceConfigurationKeyValueValueType{
					Boolean: util.Ptr(true),
				},
			},
		},
	}

	newData := DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(DeviceConfigurationKeyIdType(1)),
				Value: &DeviceConfigurationKeyValueValueType{
					Boolean: util.Ptr(false),
				},
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.DeviceConfigurationKeyValueData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.KeyId))
	assert.True(t, *item1.Value.Boolean)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.KeyId))
	assert.False(t, *item2.Value.Boolean)
}

func TestDeviceConfigurationKeyValueDescriptionListDataType_Update(t *testing.T) {
	sut := DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:     util.Ptr(DeviceConfigurationKeyIdType(0)),
				ValueType: util.Ptr(DeviceConfigurationKeyValueTypeTypeBoolean),
			},
			{
				KeyId:     util.Ptr(DeviceConfigurationKeyIdType(1)),
				ValueType: util.Ptr(DeviceConfigurationKeyValueTypeTypeBoolean),
			},
		},
	}

	newData := DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:     util.Ptr(DeviceConfigurationKeyIdType(1)),
				ValueType: util.Ptr(DeviceConfigurationKeyValueTypeTypeString),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.DeviceConfigurationKeyValueDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.KeyId))
	assert.Equal(t, DeviceConfigurationKeyValueTypeTypeBoolean, *item1.ValueType)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.KeyId))
	assert.Equal(t, DeviceConfigurationKeyValueTypeTypeString, *item2.ValueType)
}

func TestDeviceConfigurationKeyValueConstraintsListDataType_Update(t *testing.T) {
	sut := DeviceConfigurationKeyValueConstraintsListDataType{
		DeviceConfigurationKeyValueConstraintsData: []DeviceConfigurationKeyValueConstraintsDataType{
			{
				KeyId: util.Ptr(DeviceConfigurationKeyIdType(0)),
				ValueStepSize: &DeviceConfigurationKeyValueValueType{
					Boolean: util.Ptr(true),
				},
			},
			{
				KeyId: util.Ptr(DeviceConfigurationKeyIdType(1)),
				ValueStepSize: &DeviceConfigurationKeyValueValueType{
					Boolean: util.Ptr(true),
				},
			},
		},
	}

	newData := DeviceConfigurationKeyValueConstraintsListDataType{
		DeviceConfigurationKeyValueConstraintsData: []DeviceConfigurationKeyValueConstraintsDataType{
			{
				KeyId: util.Ptr(DeviceConfigurationKeyIdType(1)),
				ValueStepSize: &DeviceConfigurationKeyValueValueType{
					Boolean: util.Ptr(false),
				},
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.DeviceConfigurationKeyValueConstraintsData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.KeyId))
	assert.Equal(t, true, *item1.ValueStepSize.Boolean)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.KeyId))
	assert.Equal(t, false, *item2.ValueStepSize.Boolean)
}

func TestDeviceConfigurationKeyValueListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	testData := &DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []DeviceConfigurationKeyValueDataType{
			{
				KeyId: Ptr(DeviceConfigurationKeyIdType(1)),
			},
			{
				KeyId: Ptr(DeviceConfigurationKeyIdType(2)),
			},
		},
	}

	// Test 1: No filter - should return all data
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*DeviceConfigurationKeyValueListDataType)
	assert.Len(t, resultData.DeviceConfigurationKeyValueData, 2)

	// Test 2: Filter with specific key ID
	filter := &FilterType{
		DeviceConfigurationKeyValueListDataSelectors: &DeviceConfigurationKeyValueListDataSelectorsType{
			KeyId: Ptr(DeviceConfigurationKeyIdType(1)),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*DeviceConfigurationKeyValueListDataType)
	assert.Len(t, resultData.DeviceConfigurationKeyValueData, 1)
	assert.Equal(t, DeviceConfigurationKeyIdType(1), *resultData.DeviceConfigurationKeyValueData[0].KeyId)

	// Test 3: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestDeviceConfigurationKeyValueDescriptionListDataType_ReadPartialData(t *testing.T) {
	keyId1 := DeviceConfigurationKeyIdType(1)
	keyId2 := DeviceConfigurationKeyIdType(2)

	testData := &DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:       &keyId1,
				KeyName:     util.Ptr(DeviceConfigurationKeyNameType("maxPower")),
				Description: util.Ptr(DescriptionType("Maximum power consumption")),
				Label:       util.Ptr(LabelType("Max Power")),
				Unit:        util.Ptr(UnitOfMeasurementType("W")),
			},
			{
				KeyId:       &keyId2,
				KeyName:     util.Ptr(DeviceConfigurationKeyNameType("minPower")),
				Description: util.Ptr(DescriptionType("Minimum power consumption")),
				Label:       util.Ptr(LabelType("Min Power")),
				Unit:        util.Ptr(UnitOfMeasurementType("W")),
			},
		},
	}

	// Test 1: Filter by key ID
	filter := &FilterType{
		DeviceConfigurationKeyValueDescriptionListDataSelectors: &DeviceConfigurationKeyValueDescriptionListDataSelectorsType{
			KeyId: &keyId1,
		},
	}

	result, success := testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData := result.(*DeviceConfigurationKeyValueDescriptionListDataType)
	assert.Len(t, resultData.DeviceConfigurationKeyValueDescriptionData, 1)
	assert.Equal(t, keyId1, *resultData.DeviceConfigurationKeyValueDescriptionData[0].KeyId)
	assert.Equal(t, "maxPower", string(*resultData.DeviceConfigurationKeyValueDescriptionData[0].KeyName))

	// Test 2: Filter by key name
	filter = &FilterType{
		DeviceConfigurationKeyValueDescriptionListDataSelectors: &DeviceConfigurationKeyValueDescriptionListDataSelectorsType{
			KeyName: util.Ptr(DeviceConfigurationKeyNameType("minPower")),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*DeviceConfigurationKeyValueDescriptionListDataType)
	assert.Len(t, resultData.DeviceConfigurationKeyValueDescriptionData, 1)
	assert.Equal(t, "minPower", string(*resultData.DeviceConfigurationKeyValueDescriptionData[0].KeyName))

	// Test 3: No filter should return all data
	result, success = testData.ReadPartialData(nil)
	assert.True(t, success)
	resultData = result.(*DeviceConfigurationKeyValueDescriptionListDataType)
	assert.Len(t, resultData.DeviceConfigurationKeyValueDescriptionData, 2)

	// Test 4: UpdateList functionality
	updatedItem := DeviceConfigurationKeyValueDescriptionDataType{
		KeyId:       &keyId1,
		KeyName:     util.Ptr(DeviceConfigurationKeyNameType("maxPower")),
		Description: util.Ptr(DescriptionType("Updated maximum power consumption")),
		Label:       util.Ptr(LabelType("Updated Max Power")),
		Unit:        util.Ptr(UnitOfMeasurementType("kW")),
	}

	newData := &DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []DeviceConfigurationKeyValueDescriptionDataType{updatedItem},
	}

	updatedList, success := testData.UpdateList(false, true, newData, nil, nil, Ptr(FunctionTypeDeviceConfigurationKeyValueDescriptionListData))
	assert.True(t, success)
	updatedListData := updatedList.([]DeviceConfigurationKeyValueDescriptionDataType)
	assert.Len(t, updatedListData, 2) // Should have both items

	// Find the updated item
	var foundUpdatedItem *DeviceConfigurationKeyValueDescriptionDataType
	for _, item := range updatedListData {
		if *item.KeyId == keyId1 {
			foundUpdatedItem = &item
			break
		}
	}

	assert.NotNil(t, foundUpdatedItem)
	assert.Equal(t, "Updated maximum power consumption", string(*foundUpdatedItem.Description))
	assert.Equal(t, "kW", string(*foundUpdatedItem.Unit))

	// Test 5: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestDeviceConfigurationKeyValueConstraintsListDataType_ReadPartialData(t *testing.T) {
	keyId1 := DeviceConfigurationKeyIdType(1)
	keyId2 := DeviceConfigurationKeyIdType(2)

	testData := &DeviceConfigurationKeyValueConstraintsListDataType{
		DeviceConfigurationKeyValueConstraintsData: []DeviceConfigurationKeyValueConstraintsDataType{
			{
				KeyId: &keyId1,
			},
			{
				KeyId: &keyId2,
			},
		},
	}

	// Test 1: ReadPartialData with nil filter should return all data
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, testData, result)

	// Test 2: ReadPartialData with filter
	filter := &FilterType{
		DeviceConfigurationKeyValueConstraintsListDataSelectors: &DeviceConfigurationKeyValueConstraintsListDataSelectorsType{
			KeyId: &keyId1,
		},
	}
	_, success = testData.ReadPartialData(filter)
	assert.True(t, success)

	// Test 3: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}
