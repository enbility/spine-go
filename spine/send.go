package spine

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"

	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/golanguzb70/lrucache"
)

type reqMsgCacheData map[model.MsgCounterType]string

type Sender struct {
	msgNum uint64 // 64bit values need to be defined on top of the struct to make atomic commands work on 32bit systems

	// we cache the last 100 notify messages, so we can find the matching item for result errors being returned
	datagramNotifyCache *lrucache.LRUCache[model.MsgCounterType, model.DatagramType]

	writeHandler shipapi.ShipConnectionDataWriterInterface

	reqMsgCache reqMsgCacheData // cache for unanswered request messages, so we can filter duplicates and not send them

	muxRequestSend sync.Mutex

	muxNotifyCache sync.RWMutex
	muxReadCache   sync.RWMutex
}

var _ api.SenderInterface = (*Sender)(nil)

func NewSender(writeI shipapi.ShipConnectionDataWriterInterface) api.SenderInterface {
	cache := lrucache.New[model.MsgCounterType, model.DatagramType](100, 0)
	return &Sender{
		datagramNotifyCache: &cache,
		writeHandler:        writeI,
		reqMsgCache:         make(reqMsgCacheData),
	}
}

// return the datagram for a given msgCounter (only availbe for Notify messasges!), error if not found
func (c *Sender) DatagramForMsgCounter(msgCounter model.MsgCounterType) (model.DatagramType, error) {
	c.muxNotifyCache.RLock()
	defer c.muxNotifyCache.RUnlock()

	if datagram, ok := c.datagramNotifyCache.Get(msgCounter); ok {
		return datagram, nil
	}

	return model.DatagramType{}, errors.New("msgCounter not found")
}

func (c *Sender) sendSpineMessage(datagram model.DatagramType) error {
	// pack into datagram
	data := model.Datagram{
		Datagram: datagram,
	}

	// marshal
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if c.writeHandler == nil {
		return errors.New("outgoing interface implementation not set")
	}

	if msg == nil {
		return errors.New("message is nil")
	}

	logging.Log().Debug(datagram.PrintMessageOverview(true, "", ""))

	// write to channel
	c.writeHandler.WriteShipMessageWithPayload(msg)

	return nil
}

// Caching of outgoing and unanswered requests, so we can filter duplicates
func (c *Sender) hashForMessage(destinationAddress *model.FeatureAddressType, cmd []model.CmdType) string {
	cmdString, err := json.Marshal(cmd)
	if err != nil {
		return ""
	}

	sig := fmt.Sprintf("%s-%s", destinationAddress.String(), cmdString)
	shaBytes := sha256.Sum256([]byte(sig))
	return hex.EncodeToString(shaBytes[:])
}

func (c *Sender) msgCounterForHashFromCache(hash string) *model.MsgCounterType {
	c.muxReadCache.RLock()
	defer c.muxReadCache.RUnlock()

	for msgCounter, h := range c.reqMsgCache {
		if h == hash {
			return &msgCounter
		}
	}

	return nil
}

func (c *Sender) hasMsgCounterInCache(msgCounter model.MsgCounterType) bool {
	c.muxReadCache.RLock()
	defer c.muxReadCache.RUnlock()

	_, ok := c.reqMsgCache[msgCounter]

	return ok
}

func (c *Sender) addMsgCounterHashToCache(msgCounter model.MsgCounterType, hash string) {
	c.muxReadCache.Lock()
	defer c.muxReadCache.Unlock()

	// cleanup cache, keep only the last 20 messages
	if len(c.reqMsgCache) > 20 {
		keys := make([]uint64, 0, len(c.reqMsgCache))
		for k := range c.reqMsgCache {
			keys = append(keys, uint64(k))
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

		// oldest key is the one with the lowest msgCounterValue
		oldestKey := keys[0]
		delete(c.reqMsgCache, model.MsgCounterType(oldestKey))
	}

	c.reqMsgCache[msgCounter] = hash
}

// we need to remove the msgCounter from the cache, if we have it cached
func (c *Sender) ProcessResponseForMsgCounterReference(msgCounterRef *model.MsgCounterType) {
	if msgCounterRef != nil &&
		c.hasMsgCounterInCache(*msgCounterRef) {
		c.muxReadCache.Lock()
		defer c.muxReadCache.Unlock()

		delete(c.reqMsgCache, *msgCounterRef)
	}
}

// Sends request
func (c *Sender) Request(cmdClassifier model.CmdClassifierType, senderAddress, destinationAddress *model.FeatureAddressType, ackRequest bool, cmd []model.CmdType) (*model.MsgCounterType, error) {
	// lock the method so caching works if the method is called really simultaniously and the cache therefor was not updated yet
	c.muxRequestSend.Lock()
	defer c.muxRequestSend.Unlock()

	// check if there is an unanswered subscribe message for this destination and cmd and return that msgCounter
	hash := c.hashForMessage(destinationAddress, cmd)
	if len(hash) > 0 {
		if msgCounterCache := c.msgCounterForHashFromCache(hash); msgCounterCache != nil {
			return msgCounterCache, nil
		}
	}

	msgCounter := c.getMsgCounter()

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &SpecificationVersion,
			AddressSource:        senderAddress,
			AddressDestination:   destinationAddress,
			MsgCounter:           msgCounter,
			CmdClassifier:        &cmdClassifier,
		},
		Payload: model.PayloadType{
			Cmd: cmd,
		},
	}

	if ackRequest {
		datagram.Header.AckRequest = &ackRequest
	}

	err := c.sendSpineMessage(datagram)
	if err == nil {
		if len(hash) > 0 {
			c.addMsgCounterHashToCache(*msgCounter, hash)
		}
	}

	return msgCounter, err
}

func (c *Sender) ResultSuccess(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType) error {
	return c.result(requestHeader, senderAddress, nil)
}

func (c *Sender) ResultError(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType, err *model.ErrorType) error {
	return c.result(requestHeader, senderAddress, err)
}

// sends a result for a request
func (c *Sender) result(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType, err *model.ErrorType) error {
	cmdClassifier := model.CmdClassifierTypeResult

	addressSource := *requestHeader.AddressDestination
	addressSource.Device = senderAddress.Device

	var resultData model.ResultDataType
	if err != nil {
		resultData = model.ResultDataType{
			ErrorNumber: &err.ErrorNumber,
			Description: err.Description,
		}
	} else {
		resultData = model.ResultDataType{
			ErrorNumber: util.Ptr(model.ErrorNumberTypeNoError),
		}
	}

	cmd := model.CmdType{
		ResultData: &resultData,
	}

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &SpecificationVersion,
			AddressSource:        &addressSource,
			AddressDestination:   requestHeader.AddressSource,
			MsgCounter:           c.getMsgCounter(),
			MsgCounterReference:  requestHeader.MsgCounter,
			CmdClassifier:        &cmdClassifier,
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{cmd},
		},
	}

	return c.sendSpineMessage(datagram)
}

// Reply sends reply to original sender
func (c *Sender) Reply(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType, cmd model.CmdType) error {
	cmdClassifier := model.CmdClassifierTypeReply

	addressSource := *requestHeader.AddressDestination
	addressSource.Device = senderAddress.Device

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &SpecificationVersion,
			AddressSource:        &addressSource,
			AddressDestination:   requestHeader.AddressSource,
			MsgCounter:           c.getMsgCounter(),
			MsgCounterReference:  requestHeader.MsgCounter,
			CmdClassifier:        &cmdClassifier,
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{cmd},
		},
	}

	return c.sendSpineMessage(datagram)
}

// Notify sends notification to destination
func (c *Sender) Notify(senderAddress, destinationAddress *model.FeatureAddressType, cmd model.CmdType) (*model.MsgCounterType, error) {
	msgCounter := c.getMsgCounter()

	cmdClassifier := model.CmdClassifierTypeNotify

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &SpecificationVersion,
			AddressSource:        senderAddress,
			AddressDestination:   destinationAddress,
			MsgCounter:           msgCounter,
			CmdClassifier:        &cmdClassifier,
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{cmd},
		},
	}

	c.muxNotifyCache.Lock()
	c.datagramNotifyCache.Put(*msgCounter, datagram)
	c.muxNotifyCache.Unlock()

	return msgCounter, c.sendSpineMessage(datagram)
}

// Write sends notification to destination
func (c *Sender) Write(senderAddress, destinationAddress *model.FeatureAddressType, cmd model.CmdType) (*model.MsgCounterType, error) {
	msgCounter := c.getMsgCounter()

	cmdClassifier := model.CmdClassifierTypeWrite
	ackRequest := true

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &SpecificationVersion,
			AddressSource:        senderAddress,
			AddressDestination:   destinationAddress,
			MsgCounter:           msgCounter,
			CmdClassifier:        &cmdClassifier,
			AckRequest:           &ackRequest,
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{cmd},
		},
	}

	return msgCounter, c.sendSpineMessage(datagram)
}

// Send a subscription request to a remote server feature
func (c *Sender) Subscribe(senderAddress, destinationAddress *model.FeatureAddressType, serverFeatureType model.FeatureTypeType) (*model.MsgCounterType, error) {
	cmd := model.CmdType{
		NodeManagementSubscriptionRequestCall: NewNodeManagementSubscriptionRequestCallType(senderAddress, destinationAddress, serverFeatureType),
	}

	// we always send it to the remote NodeManagement feature, which always is at entity:[0],feature:0
	localAddress := NodeManagementAddress(senderAddress.Device)
	remoteAddress := NodeManagementAddress(destinationAddress.Device)

	return c.Request(model.CmdClassifierTypeCall, localAddress, remoteAddress, true, []model.CmdType{cmd})
}

// Send a subscription deletion request to a remote server feature
func (c *Sender) Unsubscribe(senderAddress, destinationAddress *model.FeatureAddressType) (*model.MsgCounterType, error) {
	cmd := model.CmdType{
		NodeManagementSubscriptionDeleteCall: NewNodeManagementSubscriptionDeleteCallType(senderAddress, destinationAddress),
	}

	// we always send it to the remote NodeManagement feature, which always is at entity:[0],feature:0
	localAddress := NodeManagementAddress(senderAddress.Device)
	remoteAddress := NodeManagementAddress(destinationAddress.Device)

	return c.Request(model.CmdClassifierTypeCall, localAddress, remoteAddress, true, []model.CmdType{cmd})
}

// Send a binding request to a remote server feature
func (c *Sender) Bind(senderAddress, destinationAddress *model.FeatureAddressType, serverFeatureType model.FeatureTypeType) (*model.MsgCounterType, error) {
	cmd := model.CmdType{
		NodeManagementBindingRequestCall: NewNodeManagementBindingRequestCallType(senderAddress, destinationAddress, serverFeatureType),
	}

	// we always send it to the remote NodeManagement feature, which always is at entity:[0],feature:0
	localAddress := NodeManagementAddress(senderAddress.Device)
	remoteAddress := NodeManagementAddress(destinationAddress.Device)

	return c.Request(model.CmdClassifierTypeCall, localAddress, remoteAddress, true, []model.CmdType{cmd})
}

// Send a binding request to a remote server feature
func (c *Sender) Unbind(senderAddress, destinationAddress *model.FeatureAddressType) (*model.MsgCounterType, error) {
	cmd := model.CmdType{
		NodeManagementBindingDeleteCall: NewNodeManagementBindingDeleteCallType(senderAddress, destinationAddress),
	}

	// we always send it to the remote NodeManagement feature, which always is at entity:[0],feature:0
	localAddress := NodeManagementAddress(senderAddress.Device)
	remoteAddress := NodeManagementAddress(destinationAddress.Device)

	return c.Request(model.CmdClassifierTypeCall, localAddress, remoteAddress, true, []model.CmdType{cmd})
}

func (c *Sender) getMsgCounter() *model.MsgCounterType {
	// TODO:  persistence
	i := model.MsgCounterType(atomic.AddUint64(&c.msgNum, 1))
	return &i
}
