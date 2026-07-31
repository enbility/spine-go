package model_test

import (
	"testing"

	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

// Test RFE partial write with ScaledNumberType:
// When Scale is omitted from incoming data, it should be preserved from existing data
func TestRFE_ScaledNumberType_OmitScale_PreservesExisting(t *testing.T) {
	// Existing entry: value = 100 * 10^-1 = 10.0
	existingData := []model.LoadControlLimitDataType{
		{
			LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
			IsLimitChangeable: util.Ptr(true),
			Value: &model.ScaledNumberType{
				Number: util.Ptr(model.NumberType(100)),
				Scale:  util.Ptr(model.ScaleType(-1)), // Scale is -1
			},
		},
	}

	// Incoming partial write: new Number, but Scale omitted (nil)
	// This simulates a partial write where scale is not provided
	updateData := []model.LoadControlLimitDataType{
		{
			LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
			Value: &model.ScaledNumberType{
				Number: util.Ptr(model.NumberType(250)), // New number
				Scale:  nil, // Scale is omitted!
			},
		},
	}

	// Process without partial filter (direct merge)
	result, success := model.UpdateList(false, existingData, updateData, nil, nil, nil)

	assert.True(t, success, "update should succeed")
	assert.NotNil(t, result)
	assert.Len(t, result, 1)

	// Check if Scale was preserved
	resultEntry := result[0]
	assert.NotNil(t, resultEntry.Value, "Value should not be nil")
	assert.NotNil(t, resultEntry.Value.Number, "Number should not be nil")
	assert.NotNil(t, resultEntry.Value.Scale, "Scale should be preserved from existing data")

	if resultEntry.Value.Scale != nil {
		t.Logf("Result: Number=%d, Scale=%d", *resultEntry.Value.Number, *resultEntry.Value.Scale)
		if *resultEntry.Value.Scale == -1 {
			t.Logf("✓ Scale preserved correctly as -1")
		} else {
			t.Errorf("✗ Scale was not preserved: got %d, expected -1", *resultEntry.Value.Scale)
		}
	}
}

// Test with partial filter (selector-based merge)
func TestRFE_ScaledNumberType_PartialFilter_OmitScale(t *testing.T) {
	existingData := []model.LoadControlLimitDataType{
		{
			LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
			IsLimitChangeable: util.Ptr(true),
			Value: &model.ScaledNumberType{
				Number: util.Ptr(model.NumberType(500)),
				Scale:  util.Ptr(model.ScaleType(-2)), // Existing scale
			},
		},
	}

	updateData := []model.LoadControlLimitDataType{
		{
			Value: &model.ScaledNumberType{
				Number: util.Ptr(model.NumberType(1000)),
				Scale:  nil, // Omitted
			},
		},
	}

	// Create partial filter targeting the limit
	filterPartial := model.NewFilterTypePartial()
	filterPartial.LoadControlLimitListDataSelectors = &model.LoadControlLimitListDataSelectorsType{
		LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
	}

	result, success := model.UpdateList(false, existingData, updateData, filterPartial, nil, nil)

	assert.True(t, success, "update should succeed")
	assert.NotNil(t, result)
	assert.Len(t, result, 1)

	resultEntry := result[0]
	if resultEntry.Value != nil && resultEntry.Value.Scale != nil {
		t.Logf("Partial filter result: Number=%d, Scale=%d", 
			*resultEntry.Value.Number, *resultEntry.Value.Scale)
	}
}
