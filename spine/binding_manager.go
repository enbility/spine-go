package spine

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type BindingManager struct {
	localDevice api.DeviceLocalInterface
}

func NewBindingManager(localDevice api.DeviceLocalInterface) *BindingManager {
	c := &BindingManager{
		localDevice: localDevice,
	}

	return c
}

// Add a binding between a client and server feature where one of each is local and the other one is remote
//
// Note: The device values of both addresses may not be nil
func (c *BindingManager) AddBinding(remoteDevice api.DeviceRemoteInterface, data model.BindingManagementRequestCallType) error {
	// binding already exists, we're already in the desired state
	// return success to indicate that the binding exists and simplify synchronization between local and remote device
	if c.HasBinding(data.ClientAddress, data.ServerAddress) {
		return nil
	}

	localFeature, remoteFeature, localRole, remoteRole, err := addressDetails(c.localDevice, remoteDevice, data.ClientAddress, data.ServerAddress)
	if err != nil {
		return err
	}

	// the server feature type is optional, only validate it if it is set
	if data.ServerFeatureType != nil {
		if err := c.checkRoleAndType(localFeature, localRole, *data.ServerFeatureType); err != nil {
			return err
		}
		if err := c.checkRoleAndType(remoteFeature, remoteRole, *data.ServerFeatureType); err != nil {
			return err
		}
	}

	// a local feature can only have one remote binding for now
	// see also https://github.com/enbility/spine-go/issues/25
	if localRole == model.RoleTypeServer {
		bindings := c.BindingsForFeatureAddress(*localFeature.Address())
		if len(bindings) > 0 {
			return errors.New("the server feature already has a binding")
		}
	}

	bindingEntry := model.BindingManagementEntryDataType{
		ClientAddress: data.ClientAddress,
		ServerAddress: data.ServerAddress,
	}

	nodeMgmt := c.localDevice.NodeManagement()
	bindingData := c.bindingData()
	bindingData.BindingEntry = append(bindingData.BindingEntry, bindingEntry)

	nodeMgmt.SetData(model.FunctionTypeNodeManagementBindingData, bindingData)

	payload := api.EventPayload{
		Ski:          remoteDevice.Ski(),
		EventType:    api.EventTypeBindingChange,
		ChangeType:   api.ElementChangeAdd,
		Data:         data,
		Device:       remoteDevice,
		Entity:       remoteFeature.Entity(),
		Feature:      remoteFeature,
		LocalFeature: localFeature,
	}
	c.localDevice.Events().Publish(payload)

	return nil
}

// Remove a binding between a client and server feature where one of each is local and the other one is remote
//
// Note: The device values of both addresses may not be nil
func (c *BindingManager) RemoveBinding(remoteDevice api.DeviceRemoteInterface, data model.BindingManagementDeleteCallType) error {
	bindingData := c.bindingData()

	newBindingData := &model.NodeManagementBindingDataType{
		BindingEntry: []model.BindingManagementEntryDataType{},
	}
	deletedBindings := []model.BindingManagementEntryDataType{}

	for _, item := range bindingData.BindingEntry {
		// remove a specific binding
		if data.ClientAddress.Feature != nil &&
			reflect.DeepEqual(item.ClientAddress, data.ClientAddress) &&
			reflect.DeepEqual(item.ServerAddress, data.ServerAddress) {
			deletedBindings = append(deletedBindings, item)
			continue
		}

		// remove all bindings for a specific entity with the same "role-relation"
		if data.ClientAddress.Feature == nil &&
			data.ClientAddress.Entity != nil &&
			reflect.DeepEqual(item.ClientAddress.Device, data.ClientAddress.Device) &&
			reflect.DeepEqual(item.ServerAddress.Device, data.ServerAddress.Device) &&
			reflect.DeepEqual(item.ClientAddress.Entity, data.ClientAddress.Entity) &&
			reflect.DeepEqual(item.ServerAddress.Entity, data.ServerAddress.Entity) {
			deletedBindings = append(deletedBindings, item)
			continue
		}

		// remove all bindings for a specific device with the same "role-relation"
		if data.ClientAddress.Feature == nil &&
			data.ClientAddress.Entity == nil &&
			reflect.DeepEqual(item.ClientAddress.Device, data.ClientAddress.Device) &&
			reflect.DeepEqual(item.ServerAddress.Device, data.ServerAddress.Device) {
			deletedBindings = append(deletedBindings, item)
			continue
		}

		newBindingData.BindingEntry = append(newBindingData.BindingEntry, item)
	}

	// we did not find any binding to delete, so we're already in the desired state
	// return success to indicate that the binding doesn't exist and simplify synchronization between local and remote device
	if len(deletedBindings) == 0 {
		return nil
	}

	nodeMgmt := c.localDevice.NodeManagement()

	nodeMgmt.SetData(model.FunctionTypeNodeManagementBindingData, newBindingData)

	for _, item := range deletedBindings {
		// inform about every deleted binding
		if localFeature, remoteFeature, _, _, err := addressDetails(c.localDevice, remoteDevice, item.ClientAddress, item.ServerAddress); err == nil {
			payload := api.EventPayload{
				Ski:          remoteDevice.Ski(),
				EventType:    api.EventTypeBindingChange,
				ChangeType:   api.ElementChangeRemove,
				Data:         data,
				Device:       remoteDevice,
				Entity:       remoteFeature.Entity(),
				Feature:      remoteFeature,
				LocalFeature: localFeature,
			}
			c.localDevice.Events().Publish(payload)
		}
	}

	return nil
}

// Remove all stored bindings for a given remote device
func (c *BindingManager) RemoveBindingsForRemoteDevice(remoteDevice api.DeviceRemoteInterface) {
	if remoteDevice == nil {
		return
	}

	for _, entity := range remoteDevice.Entities() {
		c.RemoveBindingsForRemoteEntity(entity)
	}
}

// Remove all stored bindings for a given remote device entity
func (c *BindingManager) RemoveBindingsForRemoteEntity(remoteEntity api.EntityRemoteInterface) {
	if remoteEntity == nil {
		return
	}

	bindingData := c.bindingData()

	remoteDeviceAddress := remoteEntity.Device().Address()
	remoteEntityAddress := remoteEntity.Address().Entity

	for _, binding := range bindingData.BindingEntry {
		// check if binding matches ClientAddress or ServerAddress
		if !isMatchingClientOrServerByDeviceAndEntity(
			binding.ClientAddress, binding.ServerAddress,
			remoteDeviceAddress, remoteEntityAddress) {
			continue
		}

		_ = c.RemoveBinding(remoteEntity.Device(), model.BindingManagementDeleteCallType{
			ClientAddress: binding.ClientAddress,
			ServerAddress: binding.ServerAddress,
		})
	}
}

// Remove all stored bindings for a given local device entity
func (c *BindingManager) RemoveBindingsForLocalEntity(localEntity api.EntityLocalInterface) {
	if localEntity == nil {
		return
	}

	bindingData := c.bindingData()

	localDeviceAddress := localEntity.Device().Address()
	localEntityAddress := localEntity.Address().Entity

	for _, binding := range bindingData.BindingEntry {
		// check if binding matches ClientAddress or ServerAddress
		if !isMatchingClientOrServerByDeviceAndEntity(
			binding.ClientAddress, binding.ServerAddress,
			localDeviceAddress, localEntityAddress) {
			continue
		}

		var remoteDevice api.DeviceRemoteInterface

		if reflect.DeepEqual(binding.ClientAddress.Device, localDeviceAddress) {
			// defense in depth in case invalid bindings are ever added
			if binding.ServerAddress == nil || binding.ServerAddress.Device == nil {
				logging.Log().Debug("skipping invalid binding with unset ServerAddress")
				continue
			}
			remoteDevice = c.localDevice.RemoteDeviceForAddress(*binding.ServerAddress.Device)
		} else {
			// defense in depth in case invalid bindings are ever added
			if binding.ClientAddress == nil || binding.ClientAddress.Device == nil {
				logging.Log().Debug("skipping invalid binding with unset ClientAddress")
				continue
			}
			remoteDevice = c.localDevice.RemoteDeviceForAddress(*binding.ClientAddress.Device)
		}

		_ = c.RemoveBinding(remoteDevice, model.BindingManagementDeleteCallType{
			ClientAddress: binding.ClientAddress,
			ServerAddress: binding.ServerAddress,
		})
	}
}

// Checks if a binding between the client and server feature exists
func (c *BindingManager) HasBinding(clientAddress, serverAddress *model.FeatureAddressType) bool {
	bindingData := c.bindingData()

	for _, item := range bindingData.BindingEntry {
		if reflect.DeepEqual(item.ClientAddress, clientAddress) &&
			reflect.DeepEqual(item.ServerAddress, serverAddress) {
			return true
		}
	}

	return false
}

// Return all stored bindings for a given remote device
func (c *BindingManager) BindingsForRemoteDevice(remoteDevice api.DeviceRemoteInterface) []model.BindingManagementEntryDataType {
	bindingData := c.bindingData()

	filteredBindings := []model.BindingManagementEntryDataType{}

	if bindingData != nil {
		for _, binding := range bindingData.BindingEntry {
			if reflect.DeepEqual(binding.ClientAddress.Device, remoteDevice.Address()) ||
				reflect.DeepEqual(binding.ServerAddress.Device, remoteDevice.Address()) {
				filteredBindings = append(filteredBindings, binding)
			}
		}
	}

	return filteredBindings
}

// Return all stored bindings for a given feature address
func (c *BindingManager) BindingsForFeatureAddress(featureAddress model.FeatureAddressType) []model.BindingManagementEntryDataType {
	bindingData := c.bindingData()

	filteredBindings := []model.BindingManagementEntryDataType{}

	if bindingData != nil {
		for _, binding := range bindingData.BindingEntry {
			if reflect.DeepEqual(*binding.ClientAddress, featureAddress) ||
				reflect.DeepEqual(*binding.ServerAddress, featureAddress) {
				filteredBindings = append(filteredBindings, binding)
			}
		}
	}

	return filteredBindings
}

func (c *BindingManager) bindingData() *model.NodeManagementBindingDataType {
	nodeMgmt := c.localDevice.NodeManagement()
	bindingDataCopy := nodeMgmt.DataCopy(model.FunctionTypeNodeManagementBindingData)
	return bindingDataCopy.(*model.NodeManagementBindingDataType)
}

func (c *BindingManager) checkRoleAndType(feature api.FeatureInterface, role model.RoleType, featureType model.FeatureTypeType) error {
	if feature.Role() != model.RoleTypeSpecial && feature.Role() != role {
		return fmt.Errorf("found feature %s is not matching required role %s", feature.Type(), role)
	}

	if feature.Type() != featureType && feature.Type() != model.FeatureTypeTypeGeneric {
		return fmt.Errorf("found feature %s is not matching required type %s", feature.Type(), featureType)
	}

	return nil
}
