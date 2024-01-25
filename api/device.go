package api

import (
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/spine-go/model"
)

/* Device */

type DeviceInterface interface {
	Address() *model.AddressDeviceType
	DeviceType() *model.DeviceTypeType
	FeatureSet() *model.NetworkManagementFeatureSetType
	DestinationData() model.NodeManagementDestinationDataType
}

type DeviceLocalInterface interface {
	DeviceInterface
	RemoveRemoteDeviceConnection(ski string)
	AddRemoteDeviceForSki(ski string, rDevice DeviceRemoteInterface)
	SetupRemoteDevice(ski string, writeI shipapi.ShipConnectionDataWriterInterface) shipapi.ShipConnectionDataReaderInterface
	RemoveRemoteDevice(ski string)
	RemoteDevices() []DeviceRemoteInterface
	RemoteDeviceForAddress(address model.AddressDeviceType) DeviceRemoteInterface
	RemoteDeviceForSki(ski string) DeviceRemoteInterface
	ProcessCmd(datagram model.DatagramType, remoteDevice DeviceRemoteInterface) error
	NodeManagement() NodeManagementInterface
	SubscriptionManager() SubscriptionManagerInterface
	BindingManager() BindingManagerInterface
	HeartbeatManager() HeartbeatManagerInterface
	AddEntity(entity EntityLocalInterface)
	RemoveEntity(entity EntityLocalInterface)
	Entities() []EntityLocalInterface
	Entity(id []model.AddressEntityType) EntityLocalInterface
	EntityForType(entityType model.EntityTypeType) EntityLocalInterface
	FeatureByAddress(address *model.FeatureAddressType) FeatureLocalInterface
	NotifySubscribers(featureAddress *model.FeatureAddressType, cmd model.CmdType)
	Information() *model.NodeManagementDetailedDiscoveryDeviceInformationType
}

type DeviceRemoteInterface interface {
	DeviceInterface
	Ski() string
	SetAddress(address *model.AddressDeviceType)
	HandleSpineMesssage(message []byte) (*model.MsgCounterType, error)
	Sender() SenderInterface
	Entity(id []model.AddressEntityType) EntityRemoteInterface
	Entities() []EntityRemoteInterface
	FeatureByAddress(address *model.FeatureAddressType) FeatureRemoteInterface
	RemoveByAddress(addr []model.AddressEntityType) EntityRemoteInterface
	FeatureByEntityTypeAndRole(entity EntityRemoteInterface, featureType model.FeatureTypeType, role model.RoleType) FeatureRemoteInterface
	UpdateDevice(description *model.NetworkManagementDeviceDescriptionDataType)
	AddEntityAndFeatures(initialData bool, data *model.NodeManagementDetailedDiscoveryDataType) ([]EntityRemoteInterface, error)
	AddEntity(entity EntityRemoteInterface) EntityRemoteInterface
	UseCases() []model.UseCaseInformationDataType
	VerifyUseCaseScenariosAndFeaturesSupport(
		usecaseActor model.UseCaseActorType,
		usecaseName model.UseCaseNameType,
		scenarios []model.UseCaseScenarioSupportType,
		serverFeatures []model.FeatureTypeType,
	) bool
	CheckEntityInformation(initialData bool, entity model.NodeManagementDetailedDiscoveryEntityInformationType) error
}
