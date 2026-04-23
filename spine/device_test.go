package spine

import (
	"encoding/json"
	"testing"

	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestDeviceSuite(t *testing.T) {
	suite.Run(t, new(DeviceTestSuite))
}

type DeviceTestSuite struct {
	suite.Suite
}

func (s *DeviceTestSuite) Test_Device() {
	deviceAddress := model.AddressDeviceType("test")
	device := NewDevice(&deviceAddress, nil, nil)

	value, err := json.Marshal(device)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)
	assert.Equal(s.T(), `"test"`, string(value))

	device = NewDevice(nil, nil, nil)

	value, err = json.Marshal(device)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)
	assert.Equal(s.T(), `""`, string(value))
}
