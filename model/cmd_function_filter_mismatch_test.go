package model_test

import (
	"testing"

	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

// TestCmdFunctionFilterMismatch demonstrates the critical security issue where
// cmd.Function doesn't match the function in the filter
func TestCmdFunctionFilterMismatch(t *testing.T) {
	t.Run("Mismatched function - Security Issue", func(t *testing.T) {
		// Create a CmdType with mismatched function and filter
		// This demonstrates the vulnerability: cmd.Function says one thing,
		// but the filter contains data for a different function
		cmd := model.CmdType{
			// This says we're dealing with measurement data
			Function: util.Ptr(model.FunctionType("measurementListData")),
			
			// But the filter has LoadControl selectors!
			Filter: []model.FilterType{
				{
					CmdControl: &model.CmdControlType{
						Partial: &model.ElementTagType{},
					},
					// This is for load control, NOT measurement!
					LoadControlLimitListDataSelectors: &model.LoadControlLimitListDataSelectorsType{
						LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
					},
				},
			},
			
			// And we have measurement data
			MeasurementListData: &model.MeasurementListDataType{
				MeasurementData: []model.MeasurementDataType{
					{
						MeasurementId: util.Ptr(model.MeasurementIdType(1)),
						Value:         model.NewScaledNumberType(100),
					},
				},
			},
		}

		// Extract the function from cmd.Data()
		cmdData, err := cmd.Data()
		assert.NoError(t, err)
		assert.NotNil(t, cmdData)
		assert.NotNil(t, cmdData.Function)
		
		// The cmd.Data() correctly identifies this as measurementListData
		assert.Equal(t, model.FunctionType("measurementListData"), *cmdData.Function)
		
		// But what about the filter?
		filterData, err := cmd.Filter[0].Data(cmd.Function)
		assert.NoError(t, err)
		assert.NotNil(t, filterData)
		assert.NotNil(t, filterData.Function)
		
		// The filter thinks this is loadControlLimitListData!
		assert.Equal(t, model.FunctionType("loadControlLimitListData"), *filterData.Function)
		
		// SECURITY ISSUE: These should match but they don't!
		assert.NotEqual(t, *cmdData.Function, *filterData.Function,
			"CRITICAL: cmd.Function (%s) does not match filter function (%s)",
			*cmdData.Function, *filterData.Function)
		
		// This could lead to:
		// 1. Wrong data being processed with wrong filters
		// 2. Type confusion vulnerabilities
		// 3. Data integrity issues
		// 4. Potential crashes or undefined behavior
	})

	t.Run("Multiple filters with different functions", func(t *testing.T) {
		// Even worse: multiple filters with different functions
		cmd := model.CmdType{
			Function: util.Ptr(model.FunctionType("measurementListData")),
			
			Filter: []model.FilterType{
				{
					CmdControl: &model.CmdControlType{
						Delete: &model.ElementTagType{},
					},
					// Delete filter for load control
					LoadControlLimitListDataSelectors: &model.LoadControlLimitListDataSelectorsType{
						LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
					},
				},
				{
					CmdControl: &model.CmdControlType{
						Partial: &model.ElementTagType{},
					},
					// Partial filter for electrical connection
					ElectricalConnectionParameterDescriptionListDataSelectors: &model.ElectricalConnectionParameterDescriptionListDataSelectorsType{
						ParameterId: util.Ptr(model.ElectricalConnectionParameterIdType(1)),
					},
				},
			},
			
			MeasurementListData: &model.MeasurementListDataType{
				MeasurementData: []model.MeasurementDataType{
					{
						MeasurementId: util.Ptr(model.MeasurementIdType(1)),
						Value:         model.NewScaledNumberType(100),
					},
				},
			},
		}

		// Check each filter
		for i, filter := range cmd.Filter {
			filterData, err := filter.Data(cmd.Function)
			assert.NoError(t, err)
			assert.NotNil(t, filterData)
			
			cmdData, _ := cmd.Data()
			if i == 0 {
				assert.Equal(t, model.FunctionType("loadControlLimitListData"), *filterData.Function,
					"Filter %d has wrong function type", i)
			} else {
				assert.Equal(t, model.FunctionType("electricalConnectionParameterDescriptionListData"), *filterData.Function,
					"Filter %d has wrong function type", i)
			}
			
			// None of them match the cmd data function!
			assert.NotEqual(t, *cmdData.Function, *filterData.Function,
				"Filter %d function mismatch with cmd.Function", i)
		}
	})

	t.Run("Attack scenario - Type confusion", func(t *testing.T) {
		// An attacker could send measurement data but with load control filters
		// This could bypass access controls or cause unexpected behavior
		
		// Legitimate measurement read request
		legitimateCmd := model.CmdType{
			Function: util.Ptr(model.FunctionType("measurementListData")),
			Filter: []model.FilterType{
				{
					CmdControl: &model.CmdControlType{
						Partial: &model.ElementTagType{},
					},
					MeasurementListDataSelectors: &model.MeasurementListDataSelectorsType{
						MeasurementId: util.Ptr(model.MeasurementIdType(1)),
					},
				},
			},
		}
		
		// Malicious request - same function but wrong filter
		maliciousCmd := model.CmdType{
			Function: util.Ptr(model.FunctionType("measurementListData")),
			Filter: []model.FilterType{
				{
					CmdControl: &model.CmdControlType{
						Partial: &model.ElementTagType{},
					},
					// Attacker uses load control filter for measurement function!
					LoadControlLimitListDataSelectors: &model.LoadControlLimitListDataSelectorsType{
						LimitId: util.Ptr(model.LoadControlLimitIdType(999)),
					},
				},
			},
		}
		
		// Both claim to have the same cmd.Function
		assert.Equal(t, *legitimateCmd.Function, *maliciousCmd.Function)
		
		// But different filter functions
		legitFilter, _ := legitimateCmd.Filter[0].Data(legitimateCmd.Function)
		maliciousFilter, _ := maliciousCmd.Filter[0].Data(maliciousCmd.Function)
		assert.NotEqual(t, legitFilter.Function, maliciousFilter.Function,
			"Attack vector: Filter function mismatch not detected!")
	})
}

// TestValidationGap demonstrates that there's NO validation in the current code
func TestValidationGap(t *testing.T) {
	t.Run("No validation exists for function mismatch", func(t *testing.T) {
		// Create a completely invalid combination
		cmd := model.CmdType{
			Function: util.Ptr(model.FunctionType("deviceDiagnosisStateData")),
			
			Filter: []model.FilterType{
				{
					CmdControl: &model.CmdControlType{
						Partial: &model.ElementTagType{},
					},
					// Using a filter for a completely different function
					IdentificationListDataSelectors: &model.IdentificationListDataSelectorsType{
						IdentificationId: util.Ptr(model.IdentificationIdType(1)),
					},
				},
			},
			
			// And data for yet another function
			MeasurementListData: &model.MeasurementListDataType{
				MeasurementData: []model.MeasurementDataType{
					{
						MeasurementId: util.Ptr(model.MeasurementIdType(1)),
					},
				},
			},
		}
		
		// All these operations succeed without any validation!
		cmdData, err := cmd.Data()
		assert.NoError(t, err, "No error despite function mismatch")
		
		filterData, err := cmd.Filter[0].Data(cmd.Function)
		assert.NoError(t, err, "No error despite filter mismatch")
		
		// We have 3 different functions all in one message!
		assert.Equal(t, model.FunctionType("measurementListData"), *cmdData.Function,
			"Data function from actual data field")
		assert.Equal(t, model.FunctionType("identificationListData"), *filterData.Function,
			"Filter function from filter selector")
		// And cmd.Function is something else entirely
		assert.Equal(t, model.FunctionType("deviceDiagnosisStateData"), *cmd.Function,
			"cmd.Function is different from both!")
		
		// This is a massive validation gap!
		t.Logf("WARNING: No validation for function consistency!")
		t.Logf("  cmd.Function: %s", *cmd.Function)
		t.Logf("  Filter function: %s", *filterData.Function)
		t.Logf("  Data function: %s", *cmdData.Function)
	})
}