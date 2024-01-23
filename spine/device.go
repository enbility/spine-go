package spine

import "github.com/enbility/spine-go/model"

type Device struct {
	address    *model.AddressDeviceType
	dType      *model.DeviceTypeType
	featureSet *model.NetworkManagementFeatureSetType
}

// Initialize a new device
// Both values are required for a local device but provided as empty strings for a remote device
// as the address is only provided via detailed discovery response
func NewDevice(address *model.AddressDeviceType, dType *model.DeviceTypeType, featureSet *model.NetworkManagementFeatureSetType) *Device {
	device := &Device{}

	if dType != nil {
		device.dType = dType
	}

	if address != nil {
		device.address = address
	}

	if featureSet != nil {
		device.featureSet = featureSet
	}

	return device
}

func (r *Device) Address() *model.AddressDeviceType {
	return r.address
}

func (r *Device) DeviceType() *model.DeviceTypeType {
	return r.dType
}

func (r *Device) FeatureSet() *model.NetworkManagementFeatureSetType {
	return r.featureSet
}

func (r *Device) DestinationData() model.NodeManagementDestinationDataType {
	return model.NodeManagementDestinationDataType{
		DeviceDescription: &model.NetworkManagementDeviceDescriptionDataType{
			DeviceAddress: &model.DeviceAddressType{
				Device: r.Address(),
			},
			DeviceType:        r.DeviceType(),
			NetworkFeatureSet: r.FeatureSet(),
		},
	}
}
