# Version Management in SPINE and spine-go

**Last Updated:** 2025-06-25  
**Status:** Active  
**Purpose:** Comprehensive analysis of version management challenges and architectural responsibilities

## Change History

### 2025-06-25
- Initial analysis of version management in SPINE
- Clarified architectural responsibilities between foundation and use case layers
- Documented protocol version validation gaps
- Analyzed real-world version compliance issues

## Executive Summary

**Critical Finding:** SPINE has fundamental version management gaps at both protocol and use case levels. However, spine-go correctly implements its responsibilities as a foundation library.

**Key Understanding:**
- **Protocol version validation** is missing (spine-go gap)
- **Use case version negotiation** belongs in use case implementations (NOT spine-go responsibility)
- **Real-world version compliance** is poor, requiring liberal parsing
- **Version infrastructure** exists but lacks selection mechanisms

## Table of Contents

1. [Architectural Responsibility Clarification](#architectural-responsibility-clarification)
2. [Protocol Version Management Issues](#protocol-version-management-issues)
3. [Use Case Version Management](#use-case-version-management)
4. [Real-World Version Compliance](#real-world-version-compliance)
5. [Implementation Recommendations](#implementation-recommendations)

---

## Architectural Responsibility Clarification

### Foundation Library vs Use Case Implementation

**Critical Understanding:** Version management spans two architectural layers with different responsibilities.

### Layer 1: Foundation Library (spine-go)
**Responsibilities:**
- ✅ **Store version information** (as opaque strings in data structures)
- ✅ **Provide AddUseCaseSupport API** to announce versions
- ✅ **Exchange version information** during discovery
- ✅ **Transport version data** between devices
- ✅ **Provide data structures** (UseCaseSupportType)
- ❌ **MISSING: Protocol version validation** (critical gap)

**NOT Responsible For:**
- ❌ Use case version parsing or semantic versioning
- ❌ Use case version compatibility checking
- ❌ Use case version negotiation algorithms
- ❌ Selecting which use case version to use
- ❌ Tracking active use case versions per connection

### Layer 2: Use Case Implementations (e.g., eebus-go)
**Responsibilities:**
- ✅ Define version format for their specific use cases
- ✅ Parse version strings (e.g., semantic versioning)
- ✅ Implement version compatibility rules
- ✅ Negotiate which version to use
- ✅ Track active versions per connection
- ✅ Handle version-specific behavior differences

### Why This Separation Matters

**spine-go as Foundation:**
```go
// spine-go provides the plumbing
type UseCaseSupportType struct {
    UseCaseName    *UseCaseNameType
    UseCaseVersion *SpecificationVersionType  // Opaque string
    // ... transport infrastructure
}

// AddUseCaseSupport stores version info
func (d *DeviceLocal) AddUseCaseSupport(
    useCaseName UseCaseNameType,
    useCaseVersion SpecificationVersionType,
    // ...
) {
    // Store and announce - no interpretation
}
```

**Use case implementation adds intelligence:**
```go
// eebus-go or similar adds business logic
type EVChargingUseCase struct {
    spine           *spine.DeviceLocal
    supportedVersions []EVChargingVersion
    negotiatedVersions map[string]EVChargingVersion
}

func (uc *EVChargingUseCase) NegotiateVersion(
    remoteVersions []SpecificationVersionType,
) (EVChargingVersion, error) {
    // Parse, compare, select - use case specific logic
    for _, remote := range remoteVersions {
        if version := uc.parseVersion(remote); version.IsCompatible(uc.localVersion) {
            return version, nil
        }
    }
    return EVChargingVersion{}, errors.New("no compatible version")
}
```

---

## Protocol Version Management Issues

### Critical Gap: No Protocol Version Validation

**SPINE Specification Requirement:**
> "The specificationVersion element SHALL be used in the header"
> "Different major versions have different compatibility groups"

**Current spine-go Implementation:**
```go
func (d *DeviceLocal) ProcessCmd(datagram model.DatagramType, ...) error {
    // NO check of datagram.Header.SpecificationVersion
    // Message processed regardless of version!
}
```

**Consequences:**
- ❌ **Silent failures** - Incompatible protocol versions processed incorrectly
- ❌ **Data corruption** - Version-specific fields misinterpreted
- ❌ **Security risks** - No validation of message format expectations

### Real-World Version String Reality

**Specification Expects:** Semantic versioning (e.g., "1.3.0")

**Reality:** Some devices send non-compliant strings:
- Empty strings: `""`
- Dots only: `"..."`
- Draft versions: `"draft"`
- RC versions: `"1.3.0-RC1"`
- Invalid formats: `"v1.3.0"`, `"1.3"`

**The Dilemma:**
- **Strict validation** would break compatibility with devices sending non-compliant strings
- **Liberal acceptance** violates specification requirements
- **Current approach** (no validation) accidentally enables broader compatibility

### Protocol Version Solution Approach

**Recommended Implementation:**
```go
type ProtocolVersionManager struct {
    localVersion       Version
    supportedVersions  []Version
    strictMode         bool  // Configuration option
    validationStats    *ValidationStats
}

func (pvm *ProtocolVersionManager) ValidateMessage(
    header *model.HeaderType,
) error {
    if header.SpecificationVersion == nil {
        return fmt.Errorf("missing specificationVersion")
    }
    
    version, err := pvm.parseVersion(*header.SpecificationVersion)
    if err != nil {
        if pvm.strictMode {
            return fmt.Errorf("invalid specificationVersion: %w", err)
        } else {
            // Liberal mode: log and continue
            pvm.validationStats.RecordNonCompliant(*header.SpecificationVersion)
            return nil
        }
    }
    
    if !pvm.isCompatible(version) {
        return fmt.Errorf("incompatible protocol version: %s", version)
    }
    
    return nil
}
```

**Benefits:**
- **Configurable strictness** for different deployment scenarios
- **Monitoring compliance** in real-world deployments
- **Migration path** to stricter validation over time

---

## Use Case Version Management

### The Specification Gap

**What SPINE Provides:**
- Storage mechanism for use case versions
- Transport for version announcements
- Discovery exchange of version information

**What SPINE Doesn't Provide:**
- Version negotiation protocol
- Selection mechanisms
- Compatibility rules
- Active version tracking

### Real-World Use Case Version Chaos

**Example Scenario:**
```json
{
  "useCaseSupport": [
    {
      "useCaseName": "optimizationOfSelfConsumptionDuringEvCharging",
      "useCaseVersion": "1.0.1",  // Legacy version
      "scenarioSupport": [1, 2, 3]
    },
    {
      "useCaseName": "optimizationOfSelfConsumptionDuringEvCharging", 
      "useCaseVersion": "2.0.0",  // New incompatible version
      "scenarioSupport": [1, 2, 3, 4, 5]
    }
  ]
}
```

**Questions with No Specification Answers:**
- Which version should be used?
- How to negotiate between devices?
- What happens if versions are incompatible?
- How to handle partial compatibility?

### Use Case Implementation Pattern

**Example: Energy Management Use Case**
```go
// In eebus-go or similar use case implementation
type EnergyManagementUseCase struct {
    spine             *spine.DeviceLocal
    supportedVersions []EMVersion
    activeVersions    map[string]EMVersion  // Per device
    negotiationRules  VersionNegotiationRules
}

func (em *EnergyManagementUseCase) OnDeviceDiscovered(
    remoteDevice spine.DeviceRemote,
) error {
    // 1. Get remote device's announced versions
    remoteVersions := remoteDevice.UseCaseSupport("energyManagement")
    
    // 2. Apply use case specific negotiation logic
    selectedVersion, err := em.negotiateVersion(remoteVersions)
    if err != nil {
        return fmt.Errorf("version negotiation failed: %w", err)
    }
    
    // 3. Track active version for this device
    em.activeVersions[remoteDevice.Address()] = selectedVersion
    
    // 4. Configure behavior based on negotiated version
    return em.configureForVersion(remoteDevice, selectedVersion)
}

func (em *EnergyManagementUseCase) negotiateVersion(
    remoteVersions []model.SpecificationVersionType,
) (EMVersion, error) {
    // Use case specific logic:
    // - Parse version strings
    // - Apply compatibility rules
    // - Select optimal version
    // - Handle conflicts
    
    for _, remote := range remoteVersions {
        for _, local := range em.supportedVersions {
            if local.IsCompatibleWith(em.parseVersion(remote)) {
                return local, nil
            }
        }
    }
    
    return EMVersion{}, errors.New("no compatible version found")
}
```

### Version Management Best Practices

#### For Use Case Implementers:

1. **Define Clear Version Semantics**
```go
type UseCaseVersion struct {
    Major int  // Breaking changes
    Minor int  // New features, backward compatible
    Patch int  // Bug fixes
}

func (v UseCaseVersion) IsCompatibleWith(other UseCaseVersion) bool {
    // Same major version = compatible
    return v.Major == other.Major
}
```

2. **Implement Robust Negotiation**
```go
func SelectBestVersion(local, remote []UseCaseVersion) (UseCaseVersion, error) {
    // Prefer highest compatible version
    // Handle edge cases (empty lists, no compatibility)
    // Document selection algorithm
}
```

3. **Handle Version-Specific Behavior**
```go
func (uc *UseCase) ProcessMessage(msg Message, version UseCaseVersion) error {
    switch version.Major {
    case 1:
        return uc.processV1(msg)
    case 2:
        return uc.processV2(msg)
    default:
        return fmt.Errorf("unsupported version: %s", version)
    }
}
```

---

## Real-World Version Compliance

### Version String Compliance Analysis

**Specification Format:** `major.minor.revision` (e.g., "1.3.0")

**Real-World Reality:**
```
Compliant Examples:
- "1.3.0" ✅
- "1.2.0" ✅
- "2.0.0" ✅

Non-Compliant Examples Observed:
- "" (empty) ❌
- "..." (dots only) ❌
- "draft" ❌
- "1.3.0-RC1" ❌
- "v1.3.0" ❌
- "1.3" (missing patch) ❌
```

### The Liberal Validation Dilemma

**Strict Validation Impact:**
- Would reject devices with non-compliant version strings
- Break compatibility with some deployed systems
- Reduce interoperability

**Liberal Validation Impact:**
- Accept non-compliant devices
- Enable broader ecosystem compatibility
- Violate specification requirements

**Recommended Approach: Configurable Validation**
```go
type ValidationMode int

const (
    StrictMode     ValidationMode = iota  // Reject non-compliant
    LiberalMode                           // Accept with warnings
    MonitoringMode                        // Accept all, log compliance
)

func (vm *VersionManager) SetValidationMode(mode ValidationMode) {
    vm.mode = mode
}
```

---

## Implementation Recommendations

### For spine-go (Foundation Layer)

#### 1. Add Protocol Version Validation (Critical)
```go
// P0: Implement protocol version checking
func (d *DeviceLocal) ProcessCmd(datagram model.DatagramType, ...) error {
    if err := d.versionManager.ValidateProtocolVersion(datagram.Header); err != nil {
        return fmt.Errorf("protocol version validation failed: %w", err)
    }
    // Continue with existing processing
}
```

#### 2. Provide Liberal Version Parsing
```go
// Support real-world version strings
func ParseVersionLiberally(versionStr string) (Version, error) {
    // Handle common non-compliant formats
    // Log compliance statistics
    // Provide migration path to strict parsing
}
```

#### 3. Add Version Monitoring
```go
// Track version compliance in deployments
type VersionStats struct {
    CompliantVersions    map[string]int
    NonCompliantVersions map[string]int
    ParseErrors         []string
}
```

### For Use Case Implementations

#### 1. Implement Version Negotiation
```go
// Each use case must implement its own negotiation
type UseCaseVersionNegotiator interface {
    NegotiateVersion(local, remote []Version) (Version, error)
    IsCompatible(v1, v2 Version) bool
    ParseVersion(versionStr string) (Version, error)
}
```

#### 2. Handle Version-Specific Behavior
```go
// Version-aware message processing
func (uc *UseCase) ProcessMessage(
    msg Message, 
    sourceVersion Version,
) error {
    // Adapt behavior based on negotiated version
}
```

#### 3. Provide Version Migration
```go
// Help users upgrade between versions
type VersionMigrator interface {
    CanMigrate(from, to Version) bool
    Migrate(data interface{}, from, to Version) (interface{}, error)
}
```

### For System Integrators

#### 1. Plan Version Strategy
- Define supported version ranges
- Test with multiple version combinations
- Plan migration paths for upgrades

#### 2. Monitor Version Compliance
- Log version strings in deployments
- Track compliance statistics
- Plan remediation for non-compliant devices

#### 3. Design for Version Evolution
- Avoid tight coupling to specific versions
- Plan for backward compatibility
- Design upgrade procedures

---

## Conclusion

### Current State Assessment

**spine-go Foundation Layer:**
- ✅ **Correctly provides** version storage and transport infrastructure
- ✅ **Appropriate scope** - doesn't overstep into use case responsibilities
- ❌ **Missing critical feature** - protocol version validation
- ✅ **Accidentally liberal** - accepts non-compliant versions (broader compatibility)

**Use Case Layer Gaps:**
- ❌ **No standard negotiation** - each implementation must build own
- ❌ **No compatibility framework** - inconsistent behavior across vendors
- ❌ **No migration tools** - difficult to upgrade between versions

### Strategic Recommendations

#### Immediate (P0):
1. **Implement protocol version validation** in spine-go
2. **Add configurable strictness** for real-world compatibility
3. **Add version monitoring** to understand compliance landscape

#### Medium-term (P1):
1. **Develop use case version frameworks** in eebus-go
2. **Create version negotiation patterns** for common scenarios
3. **Build compliance monitoring tools** for deployments

#### Long-term (P2):
1. **Advocate for specification improvements** in version management
2. **Develop industry standards** for version negotiation
3. **Create migration frameworks** for version evolution

### The Bottom Line

**Version management in SPINE requires a two-layer approach:**
- **spine-go handles protocol version validation** (missing but required)
- **Use case implementations handle use case version negotiation** (correctly delegated)

The current architecture is sound, but both layers need strengthening to handle the realities of version diversity in production deployments.

---

**Related Documents:**
- [../detailed-analysis/SPINE_SPECIFICATION_ANALYSIS.md](../detailed-analysis/SPINE_SPECIFICATION_ANALYSIS.md) - Section 7 & 8 on version analysis
- [../detailed-analysis/IMPROVEMENT_ROADMAP.md](../detailed-analysis/IMPROVEMENT_ROADMAP.md) - Protocol version validation implementation
- [../detailed-analysis/SPEC_DEVIATIONS.md](../detailed-analysis/SPEC_DEVIATIONS.md) - Version validation requirements