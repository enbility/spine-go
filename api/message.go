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

type ResponseMessage struct {
	MsgCounterReference model.MsgCounterType   // required
	Data                any                    // required
	FeatureLocal        FeatureLocalInterface  // required
	FeatureRemote       FeatureRemoteInterface // required
	EntityRemote        EntityRemoteInterface  // required
	DeviceRemote        DeviceRemoteInterface  // required
}
