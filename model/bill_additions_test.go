package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestBillListDataType_Update(t *testing.T) {
	sut := BillListDataType{
		BillData: []BillDataType{
			{
				BillId:    util.Ptr(BillIdType(0)),
				ScopeType: util.Ptr(ScopeTypeTypeACCurrent),
			},
			{
				BillId:    util.Ptr(BillIdType(1)),
				ScopeType: util.Ptr(ScopeTypeTypeACCurrent),
			},
		},
	}

	newData := BillListDataType{
		BillData: []BillDataType{
			{
				BillId:    util.Ptr(BillIdType(1)),
				ScopeType: util.Ptr(ScopeTypeTypeACPower),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.BillData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.BillId))
	assert.Equal(t, ScopeTypeTypeACCurrent, *item1.ScopeType)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.BillId))
	assert.Equal(t, ScopeTypeTypeACPower, *item2.ScopeType)
}

func TestBillConstraintsListDataType_Update(t *testing.T) {
	sut := BillConstraintsListDataType{
		BillConstraintsData: []BillConstraintsDataType{
			{
				BillId:           util.Ptr(BillIdType(0)),
				PositionCountMin: util.Ptr(BillPositionCountType(0)),
			},
			{
				BillId:           util.Ptr(BillIdType(1)),
				PositionCountMin: util.Ptr(BillPositionCountType(0)),
			},
		},
	}

	newData := BillConstraintsListDataType{
		BillConstraintsData: []BillConstraintsDataType{
			{
				BillId:           util.Ptr(BillIdType(1)),
				PositionCountMin: util.Ptr(BillPositionCountType(1)),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.BillConstraintsData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.BillId))
	assert.Equal(t, 0, int(*item1.PositionCountMin))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.BillId))
	assert.Equal(t, 1, int(*item2.PositionCountMin))
}

func TestBillDescriptionListDataType_Update(t *testing.T) {
	sut := BillDescriptionListDataType{
		BillDescriptionData: []BillDescriptionDataType{
			{
				BillId:         util.Ptr(BillIdType(0)),
				UpdateRequired: util.Ptr(false),
			},
			{
				BillId:         util.Ptr(BillIdType(1)),
				UpdateRequired: util.Ptr(false),
			},
		},
	}

	newData := BillDescriptionListDataType{
		BillDescriptionData: []BillDescriptionDataType{
			{
				BillId:         util.Ptr(BillIdType(1)),
				UpdateRequired: util.Ptr(true),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.BillDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.BillId))
	assert.Equal(t, false, *item1.UpdateRequired)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.BillId))
	assert.Equal(t, true, *item2.UpdateRequired)
}

func TestBillListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	data := &BillListDataType{
		BillData: []BillDataType{
			{
				BillId:    util.Ptr(BillIdType(1)),
				BillType:  util.Ptr(BillTypeTypeChargingSummary),
				ScopeType: util.Ptr(ScopeTypeTypeACCurrent),
				Total: &BillPositionType{
					PositionId: util.Ptr(BillPositionIdType(1)),
					Value: &BillValueType{
						Value: NewScaledNumberType(100),
					},
				},
			},
			{
				BillId:    util.Ptr(BillIdType(2)),
				BillType:  util.Ptr(BillTypeTypeChargingSummary),
				ScopeType: util.Ptr(ScopeTypeTypeACPower),
				Total: &BillPositionType{
					PositionId: util.Ptr(BillPositionIdType(2)),
					Value: &BillValueType{
						Value: NewScaledNumberType(200),
					},
				},
			},
		},
	}

	// Test 1: No filter (should return all data)
	result, success := data.ReadPartialData(nil)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData := result.(*BillListDataType)
	assert.Len(t, resultData.BillData, 2)

	// Test 2: Selector filter by BillId
	selectorFilter := &FilterType{
		BillListDataSelectors: &BillListDataSelectorsType{
			BillId: util.Ptr(BillIdType(1)),
		},
	}
	result, success = data.ReadPartialData(selectorFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*BillListDataType)
	assert.Len(t, resultData.BillData, 1)
	assert.Equal(t, BillIdType(1), *resultData.BillData[0].BillId)

	// Test 3: Selector filter by ScopeType
	scopeFilter := &FilterType{
		BillListDataSelectors: &BillListDataSelectorsType{
			ScopeType: util.Ptr(ScopeTypeTypeACPower),
		},
	}
	result, success = data.ReadPartialData(scopeFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*BillListDataType)
	assert.Len(t, resultData.BillData, 1)
	assert.Equal(t, BillIdType(2), *resultData.BillData[0].BillId)

	// Test 4: Element filter
	elementFilter := &FilterType{
		BillDataElements: &BillDataElementsType{
			BillId:   util.Ptr(ElementTagType{}),
			BillType: util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(elementFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*BillListDataType)
	assert.Len(t, resultData.BillData, 2)
	// Check that only requested fields are present
	assert.NotNil(t, resultData.BillData[0].BillId)
	assert.NotNil(t, resultData.BillData[0].BillType)
	// Check that non-requested fields are filtered out
	assert.Nil(t, resultData.BillData[0].ScopeType)
	assert.Nil(t, resultData.BillData[0].Total)

	// Test 5: Combined selector and element filter
	combinedFilter := &FilterType{
		BillListDataSelectors: &BillListDataSelectorsType{
			ScopeType: util.Ptr(ScopeTypeTypeACCurrent),
		},
		BillDataElements: &BillDataElementsType{
			BillId: util.Ptr(ElementTagType{}),
			Total:  &BillPositionElementsType{},
		},
	}
	result, success = data.ReadPartialData(combinedFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*BillListDataType)
	assert.Len(t, resultData.BillData, 1)
	assert.Equal(t, BillIdType(1), *resultData.BillData[0].BillId)
	assert.NotNil(t, resultData.BillData[0].Total)
	assert.Nil(t, resultData.BillData[0].BillType)
	assert.Nil(t, resultData.BillData[0].ScopeType)

	// Test 6: Non-matching selector (should return empty)
	nonMatchingFilter := &FilterType{
		BillListDataSelectors: &BillListDataSelectorsType{
			BillId: util.Ptr(BillIdType(999)),
		},
	}
	result, success = data.ReadPartialData(nonMatchingFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*BillListDataType)
	assert.Len(t, resultData.BillData, 0)

	// Test 7: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = data.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestBillConstraintsListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	data := &BillConstraintsListDataType{
		BillConstraintsData: []BillConstraintsDataType{
			{
				BillId:           util.Ptr(BillIdType(1)),
				PositionCountMin: util.Ptr(BillPositionCountType(1)),
				PositionCountMax: util.Ptr(BillPositionCountType(10)),
			},
			{
				BillId:           util.Ptr(BillIdType(2)),
				PositionCountMin: util.Ptr(BillPositionCountType(5)),
				PositionCountMax: util.Ptr(BillPositionCountType(20)),
			},
		},
	}

	// Test 1: No filter (should return all data)
	result, success := data.ReadPartialData(nil)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData := result.(*BillConstraintsListDataType)
	assert.Len(t, resultData.BillConstraintsData, 2)

	// Test 2: Selector filter
	selectorFilter := &FilterType{
		BillConstraintsListDataSelectors: &BillConstraintsListDataSelectorsType{
			BillId: util.Ptr(BillIdType(2)),
		},
	}
	result, success = data.ReadPartialData(selectorFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*BillConstraintsListDataType)
	assert.Len(t, resultData.BillConstraintsData, 1)
	assert.Equal(t, BillIdType(2), *resultData.BillConstraintsData[0].BillId)

	// Test 3: Element filter
	elementFilter := &FilterType{
		BillConstraintsDataElements: &BillConstraintsDataElementsType{
			BillId:           util.Ptr(ElementTagType{}),
			PositionCountMin: util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(elementFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*BillConstraintsListDataType)
	assert.Len(t, resultData.BillConstraintsData, 2)
	// Check that only requested fields are present
	assert.NotNil(t, resultData.BillConstraintsData[0].BillId)
	assert.NotNil(t, resultData.BillConstraintsData[0].PositionCountMin)
	// Check that non-requested field is filtered out
	assert.Nil(t, resultData.BillConstraintsData[0].PositionCountMax)

	// Test 4: Combined selector and element filter
	combinedFilter := &FilterType{
		BillConstraintsListDataSelectors: &BillConstraintsListDataSelectorsType{
			BillId: util.Ptr(BillIdType(1)),
		},
		BillConstraintsDataElements: &BillConstraintsDataElementsType{
			BillId:           util.Ptr(ElementTagType{}),
			PositionCountMax: util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(combinedFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*BillConstraintsListDataType)
	assert.Len(t, resultData.BillConstraintsData, 1)
	assert.Equal(t, BillIdType(1), *resultData.BillConstraintsData[0].BillId)
	assert.NotNil(t, resultData.BillConstraintsData[0].PositionCountMax)
	assert.Nil(t, resultData.BillConstraintsData[0].PositionCountMin)

	// Test 5: Non-matching selector (should return empty)
	nonMatchingFilter := &FilterType{
		BillConstraintsListDataSelectors: &BillConstraintsListDataSelectorsType{
			BillId: util.Ptr(BillIdType(999)),
		},
	}
	result, success = data.ReadPartialData(nonMatchingFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*BillConstraintsListDataType)
	assert.Len(t, resultData.BillConstraintsData, 0)

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = data.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestBillDescriptionListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	data := &BillDescriptionListDataType{
		BillDescriptionData: []BillDescriptionDataType{
			{
				BillId:            util.Ptr(BillIdType(1)),
				BillWriteable:     util.Ptr(true),
				UpdateRequired:    util.Ptr(false),
				SupportedBillType: []BillTypeType{BillTypeTypeChargingSummary},
				SessionId:         util.Ptr(SessionIdType(100)),
			},
			{
				BillId:            util.Ptr(BillIdType(2)),
				BillWriteable:     util.Ptr(false),
				UpdateRequired:    util.Ptr(true),
				SupportedBillType: []BillTypeType{BillTypeTypeChargingSummary},
				SessionId:         util.Ptr(SessionIdType(200)),
			},
		},
	}

	// Test 1: No filter (should return all data)
	result, success := data.ReadPartialData(nil)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData := result.(*BillDescriptionListDataType)
	assert.Len(t, resultData.BillDescriptionData, 2)

	// Test 2: Selector filter
	selectorFilter := &FilterType{
		BillDescriptionListDataSelectors: &BillDescriptionListDataSelectorsType{
			BillId: util.Ptr(BillIdType(2)),
		},
	}
	result, success = data.ReadPartialData(selectorFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*BillDescriptionListDataType)
	assert.Len(t, resultData.BillDescriptionData, 1)
	assert.Equal(t, BillIdType(2), *resultData.BillDescriptionData[0].BillId)

	// Test 3: Element filter
	elementFilter := &FilterType{
		BillDescriptionDataElements: &BillDescriptionDataElementsType{
			BillId:         util.Ptr(ElementTagType{}),
			BillWriteable:  util.Ptr(ElementTagType{}),
			UpdateRequired: util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(elementFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*BillDescriptionListDataType)
	assert.Len(t, resultData.BillDescriptionData, 2)
	// Check that only requested fields are present
	assert.NotNil(t, resultData.BillDescriptionData[0].BillId)
	assert.NotNil(t, resultData.BillDescriptionData[0].BillWriteable)
	assert.NotNil(t, resultData.BillDescriptionData[0].UpdateRequired)
	// Check that non-requested fields are filtered out
	assert.Nil(t, resultData.BillDescriptionData[0].SupportedBillType)
	assert.Nil(t, resultData.BillDescriptionData[0].SessionId)

	// Test 4: Combined selector and element filter
	combinedFilter := &FilterType{
		BillDescriptionListDataSelectors: &BillDescriptionListDataSelectorsType{
			BillId: util.Ptr(BillIdType(1)),
		},
		BillDescriptionDataElements: &BillDescriptionDataElementsType{
			BillId:            util.Ptr(ElementTagType{}),
			SupportedBillType: util.Ptr(ElementTagType{}),
		},
	}
	result, success = data.ReadPartialData(combinedFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*BillDescriptionListDataType)
	assert.Len(t, resultData.BillDescriptionData, 1)
	assert.Equal(t, BillIdType(1), *resultData.BillDescriptionData[0].BillId)
	assert.NotNil(t, resultData.BillDescriptionData[0].SupportedBillType)
	assert.Nil(t, resultData.BillDescriptionData[0].BillWriteable)
	assert.Nil(t, resultData.BillDescriptionData[0].SessionId)

	// Test 5: Non-matching selector (should return empty)
	nonMatchingFilter := &FilterType{
		BillDescriptionListDataSelectors: &BillDescriptionListDataSelectorsType{
			BillId: util.Ptr(BillIdType(999)),
		},
	}
	result, success = data.ReadPartialData(nonMatchingFilter)
	assert.True(t, success)
	assert.NotNil(t, result)
	resultData = result.(*BillDescriptionListDataType)
	assert.Len(t, resultData.BillDescriptionData, 0)

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = data.ReadPartialData(emptyFilter)
	assert.False(t, success)
}
