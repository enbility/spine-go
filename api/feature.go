package api

import (
	"time"

	"github.com/enbility/spine-go/model"
)

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
	RequestRemoteData(
		function model.FunctionType,
		selector any,
		elements any,
		destination FeatureRemoteInterface) (*model.MsgCounterType, *model.ErrorType)
	RequestRemoteDataBySenderAddress(
		cmd model.CmdType,
		sender SenderInterface,
		destinationSki string,
		destinationAddress *model.FeatureAddressType,
		maxDelay time.Duration) (*model.MsgCounterType, *model.ErrorType)
	FetchRequestRemoteData(
		msgCounter model.MsgCounterType,
		destination FeatureRemoteInterface) (any, *model.ErrorType)
	Subscribe(remoteAdress *model.FeatureAddressType) (*model.MsgCounterType, *model.ErrorType)
	RemoveSubscription(remoteAddress *model.FeatureAddressType)
	RemoveAllSubscriptions()
	Bind(remoteAdress *model.FeatureAddressType) (*model.MsgCounterType, *model.ErrorType)
	RemoveBinding(remoteAddress *model.FeatureAddressType)
	RemoveAllBindings()
	HandleMessage(message *Message) *model.ErrorType
}

type NodeManagementInterface interface {
	FeatureLocalInterface
}

type FeatureResultInterface interface {
	HandleResult(ResultMessage)
}

type FeatureRemoteInterface interface {
	FeatureInterface
	DataCopy(function model.FunctionType) any
	UpdateData(function model.FunctionType, data any, filterPartial *model.FilterType, filterDelete *model.FilterType)
	Sender() SenderInterface
	Device() DeviceRemoteInterface
	Entity() EntityRemoteInterface
	SetOperations(functions []model.FunctionPropertyType)
	SetMaxResponseDelay(delay *model.MaxResponseDelayType)
	MaxResponseDelayDuration() time.Duration
}
