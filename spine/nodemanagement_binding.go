package spine

import (
	"fmt"

	"github.com/ahmetb/go-linq/v3"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

func NewNodeManagementBindingRequestCallType(clientAddress *model.FeatureAddressType, serverAddress *model.FeatureAddressType, featureType model.FeatureTypeType) *model.NodeManagementBindingRequestCallType {
	return &model.NodeManagementBindingRequestCallType{
		BindingRequest: &model.BindingManagementRequestCallType{
			ClientAddress:     clientAddress,
			ServerAddress:     serverAddress,
			ServerFeatureType: &featureType,
		},
	}
}

func NewNodeManagementBindingDeleteCallType(clientAddress *model.FeatureAddressType, serverAddress *model.FeatureAddressType) *model.NodeManagementBindingDeleteCallType {
	return &model.NodeManagementBindingDeleteCallType{
		BindingDelete: &model.BindingManagementDeleteCallType{
			ClientAddress: clientAddress,
			ServerAddress: serverAddress,
		},
	}
}

// route bindings request calls to the appropriate feature implementation and add the bindings to the current list
func (r *NodeManagementImpl) processReadBindingData(message *api.Message) error {

	var remoteDeviceBindings []model.BindingManagementEntryDataType
	remoteDeviceBindingEntries := r.Device().BindingManager().Bindings(message.FeatureRemote.Device())
	linq.From(remoteDeviceBindingEntries).SelectT(func(s *api.BindingEntry) model.BindingManagementEntryDataType {
		return model.BindingManagementEntryDataType{
			BindingId:     util.Ptr(model.BindingIdType(s.Id)),
			ServerAddress: s.ServerFeature.Address(),
			ClientAddress: s.ClientFeature.Address(),
		}
	}).ToSlice(&remoteDeviceBindings)

	cmd := model.CmdType{
		NodeManagementBindingData: &model.NodeManagementBindingDataType{
			BindingEntry: remoteDeviceBindings,
		},
	}

	return message.FeatureRemote.Sender().Reply(message.RequestHeader, r.Address(), cmd)
}

func (r *NodeManagementImpl) handleMsgBindingData(message *api.Message) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeCall:
		return r.processReadBindingData(message)

	default:
		return fmt.Errorf("nodemanagement.handleBindingDeleteCall: NodeManagementBindingRequestCall CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}

func (r *NodeManagementImpl) handleMsgBindingRequestCall(message *api.Message, data *model.NodeManagementBindingRequestCallType) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeCall:
		return r.Device().BindingManager().AddBinding(message.FeatureRemote.Device(), *data.BindingRequest)

	default:
		return fmt.Errorf("nodemanagement.handleBindingRequestCall: NodeManagementBindingRequestCall CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}

func (r *NodeManagementImpl) handleMsgBindingDeleteCall(message *api.Message, data *model.NodeManagementBindingDeleteCallType) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeCall:
		return r.Device().BindingManager().RemoveBinding(*data.BindingDelete, message.FeatureRemote.Device())

	default:
		return fmt.Errorf("nodemanagement.handleBindingDeleteCall: NodeManagementBindingRequestCall CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}
