package api

import (
	"github.com/enbility/spine-go/model"
)

//go:generate mockery

type EventHandlerInterface interface {
	HandleEvent(EventPayload)
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
