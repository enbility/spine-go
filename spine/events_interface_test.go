package spine

import (
	"testing"

	"github.com/enbility/spine-go/api"
	"github.com/stretchr/testify/assert"
)

// Test that events struct implements EventsManagerInterface
func TestEventsImplementsInterface(t *testing.T) {
	// Compile-time check: *events must satisfy EventsManagerInterface
	var _ api.EventsManagerInterface = &events{}
	var _ api.EventsManagerInterface = newEvents()

	// Runtime check: newEvents returns a valid implementation
	em := newEvents()
	assert.NotNil(t, em)
}

// Test that newEvents() returns a type that satisfies the interface
func TestNewEventsReturnsInterface(t *testing.T) {
	em := newEvents()
	var _ api.EventsManagerInterface = em
	assert.NotNil(t, em)
}
