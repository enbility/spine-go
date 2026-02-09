package spine

import (
	"errors"
	"fmt"
	"strings"
)

// VersionIncompatibilityError represents an error when no compatible protocol version is found
type VersionIncompatibilityError struct {
	LocalVersions  []string
	RemoteVersions []string
}

func (e *VersionIncompatibilityError) Error() string {
	return fmt.Sprintf("no compatible protocol version found. Local: %v, Remote: %v", 
		e.LocalVersions, e.RemoteVersions)
}

// NewVersionIncompatibilityError creates a new version incompatibility error
func NewVersionIncompatibilityError(localVersions, remoteVersions []string) *VersionIncompatibilityError {
	return &VersionIncompatibilityError{
		LocalVersions:  localVersions,
		RemoteVersions: remoteVersions,
	}
}

// IsVersionIncompatibilityError checks if an error is a version incompatibility error
// This handles both direct errors and errors wrapped by the SPINE error system
func IsVersionIncompatibilityError(err error) bool {
	if err == nil {
		return false
	}
	
	// Check for direct type
	var versionErr *VersionIncompatibilityError
	if errors.As(err, &versionErr) {
		return true
	}
	
	// Check for error message containing our signature text
	// This handles cases where SPINE error system wraps our error
	errMsg := err.Error()
	return strings.Contains(errMsg, "no compatible protocol version found")
}