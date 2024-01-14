package spine

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/ahmetb/go-linq/v3"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

type SubscriptionManagerImpl struct {
	localDevice api.DeviceLocal

	subscriptionNum     uint64
	subscriptionEntries []*api.SubscriptionEntry

	mux sync.Mutex
	// TODO: add persistence
}

func NewSubscriptionManager(localDevice api.DeviceLocal) *SubscriptionManagerImpl {
	c := &SubscriptionManagerImpl{
		subscriptionNum: 0,
		localDevice:     localDevice,
	}

	return c
}

// is sent from the client (remote device) to the server (local device)
func (c *SubscriptionManagerImpl) AddSubscription(remoteDevice api.DeviceRemote, data model.SubscriptionManagementRequestCallType) error {

	serverFeature := c.localDevice.FeatureByAddress(data.ServerAddress)
	if serverFeature == nil {
		return fmt.Errorf("server feature '%s' in local device '%s' not found", data.ServerAddress, *c.localDevice.Address())
	}
	if err := c.checkRoleAndType(serverFeature, model.RoleTypeServer, *data.ServerFeatureType); err != nil {
		return err
	}

	clientFeature := remoteDevice.FeatureByAddress(data.ClientAddress)
	if clientFeature == nil {
		return fmt.Errorf("client feature '%s' in remote device '%s' not found", data.ClientAddress, *remoteDevice.Address())
	}
	if err := c.checkRoleAndType(clientFeature, model.RoleTypeClient, *data.ServerFeatureType); err != nil {
		return err
	}

	subscriptionEntry := &api.SubscriptionEntry{
		Id:            c.subscriptionId(),
		ServerFeature: serverFeature,
		ClientFeature: clientFeature,
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	for _, item := range c.subscriptionEntries {
		if reflect.DeepEqual(item.ServerFeature, serverFeature) && reflect.DeepEqual(item.ClientFeature, clientFeature) {
			return fmt.Errorf("requested subscription is already present")
		}
	}

	c.subscriptionEntries = append(c.subscriptionEntries, subscriptionEntry)

	payload := api.EventPayload{
		Ski:        remoteDevice.Ski(),
		EventType:  api.EventTypeSubscriptionChange,
		ChangeType: api.ElementChangeAdd,
		Data:       data,
		Feature:    clientFeature,
	}
	Events.Publish(payload)

	return nil
}

// Remove a specific subscription that is provided by a delete message from a remote device
func (c *SubscriptionManagerImpl) RemoveSubscription(data model.SubscriptionManagementDeleteCallType, remoteDevice api.DeviceRemote) error {
	var newSubscriptionEntries []*api.SubscriptionEntry

	// according to the spec 7.4.4
	// a. The absence of "subscriptionDelete. clientAddress. device" SHALL be treated as if it was
	//    present and set to the sender's "device" address part.
	// b. The absence of "subscriptionDelete. serverAddress. device" SHALL be treated as if it was
	//    present and set to the recipient's "device" address part.

	var clientAddress model.FeatureAddressType
	util.DeepCopy(data.ClientAddress, &clientAddress)
	if data.ClientAddress.Device == nil {
		clientAddress.Device = remoteDevice.Address()
	}

	clientFeature := remoteDevice.FeatureByAddress(data.ClientAddress)
	if clientFeature == nil {
		return fmt.Errorf("client feature '%s' in remote device '%s' not found", data.ClientAddress, *remoteDevice.Address())
	}

	serverFeature := c.localDevice.FeatureByAddress(data.ServerAddress)
	if serverFeature == nil {
		return fmt.Errorf("server feature '%s' in local device '%s' not found", data.ServerAddress, *c.localDevice.Address())
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	for _, item := range c.subscriptionEntries {
		itemAddress := item.ClientFeature.Address()

		if !reflect.DeepEqual(*itemAddress, clientAddress) &&
			!reflect.DeepEqual(item.ServerFeature, serverFeature) {
			newSubscriptionEntries = append(newSubscriptionEntries, item)
		}
	}

	if len(newSubscriptionEntries) == len(c.subscriptionEntries) {
		return errors.New("could not find requested SubscriptionId to be removed")
	}

	c.subscriptionEntries = newSubscriptionEntries

	payload := api.EventPayload{
		Ski:        remoteDevice.Ski(),
		EventType:  api.EventTypeSubscriptionChange,
		ChangeType: api.ElementChangeRemove,
		Data:       data,
		Device:     remoteDevice,
		Feature:    clientFeature,
	}
	Events.Publish(payload)

	return nil
}

// Remove all existing subscriptions for a given remote device
func (c *SubscriptionManagerImpl) RemoveSubscriptionsForDevice(remoteDevice api.DeviceRemote) {
	if remoteDevice == nil {
		return
	}

	for _, entity := range remoteDevice.Entities() {
		c.RemoveSubscriptionsForEntity(entity)
	}
}

// Remove all existing subscriptions for a given remote device entity
func (c *SubscriptionManagerImpl) RemoveSubscriptionsForEntity(remoteEntity api.EntityRemote) {
	if remoteEntity == nil {
		return
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	var newSubscriptionEntries []*api.SubscriptionEntry
	for _, item := range c.subscriptionEntries {
		if !reflect.DeepEqual(item.ClientFeature.Address().Entity, remoteEntity.Address().Entity) {
			newSubscriptionEntries = append(newSubscriptionEntries, item)
			continue
		}

		clientFeature := remoteEntity.Feature(item.ClientFeature.Address().Feature)
		payload := api.EventPayload{
			Ski:        remoteEntity.Device().Ski(),
			EventType:  api.EventTypeSubscriptionChange,
			ChangeType: api.ElementChangeRemove,
			Entity:     remoteEntity,
			Feature:    clientFeature,
		}
		Events.Publish(payload)
	}

	c.subscriptionEntries = newSubscriptionEntries
}

func (c *SubscriptionManagerImpl) Subscriptions(remoteDevice api.DeviceRemote) []*api.SubscriptionEntry {
	var result []*api.SubscriptionEntry

	c.mux.Lock()
	defer c.mux.Unlock()

	linq.From(c.subscriptionEntries).WhereT(func(s *api.SubscriptionEntry) bool {
		return s.ClientFeature.Device().Ski() == remoteDevice.Ski()
	}).ToSlice(&result)

	return result
}

func (c *SubscriptionManagerImpl) SubscriptionsOnFeature(featureAddress model.FeatureAddressType) []*api.SubscriptionEntry {
	var result []*api.SubscriptionEntry

	c.mux.Lock()
	defer c.mux.Unlock()

	linq.From(c.subscriptionEntries).WhereT(func(s *api.SubscriptionEntry) bool {
		return reflect.DeepEqual(*s.ServerFeature.Address(), featureAddress)
	}).ToSlice(&result)

	return result
}

func (c *SubscriptionManagerImpl) subscriptionId() uint64 {
	i := atomic.AddUint64(&c.subscriptionNum, 1)
	return i
}

func (c *SubscriptionManagerImpl) checkRoleAndType(feature api.Feature, role model.RoleType, featureType model.FeatureTypeType) error {
	if feature.Role() != model.RoleTypeSpecial && feature.Role() != role {
		return fmt.Errorf("found feature %s is not matching required role %s", feature.Type(), role)
	}

	if feature.Type() != featureType && feature.Type() != model.FeatureTypeTypeGeneric {
		return fmt.Errorf("found feature %s is not matching required type %s", feature.Type(), featureType)
	}

	return nil
}
