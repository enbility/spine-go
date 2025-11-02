package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestLoadControlEventListDataType_Update(t *testing.T) {
	sut := LoadControlEventListDataType{
		LoadControlEventData: []LoadControlEventDataType{
			{
				EventId:            util.Ptr(LoadControlEventIdType(0)),
				EventActionConsume: util.Ptr(LoadControlEventActionTypeNormal),
			},
			{
				EventId:            util.Ptr(LoadControlEventIdType(1)),
				EventActionConsume: util.Ptr(LoadControlEventActionTypeNormal),
			},
		},
	}

	newData := LoadControlEventListDataType{
		LoadControlEventData: []LoadControlEventDataType{
			{
				EventId:            util.Ptr(LoadControlEventIdType(1)),
				EventActionConsume: util.Ptr(LoadControlEventActionTypeIncrease),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.LoadControlEventData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.EventId))
	assert.Equal(t, LoadControlEventActionTypeNormal, *item1.EventActionConsume)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.EventId))
	assert.Equal(t, LoadControlEventActionTypeIncrease, *item2.EventActionConsume)
}

func TestLoadControlStateListDataType_Update(t *testing.T) {
	sut := LoadControlStateListDataType{
		LoadControlStateData: []LoadControlStateDataType{
			{
				EventId:           util.Ptr(LoadControlEventIdType(0)),
				EventStateConsume: util.Ptr(LoadControlEventStateTypeEventAccepted),
			},
			{
				EventId:           util.Ptr(LoadControlEventIdType(1)),
				EventStateConsume: util.Ptr(LoadControlEventStateTypeEventAccepted),
			},
		},
	}

	newData := LoadControlStateListDataType{
		LoadControlStateData: []LoadControlStateDataType{
			{
				EventId:           util.Ptr(LoadControlEventIdType(1)),
				EventStateConsume: util.Ptr(LoadControlEventStateTypeEventStopped),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.LoadControlStateData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.EventId))
	assert.Equal(t, LoadControlEventStateTypeEventAccepted, *item1.EventStateConsume)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.EventId))
	assert.Equal(t, LoadControlEventStateTypeEventStopped, *item2.EventStateConsume)
}

func TestLoadControlLimitListDataType_Update(t *testing.T) {
	sut := LoadControlLimitListDataType{
		LoadControlLimitData: []LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(LoadControlLimitIdType(0)),
				IsLimitChangeable: util.Ptr(false),
			},
			{
				LimitId:           util.Ptr(LoadControlLimitIdType(1)),
				IsLimitChangeable: util.Ptr(true),
			},
		},
	}

	newData := LoadControlLimitListDataType{
		LoadControlLimitData: []LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(LoadControlLimitIdType(1)),
				IsLimitChangeable: util.Ptr(true),
				Value:             NewScaledNumberType(10),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, util.Ptr(FunctionTypeLoadControlLimitListData))
	assert.True(t, success)

	data := sut.LoadControlLimitData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.LimitId))
	assert.False(t, *item1.IsLimitChangeable)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.LimitId))
	assert.True(t, *item2.IsLimitChangeable)
	assert.Equal(t, 10.0, item2.Value.GetValue())
}

func TestLoadControlLimitConstraintsListDataType_Update(t *testing.T) {
	sut := LoadControlLimitConstraintsListDataType{
		LoadControlLimitConstraintsData: []LoadControlLimitConstraintsDataType{
			{
				LimitId:       util.Ptr(LoadControlLimitIdType(0)),
				ValueStepSize: NewScaledNumberType(1),
			},
			{
				LimitId:       util.Ptr(LoadControlLimitIdType(1)),
				ValueStepSize: NewScaledNumberType(1),
			},
		},
	}

	newData := LoadControlLimitConstraintsListDataType{
		LoadControlLimitConstraintsData: []LoadControlLimitConstraintsDataType{
			{
				LimitId:       util.Ptr(LoadControlLimitIdType(1)),
				ValueStepSize: NewScaledNumberType(10),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.LoadControlLimitConstraintsData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.LimitId))
	assert.Equal(t, 1.0, float64(item1.ValueStepSize.GetValue()))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.LimitId))
	assert.Equal(t, 10.0, float64(item2.ValueStepSize.GetValue()))
}

func TestLoadControlLimitDescriptionListDataType_Update(t *testing.T) {
	sut := LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []LoadControlLimitDescriptionDataType{
			{
				LimitId:       util.Ptr(LoadControlLimitIdType(0)),
				LimitCategory: util.Ptr(LoadControlCategoryTypeObligation),
			},
			{
				LimitId:       util.Ptr(LoadControlLimitIdType(1)),
				LimitCategory: util.Ptr(LoadControlCategoryTypeObligation),
			},
		},
	}

	newData := LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []LoadControlLimitDescriptionDataType{
			{
				LimitId:       util.Ptr(LoadControlLimitIdType(1)),
				LimitCategory: util.Ptr(LoadControlCategoryTypeOptimization),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.LoadControlLimitDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.LimitId))
	assert.Equal(t, LoadControlCategoryTypeObligation, *item1.LimitCategory)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.LimitId))
	assert.Equal(t, LoadControlCategoryTypeOptimization, *item2.LimitCategory)
}

func TestLoadControlStateListDataType_ReadPartialData(t *testing.T) {
	eventId := LoadControlEventIdType(123)
	data := &LoadControlStateListDataType{
		LoadControlStateData: []LoadControlStateDataType{
			{EventId: &eventId},
			{EventId: Ptr(LoadControlEventIdType(456))},
		},
	}

	// Test with event ID filter
	filter := &FilterType{
		LoadControlStateListDataSelectors: &LoadControlStateListDataSelectorsType{
			EventId: &eventId,
		},
	}

	result, ok := data.ReadPartialData(filter)
	assert.True(t, ok)

	resultData := result.(*LoadControlStateListDataType)
	assert.Len(t, resultData.LoadControlStateData, 1)
	assert.Equal(t, eventId, *resultData.LoadControlStateData[0].EventId)

	// Test 2: ReadPartialData with nil filter
	_, success := data.ReadPartialData(nil)
	assert.True(t, success)

	// Test 3: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = data.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestLoadControlLimitListDataType_ReadPartialData_Enhanced(t *testing.T) {
	limitId1 := LoadControlLimitIdType(1)
	limitId2 := LoadControlLimitIdType(2)

	testData := &LoadControlLimitListDataType{
		LoadControlLimitData: []LoadControlLimitDataType{
			{
				LimitId:       &limitId1,
				IsLimitActive: util.Ptr(true),
				Value:         &ScaledNumberType{Number: util.Ptr(NumberType(1000)), Scale: util.Ptr(ScaleType(0))},
				TimePeriod: &TimePeriodType{
					StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT1H")),
					EndTime:   util.Ptr(AbsoluteOrRelativeTimeType("PT2H")),
				},
			},
			{
				LimitId:       &limitId2,
				IsLimitActive: util.Ptr(false),
				Value:         &ScaledNumberType{Number: util.Ptr(NumberType(2000)), Scale: util.Ptr(ScaleType(0))},
				TimePeriod: &TimePeriodType{
					StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT2H")),
					EndTime:   util.Ptr(AbsoluteOrRelativeTimeType("PT3H")),
				},
			},
		},
	}

	// Test 1: Filter by limit ID
	filter := &FilterType{
		LoadControlLimitListDataSelectors: &LoadControlLimitListDataSelectorsType{
			LimitId: &limitId1,
		},
	}

	result, success := testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData := result.(*LoadControlLimitListDataType)
	assert.Len(t, resultData.LoadControlLimitData, 1)
	assert.Equal(t, limitId1, *resultData.LoadControlLimitData[0].LimitId)
	assert.True(t, *resultData.LoadControlLimitData[0].IsLimitActive)

	// Test 2: No matching filter
	nonExistentId := LoadControlLimitIdType(999)
	filter = &FilterType{
		LoadControlLimitListDataSelectors: &LoadControlLimitListDataSelectorsType{
			LimitId: &nonExistentId,
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*LoadControlLimitListDataType)
	assert.Len(t, resultData.LoadControlLimitData, 0)

	// Test 3: UpdateList functionality
	updatedLimit := LoadControlLimitDataType{
		LimitId:       &limitId1,
		IsLimitActive: util.Ptr(false),
		Value:         &ScaledNumberType{Number: util.Ptr(NumberType(500)), Scale: util.Ptr(ScaleType(0))},
	}

	newData := &LoadControlLimitListDataType{
		LoadControlLimitData: []LoadControlLimitDataType{updatedLimit},
	}

	updatedList, success := testData.UpdateList(false, true, newData, nil, nil, Ptr(FunctionTypeLoadControlLimitListData))
	assert.True(t, success)
	updatedListData := updatedList.([]LoadControlLimitDataType)
	assert.Len(t, updatedListData, 2) // Should have both items

	// Find the updated item
	var foundUpdatedItem *LoadControlLimitDataType
	for _, item := range updatedListData {
		if *item.LimitId == limitId1 {
			foundUpdatedItem = &item
			break
		}
	}

	assert.NotNil(t, foundUpdatedItem)
	assert.False(t, *foundUpdatedItem.IsLimitActive)                 // Should be updated to false
	assert.Equal(t, float64(500), foundUpdatedItem.Value.GetValue()) // Should be updated to 500

	// Test 4: ReadPartialData with nil filter
	_, success = testData.ReadPartialData(nil)
	assert.True(t, success)

	// Test 5: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestLoadControlEventListDataType_ReadPartialData(t *testing.T) {
	eventId1 := LoadControlEventIdType(100)
	eventId2 := LoadControlEventIdType(200)

	testData := &LoadControlEventListDataType{
		LoadControlEventData: []LoadControlEventDataType{
			{
				EventId:            &eventId1,
				EventActionConsume: util.Ptr(LoadControlEventActionType("start")),
				EventActionProduce: util.Ptr(LoadControlEventActionType("pause")),
				TimePeriod: &TimePeriodType{
					StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT0S")),
					EndTime:   util.Ptr(AbsoluteOrRelativeTimeType("PT1H")),
				},
			},
			{
				EventId:            &eventId2,
				EventActionConsume: util.Ptr(LoadControlEventActionType("end")),
				EventActionProduce: util.Ptr(LoadControlEventActionType("stop")),
				TimePeriod: &TimePeriodType{
					StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT1H")),
					EndTime:   util.Ptr(AbsoluteOrRelativeTimeType("PT2H")),
				},
			},
		},
	}

	// Test 1: Filter by event ID
	filter := &FilterType{
		LoadControlEventListDataSelectors: &LoadControlEventListDataSelectorsType{
			EventId: &eventId1,
		},
	}

	result, success := testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData := result.(*LoadControlEventListDataType)
	assert.Len(t, resultData.LoadControlEventData, 1)
	assert.Equal(t, eventId1, *resultData.LoadControlEventData[0].EventId)
	assert.Equal(t, "start", string(*resultData.LoadControlEventData[0].EventActionConsume))

	// Test 2: No filter should return all data
	result, success = testData.ReadPartialData(nil)
	assert.True(t, success)
	resultData = result.(*LoadControlEventListDataType)
	assert.Len(t, resultData.LoadControlEventData, 2)

	// Test 3: UpdateList functionality
	newEvent := LoadControlEventDataType{
		EventId:            &eventId1,
		EventActionConsume: util.Ptr(LoadControlEventActionType("resume")),
	}

	newData := &LoadControlEventListDataType{
		LoadControlEventData: []LoadControlEventDataType{newEvent},
	}

	updatedList, success := testData.UpdateList(false, true, newData, nil, nil, Ptr(FunctionTypeLoadControlEventListData))
	assert.True(t, success)
	updatedListData := updatedList.([]LoadControlEventDataType)

	// Find the updated item
	var foundUpdatedItem *LoadControlEventDataType
	for _, item := range updatedListData {
		if *item.EventId == eventId1 {
			foundUpdatedItem = &item
			break
		}
	}

	assert.NotNil(t, foundUpdatedItem)
	assert.Equal(t, "resume", string(*foundUpdatedItem.EventActionConsume))

	// Test 4: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestLoadControlLimitConstraintsListDataType_ReadPartialData(t *testing.T) {
	limitId1 := LoadControlLimitIdType(1)
	limitId2 := LoadControlLimitIdType(2)

	testData := &LoadControlLimitConstraintsListDataType{
		LoadControlLimitConstraintsData: []LoadControlLimitConstraintsDataType{
			{
				LimitId: &limitId1,
			},
			{
				LimitId: &limitId2,
			},
		},
	}

	// Test 1: ReadPartialData with nil filter should return all data
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, testData, result)

	// Test 2: ReadPartialData with filter
	filter := &FilterType{
		LoadControlLimitConstraintsListDataSelectors: &LoadControlLimitConstraintsListDataSelectorsType{
			LimitId: &limitId1,
		},
	}
	_, success = testData.ReadPartialData(filter)
	assert.True(t, success)

	// Test 3: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestLoadControlLimitDescriptionListDataType_ReadPartialData(t *testing.T) {
	limitId1 := LoadControlLimitIdType(1)
	limitId2 := LoadControlLimitIdType(2)

	testData := &LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []LoadControlLimitDescriptionDataType{
			{
				LimitId: &limitId1,
			},
			{
				LimitId: &limitId2,
			},
		},
	}

	// Test 1: ReadPartialData with nil filter should return all data
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, testData, result)

	// Test 2: ReadPartialData with filter
	filter := &FilterType{
		LoadControlLimitDescriptionListDataSelectors: &LoadControlLimitDescriptionListDataSelectorsType{
			LimitId: &limitId1,
		},
	}
	_, success = testData.ReadPartialData(filter)
	assert.True(t, success)

	// Test 3: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}
