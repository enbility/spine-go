package spine

import (
	"encoding/json"
	"testing"

	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestVersionValidationSuite(t *testing.T) {
	suite.Run(t, new(VersionValidationSuite))
}

type VersionValidationSuite struct {
	suite.Suite

	localDevice  *DeviceLocal
	remoteDevice *DeviceRemote
}

func (s *VersionValidationSuite) WriteShipMessageWithPayload([]byte) {}

func (s *VersionValidationSuite) BeforeTest(suiteName, testName string) {
	s.localDevice = NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)

	ski := "test"
	sender := NewSender(s, nil)
	s.remoteDevice = NewDeviceRemote(s.localDevice, ski, sender)
	desc := &model.NetworkManagementDeviceDescriptionDataType{
		DeviceAddress: &model.DeviceAddressType{
			Device: util.Ptr(model.AddressDeviceType("test")),
		},
	}
	s.remoteDevice.UpdateDevice(desc)
	_ = s.localDevice.SetupRemoteDevice(ski, s)
}

func (s *VersionValidationSuite) createTestMessage(version string) []byte {
	versionType := model.SpecificationVersionType(version)
	msg := model.Datagram{
		Datagram: model.DatagramType{
			Header: model.HeaderType{
				SpecificationVersion: &versionType,
				AddressSource: &model.FeatureAddressType{
					Device:  util.Ptr(model.AddressDeviceType("test")),
					Entity:  []model.AddressEntityType{0},
					Feature: util.Ptr(model.AddressFeatureType(0)),
				},
				AddressDestination: &model.FeatureAddressType{
					Device:  util.Ptr(model.AddressDeviceType("local")),
					Entity:  []model.AddressEntityType{0},
					Feature: util.Ptr(model.AddressFeatureType(0)),
				},
				MsgCounter:    util.Ptr(model.MsgCounterType(1)),
				CmdClassifier: util.Ptr(model.CmdClassifierTypeRead),
			},
			Payload: model.PayloadType{
				Cmd: []model.CmdType{
					{
						Function: util.Ptr(model.FunctionTypeNodeManagementDetailedDiscoveryData),
					},
				},
			},
		},
	}

	data, _ := json.Marshal(msg)
	return data
}

// Test cases for version validation

func (s *VersionValidationSuite) Test_VersionValidation_CompatibleVersions() {
	// Test: Version 1.3.0 should be accepted (current version)
	msg := s.createTestMessage("1.3.0")
	msgCounter, err := s.remoteDevice.HandleSpineMesssage(msg)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	// Test: Version 1.2.0 should be accepted (older minor version)
	msg = s.createTestMessage("1.2.0")
	msgCounter, err = s.remoteDevice.HandleSpineMesssage(msg)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	// Test: Version 1.4.0 should be accepted (newer minor version)
	msg = s.createTestMessage("1.4.0")
	msgCounter, err = s.remoteDevice.HandleSpineMesssage(msg)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	// Test: Version 1.0.0 should be accepted (oldest 1.x version)
	msg = s.createTestMessage("1.0.0")
	msgCounter, err = s.remoteDevice.HandleSpineMesssage(msg)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), msgCounter)
}

func (s *VersionValidationSuite) Test_VersionValidation_IncompatibleVersions() {
	// Test: Version 2.0.0 should be rejected (different major version)
	msg := s.createTestMessage("2.0.0")
	msgCounter, err := s.remoteDevice.HandleSpineMesssage(msg)
	assert.Error(s.T(), err)
	if err != nil {
		assert.Contains(s.T(), err.Error(), "incompatible major version")
	}
	assert.Nil(s.T(), msgCounter)

	// Test: Version 3.0.0 should be rejected (different major version)
	msg = s.createTestMessage("3.0.0")
	msgCounter, err = s.remoteDevice.HandleSpineMesssage(msg)
	assert.Error(s.T(), err)
	if err != nil {
		assert.Contains(s.T(), err.Error(), "incompatible major version")
	}
	assert.Nil(s.T(), msgCounter)
}

func (s *VersionValidationSuite) Test_VersionValidation_NonCompliantVersions() {
	// Test: Empty version string should be accepted (real-world compatibility)
	msg := s.createTestMessage("")
	msgCounter, err := s.remoteDevice.HandleSpineMesssage(msg)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	// Test: Dots only version should be accepted (real-world compatibility)
	msg = s.createTestMessage("...")
	msgCounter, err = s.remoteDevice.HandleSpineMesssage(msg)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	// Test: Draft version should be accepted (real-world compatibility)
	msg = s.createTestMessage("draft")
	msgCounter, err = s.remoteDevice.HandleSpineMesssage(msg)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	// Test: RC version should be accepted (real-world compatibility)
	msg = s.createTestMessage("1.3.0-RC1")
	msgCounter, err = s.remoteDevice.HandleSpineMesssage(msg)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	// Test: Version with 'v' prefix should be accepted
	msg = s.createTestMessage("v1.3.0")
	msgCounter, err = s.remoteDevice.HandleSpineMesssage(msg)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), msgCounter)
}

func (s *VersionValidationSuite) Test_VersionValidation_NilVersion() {
	// Test: Nil version should be accepted (missing version field)
	msg := model.Datagram{
		Datagram: model.DatagramType{
			Header: model.HeaderType{
				// No SpecificationVersion field
				AddressSource: &model.FeatureAddressType{
					Device:  util.Ptr(model.AddressDeviceType("test")),
					Entity:  []model.AddressEntityType{0},
					Feature: util.Ptr(model.AddressFeatureType(0)),
				},
				AddressDestination: &model.FeatureAddressType{
					Device:  util.Ptr(model.AddressDeviceType("local")),
					Entity:  []model.AddressEntityType{0},
					Feature: util.Ptr(model.AddressFeatureType(0)),
				},
				MsgCounter:    util.Ptr(model.MsgCounterType(1)),
				CmdClassifier: util.Ptr(model.CmdClassifierTypeRead),
			},
			Payload: model.PayloadType{
				Cmd: []model.CmdType{
					{
						Function: util.Ptr(model.FunctionTypeNodeManagementDetailedDiscoveryData),
					},
				},
			},
		},
	}

	data, _ := json.Marshal(msg)
	msgCounter, err := s.remoteDevice.HandleSpineMesssage(data)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), msgCounter)
}

func (s *VersionValidationSuite) Test_VersionValidation_EdgeCases() {
	// Test: Version starting with 2 but not valid semantic version format - should be accepted
	msg := s.createTestMessage("2draft")
	msgCounter, err := s.remoteDevice.HandleSpineMesssage(msg)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	// Test: Very long version string starting with 1
	msg = s.createTestMessage("1.999999.999999")
	msgCounter, err = s.remoteDevice.HandleSpineMesssage(msg)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	// Test: Version with extra spaces
	msg = s.createTestMessage(" 1.3.0 ")
	msgCounter, err = s.remoteDevice.HandleSpineMesssage(msg)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), msgCounter)
}