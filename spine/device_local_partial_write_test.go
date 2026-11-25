package spine

import (
	"testing"
	"time"

	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestDeviceLocalPartialWriteSuite(t *testing.T) {
	suite.Run(t, new(DeviceLocalPartialWriteSuite))
}

type DeviceLocalPartialWriteSuite struct {
	suite.Suite

	sut           *DeviceLocal
	localFeature  *FeatureLocal
	remoteDevice  *DeviceRemote
	remoteFeature *FeatureRemote
	lastMessage   string
}

var _ shipapi.ShipConnectionDataWriterInterface = (*DeviceLocalPartialWriteSuite)(nil)

func (s *DeviceLocalPartialWriteSuite) WriteShipMessageWithPayload(msg []byte) {
	s.lastMessage = string(msg)
}

func (s *DeviceLocalPartialWriteSuite) SetupTest() {
	s.sut = NewDeviceLocal("brand", "model", "serial", "code", "address",
		model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)

	localEntity := NewEntityLocal(s.sut, model.EntityTypeTypeCEM,
		NewAddressEntityType([]uint{1}), time.Second*4)
	s.sut.AddEntity(localEntity)

	// Create a server feature with write support but NO partial write support.
	// DeviceClassificationManufacturerData does not implement Updater,
	// so writePartial will be false.
	s.localFeature = NewFeatureLocal(1, localEntity,
		model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)
	s.localFeature.AddFunctionType(model.FunctionTypeDeviceClassificationManufacturerData, true, true)
	localEntity.AddFeature(s.localFeature)

	ski := "test"
	_ = s.sut.SetupRemoteDevice(ski, s)
	remote := s.sut.RemoteDeviceForSki(ski)
	s.remoteDevice = remote.(*DeviceRemote)
	s.remoteDevice.address = util.Ptr(model.AddressDeviceType("remoteDevice"))

	remoteEntity := NewEntityRemote(s.remoteDevice, model.EntityTypeTypeCEM,
		[]model.AddressEntityType{1})
	s.remoteFeature = NewFeatureRemote(1, remoteEntity,
		model.FeatureTypeTypeDeviceClassification, model.RoleTypeClient)
	remoteEntity.AddFeature(s.remoteFeature)
	s.remoteDevice.AddEntity(remoteEntity)

	// Add a binding so write permission checks pass
	binding := model.BindingManagementRequestCallType{
		ClientAddress:     s.remoteFeature.Address(),
		ServerAddress:     s.localFeature.Address(),
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeDeviceClassification),
	}
	err := s.sut.BindingManager().AddBinding(s.remoteDevice, binding)
	assert.Nil(s.T(), err)
}

// Test that a partial write to a feature that does not announce partial write
// is rejected with error code 8 (RestrictedFunctionExchangeCombinationNotSupported)
func (s *DeviceLocalPartialWriteSuite) Test_PartialWriteDenied_WhenNotSupported() {
	// Verify precondition: the function supports write but NOT writePartial
	ops := s.localFeature.Operations()
	fnOps, ok := ops[model.FunctionTypeDeviceClassificationManufacturerData]
	assert.True(s.T(), ok)
	assert.True(s.T(), fnOps.Write())
	assert.False(s.T(), fnOps.WritePartial(), "precondition: writePartial must be false")

	// Build a write datagram with a partial filter
	datagram := model.DatagramType{
		Header: model.HeaderType{
			AddressSource:      s.remoteFeature.Address(),
			AddressDestination: s.localFeature.Address(),
			MsgCounter:         util.Ptr(model.MsgCounterType(1)),
			CmdClassifier:      util.Ptr(model.CmdClassifierTypeWrite),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{
				{
					Function: util.Ptr(model.FunctionTypeDeviceClassificationManufacturerData),
					Filter: []model.FilterType{
						{
							CmdControl: &model.CmdControlType{
								Partial: &model.ElementTagType{},
							},
						},
					},
					DeviceClassificationManufacturerData: &model.DeviceClassificationManufacturerDataType{
						BrandName: util.Ptr(model.DeviceClassificationStringType("TestBrand")),
					},
				},
			},
		},
	}

	err := s.sut.ProcessCmd(datagram, s.remoteDevice)
	if assert.Error(s.T(), err) {
		assert.Contains(s.T(), err.Error(), "partial write not supported")
	}
}

// Test that a write with a delete filter to a feature that does not announce
// partial write is also rejected
func (s *DeviceLocalPartialWriteSuite) Test_DeleteFilterDenied_WhenPartialNotSupported() {
	ops := s.localFeature.Operations()
	fnOps, _ := ops[model.FunctionTypeDeviceClassificationManufacturerData]
	assert.False(s.T(), fnOps.WritePartial(), "precondition: writePartial must be false")

	datagram := model.DatagramType{
		Header: model.HeaderType{
			AddressSource:      s.remoteFeature.Address(),
			AddressDestination: s.localFeature.Address(),
			MsgCounter:         util.Ptr(model.MsgCounterType(1)),
			CmdClassifier:      util.Ptr(model.CmdClassifierTypeWrite),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{
				{
					Function: util.Ptr(model.FunctionTypeDeviceClassificationManufacturerData),
					Filter: []model.FilterType{
						{
							CmdControl: &model.CmdControlType{
								Delete: &model.ElementTagType{},
							},
						},
					},
					DeviceClassificationManufacturerData: &model.DeviceClassificationManufacturerDataType{},
				},
			},
		},
	}

	err := s.sut.ProcessCmd(datagram, s.remoteDevice)
	if assert.Error(s.T(), err) {
		assert.Contains(s.T(), err.Error(), "partial write not supported")
	}
}

// Test that a full write (no filter) to a feature that does not support
// partial write is still accepted
func (s *DeviceLocalPartialWriteSuite) Test_FullWriteAllowed_WhenPartialNotSupported() {
	datagram := model.DatagramType{
		Header: model.HeaderType{
			AddressSource:      s.remoteFeature.Address(),
			AddressDestination: s.localFeature.Address(),
			MsgCounter:         util.Ptr(model.MsgCounterType(1)),
			CmdClassifier:      util.Ptr(model.CmdClassifierTypeWrite),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{
				{
					DeviceClassificationManufacturerData: &model.DeviceClassificationManufacturerDataType{
						BrandName: util.Ptr(model.DeviceClassificationStringType("TestBrand")),
					},
				},
			},
		},
	}

	err := s.sut.ProcessCmd(datagram, s.remoteDevice)
	assert.NoError(s.T(), err)
}

// Test that a partial write to a feature that DOES support partial write
// is accepted
func (s *DeviceLocalPartialWriteSuite) Test_PartialWriteAllowed_WhenSupported() {
	// Set up a LoadControl server feature which supports partial write
	// (LoadControlLimitListDataType implements Updater)
	localEntity := s.sut.Entity([]model.AddressEntityType{1})

	lcFeature := NewFeatureLocal(2, localEntity.(*EntityLocal),
		model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	lcFeature.AddFunctionType(model.FunctionTypeLoadControlLimitListData, true, true)
	localEntity.AddFeature(lcFeature)

	// Verify precondition: writePartial is true for this function
	ops := lcFeature.Operations()
	fnOps, ok := ops[model.FunctionTypeLoadControlLimitListData]
	assert.True(s.T(), ok)
	assert.True(s.T(), fnOps.WritePartial(), "precondition: writePartial must be true")

	// Add a remote LoadControl client feature so binding type matches
	remoteEntity := s.remoteDevice.Entity([]model.AddressEntityType{1})
	remoteLcFeature := NewFeatureRemote(2, remoteEntity.(*EntityRemote),
		model.FeatureTypeTypeLoadControl, model.RoleTypeClient)
	remoteEntity.AddFeature(remoteLcFeature)

	// Add binding for the new feature pair
	binding := model.BindingManagementRequestCallType{
		ClientAddress:     remoteLcFeature.Address(),
		ServerAddress:     lcFeature.Address(),
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeLoadControl),
	}
	err := s.sut.BindingManager().AddBinding(s.remoteDevice, binding)
	assert.Nil(s.T(), err)

	datagram := model.DatagramType{
		Header: model.HeaderType{
			AddressSource:      remoteLcFeature.Address(),
			AddressDestination: lcFeature.Address(),
			MsgCounter:         util.Ptr(model.MsgCounterType(2)),
			CmdClassifier:      util.Ptr(model.CmdClassifierTypeWrite),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{
				{
					Function: util.Ptr(model.FunctionTypeLoadControlLimitListData),
					Filter:   filterEmptyPartial(),
					LoadControlLimitListData: &model.LoadControlLimitListDataType{
						LoadControlLimitData: []model.LoadControlLimitDataType{
							{
								LimitId:       util.Ptr(model.LoadControlLimitIdType(1)),
								IsLimitActive: util.Ptr(true),
							},
						},
					},
				},
			},
		},
	}

	err = s.sut.ProcessCmd(datagram, s.remoteDevice)
	assert.NoError(s.T(), err)
}
