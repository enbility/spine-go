package spine

import (
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

var _ api.EntityRemote = (*EntityRemoteImpl)(nil)

type EntityRemoteImpl struct {
	*EntityImpl
	device   api.DeviceRemote
	features []api.FeatureRemote
}

func NewEntityRemoteImpl(device api.DeviceRemote, eType model.EntityTypeType, entityAddress []model.AddressEntityType) *EntityRemoteImpl {
	return &EntityRemoteImpl{
		EntityImpl: NewEntity(eType, device.Address(), entityAddress),
		device:     device,
	}
}

func (r *EntityRemoteImpl) Device() api.DeviceRemote {
	return r.device
}

func (r *EntityRemoteImpl) AddFeature(f api.FeatureRemote) {
	r.features = append(r.features, f)
}

func (r *EntityRemoteImpl) Features() []api.FeatureRemote {
	return r.features
}

func (r *EntityRemoteImpl) Feature(addressFeature *model.AddressFeatureType) api.FeatureRemote {
	for _, f := range r.features {
		if *f.Address().Feature == *addressFeature {
			return f
		}
	}
	return nil
}

func (r *EntityRemoteImpl) RemoveAllFeatures() {
	r.features = nil
}
