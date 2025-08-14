package model

import (
	"errors"
	"fmt"
)

// ValidateFunctionConsistencyStrict validates that all function references in a CmdType are consistent
// This prevents type confusion attacks where cmd.Function doesn't match filter functions
// Enforces SPINE spec requirement that cmd.Function must be present and match when filters are used
func (cmd *CmdType) ValidateFunctionConsistencyStrict() error {
	if cmd == nil {
		return errors.New("cmd is nil")
	}

	// Extract the actual function from the data
	cmdData, err := cmd.Data()
	if err != nil {
		return fmt.Errorf("failed to extract cmd data: %w", err)
	}

	if cmdData.Function == nil {
		return errors.New("cmd data has no function")
	}

	baseFunction := *cmdData.Function

	// In strict mode, cmd.Function must be present and match
	if cmd.Function == nil || *cmd.Function == "" {
		return fmt.Errorf("cmd.Function is missing or empty, expected %s", baseFunction)
	}

	if *cmd.Function != baseFunction {
		return fmt.Errorf("cmd.Function (%s) doesn't match data function (%s)",
			*cmd.Function, baseFunction)
	}

	// Check all filters - in strict mode, all must be valid
	if len(cmd.Filter) > 0 {
		for i, filter := range cmd.Filter {
			// Pass cmd.Function for partial filters without selectors
			filterData, err := filter.Data(cmd.Function)
			if err != nil {
				return fmt.Errorf("filter[%d] has invalid data: %w", i, err)
			}
			if filterData.Function == nil {
				return fmt.Errorf("filter[%d] has no function", i)
			}
			if *filterData.Function != baseFunction {
				return fmt.Errorf("filter[%d] function (%s) doesn't match data function (%s)",
					i, *filterData.Function, baseFunction)
			}
		}
	}

	return nil
}

// HasFunctionMismatch returns true if there's any function inconsistency
// This is useful for logging/metrics without failing the operation
func (cmd *CmdType) HasFunctionMismatch() bool {
	if cmd == nil {
		return false
	}

	cmdData, err := cmd.Data()
	if err != nil || cmdData.Function == nil {
		return false
	}

	baseFunction := *cmdData.Function

	// Check cmd.Function
	if cmd.Function != nil && *cmd.Function != "" && *cmd.Function != baseFunction {
		return true
	}

	// Check filters with cmd.Function as fallback
	for _, filter := range cmd.Filter {
		filterData, err := filter.Data(cmd.Function)
		if err != nil || filterData.Function == nil {
			continue
		}
		if *filterData.Function != baseFunction {
			return true
		}
	}

	return false
}

// GetInconsistentFunctions returns a list of all inconsistent function references
// Useful for detailed error reporting and debugging
func (cmd *CmdType) GetInconsistentFunctions() []string {
	if cmd == nil {
		return nil
	}

	var inconsistencies []string

	cmdData, err := cmd.Data()
	if err != nil || cmdData.Function == nil {
		return inconsistencies
	}

	baseFunction := *cmdData.Function

	// Check cmd.Function
	if cmd.Function != nil && *cmd.Function != "" && *cmd.Function != baseFunction {
		inconsistencies = append(inconsistencies,
			fmt.Sprintf("cmd.Function=%s (expected %s)", *cmd.Function, baseFunction))
	}

	// Check filters with cmd.Function as fallback
	for i, filter := range cmd.Filter {
		filterData, err := filter.Data(cmd.Function)
		if err != nil || filterData.Function == nil {
			continue
		}
		if *filterData.Function != baseFunction {
			inconsistencies = append(inconsistencies,
				fmt.Sprintf("filter[%d].Function=%s (expected %s)",
					i, *filterData.Function, baseFunction))
		}
	}

	return inconsistencies
}
