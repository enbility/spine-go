package spine

import (
	"sort"
	"strings"
)

// findHighestCommonVersion finds the highest version supported by both local and remote devices
func findHighestCommonVersion(localVersions, remoteVersions []string) (string, bool) {
	if len(localVersions) == 0 || len(remoteVersions) == 0 {
		return "", false
	}

	// Create maps for O(1) lookup
	localMap := make(map[string]bool)
	for _, v := range localVersions {
		localMap[v] = true
	}

	// Find all common versions
	var commonVersions []string
	for _, v := range remoteVersions {
		if localMap[v] {
			commonVersions = append(commonVersions, v)
		}
	}

	if len(commonVersions) == 0 {
		// No exact matches, check for compatible versions
		return findCompatibleVersion(localVersions, remoteVersions)
	}

	// Sort versions to find the highest
	sort.Slice(commonVersions, func(i, j int) bool {
		return compareVersions(commonVersions[i], commonVersions[j]) > 0
	})

	return commonVersions[0], true
}

// findCompatibleVersion tries to find compatible versions when no exact match exists
func findCompatibleVersion(localVersions, remoteVersions []string) (string, bool) {
	// For each combination, check if versions are compatible
	for _, localVer := range localVersions {
		for _, remoteVer := range remoteVersions {
			if areVersionsCompatible(localVer, remoteVer) {
				// Return the higher version if compatible
				if compareVersions(localVer, remoteVer) > 0 {
					return localVer, true
				}
				return remoteVer, true
			}
		}
	}
	return "", false
}

// areVersionsCompatible checks if two versions are compatible based on major version
func areVersionsCompatible(v1, v2 string) bool {
	major1, _, _, valid1 := parseSemanticVersion(v1)
	major2, _, _, valid2 := parseSemanticVersion(v2)

	// If either version is not valid semantic version, assume compatible
	if !valid1 || !valid2 {
		return true
	}

	// Major versions 0 and 1 are compatible with each other
	if major1 <= 1 && major2 <= 1 {
		return true
	}

	// For major versions >= 2, they must match exactly
	return major1 == major2
}

// compareVersions compares two semantic versions
// Returns: 1 if v1 > v2, -1 if v1 < v2, 0 if equal
func compareVersions(v1, v2 string) int {
	major1, minor1, patch1, valid1 := parseSemanticVersion(v1)
	major2, minor2, patch2, valid2 := parseSemanticVersion(v2)

	// Invalid versions are considered lower
	if !valid1 && !valid2 {
		return strings.Compare(v1, v2)
	}
	if !valid1 {
		return -1
	}
	if !valid2 {
		return 1
	}

	// Compare major
	if major1 != major2 {
		if major1 > major2 {
			return 1
		}
		return -1
	}

	// Compare minor
	if minor1 != minor2 {
		if minor1 > minor2 {
			return 1
		}
		return -1
	}

	// Compare patch
	if patch1 != patch2 {
		if patch1 > patch2 {
			return 1
		}
		return -1
	}

	return 0
}