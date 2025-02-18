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

	localDevice.AddRemoteDeviceForSki(remoteDevice.ski, remoteDevice)
	localDevice.AddRemoteDeviceForSki(remoteDevice2.ski, remoteDevice2)

	sut := NewNodeManagement(0, serverFeature.Entity())

	// add a binding to serverFeature from a remote device without providing the feature type
	requestMsg := api.Message{
		Cmd: model.CmdType{
			NodeManagementBindingRequestCall: &model.NodeManagementBindingRequestCallType{
				BindingRequest: &model.BindingManagementRequestCallType{
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
	sut.Device().BindingManager().RemoveBindingsForLocalEntity(localEntity)

	// add a binding to serverFeature from a remote device
	requestMsg = api.Message{
		Cmd: model.CmdType{
			NodeManagementBindingRequestCall: NewNodeManagementBindingRequestCallType(
				clientFeature.Address(), serverFeature.Address(), featureType),
		},
		CmdClassifier: model.CmdClassifierTypeCall,
		DeviceRemote:  remoteDevice,
		FeatureRemote: clientFeature,
	}

	err = sut.HandleMessage(&requestMsg)
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

	// test createBindingAddMissingDeviceAddresses
	bindingCreate := &model.BindingManagementRequestCallType{
		ClientAddress:     clientFeature.Address(),
		ServerAddress:     serverFeature.Address(),
		ServerFeatureType: util.Ptr(serverFeature.Type()),
	}

	dataMsg = api.Message{
		Cmd: model.CmdType{
			NodeManagementBindingRequestCall: &model.NodeManagementBindingRequestCallType{
				BindingRequest: bindingCreate,
			},
		},
		CmdClassifier: model.CmdClassifierTypeCall,
		DeviceRemote:  remoteDevice,
		FeatureRemote: clientFeature,
	}
	dataCreate := sut.createBindingAddMissingDeviceAddresses(&dataMsg, bindingCreate)
	assert.NotNil(t, dataCreate)

	bindingCreate.ClientAddress.Device = nil
	bindingCreate.ServerAddress.Device = serverFeature.Address().Device
	dataCreate = sut.createBindingAddMissingDeviceAddresses(&dataMsg, bindingCreate)
	assert.NotNil(t, dataCreate)
	assert.Equal(t, *dataCreate.ClientAddress.Device, *remoteDevice.Address())

	bindingCreate.ClientAddress.Device = clientFeature.Address().Device
	bindingCreate.ServerAddress.Device = nil
	dataCreate = sut.createBindingAddMissingDeviceAddresses(&dataMsg, bindingCreate)
	assert.NotNil(t, dataCreate)
	assert.Equal(t, *dataCreate.ServerAddress.Device, *localDevice.Address())

	// test deleteBindingAddMissingDeviceAddresses

	bindingDelete := &model.BindingManagementDeleteCallType{
		ClientAddress: clientFeature.Address(),
		ServerAddress: serverFeature.Address(),
	}

	dataMsg = api.Message{
		Cmd: model.CmdType{
			NodeManagementBindingDeleteCall: &model.NodeManagementBindingDeleteCallType{
				BindingDelete: bindingDelete,
			},
		},
		CmdClassifier: model.CmdClassifierTypeCall,
		DeviceRemote:  remoteDevice,
		FeatureRemote: clientFeature,
	}
	dataDelete := sut.deleteBindingAddMissingDeviceAddresses(&dataMsg, bindingDelete)
	assert.NotNil(t, dataDelete)

	bindingDelete.ClientAddress.Device = nil
	bindingDelete.ServerAddress.Device = nil
	dataDelete = sut.deleteBindingAddMissingDeviceAddresses(&dataMsg, bindingDelete)
	assert.NotNil(t, dataDelete)
	assert.Equal(t, *dataDelete.ClientAddress.Device, *localDevice.Address())
	assert.Equal(t, *dataDelete.ServerAddress.Device, *remoteDevice.Address())

	bindingDelete.ClientAddress.Device = nil
	bindingDelete.ServerAddress.Device = remoteDevice.Address()
	dataDelete = sut.deleteBindingAddMissingDeviceAddresses(&dataMsg, bindingDelete)
	assert.NotNil(t, dataDelete)
	assert.Equal(t, *dataDelete.ClientAddress.Device, *localDevice.Address())

	bindingDelete.ClientAddress.Device = remoteDevice.Address()
	bindingDelete.ServerAddress.Device = nil
	dataDelete = sut.deleteBindingAddMissingDeviceAddresses(&dataMsg, bindingDelete)
	assert.NotNil(t, dataDelete)
	assert.Equal(t, *dataDelete.ServerAddress.Device, *localDevice.Address())
}
