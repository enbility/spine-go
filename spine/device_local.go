package spine

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

var _ api.DeviceLocal = (*DeviceLocalImpl)(nil)

type DeviceLocalImpl struct {
	*DeviceImpl
	entities            []api.EntityLocal
	subscriptionManager api.SubscriptionManager
	bindingManager      api.BindingManager
	heartbeatManager    api.HeartbeatManager
	nodeManagement      *NodeManagementImpl

	remoteDevices map[string]api.DeviceRemote

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
func NewDeviceLocalImpl(
	brandName, deviceModel, serialNumber, deviceCode, deviceAddress string,
	deviceType model.DeviceTypeType,
	featureSet model.NetworkManagementFeatureSetType,
	heartbeatTimeout time.Duration) *DeviceLocalImpl {
	address := model.AddressDeviceType(deviceAddress)

	var fSet *model.NetworkManagementFeatureSetType
	if len(featureSet) != 0 {
		fSet = &featureSet
	}

	res := &DeviceLocalImpl{
		DeviceImpl:    NewDeviceImpl(&address, &deviceType, fSet),
		remoteDevices: make(map[string]api.DeviceRemote),
		brandName:     brandName,
		deviceModel:   deviceModel,
		serialNumber:  serialNumber,
		deviceCode:    deviceCode,
	}

	res.subscriptionManager = NewSubscriptionManager(res)
	res.bindingManager = NewBindingManager(res)
	res.heartbeatManager = NewHeartbeatManager(res, res.subscriptionManager, heartbeatTimeout)

	res.addDeviceInformation()
	return res
}

func (r *DeviceLocalImpl) RemoveRemoteDeviceConnection(ski string) {
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

// Helper method used by tests and AddRemoteDevice
func (r *DeviceLocalImpl) AddRemoteDeviceForSki(ski string, rDevice api.DeviceRemote) {
	r.mux.Lock()
	r.remoteDevices[ski] = rDevice
	r.mux.Unlock()
}

// Setup a new remote device with a given SKI and triggers SPINE requesting device details
func (r *DeviceLocalImpl) SetupRemoteDevice(ski string, writeI shipapi.SpineDataConnection) shipapi.SpineDataProcessing {
	sender := NewSender(writeI)
	rDevice := NewDeviceRemoteImpl(r, ski, sender)

	r.AddRemoteDeviceForSki(ski, rDevice)

	// Request Detailed Discovery Data
	_, _ = r.nodeManagement.RequestDetailedDiscovery(rDevice.ski, rDevice.address, rDevice.sender)

	// TODO: Add error handling
	// If the request returned an error, it should be retried until it does not

	// always add subscription, as it checks if it already exists
	_ = Events.subscribe(api.EventHandlerLevelCore, r)

	return rDevice
}

// React to some specific events
func (r *DeviceLocalImpl) HandleEvent(payload api.EventPayload) {
	// Subscribe to NodeManagment after DetailedDiscovery is received
	if payload.EventType != api.EventTypeDeviceChange || payload.ChangeType != api.ElementChangeAdd {
		return
	}

	if payload.Data == nil {
		return
	}

	if len(payload.Ski) == 0 {
		return
	}

	if r.RemoteDeviceForSki(payload.Ski) == nil {
		return
	}

	// the codefactor warning is invalid, as .(type) check can not be replaced with if then
	//revive:disable-next-line
	switch payload.Data.(type) {
	case *model.NodeManagementDetailedDiscoveryDataType:
		_, _ = r.nodeManagement.Subscribe(payload.Feature.Address())

		// Request Use Case Data
		_, _ = r.nodeManagement.RequestUseCaseData(payload.Device.Ski(), payload.Device.Address(), payload.Device.Sender())
	}
}

func (r *DeviceLocalImpl) RemoveRemoteDevice(ski string) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.remoteDevices[ski] == nil {
		return
	}

	// remove all subscriptions for this device
	subscriptionMgr := r.SubscriptionManager()
	subscriptionMgr.RemoveSubscriptionsForDevice(r.remoteDevices[ski])

	// make sure Heartbeat Manager is up to date
	r.HeartbeatManager().UpdateHeartbeatOnSubscriptions()

	// remove all bindings for this device
	bindingMgr := r.BindingManager()
	bindingMgr.RemoveBindingsForDevice(r.remoteDevices[ski])

	delete(r.remoteDevices, ski)

	// only unsubscribe if we don't have any remote devices left
	if len(r.remoteDevices) == 0 {
		_ = Events.unsubscribe(api.EventHandlerLevelCore, r)
	}
}

func (r *DeviceLocalImpl) RemoteDevices() []api.DeviceRemote {
	r.mux.Lock()
	defer r.mux.Unlock()

	res := make([]api.DeviceRemote, 0)
	for _, rDevice := range r.remoteDevices {
		res = append(res, rDevice)
	}

	return res
}

func (r *DeviceLocalImpl) RemoteDeviceForAddress(address model.AddressDeviceType) api.DeviceRemote {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, item := range r.remoteDevices {
		if *item.Address() == address {
			return item
		}
	}

	return nil
}

func (r *DeviceLocalImpl) RemoteDeviceForSki(ski string) api.DeviceRemote {
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.remoteDevices[ski]
}

func (r *DeviceLocalImpl) ProcessCmd(datagram model.DatagramType, remoteDevice api.DeviceRemote) error {
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
		CmdClassifier: *cmdClassifier,
		Cmd:           cmd,
		FilterPartial: filterPartial,
		FilterDelete:  filterDelete,
		FeatureRemote: remoteFeature,
		EntityRemote:  remoteEntity,
		DeviceRemote:  remoteDevice,
	}

	ackRequest := datagram.Header.AckRequest

	if localFeature == nil {
		errorMessage := "invalid feature address"
		_ = remoteFeature.Sender().ResultError(message.RequestHeader, destAddr, model.NewErrorType(model.ErrorNumberTypeDestinationUnknown, errorMessage))

		return errors.New(errorMessage)
	}

	lfType := string(localFeature.Type())
	rfType := ""
	if remoteFeature != nil {
		rfType = string(remoteFeature.Type())
	}

	logging.Log().Debug(datagram.PrintMessageOverview(false, lfType, rfType))

	// check if this is a write with an existing binding and if write is allowed on this feature
	if message.CmdClassifier == model.CmdClassifierTypeWrite {
		cmdData, err := cmd.Data()
		if err != nil || cmdData.Function == nil {
			err := model.NewErrorTypeFromString("no function found for cmd data")
			_ = remoteFeature.Sender().ResultError(message.RequestHeader, localFeature.Address(), err)
			return errors.New(err.String())
		}

		if operations, ok := localFeature.Operations()[*cmdData.Function]; !ok || !operations.Write() {
			err := model.NewErrorTypeFromString("write is not allowed on this function")
			_ = remoteFeature.Sender().ResultError(message.RequestHeader, localFeature.Address(), err)
			return errors.New(err.String())
		}

		if remoteFeature == nil ||
			!r.BindingManager().HasLocalFeatureRemoteBinding(localFeature.Address(), remoteFeature.Address()) {
			err := model.NewErrorTypeFromString("write denied due to missing binding")
			_ = remoteFeature.Sender().ResultError(message.RequestHeader, localFeature.Address(), err)
			return errors.New(err.String())
		}

	}

	err := localFeature.HandleMessage(message)
	if err != nil {
		// TODO: add error description in a useful format

		// Don't send error responses for incoming resulterror messages
		if message.CmdClassifier != model.CmdClassifierTypeResult {
			_ = remoteFeature.Sender().ResultError(message.RequestHeader, localFeature.Address(), err)
		}

		return errors.New(err.String())
	}
	if ackRequest != nil && *ackRequest {
		_ = remoteFeature.Sender().ResultSuccess(message.RequestHeader, localFeature.Address())
	}

	return nil
}

func (r *DeviceLocalImpl) NodeManagement() api.NodeManagement {
	return r.nodeManagement
}

func (r *DeviceLocalImpl) SubscriptionManager() api.SubscriptionManager {
	return r.subscriptionManager
}

func (r *DeviceLocalImpl) BindingManager() api.BindingManager {
	return r.bindingManager
}

func (r *DeviceLocalImpl) HeartbeatManager() api.HeartbeatManager {
	return r.heartbeatManager
}

func (r *DeviceLocalImpl) AddEntity(entity api.EntityLocal) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.entities = append(r.entities, entity)

	r.notifySubscribersOfEntity(entity, model.NetworkManagementStateChangeTypeAdded)
}

func (r *DeviceLocalImpl) RemoveEntity(entity api.EntityLocal) {
	entity.RemoveAllUseCaseSupports()
	entity.RemoveAllSubscriptions()
	entity.RemoveAllBindings()

	r.mux.Lock()
	defer r.mux.Unlock()

	var entities []api.EntityLocal
	for _, e := range r.entities {
		if e != entity {
			entities = append(entities, e)
		}
	}

	r.entities = entities

	r.notifySubscribersOfEntity(entity, model.NetworkManagementStateChangeTypeRemoved)
}

func (r *DeviceLocalImpl) Entities() []api.EntityLocal {
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.entities
}

func (r *DeviceLocalImpl) Entity(id []model.AddressEntityType) api.EntityLocal {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, e := range r.entities {
		if reflect.DeepEqual(id, e.Address().Entity) {
			return e
		}
	}
	return nil
}

func (r *DeviceLocalImpl) EntityForType(entityType model.EntityTypeType) api.EntityLocal {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, e := range r.entities {
		if e.EntityType() == entityType {
			return e
		}
	}
	return nil
}

func (r *DeviceLocalImpl) FeatureByAddress(address *model.FeatureAddressType) api.FeatureLocal {
	entity := r.Entity(address.Entity)
	if entity != nil {
		return entity.Feature(address.Feature)
	}
	return nil
}

func (r *DeviceLocalImpl) Information() *model.NodeManagementDetailedDiscoveryDeviceInformationType {
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

// send a notify message to all remote devices
func (r *DeviceLocalImpl) NotifyUseCaseData() {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, remoteDevice := range r.remoteDevices {
		// TODO: add error management
		_, _ = r.nodeManagement.NotifyUseCaseData(remoteDevice)
	}
}

func (r *DeviceLocalImpl) NotifySubscribers(featureAddress *model.FeatureAddressType, cmd model.CmdType) {
	subscriptions := r.SubscriptionManager().SubscriptionsOnFeature(*featureAddress)
	for _, subscription := range subscriptions {
		// TODO: error handling
		_, _ = subscription.ClientFeature.Sender().Notify(subscription.ServerFeature.Address(), subscription.ClientFeature.Address(), cmd)
	}
}

func (r *DeviceLocalImpl) notifySubscribersOfEntity(entity api.EntityLocal, state model.NetworkManagementStateChangeType) {
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

func (r *DeviceLocalImpl) addDeviceInformation() {
	entityType := model.EntityTypeTypeDeviceInformation
	entity := NewEntityLocalImpl(r, entityType, []model.AddressEntityType{model.AddressEntityType(DeviceInformationEntityId)})

	{
		r.nodeManagement = NewNodeManagementImpl(entity.NextFeatureId(), entity)
		entity.AddFeature(r.nodeManagement)
	}
	{
		f := NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)

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
