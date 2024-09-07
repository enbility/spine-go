package api

import "github.com/enbility/spine-go/model"

/* Sender */

type ComControlInterface interface {
	// This must be connected to the correct remote device !!
	SendSpineMessage(datagram model.DatagramType) error
}

type SenderInterface interface {
	// Process a received message, e.g. for handling caching data
	ProcessResponseForMsgCounterReference(msgCounterRef *model.MsgCounterType)
	// Sends a read cmd to request some data
	Request(cmdClassifier model.CmdClassifierType, senderAddress, destinationAddress *model.FeatureAddressType, ackRequest bool, cmd []model.CmdType) (*model.MsgCounterType, error)
	// Sends a result cmd with no error to indicate that a message was processed successfully
	ResultSuccess(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType) error
	// Sends a result cmd with error information to indicate that a message processing failed
	ResultError(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType, err *model.ErrorType) error
	// Sends a reply cmd to response to a read cmd
	Reply(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType, cmd model.CmdType) error
	// Sends a call cmd with a subscription request
	Subscribe(senderAddress, destinationAddress *model.FeatureAddressType, serverFeatureType model.FeatureTypeType) (*model.MsgCounterType, error)
	// Sends a call cmd with a subscription delete request
	Unsubscribe(senderAddress, destinationAddress *model.FeatureAddressType) (*model.MsgCounterType, error)
	// Sends a call cmd with a binding request
	Bind(senderAddress, destinationAddress *model.FeatureAddressType, serverFeatureType model.FeatureTypeType) (*model.MsgCounterType, error)
	// Sends a call cmd with a binding delte request
	Unbind(senderAddress, destinationAddress *model.FeatureAddressType) (*model.MsgCounterType, error)
	// Sends a notify cmd to indicate that a subscribed feature changed
	Notify(senderAddress, destinationAddress *model.FeatureAddressType, cmd model.CmdType) (*model.MsgCounterType, error)
	// Sends a write cmd, setting properties of remote features
	Write(senderAddress, destinationAddress *model.FeatureAddressType, cmd model.CmdType) (*model.MsgCounterType, error)
	// return the datagram for a given msgCounter (only availbe for Notify messages!), error if not found
	DatagramForMsgCounter(msgCounter model.MsgCounterType) (model.DatagramType, error)
}
