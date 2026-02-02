package spine

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"sync"
	"time"

	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

type FeatureLocal struct {
	*Feature

	entity              api.EntityLocalInterface
	functionDataMap     map[model.FunctionType]api.FunctionDataCmdInterface
	muxResponseCB       sync.Mutex
	responseMsgCallback map[model.MsgCounterType][]func(result api.ResponseMessage)
	resultCallbacks     []func(result api.ResponseMessage)

	writeTimeout           time.Duration
	writeApprovalCallbacks []api.WriteApprovalCallbackFunc
	muxWriteReceived       sync.Mutex
	writeApprovalReceived  map[string]map[model.MsgCounterType]int
	pendingWriteApprovals  map[string]map[model.MsgCounterType]*time.Timer

	mux sync.Mutex
}

func NewFeatureLocal(id uint, entity api.EntityLocalInterface, ftype model.FeatureTypeType, role model.RoleType) *FeatureLocal {
	res := &FeatureLocal{
		Feature: NewFeature(
			featureAddressType(id, entity.Address()),
			ftype,
			role),
		entity:                entity,
		functionDataMap:       make(map[model.FunctionType]api.FunctionDataCmdInterface),
		responseMsgCallback:   make(map[model.MsgCounterType][]func(result api.ResponseMessage)),
		writeApprovalReceived: make(map[string]map[model.MsgCounterType]int),
		pendingWriteApprovals: make(map[string]map[model.MsgCounterType]*time.Timer),
		writeTimeout:          defaultMaxResponseDelay,
	}

	for _, fd := range CreateFunctionData[api.FunctionDataCmdInterface](ftype) {
		res.functionDataMap[fd.FunctionType()] = fd
	}
	res.operations = make(map[model.FunctionType]api.OperationsInterface)

	return res
}

var _ api.FeatureLocalInterface = (*FeatureLocal)(nil)

/* FeatureLocalInterface */

func (r *FeatureLocal) Device() api.DeviceLocalInterface {
	return r.entity.Device()
}

func (r *FeatureLocal) Entity() api.EntityLocalInterface {
	return r.entity
}

// Add supported function to the feature if its role is Server or Special
func (r *FeatureLocal) AddFunctionType(function model.FunctionType, read, write bool) {
	if r.role != model.RoleTypeServer && r.role != model.RoleTypeSpecial {
		return
	}
	if r.operations[function] != nil {
		return
	}
	writePartial := false
	if write {
		// partials are not supported on all features and functions, so check if this function supports it
		if fctData := r.functionData(function); fctData != nil {
			writePartial = fctData.SupportsPartialWrite()
		}
	}
	// Partial reads are intentionally not supported (spec-compliant design decision)
	// SPINE specification section 5.3.4.5 states: "A server MAY ignore unsupported cmdOption 
	// combinations and then replies with more than the requested parts instead."
	// By setting readPartial to false, we ensure all read requests return full data,
	// which provides the safest interoperability behavior for multi-vendor scenarios.
	r.operations[function] = NewOperations(read, false, write, writePartial)

	if r.role == model.RoleTypeServer &&
		r.ftype == model.FeatureTypeTypeDeviceDiagnosis &&
		function == model.FunctionTypeDeviceDiagnosisHeartbeatData {
		// Update HeartbeatManager
		r.Entity().HeartbeatManager().SetLocalFeature(r.Entity(), r)
	}
}

func (r *FeatureLocal) Functions() []model.FunctionType {
	var fcts []model.FunctionType

	for key := range r.operations {
		fcts = append(fcts, key)
	}

	return fcts
}

// Add a callback function to be invoked when SPINE message comes in with a given msgCounterReference value
//
// Returns an error if the provided callback function for the msgCounter is already set
func (r *FeatureLocal) AddResponseCallback(msgCounterReference model.MsgCounterType, function func(msg api.ResponseMessage)) error {
	r.muxResponseCB.Lock()
	defer r.muxResponseCB.Unlock()

	if _, ok := r.responseMsgCallback[msgCounterReference]; ok {
		for _, cb := range r.responseMsgCallback[msgCounterReference] {
			if reflect.ValueOf(cb).Pointer() == reflect.ValueOf(function).Pointer() {
				return errors.New("callback already set")
			}
		}
	}

	r.responseMsgCallback[msgCounterReference] = append(r.responseMsgCallback[msgCounterReference], function)

	return nil
}

func (r *FeatureLocal) processResponseMsgCallbacks(msgCounterReference model.MsgCounterType, msg api.ResponseMessage) {
	r.muxResponseCB.Lock()
	defer r.muxResponseCB.Unlock()

	cbs, ok := r.responseMsgCallback[msgCounterReference]
	if !ok {
		return
	}

	for _, cb := range cbs {
		go cb(msg)
	}

	delete(r.responseMsgCallback, msgCounterReference)
}

// Add a callback function to be invoked when a result message comes in for this feature
func (r *FeatureLocal) AddResultCallback(function func(msg api.ResponseMessage)) {
	r.muxResponseCB.Lock()
	defer r.muxResponseCB.Unlock()

	r.resultCallbacks = append(r.resultCallbacks, function)
}

func (r *FeatureLocal) processResultCallbacks(msg api.ResponseMessage) {
	r.muxResponseCB.Lock()
	defer r.muxResponseCB.Unlock()

	for _, cb := range r.resultCallbacks {
		go cb(msg)
	}
}

func (r *FeatureLocal) AddWriteApprovalCallback(function api.WriteApprovalCallbackFunc) error {
	if r.Role() != model.RoleTypeServer {
		return errors.New("only allowed on a server feature")
	}

	r.muxResponseCB.Lock()
	defer r.muxResponseCB.Unlock()

	r.writeApprovalCallbacks = append(r.writeApprovalCallbacks, function)

	return nil
}

func (r *FeatureLocal) processWriteApprovalCallbacks(msg *api.Message) {
	r.muxResponseCB.Lock()
	defer r.muxResponseCB.Unlock()

	for _, cb := range r.writeApprovalCallbacks {
		go cb(msg)
	}
}

func (r *FeatureLocal) addPendingApproval(msg *api.Message) {
	if r.Role() != model.RoleTypeServer ||
		msg.DeviceRemote == nil ||
		msg.RequestHeader == nil ||
		msg.RequestHeader.MsgCounter == nil {
		return
	}

	ski := msg.DeviceRemote.Ski()

	newTimer := time.AfterFunc(r.writeTimeout, func() {
		r.muxResponseCB.Lock()
		delete(r.pendingWriteApprovals[ski], *msg.RequestHeader.MsgCounter)
		r.muxResponseCB.Unlock()

		err := model.NewErrorTypeFromString("write not approved in time by application")
		_ = msg.FeatureRemote.Device().Sender().ResultError(msg.RequestHeader, r.Address(), err)
	})

	r.muxResponseCB.Lock()
	if _, ok := r.pendingWriteApprovals[ski]; !ok {
		r.pendingWriteApprovals[ski] = make(map[model.MsgCounterType]*time.Timer)
	}
	r.pendingWriteApprovals[ski][*msg.RequestHeader.MsgCounter] = newTimer
	r.muxResponseCB.Unlock()
}

func (r *FeatureLocal) ApproveOrDenyWrite(msg *api.Message, err model.ErrorType) {
	if r.Role() != model.RoleTypeServer ||
		msg.DeviceRemote == nil {
		return
	}

	ski := msg.DeviceRemote.Ski()

	r.muxResponseCB.Lock()
	timer, ok := r.pendingWriteApprovals[ski][*msg.RequestHeader.MsgCounter]
	count := len(r.writeApprovalCallbacks)
	r.muxResponseCB.Unlock()

	// if there is no timer running, we are too late and error has already been sent
	if !ok || timer == nil {
		return
	}

	// do we have enough approvals?
	r.muxWriteReceived.Lock()
	defer r.muxWriteReceived.Unlock()
	if count > 1 && err.ErrorNumber == 0 {
		amount, ok := r.writeApprovalReceived[ski][*msg.RequestHeader.MsgCounter]
		if ok {
			r.writeApprovalReceived[ski][*msg.RequestHeader.MsgCounter] = amount + 1
		} else {
			r.writeApprovalReceived[ski] = make(map[model.MsgCounterType]int)
			r.writeApprovalReceived[ski][*msg.RequestHeader.MsgCounter] = 1
		}
		// do we have enough approve messages, if not exit
		if r.writeApprovalReceived[ski][*msg.RequestHeader.MsgCounter] < count {
			return
		}
	}

	timer.Stop()

	delete(r.writeApprovalReceived[ski], *msg.RequestHeader.MsgCounter)

	r.muxResponseCB.Lock()
	defer r.muxResponseCB.Unlock()
	delete(r.pendingWriteApprovals[ski], *msg.RequestHeader.MsgCounter)

	if err.ErrorNumber == 0 {
		r.processWrite(msg)
		return
	}

	_ = msg.FeatureRemote.Device().Sender().ResultError(msg.RequestHeader, r.Address(), &err)
}

func (r *FeatureLocal) SetWriteApprovalTimeout(duration time.Duration) {
	r.writeTimeout = duration
}

func (r *FeatureLocal) CleanWriteApprovalCaches(ski string) {
	r.muxResponseCB.Lock()
	defer r.muxResponseCB.Unlock()

	delete(r.pendingWriteApprovals, ski)
	delete(r.writeApprovalReceived, ski)
}

// Remove subscriptions and bindings from local cache for a remote device
// used if a remote device is getting disconnected
func (r *FeatureLocal) CleanRemoteDeviceCaches(remoteAddress *model.DeviceAddressType) {
	if remoteAddress == nil ||
		remoteAddress.Device == nil {
		return
	}

	remoteDevice := r.Device().RemoteDeviceForAddress(*remoteAddress.Device)
	r.Device().BindingManager().RemoveBindingsForRemoteDevice(remoteDevice)
	r.Device().SubscriptionManager().RemoveSubscriptionsForRemoteDevice(remoteDevice)
}

// Remove subscriptions and bindings from local cache for a remote entity
// used if a remote entity is removed
func (r *FeatureLocal) CleanRemoteEntityCaches(remoteAddress *model.EntityAddressType) {
	if remoteAddress == nil ||
		remoteAddress.Device == nil ||
		remoteAddress.Entity == nil {
		return
	}

	remoteDevice := r.Device().RemoteDeviceForAddress(*remoteAddress.Device)
	if remoteDevice == nil {
		return
	}
	remoteEntity := remoteDevice.Entity(remoteAddress.Entity)
	if remoteEntity == nil {
		return
	}
	r.Device().BindingManager().RemoveBindingsForRemoteEntity(remoteEntity)
	r.Device().SubscriptionManager().RemoveSubscriptionsForRemoteEntity(remoteEntity)
}

func (r *FeatureLocal) DataCopy(function model.FunctionType) any {
	r.mux.Lock()
	defer r.mux.Unlock()

	fctData := r.functionData(function)
	if fctData == nil {
		return nil
	}

	return fctData.DataCopyAny()
}

func (r *FeatureLocal) SetData(function model.FunctionType, data any) {
	fctData, err := r.updateData(false, function, data, nil, nil)

	if err != nil {
		logging.Log().Debug(err.String())
	}

	if fctData != nil && err == nil {
		// do not notify subscribers for the following data functions:
		// - FunctionTypeNodeManagementBindingData
		// - FunctionTypeNodeManagementSubscriptionData
		// because the send out data would have to be filtered for the recipient,
		// partial data for the models aren't supported and filtering on top of this
		// is also not supported. Also no other implementations uses this data or
		// provides it.
		ignoreNotify := []model.FunctionType{
			model.FunctionTypeNodeManagementBindingData,
			model.FunctionTypeNodeManagementSubscriptionData,
		}

		if !slices.Contains(ignoreNotify, function) {
			r.Device().NotifySubscribers(r.Address(), fctData.NotifyOrWriteCmdType(nil, nil, false, nil))
		}
	}
}

func (r *FeatureLocal) UpdateData(function model.FunctionType, data any, filterPartial *model.FilterType, filterDelete *model.FilterType) *model.ErrorType {
	fctData, err := r.updateData(false, function, data, filterPartial, filterDelete)

	if err != nil {
		logging.Log().Debug(err.String())
	}

	if fctData != nil && err == nil {
		var deleteSelector, deleteElements, partialSelector any

		cmdFunction := util.Ptr(function)
		if filterDelete != nil {
			if fDelete, err := filterDelete.Data(cmdFunction); err == nil {
				if fDelete.Selector != nil {
					deleteSelector = fDelete.Selector
				}
				if fDelete.Elements != nil {
					deleteElements = fDelete.Elements
				}
			}
		}

		if filterPartial != nil {
			if fPartial, err := filterPartial.Data(cmdFunction); err == nil && fPartial.Selector != nil {
				partialSelector = fPartial.Selector
			}
		}

		r.Device().NotifySubscribers(r.Address(), fctData.NotifyOrWriteCmdType(deleteSelector, partialSelector, partialSelector == nil, deleteElements))
	}

	return err
}

func (r *FeatureLocal) updateData(remoteWrite bool, function model.FunctionType, data any, filterPartial *model.FilterType, filterDelete *model.FilterType) (api.FunctionDataCmdInterface, *model.ErrorType) {
	r.mux.Lock()
	defer r.mux.Unlock()

	fctData := r.functionData(function)
	if fctData == nil {
		return nil, model.NewErrorType(model.ErrorNumberTypeCommandNotSupported, "data not found")
	}

	// Pass the function type to UpdateDataAny for filter context
	cmdFunction := util.Ptr(function)
	_, err := fctData.UpdateDataAny(remoteWrite, true, data, filterPartial, filterDelete, cmdFunction)

	return fctData, err
}

func (r *FeatureLocal) RequestRemoteData(
	function model.FunctionType,
	selector any,
	elements any,
	destination api.FeatureRemoteInterface) (*model.MsgCounterType, *model.ErrorType) {
	fd := r.functionData(function)
	if fd == nil {
		return nil, model.NewErrorType(model.ErrorNumberTypeCommandNotSupported, "function data not found")
	}

	cmd := fd.ReadCmdType(selector, elements)

	return r.RequestRemoteDataBySenderAddress(cmd, destination.Device().Sender(), destination.Device().Ski(), destination.Address(), destination.MaxResponseDelayDuration())
}

func (r *FeatureLocal) RequestRemoteDataBySenderAddress(
	cmd model.CmdType,
	sender api.SenderInterface,
	deviceSki string,
	destinationAddress *model.FeatureAddressType,
	maxDelay time.Duration) (*model.MsgCounterType, *model.ErrorType) {
	// Note: maxDelay parameter is informational only and not used for timeout detection
	// Read request timeouts are not implemented in spine-go (per SPINE spec MAY requirement)
	msgCounter, err := sender.Request(model.CmdClassifierTypeRead, r.Address(), destinationAddress, false, []model.CmdType{cmd})
	if err == nil {
		return msgCounter, nil
	}

	return msgCounter, model.NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
}

// check if there already is a subscription to a remote feature
func (r *FeatureLocal) HasSubscriptionToRemote(remoteAddress *model.FeatureAddressType) bool {
	// subscriptions are also valid on NodeManagement, which has role Special
	// so to cover all cases, any of the combinations of client/server roles should be checked
	asClient := r.Device().SubscriptionManager().HasSubscription(r.Address(), remoteAddress)
	asServer := r.Device().SubscriptionManager().HasSubscription(remoteAddress, r.Address())
	return asClient || asServer
}

// SubscribeToRemote to a remote feature
//
// Returns:
// - msgCounter: the message counter reference for the request, nil if the subscription already exists or an error occurred
// - error: an error if creating the subscription request failed or sending failed, or nil if the subscription already exists or sending the request was possible
func (r *FeatureLocal) SubscribeToRemote(remoteAddress *model.FeatureAddressType) (*model.MsgCounterType, *model.ErrorType) {
	if remoteAddress.Device == nil {
		return nil, model.NewErrorTypeFromString("device not found")
	}
	remoteDevice := r.entity.Device().RemoteDeviceForAddress(*remoteAddress.Device)
	if remoteDevice == nil {
		return nil, model.NewErrorTypeFromString("device not found")
	}

	if r.Role() == model.RoleTypeServer {
		return nil, model.NewErrorTypeFromString(fmt.Sprintf("the server feature '%s' cannot request a subscription", r.Feature.String()))
	}

	// check if we already have this subscription
	if r.HasSubscriptionToRemote(remoteAddress) {
		return nil, nil
	}

	remoteFeature := remoteDevice.FeatureByAddress(remoteAddress)
	remoteFeatureType := remoteFeature.Type()
	if remoteFeature.Role() == model.RoleTypeClient {
		return nil, model.NewErrorTypeFromString(fmt.Sprintf("remote feature '%s' is not a server", remoteFeature.String()))
	}

	msgCounter, err := remoteDevice.Sender().Subscribe(r.Address(), remoteAddress, remoteFeatureType)
	if err != nil {
		return nil, model.NewErrorTypeFromString(err.Error())
	}

	_ = r.AddResponseCallback(*msgCounter, func(msg api.ResponseMessage) {
		r.subscribeResponseCallback(remoteDevice, remoteAddress, remoteFeatureType, msg)
	})

	return msgCounter, nil
}

func (r *FeatureLocal) subscribeResponseCallback(
	remoteDevice api.DeviceRemoteInterface,
	remoteAddress *model.FeatureAddressType,
	fType model.FeatureTypeType,
	msg api.ResponseMessage) {
	resultData, ok := msg.Data.(*model.ResultDataType)
	if !ok || resultData.ErrorNumber == nil {
		return
	}

	// only add the subscription if it was successful
	if *resultData.ErrorNumber == 0 {
		data := model.SubscriptionManagementRequestCallType{
			ClientAddress:     r.Address(),
			ServerAddress:     remoteAddress,
			ServerFeatureType: &fType,
		}

		if err := r.Device().SubscriptionManager().AddSubscription(remoteDevice, data); err != nil {
			logging.Log().Debug("Adding accepted remote subscription failed", err)
		}
	}
}

// Remove a subscriptions to a remote feature
func (r *FeatureLocal) RemoveRemoteSubscription(remoteAddress *model.FeatureAddressType) (*model.MsgCounterType, *model.ErrorType) {
	if remoteAddress.Device == nil {
		return nil, model.NewErrorTypeFromString("device not found")
	}
	remoteDevice := r.entity.Device().RemoteDeviceForAddress(*remoteAddress.Device)
	if remoteDevice == nil {
		return nil, model.NewErrorTypeFromString("device not found")
	}

	msgCounter, err := remoteDevice.Sender().Unsubscribe(r.Address(), remoteAddress)
	if err != nil {
		return nil, model.NewErrorTypeFromString("device not found")
	}

	_ = r.AddResponseCallback(*msgCounter, func(msg api.ResponseMessage) {
		r.unsubscribeResponseCallback(remoteDevice, remoteAddress, msg)
	})

	return msgCounter, nil
}

func (r *FeatureLocal) unsubscribeResponseCallback(
	remoteDevice api.DeviceRemoteInterface,
	remoteAddress *model.FeatureAddressType,
	msg api.ResponseMessage) {
	resultData, ok := msg.Data.(*model.ResultDataType)
	if !ok || resultData.ErrorNumber == nil {
		return
	}

	// only remove the subscription if the removal was successful
	if *resultData.ErrorNumber == 0 {
		var data model.SubscriptionManagementDeleteCallType

		if r.role == model.RoleTypeServer {
			data.ClientAddress = remoteAddress
			data.ServerAddress = r.Address()
		} else {
			data.ClientAddress = r.Address()
			data.ServerAddress = remoteAddress
		}

		if err := r.Device().SubscriptionManager().RemoveSubscription(remoteDevice, data); err != nil {
			logging.Log().Debug("Removing binding to remote feature failed", err)
		}
	}
}

// check if there already is a binding to a remote feature
func (r *FeatureLocal) HasBindingToRemote(remoteAddress *model.FeatureAddressType) bool {
	if r.role == model.RoleTypeClient {
		return r.Device().BindingManager().HasBinding(r.Address(), remoteAddress)
	}

	return r.Device().BindingManager().HasBinding(remoteAddress, r.Address())
}

// Request a binding to a remote feature
//
// Returns:
// - msgCounter: the message counter reference for the request, nil if the binding already exists or an error occurred
// - error: an error if creating the binding request failed or sending failed, or nil if the binding already exists or sending the request was possible
func (r *FeatureLocal) BindToRemote(remoteAddress *model.FeatureAddressType) (*model.MsgCounterType, *model.ErrorType) {
	if remoteAddress.Device == nil {
		return nil, model.NewErrorTypeFromString("device not found")
	}
	remoteDevice := r.entity.Device().RemoteDeviceForAddress(*remoteAddress.Device)
	if remoteDevice == nil {
		return nil, model.NewErrorTypeFromString("device not found")
	}

	if r.Role() == model.RoleTypeServer {
		return nil, model.NewErrorTypeFromString(fmt.Sprintf("the server feature '%s' cannot request a binding", r.Feature.String()))
	}

	// check if we already have this binding
	if r.HasBindingToRemote(remoteAddress) {
		return nil, nil
	}

	remoteFeature := remoteDevice.FeatureByAddress(remoteAddress)
	remoteFeatureType := remoteFeature.Type()
	if remoteFeature.Role() == model.RoleTypeClient {
		return nil, model.NewErrorTypeFromString(fmt.Sprintf("remote feature '%s' is not a server", remoteFeature.String()))
	}

	msgCounter, err := remoteDevice.Sender().Bind(r.Address(), remoteAddress, remoteFeatureType)
	if err != nil {
		return nil, model.NewErrorTypeFromString(err.Error())
	}

	_ = r.AddResponseCallback(*msgCounter, func(msg api.ResponseMessage) {
		r.bindResponseCallback(remoteDevice, remoteAddress, remoteFeatureType, msg)
	})

	return msgCounter, nil
}

func (r *FeatureLocal) bindResponseCallback(
	remoteDevice api.DeviceRemoteInterface,
	remoteAddress *model.FeatureAddressType,
	fType model.FeatureTypeType,
	msg api.ResponseMessage) {
	resultData, ok := msg.Data.(*model.ResultDataType)
	if !ok || resultData.ErrorNumber == nil {
		return
	}

	// only add the binding if it was successful
	if *resultData.ErrorNumber == 0 {
		data := model.BindingManagementRequestCallType{
			ClientAddress:     r.Address(),
			ServerAddress:     remoteAddress,
			ServerFeatureType: &fType,
		}

		if err := r.Device().BindingManager().AddBinding(remoteDevice, data); err != nil {
			logging.Log().Debug("Adding accepted remote binding failed", err)
		}
	}
}

// Send a request to remove a binding with a remote feature
func (r *FeatureLocal) RemoveRemoteBinding(remoteAddress *model.FeatureAddressType) (*model.MsgCounterType, *model.ErrorType) {
	if remoteAddress.Device == nil {
		return nil, model.NewErrorTypeFromString("device not found")
	}
	remoteDevice := r.entity.Device().RemoteDeviceForAddress(*remoteAddress.Device)
	if remoteDevice == nil {
		return nil, model.NewErrorTypeFromString("device not found")
	}

	msgCounter, err := remoteDevice.Sender().Unbind(r.Address(), remoteAddress)
	if err != nil {
		return nil, model.NewErrorTypeFromString(err.Error())
	}

	_ = r.AddResponseCallback(*msgCounter, func(msg api.ResponseMessage) {
		r.unbindResponseCallback(remoteDevice, remoteAddress, msg)
	})

	return msgCounter, nil
}

func (r *FeatureLocal) unbindResponseCallback(
	remoteDevice api.DeviceRemoteInterface,
	remoteAddress *model.FeatureAddressType,
	msg api.ResponseMessage) {
	resultData, ok := msg.Data.(*model.ResultDataType)
	if !ok || resultData.ErrorNumber == nil {
		return
	}

	// only remove the binding if the removal was successful
	if *resultData.ErrorNumber == 0 {
		var data model.BindingManagementDeleteCallType

		if r.Role() == model.RoleTypeServer {
			data.ClientAddress = remoteAddress
			data.ServerAddress = r.Address()
		} else {
			data.ClientAddress = r.Address()
			data.ServerAddress = remoteAddress
		}

		if err := r.Device().BindingManager().RemoveBinding(remoteDevice, data); err != nil {
			logging.Log().Debug("Removing binding to remote feature failed", err)
		}
	}
}

func (r *FeatureLocal) HandleMessage(message *api.Message) *model.ErrorType {
	cmdData, err := message.Cmd.Data()
	if err != nil {
		return model.NewErrorType(model.ErrorNumberTypeCommandNotSupported, err.Error())
	}
	if cmdData.Function == nil {
		return model.NewErrorType(model.ErrorNumberTypeCommandNotSupported, "No function found for cmd data")
	}

	switch message.CmdClassifier {
	case model.CmdClassifierTypeResult:
		if err := r.processResult(message); err != nil {
			return err
		}
	case model.CmdClassifierTypeRead:
		if err := r.processRead(*cmdData.Function, message.RequestHeader, message.FeatureRemote); err != nil {
			return err
		}
	case model.CmdClassifierTypeReply:
		if err := r.processReply(message); err != nil {
			return err
		}
	case model.CmdClassifierTypeNotify:
		if err := r.processNotify(*cmdData.Function, cmdData.Value, message.FilterPartial, message.FilterDelete, message.FeatureRemote); err != nil {
			return err
		}
	case model.CmdClassifierTypeWrite:
		// if there is a write permission check callback set, invoke this instead of directly allowing the write
		if len(r.writeApprovalCallbacks) > 0 {
			r.addPendingApproval(message)
			r.processWriteApprovalCallbacks(message)
		} else {
			// this method handles ack and error results, so no need to return an error
			r.processWrite(message)
		}
	default:
		return model.NewErrorTypeFromString(fmt.Sprintf("CmdClassifier not implemented: %s", message.CmdClassifier))
	}

	return nil
}

func (r *FeatureLocal) processResult(message *api.Message) *model.ErrorType {
	if message.Cmd.ResultData == nil || message.Cmd.ResultData.ErrorNumber == nil {
		return model.NewErrorType(
			model.ErrorNumberTypeGeneralError,
			fmt.Sprintf("ResultData CmdClassifierType %s not implemented", message.CmdClassifier))
	}

	if *message.Cmd.ResultData.ErrorNumber != model.ErrorNumberTypeNoError {
		// error numbers explained in Resource Spec 3.11
		errorString := fmt.Sprintf("Error Result received %d", *message.Cmd.ResultData.ErrorNumber)
		if message.Cmd.ResultData.Description != nil {
			errorString += fmt.Sprintf(": %s", *message.Cmd.ResultData.Description)
		}
		logging.Log().Debug(errorString)
	}

	// we don't need to populate this message if there is no MsgCounterReference
	if message.RequestHeader == nil || message.RequestHeader.MsgCounterReference == nil {
		return nil
	}

	responseMsg := api.ResponseMessage{
		MsgCounterReference: *message.RequestHeader.MsgCounterReference,
		Data:                message.Cmd.ResultData,
		FeatureLocal:        r,
		FeatureRemote:       message.FeatureRemote,
		EntityRemote:        message.EntityRemote,
		DeviceRemote:        message.DeviceRemote,
	}

	r.processResponseMsgCallbacks(*message.RequestHeader.MsgCounterReference, responseMsg)
	r.processResultCallbacks(responseMsg)

	return nil
}

func (r *FeatureLocal) processRead(function model.FunctionType, requestHeader *model.HeaderType, featureRemote api.FeatureRemoteInterface) *model.ErrorType {
	// is this a read request to a local server/special feature?
	if r.role == model.RoleTypeClient {
		// Read requests to a client feature are not allowed
		return model.NewErrorTypeFromNumber(model.ErrorNumberTypeCommandRejected)
	}

	fd := r.functionData(function)
	if fd == nil {
		return model.NewErrorType(model.ErrorNumberTypeCommandNotSupported, "function data not found")
	}

	// SPEC-COMPLIANT BEHAVIOR: Partial filters are intentionally ignored
	// 
	// The incoming message may contain FilterPartial with element selectors,
	// selectors, or other cmdOptions, but we always reply with full data.
	// This implements SPINE specification section 5.3.4.5:
	// "A server MAY ignore unsupported cmdOption combinations and then replies 
	// with more than the requested parts instead."
	//
	// Benefits of this approach:
	// 1. Ensures interoperability - no partial read implementation variations
	// 2. Prevents data inconsistency in multi-vendor scenarios
	// 3. Provides predictable behavior for clients
	// 4. Complies with spec requirement for unsupported cmdOptions
	cmd := fd.ReplyCmdType(false) // false = full data, ignore any partial filters
	if err := featureRemote.Device().Sender().Reply(requestHeader, r.Address(), cmd); err != nil {
		return model.NewErrorTypeFromString(err.Error())
	}

	return nil
}

func (r *FeatureLocal) processReply(message *api.Message) *model.ErrorType {
	// function model.FunctionType, data any, filterPartial *model.FilterType, filterDelete *model.FilterType, featureRemote api.FeatureRemoteInterface)

	// the error is handled already in the caller
	cmdData, _ := message.Cmd.Data()
	featureRemote := message.FeatureRemote

	if _, err := featureRemote.UpdateData(true, *cmdData.Function, cmdData.Value, message.FilterPartial, message.FilterDelete); err != nil {
		return err
	}

	// the data was updated, so send an event, other event handlers may watch out for this as well
	payload := api.EventPayload{
		Ski:           featureRemote.Device().Ski(),
		EventType:     api.EventTypeDataChange,
		ChangeType:    api.ElementChangeUpdate,
		Feature:       featureRemote,
		Device:        featureRemote.Device(),
		Entity:        featureRemote.Entity(),
		LocalFeature:  r,
		Function:      *cmdData.Function,
		CmdClassifier: util.Ptr(model.CmdClassifierTypeReply),
		Data:          cmdData.Value,
	}
	r.Device().Events().Publish(payload)

	// we don't need to populate this message if there is no MsgCounterReference
	if message.RequestHeader == nil || message.RequestHeader.MsgCounterReference == nil {
		return nil
	}

	responseMsg := api.ResponseMessage{
		MsgCounterReference: *message.RequestHeader.MsgCounterReference,
		Data:                cmdData.Value,
		FeatureLocal:        r,
		FeatureRemote:       message.FeatureRemote,
		EntityRemote:        message.EntityRemote,
		DeviceRemote:        message.DeviceRemote,
	}

	r.processResponseMsgCallbacks(*message.RequestHeader.MsgCounterReference, responseMsg)

	return nil
}

func (r *FeatureLocal) processNotify(function model.FunctionType, data any, filterPartial *model.FilterType, filterDelete *model.FilterType, featureRemote api.FeatureRemoteInterface) *model.ErrorType {
	if _, err := featureRemote.UpdateData(true, function, data, filterPartial, filterDelete); err != nil {
		return err
	}

	payload := api.EventPayload{
		Ski:           featureRemote.Device().Ski(),
		EventType:     api.EventTypeDataChange,
		ChangeType:    api.ElementChangeUpdate,
		Feature:       featureRemote,
		Device:        featureRemote.Device(),
		Entity:        featureRemote.Entity(),
		LocalFeature:  r,
		Function:      function,
		CmdClassifier: util.Ptr(model.CmdClassifierTypeNotify),
		Data:          data,
	}
	r.Device().Events().Publish(payload)

	return nil
}

func (r *FeatureLocal) processWrite(msg *api.Message) {
	if err := r.executeWrite(msg); err != nil {
		_ = msg.FeatureRemote.Device().Sender().ResultError(msg.RequestHeader, r.Address(), err)
	} else if msg.RequestHeader != nil {
		ackRequest := msg.RequestHeader.AckRequest
		if ackRequest != nil && *ackRequest {
			_ = msg.FeatureRemote.Device().Sender().ResultSuccess(msg.RequestHeader, r.Address())
		}
	}
}

func (r *FeatureLocal) executeWrite(msg *api.Message) *model.ErrorType {
	cmdData, err := msg.Cmd.Data()
	if err != nil {
		return model.NewErrorType(model.ErrorNumberTypeCommandNotSupported, err.Error())
	}
	if cmdData.Function == nil {
		return model.NewErrorType(model.ErrorNumberTypeCommandNotSupported, "No function found for cmd data")
	}

	fctData, err1 := r.updateData(true, *cmdData.Function, cmdData.Value, msg.FilterPartial, msg.FilterDelete)
	if err1 != nil {
		return err1
	} else if fctData == nil {
		return model.NewErrorType(model.ErrorNumberTypeCommandNotSupported, "function not found")
	}

	r.Device().NotifySubscribers(r.Address(), fctData.NotifyOrWriteCmdType(nil, nil, false, nil))

	payload := api.EventPayload{
		Ski:           msg.FeatureRemote.Device().Ski(),
		EventType:     api.EventTypeDataChange,
		ChangeType:    api.ElementChangeUpdate,
		Feature:       msg.FeatureRemote,
		Device:        msg.FeatureRemote.Device(),
		Entity:        msg.FeatureRemote.Entity(),
		LocalFeature:  r,
		Function:      *cmdData.Function,
		CmdClassifier: util.Ptr(model.CmdClassifierTypeWrite),
		Data:          cmdData.Value,
	}
	r.Device().Events().Publish(payload)

	return nil
}

func (r *FeatureLocal) functionData(function model.FunctionType) api.FunctionDataCmdInterface {
	fd, found := r.functionDataMap[function]
	if !found {
		logging.Log().Errorf("Data was not found for function '%s'", function)
		return nil
	}
	return fd
}

func (r *FeatureLocal) Information() *model.NodeManagementDetailedDiscoveryFeatureInformationType {
	var funs []model.FunctionPropertyType
	for fun, operations := range r.operations {
		var functionType = model.FunctionType(fun)
		sf := model.FunctionPropertyType{
			Function:           &functionType,
			PossibleOperations: operations.Information(),
		}

		funs = append(funs, sf)
	}

	return model.NewFeatureInformationForNodeManagement(r.address.Entity, r.address.Feature, &r.ftype, &r.role, r.description, funs)
}
