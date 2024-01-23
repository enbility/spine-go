package spine

import (
	"encoding/json"
	"errors"
	"reflect"
	"slices"
	"sync"

	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

var _ api.DeviceRemoteInterface = (*DeviceRemote)(nil)

type DeviceRemote struct {
	*Device

	ski string

	entities      []api.EntityRemoteInterface
	entitiesMutex sync.Mutex

	sender api.SenderInterface

	localDevice api.DeviceLocalInterface
}

var _ shipapi.ShipConnectionDataReaderInterface = (*DeviceRemote)(nil)

func NewDeviceRemote(localDevice api.DeviceLocalInterface, ski string, sender api.SenderInterface) *DeviceRemote {
	res := DeviceRemote{
		Device:      NewDevice(nil, nil, nil),
		ski:         ski,
		localDevice: localDevice,
		sender:      sender,
	}
	res.addNodeManagement()

	return &res
}

// return the device SKI
func (d *DeviceRemote) Ski() string {
	return d.ski
}

func (d *DeviceRemote) SetAddress(address *model.AddressDeviceType) {
	d.address = address
}

func (d *DeviceRemote) HandleSpineMesssage(message []byte) (*model.MsgCounterType, error) {
	datagram := model.Datagram{}
	if err := json.Unmarshal([]byte(message), &datagram); err != nil {
		return nil, err
	}
	err := d.localDevice.ProcessCmd(datagram.Datagram, d)
	if err != nil {
		logging.Log().Trace(err)
	}

	return datagram.Datagram.Header.MsgCounter, nil
}

// processing incoming SPINE message from the associated SHIP connection
func (d *DeviceRemote) HandleShipPayloadMessage(message []byte) {
	if _, err := d.HandleSpineMesssage(message); err != nil {
		logging.Log().Errorf("error handling spine message", err)
	}
}

func (d *DeviceRemote) addNodeManagement() {
	deviceInformation := d.addNewEntity(model.EntityTypeTypeDeviceInformation, NewAddressEntityType([]uint{DeviceInformationEntityId}))
	nodeManagement := NewFeatureRemote(deviceInformation.NextFeatureId(), deviceInformation, model.FeatureTypeTypeNodeManagement, model.RoleTypeSpecial)
	deviceInformation.AddFeature(nodeManagement)
}

func (d *DeviceRemote) Sender() api.SenderInterface {
	return d.sender
}

// Return an entity with a given address
func (d *DeviceRemote) Entity(id []model.AddressEntityType) api.EntityRemoteInterface {
	d.entitiesMutex.Lock()
	defer d.entitiesMutex.Unlock()

	for _, e := range d.entities {
		if reflect.DeepEqual(id, e.Address().Entity) {
			return e
		}
	}
	return nil
}

// Return all entities of this device
func (d *DeviceRemote) Entities() []api.EntityRemoteInterface {
	return d.entities
}

// Return the feature for a given address
func (d *DeviceRemote) FeatureByAddress(address *model.FeatureAddressType) api.FeatureRemoteInterface {
	entity := d.Entity(address.Entity)
	if entity != nil {
		return entity.Feature(address.Feature)
	}
	return nil
}

// Remove an entity with a given address from this device
func (d *DeviceRemote) RemoveByAddress(addr []model.AddressEntityType) api.EntityRemoteInterface {
	entityForRemoval := d.Entity(addr)
	if entityForRemoval == nil {
		return nil
	}

	d.entitiesMutex.Lock()
	defer d.entitiesMutex.Unlock()

	var newEntities []api.EntityRemoteInterface
	for _, item := range d.entities {
		if !reflect.DeepEqual(item, entityForRemoval) {
			newEntities = append(newEntities, item)
		}
	}
	d.entities = newEntities

	return entityForRemoval
}

// Get the feature for a given entity, feature type and feature role
func (r *DeviceRemote) FeatureByEntityTypeAndRole(entity api.EntityRemoteInterface, featureType model.FeatureTypeType, role model.RoleType) api.FeatureRemoteInterface {
	if len(r.entities) < 1 {
		return nil
	}

	r.entitiesMutex.Lock()
	defer r.entitiesMutex.Unlock()

	for _, e := range r.entities {
		if entity != e {
			continue
		}
		for _, feature := range entity.Features() {
			if feature.Type() == featureType && feature.Role() == role {
				return feature
			}
		}
	}

	return nil
}

func (d *DeviceRemote) UpdateDevice(description *model.NetworkManagementDeviceDescriptionDataType) {
	if description != nil {
		if description.DeviceAddress != nil && description.DeviceAddress.Device != nil {
			d.address = description.DeviceAddress.Device
		}
		if description.DeviceType != nil {
			d.dType = description.DeviceType
		}
		if description.NetworkFeatureSet != nil {
			d.featureSet = description.NetworkFeatureSet
		}
	}
}

func (d *DeviceRemote) AddEntityAndFeatures(initialData bool, data *model.NodeManagementDetailedDiscoveryDataType) ([]api.EntityRemoteInterface, error) {
	rEntites := make([]api.EntityRemoteInterface, 0)

	for _, ei := range data.EntityInformation {
		if err := d.CheckEntityInformation(initialData, ei); err != nil {
			return nil, err
		}

		entityAddress := ei.Description.EntityAddress.Entity

		entity := d.Entity(entityAddress)
		if entity == nil {
			entity = d.addNewEntity(*ei.Description.EntityType, entityAddress)
			rEntites = append(rEntites, entity)
		}

		entity.SetDescription(ei.Description.Description)
		entity.RemoveAllFeatures()

		for _, fi := range data.FeatureInformation {
			if reflect.DeepEqual(fi.Description.FeatureAddress.Entity, entityAddress) {
				if f, ok := unmarshalFeature(entity, fi); ok {
					entity.AddFeature(f)
				}
			}
		}

		// TOV-TODO: check this approach
		// if err := f.announceFeatureDiscovery(entity); err != nil {
		// 	return err
		// }
	}

	return rEntites, nil
}

// check if the provided entity information is correct
// provide initialData to check if the entity is new and not an update
func (d *DeviceRemote) CheckEntityInformation(initialData bool, entity model.NodeManagementDetailedDiscoveryEntityInformationType) error {
	description := entity.Description
	if description == nil {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: invalid EntityInformation.Description")
	}

	if description.EntityAddress == nil {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: invalid EntityInformation.Description.EntityAddress")
	}

	if description.EntityAddress.Entity == nil {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: invalid EntityInformation.Description.EntityAddress.Entity")
	}

	// Consider on initial NodeManagement Detailed Discovery, the device being empty as it is not yet known
	if initialData {
		return nil
	}

	address := d.Address()
	if description.EntityAddress.Device != nil && address != nil && *description.EntityAddress.Device != *address {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: device address mismatch")
	}

	return nil
}

func (d *DeviceRemote) addNewEntity(eType model.EntityTypeType, address []model.AddressEntityType) api.EntityRemoteInterface {
	newEntity := NewEntityRemote(d, eType, address)
	return d.AddEntity(newEntity)
}

func (d *DeviceRemote) AddEntity(entity api.EntityRemoteInterface) api.EntityRemoteInterface {
	d.entitiesMutex.Lock()
	defer d.entitiesMutex.Unlock()

	d.entities = append(d.entities, entity)

	return entity
}

func (d *DeviceRemote) UseCases() []model.UseCaseInformationDataType {
	entity := d.Entity(DeviceInformationAddressEntity)

	nodemgmt := d.FeatureByEntityTypeAndRole(entity, model.FeatureTypeTypeNodeManagement, model.RoleTypeSpecial)

	data := nodemgmt.DataCopy(model.FunctionTypeNodeManagementUseCaseData).(*model.NodeManagementUseCaseDataType)
	if data != nil {
		return data.UseCaseInformation
	}

	return nil
}

// Checks if the given actor, usecasename and provided server features are available
// Note: the server features are expected to be in a single entity and entity 0 is not checked!
func (d *DeviceRemote) VerifyUseCaseScenariosAndFeaturesSupport(
	usecaseActor model.UseCaseActorType,
	usecaseName model.UseCaseNameType,
	scenarios []model.UseCaseScenarioSupportType,
	serverFeatures []model.FeatureTypeType,
) bool {
	entity := d.Entity(DeviceInformationAddressEntity)

	nodemgmt := d.FeatureByEntityTypeAndRole(entity, model.FeatureTypeTypeNodeManagement, model.RoleTypeSpecial)

	usecases := nodemgmt.DataCopy(model.FunctionTypeNodeManagementUseCaseData).(*model.NodeManagementUseCaseDataType)

	if usecases == nil || len(usecases.UseCaseInformation) == 0 {
		return false
	}

	usecaseAndScenariosFound := false
	for _, usecase := range usecases.UseCaseInformation {
		if usecase.Actor == nil || *usecase.Actor != usecaseActor {
			continue
		}

		for _, support := range usecase.UseCaseSupport {
			if support.UseCaseName == nil || *support.UseCaseName != usecaseName {
				continue
			}

			var foundScenarios []model.UseCaseScenarioSupportType
			for _, scenario := range support.ScenarioSupport {
				if slices.Contains(scenarios, scenario) {
					foundScenarios = append(foundScenarios, scenario)
				}
			}

			if len(foundScenarios) == len(scenarios) {
				usecaseAndScenariosFound = true
				break
			}
		}

		if usecaseAndScenariosFound {
			break
		}
	}

	if !usecaseAndScenariosFound {
		return false
	}

	entities := d.Entities()
	if len(entities) < 2 {
		return false
	}

	entityWithServerFeaturesFound := false

	for index, entity := range entities {
		// ignore NodeManagement entity
		if index == 0 {
			continue
		}

		var foundServerFeatures []model.FeatureTypeType
		for _, feature := range entity.Features() {
			if feature.Role() != model.RoleTypeServer {
				continue
			}

			if slices.Contains(serverFeatures, feature.Type()) {
				foundServerFeatures = append(foundServerFeatures, feature.Type())
			}
		}

		if len(serverFeatures) == len(foundServerFeatures) {
			entityWithServerFeaturesFound = true
			break
		}
	}

	return entityWithServerFeaturesFound
}

func unmarshalFeature(entity api.EntityRemoteInterface,
	featureData model.NodeManagementDetailedDiscoveryFeatureInformationType,
) (api.FeatureRemoteInterface, bool) {
	var result api.FeatureRemoteInterface

	fid := featureData.Description

	if fid == nil {
		return nil, false
	}

	result = NewFeatureRemote(uint(*fid.FeatureAddress.Feature), entity, *fid.FeatureType, *fid.Role)

	result.SetDescription(fid.Description)
	result.SetMaxResponseDelay(fid.MaxResponseDelay)
	result.SetOperations(fid.SupportedFunction)

	return result, true
}
