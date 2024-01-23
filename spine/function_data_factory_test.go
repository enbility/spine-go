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
	assert.IsType(t, &FunctionDataImpl[model.BillDescriptionListDataType]{}, result[0])
	assert.IsType(t, &FunctionDataImpl[model.BillConstraintsListDataType]{}, result[1])
	assert.IsType(t, &FunctionDataImpl[model.BillListDataType]{}, result[2])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeDeviceClassification)
	assert.Equal(t, 2, len(result))
	assert.IsType(t, &FunctionDataImpl[model.DeviceClassificationManufacturerDataType]{}, result[0])
	assert.IsType(t, &FunctionDataImpl[model.DeviceClassificationUserDataType]{}, result[1])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeDeviceConfiguration)
	assert.Equal(t, 3, len(result))
	assert.IsType(t, &FunctionDataImpl[model.DeviceConfigurationKeyValueConstraintsListDataType]{}, result[0])
	assert.IsType(t, &FunctionDataImpl[model.DeviceConfigurationKeyValueDescriptionListDataType]{}, result[1])
	assert.IsType(t, &FunctionDataImpl[model.DeviceConfigurationKeyValueListDataType]{}, result[2])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeDeviceDiagnosis)
	assert.Equal(t, 3, len(result))
	assert.IsType(t, &FunctionDataImpl[model.DeviceDiagnosisStateDataType]{}, result[0])
	assert.IsType(t, &FunctionDataImpl[model.DeviceDiagnosisHeartbeatDataType]{}, result[1])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeElectricalConnection)
	assert.Equal(t, 6, len(result))
	assert.IsType(t, &FunctionDataImpl[model.ElectricalConnectionDescriptionListDataType]{}, result[0])
	assert.IsType(t, &FunctionDataImpl[model.ElectricalConnectionParameterDescriptionListDataType]{}, result[1])
	assert.IsType(t, &FunctionDataImpl[model.ElectricalConnectionPermittedValueSetListDataType]{}, result[2])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeHvac)
	assert.Equal(t, 8, len(result))
	assert.IsType(t, &FunctionDataImpl[model.HvacOperationModeDescriptionDataType]{}, result[0])
	assert.IsType(t, &FunctionDataImpl[model.HvacOverrunDescriptionListDataType]{}, result[1])
	assert.IsType(t, &FunctionDataImpl[model.HvacOverrunListDataType]{}, result[2])
	assert.IsType(t, &FunctionDataImpl[model.HvacSystemFunctionDescriptionDataType]{}, result[3])
	assert.IsType(t, &FunctionDataImpl[model.HvacSystemFunctionListDataType]{}, result[4])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeIdentification)
	assert.Equal(t, 3, len(result))
	assert.IsType(t, &FunctionDataImpl[model.IdentificationListDataType]{}, result[0])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeIncentiveTable)
	assert.Equal(t, 3, len(result))
	assert.IsType(t, &FunctionDataImpl[model.IncentiveTableDescriptionDataType]{}, result[0])
	assert.IsType(t, &FunctionDataImpl[model.IncentiveTableConstraintsDataType]{}, result[1])
	assert.IsType(t, &FunctionDataImpl[model.IncentiveTableDataType]{}, result[2])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeLoadControl)
	assert.Equal(t, 6, len(result))
	assert.IsType(t, &FunctionDataImpl[model.LoadControlEventListDataType]{}, result[0])
	assert.IsType(t, &FunctionDataImpl[model.LoadControlLimitConstraintsListDataType]{}, result[1])
	assert.IsType(t, &FunctionDataImpl[model.LoadControlLimitDescriptionListDataType]{}, result[2])
	assert.IsType(t, &FunctionDataImpl[model.LoadControlLimitListDataType]{}, result[3])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeMeasurement)
	assert.Equal(t, 5, len(result))
	assert.IsType(t, &FunctionDataImpl[model.MeasurementListDataType]{}, result[0])
	assert.IsType(t, &FunctionDataImpl[model.MeasurementDescriptionListDataType]{}, result[1])
	assert.IsType(t, &FunctionDataImpl[model.MeasurementConstraintsListDataType]{}, result[2])
	assert.IsType(t, &FunctionDataImpl[model.MeasurementThresholdRelationListDataType]{}, result[3])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeTimeSeries)
	assert.Equal(t, 3, len(result))
	assert.IsType(t, &FunctionDataImpl[model.TimeSeriesDescriptionListDataType]{}, result[0])
	assert.IsType(t, &FunctionDataImpl[model.TimeSeriesConstraintsListDataType]{}, result[1])
	assert.IsType(t, &FunctionDataImpl[model.TimeSeriesListDataType]{}, result[2])

	result = CreateFunctionData[api.FunctionDataInterface](model.FeatureTypeTypeGeneric)
	assert.Equal(t, 124, len(result))
}

func TestFunctionDataFactory_FunctionDataCmd(t *testing.T) {
	result := CreateFunctionData[api.FunctionDataCmdInterface](model.FeatureTypeTypeDeviceClassification)
	assert.Equal(t, 2, len(result))
	assert.IsType(t, &FunctionDataCmdImpl[model.DeviceClassificationManufacturerDataType]{}, result[0])
	assert.IsType(t, &FunctionDataCmdImpl[model.DeviceClassificationUserDataType]{}, result[1])
}

func TestFunctionDataFactory_NodeMgmtFeatureType(t *testing.T) {
	result := CreateFunctionData[api.FunctionDataCmdInterface](model.FeatureTypeTypeNodeManagement)
	assert.Equal(t, 3, len(result))
}

func TestFunctionDataFactory_unknownFunctionDataType(t *testing.T) {
	assert.PanicsWithError(t, "only FunctionData and FunctionDataCmd are supported",
		func() { CreateFunctionData[int](model.FeatureTypeTypeDeviceClassification) })
}
