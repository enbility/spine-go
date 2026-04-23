package spine

import (
	"encoding/json"

	"github.com/enbility/spine-go/model"
)

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

// Add support for JSON Marshalling
//
// Instances of EntityInterface are used as arguments and return values in various API calls,
// therefor it is helpfull to be able to marshal them to JSON and thus make the API calls
// usable with various communication interfaces
func (r *Device) MarshalJSON() ([]byte, error) {
	var tempAddress string

	if r.address != nil {
		tempAddress = string(*r.address)
	}

	bytes, err := json.Marshal(tempAddress)
	if err != nil {
		return nil, err
	}

	return bytes, nil
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
