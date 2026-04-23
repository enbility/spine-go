package spine

import (
	"sync"
	"testing"
	"time"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

// testEventHandler is a simple event handler for testing
type testEventHandler struct {
	mu       sync.Mutex
	events   []api.EventPayload
	deviceID string // identifier to track which handler this is
}

func newTestEventHandler(deviceID string) *testEventHandler {
	return &testEventHandler{
		events:   make([]api.EventPayload, 0),
		deviceID: deviceID,
	}
}

func (h *testEventHandler) HandleEvent(payload api.EventPayload) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.events = append(h.events, payload)
}

func (h *testEventHandler) receivedEvents() []api.EventPayload {
	h.mu.Lock()
	defer h.mu.Unlock()
	result := make([]api.EventPayload, len(h.events))
	copy(result, h.events)
	return result
}

// TestEventsIsolationBetweenDevices validates that two DeviceLocal instances
// automatically have separate event managers and don't receive each other's events
func TestEventsIsolationBetweenDevices(t *testing.T) {
	// Create two devices - each automatically gets its own events manager
	deviceA := NewDeviceLocal("brandA", "modelA", "serialA", "codeA", "addrA",
		model.DeviceTypeTypeGeneric, model.NetworkManagementFeatureSetTypeSmart)

	deviceB := NewDeviceLocal("brandB", "modelB", "serialB", "codeB", "addrB",
		model.DeviceTypeTypeGeneric, model.NetworkManagementFeatureSetTypeSmart)

	// Verify they have different events managers (automatic isolation)
	assert.NotSame(t, deviceA.Events(), deviceB.Events())

	// Create handlers for each device
	handlerA := newTestEventHandler("deviceA")
	handlerB := newTestEventHandler("deviceB")

	// Subscribe handlers to their respective device's event managers
	_ = deviceA.Events().Subscribe(handlerA)
	_ = deviceB.Events().Subscribe(handlerB)

	// Publish an event on device A's event manager
	payloadA := api.EventPayload{
		Ski:        "ski-device-a",
		EventType:  api.EventTypeDataChange,
		ChangeType: api.ElementChangeUpdate,
	}
	deviceA.Events().Publish(payloadA)

	// Give async handlers time to process
	time.Sleep(50 * time.Millisecond)

	// Handler A should receive the event
	eventsFromA := handlerA.receivedEvents()
	assert.Len(t, eventsFromA, 1, "Handler A should receive 1 event")
	assert.Equal(t, "ski-device-a", eventsFromA[0].Ski)

	// Handler B should NOT receive the event (isolation!)
	eventsFromB := handlerB.receivedEvents()
	assert.Len(t, eventsFromB, 0, "Handler B should NOT receive events from device A")

	// Now publish on device B
	payloadB := api.EventPayload{
		Ski:        "ski-device-b",
		EventType:  api.EventTypeDeviceChange,
		ChangeType: api.ElementChangeAdd,
	}
	deviceB.Events().Publish(payloadB)

	// Give async handlers time to process
	time.Sleep(50 * time.Millisecond)

	// Handler B should now have 1 event
	eventsFromB = handlerB.receivedEvents()
	assert.Len(t, eventsFromB, 1, "Handler B should receive 1 event")
	assert.Equal(t, "ski-device-b", eventsFromB[0].Ski)

	// Handler A should still have only 1 event (no cross-talk)
	eventsFromA = handlerA.receivedEvents()
	assert.Len(t, eventsFromA, 1, "Handler A should still have only 1 event (no cross-talk)")

	// Clean up subscriptions
	_ = deviceA.Events().Unsubscribe(handlerA)
	_ = deviceB.Events().Unsubscribe(handlerB)
}
