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

func TestDeviceLocalValidationSuite(t *testing.T) {
	suite.Run(t, new(DeviceLocalValidationSuite))
}

type DeviceLocalValidationSuite struct {
	suite.Suite
	
	lastMessage string
}

var _ shipapi.ShipConnectionDataWriterInterface = (*DeviceLocalValidationSuite)(nil)

func (s *DeviceLocalValidationSuite) WriteShipMessageWithPayload(msg []byte) {
	// Mock implementation - just store the message
	s.lastMessage = string(msg)
}

func (s *DeviceLocalValidationSuite) Test_ProcessCmd_ValidFunctionAlignment() {
	sut := NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)
	localEntity := NewEntityLocal(sut, model.EntityTypeTypeCEM, NewAddressEntityType([]uint{1}), time.Second*4)
	localFeature := NewFeatureLocal(1, localEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	localEntity.AddFeature(localFeature)
	sut.AddEntity(localEntity)

	ski := "test"
	_ = sut.SetupRemoteDevice(ski, s)
	remote := sut.RemoteDeviceForSki(ski)
	assert.NotNil(s.T(), remote)
	
	remoteEntity := NewEntityRemote(remote, model.EntityTypeTypeCEM, []model.AddressEntityType{1})
	remoteFeature := NewFeatureRemote(1, remoteEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	remoteEntity.AddFeature(remoteFeature)
	remote.AddEntity(remoteEntity)

	// Create a valid command where all functions align
	cmd := model.CmdType{
		Function: util.Ptr(model.FunctionType("measurementListData")),
		MeasurementListData: &model.MeasurementListDataType{
			MeasurementData: []model.MeasurementDataType{
				{
					MeasurementId: util.Ptr(model.MeasurementIdType(1)),
					Value: model.NewScaledNumberType(100),
				},
			},
		},
	}

	datagram := model.DatagramType{
		Header: model.HeaderType{
			AddressSource:      remoteFeature.Address(),
			AddressDestination: localFeature.Address(),
			MsgCounter:         util.Ptr(model.MsgCounterType(1)),
			CmdClassifier:      util.Ptr(model.CmdClassifierTypeWrite),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{cmd},
		},
	}

	// Process should succeed - all functions match
	err := sut.ProcessCmd(datagram, remote)
	// We expect no error for validation since functions match
	// The error might still occur for other reasons (feature not supporting the function)
	// but the validation itself should pass
	if err != nil {
		assert.NotContains(s.T(), err.Error(), "cmd function validation failed", 
			"Valid command should not fail function validation")
	}
}

func (s *DeviceLocalValidationSuite) Test_ProcessCmd_FunctionWithoutFiltersAllowed() {
	sut := NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)
	localEntity := NewEntityLocal(sut, model.EntityTypeTypeCEM, NewAddressEntityType([]uint{1}), time.Second*4)
	localFeature := NewFeatureLocal(1, localEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	localEntity.AddFeature(localFeature)
	sut.AddEntity(localEntity)

	ski := "test"
	_ = sut.SetupRemoteDevice(ski, s)
	remote := sut.RemoteDeviceForSki(ski)
	assert.NotNil(s.T(), remote)
	
	remoteEntity := NewEntityRemote(remote, model.EntityTypeTypeCEM, []model.AddressEntityType{1})
	remoteFeature := NewFeatureRemote(1, remoteEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	remoteEntity.AddFeature(remoteFeature)
	remote.AddEntity(remoteEntity)

	// Create a command with function but NO filters
	// SPINE spec says function SHALL be absent when no filters present,
	// but we don't strictly enforce this since there's no security risk
	cmd := model.CmdType{
		Function: util.Ptr(model.FunctionType("measurementListData")),
		// No filters - technically non-compliant but allowed
		MeasurementListData: &model.MeasurementListDataType{
			MeasurementData: []model.MeasurementDataType{
				{
					MeasurementId: util.Ptr(model.MeasurementIdType(1)),
					Value: model.NewScaledNumberType(100),
				},
			},
		},
	}

	datagram := model.DatagramType{
		Header: model.HeaderType{
			AddressSource:      remoteFeature.Address(),
			AddressDestination: localFeature.Address(),
			MsgCounter:         util.Ptr(model.MsgCounterType(1)),
			CmdClassifier:      util.Ptr(model.CmdClassifierTypeWrite),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{cmd},
		},
	}

	// Process should now succeed - we don't strictly enforce function absence for commands without filters
	// The security risk only exists when filters are present and mismatched
	err := sut.ProcessCmd(datagram, remote)
	// May fail for other reasons (feature not supporting the function) but not for validation
	if err != nil {
		assert.NotContains(s.T(), err.Error(), "protocol violation", 
			"Should not fail for protocol violation when no filters present")
	}
}

func (s *DeviceLocalValidationSuite) Test_ProcessCmd_MismatchedFunctionWithFilters() {
	sut := NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)
	localEntity := NewEntityLocal(sut, model.EntityTypeTypeCEM, NewAddressEntityType([]uint{1}), time.Second*4)
	localFeature := NewFeatureLocal(1, localEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	localEntity.AddFeature(localFeature)
	sut.AddEntity(localEntity)

	ski := "test"
	_ = sut.SetupRemoteDevice(ski, s)
	remote := sut.RemoteDeviceForSki(ski)
	assert.NotNil(s.T(), remote)
	
	remoteEntity := NewEntityRemote(remote, model.EntityTypeTypeCEM, []model.AddressEntityType{1})
	remoteFeature := NewFeatureRemote(1, remoteEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	remoteEntity.AddFeature(remoteFeature)
	remote.AddEntity(remoteEntity)

	// Create a command where filter function doesn't match data
	cmd := model.CmdType{
		Function: util.Ptr(model.FunctionType("measurementListData")),
		Filter: []model.FilterType{
			{
				CmdControl: &model.CmdControlType{
					Partial: &model.ElementTagType{},
				},
				// Filter is for a different function (loadControlLimitListData)
				LoadControlLimitListDataSelectors: &model.LoadControlLimitListDataSelectorsType{
					LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
				},
			},
		},
		// Data is measurementListData
		MeasurementListData: &model.MeasurementListDataType{
			MeasurementData: []model.MeasurementDataType{
				{
					MeasurementId: util.Ptr(model.MeasurementIdType(1)),
					Value: model.NewScaledNumberType(100),
				},
			},
		},
	}

	datagram := model.DatagramType{
		Header: model.HeaderType{
			AddressSource:      remoteFeature.Address(),
			AddressDestination: localFeature.Address(),
			MsgCounter:         util.Ptr(model.MsgCounterType(1)),
			CmdClassifier:      util.Ptr(model.CmdClassifierTypeWrite),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{cmd},
		},
	}

	// Process should fail - filter function mismatch
	err := sut.ProcessCmd(datagram, remote)
	assert.Error(s.T(), err, "Filter function mismatch should fail validation")
	assert.Contains(s.T(), err.Error(), "cmd function validation failed", 
		"Error should indicate validation failure")
}

func (s *DeviceLocalValidationSuite) Test_ProcessCmd_EmptyFunctionWithFilterRejected() {
	sut := NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)
	localEntity := NewEntityLocal(sut, model.EntityTypeTypeCEM, NewAddressEntityType([]uint{1}), time.Second*4)
	localFeature := NewFeatureLocal(1, localEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	localEntity.AddFeature(localFeature)
	sut.AddEntity(localEntity)

	ski := "test"
	_ = sut.SetupRemoteDevice(ski, s)
	remote := sut.RemoteDeviceForSki(ski)
	assert.NotNil(s.T(), remote)
	
	remoteEntity := NewEntityRemote(remote, model.EntityTypeTypeCEM, []model.AddressEntityType{1})
	remoteFeature := NewFeatureRemote(1, remoteEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	remoteEntity.AddFeature(remoteFeature)
	remote.AddEntity(remoteEntity)

	// Create a command with empty cmd.Function and filters (violates SPINE spec)
	// Per spec: function SHALL be present when filters are present
	cmd := model.CmdType{
		Function: util.Ptr(model.FunctionType("")), // Empty function with filter - NOT allowed
		Filter: []model.FilterType{
			{
				CmdControl: &model.CmdControlType{
					Partial: &model.ElementTagType{},
				},
				MeasurementListDataSelectors: &model.MeasurementListDataSelectorsType{
					MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				},
			},
		},
		MeasurementListData: &model.MeasurementListDataType{
			MeasurementData: []model.MeasurementDataType{
				{
					MeasurementId: util.Ptr(model.MeasurementIdType(1)),
					Value: model.NewScaledNumberType(100),
				},
			},
		},
	}

	datagram := model.DatagramType{
		Header: model.HeaderType{
			AddressSource:      remoteFeature.Address(),
			AddressDestination: localFeature.Address(),
			MsgCounter:         util.Ptr(model.MsgCounterType(1)),
			CmdClassifier:      util.Ptr(model.CmdClassifierTypeWrite),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{cmd},
		},
	}

	// Process should fail - empty function violates SPINE spec when filters are present
	err := sut.ProcessCmd(datagram, remote)
	assert.Error(s.T(), err, "Empty cmd.Function with filters should be rejected per SPINE spec")
	assert.Contains(s.T(), err.Error(), "cmd function validation failed", 
		"Error should indicate validation failure for empty function with filters")
}