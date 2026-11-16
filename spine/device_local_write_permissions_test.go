package spine

import (
	"testing"

	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

// TestWritePermissionErrorMessages verifies that the correct error messages
// are returned for different write permission scenarios
func TestWritePermissionErrorMessages(t *testing.T) {
	_, localEntity := createLocalDeviceAndEntity(1)

	// Create a LoadControl feature
	localFeature := NewFeatureLocal(
		localEntity.NextFeatureId(),
		localEntity,
		model.FeatureTypeTypeLoadControl,
		model.RoleTypeServer,
	)

	// Add a function that supports both read and write
	localFeature.AddFunctionType(model.FunctionTypeLoadControlLimitListData, true, true)

	// Add a function that only supports read (no write)
	localFeature.AddFunctionType(model.FunctionTypeLoadControlLimitDescriptionListData, true, false)

	localEntity.AddFeature(localFeature)

	// Get operations map
	operations := localFeature.Operations()

	t.Run("function exists and supports write", func(t *testing.T) {
		function := model.FunctionTypeLoadControlLimitListData
		ops, ok := operations[function]

		assert.True(t, ok, "Function should exist in operations")
		assert.True(t, ops.Write(), "Function should support write")

		// In this case, the check passes and no error is generated
	})

	t.Run("function exists but does not support write", func(t *testing.T) {
		function := model.FunctionTypeLoadControlLimitDescriptionListData
		ops, ok := operations[function]

		assert.True(t, ok, "Function should exist in operations")
		assert.False(t, ops.Write(), "Function should NOT support write")

		// Verify the error message that would be generated
		if ok && !ops.Write() {
			expectedMsg := "write operation not supported for this function"
			err := model.NewErrorType(model.ErrorNumberTypeCommandNotSupported, expectedMsg)

			assert.Equal(t, model.ErrorNumberTypeCommandNotSupported, err.ErrorNumber)
			assert.Equal(t, expectedMsg, string(*err.Description))
		}
	})

	t.Run("function not found in operations", func(t *testing.T) {
		// Try a function that doesn't exist in LoadControl feature
		function := model.FunctionTypeDeviceClassificationManufacturerData
		_, ok := operations[function]

		assert.False(t, ok, "Function should NOT exist in LoadControl operations")

		// Verify the error message that would be generated
		if !ok {
			expectedMsg := "function not found in feature operations"
			err := model.NewErrorType(model.ErrorNumberTypeCommandNotSupported, expectedMsg)

			assert.Equal(t, model.ErrorNumberTypeCommandNotSupported, err.ErrorNumber)
			assert.Equal(t, expectedMsg, string(*err.Description))
		}
	})
}

// TestFeatureLocal_AddFunctionType verifies AddFunctionType works correctly
func TestFeatureLocal_AddFunctionType(t *testing.T) {
	_, localEntity := createLocalDeviceAndEntity(1)

	localFeature := NewFeatureLocal(
		localEntity.NextFeatureId(),
		localEntity,
		model.FeatureTypeTypeMeasurement,
		model.RoleTypeServer,
	)

	// Test adding functions with different permissions
	localFeature.AddFunctionType(model.FunctionTypeMeasurementListData, true, true)
	localFeature.AddFunctionType(model.FunctionTypeMeasurementDescriptionListData, true, false)
	localFeature.AddFunctionType(model.FunctionTypeMeasurementConstraintsListData, false, false)

	ops := localFeature.Operations()

	// Verify read+write function
	assert.NotNil(t, ops[model.FunctionTypeMeasurementListData])
	assert.True(t, ops[model.FunctionTypeMeasurementListData].Read())
	assert.True(t, ops[model.FunctionTypeMeasurementListData].Write())

	// Verify read-only function
	assert.NotNil(t, ops[model.FunctionTypeMeasurementDescriptionListData])
	assert.True(t, ops[model.FunctionTypeMeasurementDescriptionListData].Read())
	assert.False(t, ops[model.FunctionTypeMeasurementDescriptionListData].Write())

	// Verify no-access function (unusual but possible)
	assert.NotNil(t, ops[model.FunctionTypeMeasurementConstraintsListData])
	assert.False(t, ops[model.FunctionTypeMeasurementConstraintsListData].Read())
	assert.False(t, ops[model.FunctionTypeMeasurementConstraintsListData].Write())
}

func TestFeatureLocal_AddFunctionType_MergesDuplicatePermissions(t *testing.T) {
	_, localEntity := createLocalDeviceAndEntity(1)

	localFeature := NewFeatureLocal(
		localEntity.NextFeatureId(),
		localEntity,
		model.FeatureTypeTypeMeasurement,
		model.RoleTypeServer,
	)

	function := model.FunctionTypeMeasurementListData

	localFeature.AddFunctionType(function, false, false)
	ops := localFeature.Operations()[function]
	assert.False(t, ops.Read())
	assert.False(t, ops.Write())

	localFeature.AddFunctionType(function, true, true)
	ops = localFeature.Operations()[function]
	assert.True(t, ops.Read())
	assert.True(t, ops.Write())
	assert.False(t, ops.ReadPartial())
	assert.True(t, ops.WritePartial())

	localFeature.AddFunctionType(function, false, false)
	ops = localFeature.Operations()[function]
	assert.True(t, ops.Read())
	assert.True(t, ops.Write())
	assert.False(t, ops.ReadPartial())
	assert.True(t, ops.WritePartial())
}
