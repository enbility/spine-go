package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestSetpointListDataType_Update(t *testing.T) {
	sut := SetpointListDataType{
		SetpointData: []SetpointDataType{
			{
				SetpointId: util.Ptr(SetpointIdType(0)),
				Value:      NewScaledNumberType(1),
			},
			{
				SetpointId: util.Ptr(SetpointIdType(1)),
				Value:      NewScaledNumberType(1),
			},
		},
	}

	newData := SetpointListDataType{
		SetpointData: []SetpointDataType{
			{
				SetpointId: util.Ptr(SetpointIdType(1)),
				Value:      NewScaledNumberType(10),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.SetpointData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SetpointId))
	assert.Equal(t, 1.0, item1.Value.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SetpointId))
	assert.Equal(t, 10.0, item2.Value.GetValue())
}

func TestSetpointDescriptionListDataType_Update(t *testing.T) {
	sut := SetpointDescriptionListDataType{
		SetpointDescriptionData: []SetpointDescriptionDataType{
			{
				SetpointId:    util.Ptr(SetpointIdType(0)),
				MeasurementId: util.Ptr(SetpointIdType(0)),
				TimeTableId:   util.Ptr(SetpointIdType(0)),
				Description:   util.Ptr(DescriptionType("old")),
			},
			{
				SetpointId:    util.Ptr(SetpointIdType(1)),
				MeasurementId: util.Ptr(SetpointIdType(1)),
				TimeTableId:   util.Ptr(SetpointIdType(1)),
				Description:   util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := SetpointDescriptionListDataType{
		SetpointDescriptionData: []SetpointDescriptionDataType{
			{
				SetpointId:    util.Ptr(SetpointIdType(1)),
				MeasurementId: util.Ptr(SetpointIdType(1)),
				TimeTableId:   util.Ptr(SetpointIdType(1)),
				Description:   util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.SetpointDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SetpointId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SetpointId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestSetpointListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	testData := &SetpointListDataType{
		SetpointData: []SetpointDataType{
			{
				SetpointId:           util.Ptr(SetpointIdType(1)),
				Value:                NewScaledNumberType(25.5),
				ValueMin:             NewScaledNumberType(10.0),
				ValueMax:             NewScaledNumberType(30.0),
				IsSetpointChangeable: util.Ptr(true),
				IsSetpointActive:     util.Ptr(true),
			},
			{
				SetpointId:           util.Ptr(SetpointIdType(2)),
				Value:                NewScaledNumberType(18.0),
				ValueMin:             NewScaledNumberType(15.0),
				ValueMax:             NewScaledNumberType(25.0),
				IsSetpointChangeable: util.Ptr(false),
				IsSetpointActive:     util.Ptr(false),
			},
		},
	}

	// Test 1: No filter - should return all data
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*SetpointListDataType)
	assert.Len(t, resultData.SetpointData, 2)
	assert.Equal(t, testData, result)

	// Test 2: Filter with specific setpoint ID
	filter := &FilterType{
		SetpointListDataSelectors: &SetpointListDataSelectorsType{
			SetpointId: util.Ptr(SetpointIdType(1)),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*SetpointListDataType)
	assert.Len(t, resultData.SetpointData, 1)
	assert.Equal(t, SetpointIdType(1), *resultData.SetpointData[0].SetpointId)
	assert.Equal(t, 25.5, resultData.SetpointData[0].Value.GetValue())

	// Test 3: Filter with elements (partial fields)
	filter = &FilterType{
		SetpointDataElements: &SetpointDataElementsType{
			SetpointId:           &ElementTagType{},
			Value:                &ScaledNumberElementsType{},
			IsSetpointChangeable: &ElementTagType{},
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*SetpointListDataType)
	assert.Len(t, resultData.SetpointData, 2)
	// Check that only requested elements are included
	for _, item := range resultData.SetpointData {
		assert.NotNil(t, item.SetpointId)
		assert.NotNil(t, item.Value)
		assert.NotNil(t, item.IsSetpointChangeable)
		// Other fields should be nil since not requested in elements
		assert.Nil(t, item.ValueMin)
		assert.Nil(t, item.ValueMax)
		assert.Nil(t, item.IsSetpointActive)
	}

	// Test 4: Combined filter with selector and elements
	filter = &FilterType{
		SetpointListDataSelectors: &SetpointListDataSelectorsType{
			SetpointId: util.Ptr(SetpointIdType(2)),
		},
		SetpointDataElements: &SetpointDataElementsType{
			SetpointId:       &ElementTagType{},
			IsSetpointActive: &ElementTagType{},
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*SetpointListDataType)
	assert.Len(t, resultData.SetpointData, 1)
	assert.Equal(t, SetpointIdType(2), *resultData.SetpointData[0].SetpointId)
	assert.Equal(t, false, *resultData.SetpointData[0].IsSetpointActive)
	assert.Nil(t, resultData.SetpointData[0].Value)
	assert.Nil(t, resultData.SetpointData[0].IsSetpointChangeable)

	// Test 5: No matching filter
	filter = &FilterType{
		SetpointListDataSelectors: &SetpointListDataSelectorsType{
			SetpointId: util.Ptr(SetpointIdType(999)),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*SetpointListDataType)
	assert.Len(t, resultData.SetpointData, 0)

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestSetpointDescriptionListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	testData := &SetpointDescriptionListDataType{
		SetpointDescriptionData: []SetpointDescriptionDataType{
			{
				SetpointId:    util.Ptr(SetpointIdType(1)),
				MeasurementId: util.Ptr(SetpointIdType(10)),
				TimeTableId:   util.Ptr(SetpointIdType(100)),
				Description:   util.Ptr(DescriptionType("Setpoint 1 Description")),
			},
			{
				SetpointId:    util.Ptr(SetpointIdType(2)),
				MeasurementId: util.Ptr(SetpointIdType(20)),
				TimeTableId:   util.Ptr(SetpointIdType(200)),
				Description:   util.Ptr(DescriptionType("Setpoint 2 Description")),
			},
		},
	}

	// Test 1: No filter - should return all data
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*SetpointDescriptionListDataType)
	assert.Len(t, resultData.SetpointDescriptionData, 2)
	assert.Equal(t, testData, result)

	// Test 2: Filter with specific setpoint ID
	filter := &FilterType{
		SetpointDescriptionListDataSelectors: &SetpointDescriptionListDataSelectorsType{
			SetpointId: util.Ptr(SetpointIdType(1)),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*SetpointDescriptionListDataType)
	assert.Len(t, resultData.SetpointDescriptionData, 1)
	assert.Equal(t, SetpointIdType(1), *resultData.SetpointDescriptionData[0].SetpointId)
	assert.Equal(t, "Setpoint 1 Description", string(*resultData.SetpointDescriptionData[0].Description))

	// Test 3: Filter with elements (partial fields)
	filter = &FilterType{
		SetpointDescriptionDataElements: &SetpointDescriptionDataElementsType{
			SetpointId:  &ElementTagType{},
			Description: &ElementTagType{},
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*SetpointDescriptionListDataType)
	assert.Len(t, resultData.SetpointDescriptionData, 2)
	// Check that only requested elements are included
	for _, item := range resultData.SetpointDescriptionData {
		assert.NotNil(t, item.SetpointId)
		assert.NotNil(t, item.Description)
		// Other fields should be nil since not requested in elements
		assert.Nil(t, item.MeasurementId)
		assert.Nil(t, item.TimeTableId)
	}

	// Test 4: Combined filter with selector and elements
	filter = &FilterType{
		SetpointDescriptionListDataSelectors: &SetpointDescriptionListDataSelectorsType{
			SetpointId: util.Ptr(SetpointIdType(2)),
		},
		SetpointDescriptionDataElements: &SetpointDescriptionDataElementsType{
			SetpointId:  &ElementTagType{},
			Description: &ElementTagType{},
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*SetpointDescriptionListDataType)
	assert.Len(t, resultData.SetpointDescriptionData, 1)
	assert.Equal(t, SetpointIdType(2), *resultData.SetpointDescriptionData[0].SetpointId)
	assert.Equal(t, "Setpoint 2 Description", string(*resultData.SetpointDescriptionData[0].Description))
	assert.Nil(t, resultData.SetpointDescriptionData[0].MeasurementId)
	assert.Nil(t, resultData.SetpointDescriptionData[0].TimeTableId)

	// Test 5: No matching filter
	filter = &FilterType{
		SetpointDescriptionListDataSelectors: &SetpointDescriptionListDataSelectorsType{
			SetpointId: util.Ptr(SetpointIdType(999)),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*SetpointDescriptionListDataType)
	assert.Len(t, resultData.SetpointDescriptionData, 0)

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}
