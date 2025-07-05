package spine

import (
	"fmt"
	"testing"

	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/suite"
)

func TestDeviceRemoteVersionEdgeCases(t *testing.T) {
	suite.Run(t, new(DeviceRemoteVersionEdgeTestSuite))
}

type DeviceRemoteVersionEdgeTestSuite struct {
	suite.Suite
	localDevice  *DeviceLocal
	remoteDevice *DeviceRemote
}

func (s *DeviceRemoteVersionEdgeTestSuite) SetupTest() {
	s.localDevice = NewDeviceLocal("test", "test", "test", "test", "test", model.DeviceTypeTypeEnergyManagementSystem, "")
	s.remoteDevice = NewDeviceRemote(s.localDevice, "test-ski", nil)
}

// Edge Case 1: Remote announces future version from same compatibility group
func (s *DeviceRemoteVersionEdgeTestSuite) Test_FutureVersionSameGroup() {
	// Remote announces a future 1.x version
	s.remoteDevice.SetSupportedProtocolVersions([]string{"1.5.0", "1.6.0"})
	s.remoteDevice.UpdateEstimatedRemoteVersion()
	
	// Should pick highest 1.x version
	s.Equal("1.6.0", s.remoteDevice.EstimatedRemoteVersion())
	
	// When they send 1.5.0, it should work
	err := s.remoteDevice.UpdateDetectedVersion("1.5.0")
	s.NoError(err)
}

// Edge Case 2: Version oscillation - remote alternates between versions
func (s *DeviceRemoteVersionEdgeTestSuite) Test_VersionOscillation() {
	s.remoteDevice.SetSupportedProtocolVersions([]string{"1.2.0", "1.3.0"})
	
	// First message with 1.2.0
	err := s.remoteDevice.UpdateDetectedVersion("1.2.0")
	s.NoError(err)
	s.False(s.remoteDevice.HasVersionChanged())
	
	// Second message with 1.3.0
	err = s.remoteDevice.UpdateDetectedVersion("1.3.0")
	s.NoError(err)
	s.True(s.remoteDevice.HasVersionChanged())
	
	// Third message back to 1.2.0
	err = s.remoteDevice.UpdateDetectedVersion("1.2.0")
	s.NoError(err)
	s.True(s.remoteDevice.HasVersionChanged())
}

// Edge Case 3: Malformed version after valid detection
func (s *DeviceRemoteVersionEdgeTestSuite) Test_MalformedVersionAfterValid() {
	// Start with valid version
	err := s.remoteDevice.UpdateDetectedVersion("1.2.0")
	s.NoError(err)
	
	// Then send malformed versions
	malformedVersions := []string{"", "...", "draft", "1.2", "v1.2.0", "1.2.0-RC1"}
	
	for _, v := range malformedVersions {
		err = s.remoteDevice.UpdateDetectedVersion(v)
		s.NoError(err, "Should accept malformed version %s for compatibility", v)
		s.Equal(v, s.remoteDevice.DetectedRemoteVersion())
		s.True(s.remoteDevice.HasVersionChanged())
	}
}

// Edge Case 4: Version downgrade within compatibility group
func (s *DeviceRemoteVersionEdgeTestSuite) Test_VersionDowngrade() {
	s.remoteDevice.SetSupportedProtocolVersions([]string{"1.1.0", "1.2.0", "1.3.0"})
	s.remoteDevice.UpdateEstimatedRemoteVersion()
	
	// Estimate should be highest
	s.Equal("1.3.0", s.remoteDevice.EstimatedRemoteVersion())
	
	// But device starts with 1.3.0
	err := s.remoteDevice.UpdateDetectedVersion("1.3.0")
	s.NoError(err)
	
	// Then downgrades to 1.1.0 (still compatible)
	err = s.remoteDevice.UpdateDetectedVersion("1.1.0")
	s.NoError(err)
	s.Equal("1.1.0", s.remoteDevice.DetectedRemoteVersion())
	s.True(s.remoteDevice.HasVersionChanged())
}

// Edge Case 5: Mix of valid and invalid versions in announcement
func (s *DeviceRemoteVersionEdgeTestSuite) Test_MixedValidInvalidAnnouncement() {
	// Announce mix of everything
	versions := []string{
		"",           // empty
		"...",        // dots
		"draft",      // draft
		"0.9.0",      // old version
		"1.2.0",      // valid compatible
		"1.3.0-RC1",  // RC version (invalid format)
		"2.0.0",      // incompatible
		"v1.4.0",     // invalid prefix
		"1.5",        // missing patch
	}
	
	s.remoteDevice.SetSupportedProtocolVersions(versions)
	s.remoteDevice.UpdateEstimatedRemoteVersion()
	
	// Should only pick valid 1.x version
	s.Equal("1.2.0", s.remoteDevice.EstimatedRemoteVersion())
}

// Edge Case 6: Boundary version numbers
func (s *DeviceRemoteVersionEdgeTestSuite) Test_BoundaryVersionNumbers() {
	// Test extreme version numbers
	s.remoteDevice.SetSupportedProtocolVersions([]string{"1.999.999"})
	s.remoteDevice.UpdateEstimatedRemoteVersion()
	s.Equal("1.999.999", s.remoteDevice.EstimatedRemoteVersion())
	
	// Zero versions
	s.remoteDevice.SetSupportedProtocolVersions([]string{"0.0.0"})
	s.remoteDevice.UpdateEstimatedRemoteVersion()
	s.Equal("0.0.0", s.remoteDevice.EstimatedRemoteVersion())
}

// Edge Case 7: Version detection without prior announcement
func (s *DeviceRemoteVersionEdgeTestSuite) Test_DetectionWithoutAnnouncement() {
	// No versions announced
	s.Equal("", s.remoteDevice.EstimatedRemoteVersion())
	
	// But device sends version anyway
	err := s.remoteDevice.UpdateDetectedVersion("1.2.0")
	s.NoError(err)
	s.Equal("1.2.0", s.remoteDevice.DetectedRemoteVersion())
	
	// Should handle incompatible version too
	err = s.remoteDevice.UpdateDetectedVersion("2.0.0")
	s.Error(err)
	s.Contains(err.Error(), "incompatible")
}

// Edge Case 8: Unicode and special characters in version
func (s *DeviceRemoteVersionEdgeTestSuite) Test_SpecialCharactersInVersion() {
	specialVersions := []string{
		"1.2.0\n",      // newline
		"1.2.0 ",       // trailing space
		" 1.2.0",       // leading space
		"1.2.0\t",      // tab
		"1.2.0\x00",    // null byte
		"1.२.0",        // unicode digit
		"1.2.0™",       // trademark symbol
	}
	
	for _, v := range specialVersions {
		err := s.remoteDevice.UpdateDetectedVersion(v)
		s.NoError(err, "Should handle special version: %q", v)
	}
}

// Edge Case 9: Rapid version changes
func (s *DeviceRemoteVersionEdgeTestSuite) Test_RapidVersionChanges() {
	versions := []string{"1.0.0", "1.1.0", "1.2.0", "1.3.0"}
	
	// Rapidly change versions
	for i := 0; i < 100; i++ {
		v := versions[i%len(versions)]
		err := s.remoteDevice.UpdateDetectedVersion(v)
		s.NoError(err)
	}
	
	// Should have tracked version changes
	s.True(s.remoteDevice.HasVersionChanged())
}

// Edge Case 10: Discovery followed by incompatible regular message
func (s *DeviceRemoteVersionEdgeTestSuite) Test_DiscoveryThenIncompatible() {
	// Discovery with future version - must accept
	discoveryDatagram := &model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: util.Ptr(model.SpecificationVersionType("5.0.0")),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{{
				NodeManagementDetailedDiscoveryData: &model.NodeManagementDetailedDiscoveryDataType{},
			}},
		},
	}
	
	err := s.remoteDevice.ValidateDatagramVersion(discoveryDatagram)
	s.NoError(err, "Discovery must always be accepted")
	
	// Regular message with same version - should fail
	regularDatagram := &model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: util.Ptr(model.SpecificationVersionType("5.0.0")),
			CmdClassifier: util.Ptr(model.CmdClassifierTypeReply),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{{}},
		},
	}
	
	err = s.remoteDevice.ValidateDatagramVersion(regularDatagram)
	s.Error(err, "Regular message with incompatible version should fail")
}

// Edge Case 11: Nil and empty handling throughout
func (s *DeviceRemoteVersionEdgeTestSuite) Test_NilHandling() {
	// Nil datagram
	var nilDatagram *model.DatagramType
	err := s.remoteDevice.ValidateDatagramVersion(nilDatagram)
	s.NoError(err) // Should handle gracefully
	
	// Datagram with nil header
	datagramNilHeader := &model.DatagramType{}
	err = s.remoteDevice.ValidateDatagramVersion(datagramNilHeader)
	s.NoError(err)
	
	// Empty version after detection
	s.remoteDevice.UpdateDetectedVersion("1.2.0")
	err = s.remoteDevice.UpdateDetectedVersion("")
	s.NoError(err)
	s.Equal("", s.remoteDevice.DetectedRemoteVersion())
}

// Edge Case 12: Concurrent version updates (race condition test)
func (s *DeviceRemoteVersionEdgeTestSuite) Test_ConcurrentVersionUpdates() {
	// Run multiple goroutines updating versions
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(n int) {
			version := fmt.Sprintf("1.%d.0", n)
			_ = s.remoteDevice.UpdateDetectedVersion(version)
			done <- true
		}(i)
	}
	
	// Wait for all to complete
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// Should have a version set and no panic
	s.NotEmpty(s.remoteDevice.DetectedRemoteVersion())
}