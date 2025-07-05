package spine

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestVersionErrorsSuite(t *testing.T) {
	suite.Run(t, new(VersionErrorsSuite))
}

type VersionErrorsSuite struct {
	suite.Suite
}

// Test VersionIncompatibilityError creation and formatting
func (s *VersionErrorsSuite) TestVersionIncompatibilityError_Creation() {
	localVersions := []string{"1.3.0", "1.2.0"}
	remoteVersions := []string{"2.0.0", "2.1.0"}

	err := NewVersionIncompatibilityError(localVersions, remoteVersions)

	// Test structure
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), localVersions, err.LocalVersions)
	assert.Equal(s.T(), remoteVersions, err.RemoteVersions)

	// Test error message format
	expectedMsg := "no compatible protocol version found. Local: [1.3.0 1.2.0], Remote: [2.0.0 2.1.0]"
	assert.Equal(s.T(), expectedMsg, err.Error())
}

// Test error message formatting with different version combinations
func (s *VersionErrorsSuite) TestVersionIncompatibilityError_MessageFormatting() {
	testCases := []struct {
		name           string
		localVersions  []string
		remoteVersions []string
		expectedMsg    string
	}{
		{
			name:           "Single versions",
			localVersions:  []string{"1.3.0"},
			remoteVersions: []string{"2.0.0"},
			expectedMsg:    "no compatible protocol version found. Local: [1.3.0], Remote: [2.0.0]",
		},
		{
			name:           "Multiple local, single remote",
			localVersions:  []string{"1.2.0", "1.3.0"},
			remoteVersions: []string{"2.0.0"},
			expectedMsg:    "no compatible protocol version found. Local: [1.2.0 1.3.0], Remote: [2.0.0]",
		},
		{
			name:           "Empty local versions",
			localVersions:  []string{},
			remoteVersions: []string{"2.0.0"},
			expectedMsg:    "no compatible protocol version found. Local: [], Remote: [2.0.0]",
		},
		{
			name:           "Empty remote versions",
			localVersions:  []string{"1.3.0"},
			remoteVersions: []string{},
			expectedMsg:    "no compatible protocol version found. Local: [1.3.0], Remote: []",
		},
		{
			name:           "Both empty",
			localVersions:  []string{},
			remoteVersions: []string{},
			expectedMsg:    "no compatible protocol version found. Local: [], Remote: []",
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			err := NewVersionIncompatibilityError(tc.localVersions, tc.remoteVersions)
			assert.Equal(t, tc.expectedMsg, err.Error())
		})
	}
}

// Test direct error type detection
func (s *VersionErrorsSuite) TestIsVersionIncompatibilityError_DirectType() {
	// Test with actual VersionIncompatibilityError
	err := NewVersionIncompatibilityError([]string{"1.3.0"}, []string{"2.0.0"})
	assert.True(s.T(), IsVersionIncompatibilityError(err))

	// Test with nil
	assert.False(s.T(), IsVersionIncompatibilityError(nil))

	// Test with other error types
	otherErr := errors.New("some other error")
	assert.False(s.T(), IsVersionIncompatibilityError(otherErr))

	customErr := fmt.Errorf("custom error: %w", errors.New("wrapped"))
	assert.False(s.T(), IsVersionIncompatibilityError(customErr))
}

// Test wrapped error detection using errors.As
func (s *VersionErrorsSuite) TestIsVersionIncompatibilityError_WrappedType() {
	originalErr := NewVersionIncompatibilityError([]string{"1.3.0"}, []string{"2.0.0"})

	// Test single level wrapping
	wrappedErr := fmt.Errorf("wrapped: %w", originalErr)
	assert.True(s.T(), IsVersionIncompatibilityError(wrappedErr))

	// Test multiple level wrapping
	doubleWrappedErr := fmt.Errorf("outer: %w", wrappedErr)
	assert.True(s.T(), IsVersionIncompatibilityError(doubleWrappedErr))

	// Test wrapping with other error (should not be detected)
	otherErr := errors.New("other error")
	wrappedOtherErr := fmt.Errorf("wrapped: %w", otherErr)
	assert.False(s.T(), IsVersionIncompatibilityError(wrappedOtherErr))
}

// Test string-based error detection (for SPINE error system compatibility)
func (s *VersionErrorsSuite) TestIsVersionIncompatibilityError_StringBased() {
	// Test SPINE error format: "Error X: message"
	spineErr := errors.New("Error 1: no compatible protocol version found. Local: [1.3.0], Remote: [2.0.0]")
	assert.True(s.T(), IsVersionIncompatibilityError(spineErr))

	// Test partial message containing signature
	partialErr := errors.New("Protocol error: no compatible protocol version found between devices")
	assert.True(s.T(), IsVersionIncompatibilityError(partialErr))

	// Test case sensitivity
	caseErr := errors.New("No Compatible Protocol Version Found")
	assert.False(s.T(), IsVersionIncompatibilityError(caseErr))

	// Test similar but different message
	similarErr := errors.New("no compatible version found")
	assert.False(s.T(), IsVersionIncompatibilityError(similarErr))

	// Test substring in different context - this should NOT match
	differentContextErr := errors.New("found no compatible protocol version information")
	assert.False(s.T(), IsVersionIncompatibilityError(differentContextErr))

	// Test containing exact signature in different context - this SHOULD match
	exactSignatureErr := errors.New("Discovery failed: no compatible protocol version found during negotiation")
	assert.True(s.T(), IsVersionIncompatibilityError(exactSignatureErr))
}

// Test edge cases and boundary conditions
func (s *VersionErrorsSuite) TestIsVersionIncompatibilityError_EdgeCases() {
	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "Nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "Empty error message",
			err:      errors.New(""),
			expected: false,
		},
		{
			name:     "Only whitespace",
			err:      errors.New("   "),
			expected: false,
		},
		{
			name:     "Exact signature string",
			err:      errors.New("no compatible protocol version found"),
			expected: true,
		},
		{
			name:     "Signature with prefix",
			err:      errors.New("Error: no compatible protocol version found"),
			expected: true,
		},
		{
			name:     "Signature with suffix",
			err:      errors.New("no compatible protocol version found with details"),
			expected: true,
		},
		{
			name:     "Multiple occurrences",
			err:      errors.New("no compatible protocol version found and no compatible protocol version found again"),
			expected: true,
		},
		{
			name:     "Case sensitive - wrong case",
			err:      errors.New("NO COMPATIBLE PROTOCOL VERSION FOUND"),
			expected: false,
		},
		{
			name:     "Signature in middle of long message",
			err:      errors.New("During discovery process, encountered issue: no compatible protocol version found. Unable to proceed."),
			expected: true,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			result := IsVersionIncompatibilityError(tc.err)
			assert.Equal(t, tc.expected, result, "Error: %v", tc.err)
		})
	}
}

// Test behavior with nil slices
func (s *VersionErrorsSuite) TestVersionIncompatibilityError_NilSlices() {
	// Test with nil slices
	err := NewVersionIncompatibilityError(nil, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), err.LocalVersions)
	assert.Nil(s.T(), err.RemoteVersions)

	expectedMsg := "no compatible protocol version found. Local: [], Remote: []"
	assert.Equal(s.T(), expectedMsg, err.Error())

	// Test mixed nil and empty
	err2 := NewVersionIncompatibilityError(nil, []string{"2.0.0"})
	assert.Nil(s.T(), err2.LocalVersions)
	assert.Equal(s.T(), []string{"2.0.0"}, err2.RemoteVersions)
}

// Test integration with actual error scenarios
func (s *VersionErrorsSuite) TestVersionIncompatibilityError_IntegrationScenarios() {
	// Scenario 1: Major version incompatibility
	err1 := NewVersionIncompatibilityError([]string{"1.3.0"}, []string{"2.0.0"})
	assert.True(s.T(), IsVersionIncompatibilityError(err1))
	assert.Contains(s.T(), err1.Error(), "1.3.0")
	assert.Contains(s.T(), err1.Error(), "2.0.0")

	// Scenario 2: No common versions in complex case
	err2 := NewVersionIncompatibilityError(
		[]string{"1.0.0", "1.1.0", "1.2.0"}, 
		[]string{"2.0.0", "2.1.0", "3.0.0"},
	)
	assert.True(s.T(), IsVersionIncompatibilityError(err2))
	assert.Contains(s.T(), err2.Error(), "1.0.0 1.1.0 1.2.0")
	assert.Contains(s.T(), err2.Error(), "2.0.0 2.1.0 3.0.0")

	// Scenario 3: Empty version lists (protocol negotiation failure)
	err3 := NewVersionIncompatibilityError([]string{}, []string{})
	assert.True(s.T(), IsVersionIncompatibilityError(err3))
	assert.Contains(s.T(), err3.Error(), "Local: []")
	assert.Contains(s.T(), err3.Error(), "Remote: []")
}