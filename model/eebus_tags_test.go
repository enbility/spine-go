package model

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test structs for tag testing
type TestTagStruct struct {
	SimpleKey        *string `eebus:"key"`
	PrimaryKey       *uint   `eebus:"key,primarykey"`
	WriteCheckField  *bool   `eebus:"writecheck"`
	FunctionField    *string `eebus:"fct"`
	TypeField        *string `eebus:"typ"`
	NoEEBusTag       *string
	EmptyEEBusTag    *string `eebus:""`
	MultipleFlags    *string `eebus:"key,writecheck"`
	ValuePairTag     *string `eebus:"fct:measurement"`
	MalformedTag     *string `eebus:"bad:tag:format:too:many"`
	ComplexTag       *string `eebus:"key,fct:test,writecheck"`
}

func TestEEBusTags_EmptyTag(t *testing.T) {
	field := reflect.TypeOf(TestTagStruct{}).Field(5) // NoEEBusTag
	result := EEBusTags(field)
	
	assert.Empty(t, result)
}

func TestEEBusTags_EmptyEEBusTag(t *testing.T) {
	field := reflect.TypeOf(TestTagStruct{}).Field(6) // EmptyEEBusTag
	result := EEBusTags(field)
	
	assert.Empty(t, result)
}

func TestEEBusTags_SimpleKey(t *testing.T) {
	field := reflect.TypeOf(TestTagStruct{}).Field(0) // SimpleKey
	result := EEBusTags(field)
	
	expected := map[EEBusTag]string{
		EEBusTagKey: "true",
	}
	assert.Equal(t, expected, result)
}

func TestEEBusTags_PrimaryKey(t *testing.T) {
	field := reflect.TypeOf(TestTagStruct{}).Field(1) // PrimaryKey
	result := EEBusTags(field)
	
	expected := map[EEBusTag]string{
		EEBusTagKey:        "true",
		EEBusTagPrimaryKey: "true",
	}
	assert.Equal(t, expected, result)
}

func TestEEBusTags_WriteCheck(t *testing.T) {
	field := reflect.TypeOf(TestTagStruct{}).Field(2) // WriteCheckField
	result := EEBusTags(field)
	
	expected := map[EEBusTag]string{
		EEBusTagWriteCheck: "true",
	}
	assert.Equal(t, expected, result)
}

func TestEEBusTags_Function(t *testing.T) {
	field := reflect.TypeOf(TestTagStruct{}).Field(3) // FunctionField
	result := EEBusTags(field)
	
	expected := map[EEBusTag]string{
		EEBusTagFunction: "true",
	}
	assert.Equal(t, expected, result)
}

func TestEEBusTags_Type(t *testing.T) {
	field := reflect.TypeOf(TestTagStruct{}).Field(4) // TypeField
	result := EEBusTags(field)
	
	expected := map[EEBusTag]string{
		EEBusTagType: "true",
	}
	assert.Equal(t, expected, result)
}

func TestEEBusTags_MultipleFlags(t *testing.T) {
	field := reflect.TypeOf(TestTagStruct{}).Field(7) // MultipleFlags
	result := EEBusTags(field)
	
	expected := map[EEBusTag]string{
		EEBusTagKey:        "true",
		EEBusTagWriteCheck: "true",
	}
	assert.Equal(t, expected, result)
}

func TestEEBusTags_ValuePair(t *testing.T) {
	field := reflect.TypeOf(TestTagStruct{}).Field(8) // ValuePairTag
	result := EEBusTags(field)
	
	expected := map[EEBusTag]string{
		EEBusTagFunction: "measurement",
	}
	assert.Equal(t, expected, result)
}

func TestEEBusTags_MalformedTag(t *testing.T) {
	field := reflect.TypeOf(TestTagStruct{}).Field(9) // MalformedTag
	result := EEBusTags(field)
	
	// Should still process the valid parts and ignore malformed parts
	// The function logs an error but doesn't fail
	assert.Empty(t, result) // Malformed tag is ignored
}

func TestEEBusTags_ComplexTag(t *testing.T) {
	field := reflect.TypeOf(TestTagStruct{}).Field(10) // ComplexTag
	result := EEBusTags(field)
	
	expected := map[EEBusTag]string{
		EEBusTagKey:        "true",
		EEBusTagFunction:   "test",
		EEBusTagWriteCheck: "true",
	}
	assert.Equal(t, expected, result)
}

func TestEEBusTags_AllTags(t *testing.T) {
	// Test all defined EEBus tag constants
	tests := []struct {
		name     string
		tag      string
		expected map[EEBusTag]string
	}{
		{
			name: "all boolean tags",
			tag:  `eebus:"key,primarykey,writecheck,fct,typ"`,
			expected: map[EEBusTag]string{
				EEBusTagKey:        "true",
				EEBusTagPrimaryKey: "true",
				EEBusTagWriteCheck: "true",
				EEBusTagFunction:   "true",
				EEBusTagType:       "true",
			},
		},
		{
			name: "mixed value and boolean tags",
			tag:  `eebus:"key,fct:measurement,primarykey,typ:selector"`,
			expected: map[EEBusTag]string{
				EEBusTagKey:        "true",
				EEBusTagFunction:   "measurement",
				EEBusTagPrimaryKey: "true",
				EEBusTagType:       "selector",
			},
		},
		{
			name: "only value tags",
			tag:  `eebus:"fct:loadcontrol,typ:elements"`,
			expected: map[EEBusTag]string{
				EEBusTagFunction: "loadcontrol",
				EEBusTagType:     "elements",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a struct field dynamically with the test tag
			structType := reflect.StructOf([]reflect.StructField{
				{
					Name: "TestField",
					Type: reflect.TypeOf((*string)(nil)),
					Tag:  reflect.StructTag(tt.tag),
				},
			})
			field := structType.Field(0)
			
			result := EEBusTags(field)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEEBusTagConstants(t *testing.T) {
	// Test that all constants are defined correctly
	assert.Equal(t, EEBusTag("fct"), EEBusTagFunction)
	assert.Equal(t, EEBusTag("typ"), EEBusTagType)
	assert.Equal(t, EEBusTag("key"), EEBusTagKey)
	assert.Equal(t, EEBusTag("primarykey"), EEBusTagPrimaryKey)
	assert.Equal(t, EEBusTag("writecheck"), EEBusTagWriteCheck)
	
	assert.Equal(t, "eebus", EEBusTagName)
	
	assert.Equal(t, EEBusTagTypeType("selector"), EEBusTagTypeTypeSelector)
	assert.Equal(t, EEBusTagTypeType("elements"), EEbusTagTypeTypeElements)
}

func TestEEBusTags_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		expected map[EEBusTag]string
	}{
		{
			name:     "whitespace in tags",
			tag:      `eebus:" key , primarykey "`,
			expected: map[EEBusTag]string{
				EEBusTag(" key "): "true",
				EEBusTag(" primarykey "): "true",
			},
		},
		{
			name:     "empty value pair",
			tag:      `eebus:"fct:"`,
			expected: map[EEBusTag]string{
				EEBusTagFunction: "",
			},
		},
		{
			name:     "colon but no value",
			tag:      `eebus:"key,fct:,primarykey"`,
			expected: map[EEBusTag]string{
				EEBusTagKey:        "true",
				EEBusTagFunction:   "",
				EEBusTagPrimaryKey: "true",
			},
		},
		{
			name:     "duplicate tags",
			tag:      `eebus:"key,key,primarykey"`,
			expected: map[EEBusTag]string{
				EEBusTagKey:        "true", // Last one wins
				EEBusTagPrimaryKey: "true",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			structType := reflect.StructOf([]reflect.StructField{
				{
					Name: "TestField",
					Type: reflect.TypeOf((*string)(nil)),
					Tag:  reflect.StructTag(tt.tag),
				},
			})
			field := structType.Field(0)
			
			result := EEBusTags(field)
			assert.Equal(t, tt.expected, result)
		})
	}
}