package spine

import (
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

var _ api.EntityRemoteInterface = (*EntityRemote)(nil)

type EntityRemote struct {
	*Entity
	device   api.DeviceRemoteInterface
	features []api.FeatureRemoteInterface
}

func NewEntityRemote(device api.DeviceRemoteInterface, eType model.EntityTypeType, entityAddress []model.AddressEntityType) *EntityRemote {
	return &EntityRemote{
		Entity: NewEntity(eType, device.Address(), entityAddress),
		device: device,
	}
}

func (r *EntityRemote) Device() api.DeviceRemoteInterface {
	return r.device
}

func (r *EntityRemote) AddFeature(f api.FeatureRemoteInterface) {
	r.features = append(r.features, f)
}

func (r *EntityRemote) Features() []api.FeatureRemoteInterface {
	return r.features
}

func (r *EntityRemote) Feature(addressFeature *model.AddressFeatureType) api.FeatureRemoteInterface {
	for _, f := range r.features {
		if *f.Address().Feature == *addressFeature {
			return f
		}
	}
	return nil
}

func (r *EntityRemote) RemoveAllFeatures() {
	r.features = nil
}
