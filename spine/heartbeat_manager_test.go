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

func TestHeartbeatManagerSuite(t *testing.T) {
	suite.Run(t, new(HeartBeatManagerSuite))
}

type HeartBeatManagerSuite struct {
	suite.Suite

	localDevice api.DeviceLocalInterface
	localEntity api.EntityLocalInterface

	remoteDevice api.DeviceRemoteInterface
	sut          api.HeartbeatManagerInterface
}

func (s *HeartBeatManagerSuite) WriteShipMessageWithPayload([]byte) {}

func (s *HeartBeatManagerSuite) BeforeTest(suiteName, testName string) {
	s.localDevice = NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)
	s.localEntity = NewEntityLocal(s.localDevice, model.EntityTypeTypeCEM, []model.AddressEntityType{1}, time.Second*4)
	s.localDevice.AddEntity(s.localEntity)

	ski := "test"
	sender := NewSender(s)
	s.remoteDevice = NewDeviceRemote(s.localDevice, ski, sender)

	_ = s.localDevice.SetupRemoteDevice(ski, s)

	s.sut = s.localEntity.HeartbeatManager()
}

func (s *HeartBeatManagerSuite) Test_HeartbeatFailure() {
	s.sut.SetLocalFeature(nil, nil)

	localFeature := s.localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	s.localEntity.AddFeature(localFeature)

	s.sut.SetLocalFeature(s.localEntity, localFeature)

	running := s.sut.IsHeartbeatRunning()
	assert.Equal(s.T(), false, running)

	anotherFeature := s.localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	s.localEntity.AddFeature(anotherFeature)

	s.sut.SetLocalFeature(s.localEntity, anotherFeature)

	running = s.sut.IsHeartbeatRunning()
	assert.Equal(s.T(), false, running)

	anotherFeature.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)

	s.sut.SetLocalFeature(s.localEntity, anotherFeature)

	running = s.sut.IsHeartbeatRunning()
	assert.Equal(s.T(), true, running)
}

func (s *HeartBeatManagerSuite) Test_HeartbeatSuccess() {
	localFeature := s.localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	localFeature.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, false, false)
	s.localEntity.AddFeature(localFeature)

	s.sut.SetLocalFeature(s.localEntity, localFeature)

	remoteEntity := NewEntityRemote(s.remoteDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{1})
	s.remoteDevice.AddEntity(remoteEntity)

	remoteFeature := NewFeatureRemote(remoteEntity.NextFeatureId(), remoteEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
	remoteEntity.AddFeature(remoteFeature)

	subscrRequest := &model.SubscriptionManagementRequestCallType{
		ClientAddress:     remoteFeature.Address(),
		ServerAddress:     localFeature.Address(),
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeDeviceDiagnosis),
	}

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &SpecificationVersion,
			AddressSource: &model.FeatureAddressType{
				Device:  s.remoteDevice.Address(),
				Entity:  []model.AddressEntityType{0},
				Feature: util.Ptr(model.AddressFeatureType(0)),
			},
			AddressDestination: &model.FeatureAddressType{
				Device:  s.localDevice.Address(),
				Entity:  []model.AddressEntityType{0},
				Feature: util.Ptr(model.AddressFeatureType(0)),
			},
			MsgCounter:    util.Ptr(model.MsgCounterType(1000)),
			CmdClassifier: util.Ptr(model.CmdClassifierTypeCall),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{
				{
					NodeManagementSubscriptionRequestCall: &model.NodeManagementSubscriptionRequestCallType{
						SubscriptionRequest: subscrRequest,
					},
				},
			},
		},
	}
	err := s.localDevice.ProcessCmd(datagram, s.remoteDevice)
	assert.Nil(s.T(), err)

	data := localFeature.DataCopy(model.FunctionTypeDeviceDiagnosisHeartbeatData)
	assert.Nil(s.T(), data)

	running := s.sut.IsHeartbeatRunning()
	assert.Equal(s.T(), false, running)

	s.localDevice.RemoveEntity(s.localEntity)
	s.localEntity = NewEntityLocal(s.localDevice, model.EntityTypeTypeCEM, []model.AddressEntityType{1}, time.Second*4)
	s.localDevice.AddEntity(s.localEntity)

	localFeature = s.localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	localFeature.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)
	s.localEntity.AddFeature(localFeature)

	s.sut.SetLocalFeature(s.localEntity, localFeature)

	err = s.localDevice.ProcessCmd(datagram, s.remoteDevice)
	assert.Nil(s.T(), err)

	time.Sleep(time.Second * 5)

	running = s.sut.IsHeartbeatRunning()
	assert.Equal(s.T(), true, running)

	data = localFeature.DataCopy(model.FunctionTypeDeviceDiagnosisHeartbeatData)
	assert.NotNil(s.T(), data)

	fctData := data.(*model.DeviceDiagnosisHeartbeatDataType)
	var resultCounter uint64 = 1
	assert.LessOrEqual(s.T(), resultCounter, *fctData.HeartbeatCounter)
	resultTimeout, err := fctData.HeartbeatTimeout.GetTimeDuration()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), time.Second*4, resultTimeout)

	subscrDelRequest := &model.SubscriptionManagementDeleteCallType{
		ClientAddress: remoteFeature.Address(),
		ServerAddress: localFeature.Address(),
	}

	datagram.Payload = model.PayloadType{
		Cmd: []model.CmdType{
			{
				NodeManagementSubscriptionDeleteCall: &model.NodeManagementSubscriptionDeleteCallType{
					SubscriptionDelete: subscrDelRequest,
				},
			},
		},
	}

	err = s.localDevice.ProcessCmd(datagram, s.remoteDevice)
	assert.Nil(s.T(), err)

	isHeartbeatRunning := s.sut.IsHeartbeatRunning()
	assert.Equal(s.T(), true, isHeartbeatRunning)

	s.sut.StopHeartbeat()

	isHeartbeatRunning = s.sut.IsHeartbeatRunning()
	assert.Equal(s.T(), false, isHeartbeatRunning)
}
