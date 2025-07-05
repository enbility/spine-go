package spine

import (
	"fmt"
	"testing"

	"github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestConnectionVersionTrackingSuite(t *testing.T) {
	suite.Run(t, new(ConnectionVersionTrackingSuite))
}

type ConnectionVersionTrackingSuite struct {
	suite.Suite

	localDevice  *DeviceLocal
	remoteDevice *DeviceRemote
	senderMock   *mocks.SenderInterface
}

func (s *ConnectionVersionTrackingSuite) WriteShipMessageWithPayload([]byte) {}

func (s *ConnectionVersionTrackingSuite) BeforeTest(suiteName, testName string) {
	s.localDevice = NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)

	ski := "test"
	s.senderMock = mocks.NewSenderInterface(s.T())
	
	// Set up mock expectations for all common calls that might happen during discovery
	s.senderMock.EXPECT().ProcessResponseForMsgCounterReference(mock.Anything).Maybe()
	s.senderMock.EXPECT().Request(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(util.Ptr(model.MsgCounterType(1)), nil).Maybe()
	s.senderMock.EXPECT().Reply(mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	s.senderMock.EXPECT().ResultSuccess(mock.Anything, mock.Anything).Return(nil).Maybe()
	s.senderMock.EXPECT().ResultError(mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	s.senderMock.EXPECT().Notify(mock.Anything, mock.Anything, mock.Anything).Return(util.Ptr(model.MsgCounterType(1)), nil).Maybe()
	s.senderMock.EXPECT().Write(mock.Anything, mock.Anything, mock.Anything).Return(util.Ptr(model.MsgCounterType(1)), nil).Maybe()
	s.senderMock.EXPECT().Subscribe(mock.Anything, mock.Anything, mock.Anything).Return(util.Ptr(model.MsgCounterType(1)), nil).Maybe()
	s.senderMock.EXPECT().Unsubscribe(mock.Anything, mock.Anything).Return(util.Ptr(model.MsgCounterType(1)), nil).Maybe()
	s.senderMock.EXPECT().Bind(mock.Anything, mock.Anything, mock.Anything).Return(util.Ptr(model.MsgCounterType(1)), nil).Maybe()
	s.senderMock.EXPECT().Unbind(mock.Anything, mock.Anything).Return(util.Ptr(model.MsgCounterType(1)), nil).Maybe()
	s.senderMock.EXPECT().DatagramForMsgCounter(mock.Anything).Return(model.DatagramType{}, nil).Maybe()
	
	s.remoteDevice = NewDeviceRemote(s.localDevice, ski, s.senderMock)
	desc := &model.NetworkManagementDeviceDescriptionDataType{
		DeviceAddress: &model.DeviceAddressType{
			Device: util.Ptr(model.AddressDeviceType("test")),
		},
	}
	s.remoteDevice.UpdateDevice(desc)
	_ = s.localDevice.SetupRemoteDevice(ski, s)
}

// Helper to create discovery data with version list
func (s *ConnectionVersionTrackingSuite) createDiscoveryDataWithVersions(versions []string) *model.NodeManagementDetailedDiscoveryDataType {
	versionList := make([]model.SpecificationVersionDataType, len(versions))
	for i, v := range versions {
		versionList[i] = model.SpecificationVersionDataType(v)
	}

	return &model.NodeManagementDetailedDiscoveryDataType{
		SpecificationVersionList: &model.NodeManagementSpecificationVersionListType{
			SpecificationVersion: versionList,
		},
		DeviceInformation: &model.NodeManagementDetailedDiscoveryDeviceInformationType{
			Description: &model.NetworkManagementDeviceDescriptionDataType{
				DeviceAddress: &model.DeviceAddressType{
					Device: util.Ptr(model.AddressDeviceType("remote")),
				},
				DeviceType:        util.Ptr(model.DeviceTypeTypeEnergyManagementSystem),
				NetworkFeatureSet: util.Ptr(model.NetworkManagementFeatureSetTypeSmart),
			},
		},
		EntityInformation: []model.NodeManagementDetailedDiscoveryEntityInformationType{},
		FeatureInformation: []model.NodeManagementDetailedDiscoveryFeatureInformationType{},
	}
}

// Helper to create a discovery reply message using JSON format like real test data
func (s *ConnectionVersionTrackingSuite) createDiscoveryReplyMessage(data *model.NodeManagementDetailedDiscoveryDataType) []byte {
	// Convert version list to strings for JSON
	var versionStrings []string
	if data.SpecificationVersionList != nil {
		for _, v := range data.SpecificationVersionList.SpecificationVersion {
			versionStrings = append(versionStrings, string(v))
		}
	}
	
	// Create JSON manually to match the real test data format
	messageJSON := fmt.Sprintf(`{
		"datagram": {
			"header": {
				"specificationVersion": "1.3.0",
				"addressSource": {
					"device": "test",
					"entity": [0],
					"feature": 0
				},
				"addressDestination": {
					"device": "%s",
					"entity": [0], 
					"feature": 0
				},
				"msgCounter": 100,
				"msgCounterReference": 99,
				"cmdClassifier": "reply"
			},
			"payload": {
				"cmd": [{
					"nodeManagementDetailedDiscoveryData": {
						"specificationVersionList": {
							"specificationVersion": %s
						},
						"deviceInformation": {
							"description": {
								"deviceAddress": {
									"device": "test"
								},
								"deviceType": "EnergyManagementSystem",
								"networkFeatureSet": "smart"
							}
						},
						"entityInformation": [],
						"featureInformation": []
					}
				}]
			}
		}
	}`, *s.localDevice.Address(), formatStringArray(versionStrings))
	
	return []byte(messageJSON)
}

// Helper to format string array for JSON
func formatStringArray(strings []string) string {
	if len(strings) == 0 {
		return "[]"
	}
	result := "["
	for i, s := range strings {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf(`"%s"`, s)
	}
	result += "]"
	return result
}

// Test storing remote device versions during discovery
func (s *ConnectionVersionTrackingSuite) Test_StoreRemoteVersionsDuringDiscovery() {
	// Test: Remote device supports versions 1.2.0 and 1.3.0
	discoveryData := s.createDiscoveryDataWithVersions([]string{"1.2.0", "1.3.0"})
	
	// Create a SPINE message with discovery reply
	message := s.createDiscoveryReplyMessage(discoveryData)
	
	// Process discovery reply
	_, err := s.remoteDevice.HandleSpineMesssage(message)
	assert.NoError(s.T(), err)
	
	// Verify versions are stored
	supportedVersions := s.remoteDevice.SupportedProtocolVersions()
	assert.NotNil(s.T(), supportedVersions)
	assert.Len(s.T(), supportedVersions, 2)
	assert.Contains(s.T(), supportedVersions, "1.2.0")
	assert.Contains(s.T(), supportedVersions, "1.3.0")
}

// Test finding compatible version
func (s *ConnectionVersionTrackingSuite) Test_FindCompatibleVersion() {
	// Test: Remote supports 1.1.0, 1.2.0; Local supports 1.3.0
	// Should negotiate to 1.3.0 (since all 1.x are compatible)
	discoveryData := s.createDiscoveryDataWithVersions([]string{"1.1.0", "1.2.0"})
	
	message := s.createDiscoveryReplyMessage(discoveryData)
	_, err := s.remoteDevice.HandleSpineMesssage(message)
	assert.NoError(s.T(), err)
	
	// Verify negotiated version
	negotiatedVersion := s.remoteDevice.NegotiatedProtocolVersion()
	assert.Equal(s.T(), "1.3.0", negotiatedVersion) // Local version wins within compatible major
}

// Test no compatible version - should send error
func (s *ConnectionVersionTrackingSuite) Test_NoCompatibleVersion_SendsError() {
	// Test: Remote only supports 2.0.0, Local supports 1.3.0
	discoveryData := s.createDiscoveryDataWithVersions([]string{"2.0.0"})
	
	// Use the common helper to create the message
	message := s.createDiscoveryReplyMessage(discoveryData)
	
	// Process discovery reply with incompatible version
	_, err := s.remoteDevice.HandleSpineMesssage(message)
	
	// Should return error for incompatible version
	assert.Error(s.T(), err)
	assert.True(s.T(), IsVersionIncompatibilityError(err), "Expected VersionIncompatibilityError but got: %T", err)
}

// Test empty version list - should assume current version
func (s *ConnectionVersionTrackingSuite) Test_EmptyVersionList_AssumesCurrentVersion() {
	// Test: Remote device doesn't provide version list
	discoveryData := &model.NodeManagementDetailedDiscoveryDataType{
		// No SpecificationVersionList
		DeviceInformation: &model.NodeManagementDetailedDiscoveryDeviceInformationType{
			Description: &model.NetworkManagementDeviceDescriptionDataType{
				DeviceAddress: &model.DeviceAddressType{
					Device: util.Ptr(model.AddressDeviceType("remote")),
				},
				DeviceType:        util.Ptr(model.DeviceTypeTypeEnergyManagementSystem),
				NetworkFeatureSet: util.Ptr(model.NetworkManagementFeatureSetTypeSmart),
			},
		},
		EntityInformation: []model.NodeManagementDetailedDiscoveryEntityInformationType{},
		FeatureInformation: []model.NodeManagementDetailedDiscoveryFeatureInformationType{},
	}
	
	message := s.createDiscoveryReplyMessage(discoveryData)
	_, err := s.remoteDevice.HandleSpineMesssage(message)
	assert.NoError(s.T(), err)
	
	// Should assume current version
	supportedVersions := s.remoteDevice.SupportedProtocolVersions()
	assert.Len(s.T(), supportedVersions, 1)
	assert.Equal(s.T(), "1.3.0", supportedVersions[0])
}

// Test version compatibility across minor versions
func (s *ConnectionVersionTrackingSuite) Test_VersionCompatibilityWithinMajor() {
	// Test: Different minor versions within same major version are compatible
	testCases := []struct {
		remoteVersions []string
		expectSuccess  bool
		expectedNegotiated string
	}{
		{[]string{"1.0.0", "1.1.0"}, true, "1.3.0"}, // Local has 1.3.0, all 1.x compatible
		{[]string{"1.2.0", "1.3.0", "1.4.0"}, true, "1.3.0"}, // Use exact match: 1.3.0
		{[]string{"0.9.0", "1.0.0"}, true, "1.3.0"}, // Major 0 and 1 both ok, use local version
		{[]string{"2.0.0", "2.1.0"}, false, ""}, // Major 2 incompatible
		{[]string{"1.3.0", "2.0.0"}, true, "1.3.0"}, // Has compatible option
	}
	
	for _, tc := range testCases {
		s.T().Run(tc.remoteVersions[0], func(t *testing.T) {
			// Reset device
			s.BeforeTest("", "")
			
			discoveryData := s.createDiscoveryDataWithVersions(tc.remoteVersions)
			message := s.createDiscoveryReplyMessage(discoveryData)
			_, err := s.remoteDevice.HandleSpineMesssage(message)
			
			if tc.expectSuccess {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedNegotiated, s.remoteDevice.NegotiatedProtocolVersion())
			} else {
				assert.Error(t, err)
				assert.True(t, IsVersionIncompatibilityError(err), "Expected VersionIncompatibilityError but got: %T", err)
			}
		})
	}
}

// Test version tracking persistence across messages
func (s *ConnectionVersionTrackingSuite) Test_VersionTrackingPersistence() {
	// Set up device with specific versions
	discoveryData := s.createDiscoveryDataWithVersions([]string{"1.2.0", "1.3.0"})
	message := s.createDiscoveryReplyMessage(discoveryData)
	_, err := s.remoteDevice.HandleSpineMesssage(message)
	assert.NoError(s.T(), err)
	
	// Process multiple messages - version should persist
	for i := 0; i < 3; i++ {
		versions := s.remoteDevice.SupportedProtocolVersions()
		assert.Len(s.T(), versions, 2)
		assert.Equal(s.T(), "1.3.0", s.remoteDevice.NegotiatedProtocolVersion())
	}
}

// Test version update on re-discovery
func (s *ConnectionVersionTrackingSuite) Test_VersionUpdateOnRediscovery() {
	// Initial discovery with version 1.2.0
	discoveryData := s.createDiscoveryDataWithVersions([]string{"1.2.0"})
	message := s.createDiscoveryReplyMessage(discoveryData)
	_, err := s.remoteDevice.HandleSpineMesssage(message)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "1.3.0", s.remoteDevice.NegotiatedProtocolVersion()) // Compatible so use local version
	
	// Re-discovery with updated versions
	discoveryData = s.createDiscoveryDataWithVersions([]string{"1.2.0", "1.3.0"})
	message = s.createDiscoveryReplyMessage(discoveryData)
	_, err = s.remoteDevice.HandleSpineMesssage(message)
	assert.NoError(s.T(), err)
	
	// Should use exact match when available
	assert.Equal(s.T(), "1.3.0", s.remoteDevice.NegotiatedProtocolVersion())
}