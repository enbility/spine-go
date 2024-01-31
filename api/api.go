package api

import (
	"time"

	"github.com/enbility/spine-go/model"
)

//go:generate mockery

type EventHandlerInterface interface {
	HandleEvent(EventPayload)
}

/* PendingRequests */

type PendingRequestsInterface interface {
	Add(ski string, counter model.MsgCounterType, maxDelay time.Duration)
	SetData(ski string, counter model.MsgCounterType, data any) *model.ErrorType
	SetResult(ski string, counter model.MsgCounterType, errorResult *model.ErrorType) *model.ErrorType
	GetData(ski string, counter model.MsgCounterType) (any, *model.ErrorType)
	Remove(ski string, counter model.MsgCounterType) *model.ErrorType
}

/* Binding Manager */

// implemented by BindingManagerImpl
type BindingManagerInterface interface {
	AddBinding(remoteDevice DeviceRemoteInterface, data model.BindingManagementRequestCallType) error
	RemoveBinding(data model.BindingManagementDeleteCallType, remoteDevice DeviceRemoteInterface) error
	RemoveBindingsForDevice(remoteDevice DeviceRemoteInterface)
	RemoveBindingsForEntity(remoteEntity EntityRemoteInterface)
	Bindings(remoteDevice DeviceRemoteInterface) []*BindingEntry
	BindingsOnFeature(featureAddress model.FeatureAddressType) []*BindingEntry
	HasLocalFeatureRemoteBinding(localAddress, remoteAddress *model.FeatureAddressType) bool
}

/* Subscription Manager */

type SubscriptionManagerInterface interface {
	AddSubscription(remoteDevice DeviceRemoteInterface, data model.SubscriptionManagementRequestCallType) error
	RemoveSubscription(data model.SubscriptionManagementDeleteCallType, remoteDevice DeviceRemoteInterface) error
	RemoveSubscriptionsForDevice(remoteDevice DeviceRemoteInterface)
	RemoveSubscriptionsForEntity(remoteEntity EntityRemoteInterface)
	Subscriptions(remoteDevice DeviceRemoteInterface) []*SubscriptionEntry
	SubscriptionsOnFeature(featureAddress model.FeatureAddressType) []*SubscriptionEntry
}

/* Heartbeats */

type HeartbeatManagerInterface interface {
	IsHeartbeatRunning() bool
	SetLocalFeature(entity EntityLocalInterface, feature FeatureLocalInterface)
	StartHeartbeat() error
	StopHeartbeat()
}

type OperationsInterface interface {
	Write() bool
	Read() bool
	String() string
	Information() *model.PossibleOperationsType
}
