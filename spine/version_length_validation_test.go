package spine

import (
	"strings"
	"sync"
	"testing"

	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func TestVersionLengthValidation(t *testing.T) {
	address := model.AddressDeviceType("test-device")
	device := &DeviceRemote{
		Device:        NewDevice(&address, nil, nil),
		ski:           "test-device",
		versionsMutex: sync.RWMutex{},
	}

	tests := []struct {
		name        string
		version     string
		isDiscovery bool
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid short version",
			version:     "1.3.0",
			isDiscovery: false,
			expectError: false,
		},
		{
			name:        "Valid long version within limit",
			version:     "1.3.0-beta.1+build." + strings.Repeat("x", 100), // ~115 chars
			isDiscovery: false,
			expectError: false,
		},
		{
			name:        "Version exactly at 128 character limit",
			version:     "1.3.0-" + strings.Repeat("x", 122), // exactly 128 chars
			isDiscovery: false,
			expectError: false,
		},
		{
			name:        "Version exceeding 128 character limit",
			version:     "1.3.0-" + strings.Repeat("x", 123), // 129 chars
			isDiscovery: false,
			expectError: true,
			errorMsg:    "exceeds SPINE specification limit of 128 characters",
		},
		{
			name:        "Extremely long version string",
			version:     strings.Repeat("1.2.3-", 50), // ~300 chars
			isDiscovery: false,
			expectError: true,
			errorMsg:    "exceeds SPINE specification limit of 128 characters",
		},
		{
			name:        "Discovery message with long version - should still be rejected",
			version:     "1.3.0-" + strings.Repeat("x", 123), // 129 chars
			isDiscovery: true,
			expectError: true,
			errorMsg:    "exceeds SPINE specification limit of 128 characters",
		},
		{
			name:        "Empty version string",
			version:     "",
			isDiscovery: false,
			expectError: false,
		},
		{
			name:        "Version with whitespace trimmed to within limit",
			version:     "  1.3.0-" + strings.Repeat("x", 118) + "  ", // 128 chars after trim
			isDiscovery: false,
			expectError: false,
		},
		{
			name:        "Version with whitespace exceeding limit after trim",
			version:     "  1.3.0-" + strings.Repeat("x", 125) + "  ", // 129 chars after trim
			isDiscovery: false,
			expectError: true,
			errorMsg:    "exceeds SPINE specification limit of 128 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version := model.SpecificationVersionType(tt.version)
			err := device.validateProtocolVersion(&version, tt.isDiscovery)

			if tt.expectError {
				assert.Error(t, err, "Expected error for version: %s", tt.version)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg, "Error message should contain expected text")
				}
			} else {
				assert.NoError(t, err, "Expected no error for version: %s", tt.version)
			}
		})
	}
}

func TestVersionLengthValidation_ExactBoundaries(t *testing.T) {
	address := model.AddressDeviceType("boundary-test-device")
	device := &DeviceRemote{
		Device:        NewDevice(&address, nil, nil),
		ski:           "boundary-test-device",
		versionsMutex: sync.RWMutex{},
	}

	// Test exact boundary conditions
	boundaryTests := []struct {
		name        string
		length      int
		expectError bool
	}{
		{"127 characters", 127, false},
		{"128 characters", 128, false},
		{"129 characters", 129, true},
		{"200 characters", 200, true},
		{"1000 characters", 1000, true},
	}

	for _, tt := range boundaryTests {
		t.Run(tt.name, func(t *testing.T) {
			// Create version string of exact length
			baseVersion := "1.3.0-"
			padding := strings.Repeat("x", tt.length-len(baseVersion))
			version := model.SpecificationVersionType(baseVersion + padding)

			err := device.validateProtocolVersion(&version, false)

			if tt.expectError {
				assert.Error(t, err, "Should reject version of length %d", tt.length)
				assert.Contains(t, err.Error(), "exceeds SPINE specification limit", 
					"Error should mention specification limit")
			} else {
				assert.NoError(t, err, "Should accept version of length %d", tt.length)
			}
		})
	}
}

func TestVersionLengthValidation_RealWorldExamples(t *testing.T) {
	address := model.AddressDeviceType("real-world-test-device")
	device := &DeviceRemote{
		Device:        NewDevice(&address, nil, nil),
		ski:           "real-world-test-device",
		versionsMutex: sync.RWMutex{},
	}

	// Test real-world version examples (all should be accepted for length validation)
	realWorldVersions := []string{
		"1.3.0",
		"1.2.0",
		"1.0.0",
		"0.9.0", // Compatible major version
		"1.3.0-RC1",
		"1.3.0-beta.1",
		"1.3.0-alpha.1+build.123",
		"draft", // Non-compliant but accepted for compatibility
		"",      // Empty version
		"...",   // Dots only
		"v1.3.0", // Vendor prefix
	}

	for _, versionStr := range realWorldVersions {
		t.Run("Real world: "+versionStr, func(t *testing.T) {
			version := model.SpecificationVersionType(versionStr)
			err := device.validateProtocolVersion(&version, false)
			
			// All real-world examples should be accepted (they're all well under 128 chars)
			assert.NoError(t, err, "Real-world version should be accepted: %s", versionStr)
		})
	}
}