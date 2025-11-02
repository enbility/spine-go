package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestMessagingListDataType_Update(t *testing.T) {
	sut := MessagingListDataType{
		MessagingData: []MessagingDataType{
			{
				MessagingNumber: util.Ptr(MessagingNumberType(0)),
				Text:            util.Ptr(MessagingDataTextType("old")),
			},
			{
				MessagingNumber: util.Ptr(MessagingNumberType(1)),
				Text:            util.Ptr(MessagingDataTextType("old")),
			},
		},
	}

	newData := MessagingListDataType{
		MessagingData: []MessagingDataType{
			{
				MessagingNumber: util.Ptr(MessagingNumberType(1)),
				Text:            util.Ptr(MessagingDataTextType("new")),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.MessagingData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.MessagingNumber))
	assert.Equal(t, "old", string(*item1.Text))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.MessagingNumber))
	assert.Equal(t, "new", string(*item2.Text))
}

func TestMessagingListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	testData := &MessagingListDataType{
		MessagingData: []MessagingDataType{
			{
				MessagingNumber: util.Ptr(MessagingNumberType(1)),
				Text:            util.Ptr(MessagingDataTextType("Message 1")),
				MessagingType:   util.Ptr(MessagingTypeType("info")),
				Timestamp:       util.Ptr(AbsoluteOrRelativeTimeType("2023-01-01T12:00:00Z")),
			},
			{
				MessagingNumber: util.Ptr(MessagingNumberType(2)),
				Text:            util.Ptr(MessagingDataTextType("Message 2")),
				MessagingType:   util.Ptr(MessagingTypeType("warning")),
				Timestamp:       util.Ptr(AbsoluteOrRelativeTimeType("2023-01-01T13:00:00Z")),
			},
		},
	}

	// Test 1: No filter - should return all data
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*MessagingListDataType)
	assert.Len(t, resultData.MessagingData, 2)
	assert.Equal(t, testData, result)

	// Test 2: Filter with specific messaging number
	filter := &FilterType{
		MessagingListDataSelectors: &MessagingListDataSelectorsType{
			MessagingNumber: util.Ptr(MessagingNumberType(1)),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*MessagingListDataType)
	assert.Len(t, resultData.MessagingData, 1)
	assert.Equal(t, MessagingNumberType(1), *resultData.MessagingData[0].MessagingNumber)
	assert.Equal(t, "Message 1", string(*resultData.MessagingData[0].Text))

	// Test 3: Filter with elements (partial fields)
	filter = &FilterType{
		MessagingDataElements: &MessagingDataElementsType{
			MessagingNumber: &ElementTagType{},
			Text:            &ElementTagType{},
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*MessagingListDataType)
	assert.Len(t, resultData.MessagingData, 2)
	// Check that only requested elements are included
	for _, item := range resultData.MessagingData {
		assert.NotNil(t, item.MessagingNumber)
		assert.NotNil(t, item.Text)
		// MessagingType and Timestamp should be nil since not requested in elements
		assert.Nil(t, item.MessagingType)
		assert.Nil(t, item.Timestamp)
	}

	// Test 4: Combined filter with selector and elements
	filter = &FilterType{
		MessagingListDataSelectors: &MessagingListDataSelectorsType{
			MessagingNumber: util.Ptr(MessagingNumberType(2)),
		},
		MessagingDataElements: &MessagingDataElementsType{
			MessagingNumber: &ElementTagType{},
			MessagingType:   &ElementTagType{},
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*MessagingListDataType)
	assert.Len(t, resultData.MessagingData, 1)
	assert.Equal(t, MessagingNumberType(2), *resultData.MessagingData[0].MessagingNumber)
	assert.Equal(t, "warning", string(*resultData.MessagingData[0].MessagingType))
	assert.Nil(t, resultData.MessagingData[0].Text)
	assert.Nil(t, resultData.MessagingData[0].Timestamp)

	// Test 5: No matching filter
	filter = &FilterType{
		MessagingListDataSelectors: &MessagingListDataSelectorsType{
			MessagingNumber: util.Ptr(MessagingNumberType(999)),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*MessagingListDataType)
	assert.Len(t, resultData.MessagingData, 0)

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}
