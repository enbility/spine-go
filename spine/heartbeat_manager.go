package spine

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type HeartbeatManager struct {
	localDevice  api.DeviceLocalInterface
	localEntity  api.EntityLocalInterface
	localFeature api.FeatureLocalInterface

	heartBeatNum   uint64 // see https://github.com/golang/go/issues/11891
	stopHeartbeatC chan struct{}
	stopMux        sync.Mutex

	subscriptionManager api.SubscriptionManagerInterface
	heartBeatTimeout    *model.DurationType
}

// Create a new Heartbeat Manager which handles sending of heartbeats
func NewHeartbeatManager(localDevice api.DeviceLocalInterface, subscriptionManager api.SubscriptionManagerInterface, timeout time.Duration) *HeartbeatManager {
	h := &HeartbeatManager{
		localDevice:         localDevice,
		subscriptionManager: subscriptionManager,
		heartBeatTimeout:    model.NewDurationType(timeout),
	}

	return h
}

func (c *HeartbeatManager) IsHeartbeatRunning() bool {
	c.stopMux.Lock()
	defer c.stopMux.Unlock()

	if c.stopHeartbeatC != nil && !c.isHeartbeatClosed() {
		return true
	}

	return false
}

// check if there are any heartbeat subscriptions left, otherwise stop creating new ones
// or start creating heartbeats again if needed
func (c *HeartbeatManager) UpdateHeartbeatOnSubscriptions() {
	if c.localEntity == nil {
		return
	}

	featureAddr := c.localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	if featureAddr == nil {
		return
	}

	subscriptions := c.subscriptionManager.SubscriptionsOnFeature(*featureAddr.Address())
	if len(subscriptions) == 0 {
		// stop creating heartbeats
		c.StopHeartbeat()
	} else if !c.IsHeartbeatRunning() {
		// resume creating heartbeats
		_ = c.StartHeartbeat()
	}
}

func (c *HeartbeatManager) SetLocalFeature(entity api.EntityLocalInterface, feature api.FeatureLocalInterface) {
	c.localEntity = entity
	c.localFeature = feature
}

// Start setting heartbeat data
// Make sure the a required FeatureTypeTypeDeviceDiagnosis with the role server is present
// otherwise this will end with an error
// Note: Remote features need to have a subscription to get notifications
func (c *HeartbeatManager) StartHeartbeat() error {
	if c.localEntity == nil {
		return errors.New("unknown entity")
	}

	timeout, err := c.heartBeatTimeout.GetTimeDuration()
	if err != nil {
		return err
	}

	// stop an already running heartbeat
	c.StopHeartbeat()

	c.stopHeartbeatC = make(chan struct{})

	go c.updateHearbeatData(c.stopHeartbeatC, timeout)

	return nil
}

// Stop updating heartbeat data
// Note: No active subscribers will get any further notifications!
func (c *HeartbeatManager) StopHeartbeat() {
	if c.IsHeartbeatRunning() {
		close(c.stopHeartbeatC)
	}
}

func (c *HeartbeatManager) heartbeatData(t time.Time, counter *uint64) *model.DeviceDiagnosisHeartbeatDataType {
	timestamp := t.UTC().Format(time.RFC3339)

	return &model.DeviceDiagnosisHeartbeatDataType{
		Timestamp:        &timestamp,
		HeartbeatCounter: counter,
		HeartbeatTimeout: c.heartBeatTimeout,
	}
}

func (c *HeartbeatManager) updateHearbeatData(stopC chan struct{}, d time.Duration) {
	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:

			heartbeatData := c.heartbeatData(time.Now(), c.heartBeatCounter())

			// updating the data will automatically notify all subscribed remote features
			c.localFeature.SetData(model.FunctionTypeDeviceDiagnosisHeartbeatData, heartbeatData)

		case <-stopC:
			return
		}
	}
}

func (c *HeartbeatManager) isHeartbeatClosed() bool {
	select {
	case <-c.stopHeartbeatC:
		return true
	default:
	}

	return false
}

// TODO heartBeatCounter should be global on CEM level, not on connection level
func (c *HeartbeatManager) heartBeatCounter() *uint64 {
	i := atomic.AddUint64(&c.heartBeatNum, 1)
	return &i
}
