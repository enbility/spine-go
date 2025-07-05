package spine

import (
	"testing"

	"github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestNegotiatedVersionEnforcementSuite(t *testing.T) {
	suite.Run(t, new(NegotiatedVersionEnforcementSuite))
}

type NegotiatedVersionEnforcementSuite struct {
	suite.Suite
	
	localDevice  *DeviceLocal
	remoteDevice *DeviceRemote
	senderMock   *mocks.SenderInterface
}

func (s *NegotiatedVersionEnforcementSuite) BeforeTest(suiteName, testName string) {
	s.localDevice = NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)

	ski := "test"
	s.senderMock = mocks.NewSenderInterface(s.T())
	s.remoteDevice = NewDeviceRemote(s.localDevice, ski, s.senderMock)
	desc := &model.NetworkManagementDeviceDescriptionDataType{
		DeviceAddress: &model.DeviceAddressType{
			Device: util.Ptr(model.AddressDeviceType("test")),
		},
	}
	s.remoteDevice.UpdateDevice(desc)
}

// Test that before discovery, compatible versions are allowed
func (s *NegotiatedVersionEnforcementSuite) Test_BeforeDiscovery_CompatibleVersionsAllowed() {
	// Before any version negotiation, no negotiated version should exist
	assert.Empty(s.T(), s.remoteDevice.NegotiatedProtocolVersion())
	
	// Compatible version should be allowed
	version := model.SpecificationVersionType("1.2.0")
	err := s.remoteDevice.validateProtocolVersion(&version, false)
	assert.NoError(s.T(), err)
	
	// Another compatible version should be allowed
	version = model.SpecificationVersionType("1.4.0")
	err = s.remoteDevice.validateProtocolVersion(&version, false)
	assert.NoError(s.T(), err)
	
	// Version 0.x should also be compatible with 1.x
	version = model.SpecificationVersionType("0.9.0")
	err = s.remoteDevice.validateProtocolVersion(&version, false)
	assert.NoError(s.T(), err)
}

// Test that after discovery, compatible versions within the same group are allowed
func (s *NegotiatedVersionEnforcementSuite) Test_AfterDiscovery_CompatibleVersionsAllowed() {
	// Simulate discovery negotiating version 1.3.0
	s.remoteDevice.SetNegotiatedProtocolVersion("1.3.0")
	
	// Verify negotiated version is set
	assert.Equal(s.T(), "1.3.0", s.remoteDevice.NegotiatedProtocolVersion())
	
	// Exact negotiated version should be allowed
	version := model.SpecificationVersionType("1.3.0")
	err := s.remoteDevice.validateProtocolVersion(&version, false)
	assert.NoError(s.T(), err)
	
	// Different compatible version should ALSO be allowed (asymmetric behavior)
	version = model.SpecificationVersionType("1.2.0")
	err = s.remoteDevice.validateProtocolVersion(&version, false)
	assert.NoError(s.T(), err)
	
	// Another different compatible version should ALSO be allowed
	version = model.SpecificationVersionType("1.4.0")
	err = s.remoteDevice.validateProtocolVersion(&version, false)
	assert.NoError(s.T(), err)
	
	// Version 0.x should be compatible with 1.x
	version = model.SpecificationVersionType("0.9.0")
	err = s.remoteDevice.validateProtocolVersion(&version, false)
	assert.NoError(s.T(), err)
}

// Test that incompatible versions are rejected both before and after discovery
func (s *NegotiatedVersionEnforcementSuite) Test_IncompatibleVersionsAlwaysRejected() {
	// Before discovery - incompatible version should be rejected
	version := model.SpecificationVersionType("2.0.0")
	err := s.remoteDevice.validateProtocolVersion(&version, false)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "incompatible major version")
	
	// After discovery - incompatible version should still be rejected
	s.remoteDevice.SetNegotiatedProtocolVersion("1.3.0")
	version = model.SpecificationVersionType("2.0.0")
	err = s.remoteDevice.validateProtocolVersion(&version, false)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "incompatible major version")
	
	// Even higher major versions should be rejected
	version = model.SpecificationVersionType("3.0.0")
	err = s.remoteDevice.validateProtocolVersion(&version, false)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "incompatible major version")
}

// Test edge cases with empty/nil versions
func (s *NegotiatedVersionEnforcementSuite) Test_EdgeCases_EmptyNilVersions() {
	// Set negotiated version
	s.remoteDevice.SetNegotiatedProtocolVersion("1.3.0")
	
	// Nil version should be allowed (real-world device compatibility)
	err := s.remoteDevice.validateProtocolVersion(nil, false)
	assert.NoError(s.T(), err)
	
	// Empty version should be allowed (real-world device compatibility)
	version := model.SpecificationVersionType("")
	err = s.remoteDevice.validateProtocolVersion(&version, false)
	assert.NoError(s.T(), err)
}

// Test non-compliant version strings - these are accepted for real-world compatibility
func (s *NegotiatedVersionEnforcementSuite) Test_NonCompliantVersions_Accepted() {
	// Set negotiated version
	s.remoteDevice.SetNegotiatedProtocolVersion("1.3.0")
	
	// Non-compliant version strings should be ACCEPTED for compatibility
	version := model.SpecificationVersionType("draft")
	err := s.remoteDevice.validateProtocolVersion(&version, false)
	assert.NoError(s.T(), err) // Liberal acceptance
	
	version = model.SpecificationVersionType("v1.3.0")
	err = s.remoteDevice.validateProtocolVersion(&version, false)
	assert.NoError(s.T(), err) // Liberal acceptance
	
	version = model.SpecificationVersionType("...")
	err = s.remoteDevice.validateProtocolVersion(&version, false)
	assert.NoError(s.T(), err) // Liberal acceptance
	
	version = model.SpecificationVersionType("1.3.0-RC1")
	err = s.remoteDevice.validateProtocolVersion(&version, false)
	assert.NoError(s.T(), err) // Liberal acceptance
}

// Test that discovery messages accept any version
func (s *NegotiatedVersionEnforcementSuite) Test_DiscoveryMessages_AcceptAnyVersion() {
	// Set negotiated version
	s.remoteDevice.SetNegotiatedProtocolVersion("1.3.0")
	
	// Discovery message with different compatible version should be allowed
	version := model.SpecificationVersionType("1.2.0")
	err := s.remoteDevice.validateProtocolVersion(&version, true) // isDiscoveryMessage = true
	assert.NoError(s.T(), err)
	
	// Discovery message with incompatible version should ALSO be allowed
	version = model.SpecificationVersionType("2.0.0")
	err = s.remoteDevice.validateProtocolVersion(&version, true) // isDiscoveryMessage = true
	assert.NoError(s.T(), err) // Discovery messages accept ANY version
	
	// Even invalid versions in discovery
	version = model.SpecificationVersionType("invalid-version")
	err = s.remoteDevice.validateProtocolVersion(&version, true) // isDiscoveryMessage = true
	assert.NoError(s.T(), err)
	
	// Future versions in discovery
	version = model.SpecificationVersionType("99.0.0")
	err = s.remoteDevice.validateProtocolVersion(&version, true) // isDiscoveryMessage = true
	assert.NoError(s.T(), err)
}

// Test version detection tracking
func (s *NegotiatedVersionEnforcementSuite) Test_VersionDetectionTracking() {
	// Set negotiated version and estimated version
	s.remoteDevice.SetNegotiatedProtocolVersion("1.3.0")
	s.remoteDevice.SetSupportedProtocolVersions([]string{"1.2.0", "1.3.0"})
	s.remoteDevice.UpdateEstimatedRemoteVersion()
	
	// Initially no detected version
	assert.Empty(s.T(), s.remoteDevice.DetectedRemoteVersion())
	
	// Process a regular message - should update detected version
	version := model.SpecificationVersionType("1.2.0")
	err := s.remoteDevice.validateProtocolVersion(&version, false)
	assert.NoError(s.T(), err)
	
	// Check detected version was updated
	assert.Equal(s.T(), "1.2.0", s.remoteDevice.DetectedRemoteVersion())
	
	// Process discovery message - should NOT update detected version
	version = model.SpecificationVersionType("1.3.0")
	err = s.remoteDevice.validateProtocolVersion(&version, true) // isDiscoveryMessage = true
	assert.NoError(s.T(), err)
	
	// Detected version should remain unchanged
	assert.Equal(s.T(), "1.2.0", s.remoteDevice.DetectedRemoteVersion())
}

// Test version compatibility within major version 2.x
func (s *NegotiatedVersionEnforcementSuite) Test_MajorVersion2_Compatibility() {
	// For major version >= 2, only exact major version match is allowed
	
	// Simulate a device that supports version 2.x
	s.remoteDevice.SetSupportedProtocolVersions([]string{"2.0.0", "2.1.0"})
	// Note: In reality, we wouldn't negotiate to 2.x since we support 1.x
	// But we can still test validation of incoming 2.x messages
	
	// Major version 2.x should be rejected (our implementation is 1.x)
	version := model.SpecificationVersionType("2.0.0")
	err := s.remoteDevice.validateProtocolVersion(&version, false)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "incompatible major version")
	
	version = model.SpecificationVersionType("2.1.0")
	err = s.remoteDevice.validateProtocolVersion(&version, false)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "incompatible major version")
}