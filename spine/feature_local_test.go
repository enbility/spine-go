package spine

import (
	"errors"
	"testing"
	"time"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestDeviceClassificationSuite(t *testing.T) {
	suite.Run(t, new(DeviceClassificationTestSuite))
}

type DeviceClassificationTestSuite struct {
	suite.Suite
	senderMock                                                *mocks.SenderInterface
	localDevice                                               *DeviceLocal
	localEntity                                               *EntityLocal
	function, serverWriteFunction                             model.FunctionType
	featureType, subFeatureType                               model.FeatureTypeType
	msgCounter                                                model.MsgCounterType
	remoteFeature, remoteServerFeature, remoteSubFeature      api.FeatureRemoteInterface
	localFeature, localServerFeature, localServerFeatureWrite api.FeatureLocalInterface
}

func (suite *DeviceClassificationTestSuite) BeforeTest(suiteName, testName string) {
	suite.senderMock = mocks.NewSenderInterface(suite.T())
	suite.function = model.FunctionTypeDeviceClassificationManufacturerData
	suite.featureType = model.FeatureTypeTypeDeviceClassification
	suite.subFeatureType = model.FeatureTypeTypeLoadControl
	suite.serverWriteFunction = model.FunctionTypeLoadControlLimitListData
	suite.msgCounter = model.MsgCounterType(1)

	suite.localDevice, suite.localEntity = createLocalDeviceAndEntity(1)
	suite.localFeature, suite.localServerFeature = createLocalFeatures(suite.localEntity, suite.featureType, "")
	_, suite.localServerFeatureWrite = createLocalFeatures(suite.localEntity, suite.subFeatureType, suite.serverWriteFunction)

	remoteDevice := createRemoteDevice(suite.localDevice, suite.senderMock)
	suite.remoteFeature, suite.remoteServerFeature = createRemoteEntityAndFeature(remoteDevice, 1, suite.featureType)
	suite.remoteSubFeature, _ = createRemoteEntityAndFeature(remoteDevice, 2, suite.subFeatureType)
}

func (suite *DeviceClassificationTestSuite) TestDeviceClassification_Functions() {
	fcts := suite.localServerFeatureWrite.Functions()
	assert.NotNil(suite.T(), fcts)
	assert.Equal(suite.T(), 1, len(fcts))
	assert.Equal(suite.T(), suite.serverWriteFunction, fcts[0])
}

func (suite *DeviceClassificationTestSuite) TestDeviceClassification_ResponseCB() {
	testFct := func(msg api.ResponseMessage) {}
	msgCounter := model.MsgCounterType(100)
	err := suite.localFeature.AddResponseCallback(msgCounter, testFct)
	assert.Nil(suite.T(), err)

	err = suite.localFeature.AddResponseCallback(msgCounter, testFct)
	assert.NotNil(suite.T(), err)

	suite.localFeature.AddResultCallback(testFct)
}

func (suite *DeviceClassificationTestSuite) TestDeviceClassification_Request_Reply() {
	dummyAddress := &model.FeatureAddressType{
		Device:  util.Ptr(model.AddressDeviceType("")),
		Entity:  []model.AddressEntityType{2},
		Feature: util.Ptr(model.AddressFeatureType(100)),
	}

	suite.senderMock.EXPECT().Request(
		model.CmdClassifierTypeRead,
		suite.localFeature.Address(),
		suite.remoteFeature.Address(),
		false,
		mock.AnythingOfType("[]model.CmdType")).Return(&suite.msgCounter, nil).Maybe()
	suite.senderMock.EXPECT().Request(
		model.CmdClassifierTypeRead,
		suite.localFeature.Address(),
		dummyAddress,
		false,
		mock.AnythingOfType("[]model.CmdType")).Return(nil, errors.New("test"))

	msgCounter, err := suite.localFeature.RequestRemoteData(model.FunctionTypeBillListData, nil, nil, suite.remoteFeature)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), msgCounter)

	msgCounter, err = suite.localFeature.RequestRemoteDataBySenderAddress(model.CmdType{}, suite.senderMock, "dummyfail", dummyAddress, 0)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), msgCounter)

	// send data request
	msgCounter, err = suite.localFeature.RequestRemoteData(suite.function, nil, nil, suite.remoteFeature)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), msgCounter)

	manufacturerData := &model.DeviceClassificationManufacturerDataType{
		BrandName:    util.Ptr(model.DeviceClassificationStringType("brand name")),
		VendorName:   util.Ptr(model.DeviceClassificationStringType("vendor name")),
		DeviceName:   util.Ptr(model.DeviceClassificationStringType("device name")),
		DeviceCode:   util.Ptr(model.DeviceClassificationStringType("device code")),
		SerialNumber: util.Ptr(model.DeviceClassificationStringType("serial number")),
	}

	replyMsg := api.Message{
		Cmd: model.CmdType{
			DeviceClassificationManufacturerData: manufacturerData,
		},
		CmdClassifier: model.CmdClassifierTypeReply,
		RequestHeader: &model.HeaderType{
			MsgCounter:          util.Ptr(model.MsgCounterType(1)),
			MsgCounterReference: &suite.msgCounter,
		},
		FeatureRemote: suite.remoteFeature,
	}
	// set response
	msgErr := suite.localFeature.HandleMessage(&replyMsg)
	if assert.Nil(suite.T(), msgErr) {
		remoteData := suite.remoteFeature.DataCopy(suite.function)
		assert.IsType(suite.T(), &model.DeviceClassificationManufacturerDataType{}, remoteData, "Data has wrong type")
	}
}

func (suite *DeviceClassificationTestSuite) TestDeviceClassification_Request_Error() {
	suite.senderMock.On("Request", model.CmdClassifierTypeRead, suite.localFeature.Address(), suite.remoteFeature.Address(), false, mock.AnythingOfType("[]model.CmdType")).Return(&suite.msgCounter, nil)

	const errorNumber = model.ErrorNumberTypeGeneralError
	const errorDescription = "error occurred"

	// send data request
	msgCounter, err := suite.localFeature.RequestRemoteData(suite.function, nil, nil, suite.remoteFeature)
	assert.NotNil(suite.T(), msgCounter)
	assert.Nil(suite.T(), err)

	replyMsg := api.Message{
		Cmd: model.CmdType{
			ResultData: &model.ResultDataType{
				ErrorNumber: util.Ptr(model.ErrorNumberType(errorNumber)),
				Description: util.Ptr(model.DescriptionType(errorDescription)),
			},
		},
		CmdClassifier: model.CmdClassifierTypeResult,
		RequestHeader: &model.HeaderType{
			MsgCounter:          util.Ptr(model.MsgCounterType(1)),
			MsgCounterReference: &suite.msgCounter,
		},
		FeatureRemote: suite.remoteFeature,
		EntityRemote:  suite.remoteFeature.Entity(),
		DeviceRemote:  suite.remoteFeature.Device(),
	}

	// set response
	msgErr := suite.localFeature.HandleMessage(&replyMsg)
	if assert.Nil(suite.T(), msgErr) {
		remoteData := suite.remoteFeature.DataCopy(suite.function)
		assert.Nil(suite.T(), remoteData)
	}
}

func (suite *DeviceClassificationTestSuite) TestDeviceClassification_Subscriptions() {
	suite.senderMock.On("Subscribe", mock.Anything, mock.Anything, mock.Anything).Return(&suite.msgCounter, nil)
	suite.senderMock.On("Unsubscribe", mock.Anything, mock.Anything, mock.Anything).Return(&suite.msgCounter, nil)

	msgCounter, err := suite.localFeature.SubscribeToRemote(suite.remoteFeature.Address())
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), msgCounter)

	msgCounter, err = suite.localFeature.RemoveRemoteSubscription(suite.remoteFeature.Address())
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), msgCounter)

	suite.localFeature.Device().AddRemoteDeviceForSki(suite.remoteFeature.Device().Ski(), suite.remoteFeature.Device())

	msgCounter, err = suite.localServerFeature.SubscribeToRemote(suite.remoteFeature.Address())
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), msgCounter)

	msgCounter, err = suite.localFeature.RemoveRemoteSubscription(suite.remoteFeature.Address())
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), msgCounter)

	subscribed := suite.localFeature.HasSubscriptionToRemote(suite.remoteFeature.Address())
	assert.Equal(suite.T(), false, subscribed)

	msgCounter, err = suite.localFeature.SubscribeToRemote(suite.remoteFeature.Address())
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), msgCounter)

	subscribed = suite.localFeature.HasSubscriptionToRemote(suite.remoteFeature.Address())
	assert.Equal(suite.T(), true, subscribed)

	msgCounter, err = suite.localFeature.SubscribeToRemote(suite.remoteSubFeature.Address())
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), msgCounter)

	msgCounter, err = suite.localFeature.RemoveRemoteSubscription(suite.remoteFeature.Address())
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), msgCounter)

	suite.localFeature.RemoveAllRemoteSubscriptions()
}

func (suite *DeviceClassificationTestSuite) TestDeviceClassification_Bindings() {
	suite.senderMock.On("Bind", mock.Anything, mock.Anything, mock.Anything).Return(&suite.msgCounter, nil)
	suite.senderMock.On("Unbind", mock.Anything, mock.Anything, mock.Anything).Return(&suite.msgCounter, nil)

	msgCounter, err := suite.localFeature.BindToRemote(suite.remoteFeature.Address())
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), msgCounter)

	msgCounter, err = suite.localFeature.RemoveRemoteBinding(suite.remoteFeature.Address())
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), msgCounter)

	suite.localFeature.Device().AddRemoteDeviceForSki(suite.remoteFeature.Device().Ski(), suite.remoteFeature.Device())

	msgCounter, err = suite.localServerFeature.BindToRemote(suite.remoteFeature.Address())
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), msgCounter)

	msgCounter, err = suite.localFeature.RemoveRemoteBinding(suite.remoteFeature.Address())
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), msgCounter)

	binding := suite.localFeature.HasBindingToRemote(suite.remoteFeature.Address())
	assert.Equal(suite.T(), false, binding)

	msgCounter, err = suite.localFeature.BindToRemote(suite.remoteFeature.Address())
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), msgCounter)

	binding = suite.localFeature.HasBindingToRemote(suite.remoteFeature.Address())
	assert.Equal(suite.T(), true, binding)

	msgCounter, err = suite.localFeature.BindToRemote(suite.remoteSubFeature.Address())
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), msgCounter)

	msgCounter, err = suite.localFeature.RemoveRemoteBinding(suite.remoteFeature.Address())
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), msgCounter)

	suite.localFeature.RemoveAllRemoteBindings()
}

func (suite *DeviceClassificationTestSuite) Test_HandleMessage() {
	msg := &api.Message{
		FeatureRemote: suite.remoteServerFeature,
		CmdClassifier: model.CmdClassifierType("buggy"),
		Cmd:           model.CmdType{},
	}

	err := suite.localFeature.HandleMessage(msg)
	assert.NotNil(suite.T(), err)

	msg = &api.Message{
		FeatureRemote: suite.remoteServerFeature,
		CmdClassifier: model.CmdClassifierType("buggy"),
		Cmd: model.CmdType{
			ResultData: &model.ResultDataType{},
		},
	}

	err = suite.localFeature.HandleMessage(msg)
	assert.NotNil(suite.T(), err)
}

func (suite *DeviceClassificationTestSuite) Test_Result() {
	msg := &api.Message{
		FeatureRemote: suite.remoteServerFeature,
		CmdClassifier: model.CmdClassifierTypeResult,
		Cmd: model.CmdType{
			ResultData: &model.ResultDataType{},
		},
	}

	err := suite.localFeature.HandleMessage(msg)
	assert.NotNil(suite.T(), err)

	msg.RequestHeader = &model.HeaderType{
		MsgCounterReference: util.Ptr(model.MsgCounterType(100)),
	}
	msg.Cmd.ResultData = &model.ResultDataType{
		ErrorNumber: util.Ptr(model.ErrorNumberType(1)),
		Description: util.Ptr(model.DescriptionType("test")),
	}
	err = suite.localFeature.HandleMessage(msg)
	assert.Nil(suite.T(), err)
}

func (suite *DeviceClassificationTestSuite) Test_Read() {
	msg := &api.Message{
		FeatureRemote: suite.remoteFeature,
		CmdClassifier: model.CmdClassifierTypeRead,
		Cmd: model.CmdType{
			LoadControlLimitListData: &model.LoadControlLimitListDataType{},
		},
	}

	err := suite.localFeature.HandleMessage(msg)
	assert.NotNil(suite.T(), err)

	suite.senderMock.EXPECT().Reply(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test")).Once()
	err = suite.localServerFeature.HandleMessage(msg)
	assert.NotNil(suite.T(), err)

	msg = &api.Message{
		FeatureRemote: suite.remoteFeature,
		CmdClassifier: model.CmdClassifierTypeRead,
		Cmd: model.CmdType{
			DeviceClassificationManufacturerData: &model.DeviceClassificationManufacturerDataType{},
		},
	}
	err = suite.localFeature.HandleMessage(msg)
	assert.NotNil(suite.T(), err)

	suite.senderMock.EXPECT().Reply(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
	err = suite.localServerFeature.HandleMessage(msg)
	assert.NotNil(suite.T(), err)
}

func (suite *DeviceClassificationTestSuite) Test_Reply() {
	msg := &api.Message{
		FeatureRemote: suite.remoteServerFeature,
		CmdClassifier: model.CmdClassifierTypeReply,
		Cmd: model.CmdType{
			DeviceClassificationManufacturerData: &model.DeviceClassificationManufacturerDataType{},
		},
	}

	err := suite.localFeature.HandleMessage(msg)
	assert.Nil(suite.T(), err)

	msg.RequestHeader = &model.HeaderType{
		MsgCounterReference: util.Ptr(model.MsgCounterType(100)),
	}
	err = suite.localFeature.HandleMessage(msg)
	assert.Nil(suite.T(), err)
}

func (suite *DeviceClassificationTestSuite) Test_Notify() {
	msg := &api.Message{
		FeatureRemote: suite.remoteServerFeature,
		CmdClassifier: model.CmdClassifierTypeNotify,
		Cmd: model.CmdType{
			DeviceClassificationManufacturerData: &model.DeviceClassificationManufacturerDataType{},
		},
	}

	err := suite.localFeature.HandleMessage(msg)
	assert.Nil(suite.T(), err)
}

func (suite *DeviceClassificationTestSuite) Test_Write() {
	suite.senderMock.EXPECT().ResultSuccess(mock.Anything, mock.Anything).Return(nil).Once()

	msg := &api.Message{
		RequestHeader: &model.HeaderType{
			MsgCounter: util.Ptr(model.MsgCounterType(1)),
			AckRequest: util.Ptr(true),
		},
		CmdClassifier: model.CmdClassifierTypeWrite,
		FeatureRemote: suite.remoteSubFeature,
		Cmd: model.CmdType{
			LoadControlLimitListData: &model.LoadControlLimitListDataType{},
		},
	}

	err := suite.localServerFeatureWrite.HandleMessage(msg)
	assert.Nil(suite.T(), err)
}

func (suite *DeviceClassificationTestSuite) Test_SetWriteApprovalCallback_Invalid() {
	cb := func(msg *api.Message) {}
	err := suite.localFeature.SetWriteApprovalCallback(cb)
	assert.NotNil(suite.T(), err)
	suite.localFeature.ApproveOrDenyWrite(&api.Message{}, true, "")
}

func (suite *DeviceClassificationTestSuite) Test_AddPendingApproval_Invalid() {
	cb := func(msg *api.Message) {}
	err := suite.localServerFeatureWrite.SetWriteApprovalCallback(cb)
	assert.Nil(suite.T(), err)

	msg := &api.Message{
		CmdClassifier: model.CmdClassifierTypeWrite,
		FeatureRemote: suite.remoteSubFeature,
		Cmd: model.CmdType{
			LoadControlLimitListData: &model.LoadControlLimitListDataType{},
		},
	}

	err1 := suite.localServerFeatureWrite.HandleMessage(msg)
	assert.Nil(suite.T(), err1)
}

func (suite *DeviceClassificationTestSuite) Test_Write_Callback() {
	msg := &api.Message{
		RequestHeader: &model.HeaderType{
			MsgCounter: util.Ptr(model.MsgCounterType(1)),
		},
		CmdClassifier: model.CmdClassifierTypeWrite,
		FeatureRemote: suite.remoteSubFeature,
		Cmd: model.CmdType{
			LoadControlLimitListData: &model.LoadControlLimitListDataType{},
		},
	}

	cb := func(msg *api.Message) {
		suite.localServerFeatureWrite.ApproveOrDenyWrite(msg, true, "")
	}

	suite.localServerFeatureWrite.SetWriteApprovalCallback(cb)
	err := suite.localServerFeatureWrite.HandleMessage(msg)
	assert.Nil(suite.T(), err)

	suite.senderMock.EXPECT().ResultError(mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	cb = func(msg *api.Message) {
		suite.localServerFeatureWrite.ApproveOrDenyWrite(msg, false, "not allowed by application")
	}

	suite.localServerFeatureWrite.SetWriteApprovalCallback(cb)
	err = suite.localServerFeatureWrite.HandleMessage(msg)
	assert.Nil(suite.T(), err)
}

func (suite *DeviceClassificationTestSuite) Test_Write_Callback_Timeout() {
	suite.localServerFeatureWrite.SetWriteApprovalTimeout(time.Second * 1)

	suite.senderMock.EXPECT().ResultError(mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	msg := &api.Message{
		RequestHeader: &model.HeaderType{
			MsgCounter: util.Ptr(model.MsgCounterType(1)),
		},
		CmdClassifier: model.CmdClassifierTypeWrite,
		FeatureRemote: suite.remoteSubFeature,
		Cmd: model.CmdType{
			LoadControlLimitListData: &model.LoadControlLimitListDataType{},
		},
	}

	cb := func(msg *api.Message) {
		time.Sleep(time.Second * 2)
		suite.localServerFeatureWrite.ApproveOrDenyWrite(msg, true, "")
	}

	suite.localServerFeatureWrite.SetWriteApprovalCallback(cb)
	err := suite.localServerFeatureWrite.HandleMessage(msg)
	assert.Nil(suite.T(), err)
}
