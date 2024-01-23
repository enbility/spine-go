package spine

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

var _ api.FeatureLocalInterface = (*FeatureLocal)(nil)

type FeatureLocal struct {
	*Feature

	muxResultCB     sync.Mutex
	entity          api.EntityLocalInterface
	functionDataMap map[model.FunctionType]api.FunctionDataCmdInterface
	pendingRequests api.PendingRequestsInterface
	resultHandler   []api.FeatureResultInterface
	resultCallback  map[model.MsgCounterType]func(result api.ResultMessage)

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
		entity:          entity,
		functionDataMap: make(map[model.FunctionType]api.FunctionDataCmdInterface),
		pendingRequests: NewPendingRequest(),
		resultCallback:  make(map[model.MsgCounterType]func(result api.ResultMessage)),
	}

	for _, fd := range CreateFunctionData[api.FunctionDataCmdInterface](ftype) {
		res.functionDataMap[fd.Function()] = fd
	}
	res.operations = make(map[model.FunctionType]api.OperationsInterface)

	return res
}

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
	r.operations[function] = NewOperations(read, write)
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
	r.mux.Lock()

	fd := r.functionData(function)
	if fd == nil {
		return
	}
	fd.UpdateDataAny(data, nil, nil)

	r.mux.Unlock()

	r.Device().NotifySubscribers(r.Address(), fd.NotifyCmdType(nil, nil, false, nil))
}

func (r *FeatureLocal) UpdateData(function model.FunctionType, data any, filterPartial *model.FilterType, filterDelete *model.FilterType) {
	r.mux.Lock()
	defer r.mux.Unlock()

	fctData := r.functionData(function)
	if fctData == nil {
		return
	}

	fctData.UpdateDataAny(data, filterPartial, filterDelete)
}

func (r *FeatureLocal) AddResultHandler(handler api.FeatureResultInterface) {
	r.resultHandler = append(r.resultHandler, handler)
}

func (r *FeatureLocal) AddResultCallback(msgCounterReference model.MsgCounterType, function func(msg api.ResultMessage)) {
	r.muxResultCB.Lock()
	defer r.muxResultCB.Unlock()

	r.resultCallback[msgCounterReference] = function
}

func (r *FeatureLocal) processResultCallbacks(msgCounterReference model.MsgCounterType, msg api.ResultMessage) {
	r.muxResultCB.Lock()
	defer r.muxResultCB.Unlock()

	cb, ok := r.resultCallback[msgCounterReference]
	if !ok {
		return
	}

	go cb(msg)

	delete(r.resultCallback, msgCounterReference)
}

func (r *FeatureLocal) Information() *model.NodeManagementDetailedDiscoveryFeatureInformationType {
	var funs []model.FunctionPropertyType
	for fun, operations := range r.operations {
		var functionType model.FunctionType = model.FunctionType(fun)
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

func (r *FeatureLocal) RequestData(
	function model.FunctionType,
	selector any,
	elements any,
	destination api.FeatureRemoteInterface) (*model.MsgCounterType, *model.ErrorType) {
	fd := r.functionData(function)
	if fd == nil {
		return nil, model.NewErrorTypeFromString("function data not found")
	}

	cmd := fd.ReadCmdType(selector, elements)

	return r.RequestDataBySenderAddress(cmd, destination.Sender(), destination.Device().Ski(), destination.Address(), destination.MaxResponseDelayDuration())
}

func (r *FeatureLocal) RequestDataBySenderAddress(
	cmd model.CmdType,
	sender api.SenderInterface,
	deviceSki string,
	destinationAddress *model.FeatureAddressType,
	maxDelay time.Duration) (*model.MsgCounterType, *model.ErrorType) {

	msgCounter, err := sender.Request(model.CmdClassifierTypeRead, r.Address(), destinationAddress, false, []model.CmdType{cmd})
	if err == nil {
		r.pendingRequests.Add(deviceSki, *msgCounter, maxDelay)
		return msgCounter, nil
	}

	return msgCounter, model.NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
}

// Wait and return the response from destination for a message with the msgCounter ID
// this will block until the response is received
func (r *FeatureLocal) FetchRequestData(
	msgCounter model.MsgCounterType,
	destination api.FeatureRemoteInterface) (any, *model.ErrorType) {

	return r.pendingRequests.GetData(destination.Device().Ski(), msgCounter)
}

// Subscribe to a remote feature
func (r *FeatureLocal) Subscribe(remoteAddress *model.FeatureAddressType) (*model.MsgCounterType, *model.ErrorType) {
	if remoteAddress.Device == nil {
		return nil, model.NewErrorTypeFromString("device not found")
	}
	remoteDevice := r.entity.Device().RemoteDeviceForAddress(*remoteAddress.Device)
	if remoteDevice == nil {
		return nil, model.NewErrorTypeFromString("device not found")
	}

	if r.Role() == model.RoleTypeServer {
		return nil, model.NewErrorTypeFromString(fmt.Sprintf("the server feature '%s' cannot request a subscription", r))
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
func (r *FeatureLocal) RemoveSubscription(remoteAddress *model.FeatureAddressType) {
	remoteDevice := r.entity.Device().RemoteDeviceForAddress(*remoteAddress.Device)
	if remoteDevice == nil {
		return
	}

	if _, err := remoteDevice.Sender().Unsubscribe(r.Address(), remoteAddress); err != nil {
		return
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
}

// Remove all subscriptions to remote features
func (r *FeatureLocal) RemoveAllSubscriptions() {
	for _, item := range r.subscriptions {
		r.RemoveSubscription(item)
	}
}

// Bind to a remote feature
func (r *FeatureLocal) Bind(remoteAddress *model.FeatureAddressType) (*model.MsgCounterType, *model.ErrorType) {
	remoteDevice := r.entity.Device().RemoteDeviceForAddress(*remoteAddress.Device)
	if remoteDevice == nil {
		return nil, model.NewErrorTypeFromString("device not found")
	}

	if r.Role() == model.RoleTypeServer {
		return nil, model.NewErrorTypeFromString(fmt.Sprintf("the server feature '%s' cannot request a binding", r))
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
func (r *FeatureLocal) RemoveBinding(remoteAddress *model.FeatureAddressType) {
	remoteDevice := r.entity.Device().RemoteDeviceForAddress(*remoteAddress.Device)
	if remoteDevice == nil {
		return
	}

	if _, err := remoteDevice.Sender().Unbind(r.Address(), remoteAddress); err != nil {
		return
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
}

// Remove all subscriptions to remote features
func (r *FeatureLocal) RemoveAllBindings() {
	for _, item := range r.bindings {
		r.RemoveBinding(item)
	}
}

// Send a notification message with the current data of function to the destination
func (r *FeatureLocal) NotifyData(
	function model.FunctionType,
	deleteSelector, partialSelector any,
	partialWithoutSelector bool,
	deleteElements any,
	destination api.FeatureRemoteInterface) (*model.MsgCounterType, *model.ErrorType) {
	fd := r.functionData(function)
	if fd == nil {
		return nil, model.NewErrorTypeFromString("function data not found")
	}

	cmd := fd.NotifyCmdType(deleteSelector, partialSelector, partialWithoutSelector, deleteElements)

	msgCounter, err := destination.Sender().Request(model.CmdClassifierTypeRead, r.Address(), destination.Address(), false, []model.CmdType{cmd})
	if err != nil {
		return nil, model.NewErrorTypeFromString(err.Error())
	}
	return msgCounter, nil
}

// Send a write message with provided data of function to the destination
func (r *FeatureLocal) WriteData(
	function model.FunctionType,
	deleteSelector, partialSelector any,
	deleteElements any,
	destination api.FeatureRemoteInterface) (*model.MsgCounterType, *model.ErrorType) {
	fd := r.functionData(function)
	if fd == nil {
		return nil, model.NewErrorTypeFromString("function data not found")
	}
	cmd := fd.WriteCmdType(deleteSelector, partialSelector, deleteElements)

	msgCounter, err := destination.Sender().Write(r.Address(), destination.Address(), cmd)
	if err != nil {
		return nil, model.NewErrorTypeFromString(err.Error())
	}

	return msgCounter, nil
}

func (r *FeatureLocal) HandleMessage(message *api.Message) *model.ErrorType {
	if message.Cmd.ResultData != nil {
		return r.processResult(message)
	}

	cmdData, err := message.Cmd.Data()
	if err != nil {
		return model.NewErrorType(model.ErrorNumberTypeCommandNotSupported, err.Error())
	}
	if cmdData.Function == nil {
		return model.NewErrorType(model.ErrorNumberTypeCommandNotSupported, "No function found for cmd data")
	}

	switch message.CmdClassifier {
	case model.CmdClassifierTypeRead:
		if err := r.processRead(*cmdData.Function, message.RequestHeader, message.FeatureRemote); err != nil {
			return err
		}
	case model.CmdClassifierTypeReply:
		if err := r.processReply(*cmdData.Function, cmdData.Value, message.FilterPartial, message.FilterDelete, message.RequestHeader, message.FeatureRemote); err != nil {
			return err
		}
	case model.CmdClassifierTypeNotify:
		if err := r.processNotify(*cmdData.Function, cmdData.Value, message.FilterPartial, message.FilterDelete, message.FeatureRemote); err != nil {
			return err
		}
	case model.CmdClassifierTypeWrite:
		if err := r.processWrite(*cmdData.Function, cmdData.Value, message.FilterPartial, message.FilterDelete, message.FeatureRemote); err != nil {
			return err
		}
	default:
		return model.NewErrorTypeFromString(fmt.Sprintf("CmdClassifier not implemented: %s", message.CmdClassifier))
	}

	return nil
}

func (r *FeatureLocal) processResult(message *api.Message) *model.ErrorType {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeResult:
		if *message.Cmd.ResultData.ErrorNumber != model.ErrorNumberTypeNoError {
			// error numbers explained in Resource Spec 3.11
			errorString := fmt.Sprintf("Error Result received %d", *message.Cmd.ResultData.ErrorNumber)
			if message.Cmd.ResultData.Description != nil {
				errorString += fmt.Sprintf(": %s", *message.Cmd.ResultData.Description)
			}
			logging.Log().Debug(errorString)
		}

		// we don't need to populate this error as requests don't require a pendingRequest entry
		_ = r.pendingRequests.SetResult(message.DeviceRemote.Ski(), *message.RequestHeader.MsgCounterReference, model.NewErrorTypeFromResult(message.Cmd.ResultData))

		if message.RequestHeader.MsgCounterReference == nil {
			return nil
		}

		// call the Features Error Handler
		errorMsg := api.ResultMessage{
			MsgCounterReference: *message.RequestHeader.MsgCounterReference,
			Result:              message.Cmd.ResultData,
			FeatureLocal:        r,
			FeatureRemote:       message.FeatureRemote,
			DeviceRemote:        message.DeviceRemote,
		}

		if r.resultHandler != nil {
			for _, item := range r.resultHandler {
				go item.HandleResult(errorMsg)
			}
		}

		r.processResultCallbacks(*message.RequestHeader.MsgCounterReference, errorMsg)

		return nil

	default:
		return model.NewErrorType(
			model.ErrorNumberTypeGeneralError,
			fmt.Sprintf("ResultData CmdClassifierType %s not implemented", message.CmdClassifier))
	}
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
	if err := featureRemote.Sender().Reply(requestHeader, r.Address(), cmd); err != nil {
		return model.NewErrorTypeFromString(err.Error())
	}

	return nil
}

func (r *FeatureLocal) processReply(function model.FunctionType, data any, filterPartial *model.FilterType, filterDelete *model.FilterType, requestHeader *model.HeaderType, featureRemote api.FeatureRemoteInterface) *model.ErrorType {
	featureRemote.UpdateData(function, data, filterPartial, filterDelete)

	if requestHeader != nil && requestHeader.MsgCounterReference != nil {
		_ = r.pendingRequests.SetData(featureRemote.Device().Ski(), *requestHeader.MsgCounterReference, data)
	}

	// an error in SetData only means that there is no pendingRequest waiting for this dataset
	// so this is nothing to consider as an error to return

	// the data was updated, so send an event, other event handlers may watch out for this as well
	payload := api.EventPayload{
		Ski:           featureRemote.Device().Ski(),
		EventType:     api.EventTypeDataChange,
		ChangeType:    api.ElementChangeUpdate,
		Feature:       featureRemote,
		Device:        featureRemote.Device(),
		Entity:        featureRemote.Entity(),
		CmdClassifier: util.Ptr(model.CmdClassifierTypeReply),
		Data:          data,
	}
	Events.Publish(payload)

	return nil
}

func (r *FeatureLocal) processNotify(function model.FunctionType, data any, filterPartial *model.FilterType, filterDelete *model.FilterType, featureRemote api.FeatureRemoteInterface) *model.ErrorType {
	featureRemote.UpdateData(function, data, filterPartial, filterDelete)

	payload := api.EventPayload{
		Ski:           featureRemote.Device().Ski(),
		EventType:     api.EventTypeDataChange,
		ChangeType:    api.ElementChangeUpdate,
		Feature:       featureRemote,
		Device:        featureRemote.Device(),
		Entity:        featureRemote.Entity(),
		CmdClassifier: util.Ptr(model.CmdClassifierTypeNotify),
		Data:          data,
	}
	Events.Publish(payload)

	return nil
}

func (r *FeatureLocal) processWrite(function model.FunctionType, data any, filterPartial *model.FilterType, filterDelete *model.FilterType, featureRemote api.FeatureRemoteInterface) *model.ErrorType {
	r.UpdateData(function, data, filterPartial, filterDelete)

	payload := api.EventPayload{
		Ski:           featureRemote.Device().Ski(),
		EventType:     api.EventTypeDataChange,
		ChangeType:    api.ElementChangeUpdate,
		Feature:       featureRemote,
		Device:        featureRemote.Device(),
		Entity:        featureRemote.Entity(),
		LocalFeature:  r,
		Function:      function,
		CmdClassifier: util.Ptr(model.CmdClassifierTypeWrite),
		Data:          data,
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
