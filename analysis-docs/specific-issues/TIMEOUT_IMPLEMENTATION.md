# Timeout Implementation Analysis and Guidelines

**Last Updated:** 2025-07-05  
**Status:** Active  
**Related Documents:**
- [SPINE_SPECIFICATIONS_ANALYSIS.md](../detailed-analysis/SPINE_SPECIFICATIONS_ANALYSIS.md) - Section 10.4 Timeout Specification Ambiguities
- [SPEC_DEVIATIONS.md](../detailed-analysis/SPEC_DEVIATIONS.md) - Section 2 Error Response Timing Analysis

## Purpose

This document provides comprehensive analysis of timeout handling in spine-go, explains the current implementation decisions, and provides practical guidance for applications that need timeout detection.

## Executive Summary

**Key Finding:** spine-go's selective timeout implementation is SPEC-COMPLIANT and optimized for interoperability:
- ✅ **Write approval timeouts implemented** - Critical control path protected
- ❌ **Read request timeouts not implemented** - Optional per SPINE spec (MAY requirement)
- ✅ **Maximum interoperability** - Compatible with all possible SPINE implementations

## spine-go's Timeout Strategy

### Current Implementation Status

| Operation Type | Timeout Detection | Default Timeout | Configurable | Rationale |
|---------------|------------------|-----------------|--------------|-----------|
| **Write Approvals** | ✅ Implemented | 10 seconds | ✅ Yes | Critical control path protection |
| **Read Requests** | ❌ Not implemented | N/A | N/A | SPINE spec compliance (MAY) |
| **Write Requests** | ❌ Not implemented | N/A | N/A | No spec requirement |
| **Notifications** | ❌ Not implemented | N/A | N/A | No spec requirement |

### Why Read Request Timeouts Are Not Implemented

**SPINE Specification Analysis:**
- **Timeout detection is OPTIONAL**: Spec uses "MAY" language (RFC 2119)
- **No defined timeout behavior**: Spec doesn't specify what to do on timeout
- **No standard recovery mechanism**: No retry or cleanup procedures defined

**Interoperability Benefits:**
- **Maximum compatibility**: Works with all SPINE implementations
- **No false timeouts**: Avoids breaking slow but functional devices
- **Predictable behavior**: Always waits for response or connection loss

**Technical Rationale:**
```
Real-World Multi-Vendor Scenario:
- spine-go: No read timeouts (waits indefinitely)
- Vendor A: 10-second strict timeouts
- Vendor B: 30-second conservative timeouts  
- Vendor C: Configurable timeouts (user-defined)

Result: spine-go works with ALL vendors
```

## Write Approval Timeout Implementation

### Technical Details

**Implementation Location:** `spine/feature_local.go:194-201`

**Mechanism:**
```go
// Timer-based timeout with automatic cleanup
newTimer := time.AfterFunc(r.writeTimeout, func() {
    r.muxResponseCB.Lock()
    delete(r.pendingWriteApprovals[ski], *msg.RequestHeader.MsgCounter)
    r.muxResponseCB.Unlock()
    
    err := model.NewErrorTypeFromString("write not approved in time by application")
    _ = msg.FeatureRemote.Device().Sender().ResultError(msg.RequestHeader, r.Address(), err)
})
```

**Key Features:**
- **Default timeout**: 10 seconds (`defaultMaxResponseDelay`)
- **Configurable**: Via `SetWriteApprovalTimeout(duration)`
- **Automatic cleanup**: Removes pending approval on timeout
- **Error response**: Sends error code 1 with descriptive message
- **Timer cancellation**: Stopped when approval/denial received
- **Thread safety**: Protected by mutexes

### Configuration API

```go
// Set custom write approval timeout
feature.SetWriteApprovalTimeout(time.Second * 30) // 30 seconds

// Approve or deny write within timeout
feature.ApproveOrDenyWrite(msg, errorType)
```

## Application-Level Timeout Implementation

For applications that need timeout detection for read requests, implement at the application level where requirements are better defined.

### Basic Timeout Pattern

```go
// Example: Application-level timeout for read requests
func requestWithTimeout(device DeviceInterface, timeout time.Duration) error {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    
    done := make(chan error, 1)
    
    go func() {
        // Perform spine-go read operation
        _, err := device.RequestRemoteData(...)
        done <- err
    }()
    
    select {
    case err := <-done:
        // Operation completed
        return err
    case <-ctx.Done():
        // Timeout occurred
        return fmt.Errorf("read request timed out after %v", timeout)
    }
}
```

### Advanced Timeout with Retry

```go
// Example: Timeout with exponential backoff retry
func requestWithRetry(device DeviceInterface, maxAttempts int) error {
    for attempt := 1; attempt <= maxAttempts; attempt++ {
        timeout := time.Duration(attempt) * 10 * time.Second // Exponential timeout
        
        err := requestWithTimeout(device, timeout)
        if err == nil {
            return nil // Success
        }
        
        if !isTimeoutError(err) {
            return err // Non-timeout error, don't retry
        }
        
        if attempt < maxAttempts {
            log.Warnf("Attempt %d timed out, retrying in %v", attempt, timeout/2)
            time.Sleep(timeout / 2) // Wait before retry
        }
    }
    
    return fmt.Errorf("all %d attempts timed out", maxAttempts)
}

func isTimeoutError(err error) bool {
    return strings.Contains(err.Error(), "timed out")
}
```

### Timeout Configuration Guidelines

**Conservative Values (Recommended):**
- **Local network**: 30-60 seconds
- **Wide area network**: 60-120 seconds
- **Battery-powered devices**: 120+ seconds
- **Critical operations**: Consider no timeout

**Network Latency Considerations:**
```go
// Calculate timeout with latency buffer
baseTimeout := time.Second * 10        // Base SPINE timeout
networkLatency := time.Second * 5      // Estimated network latency
deviceProcessing := time.Second * 2    // Device processing time
totalTimeout := baseTimeout + networkLatency + deviceProcessing // 17 seconds
```

**Device-Specific Timeouts:**
```go
// Different timeouts based on device capabilities
func getTimeoutForDevice(deviceType string) time.Duration {
    switch deviceType {
    case "battery_powered":
        return time.Minute * 2  // Very conservative
    case "mains_powered":
        return time.Second * 30 // Moderate
    case "high_performance":
        return time.Second * 15 // Responsive
    default:
        return time.Minute * 1  // Safe default
    }
}
```

## Best Practices

### Timeout Implementation

1. **Use conservative values** - Avoid breaking slow but functional devices
2. **Add network latency buffer** - Account for network delays
3. **Consider device context** - Battery devices need longer timeouts
4. **Log timeout events** - Help with debugging and monitoring
5. **Implement proper cleanup** - Cancel operations on timeout

### Error Handling

```go
// Distinguish between timeout and other errors
func handleRequestError(err error) {
    if isTimeoutError(err) {
        // Timeout - device may be slow or unresponsive
        log.Warn("Device response timeout - may be slow or disconnected")
        // Consider: retry, use cached data, alert user
    } else if isConnectionError(err) {
        // Connection issue - network problem
        log.Error("Network connection error")
        // Consider: reconnect, check network
    } else {
        // Protocol or application error
        log.Error("Request failed: %v", err)
        // Consider: protocol-specific handling
    }
}
```

### Retry Logic

```go
// Safe retry with backoff
func retryWithBackoff(operation func() error, maxAttempts int) error {
    for attempt := 1; attempt <= maxAttempts; attempt++ {
        err := operation()
        if err == nil {
            return nil
        }
        
        if !shouldRetry(err) {
            return err // Don't retry non-transient errors
        }
        
        if attempt < maxAttempts {
            backoff := time.Duration(attempt*attempt) * time.Second // Quadratic backoff
            time.Sleep(backoff)
        }
    }
    return fmt.Errorf("operation failed after %d attempts", maxAttempts)
}

func shouldRetry(err error) bool {
    // Only retry timeouts and connection errors
    return isTimeoutError(err) || isConnectionError(err)
}
```

## Alternative Approaches for Unresponsive Device Detection

### 1. Heartbeat Monitoring

```go
// Use spine-go's built-in heartbeat feature
heartbeatManager := device.HeartbeatManager()
heartbeatManager.SetHeartbeatTimeout(time.Minute * 2)

// Monitor heartbeat events
device.AddEventCallback(func(event api.EventType) {
    if event.EventType == api.EventTypeDeviceHeartbeat {
        log.Info("Device heartbeat received")
    } else if event.EventType == api.EventTypeDeviceDisconnected {
        log.Warn("Device heartbeat timeout - device may be offline")
    }
})
```

### 2. Connection State Monitoring

```go
// Monitor device connection status
device.AddEventCallback(func(event api.EventType) {
    switch event.EventType {
    case api.EventTypeDeviceConnected:
        log.Info("Device connected")
    case api.EventTypeDeviceDisconnected:
        log.Warn("Device disconnected")
    case api.EventTypeDeviceDestroyed:
        log.Error("Device destroyed")
    }
})
```

### 3. Application-Level Health Checks

```go
// Periodic health check implementation
func healthCheck(device DeviceInterface) {
    ticker := time.NewTicker(time.Minute * 5) // Check every 5 minutes
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            err := requestWithTimeout(device, time.Second*30)
            if err != nil {
                log.Warnf("Health check failed: %v", err)
                // Consider: mark device as unhealthy, alert monitoring
            } else {
                log.Debug("Health check passed")
            }
        }
    }
}
```

## Testing Timeout Behavior

### Unit Test Pattern

```go
func TestApplicationTimeout(t *testing.T) {
    // Setup mock device that never responds
    mockDevice := &MockDevice{
        ShouldRespond: false,
    }
    
    // Test timeout behavior
    start := time.Now()
    err := requestWithTimeout(mockDevice, time.Second*2)
    duration := time.Since(start)
    
    // Verify timeout occurred
    assert.Error(t, err)
    assert.True(t, isTimeoutError(err))
    assert.InDelta(t, 2.0, duration.Seconds(), 0.5) // 2s ± 0.5s
}
```

### Integration Test Pattern

```go
func TestRealDeviceTimeout(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    // Test with real device
    device := setupRealDevice(t)
    defer device.Disconnect()
    
    // Test various timeout scenarios
    testCases := []struct {
        name    string
        timeout time.Duration
        expectTimeout bool
    }{
        {"aggressive", time.Second * 5, true},
        {"moderate", time.Second * 30, false},
        {"conservative", time.Minute * 2, false},
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            err := requestWithTimeout(device, tc.timeout)
            if tc.expectTimeout {
                assert.Error(t, err)
                assert.True(t, isTimeoutError(err))
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## Recommendations

### For Library Users

1. **Don't implement read timeouts unless absolutely necessary**
2. **Use spine-go's write approval timeouts for control operations**
3. **Implement application-level timeouts when needed**
4. **Use conservative timeout values (30+ seconds)**
5. **Monitor device health via heartbeat and connection state**
6. **Log timeout events for debugging**

### For System Designers

1. **Assume no protocol-level timeout detection** - Safest for multi-vendor systems
2. **Design for slow devices** - Some devices take time to respond
3. **Use heartbeat monitoring** - More reliable than timeout detection
4. **Implement proper error handling** - Distinguish timeout from other errors
5. **Consider offline operation** - Cache data for when devices are slow

### For spine-go Development

1. **Keep current timeout strategy** - Maintains interoperability
2. **Document timeout behavior clearly** - Help users understand choices
3. **Provide timeout utility functions** - Helper functions for common patterns
4. **Consider timeout middleware** - Optional application-level timeout wrapper

## Conclusion

spine-go's selective timeout implementation strikes the optimal balance between functionality and interoperability. By implementing timeouts only where critical (write approvals) and avoiding them where optional (read requests), spine-go maximizes compatibility with all possible SPINE implementations while protecting the most important operations.

Applications requiring timeout detection should implement it at the application level where requirements are better defined and recovery mechanisms can be properly designed for the specific use case.

---

## Document History

### 2025-07-05
- Initial document creation based on timeout analysis
- Comprehensive implementation guidance and examples
- Best practices for application-level timeout handling
- Testing patterns and real-world scenarios