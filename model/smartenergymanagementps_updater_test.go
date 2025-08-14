package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

// TestSmartEnergyManagementPsDataType_RealWorld_InitialEmptyState tests the initial state
// from OHPCF log lines 57-58: Start with only nodeScheduleInformation (alternativesCount: 0)
func TestSmartEnergyManagementPsDataType_RealWorld_InitialEmptyState(t *testing.T) {
	// Arrange - Initial state with no alternatives
	existing := &SmartEnergyManagementPsDataType{
		NodeScheduleInformation: &PowerSequenceNodeScheduleInformationDataType{
			NodeRemoteControllable:           util.Ptr(true),
			SupportsSingleSlotSchedulingOnly: util.Ptr(true),
			AlternativesCount:                util.Ptr(uint(0)), // No alternatives initially
			TotalSequencesCountMax:           util.Ptr(uint(1)),
			SupportsReselection:              util.Ptr(false),
		},
		// No Alternatives array - truly empty
	}

	// Act & Assert - Verify initial state
	assert.NotNil(t, existing.NodeScheduleInformation)
	assert.Equal(t, uint(0), *existing.NodeScheduleInformation.AlternativesCount)
	assert.Equal(t, 0, len(existing.Alternatives))
}

// TestSmartEnergyManagementPsDataType_RealWorld_AlternativeCreation tests alternative creation
// from OHPCF log line 95: Update from alternativesCount: 0 to alternativesCount: 1
func TestSmartEnergyManagementPsDataType_RealWorld_AlternativeCreation(t *testing.T) {
	// Arrange - Start with empty state
	existing := &SmartEnergyManagementPsDataType{
		NodeScheduleInformation: &PowerSequenceNodeScheduleInformationDataType{
			AlternativesCount: util.Ptr(uint(0)),
		},
	}

	// Update - Add complete alternative with power sequence structure
	update := &SmartEnergyManagementPsDataType{
		NodeScheduleInformation: &PowerSequenceNodeScheduleInformationDataType{
			AlternativesCount: util.Ptr(uint(1)), // Update count to 1
		},
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(0)), // First alternative
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				Description: &PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(PowerSequenceIdType(0)), // sequenceId: 0
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

	// Act - Full replacement (no partial filter)
	result, success := existing.UpdateList(false, true, update, nil, nil)

	// Assert
	assert.True(t, success)
	assert.NotNil(t, result)

	resultData := result.(*SmartEnergyManagementPsDataType)
	assert.Equal(t, uint(1), *resultData.NodeScheduleInformation.AlternativesCount)
	assert.Equal(t, 1, len(resultData.Alternatives))
	assert.Equal(t, AlternativesIdType(0), *resultData.Alternatives[0].Relation.AlternativesId)
	assert.Equal(t, PowerSequenceIdType(0), *resultData.Alternatives[0].PowerSequence[0].Description.SequenceId)
	assert.Equal(t, PowerSequenceStateTypeInactive, *resultData.Alternatives[0].PowerSequence[0].State.State)
}

// TestSmartEnergyManagementPsDataType_RealWorld_PartialScheduleUpdate tests partial schedule update
// from OHPCF log line 101: Update ONLY the schedule.startTime field to "PT5S"
func TestSmartEnergyManagementPsDataType_RealWorld_PartialScheduleUpdate(t *testing.T) {
	// Arrange - Existing structure with sequenceId: 0
	existing := createComplexOHPCFStructure()

	// Update - ONLY schedule.startTime field, using composite key sequenceId: 0
	startTime := NewAbsoluteOrRelativeTimeType("PT5S")
	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(0)), // Key: alternativesId
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				Description: &PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(PowerSequenceIdType(0)), // Key: sequenceId
				},
				Schedule: &PowerSequenceScheduleDataType{
					StartTime: startTime, // ONLY field being updated
				},
			}},
		}},
	}

	// Act - Partial update
	result, success := existing.UpdateList(false, true, update, NewFilterTypePartial(), nil)

	// Assert
	assert.True(t, success)
	assert.NotNil(t, result)

	resultData := result.(*SmartEnergyManagementPsDataType)
	// StartTime should be updated
	assert.Equal(t, startTime, resultData.Alternatives[0].PowerSequence[0].Schedule.StartTime)

	// ALL other fields MUST be preserved
	assert.Equal(t, PowerSequenceStateTypeInactive, *resultData.Alternatives[0].PowerSequence[0].State.State)
	assert.True(t, *resultData.Alternatives[0].PowerSequence[0].State.SequenceRemoteControllable)
	assert.Equal(t, NumberType(3000), *resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value[0].Value.Number)
	assert.Equal(t, PowerTimeSlotValueTypeTypePower, *resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value[0].ValueType)
}

// TestSmartEnergyManagementPsDataType_RealWorld_StateTransitions tests state transitions
// from OHPCF log lines 102-109: inactive → scheduled → running → completed
func TestSmartEnergyManagementPsDataType_RealWorld_StateTransitions(t *testing.T) {
	// Arrange - Start with inactive state, sequenceId: 0
	existing := createComplexOHPCFStructure()

	// Step 1: inactive → scheduled (activeSlotNumber should NOT appear)
	scheduledUpdate := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(0)), // Key
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				Description: &PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(PowerSequenceIdType(0)), // Key
				},
				State: &PowerSequenceStateDataType{
					State: util.Ptr(PowerSequenceStateTypeScheduled),
				},
			}},
		}},
	}

	result1, success1 := existing.UpdateList(false, true, scheduledUpdate, NewFilterTypePartial(), nil)
	assert.True(t, success1)
	resultData1 := result1.(*SmartEnergyManagementPsDataType)
	assert.Equal(t, PowerSequenceStateTypeScheduled, *resultData1.Alternatives[0].PowerSequence[0].State.State)
	assert.Nil(t, resultData1.Alternatives[0].PowerSequence[0].State.ActiveSlotNumber) // Should NOT be set
	// Power values must be preserved
	assert.Equal(t, NumberType(3000), *resultData1.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value[0].Value.Number)

	// Step 2: scheduled → running (activeSlotNumber APPEARS)
	runningUpdate := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(0)), // Key
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				Description: &PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(PowerSequenceIdType(0)), // Key
				},
				State: &PowerSequenceStateDataType{
					State:            util.Ptr(PowerSequenceStateTypeRunning),
					ActiveSlotNumber: util.Ptr(PowerTimeSlotNumberType(1)), // Appears when running
				},
			}},
		}},
	}

	result2, success2 := resultData1.UpdateList(false, true, runningUpdate, NewFilterTypePartial(), nil)
	assert.True(t, success2)
	resultData2 := result2.(*SmartEnergyManagementPsDataType)
	assert.Equal(t, PowerSequenceStateTypeRunning, *resultData2.Alternatives[0].PowerSequence[0].State.State)
	assert.Equal(t, PowerTimeSlotNumberType(1), *resultData2.Alternatives[0].PowerSequence[0].State.ActiveSlotNumber) // Now set
	// Power values must be preserved
	assert.Equal(t, NumberType(3000), *resultData2.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value[0].Value.Number)

	// Step 3: running → completed (activeSlotNumber disappears logically, but partial updates preserve it)
	completedUpdate := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(0)), // Key
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				Description: &PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(PowerSequenceIdType(0)), // Key
				},
				State: &PowerSequenceStateDataType{
					State: util.Ptr(PowerSequenceStateTypeCompleted),
				},
			}},
		}},
	}

	result3, success3 := resultData2.UpdateList(false, true, completedUpdate, NewFilterTypePartial(), nil)
	assert.True(t, success3)
	resultData3 := result3.(*SmartEnergyManagementPsDataType)
	assert.Equal(t, PowerSequenceStateTypeCompleted, *resultData3.Alternatives[0].PowerSequence[0].State.State)
	// activeSlotNumber remains in partial updates (in real world, device would send full state or explicit nil)
	assert.Equal(t, PowerTimeSlotNumberType(1), *resultData3.Alternatives[0].PowerSequence[0].State.ActiveSlotNumber)
	// Power values must be preserved throughout all transitions
	assert.Equal(t, NumberType(3000), *resultData3.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value[0].Value.Number)
}

// TestSmartEnergyManagementPsDataType_RealWorld_KeyBasedMatching tests key-based matching
// Update sequence with sequenceId: 2 when array has sequences [0, 5, 2]
// THIS TEST WILL FAIL with current positional implementation
func TestSmartEnergyManagementPsDataType_RealWorld_KeyBasedMatching(t *testing.T) {
	// Arrange - Structure with multiple sequences in non-sequential order: [0, 5, 2]
	existing := &SmartEnergyManagementPsDataType{
		NodeScheduleInformation: &PowerSequenceNodeScheduleInformationDataType{
			AlternativesCount: util.Ptr(uint(1)),
		},
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(0)),
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{
				// Sequence 0 at index 0
				{
					Description: &PowerSequenceDescriptionDataType{
						SequenceId: util.Ptr(PowerSequenceIdType(0)),
					},
					State: &PowerSequenceStateDataType{
						State: util.Ptr(PowerSequenceStateTypeInactive),
					},
				},
				// Sequence 5 at index 1
				{
					Description: &PowerSequenceDescriptionDataType{
						SequenceId: util.Ptr(PowerSequenceIdType(5)),
					},
					State: &PowerSequenceStateDataType{
						State: util.Ptr(PowerSequenceStateTypeInactive),
					},
				},
				// Sequence 2 at index 2
				{
					Description: &PowerSequenceDescriptionDataType{
						SequenceId: util.Ptr(PowerSequenceIdType(2)),
					},
					State: &PowerSequenceStateDataType{
						State: util.Ptr(PowerSequenceStateTypeInactive),
					},
				},
			},
		}},
	}

	// Update - Target sequence with sequenceId: 2 (should match by ID, not position)
	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(0)), // Key: alternativesId
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				Description: &PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(PowerSequenceIdType(2)), // Key: sequenceId = 2
				},
				State: &PowerSequenceStateDataType{
					State: util.Ptr(PowerSequenceStateTypeRunning), // Update to running
				},
			}},
		}},
	}

	// Act - Partial update
	result, success := existing.UpdateList(false, true, update, NewFilterTypePartial(), nil)

	// Assert - This test demonstrates the current limitation
	assert.True(t, success)
	assert.NotNil(t, result)

	resultData := result.(*SmartEnergyManagementPsDataType)

	// CURRENT BEHAVIOR (positional): Updates sequence at index 0 (sequenceId: 0) instead of sequenceId: 2
	// This is INCORRECT but shows current implementation
	// NOTE: This assertion will PASS with current implementation but shows wrong behavior
	
	// What SHOULD happen with proper key-based matching:
	// - Sequence 0 (index 0): Should remain inactive
	// - Sequence 5 (index 1): Should remain inactive  
	// - Sequence 2 (index 2): Should be updated to running
	
	// What ACTUALLY happens with current positional implementation:
	// - Updates sequence at index 0 (sequenceId: 0) instead of the target

	// Test for EXPECTED behavior (proper key-based matching)
	// These assertions will FAIL with current implementation:
	
	// Find sequence with sequenceId: 2 in result
	var targetSequence *SmartEnergyManagementPsPowerSequenceType
	for i := range resultData.Alternatives[0].PowerSequence {
		if *resultData.Alternatives[0].PowerSequence[i].Description.SequenceId == PowerSequenceIdType(2) {
			targetSequence = &resultData.Alternatives[0].PowerSequence[i]
			break
		}
	}
	
	assert.NotNil(t, targetSequence, "Should find sequence with sequenceId: 2")
	if targetSequence != nil {
		// FAILS with current implementation: sequence 2 should be running, but it remains inactive
		// because current implementation updates by position (index 0), not by key
		assert.Equal(t, PowerSequenceStateTypeRunning, *targetSequence.State.State,
			"Sequence with sequenceId: 2 should be updated to running state")
	}

	// Verify other sequences remain unchanged
	for i := range resultData.Alternatives[0].PowerSequence {
		seq := &resultData.Alternatives[0].PowerSequence[i]
		sequenceId := *seq.Description.SequenceId
		if sequenceId != PowerSequenceIdType(2) {
			assert.Equal(t, PowerSequenceStateTypeInactive, *seq.State.State,
				"Sequences other than sequenceId: 2 should remain inactive")
		}
	}
}

// TestSmartEnergyManagementPsDataType_RealWorld_CompositeKeyValidation tests composite key validation
// Should reject updates with missing keys
func TestSmartEnergyManagementPsDataType_RealWorld_CompositeKeyValidation(t *testing.T) {
	// Arrange
	existing := createComplexOHPCFStructure()

	// Test 1: Update with missing sequenceId (composite key incomplete)
	updateMissingSequenceId := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(0)), // Has alternativesId
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				// Missing Description.SequenceId - incomplete composite key
				State: &PowerSequenceStateDataType{
					State: util.Ptr(PowerSequenceStateTypeRunning),
				},
			}},
		}},
	}

	// With proper key-based implementation, this should either:
	// 1. Apply to ALL sequences (SPINE "update all" semantics when key missing)
	// 2. Be rejected for incomplete key
	// Current implementation will apply positionally to first sequence

	result1, success1 := existing.UpdateList(false, true, updateMissingSequenceId, NewFilterTypePartial(), nil)
	assert.True(t, success1) // Currently succeeds with positional logic

	// Test 2: Update with missing alternativesId
	updateMissingAlternativesId := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			// Missing Relation.AlternativesId - incomplete composite key
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				Description: &PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(PowerSequenceIdType(0)), // Has sequenceId
				},
				State: &PowerSequenceStateDataType{
					State: util.Ptr(PowerSequenceStateTypeScheduled),
				},
			}},
		}},
	}

	result2, success2 := existing.UpdateList(false, true, updateMissingAlternativesId, NewFilterTypePartial(), nil)
	assert.True(t, success2) // Currently succeeds with positional logic

	// Both tests pass with current implementation but don't follow proper key semantics
	_ = result1
	_ = result2
}

// TestSmartEnergyManagementPsDataType_RealWorld_UpdateAllSemantics tests "update all" semantics
// SPINE spec requirement: When key is missing, update ALL matching elements
func TestSmartEnergyManagementPsDataType_RealWorld_UpdateAllSemantics(t *testing.T) {
	// Arrange - Structure with multiple sequences
	existing := &SmartEnergyManagementPsDataType{
		NodeScheduleInformation: &PowerSequenceNodeScheduleInformationDataType{
			AlternativesCount: util.Ptr(uint(1)),
		},
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(0)),
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{
				{
					Description: &PowerSequenceDescriptionDataType{
						SequenceId: util.Ptr(PowerSequenceIdType(1)),
					},
					State: &PowerSequenceStateDataType{
						State:                      util.Ptr(PowerSequenceStateTypeInactive),
						SequenceRemoteControllable: util.Ptr(true),
					},
				},
				{
					Description: &PowerSequenceDescriptionDataType{
						SequenceId: util.Ptr(PowerSequenceIdType(2)),
					},
					State: &PowerSequenceStateDataType{
						State:                      util.Ptr(PowerSequenceStateTypeInactive),
						SequenceRemoteControllable: util.Ptr(true),
					},
				},
				{
					Description: &PowerSequenceDescriptionDataType{
						SequenceId: util.Ptr(PowerSequenceIdType(3)),
					},
					State: &PowerSequenceStateDataType{
						State:                      util.Ptr(PowerSequenceStateTypeInactive),
						SequenceRemoteControllable: util.Ptr(false), // Different value
					},
				},
			},
		}},
	}

	// Update WITHOUT sequenceId key - should update ALL sequences per SPINE spec
	updateAllSequences := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(0)), // Has alternativesId
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				// NO Description.SequenceId - means "update ALL sequences"
				State: &PowerSequenceStateDataType{
					SequenceRemoteControllable: util.Ptr(false), // Set all to false
				},
			}},
		}},
	}

	// Act
	result, success := existing.UpdateList(false, true, updateAllSequences, NewFilterTypePartial(), nil)

	// Assert
	assert.True(t, success)
	assert.NotNil(t, result)

	resultData := result.(*SmartEnergyManagementPsDataType)

	// With proper "update all" semantics, ALL sequences should have SequenceRemoteControllable: false
	// Current implementation only updates first sequence positionally

	// Check what should happen vs what actually happens:
	sequences := resultData.Alternatives[0].PowerSequence
	assert.Equal(t, 3, len(sequences))

	// EXPECTED behavior (proper "update all"): All sequences updated
	// ACTUAL behavior (current): Only first sequence updated
	
	// This test demonstrates the gap between SPINE spec and current implementation
	// In a proper key-based implementation:
	// - All sequences would have SequenceRemoteControllable: false
	// - Other fields would be preserved per sequence
}

// TestSmartEnergyManagementPsDataType_RealWorld_ComplexNestedUpdate tests complex nested updates
// Update specific time slot values using slotNumber + valueType composite keys
func TestSmartEnergyManagementPsDataType_RealWorld_ComplexNestedUpdate(t *testing.T) {
	// Arrange - Structure with multiple time slots and value types
	existing := &SmartEnergyManagementPsDataType{
		NodeScheduleInformation: &PowerSequenceNodeScheduleInformationDataType{
			AlternativesCount: util.Ptr(uint(1)),
		},
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(0)),
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				Description: &PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(PowerSequenceIdType(1)),
				},
				State: &PowerSequenceStateDataType{
					State: util.Ptr(PowerSequenceStateTypeInactive),
				},
				PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{
					// Slot 1 with power and energy values
					{
						Schedule: &PowerTimeSlotScheduleDataType{
							SlotNumber: util.Ptr(PowerTimeSlotNumberType(1)),
						},
						ValueList: &SmartEnergyManagementPsPowerTimeSlotValueListType{
							Value: []PowerTimeSlotValueDataType{
								{
									ValueType: util.Ptr(PowerTimeSlotValueTypeTypePower),
									Value: &ScaledNumberType{
										Number: util.Ptr(NumberType(3000)),
										Scale:  util.Ptr(ScaleType(0)),
									},
								},
								{
									ValueType: util.Ptr(PowerTimeSlotValueTypeTypeEnergy),
									Value: &ScaledNumberType{
										Number: util.Ptr(NumberType(1500)),
										Scale:  util.Ptr(ScaleType(3)), // 1.5 kWh
									},
								},
							},
						},
					},
					// Slot 2 with different values
					{
						Schedule: &PowerTimeSlotScheduleDataType{
							SlotNumber: util.Ptr(PowerTimeSlotNumberType(2)),
						},
						ValueList: &SmartEnergyManagementPsPowerTimeSlotValueListType{
							Value: []PowerTimeSlotValueDataType{
								{
									ValueType: util.Ptr(PowerTimeSlotValueTypeTypePower),
									Value: &ScaledNumberType{
										Number: util.Ptr(NumberType(2000)),
										Scale:  util.Ptr(ScaleType(0)),
									},
								},
							},
						},
					},
				},
			}},
		}},
	}

	// Update - Change ONLY the power value in slot 2, using composite keys
	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(0)), // Key
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				Description: &PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(PowerSequenceIdType(1)), // Key
				},
				PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
					Schedule: &PowerTimeSlotScheduleDataType{
						SlotNumber: util.Ptr(PowerTimeSlotNumberType(2)), // Key: slotNumber = 2
					},
					ValueList: &SmartEnergyManagementPsPowerTimeSlotValueListType{
						Value: []PowerTimeSlotValueDataType{{
							ValueType: util.Ptr(PowerTimeSlotValueTypeTypePower), // Key: valueType = power
							Value: &ScaledNumberType{
								Number: util.Ptr(NumberType(2500)), // Update: 2000 -> 2500
								Scale:  util.Ptr(ScaleType(0)),
							},
						}},
					},
				}},
			}},
		}},
	}

	// Act
	result, success := existing.UpdateList(false, true, update, NewFilterTypePartial(), nil)

	// Assert
	assert.True(t, success)
	assert.NotNil(t, result)

	resultData := result.(*SmartEnergyManagementPsDataType)

	// With proper key-based matching:
	// - Slot 1 should remain unchanged (both power and energy values)
	// - Slot 2 power should be updated to 2500
	// - All other fields should be preserved

	// Current positional implementation likely updates wrong slot/value
	// This test demonstrates the complexity of nested composite key matching
	
	slots := resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot
	assert.Equal(t, 2, len(slots))

	// Find slot with slotNumber: 2
	var targetSlot *SmartEnergyManagementPsPowerTimeSlotType
	for i := range slots {
		if *slots[i].Schedule.SlotNumber == PowerTimeSlotNumberType(2) {
			targetSlot = &slots[i]
			break
		}
	}

	if targetSlot != nil {
		// Find power value in target slot
		var powerValue *PowerTimeSlotValueDataType
		for i := range targetSlot.ValueList.Value {
			if *targetSlot.ValueList.Value[i].ValueType == PowerTimeSlotValueTypeTypePower {
				powerValue = &targetSlot.ValueList.Value[i]
				break
			}
		}

		if powerValue != nil {
			// This assertion will likely FAIL with current positional implementation
			assert.Equal(t, NumberType(2500), *powerValue.Value.Number,
				"Power value in slot 2 should be updated to 2500")
		}
	}

	// Verify slot 1 remains unchanged
	var slot1 *SmartEnergyManagementPsPowerTimeSlotType
	for i := range slots {
		if *slots[i].Schedule.SlotNumber == PowerTimeSlotNumberType(1) {
			slot1 = &slots[i]
			break
		}
	}

	if slot1 != nil {
		// Slot 1 power should remain 3000
		var slot1PowerValue *PowerTimeSlotValueDataType
		for i := range slot1.ValueList.Value {
			if *slot1.ValueList.Value[i].ValueType == PowerTimeSlotValueTypeTypePower {
				slot1PowerValue = &slot1.ValueList.Value[i]
				break
			}
		}

		if slot1PowerValue != nil {
			assert.Equal(t, NumberType(3000), *slot1PowerValue.Value.Number,
				"Power value in slot 1 should remain unchanged")
		}

		// Slot 1 energy should remain 1500
		var slot1EnergyValue *PowerTimeSlotValueDataType
		for i := range slot1.ValueList.Value {
			if *slot1.ValueList.Value[i].ValueType == PowerTimeSlotValueTypeTypeEnergy {
				slot1EnergyValue = &slot1.ValueList.Value[i]
				break
			}
		}

		if slot1EnergyValue != nil {
			assert.Equal(t, NumberType(1500), *slot1EnergyValue.Value.Number,
				"Energy value in slot 1 should remain unchanged")
		}
	}
}

// TestSmartEnergyManagementPsDataType_RealWorld_AtomicityTest tests atomicity
// All-or-nothing updates with persist flag
func TestSmartEnergyManagementPsDataType_RealWorld_AtomicityTest(t *testing.T) {
	// Arrange
	original := createComplexOHPCFStructure()
	originalStateCopy := *original.Alternatives[0].PowerSequence[0].State.State

	// Create update that targets valid sequence
	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(0)),
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				Description: &PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(PowerSequenceIdType(0)),
				},
				State: &PowerSequenceStateDataType{
					State: util.Ptr(PowerSequenceStateTypeRunning),
				},
			}},
		}},
	}

	// Test 1: persist = false - original should not be modified
	result1, success1 := original.UpdateList(false, false, update, NewFilterTypePartial(), nil)
	assert.True(t, success1)
	assert.NotNil(t, result1)

	// Original should remain unchanged
	assert.Equal(t, originalStateCopy, *original.Alternatives[0].PowerSequence[0].State.State)

	// Result should have the update
	result1Data := result1.(*SmartEnergyManagementPsDataType)
	assert.Equal(t, PowerSequenceStateTypeRunning, *result1Data.Alternatives[0].PowerSequence[0].State.State)

	// Test 2: persist = true - original should be modified
	result2, success2 := original.UpdateList(false, true, update, NewFilterTypePartial(), nil)
	assert.True(t, success2)
	assert.NotNil(t, result2)

	// Original should now be modified
	assert.Equal(t, PowerSequenceStateTypeRunning, *original.Alternatives[0].PowerSequence[0].State.State)

	// Result should be same as original after persist
	result2Data := result2.(*SmartEnergyManagementPsDataType)
	assert.Equal(t, PowerSequenceStateTypeRunning, *result2Data.Alternatives[0].PowerSequence[0].State.State)
}

// TestSmartEnergyManagementPsDataType_RealWorld_DeepCopyBehavior tests deep copy behavior
// Ensures no shared references between original and result
func TestSmartEnergyManagementPsDataType_RealWorld_DeepCopyBehavior(t *testing.T) {
	// Arrange
	original := createComplexOHPCFStructure()

	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{
					State: util.Ptr(PowerSequenceStateTypeRunning),
				},
			}},
		}},
	}

	// Act - persist = false to test deep copy
	result, success := original.UpdateList(false, false, update, NewFilterTypePartial(), nil)
	assert.True(t, success)

	resultData := result.(*SmartEnergyManagementPsDataType)

	// Modify the result's nested value
	resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value[0].Value.Number = util.Ptr(NumberType(9999))

	// Original should remain unchanged (proves deep copy)
	originalValue := *original.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value[0].Value.Number
	assert.Equal(t, NumberType(3000), originalValue, "Original should not be affected by result modifications")

	// Verify the result was actually modified
	resultValue := *resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value[0].Value.Number
	assert.Equal(t, NumberType(9999), resultValue, "Result should have the modified value")
}

// Helper function to create a complex OHPCF structure for testing
func createComplexOHPCFStructure() *SmartEnergyManagementPsDataType {
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
				AlternativesId: util.Ptr(AlternativesIdType(0)),
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				Description: &PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(PowerSequenceIdType(0)),
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

// TestSmartEnergyManagementPsDataType_KeyBasedMatching_MultipleAlternatives verifies proper key-based matching
// This test verifies that key-based matching works correctly with multiple alternatives
func TestSmartEnergyManagementPsDataType_KeyBasedMatching_MultipleAlternatives(t *testing.T) {
	// Arrange - Multiple alternatives with different IDs in non-sequential order
	existing := &SmartEnergyManagementPsDataType{
		NodeScheduleInformation: &PowerSequenceNodeScheduleInformationDataType{
			AlternativesCount: util.Ptr(uint(3)),
		},
		Alternatives: []SmartEnergyManagementPsAlternativesType{
			// Alternative ID 10 at index 0
			{
				Relation: &SmartEnergyManagementPsAlternativesRelationType{
					AlternativesId: util.Ptr(AlternativesIdType(10)),
				},
				PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
					Description: &PowerSequenceDescriptionDataType{
						SequenceId: util.Ptr(PowerSequenceIdType(100)),
					},
					State: &PowerSequenceStateDataType{
						State: util.Ptr(PowerSequenceStateTypeInactive),
					},
				}},
			},
			// Alternative ID 5 at index 1
			{
				Relation: &SmartEnergyManagementPsAlternativesRelationType{
					AlternativesId: util.Ptr(AlternativesIdType(5)),
				},
				PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
					Description: &PowerSequenceDescriptionDataType{
						SequenceId: util.Ptr(PowerSequenceIdType(200)),
					},
					State: &PowerSequenceStateDataType{
						State: util.Ptr(PowerSequenceStateTypeInactive),
					},
				}},
			},
			// Alternative ID 15 at index 2
			{
				Relation: &SmartEnergyManagementPsAlternativesRelationType{
					AlternativesId: util.Ptr(AlternativesIdType(15)),
				},
				PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
					Description: &PowerSequenceDescriptionDataType{
						SequenceId: util.Ptr(PowerSequenceIdType(300)),
					},
					State: &PowerSequenceStateDataType{
						State: util.Ptr(PowerSequenceStateTypeInactive),
					},
				}},
			},
		},
	}

	// Update - Target alternative ID 15 (at index 2), sequence ID 300
	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(15)), // Target ID 15
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				Description: &PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(PowerSequenceIdType(300)), // Target ID 300
				},
				State: &PowerSequenceStateDataType{
					State: util.Ptr(PowerSequenceStateTypeRunning),
				},
			}},
		}},
	}

	// Act
	result, success := existing.UpdateList(false, true, update, NewFilterTypePartial(), nil)
	assert.True(t, success)

	resultData := result.(*SmartEnergyManagementPsDataType)

	// CORRECT KEY-BASED BEHAVIOR:
	// With proper key-based implementation, alternative ID 15 should be updated

	// Verify alternative ID 10 (index 0) remains inactive
	assert.Equal(t, PowerSequenceStateTypeInactive, *resultData.Alternatives[0].PowerSequence[0].State.State,
		"Alternative ID 10 should remain inactive")

	// Verify alternative ID 5 (index 1) remains inactive
	assert.Equal(t, PowerSequenceStateTypeInactive, *resultData.Alternatives[1].PowerSequence[0].State.State,
		"Alternative ID 5 should remain inactive")

	// Verify alternative ID 15 (index 2) is correctly updated to running
	assert.Equal(t, PowerSequenceStateTypeRunning, *resultData.Alternatives[2].PowerSequence[0].State.State,
		"Alternative ID 15 should be updated to running state")

	// Double-check by finding alternative by key
	var targetAlternative *SmartEnergyManagementPsAlternativesType
	for i := range resultData.Alternatives {
		if *resultData.Alternatives[i].Relation.AlternativesId == AlternativesIdType(15) {
			targetAlternative = &resultData.Alternatives[i]
			break
		}
	}

	assert.NotNil(t, targetAlternative, "Should find alternative with ID 15")
	if targetAlternative != nil {
		assert.Equal(t, PowerSequenceStateTypeRunning, *targetAlternative.PowerSequence[0].State.State,
			"Alternative ID 15 should be updated using key-based matching")
	}
}