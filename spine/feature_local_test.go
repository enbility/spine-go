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
	senderMock                                                *mocks.SenderInterface
	localDevice                                               *DeviceLocal
	localEntity                                               *EntityLocal
	function, serverWriteFunction                             model.FunctionType
	featureType, subFeatureType                               model.FeatureTypeType
	msgCounter                                                model.MsgCounterType
	remoteFeature, remoteServerFeature, remoteSubFeature      api.FeatureRemoteInterface
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

	remoteDevice := createRemoteDevice(s.localDevice, s.senderMock)
	s.remoteFeature, s.remoteServerFeature = createRemoteEntityAndFeature(remoteDevice, 1, s.featureType)
	s.remoteSubFeature, _ = createRemoteEntityAndFeature(remoteDevice, 2, s.subFeatureType)
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

	msgCounter, err = s.localServerFeature.SubscribeToRemote(s.remoteFeature.Address())
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	msgCounter, err = s.localFeature.RemoveRemoteSubscription(s.remoteFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	subscribed := s.localFeature.HasSubscriptionToRemote(s.remoteFeature.Address())
	assert.Equal(s.T(), false, subscribed)

	msgCounter, err = s.localFeature.SubscribeToRemote(s.remoteFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	subscribed = s.localFeature.HasSubscriptionToRemote(s.remoteFeature.Address())
	assert.Equal(s.T(), true, subscribed)

	msgCounter, err = s.localFeature.SubscribeToRemote(s.remoteSubFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	msgCounter, err = s.localFeature.RemoveRemoteSubscription(s.remoteFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	s.localFeature.RemoveAllRemoteSubscriptions()
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

	s.localFeature.Device().AddRemoteDeviceForSki(s.remoteFeature.Device().Ski(), s.remoteFeature.Device())

	msgCounter, err = s.localServerFeature.BindToRemote(s.remoteFeature.Address())
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	msgCounter, err = s.localFeature.RemoveRemoteBinding(s.remoteFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	binding := s.localFeature.HasBindingToRemote(s.remoteFeature.Address())
	assert.Equal(s.T(), false, binding)

	msgCounter, err = s.localFeature.BindToRemote(s.remoteFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	binding = s.localFeature.HasBindingToRemote(s.remoteFeature.Address())
	assert.Equal(s.T(), true, binding)

	msgCounter, err = s.localFeature.BindToRemote(s.remoteSubFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	msgCounter, err = s.localFeature.RemoveRemoteBinding(s.remoteFeature.Address())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msgCounter)

	s.localFeature.RemoveAllRemoteBindings()
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
		},
		CmdClassifier: model.CmdClassifierTypeWrite,
		FeatureRemote: s.remoteSubFeature,
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

	err := tempLSFWrite.HandleMessage(msg)
	assert.Nil(s.T(), err)

	// callback is called asynchronously
	time.Sleep(time.Millisecond * 200)
}

func (s *LocalFeatureTestSuite) Test_Write_Callback_One_Fail() {
	msg := &api.Message{
		RequestHeader: &model.HeaderType{
			MsgCounter: util.Ptr(model.MsgCounterType(1)),
		},
		CmdClassifier: model.CmdClassifierTypeWrite,
		FeatureRemote: s.remoteSubFeature,
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
		},
		CmdClassifier: model.CmdClassifierTypeWrite,
		FeatureRemote: s.remoteSubFeature,
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

	err := tempLSFWrite.HandleMessage(msg)
	assert.Nil(s.T(), err)

	// callback is called asynchronously
	time.Sleep(time.Millisecond * 200)
}

func (s *LocalFeatureTestSuite) Test_Write_Callback_Two_Fail() {
	msg := &api.Message{
		RequestHeader: &model.HeaderType{
			MsgCounter: util.Ptr(model.MsgCounterType(1)),
		},
		CmdClassifier: model.CmdClassifierTypeWrite,
		FeatureRemote: s.remoteSubFeature,
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
		},
		CmdClassifier: model.CmdClassifierTypeWrite,
		FeatureRemote: s.remoteSubFeature,
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
