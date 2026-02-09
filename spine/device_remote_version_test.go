package spine

import (
	"testing"

	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/suite"
)

func TestDeviceRemoteVersionDetection(t *testing.T) {
	suite.Run(t, new(DeviceRemoteVersionTestSuite))
}

type DeviceRemoteVersionTestSuite struct {
	suite.Suite
	localDevice  *DeviceLocal
	remoteDevice *DeviceRemote
}

func (s *DeviceRemoteVersionTestSuite) SetupTest() {
	s.localDevice = NewDeviceLocal("test", "test", "test", "test", "test", model.DeviceTypeTypeEnergyManagementSystem, "")
	s.remoteDevice = NewDeviceRemote(s.localDevice, "test-ski", nil)
}

// Test 1: Basic version detection - remote announces and uses same version
func (s *DeviceRemoteVersionTestSuite) Test_BasicVersionDetection() {
	// Remote announces 1.2.0
	s.remoteDevice.SetSupportedProtocolVersions([]string{"1.2.0"})
	
	// Calculate estimated version
	s.remoteDevice.UpdateEstimatedRemoteVersion()
	s.Equal("1.2.0", s.remoteDevice.EstimatedRemoteVersion())
	
	// Simulate receiving a message with version 1.2.0
	err := s.remoteDevice.UpdateDetectedVersion("1.2.0")
	s.NoError(err)
	s.Equal("1.2.0", s.remoteDevice.DetectedRemoteVersion())
}

// Test 2: Multiple versions in same compatibility group
func (s *DeviceRemoteVersionTestSuite) Test_MultipleVersionsSameGroup() {
	// Remote announces multiple 1.x versions
	s.remoteDevice.SetSupportedProtocolVersions([]string{"1.0.0", "1.1.0", "1.2.0", "1.3.0"})
	
	// Should estimate highest in group
	s.remoteDevice.UpdateEstimatedRemoteVersion()
	s.Equal("1.3.0", s.remoteDevice.EstimatedRemoteVersion())
	
	// But remote might use any of them
	err := s.remoteDevice.UpdateDetectedVersion("1.2.0")
	s.NoError(err)
	s.Equal("1.2.0", s.remoteDevice.DetectedRemoteVersion())
}

// Test 3: Versions across compatibility groups
func (s *DeviceRemoteVersionTestSuite) Test_CrossCompatibilityGroups() {
	// Remote announces versions from different groups
	s.remoteDevice.SetSupportedProtocolVersions([]string{"1.2.0", "2.0.0", "2.1.0"})
	
	// Should pick highest from our compatible group (1.x)
	s.remoteDevice.UpdateEstimatedRemoteVersion()
	s.Equal("1.2.0", s.remoteDevice.EstimatedRemoteVersion())
	
	// Remote uses their highest compatible version
	err := s.remoteDevice.UpdateDetectedVersion("1.2.0")
	s.NoError(err)
	s.Equal("1.2.0", s.remoteDevice.DetectedRemoteVersion())
}

// Test 4: No compatible versions
func (s *DeviceRemoteVersionTestSuite) Test_NoCompatibleVersions() {
	// Remote only supports incompatible versions
	s.remoteDevice.SetSupportedProtocolVersions([]string{"2.0.0", "3.0.0"})
	
	// No estimation possible
	s.remoteDevice.UpdateEstimatedRemoteVersion()
	s.Equal("", s.remoteDevice.EstimatedRemoteVersion())
	
	// Any version they send will be incompatible
	err := s.remoteDevice.UpdateDetectedVersion("2.0.0")
	s.Error(err)
	s.Contains(err.Error(), "incompatible")
}

// Test 5: Invalid version formats
func (s *DeviceRemoteVersionTestSuite) Test_InvalidVersionFormats() {
	// Remote announces mix of valid and invalid
	s.remoteDevice.SetSupportedProtocolVersions([]string{"draft", "1.2.0", "...", ""})
	
	// Should only consider valid versions
	s.remoteDevice.UpdateEstimatedRemoteVersion()
	s.Equal("1.2.0", s.remoteDevice.EstimatedRemoteVersion())
	
	// Receiving invalid version after valid announcement
	err := s.remoteDevice.UpdateDetectedVersion("...")
	s.NoError(err) // Liberal acceptance of invalid formats
	s.Equal("...", s.remoteDevice.DetectedRemoteVersion())
}

// Test 6: Version change detection
func (s *DeviceRemoteVersionTestSuite) Test_VersionChangeDetection() {
	s.remoteDevice.SetSupportedProtocolVersions([]string{"1.2.0", "1.3.0"})
	s.remoteDevice.UpdateEstimatedRemoteVersion()
	
	// First detection
	err := s.remoteDevice.UpdateDetectedVersion("1.2.0")
	s.NoError(err)
	
	// Version change - should log warning but accept
	err = s.remoteDevice.UpdateDetectedVersion("1.3.0")
	s.NoError(err) // Compatible change allowed
	s.Equal("1.3.0", s.remoteDevice.DetectedRemoteVersion())
	
	// Check version change was tracked
	s.True(s.remoteDevice.HasVersionChanged())
}

// Test 7: Empty/nil version handling
func (s *DeviceRemoteVersionTestSuite) Test_EmptyVersionHandling() {
	s.remoteDevice.SetSupportedProtocolVersions([]string{"1.2.0"})
	s.remoteDevice.UpdateEstimatedRemoteVersion()
	
	// Empty version string
	err := s.remoteDevice.UpdateDetectedVersion("")
	s.NoError(err) // Accept for compatibility
	s.Equal("", s.remoteDevice.DetectedRemoteVersion())
}

// Test 8: No announced versions
func (s *DeviceRemoteVersionTestSuite) Test_NoAnnouncedVersions() {
	// Remote didn't announce any versions
	s.remoteDevice.SetSupportedProtocolVersions([]string{})
	s.remoteDevice.UpdateEstimatedRemoteVersion()
	s.Equal("", s.remoteDevice.EstimatedRemoteVersion())
	
	// But sends messages with version
	err := s.remoteDevice.UpdateDetectedVersion("1.2.0")
	s.NoError(err) // Liberal acceptance
	s.Equal("1.2.0", s.remoteDevice.DetectedRemoteVersion())
}

// Test 9: Discovery message version validation
func (s *DeviceRemoteVersionTestSuite) Test_DiscoveryMessageAlwaysAccepted() {
	tests := []struct {
		name    string
		version string
	}{
		{"Compatible", "1.3.0"},
		{"Incompatible", "99.0.0"},
		{"Invalid", "draft"},
		{"Empty", ""},
		{"Malformed", "..."},
	}
	
	for _, tt := range tests {
		s.Run(tt.name, func() {
			datagram := s.createDiscoveryDatagram(tt.version)
			err := s.remoteDevice.ValidateDatagramVersion(datagram)
			s.NoError(err, "Discovery message should never be rejected")
		})
	}
}

// Test 10: Regular message version validation
func (s *DeviceRemoteVersionTestSuite) Test_RegularMessageVersionValidation() {
	// Setup detected version
	s.remoteDevice.UpdateDetectedVersion("1.2.0")
	
	tests := []struct {
		name      string
		version   string
		wantError bool
	}{
		{"Matching", "1.2.0", false},
		{"Compatible", "1.3.0", false}, // Allow compatible changes
		{"Incompatible", "2.0.0", true},
		{"Invalid", "draft", false}, // Liberal acceptance
		{"Empty", "", false},        // Liberal acceptance
	}
	
	for _, tt := range tests {
		s.Run(tt.name, func() {
			datagram := s.createRegularDatagram(tt.version)
			err := s.remoteDevice.ValidateDatagramVersion(datagram)
			if tt.wantError {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

// Helper methods
func (s *DeviceRemoteVersionTestSuite) createDiscoveryDatagram(version string) *model.DatagramType {
	var versionPtr *model.SpecificationVersionType
	if version != "" {
		v := model.SpecificationVersionType(version)
		versionPtr = &v
	}
	
	return &model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: versionPtr,
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{{
				NodeManagementDetailedDiscoveryData: &model.NodeManagementDetailedDiscoveryDataType{},
			}},
		},
	}
}

func (s *DeviceRemoteVersionTestSuite) createRegularDatagram(version string) *model.DatagramType {
	var versionPtr *model.SpecificationVersionType
	if version != "" {
		v := model.SpecificationVersionType(version)
		versionPtr = &v
	}
	
	return &model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: versionPtr,
			CmdClassifier: util.Ptr(model.CmdClassifierTypeReply),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{{}},
		},
	}
}