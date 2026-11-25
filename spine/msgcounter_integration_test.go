package spine

import (
	"encoding/json"
	"testing"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// MsgCounterIntegrationSuite tests msgCounter behavior in multi-device scenarios
type MsgCounterIntegrationSuite struct {
	suite.Suite
	
	localDevice   api.DeviceLocalInterface
	remoteDevice1 api.DeviceRemoteInterface
	remoteDevice2 api.DeviceRemoteInterface
	sentMessages  [][]byte
}

func (s *MsgCounterIntegrationSuite) SetupTest() {
	s.sentMessages = [][]byte{}
	
	// Setup local device
	s.localDevice = NewDeviceLocal("brand", "model", "serial", "code", "local", 
		model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)
	
	// Setup remote device 1
	ski1 := "device1"
	sender1 := NewSender(s)
	s.remoteDevice1 = NewDeviceRemote(s.localDevice, ski1, sender1)
	desc1 := &model.NetworkManagementDeviceDescriptionDataType{
		DeviceAddress: &model.DeviceAddressType{
			Device: util.Ptr(model.AddressDeviceType("device1")),
		},
	}
	s.remoteDevice1.UpdateDevice(desc1)
	_ = s.localDevice.SetupRemoteDevice(ski1, s)
	
	// Setup remote device 2
	ski2 := "device2"
	sender2 := NewSender(s)
	s.remoteDevice2 = NewDeviceRemote(s.localDevice, ski2, sender2)
	desc2 := &model.NetworkManagementDeviceDescriptionDataType{
		DeviceAddress: &model.DeviceAddressType{
			Device: util.Ptr(model.AddressDeviceType("device2")),
		},
	}
	s.remoteDevice2.UpdateDevice(desc2)
	_ = s.localDevice.SetupRemoteDevice(ski2, s)
}

func (s *MsgCounterIntegrationSuite) WriteShipMessageWithPayload(message []byte) {
	s.sentMessages = append(s.sentMessages, message)
}

func (s *MsgCounterIntegrationSuite) CloseDataConnection(err error, removeI bool) {}
func (s *MsgCounterIntegrationSuite) CloseRemoteConnection(ski string, writeI bool) {}
func (s *MsgCounterIntegrationSuite) IsDataConnectionClosed() bool { return false }
func (s *MsgCounterIntegrationSuite) RemoteSKI() string { return "test-ski" }

func TestMsgCounterIntegrationSuite(t *testing.T) {
	suite.Run(t, new(MsgCounterIntegrationSuite))
}

// Test that incoming msgCounters are extracted correctly
func (s *MsgCounterIntegrationSuite) Test_IncomingMsgCounter_Extraction() {
	// Create test message with specific msgCounter
	testMsgCounter := model.MsgCounterType(42)
	
	datagram := model.Datagram{
		Datagram: model.DatagramType{
			Header: model.HeaderType{
				SpecificationVersion: &SpecificationVersion,
				AddressSource: &model.FeatureAddressType{
					Device:  util.Ptr(model.AddressDeviceType("device1")),
					Feature: util.Ptr(model.AddressFeatureType(0)),
				},
				AddressDestination: &model.FeatureAddressType{
					Device:  util.Ptr(model.AddressDeviceType("local")),
					Feature: util.Ptr(model.AddressFeatureType(0)),
				},
				MsgCounter:   &testMsgCounter,
				CmdClassifier: util.Ptr(model.CmdClassifierTypeNotify),
			},
			Payload: model.PayloadType{
				Cmd: []model.CmdType{
					{
						NodeManagementDetailedDiscoveryData: &model.NodeManagementDetailedDiscoveryDataType{},
					},
				},
			},
		},
	}
	
	message, err := json.Marshal(datagram)
	assert.NoError(s.T(), err)
	
	// Handle the message
	receivedCounter, err := s.remoteDevice1.HandleSpineMesssage(message)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), receivedCounter)
	assert.Equal(s.T(), testMsgCounter, *receivedCounter)
}

// Test that msgCounters from different devices are independent
func (s *MsgCounterIntegrationSuite) Test_MultiDevice_Independent_Counters() {
	// Device 1 sends messages with counters 10, 11, 12
	counters1 := []model.MsgCounterType{10, 11, 12}
	for _, counter := range counters1 {
		msg := s.createTestMessage("device1", "local", counter)
		receivedCounter, err := s.remoteDevice1.HandleSpineMesssage(msg)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), counter, *receivedCounter)
	}
	
	// Device 2 sends messages with counters 20, 21, 22
	counters2 := []model.MsgCounterType{20, 21, 22}
	for _, counter := range counters2 {
		msg := s.createTestMessage("device2", "local", counter)
		receivedCounter, err := s.remoteDevice2.HandleSpineMesssage(msg)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), counter, *receivedCounter)
	}
}

// Test device reset scenario - msgCounter drops from high to low value
// This test demonstrates the missing implementation of SPINE spec requirement
// TC_SPINE_DATA_003: process incoming datagrams as usual even if the msgCounter is lower than the last received (reset/overflow).
func (s *MsgCounterIntegrationSuite) Test_DeviceReset_Detection_Gap() {
	// Device sends messages with increasing counters
	normalCounters := []model.MsgCounterType{100, 101, 102}
	for _, counter := range normalCounters {
		msg := s.createTestMessage("device1", "local", counter)
		_, err := s.remoteDevice1.HandleSpineMesssage(msg)
		assert.NoError(s.T(), err)
	}
	
	// Device resets - msgCounter drops to low value
	resetCounter := model.MsgCounterType(1)
	msg := s.createTestMessage("device1", "local", resetCounter)
	receivedCounter, err := s.remoteDevice1.HandleSpineMesssage(msg)
	
	// Current implementation: message is processed normally
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), resetCounter, *receivedCounter)
	
	// SPEC REQUIREMENT (not implemented):
	// "If a SPINE device 'A' receives a message 'X' from SPINE device 'B' with a 
	// msgCounter less or equal than the last msgCounter received from device 'B', 
	// 'A' SHALL process the message 'X' as usual."
	// "Afterwards, device 'A' SHALL use the unexpectedly low msgCounter value as 
	// the last msgCounter received from device 'B'."
	
	// IMPORTANT NOTE: This tracking requirement is NOT FUNCTIONALLY CRITICAL
	// The spec mandates tracking but provides NO functional use for the tracked data:
	// - Messages are processed identically regardless of msgCounter value
	// - The only specified use is optional: "MAY report this to the user"
	// - No duplicate detection, replay prevention, or ordering enforcement
	// This is effectively a diagnostic-only requirement with no functional benefit
	
	// This test demonstrates that:
	// 1. No tracking of last received msgCounter per device
	// 2. No detection of device resets
	// 3. No special handling for unexpectedly low msgCounter values
}

// Test duplicate message processing - SPEC COMPLIANT BEHAVIOR
func (s *MsgCounterIntegrationSuite) Test_Duplicate_Message_Processing_Compliant() {
	// Send same message (same msgCounter) multiple times
	duplicateCounter := model.MsgCounterType(50)
	
	for i := 0; i < 3; i++ {
		msg := s.createTestMessage("device1", "local", duplicateCounter)
		receivedCounter, err := s.remoteDevice1.HandleSpineMesssage(msg)
		
		// All duplicates are processed normally
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), duplicateCounter, *receivedCounter)
	}
	
	// SPEC COMPLIANT: Current implementation processes all messages regardless of msgCounter
	// Per SPINE spec 5.2.3.1: "msgCounter less or equal... SHALL process the message 'X' as usual"
	// This means duplicate messages (same msgCounter) MUST be processed normally
	// There is NO deduplication requirement in SPINE - this is correct behavior
}

// Test msgCounterReference handling for request/response correlation
func (s *MsgCounterIntegrationSuite) Test_MsgCounterReference_Correlation() {
	// Get local feature
	localEntity := s.localDevice.Entities()[0]
	localFeature := localEntity.Features()[0]
	
	// Create remote feature
	remoteEntity := NewEntityRemote(s.remoteDevice1, model.EntityTypeTypeEVSE, []model.AddressEntityType{1})
	remoteFeature := NewFeatureRemote(0, remoteEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	remoteEntity.AddFeature(remoteFeature)
	
	// Send request
	cmd := model.CmdType{
		NodeManagementDetailedDiscoveryData: &model.NodeManagementDetailedDiscoveryDataType{},
	}
	
	// Get sender from local device
	sender := NewSender(s)
	msgCounter, err := sender.Request(
		model.CmdClassifierTypeRead,
		localFeature.Address(),
		remoteFeature.Address(),
		false,
		[]model.CmdType{cmd},
	)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), msgCounter)
	
	// Create response with msgCounterReference
	responseDatagram := model.Datagram{
		Datagram: model.DatagramType{
			Header: model.HeaderType{
				SpecificationVersion: &SpecificationVersion,
				AddressSource: &model.FeatureAddressType{
					Device:  util.Ptr(model.AddressDeviceType("device1")),
					Feature: util.Ptr(model.AddressFeatureType(0)),
				},
				AddressDestination: &model.FeatureAddressType{
					Device:  util.Ptr(model.AddressDeviceType("local")),
					Feature: util.Ptr(model.AddressFeatureType(0)),
				},
				MsgCounter:          util.Ptr(model.MsgCounterType(200)),
				MsgCounterReference: msgCounter, // Reference to request
				CmdClassifier:       util.Ptr(model.CmdClassifierTypeReply),
			},
			Payload: model.PayloadType{
				Cmd: []model.CmdType{cmd},
			},
		},
	}
	
	responseMessage, err := json.Marshal(responseDatagram)
	assert.NoError(s.T(), err)
	
	// Handle response - should process msgCounterReference
	_, err = s.remoteDevice1.HandleSpineMesssage(responseMessage)
	assert.NoError(s.T(), err)
}

// Helper method to create test messages
func (s *MsgCounterIntegrationSuite) createTestMessage(source, dest string, msgCounter model.MsgCounterType) []byte {
	datagram := model.Datagram{
		Datagram: model.DatagramType{
			Header: model.HeaderType{
				SpecificationVersion: &SpecificationVersion,
				AddressSource: &model.FeatureAddressType{
					Device:  util.Ptr(model.AddressDeviceType(source)),
					Feature: util.Ptr(model.AddressFeatureType(0)),
				},
				AddressDestination: &model.FeatureAddressType{
					Device:  util.Ptr(model.AddressDeviceType(dest)),
					Feature: util.Ptr(model.AddressFeatureType(0)),
				},
				MsgCounter:    &msgCounter,
				CmdClassifier: util.Ptr(model.CmdClassifierTypeNotify),
			},
			Payload: model.PayloadType{
				Cmd: []model.CmdType{
					{
						NodeManagementDetailedDiscoveryData: &model.NodeManagementDetailedDiscoveryDataType{},
					},
				},
			},
		},
	}
	
	message, _ := json.Marshal(datagram)
	return message
}