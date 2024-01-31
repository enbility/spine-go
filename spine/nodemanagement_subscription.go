package spine

import (
	"fmt"

	"github.com/ahmetb/go-linq/v3"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
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

	var remoteDeviceSubscriptions []model.SubscriptionManagementEntryDataType
	remoteDeviceSubscriptionEntries := r.Device().SubscriptionManager().Subscriptions(message.FeatureRemote.Device())
	linq.From(remoteDeviceSubscriptionEntries).SelectT(func(s *api.SubscriptionEntry) model.SubscriptionManagementEntryDataType {
		return model.SubscriptionManagementEntryDataType{
			SubscriptionId: util.Ptr(model.SubscriptionIdType(s.Id)),
			ServerAddress:  s.ServerFeature.Address(),
			ClientAddress:  s.ClientFeature.Address(),
		}
	}).ToSlice(&remoteDeviceSubscriptions)

	cmd := model.CmdType{
		NodeManagementSubscriptionData: &model.NodeManagementSubscriptionDataType{
			SubscriptionEntry: remoteDeviceSubscriptions,
		},
	}

	return message.FeatureRemote.Device().Sender().Reply(message.RequestHeader, r.Address(), cmd)
}

func (r *NodeManagement) handleMsgSubscriptionData(message *api.Message) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeCall:
		return r.processReadSubscriptionData(message)

	default:
		return fmt.Errorf("nodemanagement.handleSubscriptionDeleteCall: NodeManagementSubscriptionRequestCall CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}

func (r *NodeManagement) handleMsgSubscriptionRequestCall(message *api.Message, data *model.NodeManagementSubscriptionRequestCallType) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeCall:
		subscriptionMgr := r.Device().SubscriptionManager()

		return subscriptionMgr.AddSubscription(message.FeatureRemote.Device(), *data.SubscriptionRequest)

	default:
		return fmt.Errorf("nodemanagement.handleSubscriptionRequestCall: NodeManagementSubscriptionRequestCall CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}

func (r *NodeManagement) handleMsgSubscriptionDeleteCall(message *api.Message, data *model.NodeManagementSubscriptionDeleteCallType) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeCall:
		subscriptionMgr := r.Device().SubscriptionManager()

		return subscriptionMgr.RemoveSubscription(*data.SubscriptionDelete, message.FeatureRemote.Device())

	default:
		return fmt.Errorf("nodemanagement.handleSubscriptionDeleteCall: NodeManagementSubscriptionRequestCall CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}
