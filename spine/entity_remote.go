package spine

import (
	"sync"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type EntityRemote struct {
	*Entity
	device   api.DeviceRemoteInterface
	features []api.FeatureRemoteInterface

	mux sync.Mutex
}

func NewEntityRemote(device api.DeviceRemoteInterface, eType model.EntityTypeType, entityAddress []model.AddressEntityType) *EntityRemote {
	return &EntityRemote{
		Entity: NewEntity(eType, device.Address(), entityAddress),
		device: device,
	}
}

var _ api.EntityRemoteInterface = (*EntityRemote)(nil)

/* EntityRemoteInterface */

func (r *EntityRemote) Device() api.DeviceRemoteInterface {
	return r.device
}

func (r *EntityRemote) UpdateDeviceAddress(address model.AddressDeviceType) {
	r.address.Device = &address
}

func (r *EntityRemote) AddFeature(f api.FeatureRemoteInterface) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.features = append(r.features, f)
}

func (r *EntityRemote) RemoveAllFeatures() {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.features = nil
}

func (r *EntityRemote) FeatureOfTypeAndRole(featureType model.FeatureTypeType, role model.RoleType) api.FeatureRemoteInterface {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, f := range r.features {
		if f.Type() == featureType && f.Role() == role {
			return f
		}
	}

	return nil
}

func (r *EntityRemote) FeatureOfAddress(addressFeature *model.AddressFeatureType) api.FeatureRemoteInterface {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, f := range r.features {
		if addressFeature != nil && f.Address() != nil &&
			*f.Address().Feature == *addressFeature {
			return f
		}
	}

	return nil
}

func (r *EntityRemote) Features() []api.FeatureRemoteInterface {
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.features
}
