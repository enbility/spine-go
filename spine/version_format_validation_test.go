package spine

import (
	"testing"

	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestVersionFormatValidationSuite(t *testing.T) {
	suite.Run(t, new(VersionFormatValidationSuite))
}

type VersionFormatValidationSuite struct {
	suite.Suite

	localDevice  *DeviceLocal
	remoteDevice *DeviceRemote
}

func (s *VersionFormatValidationSuite) WriteShipMessageWithPayload([]byte) {}

func (s *VersionFormatValidationSuite) BeforeTest(suiteName, testName string) {
	s.localDevice = NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)

	ski := "test"
	sender := NewSender(s, nil)
	s.remoteDevice = NewDeviceRemote(s.localDevice, ski, sender)
}

// Test helper to validate version strings
func (s *VersionFormatValidationSuite) testVersionValidation(version string, expectError bool, errorContains string) {
	versionType := model.SpecificationVersionType(version)
	err := s.remoteDevice.validateProtocolVersion(&versionType, false)
	
	if expectError {
		assert.Error(s.T(), err)
		if errorContains != "" {
			assert.Contains(s.T(), err.Error(), errorContains)
		}
	} else {
		assert.NoError(s.T(), err)
	}
}

// Test valid semantic version formats
func (s *VersionFormatValidationSuite) Test_ValidSemanticVersions() {
	// Test: Standard semantic versions should be validated
	s.testVersionValidation("1.0.0", false, "")
	s.testVersionValidation("1.2.3", false, "")
	s.testVersionValidation("1.3.0", false, "")
	s.testVersionValidation("1.99.99", false, "")
	
	// Test: Major version 2+ should be rejected
	s.testVersionValidation("2.0.0", true, "incompatible major version")
	s.testVersionValidation("3.1.4", true, "incompatible major version")
	s.testVersionValidation("10.0.0", true, "incompatible major version")
}

// Test invalid format versions - should be assumed compatible
func (s *VersionFormatValidationSuite) Test_InvalidFormatAssumedCompatible() {
	// Test: Empty string
	s.testVersionValidation("", false, "")
	
	// Test: Non-semantic versions
	s.testVersionValidation("...", false, "")
	s.testVersionValidation("draft", false, "")
	s.testVersionValidation("unknown", false, "")
	
	// Test: Missing parts
	s.testVersionValidation("1", false, "")
	s.testVersionValidation("1.2", false, "")
	
	// Test: Extra parts
	s.testVersionValidation("1.2.3.4", false, "")
	
	// Test: Non-numeric parts
	s.testVersionValidation("1.2.x", false, "")
	s.testVersionValidation("a.b.c", false, "")
	
	// Test: Prefixes/suffixes
	s.testVersionValidation("v1.3.0", false, "")
	s.testVersionValidation("1.3.0-RC1", false, "")
	s.testVersionValidation("1.3.0-beta", false, "")
	s.testVersionValidation("ver1.3.0", false, "")
	
	// Test: Special characters
	s.testVersionValidation("1.3.0!", false, "")
	s.testVersionValidation("@1.3.0", false, "")
	
	// Test: Whitespace
	s.testVersionValidation(" 1.3.0", false, "")
	s.testVersionValidation("1.3.0 ", false, "")
	s.testVersionValidation(" 1.3.0 ", false, "")
	s.testVersionValidation("1. 3.0", false, "")
}

// Test edge cases for major version detection
func (s *VersionFormatValidationSuite) Test_MajorVersionEdgeCases() {
	// Test: Valid format with major version 0 (should be compatible)
	s.testVersionValidation("0.1.0", false, "")
	s.testVersionValidation("0.99.99", false, "")
	
	// Test: Large but valid version numbers
	s.testVersionValidation("1.999.999", false, "")
	s.testVersionValidation("999.0.0", true, "incompatible major version")
	
	// Test: Leading zeros (technically invalid but might occur)
	s.testVersionValidation("01.02.03", false, "")
	s.testVersionValidation("001.002.003", false, "")
	s.testVersionValidation("02.0.0", true, "incompatible major version")
}

// Test versions that look like they might be 2.x but aren't
func (s *VersionFormatValidationSuite) Test_AmbiguousVersions() {
	// These should NOT be rejected as they're not valid semantic versions
	s.testVersionValidation("2", false, "")
	s.testVersionValidation("2draft", false, "")
	s.testVersionValidation("2.x.x", false, "")
	s.testVersionValidation("2.0", false, "")
	s.testVersionValidation("2.0.0.1", false, "")
	s.testVersionValidation("v2.0.0", false, "")
	s.testVersionValidation("12.0.0", true, "incompatible major version") // Valid format, but major version 12 is incompatible
}

// Test that we handle nil versions
func (s *VersionFormatValidationSuite) Test_NilVersion() {
	err := s.remoteDevice.validateProtocolVersion(nil, false)
	assert.NoError(s.T(), err)
}

// Test numeric overflow protection
func (s *VersionFormatValidationSuite) Test_NumericOverflow() {
	// Very large numbers should be treated as invalid format
	s.testVersionValidation("999999999999999999999.0.0", false, "")
	s.testVersionValidation("1.999999999999999999999.0", false, "")
	s.testVersionValidation("1.0.999999999999999999999", false, "")
}

// Test negative numbers
func (s *VersionFormatValidationSuite) Test_NegativeNumbers() {
	// Negative numbers make it invalid format
	s.testVersionValidation("-1.0.0", false, "")
	s.testVersionValidation("1.-2.0", false, "")
	s.testVersionValidation("1.0.-3", false, "")
}

// Test valid semantic versions that should trigger trace logging
func (s *VersionFormatValidationSuite) Test_ValidSemanticVersionsTraceLogging() {
	// These versions should parse as valid semantic versions and trigger trace logging
	// because they have major version 0 or 1 (compatible)
	
	// Test major version 0 (compatible)
	s.testVersionValidation("0.1.0", false, "")
	s.testVersionValidation("0.9.0", false, "")
	s.testVersionValidation("0.99.99", false, "")
	
	// Test major version 1 (compatible) - current SPINE version group
	s.testVersionValidation("1.0.0", false, "")
	s.testVersionValidation("1.1.0", false, "")
	s.testVersionValidation("1.2.0", false, "")
	s.testVersionValidation("1.3.0", false, "")
	s.testVersionValidation("1.99.99", false, "")
	
	// These should all pass validation and trigger the trace logging at line 416:
	// logging.Log().Trace("Valid protocol version from", d.address, ":", major, ".", minor, ".", patch)
}