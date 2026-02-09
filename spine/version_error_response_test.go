package spine

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestVersionErrorResponseSuite(t *testing.T) {
	suite.Run(t, new(VersionErrorResponseSuite))
}

type VersionErrorResponseSuite struct {
	suite.Suite

	localDevice     *DeviceLocal
	remoteDevice    *DeviceRemote
	senderMock      *mocks.SenderInterface
	capturedResults []model.DatagramType // Track ResultError calls
}

func (s *VersionErrorResponseSuite) BeforeTest(suiteName, testName string) {
	s.capturedResults = []model.DatagramType{} // Reset captured error responses
	
	s.localDevice = NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)

	ski := "test"
	s.senderMock = mocks.NewSenderInterface(s.T())
	
	// Set up basic mock expectations for common calls
	s.senderMock.EXPECT().ProcessResponseForMsgCounterReference(mock.Anything).Maybe()
	s.senderMock.EXPECT().Request(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(util.Ptr(model.MsgCounterType(1)), nil).Maybe()
	s.senderMock.EXPECT().Reply(mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	s.senderMock.EXPECT().ResultSuccess(mock.Anything, mock.Anything).Return(nil).Maybe()
	s.senderMock.EXPECT().Notify(mock.Anything, mock.Anything, mock.Anything).Return(util.Ptr(model.MsgCounterType(1)), nil).Maybe()
	s.senderMock.EXPECT().Write(mock.Anything, mock.Anything, mock.Anything).Return(util.Ptr(model.MsgCounterType(1)), nil).Maybe()
	s.senderMock.EXPECT().Subscribe(mock.Anything, mock.Anything, mock.Anything).Return(util.Ptr(model.MsgCounterType(1)), nil).Maybe()
	s.senderMock.EXPECT().Unsubscribe(mock.Anything, mock.Anything).Return(util.Ptr(model.MsgCounterType(1)), nil).Maybe()
	s.senderMock.EXPECT().Bind(mock.Anything, mock.Anything, mock.Anything).Return(util.Ptr(model.MsgCounterType(1)), nil).Maybe()
	s.senderMock.EXPECT().Unbind(mock.Anything, mock.Anything).Return(util.Ptr(model.MsgCounterType(1)), nil).Maybe()
	s.senderMock.EXPECT().DatagramForMsgCounter(mock.Anything).Return(model.DatagramType{}, nil).Maybe()
	
	// Set up ResultError expectation with capture function
	s.senderMock.EXPECT().ResultError(mock.Anything, mock.Anything, mock.Anything).Maybe().Run(func(args mock.Arguments) {
		// Extract typed arguments
		requestHeader := args[0].(*model.HeaderType)
		senderAddress := args[1].(*model.FeatureAddressType)
		err := args[2].(*model.ErrorType)
		
		// Capture the error response for verification
		cmdClassifier := model.CmdClassifierTypeResult
		resultData := model.ResultDataType{
			ErrorNumber: &err.ErrorNumber,
			Description: err.Description,
		}
		cmd := model.CmdType{
			ResultData: &resultData,
		}
		datagram := model.DatagramType{
			Header: model.HeaderType{
				SpecificationVersion: &SpecificationVersion,
				AddressSource:        senderAddress,
				AddressDestination:   requestHeader.AddressSource,
				MsgCounter:           util.Ptr(model.MsgCounterType(1000)),
				MsgCounterReference:  requestHeader.MsgCounter,
				CmdClassifier:        &cmdClassifier,
			},
			Payload: model.PayloadType{
				Cmd: []model.CmdType{cmd},
			},
		}
		s.capturedResults = append(s.capturedResults, datagram)
	}).Return(nil)

	s.remoteDevice = NewDeviceRemote(s.localDevice, ski, s.senderMock)
	desc := &model.NetworkManagementDeviceDescriptionDataType{
		DeviceAddress: &model.DeviceAddressType{
			Device: util.Ptr(model.AddressDeviceType("test")),
		},
	}
	s.remoteDevice.UpdateDevice(desc)
	_ = s.localDevice.SetupRemoteDevice(ski, s)
}

func (s *VersionErrorResponseSuite) WriteShipMessageWithPayload([]byte) {}

func (s *VersionErrorResponseSuite) createTestMessage(version string, cmdClassifier model.CmdClassifierType, ackRequest *bool) []byte {
	versionType := model.SpecificationVersionType(version)
	msg := model.Datagram{
		Datagram: model.DatagramType{
			Header: model.HeaderType{
				SpecificationVersion: &versionType,
				AddressSource: &model.FeatureAddressType{
					Device:  util.Ptr(model.AddressDeviceType("test")),
					Entity:  []model.AddressEntityType{0},
					Feature: util.Ptr(model.AddressFeatureType(0)),
				},
				AddressDestination: &model.FeatureAddressType{
					Device:  util.Ptr(model.AddressDeviceType("local")),
					Entity:  []model.AddressEntityType{0},
					Feature: util.Ptr(model.AddressFeatureType(0)),
				},
				MsgCounter:    util.Ptr(model.MsgCounterType(1)),
				CmdClassifier: &cmdClassifier,
				AckRequest:    ackRequest,
			},
			Payload: model.PayloadType{
				Cmd: []model.CmdType{
					{
						Function: util.Ptr(model.FunctionTypeNodeManagementDetailedDiscoveryData),
					},
				},
			},
		},
	}

	data, _ := json.Marshal(msg)
	return data
}

// Test cases for version error responses

func (s *VersionErrorResponseSuite) Test_VersionError_ReadRequest_SendsErrorResponse() {
	// Test: READ request with incompatible version should receive error response
	msg := s.createTestMessage("2.0.0", model.CmdClassifierTypeRead, nil)
	
	msgCounter, err := s.remoteDevice.HandleSpineMesssage(msg)
	
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "incompatible major version")
	assert.Nil(s.T(), msgCounter)
	
	// Verify error response was captured
	assert.Len(s.T(), s.capturedResults, 1, "Should send error response for READ")
	response := s.capturedResults[0]
	
	assert.Equal(s.T(), model.CmdClassifierTypeResult, *response.Header.CmdClassifier)
	assert.NotNil(s.T(), response.Header.MsgCounterReference)
	assert.Equal(s.T(), model.MsgCounterType(1), *response.Header.MsgCounterReference)
	
	// Check error details
	assert.Len(s.T(), response.Payload.Cmd, 1)
	assert.NotNil(s.T(), response.Payload.Cmd[0].ResultData)
	assert.Equal(s.T(), model.ErrorNumberTypeGeneralError, *response.Payload.Cmd[0].ResultData.ErrorNumber)
	assert.Contains(s.T(), *response.Payload.Cmd[0].ResultData.Description, "version")
}

func (s *VersionErrorResponseSuite) Test_VersionError_WriteRequest_SendsErrorResponse() {
	// Test: WRITE request with incompatible version should receive error response
	msg := s.createTestMessage("3.0.0", model.CmdClassifierTypeWrite, nil)
	
	msgCounter, err := s.remoteDevice.HandleSpineMesssage(msg)
	
	assert.Error(s.T(), err)
	assert.Nil(s.T(), msgCounter)
	
	// Verify error response was captured
	assert.Len(s.T(), s.capturedResults, 1, "Should send error response for WRITE")
	response := s.capturedResults[0]
	
	assert.Equal(s.T(), model.CmdClassifierTypeResult, *response.Header.CmdClassifier)
	assert.NotNil(s.T(), response.Payload.Cmd[0].ResultData)
	assert.Equal(s.T(), model.ErrorNumberTypeGeneralError, *response.Payload.Cmd[0].ResultData.ErrorNumber)
}

func (s *VersionErrorResponseSuite) Test_VersionError_NotifyWithAckRequest_SendsErrorResponse() {
	// Test: NOTIFY with ackRequest=true should receive error response
	ackRequest := true
	msg := s.createTestMessage("2.0.0", model.CmdClassifierTypeNotify, &ackRequest)
	
	msgCounter, err := s.remoteDevice.HandleSpineMesssage(msg)
	
	assert.Error(s.T(), err)
	assert.Nil(s.T(), msgCounter)
	
	// Verify error response was captured
	assert.Len(s.T(), s.capturedResults, 1, "Should send error response for NOTIFY with ackRequest")
	response := s.capturedResults[0]
	
	assert.Equal(s.T(), model.CmdClassifierTypeResult, *response.Header.CmdClassifier)
}

func (s *VersionErrorResponseSuite) Test_VersionError_NotifyWithoutAckRequest_NoErrorResponse() {
	// Test: NOTIFY without ackRequest should NOT receive error response
	ackRequest := false
	msg := s.createTestMessage("2.0.0", model.CmdClassifierTypeNotify, &ackRequest)
	
	msgCounter, err := s.remoteDevice.HandleSpineMesssage(msg)
	
	assert.Error(s.T(), err)
	assert.Nil(s.T(), msgCounter)
	
	// Verify NO error response was captured
	assert.Len(s.T(), s.capturedResults, 0, "Should NOT send error response for NOTIFY without ackRequest")
}

func (s *VersionErrorResponseSuite) Test_VersionError_ReplyMessage_NoErrorResponse() {
	// Test: REPLY messages should NOT receive error response (avoid loops)
	msg := s.createTestMessage("2.0.0", model.CmdClassifierTypeReply, nil)
	
	msgCounter, err := s.remoteDevice.HandleSpineMesssage(msg)
	
	assert.Error(s.T(), err)
	assert.Nil(s.T(), msgCounter)
	
	// Verify NO error response was captured
	assert.Len(s.T(), s.capturedResults, 0, "Should NOT send error response for REPLY")
}

func (s *VersionErrorResponseSuite) Test_VersionError_ResultMessage_NoErrorResponse() {
	// Test: RESULT messages should NOT receive error response (avoid loops)
	msg := s.createTestMessage("2.0.0", model.CmdClassifierTypeResult, nil)
	
	msgCounter, err := s.remoteDevice.HandleSpineMesssage(msg)
	
	assert.Error(s.T(), err)
	assert.Nil(s.T(), msgCounter)
	
	// Verify NO error response was captured
	assert.Len(s.T(), s.capturedResults, 0, "Should NOT send error response for RESULT")
}

func (s *VersionErrorResponseSuite) Test_CompatibleVersion_NoErrorResponse() {
	// Test: Compatible version should NOT send error response
	msg := s.createTestMessage("1.3.0", model.CmdClassifierTypeRead, nil)
	
	msgCounter, err := s.remoteDevice.HandleSpineMesssage(msg)
	
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), msgCounter)
	
	// Verify NO version-related error response was captured
	// We specifically check for errors containing "version" in the description
	versionErrorFound := false
	for _, response := range s.capturedResults {
		if len(response.Payload.Cmd) > 0 && 
		   response.Payload.Cmd[0].ResultData != nil && 
		   response.Payload.Cmd[0].ResultData.ErrorNumber != nil &&
		   *response.Payload.Cmd[0].ResultData.ErrorNumber != model.ErrorNumberTypeNoError {
			// Check if this is a version-related error
			if response.Payload.Cmd[0].ResultData.Description != nil &&
			   strings.Contains(string(*response.Payload.Cmd[0].ResultData.Description), "version") {
				versionErrorFound = true
				break
			}
		}
	}
	assert.False(s.T(), versionErrorFound, "Should NOT send version error response for compatible version")
}