package spine

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

var notFoundError = errors.New("data not found")

// check if a client or server feature address matches
// a combination of a deviceAddress and entityAddress
func isMatchingClientOrServerByDeviceAndEntity(
	clientAddress, serverAddress *model.FeatureAddressType,
	deviceAddress *model.AddressDeviceType,
	entityAddress []model.AddressEntityType,
) bool {
	if deviceAddress == nil || entityAddress == nil {
		return false
	}

	if clientAddress != nil &&
		reflect.DeepEqual(clientAddress.Device, deviceAddress) &&
		reflect.DeepEqual(clientAddress.Entity, entityAddress) {
		return true
	}

	if serverAddress != nil &&
		reflect.DeepEqual(serverAddress.Device, deviceAddress) &&
		reflect.DeepEqual(serverAddress.Entity, entityAddress) {
		return true
	}

	return false
}

// return details for a given remoteDevice of a client and server address
//
// Note: when the feature address and/or entity address is not given,
// it wll return all applicable features and entities
//
// returns an error if any of the addressed features are not found or an
// invalid combination of addresses is given
func addressDetails(
	localDevice api.DeviceLocalInterface,
	remoteDevice api.DeviceRemoteInterface,
	clientAddress, serverAddress *model.FeatureAddressType) (
	localFeature api.FeatureLocalInterface, remoteFeature api.FeatureRemoteInterface,
	localRole, remoteRole model.RoleType, err error) {
	err = nil

	if clientAddress == nil || serverAddress == nil ||
		clientAddress.Device == nil || serverAddress.Device == nil {
		err = errors.New("clientAddress and serverAddress must not be nil")
		return
	}

	// is the local feature the client and the remote feature the server?
	if reflect.DeepEqual(clientAddress.Device, localDevice.Address()) &&
		reflect.DeepEqual(serverAddress.Device, remoteDevice.Address()) {
		localRole = model.RoleTypeClient
		localFeature = localDevice.FeatureByAddress(clientAddress)

		remoteRole = model.RoleTypeServer
		remoteFeature = remoteDevice.FeatureByAddress(serverAddress)
	} else if reflect.DeepEqual(serverAddress.Device, localDevice.Address()) &&
		reflect.DeepEqual(clientAddress.Device, remoteDevice.Address()) {
		// the local device is the server and the remote feature the client
		localRole = model.RoleTypeServer
		localFeature = localDevice.FeatureByAddress(serverAddress)

		remoteRole = model.RoleTypeClient
		remoteFeature = remoteDevice.FeatureByAddress(clientAddress)
	} else {
		err = errors.New("invalid addresses")
		return
	}

	if localFeature == nil {
		err = fmt.Errorf("feature '%s' in local device '%s' not found", serverAddress, *localDevice.Address())
	}
	if remoteFeature == nil {
		err = fmt.Errorf("feature '%s' in remote device '%s' not found", clientAddress, *remoteDevice.Address())
	}

	return
}

func dataCopyOfType[T any](rdata any) (T, error) {
	x := any(*new(T))

	if rdata == nil {
		return x.(T), notFoundError
	}

	v := reflect.ValueOf(rdata)
	kind := v.Kind()
	if kind == reflect.Ptr && v.IsNil() {
		return x.(T), notFoundError
	}

	data, ok := rdata.(T)
	if !ok {
		return x.(T), notFoundError
	}

	return data, nil
}

// Note: the type has to be a pointer!
func LocalFeatureDataCopyOfType[T any](feature api.FeatureLocalInterface, function model.FunctionType) (T, error) {
	return dataCopyOfType[T](feature.DataCopy(function))
}

// Note: the type has to be a pointer!
func RemoteFeatureDataCopyOfType[T any](remote api.FeatureRemoteInterface, function model.FunctionType) (T, error) {
	return dataCopyOfType[T](remote.DataCopy(function))
}
