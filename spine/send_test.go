package spine

import (
	"encoding/json"
	"sync"
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
	sut := NewSender(temp, nil)

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
	sut := NewSender(temp, nil)

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
	sut := NewSender(temp, nil)

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
	sut := NewSender(temp, nil)

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
	sut := NewSender(temp, nil)

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
	sut := NewSender(temp, nil)

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
	sut := NewSender(temp, nil)

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
	sut := NewSender(temp, nil)

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

// Comprehensive msgCounter Verification Tests

// TestSender_MsgCounter_ThreadSafety verifies thread-safe msgCounter generation
func TestSender_MsgCounter_ThreadSafety(t *testing.T) {
	temp := &WriteMessageHandler{}
	sut := NewSender(temp, nil)
	senderImpl := sut.(*Sender)

	const numGoroutines = 100
	const msgsPerGoroutine = 100
	totalMessages := numGoroutines * msgsPerGoroutine

	// Channel to collect all msgCounters
	countersChan := make(chan model.MsgCounterType, totalMessages)
	
	// WaitGroup to synchronize goroutines
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Launch concurrent goroutines
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < msgsPerGoroutine; j++ {
				counter := senderImpl.getMsgCounter()
				countersChan <- *counter
			}
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(countersChan)

	// Collect all counters
	counters := make([]model.MsgCounterType, 0, totalMessages)
	for counter := range countersChan {
		counters = append(counters, counter)
	}

	// Verify we got all counters
	assert.Equal(t, totalMessages, len(counters), "Should have all counters")

	// Check for uniqueness
	seen := make(map[model.MsgCounterType]bool)
	for _, counter := range counters {
		assert.False(t, seen[counter], "msgCounter %d should be unique", counter)
		seen[counter] = true
	}

	// All values should be between 1 and totalMessages
	for _, counter := range counters {
		assert.GreaterOrEqual(t, counter, model.MsgCounterType(1))
		assert.LessOrEqual(t, counter, model.MsgCounterType(totalMessages))
	}
}

// TestSender_MsgCounter_Uniqueness verifies msgCounters are unique within window
func TestSender_MsgCounter_Uniqueness(t *testing.T) {
	temp := &WriteMessageHandler{}
	sut := NewSender(temp, nil)
	senderImpl := sut.(*Sender)

	const numMessages = 10000
	counters := make([]model.MsgCounterType, numMessages)

	// Generate many msgCounters
	for i := 0; i < numMessages; i++ {
		counter := senderImpl.getMsgCounter()
		counters[i] = *counter
	}

	// Check for uniqueness
	seen := make(map[model.MsgCounterType]bool)
	for i, counter := range counters {
		assert.False(t, seen[counter], "msgCounter %d at position %d should be unique", counter, i)
		seen[counter] = true
	}

	// Verify ascending order (allowing gaps per spec)
	for i := 1; i < numMessages; i++ {
		assert.Greater(t, counters[i], counters[i-1], 
			"msgCounter at position %d (%d) should be greater than position %d (%d)", 
			i, counters[i], i-1, counters[i-1])
	}
}

// TestSender_MsgCounter_StartingValue verifies msgCounter starts from 1
func TestSender_MsgCounter_StartingValue(t *testing.T) {
	// Create multiple new senders to verify consistent behavior
	for i := 0; i < 5; i++ {
		temp := &WriteMessageHandler{}
		sut := NewSender(temp, nil)
		senderImpl := sut.(*Sender)
		
		counter := senderImpl.getMsgCounter()
		assert.Equal(t, model.MsgCounterType(1), *counter, 
			"First msgCounter for new sender %d should always be 1", i)
	}
}

// TestSender_MsgCounter_OverflowSimulation simulates overflow behavior at implementation level
func TestSender_MsgCounter_OverflowSimulation(t *testing.T) {
	temp := &WriteMessageHandler{}
	sut := NewSender(temp, nil)
	senderImpl := sut.(*Sender)
	
	// Set msgNum to max value - 1 to test overflow
	maxValue := ^uint64(0) - 1 // 2^64-2
	senderImpl.msgNum = maxValue
	
	// Next counter should be max value (2^64-1)
	counter1 := senderImpl.getMsgCounter()
	assert.Equal(t, model.MsgCounterType(maxValue+1), *counter1)
	assert.Equal(t, model.MsgCounterType(18446744073709551615), *counter1)
	
	// Next counter should overflow to 0
	counter2 := senderImpl.getMsgCounter()
	assert.Equal(t, model.MsgCounterType(0), *counter2, 
		"msgCounter should overflow from max (2^64-1) to 0 per SPINE spec")
	
	// Verify continued counting after overflow
	counter3 := senderImpl.getMsgCounter()
	assert.Equal(t, model.MsgCounterType(1), *counter3)
}
