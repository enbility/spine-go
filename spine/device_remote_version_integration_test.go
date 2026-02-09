package spine

import (
	"encoding/json"
	"testing"

	"github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestVersionDetectionIntegration(t *testing.T) {
	suite.Run(t, new(VersionDetectionIntegrationSuite))
}

type VersionDetectionIntegrationSuite struct {
	suite.Suite
	localDevice  *DeviceLocal
	remoteDevice *DeviceRemote
	senderMock   *mocks.SenderInterface
}

func (s *VersionDetectionIntegrationSuite) SetupTest() {
	s.localDevice = NewDeviceLocal("test", "test", "test", "test", "local", model.DeviceTypeTypeEnergyManagementSystem, "")
	s.senderMock = mocks.NewSenderInterface(s.T())
	s.remoteDevice = NewDeviceRemote(s.localDevice, "test-ski", s.senderMock)
}

// Test complete flow: discovery -> version detection -> validation
func (s *VersionDetectionIntegrationSuite) Test_CompleteVersionDetectionFlow() {
	// Step 1: Simulate discovery by setting supported versions
	s.remoteDevice.SetSupportedProtocolVersions([]string{"1.2.0"})
	s.remoteDevice.UpdateEstimatedRemoteVersion()
	
	// Verify version estimation
	s.Equal("1.2.0", s.remoteDevice.EstimatedRemoteVersion())
	s.Equal([]string{"1.2.0"}, s.remoteDevice.SupportedProtocolVersions())
	
	// No detected version yet
	s.Equal("", s.remoteDevice.DetectedRemoteVersion())

	// Step 2: Send regular message with version
	regularMsg := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: util.Ptr(model.SpecificationVersionType("1.2.0")),
			AddressSource: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("remote")),
			},
			AddressDestination: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("local")),
			},
			MsgCounter:    util.Ptr(model.MsgCounterType(2)),
			CmdClassifier: util.Ptr(model.CmdClassifierTypeRead),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{{
				// Some read command
			}},
		},
	}

	// Process regular message
	msgBytes, err := json.Marshal(model.Datagram{Datagram: regularMsg})
	s.Require().NoError(err)
	
	msgCounter, err := s.remoteDevice.HandleSpineMesssage(msgBytes)
	s.NoError(err)
	s.NotNil(msgCounter)

	// Now version should be detected
	s.Equal("1.2.0", s.remoteDevice.DetectedRemoteVersion())
	s.False(s.remoteDevice.HasVersionChanged())
}

// Test version mismatch between estimate and actual
func (s *VersionDetectionIntegrationSuite) Test_VersionMismatchDetection() {
	// Remote announces multiple versions
	s.remoteDevice.SetSupportedProtocolVersions([]string{"1.1.0", "1.2.0", "1.3.0"})
	s.remoteDevice.UpdateEstimatedRemoteVersion()

	// We estimate they'll use highest (1.3.0)
	s.Equal("1.3.0", s.remoteDevice.EstimatedRemoteVersion())

	// But they actually use 1.2.0
	regularMsg := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: util.Ptr(model.SpecificationVersionType("1.2.0")),
			AddressSource: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("remote")),
			},
			AddressDestination: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("local")),
			},
			MsgCounter:    util.Ptr(model.MsgCounterType(2)),
			CmdClassifier: util.Ptr(model.CmdClassifierTypeRead),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{{}},
		},
	}

	msgBytes2, err2 := json.Marshal(model.Datagram{Datagram: regularMsg})
	s.Require().NoError(err2)
	
	_, err2 = s.remoteDevice.HandleSpineMesssage(msgBytes2)
	s.NoError(err2)

	// Detected differs from estimate
	s.Equal("1.2.0", s.remoteDevice.DetectedRemoteVersion())
	s.Equal("1.3.0", s.remoteDevice.EstimatedRemoteVersion())
}

// Test incompatible version rejection after discovery
func (s *VersionDetectionIntegrationSuite) Test_IncompatibleVersionRejection() {
	// Remote only supports 1.x
	s.remoteDevice.SetSupportedProtocolVersions([]string{"1.2.0"})
	s.remoteDevice.UpdateEstimatedRemoteVersion()
	
	// Discovery can use any version
	discoveryMsg := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: util.Ptr(model.SpecificationVersionType("99.0.0")), // Future version in discovery
			AddressSource: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("remote")),
			},
			AddressDestination: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("local")),
			},
			MsgCounter:    util.Ptr(model.MsgCounterType(1)),
			CmdClassifier: util.Ptr(model.CmdClassifierTypeReply),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{{
				NodeManagementDetailedDiscoveryData: &model.NodeManagementDetailedDiscoveryDataType{},
			}},
		},
	}

	msgBytes, err := json.Marshal(model.Datagram{Datagram: discoveryMsg})
	s.Require().NoError(err)
	
	_, err = s.remoteDevice.HandleSpineMesssage(msgBytes)
	s.NoError(err) // Discovery accepted despite incompatible version

	// But regular message with incompatible version fails
	regularMsg := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: util.Ptr(model.SpecificationVersionType("2.0.0")),
			AddressSource: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("remote")),
			},
			AddressDestination: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("local")),
			},
			MsgCounter:    util.Ptr(model.MsgCounterType(2)),
			CmdClassifier: util.Ptr(model.CmdClassifierTypeRead),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{{}},
		},
	}

	// Expect error response to be sent
	s.senderMock.EXPECT().ResultError(
		&regularMsg.Header,
		mock.Anything,
		mock.MatchedBy(func(err *model.ErrorType) bool {
			return err.ErrorNumber == model.ErrorNumberTypeGeneralError
		}),
	).Return(nil)

	msgBytes3, err3 := json.Marshal(model.Datagram{Datagram: regularMsg})
	s.Require().NoError(err3)
	
	_, err3 = s.remoteDevice.HandleSpineMesssage(msgBytes3)
	s.Error(err3)
	s.Contains(err3.Error(), "incompatible")
}