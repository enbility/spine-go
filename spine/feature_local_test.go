package spine

import (
	"errors"
	"reflect"
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

func TestLocalFeatureTestSuite(t *testing.T) {
	suite.Run(t, new(LocalFeatureTestSuite))
}

type LocalFeatureTestSuite struct {
	suite.Suite
	senderMock                    *mocks.SenderInterface
	localDevice                   *DeviceLocal
	localEntity                   *EntityLocal
	function, serverWriteFunction model.FunctionType
	featureType, subFeatureType   model.FeatureTypeType
	msgCounter                    model.MsgCounterType
	remoteFeature, remoteServerFeature,
	remote2Feature, remote2ServerFeature,
	remoteSubFeature, remoteServerSubFeature api.FeatureRemoteInterface
	localFeature, localServerFeature, localServerFeatureWrite api.FeatureLocalInterface
}

func (s *LocalFeatureTestSuite) BeforeTest(suiteName, testName string) {
	s.senderMock = mocks.NewSenderInterface(s.T())
	s.function = model.FunctionTypeDeviceClassificationManufacturerData
	s.featureType = model.FeatureTypeTypeDeviceClassification
	s.subFeatureType = model.FeatureTypeTypeLoadControl
	s.serverWriteFunction = model.FunctionTypeLoadControlLimitListData
	s.msgCounter = model.MsgCounterType(1)

	s.localDevice, s.localEntity = createLocalDeviceAndEntity(1)
	s.localFeature, s.localServerFeature = createLocalFeatures(s.localEntity, s.featureType, "")
	_, s.localServerFeatureWrite = createLocalFeatures(s.localEntity, s.subFeatureType, s.serverWriteFunction)

	remoteDevice := createRemoteDevice(s.localDevice, "ski", s.senderMock)
	remoteDevice2 := createRemoteDevice(s.localDevice, "iks", s.senderMock)
	s.remoteFeature, s.remoteServerFeature = createRemoteEntityAndFeature(remoteDevice, 1, s.featureType, s.function)
	s.remoteSubFeature, s.remoteServerSubFeature = createRemoteEntityAndFeature(remoteDevice, 2, s.subFeatureType, s.serverWriteFunction)
	s.remote2Feature, s.remote2ServerFeature = createRemoteEntityAndFeature(remoteDevice2, 1, s.featureType, s.function)
}

func (s *LocalFeatureTestSuite) TestDeviceClassification_Functions() {
	fcts := s.localServerFeatureWrite.Functions()
	assert.NotNil(s.T(), fcts)
	assert.Equal(s.T(), 1, len(fcts))
	assert.Equal(s.T(), s.serverWriteFunction, fcts[0])
}

func (s *LocalFeatureTestSuite) TestDeviceClassification_ResponseCB() {
	testFct := func(msg api.ResponseMessage) {}
	msgCounter := model.MsgCounterType(100)
	err := s.localFeature.AddResponseCallback(msgCounter, testFct)
	assert.Nil(s.T(), err)

	err = s.localFeature.AddResponseCallback(msgCounter, testFct)
	assert.NotNil(s.T(), err)

	s.localFeature.AddResultCallback(testFct)
}

func (s *LocalFeatureTestSuite) TestDeviceClassification_Request_Reply() {
	dummyAddress := &model.FeatureAddressType{
		Device:  util.Ptr(model.AddressDeviceType("")),
		Entity:  []model.AddressEntityType{2},
		Feature: util.Ptr(model.AddressFeatureType(100)),
	}

	s.senderMock.EXPECT().Request(
		model.CmdClassifierTypeRead,
		s.localFeature.Address(),
		s.remoteFeature.Address(),
		false,
		mock.AnythingOfType("[]model.CmdType")).Return(&s.msgCounter, nil).Maybe()
	s.senderMock.EXPECT().Request(
		model.CmdClassifierTypeRead,
		s.localFeature.Address(),
		dummyAddress,
		false,
		mock.AnythingOfType("[]model.CmdType")).Return(nil, errors.New("test"))

	msgCounter, err := s.localFeature.RequestRemoteData(model.FunctionTypeBillListData, nil, nil, s.remoteFeature)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	msgCounter, err = s.localFeature.RequestRemoteDataBySenderAddress(model.CmdType{}, s.senderMock, "dummyfail", dummyAddress, 0)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	// send data request
	msgCounter, err = s.localFeature.RequestRemoteData(s.function, nil, nil, s.remoteFeature)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

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
			MsgCounterReference: &s.msgCounter,
		},
		FeatureRemote: s.remoteFeature,
	}
	// set response
	msgErr := s.localFeature.HandleMessage(&replyMsg)
	if assert.Nil(s.T(), msgErr) {
		remoteData := s.remoteFeature.DataCopy(s.function)
		assert.IsType(s.T(), &model.DeviceClassificationManufacturerDataType{}, remoteData, "Data has wrong type")
	}
}

func (s *LocalFeatureTestSuite) TestDeviceClassification_Request_Error() {
	s.senderMock.On("Request", model.CmdClassifierTypeRead, s.localFeature.Address(), s.remoteFeature.Address(), false, mock.AnythingOfType("[]model.CmdType")).Return(&s.msgCounter, nil)

	const errorNumber = model.ErrorNumberTypeGeneralError
	const errorDescription = "error occurred"

	// send data request
	msgCounter, err := s.localFeature.RequestRemoteData(s.function, nil, nil, s.remoteFeature)
	assert.NotNil(s.T(), msgCounter)
	assert.Nil(s.T(), err)

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
			MsgCounterReference: &s.msgCounter,
		},
		FeatureRemote: s.remoteFeature,
		EntityRemote:  s.remoteFeature.Entity(),
		DeviceRemote:  s.remoteFeature.Device(),
	}

	// set response
	msgErr := s.localFeature.HandleMessage(&replyMsg)
	if assert.Nil(s.T(), msgErr) {
		remoteData := s.remoteFeature.DataCopy(s.function)
		assert.Nil(s.T(), remoteData)
	}
}

func (s *LocalFeatureTestSuite) TestDeviceClassification_Subscriptions() {
	s.senderMock.On("Subscribe", mock.Anything, mock.Anything, mock.Anything).Return(&s.msgCounter, nil)
	s.senderMock.On("Unsubscribe", mock.Anything, mock.Anything, mock.Anything).Return(&s.msgCounter, nil)

	msgCounter, err := s.localFeature.SubscribeToRemote(s.remoteFeature.Address())
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	msgCounter, err = s.localFeature.RemoveRemoteSubscription(s.remoteFeature.Address())
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	s.localFeature.Device().AddRemoteDeviceForSki(s.remoteFeature.Device().Ski(), s.remoteFeature.Device())
	s.localFeature.Device().AddRemoteDeviceForSki(s.remote2Feature.Device().Ski(), s.remote2Feature.Device())

	msgCounter, err = s.localServerFeature.SubscribeToRemote(s.remoteFeature.Address())
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	msgCounter, err = s.localFeature.RemoveRemoteSubscription(s.remoteFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	subscribed := s.localFeature.HasSubscriptionToRemote(s.remoteFeature.Address())
	assert.Equal(s.T(), false, subscribed)

	msgCounter, err = s.localFeature.SubscribeToRemote(s.remoteServerFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	subscribed = s.localFeature.HasSubscriptionToRemote(s.remoteServerFeature.Address())
	assert.False(s.T(), subscribed)

	lf := s.localFeature.(*FeatureLocal)
	msg := s.responseMsg(s.localFeature, s.remoteServerFeature, *msgCounter, 0)
	lf.subscribeResponseCallback(s.remoteServerFeature.Device(), s.remoteServerFeature.Address(), s.remoteServerFeature.Type(), msg)

	subscribed = s.localFeature.HasSubscriptionToRemote(s.remoteServerFeature.Address())
	assert.True(s.T(), subscribed)

	msgCounter, err = s.localFeature.SubscribeToRemote(s.remoteServerFeature.Address())
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	msgCounter, err = s.localFeature.SubscribeToRemote(s.remote2ServerFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	subscribed = s.localFeature.HasSubscriptionToRemote(s.remote2ServerFeature.Address())
	assert.False(s.T(), subscribed)

	msg = s.responseMsg(s.localFeature, s.remote2ServerFeature, *msgCounter, 0)
	lf.subscribeResponseCallback(s.remote2ServerFeature.Device(), s.remote2ServerFeature.Address(), s.remote2ServerFeature.Type(), msg)

	subscribed = s.localFeature.HasSubscriptionToRemote(s.remote2ServerFeature.Address())
	assert.True(s.T(), subscribed)

	msgCounter, err = s.localFeature.RemoveRemoteSubscription(s.remoteServerFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	msg = s.responseMsg(s.localFeature, s.remoteServerFeature, *msgCounter, 0)
	lf.unsubscribeResponseCallback(s.remoteServerFeature.Device(), s.remoteServerFeature.Address(), msg)

	subscribed = s.localFeature.HasSubscriptionToRemote(s.remoteServerFeature.Address())
	assert.False(s.T(), subscribed)

	subscribed = s.localFeature.HasSubscriptionToRemote(s.remote2ServerFeature.Address())
	assert.True(s.T(), subscribed)

	subscriptionAdd := model.SubscriptionManagementRequestCallType{
		ClientAddress: s.remoteFeature.Address(),
		ServerAddress: s.localServerFeature.Address(),
	}
	s.localDevice.SubscriptionManager().AddSubscription(s.remoteFeature.Device(), subscriptionAdd)

	subscribed = s.localServerFeature.HasSubscriptionToRemote(s.remoteFeature.Address())
	assert.True(s.T(), subscribed)

	msgCounter, err = s.localServerFeature.RemoveRemoteSubscription(s.remoteFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	lf = s.localServerFeature.(*FeatureLocal)
	msg = s.responseMsg(s.localServerFeature, s.remoteFeature, *msgCounter, 0)
	lf.unsubscribeResponseCallback(s.remoteFeature.Device(), s.remoteFeature.Address(), msg)

	subscribed = s.localServerFeature.HasSubscriptionToRemote(s.remoteFeature.Address())
	assert.False(s.T(), subscribed)
}

func (s *LocalFeatureTestSuite) TestDeviceClassification_Bindings() {
	s.senderMock.On("Bind", mock.Anything, mock.Anything, mock.Anything).Return(&s.msgCounter, nil)
	s.senderMock.On("Unbind", mock.Anything, mock.Anything, mock.Anything).Return(&s.msgCounter, nil)

	msgCounter, err := s.localFeature.BindToRemote(s.remoteFeature.Address())
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	msgCounter, err = s.localFeature.RemoveRemoteBinding(s.remoteFeature.Address())
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	s.localFeature.Device().AddRemoteDeviceForSki(s.remoteServerFeature.Device().Ski(), s.remoteServerFeature.Device())

	msgCounter, err = s.localServerFeature.BindToRemote(s.remoteServerFeature.Address())
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	msgCounter, err = s.localFeature.RemoveRemoteBinding(s.remoteServerFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	binding := s.localFeature.HasBindingToRemote(s.remoteServerFeature.Address())
	assert.False(s.T(), binding)

	msgCounter, err = s.localFeature.BindToRemote(s.remoteServerFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	lf := s.localFeature.(*FeatureLocal)
	msg := s.responseMsg(s.localFeature, s.remoteServerFeature, *msgCounter, 0)
	lf.bindResponseCallback(s.remoteServerFeature.Device(), s.remoteServerFeature.Address(), s.remoteServerFeature.Type(), msg)

	binding = s.localFeature.HasBindingToRemote(s.remoteServerFeature.Address())
	assert.True(s.T(), binding)

	msgCounter, err = s.localFeature.BindToRemote(s.remoteServerFeature.Address())
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	msgCounter, err = s.localFeature.BindToRemote(s.remoteSubFeature.Address())
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	msgCounter, err = s.localFeature.RemoveRemoteBinding(s.remoteServerFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	msg = s.responseMsg(s.localFeature, s.remoteServerFeature, *msgCounter, 0)
	lf.unbindResponseCallback(s.remoteServerFeature.Device(), s.remoteServerFeature.Address(), msg)

	binding = s.localFeature.HasBindingToRemote(s.remoteServerFeature.Address())
	assert.False(s.T(), binding)

	bindingAdd := model.BindingManagementRequestCallType{
		ClientAddress: s.remoteFeature.Address(),
		ServerAddress: s.localServerFeature.Address(),
	}
	s.localDevice.BindingManager().AddBinding(s.remoteFeature.Device(), bindingAdd)

	binding = s.localServerFeature.HasBindingToRemote(s.remoteFeature.Address())
	assert.True(s.T(), binding)

	msgCounter, err = s.localServerFeature.RemoveRemoteBinding(s.remoteFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	lf = s.localServerFeature.(*FeatureLocal)
	msg = s.responseMsg(s.localServerFeature, s.remoteFeature, *msgCounter, 0)
	lf.unbindResponseCallback(s.remoteFeature.Device(), s.remoteFeature.Address(), msg)

	binding = s.localServerFeature.HasBindingToRemote(s.remoteFeature.Address())
	assert.False(s.T(), binding)
}

func (s *LocalFeatureTestSuite) responseMsg(featureLocal api.FeatureLocalInterface, featureRemote api.FeatureRemoteInterface, msgCounter model.MsgCounterType, errorNumber uint) api.ResponseMessage {
	resultData := &model.ResultDataType{
		ErrorNumber: util.Ptr(model.ErrorNumberType(errorNumber)),
	}

	msg := api.ResponseMessage{
		MsgCounterReference: msgCounter,
		Data:                resultData,
		FeatureLocal:        featureLocal,
		FeatureRemote:       featureRemote,
		EntityRemote:        featureRemote.Entity(),
		DeviceRemote:        featureRemote.Device(),
	}
	return msg
}

func (s *LocalFeatureTestSuite) Test_CleanRemoteDeviceCaches() {
	s.senderMock.On("Subscribe", mock.Anything, mock.Anything, mock.Anything).Return(&s.msgCounter, nil)
	s.senderMock.On("Bind", mock.Anything, mock.Anything, mock.Anything).Return(&s.msgCounter, nil)

	s.localFeature.CleanWriteApprovalCaches("testdummytest")
	s.localFeature.CleanRemoteDeviceCaches(nil)

	address := &model.DeviceAddressType{}
	s.localFeature.CleanRemoteDeviceCaches(address)

	address.Device = util.Ptr(model.AddressDeviceType("dummy"))
	s.localFeature.CleanRemoteDeviceCaches(address)

	address.Device = s.remoteServerFeature.Address().Device
	s.localFeature.CleanRemoteDeviceCaches(address)

	s.localFeature.Device().AddRemoteDeviceForSki(s.remoteServerFeature.Device().Ski(), s.remoteServerFeature.Device())
	s.localFeature.Device().AddRemoteDeviceForSki(s.remote2ServerFeature.Device().Ski(), s.remote2ServerFeature.Device())

	msgCounter, err := s.localFeature.SubscribeToRemote(s.remoteServerFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	lf := s.localFeature.(*FeatureLocal)
	msg := s.responseMsg(s.localFeature, s.remoteServerFeature, *msgCounter, 0)
	lf.subscribeResponseCallback(s.remoteServerFeature.Device(), s.remoteServerFeature.Address(), s.remoteServerFeature.Type(), msg)

	msgCounter, err = s.localFeature.SubscribeToRemote(s.remoteFeature.Address())
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	msgCounter, err = s.localFeature.SubscribeToRemote(s.remoteServerSubFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	msg = s.responseMsg(s.localFeature, s.remoteServerSubFeature, *msgCounter, 0)
	lf.subscribeResponseCallback(s.remoteServerSubFeature.Device(), s.remoteServerSubFeature.Address(), s.remoteServerSubFeature.Type(), msg)

	value := s.localFeature.HasSubscriptionToRemote(s.remoteServerFeature.Address())
	assert.True(s.T(), value)

	value = s.localFeature.HasSubscriptionToRemote(s.remote2ServerFeature.Address())
	assert.False(s.T(), value)

	msgCounter, err = s.localFeature.BindToRemote(s.remote2ServerFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	msg = s.responseMsg(s.localFeature, s.remote2ServerFeature, *msgCounter, 0)
	lf.bindResponseCallback(s.remote2ServerFeature.Device(), s.remote2ServerFeature.Address(), s.remote2ServerFeature.Type(), msg)

	msgCounter, err = s.localFeature.BindToRemote(s.remoteServerFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	msg = s.responseMsg(s.localFeature, s.remoteServerFeature, *msgCounter, 0)
	lf.bindResponseCallback(s.remoteServerFeature.Device(), s.remoteServerFeature.Address(), s.remoteServerFeature.Type(), msg)

	msgCounter, err = s.localFeature.BindToRemote(s.remoteServerSubFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	msg = s.responseMsg(s.localFeature, s.remoteServerSubFeature, *msgCounter, 7)
	lf.bindResponseCallback(s.remoteServerSubFeature.Device(), s.remoteServerSubFeature.Address(), s.remoteServerSubFeature.Type(), msg)

	value = s.localFeature.HasBindingToRemote(s.remote2ServerFeature.Address())
	assert.True(s.T(), value)

	value = s.localFeature.HasBindingToRemote(s.remoteServerFeature.Address())
	assert.True(s.T(), value)

	value = s.localFeature.HasBindingToRemote(s.remoteServerSubFeature.Address())
	assert.False(s.T(), value)

	s.localFeature.CleanRemoteDeviceCaches(address)

	value = s.localFeature.HasSubscriptionToRemote(s.remote2ServerFeature.Address())
	assert.False(s.T(), value)

	value = s.localFeature.HasSubscriptionToRemote(s.remoteServerFeature.Address())
	assert.False(s.T(), value)

	value = s.localFeature.HasSubscriptionToRemote(s.remoteServerSubFeature.Address())
	assert.False(s.T(), value)

	value = s.localFeature.HasBindingToRemote(s.remote2ServerFeature.Address())
	assert.True(s.T(), value)

	value = s.localFeature.HasBindingToRemote(s.remoteServerFeature.Address())
	assert.False(s.T(), value)

	value = s.localFeature.HasBindingToRemote(s.remoteServerSubFeature.Address())
	assert.False(s.T(), value)
}

func (s *LocalFeatureTestSuite) Test_CleanRemoteEntityCaches() {
	s.senderMock.On("Subscribe", mock.Anything, mock.Anything, mock.Anything).Return(&s.msgCounter, nil)
	s.senderMock.On("Bind", mock.Anything, mock.Anything, mock.Anything).Return(&s.msgCounter, nil)

	s.localFeature.CleanWriteApprovalCaches("testdummytest")
	s.localFeature.CleanRemoteEntityCaches(nil)

	address := &model.EntityAddressType{}
	s.localFeature.CleanRemoteEntityCaches(address)

	address.Device = util.Ptr(model.AddressDeviceType("dummy"))
	s.localFeature.CleanRemoteEntityCaches(address)

	address.Entity = []model.AddressEntityType{10}
	s.localFeature.CleanRemoteEntityCaches(address)

	address.Device = s.remoteServerFeature.Address().Device
	s.localFeature.CleanRemoteEntityCaches(address)

	address.Entity = s.remoteServerFeature.Address().Entity
	s.localFeature.CleanRemoteEntityCaches(address)

	s.localFeature.Device().AddRemoteDeviceForSki(s.remoteServerFeature.Device().Ski(), s.remoteServerFeature.Device())

	msgCounter, err := s.localFeature.SubscribeToRemote(s.remoteServerFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	lf := s.localFeature.(*FeatureLocal)
	msg := s.responseMsg(s.localFeature, s.remoteServerFeature, *msgCounter, 0)
	lf.subscribeResponseCallback(s.remoteServerFeature.Device(), s.remoteServerFeature.Address(), s.remoteServerFeature.Type(), msg)

	msgCounter, err = s.localFeature.SubscribeToRemote(s.remoteServerSubFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	msg = s.responseMsg(s.localFeature, s.remoteServerSubFeature, *msgCounter, 7)
	lf.subscribeResponseCallback(s.remoteServerSubFeature.Device(), s.remoteServerSubFeature.Address(), s.remoteServerSubFeature.Type(), msg)

	binding := s.localFeature.HasSubscriptionToRemote(s.remoteServerFeature.Address())
	assert.True(s.T(), binding)

	binding = s.localFeature.HasSubscriptionToRemote(s.remoteServerSubFeature.Address())
	assert.False(s.T(), binding)

	msgCounter, err = s.localFeature.BindToRemote(s.remoteServerFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	msg = s.responseMsg(s.localFeature, s.remoteServerFeature, *msgCounter, 0)
	lf.bindResponseCallback(s.remoteServerFeature.Device(), s.remoteServerFeature.Address(), s.remoteServerFeature.Type(), msg)

	msgCounter, err = s.localFeature.BindToRemote(s.remoteServerSubFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	msg = s.responseMsg(s.localFeature, s.remoteServerSubFeature, *msgCounter, 7)
	lf.bindResponseCallback(s.remoteServerSubFeature.Device(), s.remoteServerSubFeature.Address(), s.remoteServerSubFeature.Type(), msg)

	binding = s.localFeature.HasBindingToRemote(s.remoteServerFeature.Address())
	assert.True(s.T(), binding)

	binding = s.localFeature.HasBindingToRemote(s.remoteServerSubFeature.Address())
	assert.False(s.T(), binding)

	s.localFeature.CleanRemoteEntityCaches(address)

	binding = s.localFeature.HasSubscriptionToRemote(s.remoteServerFeature.Address())
	assert.False(s.T(), binding)

	binding = s.localFeature.HasSubscriptionToRemote(s.remoteServerSubFeature.Address())
	assert.False(s.T(), binding)

	binding = s.localFeature.HasBindingToRemote(s.remoteServerFeature.Address())
	assert.False(s.T(), binding)

	binding = s.localFeature.HasBindingToRemote(s.remoteServerSubFeature.Address())
	assert.False(s.T(), binding)
}

func (s *LocalFeatureTestSuite) Test_HandleMessage() {
	msg := &api.Message{
		FeatureRemote: s.remoteServerFeature,
		CmdClassifier: model.CmdClassifierType("buggy"),
		Cmd:           model.CmdType{},
	}

	err := s.localFeature.HandleMessage(msg)
	assert.NotNil(s.T(), err)

	msg = &api.Message{
		FeatureRemote: s.remoteServerFeature,
		CmdClassifier: model.CmdClassifierType("buggy"),
		Cmd: model.CmdType{
			ResultData: &model.ResultDataType{},
		},
	}

	err = s.localFeature.HandleMessage(msg)
	assert.NotNil(s.T(), err)
}

func (s *LocalFeatureTestSuite) Test_Result() {
	msg := &api.Message{
		FeatureRemote: s.remoteServerFeature,
		CmdClassifier: model.CmdClassifierTypeResult,
		Cmd: model.CmdType{
			ResultData: &model.ResultDataType{},
		},
	}

	err := s.localFeature.HandleMessage(msg)
	assert.NotNil(s.T(), err)

	msg.RequestHeader = &model.HeaderType{
		MsgCounterReference: util.Ptr(model.MsgCounterType(100)),
	}
	msg.Cmd.ResultData = &model.ResultDataType{
		ErrorNumber: util.Ptr(model.ErrorNumberType(1)),
		Description: util.Ptr(model.DescriptionType("test")),
	}
	err = s.localFeature.HandleMessage(msg)
	assert.Nil(s.T(), err)
}

func (s *LocalFeatureTestSuite) Test_Read() {
	msg := &api.Message{
		FeatureRemote: s.remoteFeature,
		CmdClassifier: model.CmdClassifierTypeRead,
		Cmd: model.CmdType{
			LoadControlLimitListData: &model.LoadControlLimitListDataType{},
		},
	}

	err := s.localFeature.HandleMessage(msg)
	assert.NotNil(s.T(), err)

	s.senderMock.EXPECT().Reply(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test")).Once()
	err = s.localServerFeature.HandleMessage(msg)
	assert.NotNil(s.T(), err)

	msg = &api.Message{
		FeatureRemote: s.remoteFeature,
		CmdClassifier: model.CmdClassifierTypeRead,
		Cmd: model.CmdType{
			DeviceClassificationManufacturerData: &model.DeviceClassificationManufacturerDataType{},
		},
	}
	err = s.localFeature.HandleMessage(msg)
	assert.NotNil(s.T(), err)

	s.senderMock.EXPECT().Reply(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test"))
	err = s.localServerFeature.HandleMessage(msg)
	assert.NotNil(s.T(), err)
}

func (s *LocalFeatureTestSuite) Test_Reply() {
	msg := &api.Message{
		FeatureRemote: s.remoteServerFeature,
		CmdClassifier: model.CmdClassifierTypeReply,
		Cmd: model.CmdType{
			DeviceClassificationManufacturerData: &model.DeviceClassificationManufacturerDataType{},
		},
	}

	err := s.localFeature.HandleMessage(msg)
	assert.Nil(s.T(), err)

	msg.RequestHeader = &model.HeaderType{
		MsgCounterReference: util.Ptr(model.MsgCounterType(100)),
	}
	err = s.localFeature.HandleMessage(msg)
	assert.Nil(s.T(), err)
}

func (s *LocalFeatureTestSuite) Test_Notify() {
	msg := &api.Message{
		FeatureRemote: s.remoteServerFeature,
		CmdClassifier: model.CmdClassifierTypeNotify,
		Cmd: model.CmdType{
			DeviceClassificationManufacturerData: &model.DeviceClassificationManufacturerDataType{},
		},
	}

	err := s.localFeature.HandleMessage(msg)
	assert.Nil(s.T(), err)
}

func (s *LocalFeatureTestSuite) Test_Write() {
	s.senderMock.EXPECT().ResultSuccess(mock.Anything, mock.Anything).Return(nil).Once()

	msg := &api.Message{
		RequestHeader: &model.HeaderType{
			MsgCounter: util.Ptr(model.MsgCounterType(1)),
			AckRequest: util.Ptr(true),
		},
		CmdClassifier: model.CmdClassifierTypeWrite,
		FeatureRemote: s.remoteSubFeature,
		Cmd: model.CmdType{
			LoadControlLimitListData: &model.LoadControlLimitListDataType{},
		},
	}

	err := s.localServerFeatureWrite.HandleMessage(msg)
	assert.Nil(s.T(), err)
}

func (s *LocalFeatureTestSuite) Test_SetWriteApprovalCallback_Invalid() {
	cb := func(msg *api.Message) {}
	err := s.localFeature.AddWriteApprovalCallback(cb)
	assert.NotNil(s.T(), err)
	result := model.ErrorType{
		ErrorNumber: 0,
	}
	s.localFeature.ApproveOrDenyWrite(&api.Message{}, result)
}

func (s *LocalFeatureTestSuite) Test_AddPendingApproval_Invalid() {
	cb := func(msg *api.Message) {}
	err := s.localServerFeatureWrite.AddWriteApprovalCallback(cb)
	assert.Nil(s.T(), err)

	msg := &api.Message{
		CmdClassifier: model.CmdClassifierTypeWrite,
		FeatureRemote: s.remoteSubFeature,
		Cmd: model.CmdType{
			LoadControlLimitListData: &model.LoadControlLimitListDataType{},
		},
	}

	err1 := s.localServerFeatureWrite.HandleMessage(msg)
	assert.Nil(s.T(), err1)
}

func (s *LocalFeatureTestSuite) Test_Write_Callback_One() {
	counter := model.MsgCounterType(1)
	msg := &api.Message{
		RequestHeader: &model.HeaderType{
			MsgCounter: util.Ptr(counter),
			AckRequest: util.Ptr(true),
		},
		CmdClassifier: model.CmdClassifierTypeWrite,
		FeatureRemote: s.remoteSubFeature,
		DeviceRemote:  s.remoteSubFeature.Device(),
		Cmd: model.CmdType{
			LoadControlLimitListData: &model.LoadControlLimitListDataType{},
		},
	}

	tempLSFWrite := s.localServerFeatureWrite

	cb1 := func(msg *api.Message) {
		result := model.ErrorType{
			ErrorNumber: 0,
		}
		tempLSFWrite.ApproveOrDenyWrite(msg, result)
	}
	tempLSFWrite.AddWriteApprovalCallback(cb1)

	s.senderMock.EXPECT().ResultSuccess(mock.Anything, mock.Anything).Return(nil).Once()
	err := tempLSFWrite.HandleMessage(msg)
	assert.Nil(s.T(), err)

	// callback is called asynchronously
	time.Sleep(time.Millisecond * 200)
}

func (s *LocalFeatureTestSuite) Test_Write_Callback_One_Fail() {
	counter := model.MsgCounterType(1)
	msg := &api.Message{
		RequestHeader: &model.HeaderType{
			MsgCounter: util.Ptr(counter),
			AckRequest: util.Ptr(true),
		},
		CmdClassifier: model.CmdClassifierTypeWrite,
		FeatureRemote: s.remoteSubFeature,
		DeviceRemote:  s.remoteSubFeature.Device(),
		Cmd: model.CmdType{
			LoadControlLimitListData: &model.LoadControlLimitListDataType{},
		},
	}

	tempLSFWrite := s.localServerFeatureWrite

	cb1 := func(msg *api.Message) {
		result := model.ErrorType{
			ErrorNumber: 7,
			Description: util.Ptr(model.DescriptionType("not allowed by application")),
		}
		tempLSFWrite.ApproveOrDenyWrite(msg, result)
	}
	tempLSFWrite.AddWriteApprovalCallback(cb1)

	s.senderMock.EXPECT().ResultError(mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	err := tempLSFWrite.HandleMessage(msg)
	assert.Nil(s.T(), err)

	// callback is called asynchronously
	time.Sleep(time.Millisecond * 200)
}

func (s *LocalFeatureTestSuite) Test_Write_Callback_Two() {
	msg := &api.Message{
		RequestHeader: &model.HeaderType{
			MsgCounter: util.Ptr(model.MsgCounterType(1)),
			AckRequest: util.Ptr(true),
		},
		CmdClassifier: model.CmdClassifierTypeWrite,
		FeatureRemote: s.remoteSubFeature,
		DeviceRemote:  s.remoteSubFeature.Device(),
		Cmd: model.CmdType{
			LoadControlLimitListData: &model.LoadControlLimitListDataType{},
		},
	}

	tempLSFWrite := s.localServerFeatureWrite

	cb1 := func(msg *api.Message) {
		result := model.ErrorType{
			ErrorNumber: 0,
		}
		tempLSFWrite.ApproveOrDenyWrite(msg, result)
	}
	tempLSFWrite.AddWriteApprovalCallback(cb1)

	cb2 := func(msg *api.Message) {
		result := model.ErrorType{
			ErrorNumber: 0,
		}
		tempLSFWrite.ApproveOrDenyWrite(msg, result)
	}
	tempLSFWrite.AddWriteApprovalCallback(cb2)

	s.senderMock.EXPECT().ResultSuccess(mock.Anything, mock.Anything).Return(nil).Once()
	err := tempLSFWrite.HandleMessage(msg)
	assert.Nil(s.T(), err)

	// callback is called asynchronously
	time.Sleep(time.Millisecond * 200)
}

func (s *LocalFeatureTestSuite) Test_Write_Callback_Two_Fail() {
	msg := &api.Message{
		RequestHeader: &model.HeaderType{
			MsgCounter: util.Ptr(model.MsgCounterType(1)),
			AckRequest: util.Ptr(true),
		},
		CmdClassifier: model.CmdClassifierTypeWrite,
		FeatureRemote: s.remoteSubFeature,
		DeviceRemote:  s.remoteSubFeature.Device(),
		Cmd: model.CmdType{
			LoadControlLimitListData: &model.LoadControlLimitListDataType{},
		},
	}

	tempLSFWrite := s.localServerFeatureWrite

	cb1 := func(msg *api.Message) {
		result := model.ErrorType{
			ErrorNumber: 0,
		}
		tempLSFWrite.ApproveOrDenyWrite(msg, result)
	}
	tempLSFWrite.AddWriteApprovalCallback(cb1)

	cb2 := func(msg *api.Message) {
		result := model.ErrorType{
			ErrorNumber: 7,
			Description: util.Ptr(model.DescriptionType("not allowed by application")),
		}
		tempLSFWrite.ApproveOrDenyWrite(msg, result)
	}
	tempLSFWrite.AddWriteApprovalCallback(cb2)

	s.senderMock.EXPECT().ResultError(mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	err := tempLSFWrite.HandleMessage(msg)
	assert.Nil(s.T(), err)

	// callback is called asynchronously
	time.Sleep(time.Millisecond * 200)
}

func (s *LocalFeatureTestSuite) Test_Write_Callback_Timeout() {
	s.localServerFeatureWrite.SetWriteApprovalTimeout(time.Millisecond * 500)

	s.senderMock.EXPECT().ResultError(mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	msg := &api.Message{
		RequestHeader: &model.HeaderType{
			MsgCounter: util.Ptr(model.MsgCounterType(1)),
			AckRequest: util.Ptr(true),
		},
		CmdClassifier: model.CmdClassifierTypeWrite,
		FeatureRemote: s.remoteSubFeature,
		DeviceRemote:  s.remoteSubFeature.Device(),
		Cmd: model.CmdType{
			LoadControlLimitListData: &model.LoadControlLimitListDataType{},
		},
	}

	tempLSFWrite := s.localServerFeatureWrite

	cb := func(msg *api.Message) {
		time.Sleep(time.Second * 1)
		result := model.ErrorType{
			ErrorNumber: 0,
		}
		tempLSFWrite.ApproveOrDenyWrite(msg, result)
	}

	tempLSFWrite.AddWriteApprovalCallback(cb)
	err := tempLSFWrite.HandleMessage(msg)
	assert.Nil(s.T(), err)

	// callback is called asynchronously
	time.Sleep(time.Second * 1)
}

func (s *LocalFeatureTestSuite) Test_Set_Update() {
	partial := model.NewFilterTypePartial()

	noData := s.localServerFeatureWrite.DataCopy("dummy")
	assert.Nil(s.T(), noData)

	data := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(1)),
				IsLimitActive: util.Ptr(false),
			},
		},
	}

	s.localServerFeatureWrite.SetData(s.serverWriteFunction, data)

	dataCopy := s.localServerFeatureWrite.DataCopy(s.serverWriteFunction)
	equal := reflect.DeepEqual(data, dataCopy)
	assert.True(s.T(), equal)

	newData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(1)),
				IsLimitActive: util.Ptr(true),
			},
		},
	}

	err := s.localServerFeatureWrite.UpdateData(s.serverWriteFunction, newData, nil, nil)
	assert.Nil(s.T(), err)

	dataCopy = s.localServerFeatureWrite.DataCopy(s.serverWriteFunction)
	equal = reflect.DeepEqual(data, dataCopy)
	assert.False(s.T(), equal)

	newData = &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(1)),
				IsLimitActive:     util.Ptr(true),
				IsLimitChangeable: util.Ptr(true),
			},
		},
	}
	err = s.localServerFeatureWrite.UpdateData(s.serverWriteFunction, newData, partial, nil)
	assert.Nil(s.T(), err)

	dataCopy = s.localServerFeatureWrite.DataCopy(s.serverWriteFunction)
	modelData, ok := dataCopy.(*model.LoadControlLimitListDataType)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), 1, len(modelData.LoadControlLimitData))
	assert.Equal(s.T(), model.LoadControlLimitIdType(1), *modelData.LoadControlLimitData[0].LimitId)
	assert.True(s.T(), *modelData.LoadControlLimitData[0].IsLimitActive)
	assert.True(s.T(), *modelData.LoadControlLimitData[0].IsLimitChangeable)

	moreData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(2)),
				IsLimitActive:     util.Ptr(false),
				IsLimitChangeable: util.Ptr(false),
			},
		},
	}
	err = s.localServerFeatureWrite.UpdateData(s.serverWriteFunction, moreData, partial, nil)
	assert.Nil(s.T(), err)

	dataCopy = s.localServerFeatureWrite.DataCopy(s.serverWriteFunction)
	modelData, ok = dataCopy.(*model.LoadControlLimitListDataType)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), 2, len(modelData.LoadControlLimitData))
	assert.Equal(s.T(), model.LoadControlLimitIdType(2), *modelData.LoadControlLimitData[1].LimitId)
	assert.False(s.T(), *modelData.LoadControlLimitData[1].IsLimitActive)
	assert.False(s.T(), *modelData.LoadControlLimitData[1].IsLimitChangeable)

	updateData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(2)),
				IsLimitChangeable: util.Ptr(true),
				TimePeriod:        model.NewTimePeriodTypeWithRelativeEndTime(time.Minute * 2),
			},
		},
	}
	err = s.localServerFeatureWrite.UpdateData(s.serverWriteFunction, updateData, partial, nil)
	assert.Nil(s.T(), err)

	dataCopy = s.localServerFeatureWrite.DataCopy(s.serverWriteFunction)
	modelData, ok = dataCopy.(*model.LoadControlLimitListDataType)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), 2, len(modelData.LoadControlLimitData))
	assert.Equal(s.T(), model.LoadControlLimitIdType(2), *modelData.LoadControlLimitData[1].LimitId)
	assert.False(s.T(), *modelData.LoadControlLimitData[1].IsLimitActive)
	assert.True(s.T(), *modelData.LoadControlLimitData[1].IsLimitChangeable)
	assert.NotNil(s.T(), modelData.LoadControlLimitData[1].TimePeriod)

	deleteFilter := &model.FilterType{
		LoadControlLimitListDataSelectors: &model.LoadControlLimitListDataSelectorsType{
			LimitId: util.Ptr(model.LoadControlLimitIdType(2)),
		},
		LoadControlLimitDataElements: &model.LoadControlLimitDataElementsType{
			TimePeriod: &model.TimePeriodElementsType{},
		},
	}
	updateData = &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(2)),
				IsLimitChangeable: util.Ptr(false),
			},
		},
	}
	err = s.localServerFeatureWrite.UpdateData(s.serverWriteFunction, updateData, partial, deleteFilter)
	assert.Nil(s.T(), err)

	dataCopy = s.localServerFeatureWrite.DataCopy(s.serverWriteFunction)
	modelData, ok = dataCopy.(*model.LoadControlLimitListDataType)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), 2, len(modelData.LoadControlLimitData))
	assert.False(s.T(), *modelData.LoadControlLimitData[1].IsLimitChangeable)
	assert.Nil(s.T(), modelData.LoadControlLimitData[1].TimePeriod)
}

// Test that read requests with partial filters return full data (spec-compliant behavior)
func (s *LocalFeatureTestSuite) Test_Read_WithPartialFilter_ReturnsFullData() {
	// Set up test data in server feature
	testData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(1)),
				IsLimitActive: util.Ptr(false),
				Value:         model.NewScaledNumberType(1000),
			},
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(2)),
				IsLimitActive: util.Ptr(true),
				Value:         model.NewScaledNumberType(2000),
			},
		},
	}
	s.localServerFeatureWrite.SetData(s.serverWriteFunction, testData)

	// Create partial filter (requesting only specific elements)
	partialFilter := &model.FilterType{
		CmdControl: &model.CmdControlType{
			Partial: &model.ElementTagType{},
		},
		LoadControlLimitDataElements: &model.LoadControlLimitDataElementsType{
			LimitId: &model.ElementTagType{},
		},
	}

	// Create read message with partial filter
	msg := &api.Message{
		FeatureRemote: s.remoteFeature,
		CmdClassifier: model.CmdClassifierTypeRead,
		FilterPartial: partialFilter,
		Cmd: model.CmdType{
			LoadControlLimitListData: &model.LoadControlLimitListDataType{},
			Filter: []model.FilterType{*partialFilter},
		},
	}

	// Expect full data reply (should NOT respect partial filter)
	s.senderMock.EXPECT().Reply(
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(cmd model.CmdType) bool {
			// Verify reply contains full data, not partial
			if cmd.LoadControlLimitListData == nil {
				return false
			}
			// Should contain all data, not just LimitId
			data := cmd.LoadControlLimitListData
			if len(data.LoadControlLimitData) != 2 {
				return false
			}
			// Both entries should have all fields (full data)
			entry1 := data.LoadControlLimitData[0]
			entry2 := data.LoadControlLimitData[1]
			return entry1.LimitId != nil && entry1.IsLimitActive != nil && entry1.Value != nil &&
				entry2.LimitId != nil && entry2.IsLimitActive != nil && entry2.Value != nil
		}),
	).Return(nil)

	// Handle the message
	err := s.localServerFeatureWrite.HandleMessage(msg)
	assert.Nil(s.T(), err)
}

// Test that read requests with selector filters return all data (ignore selectors)
func (s *LocalFeatureTestSuite) Test_Read_WithSelectorFilter_ReturnsAllData() {
	// Set up test data with multiple entries
	testData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(1)),
				IsLimitActive: util.Ptr(false),
			},
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(2)),
				IsLimitActive: util.Ptr(true),
			},
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(3)),
				IsLimitActive: util.Ptr(false),
			},
		},
	}
	s.localServerFeatureWrite.SetData(s.serverWriteFunction, testData)

	// Create selector filter (requesting only specific item)
	selectorFilter := &model.FilterType{
		CmdControl: &model.CmdControlType{
			Partial: &model.ElementTagType{},
		},
		LoadControlLimitListDataSelectors: &model.LoadControlLimitListDataSelectorsType{
			LimitId: util.Ptr(model.LoadControlLimitIdType(1)), // Only request item with ID 1
		},
	}

	// Create read message with selector filter
	msg := &api.Message{
		FeatureRemote: s.remoteFeature,
		CmdClassifier: model.CmdClassifierTypeRead,
		FilterPartial: selectorFilter,
		Cmd: model.CmdType{
			LoadControlLimitListData: &model.LoadControlLimitListDataType{},
			Filter: []model.FilterType{*selectorFilter},
		},
	}

	// Expect all data (should ignore selector filter)
	s.senderMock.EXPECT().Reply(
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(cmd model.CmdType) bool {
			// Verify reply contains ALL data, not just selected item
			if cmd.LoadControlLimitListData == nil {
				return false
			}
			data := cmd.LoadControlLimitListData
			// Should contain all 3 entries, not just the one with ID 1
			return len(data.LoadControlLimitData) == 3
		}),
	).Return(nil)

	// Handle the message
	err := s.localServerFeatureWrite.HandleMessage(msg)
	assert.Nil(s.T(), err)
}

// Test that read requests with combined element and selector filters return full data
func (s *LocalFeatureTestSuite) Test_Read_WithCombinedFilters_ReturnsFullData() {
	// Set up test data
	testData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(1)),
				IsLimitActive:     util.Ptr(false),
				IsLimitChangeable: util.Ptr(true),
				Value:             model.NewScaledNumberType(1000),
			},
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(2)),
				IsLimitActive:     util.Ptr(true),
				IsLimitChangeable: util.Ptr(false),
				Value:             model.NewScaledNumberType(2000),
			},
		},
	}
	s.localServerFeatureWrite.SetData(s.serverWriteFunction, testData)

	// Create combined filter (selector + elements)
	combinedFilter := &model.FilterType{
		CmdControl: &model.CmdControlType{
			Partial: &model.ElementTagType{},
		},
		LoadControlLimitListDataSelectors: &model.LoadControlLimitListDataSelectorsType{
			LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
		},
		LoadControlLimitDataElements: &model.LoadControlLimitDataElementsType{
			LimitId:       &model.ElementTagType{},
			IsLimitActive: &model.ElementTagType{},
		},
	}

	// Create read message with combined filter
	msg := &api.Message{
		FeatureRemote: s.remoteFeature,
		CmdClassifier: model.CmdClassifierTypeRead,
		FilterPartial: combinedFilter,
		Cmd: model.CmdType{
			LoadControlLimitListData: &model.LoadControlLimitListDataType{},
			Filter: []model.FilterType{*combinedFilter},
		},
	}

	// Expect full data (should ignore both selector and element filters)
	s.senderMock.EXPECT().Reply(
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(cmd model.CmdType) bool {
			// Verify reply contains full data
			if cmd.LoadControlLimitListData == nil {
				return false
			}
			data := cmd.LoadControlLimitListData
			if len(data.LoadControlLimitData) != 2 {
				return false
			}
			// Both entries should have all fields
			entry1 := data.LoadControlLimitData[0]
			entry2 := data.LoadControlLimitData[1]
			return entry1.LimitId != nil && entry1.IsLimitActive != nil && entry1.IsLimitChangeable != nil && entry1.Value != nil &&
				entry2.LimitId != nil && entry2.IsLimitActive != nil && entry2.IsLimitChangeable != nil && entry2.Value != nil
		}),
	).Return(nil)

	// Handle the message
	err := s.localServerFeatureWrite.HandleMessage(msg)
	assert.Nil(s.T(), err)
}

// Test that no errors are returned when partial filters are provided
func (s *LocalFeatureTestSuite) Test_Read_WithPartialFilter_NoErrors() {
	// Set up minimal test data
	testData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(1)),
				IsLimitActive: util.Ptr(false),
			},
		},
	}
	s.localServerFeatureWrite.SetData(s.serverWriteFunction, testData)

	// Create various partial filters to test
	partialFilter := &model.FilterType{
		CmdControl: &model.CmdControlType{
			Partial: &model.ElementTagType{},
		},
		LoadControlLimitDataElements: &model.LoadControlLimitDataElementsType{
			LimitId: &model.ElementTagType{},
		},
	}

	// Create read message with partial filter
	msg := &api.Message{
		FeatureRemote: s.remoteFeature,
		CmdClassifier: model.CmdClassifierTypeRead,
		FilterPartial: partialFilter,
		Cmd: model.CmdType{
			LoadControlLimitListData: &model.LoadControlLimitListDataType{},
			Filter: []model.FilterType{*partialFilter},
		},
	}

	// Expect successful reply (no errors)
	s.senderMock.EXPECT().Reply(mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Handle the message - should not return any errors
	err := s.localServerFeatureWrite.HandleMessage(msg)
	assert.Nil(s.T(), err)
}

// Test that partial read capability is correctly reported as false
func (s *LocalFeatureTestSuite) Test_Operations_NoPartialReadSupport() {
	operations := s.localServerFeatureWrite.Operations()
	
	// Verify that partial read is not supported
	operation, exists := operations[s.serverWriteFunction]
	assert.True(s.T(), exists)
	assert.False(s.T(), operation.ReadPartial())
}

// TestSubscribeToRemote_NilFeature verifies that SubscribeToRemote
// returns a proper error when FeatureByAddress() returns nil for the
// requested address, instead of panicking on a nil-pointer dereference.
//
// This protects callers (notably eebus-go's eg/lpc.connected() event
// handler) from crashing in the Publish() goroutine when a remote
// device advertises a feature address that is not yet materialized in
// the local entity model.
func (s *LocalFeatureTestSuite) TestSubscribeToRemote_NilFeature() {
	s.localFeature.Device().AddRemoteDeviceForSki(s.remoteFeature.Device().Ski(), s.remoteFeature.Device())

	// Build an address that targets the remote device but a non-existent
	// (entity, feature) tuple. FeatureByAddress() will return nil for it.
	bogusFeatureAddr := &model.FeatureAddressType{
		Device:  s.remoteFeature.Address().Device,
		Entity:  []model.AddressEntityType{99999},
		Feature: util.Ptr(model.AddressFeatureType(99999)),
	}

	// Must not panic; must return a proper ErrorType.
	assert.NotPanics(s.T(), func() {
		msgCounter, err := s.localFeature.SubscribeToRemote(bogusFeatureAddr)
		assert.Nil(s.T(), msgCounter)
		assert.NotNil(s.T(), err)
	})
}

// TestBindToRemote_NilFeature mirrors TestSubscribeToRemote_NilFeature
// for the BindToRemote() code path.
func (s *LocalFeatureTestSuite) TestBindToRemote_NilFeature() {
	s.localFeature.Device().AddRemoteDeviceForSki(s.remoteFeature.Device().Ski(), s.remoteFeature.Device())

	bogusFeatureAddr := &model.FeatureAddressType{
		Device:  s.remoteFeature.Address().Device,
		Entity:  []model.AddressEntityType{99999},
		Feature: util.Ptr(model.AddressFeatureType(99999)),
	}

	assert.NotPanics(s.T(), func() {
		msgCounter, err := s.localFeature.BindToRemote(bogusFeatureAddr)
		assert.Nil(s.T(), msgCounter)
		assert.NotNil(s.T(), err)
	})
}
