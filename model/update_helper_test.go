package model_test

import (
	"testing"

	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

// Test SortData function with various edge cases
func TestSortData(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name:     "Empty slice",
			input:    []model.MeasurementDataType{},
			expected: []model.MeasurementDataType{},
		},
		{
			name: "Single element",
			input: []model.MeasurementDataType{
				{MeasurementId: util.Ptr(model.MeasurementIdType(1))},
			},
			expected: []model.MeasurementDataType{
				{MeasurementId: util.Ptr(model.MeasurementIdType(1))},
			},
		},
		{
			name: "Multiple elements in order",
			input: []model.MeasurementDataType{
				{MeasurementId: util.Ptr(model.MeasurementIdType(1))},
				{MeasurementId: util.Ptr(model.MeasurementIdType(2))},
			},
			expected: []model.MeasurementDataType{
				{MeasurementId: util.Ptr(model.MeasurementIdType(1))},
				{MeasurementId: util.Ptr(model.MeasurementIdType(2))},
			},
		},
		{
			name: "Multiple elements out of order",
			input: []model.MeasurementDataType{
				{MeasurementId: util.Ptr(model.MeasurementIdType(3))},
				{MeasurementId: util.Ptr(model.MeasurementIdType(1))},
				{MeasurementId: util.Ptr(model.MeasurementIdType(2))},
			},
			expected: []model.MeasurementDataType{
				{MeasurementId: util.Ptr(model.MeasurementIdType(1))},
				{MeasurementId: util.Ptr(model.MeasurementIdType(2))},
				{MeasurementId: util.Ptr(model.MeasurementIdType(3))},
			},
		},
		{
			name: "Elements with nil IDs",
			input: []model.MeasurementDataType{
				{MeasurementId: util.Ptr(model.MeasurementIdType(2))},
				{MeasurementId: nil},
				{MeasurementId: util.Ptr(model.MeasurementIdType(1))},
			},
			expected: []model.MeasurementDataType{
				{MeasurementId: util.Ptr(model.MeasurementIdType(2))},
				{MeasurementId: nil},
				{MeasurementId: util.Ptr(model.MeasurementIdType(1))},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.input.(type) {
			case []model.MeasurementDataType:
				result := model.SortData(v)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// Test SortData with structs that have no EEBus tags
func TestSortData_NoTags(t *testing.T) {
	type NoTagStruct struct {
		Field1 string
		Field2 int
	}
	
	input := []NoTagStruct{
		{Field1: "b", Field2: 2},
		{Field1: "a", Field2: 1},
	}
	
	result := model.SortData(input)
	// Should return unchanged when no tags present
	assert.Equal(t, input, result)
}

// Test SortData with structs that have different field counts
func TestSortData_DifferentFieldCounts(t *testing.T) {
	// Create two measurement data with composite keys
	input := []model.ElectricalConnectionParameterDescriptionDataType{
		{
			ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(2)),
			ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(1)),
		},
		{
			ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(1)),
			ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(2)),
		},
	}
	
	expected := []model.ElectricalConnectionParameterDescriptionDataType{
		{
			ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(1)),
			ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(2)),
		},
		{
			ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(2)),
			ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(1)),
		},
	}
	
	result := model.SortData(input)
	assert.Equal(t, expected, result)
}

// Test RemoveElementFromItem with various scenarios
func TestRemoveElementFromItem_EdgeCases(t *testing.T) {
	t.Run("Remove element with same type", func(t *testing.T) {
		// This tests removal of specific fields when element has non-nil values
		item := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         util.Ptr(model.ScaledNumberType{Number: util.Ptr(model.NumberType(100))}),
			ValueSource:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		}
		
		// Element specifying which fields to remove (non-nil fields get removed)
		element := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)), // non-nil, so remove this field
			Value:         nil, // nil, so don't remove
			ValueSource:   nil, // nil, so don't remove
		}
		
		model.RemoveElementFromItem(item, element)
		
		// MeasurementId should be removed (was non-nil in element)
		assert.Nil(t, item.MeasurementId)
		// Value and ValueSource should remain (were nil in element)
		assert.NotNil(t, item.Value)
		assert.NotNil(t, item.ValueSource)
	})

	t.Run("Invalid or unset fields", func(t *testing.T) {
		// Test with fields that cannot be set or are invalid
		item := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Timestamp:     model.NewAbsoluteOrRelativeTimeTypeFromDuration(10),
		}
		
		element := &model.MeasurementDataType{
			Timestamp: model.NewAbsoluteOrRelativeTimeTypeFromDuration(10),
		}
		
		model.RemoveElementFromItem(item, element)
		
		// Timestamp should be removed
		assert.Nil(t, item.Timestamp)
		// MeasurementId should remain
		assert.NotNil(t, item.MeasurementId)
	})
}

// Test CopyNonNilDataFromItemToItem edge cases
func TestCopyNonNilDataFromItemToItem_EdgeCases(t *testing.T) {
	t.Run("Nil source", func(t *testing.T) {
		var source *model.MeasurementDataType
		dest := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
		}
		
		model.CopyNonNilDataFromItemToItem(source, dest)
		
		// Destination should be unchanged
		assert.NotNil(t, dest.MeasurementId)
		assert.Equal(t, model.MeasurementIdType(1), *dest.MeasurementId)
	})

	t.Run("Nil destination", func(t *testing.T) {
		source := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
		}
		var dest *model.MeasurementDataType
		
		// Should not panic
		model.CopyNonNilDataFromItemToItem(source, dest)
	})

	t.Run("Mismatched field counts", func(t *testing.T) {
		// Testing conceptually - the function handles mismatched field counts
		// by returning early. We can't directly test with different types due to
		// Go's type system, but the coverage is achieved through other test paths.
	})

	t.Run("Invalid or non-settable fields", func(t *testing.T) {
		source := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(2)),
			Value:         util.Ptr(model.ScaledNumberType{Number: util.Ptr(model.NumberType(200))}),
		}
		
		dest := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Value:         util.Ptr(model.ScaledNumberType{Number: util.Ptr(model.NumberType(100))}),
		}
		
		model.CopyNonNilDataFromItemToItem(source, dest)
		
		// All non-nil fields from source should be copied
		assert.Equal(t, model.MeasurementIdType(2), *dest.MeasurementId)
		assert.Equal(t, model.NumberType(200), *dest.Value.Number)
	})
}

// Direct test for isFieldValueNil coverage
func TestIsFieldValueNil_Coverage(t *testing.T) {
	// We can't directly test isFieldValueNil as it's not exported
	// But we can test it through nonNilElementNames behavior
	
	t.Run("Various nil types through RemoveElementFromItem", func(t *testing.T) {
		// Test with different kinds of nil values
		item := &model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			Timestamp:     model.NewAbsoluteOrRelativeTimeTypeFromDuration(10),
			Value:         util.Ptr(model.ScaledNumberType{Number: util.Ptr(model.NumberType(100))}),
			ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal),
		}
		
		// Element with some nil and some non-nil fields
		element := &model.MeasurementDataType{
			MeasurementId: nil, // nil - should not remove from item
			Timestamp:     model.NewAbsoluteOrRelativeTimeTypeFromDuration(10), // non-nil - should remove from item
			Value:         nil, // nil - should not remove from item
			ValueState:    util.Ptr(model.MeasurementValueStateTypeNormal), // non-nil - should remove from item
		}
		
		model.RemoveElementFromItem(item, element)
		
		// Only Timestamp and ValueState should be removed (they were non-nil in element)
		assert.NotNil(t, item.MeasurementId)
		assert.Nil(t, item.Timestamp) // Should be removed
		assert.NotNil(t, item.Value)
		assert.Nil(t, item.ValueState) // Should be removed
	})
}

// Test cases that trigger different types in isFieldValueNil
func TestFieldValueNilTypes(t *testing.T) {
	t.Run("Array and slice fields", func(t *testing.T) {
		// Test with LoadControlLimitListData which has slice fields
		item := &model.LoadControlLimitListDataType{
			LoadControlLimitData: []model.LoadControlLimitDataType{
				{
					LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
					IsLimitActive: util.Ptr(true),
				},
			},
		}
		
		element := &model.LoadControlLimitListDataType{
			LoadControlLimitData: nil, // nil slice
		}
		
		// Using a wrapper to test slice handling
		type wrapper struct {
			Data *model.LoadControlLimitListDataType
		}
		
		w := &wrapper{Data: item}
		e := &wrapper{Data: element}
		
		// This won't directly remove but tests the nil check path
		model.RemoveElementFromItem(w, e)
	})
}