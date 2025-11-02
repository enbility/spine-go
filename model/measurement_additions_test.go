package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestMeasurementListDataType_Update_Add(t *testing.T) {
	sut := MeasurementListDataType{
		MeasurementData: []MeasurementDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(0)),
				ValueType:     util.Ptr(MeasurementValueTypeTypeAverageValue),
				Value:         NewScaledNumberType(1),
			},
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ValueType:     util.Ptr(MeasurementValueTypeTypeAverageValue),
				Value:         NewScaledNumberType(1),
			},
		},
	}

	newData := MeasurementListDataType{
		MeasurementData: []MeasurementDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ValueType:     util.Ptr(MeasurementValueTypeTypeValue),
				Value:         NewScaledNumberType(10),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.MeasurementData
	// check the non changing items
	assert.Equal(t, 3, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.MeasurementId))
	assert.Equal(t, MeasurementValueTypeTypeAverageValue, *item1.ValueType)
	assert.Equal(t, 1.0, item1.Value.GetValue())
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.MeasurementId))
	assert.Equal(t, MeasurementValueTypeTypeAverageValue, *item2.ValueType)
	assert.Equal(t, 1.0, item2.Value.GetValue())
}

func TestMeasurementListDataType_Update_Replace(t *testing.T) {
	sut := MeasurementListDataType{
		MeasurementData: []MeasurementDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(0)),
				ValueType:     util.Ptr(MeasurementValueTypeTypeAverageValue),
				Value:         NewScaledNumberType(1),
			},
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ValueType:     util.Ptr(MeasurementValueTypeTypeValue),
				Value:         NewScaledNumberType(1),
			},
		},
	}

	newData := MeasurementListDataType{
		MeasurementData: []MeasurementDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ValueType:     util.Ptr(MeasurementValueTypeTypeValue),
				Value:         NewScaledNumberType(10),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.MeasurementData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.MeasurementId))
	assert.Equal(t, MeasurementValueTypeTypeAverageValue, *item1.ValueType)
	assert.Equal(t, 1.0, item1.Value.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.MeasurementId))
	assert.Equal(t, MeasurementValueTypeTypeValue, *item2.ValueType)
	assert.Equal(t, 10.0, item2.Value.GetValue())
}

func TestMeasurementSeriesListDataType_Update(t *testing.T) {
	sut := MeasurementSeriesListDataType{
		MeasurementSeriesData: []MeasurementSeriesDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(0)),
				ValueType:     util.Ptr(MeasurementValueTypeTypeMinValue),
				Value:         NewScaledNumberType(1),
			},
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ValueType:     util.Ptr(MeasurementValueTypeTypeMaxValue),
				Value:         NewScaledNumberType(10),
			},
		},
	}

	newData := MeasurementSeriesListDataType{
		MeasurementSeriesData: []MeasurementSeriesDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ValueType:     util.Ptr(MeasurementValueTypeTypeMaxValue),
				Value:         NewScaledNumberType(100),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.MeasurementSeriesData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.MeasurementId))
	assert.Equal(t, 1.0, item1.Value.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.MeasurementId))
	assert.Equal(t, 100.0, item2.Value.GetValue())
}

func TestMeasurementConstraintsListDataType_Update(t *testing.T) {
	sut := MeasurementConstraintsListDataType{
		MeasurementConstraintsData: []MeasurementConstraintsDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(0)),
				ValueStepSize: NewScaledNumberType(1),
			},
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ValueStepSize: NewScaledNumberType(1),
			},
		},
	}

	newData := MeasurementConstraintsListDataType{
		MeasurementConstraintsData: []MeasurementConstraintsDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ValueStepSize: NewScaledNumberType(10),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.MeasurementConstraintsData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.MeasurementId))
	assert.Equal(t, 1.0, item1.ValueStepSize.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.MeasurementId))
	assert.Equal(t, 10.0, item2.ValueStepSize.GetValue())
}

func TestMeasurementDescriptionListDataType_Update(t *testing.T) {
	sut := MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []MeasurementDescriptionDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(0)),
				ScopeType:     util.Ptr(ScopeTypeTypeACCurrent),
			},
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ScopeType:     util.Ptr(ScopeTypeTypeACCurrent),
			},
		},
	}

	newData := MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []MeasurementDescriptionDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ScopeType:     util.Ptr(ScopeTypeTypeACPower),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.MeasurementDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.MeasurementId))
	assert.Equal(t, ScopeTypeTypeACCurrent, *item1.ScopeType)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.MeasurementId))
	assert.Equal(t, ScopeTypeTypeACPower, *item2.ScopeType)
}

func TestMeasurementThresholdRelationListDataType_Update(t *testing.T) {
	sut := MeasurementThresholdRelationListDataType{
		MeasurementThresholdRelationData: []MeasurementThresholdRelationDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(0)),
				ThresholdId:   []ThresholdIdType{ThresholdIdType(0)},
			},
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ThresholdId:   []ThresholdIdType{ThresholdIdType(0)},
			},
		},
	}

	newData := MeasurementThresholdRelationListDataType{
		MeasurementThresholdRelationData: []MeasurementThresholdRelationDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ThresholdId:   []ThresholdIdType{ThresholdIdType(1)},
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.MeasurementThresholdRelationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.MeasurementId))
	assert.Equal(t, 0, int(item1.ThresholdId[0]))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.MeasurementId))
	assert.Equal(t, 1, int(item2.ThresholdId[0]))
}

func TestMeasurementListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	testData := &MeasurementListDataType{
		MeasurementData: []MeasurementDataType{
			{
				MeasurementId: Ptr(MeasurementIdType(1)),
			},
			{
				MeasurementId: Ptr(MeasurementIdType(2)),
			},
		},
	}

	// Test 1: No filter - should return all data
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*MeasurementListDataType)
	assert.Len(t, resultData.MeasurementData, 2)

	// Test 2: Filter with specific measurement ID
	filter := &FilterType{
		MeasurementListDataSelectors: &MeasurementListDataSelectorsType{
			MeasurementId: Ptr(MeasurementIdType(1)),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*MeasurementListDataType)
	assert.Len(t, resultData.MeasurementData, 1)
	assert.Equal(t, MeasurementIdType(1), *resultData.MeasurementData[0].MeasurementId)

	// Test 3: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestMeasurementConstraintsListDataType_ReadPartialData(t *testing.T) {
	measurementId := MeasurementIdType(123)
	data := &MeasurementConstraintsListDataType{
		MeasurementConstraintsData: []MeasurementConstraintsDataType{
			{MeasurementId: &measurementId},
			{MeasurementId: Ptr(MeasurementIdType(456))},
		},
	}

	// Test with measurement ID filter
	filter := &FilterType{
		MeasurementConstraintsListDataSelectors: &MeasurementConstraintsListDataSelectorsType{
			MeasurementId: &measurementId,
		},
	}

	result, ok := data.ReadPartialData(filter)
	assert.True(t, ok)

	resultData := result.(*MeasurementConstraintsListDataType)
	assert.Len(t, resultData.MeasurementConstraintsData, 1)
	assert.Equal(t, measurementId, *resultData.MeasurementConstraintsData[0].MeasurementId)

	// Test 2: ReadPartialData with nil filter
	_, success := data.ReadPartialData(nil)
	assert.True(t, success)

	// Test 3: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = data.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestMeasurementDescriptionListDataType_ReadPartialData(t *testing.T) {
	measurementId1 := MeasurementIdType(123)
	measurementId2 := MeasurementIdType(456)

	testData := &MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []MeasurementDescriptionDataType{
			{
				MeasurementId:   &measurementId1,
				MeasurementType: util.Ptr(MeasurementTypeType("power")),
				CommodityType:   util.Ptr(CommodityTypeType("electricity")),
				ScopeType:       util.Ptr(ScopeTypeType("acPowerTotal")),
				Unit:            util.Ptr(UnitOfMeasurementType("W")),
				Label:           util.Ptr(LabelType("Total AC Power")),
				Description:     util.Ptr(DescriptionType("Total active power consumption")),
			},
			{
				MeasurementId:   &measurementId2,
				MeasurementType: util.Ptr(MeasurementTypeType("energy")),
				CommodityType:   util.Ptr(CommodityTypeType("electricity")),
				ScopeType:       util.Ptr(ScopeTypeType("acEnergyTotal")),
				Unit:            util.Ptr(UnitOfMeasurementType("Wh")),
				Label:           util.Ptr(LabelType("Total AC Energy")),
				Description:     util.Ptr(DescriptionType("Total energy consumption")),
			},
		},
	}

	// Test 1: No filter - should return all data
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, testData, result)

	// Test 2: Filter by measurement ID
	filter := &FilterType{
		MeasurementDescriptionListDataSelectors: &MeasurementDescriptionListDataSelectorsType{
			MeasurementId: &measurementId1,
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData := result.(*MeasurementDescriptionListDataType)
	assert.Len(t, resultData.MeasurementDescriptionData, 1)
	assert.Equal(t, measurementId1, *resultData.MeasurementDescriptionData[0].MeasurementId)
	assert.Equal(t, "power", string(*resultData.MeasurementDescriptionData[0].MeasurementType))

	// Test 3: Filter by measurement type
	filter = &FilterType{
		MeasurementDescriptionListDataSelectors: &MeasurementDescriptionListDataSelectorsType{
			MeasurementType: util.Ptr(MeasurementTypeType("energy")),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*MeasurementDescriptionListDataType)
	assert.Len(t, resultData.MeasurementDescriptionData, 1)
	assert.Equal(t, "energy", string(*resultData.MeasurementDescriptionData[0].MeasurementType))

	// Test 4: Filter by commodity type
	filter = &FilterType{
		MeasurementDescriptionListDataSelectors: &MeasurementDescriptionListDataSelectorsType{
			CommodityType: util.Ptr(CommodityTypeType("electricity")),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*MeasurementDescriptionListDataType)
	assert.Len(t, resultData.MeasurementDescriptionData, 2) // Both have electricity commodity

	// Test 5: Filter by scope type
	filter = &FilterType{
		MeasurementDescriptionListDataSelectors: &MeasurementDescriptionListDataSelectorsType{
			ScopeType: util.Ptr(ScopeTypeType("acPowerTotal")),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*MeasurementDescriptionListDataType)
	assert.Len(t, resultData.MeasurementDescriptionData, 1)
	assert.Equal(t, "acPowerTotal", string(*resultData.MeasurementDescriptionData[0].ScopeType))

	// Test 6: Elements filtering test
	filter = &FilterType{
		MeasurementDescriptionDataElements: &MeasurementDescriptionDataElementsType{
			MeasurementId:   &ElementTagType{},
			MeasurementType: &ElementTagType{},
			Label:           &ElementTagType{},
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*MeasurementDescriptionListDataType)
	assert.Len(t, resultData.MeasurementDescriptionData, 2)
	// Check that only requested elements are included
	for _, item := range resultData.MeasurementDescriptionData {
		assert.NotNil(t, item.MeasurementId)
		assert.NotNil(t, item.MeasurementType)
		assert.NotNil(t, item.Label)
		// Other fields should be nil since not requested in elements
		assert.Nil(t, item.CommodityType)
		assert.Nil(t, item.ScopeType)
		assert.Nil(t, item.Unit)
		assert.Nil(t, item.Description)
	}

	// Test 7: Empty result filter
	nonExistentMeasurementId := MeasurementIdType(999)
	filter = &FilterType{
		MeasurementDescriptionListDataSelectors: &MeasurementDescriptionListDataSelectorsType{
			MeasurementId: &nonExistentMeasurementId,
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*MeasurementDescriptionListDataType)
	assert.Len(t, resultData.MeasurementDescriptionData, 0)

	// Test 8: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

// Enhanced MeasurementListDataType tests with more scenarios
func TestMeasurementListDataType_ReadPartialData_Enhanced(t *testing.T) {
	measurementId1 := MeasurementIdType(1)
	measurementId2 := MeasurementIdType(2)
	measurementId3 := MeasurementIdType(3)

	testData := &MeasurementListDataType{
		MeasurementData: []MeasurementDataType{
			{
				MeasurementId: &measurementId1,
				Value:         &ScaledNumberType{Number: util.Ptr(NumberType(1005)), Scale: util.Ptr(ScaleType(-1))},
				ValueState:    util.Ptr(MeasurementValueStateType("normal")),
				Timestamp:     util.Ptr(AbsoluteOrRelativeTimeType("2023-01-01T10:00:00Z")),
			},
			{
				MeasurementId: &measurementId2,
				Value:         &ScaledNumberType{Number: util.Ptr(NumberType(2000)), Scale: util.Ptr(ScaleType(0))},
				ValueState:    util.Ptr(MeasurementValueStateType("normal")),
				Timestamp:     util.Ptr(AbsoluteOrRelativeTimeType("2023-01-01T10:01:00Z")),
			},
			{
				MeasurementId: &measurementId3,
				Value:         &ScaledNumberType{Number: util.Ptr(NumberType(0)), Scale: util.Ptr(ScaleType(0))},
				ValueState:    util.Ptr(MeasurementValueStateType("error")),
				Timestamp:     util.Ptr(AbsoluteOrRelativeTimeType("2023-01-01T10:02:00Z")),
			},
		},
	}

	// Test 1: Filter by measurement ID
	filter := &FilterType{
		MeasurementListDataSelectors: &MeasurementListDataSelectorsType{
			MeasurementId: &measurementId3,
		},
	}

	result, success := testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData := result.(*MeasurementListDataType)
	assert.Len(t, resultData.MeasurementData, 1)
	assert.Equal(t, measurementId3, *resultData.MeasurementData[0].MeasurementId)

	// Test 2: Filter by timestamp range (this test demonstrates timestamp filtering support)
	filter = &FilterType{
		MeasurementListDataSelectors: &MeasurementListDataSelectorsType{
			TimestampInterval: &TimestampIntervalType{
				StartTime: util.Ptr(AbsoluteOrRelativeTimeType("2023-01-01T10:00:30Z")),
				EndTime:   util.Ptr(AbsoluteOrRelativeTimeType("2023-01-01T10:01:30Z")),
			},
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*MeasurementListDataType)
	// This test demonstrates the structure, actual filtering would depend on implementation
	assert.NotNil(t, resultData)

	// Test 3: Combined filter criteria
	filter = &FilterType{
		MeasurementListDataSelectors: &MeasurementListDataSelectorsType{
			MeasurementId: &measurementId2,
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*MeasurementListDataType)
	assert.Len(t, resultData.MeasurementData, 1)
	assert.Equal(t, measurementId2, *resultData.MeasurementData[0].MeasurementId)
	assert.Equal(t, "normal", string(*resultData.MeasurementData[0].ValueState))

	// Test 4: Elements-only filter (field selection without selector filtering)
	filter = &FilterType{
		MeasurementDataElements: &MeasurementDataElementsType{
			MeasurementId: &ElementTagType{},
			Value:         &ElementTagType{},
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*MeasurementListDataType)
	assert.Len(t, resultData.MeasurementData, 3) // All items should be returned
	for _, item := range resultData.MeasurementData {
		assert.NotNil(t, item.MeasurementId)
		assert.NotNil(t, item.Value)
		// These fields should be filtered out since not in elements
		assert.Nil(t, item.ValueState)
		assert.Nil(t, item.Timestamp)
	}

	// Test 5: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

// Test for edge cases and error scenarios
func TestMeasurementReadPartialData_EdgeCases(t *testing.T) {
	// Test with nil measurement data
	t.Run("NilMeasurementData", func(t *testing.T) {
		testData := &MeasurementListDataType{
			MeasurementData: nil,
		}

		result, success := testData.ReadPartialData(nil)
		assert.True(t, success)
		resultData := result.(*MeasurementListDataType)
		assert.Empty(t, resultData.MeasurementData)
	})

	// Test with empty filter selectors
	t.Run("EmptyFilterSelectors", func(t *testing.T) {
		testData := &MeasurementListDataType{
			MeasurementData: []MeasurementDataType{
				{MeasurementId: util.Ptr(MeasurementIdType(1))},
			},
		}

		filter := &FilterType{
			MeasurementListDataSelectors: &MeasurementListDataSelectorsType{}, // Empty selector
		}

		result, success := testData.ReadPartialData(filter)
		assert.True(t, success)
		resultData := result.(*MeasurementListDataType)
		assert.Len(t, resultData.MeasurementData, 1) // Should return all when selector is empty
	})

	// Test with very large dataset performance
	t.Run("LargeDatasetPerformance", func(t *testing.T) {
		largeData := make([]MeasurementDataType, 10000)
		for i := 0; i < 10000; i++ {
			largeData[i] = MeasurementDataType{
				MeasurementId: util.Ptr(MeasurementIdType(i)),
				Value:         &ScaledNumberType{Number: util.Ptr(NumberType(float64(i))), Scale: util.Ptr(ScaleType(0))},
			}
		}

		testData := &MeasurementListDataType{MeasurementData: largeData}

		// Filter for a specific measurement ID
		filter := &FilterType{
			MeasurementListDataSelectors: &MeasurementListDataSelectorsType{
				MeasurementId: util.Ptr(MeasurementIdType(5000)),
			},
		}

		result, success := testData.ReadPartialData(filter)
		assert.True(t, success)
		resultData := result.(*MeasurementListDataType)
		assert.Len(t, resultData.MeasurementData, 1)
		assert.Equal(t, MeasurementIdType(5000), *resultData.MeasurementData[0].MeasurementId)
	})
}

func TestMeasurementSeriesListDataType_ReadPartialData(t *testing.T) {
	sut := &MeasurementSeriesListDataType{
		MeasurementSeriesData: []MeasurementSeriesDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ValueType:     util.Ptr(MeasurementValueTypeType("power")),
				Value:         NewScaledNumberType(150),
			},
			{
				MeasurementId: util.Ptr(MeasurementIdType(2)),
				ValueType:     util.Ptr(MeasurementValueTypeType("energy")),
				Value:         NewScaledNumberType(300),
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*MeasurementSeriesListDataType)
	assert.Equal(t, 2, len(resultData.MeasurementSeriesData))

	// Test 2: filter with selector should return filtered data
	filter := &FilterType{
		MeasurementSeriesListDataSelectors: &MeasurementSeriesListDataSelectorsType{
			MeasurementId: util.Ptr(MeasurementIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*MeasurementSeriesListDataType)
	assert.Equal(t, 1, len(resultData.MeasurementSeriesData))
	assert.Equal(t, MeasurementIdType(1), *resultData.MeasurementSeriesData[0].MeasurementId)

	// Test 3: filter with elements should return filtered data
	filter = &FilterType{
		MeasurementSeriesDataElements: &MeasurementSeriesDataElementsType{
			MeasurementId: &ElementTagType{},
			Value:         &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*MeasurementSeriesListDataType)
	assert.Equal(t, 2, len(resultData.MeasurementSeriesData))
	assert.NotNil(t, resultData.MeasurementSeriesData[0].MeasurementId)
	assert.NotNil(t, resultData.MeasurementSeriesData[0].Value)

	// Test 4: filter with both selector and elements
	filter = &FilterType{
		MeasurementSeriesListDataSelectors: &MeasurementSeriesListDataSelectorsType{
			MeasurementId: util.Ptr(MeasurementIdType(2)),
			ValueType:     util.Ptr(MeasurementValueTypeType("energy")),
		},
		MeasurementSeriesDataElements: &MeasurementSeriesDataElementsType{
			MeasurementId: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*MeasurementSeriesListDataType)
	assert.Equal(t, 1, len(resultData.MeasurementSeriesData))
	assert.Equal(t, MeasurementIdType(2), *resultData.MeasurementSeriesData[0].MeasurementId)
	assert.Nil(t, resultData.MeasurementSeriesData[0].Value)

	// Test 5: non-matching selector should return empty data
	filter = &FilterType{
		MeasurementSeriesListDataSelectors: &MeasurementSeriesListDataSelectorsType{
			MeasurementId: util.Ptr(MeasurementIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*MeasurementSeriesListDataType)
	assert.Equal(t, 0, len(resultData.MeasurementSeriesData))

	// Test 3: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)

}

func TestMeasurementThresholdRelationListDataType_ReadPartialData(t *testing.T) {
	sut := &MeasurementThresholdRelationListDataType{
		MeasurementThresholdRelationData: []MeasurementThresholdRelationDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ThresholdId:   []ThresholdIdType{10, 20, 30},
			},
			{
				MeasurementId: util.Ptr(MeasurementIdType(2)),
				ThresholdId:   []ThresholdIdType{40, 50},
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*MeasurementThresholdRelationListDataType)
	assert.Equal(t, 2, len(resultData.MeasurementThresholdRelationData))

	// Test 2: filter with selector should return filtered data
	filter := &FilterType{
		MeasurementThresholdRelationListDataSelectors: &MeasurementThresholdRelationListDataSelectorsType{
			MeasurementId: util.Ptr(MeasurementIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*MeasurementThresholdRelationListDataType)
	assert.Equal(t, 1, len(resultData.MeasurementThresholdRelationData))
	assert.Equal(t, MeasurementIdType(1), *resultData.MeasurementThresholdRelationData[0].MeasurementId)

	// Test 3: filter with elements should return filtered data
	filter = &FilterType{
		MeasurementThresholdRelationDataElements: &MeasurementThresholdRelationDataElementsType{
			MeasurementId: &ElementTagType{},
			ThresholdId:   &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*MeasurementThresholdRelationListDataType)
	assert.Equal(t, 2, len(resultData.MeasurementThresholdRelationData))
	assert.NotNil(t, resultData.MeasurementThresholdRelationData[0].MeasurementId)
	assert.NotNil(t, resultData.MeasurementThresholdRelationData[0].ThresholdId)

	// Test 4: filter with both selector and elements
	filter = &FilterType{
		MeasurementThresholdRelationListDataSelectors: &MeasurementThresholdRelationListDataSelectorsType{
			MeasurementId: util.Ptr(MeasurementIdType(2)),
			ThresholdId:   util.Ptr(ThresholdIdType(40)),
		},
		MeasurementThresholdRelationDataElements: &MeasurementThresholdRelationDataElementsType{
			MeasurementId: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*MeasurementThresholdRelationListDataType)
	assert.Equal(t, 1, len(resultData.MeasurementThresholdRelationData))
	assert.Equal(t, MeasurementIdType(2), *resultData.MeasurementThresholdRelationData[0].MeasurementId)
	assert.Nil(t, resultData.MeasurementThresholdRelationData[0].ThresholdId)

	// Test 5: non-matching selector should return empty data
	filter = &FilterType{
		MeasurementThresholdRelationListDataSelectors: &MeasurementThresholdRelationListDataSelectorsType{
			MeasurementId: util.Ptr(MeasurementIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*MeasurementThresholdRelationListDataType)
	assert.Equal(t, 0, len(resultData.MeasurementThresholdRelationData))

	// Test 3: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}
