package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestThresholdListDataType_Update(t *testing.T) {
	sut := ThresholdListDataType{
		ThresholdData: []ThresholdDataType{
			{
				ThresholdId:    util.Ptr(ThresholdIdType(0)),
				ThresholdValue: NewScaledNumberType(1),
			},
			{
				ThresholdId:    util.Ptr(ThresholdIdType(1)),
				ThresholdValue: NewScaledNumberType(1),
			},
		},
	}

	newData := ThresholdListDataType{
		ThresholdData: []ThresholdDataType{
			{
				ThresholdId:    util.Ptr(ThresholdIdType(1)),
				ThresholdValue: NewScaledNumberType(10),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.ThresholdData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ThresholdId))
	assert.Equal(t, 1.0, item1.ThresholdValue.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ThresholdId))
	assert.Equal(t, 10.0, item2.ThresholdValue.GetValue())
}

func TestThresholdConstraintsListDataType_Update(t *testing.T) {
	sut := ThresholdConstraintsListDataType{
		ThresholdConstraintsData: []ThresholdConstraintsDataType{
			{
				ThresholdId:       util.Ptr(ThresholdIdType(0)),
				ThresholdRangeMin: NewScaledNumberType(1),
			},
			{
				ThresholdId:       util.Ptr(ThresholdIdType(1)),
				ThresholdRangeMin: NewScaledNumberType(1),
			},
		},
	}

	newData := ThresholdConstraintsListDataType{
		ThresholdConstraintsData: []ThresholdConstraintsDataType{
			{
				ThresholdId:       util.Ptr(ThresholdIdType(1)),
				ThresholdRangeMin: NewScaledNumberType(10),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.ThresholdConstraintsData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ThresholdId))
	assert.Equal(t, 1.0, item1.ThresholdRangeMin.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ThresholdId))
	assert.Equal(t, 10.0, item2.ThresholdRangeMin.GetValue())
}

func TestThresholdDescriptionListDataType_Update(t *testing.T) {
	sut := ThresholdDescriptionListDataType{
		ThresholdDescriptionData: []ThresholdDescriptionDataType{
			{
				ThresholdId: util.Ptr(ThresholdIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				ThresholdId: util.Ptr(ThresholdIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := ThresholdDescriptionListDataType{
		ThresholdDescriptionData: []ThresholdDescriptionDataType{
			{
				ThresholdId: util.Ptr(ThresholdIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.ThresholdDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ThresholdId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ThresholdId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestThresholdListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	testData := &ThresholdListDataType{
		ThresholdData: []ThresholdDataType{
			{
				ThresholdId:    util.Ptr(ThresholdIdType(1)),
				ThresholdValue: NewScaledNumberType(25.5),
			},
			{
				ThresholdId:    util.Ptr(ThresholdIdType(2)),
				ThresholdValue: NewScaledNumberType(75.0),
			},
		},
	}

	// Test 1: No filter - should return all data
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*ThresholdListDataType)
	assert.Len(t, resultData.ThresholdData, 2)
	assert.Equal(t, testData, result)

	// Test 2: Filter with specific threshold ID
	filter := &FilterType{
		ThresholdListDataSelectors: &ThresholdListDataSelectorsType{
			ThresholdId: util.Ptr(ThresholdIdType(1)),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*ThresholdListDataType)
	assert.Len(t, resultData.ThresholdData, 1)
	assert.Equal(t, ThresholdIdType(1), *resultData.ThresholdData[0].ThresholdId)
	assert.Equal(t, 25.5, resultData.ThresholdData[0].ThresholdValue.GetValue())

	// Test 3: Filter with elements (partial fields)
	filter = &FilterType{
		ThresholdDataElements: &ThresholdDataElementsType{
			ThresholdId: &ElementTagType{},
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*ThresholdListDataType)
	assert.Len(t, resultData.ThresholdData, 2)
	// Check that only requested elements are included
	for _, item := range resultData.ThresholdData {
		assert.NotNil(t, item.ThresholdId)
		// ThresholdValue should be nil since not requested in elements
		assert.Nil(t, item.ThresholdValue)
	}

	// Test 4: No matching filter
	filter = &FilterType{
		ThresholdListDataSelectors: &ThresholdListDataSelectorsType{
			ThresholdId: util.Ptr(ThresholdIdType(999)),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*ThresholdListDataType)
	assert.Len(t, resultData.ThresholdData, 0)

	// Test 5: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestThresholdConstraintsListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	testData := &ThresholdConstraintsListDataType{
		ThresholdConstraintsData: []ThresholdConstraintsDataType{
			{
				ThresholdId:       util.Ptr(ThresholdIdType(1)),
				ThresholdRangeMin: NewScaledNumberType(10.0),
				ThresholdRangeMax: NewScaledNumberType(50.0),
				ThresholdStepSize: NewScaledNumberType(1.0),
			},
			{
				ThresholdId:       util.Ptr(ThresholdIdType(2)),
				ThresholdRangeMin: NewScaledNumberType(20.0),
				ThresholdRangeMax: NewScaledNumberType(80.0),
				ThresholdStepSize: NewScaledNumberType(5.0),
			},
		},
	}

	// Test 1: No filter - should return all data
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*ThresholdConstraintsListDataType)
	assert.Len(t, resultData.ThresholdConstraintsData, 2)
	assert.Equal(t, testData, result)

	// Test 2: Filter with specific threshold ID
	filter := &FilterType{
		ThresholdConstraintsListDataSelectors: &ThresholdConstraintsListDataSelectorsType{
			ThresholdId: util.Ptr(ThresholdIdType(2)),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*ThresholdConstraintsListDataType)
	assert.Len(t, resultData.ThresholdConstraintsData, 1)
	assert.Equal(t, ThresholdIdType(2), *resultData.ThresholdConstraintsData[0].ThresholdId)
	assert.Equal(t, 20.0, resultData.ThresholdConstraintsData[0].ThresholdRangeMin.GetValue())

	// Test 3: Filter with elements (partial fields)
	filter = &FilterType{
		ThresholdConstraintsDataElements: &ThresholdConstraintsDataElementsType{
			ThresholdId:       &ElementTagType{},
			ThresholdRangeMin: &ScaledNumberElementsType{},
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*ThresholdConstraintsListDataType)
	assert.Len(t, resultData.ThresholdConstraintsData, 2)
	// Check that only requested elements are included
	for _, item := range resultData.ThresholdConstraintsData {
		assert.NotNil(t, item.ThresholdId)
		assert.NotNil(t, item.ThresholdRangeMin)
		// Other fields should be nil since not requested in elements
		assert.Nil(t, item.ThresholdRangeMax)
		assert.Nil(t, item.ThresholdStepSize)
	}

	// Test 4: No matching filter
	filter = &FilterType{
		ThresholdConstraintsListDataSelectors: &ThresholdConstraintsListDataSelectorsType{
			ThresholdId: util.Ptr(ThresholdIdType(999)),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*ThresholdConstraintsListDataType)
	assert.Len(t, resultData.ThresholdConstraintsData, 0)

	// Test 5: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestThresholdDescriptionListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	testData := &ThresholdDescriptionListDataType{
		ThresholdDescriptionData: []ThresholdDescriptionDataType{
			{
				ThresholdId: util.Ptr(ThresholdIdType(1)),
				ScopeType:   util.Ptr(ScopeTypeType("device")),
				Label:       util.Ptr(LabelType("Temperature Threshold")),
				Description: util.Ptr(DescriptionType("High temperature alarm threshold")),
			},
			{
				ThresholdId: util.Ptr(ThresholdIdType(2)),
				ScopeType:   util.Ptr(ScopeTypeType("entity")),
				Label:       util.Ptr(LabelType("Pressure Threshold")),
				Description: util.Ptr(DescriptionType("Low pressure warning threshold")),
			},
		},
	}

	// Test 1: No filter - should return all data
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*ThresholdDescriptionListDataType)
	assert.Len(t, resultData.ThresholdDescriptionData, 2)
	assert.Equal(t, testData, result)

	// Test 2: Filter with specific threshold ID
	filter := &FilterType{
		ThresholdDescriptionListDataSelectors: &ThresholdDescriptionListDataSelectorsType{
			ThresholdId: util.Ptr(ThresholdIdType(1)),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*ThresholdDescriptionListDataType)
	assert.Len(t, resultData.ThresholdDescriptionData, 1)
	assert.Equal(t, ThresholdIdType(1), *resultData.ThresholdDescriptionData[0].ThresholdId)
	assert.Equal(t, "Temperature Threshold", string(*resultData.ThresholdDescriptionData[0].Label))

	// Test 3: Filter with elements (partial fields)
	filter = &FilterType{
		ThresholdDescriptionDataElements: &ThresholdDescriptionDataElementsType{
			ThresholdId: &ElementTagType{},
			Label:       &ElementTagType{},
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*ThresholdDescriptionListDataType)
	assert.Len(t, resultData.ThresholdDescriptionData, 2)
	// Check that only requested elements are included
	for _, item := range resultData.ThresholdDescriptionData {
		assert.NotNil(t, item.ThresholdId)
		assert.NotNil(t, item.Label)
		// Other fields should be nil since not requested in elements
		assert.Nil(t, item.ScopeType)
		assert.Nil(t, item.Description)
	}

	// Test 4: No matching filter
	filter = &FilterType{
		ThresholdDescriptionListDataSelectors: &ThresholdDescriptionListDataSelectorsType{
			ThresholdId: util.Ptr(ThresholdIdType(999)),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*ThresholdDescriptionListDataType)
	assert.Len(t, resultData.ThresholdDescriptionData, 0)

	// Test 5: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}
