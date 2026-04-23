package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

// Test edge cases for primary key filtering

// Test structure with slice field (not pointer)
type TestSliceFieldData struct {
	Id    *uint `eebus:"key"`
	Name  *string
	Items []string // Non-pointer slice field
}

// Test hasPrimaryKeyOnly with various field types (edge cases)
func TestHasPrimaryKeyOnly_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		// Slice field tests
		{
			name: "key_with_empty_slice",
			input: TestSliceFieldData{
				Id:    util.Ptr(uint(1)),
				Items: []string{},
			},
			expected: true, // Empty slice is considered no data
		},
		{
			name: "key_with_nil_slice",
			input: TestSliceFieldData{
				Id:    util.Ptr(uint(1)),
				Items: nil,
			},
			expected: true, // Nil slice is considered no data
		},
		{
			name: "key_with_populated_slice",
			input: TestSliceFieldData{
				Id:    util.Ptr(uint(1)),
				Items: []string{"item1", "item2"},
			},
			expected: false, // Has actual data in slice
		},
		{
			name: "key_with_name_and_items",
			input: TestSliceFieldData{
				Id:    util.Ptr(uint(1)),
				Name:  util.Ptr("test"),
				Items: []string{"item1"},
			},
			expected: false, // Has both pointer and slice data
		},

		// ElectricalConnection real-world case
		{
			name: "electrical_connection_with_data",
			input: ElectricalConnectionPermittedValueSetDataType{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(1)),
				PermittedValueSet: []ScaledNumberSetType{
					{
						Range: []ScaledNumberRangeType{
							{
								Min: NewScaledNumberType(2),
								Max: NewScaledNumberType(16),
							},
						},
					},
				},
			},
			expected: false, // Has actual data in PermittedValueSet
		},
		{
			name: "electrical_connection_keys_only",
			input: ElectricalConnectionPermittedValueSetDataType{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(1)),
			},
			expected: false, // With primarykey tag, this has sub-identifier so NOT primary key only
		},
		{
			name: "electrical_connection_with_empty_slice",
			input: ElectricalConnectionPermittedValueSetDataType{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(1)),
				PermittedValueSet:      []ScaledNumberSetType{},
			},
			expected: false, // Has sub-identifier, not just primary key
		},
		{
			name: "electrical_connection_primary_only",
			input: ElectricalConnectionPermittedValueSetDataType{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
			},
			expected: true, // Only primary key, no sub-identifier
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasPrimaryKeyOnly(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test the specific case from the issue
func TestUpdateList_RealWorldScenario_NoFilteringWithSliceData(t *testing.T) {
	// This test ensures that entries with slice data are not filtered out
	existingData := []ElectricalConnectionPermittedValueSetDataType{}

	newData := []ElectricalConnectionPermittedValueSetDataType{
		{
			ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
			ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(1)),
			PermittedValueSet: []ScaledNumberSetType{
				{
					Range: []ScaledNumberRangeType{
						{
							Min: NewScaledNumberType(2),
							Max: NewScaledNumberType(16),
						},
					},
				},
			},
		},
		{
			ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
			ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(2)),
			// No PermittedValueSet - with primarykey tag, this is NOT filtered
			// because it has both primary and sub key
		},
	}

	result, success := UpdateList(false, existingData, newData, nil, nil, nil)

	assert.True(t, success)
	assert.Len(t, result, 2) // Both entries pass through with primarykey tag
	assert.Equal(t, util.Ptr(ElectricalConnectionParameterIdType(1)), result[0].ParameterId)
	assert.Len(t, result[0].PermittedValueSet, 1) // Should have the data
	assert.Equal(t, util.Ptr(ElectricalConnectionParameterIdType(2)), result[1].ParameterId)
	assert.Empty(t, result[1].PermittedValueSet) // No data but not filtered
}
