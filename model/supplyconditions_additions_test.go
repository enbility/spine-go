package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestSupplyConditionListDataType_Update(t *testing.T) {
	sut := SupplyConditionListDataType{
		SupplyConditionData: []SupplyConditionDataType{
			{
				ConditionId: util.Ptr(ConditionIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				ConditionId: util.Ptr(ConditionIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := SupplyConditionListDataType{
		SupplyConditionData: []SupplyConditionDataType{
			{
				ConditionId: util.Ptr(ConditionIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.SupplyConditionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ConditionId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ConditionId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestSupplyConditionDescriptionListDataType_Update(t *testing.T) {
	sut := SupplyConditionDescriptionListDataType{
		SupplyConditionDescriptionData: []SupplyConditionDescriptionDataType{
			{
				ConditionId: util.Ptr(ConditionIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				ConditionId: util.Ptr(ConditionIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := SupplyConditionDescriptionListDataType{
		SupplyConditionDescriptionData: []SupplyConditionDescriptionDataType{
			{
				ConditionId: util.Ptr(ConditionIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.SupplyConditionDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ConditionId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ConditionId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestSupplyConditionThresholdRelationListDataType_Update(t *testing.T) {
	sut := SupplyConditionThresholdRelationListDataType{
		SupplyConditionThresholdRelationData: []SupplyConditionThresholdRelationDataType{
			{
				ConditionId: util.Ptr(ConditionIdType(0)),
				ThresholdId: []ThresholdIdType{0},
			},
			{
				ConditionId: util.Ptr(ConditionIdType(1)),
				ThresholdId: []ThresholdIdType{0},
			},
		},
	}

	newData := SupplyConditionThresholdRelationListDataType{
		SupplyConditionThresholdRelationData: []SupplyConditionThresholdRelationDataType{
			{
				ConditionId: util.Ptr(ConditionIdType(1)),
				ThresholdId: []ThresholdIdType{1},
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.SupplyConditionThresholdRelationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ConditionId))
	assert.Equal(t, 0, int(item1.ThresholdId[0]))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ConditionId))
	assert.Equal(t, 1, int(item2.ThresholdId[0]))
}

func TestSupplyConditionListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	data := &SupplyConditionListDataType{
		SupplyConditionData: []SupplyConditionDataType{
			{
				ConditionId:         util.Ptr(ConditionIdType(1)),
				Timestamp:           util.Ptr(AbsoluteOrRelativeTimeType("2023-01-01T10:00:00Z")),
				EventType:           util.Ptr(SupplyConditionEventTypeTypeThesholdExceeded),
				Originator:          util.Ptr(SupplyConditionOriginatorTypeExternDSO),
				ThresholdId:         util.Ptr(ThresholdIdType(10)),
				ThresholdPercentage: NewScaledNumberType(75),
				Description:         util.Ptr(DescriptionType("Condition 1")),
			},
			{
				ConditionId:         util.Ptr(ConditionIdType(2)),
				Timestamp:           util.Ptr(AbsoluteOrRelativeTimeType("2023-01-01T11:00:00Z")),
				EventType:           util.Ptr(SupplyConditionEventTypeTypeFallenBelowThreshold),
				Originator:          util.Ptr(SupplyConditionOriginatorTypeExternSupplier),
				ThresholdId:         util.Ptr(ThresholdIdType(20)),
				ThresholdPercentage: NewScaledNumberType(50),
				Description:         util.Ptr(DescriptionType("Condition 2")),
			},
		},
	}

	// Test 1: No filter (should return all data)
	result, success := data.ReadPartialData(nil)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData := result.(*SupplyConditionListDataType)
	assert.Len(t, resultData.SupplyConditionData, 2)

	// Test 2: Selector filter by ConditionId
	selectorFilter := &FilterType{
		SupplyConditionListDataSelectors: &SupplyConditionListDataSelectorsType{
			ConditionId: util.Ptr(ConditionIdType(1)),
		},
	}
	result, success = data.ReadPartialData(selectorFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*SupplyConditionListDataType)
	assert.Len(t, resultData.SupplyConditionData, 1)
	assert.Equal(t, ConditionIdType(1), *resultData.SupplyConditionData[0].ConditionId)

	// Test 3: Selector filter by EventType
	eventFilter := &FilterType{
		SupplyConditionListDataSelectors: &SupplyConditionListDataSelectorsType{
			EventType: util.Ptr(SupplyConditionEventTypeTypeFallenBelowThreshold),
		},
	}
	result, success = data.ReadPartialData(eventFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*SupplyConditionListDataType)
	assert.Len(t, resultData.SupplyConditionData, 1)
	assert.Equal(t, ConditionIdType(2), *resultData.SupplyConditionData[0].ConditionId)

	// Test 4: Element filter
	elementFilter := &FilterType{
		SupplyConditionDataElements: &SupplyConditionDataElementsType{
			ConditionId: util.Ptr(ElementTagType{}),
			EventType:   util.Ptr(ElementTagType{}),
			Timestamp:   util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(elementFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*SupplyConditionListDataType)
	assert.Len(t, resultData.SupplyConditionData, 2)
	// Check that only requested fields are present
	assert.NotNil(t, resultData.SupplyConditionData[0].ConditionId)
	assert.NotNil(t, resultData.SupplyConditionData[0].EventType)
	assert.NotNil(t, resultData.SupplyConditionData[0].Timestamp)
	// Check that non-requested fields are filtered out
	assert.Nil(t, resultData.SupplyConditionData[0].Originator)
	assert.Nil(t, resultData.SupplyConditionData[0].ThresholdId)
	assert.Nil(t, resultData.SupplyConditionData[0].Description)

	// Test 5: Combined selector and element filter
	combinedFilter := &FilterType{
		SupplyConditionListDataSelectors: &SupplyConditionListDataSelectorsType{
			Originator: util.Ptr(SupplyConditionOriginatorTypeExternSupplier),
		},
		SupplyConditionDataElements: &SupplyConditionDataElementsType{
			ConditionId:         util.Ptr(ElementTagType{}),
			ThresholdPercentage: &ScaledNumberElementsType{},
		},
	}
	result, success = data.ReadPartialData(combinedFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*SupplyConditionListDataType)
	assert.Len(t, resultData.SupplyConditionData, 1)
	assert.Equal(t, ConditionIdType(2), *resultData.SupplyConditionData[0].ConditionId)
	assert.NotNil(t, resultData.SupplyConditionData[0].ThresholdPercentage)
	assert.Nil(t, resultData.SupplyConditionData[0].EventType)
	assert.Nil(t, resultData.SupplyConditionData[0].Timestamp)

	// Test 6: Non-matching selector (should return empty)
	nonMatchingFilter := &FilterType{
		SupplyConditionListDataSelectors: &SupplyConditionListDataSelectorsType{
			ConditionId: util.Ptr(ConditionIdType(999)),
		},
	}
	result, success = data.ReadPartialData(nonMatchingFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*SupplyConditionListDataType)
	assert.Len(t, resultData.SupplyConditionData, 0)

	// Test 7: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = data.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestSupplyConditionDescriptionListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	data := &SupplyConditionDescriptionListDataType{
		SupplyConditionDescriptionData: []SupplyConditionDescriptionDataType{
			{
				ConditionId:             util.Ptr(ConditionIdType(1)),
				CommodityType:           util.Ptr(CommodityTypeTypeElectricity),
				PositiveEnergyDirection: util.Ptr(EnergyDirectionTypeConsume),
				Label:                   util.Ptr(LabelType("Condition 1")),
				Description:             util.Ptr(DescriptionType("First condition")),
			},
			{
				ConditionId:             util.Ptr(ConditionIdType(2)),
				CommodityType:           util.Ptr(CommodityTypeTypeGas),
				PositiveEnergyDirection: util.Ptr(EnergyDirectionTypeProduce),
				Label:                   util.Ptr(LabelType("Condition 2")),
				Description:             util.Ptr(DescriptionType("Second condition")),
			},
		},
	}

	// Test 1: No filter (should return all data)
	result, success := data.ReadPartialData(nil)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData := result.(*SupplyConditionDescriptionListDataType)
	assert.Len(t, resultData.SupplyConditionDescriptionData, 2)

	// Test 2: Selector filter
	selectorFilter := &FilterType{
		SupplyConditionDescriptionListDataSelectors: &SupplyConditionDescriptionListDataSelectorsType{
			ConditionId: util.Ptr(ConditionIdType(2)),
		},
	}
	result, success = data.ReadPartialData(selectorFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*SupplyConditionDescriptionListDataType)
	assert.Len(t, resultData.SupplyConditionDescriptionData, 1)
	assert.Equal(t, ConditionIdType(2), *resultData.SupplyConditionDescriptionData[0].ConditionId)

	// Test 3: Element filter
	elementFilter := &FilterType{
		SupplyConditionDescriptionDataElements: &SupplyConditionDescriptionDataElementsType{
			ConditionId:             util.Ptr(ElementTagType{}),
			CommodityType:           util.Ptr(ElementTagType{}),
			PositiveEnergyDirection: util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(elementFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*SupplyConditionDescriptionListDataType)
	assert.Len(t, resultData.SupplyConditionDescriptionData, 2)
	// Check that only requested fields are present
	assert.NotNil(t, resultData.SupplyConditionDescriptionData[0].ConditionId)
	assert.NotNil(t, resultData.SupplyConditionDescriptionData[0].CommodityType)
	assert.NotNil(t, resultData.SupplyConditionDescriptionData[0].PositiveEnergyDirection)
	// Check that non-requested fields are filtered out
	assert.Nil(t, resultData.SupplyConditionDescriptionData[0].Label)
	assert.Nil(t, resultData.SupplyConditionDescriptionData[0].Description)

	// Test 4: Combined selector and element filter
	combinedFilter := &FilterType{
		SupplyConditionDescriptionListDataSelectors: &SupplyConditionDescriptionListDataSelectorsType{
			ConditionId: util.Ptr(ConditionIdType(1)),
		},
		SupplyConditionDescriptionDataElements: &SupplyConditionDescriptionDataElementsType{
			ConditionId: util.Ptr(ElementTagType{}),
			Label:       util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(combinedFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*SupplyConditionDescriptionListDataType)
	assert.Len(t, resultData.SupplyConditionDescriptionData, 1)
	assert.Equal(t, ConditionIdType(1), *resultData.SupplyConditionDescriptionData[0].ConditionId)
	assert.NotNil(t, resultData.SupplyConditionDescriptionData[0].Label)
	assert.Nil(t, resultData.SupplyConditionDescriptionData[0].CommodityType)
	assert.Nil(t, resultData.SupplyConditionDescriptionData[0].Description)

	// Test 5: Non-matching selector (should return empty)
	nonMatchingFilter := &FilterType{
		SupplyConditionDescriptionListDataSelectors: &SupplyConditionDescriptionListDataSelectorsType{
			ConditionId: util.Ptr(ConditionIdType(999)),
		},
	}
	result, success = data.ReadPartialData(nonMatchingFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*SupplyConditionDescriptionListDataType)
	assert.Len(t, resultData.SupplyConditionDescriptionData, 0)

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = data.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestSupplyConditionThresholdRelationListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	data := &SupplyConditionThresholdRelationListDataType{
		SupplyConditionThresholdRelationData: []SupplyConditionThresholdRelationDataType{
			{
				ConditionId: util.Ptr(ConditionIdType(1)),
				ThresholdId: []ThresholdIdType{10, 11, 12},
			},
			{
				ConditionId: util.Ptr(ConditionIdType(2)),
				ThresholdId: []ThresholdIdType{20, 21},
			},
		},
	}

	// Test 1: No filter (should return all data)
	result, success := data.ReadPartialData(nil)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData := result.(*SupplyConditionThresholdRelationListDataType)
	assert.Len(t, resultData.SupplyConditionThresholdRelationData, 2)

	// Test 2: Selector filter by ConditionId
	selectorFilter := &FilterType{
		SupplyConditionThresholdRelationListDataSelectors: &SupplyConditionThresholdRelationListDataSelectorsType{
			ConditionId: util.Ptr(ConditionIdType(1)),
		},
	}
	result, success = data.ReadPartialData(selectorFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*SupplyConditionThresholdRelationListDataType)
	assert.Len(t, resultData.SupplyConditionThresholdRelationData, 1)
	assert.Equal(t, ConditionIdType(1), *resultData.SupplyConditionThresholdRelationData[0].ConditionId)

	// Test 3: Selector filter by ThresholdId
	thresholdFilter := &FilterType{
		SupplyConditionThresholdRelationListDataSelectors: &SupplyConditionThresholdRelationListDataSelectorsType{
			ThresholdId: util.Ptr(ThresholdIdType(20)),
		},
	}
	result, success = data.ReadPartialData(thresholdFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*SupplyConditionThresholdRelationListDataType)
	assert.Len(t, resultData.SupplyConditionThresholdRelationData, 1)
	assert.Equal(t, ConditionIdType(2), *resultData.SupplyConditionThresholdRelationData[0].ConditionId)

	// Test 4: Element filter
	elementFilter := &FilterType{
		SupplyConditionThresholdRelationDataElements: &SupplyConditionThresholdRelationDataElementsType{
			ConditionId: util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(elementFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*SupplyConditionThresholdRelationListDataType)
	assert.Len(t, resultData.SupplyConditionThresholdRelationData, 2)
	// Check that only requested fields are present
	assert.NotNil(t, resultData.SupplyConditionThresholdRelationData[0].ConditionId)
	// Check that non-requested field is filtered out
	assert.Nil(t, resultData.SupplyConditionThresholdRelationData[0].ThresholdId)

	// Test 5: Combined selector and element filter
	combinedFilter := &FilterType{
		SupplyConditionThresholdRelationListDataSelectors: &SupplyConditionThresholdRelationListDataSelectorsType{
			ConditionId: util.Ptr(ConditionIdType(2)),
		},
		SupplyConditionThresholdRelationDataElements: &SupplyConditionThresholdRelationDataElementsType{
			ConditionId: util.Ptr(ElementTagType{}),
			ThresholdId: util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(combinedFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*SupplyConditionThresholdRelationListDataType)
	assert.Len(t, resultData.SupplyConditionThresholdRelationData, 1)
	assert.Equal(t, ConditionIdType(2), *resultData.SupplyConditionThresholdRelationData[0].ConditionId)
	assert.NotNil(t, resultData.SupplyConditionThresholdRelationData[0].ThresholdId)

	// Test 6: Non-matching selector (should return empty)
	nonMatchingFilter := &FilterType{
		SupplyConditionThresholdRelationListDataSelectors: &SupplyConditionThresholdRelationListDataSelectorsType{
			ConditionId: util.Ptr(ConditionIdType(999)),
		},
	}
	result, success = data.ReadPartialData(nonMatchingFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*SupplyConditionThresholdRelationListDataType)
	assert.Len(t, resultData.SupplyConditionThresholdRelationData, 0)

	// Test 7: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = data.ReadPartialData(emptyFilter)
	assert.False(t, success)
}
