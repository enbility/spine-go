package model

import (
	"reflect"
	"strings"
)

// RelationshipInfo describes a foreign key relationship from one type to another
type RelationshipInfo struct {
	// FieldName is the local field that contains the foreign key
	FieldName string

	// TargetType is the name of the target type being referenced
	TargetType string

	// TargetField is the field name in the target type that this foreign key references
	TargetField string

	// IsComposite indicates if this is part of a composite foreign key
	IsComposite bool
}

// GetRelationships extracts all relationship metadata from a type's struct tags
// It looks for fields with eebus:"ref:TargetType.TargetField" tags
//
// Example usage:
//
//	type LoadControlLimitDescriptionDataType struct {
//	    LimitId       *LoadControlLimitIdType `eebus:"key,primarykey"`
//	    MeasurementId *MeasurementIdType      `eebus:"ref:MeasurementDescriptionDataType.MeasurementId"`
//	}
//
//	rels := GetRelationships(LoadControlLimitDescriptionDataType{})
//	// Returns: [{FieldName: "MeasurementId", TargetType: "MeasurementDescriptionDataType", TargetField: "MeasurementId"}]
func GetRelationships(dataType any) []RelationshipInfo {
	var result []RelationshipInfo

	t := reflect.TypeOf(dataType)
	if t == nil {
		return result
	}

	// Handle pointer types
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return result
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tags := EEBusTags(field)

		// Check for ref tag
		refValue, hasRef := tags[EEBusTagRef]
		if !hasRef || refValue == "" || refValue == "true" {
			continue
		}

		// Parse ref value: "TargetType.TargetField"
		parts := strings.Split(refValue, ".")
		if len(parts) != 2 {
			continue
		}

		// Check if this is part of a composite key
		_, hasKey := tags[EEBusTagKey]
		_, hasPrimaryKey := tags[EEBusTagPrimaryKey]
		isComposite := hasKey || hasPrimaryKey

		result = append(result, RelationshipInfo{
			FieldName:   field.Name,
			TargetType:  parts[0],
			TargetField: parts[1],
			IsComposite: isComposite,
		})
	}

	return result
}

// GetRelationshipsByFieldName returns relationship info for a specific field
func GetRelationshipsByFieldName(dataType any, fieldName string) *RelationshipInfo {
	rels := GetRelationships(dataType)
	for _, rel := range rels {
		if rel.FieldName == fieldName {
			return &rel
		}
	}
	return nil
}

// HasRelationships returns true if the type has any relationship metadata
func HasRelationships(dataType any) bool {
	return len(GetRelationships(dataType)) > 0
}
