package spine

import (
	"reflect"
	"testing"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNodemanagement_BindingCalls(t *testing.T) {
	const bindingEntityId uint = 1
	const featureType = model.FeatureTypeTypeLoadControl
	const featureType2 = model.FeatureTypeTypeMeasurement
	const clientFeatureType = model.FeatureTypeTypeGeneric

	senderMock := mocks.NewSenderInterface(t)

	localDevice, localEntity := createLocalDeviceAndEntity(bindingEntityId)
	_, serverFeature := createLocalFeatures(localEntity, featureType, "")
	_, serverFeature2 := createLocalFeatures(localEntity, featureType2, "")

	remoteDevice := createRemoteDevice(localDevice, "ski", senderMock)
	clientFeature, _ := createRemoteEntityAndFeature(remoteDevice, bindingEntityId, clientFeatureType, "")

	remoteDevice2 := createRemoteDevice(localDevice, "ski2", senderMock)
	clientFeature2, _ := createRemoteEntityAndFeature(remoteDevice2, bindingEntityId, clientFeatureType, "")

	sut := NewNodeManagement(0, serverFeature.Entity())

	// add a binding to serverFeature from a remote device
	requestMsg := api.Message{
		Cmd: model.CmdType{
			NodeManagementBindingRequestCall: NewNodeManagementBindingRequestCallType(
				clientFeature.Address(), serverFeature.Address(), featureType),
		},
		CmdClassifier: model.CmdClassifierTypeCall,
		DeviceRemote:  remoteDevice,
		FeatureRemote: clientFeature,
	}

	err := sut.HandleMessage(&requestMsg)
	assert.Nil(t, err)

	// add a binding to serverFeature2 from remoteDevice2
	requestMsg = api.Message{
		Cmd: model.CmdType{
			NodeManagementBindingRequestCall: NewNodeManagementBindingRequestCallType(
				clientFeature2.Address(), serverFeature2.Address(), featureType2),
		},
		CmdClassifier: model.CmdClassifierTypeCall,
		DeviceRemote:  remoteDevice2,
		FeatureRemote: clientFeature2,
	}

	err = sut.HandleMessage(&requestMsg)
	assert.Nil(t, err)

	// remoteDevice reads its bindings
	// we should get a reply with one binding entry
	senderMock.On("Reply", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		cmd := args.Get(2).(model.CmdType)
		assert.Equal(t, 1, len(cmd.NodeManagementBindingData.BindingEntry))
		assert.True(t, reflect.DeepEqual(cmd.NodeManagementBindingData.BindingEntry[0].ClientAddress, clientFeature.Address()))
		assert.True(t, reflect.DeepEqual(cmd.NodeManagementBindingData.BindingEntry[0].ServerAddress, serverFeature.Address()))
	}).Return(nil).Once()

	dataMsg := api.Message{
		Cmd: model.CmdType{
			NodeManagementBindingData: &model.NodeManagementBindingDataType{},
		},
		CmdClassifier: model.CmdClassifierTypeRead,
		DeviceRemote:  remoteDevice,
		FeatureRemote: clientFeature,
	}
	err = sut.HandleMessage(&dataMsg)
	assert.Nil(t, err)

	// now delete the binding of remoteDevice
	deleteMsg := api.Message{
		Cmd: model.CmdType{
			NodeManagementBindingDeleteCall: NewNodeManagementBindingDeleteCallType(
				clientFeature.Address(), serverFeature.Address()),
		},
		CmdClassifier: model.CmdClassifierTypeCall,
		DeviceRemote:  remoteDevice,
		FeatureRemote: clientFeature,
	}

	err = sut.HandleMessage(&deleteMsg)
	assert.Nil(t, err)

	// when reading its bindings, we should get an empty list
	senderMock.On("Reply", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		cmd := args.Get(2).(model.CmdType)
		assert.Equal(t, 0, len(cmd.NodeManagementBindingData.BindingEntry))
	}).Return(nil).Once()

	dataMsg = api.Message{
		Cmd: model.CmdType{
			NodeManagementBindingData: &model.NodeManagementBindingDataType{},
		},
		CmdClassifier: model.CmdClassifierTypeRead,
		DeviceRemote:  remoteDevice,
		FeatureRemote: clientFeature,
	}
	err = sut.HandleMessage(&dataMsg)
	assert.Nil(t, err)

	// when reading remoteDevice2 bindings, we should get one entry
	senderMock.On("Reply", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		cmd := args.Get(2).(model.CmdType)
		assert.Equal(t, 1, len(cmd.NodeManagementBindingData.BindingEntry))
	}).Return(nil).Once()

	dataMsg = api.Message{
		Cmd: model.CmdType{
			NodeManagementBindingData: &model.NodeManagementBindingDataType{},
		},
		CmdClassifier: model.CmdClassifierTypeRead,
		DeviceRemote:  remoteDevice2,
		FeatureRemote: clientFeature2,
	}
	err = sut.HandleMessage(&dataMsg)
	assert.Nil(t, err)
}

func TestNodemanagement_SubscriptionCalls(t *testing.T) {
	const subscriptionEntityId uint = 1
	const featureType = model.FeatureTypeTypeDeviceClassification
	const clientFeatureType = model.FeatureTypeTypeGeneric

	senderMock := mocks.NewSenderInterface(t)

	localDevice, localEntity := createLocalDeviceAndEntity(subscriptionEntityId)
	_, serverFeature := createLocalFeatures(localEntity, featureType, "")

	remoteDevice := createRemoteDevice(localDevice, "ski", senderMock)
	clientFeature, _ := createRemoteEntityAndFeature(remoteDevice, subscriptionEntityId, clientFeatureType, "")

	remoteDevice2 := createRemoteDevice(localDevice, "ski2", senderMock)
	clientFeature2, _ := createRemoteEntityAndFeature(remoteDevice2, subscriptionEntityId, clientFeatureType, "")

	sut := NewNodeManagement(0, serverFeature.Entity())

	// add a subscription from remoteDevice to serverFeature
	requestMsg := api.Message{
		Cmd: model.CmdType{
			NodeManagementSubscriptionRequestCall: NewNodeManagementSubscriptionRequestCallType(
				clientFeature.Address(), serverFeature.Address(), featureType),
		},
		CmdClassifier: model.CmdClassifierTypeCall,
		DeviceRemote:  remoteDevice,
		FeatureRemote: clientFeature,
	}

	err := sut.HandleMessage(&requestMsg)
	assert.Nil(t, err)

	// add another subscription from remoteDevice2 to serverFeature
	requestMsg = api.Message{
		Cmd: model.CmdType{
			NodeManagementSubscriptionRequestCall: NewNodeManagementSubscriptionRequestCallType(
				clientFeature2.Address(), serverFeature.Address(), featureType),
		},
		CmdClassifier: model.CmdClassifierTypeCall,
		DeviceRemote:  remoteDevice2,
		FeatureRemote: clientFeature2,
	}

	err = sut.HandleMessage(&requestMsg)
	assert.Nil(t, err)

	// reading the subscription list of remoteDevice should return one entry
	senderMock.On("Reply", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		cmd := args.Get(2).(model.CmdType)
		assert.Equal(t, 1, len(cmd.NodeManagementSubscriptionData.SubscriptionEntry))
		assert.True(t, reflect.DeepEqual(cmd.NodeManagementSubscriptionData.SubscriptionEntry[0].ClientAddress, clientFeature.Address()))
		assert.True(t, reflect.DeepEqual(cmd.NodeManagementSubscriptionData.SubscriptionEntry[0].ServerAddress, serverFeature.Address()))
	}).Return(nil).Once()

	dataMsg := api.Message{
		Cmd: model.CmdType{
			NodeManagementSubscriptionData: &model.NodeManagementSubscriptionDataType{},
		},
		CmdClassifier: model.CmdClassifierTypeRead,
		DeviceRemote:  remoteDevice,
		FeatureRemote: clientFeature,
	}
	err = sut.HandleMessage(&dataMsg)
	assert.Nil(t, err)

	// delete the subscription from remoteDevice
	deleteMsg := api.Message{
		Cmd: model.CmdType{
			NodeManagementSubscriptionDeleteCall: NewNodeManagementSubscriptionDeleteCallType(
				clientFeature.Address(), serverFeature.Address()),
		},
		CmdClassifier: model.CmdClassifierTypeCall,
		DeviceRemote:  remoteDevice,
		FeatureRemote: clientFeature,
	}

	err = sut.HandleMessage(&deleteMsg)
	assert.Nil(t, err)

	// reading the subscription list of remoteDevice should return an emoty list
	senderMock.On("Reply", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		cmd := args.Get(2).(model.CmdType)
		assert.Equal(t, 0, len(cmd.NodeManagementSubscriptionData.SubscriptionEntry))
	}).Return(nil).Once()

	dataMsg = api.Message{
		Cmd: model.CmdType{
			NodeManagementSubscriptionData: &model.NodeManagementSubscriptionDataType{},
		},
		CmdClassifier: model.CmdClassifierTypeRead,
		DeviceRemote:  remoteDevice,
		FeatureRemote: clientFeature,
	}
	err = sut.HandleMessage(&dataMsg)
	assert.Nil(t, err)

	// reading the subscription list of remoteDevice2 should return one entry
	senderMock.On("Reply", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		cmd := args.Get(2).(model.CmdType)
		assert.Equal(t, 1, len(cmd.NodeManagementSubscriptionData.SubscriptionEntry))
	}).Return(nil).Once()

	dataMsg = api.Message{
		Cmd: model.CmdType{
			NodeManagementSubscriptionData: &model.NodeManagementSubscriptionDataType{},
		},
		CmdClassifier: model.CmdClassifierTypeRead,
		DeviceRemote:  remoteDevice2,
		FeatureRemote: clientFeature2,
	}
	err = sut.HandleMessage(&dataMsg)
	assert.Nil(t, err)
}
