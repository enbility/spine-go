package spine

// TC_SPINE_RTS_001 / TC_SPINE_RTS_002 — forward compatibility: tolerate arbitrary/unknown client
// Feature types in nodeManagementDetailedDiscoveryData replies.
//
// Spec scenario: test tool (EG) sends a discovery reply where its client Feature is announced
// with an unexpected featureType. DUT (CS/server) MUST process the reply without error and
// proceed with commissioning (binding, subscription).
//
// PAR_replyDiscoveryKnownArbitraryType  → featureType = "DeviceClassification", role = "client"
// PAR_replyDiscoveryUnknownFutureType   → featureType = "CompletelyUnknownFeatureType", role = "client",
//                                         specificationVersion = "1.999.999"

import (
	"testing"
	"time"

	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	rts001_discovery_reply_file = "./testdata/rts001_discovery_reply_arbitrary_type.json"
	rts002_discovery_reply_file = "./testdata/rts002_discovery_reply_unknown_type.json"
)

func TestRTSSuite(t *testing.T) {
	suite.Run(t, new(RTSSuite))
}

type RTSSuite struct {
	suite.Suite
	sut          *DeviceLocal
	writeHandler *WriteMessageHandler
	remoteDevice *DeviceRemote
}

func (s *RTSSuite) BeforeTest(_, _ string) {
	s.sut = NewDeviceLocal("brand", "model", "serial", "code", "HEMS",
		model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)
	s.writeHandler = &WriteMessageHandler{}
	_ = s.sut.SetupRemoteDevice("eg-ski", s.writeHandler)
	s.remoteDevice = s.sut.RemoteDeviceForSki("eg-ski").(*DeviceRemote)
}

// TC_SPINE_RTS_001: Discovery reply with arbitrary KNOWN featureType on client Feature.
// DUT MUST process without error — entity and features stored, DeviceDiagnosis server
// still reachable for subscription.
func (s *RTSSuite) TestRTS001_DiscoveryReply_ArbitraryKnownClientType() {
	_, err := s.remoteDevice.HandleSpineMesssage(loadFileData(s.T(), rts001_discovery_reply_file))

	// Must not error — DUT silently tolerates unexpected client featureType
	assert.Nil(s.T(), err, "TC_SPINE_RTS_001: HandleSpineMesssage must not error for arbitrary known featureType")

	// Remote entity must be present
	remoteEntity := s.remoteDevice.Entity([]model.AddressEntityType{1})
	assert.NotNil(s.T(), remoteEntity, "TC_SPINE_RTS_001: remote entity must be discoverable after reply")

	// DeviceDiagnosis server feature (used by CS to subscribe for heartbeats) must be reachable
	ddServer := s.remoteDevice.FeatureByEntityTypeAndRole(
		remoteEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	assert.NotNil(s.T(), ddServer,
		"TC_SPINE_RTS_001: DeviceDiagnosis server feature must be found — CS needs it to subscribe for heartbeats")
}

// TC_SPINE_RTS_002: Discovery reply with completely UNKNOWN future featureType on client Feature,
// plus higher specificationVersion ("1.999.999").
// Expected behaviour: unknown featureType silently ignored/stored, commissioning proceeds.
func (s *RTSSuite) TestRTS002_DiscoveryReply_UnknownFutureClientType() {
	assert.NotPanics(s.T(), func() {
		_, _ = s.remoteDevice.HandleSpineMesssage(loadFileData(s.T(), rts002_discovery_reply_file))
	}, "TC_SPINE_RTS_002: HandleSpineMesssage must not panic for unknown future featureType")

	// Remote entity must be present after processing
	remoteEntity := s.remoteDevice.Entity([]model.AddressEntityType{1})
	assert.NotNil(s.T(), remoteEntity,
		"TC_SPINE_RTS_002: remote entity must be discoverable after reply with unknown featureType")

	// Known features (DeviceDiagnosis server) must still be reachable
	ddServer := s.remoteDevice.FeatureByEntityTypeAndRole(
		remoteEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	assert.NotNil(s.T(), ddServer,
		"TC_SPINE_RTS_002: DeviceDiagnosis server feature must survive alongside unknown featureType")
}

// commissionFromDiscoveryReply processes a discovery reply and then runs the commissioning the
// EG performs next: a binding and a subscription request from its client Feature (entity 1,
// feature 1 — the one announced with the unexpected featureType) to the DUT's LoadControl
// server Feature.
func (s *RTSSuite) commissionFromDiscoveryReply(replyFile string) (bindErr, subscribeErr error) {
	entity := NewEntityLocal(s.sut, model.EntityTypeTypeCEM, []model.AddressEntityType{1}, time.Second*4)
	s.sut.AddEntity(entity)
	localServerFeature := entity.GetOrAddFeature(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)

	_, err := s.remoteDevice.HandleSpineMesssage(loadFileData(s.T(), replyFile))
	assert.Nil(s.T(), err)

	clientAddress := &model.FeatureAddressType{
		Device:  util.Ptr(model.AddressDeviceType("EnergyGuard")),
		Entity:  []model.AddressEntityType{1},
		Feature: util.Ptr(model.AddressFeatureType(1)),
	}

	bindErr = s.sut.BindingManager().AddBinding(s.remoteDevice, model.BindingManagementRequestCallType{
		ClientAddress:     clientAddress,
		ServerAddress:     localServerFeature.Address(),
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeLoadControl),
	})

	subscribeErr = s.sut.SubscriptionManager().AddSubscription(s.remoteDevice, model.SubscriptionManagementRequestCallType{
		ClientAddress:     clientAddress,
		ServerAddress:     localServerFeature.Address(),
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeLoadControl),
	})

	return bindErr, subscribeErr
}

// TC_SPINE_RTS_001: commissioning must succeed even though the EG announced its client Feature
// as "DeviceClassification". serverFeatureType describes the server side only.
func (s *RTSSuite) TestRTS001_Commissioning_ArbitraryKnownClientType() {
	bindErr, subscribeErr := s.commissionFromDiscoveryReply(rts001_discovery_reply_file)

	assert.Nil(s.T(), bindErr,
		"TC_SPINE_RTS_001: binding must be accepted for a client Feature announced as DeviceClassification")
	assert.Nil(s.T(), subscribeErr,
		"TC_SPINE_RTS_001: subscription must be accepted for a client Feature announced as DeviceClassification")
}

// TC_SPINE_RTS_002: the same for a client Feature type this release does not know.
func (s *RTSSuite) TestRTS002_Commissioning_UnknownFutureClientType() {
	bindErr, subscribeErr := s.commissionFromDiscoveryReply(rts002_discovery_reply_file)

	assert.Nil(s.T(), bindErr,
		"TC_SPINE_RTS_002: binding must be accepted for a client Feature with an unknown featureType")
	assert.Nil(s.T(), subscribeErr,
		"TC_SPINE_RTS_002: subscription must be accepted for a client Feature with an unknown featureType")
}
