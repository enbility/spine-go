package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestStateInformationListDataType_Update(t *testing.T) {
	sut := StateInformationListDataType{
		StateInformationData: []StateInformationDataType{
			{
				StateInformationId: util.Ptr(StateInformationIdType(0)),
				IsActive:           util.Ptr(true),
			},
			{
				StateInformationId: util.Ptr(StateInformationIdType(1)),
				IsActive:           util.Ptr(false),
			},
		},
	}

	newData := StateInformationListDataType{
		StateInformationData: []StateInformationDataType{
			{
				StateInformationId: util.Ptr(StateInformationIdType(1)),
				IsActive:           util.Ptr(true),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.StateInformationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.StateInformationId))
	assert.Equal(t, true, *item1.IsActive)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.StateInformationId))
	assert.Equal(t, true, *item2.IsActive)
}

func TestStateInformationListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	testData := &StateInformationListDataType{
		StateInformationData: []StateInformationDataType{
			{
				StateInformationId: util.Ptr(StateInformationIdType(1)),
				StateInformation:   util.Ptr(StateInformationType("active")),
				IsActive:           util.Ptr(true),
				Category:           util.Ptr(StateInformationCategoryType("operational")),
				TimeOfLastChange:   util.Ptr(AbsoluteOrRelativeTimeType("2023-01-01T12:00:00Z")),
			},
			{
				StateInformationId: util.Ptr(StateInformationIdType(2)),
				StateInformation:   util.Ptr(StateInformationType("inactive")),
				IsActive:           util.Ptr(false),
				Category:           util.Ptr(StateInformationCategoryType("maintenance")),
				TimeOfLastChange:   util.Ptr(AbsoluteOrRelativeTimeType("2023-01-01T13:00:00Z")),
			},
		},
	}

	// Test 1: No filter - should return all data
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*StateInformationListDataType)
	assert.Len(t, resultData.StateInformationData, 2)
	assert.Equal(t, testData, result)

	// Test 2: Filter with specific state information ID
	filter := &FilterType{
		StateInformationListDataSelectors: &StateInformationListDataSelectorsType{
			StateInformationId: util.Ptr(StateInformationIdType(1)),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*StateInformationListDataType)
	assert.Len(t, resultData.StateInformationData, 1)
	assert.Equal(t, StateInformationIdType(1), *resultData.StateInformationData[0].StateInformationId)
	assert.Equal(t, "active", string(*resultData.StateInformationData[0].StateInformation))

	// Test 3: Filter with IsActive
	filter = &FilterType{
		StateInformationListDataSelectors: &StateInformationListDataSelectorsType{
			IsActive: util.Ptr(false),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*StateInformationListDataType)
	assert.Len(t, resultData.StateInformationData, 1)
	assert.Equal(t, false, *resultData.StateInformationData[0].IsActive)
	assert.Equal(t, "maintenance", string(*resultData.StateInformationData[0].Category))

	// Test 4: Filter with elements (partial fields)
	filter = &FilterType{
		StateInformationDataElements: &StateInformationDataElementsType{
			StateInformationId: &ElementTagType{},
			IsActive:           &ElementTagType{},
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*StateInformationListDataType)
	assert.Len(t, resultData.StateInformationData, 2)
	// Check that only requested elements are included
	for _, item := range resultData.StateInformationData {
		assert.NotNil(t, item.StateInformationId)
		assert.NotNil(t, item.IsActive)
		// Other fields should be nil since not requested in elements
		assert.Nil(t, item.StateInformation)
		assert.Nil(t, item.Category)
		assert.Nil(t, item.TimeOfLastChange)
	}

	// Test 5: Combined filter with selector and elements
	filter = &FilterType{
		StateInformationListDataSelectors: &StateInformationListDataSelectorsType{
			Category: util.Ptr(StateInformationCategoryType("operational")),
		},
		StateInformationDataElements: &StateInformationDataElementsType{
			StateInformationId: &ElementTagType{},
			Category:           &ElementTagType{},
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*StateInformationListDataType)
	assert.Len(t, resultData.StateInformationData, 1)
	assert.Equal(t, StateInformationIdType(1), *resultData.StateInformationData[0].StateInformationId)
	assert.Equal(t, "operational", string(*resultData.StateInformationData[0].Category))
	assert.Nil(t, resultData.StateInformationData[0].StateInformation)
	assert.Nil(t, resultData.StateInformationData[0].IsActive)

	// Test 6: No matching filter
	filter = &FilterType{
		StateInformationListDataSelectors: &StateInformationListDataSelectorsType{
			StateInformationId: util.Ptr(StateInformationIdType(999)),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*StateInformationListDataType)
	assert.Len(t, resultData.StateInformationData, 0)

	// Test 7: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}
