package api

import "github.com/enbility/spine-go/model"

/* Entity */

// This interface defines the functions being common to local and remote entites
// An entity corresponds to a SPINE entity, see SPINE Introduction Chapter 2.2
type EntityInterface interface {
	// Get the entity address
	Address() *model.EntityAddressType
	// Get the entity type
	EntityType() model.EntityTypeType
	// Get the entity description
	Description() *model.DescriptionType
	// Set the entity description
	SetDescription(d *model.DescriptionType)
	// Get the next incremental feature id
	NextFeatureId() uint
}

// This interface defines all the required functions need to implement a local entity
type EntityLocalInterface interface {
	EntityInterface
	// Get the associated DeviceLocalInterface implementation
	Device() DeviceLocalInterface

	// Get the hearbeat manager for this entity
	HeartbeatManager() HeartbeatManagerInterface

	// Add a new feature with a given FeatureLocalInterface implementation
	AddFeature(f FeatureLocalInterface)
	// Get a FeatureLocalInterface implementation for a given feature type and role or create it if it does not exist yet and return it
	GetOrAddFeature(featureType model.FeatureTypeType, role model.RoleType) FeatureLocalInterface
	// Get a FeatureLocalInterface implementation for a given feature type and role
	FeatureOfTypeAndRole(featureType model.FeatureTypeType, role model.RoleType) FeatureLocalInterface
	// Get a FeatureLocalInterface implementation for a given feature address
	FeatureOfAddress(addressFeature *model.AddressFeatureType) FeatureLocalInterface
	// Get all FeatureLocalInterface implementations
	Features() []FeatureLocalInterface

	// Add a new usecase
	AddUseCaseSupport(
		actor model.UseCaseActorType,
		useCaseName model.UseCaseNameType,
		useCaseVersion model.SpecificationVersionType,
		useCaseDocumemtSubRevision string,
		useCaseAvailable bool,
		scenarios []model.UseCaseScenarioSupportType,
	)
	// Check if a use case is already added
	HasUseCaseSupport(
		actor model.UseCaseActorType,
		useCaseName model.UseCaseNameType) bool
	// Remove support for a usecase
	RemoveUseCaseSupport(
		actor model.UseCaseActorType,
		useCaseName model.UseCaseNameType,
	)
	// Set the availability of a usecase. This may only be used for usescases
	// that act as a client within the usecase!
	SetUseCaseAvailability(actor model.UseCaseActorType, useCaseName model.UseCaseNameType, available bool)
	// Remove all usecases
	RemoveAllUseCaseSupports()

	// Remove all subscriptions
	RemoveAllSubscriptions()

	// Remove all bindings
	RemoveAllBindings()

	// Get the SPINE data structure for NodeManagementDetailDiscoveryData messages for this entity
	Information() *model.NodeManagementDetailedDiscoveryEntityInformationType
}

// This interface defines all the required functions need to implement a remote entity
type EntityRemoteInterface interface {
	EntityInterface

	// Get the associated DeviceRemoteInterface implementation
	Device() DeviceRemoteInterface

	// Update the device address (only used for the DeviceInformation entity when receiving the DetailDiscovery reply)
	UpdateDeviceAddress(address model.AddressDeviceType)

	// Add a new feature with a given FeatureLocalInterface implementation
	AddFeature(f FeatureRemoteInterface)

	// Remove all features
	RemoveAllFeatures()

	// Get a FeatureRemoteInterface implementation for a given feature type and role
	FeatureOfTypeAndRole(featureType model.FeatureTypeType, role model.RoleType) FeatureRemoteInterface
	// Get a FeatureRemoteInterface implementation for a given feature address
	FeatureOfAddress(addressFeature *model.AddressFeatureType) FeatureRemoteInterface
	// Get all FeatureRemoteInterface implementations
	Features() []FeatureRemoteInterface
}
