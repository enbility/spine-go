package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestTariffListDataType_Update(t *testing.T) {
	sut := TariffListDataType{
		TariffData: []TariffDataType{
			{
				TariffId:     util.Ptr(TariffIdType(0)),
				ActiveTierId: []TierIdType{0},
			},
			{
				TariffId:     util.Ptr(TariffIdType(1)),
				ActiveTierId: []TierIdType{0},
			},
		},
	}

	newData := TariffListDataType{
		TariffData: []TariffDataType{
			{
				TariffId:     util.Ptr(TariffIdType(1)),
				ActiveTierId: []TierIdType{1},
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.TariffData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TariffId))
	assert.Equal(t, 0, int(item1.ActiveTierId[0]))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TariffId))
	assert.Equal(t, 1, int(item2.ActiveTierId[0]))
}

func TestTariffTierRelationListDataType_Update(t *testing.T) {
	sut := TariffTierRelationListDataType{
		TariffTierRelationData: []TariffTierRelationDataType{
			{
				TariffId: util.Ptr(TariffIdType(0)),
				TierId:   []TierIdType{0},
			},
			{
				TariffId: util.Ptr(TariffIdType(1)),
				TierId:   []TierIdType{0},
			},
		},
	}

	newData := TariffTierRelationListDataType{
		TariffTierRelationData: []TariffTierRelationDataType{
			{
				TariffId: util.Ptr(TariffIdType(1)),
				TierId:   []TierIdType{1},
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.TariffTierRelationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TariffId))
	assert.Equal(t, 0, int(item1.TierId[0]))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TariffId))
	assert.Equal(t, 1, int(item2.TierId[0]))
}

func TestTariffBoundaryRelationListDataType_Update(t *testing.T) {
	sut := TariffBoundaryRelationListDataType{
		TariffBoundaryRelationData: []TariffBoundaryRelationDataType{
			{
				TariffId:   util.Ptr(TariffIdType(0)),
				BoundaryId: []TierBoundaryIdType{0},
			},
			{
				TariffId:   util.Ptr(TariffIdType(1)),
				BoundaryId: []TierBoundaryIdType{0},
			},
		},
	}

	newData := TariffBoundaryRelationListDataType{
		TariffBoundaryRelationData: []TariffBoundaryRelationDataType{
			{
				TariffId:   util.Ptr(TariffIdType(1)),
				BoundaryId: []TierBoundaryIdType{1},
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.TariffBoundaryRelationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TariffId))
	assert.Equal(t, 0, int(item1.BoundaryId[0]))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TariffId))
	assert.Equal(t, 1, int(item2.BoundaryId[0]))
}

func TestTariffDescriptionListDataType_Update(t *testing.T) {
	sut := TariffDescriptionListDataType{
		TariffDescriptionData: []TariffDescriptionDataType{
			{
				TariffId:    util.Ptr(TariffIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				TariffId:    util.Ptr(TariffIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := TariffDescriptionListDataType{
		TariffDescriptionData: []TariffDescriptionDataType{
			{
				TariffId:    util.Ptr(TariffIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.TariffDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TariffId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TariffId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestTierBoundaryListDataType_Update(t *testing.T) {
	sut := TierBoundaryListDataType{
		TierBoundaryData: []TierBoundaryDataType{
			{
				BoundaryId:         util.Ptr(TierBoundaryIdType(0)),
				LowerBoundaryValue: NewScaledNumberType(1),
			},
			{
				BoundaryId:         util.Ptr(TierBoundaryIdType(1)),
				LowerBoundaryValue: NewScaledNumberType(1),
			},
		},
	}

	newData := TierBoundaryListDataType{
		TierBoundaryData: []TierBoundaryDataType{
			{
				BoundaryId:         util.Ptr(TierBoundaryIdType(1)),
				LowerBoundaryValue: NewScaledNumberType(10),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.TierBoundaryData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.BoundaryId))
	assert.Equal(t, 1.0, item1.LowerBoundaryValue.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.BoundaryId))
	assert.Equal(t, 10.0, item2.LowerBoundaryValue.GetValue())
}

func TestTierBoundaryDescriptionListDataType_Update(t *testing.T) {
	sut := TierBoundaryDescriptionListDataType{
		TierBoundaryDescriptionData: []TierBoundaryDescriptionDataType{
			{
				BoundaryId:  util.Ptr(TierBoundaryIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				BoundaryId:  util.Ptr(TierBoundaryIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := TierBoundaryDescriptionListDataType{
		TierBoundaryDescriptionData: []TierBoundaryDescriptionDataType{
			{
				BoundaryId:  util.Ptr(TierBoundaryIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.TierBoundaryDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.BoundaryId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.BoundaryId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestCommodityListDataType_Update(t *testing.T) {
	sut := CommodityListDataType{
		CommodityData: []CommodityDataType{
			{
				CommodityId: util.Ptr(CommodityIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				CommodityId: util.Ptr(CommodityIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := CommodityListDataType{
		CommodityData: []CommodityDataType{
			{
				CommodityId: util.Ptr(CommodityIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.CommodityData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.CommodityId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.CommodityId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestTierListDataType_Update(t *testing.T) {
	sut := TierListDataType{
		TierData: []TierDataType{
			{
				TierId:            util.Ptr(TierIdType(0)),
				ActiveIncentiveId: []IncentiveIdType{0},
			},
			{
				TierId:            util.Ptr(TierIdType(1)),
				ActiveIncentiveId: []IncentiveIdType{0},
			},
		},
	}

	newData := TierListDataType{
		TierData: []TierDataType{
			{
				TierId:            util.Ptr(TierIdType(1)),
				ActiveIncentiveId: []IncentiveIdType{1},
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.TierData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TierId))
	assert.Equal(t, 0, int(item1.ActiveIncentiveId[0]))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TierId))
	assert.Equal(t, 1, int(item2.ActiveIncentiveId[0]))
}

func TestTierIncentiveRelationListDataType_Update(t *testing.T) {
	sut := TierIncentiveRelationListDataType{
		TierIncentiveRelationData: []TierIncentiveRelationDataType{
			{
				TierId:      util.Ptr(TierIdType(0)),
				IncentiveId: []IncentiveIdType{0},
			},
			{
				TierId:      util.Ptr(TierIdType(1)),
				IncentiveId: []IncentiveIdType{0},
			},
		},
	}

	newData := TierIncentiveRelationListDataType{
		TierIncentiveRelationData: []TierIncentiveRelationDataType{
			{
				TierId:      util.Ptr(TierIdType(1)),
				IncentiveId: []IncentiveIdType{1},
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.TierIncentiveRelationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TierId))
	assert.Equal(t, 0, int(item1.IncentiveId[0]))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TierId))
	assert.Equal(t, 1, int(item2.IncentiveId[0]))
}

func TestTierDescriptionListDataType_Update(t *testing.T) {
	sut := TierDescriptionListDataType{
		TierDescriptionData: []TierDescriptionDataType{
			{
				TierId:      util.Ptr(TierIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				TierId:      util.Ptr(TierIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := TierDescriptionListDataType{
		TierDescriptionData: []TierDescriptionDataType{
			{
				TierId:      util.Ptr(TierIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.TierDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TierId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TierId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestIncentiveListDataType_Update(t *testing.T) {
	sut := IncentiveListDataType{
		IncentiveData: []IncentiveDataType{
			{
				IncentiveId: util.Ptr(IncentiveIdType(0)),
				Value:       NewScaledNumberType(1),
			},
			{
				IncentiveId: util.Ptr(IncentiveIdType(1)),
				Value:       NewScaledNumberType(1),
			},
		},
	}

	newData := IncentiveListDataType{
		IncentiveData: []IncentiveDataType{
			{
				IncentiveId: util.Ptr(IncentiveIdType(1)),
				Value:       NewScaledNumberType(10),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.IncentiveData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.IncentiveId))
	assert.Equal(t, 1.0, item1.Value.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.IncentiveId))
	assert.Equal(t, 10.0, item2.Value.GetValue())
}

func TestIncentiveDescriptionListDataType_Update(t *testing.T) {
	sut := IncentiveDescriptionListDataType{
		IncentiveDescriptionData: []IncentiveDescriptionDataType{
			{
				IncentiveId: util.Ptr(IncentiveIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				IncentiveId: util.Ptr(IncentiveIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := IncentiveDescriptionListDataType{
		IncentiveDescriptionData: []IncentiveDescriptionDataType{
			{
				IncentiveId: util.Ptr(IncentiveIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.IncentiveDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.IncentiveId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.IncentiveId))
	assert.Equal(t, "new", string(*item2.Description))
}

// ReadPartialData tests

func TestTariffListDataType_ReadPartialData(t *testing.T) {
	sut := &TariffListDataType{
		TariffData: []TariffDataType{
			{
				TariffId:     util.Ptr(TariffIdType(1)),
				ActiveTierId: []TierIdType{1, 2},
			},
			{
				TariffId:     util.Ptr(TariffIdType(2)),
				ActiveTierId: []TierIdType{3},
			},
		},
	}

	// Test 1: Testing nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, sut, result)

	// Test 2: Testing filter with selector should return filtered data
	filter := FilterType{
		TariffListDataSelectors: &TariffListDataSelectorsType{
			TariffId: util.Ptr(TariffIdType(1)),
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult := result.(*TariffListDataType)
	assert.Equal(t, 1, len(filteredResult.TariffData))
	assert.Equal(t, TariffIdType(1), *filteredResult.TariffData[0].TariffId)

	// Test 3: Testing filter with elements should return filtered data
	filter = FilterType{
		TariffDataElements: &TariffDataElementsType{
			TariffId: &ElementTagType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TariffListDataType)
	assert.Equal(t, 2, len(filteredResult.TariffData))
	assert.NotNil(t, filteredResult.TariffData[0].TariffId)
	assert.Nil(t, filteredResult.TariffData[0].ActiveTierId) // Should be nil due to elements filtering

	// Test 4: Testing filter with selectors and elements should return filtered data
	filter = FilterType{
		TariffListDataSelectors: &TariffListDataSelectorsType{
			TariffId: util.Ptr(TariffIdType(2)),
		},
		TariffDataElements: &TariffDataElementsType{
			TariffId: &ElementTagType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TariffListDataType)
	assert.Equal(t, 1, len(filteredResult.TariffData))
	assert.Equal(t, TariffIdType(2), *filteredResult.TariffData[0].TariffId)
	assert.Nil(t, filteredResult.TariffData[0].ActiveTierId) // Should be nil due to elements filtering

	// Test 5: Testing non-matching selectors should return empty data
	filter = FilterType{
		TariffListDataSelectors: &TariffListDataSelectorsType{
			TariffId: util.Ptr(TariffIdType(99)), // Non-existing ID
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TariffListDataType)
	assert.Equal(t, 0, len(filteredResult.TariffData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestTariffTierRelationListDataType_ReadPartialData(t *testing.T) {
	sut := &TariffTierRelationListDataType{
		TariffTierRelationData: []TariffTierRelationDataType{
			{
				TariffId: util.Ptr(TariffIdType(1)),
				TierId:   []TierIdType{1, 2},
			},
			{
				TariffId: util.Ptr(TariffIdType(2)),
				TierId:   []TierIdType{3},
			},
		},
	}

	// Test 1: Testing nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, sut, result)

	// Test 2: Testing filter with selector should return filtered data
	filter := FilterType{
		TariffTierRelationListDataSelectors: &TariffTierRelationListDataSelectorsType{
			TariffId: util.Ptr(TariffIdType(1)),
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult := result.(*TariffTierRelationListDataType)
	assert.Equal(t, 1, len(filteredResult.TariffTierRelationData))
	assert.Equal(t, TariffIdType(1), *filteredResult.TariffTierRelationData[0].TariffId)

	// Test 3: Testing filter with elements should return filtered data
	filter = FilterType{
		TariffTierRelationDataElements: &TariffTierRelationDataElementsType{
			TariffId: &ElementTagType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TariffTierRelationListDataType)
	assert.Equal(t, 2, len(filteredResult.TariffTierRelationData))
	assert.NotNil(t, filteredResult.TariffTierRelationData[0].TariffId)
	assert.Nil(t, filteredResult.TariffTierRelationData[0].TierId) // Should be nil due to elements filtering

	// Test 4: Testing filter with selectors and elements should return filtered data
	filter = FilterType{
		TariffTierRelationListDataSelectors: &TariffTierRelationListDataSelectorsType{
			TariffId: util.Ptr(TariffIdType(2)),
		},
		TariffTierRelationDataElements: &TariffTierRelationDataElementsType{
			TariffId: &ElementTagType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TariffTierRelationListDataType)
	assert.Equal(t, 1, len(filteredResult.TariffTierRelationData))
	assert.Equal(t, TariffIdType(2), *filteredResult.TariffTierRelationData[0].TariffId)
	assert.Nil(t, filteredResult.TariffTierRelationData[0].TierId) // Should be nil due to elements filtering

	// Test 5: Testing non-matching selectors should return empty data
	filter = FilterType{
		TariffTierRelationListDataSelectors: &TariffTierRelationListDataSelectorsType{
			TariffId: util.Ptr(TariffIdType(99)), // Non-existing ID
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TariffTierRelationListDataType)
	assert.Equal(t, 0, len(filteredResult.TariffTierRelationData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestTariffBoundaryRelationListDataType_ReadPartialData(t *testing.T) {
	sut := &TariffBoundaryRelationListDataType{
		TariffBoundaryRelationData: []TariffBoundaryRelationDataType{
			{
				TariffId:   util.Ptr(TariffIdType(1)),
				BoundaryId: []TierBoundaryIdType{1, 2},
			},
			{
				TariffId:   util.Ptr(TariffIdType(2)),
				BoundaryId: []TierBoundaryIdType{3},
			},
		},
	}

	// Test 1: Testing nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, sut, result)

	// Test 2: Testing filter with selector should return filtered data
	filter := FilterType{
		TariffBoundaryRelationListDataSelectors: &TariffBoundaryRelationListDataSelectorsType{
			TariffId: util.Ptr(TariffIdType(1)),
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult := result.(*TariffBoundaryRelationListDataType)
	assert.Equal(t, 1, len(filteredResult.TariffBoundaryRelationData))
	assert.Equal(t, TariffIdType(1), *filteredResult.TariffBoundaryRelationData[0].TariffId)

	// Test 3: Testing filter with elements should return filtered data
	filter = FilterType{
		TariffBoundaryRelationDataElements: &TariffBoundaryRelationDataElementsType{
			TariffId: &ElementTagType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TariffBoundaryRelationListDataType)
	assert.Equal(t, 2, len(filteredResult.TariffBoundaryRelationData))
	assert.NotNil(t, filteredResult.TariffBoundaryRelationData[0].TariffId)
	assert.Nil(t, filteredResult.TariffBoundaryRelationData[0].BoundaryId) // Should be nil due to elements filtering

	// Test 4: Testing filter with selectors and elements should return filtered data
	filter = FilterType{
		TariffBoundaryRelationListDataSelectors: &TariffBoundaryRelationListDataSelectorsType{
			TariffId: util.Ptr(TariffIdType(2)),
		},
		TariffBoundaryRelationDataElements: &TariffBoundaryRelationDataElementsType{
			TariffId: &ElementTagType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TariffBoundaryRelationListDataType)
	assert.Equal(t, 1, len(filteredResult.TariffBoundaryRelationData))
	assert.Equal(t, TariffIdType(2), *filteredResult.TariffBoundaryRelationData[0].TariffId)
	assert.Nil(t, filteredResult.TariffBoundaryRelationData[0].BoundaryId) // Should be nil due to elements filtering

	// Test 5: Testing non-matching selectors should return empty data
	filter = FilterType{
		TariffBoundaryRelationListDataSelectors: &TariffBoundaryRelationListDataSelectorsType{
			TariffId: util.Ptr(TariffIdType(99)), // Non-existing ID
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TariffBoundaryRelationListDataType)
	assert.Equal(t, 0, len(filteredResult.TariffBoundaryRelationData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestTariffDescriptionListDataType_ReadPartialData(t *testing.T) {
	sut := &TariffDescriptionListDataType{
		TariffDescriptionData: []TariffDescriptionDataType{
			{
				TariffId:    util.Ptr(TariffIdType(1)),
				Description: util.Ptr(DescriptionType("Tariff 1")),
			},
			{
				TariffId:    util.Ptr(TariffIdType(2)),
				Description: util.Ptr(DescriptionType("Tariff 2")),
			},
		},
	}

	// Test 1: Testing nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, sut, result)

	// Test 2: Testing filter with selector should return filtered data
	filter := FilterType{
		TariffDescriptionListDataSelectors: &TariffDescriptionListDataSelectorsType{
			TariffId: util.Ptr(TariffIdType(1)),
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult := result.(*TariffDescriptionListDataType)
	assert.Equal(t, 1, len(filteredResult.TariffDescriptionData))
	assert.Equal(t, TariffIdType(1), *filteredResult.TariffDescriptionData[0].TariffId)

	// Test 3: Testing filter with elements should return filtered data
	filter = FilterType{
		TariffDescriptionDataElements: &TariffDescriptionDataElementsType{
			TariffId: &ElementTagType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TariffDescriptionListDataType)
	assert.Equal(t, 2, len(filteredResult.TariffDescriptionData))
	assert.NotNil(t, filteredResult.TariffDescriptionData[0].TariffId)
	assert.Nil(t, filteredResult.TariffDescriptionData[0].Description) // Should be nil due to elements filtering

	// Test 4: Testing filter with selectors and elements should return filtered data
	filter = FilterType{
		TariffDescriptionListDataSelectors: &TariffDescriptionListDataSelectorsType{
			TariffId: util.Ptr(TariffIdType(2)),
		},
		TariffDescriptionDataElements: &TariffDescriptionDataElementsType{
			Description: &ElementTagType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TariffDescriptionListDataType)
	assert.Equal(t, 1, len(filteredResult.TariffDescriptionData))
	assert.Equal(t, "Tariff 2", string(*filteredResult.TariffDescriptionData[0].Description))
	assert.Nil(t, filteredResult.TariffDescriptionData[0].TariffId) // Should be nil due to elements filtering

	// Test 5: Testing non-matching selectors should return empty data
	filter = FilterType{
		TariffDescriptionListDataSelectors: &TariffDescriptionListDataSelectorsType{
			TariffId: util.Ptr(TariffIdType(99)), // Non-existing ID
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TariffDescriptionListDataType)
	assert.Equal(t, 0, len(filteredResult.TariffDescriptionData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestTierBoundaryListDataType_ReadPartialData(t *testing.T) {
	sut := &TierBoundaryListDataType{
		TierBoundaryData: []TierBoundaryDataType{
			{
				BoundaryId:         util.Ptr(TierBoundaryIdType(1)),
				LowerBoundaryValue: NewScaledNumberType(10),
				UpperBoundaryValue: NewScaledNumberType(20),
			},
			{
				BoundaryId:         util.Ptr(TierBoundaryIdType(2)),
				LowerBoundaryValue: NewScaledNumberType(30),
				UpperBoundaryValue: NewScaledNumberType(40),
			},
		},
	}

	// Test 1: Testing nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, sut, result)

	// Test 2: Testing filter with selector should return filtered data
	filter := FilterType{
		TierBoundaryListDataSelectors: &TierBoundaryListDataSelectorsType{
			BoundaryId: util.Ptr(TierBoundaryIdType(1)),
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult := result.(*TierBoundaryListDataType)
	assert.Equal(t, 1, len(filteredResult.TierBoundaryData))
	assert.Equal(t, TierBoundaryIdType(1), *filteredResult.TierBoundaryData[0].BoundaryId)

	// Test 3: Testing filter with elements should return filtered data
	filter = FilterType{
		TierBoundaryDataElements: &TierBoundaryDataElementsType{
			BoundaryId:         &ElementTagType{},
			LowerBoundaryValue: &ScaledNumberElementsType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TierBoundaryListDataType)
	assert.Equal(t, 2, len(filteredResult.TierBoundaryData))
	assert.NotNil(t, filteredResult.TierBoundaryData[0].BoundaryId)
	assert.NotNil(t, filteredResult.TierBoundaryData[0].LowerBoundaryValue)
	assert.Nil(t, filteredResult.TierBoundaryData[0].UpperBoundaryValue) // Should be nil due to elements filtering

	// Test 4: Testing filter with selectors and elements should return filtered data
	filter = FilterType{
		TierBoundaryListDataSelectors: &TierBoundaryListDataSelectorsType{
			BoundaryId: util.Ptr(TierBoundaryIdType(2)),
		},
		TierBoundaryDataElements: &TierBoundaryDataElementsType{
			UpperBoundaryValue: &ScaledNumberElementsType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TierBoundaryListDataType)
	assert.Equal(t, 1, len(filteredResult.TierBoundaryData))
	assert.Equal(t, 40.0, filteredResult.TierBoundaryData[0].UpperBoundaryValue.GetValue())
	assert.Nil(t, filteredResult.TierBoundaryData[0].BoundaryId)         // Should be nil due to elements filtering
	assert.Nil(t, filteredResult.TierBoundaryData[0].LowerBoundaryValue) // Should be nil due to elements filtering

	// Test 5: Testing non-matching selectors should return empty data
	filter = FilterType{
		TierBoundaryListDataSelectors: &TierBoundaryListDataSelectorsType{
			BoundaryId: util.Ptr(TierBoundaryIdType(99)), // Non-existing ID
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TierBoundaryListDataType)
	assert.Equal(t, 0, len(filteredResult.TierBoundaryData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestTierBoundaryDescriptionListDataType_ReadPartialData(t *testing.T) {
	sut := &TierBoundaryDescriptionListDataType{
		TierBoundaryDescriptionData: []TierBoundaryDescriptionDataType{
			{
				BoundaryId:  util.Ptr(TierBoundaryIdType(1)),
				Description: util.Ptr(DescriptionType("Boundary 1")),
			},
			{
				BoundaryId:  util.Ptr(TierBoundaryIdType(2)),
				Description: util.Ptr(DescriptionType("Boundary 2")),
			},
		},
	}

	// Test 1: Testing nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, sut, result)

	// Test 2: Testing filter with selector should return filtered data
	filter := FilterType{
		TierBoundaryDescriptionListDataSelectors: &TierBoundaryDescriptionListDataSelectorsType{
			BoundaryId: util.Ptr(TierBoundaryIdType(1)),
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult := result.(*TierBoundaryDescriptionListDataType)
	assert.Equal(t, 1, len(filteredResult.TierBoundaryDescriptionData))
	assert.Equal(t, TierBoundaryIdType(1), *filteredResult.TierBoundaryDescriptionData[0].BoundaryId)

	// Test 3: Testing filter with elements should return filtered data
	filter = FilterType{
		TierBoundaryDescriptionDataElements: &TierBoundaryDescriptionDataElementsType{
			BoundaryId: &ElementTagType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TierBoundaryDescriptionListDataType)
	assert.Equal(t, 2, len(filteredResult.TierBoundaryDescriptionData))
	assert.NotNil(t, filteredResult.TierBoundaryDescriptionData[0].BoundaryId)
	assert.Nil(t, filteredResult.TierBoundaryDescriptionData[0].Description)

	// Test 4: Testing filter with selectors and elements should return filtered data
	filter = FilterType{
		TierBoundaryDescriptionListDataSelectors: &TierBoundaryDescriptionListDataSelectorsType{
			BoundaryId: util.Ptr(TierBoundaryIdType(2)),
		},
		TierBoundaryDescriptionDataElements: &TierBoundaryDescriptionDataElementsType{
			Description: &ElementTagType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TierBoundaryDescriptionListDataType)
	assert.Equal(t, 1, len(filteredResult.TierBoundaryDescriptionData))
	assert.Equal(t, "Boundary 2", string(*filteredResult.TierBoundaryDescriptionData[0].Description))
	assert.Nil(t, filteredResult.TierBoundaryDescriptionData[0].BoundaryId)

	// Test 5: Testing non-matching selectors should return empty data
	filter = FilterType{
		TierBoundaryDescriptionListDataSelectors: &TierBoundaryDescriptionListDataSelectorsType{
			BoundaryId: util.Ptr(TierBoundaryIdType(99)),
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TierBoundaryDescriptionListDataType)
	assert.Equal(t, 0, len(filteredResult.TierBoundaryDescriptionData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestCommodityListDataType_ReadPartialData(t *testing.T) {
	sut := &CommodityListDataType{
		CommodityData: []CommodityDataType{
			{
				CommodityId: util.Ptr(CommodityIdType(1)),
				Description: util.Ptr(DescriptionType("Commodity 1")),
			},
			{
				CommodityId: util.Ptr(CommodityIdType(2)),
				Description: util.Ptr(DescriptionType("Commodity 2")),
			},
		},
	}

	// Test 1: Testing nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, sut, result)

	// Test 2: Testing filter with selector should return filtered data
	filter := FilterType{
		CommodityListDataSelectors: &CommodityListDataSelectorsType{
			CommodityId: util.Ptr(CommodityIdType(1)),
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult := result.(*CommodityListDataType)
	assert.Equal(t, 1, len(filteredResult.CommodityData))
	assert.Equal(t, CommodityIdType(1), *filteredResult.CommodityData[0].CommodityId)

	// Test 3: Testing filter with elements should return filtered data
	filter = FilterType{
		CommodityDataElements: &CommodityDataElementsType{
			CommodityId: &ElementTagType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*CommodityListDataType)
	assert.Equal(t, 2, len(filteredResult.CommodityData))
	assert.NotNil(t, filteredResult.CommodityData[0].CommodityId)
	assert.Nil(t, filteredResult.CommodityData[0].Description)

	// Test 4: Testing filter with selectors and elements should return filtered data
	filter = FilterType{
		CommodityListDataSelectors: &CommodityListDataSelectorsType{
			CommodityId: util.Ptr(CommodityIdType(2)),
		},
		CommodityDataElements: &CommodityDataElementsType{
			Description: &ElementTagType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*CommodityListDataType)
	assert.Equal(t, 1, len(filteredResult.CommodityData))
	assert.Equal(t, "Commodity 2", string(*filteredResult.CommodityData[0].Description))
	assert.Nil(t, filteredResult.CommodityData[0].CommodityId)

	// Test 5: Testing non-matching selectors should return empty data
	filter = FilterType{
		CommodityListDataSelectors: &CommodityListDataSelectorsType{
			CommodityId: util.Ptr(CommodityIdType(99)),
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*CommodityListDataType)
	assert.Equal(t, 0, len(filteredResult.CommodityData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestTierListDataType_ReadPartialData(t *testing.T) {
	sut := &TierListDataType{
		TierData: []TierDataType{
			{
				TierId:            util.Ptr(TierIdType(1)),
				ActiveIncentiveId: []IncentiveIdType{1, 2},
			},
			{
				TierId:            util.Ptr(TierIdType(2)),
				ActiveIncentiveId: []IncentiveIdType{3},
			},
		},
	}

	// Test 1: Testing nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, sut, result)

	// Test 2: Testing filter with selector should return filtered data
	filter := FilterType{
		TierListDataSelectors: &TierListDataSelectorsType{
			TierId: util.Ptr(TierIdType(1)),
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult := result.(*TierListDataType)
	assert.Equal(t, 1, len(filteredResult.TierData))
	assert.Equal(t, TierIdType(1), *filteredResult.TierData[0].TierId)

	// Test 3: Testing filter with elements should return filtered data
	filter = FilterType{
		TierDataElements: &TierDataElementsType{
			TierId: &ElementTagType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TierListDataType)
	assert.Equal(t, 2, len(filteredResult.TierData))
	assert.NotNil(t, filteredResult.TierData[0].TierId)
	assert.Nil(t, filteredResult.TierData[0].ActiveIncentiveId)

	// Test 4: Testing filter with selectors and elements should return filtered data
	filter = FilterType{
		TierListDataSelectors: &TierListDataSelectorsType{
			TierId: util.Ptr(TierIdType(2)),
		},
		TierDataElements: &TierDataElementsType{
			ActiveIncentiveId: &ElementTagType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TierListDataType)
	assert.Equal(t, 1, len(filteredResult.TierData))
	assert.Equal(t, 1, len(filteredResult.TierData[0].ActiveIncentiveId))
	assert.Equal(t, IncentiveIdType(3), filteredResult.TierData[0].ActiveIncentiveId[0])
	assert.Nil(t, filteredResult.TierData[0].TierId)

	// Test 5: Testing non-matching selectors should return empty data
	filter = FilterType{
		TierListDataSelectors: &TierListDataSelectorsType{
			TierId: util.Ptr(TierIdType(99)),
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TierListDataType)
	assert.Equal(t, 0, len(filteredResult.TierData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestTierIncentiveRelationListDataType_ReadPartialData(t *testing.T) {
	sut := &TierIncentiveRelationListDataType{
		TierIncentiveRelationData: []TierIncentiveRelationDataType{
			{
				TierId:      util.Ptr(TierIdType(1)),
				IncentiveId: []IncentiveIdType{1, 2},
			},
			{
				TierId:      util.Ptr(TierIdType(2)),
				IncentiveId: []IncentiveIdType{3},
			},
		},
	}

	// Test 1: Testing nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, sut, result)

	// Test 2: Testing filter with selector should return filtered data
	filter := FilterType{
		TierIncentiveRelationListDataSelectors: &TierIncentiveRelationListDataSelectorsType{
			TierId: util.Ptr(TierIdType(1)),
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult := result.(*TierIncentiveRelationListDataType)
	assert.Equal(t, 1, len(filteredResult.TierIncentiveRelationData))
	assert.Equal(t, TierIdType(1), *filteredResult.TierIncentiveRelationData[0].TierId)

	// Test 3: Testing filter with elements should return filtered data
	filter = FilterType{
		TierIncentiveRelationDataElements: &TierIncentiveRelationDataElementsType{
			TierId: &ElementTagType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TierIncentiveRelationListDataType)
	assert.Equal(t, 2, len(filteredResult.TierIncentiveRelationData))
	assert.NotNil(t, filteredResult.TierIncentiveRelationData[0].TierId)
	assert.Nil(t, filteredResult.TierIncentiveRelationData[0].IncentiveId)

	// Test 4: Testing filter with selectors and elements should return filtered data
	filter = FilterType{
		TierIncentiveRelationListDataSelectors: &TierIncentiveRelationListDataSelectorsType{
			TierId: util.Ptr(TierIdType(2)),
		},
		TierIncentiveRelationDataElements: &TierIncentiveRelationDataElementsType{
			IncentiveId: &ElementTagType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TierIncentiveRelationListDataType)
	assert.Equal(t, 1, len(filteredResult.TierIncentiveRelationData))
	assert.Equal(t, 1, len(filteredResult.TierIncentiveRelationData[0].IncentiveId))
	assert.Equal(t, IncentiveIdType(3), filteredResult.TierIncentiveRelationData[0].IncentiveId[0])
	assert.Nil(t, filteredResult.TierIncentiveRelationData[0].TierId)

	// Test 5: Testing non-matching selectors should return empty data
	filter = FilterType{
		TierIncentiveRelationListDataSelectors: &TierIncentiveRelationListDataSelectorsType{
			TierId: util.Ptr(TierIdType(99)),
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TierIncentiveRelationListDataType)
	assert.Equal(t, 0, len(filteredResult.TierIncentiveRelationData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestTierDescriptionListDataType_ReadPartialData(t *testing.T) {
	sut := &TierDescriptionListDataType{
		TierDescriptionData: []TierDescriptionDataType{
			{
				TierId:      util.Ptr(TierIdType(1)),
				Description: util.Ptr(DescriptionType("Tier 1")),
			},
			{
				TierId:      util.Ptr(TierIdType(2)),
				Description: util.Ptr(DescriptionType("Tier 2")),
			},
		},
	}

	// Test 1: Testing nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, sut, result)

	// Test 2: Testing filter with selector should return filtered data
	filter := FilterType{
		TierDescriptionListDataSelectors: &TierDescriptionListDataSelectorsType{
			TierId: util.Ptr(TierIdType(1)),
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult := result.(*TierDescriptionListDataType)
	assert.Equal(t, 1, len(filteredResult.TierDescriptionData))
	assert.Equal(t, TierIdType(1), *filteredResult.TierDescriptionData[0].TierId)

	// Test 3: Testing filter with elements should return filtered data
	filter = FilterType{
		TierDescriptionDataElements: &TierDescriptionDataElementsType{
			TierId: &ElementTagType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TierDescriptionListDataType)
	assert.Equal(t, 2, len(filteredResult.TierDescriptionData))
	assert.NotNil(t, filteredResult.TierDescriptionData[0].TierId)
	assert.Nil(t, filteredResult.TierDescriptionData[0].Description)

	// Test 4: Testing filter with selectors and elements should return filtered data
	filter = FilterType{
		TierDescriptionListDataSelectors: &TierDescriptionListDataSelectorsType{
			TierId: util.Ptr(TierIdType(2)),
		},
		TierDescriptionDataElements: &TierDescriptionDataElementsType{
			Description: &ElementTagType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TierDescriptionListDataType)
	assert.Equal(t, 1, len(filteredResult.TierDescriptionData))
	assert.Equal(t, "Tier 2", string(*filteredResult.TierDescriptionData[0].Description))
	assert.Nil(t, filteredResult.TierDescriptionData[0].TierId)

	// Test 5: Testing non-matching selectors should return empty data
	filter = FilterType{
		TierDescriptionListDataSelectors: &TierDescriptionListDataSelectorsType{
			TierId: util.Ptr(TierIdType(99)),
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*TierDescriptionListDataType)
	assert.Equal(t, 0, len(filteredResult.TierDescriptionData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestIncentiveListDataType_ReadPartialData(t *testing.T) {
	sut := &IncentiveListDataType{
		IncentiveData: []IncentiveDataType{
			{
				IncentiveId: util.Ptr(IncentiveIdType(1)),
				Value:       NewScaledNumberType(10.5),
			},
			{
				IncentiveId: util.Ptr(IncentiveIdType(2)),
				Value:       NewScaledNumberType(20.7),
			},
		},
	}

	// Test 1: Testing nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, sut, result)

	// Test 2: Testing filter with selector should return filtered data
	filter := FilterType{
		IncentiveListDataSelectors: &IncentiveListDataSelectorsType{
			IncentiveId: util.Ptr(IncentiveIdType(1)),
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult := result.(*IncentiveListDataType)
	assert.Equal(t, 1, len(filteredResult.IncentiveData))
	assert.Equal(t, IncentiveIdType(1), *filteredResult.IncentiveData[0].IncentiveId)

	// Test 3: Testing filter with elements should return filtered data
	filter = FilterType{
		IncentiveDataElements: &IncentiveDataElementsType{
			IncentiveId: &ElementTagType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*IncentiveListDataType)
	assert.Equal(t, 2, len(filteredResult.IncentiveData))
	assert.NotNil(t, filteredResult.IncentiveData[0].IncentiveId)
	assert.Nil(t, filteredResult.IncentiveData[0].Value)

	// Test 4: Testing filter with selectors and elements should return filtered data
	filter = FilterType{
		IncentiveListDataSelectors: &IncentiveListDataSelectorsType{
			IncentiveId: util.Ptr(IncentiveIdType(2)),
		},
		IncentiveDataElements: &IncentiveDataElementsType{
			Value: &ScaledNumberElementsType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*IncentiveListDataType)
	assert.Equal(t, 1, len(filteredResult.IncentiveData))
	assert.InDelta(t, 20.7, filteredResult.IncentiveData[0].Value.GetValue(), 0.001)
	assert.Nil(t, filteredResult.IncentiveData[0].IncentiveId)

	// Test 5: Testing non-matching selectors should return empty data
	filter = FilterType{
		IncentiveListDataSelectors: &IncentiveListDataSelectorsType{
			IncentiveId: util.Ptr(IncentiveIdType(99)),
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*IncentiveListDataType)
	assert.Equal(t, 0, len(filteredResult.IncentiveData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestIncentiveDescriptionListDataType_ReadPartialData(t *testing.T) {
	sut := &IncentiveDescriptionListDataType{
		IncentiveDescriptionData: []IncentiveDescriptionDataType{
			{
				IncentiveId: util.Ptr(IncentiveIdType(1)),
				Description: util.Ptr(DescriptionType("Incentive 1")),
			},
			{
				IncentiveId: util.Ptr(IncentiveIdType(2)),
				Description: util.Ptr(DescriptionType("Incentive 2")),
			},
		},
	}

	// Test 1: Testing nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, sut, result)

	// Test 2: Testing filter with selector should return filtered data
	filter := FilterType{
		IncentiveDescriptionListDataSelectors: &IncentiveDescriptionListDataSelectorsType{
			IncentiveId: util.Ptr(IncentiveIdType(1)),
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult := result.(*IncentiveDescriptionListDataType)
	assert.Equal(t, 1, len(filteredResult.IncentiveDescriptionData))
	assert.Equal(t, IncentiveIdType(1), *filteredResult.IncentiveDescriptionData[0].IncentiveId)

	// Test 3: Testing filter with elements should return filtered data
	filter = FilterType{
		IncentiveDescriptionDataElements: &IncentiveDescriptionDataElementsType{
			IncentiveId: &ElementTagType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*IncentiveDescriptionListDataType)
	assert.Equal(t, 2, len(filteredResult.IncentiveDescriptionData))
	assert.NotNil(t, filteredResult.IncentiveDescriptionData[0].IncentiveId)
	assert.Nil(t, filteredResult.IncentiveDescriptionData[0].Description)

	// Test 4: Testing filter with selectors and elements should return filtered data
	filter = FilterType{
		IncentiveDescriptionListDataSelectors: &IncentiveDescriptionListDataSelectorsType{
			IncentiveId: util.Ptr(IncentiveIdType(2)),
		},
		IncentiveDescriptionDataElements: &IncentiveDescriptionDataElementsType{
			Description: &ElementTagType{},
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*IncentiveDescriptionListDataType)
	assert.Equal(t, 1, len(filteredResult.IncentiveDescriptionData))
	assert.Equal(t, "Incentive 2", string(*filteredResult.IncentiveDescriptionData[0].Description))
	assert.Nil(t, filteredResult.IncentiveDescriptionData[0].IncentiveId)

	// Test 5: Testing non-matching selectors should return empty data
	filter = FilterType{
		IncentiveDescriptionListDataSelectors: &IncentiveDescriptionListDataSelectorsType{
			IncentiveId: util.Ptr(IncentiveIdType(99)),
		},
	}

	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	filteredResult = result.(*IncentiveDescriptionListDataType)
	assert.Equal(t, 0, len(filteredResult.IncentiveDescriptionData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}
