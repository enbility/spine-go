package spine

import (
	"testing"
	"time"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestUtilsSuite(t *testing.T) {
	suite.Run(t, new(UtilsSuite))
}

type UtilsSuite struct {
	suite.Suite

	localDevice  api.DeviceLocalInterface
	remoteDevice api.DeviceRemoteInterface
}

func (s *UtilsSuite) WriteShipMessageWithPayload([]byte) {}

func (s *UtilsSuite) Test_DataCopyOfType() {
	s.localDevice = NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)
	localEntity := NewEntityLocal(s.localDevice, model.EntityTypeTypeCEM, NewAddressEntityType([]uint{1}))
	s.localDevice.AddEntity(localEntity)

	localFeature := NewFeatureLocal(1, localEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	localFeature.AddFunctionType(model.FunctionTypeElectricalConnectionDescriptionListData, true, false)
	localEntity.AddFeature(localFeature)
	assert.Equal(s.T(), 1, len(localEntity.Features()))

	// this function does not exist
	_, err := LocalFeatureDataCopyOfType[*model.NodeManagementUseCaseDataType](localFeature, "dummy")
	assert.NotNil(s.T(), err)

	_, err = LocalFeatureDataCopyOfType[*model.ElectricalConnectionDescriptionListDataType](
		localFeature,
		model.FunctionTypeElectricalConnectionDescriptionListData,
	)
	assert.NotNil(s.T(), err)

	data := &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{},
	}
	localFeature.updateData(model.FunctionTypeElectricalConnectionDescriptionListData, data, nil, nil)

	_, err = LocalFeatureDataCopyOfType[*model.ElectricalConnectionDescriptionListDataType](
		localFeature,
		model.FunctionTypeElectricalConnectionDescriptionListData,
	)
	assert.Nil(s.T(), err)

	_, err = LocalFeatureDataCopyOfType[string](
		localFeature,
		model.FunctionTypeElectricalConnectionDescriptionListData,
	)
	assert.NotNil(s.T(), err)

	ski := "test"
	sender := NewSender(s)

	s.remoteDevice = NewDeviceRemote(s.localDevice, ski, sender)
	remoteEntity := NewEntityRemote(s.remoteDevice, model.EntityTypeTypeCEM, NewAddressEntityType([]uint{1}))
	s.remoteDevice.AddEntity(remoteEntity)

	remoteFeature := NewFeatureRemote(1, remoteEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	remoteEntity.AddFeature(remoteFeature)

	// this function does not exist
	_, err = RemoteFeatureDataCopyOfType[*model.NodeManagementUseCaseDataType](remoteFeature, "dummy")
	assert.NotNil(s.T(), err)
}
