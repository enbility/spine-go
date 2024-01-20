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
	Ski           string            // required
	EventType     EventType         // required
	ChangeType    ElementChangeType // required
	Device        DeviceRemote      // required for DetailedDiscovery Call
	Entity        EntityRemote      // required for DetailedDiscovery Call and Notify
	Feature       FeatureRemote
	LocalFeature  FeatureLocal             // required for write commands
	Function      model.FunctionType       // required for write commands
	CmdClassifier *model.CmdClassifierType // optional, used together with EventType EventTypeDataChange
	Data          any
}
