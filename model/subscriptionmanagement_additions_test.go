package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestSubscriptionManagementEntryListDataType_Update(t *testing.T) {
	sut := SubscriptionManagementEntryListDataType{
		SubscriptionManagementEntryData: []SubscriptionManagementEntryDataType{
			{
				SubscriptionId: util.Ptr(SubscriptionIdType(0)),
				Description:    util.Ptr(DescriptionType("old")),
			},
			{
				SubscriptionId: util.Ptr(SubscriptionIdType(1)),
				Description:    util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := SubscriptionManagementEntryListDataType{
		SubscriptionManagementEntryData: []SubscriptionManagementEntryDataType{
			{
				SubscriptionId: util.Ptr(SubscriptionIdType(1)),
				Description:    util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.SubscriptionManagementEntryData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SubscriptionId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SubscriptionId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestSubscriptionManagementEntryListDataType_ReadPartialData(t *testing.T) {
	testData := &SubscriptionManagementEntryListDataType{
		SubscriptionManagementEntryData: []SubscriptionManagementEntryDataType{
			{
				SubscriptionId: util.Ptr(SubscriptionIdType(1)),
			},
			{
				SubscriptionId: util.Ptr(SubscriptionIdType(2)),
			},
		},
	}

	// Test 1: ReadPartialData with nil filter should return all data
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, testData, result)

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		SubscriptionManagementEntryListDataSelectors: &SubscriptionManagementEntryListDataSelectorsType{
			SubscriptionId: util.Ptr(SubscriptionIdType(1)),
		},
	}
	result, success = testData.ReadPartialData(&filter)
	assert.True(t, success)
	resultData := result.(*SubscriptionManagementEntryListDataType)
	assert.Equal(t, 1, len(resultData.SubscriptionManagementEntryData))
	if len(resultData.SubscriptionManagementEntryData) > 0 {
		assert.Equal(t, SubscriptionIdType(1), *resultData.SubscriptionManagementEntryData[0].SubscriptionId)
	}

	// Test 3: filter with non-matching selector should return empty data
	filter = FilterType{
		SubscriptionManagementEntryListDataSelectors: &SubscriptionManagementEntryListDataSelectorsType{
			SubscriptionId: util.Ptr(SubscriptionIdType(3)),
		},
	}
	result, success = testData.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*SubscriptionManagementEntryListDataType)
	assert.Equal(t, 0, len(resultData.SubscriptionManagementEntryData))

	// Test 4: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}
