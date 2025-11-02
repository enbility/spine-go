package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestTimeSeriesListDataType_Update(t *testing.T) {
	sut := TimeSeriesListDataType{
		TimeSeriesData: []TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(0)),
				TimeSeriesSlot: []TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(0)),
					},
				},
			},
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(1)),
				TimeSeriesSlot: []TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(0)),
					},
				},
			},
		},
	}

	newData := TimeSeriesListDataType{
		TimeSeriesData: []TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(1)),
				TimeSeriesSlot: []TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(1)),
					},
				},
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.TimeSeriesData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TimeSeriesId))
	assert.Equal(t, 0, int(*item1.TimeSeriesSlot[0].TimeSeriesSlotId))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TimeSeriesId))
	assert.Equal(t, 1, int(*item2.TimeSeriesSlot[0].TimeSeriesSlotId))
}

func TestTimeSeriesListDataType_Update_02(t *testing.T) {
	sut := TimeSeriesListDataType{
		TimeSeriesData: []TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(1)),
				TimePeriod: &TimePeriodType{
					StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT0S")),
					EndTime:   util.Ptr(AbsoluteOrRelativeTimeType("P6D")),
				},
				TimeSeriesSlot: []TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(0)),
						TimePeriod: &TimePeriodType{
							StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT0S")),
							EndTime:   util.Ptr(AbsoluteOrRelativeTimeType("P6D")),
						},
						MaxValue: NewScaledNumberType(10000),
					},
				},
			},
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(2)),
				TimePeriod: &TimePeriodType{
					StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT0S")),
				},
				TimeSeriesSlot: []TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(0)),
						Duration:         util.Ptr(DurationType("P1DT6H46M33S")),
						MaxValue:         NewScaledNumberType(0),
					},
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(1)),
						Duration:         util.Ptr(DurationType("PT7H37M53S")),
						MaxValue:         NewScaledNumberType(4410),
					},
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(2)),
						Duration:         util.Ptr(DurationType("PT38M")),
						MaxValue:         NewScaledNumberType(0),
					},
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(3)),
						Duration:         util.Ptr(DurationType("PT32M")),
						MaxValue:         NewScaledNumberType(4410),
					},
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(4)),
						Duration:         util.Ptr(DurationType("P1D")),
						MaxValue:         NewScaledNumberType(0),
					},
				},
			},
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(3)),
				TimePeriod: &TimePeriodType{
					StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT0S")),
				},
				TimeSeriesSlot: []TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(1)),
						Duration:         util.Ptr(DurationType("P1DT15H24M57S")),
						Value:            NewScaledNumberType(44229),
						MaxValue:         NewScaledNumberType(49629),
					},
				},
			},
		},
	}

	newData := TimeSeriesListDataType{
		TimeSeriesData: []TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(3)),
				TimePeriod: &TimePeriodType{
					StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT0S")),
				},
				TimeSeriesSlot: []TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(1)),
						Duration:         util.Ptr(DurationType("P1DT15H16M50S")),
						Value:            NewScaledNumberType(11539),
						MaxValue:         NewScaledNumberType(49629),
					},
				},
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.TimeSeriesData
	// check the non changing items
	assert.Equal(t, 3, len(data))
	item1 := data[0]
	assert.Equal(t, 1, int(*item1.TimeSeriesId))
	assert.Equal(t, 0, int(*item1.TimeSeriesSlot[0].TimeSeriesSlotId))
	// check properties of updated item
	item2 := data[2]
	assert.Equal(t, 3, int(*item2.TimeSeriesId))
	assert.Equal(t, 1, int(*item2.TimeSeriesSlot[0].TimeSeriesSlotId))
	assert.Equal(t, 11539, int(item2.TimeSeriesSlot[0].Value.GetValue()))
}

func TestTimeSeriesDescriptionListDataType_Update(t *testing.T) {
	sut := TimeSeriesDescriptionListDataType{
		TimeSeriesDescriptionData: []TimeSeriesDescriptionDataType{
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(0)),
				Description:  util.Ptr(DescriptionType("old")),
			},
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(1)),
				Description:  util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := TimeSeriesDescriptionListDataType{
		TimeSeriesDescriptionData: []TimeSeriesDescriptionDataType{
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(1)),
				Description:  util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.TimeSeriesDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TimeSeriesId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TimeSeriesId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestTimeSeriesConstraintsListDataType_Update(t *testing.T) {
	sut := TimeSeriesConstraintsListDataType{
		TimeSeriesConstraintsData: []TimeSeriesConstraintsDataType{
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(0)),
				SlotValueMin: NewScaledNumberType(1),
			},
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(1)),
				SlotValueMin: NewScaledNumberType(1),
			},
		},
	}

	newData := TimeSeriesConstraintsListDataType{
		TimeSeriesConstraintsData: []TimeSeriesConstraintsDataType{
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(1)),
				SlotValueMin: NewScaledNumberType(10),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.TimeSeriesConstraintsData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TimeSeriesId))
	assert.Equal(t, 1.0, item1.SlotValueMin.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TimeSeriesId))
	assert.Equal(t, 10.0, item2.SlotValueMin.GetValue())
}

func TestTimeSeriesListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	data := &TimeSeriesListDataType{
		TimeSeriesData: []TimeSeriesDataType{
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(1)),
				TimePeriod: &TimePeriodType{
					StartTime: util.Ptr(AbsoluteOrRelativeTimeType("start1")),
					EndTime:   util.Ptr(AbsoluteOrRelativeTimeType("end1")),
				},
				TimeSeriesSlot: []TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(1)),
						Value:            NewScaledNumberType(10),
					},
				},
			},
			{
				TimeSeriesId: util.Ptr(TimeSeriesIdType(2)),
				TimePeriod: &TimePeriodType{
					StartTime: util.Ptr(AbsoluteOrRelativeTimeType("start2")),
					EndTime:   util.Ptr(AbsoluteOrRelativeTimeType("end2")),
				},
				TimeSeriesSlot: []TimeSeriesSlotType{
					{
						TimeSeriesSlotId: util.Ptr(TimeSeriesSlotIdType(2)),
						Value:            NewScaledNumberType(20),
					},
				},
			},
		},
	}

	// Test 1: No filter (should return all data)
	result, success := data.ReadPartialData(nil)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData := result.(*TimeSeriesListDataType)
	assert.Len(t, resultData.TimeSeriesData, 2)

	// Test 2: Selector filter
	selectorFilter := &FilterType{
		TimeSeriesListDataSelectors: &TimeSeriesListDataSelectorsType{
			TimeSeriesId: util.Ptr(TimeSeriesIdType(1)),
		},
	}
	result, success = data.ReadPartialData(selectorFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeSeriesListDataType)
	assert.Len(t, resultData.TimeSeriesData, 1)
	assert.Equal(t, TimeSeriesIdType(1), *resultData.TimeSeriesData[0].TimeSeriesId)

	// Test 3: Element filter
	elementFilter := &FilterType{
		TimeSeriesDataElements: &TimeSeriesDataElementsType{
			TimeSeriesId: util.Ptr(ElementTagType{}),
			TimePeriod:   &TimePeriodElementsType{},
		},
	}
	result, success = data.ReadPartialData(elementFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeSeriesListDataType)
	assert.Len(t, resultData.TimeSeriesData, 2)
	// Check that only requested fields are present
	assert.NotNil(t, resultData.TimeSeriesData[0].TimeSeriesId)
	assert.NotNil(t, resultData.TimeSeriesData[0].TimePeriod)
	// TimeSeriesSlot should be nil since not requested in elements
	assert.Nil(t, resultData.TimeSeriesData[0].TimeSeriesSlot)

	// Test 4: Combined selector and element filter
	combinedFilter := &FilterType{
		TimeSeriesListDataSelectors: &TimeSeriesListDataSelectorsType{
			TimeSeriesId: util.Ptr(TimeSeriesIdType(2)),
		},
		TimeSeriesDataElements: &TimeSeriesDataElementsType{
			TimeSeriesId: util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(combinedFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeSeriesListDataType)
	assert.Len(t, resultData.TimeSeriesData, 1)
	assert.Equal(t, TimeSeriesIdType(2), *resultData.TimeSeriesData[0].TimeSeriesId)
	assert.Nil(t, resultData.TimeSeriesData[0].TimePeriod)
	// TimeSeriesSlot should also be nil since not requested
	assert.Nil(t, resultData.TimeSeriesData[0].TimeSeriesSlot)

	// Test 5: Element filter with different combination
	simpleElementFilter := &FilterType{
		TimeSeriesDataElements: &TimeSeriesDataElementsType{
			TimeSeriesId: util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(simpleElementFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeSeriesListDataType)
	assert.Len(t, resultData.TimeSeriesData, 2)
	// Only TimeSeriesId should be present
	assert.NotNil(t, resultData.TimeSeriesData[0].TimeSeriesId)
	assert.Nil(t, resultData.TimeSeriesData[0].TimePeriod)
	assert.Nil(t, resultData.TimeSeriesData[0].TimeSeriesSlot)

	// Test 6: Non-matching selector (should return empty)
	nonMatchingFilter := &FilterType{
		TimeSeriesListDataSelectors: &TimeSeriesListDataSelectorsType{
			TimeSeriesId: util.Ptr(TimeSeriesIdType(999)),
		},
	}
	result, success = data.ReadPartialData(nonMatchingFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeSeriesListDataType)
	assert.Len(t, resultData.TimeSeriesData, 0)

	// Test 7: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = data.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestTimeSeriesDescriptionListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	data := &TimeSeriesDescriptionListDataType{
		TimeSeriesDescriptionData: []TimeSeriesDescriptionDataType{
			{
				TimeSeriesId:        util.Ptr(TimeSeriesIdType(1)),
				TimeSeriesType:      util.Ptr(TimeSeriesTypeType("type1")),
				TimeSeriesWriteable: util.Ptr(true),
				MeasurementId:       util.Ptr(MeasurementIdType(100)),
				Label:               util.Ptr(LabelType("label1")),
				ScopeType:           util.Ptr(ScopeTypeType("scope1")),
			},
			{
				TimeSeriesId:        util.Ptr(TimeSeriesIdType(2)),
				TimeSeriesType:      util.Ptr(TimeSeriesTypeType("type2")),
				TimeSeriesWriteable: util.Ptr(false),
				MeasurementId:       util.Ptr(MeasurementIdType(200)),
				Label:               util.Ptr(LabelType("label2")),
				ScopeType:           util.Ptr(ScopeTypeType("scope2")),
			},
		},
	}

	// Test 1: No filter (should return all data)
	result, success := data.ReadPartialData(nil)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData := result.(*TimeSeriesDescriptionListDataType)
	assert.Len(t, resultData.TimeSeriesDescriptionData, 2)

	// Test 2: Selector filter by TimeSeriesId
	selectorFilter := &FilterType{
		TimeSeriesDescriptionListDataSelectors: &TimeSeriesDescriptionListDataSelectorsType{
			TimeSeriesId: util.Ptr(TimeSeriesIdType(1)),
		},
	}
	result, success = data.ReadPartialData(selectorFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeSeriesDescriptionListDataType)
	assert.Len(t, resultData.TimeSeriesDescriptionData, 1)
	assert.Equal(t, TimeSeriesIdType(1), *resultData.TimeSeriesDescriptionData[0].TimeSeriesId)

	// Test 3: Selector filter by TimeSeriesType
	typeFilter := &FilterType{
		TimeSeriesDescriptionListDataSelectors: &TimeSeriesDescriptionListDataSelectorsType{
			TimeSeriesType: util.Ptr(TimeSeriesTypeType("type2")),
		},
	}
	result, success = data.ReadPartialData(typeFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeSeriesDescriptionListDataType)
	assert.Len(t, resultData.TimeSeriesDescriptionData, 1)
	assert.Equal(t, TimeSeriesTypeType("type2"), *resultData.TimeSeriesDescriptionData[0].TimeSeriesType)

	// Test 4: Element filter
	elementFilter := &FilterType{
		TimeSeriesDescriptionDataElements: &TimeSeriesDescriptionDataElementsType{
			TimeSeriesId:   util.Ptr(ElementTagType{}),
			TimeSeriesType: util.Ptr(ElementTagType{}),
			Label:          util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(elementFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeSeriesDescriptionListDataType)
	assert.Len(t, resultData.TimeSeriesDescriptionData, 2)
	// Check that only requested fields are present
	assert.NotNil(t, resultData.TimeSeriesDescriptionData[0].TimeSeriesId)
	assert.NotNil(t, resultData.TimeSeriesDescriptionData[0].TimeSeriesType)
	assert.NotNil(t, resultData.TimeSeriesDescriptionData[0].Label)
	// Check that non-requested field is filtered out
	assert.Nil(t, resultData.TimeSeriesDescriptionData[0].TimeSeriesWriteable)
	assert.Nil(t, resultData.TimeSeriesDescriptionData[0].MeasurementId)
	assert.Nil(t, resultData.TimeSeriesDescriptionData[0].ScopeType)

	// Test 5: Combined selector and element filter
	combinedFilter := &FilterType{
		TimeSeriesDescriptionListDataSelectors: &TimeSeriesDescriptionListDataSelectorsType{
			MeasurementId: util.Ptr(MeasurementIdType(200)),
		},
		TimeSeriesDescriptionDataElements: &TimeSeriesDescriptionDataElementsType{
			TimeSeriesId:  util.Ptr(ElementTagType{}),
			MeasurementId: util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(combinedFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeSeriesDescriptionListDataType)
	assert.Len(t, resultData.TimeSeriesDescriptionData, 1)
	assert.Equal(t, TimeSeriesIdType(2), *resultData.TimeSeriesDescriptionData[0].TimeSeriesId)
	assert.Equal(t, MeasurementIdType(200), *resultData.TimeSeriesDescriptionData[0].MeasurementId)
	assert.Nil(t, resultData.TimeSeriesDescriptionData[0].TimeSeriesType)
	assert.Nil(t, resultData.TimeSeriesDescriptionData[0].Label)
	assert.Nil(t, resultData.TimeSeriesDescriptionData[0].ScopeType)

	// Test 6: Non-matching selector (should return empty)
	nonMatchingFilter := &FilterType{
		TimeSeriesDescriptionListDataSelectors: &TimeSeriesDescriptionListDataSelectorsType{
			TimeSeriesId: util.Ptr(TimeSeriesIdType(999)),
		},
	}
	result, success = data.ReadPartialData(nonMatchingFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeSeriesDescriptionListDataType)
	assert.Len(t, resultData.TimeSeriesDescriptionData, 0)

	// Test 7: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = data.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestTimeSeriesConstraintsListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	data := &TimeSeriesConstraintsListDataType{
		TimeSeriesConstraintsData: []TimeSeriesConstraintsDataType{
			{
				TimeSeriesId:    util.Ptr(TimeSeriesIdType(1)),
				SlotCountMin:    util.Ptr(TimeSeriesSlotCountType(1)),
				SlotCountMax:    util.Ptr(TimeSeriesSlotCountType(10)),
				SlotDurationMin: util.Ptr(DurationType("PT1M")),
				SlotDurationMax: util.Ptr(DurationType("PT1H")),
				SlotValueMin:    NewScaledNumberType(0),
				SlotValueMax:    NewScaledNumberType(100),
			},
			{
				TimeSeriesId:    util.Ptr(TimeSeriesIdType(2)),
				SlotCountMin:    util.Ptr(TimeSeriesSlotCountType(5)),
				SlotCountMax:    util.Ptr(TimeSeriesSlotCountType(20)),
				SlotDurationMin: util.Ptr(DurationType("PT5M")),
				SlotDurationMax: util.Ptr(DurationType("PT2H")),
				SlotValueMin:    NewScaledNumberType(10),
				SlotValueMax:    NewScaledNumberType(200),
			},
		},
	}

	// Test 1: No filter (should return all data)
	result, success := data.ReadPartialData(nil)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData := result.(*TimeSeriesConstraintsListDataType)
	assert.Len(t, resultData.TimeSeriesConstraintsData, 2)

	// Test 2: Selector filter
	selectorFilter := &FilterType{
		TimeSeriesConstraintsListDataSelectors: &TimeSeriesConstraintsListDataSelectorsType{
			TimeSeriesId: util.Ptr(TimeSeriesIdType(1)),
		},
	}
	result, success = data.ReadPartialData(selectorFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeSeriesConstraintsListDataType)
	assert.Len(t, resultData.TimeSeriesConstraintsData, 1)
	assert.Equal(t, TimeSeriesIdType(1), *resultData.TimeSeriesConstraintsData[0].TimeSeriesId)

	// Test 3: Element filter
	elementFilter := &FilterType{
		TimeSeriesConstraintsDataElements: &TimeSeriesConstraintsDataElementsType{
			TimeSeriesId: util.Ptr(ElementTagType{}),
			SlotCountMin: util.Ptr(ElementTagType{}),
			SlotCountMax: util.Ptr(ElementTagType{}),
			SlotValueMin: util.Ptr(ElementTagType{}),
			SlotValueMax: util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(elementFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeSeriesConstraintsListDataType)
	assert.Len(t, resultData.TimeSeriesConstraintsData, 2)
	// Check that only requested fields are present
	assert.NotNil(t, resultData.TimeSeriesConstraintsData[0].TimeSeriesId)
	assert.NotNil(t, resultData.TimeSeriesConstraintsData[0].SlotCountMin)
	assert.NotNil(t, resultData.TimeSeriesConstraintsData[0].SlotCountMax)
	assert.NotNil(t, resultData.TimeSeriesConstraintsData[0].SlotValueMin)
	assert.NotNil(t, resultData.TimeSeriesConstraintsData[0].SlotValueMax)
	// Check that non-requested fields are filtered out
	assert.Nil(t, resultData.TimeSeriesConstraintsData[0].SlotDurationMin)
	assert.Nil(t, resultData.TimeSeriesConstraintsData[0].SlotDurationMax)
	assert.Nil(t, resultData.TimeSeriesConstraintsData[0].EarliestTimeSeriesStartTime)
	assert.Nil(t, resultData.TimeSeriesConstraintsData[0].LatestTimeSeriesEndTime)
	assert.Nil(t, resultData.TimeSeriesConstraintsData[0].SlotDurationStepSize)
	assert.Nil(t, resultData.TimeSeriesConstraintsData[0].SlotValueStepSize)

	// Test 4: Combined selector and element filter
	combinedFilter := &FilterType{
		TimeSeriesConstraintsListDataSelectors: &TimeSeriesConstraintsListDataSelectorsType{
			TimeSeriesId: util.Ptr(TimeSeriesIdType(2)),
		},
		TimeSeriesConstraintsDataElements: &TimeSeriesConstraintsDataElementsType{
			TimeSeriesId:    util.Ptr(ElementTagType{}),
			SlotDurationMin: util.Ptr(ElementTagType{}),
			SlotDurationMax: util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(combinedFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeSeriesConstraintsListDataType)
	assert.Len(t, resultData.TimeSeriesConstraintsData, 1)
	assert.Equal(t, TimeSeriesIdType(2), *resultData.TimeSeriesConstraintsData[0].TimeSeriesId)
	assert.NotNil(t, resultData.TimeSeriesConstraintsData[0].SlotDurationMin)
	assert.NotNil(t, resultData.TimeSeriesConstraintsData[0].SlotDurationMax)
	assert.Nil(t, resultData.TimeSeriesConstraintsData[0].SlotCountMin)
	assert.Nil(t, resultData.TimeSeriesConstraintsData[0].SlotCountMax)
	assert.Nil(t, resultData.TimeSeriesConstraintsData[0].SlotValueMin)
	assert.Nil(t, resultData.TimeSeriesConstraintsData[0].SlotValueMax)

	// Test 5: Non-matching selector (should return empty)
	nonMatchingFilter := &FilterType{
		TimeSeriesConstraintsListDataSelectors: &TimeSeriesConstraintsListDataSelectorsType{
			TimeSeriesId: util.Ptr(TimeSeriesIdType(999)),
		},
	}
	result, success = data.ReadPartialData(nonMatchingFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeSeriesConstraintsListDataType)
	assert.Len(t, resultData.TimeSeriesConstraintsData, 0)

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = data.ReadPartialData(emptyFilter)
	assert.False(t, success)
}
