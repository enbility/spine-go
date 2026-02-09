package spine

import (
	"testing"

	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/suite"
)

func TestSendVersionSuite(t *testing.T) {
	suite.Run(t, new(SendVersionSuite))
}

type SendVersionSuite struct {
	suite.Suite
	localDevice  *DeviceLocal
	remoteDevice *DeviceRemote
	writeHandler *WriteMessageHandler
}

func (s *SendVersionSuite) SetupTest() {
	s.localDevice = NewDeviceLocal("test", "test", "test", "test", "local", model.DeviceTypeTypeEnergyManagementSystem, "")
	s.writeHandler = &WriteMessageHandler{}
	
	// Create sender with local device for version lookups
	sender := NewSender(s.writeHandler, s.localDevice)
	s.remoteDevice = NewDeviceRemote(s.localDevice, "test-ski", sender)
	
	// Set the device address for lookup
	s.remoteDevice.address = util.Ptr(model.AddressDeviceType("test-ski"))
	
	// Add the remote device to local device's registry
	s.localDevice.AddRemoteDeviceForSki("test-ski", s.remoteDevice)
}

// Test that outgoing messages use our local negotiated version
func (s *SendVersionSuite) Test_OutgoingMessage_UsesNegotiatedVersion() {
	// Set up version negotiation
	s.remoteDevice.SetSupportedProtocolVersions([]string{"1.2.0"})
	s.remoteDevice.SetNegotiatedProtocolVersion("1.2.0") // We negotiated to use 1.2.0 when talking to this device
	
	// Verify the negotiation is set
	s.Equal("1.2.0", s.remoteDevice.NegotiatedProtocolVersion())
	
	// Also verify the device can be looked up by address
	s.NotNil(s.remoteDevice.Address())
	deviceFound := s.localDevice.RemoteDeviceForAddress(*s.remoteDevice.Address())
	s.Equal(s.remoteDevice, deviceFound)
	
	// Create addresses
	sourceAddr := &model.FeatureAddressType{
		Device:  util.Ptr(model.AddressDeviceType("local")),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(0)),
	}
	
	destAddr := &model.FeatureAddressType{
		Device:  util.Ptr(model.AddressDeviceType("test-ski")),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(1)),
	}
	
	// Send a request
	sender := s.remoteDevice.Sender()
	_, err := sender.Request(model.CmdClassifierTypeRead, sourceAddr, destAddr, false, []model.CmdType{})
	s.NoError(err)
	
	// Check that the message was sent with negotiated version
	lastMsg := s.writeHandler.LastMessage()
	s.NotNil(lastMsg)
	s.Contains(string(lastMsg), "\"specificationVersion\":\"1.2.0\"")
}

// Test that messages to unknown devices use default version
func (s *SendVersionSuite) Test_UnknownDevice_UsesDefaultVersion() {
	// Create address to unknown device
	sourceAddr := &model.FeatureAddressType{
		Device:  util.Ptr(model.AddressDeviceType("local")),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(0)),
	}
	
	destAddr := &model.FeatureAddressType{
		Device:  util.Ptr(model.AddressDeviceType("unknown-device")),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(1)),
	}
	
	// Create a new sender without local device context
	writeHandler2 := &WriteMessageHandler{}
	sender := NewSender(writeHandler2, nil)
	
	// Send a request
	_, err := sender.Request(model.CmdClassifierTypeRead, sourceAddr, destAddr, false, []model.CmdType{})
	s.NoError(err)
	
	// Check that the message was sent with default version (1.3.0)
	lastMsg := writeHandler2.LastMessage()
	s.NotNil(lastMsg)
	s.Contains(string(lastMsg), "\"specificationVersion\":\"1.3.0\"")
}

// Test all message types use negotiated version
func (s *SendVersionSuite) Test_AllMessageTypes_UseNegotiatedVersion() {
	// Setup negotiated version
	s.remoteDevice.SetSupportedProtocolVersions([]string{"1.1.0"})
	s.remoteDevice.SetNegotiatedProtocolVersion("1.1.0")
	
	sourceAddr := &model.FeatureAddressType{
		Device:  util.Ptr(model.AddressDeviceType("local")),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(0)),
	}
	
	destAddr := &model.FeatureAddressType{
		Device:  util.Ptr(model.AddressDeviceType("test-ski")),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(1)),
	}
	
	sender := s.remoteDevice.Sender()
	
	// Test Request
	_, err := sender.Request(model.CmdClassifierTypeRead, sourceAddr, destAddr, false, []model.CmdType{})
	s.NoError(err)
	lastMsg := s.writeHandler.LastMessage()
	s.Contains(string(lastMsg), "\"specificationVersion\":\"1.1.0\"")
	
	// Test Notify
	s.writeHandler.sentMessages = nil // Clear previous messages
	_, err = sender.Notify(sourceAddr, destAddr, model.CmdType{})
	s.NoError(err)
	lastMsg = s.writeHandler.LastMessage()
	s.Contains(string(lastMsg), "\"specificationVersion\":\"1.1.0\"")
	
	// Test Write
	s.writeHandler.sentMessages = nil
	_, err = sender.Write(sourceAddr, destAddr, model.CmdType{})
	s.NoError(err)
	lastMsg = s.writeHandler.LastMessage()
	s.Contains(string(lastMsg), "\"specificationVersion\":\"1.1.0\"")
	
	// Test Reply
	requestHeader := &model.HeaderType{
		AddressSource: destAddr,
		AddressDestination: sourceAddr,
		MsgCounter: util.Ptr(model.MsgCounterType(1)),
	}
	
	s.writeHandler.sentMessages = nil
	err = sender.Reply(requestHeader, sourceAddr, model.CmdType{})
	s.NoError(err)
	lastMsg = s.writeHandler.LastMessage()
	s.Contains(string(lastMsg), "\"specificationVersion\":\"1.1.0\"")
	
	// Test Result
	s.writeHandler.sentMessages = nil
	err = sender.ResultSuccess(requestHeader, sourceAddr)
	s.NoError(err)
	lastMsg = s.writeHandler.LastMessage()
	s.Contains(string(lastMsg), "\"specificationVersion\":\"1.1.0\"")
}

// Test discovery messages use negotiated version
func (s *SendVersionSuite) Test_DiscoveryMessage_UsesNegotiatedVersion() {
	// With negotiated version, discovery should also use it
	s.remoteDevice.SetSupportedProtocolVersions([]string{"1.0.0"})
	s.remoteDevice.SetNegotiatedProtocolVersion("1.0.0")
	
	// Node management addresses (used for discovery)
	sourceAddr := &model.FeatureAddressType{
		Device:  util.Ptr(model.AddressDeviceType("local")),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(0)), // NodeManagement feature
	}
	
	destAddr := &model.FeatureAddressType{
		Device:  util.Ptr(model.AddressDeviceType("test-ski")),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(0)), // NodeManagement feature
	}
	
	// Send discovery request
	sender := s.remoteDevice.Sender()
	cmd := []model.CmdType{{
		NodeManagementDetailedDiscoveryData: &model.NodeManagementDetailedDiscoveryDataType{},
	}}
	
	s.writeHandler.sentMessages = nil
	_, err := sender.Request(model.CmdClassifierTypeRead, sourceAddr, destAddr, false, cmd)
	s.NoError(err)
	
	// Discovery messages should use negotiated version
	lastMsg := s.writeHandler.LastMessage()
	s.NotNil(lastMsg)
	s.Contains(string(lastMsg), "\"specificationVersion\":\"1.0.0\"")
}

// Test version detection affects outgoing messages
func (s *SendVersionSuite) Test_VersionDetection_AffectsOutgoing() {
	// Remote announced 1.2.0 and 1.3.0, we estimate 1.3.0
	s.remoteDevice.SetSupportedProtocolVersions([]string{"1.2.0", "1.3.0"})
	s.remoteDevice.UpdateEstimatedRemoteVersion()
	s.remoteDevice.SetNegotiatedProtocolVersion("1.3.0") // We can use 1.3.0
	
	// But remote actually uses 1.2.0 (detected)
	s.remoteDevice.UpdateDetectedVersion("1.2.0")
	
	sourceAddr := &model.FeatureAddressType{
		Device:  util.Ptr(model.AddressDeviceType("local")),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(1)),
	}
	
	destAddr := &model.FeatureAddressType{
		Device:  util.Ptr(model.AddressDeviceType("test-ski")),
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(1)),
	}
	
	// We should still use our negotiated version (1.3.0)
	sender := s.remoteDevice.Sender()
	s.writeHandler.sentMessages = nil
	_, err := sender.Request(model.CmdClassifierTypeRead, sourceAddr, destAddr, false, []model.CmdType{})
	s.NoError(err)
	
	lastMsg := s.writeHandler.LastMessage()
	s.NotNil(lastMsg)
	s.Contains(string(lastMsg), "\"specificationVersion\":\"1.3.0\"")
}