package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestIdentificationListDataType_Update(t *testing.T) {
	sut := IdentificationListDataType{
		IdentificationData: []IdentificationDataType{
			{
				IdentificationId:   util.Ptr(IdentificationIdType(0)),
				IdentificationType: util.Ptr(IdentificationTypeTypeEui48),
			},
			{
				IdentificationId:   util.Ptr(IdentificationIdType(1)),
				IdentificationType: util.Ptr(IdentificationTypeTypeEui48),
			},
		},
	}

	newData := IdentificationListDataType{
		IdentificationData: []IdentificationDataType{
			{
				IdentificationId:   util.Ptr(IdentificationIdType(1)),
				IdentificationType: util.Ptr(IdentificationTypeTypeEui64),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.IdentificationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.IdentificationId))
	assert.Equal(t, IdentificationTypeTypeEui48, *item1.IdentificationType)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.IdentificationId))
	assert.Equal(t, IdentificationTypeTypeEui64, *item2.IdentificationType)
}

func TestSessionIdentificationListDataType_Update(t *testing.T) {
	sut := SessionIdentificationListDataType{
		SessionIdentificationData: []SessionIdentificationDataType{
			{
				IdentificationId: util.Ptr(IdentificationIdType(0)),
				SessionId:        util.Ptr(SessionIdType(1)),
				IsLatestSession:  util.Ptr(false),
			},
			{
				IdentificationId: util.Ptr(IdentificationIdType(1)),
				SessionId:        util.Ptr(SessionIdType(2)),
				IsLatestSession:  util.Ptr(true),
			},
		},
	}

	newData := SessionIdentificationListDataType{
		SessionIdentificationData: []SessionIdentificationDataType{
			{
				IdentificationId: util.Ptr(IdentificationIdType(1)),
				SessionId:        util.Ptr(SessionIdType(2)),
				IsLatestSession:  util.Ptr(false),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.SessionIdentificationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.IdentificationId))
	assert.Equal(t, false, *item1.IsLatestSession)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.IdentificationId))
	assert.Equal(t, false, *item2.IsLatestSession)
}

func TestSessionMeasurementRelationListDataType_Update(t *testing.T) {
	sut := SessionMeasurementRelationListDataType{
		SessionMeasurementRelationData: []SessionMeasurementRelationDataType{
			{
				SessionId: util.Ptr(SessionIdType(0)),
				MeasurementId: []MeasurementIdType{
					0, 1,
				},
			},
			{
				SessionId: util.Ptr(SessionIdType(1)),
				MeasurementId: []MeasurementIdType{
					2, 3,
				},
			},
		},
	}

	newData := SessionMeasurementRelationListDataType{
		SessionMeasurementRelationData: []SessionMeasurementRelationDataType{
			{
				SessionId: util.Ptr(SessionIdType(1)),
				MeasurementId: []MeasurementIdType{
					2, 3, 4,
				},
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.SessionMeasurementRelationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SessionId))
	assert.Equal(t, []MeasurementIdType{0, 1}, item1.MeasurementId)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SessionId))
	assert.Equal(t, []MeasurementIdType{2, 3, 4}, item2.MeasurementId)
}

func TestSessionIdentificationListDataType_ReadPartialData(t *testing.T) {
	sessionId := SessionIdType(123)
	data := &SessionIdentificationListDataType{
		SessionIdentificationData: []SessionIdentificationDataType{
			{SessionId: &sessionId},
			{SessionId: Ptr(SessionIdType(456))},
		},
	}

	// Test without filter - should return all data
	result, ok := data.ReadPartialData(nil)
	assert.True(t, ok)
	assert.Equal(t, data, result)

	// Test with session ID filter
	filter := &FilterType{
		SessionIdentificationListDataSelectors: &SessionIdentificationListDataSelectorsType{
			SessionId: &sessionId,
		},
	}

	result, ok = data.ReadPartialData(filter)
	assert.True(t, ok)

	resultData := result.(*SessionIdentificationListDataType)
	assert.Len(t, resultData.SessionIdentificationData, 1)
	assert.Equal(t, sessionId, *resultData.SessionIdentificationData[0].SessionId)

	// Test 3: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success := data.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestSessionMeasurementRelationListDataType_ReadPartialData(t *testing.T) {
	sessionId := SessionIdType(123)
	measurementId1 := MeasurementIdType(1)
	measurementId2 := MeasurementIdType(2)

	data := &SessionMeasurementRelationListDataType{
		SessionMeasurementRelationData: []SessionMeasurementRelationDataType{
			{
				SessionId:     &sessionId,
				MeasurementId: []MeasurementIdType{measurementId1, measurementId2},
			},
			{
				SessionId:     Ptr(SessionIdType(456)),
				MeasurementId: []MeasurementIdType{measurementId2},
			},
		},
	}

	// Test with measurement ID filter
	filter := &FilterType{
		SessionMeasurementRelationListDataSelectors: &SessionMeasurementRelationListDataSelectorsType{
			MeasurementId: &measurementId1,
		},
	}

	result, ok := data.ReadPartialData(filter)
	assert.True(t, ok)

	resultData := result.(*SessionMeasurementRelationListDataType)
	assert.Len(t, resultData.SessionMeasurementRelationData, 1)
	assert.Equal(t, sessionId, *resultData.SessionMeasurementRelationData[0].SessionId)

	// Test 2: ReadPartialData with nil filter
	_, success := data.ReadPartialData(nil)
	assert.True(t, success)

	// Test 3: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = data.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestIdentificationListDataType_ReadPartialData(t *testing.T) {
	identificationId1 := IdentificationIdType(1)
	identificationId2 := IdentificationIdType(2)

	testData := &IdentificationListDataType{
		IdentificationData: []IdentificationDataType{
			{
				IdentificationId: &identificationId1,
			},
			{
				IdentificationId: &identificationId2,
			},
		},
	}

	// Test 1: ReadPartialData with nil filter should return all data
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, testData, result)

	// Test 2: ReadPartialData with filter
	filter := &FilterType{
		IdentificationListDataSelectors: &IdentificationListDataSelectorsType{
			IdentificationId: &identificationId1,
		},
	}
	_, success = testData.ReadPartialData(filter)
	assert.True(t, success)

	// Test 3: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}
