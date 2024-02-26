package api

import (
	"time"

	"github.com/enbility/spine-go/model"
)

/* Feature */

// This interface defines the functions being common to local and remote features
// A feature corresponds to a SPINE feature, see SPINE Introduction Chapter 2.2
type FeatureInterface interface {
	// Get the feature address
	Address() *model.FeatureAddressType
	// Get the feature type
	Type() model.FeatureTypeType
	// Get the feature role
	Role() model.RoleType
	// Get the feature operations
	Operations() map[model.FunctionType]OperationsInterface
	// Get the feature description
	Description() *model.DescriptionType
	// Set the feature description with a given type
	SetDescription(desc *model.DescriptionType)
	// Set the feature description with a given string
	SetDescriptionString(s string)
	// Return a descriptive feature summary as a string
	String() string
}

// This interface defines all the required functions need to implement a local feature
type FeatureLocalInterface interface {
	FeatureInterface
	// Get the associated DeviceLocalInterface implementation
	Device() DeviceLocalInterface
	// Get the associated EntityLocalInterface implementation
	Entity() EntityLocalInterface

	// Add a function type with
	AddFunctionType(function model.FunctionType, read, write bool)
	// Add a FeatureResultInterface implementation to be able to get incoming result messages for this feature
	AddResultHandler(handler FeatureResultInterface)
	// Add a callback function to be invoked when SPINE message comes in with a given msgCounterReference value
	AddResultCallback(msgCounterReference model.MsgCounterType, function func(msg ResultMessage))

	// return all functions
	Functions() []model.FunctionType

	// Get a copy of the features data for a given function type
	DataCopy(function model.FunctionType) any
	// Set the features data for a given function type
	SetData(function model.FunctionType, data any)

	// Trigger a read request message for a given FeatureRemoteInterface implementation
	RequestRemoteData(
		function model.FunctionType,
		selector any,
		elements any,
		destination FeatureRemoteInterface) (*model.MsgCounterType, *model.ErrorType)
	// Trigger a read request message for a remote ski and feature address
	RequestRemoteDataBySenderAddress(
		cmd model.CmdType,
		sender SenderInterface,
		destinationSki string,
		destinationAddress *model.FeatureAddressType,
		maxDelay time.Duration) (*model.MsgCounterType, *model.ErrorType)
	// Trigger a blocking read request message for a given FeatureRemoteInterface implementation
	FetchRequestRemoteData(
		msgCounter model.MsgCounterType,
		destination FeatureRemoteInterface) (any, *model.ErrorType)

	// Check if there already is a subscription to a given feature remote address
	HasSubscriptionToRemote(remoteAddress *model.FeatureAddressType) bool
	// Trigger a subscription request to a given feature remote address
	SubscribeToRemote(remoteAddress *model.FeatureAddressType) (*model.MsgCounterType, *model.ErrorType)
	// Trigger a subscription removal request for a given feature remote address
	RemoveRemoteSubscription(remoteAddress *model.FeatureAddressType) (*model.MsgCounterType, *model.ErrorType)
	// Trigger subscription removal requests for all subscriptions of this feature
	RemoveAllRemoteSubscriptions()

	// Check if there already is a binding to a given feature remote address
	HasBindingToRemote(remoteAddress *model.FeatureAddressType) bool
	// Trigger a binding request to a given feature remote address
	BindToRemote(remoteAddress *model.FeatureAddressType) (*model.MsgCounterType, *model.ErrorType)
	// Trigger a binding removal request for a given feature remote address
	RemoveRemoteBinding(remoteAddress *model.FeatureAddressType) (*model.MsgCounterType, *model.ErrorType)
	// Trigger binding removal requests for all subscriptions of this feature
	RemoveAllRemoteBindings()

	// Handle an incoming SPINE message for this feature
	HandleMessage(message *Message) *model.ErrorType

	// Get the SPINE data structure for NodeManagementDetailDiscoveryData messages for this feature
	Information() *model.NodeManagementDetailedDiscoveryFeatureInformationType
}

// Interface for local NodeManagement feature
type NodeManagementInterface interface {
	FeatureLocalInterface
}

// Interface for working with SPINE result messages
type FeatureResultInterface interface {
	// Handle a incoming SPINE result message
	HandleResult(ResultMessage)
}

// This interface defines all the required functions need to implement a remote feature
type FeatureRemoteInterface interface {
	FeatureInterface

	// Get the associated DeviceRemoteInterface implementation
	Device() DeviceRemoteInterface
	// Get the associated EntityRemoteInterface implementation
	Entity() EntityRemoteInterface

	// Get a copy of the features data for a given function type
	DataCopy(function model.FunctionType) any
	// Set the features data for a given function type
	UpdateData(function model.FunctionType, data any, filterPartial *model.FilterType, filterDelete *model.FilterType) *model.ErrorType

	// Set the supported operations of the feature for a set of functions
	SetOperations(functions []model.FunctionPropertyType)

	// Define the maximum response duration
	SetMaxResponseDelay(delay *model.MaxResponseDelayType)
	// Get the maximum allowed response duration
	MaxResponseDelayDuration() time.Duration
}
