package spine

import (
	"fmt"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

func NewNodeManagementSubscriptionRequestCallType(clientAddress *model.FeatureAddressType, serverAddress *model.FeatureAddressType, featureType model.FeatureTypeType) *model.NodeManagementSubscriptionRequestCallType {
	return &model.NodeManagementSubscriptionRequestCallType{
		SubscriptionRequest: &model.SubscriptionManagementRequestCallType{
			ClientAddress:     clientAddress,
			ServerAddress:     serverAddress,
			ServerFeatureType: &featureType,
		},
	}
}

func NewNodeManagementSubscriptionDeleteCallType(clientAddress *model.FeatureAddressType, serverAddress *model.FeatureAddressType) *model.NodeManagementSubscriptionDeleteCallType {
	return &model.NodeManagementSubscriptionDeleteCallType{
		SubscriptionDelete: &model.SubscriptionManagementDeleteCallType{
			ClientAddress: clientAddress,
			ServerAddress: serverAddress,
		},
	}
}

// route subscription request calls to the appropriate feature implementation and add the subscription to the current list
func (r *NodeManagement) processReadSubscriptionData(message *api.Message) error {
	remoteDeviceSubscriptionEntries := r.Device().SubscriptionManager().SubscriptionsForRemoteDevice(message.FeatureRemote.Device())

	cmd := model.CmdType{
		NodeManagementSubscriptionData: &model.NodeManagementSubscriptionDataType{
			SubscriptionEntry: remoteDeviceSubscriptionEntries,
		},
	}

	return message.DeviceRemote.Sender().Reply(message.RequestHeader, r.Address(), cmd)
}

func (r *NodeManagement) handleMsgSubscriptionData(message *api.Message) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeRead:
		return r.processReadSubscriptionData(message)

	default:
		return fmt.Errorf("nodemanagement.handleSubscriptionDeleteCall: NodeManagementSubscriptionRequestCall CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}

func (r *NodeManagement) handleMsgSubscriptionRequestCall(message *api.Message, data *model.NodeManagementSubscriptionRequestCallType) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeCall:
		subscriptionMgr := r.Device().SubscriptionManager()

		readData := r.createSubscriptionAddMissingDeviceAddresses(message, data.SubscriptionRequest)

		return subscriptionMgr.AddSubscription(message.FeatureRemote.Device(), *readData)

	default:
		return fmt.Errorf("nodemanagement.handleSubscriptionRequestCall: NodeManagementSubscriptionRequestCall CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}

func (r *NodeManagement) handleMsgSubscriptionDeleteCall(message *api.Message, data *model.NodeManagementSubscriptionDeleteCallType) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeCall:
		subscriptionMgr := r.Device().SubscriptionManager()

		deleteData := r.deleteSubscriptionAddMissingDeviceAddresses(message, data.SubscriptionDelete)

		return subscriptionMgr.RemoveSubscription(message.FeatureRemote.Device(), *deleteData)

	default:
		return fmt.Errorf("nodemanagement.handleSubscriptionDeleteCall: NodeManagementSubscriptionRequestCall CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}

// adds potentially missing device addresses to the subscription data according to SPINE protocol spec 7.4.2
func (r *NodeManagement) createSubscriptionAddMissingDeviceAddresses(message *api.Message, data *model.SubscriptionManagementRequestCallType) *model.SubscriptionManagementRequestCallType {
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

// adds potentially missing device addresses to the subscription data according to SPINE protocol spec 7.4.4
func (r *NodeManagement) deleteSubscriptionAddMissingDeviceAddresses(message *api.Message, data *model.SubscriptionManagementDeleteCallType) *model.SubscriptionManagementDeleteCallType {
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
