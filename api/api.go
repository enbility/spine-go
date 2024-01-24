package api

import (
	"time"

	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/spine-go/model"
)

//go:generate mockery

type EventHandlerInterface interface {
	HandleEvent(EventPayload)
}

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

/* Feature */

type FeatureInterface interface {
	Address() *model.FeatureAddressType
	Type() model.FeatureTypeType
	Role() model.RoleType
	Operations() map[model.FunctionType]OperationsInterface
	Description() *model.DescriptionType
	SetDescription(desc *model.DescriptionType)
	SetDescriptionString(s string)
	String() string
}

type FeatureRemoteInterface interface {
	FeatureInterface
	DataCopy(function model.FunctionType) any
	SetData(function model.FunctionType, data any)
	UpdateData(function model.FunctionType, data any, filterPartial *model.FilterType, filterDelete *model.FilterType)
	Sender() SenderInterface
	Device() DeviceRemoteInterface
	Entity() EntityRemoteInterface
	SetOperations(functions []model.FunctionPropertyType)
	SetMaxResponseDelay(delay *model.MaxResponseDelayType)
	MaxResponseDelayDuration() time.Duration
}

type FeatureLocalInterface interface {
	FeatureInterface
	Device() DeviceLocalInterface
	Entity() EntityLocalInterface
	DataCopy(function model.FunctionType) any
	SetData(function model.FunctionType, data any)
	AddResultHandler(handler FeatureResultInterface)
	AddResultCallback(msgCounterReference model.MsgCounterType, function func(msg ResultMessage))
	Information() *model.NodeManagementDetailedDiscoveryFeatureInformationType
	AddFunctionType(function model.FunctionType, read, write bool)
	RequestData(
		function model.FunctionType,
		selector any,
		elements any,
		destination FeatureRemoteInterface) (*model.MsgCounterType, *model.ErrorType)
	RequestDataBySenderAddress(
		cmd model.CmdType,
		sender SenderInterface,
		destinationSki string,
		destinationAddress *model.FeatureAddressType,
		maxDelay time.Duration) (*model.MsgCounterType, *model.ErrorType)
	FetchRequestData(
		msgCounter model.MsgCounterType,
		destination FeatureRemoteInterface) (any, *model.ErrorType)
	Subscribe(remoteAdress *model.FeatureAddressType) (*model.MsgCounterType, *model.ErrorType)
	// SubscribeAndWait(remoteDevice DeviceRemote, remoteAdress *model.FeatureAddressType) *ErrorType // Subscribes the local feature to the given destination feature; the go routine will block until the response is processed
	RemoveSubscription(remoteAddress *model.FeatureAddressType)
	RemoveAllSubscriptions()
	Bind(remoteAdress *model.FeatureAddressType) (*model.MsgCounterType, *model.ErrorType)
	// BindAndWait(remoteDevice DeviceRemote, remoteAddress *model.FeatureAddressType) *ErrorType
	RemoveBinding(remoteAddress *model.FeatureAddressType)
	RemoveAllBindings()
	NotifyData(
		function model.FunctionType,
		deleteSelector, partialSelector any,
		partialWithoutSelector bool,
		deleteElements any,
		destination FeatureRemoteInterface) (*model.MsgCounterType, *model.ErrorType)
	WriteData(
		function model.FunctionType,
		deleteSelector, partialSelector any,
		deleteElements any,
		destination FeatureRemoteInterface) (*model.MsgCounterType, *model.ErrorType)
	HandleMessage(message *Message) *model.ErrorType
}

type NodeManagementInterface interface {
	FeatureLocalInterface
}

type FeatureResultInterface interface {
	HandleResult(ResultMessage)
}

/* Functions */

type FunctionDataCmdInterface interface {
	FunctionDataInterface
	ReadCmdType(partialSelector any, elements any) model.CmdType
	ReplyCmdType(partial bool) model.CmdType
	NotifyCmdType(deleteSelector, partialSelector any, partialWithoutSelector bool, deleteElements any) model.CmdType
	WriteCmdType(deleteSelector, partialSelector any, deleteElements any) model.CmdType
}

type FunctionDataInterface interface {
	Function() model.FunctionType
	DataCopyAny() any
	UpdateDataAny(data any, filterPartial *model.FilterType, filterDelete *model.FilterType)
}

/* Sender */

type ComControlInterface interface {
	// This must be connected to the correct remote device !!
	SendSpineMessage(datagram model.DatagramType) error
}

type SenderInterface interface {
	// Sends a read cmd to request some data
	Request(cmdClassifier model.CmdClassifierType, senderAddress, destinationAddress *model.FeatureAddressType, ackRequest bool, cmd []model.CmdType) (*model.MsgCounterType, error)
	// Sends a result cmd with no error to indicate that a message was processed successfully
	ResultSuccess(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType) error
	// Sends a result cmd with error information to indicate that a message processing failed
	ResultError(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType, err *model.ErrorType) error
	// Sends a reply cmd to response to a read cmd
	Reply(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType, cmd model.CmdType) error
	// Sends a call cmd with a subscription request
	Subscribe(senderAddress, destinationAddress *model.FeatureAddressType, serverFeatureType model.FeatureTypeType) (*model.MsgCounterType, error)
	// Sends a call cmd with a subscription delete request
	Unsubscribe(senderAddress, destinationAddress *model.FeatureAddressType) (*model.MsgCounterType, error)
	// Sends a call cmd with a binding request
	Bind(senderAddress, destinationAddress *model.FeatureAddressType, serverFeatureType model.FeatureTypeType) (*model.MsgCounterType, error)
	// Sends a call cmd with a binding delte request
	Unbind(senderAddress, destinationAddress *model.FeatureAddressType) (*model.MsgCounterType, error)
	// Sends a notify cmd to indicate that a subscribed feature changed
	Notify(senderAddress, destinationAddress *model.FeatureAddressType, cmd model.CmdType) (*model.MsgCounterType, error)
	// Sends a write cmd, setting properties of remote features
	Write(senderAddress, destinationAddress *model.FeatureAddressType, cmd model.CmdType) (*model.MsgCounterType, error)
	// return the datagram for a given msgCounter (only availbe for Notify messasges!), error if not found
	DatagramForMsgCounter(msgCounter model.MsgCounterType) (model.DatagramType, error)
}

/* PendingRequests */

type PendingRequestsInterface interface {
	Add(ski string, counter model.MsgCounterType, maxDelay time.Duration)
	SetData(ski string, counter model.MsgCounterType, data any) *model.ErrorType
	SetResult(ski string, counter model.MsgCounterType, errorResult *model.ErrorType) *model.ErrorType
	GetData(ski string, counter model.MsgCounterType) (any, *model.ErrorType)
	Remove(ski string, counter model.MsgCounterType) *model.ErrorType
}

/* Bindings */

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
	UpdateHeartbeatOnSubscriptions()
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
