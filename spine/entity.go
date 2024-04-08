package spine

import (
	"sync"

	"github.com/ahmetb/go-linq/v3"
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
	if entityAddress[0] == 0 {
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
	linq.From(entityIds).SelectT(func(i uint) model.AddressEntityType { return model.AddressEntityType(i) }).ToSlice(&addressEntity)
	return addressEntity
}

func newFeatureIdGenerator(id uint) func() uint {
	return func() uint {
		defer func() { id += 1 }()
		return id
	}
}
