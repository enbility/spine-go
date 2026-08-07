package spine

// lifecycle_issues_demo_test.go
//
// Verifies that lifecycle issues in spine-go's start/stop handling are fixed.
// Run with:
//   go test -v -run TestLifecycleIssues ./spine/
//   go test -v -race -run TestLifecycleIssues_Race ./spine/

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ---------------------------------------------------------------------------
// Shared helpers
// ---------------------------------------------------------------------------

func countGoroutines() int {
	return runtime.NumGoroutine()
}

func stabilize() {
	runtime.Gosched()
	time.Sleep(50 * time.Millisecond)
}

type trackingWriter struct {
	mu       sync.Mutex
	messages [][]byte
	calls    atomic.Int64
}

func (w *trackingWriter) WriteShipMessageWithPayload(msg []byte) {
	w.calls.Add(1)
	w.mu.Lock()
	defer w.mu.Unlock()
	w.messages = append(w.messages, msg)
}

func (w *trackingWriter) callCount() int64 {
	return w.calls.Load()
}

type slowEventHandler struct {
	mu       sync.Mutex
	events   []api.EventPayload
	started  atomic.Int64
	finished atomic.Int64
	delay    time.Duration
}

func newSlowEventHandler(delay time.Duration) *slowEventHandler {
	return &slowEventHandler{delay: delay}
}

func (h *slowEventHandler) HandleEvent(payload api.EventPayload) {
	h.started.Add(1)
	time.Sleep(h.delay)
	h.mu.Lock()
	h.events = append(h.events, payload)
	h.mu.Unlock()
	h.finished.Add(1)
}

// ---------------------------------------------------------------------------
// Test Suite
// ---------------------------------------------------------------------------

type LifecycleIssuesSuite struct {
	suite.Suite
}

func TestLifecycleIssues(t *testing.T) {
	suite.Run(t, new(LifecycleIssuesSuite))
}

// =========================================================================
// Issue 1: Close() stops all goroutines — no leak when used
// =========================================================================

func (s *LifecycleIssuesSuite) TestIssue1_CloseStopsAllGoroutines() {
	stabilize()
	before := countGoroutines()

	device := NewDeviceLocal(
		"Demo", "Model", "Serial", "Code", "device-1",
		model.DeviceTypeTypeEnergyManagementSystem,
		model.NetworkManagementFeatureSetTypeSmart)

	entity := NewEntityLocal(device, model.EntityTypeTypeCEM,
		[]model.AddressEntityType{1}, 4*time.Second)
	device.AddEntity(entity)

	diagFeature := NewFeatureLocal(
		entity.NextFeatureId(), entity,
		model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	diagFeature.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)
	entity.AddFeature(diagFeature)

	assert.True(s.T(), entity.HeartbeatManager().IsHeartbeatRunning())

	// Close() should clean up everything
	device.Close()

	stabilize()
	after := countGoroutines()

	assert.LessOrEqual(s.T(), after, before+1,
		"no goroutines should leak after Close()")
}

// =========================================================================
// Issue 2: CleanWriteApprovalCaches stops timers — no post-cleanup fire
// =========================================================================

func (s *LifecycleIssuesSuite) TestIssue2_TimersStoppedOnCleanup() {
	writer := &trackingWriter{}

	device, localEntity := createLocalDeviceAndEntity(1)
	_, serverFeature := createLocalFeatures(localEntity, model.FeatureTypeTypeLoadControl, model.FunctionTypeLoadControlLimitListData)

	serverFeature.AddWriteApprovalCallback(func(msg *api.Message) {})

	ski := "remote-1"
	sender := NewSender(writer)
	remoteDevice := createRemoteDevice(device, ski, sender)
	device.AddRemoteDeviceForSki(ski, remoteDevice)

	remoteFeature, _ := createRemoteEntityAndFeature(remoteDevice, 1,
		model.FeatureTypeTypeLoadControl, model.FunctionTypeLoadControlLimitListData)

	serverFeature.SetWriteApprovalTimeout(200 * time.Millisecond)

	msgCounter := model.MsgCounterType(42)
	msg := &api.Message{
		RequestHeader: &model.HeaderType{
			MsgCounter: util.Ptr(msgCounter),
			AddressSource: &model.FeatureAddressType{
				Device:  remoteDevice.Address(),
				Entity:  remoteDevice.Entity([]model.AddressEntityType{1}).Address().Entity,
				Feature: remoteFeature.Address().Feature,
			},
			AddressDestination: serverFeature.Address(),
		},
		DeviceRemote:  remoteDevice,
		EntityRemote:  remoteDevice.Entity([]model.AddressEntityType{1}),
		FeatureRemote: remoteFeature,
	}

	serverFeature.(*FeatureLocal).addPendingApproval(msg)

	callsBefore := writer.callCount()
	device.RemoveRemoteDevice(ski)

	// Wait past timer expiry
	time.Sleep(400 * time.Millisecond)

	callsAfter := writer.callCount()
	assert.Equal(s.T(), callsBefore, callsAfter,
		"timer should not fire after CleanWriteApprovalCaches stops it")
}

// =========================================================================
// Issue 3: StopHeartbeat waits for goroutine — synchronous stop
// =========================================================================

func (s *LifecycleIssuesSuite) TestIssue3_HeartbeatJoinedOnRemoval() {
	device := NewDeviceLocal(
		"Demo", "Model", "Serial", "Code", "device-1",
		model.DeviceTypeTypeEnergyManagementSystem,
		model.NetworkManagementFeatureSetTypeSmart)

	entity := NewEntityLocal(device, model.EntityTypeTypeCEM,
		[]model.AddressEntityType{1}, 4*time.Second)
	device.AddEntity(entity)

	diagFeature := NewFeatureLocal(
		entity.NextFeatureId(), entity,
		model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	diagFeature.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)
	entity.AddFeature(diagFeature)

	assert.True(s.T(), entity.HeartbeatManager().IsHeartbeatRunning())

	goroutinesBefore := countGoroutines()

	// RemoveEntity -> StopHeartbeat now blocks until goroutine exits
	device.RemoveEntity(entity)

	// Goroutine should be gone immediately (no sleep needed)
	goroutinesAfter := countGoroutines()
	assert.LessOrEqual(s.T(), goroutinesAfter, goroutinesBefore,
		"goroutine should exit synchronously when StopHeartbeat returns")
}

// =========================================================================
// Issue 4: No data race on stopHeartbeatC (run with -race)
// =========================================================================

type LifecycleRaceSuite struct {
	suite.Suite
}

func TestLifecycleIssues_Race(t *testing.T) {
	suite.Run(t, new(LifecycleRaceSuite))
}

func (s *LifecycleRaceSuite) TestIssue4_NoRaceOnStopHeartbeatChannel() {
	device := NewDeviceLocal(
		"Demo", "Model", "Serial", "Code", "device-1",
		model.DeviceTypeTypeEnergyManagementSystem,
		model.NetworkManagementFeatureSetTypeSmart)

	entity := NewEntityLocal(device, model.EntityTypeTypeCEM,
		[]model.AddressEntityType{1}, 4*time.Second)
	device.AddEntity(entity)

	diagFeature := NewFeatureLocal(
		entity.NextFeatureId(), entity,
		model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	diagFeature.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)
	entity.AddFeature(diagFeature)

	hbm := entity.HeartbeatManager()

	// Concurrent StartHeartbeat/IsHeartbeatRunning — should NOT trigger race detector.
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_ = hbm.StartHeartbeat()
		}()
		go func() {
			defer wg.Done()
			_ = hbm.IsHeartbeatRunning()
		}()
	}
	wg.Wait()

	hbm.StopHeartbeat()
	assert.False(s.T(), hbm.IsHeartbeatRunning())
}

// =========================================================================
// Issue 5: Ticker properly stopped — no resource leak
// =========================================================================

func (s *LifecycleIssuesSuite) TestIssue5_TickerStoppedOnExit() {
	device := NewDeviceLocal(
		"Demo", "Model", "Serial", "Code", "device-1",
		model.DeviceTypeTypeEnergyManagementSystem,
		model.NetworkManagementFeatureSetTypeSmart)

	entity := NewEntityLocal(device, model.EntityTypeTypeCEM,
		[]model.AddressEntityType{1}, 4*time.Second)
	device.AddEntity(entity)

	diagFeature := NewFeatureLocal(
		entity.NextFeatureId(), entity,
		model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	diagFeature.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)
	entity.AddFeature(diagFeature)

	stabilize()
	baseline := countGoroutines()

	// 5 start/stop cycles — goroutine count should be stable
	for i := 0; i < 5; i++ {
		_ = entity.HeartbeatManager().StartHeartbeat()
		time.Sleep(20 * time.Millisecond)
		entity.HeartbeatManager().StopHeartbeat()
	}

	stabilize()
	after := countGoroutines()

	assert.LessOrEqual(s.T(), after, baseline+1,
		"goroutine count should not grow after repeated start/stop (ticker properly stopped)")
}

// =========================================================================
// Issue 6: Event handlers can be drained after teardown
// =========================================================================

func (s *LifecycleIssuesSuite) TestIssue6_EventHandlersDrainedAfterTeardown() {
	writer := &trackingWriter{}
	handler := newSlowEventHandler(300 * time.Millisecond)

	device := NewDeviceLocal(
		"Demo", "Model", "Serial", "Code", "device-1",
		model.DeviceTypeTypeEnergyManagementSystem,
		model.NetworkManagementFeatureSetTypeSmart)

	_ = device.events.Subscribe(handler)
	defer func() { _ = device.events.Unsubscribe(handler) }()

	ski := "remote-1"
	_ = device.SetupRemoteDevice(ski, writer)

	device.RemoveRemoteDeviceConnection(ski)

	// Handler is still running — drain waits for it
	device.events.drain()

	started := handler.started.Load()
	finished := handler.finished.Load()
	assert.Equal(s.T(), started, finished,
		"all event handlers should be finished after drain()")
}

// =========================================================================
// Issue 7: No double response from ApproveOrDenyWrite
// =========================================================================

func (s *LifecycleIssuesSuite) TestIssue7_NoDoubleResponse() {
	writer := &trackingWriter{}

	device, localEntity := createLocalDeviceAndEntity(1)
	_, serverFeature := createLocalFeatures(localEntity, model.FeatureTypeTypeLoadControl, model.FunctionTypeLoadControlLimitListData)

	serverFeature.AddWriteApprovalCallback(func(msg *api.Message) {})

	ski := "remote-1"
	sender := NewSender(writer)
	remoteDevice := createRemoteDevice(device, ski, sender)
	device.AddRemoteDeviceForSki(ski, remoteDevice)

	remoteFeature, _ := createRemoteEntityAndFeature(remoteDevice, 1,
		model.FeatureTypeTypeLoadControl, model.FunctionTypeLoadControlLimitListData)

	serverFeature.SetWriteApprovalTimeout(80 * time.Millisecond)

	msgCounter := model.MsgCounterType(99)
	msg := &api.Message{
		RequestHeader: &model.HeaderType{
			MsgCounter: util.Ptr(msgCounter),
			AddressSource: &model.FeatureAddressType{
				Device:  remoteDevice.Address(),
				Entity:  remoteDevice.Entity([]model.AddressEntityType{1}).Address().Entity,
				Feature: remoteFeature.Address().Feature,
			},
			AddressDestination: serverFeature.Address(),
		},
		DeviceRemote:  remoteDevice,
		EntityRemote:  remoteDevice.Entity([]model.AddressEntityType{1}),
		FeatureRemote: remoteFeature,
	}

	serverFeature.(*FeatureLocal).addPendingApproval(msg)

	// Wait close to timeout, then approve
	time.Sleep(70 * time.Millisecond)
	serverFeature.ApproveOrDenyWrite(msg, model.ErrorType{ErrorNumber: model.ErrorNumberType(0)})

	// Wait for any timer to fire
	time.Sleep(100 * time.Millisecond)

	assert.LessOrEqual(s.T(), writer.callCount(), int64(1),
		"at most 1 response should be sent, not both timer error and approval success")

	device.RemoveRemoteDevice(ski)
}

// =========================================================================
// Issue 9: SetupRemoteDevice cleans up existing device
// =========================================================================

func (s *LifecycleIssuesSuite) TestIssue9_SetupRemoteDeviceCleansUpExisting() {
	writer := &trackingWriter{}

	device := NewDeviceLocal(
		"Demo", "Model", "Serial", "Code", "device-1",
		model.DeviceTypeTypeEnergyManagementSystem,
		model.NetworkManagementFeatureSetTypeSmart)

	ski := "remote-1"

	reader1 := device.SetupRemoteDevice(ski, writer)
	assert.NotNil(s.T(), reader1)

	device1 := device.RemoteDeviceForSki(ski)
	assert.NotNil(s.T(), device1)

	// Second setup for same SKI — should clean up old device first
	reader2 := device.SetupRemoteDevice(ski, writer)
	assert.NotNil(s.T(), reader2)

	device2 := device.RemoteDeviceForSki(ski)
	assert.NotNil(s.T(), device2)

	// New device should be different object from old one
	assert.NotEqual(s.T(), device1, device2,
		"second SetupRemoteDevice should create a new device, not reuse the old one")

	device.RemoveRemoteDeviceConnection(ski)
}

// =========================================================================
// Combined: Full lifecycle with Close() — all clean
// =========================================================================

func (s *LifecycleIssuesSuite) TestCombined_FullLifecycleWithClose() {
	writer := &trackingWriter{}
	handler := newSlowEventHandler(200 * time.Millisecond)

	stabilize()
	goroutinesBefore := countGoroutines()

	device := NewDeviceLocal(
		"Demo", "Model", "Serial", "Code", "device-1",
		model.DeviceTypeTypeEnergyManagementSystem,
		model.NetworkManagementFeatureSetTypeSmart)

	_ = device.events.Subscribe(handler)
	defer func() {
		_ = device.events.Unsubscribe(handler)
		stabilize()
	}()

	for i := uint(1); i <= 3; i++ {
		entity := NewEntityLocal(device, model.EntityTypeTypeCEM,
			[]model.AddressEntityType{model.AddressEntityType(i)}, 4*time.Second)
		device.AddEntity(entity)

		diagFeature := NewFeatureLocal(
			entity.NextFeatureId(), entity,
			model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
		diagFeature.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)
		entity.AddFeature(diagFeature)
	}

	for i := 1; i <= 2; i++ {
		ski := fmt.Sprintf("remote-%d", i)
		_ = device.SetupRemoteDevice(ski, writer)
	}

	stabilize()
	goroutinesRunning := countGoroutines()
	assert.Greater(s.T(), goroutinesRunning, goroutinesBefore,
		"heartbeat goroutines should be running")

	// Single Close() call cleans up everything
	device.Close()

	stabilize()
	goroutinesAfter := countGoroutines()

	assert.LessOrEqual(s.T(), goroutinesAfter, goroutinesBefore+1,
		"all goroutines should be cleaned up after Close()")
	assert.Empty(s.T(), device.RemoteDevices(),
		"all remote devices should be removed after Close()")

	// Only entity[0] should remain
	entities := device.Entities()
	assert.Equal(s.T(), 1, len(entities),
		"only DeviceInformation entity should remain after Close()")
}
