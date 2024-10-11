package spine

import (
	"testing"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func TestFunctionDataFactory_FunctionData(t *testing.T) {
	result := CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeBill)
	assert.Equal(t, 3, len(result))
	assert.IsType(t, &FunctionData[model.BillDescriptionListDataType]{}, result[0])
	assert.IsType(t, &FunctionData[model.BillConstraintsListDataType]{}, result[1])
	assert.IsType(t, &FunctionData[model.BillListDataType]{}, result[2])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeDeviceClassification)
	assert.Equal(t, 2, len(result))
	assert.IsType(t, &FunctionData[model.DeviceClassificationManufacturerDataType]{}, result[0])
	assert.IsType(t, &FunctionData[model.DeviceClassificationUserDataType]{}, result[1])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeDeviceConfiguration)
	assert.Equal(t, 3, len(result))
	assert.IsType(t, &FunctionData[model.DeviceConfigurationKeyValueConstraintsListDataType]{}, result[0])
	assert.IsType(t, &FunctionData[model.DeviceConfigurationKeyValueDescriptionListDataType]{}, result[1])
	assert.IsType(t, &FunctionData[model.DeviceConfigurationKeyValueListDataType]{}, result[2])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeDeviceDiagnosis)
	assert.Equal(t, 3, len(result))
	assert.IsType(t, &FunctionData[model.DeviceDiagnosisStateDataType]{}, result[0])
	assert.IsType(t, &FunctionData[model.DeviceDiagnosisHeartbeatDataType]{}, result[1])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeElectricalConnection)
	assert.Equal(t, 6, len(result))
	assert.IsType(t, &FunctionData[model.ElectricalConnectionDescriptionListDataType]{}, result[0])
	assert.IsType(t, &FunctionData[model.ElectricalConnectionParameterDescriptionListDataType]{}, result[1])
	assert.IsType(t, &FunctionData[model.ElectricalConnectionPermittedValueSetListDataType]{}, result[2])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeHvac)
	assert.Equal(t, 8, len(result))
	assert.IsType(t, &FunctionData[model.HvacOperationModeDescriptionListDataType]{}, result[0])
	assert.IsType(t, &FunctionData[model.HvacOverrunDescriptionListDataType]{}, result[1])
	assert.IsType(t, &FunctionData[model.HvacOverrunListDataType]{}, result[2])
	assert.IsType(t, &FunctionData[model.HvacSystemFunctionDescriptionDataType]{}, result[3])
	assert.IsType(t, &FunctionData[model.HvacSystemFunctionListDataType]{}, result[4])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeIdentification)
	assert.Equal(t, 3, len(result))
	assert.IsType(t, &FunctionData[model.IdentificationListDataType]{}, result[0])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeIncentiveTable)
	assert.Equal(t, 3, len(result))
	assert.IsType(t, &FunctionData[model.IncentiveTableDescriptionDataType]{}, result[0])
	assert.IsType(t, &FunctionData[model.IncentiveTableConstraintsDataType]{}, result[1])
	assert.IsType(t, &FunctionData[model.IncentiveTableDataType]{}, result[2])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeLoadControl)
	assert.Equal(t, 6, len(result))
	assert.IsType(t, &FunctionData[model.LoadControlEventListDataType]{}, result[0])
	assert.IsType(t, &FunctionData[model.LoadControlLimitConstraintsListDataType]{}, result[1])
	assert.IsType(t, &FunctionData[model.LoadControlLimitDescriptionListDataType]{}, result[2])
	assert.IsType(t, &FunctionData[model.LoadControlLimitListDataType]{}, result[3])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeMeasurement)
	assert.Equal(t, 5, len(result))
	assert.IsType(t, &FunctionData[model.MeasurementListDataType]{}, result[0])
	assert.IsType(t, &FunctionData[model.MeasurementDescriptionListDataType]{}, result[1])
	assert.IsType(t, &FunctionData[model.MeasurementConstraintsListDataType]{}, result[2])
	assert.IsType(t, &FunctionData[model.MeasurementThresholdRelationListDataType]{}, result[3])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeTimeSeries)
	assert.Equal(t, 3, len(result))
	assert.IsType(t, &FunctionData[model.TimeSeriesDescriptionListDataType]{}, result[0])
	assert.IsType(t, &FunctionData[model.TimeSeriesConstraintsListDataType]{}, result[1])
	assert.IsType(t, &FunctionData[model.TimeSeriesListDataType]{}, result[2])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeGeneric)
	assert.Equal(t, 124, len(result))
}

func TestFunctionDataFactory_FunctionDataCmd(t *testing.T) {
	result := CreateFunctionData[api.FunctionDataCmdInterface](model.FeatureTypeTypeDeviceClassification)
	assert.Equal(t, 2, len(result))
	assert.IsType(t, &FunctionDataCmd[model.DeviceClassificationManufacturerDataType]{}, result[0])
	assert.IsType(t, &FunctionDataCmd[model.DeviceClassificationUserDataType]{}, result[1])
}

func TestFunctionDataFactory_NodeMgmtFeatureType(t *testing.T) {
	result := CreateFunctionData[api.FunctionDataCmdInterface](model.FeatureTypeTypeNodeManagement)
	assert.Equal(t, 3, len(result))
}

func TestFunctionDataFactory_unknownFunctionDataType(t *testing.T) {
	assert.PanicsWithError(t, "only FunctionData and FunctionDataCmd are supported",
		func() { CreateFunctionData[int](model.FeatureTypeTypeDeviceClassification) })
}
