package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

const emptyResult string = "Send:  Unknown 0 unknown"

func TestTestPrintMessageOverview_Emtpty(t *testing.T) {
	datagram := &DatagramType{}

	result := datagram.PrintMessageOverview(true, "", "")
	assert.Equal(t, emptyResult, result)
}

func TestPrintMessageOverview_Read_Send(t *testing.T) {
	datagram := &DatagramType{
		Header: HeaderType{
			AddressSource: &FeatureAddressType{
				Device: util.Ptr(AddressDeviceType("localdevice")),
			},
			AddressDestination: &FeatureAddressType{
				Device: util.Ptr(AddressDeviceType("localdevice")),
			},
			MsgCounter:    util.Ptr(MsgCounterType(1)),
			CmdClassifier: util.Ptr(CmdClassifierTypeRead),
		},
		Payload: PayloadType{
			Cmd: []CmdType{
				{},
			},
		},
	}

	result := datagram.PrintMessageOverview(false, "", "")
	assert.NotEqual(t, emptyResult, result)
}

func TestPrintMessageOverview_Read_Recv(t *testing.T) {
	datagram := &DatagramType{
		Header: HeaderType{
			AddressSource: &FeatureAddressType{},
			AddressDestination: &FeatureAddressType{
				Device: util.Ptr(AddressDeviceType("localdevice")),
			},
			MsgCounter:    util.Ptr(MsgCounterType(1)),
			CmdClassifier: util.Ptr(CmdClassifierTypeRead),
		},
		Payload: PayloadType{
			Cmd: []CmdType{
				{},
			},
		},
	}

	result := datagram.PrintMessageOverview(true, "", "")
	assert.NotEqual(t, emptyResult, result)
}

func TestPrintMessageOverview_Reply_Recv(t *testing.T) {
	datagram := &DatagramType{
		Header: HeaderType{
			AddressSource: &FeatureAddressType{},
			AddressDestination: &FeatureAddressType{
				Device: util.Ptr(AddressDeviceType("localdevice")),
			},
			MsgCounter:          util.Ptr(MsgCounterType(1)),
			MsgCounterReference: util.Ptr(MsgCounterType(1)),
			CmdClassifier:       util.Ptr(CmdClassifierTypeReply),
		},
		Payload: PayloadType{
			Cmd: []CmdType{
				{},
			},
		},
	}

	result := datagram.PrintMessageOverview(true, "", "")
	assert.NotEqual(t, emptyResult, result)
}

func TestPrintMessageOverview_Result_Recv(t *testing.T) {
	datagram := &DatagramType{
		Header: HeaderType{
			AddressSource: &FeatureAddressType{},
			AddressDestination: &FeatureAddressType{
				Device: util.Ptr(AddressDeviceType("localdevice")),
			},
			MsgCounter:          util.Ptr(MsgCounterType(1)),
			MsgCounterReference: util.Ptr(MsgCounterType(1)),
			CmdClassifier:       util.Ptr(CmdClassifierTypeResult),
		},
		Payload: PayloadType{
			Cmd: []CmdType{
				{
					ResultData: &ResultDataType{
						ErrorNumber: util.Ptr(ErrorNumberType(1)),
					},
				},
			},
		},
	}

	result := datagram.PrintMessageOverview(true, "", "")
	assert.NotEqual(t, emptyResult, result)
}

func TestPrintMessageOverview_Write_Recv(t *testing.T) {
	datagram := &DatagramType{
		Header: HeaderType{
			AddressSource: &FeatureAddressType{},
			AddressDestination: &FeatureAddressType{
				Device: util.Ptr(AddressDeviceType("localdevice")),
			},
			MsgCounter:          util.Ptr(MsgCounterType(1)),
			MsgCounterReference: util.Ptr(MsgCounterType(1)),
			CmdClassifier:       util.Ptr(CmdClassifierTypeWrite),
		},
		Payload: PayloadType{
			Cmd: []CmdType{
				{},
			},
		},
	}

	result := datagram.PrintMessageOverview(true, "", "")
	assert.NotEqual(t, emptyResult, result)
}
