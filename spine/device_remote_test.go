package spine

import (
	"testing"
	"time"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	nm_usecaseinformationlistdata_recv_reply_file_path = "../spine/testdata/nm_usecaseinformationlistdata_recv_reply.json"
)

func TestDeviceRemoteSuite(t *testing.T) {
	suite.Run(t, new(DeviceRemoteSuite))
}

type DeviceRemoteSuite struct {
	suite.Suite

	localDevice  api.DeviceLocalInterface
	remoteDevice api.DeviceRemoteInterface
	remoteEntity api.EntityRemoteInterface
}

func (s *DeviceRemoteSuite) WriteShipMessageWithPayload([]byte) {}

func (s *DeviceRemoteSuite) SetupSuite() {}

func (s *DeviceRemoteSuite) BeforeTest(suiteName, testName string) {
	s.localDevice = NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)

	ski := "test"
	sender := NewSender(s)
	s.remoteDevice = NewDeviceRemote(s.localDevice, ski, sender)
	desc := &model.NetworkManagementDeviceDescriptionDataType{
		DeviceAddress: &model.DeviceAddressType{
			Device: util.Ptr(model.AddressDeviceType("test")),
		},
	}
	s.remoteDevice.UpdateDevice(desc)
	_ = s.localDevice.SetupRemoteDevice(ski, s)

	s.remoteEntity = NewEntityRemote(s.remoteDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{1})

	feature := NewFeatureRemote(0, s.remoteEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	s.remoteEntity.AddFeature(feature)

	s.remoteDevice.AddEntity(s.remoteEntity)
}

func (s *DeviceRemoteSuite) Test_RemoveByAddress() {
	assert.Equal(s.T(), 2, len(s.remoteDevice.Entities()))

	s.remoteDevice.RemoveEntityByAddress([]model.AddressEntityType{2})
	assert.Equal(s.T(), 2, len(s.remoteDevice.Entities()))

	s.remoteDevice.RemoveEntityByAddress([]model.AddressEntityType{1})
	assert.Equal(s.T(), 1, len(s.remoteDevice.Entities()))
}

func (s *DeviceRemoteSuite) Test_FeatureByEntityTypeAndRole() {
	entity := s.remoteDevice.Entity([]model.AddressEntityType{1})
	assert.NotNil(s.T(), entity)

	assert.Equal(s.T(), 1, len(entity.Features()))

	feature := s.remoteDevice.FeatureByEntityTypeAndRole(entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
	assert.Nil(s.T(), feature)

	feature = s.remoteDevice.FeatureByEntityTypeAndRole(entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	assert.NotNil(s.T(), feature)

	s.remoteDevice.RemoveEntityByAddress([]model.AddressEntityType{1})
	assert.Equal(s.T(), 1, len(s.remoteDevice.Entities()))

	_ = s.remoteDevice.Entity([]model.AddressEntityType{0})
	s.remoteDevice.RemoveEntityByAddress([]model.AddressEntityType{0})
	assert.Equal(s.T(), 0, len(s.remoteDevice.Entities()))

	feature = s.remoteDevice.FeatureByEntityTypeAndRole(entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	assert.Nil(s.T(), feature)
}

func (s *DeviceRemoteSuite) Test_Usecases() {
	uc := s.remoteDevice.UseCases()
	assert.Nil(s.T(), uc)

	_, _ = s.remoteDevice.HandleSpineMesssage(loadFileData(s.T(), nm_usecaseinformationlistdata_recv_reply_file_path))

	uc = s.remoteDevice.UseCases()
	assert.NotNil(s.T(), uc)
}
