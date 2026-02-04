package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

// TestSmartEnergyManagementPsDataType_PersistFlagBehavior validates correct persist flag behavior
func TestSmartEnergyManagementPsDataType_PersistFlagBehavior(t *testing.T) {
	// Test 1: persist=true should modify original
	t.Run("persist_true_modifies_original", func(t *testing.T) {
		// Arrange
		original := &SmartEnergyManagementPsDataType{
			Alternatives: []SmartEnergyManagementPsAlternativesType{{
				Relation: &SmartEnergyManagementPsAlternativesRelationType{
					AlternativesId: util.Ptr(AlternativesIdType(0)),
				},
				PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
					Description: &PowerSequenceDescriptionDataType{
						SequenceId: util.Ptr(PowerSequenceIdType(0)),
					},
					State: &PowerSequenceStateDataType{
						State: util.Ptr(PowerSequenceStateTypeInactive),
					},
				}},
			}},
		}

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

		// Act - persist=true
		result, success := original.UpdateList(false, true, update, NewFilterTypePartial(), nil, nil)

		// Assert
		assert.True(t, success)
		assert.NotNil(t, result)
		
		// Original SHOULD be modified when persist=true
		assert.Equal(t, PowerSequenceStateTypeRunning, *original.Alternatives[0].PowerSequence[0].State.State,
			"Original should be modified when persist=true")
		
		// Result should point to the modified original
		resultData := result.(*SmartEnergyManagementPsDataType)
		assert.Equal(t, PowerSequenceStateTypeRunning, *resultData.Alternatives[0].PowerSequence[0].State.State)
	})

	// Test 2: persist=false should NOT modify original
	t.Run("persist_false_preserves_original", func(t *testing.T) {
		// Arrange
		original := &SmartEnergyManagementPsDataType{
			Alternatives: []SmartEnergyManagementPsAlternativesType{{
				Relation: &SmartEnergyManagementPsAlternativesRelationType{
					AlternativesId: util.Ptr(AlternativesIdType(0)),
				},
				PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
					Description: &PowerSequenceDescriptionDataType{
						SequenceId: util.Ptr(PowerSequenceIdType(0)),
					},
					State: &PowerSequenceStateDataType{
						State: util.Ptr(PowerSequenceStateTypeInactive),
					},
				}},
			}},
		}

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

		// Act - persist=false
		result, success := original.UpdateList(false, false, update, NewFilterTypePartial(), nil, nil)

		// Assert
		assert.True(t, success)
		assert.NotNil(t, result)
		
		// Original should NOT be modified when persist=false
		assert.Equal(t, PowerSequenceStateTypeInactive, *original.Alternatives[0].PowerSequence[0].State.State,
			"Original should NOT be modified when persist=false")
		
		// Result should be a new instance with the updates
		resultData := result.(*SmartEnergyManagementPsDataType)
		assert.Equal(t, PowerSequenceStateTypeRunning, *resultData.Alternatives[0].PowerSequence[0].State.State,
			"Result should contain the updates")
	})

	// Test 3: Verify memory independence when persist=false
	t.Run("persist_false_creates_independent_copy", func(t *testing.T) {
		// Arrange
		original := &SmartEnergyManagementPsDataType{
			Alternatives: []SmartEnergyManagementPsAlternativesType{{
				Relation: &SmartEnergyManagementPsAlternativesRelationType{
					AlternativesId: util.Ptr(AlternativesIdType(0)),
				},
				PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
					Description: &PowerSequenceDescriptionDataType{
						SequenceId: util.Ptr(PowerSequenceIdType(0)),
					},
					State: &PowerSequenceStateDataType{
						State: util.Ptr(PowerSequenceStateTypeInactive),
					},
				}},
			}},
		}

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
						State: util.Ptr(PowerSequenceStateTypeScheduled),
					},
				}},
			}},
		}

		// Act - persist=false
		result, success := original.UpdateList(false, false, update, NewFilterTypePartial(), nil, nil)

		// Assert
		assert.True(t, success)
		resultData := result.(*SmartEnergyManagementPsDataType)
		
		// Further modify the result
		*resultData.Alternatives[0].PowerSequence[0].State.State = PowerSequenceStateTypeCompleted
		
		// Original should still have its original value
		assert.Equal(t, PowerSequenceStateTypeInactive, *original.Alternatives[0].PowerSequence[0].State.State,
			"Original should remain unchanged after modifying result")
	})
}