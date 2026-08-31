package spine

import (
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

const NodeManagementFeatureId uint = 0

func NodeManagementAddress(deviceAddress *model.AddressDeviceType) *model.FeatureAddressType {
	return &model.FeatureAddressType{
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(NodeManagementFeatureId)),
		Device:  deviceAddress,
	}
}

var _ api.NodeManagementInterface = (*NodeManagement)(nil)

// maximum number of requests deferred per remote device while its detailed
// discovery is still pending
const maxDeferredRequests = 32

// maximum time a deferred request waits for the discovery to complete,
// variable for testing
var deferredRequestTimeout = defaultMaxResponseDelay

// a request waiting for the sending remote device to be discovered
type deferredRequest struct {
	message *api.Message
	timer   *time.Timer
}

type NodeManagement struct {
	*FeatureLocal

	// requests received before the sending remote device was discovered, by SKI
	deferred    map[string][]*deferredRequest
	deferredMux sync.Mutex
}

func NewNodeManagement(id uint, entity api.EntityLocalInterface) *NodeManagement {
	f := &NodeManagement{
		FeatureLocal: NewFeatureLocal(
			id, entity,
			model.FeatureTypeTypeNodeManagement,
			model.RoleTypeSpecial),
		deferred: make(map[string][]*deferredRequest),
	}

	f.AddFunctionType(model.FunctionTypeNodeManagementDetailedDiscoveryData, true, false)
	f.AddFunctionType(model.FunctionTypeNodeManagementUseCaseData, true, false)
	f.AddFunctionType(model.FunctionTypeNodeManagementSubscriptionData, true, false)
	f.AddFunctionType(model.FunctionTypeNodeManagementSubscriptionRequestCall, false, false)
	f.AddFunctionType(model.FunctionTypeNodeManagementSubscriptionDeleteCall, false, false)
	f.AddFunctionType(model.FunctionTypeNodeManagementBindingData, true, false)
	f.AddFunctionType(model.FunctionTypeNodeManagementBindingRequestCall, false, false)
	f.AddFunctionType(model.FunctionTypeNodeManagementBindingDeleteCall, false, false)
	if f.Device().FeatureSet() != nil && *f.Device().FeatureSet() != model.NetworkManagementFeatureSetTypeSimple {
		f.AddFunctionType(model.FunctionTypeNodeManagementDestinationListData, true, false)
	}

	return f
}

func (r *NodeManagement) Device() api.DeviceLocalInterface {
	return r.entity.Device()
}

func (r *NodeManagement) HandleMessage(message *api.Message) *model.ErrorType {
	switch {
	case message.Cmd.ResultData != nil:
		if err := r.processResult(message); err != nil {
			return err
		}

	case message.Cmd.NodeManagementDetailedDiscoveryData != nil:
		if err := r.handleMsgDetailedDiscoveryData(message, message.Cmd.NodeManagementDetailedDiscoveryData); err != nil {
			return model.NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
		}

	case message.Cmd.NodeManagementSubscriptionRequestCall != nil:
		if err := r.handleMsgSubscriptionRequestCall(message, message.Cmd.NodeManagementSubscriptionRequestCall); err != nil {
			return model.NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
		}

	case message.Cmd.NodeManagementSubscriptionDeleteCall != nil:
		if err := r.handleMsgSubscriptionDeleteCall(message, message.Cmd.NodeManagementSubscriptionDeleteCall); err != nil {
			return model.NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
		}

	case message.Cmd.NodeManagementSubscriptionData != nil:
		if err := r.handleMsgSubscriptionData(message); err != nil {
			return model.NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
		}

	case message.Cmd.NodeManagementBindingRequestCall != nil:
		if err := r.handleMsgBindingRequestCall(message, message.Cmd.NodeManagementBindingRequestCall); err != nil {
			return model.NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
		}

	case message.Cmd.NodeManagementBindingDeleteCall != nil:
		if err := r.handleMsgBindingDeleteCall(message, message.Cmd.NodeManagementBindingDeleteCall); err != nil {
			return model.NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
		}

	case message.Cmd.NodeManagementBindingData != nil:
		if err := r.handleMsgBindingData(message); err != nil {
			return model.NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
		}

	case message.Cmd.NodeManagementUseCaseData != nil:
		if err := r.handleMsgUseCaseData(message, message.Cmd.NodeManagementUseCaseData); err != nil {
			return model.NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
		}

	case message.Cmd.NodeManagementDestinationListData != nil:
		if err := r.handleMsgDestinationListData(message, message.Cmd.NodeManagementDestinationListData); err != nil {
			return model.NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
		}

	default:
		return model.NewErrorType(model.ErrorNumberTypeCommandNotSupported, fmt.Sprintf("nodemanagement.Handle: Cmd data not implemented: %s", message.Cmd.DataName()))
	}

	return nil
}

// subscription and binding requests need the sending device's address and features,
// which are both only known once its detailed discovery reply was processed
func deferrableRequest(message *api.Message) bool {
	return message.CmdClassifier == model.CmdClassifierTypeCall &&
		(message.Cmd.NodeManagementSubscriptionRequestCall != nil ||
			message.Cmd.NodeManagementBindingRequestCall != nil)
}

// Queue a request which arrived before the sending remote device was discovered.
// SPINE 7.4.1 rule 8 permits subscription requests at any time and 5.2.5.3 permits a
// delayed result, so the request is processed and answered once discovery completed
// instead of being rejected. Binding requests are handled alike.
//
// Returns true if the request was deferred and must not be processed or answered now.
func (r *NodeManagement) deferRequest(message *api.Message) bool {
	if !deferrableRequest(message) || message.DeviceRemote.Address() != nil {
		return false
	}

	ski := message.DeviceRemote.Ski()

	r.deferredMux.Lock()
	defer r.deferredMux.Unlock()

	if len(r.deferred[ski]) >= maxDeferredRequests {
		return false
	}

	entry := &deferredRequest{message: message}
	entry.timer = time.AfterFunc(deferredRequestTimeout, func() {
		r.timeoutDeferredRequest(ski, entry)
	})

	r.deferred[ski] = append(r.deferred[ski], entry)

	return true
}

// Answer a request which was not processed within the maximum response delay
func (r *NodeManagement) timeoutDeferredRequest(ski string, entry *deferredRequest) {
	r.deferredMux.Lock()
	index := slices.Index(r.deferred[ski], entry)
	if index >= 0 {
		r.deferred[ski] = slices.Delete(r.deferred[ski], index, index+1)
		if len(r.deferred[ski]) == 0 {
			delete(r.deferred, ski)
		}
	}
	r.deferredMux.Unlock()

	// the request was processed or dropped in the meantime
	if index < 0 {
		return
	}

	err := model.NewErrorType(model.ErrorNumberTypeTimeout, "device not discovered")
	_ = entry.message.DeviceRemote.Sender().ResultError(entry.message.RequestHeader, r.Address(), err)
}

// Process the requests deferred for a remote device and send their delayed results
func (r *NodeManagement) processDeferredRequests(remoteDevice api.DeviceRemoteInterface) {
	r.deferredMux.Lock()
	entries := r.deferred[remoteDevice.Ski()]
	delete(r.deferred, remoteDevice.Ski())
	r.deferredMux.Unlock()

	sender := remoteDevice.Sender()

	for _, entry := range entries {
		entry.timer.Stop()

		message := entry.message

		// the remote features were replaced by the discovery reply
		message.FeatureRemote = remoteDevice.FeatureByAddress(message.RequestHeader.AddressSource)
		if message.FeatureRemote == nil {
			err := model.NewErrorType(model.ErrorNumberTypeDestinationUnknown, "invalid remote feature address")
			_ = sender.ResultError(message.RequestHeader, r.Address(), err)
			continue
		}

		if err := r.HandleMessage(message); err != nil {
			_ = sender.ResultError(message.RequestHeader, r.Address(), err)
			continue
		}

		if ackRequest := message.RequestHeader.AckRequest; ackRequest != nil && *ackRequest {
			// return success as defined in SPINE chapter 5.2.4
			_ = sender.ResultSuccess(message.RequestHeader, r.Address())
		}
	}
}

// Drop the requests deferred for a remote device which is no longer connected
func (r *NodeManagement) removeDeferredRequests(ski string) {
	r.deferredMux.Lock()
	defer r.deferredMux.Unlock()

	for _, entry := range r.deferred[ski] {
		entry.timer.Stop()
	}

	delete(r.deferred, ski)
}
