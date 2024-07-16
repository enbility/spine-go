package spine

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

const (
	wallbox_detaileddiscoverydata_recv_reply_file_path  = "./testdata/wallbox_detaileddiscoverydata_recv_reply.json"
	wallbox_detaileddiscoverydata_recv_notify_file_path = "./testdata/wallbox_detaileddiscoverydata_recv_notify.json"
)

type WriteMessageHandler struct {
	sentMessages [][]byte

	mux sync.Mutex
}

var _ shipapi.ShipConnectionDataWriterInterface = (*WriteMessageHandler)(nil)

func (t *WriteMessageHandler) WriteShipMessageWithPayload(message []byte) {
	t.mux.Lock()
	defer t.mux.Unlock()

	t.sentMessages = append(t.sentMessages, message)
}

func (t *WriteMessageHandler) LastMessage() []byte {
	t.mux.Lock()
	defer t.mux.Unlock()

	if len(t.sentMessages) == 0 {
		return nil
	}

	return t.sentMessages[len(t.sentMessages)-1]
}

func (t *WriteMessageHandler) MessageWithReference(msgCounterReference *model.MsgCounterType) []byte {
	t.mux.Lock()
	defer t.mux.Unlock()

	var datagram model.Datagram

	for _, msg := range t.sentMessages {
		if err := json.Unmarshal(msg, &datagram); err != nil {
			return nil
		}
		if datagram.Datagram.Header.MsgCounterReference == nil {
			continue
		}
		if uint(*datagram.Datagram.Header.MsgCounterReference) != uint(*msgCounterReference) {
			continue
		}
		if datagram.Datagram.Payload.Cmd[0].ResultData != nil {
			continue
		}

		return msg
	}

	return nil
}

func (t *WriteMessageHandler) ResultWithReference(msgCounterReference *model.MsgCounterType) []byte {
	t.mux.Lock()
	defer t.mux.Unlock()

	var datagram model.Datagram

	for _, msg := range t.sentMessages {
		if err := json.Unmarshal(msg, &datagram); err != nil {
			return nil
		}
		if datagram.Datagram.Header.MsgCounterReference == nil {
			continue
		}
		if uint(*datagram.Datagram.Header.MsgCounterReference) != uint(*msgCounterReference) {
			continue
		}
		if datagram.Datagram.Payload.Cmd[0].ResultData == nil {
			continue
		}

		return msg
	}

	return nil
}

func loadFileData(t *testing.T, fileName string) []byte {
	fileData, err := os.ReadFile(fileName) // #nosec G304
	if err != nil {
		t.Fatal(err)
	}

	return fileData
}

func checkSentData(t *testing.T, sendBytes []byte, msgSendFilePrefix string) {
	msgSendExpectedBytes, err := os.ReadFile(msgSendFilePrefix + "_expected.json") // #nosec G304
	if err != nil {
		t.Fatal(err)
	}

	msgSendActualFileName := msgSendFilePrefix + "_actual.json"
	equal := jsonDatagramEqual(t, msgSendExpectedBytes, sendBytes)
	if !equal {
		saveJsonToFile(t, sendBytes, msgSendActualFileName)
	}
	assert.Truef(t, equal, "Assert equal failed! Check '%s' ", msgSendActualFileName)
}

func jsonDatagramEqual(t *testing.T, expectedJson, actualJson []byte) bool {
	var actualDatagram model.Datagram
	if err := json.Unmarshal(actualJson, &actualDatagram); err != nil {
		t.Fatal(err)
	}
	var expectedDatagram model.Datagram
	if err := json.Unmarshal(expectedJson, &expectedDatagram); err != nil {
		t.Fatal(err)
	}

	less := func(a, b model.FunctionPropertyType) bool { return string(*a.Function) < string(*b.Function) }
	return cmp.Equal(expectedDatagram, actualDatagram, cmpopts.SortSlices(less))
}

func saveJsonToFile(t *testing.T, data json.RawMessage, fileName string) {
	jsonIndent, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(fileName, jsonIndent, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}

func waitForAck(t *testing.T, msgCounterReference *model.MsgCounterType, writeHandler *WriteMessageHandler) {
	var datagram model.Datagram

	msg := writeHandler.ResultWithReference(msgCounterReference)
	if msg == nil {
		t.Fatal("acknowledge message was not sent!!")
	}

	if err := json.Unmarshal(msg, &datagram); err != nil {
		t.Fatal(err)
	}

	cmd := datagram.Datagram.Payload.Cmd[0]
	if cmd.ResultData != nil {
		if cmd.ResultData.ErrorNumber != nil && uint(*cmd.ResultData.ErrorNumber) != uint(model.ErrorNumberTypeNoError) {
			t.Fatal(fmt.Errorf("error '%d' result data received", uint(*cmd.ResultData.ErrorNumber)))
		}
	}
}

func createLocalDeviceAndEntity(entityId uint) (*DeviceLocal, *EntityLocal) {
	localDevice := NewDeviceLocal("Vendor", "DeviceName", "SerialNumber", "DeviceCode", "Address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)
	localDevice.address = util.Ptr(model.AddressDeviceType("Address"))

	localEntity := NewEntityLocal(localDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{model.AddressEntityType(entityId)})
	localDevice.AddEntity(localEntity)

	return localDevice, localEntity
}

func createLocalFeatures(localEntity *EntityLocal, featureType model.FeatureTypeType, writeFunction model.FunctionType) (api.FeatureLocalInterface, api.FeatureLocalInterface) {
	localFeature := NewFeatureLocal(localEntity.NextFeatureId(), localEntity, featureType, model.RoleTypeClient)
	localEntity.AddFeature(localFeature)
	localServerFeature := NewFeatureLocal(localEntity.NextFeatureId(), localEntity, featureType, model.RoleTypeServer)
	if len(string(writeFunction)) > 0 {
		localServerFeature.AddFunctionType(writeFunction, true, true)
	}
	localEntity.AddFeature(localServerFeature)

	return localFeature, localServerFeature
}

func createRemoteDevice(localDevice *DeviceLocal, sender api.SenderInterface) *DeviceRemote {
	remoteDevice := NewDeviceRemote(localDevice, "ski", sender)
	remoteDevice.address = util.Ptr(model.AddressDeviceType("Address"))

	return remoteDevice
}

func createRemoteEntityAndFeature(remoteDevice *DeviceRemote, entityId uint, featureType model.FeatureTypeType, functionType model.FunctionType) (api.FeatureRemoteInterface, api.FeatureRemoteInterface) {
	remoteEntity := NewEntityRemote(remoteDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{model.AddressEntityType(entityId)})
	remoteDevice.AddEntity(remoteEntity)
	remoteFeature := NewFeatureRemote(remoteEntity.NextFeatureId(), remoteEntity, featureType, model.RoleTypeClient)
	remoteEntity.AddFeature(remoteFeature)
	remoteServerFeature := NewFeatureRemote(remoteEntity.NextFeatureId(), remoteEntity, featureType, model.RoleTypeServer)
	remoteServerFeature.SetOperations([]model.FunctionPropertyType{
		{
			Function: util.Ptr(functionType),
			PossibleOperations: &model.PossibleOperationsType{
				Read: &model.PossibleOperationsReadType{
					Partial: &model.ElementTagType{},
				},
			},
		},
	})
	remoteEntity.AddFeature(remoteServerFeature)

	return remoteFeature, remoteServerFeature
}
