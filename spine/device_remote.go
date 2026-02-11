package spine

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"sync"

	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type DeviceRemote struct {
	*Device

	ski string

	entities      []api.EntityRemoteInterface
	entitiesMutex sync.Mutex

	sender api.SenderInterface

	localDevice api.DeviceLocalInterface
}

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

func (d *DeviceRemote) addNodeManagement() {
	deviceInformation := d.addNewEntity(model.EntityTypeTypeDeviceInformation, NewAddressEntityType([]uint{DeviceInformationEntityId}))
	nodeManagement := NewFeatureRemote(deviceInformation.NextFeatureId(), deviceInformation, model.FeatureTypeTypeNodeManagement, model.RoleTypeSpecial)
	deviceInformation.AddFeature(nodeManagement)
}

var _ shipapi.ShipConnectionDataReaderInterface = (*DeviceRemote)(nil)

/* ShipConnectionDataReaderInterface */

// processing incoming SPINE message from the associated SHIP connection
func (d *DeviceRemote) HandleShipPayloadMessage(message []byte) {
	if _, err := d.HandleSpineMesssage(message); err != nil {
		logging.Log().Errorf("error handling spine message", err)
	}
}

var _ api.DeviceRemoteInterface = (*DeviceRemote)(nil)

/* DeviceRemoteInterface */

// return the device SKI
func (d *DeviceRemote) Ski() string {
	return d.ski
}

func (d *DeviceRemote) AddEntity(entity api.EntityRemoteInterface) {
	d.entitiesMutex.Lock()
	defer d.entitiesMutex.Unlock()

	d.entities = append(d.entities, entity)
}

func (d *DeviceRemote) addNewEntity(eType model.EntityTypeType, address []model.AddressEntityType) api.EntityRemoteInterface {
	newEntity := NewEntityRemote(d, eType, address)
	d.AddEntity(newEntity)
	return newEntity
}

// Remove an entity with a given address from this device
func (d *DeviceRemote) RemoveEntityByAddress(addr []model.AddressEntityType) api.EntityRemoteInterface {
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
	d.entitiesMutex.Lock()
	defer d.entitiesMutex.Unlock()

	return d.entities
}

// Return the feature for a given address
func (d *DeviceRemote) FeatureByAddress(address *model.FeatureAddressType) api.FeatureRemoteInterface {
	entity := d.Entity(address.Entity)
	if entity != nil {
		return entity.FeatureOfAddress(address.Feature)
	}
	return nil
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

func (d *DeviceRemote) HandleSpineMesssage(message []byte) (*model.MsgCounterType, error) {
	fixedMessage := fixupSliceFields(message)

	datagram := model.Datagram{}
	if err := json.Unmarshal([]byte(fixedMessage), &datagram); err != nil {
		return nil, err
	}

	if datagram.Datagram.Header.MsgCounterReference != nil {
		d.sender.ProcessResponseForMsgCounterReference(datagram.Datagram.Header.MsgCounterReference)
	}

	err := d.localDevice.ProcessCmd(datagram.Datagram, d)
	if err != nil {
		logging.Log().Trace(err)
	}

	return datagram.Datagram.Header.MsgCounter, nil
}

func (d *DeviceRemote) Sender() api.SenderInterface {
	return d.sender
}

func (d *DeviceRemote) UseCases() []model.UseCaseInformationDataType {
	entity := d.Entity(DeviceInformationAddressEntity)

	nodemgmt := d.FeatureByEntityTypeAndRole(entity, model.FeatureTypeTypeNodeManagement, model.RoleTypeSpecial)

	data, ok := nodemgmt.DataCopy(model.FunctionTypeNodeManagementUseCaseData).(*model.NodeManagementUseCaseDataType)
	if ok && data != nil {
		return data.UseCaseInformation
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

func (d *DeviceRemote) AddEntityAndFeatures(
	initialData bool,
	data *model.NodeManagementDetailedDiscoveryDataType,
	entityAddressToAdd *model.EntityAddressType,
) ([]api.EntityRemoteInterface, error) {
	rEntites := make([]api.EntityRemoteInterface, 0)

	for _, ei := range data.EntityInformation {
		if err := d.CheckEntityInformation(initialData, ei); err != nil {
			return nil, err
		}

		entityAddress := ei.Description.EntityAddress.Entity
		// if entityAddressToAdd, make sure we are adding the correct entity
		if entityAddressToAdd != nil && !reflect.DeepEqual(entityAddress, entityAddressToAdd.Entity) {
			continue
		}

		entity := d.Entity(entityAddress)
		if entity == nil {
			entity = d.addNewEntity(*ei.Description.EntityType, entityAddress)
			rEntites = append(rEntites, entity)
		}

		// make sure the device address is set, which is not on entity 0 on startup !
		if entity.Address().Device == nil || len(*entity.Address().Device) == 0 {
			if data.DeviceInformation != nil &&
				data.DeviceInformation.Description != nil &&
				data.DeviceInformation.Description.DeviceAddress != nil &&
				data.DeviceInformation.Description.DeviceAddress.Device != nil {
				entity.UpdateDeviceAddress(*data.DeviceInformation.Description.DeviceAddress.Device)
			}
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

// fixupSliceFields walks the JSON structure and converts {} back to [] for fields
// that are defined as slices in the spine-go model.
func fixupSliceFields(jsonData []byte) []byte {
	// Quick check: if there's no empty object "{}" that could be a wrongly-converted
	// slice, skip the expensive reflection walk entirely.
	// Note: This may trigger on "{}" inside strings, but that's harmless - the actual
	// fix logic only converts empty maps that are values of slice-typed fields.
	if !bytes.Contains(jsonData, []byte("{}")) {
		return jsonData
	}

	// Parse into generic structure
	var generic interface{}
	if err := json.Unmarshal(jsonData, &generic); err != nil {
		// If parsing fails, return as-is
		return jsonData
	}

	// Get the type of model.Datagram for schema reference
	datagramType := reflect.TypeOf(model.Datagram{})

	// Walk and fix the structure
	fixed := fixupSliceFieldsRecursive(generic, datagramType)

	// Re-marshal
	result, err := json.Marshal(fixed)
	if err != nil {
		return jsonData
	}

	return result
}

// fixupSliceFieldsRecursive recursively walks the JSON structure and fixes slice fields.
// modelType is the expected Go type for this level of the structure.
func fixupSliceFieldsRecursive(v interface{}, modelType reflect.Type) interface{} {
	// Dereference pointer types
	for modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	switch val := v.(type) {
	case map[string]interface{}:
		// For struct types, check each field against the model
		if modelType.Kind() == reflect.Struct {
			result := make(map[string]interface{})
			for key, value := range val {
				// Find the field in the model type by JSON tag
				fieldType := findFieldTypeByJSONTag(modelType, key)
				if fieldType != nil {
					// Check if this field is a slice and the value is an empty map
					actualFieldType := *fieldType
					for actualFieldType.Kind() == reflect.Ptr {
						actualFieldType = actualFieldType.Elem()
					}

					if actualFieldType.Kind() == reflect.Slice {
						// This is a slice field
						if emptyMap, ok := value.(map[string]interface{}); ok && len(emptyMap) == 0 {
							// Empty map {} should be empty slice []
							result[key] = []interface{}{}
							continue
						}
					}

					// Recurse with the field's type
					result[key] = fixupSliceFieldsRecursive(value, *fieldType)
				} else {
					// Field not found in model, keep as-is but still recurse
					result[key] = fixupSliceFieldsRecursive(value, reflect.TypeOf((*interface{})(nil)).Elem())
				}
			}
			return result
		}

		// For non-struct types (like interface{}), just recurse on values
		result := make(map[string]interface{})
		for key, value := range val {
			result[key] = fixupSliceFieldsRecursive(value, reflect.TypeOf((*interface{})(nil)).Elem())
		}
		return result

	case []interface{}:
		// For arrays, get the element type and recurse
		var elemType reflect.Type
		if modelType.Kind() == reflect.Slice {
			elemType = modelType.Elem()
		} else {
			elemType = reflect.TypeOf((*interface{})(nil)).Elem()
		}

		result := make([]interface{}, len(val))
		for i, elem := range val {
			result[i] = fixupSliceFieldsRecursive(elem, elemType)
		}
		return result

	default:
		// Primitive value, return as-is
		return val
	}
}

// findFieldTypeByJSONTag finds a struct field by its JSON tag name and returns its type.
func findFieldTypeByJSONTag(structType reflect.Type, jsonName string) *reflect.Type {
	for structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	if structType.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" {
			continue
		}

		// JSON tag format: "fieldName,omitempty"
		tagName := strings.Split(tag, ",")[0]
		if tagName == jsonName {
			return &field.Type
		}
	}

	return nil
}
