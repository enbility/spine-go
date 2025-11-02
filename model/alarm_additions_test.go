package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestAlarmListDataType_Update(t *testing.T) {
	sut := AlarmListDataType{
		AlarmListData: []AlarmDataType{
			{
				AlarmId:     util.Ptr(AlarmIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				AlarmId:     util.Ptr(AlarmIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := AlarmListDataType{
		AlarmListData: []AlarmDataType{
			{
				AlarmId:     util.Ptr(AlarmIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.AlarmListData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.AlarmId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.AlarmId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestAlarmListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	testData := &AlarmListDataType{
		AlarmListData: []AlarmDataType{
			{
				AlarmId:     util.Ptr(AlarmIdType(1)),
				Description: util.Ptr(DescriptionType("Alarm 1")),
				AlarmType:   util.Ptr(AlarmTypeType("warning")),
				Label:       util.Ptr(LabelType("Test Alarm 1")),
				ScopeType:   util.Ptr(ScopeTypeType("device")),
			},
			{
				AlarmId:     util.Ptr(AlarmIdType(2)),
				Description: util.Ptr(DescriptionType("Alarm 2")),
				AlarmType:   util.Ptr(AlarmTypeType("error")),
				Label:       util.Ptr(LabelType("Test Alarm 2")),
				ScopeType:   util.Ptr(ScopeTypeType("entity")),
			},
		},
	}

	// Test 1: No filter - should return all data
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*AlarmListDataType)
	assert.Len(t, resultData.AlarmListData, 2)
	assert.Equal(t, testData, result)

	// Test 2: Filter with specific alarm ID
	filter := &FilterType{
		AlarmListDataSelectors: &AlarmListDataSelectorsType{
			AlarmId: util.Ptr(AlarmIdType(1)),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*AlarmListDataType)
	assert.Len(t, resultData.AlarmListData, 1)
	assert.Equal(t, AlarmIdType(1), *resultData.AlarmListData[0].AlarmId)
	assert.Equal(t, "Alarm 1", string(*resultData.AlarmListData[0].Description))

	// Test 3: Filter with scope type
	filter = &FilterType{
		AlarmListDataSelectors: &AlarmListDataSelectorsType{
			ScopeType: util.Ptr(ScopeTypeType("device")),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*AlarmListDataType)
	assert.Len(t, resultData.AlarmListData, 1)
	assert.Equal(t, "device", string(*resultData.AlarmListData[0].ScopeType))
	assert.Equal(t, "Test Alarm 1", string(*resultData.AlarmListData[0].Label))

	// Test 4: Filter with elements (partial fields)
	filter = &FilterType{
		AlarmDataElements: &AlarmDataElementsType{
			AlarmId:     &ElementTagType{},
			Description: &ElementTagType{},
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*AlarmListDataType)
	assert.Len(t, resultData.AlarmListData, 2)
	// Check that only requested elements are included
	for _, item := range resultData.AlarmListData {
		assert.NotNil(t, item.AlarmId)
		assert.NotNil(t, item.Description)
		// AlarmType and Label should be nil since not requested in elements
		assert.Nil(t, item.AlarmType)
		assert.Nil(t, item.Label)
	}

	// Test 5: Combined filter with selector and elements
	filter = &FilterType{
		AlarmListDataSelectors: &AlarmListDataSelectorsType{
			AlarmId: util.Ptr(AlarmIdType(2)),
		},
		AlarmDataElements: &AlarmDataElementsType{
			AlarmId: &ElementTagType{},
			Label:   &ElementTagType{},
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*AlarmListDataType)
	assert.Len(t, resultData.AlarmListData, 1)
	assert.Equal(t, AlarmIdType(2), *resultData.AlarmListData[0].AlarmId)
	assert.Equal(t, "Test Alarm 2", string(*resultData.AlarmListData[0].Label))
	assert.Nil(t, resultData.AlarmListData[0].Description)
	assert.Nil(t, resultData.AlarmListData[0].AlarmType)

	// Test 6: No matching filter
	filter = &FilterType{
		AlarmListDataSelectors: &AlarmListDataSelectorsType{
			AlarmId: util.Ptr(AlarmIdType(999)),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*AlarmListDataType)
	assert.Len(t, resultData.AlarmListData, 0)

	// Test 7: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}
