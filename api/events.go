package api

import "github.com/enbility/spine-go/model"

type EventHandlerLevel uint

const (
	EventHandlerLevelCore        EventHandlerLevel = iota // Shall only be used by the core stack
	EventHandlerLevelApplication                          // Shall only be used by applications
)

type ElementChangeType uint16

const (
	ElementChangeAdd ElementChangeType = iota
	ElementChangeUpdate
	ElementChangeRemove
)

type EventType uint16

const (
	EventTypeDeviceChange       EventType = iota // Sent after successful response of NodeManagementDetailedDiscovery
	EventTypeEntityChange                        // Sent after successful response of NodeManagementDetailedDiscovery
	EventTypeSubscriptionChange                  // Sent after successful subscription request from remote
	EventTypeBindingChange                       // Sent after successful binding request from remote
	EventTypeDataChange                          // Sent after remote provided new data items for a function
)

type EventPayload struct {
	Ski           string                // required
	EventType     EventType             // required
	ChangeType    ElementChangeType     // required
	Device        DeviceRemoteInterface // required for DetailedDiscovery Call
	Entity        EntityRemoteInterface // required for DetailedDiscovery Call and Notify
	Feature       FeatureRemoteInterface
	LocalFeature  FeatureLocalInterface    // required for write commands
	Function      model.FunctionType       // required for write commands
	CmdClassifier *model.CmdClassifierType // optional, used together with EventType EventTypeDataChange
	Data          any
}

// EventsManagerInterface defines the interface for managing event subscriptions and publishing.
// This interface allows for dependency injection of the events manager, enabling
// test isolation when multiple DeviceLocal instances exist in the same process.
type EventsManagerInterface interface {
	// Subscribe registers an event handler to receive events at the application level.
	Subscribe(handler EventHandlerInterface) error
	// Unsubscribe removes an event handler from receiving events.
	Unsubscribe(handler EventHandlerInterface) error
	// Publish sends an event to all registered handlers.
	Publish(payload EventPayload)
}
