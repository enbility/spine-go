package spine

import (
	"fmt"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

type Feature struct {
	address     *model.FeatureAddressType
	ftype       model.FeatureTypeType
	description *model.DescriptionType
	role        model.RoleType
	operations  map[model.FunctionType]api.OperationsInterface
}

var _ api.FeatureInterface = (*Feature)(nil)

func NewFeature(address *model.FeatureAddressType, ftype model.FeatureTypeType, role model.RoleType) *Feature {
	res := &Feature{
		address: address,
		ftype:   ftype,
		role:    role,
	}

	return res
}

func (r *Feature) Address() *model.FeatureAddressType {
	return r.address
}

func (r *Feature) Type() model.FeatureTypeType {
	return r.ftype
}

func (r *Feature) Role() model.RoleType {
	return r.role
}

func (r *Feature) Operations() map[model.FunctionType]api.OperationsInterface {
	return r.operations
}

func (r *Feature) Description() *model.DescriptionType {
	return r.description
}

func (r *Feature) SetDescription(d *model.DescriptionType) {
	r.description = d
}

func (r *Feature) SetDescriptionString(s string) {
	r.description = util.Ptr(model.DescriptionType(s))
}

func (r *Feature) String() string {
	if r == nil {
		return ""
	}
	return fmt.Sprintf("Id: %d (%s)", *r.Address().Feature, string(r.ftype))
}

func featureAddressType(id uint, entityAddress *model.EntityAddressType) *model.FeatureAddressType {
	res := model.FeatureAddressType{
		Device:  entityAddress.Device,
		Entity:  entityAddress.Entity,
		Feature: util.Ptr(model.AddressFeatureType(id)),
	}

	return &res
}
