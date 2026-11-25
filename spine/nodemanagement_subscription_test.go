package spine

import (
	"reflect"
	"testing"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TC_SPINE_SUBS_001: accept a subscription request to the primary NodeManagement feature.
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

	localDevice.AddRemoteDeviceForSki(remoteDevice.ski, remoteDevice)
	localDevice.AddRemoteDeviceForSki(remoteDevice2.ski, remoteDevice2)

	sut := NewNodeManagement(0, serverFeature.Entity())

	// add a subscription to serverFeature from a remote device without providing the feature type
	requestMsg := api.Message{
		Cmd: model.CmdType{
			NodeManagementSubscriptionRequestCall: &model.NodeManagementSubscriptionRequestCallType{
				SubscriptionRequest: &model.SubscriptionManagementRequestCallType{
					ClientAddress: clientFeature.Address(),
					ServerAddress: serverFeature.Address(),
				},
			},
		},
		CmdClassifier: model.CmdClassifierTypeCall,
		DeviceRemote:  remoteDevice,
		FeatureRemote: clientFeature,
	}

	err := sut.HandleMessage(&requestMsg)
	assert.Nil(t, err)

	// remove the binding again
	sut.Device().SubscriptionManager().RemoveSubscriptionsForLocalEntity(localEntity)

	// add a subscription from remoteDevice to serverFeature
	requestMsg = api.Message{
		Cmd: model.CmdType{
			NodeManagementSubscriptionRequestCall: NewNodeManagementSubscriptionRequestCallType(
				clientFeature.Address(), serverFeature.Address(), featureType),
		},
		CmdClassifier: model.CmdClassifierTypeCall,
		DeviceRemote:  remoteDevice,
		FeatureRemote: clientFeature,
	}

	err = sut.HandleMessage(&requestMsg)
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

	// test createSubscriptionAddMissingDeviceAddresses
	subscriptionCreate := &model.SubscriptionManagementRequestCallType{
		ClientAddress:     clientFeature.Address(),
		ServerAddress:     serverFeature.Address(),
		ServerFeatureType: util.Ptr(serverFeature.Type()),
	}

	dataMsg = api.Message{
		Cmd: model.CmdType{
			NodeManagementSubscriptionRequestCall: &model.NodeManagementSubscriptionRequestCallType{
				SubscriptionRequest: subscriptionCreate,
			},
		},
		CmdClassifier: model.CmdClassifierTypeCall,
		DeviceRemote:  remoteDevice,
		FeatureRemote: clientFeature,
	}
	dataCreate := sut.createSubscriptionAddMissingDeviceAddresses(&dataMsg, subscriptionCreate)
	assert.NotNil(t, dataCreate)

	subscriptionCreate.ClientAddress.Device = nil
	subscriptionCreate.ServerAddress.Device = serverFeature.Address().Device
	dataCreate = sut.createSubscriptionAddMissingDeviceAddresses(&dataMsg, subscriptionCreate)
	assert.NotNil(t, dataCreate)
	assert.Equal(t, *dataCreate.ClientAddress.Device, *remoteDevice.Address())

	subscriptionCreate.ClientAddress.Device = clientFeature.Address().Device
	subscriptionCreate.ServerAddress.Device = nil
	dataCreate = sut.createSubscriptionAddMissingDeviceAddresses(&dataMsg, subscriptionCreate)
	assert.NotNil(t, dataCreate)
	assert.Equal(t, *dataCreate.ServerAddress.Device, *localDevice.Address())

	// test deleteSubscriptionAddMissingDeviceAddresses

	subscriptionDelete := &model.SubscriptionManagementDeleteCallType{
		ClientAddress: clientFeature.Address(),
		ServerAddress: serverFeature.Address(),
	}

	dataMsg = api.Message{
		Cmd: model.CmdType{
			NodeManagementSubscriptionDeleteCall: &model.NodeManagementSubscriptionDeleteCallType{
				SubscriptionDelete: subscriptionDelete,
			},
		},
		CmdClassifier: model.CmdClassifierTypeCall,
		DeviceRemote:  remoteDevice,
		FeatureRemote: clientFeature,
	}
	dataDelete := sut.deleteSubscriptionAddMissingDeviceAddresses(&dataMsg, subscriptionDelete)
	assert.NotNil(t, dataDelete)

	subscriptionDelete.ClientAddress.Device = nil
	subscriptionDelete.ServerAddress.Device = nil
	dataDelete = sut.deleteSubscriptionAddMissingDeviceAddresses(&dataMsg, subscriptionDelete)
	assert.NotNil(t, dataDelete)
	assert.Equal(t, *dataDelete.ClientAddress.Device, *localDevice.Address())
	assert.Equal(t, *dataDelete.ServerAddress.Device, *remoteDevice.Address())

	subscriptionDelete.ClientAddress.Device = nil
	subscriptionDelete.ServerAddress.Device = remoteDevice.Address()
	dataDelete = sut.deleteSubscriptionAddMissingDeviceAddresses(&dataMsg, subscriptionDelete)
	assert.NotNil(t, dataDelete)
	assert.Equal(t, *dataDelete.ClientAddress.Device, *localDevice.Address())

	subscriptionDelete.ClientAddress.Device = remoteDevice.Address()
	subscriptionDelete.ServerAddress.Device = nil
	dataDelete = sut.deleteSubscriptionAddMissingDeviceAddresses(&dataMsg, subscriptionDelete)
	assert.NotNil(t, dataDelete)
	assert.Equal(t, *dataDelete.ServerAddress.Device, *localDevice.Address())
}
