package spine

import (
	"errors"
	"fmt"
	"reflect"
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

	bindings      []*model.FeatureAddressType // bindings to remote features
	subscriptions []*model.FeatureAddressType // subscriptions to remote features

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
	// partial reads are currently not supported!
	r.operations[function] = NewOperations(read, false, write, writePartial)

	if r.role == model.RoleTypeServer &&
		r.ftype == model.FeatureTypeTypeDeviceDiagnosis &&
		function == model.FunctionTypeDeviceDiagnosisHeartbeatData {
		// Update HeartbeatManager
		r.Device().HeartbeatManager().SetLocalFeature(r.Entity(), r)
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
// Returns an error if there is already a callback for the msgCounter set
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

func (r *FeatureLocal) CleanCaches(ski string) {
	r.muxResponseCB.Lock()
	defer r.muxResponseCB.Unlock()

	delete(r.pendingWriteApprovals, ski)
	delete(r.writeApprovalReceived, ski)
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
		r.Device().NotifySubscribers(r.Address(), fctData.NotifyOrWriteCmdType(nil, nil, false, nil))
	}
}

func (r *FeatureLocal) UpdateData(function model.FunctionType, data any, filterPartial *model.FilterType, filterDelete *model.FilterType) *model.ErrorType {
	fctData, err := r.updateData(false, function, data, filterPartial, filterDelete)

	if err != nil {
		logging.Log().Debug(err.String())
	}

	if fctData != nil && err == nil {
		var deleteSelector, deleteElements, partialSelector any

		if filterDelete != nil {
			if fDelete, err := filterDelete.Data(); err == nil {
				if fDelete.Selector != nil {
					deleteSelector = fDelete.Selector
				}
				if fDelete.Elements != nil {
					deleteElements = fDelete.Elements
				}
			}
		}

		if filterPartial != nil {
			if fPartial, err := filterPartial.Data(); err == nil && fPartial.Selector != nil {
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
		return nil, model.NewErrorTypeFromString("data not found")
	}

	_, err := fctData.UpdateDataAny(remoteWrite, true, data, filterPartial, filterDelete)

	return fctData, err
}

func (r *FeatureLocal) RequestRemoteData(
	function model.FunctionType,
	selector any,
	elements any,
	destination api.FeatureRemoteInterface) (*model.MsgCounterType, *model.ErrorType) {
	fd := r.functionData(function)
	if fd == nil {
		return nil, model.NewErrorTypeFromString("function data not found")
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
	msgCounter, err := sender.Request(model.CmdClassifierTypeRead, r.Address(), destinationAddress, false, []model.CmdType{cmd})
	if err == nil {
		return msgCounter, nil
	}

	return msgCounter, model.NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
}

// check if there already is a subscription to a remote feature
func (r *FeatureLocal) HasSubscriptionToRemote(remoteAddress *model.FeatureAddressType) bool {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, item := range r.subscriptions {
		if reflect.DeepEqual(*remoteAddress, *item) {
			return true
		}
	}

	return false
}

// SubscribeToRemote to a remote feature
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

	msgCounter, err := remoteDevice.Sender().Subscribe(r.Address(), remoteAddress, r.ftype)
	if err != nil {
		return nil, model.NewErrorTypeFromString(err.Error())
	}

	r.mux.Lock()
	r.subscriptions = append(r.subscriptions, remoteAddress)
	r.mux.Unlock()

	return msgCounter, nil
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

	var subscriptions []*model.FeatureAddressType

	r.mux.Lock()
	defer r.mux.Unlock()

	for _, item := range r.subscriptions {
		if reflect.DeepEqual(item, remoteAddress) {
			continue
		}

		subscriptions = append(subscriptions, item)
	}

	r.subscriptions = subscriptions

	return msgCounter, nil
}

// Remove all subscriptions to remote features
func (r *FeatureLocal) RemoveAllRemoteSubscriptions() {
	for _, item := range r.subscriptions {
		_, _ = r.RemoveRemoteSubscription(item)
	}
}

// check if there already is a binding to a remote feature
func (r *FeatureLocal) HasBindingToRemote(remoteAddress *model.FeatureAddressType) bool {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, item := range r.bindings {
		if reflect.DeepEqual(*remoteAddress, *item) {
			return true
		}
	}

	return false
}

// BindToRemote to a remote feature
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

	msgCounter, err := remoteDevice.Sender().Bind(r.Address(), remoteAddress, r.ftype)
	if err != nil {
		return nil, model.NewErrorTypeFromString(err.Error())
	}

	r.mux.Lock()
	r.bindings = append(r.bindings, remoteAddress)
	r.mux.Unlock()

	return msgCounter, nil
}

// Remove a binding to a remote feature
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

	var bindings []*model.FeatureAddressType

	r.mux.Lock()
	defer r.mux.Unlock()

	for _, item := range r.bindings {
		if reflect.DeepEqual(item, remoteAddress) {
			continue
		}

		bindings = append(bindings, item)
	}

	r.bindings = bindings

	return msgCounter, nil
}

// Remove all subscriptions to remote features
func (r *FeatureLocal) RemoveAllRemoteBindings() {
	for _, item := range r.bindings {
		_, _ = r.RemoveRemoteBinding(item)
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
		return model.NewErrorTypeFromString("function data not found")
	}

	cmd := fd.ReplyCmdType(false)
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
	Events.Publish(payload)

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
	Events.Publish(payload)

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
		return model.NewErrorTypeFromString("function not found")
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
	Events.Publish(payload)

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

	res := model.NodeManagementDetailedDiscoveryFeatureInformationType{
		Description: &model.NetworkManagementFeatureDescriptionDataType{
			FeatureAddress:    r.Address(),
			FeatureType:       &r.ftype,
			Role:              &r.role,
			Description:       r.description,
			SupportedFunction: funs,
		},
	}

	return &res
}
