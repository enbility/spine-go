package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestSmartEnergyManagementPsDataType_UpdateList_BasicReplacement(t *testing.T) {
	// Arrange - existing data structure
	existing := &SmartEnergyManagementPsDataType{
		NodeScheduleInformation: &PowerSequenceNodeScheduleInformationDataType{
			NodeRemoteControllable:           util.Ptr(true),
			SupportsSingleSlotSchedulingOnly: util.Ptr(true),
			AlternativesCount:                util.Ptr(uint(1)),
			SupportsReselection:              util.Ptr(false),
		},
	}

	// New data for full replacement
	newData := &SmartEnergyManagementPsDataType{
		NodeScheduleInformation: &PowerSequenceNodeScheduleInformationDataType{
			NodeRemoteControllable:           util.Ptr(false),
			SupportsSingleSlotSchedulingOnly: util.Ptr(false),
			AlternativesCount:                util.Ptr(uint(2)),
			SupportsReselection:              util.Ptr(true),
		},
	}

	// Act - this should fail initially as UpdateList is not implemented
	result, success := existing.UpdateList(false, true, newData, nil, nil)

	// Assert
	assert.True(t, success)
	assert.NotNil(t, result)

	// Verify the data was replaced
	resultData := result.(*SmartEnergyManagementPsDataType)
	assert.Equal(t, false, *resultData.NodeScheduleInformation.NodeRemoteControllable)
	assert.Equal(t, false, *resultData.NodeScheduleInformation.SupportsSingleSlotSchedulingOnly)
	assert.Equal(t, uint(2), *resultData.NodeScheduleInformation.AlternativesCount)
	assert.Equal(t, true, *resultData.NodeScheduleInformation.SupportsReselection)
}

func TestSmartEnergyManagementPsDataType_UpdateList_InvalidType(t *testing.T) {
	// Arrange
	existing := &SmartEnergyManagementPsDataType{}
	invalidData := "not a SmartEnergyManagementPsDataType"

	// Act
	result, success := existing.UpdateList(false, true, invalidData, nil, nil)

	// Assert
	assert.False(t, success)
	assert.Equal(t, existing, result)
}

// Test OHPCF-005: Update sequence state
func TestSmartEnergyManagementPsDataType_UpdateList_SequenceState(t *testing.T) {
	// Arrange - heat pump with one sequence in "inactive" state
	existing := createOHPCFStructure()

	// Create update to change state to "running"
	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{
					State: util.Ptr(PowerSequenceStateTypeRunning),
				},
			}},
		}},
	}

	// Act - this should update only the state field
	result, success := existing.UpdateList(false, true, update, NewFilterTypePartial(), nil)

	// Assert
	assert.True(t, success)
	assert.NotNil(t, result)

	resultData := result.(*SmartEnergyManagementPsDataType)
	// State should be updated
	assert.Equal(t, PowerSequenceStateTypeRunning, *resultData.Alternatives[0].PowerSequence[0].State.State)
	// Other fields should remain unchanged
	assert.Equal(t, PowerSequenceIdType(1), *resultData.Alternatives[0].PowerSequence[0].Description.SequenceId)
	assert.True(t, *resultData.Alternatives[0].PowerSequence[0].State.SequenceRemoteControllable)
}

// Test OHPCF-004: Update schedule start time
func TestSmartEnergyManagementPsDataType_UpdateList_ScheduleStartTime(t *testing.T) {
	// Arrange - heat pump ready to be scheduled
	existing := createOHPCFStructure()
	startTime := NewAbsoluteOrRelativeTimeType("2024-01-15T10:00:00Z")

	// Create update to set schedule start time
	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				Schedule: &PowerSequenceScheduleDataType{
					StartTime: startTime,
				},
			}},
		}},
	}

	// Act
	result, success := existing.UpdateList(false, true, update, NewFilterTypePartial(), nil)

	// Assert
	assert.True(t, success)
	assert.NotNil(t, result)

	resultData := result.(*SmartEnergyManagementPsDataType)
	// StartTime should be updated
	assert.Equal(t, startTime, resultData.Alternatives[0].PowerSequence[0].Schedule.StartTime)
	// State should remain unchanged
	assert.Equal(t, PowerSequenceStateTypeInactive, *resultData.Alternatives[0].PowerSequence[0].State.State)
}

// Test that updates don't modify the original when persist is false
func TestSmartEnergyManagementPsDataType_UpdateList_NoPersist(t *testing.T) {
	// Arrange
	existing := createOHPCFStructure()
	originalState := *existing.Alternatives[0].PowerSequence[0].State.State

	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{
					State: util.Ptr(PowerSequenceStateTypeRunning),
				},
			}},
		}},
	}

	// Act - persist = false
	result, success := existing.UpdateList(false, false, update, NewFilterTypePartial(), nil)

	// Assert
	assert.True(t, success)
	assert.NotNil(t, result)

	// Original should not be modified
	assert.Equal(t, originalState, *existing.Alternatives[0].PowerSequence[0].State.State)

	// Result should have the update
	resultData := result.(*SmartEnergyManagementPsDataType)
	assert.Equal(t, PowerSequenceStateTypeRunning, *resultData.Alternatives[0].PowerSequence[0].State.State)
}

// Test edge case: update with no alternatives
func TestSmartEnergyManagementPsDataType_UpdateList_NoAlternatives(t *testing.T) {
	// Arrange - structure with no alternatives
	existing := &SmartEnergyManagementPsDataType{
		NodeScheduleInformation: &PowerSequenceNodeScheduleInformationDataType{
			NodeRemoteControllable: util.Ptr(true),
		},
	}

	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{
					State: util.Ptr(PowerSequenceStateTypeRunning),
				},
			}},
		}},
	}

	// Act
	result, success := existing.UpdateList(false, true, update, NewFilterTypePartial(), nil)

	// Assert - should succeed but not crash
	assert.True(t, success)
	assert.NotNil(t, result)

	// Original structure should remain unchanged
	resultData := result.(*SmartEnergyManagementPsDataType)
	assert.Equal(t, 0, len(resultData.Alternatives))
}

// Test multiple field updates in single operation
func TestSmartEnergyManagementPsDataType_UpdateList_MultipleFields(t *testing.T) {
	// Arrange
	existing := createOHPCFStructure()
	startTime := NewAbsoluteOrRelativeTimeType("2024-01-15T10:00:00Z")

	// Update both state and schedule
	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{
					State: util.Ptr(PowerSequenceStateTypeScheduled),
				},
				Schedule: &PowerSequenceScheduleDataType{
					StartTime: startTime,
				},
			}},
		}},
	}

	// Act
	result, success := existing.UpdateList(false, true, update, NewFilterTypePartial(), nil)

	// Assert
	assert.True(t, success)
	assert.NotNil(t, result)

	resultData := result.(*SmartEnergyManagementPsDataType)
	// Both fields should be updated
	assert.Equal(t, PowerSequenceStateTypeScheduled, *resultData.Alternatives[0].PowerSequence[0].State.State)
	assert.Equal(t, startTime, resultData.Alternatives[0].PowerSequence[0].Schedule.StartTime)
	// Other fields remain unchanged
	assert.True(t, *resultData.Alternatives[0].PowerSequence[0].State.SequenceRemoteControllable)
}

// Test that delete filter returns false (not supported)
func TestSmartEnergyManagementPsDataType_UpdateList_DeleteNotSupported(t *testing.T) {
	// Arrange
	existing := createOHPCFStructure()

	// Act - create a delete filter
	deleteFilter := &FilterType{
		CmdControl: &CmdControlType{Delete: &ElementTagType{}},
	}
	result, success := existing.UpdateList(false, true, existing, nil, deleteFilter)

	// Assert
	assert.False(t, success)
	assert.Equal(t, existing, result)
}

// PHASE 2 NOTIFICATION SUPPORT TESTS

// Test OHPCF-005: Device notifies state change from inactive to scheduled
func TestSmartEnergyManagementPsDataType_UpdateList_StateChangeNotification(t *testing.T) {
	// Arrange - heat pump in inactive state
	existing := createOHPCFStructure()

	// Device sends notification that state changed to scheduled
	notification := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{
					State: util.Ptr(PowerSequenceStateTypeScheduled),
				},
			}},
		}},
	}

	// Act - handle partial update notification
	result, success := existing.UpdateList(false, true, notification, NewFilterTypePartial(), nil)

	// Assert
	assert.True(t, success)
	assert.NotNil(t, result)

	resultData := result.(*SmartEnergyManagementPsDataType)
	// State should be updated to scheduled
	assert.Equal(t, PowerSequenceStateTypeScheduled, *resultData.Alternatives[0].PowerSequence[0].State.State)
	// Other state fields should remain unchanged
	assert.True(t, *resultData.Alternatives[0].PowerSequence[0].State.SequenceRemoteControllable)
	assert.Nil(t, resultData.Alternatives[0].PowerSequence[0].State.ActiveSlotNumber)
}

// Test device updates activeSlotNumber when execution starts
func TestSmartEnergyManagementPsDataType_UpdateList_ActiveSlotNumberUpdate(t *testing.T) {
	// Arrange - heat pump in running state
	existing := createOHPCFStructure()
	existing.Alternatives[0].PowerSequence[0].State.State = util.Ptr(PowerSequenceStateTypeRunning)

	// Device notifies that slot 1 is now active
	notification := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{
					ActiveSlotNumber: util.Ptr(PowerTimeSlotNumberType(1)),
				},
			}},
		}},
	}

	// Act
	result, success := existing.UpdateList(false, true, notification, NewFilterTypePartial(), nil)

	// Assert
	assert.True(t, success)
	assert.NotNil(t, result)

	resultData := result.(*SmartEnergyManagementPsDataType)
	// ActiveSlotNumber should be updated
	assert.NotNil(t, resultData.Alternatives[0].PowerSequence[0].State.ActiveSlotNumber)
	assert.Equal(t, PowerTimeSlotNumberType(1), *resultData.Alternatives[0].PowerSequence[0].State.ActiveSlotNumber)
	// State should remain running
	assert.Equal(t, PowerSequenceStateTypeRunning, *resultData.Alternatives[0].PowerSequence[0].State.State)
}

// Test device updates elapsed time during execution
func TestSmartEnergyManagementPsDataType_UpdateList_ElapsedTimeUpdate(t *testing.T) {
	// Arrange - heat pump running slot 1
	existing := createOHPCFStructure()
	existing.Alternatives[0].PowerSequence[0].State.State = util.Ptr(PowerSequenceStateTypeRunning)
	existing.Alternatives[0].PowerSequence[0].State.ActiveSlotNumber = util.Ptr(PowerTimeSlotNumberType(1))

	// Device notifies elapsed time progress (e.g., 5 minutes into execution)
	elapsedTime := DurationType("PT5M") // 5 minutes in ISO 8601 duration format
	notification := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{
					ElapsedSlotTime: &elapsedTime,
				},
			}},
		}},
	}

	// Act
	result, success := existing.UpdateList(false, true, notification, NewFilterTypePartial(), nil)

	// Assert
	assert.True(t, success)
	assert.NotNil(t, result)

	resultData := result.(*SmartEnergyManagementPsDataType)
	// ElapsedSlotTime should be updated
	assert.NotNil(t, resultData.Alternatives[0].PowerSequence[0].State.ElapsedSlotTime)
	assert.Equal(t, DurationType("PT5M"), *resultData.Alternatives[0].PowerSequence[0].State.ElapsedSlotTime)
	// Other fields should remain unchanged
	assert.Equal(t, PowerSequenceStateTypeRunning, *resultData.Alternatives[0].PowerSequence[0].State.State)
	assert.Equal(t, PowerTimeSlotNumberType(1), *resultData.Alternatives[0].PowerSequence[0].State.ActiveSlotNumber)
}

// Test device sends multiple runtime updates in one message
func TestSmartEnergyManagementPsDataType_UpdateList_MultipleRuntimeUpdates(t *testing.T) {
	// Arrange - heat pump scheduled but not yet running
	existing := createOHPCFStructure()
	existing.Alternatives[0].PowerSequence[0].State.State = util.Ptr(PowerSequenceStateTypeScheduled)

	// Device sends comprehensive runtime update when starting execution
	elapsedTime := DurationType("PT0S")    // Just started
	remainingTime := DurationType("PT30M") // 30 minutes remaining
	notification := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{
					State:             util.Ptr(PowerSequenceStateTypeRunning),
					ActiveSlotNumber:  util.Ptr(PowerTimeSlotNumberType(1)),
					ElapsedSlotTime:   &elapsedTime,
					RemainingSlotTime: &remainingTime,
				},
			}},
		}},
	}

	// Act
	result, success := existing.UpdateList(false, true, notification, NewFilterTypePartial(), nil)

	// Assert
	assert.True(t, success)
	assert.NotNil(t, result)

	resultData := result.(*SmartEnergyManagementPsDataType)
	// All runtime fields should be updated
	assert.Equal(t, PowerSequenceStateTypeRunning, *resultData.Alternatives[0].PowerSequence[0].State.State)
	assert.Equal(t, PowerTimeSlotNumberType(1), *resultData.Alternatives[0].PowerSequence[0].State.ActiveSlotNumber)
	assert.Equal(t, DurationType("PT0S"), *resultData.Alternatives[0].PowerSequence[0].State.ElapsedSlotTime)
	assert.Equal(t, DurationType("PT30M"), *resultData.Alternatives[0].PowerSequence[0].State.RemainingSlotTime)
	// SequenceRemoteControllable should remain unchanged
	assert.True(t, *resultData.Alternatives[0].PowerSequence[0].State.SequenceRemoteControllable)
}

// Edge case: Update to non-existent sequence
func TestSmartEnergyManagementPsDataType_UpdateList_NonExistentSequence(t *testing.T) {
	// Arrange - heat pump with one sequence (id=1)
	existing := createOHPCFStructure()

	// Device tries to update sequence with id=2 (doesn't exist)
	notification := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{
				{}, // Empty first sequence
				{ // Second sequence (index 1)
					State: &PowerSequenceStateDataType{
						State: util.Ptr(PowerSequenceStateTypeRunning),
					},
				},
			},
		}},
	}

	// Act
	result, success := existing.UpdateList(false, true, notification, NewFilterTypePartial(), nil)

	// Assert - should succeed but not crash, existing data unchanged
	assert.True(t, success)
	assert.NotNil(t, result)

	resultData := result.(*SmartEnergyManagementPsDataType)
	// Should still have only one sequence
	assert.Equal(t, 1, len(resultData.Alternatives[0].PowerSequence))
	// Original sequence should be unchanged
	assert.Equal(t, PowerSequenceStateTypeInactive, *resultData.Alternatives[0].PowerSequence[0].State.State)
}

// Edge case: Invalid slot number update
func TestSmartEnergyManagementPsDataType_UpdateList_InvalidSlotNumber(t *testing.T) {
	// Arrange - heat pump with only one slot
	existing := createOHPCFStructure()
	existing.Alternatives[0].PowerSequence[0].State.State = util.Ptr(PowerSequenceStateTypeRunning)

	// Device sends invalid slot number (99, but only slot 1 exists)
	notification := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{
					ActiveSlotNumber: util.Ptr(PowerTimeSlotNumberType(99)),
				},
			}},
		}},
	}

	// Act - should accept the update (validation is higher layer responsibility)
	result, success := existing.UpdateList(false, true, notification, NewFilterTypePartial(), nil)

	// Assert
	assert.True(t, success)
	assert.NotNil(t, result)

	resultData := result.(*SmartEnergyManagementPsDataType)
	// ActiveSlotNumber should be updated even if invalid (validation happens elsewhere)
	assert.NotNil(t, resultData.Alternatives[0].PowerSequence[0].State.ActiveSlotNumber)
	assert.Equal(t, PowerTimeSlotNumberType(99), *resultData.Alternatives[0].PowerSequence[0].State.ActiveSlotNumber)
}

// Test state transition: scheduled -> running -> completed
func TestSmartEnergyManagementPsDataType_UpdateList_StateTransitionFlow(t *testing.T) {
	// Arrange - start with scheduled state
	existing := createOHPCFStructure()
	existing.Alternatives[0].PowerSequence[0].State.State = util.Ptr(PowerSequenceStateTypeScheduled)

	// Step 1: Device starts running
	runningNotification := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{
					State:            util.Ptr(PowerSequenceStateTypeRunning),
					ActiveSlotNumber: util.Ptr(PowerTimeSlotNumberType(1)),
				},
			}},
		}},
	}

	result1, success1 := existing.UpdateList(false, true, runningNotification, NewFilterTypePartial(), nil)
	assert.True(t, success1)

	// Step 2: Device completes execution
	// Note: In partial updates, nil means "don't update". To clear a field, we'd need a different mechanism.
	// For now, the device just updates the state, and ActiveSlotNumber remains (though it's irrelevant when completed)
	completedNotification := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{
					State: util.Ptr(PowerSequenceStateTypeCompleted),
				},
			}},
		}},
	}

	result2, success2 := result1.(*SmartEnergyManagementPsDataType).UpdateList(false, true, completedNotification, NewFilterTypePartial(), nil)

	// Assert final state
	assert.True(t, success2)
	assert.NotNil(t, result2)

	finalData := result2.(*SmartEnergyManagementPsDataType)
	assert.Equal(t, PowerSequenceStateTypeCompleted, *finalData.Alternatives[0].PowerSequence[0].State.State)
	// ActiveSlotNumber remains from previous state (this is expected behavior for partial updates)
	assert.Equal(t, PowerTimeSlotNumberType(1), *finalData.Alternatives[0].PowerSequence[0].State.ActiveSlotNumber)
}

// Test notification with only elapsed time (no other changes)
func TestSmartEnergyManagementPsDataType_UpdateList_ElapsedTimeOnly(t *testing.T) {
	// Arrange - heat pump running with active slot
	existing := createOHPCFStructure()
	existing.Alternatives[0].PowerSequence[0].State.State = util.Ptr(PowerSequenceStateTypeRunning)
	existing.Alternatives[0].PowerSequence[0].State.ActiveSlotNumber = util.Ptr(PowerTimeSlotNumberType(1))
	initialElapsedTime := DurationType("PT5M") // 5 minutes
	existing.Alternatives[0].PowerSequence[0].State.ElapsedSlotTime = &initialElapsedTime

	// Device sends progress update (now 10 minutes elapsed)
	updatedElapsedTime := DurationType("PT10M") // 10 minutes
	notification := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{
					ElapsedSlotTime: &updatedElapsedTime,
				},
			}},
		}},
	}

	// Act
	result, success := existing.UpdateList(false, true, notification, NewFilterTypePartial(), nil)

	// Assert
	assert.True(t, success)
	assert.NotNil(t, result)

	resultData := result.(*SmartEnergyManagementPsDataType)
	// Only elapsed time should be updated
	assert.Equal(t, DurationType("PT10M"), *resultData.Alternatives[0].PowerSequence[0].State.ElapsedSlotTime)
	// All other fields should remain unchanged
	assert.Equal(t, PowerSequenceStateTypeRunning, *resultData.Alternatives[0].PowerSequence[0].State.State)
	assert.Equal(t, PowerTimeSlotNumberType(1), *resultData.Alternatives[0].PowerSequence[0].State.ActiveSlotNumber)
	assert.True(t, *resultData.Alternatives[0].PowerSequence[0].State.SequenceRemoteControllable)
}

// Test complete notification scenario from device perspective
func TestSmartEnergyManagementPsDataType_UpdateList_CompleteNotificationScenario(t *testing.T) {
	// Arrange - Energy manager has sent a schedule to heat pump
	heatPump := createOHPCFStructure()
	startTime := NewAbsoluteOrRelativeTimeType("2024-01-15T14:00:00Z")
	heatPump.Alternatives[0].PowerSequence[0].Schedule.StartTime = startTime
	heatPump.Alternatives[0].PowerSequence[0].State.State = util.Ptr(PowerSequenceStateTypeScheduled)

	// Scenario: Heat pump executes the scheduled sequence

	// 1. Heat pump starts execution at scheduled time
	startNotification := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{
					State:             util.Ptr(PowerSequenceStateTypeRunning),
					ActiveSlotNumber:  util.Ptr(PowerTimeSlotNumberType(1)),
					ElapsedSlotTime:   util.Ptr(DurationType("PT0S")),
					RemainingSlotTime: util.Ptr(DurationType("PT30M")),
				},
			}},
		}},
	}

	result1, success1 := heatPump.UpdateList(false, true, startNotification, NewFilterTypePartial(), nil)
	assert.True(t, success1)
	heatPump = result1.(*SmartEnergyManagementPsDataType)

	// Verify state after start
	assert.Equal(t, PowerSequenceStateTypeRunning, *heatPump.Alternatives[0].PowerSequence[0].State.State)
	assert.Equal(t, PowerTimeSlotNumberType(1), *heatPump.Alternatives[0].PowerSequence[0].State.ActiveSlotNumber)
	assert.Equal(t, DurationType("PT0S"), *heatPump.Alternatives[0].PowerSequence[0].State.ElapsedSlotTime)

	// 2. Heat pump sends progress update after 15 minutes
	progressNotification := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{
					ElapsedSlotTime:   util.Ptr(DurationType("PT15M")),
					RemainingSlotTime: util.Ptr(DurationType("PT15M")),
				},
			}},
		}},
	}

	result2, success2 := heatPump.UpdateList(false, true, progressNotification, NewFilterTypePartial(), nil)
	assert.True(t, success2)
	heatPump = result2.(*SmartEnergyManagementPsDataType)

	// Verify progress update
	assert.Equal(t, PowerSequenceStateTypeRunning, *heatPump.Alternatives[0].PowerSequence[0].State.State) // Unchanged
	assert.Equal(t, DurationType("PT15M"), *heatPump.Alternatives[0].PowerSequence[0].State.ElapsedSlotTime)
	assert.Equal(t, DurationType("PT15M"), *heatPump.Alternatives[0].PowerSequence[0].State.RemainingSlotTime)

	// 3. Heat pump completes execution
	completeNotification := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{
					State:             util.Ptr(PowerSequenceStateTypeCompleted),
					ElapsedSlotTime:   util.Ptr(DurationType("PT30M")),
					RemainingSlotTime: util.Ptr(DurationType("PT0S")),
				},
			}},
		}},
	}

	result3, success3 := heatPump.UpdateList(false, true, completeNotification, NewFilterTypePartial(), nil)
	assert.True(t, success3)
	heatPump = result3.(*SmartEnergyManagementPsDataType)

	// Verify final state
	assert.Equal(t, PowerSequenceStateTypeCompleted, *heatPump.Alternatives[0].PowerSequence[0].State.State)
	assert.Equal(t, DurationType("PT30M"), *heatPump.Alternatives[0].PowerSequence[0].State.ElapsedSlotTime)
	assert.Equal(t, DurationType("PT0S"), *heatPump.Alternatives[0].PowerSequence[0].State.RemainingSlotTime)
	// ActiveSlotNumber remains set (partial updates don't clear fields)
	assert.Equal(t, PowerTimeSlotNumberType(1), *heatPump.Alternatives[0].PowerSequence[0].State.ActiveSlotNumber)

	// Verify other fields remained unchanged throughout
	assert.Equal(t, startTime, heatPump.Alternatives[0].PowerSequence[0].Schedule.StartTime)
	assert.True(t, *heatPump.Alternatives[0].PowerSequence[0].State.SequenceRemoteControllable)
}

// Helper function to create a typical OHPCF structure
func createOHPCFStructure() *SmartEnergyManagementPsDataType {
	return &SmartEnergyManagementPsDataType{
		NodeScheduleInformation: &PowerSequenceNodeScheduleInformationDataType{
			NodeRemoteControllable:           util.Ptr(true),
			SupportsSingleSlotSchedulingOnly: util.Ptr(true),
			AlternativesCount:                util.Ptr(uint(1)),
			TotalSequencesCountMax:           util.Ptr(uint(1)),
			SupportsReselection:              util.Ptr(false),
		},
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(1)),
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				Description: &PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(PowerSequenceIdType(1)),
					PowerUnit:  util.Ptr(UnitOfMeasurementTypeW),
				},
				State: &PowerSequenceStateDataType{
					State:                      util.Ptr(PowerSequenceStateTypeInactive),
					SequenceRemoteControllable: util.Ptr(true),
				},
				Schedule: &PowerSequenceScheduleDataType{},
				PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
					Schedule: &PowerTimeSlotScheduleDataType{
						SlotNumber: util.Ptr(PowerTimeSlotNumberType(1)),
					},
					ValueList: &SmartEnergyManagementPsPowerTimeSlotValueListType{
						Value: []PowerTimeSlotValueDataType{{
							ValueType: util.Ptr(PowerTimeSlotValueTypeTypePower),
							Value: &ScaledNumberType{
								Number: util.Ptr(NumberType(3000)),
								Scale:  util.Ptr(ScaleType(0)),
							},
						}},
					},
				}},
			}},
		}},
	}
}
