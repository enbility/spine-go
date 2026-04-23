package spine

import (
	"fmt"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
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
func (r *NodeManagement) processReadBindingData(message *api.Message) error {
	bindingMgr := r.Device().BindingManager()
	remoteDeviceBindingEntries := bindingMgr.BindingsForRemoteDevice(message.FeatureRemote.Device())

	cmd := model.CmdType{
		NodeManagementBindingData: &model.NodeManagementBindingDataType{
			BindingEntry: remoteDeviceBindingEntries,
		},
	}

	return message.DeviceRemote.Sender().Reply(message.RequestHeader, r.Address(), cmd)
}

func (r *NodeManagement) handleMsgBindingData(message *api.Message) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeRead:
		return r.processReadBindingData(message)

	default:
		return fmt.Errorf("nodemanagement.handleBindingDeleteCall: NodeManagementBindingRequestCall CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}

func (r *NodeManagement) handleMsgBindingRequestCall(message *api.Message, data *model.NodeManagementBindingRequestCallType) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeCall:
		bindingMgr := r.Device().BindingManager()

		createData := r.createBindingAddMissingDeviceAddresses(message, data.BindingRequest)

		return bindingMgr.AddBinding(message.FeatureRemote.Device(), *createData)

	default:
		return fmt.Errorf("nodemanagement.handleBindingRequestCall: NodeManagementBindingRequestCall CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}

func (r *NodeManagement) handleMsgBindingDeleteCall(message *api.Message, data *model.NodeManagementBindingDeleteCallType) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeCall:
		bindingMgr := r.Device().BindingManager()

		deleteData := r.deleteBindingAddMissingDeviceAddresses(message, data.BindingDelete)

		return bindingMgr.RemoveBinding(message.FeatureRemote.Device(), *deleteData)

	default:
		return fmt.Errorf("nodemanagement.handleBindingDeleteCall: NodeManagementBindingRequestCall CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}

// adds potentially missing device addresses to the binding data according to SPINE protocol spec 7.3.2
func (r *NodeManagement) createBindingAddMissingDeviceAddresses(message *api.Message, data *model.BindingManagementRequestCallType) *model.BindingManagementRequestCallType {
	// any device address missing rule according to the spec:
	// If absent, the receiver has to identify the device via some other method.

	// subscriptions can only be requested by clients, so the server must be the recipient
	if data.ClientAddress.Device == nil {
		data.ClientAddress.Device = message.DeviceRemote.Address()
	}
	if data.ServerAddress.Device == nil {
		data.ServerAddress.Device = r.Device().Address()
	}

	return data
}

// adds potentially missing device addresses to the binding data according to SPINE protocol spec 7.3.4
func (r *NodeManagement) deleteBindingAddMissingDeviceAddresses(message *api.Message, data *model.BindingManagementDeleteCallType) *model.BindingManagementDeleteCallType {
	if data.ClientAddress.Device == nil && data.ServerAddress.Device == nil {
		// if both are missing, then client has to be the recipient, and server the sender
		data.ClientAddress.Device = r.Device().Address()
		data.ServerAddress.Device = message.DeviceRemote.Address()
	} else if data.ClientAddress.Device == nil {
		// only the recipient address may be missing
		data.ClientAddress.Device = r.Device().Address()
	} else if data.ServerAddress.Device == nil {
		// only the recipient address may be missing
		data.ServerAddress.Device = r.Device().Address()
	}

	return data
}
