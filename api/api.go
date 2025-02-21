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
	// Add a binding between a client and server feature where one of each is local and the other one is remote
	AddBinding(remoteDevice DeviceRemoteInterface, data model.BindingManagementRequestCallType) error
	// Remove a binding between a client and server feature where one of each is local and the other one is remote
	RemoveBinding(remoteDevice DeviceRemoteInterface, data model.BindingManagementDeleteCallType) error
	// Remove all stored bindings for a given remote device
	RemoveBindingsForRemoteDevice(remoteDevice DeviceRemoteInterface)
	// Remove all stored bindings for a given remote device entity
	RemoveBindingsForRemoteEntity(remoteEntity EntityRemoteInterface)
	// Remove all stored bindings for a given local device entity
	RemoveBindingsForLocalEntity(localEntity EntityLocalInterface)
	// Checks if a binding between the client and server feature exists
	HasBinding(clientAddress, serverAddress *model.FeatureAddressType) bool
	// Return all stored bindings for a given remote device
	BindingsForRemoteDevice(remoteDevice DeviceRemoteInterface) []model.BindingManagementEntryDataType
	// Return all stored bindings for a given feature address
	BindingsForFeatureAddress(localAddress model.FeatureAddressType) []model.BindingManagementEntryDataType
}

/* Subscription Manager */

type SubscriptionManagerInterface interface {
	// Add a subscription between a client and server feature where one of each is local and the other one is remote
	AddSubscription(remoteDevice DeviceRemoteInterface, data model.SubscriptionManagementRequestCallType) error
	// Remove a subscription between a client and server feature where one of each is local and the other one is remote
	RemoveSubscription(remoteDevice DeviceRemoteInterface, data model.SubscriptionManagementDeleteCallType) error
	// Remove all stored subscription for a given remote device
	RemoveSubscriptionsForRemoteDevice(remoteDevice DeviceRemoteInterface)
	// Remove all stored subscription for a given remote device entity
	RemoveSubscriptionsForRemoteEntity(remoteEntity EntityRemoteInterface)
	// Remove all stored subscription for a given local device entity
	RemoveSubscriptionsForLocalEntity(localEntity EntityLocalInterface)
	// Checks if a subscription between the client and server feature exists
	HasSubscription(clientAddress, serverAddress *model.FeatureAddressType) bool
	// Return all stored subscriptions for a given remote device
	SubscriptionsForRemoteDevice(remoteDevice DeviceRemoteInterface) []model.SubscriptionManagementEntryDataType
	// Return all stored subscriptions for a given feature address
	SubscriptionsForFeatureAddress(localAddress model.FeatureAddressType) []model.SubscriptionManagementEntryDataType
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
	WritePartial() bool
	Read() bool
	ReadPartial() bool
	String() string
	Information() *model.PossibleOperationsType
}
