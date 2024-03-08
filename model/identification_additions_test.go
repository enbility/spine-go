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
	sut.UpdateList(false, &newData, NewFilterTypePartial(), nil)

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
	sut.UpdateList(false, &newData, NewFilterTypePartial(), nil)

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
	sut.UpdateList(false, &newData, NewFilterTypePartial(), nil)

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
