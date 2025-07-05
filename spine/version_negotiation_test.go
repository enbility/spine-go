package spine

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestVersionNegotiationSuite(t *testing.T) {
	suite.Run(t, new(VersionNegotiationSuite))
}

type VersionNegotiationSuite struct {
	suite.Suite
}

// Test compareVersions function with various version combinations
func (s *VersionNegotiationSuite) TestCompareVersions() {
	testCases := []struct {
		name     string
		v1       string
		v2       string
		expected int // 1 if v1 > v2, -1 if v1 < v2, 0 if equal
	}{
		// Equal versions
		{"Same version", "1.3.0", "1.3.0", 0},
		{"Same complex version", "1.2.3", "1.2.3", 0},

		// Major version differences
		{"Major v1 higher", "2.0.0", "1.0.0", 1},
		{"Major v2 higher", "1.0.0", "2.0.0", -1},
		{"Major large diff", "10.0.0", "1.0.0", 1},

		// Minor version differences (same major)
		{"Minor v1 higher", "1.3.0", "1.2.0", 1},
		{"Minor v2 higher", "1.2.0", "1.3.0", -1},
		{"Minor large diff", "1.10.0", "1.2.0", 1},

		// Patch version differences (same major.minor)
		{"Patch v1 higher", "1.3.1", "1.3.0", 1},
		{"Patch v2 higher", "1.3.0", "1.3.1", -1},
		{"Patch large diff", "1.3.10", "1.3.2", 1},

		// Mixed differences
		{"Higher major wins over minor", "2.0.0", "1.9.9", 1},
		{"Higher major wins over patch", "2.0.0", "1.3.99", 1},
		{"Higher minor wins over patch", "1.3.0", "1.2.99", 1},

		// Version 0.x.x handling
		{"Zero major versions", "0.2.0", "0.1.0", 1},
		{"Zero vs one major", "1.0.0", "0.9.0", 1},

		// Invalid versions (fall back to string comparison)
		{"Both invalid", "invalid1", "invalid2", -1}, // "invalid1" < "invalid2"
		{"Both invalid reverse", "invalid2", "invalid1", 1},
		{"v1 invalid", "invalid", "1.0.0", -1},
		{"v2 invalid", "1.0.0", "invalid", 1},
		{"Non-semantic format", "v1.0.0", "1.0.0", -1}, // invalid vs valid, invalid loses
		{"Partial version", "1.0", "1.0.0", -1},       // invalid vs valid
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			result := compareVersions(tc.v1, tc.v2)
			assert.Equal(t, tc.expected, result, "compareVersions(%s, %s)", tc.v1, tc.v2)

			// Test symmetry for non-equal cases
			if tc.expected != 0 {
				reverseResult := compareVersions(tc.v2, tc.v1)
				assert.Equal(t, -tc.expected, reverseResult, "Symmetry test failed")
			}
		})
	}
}

// Test areVersionsCompatible function
func (s *VersionNegotiationSuite) TestAreVersionsCompatible() {
	testCases := []struct {
		name       string
		v1         string
		v2         string
		compatible bool
	}{
		// Major version 0 and 1 compatibility
		{"Both major 0", "0.1.0", "0.2.0", true},
		{"Both major 1", "1.1.0", "1.2.0", true},
		{"Major 0 and 1", "0.9.0", "1.0.0", true},
		{"Major 1 and 0", "1.3.0", "0.5.0", true},

		// Major version 2+ must match exactly
		{"Both major 2", "2.0.0", "2.1.0", true},
		{"Both major 3", "3.1.0", "3.2.0", true},
		{"Major 2 and 3", "2.0.0", "3.0.0", false},
		{"Major 1 and 2", "1.3.0", "2.0.0", false},
		{"Major 0 and 2", "0.9.0", "2.0.0", false},

		// Large major version differences
		{"Major 5 and 10", "5.0.0", "10.0.0", false},
		{"Major 2 and 100", "2.0.0", "100.0.0", false},

		// Invalid versions (assume compatible)
		{"Both invalid", "invalid1", "invalid2", true},
		{"v1 invalid", "invalid", "1.0.0", true},
		{"v2 invalid", "1.0.0", "invalid", true},
		{"Empty strings", "", "", true},
		{"One empty", "", "1.0.0", true},
		{"Non-semantic", "v1.0.0", "1.0.0", true},
		{"Partial versions", "1.0", "1.0.0", true},

		// Edge cases with version numbers
		{"Same version", "1.3.0", "1.3.0", true},
		{"Zero versions", "0.0.0", "0.0.0", true},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			result := areVersionsCompatible(tc.v1, tc.v2)
			assert.Equal(t, tc.compatible, result, "areVersionsCompatible(%s, %s)", tc.v1, tc.v2)

			// Test symmetry
			reverseResult := areVersionsCompatible(tc.v2, tc.v1)
			assert.Equal(t, tc.compatible, reverseResult, "Symmetry test failed")
		})
	}
}

// Test findCompatibleVersion function
func (s *VersionNegotiationSuite) TestFindCompatibleVersion() {
	testCases := []struct {
		name            string
		localVersions   []string
		remoteVersions  []string
		expectedVersion string
		expectedFound   bool
		description     string
	}{
		{
			name:            "Compatible major 0 and 1",
			localVersions:   []string{"1.3.0"},
			remoteVersions:  []string{"0.9.0"},
			expectedVersion: "1.3.0", // Higher version wins
			expectedFound:   true,
			description:     "Major versions 0 and 1 are compatible",
		},
		{
			name:            "Compatible same major version",
			localVersions:   []string{"1.2.0"},
			remoteVersions:  []string{"1.3.0"},
			expectedVersion: "1.3.0", // Higher version wins
			expectedFound:   true,
			description:     "Same major version 1",
		},
		{
			name:            "Incompatible major versions",
			localVersions:   []string{"1.3.0"},
			remoteVersions:  []string{"2.0.0"},
			expectedVersion: "",
			expectedFound:   false,
			description:     "Major version 1 vs 2 not compatible",
		},
		{
			name:            "Multiple versions with compatibility",
			localVersions:   []string{"1.2.0", "1.3.0"},
			remoteVersions:  []string{"0.9.0", "1.1.0"},
			expectedVersion: "1.2.0", // First compatible found (1.2.0 vs 0.9.0)
			expectedFound:   true,
			description:     "Multiple versions, find first compatible",
		},
		{
			name:            "Multiple versions no compatibility",
			localVersions:   []string{"1.2.0", "1.3.0"},
			remoteVersions:  []string{"2.0.0", "3.0.0"},
			expectedVersion: "",
			expectedFound:   false,
			description:     "No compatible versions between 1.x and 2.x/3.x",
		},
		{
			name:            "Empty local versions",
			localVersions:   []string{},
			remoteVersions:  []string{"1.3.0"},
			expectedVersion: "",
			expectedFound:   false,
			description:     "Empty local versions",
		},
		{
			name:            "Empty remote versions",
			localVersions:   []string{"1.3.0"},
			remoteVersions:  []string{},
			expectedVersion: "",
			expectedFound:   false,
			description:     "Empty remote versions",
		},
		{
			name:            "Invalid versions assume compatible",
			localVersions:   []string{"invalid1"},
			remoteVersions:  []string{"invalid2"},
			expectedVersion: "invalid2", // String comparison: "invalid1" < "invalid2"
			expectedFound:   true,
			description:     "Invalid versions are assumed compatible",
		},
		{
			name:            "Mixed valid and invalid",
			localVersions:   []string{"1.3.0"},
			remoteVersions:  []string{"invalid"},
			expectedVersion: "1.3.0", // Valid version wins
			expectedFound:   true,
			description:     "Mixed valid/invalid versions",
		},
		{
			name:            "Major version 2 compatibility",
			localVersions:   []string{"2.1.0"},
			remoteVersions:  []string{"2.0.0"},
			expectedVersion: "2.1.0", // Higher version wins within major 2
			expectedFound:   true,
			description:     "Same major version 2",
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			version, found := findCompatibleVersion(tc.localVersions, tc.remoteVersions)
			assert.Equal(t, tc.expectedFound, found, "Expected found=%t for %s", tc.expectedFound, tc.description)
			if tc.expectedFound {
				assert.Equal(t, tc.expectedVersion, version, "Expected version=%s for %s", tc.expectedVersion, tc.description)
			}
		})
	}
}

// Test findHighestCommonVersion function
func (s *VersionNegotiationSuite) TestFindHighestCommonVersion() {
	testCases := []struct {
		name            string
		localVersions   []string
		remoteVersions  []string
		expectedVersion string
		expectedFound   bool
		description     string
	}{
		{
			name:            "Exact match single version",
			localVersions:   []string{"1.3.0"},
			remoteVersions:  []string{"1.3.0"},
			expectedVersion: "1.3.0",
			expectedFound:   true,
			description:     "Single exact match",
		},
		{
			name:            "Exact match multiple versions",
			localVersions:   []string{"1.2.0", "1.3.0"},
			remoteVersions:  []string{"1.3.0", "1.4.0"},
			expectedVersion: "1.3.0",
			expectedFound:   true,
			description:     "One exact match among multiple",
		},
		{
			name:            "Multiple exact matches - returns highest",
			localVersions:   []string{"1.2.0", "1.3.0", "1.4.0"},
			remoteVersions:  []string{"1.2.0", "1.3.0"},
			expectedVersion: "1.3.0", // Highest common version
			expectedFound:   true,
			description:     "Multiple exact matches, should return highest",
		},
		{
			name:            "No exact match but compatible",
			localVersions:   []string{"1.3.0"},
			remoteVersions:  []string{"1.2.0"},
			expectedVersion: "1.3.0", // Higher compatible version
			expectedFound:   true,
			description:     "No exact match, falls back to compatibility",
		},
		{
			name:            "No exact match, incompatible",
			localVersions:   []string{"1.3.0"},
			remoteVersions:  []string{"2.0.0"},
			expectedVersion: "",
			expectedFound:   false,
			description:     "No exact match and incompatible",
		},
		{
			name:            "Complex scenario with exact and compatible",
			localVersions:   []string{"1.1.0", "1.3.0", "2.0.0"},
			remoteVersions:  []string{"1.3.0", "1.4.0", "3.0.0"},
			expectedVersion: "1.3.0", // Exact match takes precedence
			expectedFound:   true,
			description:     "Exact match preferred over compatible",
		},
		{
			name:            "Empty local versions",
			localVersions:   []string{},
			remoteVersions:  []string{"1.3.0"},
			expectedVersion: "",
			expectedFound:   false,
			description:     "Empty local versions",
		},
		{
			name:            "Empty remote versions",
			localVersions:   []string{"1.3.0"},
			remoteVersions:  []string{},
			expectedVersion: "",
			expectedFound:   false,
			description:     "Empty remote versions",
		},
		{
			name:            "Both empty",
			localVersions:   []string{},
			remoteVersions:  []string{},
			expectedVersion: "",
			expectedFound:   false,
			description:     "Both versions lists empty",
		},
		{
			name:            "Version ordering test",
			localVersions:   []string{"1.1.0", "1.2.0", "1.10.0"}, // 1.10.0 > 1.2.0 > 1.1.0
			remoteVersions:  []string{"1.2.0", "1.10.0"},
			expectedVersion: "1.10.0", // Highest common version
			expectedFound:   true,
			description:     "Proper semantic version ordering",
		},
		{
			name:            "Invalid versions in mix",
			localVersions:   []string{"1.3.0", "invalid"},
			remoteVersions:  []string{"invalid", "2.0.0"},
			expectedVersion: "invalid", // Exact match on invalid version
			expectedFound:   true,
			description:     "Invalid versions can still be exact matches",
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			version, found := findHighestCommonVersion(tc.localVersions, tc.remoteVersions)
			assert.Equal(t, tc.expectedFound, found, "Expected found=%t for %s", tc.expectedFound, tc.description)
			if tc.expectedFound {
				assert.Equal(t, tc.expectedVersion, version, "Expected version=%s for %s", tc.expectedVersion, tc.description)
			}
		})
	}
}

// Test edge cases and boundary conditions for all functions
func (s *VersionNegotiationSuite) TestEdgeCases() {
	s.T().Run("Nil slices", func(t *testing.T) {
		// Test with nil slices
		version, found := findHighestCommonVersion(nil, nil)
		assert.False(t, found)
		assert.Empty(t, version)

		version, found = findCompatibleVersion(nil, nil)
		assert.False(t, found)
		assert.Empty(t, version)
	})

	s.T().Run("Large version numbers", func(t *testing.T) {
		// Test with large version numbers
		result := compareVersions("999.999.999", "1000.0.0")
		assert.Equal(t, -1, result)

		compatible := areVersionsCompatible("999.0.0", "999.1.0")
		assert.True(t, compatible)
	})

	s.T().Run("Special characters in versions", func(t *testing.T) {
		// Test with special characters (invalid formats)
		compatible := areVersionsCompatible("1.3.0-beta", "1.3.0-alpha")
		assert.True(t, compatible) // Invalid formats assume compatible

		result := compareVersions("1.3.0-beta", "1.3.0-alpha")
		assert.Equal(t, 1, result) // String comparison: "beta" > "alpha"
	})

	s.T().Run("Zero padding in versions", func(t *testing.T) {
		// Test versions with zero padding
		compatible := areVersionsCompatible("01.03.00", "1.3.0")
		assert.True(t, compatible) // Both are valid, zero-padded numbers parse correctly

		result := compareVersions("01.03.00", "1.3.0")
		assert.Equal(t, 0, result) // Zero-padded versions are equivalent when parsed
	})
}

// Test performance characteristics with large inputs
func (s *VersionNegotiationSuite) TestPerformance() {
	s.T().Run("Large version lists", func(t *testing.T) {
		// Create large version lists
		localVersions := make([]string, 100)
		remoteVersions := make([]string, 100)

		for i := 0; i < 100; i++ {
			localVersions[i] = fmt.Sprintf("1.%d.0", i)
			remoteVersions[i] = fmt.Sprintf("1.%d.0", i+50) // Partial overlap
		}

		// Should find highest common version efficiently
		version, found := findHighestCommonVersion(localVersions, remoteVersions)
		assert.True(t, found)
		assert.Equal(t, "1.99.0", version) // Highest overlap
	})
}

// Test real-world scenarios based on SPINE specification
func (s *VersionNegotiationSuite) TestRealWorldScenarios() {
	s.T().Run("SPINE 1.3.0 discovery", func(t *testing.T) {
		// Typical SPINE discovery scenario
		localVersions := []string{"1.3.0"}   // Modern device
		remoteVersions := []string{"1.2.0"}  // Older device

		version, found := findHighestCommonVersion(localVersions, remoteVersions)
		assert.True(t, found)
		assert.Equal(t, "1.3.0", version) // Should negotiate to higher compatible version
	})

	s.T().Run("Future SPINE major version", func(t *testing.T) {
		// Future SPINE major version incompatibility
		localVersions := []string{"1.3.0"}   // Current device
		remoteVersions := []string{"2.0.0"}  // Future device

		version, found := findHighestCommonVersion(localVersions, remoteVersions)
		assert.False(t, found) // Should be incompatible
		assert.Empty(t, version)
	})

	s.T().Run("Multiple SPINE versions support", func(t *testing.T) {
		// Device supporting multiple SPINE versions
		localVersions := []string{"1.1.0", "1.2.0", "1.3.0"}
		remoteVersions := []string{"1.2.0", "1.3.0", "1.4.0"}

		version, found := findHighestCommonVersion(localVersions, remoteVersions)
		assert.True(t, found)
		assert.Equal(t, "1.3.0", version) // Should pick highest common
	})
}