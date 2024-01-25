package api

import "github.com/enbility/spine-go/model"

/* Entity */

type EntityInterface interface {
	EntityType() model.EntityTypeType
	Address() *model.EntityAddressType
	Description() *model.DescriptionType
	SetDescription(d *model.DescriptionType)
	NextFeatureId() uint
}

type EntityLocalInterface interface {
	EntityInterface
	Device() DeviceLocalInterface
	AddFeature(f FeatureLocalInterface)
	GetOrAddFeature(featureType model.FeatureTypeType, role model.RoleType) FeatureLocalInterface
	FeatureOfTypeAndRole(featureType model.FeatureTypeType, role model.RoleType) FeatureLocalInterface
	Features() []FeatureLocalInterface
	Feature(addressFeature *model.AddressFeatureType) FeatureLocalInterface
	Information() *model.NodeManagementDetailedDiscoveryEntityInformationType
	AddUseCaseSupport(
		actor model.UseCaseActorType,
		useCaseName model.UseCaseNameType,
		useCaseVersion model.SpecificationVersionType,
		useCaseDocumemtSubRevision string,
		useCaseAvailable bool,
		scenarios []model.UseCaseScenarioSupportType,
	)
	HasUseCaseSupport(
		actor model.UseCaseActorType,
		useCaseName model.UseCaseNameType) bool
	RemoveUseCaseSupport(
		actor model.UseCaseActorType,
		useCaseName model.UseCaseNameType,
	)
	RemoveAllUseCaseSupports()
	RemoveAllSubscriptions()
	RemoveAllBindings()
}

type EntityRemoteInterface interface {
	EntityInterface
	Device() DeviceRemoteInterface
	AddFeature(f FeatureRemoteInterface)
	Features() []FeatureRemoteInterface
	Feature(addressFeature *model.AddressFeatureType) FeatureRemoteInterface
	UpdateDeviceAddress(address model.AddressDeviceType)
	RemoveAllFeatures()
}
