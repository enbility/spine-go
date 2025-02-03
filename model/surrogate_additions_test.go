package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestSurrogateDescriptionListDataType_Update(t *testing.T) {
	sut := SurrogateDescriptionListDataType{
		SurrogateDescriptionData: []SurrogateDescriptionDataType{
			{
				SurrogateId:    util.Ptr(SurrogateIdType(0)),
				SurrogateScope: util.Ptr(SurrogateScopeType("old")),
			},
			{
				SurrogateId:    util.Ptr(SurrogateIdType(1)),
				SurrogateScope: util.Ptr(SurrogateScopeType("test")),
			},
		},
	}

	newData := SurrogateDescriptionListDataType{
		SurrogateDescriptionData: []SurrogateDescriptionDataType{
			{
				SurrogateId:    util.Ptr(SurrogateIdType(1)),
				SurrogateScope: util.Ptr(SurrogateScopeType("new")),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil)
	assert.True(t, success)

	data := sut.SurrogateDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SurrogateId))
	assert.Equal(t, "old", string(*item1.SurrogateScope))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SurrogateId))
	assert.Equal(t, "new", string(*item2.SurrogateScope))
}
