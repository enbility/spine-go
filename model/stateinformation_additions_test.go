package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestStateInformationListDataType_Update(t *testing.T) {
	sut := StateInformationListDataType{
		StateInformationData: []StateInformationDataType{
			{
				StateInformationId: util.Ptr(StateInformationIdType(0)),
				IsActive:           util.Ptr(true),
			},
			{
				StateInformationId: util.Ptr(StateInformationIdType(1)),
				IsActive:           util.Ptr(false),
			},
		},
	}

	newData := StateInformationListDataType{
		StateInformationData: []StateInformationDataType{
			{
				StateInformationId: util.Ptr(StateInformationIdType(1)),
				IsActive:           util.Ptr(true),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.StateInformationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.StateInformationId))
	assert.Equal(t, true, *item1.IsActive)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.StateInformationId))
	assert.Equal(t, true, *item2.IsActive)
}
