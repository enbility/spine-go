package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestHvacSystemFunctionListDataType_Update(t *testing.T) {
	sut := HvacSystemFunctionListDataType{
		HvacSystemFunctionData: []HvacSystemFunctionDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(0)),
				IsOverrunActive:  util.Ptr(false),
			},
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				IsOverrunActive:  util.Ptr(false),
			},
		},
	}

	newData := HvacSystemFunctionListDataType{
		HvacSystemFunctionData: []HvacSystemFunctionDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				IsOverrunActive:  util.Ptr(true),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.HvacSystemFunctionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SystemFunctionId))
	assert.Equal(t, false, *item1.IsOverrunActive)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SystemFunctionId))
	assert.Equal(t, true, *item2.IsOverrunActive)
}

func TestHvacSystemFunctionOperationModeRelationListDataType_Update(t *testing.T) {
	sut := HvacSystemFunctionOperationModeRelationListDataType{
		HvacSystemFunctionOperationModeRelationData: []HvacSystemFunctionOperationModeRelationDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(0)),
				OperationModeId:  []HvacOperationModeIdType{0},
			},
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				OperationModeId:  []HvacOperationModeIdType{0},
			},
		},
	}

	newData := HvacSystemFunctionOperationModeRelationListDataType{
		HvacSystemFunctionOperationModeRelationData: []HvacSystemFunctionOperationModeRelationDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				OperationModeId:  []HvacOperationModeIdType{1},
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.HvacSystemFunctionOperationModeRelationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SystemFunctionId))
	assert.Equal(t, 0, int(item1.OperationModeId[0]))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SystemFunctionId))
	assert.Equal(t, 1, int(item2.OperationModeId[0]))
}

func TestHvacSystemFunctionSetpointRelationListDataType_Update(t *testing.T) {
	sut := HvacSystemFunctionSetpointRelationListDataType{
		HvacSystemFunctionSetpointRelationData: []HvacSystemFunctionSetpointRelationDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(0)),
				OperationModeId:  util.Ptr(HvacOperationModeIdType(0)),
			},
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				OperationModeId:  util.Ptr(HvacOperationModeIdType(0)),
			},
		},
	}

	newData := HvacSystemFunctionSetpointRelationListDataType{
		HvacSystemFunctionSetpointRelationData: []HvacSystemFunctionSetpointRelationDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				OperationModeId:  util.Ptr(HvacOperationModeIdType(1)),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.HvacSystemFunctionSetpointRelationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SystemFunctionId))
	assert.Equal(t, 0, int(*item1.OperationModeId))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SystemFunctionId))
	assert.Equal(t, 1, int(*item2.OperationModeId))
}

func TestHvacSystemFunctionPowerSequenceRelationListDataType_Update(t *testing.T) {
	sut := HvacSystemFunctionPowerSequenceRelationListDataType{
		HvacSystemFunctionPowerSequenceRelationData: []HvacSystemFunctionPowerSequenceRelationDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(0)),
				SequenceId:       []PowerSequenceIdType{0},
			},
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				SequenceId:       []PowerSequenceIdType{0},
			},
		},
	}

	newData := HvacSystemFunctionPowerSequenceRelationListDataType{
		HvacSystemFunctionPowerSequenceRelationData: []HvacSystemFunctionPowerSequenceRelationDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				SequenceId:       []PowerSequenceIdType{1},
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.HvacSystemFunctionPowerSequenceRelationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SystemFunctionId))
	assert.Equal(t, 0, int(item1.SequenceId[0]))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SystemFunctionId))
	assert.Equal(t, 1, int(item2.SequenceId[0]))
}

func TestHvacSystemFunctionDescriptionListDataType_Update(t *testing.T) {
	sut := HvacSystemFunctionDescriptionListDataType{
		HvacSystemFunctionDescriptionData: []HvacSystemFunctionDescriptionDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(0)),
				Description:      util.Ptr(DescriptionType("old")),
			},
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				Description:      util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := HvacSystemFunctionDescriptionListDataType{
		HvacSystemFunctionDescriptionData: []HvacSystemFunctionDescriptionDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				Description:      util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.HvacSystemFunctionDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SystemFunctionId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SystemFunctionId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestHvacOperationModeDescriptionListDataType_Update(t *testing.T) {
	sut := HvacOperationModeDescriptionListDataType{
		HvacOperationModeDescriptionData: []HvacOperationModeDescriptionDataType{
			{
				OperationModeId: util.Ptr(HvacOperationModeIdType(0)),
				Description:     util.Ptr(DescriptionType("old")),
			},
			{
				OperationModeId: util.Ptr(HvacOperationModeIdType(1)),
				Description:     util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := HvacOperationModeDescriptionListDataType{
		HvacOperationModeDescriptionData: []HvacOperationModeDescriptionDataType{
			{
				OperationModeId: util.Ptr(HvacOperationModeIdType(1)),
				Description:     util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.HvacOperationModeDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.OperationModeId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.OperationModeId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestHvacOverrunListDataType_Update(t *testing.T) {
	sut := HvacOverrunListDataType{
		HvacOverrunData: []HvacOverrunDataType{
			{
				OverrunId:                 util.Ptr(HvacOverrunIdType(0)),
				IsOverrunStatusChangeable: util.Ptr(false),
			},
			{
				OverrunId:                 util.Ptr(HvacOverrunIdType(1)),
				IsOverrunStatusChangeable: util.Ptr(false),
			},
		},
	}

	newData := HvacOverrunListDataType{
		HvacOverrunData: []HvacOverrunDataType{
			{
				OverrunId:                 util.Ptr(HvacOverrunIdType(1)),
				IsOverrunStatusChangeable: util.Ptr(true),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.HvacOverrunData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.OverrunId))
	assert.Equal(t, false, *item1.IsOverrunStatusChangeable)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.OverrunId))
	assert.Equal(t, true, *item2.IsOverrunStatusChangeable)
}

func TestHvacOverrunDescriptionListDataType_Update(t *testing.T) {
	sut := HvacOverrunDescriptionListDataType{
		HvacOverrunDescriptionData: []HvacOverrunDescriptionDataType{
			{
				OverrunId:   util.Ptr(HvacOverrunIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				OverrunId:   util.Ptr(HvacOverrunIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := HvacOverrunDescriptionListDataType{
		HvacOverrunDescriptionData: []HvacOverrunDescriptionDataType{
			{
				OverrunId:   util.Ptr(HvacOverrunIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.HvacOverrunDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.OverrunId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.OverrunId))
	assert.Equal(t, "new", string(*item2.Description))
}

// ReadPartialData tests

func TestHvacOverrunDescriptionListDataType_ReadPartialData(t *testing.T) {
	sut := HvacOverrunDescriptionListDataType{
		HvacOverrunDescriptionData: []HvacOverrunDescriptionDataType{
			{
				OverrunId:                util.Ptr(HvacOverrunIdType(1)),
				OverrunType:              util.Ptr(HvacOverrunTypeTypeParty),
				AffectedSystemFunctionId: []HvacSystemFunctionIdType{10, 20},
				Label:                    util.Ptr(LabelType("Test Label 1")),
				Description:              util.Ptr(DescriptionType("Test Description 1")),
			},
			{
				OverrunId:                util.Ptr(HvacOverrunIdType(2)),
				OverrunType:              util.Ptr(HvacOverrunTypeTypeOneTimeDhw),
				AffectedSystemFunctionId: []HvacSystemFunctionIdType{30},
				Label:                    util.Ptr(LabelType("Test Label 2")),
				Description:              util.Ptr(DescriptionType("Test Description 2")),
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*HvacOverrunDescriptionListDataType)
	assert.Equal(t, 2, len(resultData.HvacOverrunDescriptionData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		HvacOverrunDescriptionListDataSelectors: &HvacOverrunDescriptionListDataSelectorsType{
			OverrunId: util.Ptr(HvacOverrunIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacOverrunDescriptionListDataType)
	assert.Equal(t, 1, len(resultData.HvacOverrunDescriptionData))
	assert.Equal(t, HvacOverrunIdType(1), *resultData.HvacOverrunDescriptionData[0].OverrunId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		HvacOverrunDescriptionDataElements: &HvacOverrunDescriptionDataElementsType{
			OverrunId:   &ElementTagType{},
			Description: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacOverrunDescriptionListDataType)
	assert.Equal(t, 2, len(resultData.HvacOverrunDescriptionData))
	assert.NotNil(t, resultData.HvacOverrunDescriptionData[0].OverrunId)
	assert.NotNil(t, resultData.HvacOverrunDescriptionData[0].Description)
	assert.Nil(t, resultData.HvacOverrunDescriptionData[0].OverrunType)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		HvacOverrunDescriptionListDataSelectors: &HvacOverrunDescriptionListDataSelectorsType{
			OverrunId: util.Ptr(HvacOverrunIdType(2)),
		},
		HvacOverrunDescriptionDataElements: &HvacOverrunDescriptionDataElementsType{
			OverrunId: &ElementTagType{},
			Label:     &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacOverrunDescriptionListDataType)
	assert.Equal(t, 1, len(resultData.HvacOverrunDescriptionData))
	assert.Equal(t, HvacOverrunIdType(2), *resultData.HvacOverrunDescriptionData[0].OverrunId)
	assert.NotNil(t, resultData.HvacOverrunDescriptionData[0].Label)
	assert.Nil(t, resultData.HvacOverrunDescriptionData[0].Description)

	// Test 5: empty filter should return original data
	result, success = sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData = result.(*HvacOverrunDescriptionListDataType)
	assert.Equal(t, 2, len(resultData.HvacOverrunDescriptionData))

	// Test 6: non-matching selector should return empty data
	filter = FilterType{
		HvacOverrunDescriptionListDataSelectors: &HvacOverrunDescriptionListDataSelectorsType{
			OverrunId: util.Ptr(HvacOverrunIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacOverrunDescriptionListDataType)
	assert.Equal(t, 0, len(resultData.HvacOverrunDescriptionData))

	// Test 7: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestHvacSystemFunctionListDataType_ReadPartialData(t *testing.T) {
	sut := &HvacSystemFunctionListDataType{
		HvacSystemFunctionData: []HvacSystemFunctionDataType{
			{
				SystemFunctionId:            util.Ptr(HvacSystemFunctionIdType(1)),
				CurrentOperationModeId:      util.Ptr(HvacOperationModeIdType(10)),
				IsOperationModeIdChangeable: util.Ptr(true),
				CurrentSetpointId:           util.Ptr(SetpointIdType(100)),
				IsSetpointIdChangeable:      util.Ptr(false),
				IsOverrunActive:             util.Ptr(true),
			},
			{
				SystemFunctionId:            util.Ptr(HvacSystemFunctionIdType(2)),
				CurrentOperationModeId:      util.Ptr(HvacOperationModeIdType(20)),
				IsOperationModeIdChangeable: util.Ptr(false),
				CurrentSetpointId:           util.Ptr(SetpointIdType(200)),
				IsSetpointIdChangeable:      util.Ptr(true),
				IsOverrunActive:             util.Ptr(false),
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*HvacSystemFunctionListDataType)
	assert.Equal(t, 2, len(resultData.HvacSystemFunctionData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		HvacSystemFunctionListDataSelectors: &HvacSystemFunctionListDataSelectorsType{
			SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacSystemFunctionListDataType)
	assert.Equal(t, 1, len(resultData.HvacSystemFunctionData))
	assert.Equal(t, HvacSystemFunctionIdType(1), *resultData.HvacSystemFunctionData[0].SystemFunctionId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		HvacSystemFunctionDataElements: &HvacSystemFunctionDataElementsType{
			SystemFunctionId: &ElementTagType{},
			IsOverrunActive:  &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacSystemFunctionListDataType)
	assert.Equal(t, 2, len(resultData.HvacSystemFunctionData))
	assert.NotNil(t, resultData.HvacSystemFunctionData[0].SystemFunctionId)
	assert.NotNil(t, resultData.HvacSystemFunctionData[0].IsOverrunActive)
	assert.Nil(t, resultData.HvacSystemFunctionData[0].CurrentOperationModeId)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		HvacSystemFunctionListDataSelectors: &HvacSystemFunctionListDataSelectorsType{
			SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(2)),
		},
		HvacSystemFunctionDataElements: &HvacSystemFunctionDataElementsType{
			SystemFunctionId:  &ElementTagType{},
			CurrentSetpointId: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacSystemFunctionListDataType)
	assert.Equal(t, 1, len(resultData.HvacSystemFunctionData))
	assert.Equal(t, HvacSystemFunctionIdType(2), *resultData.HvacSystemFunctionData[0].SystemFunctionId)
	assert.NotNil(t, resultData.HvacSystemFunctionData[0].CurrentSetpointId)
	assert.Nil(t, resultData.HvacSystemFunctionData[0].IsOverrunActive)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		HvacSystemFunctionListDataSelectors: &HvacSystemFunctionListDataSelectorsType{
			SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacSystemFunctionListDataType)
	assert.Equal(t, 0, len(resultData.HvacSystemFunctionData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestHvacSystemFunctionOperationModeRelationListDataType_ReadPartialData(t *testing.T) {
	sut := &HvacSystemFunctionOperationModeRelationListDataType{
		HvacSystemFunctionOperationModeRelationData: []HvacSystemFunctionOperationModeRelationDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				OperationModeId:  []HvacOperationModeIdType{10, 20},
			},
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(2)),
				OperationModeId:  []HvacOperationModeIdType{30},
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*HvacSystemFunctionOperationModeRelationListDataType)
	assert.Equal(t, 2, len(resultData.HvacSystemFunctionOperationModeRelationData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		HvacSystemFunctionOperationModeRelationListDataSelectors: &HvacSystemFunctionOperationModeRelationListDataSelectorsType{
			SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacSystemFunctionOperationModeRelationListDataType)
	assert.Equal(t, 1, len(resultData.HvacSystemFunctionOperationModeRelationData))
	assert.Equal(t, HvacSystemFunctionIdType(1), *resultData.HvacSystemFunctionOperationModeRelationData[0].SystemFunctionId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		HvacSystemFunctionOperationModeRelationDataElements: &HvacSystemFunctionOperationModeRelationDataElementsType{
			SystemFunctionId: &ElementTagType{},
			OperationModeId:  &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacSystemFunctionOperationModeRelationListDataType)
	assert.Equal(t, 2, len(resultData.HvacSystemFunctionOperationModeRelationData))
	assert.NotNil(t, resultData.HvacSystemFunctionOperationModeRelationData[0].SystemFunctionId)
	assert.NotNil(t, resultData.HvacSystemFunctionOperationModeRelationData[0].OperationModeId)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		HvacSystemFunctionOperationModeRelationListDataSelectors: &HvacSystemFunctionOperationModeRelationListDataSelectorsType{
			SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(2)),
		},
		HvacSystemFunctionOperationModeRelationDataElements: &HvacSystemFunctionOperationModeRelationDataElementsType{
			SystemFunctionId: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacSystemFunctionOperationModeRelationListDataType)
	assert.Equal(t, 1, len(resultData.HvacSystemFunctionOperationModeRelationData))
	assert.Equal(t, HvacSystemFunctionIdType(2), *resultData.HvacSystemFunctionOperationModeRelationData[0].SystemFunctionId)
	assert.Nil(t, resultData.HvacSystemFunctionOperationModeRelationData[0].OperationModeId)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		HvacSystemFunctionOperationModeRelationListDataSelectors: &HvacSystemFunctionOperationModeRelationListDataSelectorsType{
			SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacSystemFunctionOperationModeRelationListDataType)
	assert.Equal(t, 0, len(resultData.HvacSystemFunctionOperationModeRelationData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestHvacSystemFunctionSetpointRelationListDataType_ReadPartialData(t *testing.T) {
	sut := &HvacSystemFunctionSetpointRelationListDataType{
		HvacSystemFunctionSetpointRelationData: []HvacSystemFunctionSetpointRelationDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				OperationModeId:  util.Ptr(HvacOperationModeIdType(10)),
				SetpointId:       []SetpointIdType{100, 200},
			},
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(2)),
				OperationModeId:  util.Ptr(HvacOperationModeIdType(20)),
				SetpointId:       []SetpointIdType{300},
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*HvacSystemFunctionSetpointRelationListDataType)
	assert.Equal(t, 2, len(resultData.HvacSystemFunctionSetpointRelationData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		HvacSystemFunctionSetpointRelationListDataSelectors: &HvacSystemFunctionSetpointRelationListDataSelectorsType{
			SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacSystemFunctionSetpointRelationListDataType)
	assert.Equal(t, 1, len(resultData.HvacSystemFunctionSetpointRelationData))
	assert.Equal(t, HvacSystemFunctionIdType(1), *resultData.HvacSystemFunctionSetpointRelationData[0].SystemFunctionId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		HvacSystemFunctionSetpointRelationDataElements: &HvacSystemFunctionSetpointRelationDataElementsType{
			SystemFunctionId: &ElementTagType{},
			SetpointId:       &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacSystemFunctionSetpointRelationListDataType)
	assert.Equal(t, 2, len(resultData.HvacSystemFunctionSetpointRelationData))
	assert.NotNil(t, resultData.HvacSystemFunctionSetpointRelationData[0].SystemFunctionId)
	assert.NotNil(t, resultData.HvacSystemFunctionSetpointRelationData[0].SetpointId)
	assert.Nil(t, resultData.HvacSystemFunctionSetpointRelationData[0].OperationModeId)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		HvacSystemFunctionSetpointRelationListDataSelectors: &HvacSystemFunctionSetpointRelationListDataSelectorsType{
			SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(2)),
			OperationModeId:  util.Ptr(HvacOperationModeIdType(20)),
		},
		HvacSystemFunctionSetpointRelationDataElements: &HvacSystemFunctionSetpointRelationDataElementsType{
			SystemFunctionId: &ElementTagType{},
			OperationModeId:  &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacSystemFunctionSetpointRelationListDataType)
	assert.Equal(t, 1, len(resultData.HvacSystemFunctionSetpointRelationData))
	assert.Equal(t, HvacSystemFunctionIdType(2), *resultData.HvacSystemFunctionSetpointRelationData[0].SystemFunctionId)
	assert.NotNil(t, resultData.HvacSystemFunctionSetpointRelationData[0].OperationModeId)
	assert.Nil(t, resultData.HvacSystemFunctionSetpointRelationData[0].SetpointId)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		HvacSystemFunctionSetpointRelationListDataSelectors: &HvacSystemFunctionSetpointRelationListDataSelectorsType{
			SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacSystemFunctionSetpointRelationListDataType)
	assert.Equal(t, 0, len(resultData.HvacSystemFunctionSetpointRelationData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestHvacSystemFunctionPowerSequenceRelationListDataType_ReadPartialData(t *testing.T) {
	sut := &HvacSystemFunctionPowerSequenceRelationListDataType{
		HvacSystemFunctionPowerSequenceRelationData: []HvacSystemFunctionPowerSequenceRelationDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				SequenceId:       []PowerSequenceIdType{10, 20},
			},
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(2)),
				SequenceId:       []PowerSequenceIdType{30},
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*HvacSystemFunctionPowerSequenceRelationListDataType)
	assert.Equal(t, 2, len(resultData.HvacSystemFunctionPowerSequenceRelationData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		HvacSystemFunctionPowerSequenceRelationListDataSelectors: &HvacSystemFunctionPowerSequenceRelationListDataSelectorsType{
			SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacSystemFunctionPowerSequenceRelationListDataType)
	assert.Equal(t, 1, len(resultData.HvacSystemFunctionPowerSequenceRelationData))
	assert.Equal(t, HvacSystemFunctionIdType(1), *resultData.HvacSystemFunctionPowerSequenceRelationData[0].SystemFunctionId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		HvacSystemFunctionPowerSequenceRelationDataElements: &HvacSystemFunctionPowerSequenceRelationDataElementsType{
			SystemFunctionId: &ElementTagType{},
			SequenceId:       &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacSystemFunctionPowerSequenceRelationListDataType)
	assert.Equal(t, 2, len(resultData.HvacSystemFunctionPowerSequenceRelationData))
	assert.NotNil(t, resultData.HvacSystemFunctionPowerSequenceRelationData[0].SystemFunctionId)
	assert.NotNil(t, resultData.HvacSystemFunctionPowerSequenceRelationData[0].SequenceId)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		HvacSystemFunctionPowerSequenceRelationListDataSelectors: &HvacSystemFunctionPowerSequenceRelationListDataSelectorsType{
			SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(2)),
		},
		HvacSystemFunctionPowerSequenceRelationDataElements: &HvacSystemFunctionPowerSequenceRelationDataElementsType{
			SystemFunctionId: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacSystemFunctionPowerSequenceRelationListDataType)
	assert.Equal(t, 1, len(resultData.HvacSystemFunctionPowerSequenceRelationData))
	assert.Equal(t, HvacSystemFunctionIdType(2), *resultData.HvacSystemFunctionPowerSequenceRelationData[0].SystemFunctionId)
	assert.Nil(t, resultData.HvacSystemFunctionPowerSequenceRelationData[0].SequenceId)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		HvacSystemFunctionPowerSequenceRelationListDataSelectors: &HvacSystemFunctionPowerSequenceRelationListDataSelectorsType{
			SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacSystemFunctionPowerSequenceRelationListDataType)
	assert.Equal(t, 0, len(resultData.HvacSystemFunctionPowerSequenceRelationData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestHvacSystemFunctionDescriptionListDataType_ReadPartialData(t *testing.T) {
	sut := &HvacSystemFunctionDescriptionListDataType{
		HvacSystemFunctionDescriptionData: []HvacSystemFunctionDescriptionDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				Label:            util.Ptr(LabelType("Heating")),
				Description:      util.Ptr(DescriptionType("Primary heating system")),
			},
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(2)),
				Label:            util.Ptr(LabelType("Cooling")),
				Description:      util.Ptr(DescriptionType("Primary cooling system")),
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*HvacSystemFunctionDescriptionListDataType)
	assert.Equal(t, 2, len(resultData.HvacSystemFunctionDescriptionData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		HvacSystemFunctionDescriptionListDataSelectors: &HvacSystemFunctionDescriptionListDataSelectorsType{
			SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacSystemFunctionDescriptionListDataType)
	assert.Equal(t, 1, len(resultData.HvacSystemFunctionDescriptionData))
	assert.Equal(t, HvacSystemFunctionIdType(1), *resultData.HvacSystemFunctionDescriptionData[0].SystemFunctionId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		HvacSystemFunctionDescriptionDataElements: &HvacSystemFunctionDescriptionDataElementsType{
			SystemFunctionId: &ElementTagType{},
			Label:            &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacSystemFunctionDescriptionListDataType)
	assert.Equal(t, 2, len(resultData.HvacSystemFunctionDescriptionData))
	assert.NotNil(t, resultData.HvacSystemFunctionDescriptionData[0].SystemFunctionId)
	assert.NotNil(t, resultData.HvacSystemFunctionDescriptionData[0].Label)
	assert.Nil(t, resultData.HvacSystemFunctionDescriptionData[0].Description)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		HvacSystemFunctionDescriptionListDataSelectors: &HvacSystemFunctionDescriptionListDataSelectorsType{
			SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(2)),
		},
		HvacSystemFunctionDescriptionDataElements: &HvacSystemFunctionDescriptionDataElementsType{
			SystemFunctionId: &ElementTagType{},
			Description:      &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacSystemFunctionDescriptionListDataType)
	assert.Equal(t, 1, len(resultData.HvacSystemFunctionDescriptionData))
	assert.Equal(t, HvacSystemFunctionIdType(2), *resultData.HvacSystemFunctionDescriptionData[0].SystemFunctionId)
	assert.NotNil(t, resultData.HvacSystemFunctionDescriptionData[0].Description)
	assert.Nil(t, resultData.HvacSystemFunctionDescriptionData[0].Label)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		HvacSystemFunctionDescriptionListDataSelectors: &HvacSystemFunctionDescriptionListDataSelectorsType{
			SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacSystemFunctionDescriptionListDataType)
	assert.Equal(t, 0, len(resultData.HvacSystemFunctionDescriptionData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestHvacOperationModeDescriptionListDataType_ReadPartialData(t *testing.T) {
	sut := &HvacOperationModeDescriptionListDataType{
		HvacOperationModeDescriptionData: []HvacOperationModeDescriptionDataType{
			{
				OperationModeId: util.Ptr(HvacOperationModeIdType(10)),
				Label:           util.Ptr(LabelType("Auto")),
				Description:     util.Ptr(DescriptionType("Automatic mode")),
			},
			{
				OperationModeId: util.Ptr(HvacOperationModeIdType(20)),
				Label:           util.Ptr(LabelType("Manual")),
				Description:     util.Ptr(DescriptionType("Manual mode")),
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*HvacOperationModeDescriptionListDataType)
	assert.Equal(t, 2, len(resultData.HvacOperationModeDescriptionData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		HvacOperationModeDescriptionListDataSelectors: &HvacOperationModeDescriptionListDataSelectorsType{
			OperationModeId: util.Ptr(HvacOperationModeIdType(10)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacOperationModeDescriptionListDataType)
	assert.Equal(t, 1, len(resultData.HvacOperationModeDescriptionData))
	assert.Equal(t, HvacOperationModeIdType(10), *resultData.HvacOperationModeDescriptionData[0].OperationModeId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		HvacOperationModeDescriptionDataElements: &HvacOperationModeDescriptionDataElementsType{
			OperationModeId: &ElementTagType{},
			Label:           &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacOperationModeDescriptionListDataType)
	assert.Equal(t, 2, len(resultData.HvacOperationModeDescriptionData))
	assert.NotNil(t, resultData.HvacOperationModeDescriptionData[0].OperationModeId)
	assert.NotNil(t, resultData.HvacOperationModeDescriptionData[0].Label)
	assert.Nil(t, resultData.HvacOperationModeDescriptionData[0].Description)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		HvacOperationModeDescriptionListDataSelectors: &HvacOperationModeDescriptionListDataSelectorsType{
			OperationModeId: util.Ptr(HvacOperationModeIdType(20)),
		},
		HvacOperationModeDescriptionDataElements: &HvacOperationModeDescriptionDataElementsType{
			OperationModeId: &ElementTagType{},
			Description:     &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacOperationModeDescriptionListDataType)
	assert.Equal(t, 1, len(resultData.HvacOperationModeDescriptionData))
	assert.Equal(t, HvacOperationModeIdType(20), *resultData.HvacOperationModeDescriptionData[0].OperationModeId)
	assert.NotNil(t, resultData.HvacOperationModeDescriptionData[0].Description)
	assert.Nil(t, resultData.HvacOperationModeDescriptionData[0].Label)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		HvacOperationModeDescriptionListDataSelectors: &HvacOperationModeDescriptionListDataSelectorsType{
			OperationModeId: util.Ptr(HvacOperationModeIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacOperationModeDescriptionListDataType)
	assert.Equal(t, 0, len(resultData.HvacOperationModeDescriptionData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestHvacOverrunListDataType_ReadPartialData(t *testing.T) {
	sut := &HvacOverrunListDataType{
		HvacOverrunData: []HvacOverrunDataType{
			{
				OverrunId:                 util.Ptr(HvacOverrunIdType(1)),
				OverrunStatus:             util.Ptr(HvacOverrunStatusType("active")),
				TimeTableId:               util.Ptr(TimeTableIdType(10)),
				IsOverrunStatusChangeable: util.Ptr(true),
			},
			{
				OverrunId:                 util.Ptr(HvacOverrunIdType(2)),
				OverrunStatus:             util.Ptr(HvacOverrunStatusType("inactive")),
				TimeTableId:               util.Ptr(TimeTableIdType(20)),
				IsOverrunStatusChangeable: util.Ptr(false),
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*HvacOverrunListDataType)
	assert.Equal(t, 2, len(resultData.HvacOverrunData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		HvacOverrunListDataSelectors: &HvacOverrunListDataSelectorsType{
			OverrunId: util.Ptr(HvacOverrunIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacOverrunListDataType)
	assert.Equal(t, 1, len(resultData.HvacOverrunData))
	assert.Equal(t, HvacOverrunIdType(1), *resultData.HvacOverrunData[0].OverrunId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		HvacOverrunDataElements: &HvacOverrunDataElementsType{
			OverrunId:     &ElementTagType{},
			OverrunStatus: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacOverrunListDataType)
	assert.Equal(t, 2, len(resultData.HvacOverrunData))
	assert.NotNil(t, resultData.HvacOverrunData[0].OverrunId)
	assert.NotNil(t, resultData.HvacOverrunData[0].OverrunStatus)
	assert.Nil(t, resultData.HvacOverrunData[0].TimeTableId)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		HvacOverrunListDataSelectors: &HvacOverrunListDataSelectorsType{
			OverrunId: util.Ptr(HvacOverrunIdType(2)),
		},
		HvacOverrunDataElements: &HvacOverrunDataElementsType{
			OverrunId:                 &ElementTagType{},
			IsOverrunStatusChangeable: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacOverrunListDataType)
	assert.Equal(t, 1, len(resultData.HvacOverrunData))
	assert.Equal(t, HvacOverrunIdType(2), *resultData.HvacOverrunData[0].OverrunId)
	assert.NotNil(t, resultData.HvacOverrunData[0].IsOverrunStatusChangeable)
	assert.Nil(t, resultData.HvacOverrunData[0].OverrunStatus)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		HvacOverrunListDataSelectors: &HvacOverrunListDataSelectorsType{
			OverrunId: util.Ptr(HvacOverrunIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*HvacOverrunListDataType)
	assert.Equal(t, 0, len(resultData.HvacOverrunData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}
