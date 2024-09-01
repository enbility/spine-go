package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestVersionSuite(t *testing.T) {
	suite.Run(t, new(VersionSuite))
}

type VersionSuite struct {
	suite.Suite
}

func (s *VersionSuite) Test_UpdateList() {
	sut := SpecificationVersionListDataType{
		SpecificationVersionData: []SpecificationVersionDataType{
			SpecificationVersionDataType("1.0.0"),
		},
	}

	newData := SpecificationVersionListDataType{
		SpecificationVersionData: []SpecificationVersionDataType{
			SpecificationVersionDataType("1.0.1"),
		},
	}

	data := sut.SpecificationVersionData
	// check properties of updated item
	item1 := data[0]
	assert.Equal(s.T(), "1.0.0", string(item1))

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil)
	assert.True(s.T(), success)

	data = sut.SpecificationVersionData
	// check properties of updated item
	item1 = data[0]
	assert.Equal(s.T(), "1.0.1", string(item1))
}
