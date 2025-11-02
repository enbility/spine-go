package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestTimeTableListDataType_Update(t *testing.T) {
	sut := TimeTableListDataType{
		TimeTableData: []TimeTableDataType{
			{
				TimeTableId: util.Ptr(TimeTableIdType(0)),
				RecurrenceInformation: &RecurrenceInformationType{
					ExecutionCount: util.Ptr(uint(1)),
				},
			},
			{
				TimeTableId: util.Ptr(TimeTableIdType(1)),
				RecurrenceInformation: &RecurrenceInformationType{
					ExecutionCount: util.Ptr(uint(1)),
				},
			},
		},
	}

	newData := TimeTableListDataType{
		TimeTableData: []TimeTableDataType{
			{
				TimeTableId: util.Ptr(TimeTableIdType(1)),
				RecurrenceInformation: &RecurrenceInformationType{
					ExecutionCount: util.Ptr(uint(10)),
				},
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.TimeTableData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TimeTableId))
	assert.Equal(t, 1, int(*item1.RecurrenceInformation.ExecutionCount))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TimeTableId))
	assert.Equal(t, 10, int(*item2.RecurrenceInformation.ExecutionCount))
}

func TestTimeTableConstraintsListDataType_Update(t *testing.T) {
	sut := TimeTableConstraintsListDataType{
		TimeTableConstraintsData: []TimeTableConstraintsDataType{
			{
				TimeTableId:  util.Ptr(TimeTableIdType(0)),
				SlotCountMin: util.Ptr(TimeSlotCountType(1)),
			},
			{
				TimeTableId:  util.Ptr(TimeTableIdType(1)),
				SlotCountMin: util.Ptr(TimeSlotCountType(1)),
			},
		},
	}

	newData := TimeTableConstraintsListDataType{
		TimeTableConstraintsData: []TimeTableConstraintsDataType{
			{
				TimeTableId:  util.Ptr(TimeTableIdType(1)),
				SlotCountMin: util.Ptr(TimeSlotCountType(10)),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.TimeTableConstraintsData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TimeTableId))
	assert.Equal(t, 1, int(*item1.SlotCountMin))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TimeTableId))
	assert.Equal(t, 10, int(*item2.SlotCountMin))
}

func TestTimeTableDescriptionListDataType_Update(t *testing.T) {
	sut := TimeTableDescriptionListDataType{
		TimeTableDescriptionData: []TimeTableDescriptionDataType{
			{
				TimeTableId: util.Ptr(TimeTableIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				TimeTableId: util.Ptr(TimeTableIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := TimeTableDescriptionListDataType{
		TimeTableDescriptionData: []TimeTableDescriptionDataType{
			{
				TimeTableId: util.Ptr(TimeTableIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.TimeTableDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TimeTableId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TimeTableId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestTimeTableListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	data := &TimeTableListDataType{
		TimeTableData: []TimeTableDataType{
			{
				TimeTableId: util.Ptr(TimeTableIdType(1)),
				TimeSlotId:  util.Ptr(TimeSlotIdType(10)),
				RecurrenceInformation: &RecurrenceInformationType{
					ExecutionCount: util.Ptr(uint(5)),
				},
				StartTime: &AbsoluteOrRecurringTimeType{
					Time: util.Ptr(TimeType("12:00:00")),
				},
				EndTime: &AbsoluteOrRecurringTimeType{
					Time: util.Ptr(TimeType("13:00:00")),
				},
			},
			{
				TimeTableId: util.Ptr(TimeTableIdType(2)),
				TimeSlotId:  util.Ptr(TimeSlotIdType(20)),
				RecurrenceInformation: &RecurrenceInformationType{
					ExecutionCount: util.Ptr(uint(3)),
				},
				StartTime: &AbsoluteOrRecurringTimeType{
					Time: util.Ptr(TimeType("14:00:00")),
				},
				EndTime: &AbsoluteOrRecurringTimeType{
					Time: util.Ptr(TimeType("15:00:00")),
				},
			},
		},
	}

	// Test 1: No filter (should return all data)
	result, success := data.ReadPartialData(nil)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData := result.(*TimeTableListDataType)
	assert.Len(t, resultData.TimeTableData, 2)

	// Test 2: Selector filter by TimeTableId
	selectorFilter := &FilterType{
		TimeTableListDataSelectors: &TimeTableListDataSelectorsType{
			TimeTableId: util.Ptr(TimeTableIdType(1)),
		},
	}
	result, success = data.ReadPartialData(selectorFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeTableListDataType)
	assert.Len(t, resultData.TimeTableData, 1)
	assert.Equal(t, TimeTableIdType(1), *resultData.TimeTableData[0].TimeTableId)

	// Test 3: Selector filter by TimeSlotId
	slotFilter := &FilterType{
		TimeTableListDataSelectors: &TimeTableListDataSelectorsType{
			TimeSlotId: util.Ptr(TimeSlotIdType(20)),
		},
	}
	result, success = data.ReadPartialData(slotFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeTableListDataType)
	assert.Len(t, resultData.TimeTableData, 1)
	assert.Equal(t, TimeTableIdType(2), *resultData.TimeTableData[0].TimeTableId)

	// Test 4: Element filter
	elementFilter := &FilterType{
		TimeTableDataElements: &TimeTableDataElementsType{
			TimeTableId: util.Ptr(ElementTagType{}),
			TimeSlotId:  util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(elementFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeTableListDataType)
	assert.Len(t, resultData.TimeTableData, 2)
	// Check that only requested fields are present
	assert.NotNil(t, resultData.TimeTableData[0].TimeTableId)
	assert.NotNil(t, resultData.TimeTableData[0].TimeSlotId)
	// Check that non-requested fields are filtered out
	assert.Nil(t, resultData.TimeTableData[0].RecurrenceInformation)
	assert.Nil(t, resultData.TimeTableData[0].StartTime)
	assert.Nil(t, resultData.TimeTableData[0].EndTime)

	// Test 5: Combined selector and element filter
	combinedFilter := &FilterType{
		TimeTableListDataSelectors: &TimeTableListDataSelectorsType{
			TimeTableId: util.Ptr(TimeTableIdType(2)),
		},
		TimeTableDataElements: &TimeTableDataElementsType{
			TimeTableId: util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(combinedFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeTableListDataType)
	assert.Len(t, resultData.TimeTableData, 1)
	assert.Equal(t, TimeTableIdType(2), *resultData.TimeTableData[0].TimeTableId)
	assert.Nil(t, resultData.TimeTableData[0].TimeSlotId)
	assert.Nil(t, resultData.TimeTableData[0].RecurrenceInformation)

	// Test 6: Non-matching selector (should return empty)
	nonMatchingFilter := &FilterType{
		TimeTableListDataSelectors: &TimeTableListDataSelectorsType{
			TimeTableId: util.Ptr(TimeTableIdType(999)),
		},
	}
	result, success = data.ReadPartialData(nonMatchingFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeTableListDataType)
	assert.Len(t, resultData.TimeTableData, 0)

	// Test 7: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = data.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestTimeTableConstraintsListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	data := &TimeTableConstraintsListDataType{
		TimeTableConstraintsData: []TimeTableConstraintsDataType{
			{
				TimeTableId:          util.Ptr(TimeTableIdType(1)),
				SlotCountMin:         util.Ptr(TimeSlotCountType(1)),
				SlotCountMax:         util.Ptr(TimeSlotCountType(10)),
				SlotDurationMin:      util.Ptr(DurationType("PT1M")),
				SlotDurationMax:      util.Ptr(DurationType("PT1H")),
				SlotDurationStepSize: util.Ptr(DurationType("PT5M")),
				FirstSlotBeginsAt:    util.Ptr(TimeType("08:00:00")),
			},
			{
				TimeTableId:          util.Ptr(TimeTableIdType(2)),
				SlotCountMin:         util.Ptr(TimeSlotCountType(2)),
				SlotCountMax:         util.Ptr(TimeSlotCountType(20)),
				SlotDurationMin:      util.Ptr(DurationType("PT10M")),
				SlotDurationMax:      util.Ptr(DurationType("PT2H")),
				SlotDurationStepSize: util.Ptr(DurationType("PT10M")),
				FirstSlotBeginsAt:    util.Ptr(TimeType("09:00:00")),
			},
		},
	}

	// Test 1: No filter (should return all data)
	result, success := data.ReadPartialData(nil)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData := result.(*TimeTableConstraintsListDataType)
	assert.Len(t, resultData.TimeTableConstraintsData, 2)

	// Test 2: Selector filter
	selectorFilter := &FilterType{
		TimeTableConstraintsListDataSelectors: &TimeTableConstraintsListDataSelectorsType{
			TimeTableId: util.Ptr(TimeTableIdType(1)),
		},
	}
	result, success = data.ReadPartialData(selectorFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeTableConstraintsListDataType)
	assert.Len(t, resultData.TimeTableConstraintsData, 1)
	assert.Equal(t, TimeTableIdType(1), *resultData.TimeTableConstraintsData[0].TimeTableId)

	// Test 3: Element filter
	elementFilter := &FilterType{
		TimeTableConstraintsDataElements: &TimeTableConstraintsDataElementsType{
			TimeTableId:     util.Ptr(ElementTagType{}),
			SlotCountMin:    util.Ptr(ElementTagType{}),
			SlotCountMax:    util.Ptr(ElementTagType{}),
			SlotDurationMin: util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(elementFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeTableConstraintsListDataType)
	assert.Len(t, resultData.TimeTableConstraintsData, 2)
	// Check that only requested fields are present
	assert.NotNil(t, resultData.TimeTableConstraintsData[0].TimeTableId)
	assert.NotNil(t, resultData.TimeTableConstraintsData[0].SlotCountMin)
	assert.NotNil(t, resultData.TimeTableConstraintsData[0].SlotCountMax)
	assert.NotNil(t, resultData.TimeTableConstraintsData[0].SlotDurationMin)
	// Check that non-requested fields are filtered out
	assert.Nil(t, resultData.TimeTableConstraintsData[0].SlotDurationMax)
	assert.Nil(t, resultData.TimeTableConstraintsData[0].SlotDurationStepSize)
	assert.Nil(t, resultData.TimeTableConstraintsData[0].SlotShiftStepSize)
	assert.Nil(t, resultData.TimeTableConstraintsData[0].FirstSlotBeginsAt)

	// Test 4: Combined selector and element filter
	combinedFilter := &FilterType{
		TimeTableConstraintsListDataSelectors: &TimeTableConstraintsListDataSelectorsType{
			TimeTableId: util.Ptr(TimeTableIdType(2)),
		},
		TimeTableConstraintsDataElements: &TimeTableConstraintsDataElementsType{
			TimeTableId:       util.Ptr(ElementTagType{}),
			FirstSlotBeginsAt: util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(combinedFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeTableConstraintsListDataType)
	assert.Len(t, resultData.TimeTableConstraintsData, 1)
	assert.Equal(t, TimeTableIdType(2), *resultData.TimeTableConstraintsData[0].TimeTableId)
	assert.NotNil(t, resultData.TimeTableConstraintsData[0].FirstSlotBeginsAt)
	assert.Nil(t, resultData.TimeTableConstraintsData[0].SlotCountMin)
	assert.Nil(t, resultData.TimeTableConstraintsData[0].SlotDurationMin)

	// Test 5: Non-matching selector (should return empty)
	nonMatchingFilter := &FilterType{
		TimeTableConstraintsListDataSelectors: &TimeTableConstraintsListDataSelectorsType{
			TimeTableId: util.Ptr(TimeTableIdType(999)),
		},
	}
	result, success = data.ReadPartialData(nonMatchingFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeTableConstraintsListDataType)
	assert.Len(t, resultData.TimeTableConstraintsData, 0)

	// Test 7: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = data.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestTimeTableDescriptionListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	data := &TimeTableDescriptionListDataType{
		TimeTableDescriptionData: []TimeTableDescriptionDataType{
			{
				TimeTableId:             util.Ptr(TimeTableIdType(1)),
				TimeSlotCountChangeable: util.Ptr(true),
				TimeSlotTimesChangeable: util.Ptr(false),
				TimeSlotTimeMode:        util.Ptr(TimeSlotTimeModeType("absolute")),
				Label:                   util.Ptr(LabelType("Table 1")),
				Description:             util.Ptr(DescriptionType("First time table")),
			},
			{
				TimeTableId:             util.Ptr(TimeTableIdType(2)),
				TimeSlotCountChangeable: util.Ptr(false),
				TimeSlotTimesChangeable: util.Ptr(true),
				TimeSlotTimeMode:        util.Ptr(TimeSlotTimeModeType("relative")),
				Label:                   util.Ptr(LabelType("Table 2")),
				Description:             util.Ptr(DescriptionType("Second time table")),
			},
		},
	}

	// Test 1: No filter (should return all data)
	result, success := data.ReadPartialData(nil)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData := result.(*TimeTableDescriptionListDataType)
	assert.Len(t, resultData.TimeTableDescriptionData, 2)

	// Test 2: Selector filter
	selectorFilter := &FilterType{
		TimeTableDescriptionListDataSelectors: &TimeTableDescriptionListDataSelectorsType{
			TimeTableId: util.Ptr(TimeTableIdType(1)),
		},
	}
	result, success = data.ReadPartialData(selectorFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeTableDescriptionListDataType)
	assert.Len(t, resultData.TimeTableDescriptionData, 1)
	assert.Equal(t, TimeTableIdType(1), *resultData.TimeTableDescriptionData[0].TimeTableId)

	// Test 3: Element filter
	elementFilter := &FilterType{
		TimeTableDescriptionDataElements: &TimeTableDescriptionDataElementsType{
			TimeTableId:      util.Ptr(ElementTagType{}),
			TimeSlotTimeMode: util.Ptr(ElementTagType{}),
			Label:            util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(elementFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeTableDescriptionListDataType)
	assert.Len(t, resultData.TimeTableDescriptionData, 2)
	// Check that only requested fields are present
	assert.NotNil(t, resultData.TimeTableDescriptionData[0].TimeTableId)
	assert.NotNil(t, resultData.TimeTableDescriptionData[0].TimeSlotTimeMode)
	assert.NotNil(t, resultData.TimeTableDescriptionData[0].Label)
	// Check that non-requested fields are filtered out
	assert.Nil(t, resultData.TimeTableDescriptionData[0].TimeSlotCountChangeable)
	assert.Nil(t, resultData.TimeTableDescriptionData[0].TimeSlotTimesChangeable)
	assert.Nil(t, resultData.TimeTableDescriptionData[0].Description)

	// Test 4: Combined selector and element filter
	combinedFilter := &FilterType{
		TimeTableDescriptionListDataSelectors: &TimeTableDescriptionListDataSelectorsType{
			TimeTableId: util.Ptr(TimeTableIdType(2)),
		},
		TimeTableDescriptionDataElements: &TimeTableDescriptionDataElementsType{
			TimeTableId: util.Ptr(ElementTagType{}),
			Description: util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(combinedFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeTableDescriptionListDataType)
	assert.Len(t, resultData.TimeTableDescriptionData, 1)
	assert.Equal(t, TimeTableIdType(2), *resultData.TimeTableDescriptionData[0].TimeTableId)
	assert.NotNil(t, resultData.TimeTableDescriptionData[0].Description)
	assert.Nil(t, resultData.TimeTableDescriptionData[0].Label)
	assert.Nil(t, resultData.TimeTableDescriptionData[0].TimeSlotTimeMode)

	// Test 5: Non-matching selector (should return empty)
	nonMatchingFilter := &FilterType{
		TimeTableDescriptionListDataSelectors: &TimeTableDescriptionListDataSelectorsType{
			TimeTableId: util.Ptr(TimeTableIdType(999)),
		},
	}
	result, success = data.ReadPartialData(nonMatchingFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TimeTableDescriptionListDataType)
	assert.Len(t, resultData.TimeTableDescriptionData, 0)

	// Test 7: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = data.ReadPartialData(emptyFilter)
	assert.False(t, success)
}
