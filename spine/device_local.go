package spine

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"sync"

	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

type DeviceLocal struct {
	*Device
	entities            []api.EntityLocalInterface
	subscriptionManager api.SubscriptionManagerInterface
	bindingManager      api.BindingManagerInterface
	nodeManagement      *NodeManagement

	remoteDevices map[string]api.DeviceRemoteInterface

	brandName    string
	deviceModel  string
	deviceCode   string
	serialNumber string

	mux sync.Mutex
}

// BrandName is the brand
// DeviceModel is the model
// SerialNumber is the serial number
// DeviceCode is the SHIP id (accessMethods.id)
// DeviceAddress is the SPINE device address
func NewDeviceLocal(
	brandName, deviceModel, serialNumber, deviceCode, deviceAddress string,
	deviceType model.DeviceTypeType,
	featureSet model.NetworkManagementFeatureSetType) *DeviceLocal {
	address := model.AddressDeviceType(deviceAddress)

	var fSet *model.NetworkManagementFeatureSetType
	if len(featureSet) != 0 {
		fSet = &featureSet
	}

	res := &DeviceLocal{
		Device:        NewDevice(&address, &deviceType, fSet),
		remoteDevices: make(map[string]api.DeviceRemoteInterface),
		brandName:     brandName,
		deviceModel:   deviceModel,
		serialNumber:  serialNumber,
		deviceCode:    deviceCode,
	}

	res.subscriptionManager = NewSubscriptionManager(res)
	res.bindingManager = NewBindingManager(res)

	res.addDeviceInformation()
	return res
}

var _ api.EventHandlerInterface = (*DeviceLocal)(nil)

/* EventHandlerInterface */

// React to some specific events
func (r *DeviceLocal) HandleEvent(payload api.EventPayload) {
	// Subscribe to NodeManagement after DetailedDiscovery is received
	if payload.EventType != api.EventTypeDeviceChange || payload.ChangeType != api.ElementChangeAdd {
		return
	}

	if payload.Data == nil {
		return
	}

	if len(payload.Ski) == 0 {
		return
	}

	remoteDevice := r.RemoteDeviceForSki(payload.Ski)
	if remoteDevice == nil {
		return
	}

	// the codefactor warning is invalid, as .(type) check can not be replaced with if then
	//revive:disable-next-line
	switch payload.Data.(type) {
	case *model.NodeManagementDetailedDiscoveryDataType:
		address := payload.Feature.Address()
		if address.Device == nil {
			address.Device = remoteDevice.Address()
		}
		_, _ = r.nodeManagement.SubscribeToRemote(address)

		// Request Use Case Data
		_, _ = r.nodeManagement.RequestUseCaseData(payload.Device.Ski(), remoteDevice.Address(), payload.Device.Sender())
	}
}

var _ api.DeviceLocalInterface = (*DeviceLocal)(nil)

/* DeviceLocalInterface */

// Setup a new remote device with a given SKI and triggers SPINE requesting device details
func (r *DeviceLocal) SetupRemoteDevice(ski string, writeI shipapi.ShipConnectionDataWriterInterface) shipapi.ShipConnectionDataReaderInterface {
	sender := NewSender(writeI)
	rDevice := NewDeviceRemote(r, ski, sender)

	r.AddRemoteDeviceForSki(ski, rDevice)

	// always add subscription, as it checks if it already exists
	_ = Events.subscribe(api.EventHandlerLevelCore, r)

	// Request Detailed Discovery Data
	_, _ = r.RequestRemoteDetailedDiscoveryData(rDevice)

	// TODO: Add error handling
	// If the request returned an error, it should be retried until it does not

	return rDevice
}

func (r *DeviceLocal) RequestRemoteDetailedDiscoveryData(rDevice api.DeviceRemoteInterface) (*model.MsgCounterType, *model.ErrorType) {
	// Request Detailed Discovery Data
	return r.nodeManagement.RequestDetailedDiscovery(rDevice.Ski(), rDevice.Address(), rDevice.Sender())
}

// Helper method used by tests and AddRemoteDevice
func (r *DeviceLocal) AddRemoteDeviceForSki(ski string, rDevice api.DeviceRemoteInterface) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.remoteDevices[ski] = rDevice
}

func (r *DeviceLocal) RemoveRemoteDeviceConnection(ski string) {
	remoteDevice := r.RemoteDeviceForSki(ski)

	r.RemoveRemoteDevice(ski)

	// inform about the disconnection
	payload := api.EventPayload{
		Ski:        ski,
		EventType:  api.EventTypeDeviceChange,
		ChangeType: api.ElementChangeRemove,
		Device:     remoteDevice,
	}
	Events.Publish(payload)
}

func (r *DeviceLocal) RemoveRemoteDevice(ski string) {
	remoteDevice := r.RemoteDeviceForSki(ski)
	if remoteDevice == nil {
		return
	}

	// remove all subscriptions for this device
	subscriptionMgr := r.SubscriptionManager()
	subscriptionMgr.RemoveSubscriptionsForDevice(r.remoteDevices[ski])

	// remove all bindings for this device
	bindingMgr := r.BindingManager()
	bindingMgr.RemoveBindingsForDevice(r.remoteDevices[ski])

	delete(r.remoteDevices, ski)

	// only unsubscribe if we don't have any remote devices left
	if len(r.remoteDevices) == 0 {
		_ = Events.unsubscribe(api.EventHandlerLevelCore, r)
	}

	remoteDeviceAddress := &model.DeviceAddressType{
		Device: remoteDevice.Address(),
	}
	// remove all data caches for this device
	for _, entity := range r.entities {
		for _, feature := range entity.Features() {
			feature.CleanWriteApprovalCaches(ski)
			feature.CleanRemoteDeviceCaches(remoteDeviceAddress)
		}
	}
}

func (r *DeviceLocal) RemoteDevices() []api.DeviceRemoteInterface {
	r.mux.Lock()
	defer r.mux.Unlock()

	res := make([]api.DeviceRemoteInterface, 0)
	for _, rDevice := range r.remoteDevices {
		res = append(res, rDevice)
	}

	return res
}

func (r *DeviceLocal) RemoteDeviceForAddress(address model.AddressDeviceType) api.DeviceRemoteInterface {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, item := range r.remoteDevices {
		if item.Address() != nil && *item.Address() == address {
			return item
		}
	}

	return nil
}

func (r *DeviceLocal) RemoteDeviceForSki(ski string) api.DeviceRemoteInterface {
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.remoteDevices[ski]
}

func (r *DeviceLocal) AddEntity(entity api.EntityLocalInterface) {
	r.mux.Lock()

	r.entities = append(r.entities, entity)

	r.mux.Unlock()

	r.notifySubscribersOfEntity(entity, model.NetworkManagementStateChangeTypeAdded)
}

func (r *DeviceLocal) RemoveEntity(entity api.EntityLocalInterface) {
	entity.RemoveAllUseCaseSupports()
	entity.RemoveAllSubscriptions()
	entity.RemoveAllBindings()

	r.mux.Lock()

	var entities []api.EntityLocalInterface
	for _, e := range r.entities {
		if e != entity {
			entities = append(entities, e)
		}
	}

	r.entities = entities

	r.mux.Unlock()

	r.notifySubscribersOfEntity(entity, model.NetworkManagementStateChangeTypeRemoved)
}

func (r *DeviceLocal) Entities() []api.EntityLocalInterface {
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.entities
}

func (r *DeviceLocal) Entity(id []model.AddressEntityType) api.EntityLocalInterface {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, e := range r.entities {
		if reflect.DeepEqual(id, e.Address().Entity) {
			return e
		}
	}
	return nil
}

func (r *DeviceLocal) EntityForType(entityType model.EntityTypeType) api.EntityLocalInterface {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, e := range r.entities {
		if e.EntityType() == entityType {
			return e
		}
	}
	return nil
}

func (r *DeviceLocal) FeatureByAddress(address *model.FeatureAddressType) api.FeatureLocalInterface {
	entity := r.Entity(address.Entity)
	if entity != nil {
		return entity.FeatureOfAddress(address.Feature)
	}
	return nil
}

func (r *DeviceLocal) CleanRemoteEntityCaches(remoteAddress *model.EntityAddressType) {
	for _, entity := range r.entities {
		for _, feature := range entity.Features() {
			feature.CleanRemoteEntityCaches(remoteAddress)
		}
	}
}

func (r *DeviceLocal) ProcessCmd(datagram model.DatagramType, remoteDevice api.DeviceRemoteInterface) error {
	destAddr := datagram.Header.AddressDestination
	localFeature := r.FeatureByAddress(destAddr)

	cmdClassifier := datagram.Header.CmdClassifier
	if len(datagram.Payload.Cmd) == 0 {
		return errors.New("no payload cmd content available")
	}
	cmd := datagram.Payload.Cmd[0]

	// TODO check if cmd.Function is the same as the provided cmd value
	filterPartial, filterDelete := cmd.ExtractFilter()

	remoteEntity := remoteDevice.Entity(datagram.Header.AddressSource.Entity)
	remoteFeature := remoteDevice.FeatureByAddress(datagram.Header.AddressSource)
	if remoteFeature == nil {
		return fmt.Errorf("invalid remote feature address: '%s'", datagram.Header.AddressSource)
	}

	message := &api.Message{
		RequestHeader: &datagram.Header,
		Cmd:           cmd,
		FilterPartial: filterPartial,
		FilterDelete:  filterDelete,
		FeatureRemote: remoteFeature,
		EntityRemote:  remoteEntity,
		DeviceRemote:  remoteDevice,
	}

	if cmdClassifier != nil {
		message.CmdClassifier = *cmdClassifier
	} else {
		errorMessage := "cmdClassifier may not be empty"

		_ = remoteFeature.Device().Sender().ResultError(message.RequestHeader, destAddr, model.NewErrorType(model.ErrorNumberTypeDestinationUnknown, errorMessage))

		return errors.New(errorMessage)
	}

	if localFeature == nil {
		errorMessage := "invalid feature address"
		_ = remoteFeature.Device().Sender().ResultError(message.RequestHeader, destAddr, model.NewErrorType(model.ErrorNumberTypeDestinationUnknown, errorMessage))

		return errors.New(errorMessage)
	}

	lfType := string(localFeature.Type())
	rfType := string(remoteFeature.Type())

	logging.Log().Debug(datagram.PrintMessageOverview(false, lfType, rfType))

	// check if this is a write with an existing binding and if write is allowed on this feature
	if message.CmdClassifier == model.CmdClassifierTypeWrite {
		cmdData, err := cmd.Data()
		if err != nil || cmdData.Function == nil {
			err := model.NewErrorTypeFromString("no function found for cmd data")
			_ = remoteFeature.Device().Sender().ResultError(message.RequestHeader, localFeature.Address(), err)
			return errors.New(err.String())
		}

		if operations, ok := localFeature.Operations()[*cmdData.Function]; !ok || !operations.Write() {
			err := model.NewErrorTypeFromString("write is not allowed on this function")
			_ = remoteFeature.Device().Sender().ResultError(message.RequestHeader, localFeature.Address(), err)
			return errors.New(err.String())
		}

		if !r.BindingManager().HasLocalFeatureRemoteBinding(localFeature.Address(), remoteFeature.Address()) {
			err := model.NewErrorTypeFromString("write denied due to missing binding")
			_ = remoteFeature.Device().Sender().ResultError(message.RequestHeader, localFeature.Address(), err)
			return errors.New(err.String())
		}
	}

	err := localFeature.HandleMessage(message)
	if err != nil {
		// TODO: add error description in a useful format

		// Don't send error responses for incoming resulterror messages
		if message.CmdClassifier != model.CmdClassifierTypeResult {
			_ = remoteFeature.Device().Sender().ResultError(message.RequestHeader, localFeature.Address(), err)
		}

		// if this is an error for a notify message, automatically trigger a read request for the same
		if message.CmdClassifier == model.CmdClassifierTypeNotify {
			// set the command function to empty

			if cmdData, err := message.Cmd.Data(); err == nil {
				_, _ = localFeature.RequestRemoteData(*cmdData.Function, nil, nil, remoteFeature)
			}
		}

		return errors.New(err.String())
	}

	ackRequest := message.RequestHeader.AckRequest
	ackClassifiers := []model.CmdClassifierType{
		model.CmdClassifierTypeCall,
		model.CmdClassifierTypeReply,
		model.CmdClassifierTypeNotify}

	if ackRequest != nil && *ackRequest && slices.Contains(ackClassifiers, message.CmdClassifier) {
		// return success as defined in SPINE chapter 5.2.4
		_ = remoteFeature.Device().Sender().ResultSuccess(message.RequestHeader, localFeature.Address())
	}

	return nil
}

func (r *DeviceLocal) NodeManagement() api.NodeManagementInterface {
	return r.nodeManagement
}

func (r *DeviceLocal) SubscriptionManager() api.SubscriptionManagerInterface {
	return r.subscriptionManager
}

func (r *DeviceLocal) BindingManager() api.BindingManagerInterface {
	return r.bindingManager
}

func (r *DeviceLocal) Information() *model.NodeManagementDetailedDiscoveryDeviceInformationType {
	res := model.NodeManagementDetailedDiscoveryDeviceInformationType{
		Description: &model.NetworkManagementDeviceDescriptionDataType{
			DeviceAddress: &model.DeviceAddressType{
				Device: r.address,
			},
			DeviceType:        r.dType,
			NetworkFeatureSet: r.featureSet,
		},
	}
	return &res
}

func (r *DeviceLocal) NotifySubscribers(featureAddress *model.FeatureAddressType, cmd model.CmdType) {
	subscriptions := r.SubscriptionManager().SubscriptionsOnFeature(*featureAddress)
	for _, subscription := range subscriptions {
		// TODO: error handling
		_, _ = subscription.ClientFeature.Device().Sender().Notify(subscription.ServerFeature.Address(), subscription.ClientFeature.Address(), cmd)
	}
}

func (r *DeviceLocal) notifySubscribersOfEntity(entity api.EntityLocalInterface, state model.NetworkManagementStateChangeType) {
	deviceInformation := r.Information()
	entityInformation := *entity.Information()
	entityInformation.Description.LastStateChange = &state

	var featureInformation []model.NodeManagementDetailedDiscoveryFeatureInformationType
	if state == model.NetworkManagementStateChangeTypeAdded {
		for _, f := range entity.Features() {
			featureInformation = append(featureInformation, *f.Information())
		}
	}

	cmd := model.CmdType{
		Function: util.Ptr(model.FunctionTypeNodeManagementDetailedDiscoveryData),
		Filter:   filterEmptyPartial(),
		NodeManagementDetailedDiscoveryData: &model.NodeManagementDetailedDiscoveryDataType{
			SpecificationVersionList: &model.NodeManagementSpecificationVersionListType{
				SpecificationVersion: []model.SpecificationVersionDataType{model.SpecificationVersionDataType(SpecificationVersion)},
			},
			DeviceInformation:  deviceInformation,
			EntityInformation:  []model.NodeManagementDetailedDiscoveryEntityInformationType{entityInformation},
			FeatureInformation: featureInformation,
		},
	}

	r.NotifySubscribers(r.nodeManagement.Address(), cmd)
}

func (r *DeviceLocal) addDeviceInformation() {
	entityType := model.EntityTypeTypeDeviceInformation
	entity := NewEntityLocal(r, entityType, []model.AddressEntityType{model.AddressEntityType(DeviceInformationEntityId)}, 0)

	{
		r.nodeManagement = NewNodeManagement(entity.NextFeatureId(), entity)
		entity.AddFeature(r.nodeManagement)
	}
	{
		f := NewFeatureLocal(entity.NextFeatureId(), entity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)

		f.AddFunctionType(model.FunctionTypeDeviceClassificationManufacturerData, true, false)

		manufacturerData := &model.DeviceClassificationManufacturerDataType{
			BrandName:    util.Ptr(model.DeviceClassificationStringType(r.brandName)),
			VendorName:   util.Ptr(model.DeviceClassificationStringType(r.brandName)),
			DeviceName:   util.Ptr(model.DeviceClassificationStringType(r.deviceModel)),
			DeviceCode:   util.Ptr(model.DeviceClassificationStringType(r.deviceCode)),
			SerialNumber: util.Ptr(model.DeviceClassificationStringType(r.serialNumber)),
		}
		f.SetData(model.FunctionTypeDeviceClassificationManufacturerData, manufacturerData)

		entity.AddFeature(f)
	}

	r.entities = append(r.entities, entity)
}
