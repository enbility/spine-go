package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRelationships_LoadControlLimitDescription(t *testing.T) {
	data := LoadControlLimitDescriptionDataType{}
	rels := GetRelationships(data)

	assert.Len(t, rels, 1, "LoadControlLimitDescriptionDataType should have 1 relationship")
	assert.Equal(t, "MeasurementId", rels[0].FieldName)
	assert.Equal(t, "MeasurementDescriptionDataType", rels[0].TargetType)
	assert.Equal(t, "MeasurementId", rels[0].TargetField)
	assert.False(t, rels[0].IsComposite, "MeasurementId is not part of a composite key")
}

func TestGetRelationships_LoadControlLimitData(t *testing.T) {
	data := LoadControlLimitDataType{}
	rels := GetRelationships(data)

	assert.Len(t, rels, 1, "LoadControlLimitDataType should have 1 relationship")
	assert.Equal(t, "LimitId", rels[0].FieldName)
	assert.Equal(t, "LoadControlLimitDescriptionDataType", rels[0].TargetType)
	assert.Equal(t, "LimitId", rels[0].TargetField)
	assert.True(t, rels[0].IsComposite, "LimitId is a primary key")
}

func TestGetRelationships_MeasurementData(t *testing.T) {
	data := MeasurementDataType{}
	rels := GetRelationships(data)

	assert.Len(t, rels, 1, "MeasurementDataType should have 1 relationship")
	assert.Equal(t, "MeasurementId", rels[0].FieldName)
	assert.Equal(t, "MeasurementDescriptionDataType", rels[0].TargetType)
	assert.Equal(t, "MeasurementId", rels[0].TargetField)
	assert.True(t, rels[0].IsComposite, "MeasurementId is a primary key")
}

func TestGetRelationships_ElectricalConnectionParameterDescription(t *testing.T) {
	data := ElectricalConnectionParameterDescriptionDataType{}
	rels := GetRelationships(data)

	assert.Len(t, rels, 2, "ElectricalConnectionParameterDescriptionDataType should have 2 relationships")

	// Check both relationships exist
	var foundElectricalConnectionId, foundMeasurementId bool
	for _, rel := range rels {
		if rel.FieldName == "ElectricalConnectionId" {
			assert.Equal(t, "ElectricalConnectionDescriptionDataType", rel.TargetType)
			assert.Equal(t, "ElectricalConnectionId", rel.TargetField)
			assert.True(t, rel.IsComposite, "ElectricalConnectionId is part of the composite key")
			foundElectricalConnectionId = true
		}
		if rel.FieldName == "MeasurementId" {
			assert.Equal(t, "MeasurementDescriptionDataType", rel.TargetType)
			assert.Equal(t, "MeasurementId", rel.TargetField)
			assert.False(t, rel.IsComposite, "MeasurementId is not part of the composite key")
			foundMeasurementId = true
		}
	}
	assert.True(t, foundElectricalConnectionId, "Should have ElectricalConnectionId relationship")
	assert.True(t, foundMeasurementId, "Should have MeasurementId relationship")
}

func TestGetRelationships_ElectricalConnectionCharacteristic_CompositeKey(t *testing.T) {
	data := ElectricalConnectionCharacteristicDataType{}
	rels := GetRelationships(data)

	assert.Len(t, rels, 2, "ElectricalConnectionCharacteristicDataType should have 2 relationship fields (composite foreign key)")

	// Check ElectricalConnectionId relationship
	var ecIdRel *RelationshipInfo
	for i := range rels {
		if rels[i].FieldName == "ElectricalConnectionId" {
			ecIdRel = &rels[i]
			break
		}
	}
	assert.NotNil(t, ecIdRel, "Should have ElectricalConnectionId relationship")
	assert.Equal(t, "ElectricalConnectionParameterDescriptionDataType", ecIdRel.TargetType)
	assert.Equal(t, "ElectricalConnectionId", ecIdRel.TargetField)
	assert.True(t, ecIdRel.IsComposite, "ElectricalConnectionId is part of composite key")

	// Check ParameterId relationship
	var paramIdRel *RelationshipInfo
	for i := range rels {
		if rels[i].FieldName == "ParameterId" {
			paramIdRel = &rels[i]
			break
		}
	}
	assert.NotNil(t, paramIdRel, "Should have ParameterId relationship")
	assert.Equal(t, "ElectricalConnectionParameterDescriptionDataType", paramIdRel.TargetType)
	assert.Equal(t, "ParameterId", paramIdRel.TargetField)
	assert.True(t, paramIdRel.IsComposite, "ParameterId is part of composite key")
}

func TestGetRelationships_NoRelationships(t *testing.T) {
	// MeasurementDescriptionDataType is a target, not a source, so it should have no relationships
	data := MeasurementDescriptionDataType{}
	rels := GetRelationships(data)

	assert.Len(t, rels, 0, "MeasurementDescriptionDataType should have no outgoing relationships")
}

func TestGetRelationshipsByFieldName(t *testing.T) {
	data := LoadControlLimitDescriptionDataType{}

	// Test existing field
	rel := GetRelationshipsByFieldName(data, "MeasurementId")
	assert.NotNil(t, rel)
	assert.Equal(t, "MeasurementDescriptionDataType", rel.TargetType)

	// Test non-existent field
	rel = GetRelationshipsByFieldName(data, "NonExistentField")
	assert.Nil(t, rel)
}

func TestHasRelationships(t *testing.T) {
	// Type with relationships
	data := LoadControlLimitDescriptionDataType{}
	assert.True(t, HasRelationships(data))

	// Type without relationships
	desc := MeasurementDescriptionDataType{}
	assert.False(t, HasRelationships(desc))
}

func TestGetRelationships_WithPointer(t *testing.T) {
	// Test that function works with pointer types too
	data := &LoadControlLimitDescriptionDataType{}
	rels := GetRelationships(data)

	assert.Len(t, rels, 1)
	assert.Equal(t, "MeasurementId", rels[0].FieldName)
}

func TestGetRelationships_InvalidType(t *testing.T) {
	// Test with nil
	rels := GetRelationships(nil)
	assert.Len(t, rels, 0)

	// Test with non-struct type
	rels = GetRelationships("not a struct")
	assert.Len(t, rels, 0)
}
