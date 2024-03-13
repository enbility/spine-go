package spine

import (
	"testing"

	"github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestEntityRemoteSuite(t *testing.T) {
	suite.Run(t, new(EntityRemoteTestSuite))
}

type EntityRemoteTestSuite struct {
	suite.Suite
}

func (s *EntityRemoteTestSuite) Test_Entity() {
	deviceAddress := model.AddressDeviceType("")
	device := mocks.NewDeviceRemoteInterface(s.T())
	device.EXPECT().Address().Return(&deviceAddress)
	device.EXPECT().AddEntity(mock.Anything)

	entity := NewEntityRemote(device, model.EntityTypeTypeCEM, NewAddressEntityType([]uint{1}))
	device.AddEntity(entity)

	assert.Equal(s.T(), device, entity.Device())

	entity.UpdateDeviceAddress(model.AddressDeviceType("new"))

	f := NewFeatureRemote(1, entity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient)
	entity.AddFeature(f)
	assert.Equal(s.T(), 1, len(entity.Features()))

	f1 := entity.FeatureOfAddress(nil)
	assert.Nil(s.T(), f1)

	f1 = entity.FeatureOfAddress(f.Address().Feature)
	assert.NotNil(s.T(), f1)

	fakeAddress := model.AddressFeatureType(5)
	f1 = entity.FeatureOfAddress(&fakeAddress)
	assert.Nil(s.T(), f1)

	f1 = entity.FeatureOfTypeAndRole(model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient)
	assert.NotNil(s.T(), f1)

	f1 = entity.FeatureOfTypeAndRole(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	assert.Nil(s.T(), f1)

	entity.RemoveAllFeatures()
	assert.Equal(s.T(), 0, len(entity.Features()))
}
