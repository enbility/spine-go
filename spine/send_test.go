package spine

import (
	"encoding/json"
	"testing"

	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func Test_SendSpineMessage(t *testing.T) {
	sut := &Sender{}

	datagram := model.DatagramType{}
	err := sut.sendSpineMessage(datagram)
	assert.NotNil(t, err)

	temp := &WriteMessageHandler{}
	sut.writeHandler = temp
	err = sut.sendSpineMessage(datagram)
	assert.Nil(t, err)
}

func Test_Cache(t *testing.T) {
	temp := &WriteMessageHandler{}
	sut := NewSender(temp)

	cmdClassifier := model.CmdClassifierTypeRead
	senderAddress := featureAddressType(1, NewEntityAddressType("Sender", []uint{1}))
	destinationAddress := featureAddressType(2, NewEntityAddressType("destination", []uint{1}))
	cmd := []model.CmdType{
		{
			ResultData: &model.ResultDataType{ErrorNumber: util.Ptr(model.ErrorNumberType(model.ErrorNumberTypeNoError))},
		},
	}

	msgCounter, err := sut.Request(cmdClassifier, senderAddress, destinationAddress, false, cmd)
	assert.NoError(t, err)
	assert.NotNil(t, msgCounter)

	msgCounter2, err := sut.Request(cmdClassifier, senderAddress, destinationAddress, false, cmd)
	assert.NoError(t, err)
	assert.NotNil(t, msgCounter2)
	assert.Equal(t, *msgCounter, *msgCounter2)

	sut.ProcessResponseForMsgCounterReference(msgCounter)

	msgCounter3, err := sut.Request(cmdClassifier, senderAddress, destinationAddress, false, cmd)
	assert.NoError(t, err)
	assert.NotNil(t, msgCounter3)
	assert.NotEqual(t, *msgCounter, *msgCounter3)

	for i := 0; i < 50; i++ {
		expMsgCounter4 := model.MsgCounterType(i + 3)
		destinationAddress = featureAddressType(2, NewEntityAddressType("destination", []uint{1}))
		cmd = []model.CmdType{
			{
				ResultData: &model.ResultDataType{ErrorNumber: util.Ptr(model.ErrorNumberType(i + 1))},
			},
		}

		msgCounter4, err := sut.Request(cmdClassifier, senderAddress, destinationAddress, false, cmd)
		assert.Nil(t, err)
		assert.NotNil(t, expMsgCounter4)
		assert.NotEqual(t, *msgCounter, *msgCounter4)
		assert.NotEqual(t, *msgCounter3, *msgCounter4)
	}
}

func TestSender_Reply_MsgCounter(t *testing.T) {
	temp := &WriteMessageHandler{}
	sut := NewSender(temp)

	senderAddress := featureAddressType(1, NewEntityAddressType("Sender", []uint{1}))
	destinationAddress := featureAddressType(2, NewEntityAddressType("destination", []uint{1}))
	requestHeader := &model.HeaderType{
		AddressSource:      senderAddress,
		AddressDestination: destinationAddress,
		MsgCounter:         util.Ptr(model.MsgCounterType(10)),
	}
	cmd := model.CmdType{
		ResultData: &model.ResultDataType{ErrorNumber: util.Ptr(model.ErrorNumberType(model.ErrorNumberTypeNoError))},
	}

	err := sut.Reply(requestHeader, senderAddress, cmd)
	assert.NoError(t, err)

	// Act
	err = sut.Reply(requestHeader, senderAddress, cmd)
	assert.NoError(t, err)
	expectedMsgCounter := 2 //because Notify was called twice

	sentBytes := temp.LastMessage()
	var sentDatagram model.Datagram
	assert.NoError(t, json.Unmarshal(sentBytes, &sentDatagram))
	assert.Equal(t, expectedMsgCounter, int(*sentDatagram.Datagram.Header.MsgCounter))
}

func TestSender_Notify_MsgCounter(t *testing.T) {
	temp := &WriteMessageHandler{}
	sut := NewSender(temp)

	senderAddress := featureAddressType(1, NewEntityAddressType("Sender", []uint{1}))
	destinationAddress := featureAddressType(2, NewEntityAddressType("destination", []uint{1}))
	cmd := model.CmdType{
		ResultData: &model.ResultDataType{ErrorNumber: util.Ptr(model.ErrorNumberType(model.ErrorNumberTypeNoError))},
	}

	_, err := sut.Notify(senderAddress, destinationAddress, cmd)
	assert.NoError(t, err)

	// Act
	_, err = sut.Notify(senderAddress, destinationAddress, cmd)
	assert.NoError(t, err)
	expectedMsgCounter := 2 //because Notify was called twice

	sentBytes := temp.LastMessage()
	var sentDatagram model.Datagram
	assert.NoError(t, json.Unmarshal(sentBytes, &sentDatagram))
	assert.Equal(t, expectedMsgCounter, int(*sentDatagram.Datagram.Header.MsgCounter))

	_, err = sut.DatagramForMsgCounter(model.MsgCounterType(2))
	assert.NoError(t, err)

	_, err = sut.DatagramForMsgCounter(model.MsgCounterType(3))
	assert.Error(t, err)
}

func TestSender_Write_MsgCounter(t *testing.T) {
	temp := &WriteMessageHandler{}
	sut := NewSender(temp)

	senderAddress := featureAddressType(1, NewEntityAddressType("Sender", []uint{1}))
	destinationAddress := featureAddressType(2, NewEntityAddressType("destination", []uint{1}))
	cmd := model.CmdType{
		ResultData: &model.ResultDataType{ErrorNumber: util.Ptr(model.ErrorNumberType(model.ErrorNumberTypeNoError))},
	}

	_, err := sut.Write(senderAddress, destinationAddress, cmd)
	assert.NoError(t, err)

	// Act
	_, err = sut.Write(senderAddress, destinationAddress, cmd)
	assert.NoError(t, err)
	expectedMsgCounter := 2 //because Write was called twice

	sentBytes := temp.LastMessage()
	var sentDatagram model.Datagram
	assert.NoError(t, json.Unmarshal(sentBytes, &sentDatagram))
	assert.Equal(t, expectedMsgCounter, int(*sentDatagram.Datagram.Header.MsgCounter))
}

func TestSender_Subscribe_MsgCounter(t *testing.T) {
	temp := &WriteMessageHandler{}
	sut := NewSender(temp)

	senderAddress := featureAddressType(1, NewEntityAddressType("Sender", []uint{1}))
	destinationAddress := featureAddressType(2, NewEntityAddressType("destination", []uint{1}))

	_, err := sut.Subscribe(senderAddress, destinationAddress, model.FeatureTypeTypeLoadControl)
	assert.NoError(t, err)

	// Act
	_, err = sut.Subscribe(senderAddress, destinationAddress, model.FeatureTypeTypeLoadControl)
	assert.NoError(t, err)
	expectedMsgCounter := 1 //because Subscribe was called twice and it was cached

	sentBytes := temp.LastMessage()
	var sentDatagram model.Datagram
	assert.NoError(t, json.Unmarshal(sentBytes, &sentDatagram))
	assert.Equal(t, expectedMsgCounter, int(*sentDatagram.Datagram.Header.MsgCounter))

	msgCounter := model.MsgCounterType(expectedMsgCounter)
	sut.ProcessResponseForMsgCounterReference(&msgCounter)

	_, err = sut.Subscribe(senderAddress, destinationAddress, model.FeatureTypeTypeLoadControl)
	assert.NoError(t, err)
	expectedMsgCounter = 2 //because Subscribe was called again

	sentBytes = temp.LastMessage()
	assert.NoError(t, json.Unmarshal(sentBytes, &sentDatagram))
	assert.Equal(t, expectedMsgCounter, int(*sentDatagram.Datagram.Header.MsgCounter))
}

func TestSender_Unsubscribe_MsgCounter(t *testing.T) {
	temp := &WriteMessageHandler{}
	sut := NewSender(temp)

	senderAddress := featureAddressType(1, NewEntityAddressType("Sender", []uint{1}))
	destinationAddress := featureAddressType(2, NewEntityAddressType("destination", []uint{1}))

	_, err := sut.Unsubscribe(senderAddress, destinationAddress)
	assert.NoError(t, err)

	// Act
	_, err = sut.Unsubscribe(senderAddress, destinationAddress)
	assert.NoError(t, err)
	expectedMsgCounter := 1 //because Unsubscribe was called twice and it was cached

	sentBytes := temp.LastMessage()
	var sentDatagram model.Datagram
	assert.NoError(t, json.Unmarshal(sentBytes, &sentDatagram))
	assert.Equal(t, expectedMsgCounter, int(*sentDatagram.Datagram.Header.MsgCounter))

	msgCounter := model.MsgCounterType(expectedMsgCounter)
	sut.ProcessResponseForMsgCounterReference(&msgCounter)

	_, err = sut.Unsubscribe(senderAddress, destinationAddress)
	assert.NoError(t, err)
	expectedMsgCounter = 2 //because Unsubscribe was called again

	sentBytes = temp.LastMessage()
	assert.NoError(t, json.Unmarshal(sentBytes, &sentDatagram))
	assert.Equal(t, expectedMsgCounter, int(*sentDatagram.Datagram.Header.MsgCounter))
}

func TestSender_Bind_MsgCounter(t *testing.T) {
	temp := &WriteMessageHandler{}
	sut := NewSender(temp)

	senderAddress := featureAddressType(1, NewEntityAddressType("Sender", []uint{1}))
	destinationAddress := featureAddressType(2, NewEntityAddressType("destination", []uint{1}))

	_, err := sut.Bind(senderAddress, destinationAddress, model.FeatureTypeTypeLoadControl)
	assert.NoError(t, err)

	// Act
	_, err = sut.Bind(senderAddress, destinationAddress, model.FeatureTypeTypeLoadControl)
	assert.NoError(t, err)
	expectedMsgCounter := 1 //because Bind was called twice and it was cached

	sentBytes := temp.LastMessage()
	var sentDatagram model.Datagram
	assert.NoError(t, json.Unmarshal(sentBytes, &sentDatagram))
	assert.Equal(t, expectedMsgCounter, int(*sentDatagram.Datagram.Header.MsgCounter))

	msgCounter := model.MsgCounterType(expectedMsgCounter)
	sut.ProcessResponseForMsgCounterReference(&msgCounter)

	_, err = sut.Bind(senderAddress, destinationAddress, model.FeatureTypeTypeLoadControl)
	assert.NoError(t, err)
	expectedMsgCounter = 2 //because Bind was called again

	sentBytes = temp.LastMessage()
	assert.NoError(t, json.Unmarshal(sentBytes, &sentDatagram))
	assert.Equal(t, expectedMsgCounter, int(*sentDatagram.Datagram.Header.MsgCounter))
}

func TestSender_Unbind_MsgCounter(t *testing.T) {
	temp := &WriteMessageHandler{}
	sut := NewSender(temp)

	senderAddress := featureAddressType(1, NewEntityAddressType("Sender", []uint{1}))
	destinationAddress := featureAddressType(2, NewEntityAddressType("destination", []uint{1}))

	_, err := sut.Unbind(senderAddress, destinationAddress)
	assert.NoError(t, err)

	// Act
	_, err = sut.Unbind(senderAddress, destinationAddress)
	assert.NoError(t, err)
	expectedMsgCounter := 1 //because Unbind was called twice and it was cached

	sentBytes := temp.LastMessage()
	var sentDatagram model.Datagram
	assert.NoError(t, json.Unmarshal(sentBytes, &sentDatagram))
	assert.Equal(t, expectedMsgCounter, int(*sentDatagram.Datagram.Header.MsgCounter))

	msgCounter := model.MsgCounterType(expectedMsgCounter)
	sut.ProcessResponseForMsgCounterReference(&msgCounter)

	_, err = sut.Unbind(senderAddress, destinationAddress)
	assert.NoError(t, err)
	expectedMsgCounter = 2 //because Unbind was called again

	sentBytes = temp.LastMessage()
	assert.NoError(t, json.Unmarshal(sentBytes, &sentDatagram))
	assert.Equal(t, expectedMsgCounter, int(*sentDatagram.Datagram.Header.MsgCounter))
}
