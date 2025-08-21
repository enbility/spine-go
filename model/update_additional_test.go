package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

// Additional test data types for comprehensive testing
type TestCompositeKeyWithPrimaryTag struct {
	PrimaryId *uint   `eebus:"key,primarykey"`
	SubId     *string `eebus:"key"`
	Value     *int
	Metadata  *string
}

type TestNoKeyData struct {
	Value *int
	Name  *string
}

func TestFieldNamesWithEEBusTag(t *testing.T) {
	tests := []struct {
		name     string
		tag      EEBusTag
		item     interface{}
		expected []string
	}{
		{
			name:     "find key fields in composite key struct",
			tag:      EEBusTagKey,
			item:     TestCompositeKeyWithPrimaryTag{},
			expected: []string{"PrimaryId", "SubId"},
		},
		{
			name:     "find primarykey fields in composite key struct",
			tag:      EEBusTagPrimaryKey,
			item:     TestCompositeKeyWithPrimaryTag{},
			expected: []string{"PrimaryId"},
		},
		{
			name:     "find key fields in no-key struct",
			tag:      EEBusTagKey,
			item:     TestNoKeyData{},
			expected: []string{},
		},
		{
			name:     "find writecheck fields",
			tag:      EEBusTagWriteCheck,
			item:     TestUpdateData{},
			expected: []string{"IsChangeable"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fieldNamesWithEEBusTag(tt.tag, tt.item)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

func TestFieldNamesWithEEBusTag_NonStruct(t *testing.T) {
	// Test with non-struct types
	result := fieldNamesWithEEBusTag(EEBusTagKey, "not a struct")
	assert.Empty(t, result)

	result = fieldNamesWithEEBusTag(EEBusTagKey, 42)
	assert.Empty(t, result)

	result = fieldNamesWithEEBusTag(EEBusTagKey, nil)
	assert.Empty(t, result)
}

func TestHasIdentifiers(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		expected bool
	}{
		{
			name:     "has all identifiers",
			data:     TestCompositeKeyWithPrimaryTag{PrimaryId: util.Ptr(uint(1)), SubId: util.Ptr("test")},
			expected: true,
		},
		{
			name:     "missing one identifier",
			data:     TestCompositeKeyWithPrimaryTag{PrimaryId: util.Ptr(uint(1))}, // SubId is nil
			expected: false,
		},
		{
			name:     "no key fields",
			data:     TestNoKeyData{Value: util.Ptr(1)},
			expected: true, // No key fields means no requirement
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasIdentifiers(tt.data)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHasPrimaryKeyOnly_CompositeKeyWithPrimaryTag(t *testing.T) {
	tests := []struct {
		name     string
		item     interface{}
		expected bool
	}{
		{
			name:     "composite key - only primary key",
			item:     TestCompositeKeyWithPrimaryTag{PrimaryId: util.Ptr(uint(1))},
			expected: true,
		},
		{
			name:     "composite key - primary key plus sub key",
			item:     TestCompositeKeyWithPrimaryTag{PrimaryId: util.Ptr(uint(1)), SubId: util.Ptr("test")},
			expected: false,
		},
		{
			name:     "composite key - primary key plus data",
			item:     TestCompositeKeyWithPrimaryTag{PrimaryId: util.Ptr(uint(1)), Value: util.Ptr(100)},
			expected: false,
		},
		{
			name:     "composite key - no fields",
			item:     TestCompositeKeyWithPrimaryTag{},
			expected: false,
		},
		{
			name:     "no key fields",
			item:     TestNoKeyData{Value: util.Ptr(1)},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasPrimaryKeyOnly(tt.item)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHasOnlySingleKey(t *testing.T) {
	// Use the existing TestSingleKeyData from update_primary_key_filter_test.go
	tests := []struct {
		name     string
		item     interface{}
		keyField string
		expected bool
	}{
		{
			name:     "only key field has value",
			item:     TestSingleKeyData{Id: util.Ptr(uint(1))},
			keyField: "Id",
			expected: true,
		},
		{
			name:     "key field plus other field",
			item:     TestSingleKeyData{Id: util.Ptr(uint(1)), Value: util.Ptr(100)},
			keyField: "Id",
			expected: false,
		},
		{
			name:     "no key field value",
			item:     TestSingleKeyData{Value: util.Ptr(100)},
			keyField: "Id",
			expected: false,
		},
		{
			name:     "empty struct",
			item:     TestSingleKeyData{},
			keyField: "Id",
			expected: false,
		},
		{
			name:     "non-struct item",
			item:     "not a struct",
			keyField: "Id",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasOnlySingleKey(tt.item, tt.keyField)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHasPrimaryKeyOnlyNew(t *testing.T) {
	primaryKeyFields := []string{"PrimaryId"}

	tests := []struct {
		name     string
		item     interface{}
		expected bool
	}{
		{
			name:     "only primary key field",
			item:     TestCompositeKeyWithPrimaryTag{PrimaryId: util.Ptr(uint(1))},
			expected: true,
		},
		{
			name:     "primary key plus sub key",
			item:     TestCompositeKeyWithPrimaryTag{PrimaryId: util.Ptr(uint(1)), SubId: util.Ptr("test")},
			expected: false,
		},
		{
			name:     "primary key plus data field",
			item:     TestCompositeKeyWithPrimaryTag{PrimaryId: util.Ptr(uint(1)), Value: util.Ptr(100)},
			expected: false,
		},
		{
			name:     "no primary key",
			item:     TestCompositeKeyWithPrimaryTag{SubId: util.Ptr("test")},
			expected: false,
		},
		{
			name:     "empty struct",
			item:     TestCompositeKeyWithPrimaryTag{},
			expected: false,
		},
		{
			name:     "non-struct",
			item:     "not a struct",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasPrimaryKeyOnlyNew(tt.item, primaryKeyFields)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilterPrimaryKeyOnlyEntries_CompositeKeyWithPrimaryTag(t *testing.T) {
	tests := []struct {
		name           string
		data           []TestCompositeKeyWithPrimaryTag
		expectedResult []TestCompositeKeyWithPrimaryTag
		expectedLog    bool // Whether debug log should be called
	}{
		{
			name:           "empty data",
			data:           []TestCompositeKeyWithPrimaryTag{},
			expectedResult: []TestCompositeKeyWithPrimaryTag{},
			expectedLog:    false,
		},
		{
			name: "no primary-key-only entries",
			data: []TestCompositeKeyWithPrimaryTag{
				{PrimaryId: util.Ptr(uint(1)), SubId: util.Ptr("test"), Value: util.Ptr(100)},
				{PrimaryId: util.Ptr(uint(2)), Value: util.Ptr(200)},
			},
			expectedResult: []TestCompositeKeyWithPrimaryTag{
				{PrimaryId: util.Ptr(uint(1)), SubId: util.Ptr("test"), Value: util.Ptr(100)},
				{PrimaryId: util.Ptr(uint(2)), Value: util.Ptr(200)},
			},
			expectedLog: false,
		},
		{
			name: "mixed entries",
			data: []TestCompositeKeyWithPrimaryTag{
				{PrimaryId: util.Ptr(uint(1))}, // Only primary key - should be filtered
				{PrimaryId: util.Ptr(uint(2)), Value: util.Ptr(200)}, // Has data - should remain
				{PrimaryId: util.Ptr(uint(3))}, // Only primary key - should be filtered
			},
			expectedResult: []TestCompositeKeyWithPrimaryTag{
				{PrimaryId: util.Ptr(uint(2)), Value: util.Ptr(200)},
			},
			expectedLog: true,
		},
		{
			name: "all primary-key-only entries",
			data: []TestCompositeKeyWithPrimaryTag{
				{PrimaryId: util.Ptr(uint(1))},
				{PrimaryId: util.Ptr(uint(2))},
			},
			expectedResult: nil, // filterPrimaryKeyOnlyEntries returns nil slice when all are filtered
			expectedLog:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterPrimaryKeyOnlyEntries(tt.data)
			if tt.expectedResult == nil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestSortData(t *testing.T) {
	tests := []struct {
		name     string
		data     []TestSingleKeyData
		expected []TestSingleKeyData
	}{
		{
			name:     "empty data",
			data:     []TestSingleKeyData{},
			expected: []TestSingleKeyData{},
		},
		{
			name: "already sorted",
			data: []TestSingleKeyData{
				{Id: util.Ptr(uint(1)), Value: util.Ptr(100)},
				{Id: util.Ptr(uint(2)), Value: util.Ptr(200)},
			},
			expected: []TestSingleKeyData{
				{Id: util.Ptr(uint(1)), Value: util.Ptr(100)},
				{Id: util.Ptr(uint(2)), Value: util.Ptr(200)},
			},
		},
		{
			name: "unsorted data",
			data: []TestSingleKeyData{
				{Id: util.Ptr(uint(3)), Value: util.Ptr(300)},
				{Id: util.Ptr(uint(1)), Value: util.Ptr(100)},
				{Id: util.Ptr(uint(2)), Value: util.Ptr(200)},
			},
			expected: []TestSingleKeyData{
				{Id: util.Ptr(uint(1)), Value: util.Ptr(100)},
				{Id: util.Ptr(uint(2)), Value: util.Ptr(200)},
				{Id: util.Ptr(uint(3)), Value: util.Ptr(300)},
			},
		},
		{
			name: "data without keys - no sorting",
			data: []TestSingleKeyData{
				{Value: util.Ptr(300)},
				{Value: util.Ptr(100)},
			},
			expected: []TestSingleKeyData{
				{Value: util.Ptr(300)}, // Order preserved when no keys
				{Value: util.Ptr(100)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortData(tt.data)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSortData_NoKeyStruct(t *testing.T) {
	// Test with struct that has no key fields
	data := []TestNoKeyData{
		{Value: util.Ptr(3), Name: util.Ptr("third")},
		{Value: util.Ptr(1), Name: util.Ptr("first")},
	}
	expected := data // Should remain unchanged
	
	result := SortData(data)
	assert.Equal(t, expected, result)
}

func TestCopyNonNilDataFromItemToItem(t *testing.T) {
	// Use the existing TestSingleKeyData from update_primary_key_filter_test.go
	tests := []struct {
		name        string
		source      *TestSingleKeyData
		destination *TestSingleKeyData
		expected    *TestSingleKeyData
	}{
		{
			name:        "nil source",
			source:      nil,
			destination: &TestSingleKeyData{Id: util.Ptr(uint(1))},
			expected:    &TestSingleKeyData{Id: util.Ptr(uint(1))}, // Unchanged
		},
		{
			name:        "nil destination",
			source:      &TestSingleKeyData{Id: util.Ptr(uint(1))},
			destination: nil,
			expected:    nil,
		},
		{
			name:        "copy non-nil fields only",
			source:      &TestSingleKeyData{Id: util.Ptr(uint(2))}, // Value is nil
			destination: &TestSingleKeyData{Id: util.Ptr(uint(1)), Value: util.Ptr(100)},
			expected:    &TestSingleKeyData{Id: util.Ptr(uint(2)), Value: util.Ptr(100)}, // Only Id copied
		},
		{
			name:        "copy all non-nil fields",
			source:      &TestSingleKeyData{Id: util.Ptr(uint(2)), Value: util.Ptr(200)},
			destination: &TestSingleKeyData{Id: util.Ptr(uint(1)), Value: util.Ptr(100)},
			expected:    &TestSingleKeyData{Id: util.Ptr(uint(2)), Value: util.Ptr(200)}, // All copied
		},
		{
			name:        "empty source",
			source:      &TestSingleKeyData{},
			destination: &TestSingleKeyData{Id: util.Ptr(uint(1)), Value: util.Ptr(100)},
			expected:    &TestSingleKeyData{Id: util.Ptr(uint(1)), Value: util.Ptr(100)}, // Unchanged
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CopyNonNilDataFromItemToItem(tt.source, tt.destination)
			if tt.expected == nil {
				assert.Nil(t, tt.destination)
			} else {
				assert.Equal(t, tt.expected, tt.destination)
			}
		})
	}
}

func TestUpdateList_PrimaryKeyFiltering(t *testing.T) {
	// Use real SPINE types to test the actual functionality
	existingData := []MeasurementDataType{
		{
			MeasurementId: util.Ptr(MeasurementIdType(1)),
			ValueType:     util.Ptr(MeasurementValueTypeType("power")),
			Value:         NewScaledNumberType(100),
		},
	}

	// Mix of valid and primary-key-only entries
	newData := []MeasurementDataType{
		{MeasurementId: util.Ptr(MeasurementIdType(1))}, // Primary key only - should be filtered
		{
			MeasurementId: util.Ptr(MeasurementIdType(2)),
			ValueType:     util.Ptr(MeasurementValueTypeType("voltage")),
			Value:         NewScaledNumberType(220),
		}, // Valid - should be added
	}

	result, success := UpdateList(false, existingData, newData, nil, nil)
	assert.True(t, success)
	assert.Len(t, result, 2)
	
	// Check first item - should be unchanged since primary-key-only update was filtered
	assert.Equal(t, util.Ptr(MeasurementIdType(1)), result[0].MeasurementId)
	assert.Equal(t, util.Ptr(MeasurementValueTypeType("power")), result[0].ValueType)
	assert.NotNil(t, result[0].Value)
	
	// Check second item - should be newly added
	assert.Equal(t, util.Ptr(MeasurementIdType(2)), result[1].MeasurementId)
	assert.Equal(t, util.Ptr(MeasurementValueTypeType("voltage")), result[1].ValueType)
	assert.NotNil(t, result[1].Value)
}

func TestUpdateList_AllPrimaryKeyOnly(t *testing.T) {
	// Use simple single-key type 
	existingData := []TestSingleKeyData{
		{Id: util.Ptr(uint(1)), Value: util.Ptr(100)},
	}

	// All entries are key-only
	newData := []TestSingleKeyData{
		{Id: util.Ptr(uint(1))},
		{Id: util.Ptr(uint(2))},
	}

	// Should return existing data unchanged since all new data was filtered
	result, success := UpdateList(false, existingData, newData, nil, nil)
	assert.True(t, success)
	assert.Equal(t, existingData, result)
}

func TestIsFieldValueNil(t *testing.T) {
	tests := []struct {
		name     string
		field    interface{}
		expected bool
	}{
		{
			name:     "nil pointer",
			field:    (*string)(nil),
			expected: true,
		},
		{
			name:     "non-nil pointer",
			field:    util.Ptr("test"),
			expected: false,
		},
		{
			name:     "nil slice",
			field:    ([]string)(nil),
			expected: true,
		},
		{
			name:     "empty slice",
			field:    []string{},
			expected: false,
		},
		{
			name:     "nil map",
			field:    (map[string]int)(nil),
			expected: true,
		},
		{
			name:     "primitive value",
			field:    42,
			expected: false,
		},
		{
			name:     "nil interface",
			field:    nil,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isFieldValueNil(tt.field)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test hasPrimaryKeyOnly with different field types
func TestHasPrimaryKeyOnly_FieldTypes(t *testing.T) {
	type TestFieldTypes struct {
		PrimaryId *uint                 `eebus:"key,primarykey"`
		StringVal *string               // Pointer type
		SliceVal  []string              // Slice type
		MapVal    map[string]int        // Map type
		IntVal    int                   // Non-pointer type
	}

	tests := []struct {
		name     string
		item     TestFieldTypes
		expected bool
	}{
		{
			name:     "only primary key",
			item:     TestFieldTypes{PrimaryId: util.Ptr(uint(1))},
			expected: true,
		},
		{
			name:     "primary key + string pointer",
			item:     TestFieldTypes{PrimaryId: util.Ptr(uint(1)), StringVal: util.Ptr("test")},
			expected: false,
		},
		{
			name:     "primary key + slice",
			item:     TestFieldTypes{PrimaryId: util.Ptr(uint(1)), SliceVal: []string{"test"}},
			expected: false,
		},
		{
			name:     "primary key + map",
			item:     TestFieldTypes{PrimaryId: util.Ptr(uint(1)), MapVal: map[string]int{"test": 1}},
			expected: false,
		},
		{
			name:     "primary key + non-zero int",
			item:     TestFieldTypes{PrimaryId: util.Ptr(uint(1)), IntVal: 42},
			expected: false,
		},
		{
			name:     "primary key + zero int (zero value treated as no value)",
			item:     TestFieldTypes{PrimaryId: util.Ptr(uint(1)), IntVal: 0},
			expected: true, // Zero is the zero value for int, so it's treated as "no value"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasPrimaryKeyOnly(tt.item)
			assert.Equal(t, tt.expected, result)
		})
	}
}