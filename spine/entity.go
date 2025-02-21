package spine

import (
	"encoding/json"
	"sync"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

const DeviceInformationEntityId uint = 0

var DeviceInformationAddressEntity = []model.AddressEntityType{model.AddressEntityType(DeviceInformationEntityId)}

type Entity struct {
	eType        model.EntityTypeType
	address      *model.EntityAddressType
	description  *model.DescriptionType
	fIdGenerator func() uint

	muxGenerator sync.Mutex
}

var _ api.EntityInterface = (*Entity)(nil)

func NewEntity(eType model.EntityTypeType, deviceAddress *model.AddressDeviceType, entityAddress []model.AddressEntityType) *Entity {
	entity := &Entity{
		eType: eType,
		address: &model.EntityAddressType{
			Device: deviceAddress,
			Entity: entityAddress,
		},
	}
	if entityAddress != nil && entityAddress[0] == 0 {
		// Entity 0 Feature addresses start with 0
		entity.fIdGenerator = newFeatureIdGenerator(0)
	} else {
		// Entity 1 and further Feature addresses start with 1
		entity.fIdGenerator = newFeatureIdGenerator(1)
	}

	return entity
}

func (r *Entity) Address() *model.EntityAddressType {
	return r.address
}

// Add support for JSON Marshalling
//
// Instances of EntityInterface are used as arguments and return values in various API calls,
// therefor it is helpfull to be able to marshal them to JSON and thus make the API calls
// usable with various communication interfaces
func (r *Entity) MarshalJSON() ([]byte, error) {
	// we do not want to omit address fields, if they are nil
	// and field names should not be lowercased
	type tempAddressType struct {
		Device model.AddressDeviceType
		Entity []model.AddressEntityType
	}
	var tempAddress tempAddressType

	if r.address.Device != nil {
		tempAddress.Device = *r.address.Device
	}
	tempAddress.Entity = r.address.Entity

	bytes, err := json.Marshal(tempAddress)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (r *Entity) EntityType() model.EntityTypeType {
	return r.eType
}

func (r *Entity) Description() *model.DescriptionType {
	return r.description
}

func (r *Entity) SetDescription(d *model.DescriptionType) {
	r.description = d
}

func (r *Entity) NextFeatureId() uint {
	r.muxGenerator.Lock()
	defer r.muxGenerator.Unlock()

	return r.fIdGenerator()
}

func EntityAddressType(deviceAddress *model.AddressDeviceType, entityAddress []model.AddressEntityType) *model.EntityAddressType {
	return &model.EntityAddressType{
		Device: deviceAddress,
		Entity: entityAddress,
	}
}

func NewEntityAddressType(deviceName string, entityIds []uint) *model.EntityAddressType {
	return &model.EntityAddressType{
		Device: util.Ptr(model.AddressDeviceType(deviceName)),
		Entity: NewAddressEntityType(entityIds),
	}
}

func NewAddressEntityType(entityIds []uint) []model.AddressEntityType {
	var addressEntity []model.AddressEntityType
	for _, item := range entityIds {
		addressEntity = append(addressEntity, model.AddressEntityType(item))
	}
	return addressEntity
}

func newFeatureIdGenerator(id uint) func() uint {
	return func() uint {
		defer func() { id += 1 }()
		return id
	}
}
