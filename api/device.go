package api

import (
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/spine-go/model"
)

/* Device */

// This interface defines the functions being common to local and remote devices
// A device corresponds to a SPINE device, see SPINE Introduction Chapter 2.2
type DeviceInterface interface {
	// Get the device address
	Address() *model.AddressDeviceType
	// Get the device type
	DeviceType() *model.DeviceTypeType
	// Get the device feature set
	FeatureSet() *model.NetworkManagementFeatureSetType
	// Get the device destination data
	DestinationData() model.NodeManagementDestinationDataType
}

// This interface defines all the required functions need to implement a local device
type DeviceLocalInterface interface {
	DeviceInterface
	// Setup a new remote device with a given SKI and triggers SPINE requesting device details
	SetupRemoteDevice(ski string, writeI shipapi.ShipConnectionDataWriterInterface) shipapi.ShipConnectionDataReaderInterface
	// Add a DeviceRemoteInterface implementation (used in SetupRemoteDevice and in tests)
	AddRemoteDeviceForSki(ski string, rDevice DeviceRemoteInterface)

	// Request NodeManagmentDetailedDiscovery Data from a remote device
	RequestRemoteDetailedDiscoveryData(rDevice DeviceRemoteInterface) (*model.MsgCounterType, *model.ErrorType)

	// Remove a remote device and its connection
	RemoveRemoteDeviceConnection(ski string)
	// Remove a remote device (used in RemoveRemoteDeviceConnection and in tests)
	RemoveRemoteDevice(ski string)

	// Get a list of connected remote devices DeviceRemoteInterface implementations
	RemoteDevices() []DeviceRemoteInterface
	// Get a connected remote device DeviceRemoteInterface implementation for a given device address
	RemoteDeviceForAddress(address model.AddressDeviceType) DeviceRemoteInterface
	// Get a connected remote device DeviceRemoteInterface implementation for a SKI
	RemoteDeviceForSki(ski string) DeviceRemoteInterface

	// Add a new entity to the device
	// It should trigger a notification of all remote devices about this new entity
	AddEntity(entity EntityLocalInterface)
	// Remove a entity from the device
	// It should trigger a notification of all remote devices about this removed entity
	RemoveEntity(entity EntityLocalInterface)
	// Get a list of all entities EntityLocalInterface implementations
	Entities() []EntityLocalInterface
	// Get an entity EntityLocalInterface implementation for a given entity address
	Entity(id []model.AddressEntityType) EntityLocalInterface
	// Get the first entity EntityLocalInterface implementation for a given entity type
	EntityForType(entityType model.EntityTypeType) EntityLocalInterface

	// Get a FeatureLocalInterface implementation for a given feature address
	FeatureByAddress(address *model.FeatureAddressType) FeatureLocalInterface

	// Process incoming SPINE datagram
	ProcessCmd(datagram model.DatagramType, remoteDevice DeviceRemoteInterface) error

	// Get the node management
	NodeManagement() NodeManagementInterface

	// Get the bindings manager
	BindingManager() BindingManagerInterface

	// Get the subscription manager
	SubscriptionManager() SubscriptionManagerInterface
	// Send a notify message to remote device subscribing to a specific feature
	NotifySubscribers(featureAddress *model.FeatureAddressType, cmd model.CmdType)

	// Get the hearbeat manager
	HeartbeatManager() HeartbeatManagerInterface

	// Get the SPINE data structure for NodeManagementDetailDiscoveryData messages for this device
	Information() *model.NodeManagementDetailedDiscoveryDeviceInformationType
}

// This interface defines all the required functions need to implement a remote device
type DeviceRemoteInterface interface {
	DeviceInterface
	// Get the SKI of a remote device
	Ski() string

	// Add a new entity EntityRemoteInterface implementation
	AddEntity(entity EntityRemoteInterface)

	// Remove an entity for a given address and return the entity that was removed
	RemoveEntityByAddress(addr []model.AddressEntityType) EntityRemoteInterface

	// Get an entity EntityRemoteInterface implementation for a given address
	Entity(id []model.AddressEntityType) EntityRemoteInterface
	// Get all entities EntityRemoteInterface implementations
	Entities() []EntityRemoteInterface

	// Get a feature FeatureRemoteInterface implementation for a given address
	FeatureByAddress(address *model.FeatureAddressType) FeatureRemoteInterface
	// Get a feature FeatureRemoteInterface implementation from a given entity EntityRemoteInterface implementation by the feature type and feature role
	FeatureByEntityTypeAndRole(entity EntityRemoteInterface, featureType model.FeatureTypeType, role model.RoleType) FeatureRemoteInterface

	// Process incoming data payload
	HandleSpineMesssage(message []byte) (*model.MsgCounterType, error)

	// Get the SenderInterface implementation
	Sender() SenderInterface

	// Get the devices usecase data
	UseCases() []model.UseCaseInformationDataType
	// Verify if the device supports a usecase depending on the usecase actor, name, scenarios and requires server features
	VerifyUseCaseScenariosAndFeaturesSupport(
		usecaseActor model.UseCaseActorType,
		usecaseName model.UseCaseNameType,
		scenarios []model.UseCaseScenarioSupportType,
		serverFeatures []model.FeatureTypeType,
	) bool

	// Update the devices address, type and featureset based on NetworkManagementDeviceDescriptionData
	UpdateDevice(description *model.NetworkManagementDeviceDescriptionDataType)

	// Add entities and their features using provided NodeManagementDetailedDiscoveryData
	AddEntityAndFeatures(initialData bool, data *model.NodeManagementDetailedDiscoveryDataType) ([]EntityRemoteInterface, error)

	// Helper method for checking incoming NodeManagementDetailedDiscoveryEntityInformation data
	CheckEntityInformation(initialData bool, entity model.NodeManagementDetailedDiscoveryEntityInformationType) error
}
