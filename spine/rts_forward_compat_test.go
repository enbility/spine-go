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

	"github.com/enbility/spine-go/model"
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
// Current behaviour: spine-go PANICS in NewFeatureRemote → function_data_factory.go:305
// Expected behaviour: unknown featureType silently ignored/stored, commissioning proceeds.
func (s *RTSSuite) TestRTS002_DiscoveryReply_UnknownFutureClientType() {
	// This test currently panics — captures the spec violation
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
