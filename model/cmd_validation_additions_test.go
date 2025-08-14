package model_test

import (
	"testing"

	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestValidateFunctionConsistencyStrict(t *testing.T) {
	tests := []struct {
		name        string
		cmd         *model.CmdType
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid - all functions match",
			cmd: &model.CmdType{
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
				MeasurementListData: &model.MeasurementListDataType{
					MeasurementData: []model.MeasurementDataType{
						{
							MeasurementId: util.Ptr(model.MeasurementIdType(1)),
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Invalid - empty cmd.Function",
			cmd: &model.CmdType{
				Function: util.Ptr(model.FunctionType("")),
				MeasurementListData: &model.MeasurementListDataType{
					MeasurementData: []model.MeasurementDataType{
						{
							MeasurementId: util.Ptr(model.MeasurementIdType(1)),
						},
					},
				},
			},
			expectError: true,
			errorMsg:    "cmd.Function is missing or empty",
		},
		{
			name: "Invalid - nil cmd.Function",
			cmd: &model.CmdType{
				Function: nil,
				MeasurementListData: &model.MeasurementListDataType{
					MeasurementData: []model.MeasurementDataType{
						{
							MeasurementId: util.Ptr(model.MeasurementIdType(1)),
						},
					},
				},
			},
			expectError: true,
			errorMsg:    "cmd.Function is missing or empty",
		},
		{
			name: "Invalid - cmd.Function doesn't match data",
			cmd: &model.CmdType{
				Function: util.Ptr(model.FunctionType("loadControlLimitListData")),
				MeasurementListData: &model.MeasurementListDataType{
					MeasurementData: []model.MeasurementDataType{
						{
							MeasurementId: util.Ptr(model.MeasurementIdType(1)),
						},
					},
				},
			},
			expectError: true,
			errorMsg:    "cmd.Function (loadControlLimitListData) doesn't match data function (measurementListData)",
		},
		{
			name: "Invalid - filter function doesn't match data",
			cmd: &model.CmdType{
				Function: util.Ptr(model.FunctionType("measurementListData")),
				Filter: []model.FilterType{
					{
						CmdControl: &model.CmdControlType{
							Partial: &model.ElementTagType{},
						},
						LoadControlLimitListDataSelectors: &model.LoadControlLimitListDataSelectorsType{
							LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
						},
					},
				},
				MeasurementListData: &model.MeasurementListDataType{
					MeasurementData: []model.MeasurementDataType{
						{
							MeasurementId: util.Ptr(model.MeasurementIdType(1)),
						},
					},
				},
			},
			expectError: true,
			errorMsg:    "filter[0] function (loadControlLimitListData) doesn't match data function (measurementListData)",
		},
		{
			name: "Valid - partial filter without selectors (means all fields)",
			cmd: &model.CmdType{
				Function: util.Ptr(model.FunctionType("measurementListData")),
				Filter: []model.FilterType{
					{
						CmdControl: &model.CmdControlType{
							Partial: &model.ElementTagType{},
						},
						// No selector or element fields with function tags - valid SPINE, means "all fields"
					},
				},
				MeasurementListData: &model.MeasurementListDataType{
					MeasurementData: []model.MeasurementDataType{
						{
							MeasurementId: util.Ptr(model.MeasurementIdType(1)),
						},
					},
				},
			},
			expectError: false, // This is actually valid SPINE - partial filter without selectors
		},
		{
			name: "Invalid - multiple filter mismatches",
			cmd: &model.CmdType{
				Filter: []model.FilterType{
					{
						CmdControl: &model.CmdControlType{
							Delete: &model.ElementTagType{},
						},
						BillListDataSelectors: &model.BillListDataSelectorsType{
							BillId: util.Ptr(model.BillIdType(1)),
						},
					},
					{
						CmdControl: &model.CmdControlType{
							Partial: &model.ElementTagType{},
						},
						LoadControlLimitListDataSelectors: &model.LoadControlLimitListDataSelectorsType{
							LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
						},
					},
				},
				MeasurementListData: &model.MeasurementListDataType{
					MeasurementData: []model.MeasurementDataType{
						{
							MeasurementId: util.Ptr(model.MeasurementIdType(1)),
						},
					},
				},
			},
			expectError: true,
			errorMsg:    "cmd.Function is missing or empty",
		},
		{
			name:        "Nil cmd",
			cmd:         nil,
			expectError: true,
			errorMsg:    "cmd is nil",
		},
		{
			name: "No data in cmd",
			cmd: &model.CmdType{
				Function: util.Ptr(model.FunctionType("measurementListData")),
			},
			expectError: true,
			errorMsg:    "failed to extract cmd data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.ValidateFunctionConsistencyStrict()
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" && err != nil {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetInconsistentFunctions(t *testing.T) {
	tests := []struct {
		name              string
		cmd               *model.CmdType
		expectedInconsist []string
	}{
		{
			name: "No inconsistencies",
			cmd: &model.CmdType{
				Function: util.Ptr(model.FunctionType("measurementListData")),
				MeasurementListData: &model.MeasurementListDataType{
					MeasurementData: []model.MeasurementDataType{
						{
							MeasurementId: util.Ptr(model.MeasurementIdType(1)),
						},
					},
				},
			},
			expectedInconsist: []string{},
		},
		{
			name: "cmd.Function mismatch",
			cmd: &model.CmdType{
				Function: util.Ptr(model.FunctionType("loadControlLimitListData")),
				MeasurementListData: &model.MeasurementListDataType{
					MeasurementData: []model.MeasurementDataType{
						{
							MeasurementId: util.Ptr(model.MeasurementIdType(1)),
						},
					},
				},
			},
			expectedInconsist: []string{
				"cmd.Function=loadControlLimitListData (expected measurementListData)",
			},
		},
		{
			name: "Filter function mismatch",
			cmd: &model.CmdType{
				Function: util.Ptr(model.FunctionType("measurementListData")),
				Filter: []model.FilterType{
					{
						CmdControl: &model.CmdControlType{
							Partial: &model.ElementTagType{},
						},
						LoadControlLimitListDataSelectors: &model.LoadControlLimitListDataSelectorsType{
							LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
						},
					},
				},
				MeasurementListData: &model.MeasurementListDataType{
					MeasurementData: []model.MeasurementDataType{
						{
							MeasurementId: util.Ptr(model.MeasurementIdType(1)),
						},
					},
				},
			},
			expectedInconsist: []string{
				"filter[0].Function=loadControlLimitListData (expected measurementListData)",
			},
		},
		{
			name:              "Nil cmd",
			cmd:               nil,
			expectedInconsist: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cmd.GetInconsistentFunctions()
			if tt.expectedInconsist == nil {
				assert.Nil(t, result)
			} else {
				assert.ElementsMatch(t, tt.expectedInconsist, result)
			}
		})
	}
}

func TestHasFunctionMismatch(t *testing.T) {
	tests := []struct {
		name        string
		cmd         *model.CmdType
		hasMismatch bool
	}{
		{
			name: "No mismatch",
			cmd: &model.CmdType{
				Function: util.Ptr(model.FunctionType("measurementListData")),
				MeasurementListData: &model.MeasurementListDataType{
					MeasurementData: []model.MeasurementDataType{
						{
							MeasurementId: util.Ptr(model.MeasurementIdType(1)),
						},
					},
				},
			},
			hasMismatch: false,
		},
		{
			name: "Has mismatch",
			cmd: &model.CmdType{
				Function: util.Ptr(model.FunctionType("loadControlLimitListData")),
				MeasurementListData: &model.MeasurementListDataType{
					MeasurementData: []model.MeasurementDataType{
						{
							MeasurementId: util.Ptr(model.MeasurementIdType(1)),
						},
					},
				},
			},
			hasMismatch: true,
		},
		{
			name:        "Nil cmd",
			cmd:         nil,
			hasMismatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cmd.HasFunctionMismatch()
			assert.Equal(t, tt.hasMismatch, result)
		})
	}
}