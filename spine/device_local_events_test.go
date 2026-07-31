package spine

import (
	"testing"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

// Test that DeviceLocal has Events() method returning its own events manager
func TestDeviceLocalHasOwnEventsManager(t *testing.T) {
	device := NewDeviceLocal("brand", "model", "serial", "code", "addr",
		model.DeviceTypeTypeGeneric, model.NetworkManagementFeatureSetTypeSmart)

	assert.NotNil(t, device.Events())
}

// Test that each DeviceLocal gets its own events manager (automatic isolation)
func TestDeviceLocalEventsAreIsolated(t *testing.T) {
	deviceA := NewDeviceLocal("brandA", "modelA", "serialA", "codeA", "addrA",
		model.DeviceTypeTypeGeneric, model.NetworkManagementFeatureSetTypeSmart)

	deviceB := NewDeviceLocal("brandB", "modelB", "serialB", "codeB", "addrB",
		model.DeviceTypeTypeGeneric, model.NetworkManagementFeatureSetTypeSmart)

	// Each device should have its own events manager (not same pointer)
	assert.NotSame(t, deviceA.Events(), deviceB.Events())
}

// Test that DeviceLocalInterface includes Events() method
func TestDeviceLocalInterfaceHasEvents(t *testing.T) {
	device := NewDeviceLocal("brand", "model", "serial", "code", "addr",
		model.DeviceTypeTypeGeneric, model.NetworkManagementFeatureSetTypeSmart)

	// Use interface type to verify the method exists on the interface
	var deviceInterface api.DeviceLocalInterface = device
	assert.NotNil(t, deviceInterface.Events())
}
