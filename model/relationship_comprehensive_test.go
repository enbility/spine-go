package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewRelationships_TimeSeries tests TimeSeries relationships
func TestNewRelationships_TimeSeries(t *testing.T) {
	data := TimeSeriesDescriptionDataType{}
	rels := GetRelationships(data)

	assert.Len(t, rels, 1)
	assert.Equal(t, "MeasurementId", rels[0].FieldName)
	assert.Equal(t, "MeasurementDescriptionDataType", rels[0].TargetType)
}

// TestNewRelationships_LoadControlState tests LoadControl state relationships
func TestNewRelationships_LoadControlState(t *testing.T) {
	data := LoadControlStateDataType{}
	rels := GetRelationships(data)

	assert.Len(t, rels, 1)
	assert.Equal(t, "EventId", rels[0].FieldName)
	assert.Equal(t, "LoadControlEventDataType", rels[0].TargetType)
	assert.True(t, rels[0].IsComposite)
}

// TestNewRelationships_LoadControlConstraints tests LoadControl constraints relationships
func TestNewRelationships_LoadControlConstraints(t *testing.T) {
	data := LoadControlLimitConstraintsDataType{}
	rels := GetRelationships(data)

	assert.Len(t, rels, 1)
	assert.Equal(t, "LimitId", rels[0].FieldName)
	assert.Equal(t, "LoadControlLimitDescriptionDataType", rels[0].TargetType)
	assert.True(t, rels[0].IsComposite)
}

// TestNewRelationships_Alarm tests Alarm relationships
func TestNewRelationships_Alarm(t *testing.T) {
	data := AlarmDataType{}
	rels := GetRelationships(data)

	assert.Len(t, rels, 1)
	assert.Equal(t, "ThresholdId", rels[0].FieldName)
	assert.Equal(t, "ThresholdDescriptionDataType", rels[0].TargetType)
}

// TestNewRelationships_Bill tests Bill relationships
func TestNewRelationships_Bill(t *testing.T) {
	data := BillDescriptionDataType{}
	rels := GetRelationships(data)

	assert.Len(t, rels, 1)
	assert.Equal(t, "SessionId", rels[0].FieldName)
	assert.Equal(t, "SessionIdentificationDataType", rels[0].TargetType)
}

// TestNewRelationships_TariffDescription tests Tariff description relationships
func TestNewRelationships_TariffDescription(t *testing.T) {
	data := TariffDescriptionDataType{}
	rels := GetRelationships(data)

	assert.Len(t, rels, 1)
	assert.Equal(t, "MeasurementId", rels[0].FieldName)
	assert.Equal(t, "MeasurementDescriptionDataType", rels[0].TargetType)
}

// TestNewRelationships_TierBoundaryDescription tests complex multi-reference relationships
func TestNewRelationships_TierBoundaryDescription(t *testing.T) {
	data := TierBoundaryDescriptionDataType{}
	rels := GetRelationships(data)

	// Should have 3 TierId references
	assert.Len(t, rels, 3)

	// All should reference TierDescriptionDataType.TierId
	for _, rel := range rels {
		assert.Equal(t, "TierDescriptionDataType", rel.TargetType)
		assert.Equal(t, "TierId", rel.TargetField)
	}

	// Check specific field names
	fieldNames := []string{}
	for _, rel := range rels {
		fieldNames = append(fieldNames, rel.FieldName)
	}
	assert.Contains(t, fieldNames, "ValidForTierId")
	assert.Contains(t, fieldNames, "SwitchToTierIdWhenLower")
	assert.Contains(t, fieldNames, "SwitchToTierIdWhenHigher")
}

// TestNewRelationships_TierBoundaryData tests TierBoundary data relationships
func TestNewRelationships_TierBoundaryData(t *testing.T) {
	data := TierBoundaryDataType{}
	rels := GetRelationships(data)

	assert.Len(t, rels, 1)
	assert.Equal(t, "TimeTableId", rels[0].FieldName)
	assert.Equal(t, "TimeTableDescriptionDataType", rels[0].TargetType)
}

// TestNewRelationships_SessionIdentification tests SessionIdentification relationships
func TestNewRelationships_SessionIdentification(t *testing.T) {
	data := SessionIdentificationDataType{}
	rels := GetRelationships(data)

	assert.Len(t, rels, 1)
	assert.Equal(t, "IdentificationId", rels[0].FieldName)
	assert.Equal(t, "IdentificationDataType", rels[0].TargetType)
}

// TestNewRelationships_HvacSystemFunction tests HVAC relationships
func TestNewRelationships_HvacSystemFunction(t *testing.T) {
	data := HvacSystemFunctionDataType{}
	rels := GetRelationships(data)

	assert.Len(t, rels, 2)

	// Check both relationships exist
	var foundOperationMode, foundSetpoint bool
	for _, rel := range rels {
		if rel.FieldName == "CurrentOperationModeId" {
			assert.Equal(t, "HvacOperationModeDescriptionDataType", rel.TargetType)
			foundOperationMode = true
		}
		if rel.FieldName == "CurrentSetpointId" {
			assert.Equal(t, "SetpointDescriptionDataType", rel.TargetType)
			foundSetpoint = true
		}
	}
	assert.True(t, foundOperationMode, "Should have CurrentOperationModeId relationship")
	assert.True(t, foundSetpoint, "Should have CurrentSetpointId relationship")
}

// TestNewRelationships_SupplyCondition tests SupplyCondition relationships
func TestNewRelationships_SupplyCondition(t *testing.T) {
	data := SupplyConditionDataType{}
	rels := GetRelationships(data)

	assert.Len(t, rels, 1)
	assert.Equal(t, "ThresholdId", rels[0].FieldName)
	assert.Equal(t, "ThresholdDescriptionDataType", rels[0].TargetType)
}

// TestNewRelationships_TaskManagement tests TaskManagement cross-feature relationships
func TestNewRelationships_TaskManagement(t *testing.T) {
	// Test HVAC-related task
	hvacData := TaskManagementHvacRelatedType{}
	hvacRels := GetRelationships(hvacData)
	assert.Len(t, hvacRels, 1)
	assert.Equal(t, "OverrunId", hvacRels[0].FieldName)
	assert.Equal(t, "HvacOverrunDescriptionDataType", hvacRels[0].TargetType)

	// Test LoadControl-related task
	lcData := TaskManagementLoadControlReleatedType{}
	lcRels := GetRelationships(lcData)
	assert.Len(t, lcRels, 1)
	assert.Equal(t, "EventId", lcRels[0].FieldName)
	assert.Equal(t, "LoadControlEventDataType", lcRels[0].TargetType)

	// Test PowerSequences-related task
	psData := TaskManagementPowerSequencesRelatedType{}
	psRels := GetRelationships(psData)
	assert.Len(t, psRels, 1)
	assert.Equal(t, "SequenceId", psRels[0].FieldName)
	assert.Equal(t, "PowerSequenceDescriptionDataType", psRels[0].TargetType)
}
