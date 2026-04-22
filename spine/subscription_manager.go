package spine

import (
	"fmt"
	"reflect"

	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type SubscriptionManager struct {
	localDevice api.DeviceLocalInterface
}

func NewSubscriptionManager(localDevice api.DeviceLocalInterface) *SubscriptionManager {
	c := &SubscriptionManager{
		localDevice: localDevice,
	}

	return c
}

// Add a subscription between a client and server feature where one of each is local and the other one is remote
//
// Note: The device values of both addresses may not be nil
func (c *SubscriptionManager) AddSubscription(remoteDevice api.DeviceRemoteInterface, data model.SubscriptionManagementRequestCallType) error {
	// subscription already exists, we're already in the desired state
	// return success to indicate that the subscription exists and simplify synchronization between local and remote device
	if c.HasSubscription(data.ClientAddress, data.ServerAddress) {
		return nil
	}

	localFeature, remoteFeature, localRole, remoteRole, err := addressDetails(c.localDevice, remoteDevice, data.ClientAddress, data.ServerAddress)
	if err != nil {
		return err
	}

	// the server feature type is optional, only validate it if it is set
	serverFeatureType := data.ServerFeatureType
	if serverFeatureType != nil {
		if err := c.checkRoleAndType(localFeature, localRole, *serverFeatureType); err != nil {
			return err
		}
		if err := c.checkRoleAndType(remoteFeature, remoteRole, *serverFeatureType); err != nil {
			return err
		}
	}

	subscriptionEntry := model.SubscriptionManagementEntryDataType{
		ClientAddress: data.ClientAddress,
		ServerAddress: data.ServerAddress,
	}

	nodeMgmt := c.localDevice.NodeManagement()
	subscriptionData := c.subscriptionData()
	subscriptionData.SubscriptionEntry = append(subscriptionData.SubscriptionEntry, subscriptionEntry)

	nodeMgmt.SetData(model.FunctionTypeNodeManagementSubscriptionData, subscriptionData)

	payload := api.EventPayload{
		Ski:          remoteDevice.Ski(),
		EventType:    api.EventTypeSubscriptionChange,
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

// Remove a subscription between a client and server feature where one of each is local and the other one is remote
//
// Note: The device values of both addresses may not be nil
func (c *SubscriptionManager) RemoveSubscription(remoteDevice api.DeviceRemoteInterface, data model.SubscriptionManagementDeleteCallType) error {
	subscriptionData := c.subscriptionData()

	newSubscriptionData := &model.NodeManagementSubscriptionDataType{
		SubscriptionEntry: []model.SubscriptionManagementEntryDataType{},
	}
	deletedSubscriptions := []model.SubscriptionManagementEntryDataType{}

	for _, item := range subscriptionData.SubscriptionEntry {
		// remove a specific subscription
		if data.ClientAddress.Feature != nil &&
			reflect.DeepEqual(item.ClientAddress, data.ClientAddress) &&
			reflect.DeepEqual(item.ServerAddress, data.ServerAddress) {
			deletedSubscriptions = append(deletedSubscriptions, item)
			continue
		}

		// remove all subscriptions for a specific entity with the same "role-relation"
		if data.ClientAddress.Feature == nil &&
			data.ClientAddress.Entity != nil &&
			reflect.DeepEqual(item.ClientAddress.Device, data.ClientAddress.Device) &&
			reflect.DeepEqual(item.ServerAddress.Device, data.ServerAddress.Device) &&
			reflect.DeepEqual(item.ClientAddress.Entity, data.ClientAddress.Entity) &&
			reflect.DeepEqual(item.ServerAddress.Entity, data.ServerAddress.Entity) {
			deletedSubscriptions = append(deletedSubscriptions, item)
			continue
		}

		// remove all subscriptions for a specific device with the same "role-relation"
		if data.ClientAddress.Feature == nil &&
			data.ClientAddress.Entity == nil &&
			reflect.DeepEqual(item.ClientAddress.Device, data.ClientAddress.Device) &&
			reflect.DeepEqual(item.ServerAddress.Device, data.ServerAddress.Device) {
			deletedSubscriptions = append(deletedSubscriptions, item)
			continue
		}

		newSubscriptionData.SubscriptionEntry = append(newSubscriptionData.SubscriptionEntry, item)
	}

	// we did not find any subscription to delete, so we're already in the desired state
	// return success to indicate that the subscription doesn't exist and simplify synchronization between local and remote device
	if len(deletedSubscriptions) == 0 {
		return nil
	}

	nodeMgmt := c.localDevice.NodeManagement()

	nodeMgmt.SetData(model.FunctionTypeNodeManagementSubscriptionData, newSubscriptionData)

	// inform about every deleted subscription
	for _, item := range deletedSubscriptions {
		if localFeature, remoteFeature, _, _, err := addressDetails(c.localDevice, remoteDevice, item.ClientAddress, item.ServerAddress); err == nil {
			payload := api.EventPayload{
				Ski:          remoteDevice.Ski(),
				EventType:    api.EventTypeSubscriptionChange,
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

// Remove all existing subscriptions for a given remote device
func (c *SubscriptionManager) RemoveSubscriptionsForRemoteDevice(remoteDevice api.DeviceRemoteInterface) {
	if remoteDevice == nil {
		return
	}

	for _, entity := range remoteDevice.Entities() {
		c.RemoveSubscriptionsForRemoteEntity(entity)
	}
}

// Remove all existing subscriptions for a given remote device entity
func (c *SubscriptionManager) RemoveSubscriptionsForRemoteEntity(remoteEntity api.EntityRemoteInterface) {
	if remoteEntity == nil {
		return
	}

	subscriptionData := c.subscriptionData()

	remoteDeviceAddress := remoteEntity.Device().Address()
	remoteEntityAddress := remoteEntity.Address().Entity

	for _, subscription := range subscriptionData.SubscriptionEntry {
		// check if subscription matches ClientAddress or ServerAddress
		if !isMatchingClientOrServerByDeviceAndEntity(
			subscription.ClientAddress, subscription.ServerAddress,
			remoteDeviceAddress, remoteEntityAddress) {
			continue
		}

		_ = c.RemoveSubscription(remoteEntity.Device(), model.SubscriptionManagementDeleteCallType{
			ClientAddress: subscription.ClientAddress,
			ServerAddress: subscription.ServerAddress,
		})
	}
}

// Remove all existing subscriptions for a given local device entity
func (c *SubscriptionManager) RemoveSubscriptionsForLocalEntity(localEntity api.EntityLocalInterface) {
	if localEntity == nil {
		return
	}

	subscriptionData := c.subscriptionData()

	localDeviceAddress := localEntity.Device().Address()
	localEntityAddress := localEntity.Address().Entity

	for _, subscription := range subscriptionData.SubscriptionEntry {
		// check if subscription matches ClientAddress or ServerAddress
		if !isMatchingClientOrServerByDeviceAndEntity(
			subscription.ClientAddress, subscription.ServerAddress,
			localDeviceAddress, localEntityAddress) {
			continue
		}

		var remoteDevice api.DeviceRemoteInterface

		if reflect.DeepEqual(subscription.ClientAddress.Device, localDeviceAddress) {
			// defense in depth in case invalid subscriptions are ever added
			if subscription.ServerAddress == nil || subscription.ServerAddress.Device == nil {
				logging.Log().Debug("skipping invalid subscription with unset ServerAddress")
				continue
			}
			remoteDevice = c.localDevice.RemoteDeviceForAddress(*subscription.ServerAddress.Device)
		} else {
			// defense in depth in case invalid subscriptions are ever added
			if subscription.ClientAddress == nil || subscription.ClientAddress.Device == nil {
				logging.Log().Debug("skipping invalid subscription with unset ClientAddress")
				continue
			}
			remoteDevice = c.localDevice.RemoteDeviceForAddress(*subscription.ClientAddress.Device)
		}

		_ = c.RemoveSubscription(remoteDevice, model.SubscriptionManagementDeleteCallType{
			ClientAddress: subscription.ClientAddress,
			ServerAddress: subscription.ServerAddress,
		})
	}
}

// Checks if a subscription between the client and server feature exists
func (c *SubscriptionManager) HasSubscription(clientAddress, serverAddress *model.FeatureAddressType) bool {
	subscriptionData := c.subscriptionData()

	for _, item := range subscriptionData.SubscriptionEntry {
		if reflect.DeepEqual(item.ClientAddress, clientAddress) &&
			reflect.DeepEqual(item.ServerAddress, serverAddress) {
			return true
		}
	}

	return false
}

// Return all stored subscriptions for a given remote device
func (c *SubscriptionManager) SubscriptionsForRemoteDevice(remoteDevice api.DeviceRemoteInterface) []model.SubscriptionManagementEntryDataType {
	subscriptionData := c.subscriptionData()

	filteredSubscriptions := []model.SubscriptionManagementEntryDataType{}

	if subscriptionData != nil {
		for _, subscription := range subscriptionData.SubscriptionEntry {
			if reflect.DeepEqual(subscription.ClientAddress.Device, remoteDevice.Address()) ||
				reflect.DeepEqual(subscription.ServerAddress.Device, remoteDevice.Address()) {
				filteredSubscriptions = append(filteredSubscriptions, subscription)
			}
		}
	}

	return filteredSubscriptions
}

// Return all stored subscriptions for a given feature address
func (c *SubscriptionManager) SubscriptionsForFeatureAddress(featureAddress model.FeatureAddressType) []model.SubscriptionManagementEntryDataType {
	subscriptionData := c.subscriptionData()

	filteredSubscriptions := []model.SubscriptionManagementEntryDataType{}

	if subscriptionData != nil {
		for _, subscription := range subscriptionData.SubscriptionEntry {
			if reflect.DeepEqual(*subscription.ClientAddress, featureAddress) ||
				reflect.DeepEqual(*subscription.ServerAddress, featureAddress) {
				filteredSubscriptions = append(filteredSubscriptions, subscription)
			}
		}
	}

	return filteredSubscriptions
}

func (c *SubscriptionManager) subscriptionData() *model.NodeManagementSubscriptionDataType {
	nodeMgmt := c.localDevice.NodeManagement()
	subscriptionDataCopy := nodeMgmt.DataCopy(model.FunctionTypeNodeManagementSubscriptionData)
	return subscriptionDataCopy.(*model.NodeManagementSubscriptionDataType)
}

func (c *SubscriptionManager) checkRoleAndType(feature api.FeatureInterface, role model.RoleType, featureType model.FeatureTypeType) error {
	if feature.Role() != model.RoleTypeSpecial && feature.Role() != role {
		return fmt.Errorf("found feature %s is not matching required role %s", feature.Type(), role)
	}

	if feature.Type() != featureType && feature.Type() != model.FeatureTypeTypeGeneric {
		return fmt.Errorf("found feature %s is not matching required type %s", feature.Type(), featureType)
	}

	return nil
}
