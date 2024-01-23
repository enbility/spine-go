package api

import "github.com/enbility/spine-go/model"

type Message struct {
	RequestHeader *model.HeaderType
	CmdClassifier model.CmdClassifierType
	Cmd           model.CmdType
	FilterPartial *model.FilterType
	FilterDelete  *model.FilterType
	FeatureRemote FeatureRemoteInterface
	EntityRemote  EntityRemoteInterface
	DeviceRemote  DeviceRemoteInterface
}

type ResultMessage struct {
	MsgCounterReference model.MsgCounterType   // required
	Result              *model.ResultDataType  // required, may not be nil
	FeatureLocal        FeatureLocalInterface  // required, may not be nil
	FeatureRemote       FeatureRemoteInterface // required, may not be nil
	EntityRemote        EntityRemoteInterface  // required, may not be nil
	DeviceRemote        DeviceRemoteInterface  // required, may not be nil
}
