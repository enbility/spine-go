package spine

import (
	"testing"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

// TestFeatureLocal_UnknownFunction_ErrorCode6 verifies that spine-go returns
// error code 6 (CommandNotSupported) for unknown functions as per SPINE specification.
// This test confirms the fix from returning error code 1 (GeneralError) to error code 6.
// TC_SPINE_COMP_001: reject a request for an unknown/unsupported Function with errorNumber > 0 (CommandNotSupported).
func TestFeatureLocal_UnknownFunction_ErrorCode6(t *testing.T) {
	// Setup
	_, localEntity := createLocalDeviceAndEntity(1)
	
	// Create a Measurement server feature
	localFeature := NewFeatureLocal(
		localEntity.NextFeatureId(),
		localEntity,
		model.FeatureTypeTypeMeasurement,
		model.RoleTypeServer,
	)
	localEntity.AddFeature(localFeature)
	
	t.Run("HandleMessage with empty cmd returns error 6", func(t *testing.T) {
		message := &api.Message{
			CmdClassifier: model.CmdClassifierTypeRead,
			Cmd:           model.CmdType{},
			RequestHeader: &model.HeaderType{
				MsgCounter: util.Ptr(model.MsgCounterType(1)),
			},
		}
		
		err := localFeature.HandleMessage(message)
		
		// Verify error code 6 is returned
		assert.NotNil(t, err)
		assert.Equal(t, model.ErrorNumberTypeCommandNotSupported, err.ErrorNumber,
			"Empty cmd should return CommandNotSupported (6), not GeneralError (1)")
		assert.Equal(t, "Data not found in Cmd", string(*err.Description))
	})
	
	t.Run("HandleMessage with no function in cmd returns error 6", func(t *testing.T) {
		// Create a cmd with valid data but nil Function after Data() processing
		message := &api.Message{
			CmdClassifier: model.CmdClassifierTypeRead,
			Cmd: model.CmdType{
				ResultData: &model.ResultDataType{}, // This will result in nil Function
			},
			RequestHeader: &model.HeaderType{
				MsgCounter: util.Ptr(model.MsgCounterType(2)),
			},
		}
		
		err := localFeature.HandleMessage(message)
		
		// Should return error 6
		assert.NotNil(t, err)
		assert.Equal(t, model.ErrorNumberTypeCommandNotSupported, err.ErrorNumber)
		assert.Equal(t, "function data not found", string(*err.Description))
	})
}

// TestFeatureLocal_processRead_UnknownFunction verifies processRead behavior
func TestFeatureLocal_processRead_UnknownFunction(t *testing.T) {
	_, localEntity := createLocalDeviceAndEntity(1)
	
	t.Run("server feature with unknown function", func(t *testing.T) {
		localFeature := NewFeatureLocal(
			localEntity.NextFeatureId(),
			localEntity,
			model.FeatureTypeTypeMeasurement,
			model.RoleTypeServer,
		)
		localEntity.AddFeature(localFeature)
		
		// Try to read an unsupported function
		err := localFeature.processRead(
			model.FunctionTypeDeviceClassificationManufacturerData, // Not supported by Measurement
			nil,
			nil,
		)
		
		// Should return error code 6
		assert.NotNil(t, err)
		assert.Equal(t, model.ErrorNumberTypeCommandNotSupported, err.ErrorNumber)
		assert.Equal(t, "function data not found", string(*err.Description))
	})
	
	t.Run("client feature rejects any read", func(t *testing.T) {
		localFeature := NewFeatureLocal(
			localEntity.NextFeatureId(),
			localEntity,
			model.FeatureTypeTypeMeasurement,
			model.RoleTypeClient,
		)
		localEntity.AddFeature(localFeature)
		
		// Client features reject all reads
		err := localFeature.processRead(
			model.FunctionTypeMeasurementListData,
			nil,
			nil,
		)
		
		// Should return error code 7 (CommandRejected)
		assert.NotNil(t, err)
		assert.Equal(t, model.ErrorNumberTypeCommandRejected, err.ErrorNumber)
	})
}

// TestFeatureLocal_executeWrite_UnknownFunction tests write handling
func TestFeatureLocal_executeWrite_UnknownFunction(t *testing.T) {
	_, localEntity := createLocalDeviceAndEntity(1)
	
	localFeature := NewFeatureLocal(
		localEntity.NextFeatureId(),
		localEntity,
		model.FeatureTypeTypeMeasurement,
		model.RoleTypeServer,
	)
	localEntity.AddFeature(localFeature)
	
	t.Run("write unknown function returns error 6", func(t *testing.T) {
		message := &api.Message{
			CmdClassifier: model.CmdClassifierTypeWrite,
			Cmd: model.CmdType{
				DeviceClassificationManufacturerData: &model.DeviceClassificationManufacturerDataType{
					DeviceName: util.Ptr(model.DeviceClassificationStringType("Test")),
				},
			},
		}
		
		err := localFeature.executeWrite(message)
		
		// Should return error code 6
		assert.NotNil(t, err)
		assert.Equal(t, model.ErrorNumberTypeCommandNotSupported, err.ErrorNumber)
		assert.Equal(t, "data not found", string(*err.Description))
	})
	
	t.Run("write empty cmd returns error 6", func(t *testing.T) {
		message := &api.Message{
			CmdClassifier: model.CmdClassifierTypeWrite,
			Cmd:           model.CmdType{},
		}
		
		err := localFeature.executeWrite(message)
		
		// Should return error code 6
		assert.NotNil(t, err)
		assert.Equal(t, model.ErrorNumberTypeCommandNotSupported, err.ErrorNumber)
		assert.Contains(t, string(*err.Description), "Data not found in Cmd")
	})
}

// TestFeatureLocal_functionData verifies functionData returns nil for unknown functions
func TestFeatureLocal_functionData(t *testing.T) {
	_, localEntity := createLocalDeviceAndEntity(1)
	
	tests := []struct {
		name         string
		featureType  model.FeatureTypeType
		knownFunc    model.FunctionType
		unknownFunc  model.FunctionType
	}{
		{
			name:        "LoadControl feature",
			featureType: model.FeatureTypeTypeLoadControl,
			knownFunc:   model.FunctionTypeLoadControlLimitListData,
			unknownFunc: model.FunctionTypeDeviceClassificationManufacturerData,
		},
		{
			name:        "Measurement feature",
			featureType: model.FeatureTypeTypeMeasurement,
			knownFunc:   model.FunctionTypeMeasurementListData,
			unknownFunc: model.FunctionTypeLoadControlLimitListData,
		},
		{
			name:        "DeviceClassification feature",
			featureType: model.FeatureTypeTypeDeviceClassification,
			knownFunc:   model.FunctionTypeDeviceClassificationManufacturerData,
			unknownFunc: model.FunctionTypeMeasurementListData,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			localFeature := NewFeatureLocal(
				localEntity.NextFeatureId(),
				localEntity,
				tt.featureType,
				model.RoleTypeServer,
			)
			localEntity.AddFeature(localFeature)
			
			// Test known function
			fd := localFeature.functionData(tt.knownFunc)
			assert.NotNil(t, fd, "functionData should return data for known function %s", tt.knownFunc)
			
			// Test unknown function
			fd = localFeature.functionData(tt.unknownFunc)
			assert.Nil(t, fd, "functionData should return nil for unknown function %s", tt.unknownFunc)
		})
	}
}

// TestErrorTypeCreation demonstrates the fix: using NewErrorType instead of NewErrorTypeFromString
func TestErrorTypeCreation(t *testing.T) {
	t.Run("NewErrorTypeFromString always creates GeneralError", func(t *testing.T) {
		// This is what was causing the problem
		err := model.NewErrorTypeFromString("function not found")
		assert.Equal(t, model.ErrorNumberTypeGeneralError, err.ErrorNumber,
			"NewErrorTypeFromString always returns GeneralError (1)")
	})
	
	t.Run("NewErrorType with explicit error code - the fix", func(t *testing.T) {
		// This is the fix: use NewErrorType with explicit error code
		err := model.NewErrorType(model.ErrorNumberTypeCommandNotSupported, "function not found")
		assert.Equal(t, model.ErrorNumberTypeCommandNotSupported, err.ErrorNumber,
			"NewErrorType with explicit code returns the correct error number")
	})
	
	// Summary: The fix was to replace NewErrorTypeFromString with NewErrorType
	// for all function-related errors to ensure error code 6 is returned
}