package spine

import (
	"sync"

	"github.com/enbility/spine-go/api"
)

var Events events

type eventHandlerItem struct {
	Level   api.EventHandlerLevel
	Handler api.EventHandler
}

type events struct {
	mu       sync.Mutex
	muHandle sync.Mutex

	handlers []eventHandlerItem // event handling outside of the core stack
}

// will be used in EEBUS core directly to access the level EventHandlerLevelCore
func (r *events) subscribe(level api.EventHandlerLevel, handler api.EventHandler) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	exists := false
	for _, item := range r.handlers {
		if item.Level == level && item.Handler == handler {
			exists = true
			break
		}
	}

	if !exists {
		newHandlerItem := eventHandlerItem{
			Level:   level,
			Handler: handler,
		}
		r.handlers = append(r.handlers, newHandlerItem)
	}

	return nil
}

// Subscribe to message events and handle them in
// the Eventhandler interface implementation
//
// returns an error if EventHandlerLevelCore is used as
// that is only allowed for internal use
func (r *events) Subscribe(handler api.EventHandler) error {
	return r.subscribe(api.EventHandlerLevelApplication, handler)
}

// will be used in EEBUS core directly to access the level EventHandlerLevelCore
func (r *events) unsubscribe(level api.EventHandlerLevel, handler api.EventHandler) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var newHandlers []eventHandlerItem
	for _, item := range r.handlers {
		if item.Level != level && item.Handler != handler {
			newHandlers = append(newHandlers, item)
		}
	}

	r.handlers = newHandlers

	return nil
}

// Unsubscribe from getting events
func (r *events) Unsubscribe(handler api.EventHandler) error {
	return r.unsubscribe(api.EventHandlerLevelApplication, handler)
}

// Publish an event to all subscribers
func (r *events) Publish(payload api.EventPayload) {
	r.mu.Lock()
	var handler []eventHandlerItem
	copy(r.handlers, handler)
	r.mu.Unlock()

	// Use different locks, so unpublish is possible in the event handlers
	r.muHandle.Lock()
	// process subscribers by level
	handlerLevels := []api.EventHandlerLevel{
		api.EventHandlerLevelCore,
		api.EventHandlerLevelApplication,
	}

	for _, level := range handlerLevels {
		for _, item := range r.handlers {
			if item.Level != level {
				continue
			}

			// do not run this asynchronously, to make sure all required
			// and expected actions are taken
			item.Handler.HandleEvent(payload)
		}
	}
	r.muHandle.Unlock()
}
