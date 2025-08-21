package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

// Test structures for different key configurations

// Single primary key structure
type TestSingleKeyData struct {
	Id    *uint   `eebus:"key"`
	Name  *string
	Value *int
}

// Composite key structure
type TestCompositeKeyData struct {
	Id1    *uint   `eebus:"key"`
	Id2    *string `eebus:"key"`
	Data   *string
	Status *bool
}

// Structure similar to MeasurementDataType
type TestMeasurementLikeData struct {
	MeasurementId *uint    `eebus:"key"`
	ValueType     *string  `eebus:"key"`
	Value         *float64
	State         *string
}

// Test backward compatibility for single key types
func TestSingleKeyTypes_BackwardCompatibility(t *testing.T) {
	// Single key types should still work without primarykey tag
	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		{
			name: "single_key_only",
			input: TestSingleKeyData{
				Id: util.Ptr(uint(1)),
			},
			expected: true,
		},
		{
			name: "single_key_with_data",
			input: TestSingleKeyData{
				Id:    util.Ptr(uint(1)),
				Value: util.Ptr(42),
			},
			expected: false,
		},
		{
			name: "all_fields_nil",
			input: TestSingleKeyData{},
			expected: false,
		},
		{
			name: "only_non_key_fields",
			input: TestSingleKeyData{
				Name:  util.Ptr("test"),
				Value: util.Ptr(42),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Single key types should work with hasPrimaryKeyOnly
			result := hasPrimaryKeyOnly(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test filterPrimaryKeyOnlyEntries function
func TestFilterPrimaryKeyOnlyEntries(t *testing.T) {
	t.Run("empty_slice", func(t *testing.T) {
		input := []TestSingleKeyData{}
		result := filterPrimaryKeyOnlyEntries(input)
		assert.Equal(t, input, result)
	})

	t.Run("nil_slice", func(t *testing.T) {
		var input []TestSingleKeyData
		result := filterPrimaryKeyOnlyEntries(input)
		assert.Equal(t, input, result)
	})

	t.Run("all_key_only_entries", func(t *testing.T) {
		input := []TestSingleKeyData{
			{Id: util.Ptr(uint(1))},
			{Id: util.Ptr(uint(2))},
			{Id: util.Ptr(uint(3))},
		}
		result := filterPrimaryKeyOnlyEntries(input)
		assert.Empty(t, result)
	})

	t.Run("mixed_entries", func(t *testing.T) {
		input := []TestSingleKeyData{
			{Id: util.Ptr(uint(1))},                        // Should be filtered
			{Id: util.Ptr(uint(2)), Value: util.Ptr(42)},   // Should pass
			{Id: util.Ptr(uint(3))},                        // Should be filtered
			{Id: util.Ptr(uint(4)), Name: util.Ptr("test")}, // Should pass
		}
		expected := []TestSingleKeyData{
			{Id: util.Ptr(uint(2)), Value: util.Ptr(42)},
			{Id: util.Ptr(uint(4)), Name: util.Ptr("test")},
		}
		result := filterPrimaryKeyOnlyEntries(input)
		assert.Equal(t, expected, result)
	})

	t.Run("composite_key_without_primarykey_tag", func(t *testing.T) {
		// TestCompositeKeyData doesn't have primarykey tags, so it's treated as 
		// a composite key type that hasn't been migrated - no filtering occurs
		input := []TestCompositeKeyData{
			{Id1: util.Ptr(uint(1)), Id2: util.Ptr("A")},                      
			{Id1: util.Ptr(uint(2))},                                          
			{Id1: util.Ptr(uint(3)), Id2: util.Ptr("C"), Data: util.Ptr("x")}, 
		}
		// Without primarykey tag, none are filtered
		result := filterPrimaryKeyOnlyEntries(input)
		assert.Equal(t, input, result)
	})
}

// Test UpdateList integration with primary key filtering
func TestUpdateList_WithPrimaryKeyFiltering(t *testing.T) {
	t.Run("filters_key_only_entries_before_merge", func(t *testing.T) {
		existingData := []TestSingleKeyData{
			{Id: util.Ptr(uint(1)), Name: util.Ptr("existing"), Value: util.Ptr(10)},
		}
		
		newData := []TestSingleKeyData{
			{Id: util.Ptr(uint(1))},                        // Key only - should be filtered
			{Id: util.Ptr(uint(2)), Value: util.Ptr(20)},   // New entry with data
		}

		result, success := UpdateList(false, existingData, newData, nil, nil)
		
		assert.True(t, success)
		assert.Len(t, result, 2)
		
		// Existing entry should be unchanged (key-only update was filtered)
		assert.Equal(t, util.Ptr(uint(1)), result[0].Id)
		assert.Equal(t, util.Ptr("existing"), result[0].Name)
		assert.Equal(t, util.Ptr(10), result[0].Value)
		
		// New entry should be added
		assert.Equal(t, util.Ptr(uint(2)), result[1].Id)
		assert.Equal(t, util.Ptr(20), result[1].Value)
	})

	t.Run("all_filtered_returns_existing", func(t *testing.T) {
		existingData := []TestSingleKeyData{
			{Id: util.Ptr(uint(1)), Value: util.Ptr(10)},
		}
		
		newData := []TestSingleKeyData{
			{Id: util.Ptr(uint(2))}, // Only key-only entries
			{Id: util.Ptr(uint(3))},
		}

		result, success := UpdateList(false, existingData, newData, nil, nil)
		
		assert.True(t, success)
		assert.Equal(t, existingData, result) // Should return unchanged existing data
	})
}

// Test real-world MeasurementData scenario
func TestUpdateList_MeasurementDataDuplicateIssue(t *testing.T) {
	// Start with empty data
	var existingData []MeasurementDataType

	// First message - structure only (should be filtered out)
	firstMessage := []MeasurementDataType{
		{MeasurementId: util.Ptr(MeasurementIdType(0))},
		{MeasurementId: util.Ptr(MeasurementIdType(4))},
		{MeasurementId: util.Ptr(MeasurementIdType(7))},
	}

	// Process first message
	result, success := UpdateList(false, existingData, firstMessage, nil, nil)
	assert.True(t, success)
	assert.Empty(t, result) // All entries should be filtered out

	// Second message - actual data
	secondMessage := []MeasurementDataType{
		{
			MeasurementId: util.Ptr(MeasurementIdType(4)),
			ValueType:     util.Ptr(MeasurementValueTypeType("value")),
			Value:         NewScaledNumberType(0),
			ValueSource:   util.Ptr(MeasurementValueSourceType("measuredValue")),
			ValueState:    util.Ptr(MeasurementValueStateType("normal")),
		},
	}

	// Process second message
	result, success = UpdateList(false, result, secondMessage, nil, nil)
	assert.True(t, success)
	assert.Len(t, result, 1) // Should have exactly one entry
	
	// Verify no duplicates
	assert.Equal(t, util.Ptr(MeasurementIdType(4)), result[0].MeasurementId)
	assert.Equal(t, util.Ptr(MeasurementValueTypeType("value")), result[0].ValueType)
}

// Test with actual SPINE data types
func TestUpdateList_BillDataFiltering(t *testing.T) {
	existingData := []BillDataType{
		{
			BillId:   util.Ptr(BillIdType(1)),
			BillType: util.Ptr(BillTypeType("summary")),
			ScopeType: util.Ptr(ScopeTypeType("invoice")),
		},
	}

	newData := []BillDataType{
		{BillId: util.Ptr(BillIdType(1))}, // Key only - should be filtered
		{
			BillId:   util.Ptr(BillIdType(2)),
			BillType: util.Ptr(BillTypeType("detail")),
			ScopeType: util.Ptr(ScopeTypeType("invoice")),
		},
	}

	result, success := UpdateList(false, existingData, newData, nil, nil)
	
	assert.True(t, success)
	assert.Len(t, result, 2)
	
	// First entry unchanged
	assert.Equal(t, util.Ptr(BillIdType(1)), result[0].BillId)
	assert.Equal(t, util.Ptr(BillTypeType("summary")), result[0].BillType)
	assert.Equal(t, util.Ptr(ScopeTypeType("invoice")), result[0].ScopeType)
	
	// Second entry added
	assert.Equal(t, util.Ptr(BillIdType(2)), result[1].BillId)
	assert.Equal(t, util.Ptr(BillTypeType("detail")), result[1].BillType)
}