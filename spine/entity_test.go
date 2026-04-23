package spine

import (
	"encoding/json"
	"testing"

	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestEntitySuite(t *testing.T) {
	suite.Run(t, new(EntityTestSuite))
}

type EntityTestSuite struct {
	suite.Suite
}

func (s *EntityTestSuite) Test_Entity() {
	deviceAddress := model.AddressDeviceType("test")
	entity := NewEntity(model.EntityTypeTypeCEM, &deviceAddress, NewAddressEntityType([]uint{1, 1}))

	value, err := json.Marshal(entity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)
	assert.Equal(s.T(), `{"Device":"test","Entity":[1,1]}`, string(value))

	entity = NewEntity(model.EntityTypeTypeCEM, &deviceAddress, nil)

	value, err = json.Marshal(entity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)
	assert.Equal(s.T(), `{"Device":"test","Entity":null}`, string(value))

	entity = NewEntity(model.EntityTypeTypeCEM, nil, nil)

	value, err = json.Marshal(entity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)
	assert.Equal(s.T(), `{"Device":"","Entity":null}`, string(value))
}
