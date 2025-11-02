package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestTaskManagementJobListDataType_Update(t *testing.T) {
	sut := TaskManagementJobListDataType{
		TaskManagementJobData: []TaskManagementJobDataType{
			{
				JobId:    util.Ptr(TaskManagementJobIdType(0)),
				JobState: util.Ptr(TaskManagementJobStateTypeActive),
			},
			{
				JobId:    util.Ptr(TaskManagementJobIdType(1)),
				JobState: util.Ptr(TaskManagementJobStateTypeActive),
			},
		},
	}

	newData := TaskManagementJobListDataType{
		TaskManagementJobData: []TaskManagementJobDataType{
			{
				JobId:    util.Ptr(TaskManagementJobIdType(1)),
				JobState: util.Ptr(TaskManagementJobStateTypeCompleted),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.TaskManagementJobData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.JobId))
	assert.Equal(t, TaskManagementJobStateTypeActive, *item1.JobState)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.JobId))
	assert.Equal(t, TaskManagementJobStateTypeCompleted, *item2.JobState)
}

func TestTaskManagementJobRelationListDataType_Update(t *testing.T) {
	sut := TaskManagementJobRelationListDataType{
		TaskManagementJobRelationData: []TaskManagementJobRelationDataType{
			{
				JobId: util.Ptr(TaskManagementJobIdType(0)),
				LoadControlReleated: &TaskManagementLoadControlReleatedType{
					EventId: util.Ptr(LoadControlEventIdType(0)),
				},
			},
			{
				JobId: util.Ptr(TaskManagementJobIdType(1)),
				LoadControlReleated: &TaskManagementLoadControlReleatedType{
					EventId: util.Ptr(LoadControlEventIdType(0)),
				},
			},
		},
	}

	newData := TaskManagementJobRelationListDataType{
		TaskManagementJobRelationData: []TaskManagementJobRelationDataType{
			{
				JobId: util.Ptr(TaskManagementJobIdType(1)),
				LoadControlReleated: &TaskManagementLoadControlReleatedType{
					EventId: util.Ptr(LoadControlEventIdType(1)),
				},
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.TaskManagementJobRelationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.JobId))
	assert.Equal(t, 0, int(*item1.LoadControlReleated.EventId))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.JobId))
	assert.Equal(t, 1, int(*item2.LoadControlReleated.EventId))
}

func TestTaskManagementJobDescriptionListDataType_Update(t *testing.T) {
	sut := TaskManagementJobDescriptionListDataType{
		TaskManagementJobDescriptionData: []TaskManagementJobDescriptionDataType{
			{
				JobId:       util.Ptr(TaskManagementJobIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				JobId:       util.Ptr(TaskManagementJobIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := TaskManagementJobDescriptionListDataType{
		TaskManagementJobDescriptionData: []TaskManagementJobDescriptionDataType{
			{
				JobId:       util.Ptr(TaskManagementJobIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.TaskManagementJobDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.JobId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.JobId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestTaskManagementJobListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	data := &TaskManagementJobListDataType{
		TaskManagementJobData: []TaskManagementJobDataType{
			{
				JobId:         util.Ptr(TaskManagementJobIdType(1)),
				Timestamp:     util.Ptr(AbsoluteOrRelativeTimeType("2023-01-01T10:00:00Z")),
				JobState:      util.Ptr(TaskManagementJobStateTypeActive),
				ElapsedTime:   util.Ptr(DurationType("PT1H")),
				RemainingTime: util.Ptr(DurationType("PT2H")),
			},
			{
				JobId:         util.Ptr(TaskManagementJobIdType(2)),
				Timestamp:     util.Ptr(AbsoluteOrRelativeTimeType("2023-01-01T11:00:00Z")),
				JobState:      util.Ptr(TaskManagementJobStateTypeCompleted),
				ElapsedTime:   util.Ptr(DurationType("PT30M")),
				RemainingTime: util.Ptr(DurationType("PT0S")),
			},
		},
	}

	// Test 1: No filter (should return all data)
	result, success := data.ReadPartialData(nil)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData := result.(*TaskManagementJobListDataType)
	assert.Len(t, resultData.TaskManagementJobData, 2)

	// Test 2: Selector filter by JobId
	selectorFilter := &FilterType{
		TaskManagementJobListDataSelectors: &TaskManagementJobListDataSelectorsType{
			JobId: util.Ptr(TaskManagementJobIdType(1)),
		},
	}
	result, success = data.ReadPartialData(selectorFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TaskManagementJobListDataType)
	assert.Len(t, resultData.TaskManagementJobData, 1)
	assert.Equal(t, TaskManagementJobIdType(1), *resultData.TaskManagementJobData[0].JobId)

	// Test 3: Selector filter by JobState
	stateFilter := &FilterType{
		TaskManagementJobListDataSelectors: &TaskManagementJobListDataSelectorsType{
			JobState: util.Ptr(TaskManagementJobStateTypeCompleted),
		},
	}
	result, success = data.ReadPartialData(stateFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TaskManagementJobListDataType)
	assert.Len(t, resultData.TaskManagementJobData, 1)
	assert.Equal(t, TaskManagementJobIdType(2), *resultData.TaskManagementJobData[0].JobId)

	// Test 4: Element filter
	elementFilter := &FilterType{
		TaskManagementJobDataElements: &TaskManagementJobDataElementsType{
			JobId:     util.Ptr(ElementTagType{}),
			JobState:  util.Ptr(ElementTagType{}),
			Timestamp: util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(elementFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TaskManagementJobListDataType)
	assert.Len(t, resultData.TaskManagementJobData, 2)
	// Check that only requested fields are present
	assert.NotNil(t, resultData.TaskManagementJobData[0].JobId)
	assert.NotNil(t, resultData.TaskManagementJobData[0].JobState)
	assert.NotNil(t, resultData.TaskManagementJobData[0].Timestamp)
	// Check that non-requested fields are filtered out
	assert.Nil(t, resultData.TaskManagementJobData[0].ElapsedTime)
	assert.Nil(t, resultData.TaskManagementJobData[0].RemainingTime)

	// Test 5: Combined selector and element filter
	combinedFilter := &FilterType{
		TaskManagementJobListDataSelectors: &TaskManagementJobListDataSelectorsType{
			JobState: util.Ptr(TaskManagementJobStateTypeActive),
		},
		TaskManagementJobDataElements: &TaskManagementJobDataElementsType{
			JobId:         util.Ptr(ElementTagType{}),
			RemainingTime: util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(combinedFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TaskManagementJobListDataType)
	assert.Len(t, resultData.TaskManagementJobData, 1)
	assert.Equal(t, TaskManagementJobIdType(1), *resultData.TaskManagementJobData[0].JobId)
	assert.NotNil(t, resultData.TaskManagementJobData[0].RemainingTime)
	assert.Nil(t, resultData.TaskManagementJobData[0].JobState)
	assert.Nil(t, resultData.TaskManagementJobData[0].ElapsedTime)

	// Test 6: Non-matching selector (should return empty)
	nonMatchingFilter := &FilterType{
		TaskManagementJobListDataSelectors: &TaskManagementJobListDataSelectorsType{
			JobId: util.Ptr(TaskManagementJobIdType(999)),
		},
	}
	result, success = data.ReadPartialData(nonMatchingFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TaskManagementJobListDataType)
	assert.Len(t, resultData.TaskManagementJobData, 0)

	// Test 7: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = data.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestTaskManagementJobRelationListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	data := &TaskManagementJobRelationListDataType{
		TaskManagementJobRelationData: []TaskManagementJobRelationDataType{
			{
				JobId:                util.Ptr(TaskManagementJobIdType(1)),
				DirectControlRelated: &TaskManagementDirectControlRelatedType{},
				LoadControlReleated: &TaskManagementLoadControlReleatedType{
					EventId: util.Ptr(LoadControlEventIdType(10)),
				},
			},
			{
				JobId: util.Ptr(TaskManagementJobIdType(2)),
				HvacRelated: &TaskManagementHvacRelatedType{
					OverrunId: util.Ptr(HvacOverrunIdType(1)),
				},
				PowerSequencesRelated: &TaskManagementPowerSequencesRelatedType{
					SequenceId: util.Ptr(PowerSequenceIdType(5)),
				},
			},
		},
	}

	// Test 1: No filter (should return all data)
	result, success := data.ReadPartialData(nil)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData := result.(*TaskManagementJobRelationListDataType)
	assert.Len(t, resultData.TaskManagementJobRelationData, 2)

	// Test 2: Selector filter
	selectorFilter := &FilterType{
		TaskManagementJobRelationListDataSelectors: &TaskManagementJobRelationListDataSelectorsType{
			JobId: util.Ptr(TaskManagementJobIdType(1)),
		},
	}
	result, success = data.ReadPartialData(selectorFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TaskManagementJobRelationListDataType)
	assert.Len(t, resultData.TaskManagementJobRelationData, 1)
	assert.Equal(t, TaskManagementJobIdType(1), *resultData.TaskManagementJobRelationData[0].JobId)

	// Test 3: Element filter
	elementFilter := &FilterType{
		TaskManagementJobRelationDataElements: &TaskManagementJobRelationDataElementsType{
			JobId:       util.Ptr(ElementTagType{}),
			HvacRelated: &TaskManagementHvacRelatedElementsType{},
		},
	}
	result, success = data.ReadPartialData(elementFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TaskManagementJobRelationListDataType)
	assert.Len(t, resultData.TaskManagementJobRelationData, 2)
	// Check that only requested fields are present
	assert.NotNil(t, resultData.TaskManagementJobRelationData[0].JobId)
	// Check that non-requested fields are filtered out
	assert.Nil(t, resultData.TaskManagementJobRelationData[0].DirectControlRelated)
	assert.Nil(t, resultData.TaskManagementJobRelationData[0].LoadControlReleated)
	assert.Nil(t, resultData.TaskManagementJobRelationData[0].PowerSequencesRelated)

	// Test 4: Combined selector and element filter
	combinedFilter := &FilterType{
		TaskManagementJobRelationListDataSelectors: &TaskManagementJobRelationListDataSelectorsType{
			JobId: util.Ptr(TaskManagementJobIdType(2)),
		},
		TaskManagementJobRelationDataElements: &TaskManagementJobRelationDataElementsType{
			JobId:                 util.Ptr(ElementTagType{}),
			PowerSequencesRelated: &TaskManagementPowerSequencesRelatedElementsType{},
		},
	}
	result, success = data.ReadPartialData(combinedFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TaskManagementJobRelationListDataType)
	assert.Len(t, resultData.TaskManagementJobRelationData, 1)
	assert.Equal(t, TaskManagementJobIdType(2), *resultData.TaskManagementJobRelationData[0].JobId)
	assert.NotNil(t, resultData.TaskManagementJobRelationData[0].PowerSequencesRelated)
	assert.Nil(t, resultData.TaskManagementJobRelationData[0].HvacRelated)
	assert.Nil(t, resultData.TaskManagementJobRelationData[0].DirectControlRelated)

	// Test 5: Non-matching selector (should return empty)
	nonMatchingFilter := &FilterType{
		TaskManagementJobRelationListDataSelectors: &TaskManagementJobRelationListDataSelectorsType{
			JobId: util.Ptr(TaskManagementJobIdType(999)),
		},
	}
	result, success = data.ReadPartialData(nonMatchingFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TaskManagementJobRelationListDataType)
	assert.Len(t, resultData.TaskManagementJobRelationData, 0)

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = data.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestTaskManagementJobDescriptionListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	data := &TaskManagementJobDescriptionListDataType{
		TaskManagementJobDescriptionData: []TaskManagementJobDescriptionDataType{
			{
				JobId:       util.Ptr(TaskManagementJobIdType(1)),
				JobSource:   util.Ptr(TaskManagementJobSourceType("user")),
				Label:       util.Ptr(LabelType("Job 1")),
				Description: util.Ptr(DescriptionType("First job")),
			},
			{
				JobId:       util.Ptr(TaskManagementJobIdType(2)),
				JobSource:   util.Ptr(TaskManagementJobSourceType("system")),
				Label:       util.Ptr(LabelType("Job 2")),
				Description: util.Ptr(DescriptionType("Second job")),
			},
		},
	}

	// Test 1: No filter (should return all data)
	result, success := data.ReadPartialData(nil)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData := result.(*TaskManagementJobDescriptionListDataType)
	assert.Len(t, resultData.TaskManagementJobDescriptionData, 2)

	// Test 2: Selector filter by JobId
	selectorFilter := &FilterType{
		TaskManagementJobDescriptionListDataSelectors: &TaskManagementJobDescriptionListDataSelectorsType{
			JobId: util.Ptr(TaskManagementJobIdType(1)),
		},
	}
	result, success = data.ReadPartialData(selectorFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TaskManagementJobDescriptionListDataType)
	assert.Len(t, resultData.TaskManagementJobDescriptionData, 1)
	assert.Equal(t, TaskManagementJobIdType(1), *resultData.TaskManagementJobDescriptionData[0].JobId)

	// Test 3: Selector filter by JobSource
	sourceFilter := &FilterType{
		TaskManagementJobDescriptionListDataSelectors: &TaskManagementJobDescriptionListDataSelectorsType{
			JobSource: util.Ptr(TaskManagementJobSourceType("system")),
		},
	}
	result, success = data.ReadPartialData(sourceFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TaskManagementJobDescriptionListDataType)
	assert.Len(t, resultData.TaskManagementJobDescriptionData, 1)
	assert.Equal(t, TaskManagementJobIdType(2), *resultData.TaskManagementJobDescriptionData[0].JobId)

	// Test 4: Element filter
	elementFilter := &FilterType{
		TaskManagementJobDescriptionDataElements: &TaskManagementJobDescriptionDataElementsType{
			JobId:     util.Ptr(ElementTagType{}),
			JobSource: util.Ptr(ElementTagType{}),
			Label:     util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(elementFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TaskManagementJobDescriptionListDataType)
	assert.Len(t, resultData.TaskManagementJobDescriptionData, 2)
	// Check that only requested fields are present
	assert.NotNil(t, resultData.TaskManagementJobDescriptionData[0].JobId)
	assert.NotNil(t, resultData.TaskManagementJobDescriptionData[0].JobSource)
	assert.NotNil(t, resultData.TaskManagementJobDescriptionData[0].Label)
	// Check that non-requested field is filtered out
	assert.Nil(t, resultData.TaskManagementJobDescriptionData[0].Description)

	// Test 5: Combined selector and element filter
	combinedFilter := &FilterType{
		TaskManagementJobDescriptionListDataSelectors: &TaskManagementJobDescriptionListDataSelectorsType{
			JobSource: util.Ptr(TaskManagementJobSourceType("user")),
		},
		TaskManagementJobDescriptionDataElements: &TaskManagementJobDescriptionDataElementsType{
			JobId:       util.Ptr(ElementTagType{}),
			Description: util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(combinedFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TaskManagementJobDescriptionListDataType)
	assert.Len(t, resultData.TaskManagementJobDescriptionData, 1)
	assert.Equal(t, TaskManagementJobIdType(1), *resultData.TaskManagementJobDescriptionData[0].JobId)
	assert.NotNil(t, resultData.TaskManagementJobDescriptionData[0].Description)
	assert.Nil(t, resultData.TaskManagementJobDescriptionData[0].JobSource)
	assert.Nil(t, resultData.TaskManagementJobDescriptionData[0].Label)

	// Test 6: Non-matching selector (should return empty)
	nonMatchingFilter := &FilterType{
		TaskManagementJobDescriptionListDataSelectors: &TaskManagementJobDescriptionListDataSelectorsType{
			JobId: util.Ptr(TaskManagementJobIdType(999)),
		},
	}
	result, success = data.ReadPartialData(nonMatchingFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*TaskManagementJobDescriptionListDataType)
	assert.Len(t, resultData.TaskManagementJobDescriptionData, 0)

	// Test 7: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = data.ReadPartialData(emptyFilter)
	assert.False(t, success)
}
