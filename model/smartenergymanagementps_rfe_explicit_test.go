package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

// filterPartialWithSelectors builds a FilterType carrying the given selectors.
func filterPartialWithSelectors(sel *SmartEnergyManagementPsDataSelectorsType) *FilterType {
	f := NewFilterTypePartial()
	f.SmartEnergyManagementPsDataSelectors = sel
	return f
}

// baseExplicitRFEData builds a two-alternative, two-sequence-each, two-slot-each fixture.
func baseExplicitRFEData() *SmartEnergyManagementPsDataType {
	mkSeq := func(seqId uint, state PowerSequenceStateType) SmartEnergyManagementPsPowerSequenceType {
		return SmartEnergyManagementPsPowerSequenceType{
			Description: &PowerSequenceDescriptionDataType{
				SequenceId: util.Ptr(PowerSequenceIdType(seqId)),
			},
			State: &PowerSequenceStateDataType{
				State: util.Ptr(state),
			},
			PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{
				{
					Schedule: &PowerTimeSlotScheduleDataType{
						SlotNumber:      util.Ptr(PowerTimeSlotNumberType(1)),
						DefaultDuration: util.Ptr(DurationType("PT10M")),
					},
					ValueList: &SmartEnergyManagementPsPowerTimeSlotValueListType{
						Value: []PowerTimeSlotValueDataType{
							{ValueType: util.Ptr(PowerTimeSlotValueTypeTypePower), Value: &ScaledNumberType{Number: util.Ptr(NumberType(1000)), Scale: util.Ptr(ScaleType(0))}},
							{ValueType: util.Ptr(PowerTimeSlotValueTypeTypePowerMax), Value: &ScaledNumberType{Number: util.Ptr(NumberType(2000)), Scale: util.Ptr(ScaleType(0))}},
						},
					},
				},
				{
					Schedule: &PowerTimeSlotScheduleDataType{
						SlotNumber:      util.Ptr(PowerTimeSlotNumberType(2)),
						DefaultDuration: util.Ptr(DurationType("PT20M")),
					},
					ValueList: &SmartEnergyManagementPsPowerTimeSlotValueListType{
						Value: []PowerTimeSlotValueDataType{
							{ValueType: util.Ptr(PowerTimeSlotValueTypeTypePower), Value: &ScaledNumberType{Number: util.Ptr(NumberType(500)), Scale: util.Ptr(ScaleType(0))}},
						},
					},
				},
			},
		}
	}

	return &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{
			{
				Relation: &SmartEnergyManagementPsAlternativesRelationType{
					AlternativesId: util.Ptr(AlternativesIdType(1)),
					SequenceId:     []PowerSequenceIdType{10, 11},
				},
				PowerSequence: []SmartEnergyManagementPsPowerSequenceType{
					mkSeq(10, PowerSequenceStateTypeInactive),
					mkSeq(11, PowerSequenceStateTypeInactive),
				},
			},
			{
				Relation: &SmartEnergyManagementPsAlternativesRelationType{
					AlternativesId: util.Ptr(AlternativesIdType(2)),
					SequenceId:     []PowerSequenceIdType{20, 21},
				},
				PowerSequence: []SmartEnergyManagementPsPowerSequenceType{
					mkSeq(20, PowerSequenceStateTypeInactive),
					mkSeq(21, PowerSequenceStateTypeInactive),
				},
			},
		},
	}
}

// TestExplicitRFE_AlternativesRelation_AlternativesId verifies that
// AlternativesRelation.AlternativesId restricts updates to the matching alternative only.
func TestExplicitRFE_AlternativesRelation_AlternativesId(t *testing.T) {
	existing := baseExplicitRFEData()

	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{State: util.Ptr(PowerSequenceStateTypeScheduled)},
			}},
		}},
	}

	sel := &SmartEnergyManagementPsDataSelectorsType{
		AlternativesRelation: &PowerSequenceAlternativesRelationListDataSelectorsType{
			AlternativesId: util.Ptr(AlternativesIdType(1)),
		},
	}

	result, ok := existing.UpdateList(false, false, update, filterPartialWithSelectors(sel), nil, nil)
	assert.True(t, ok)

	rd := result.(*SmartEnergyManagementPsDataType)

	// Alt 1 sequences updated
	assert.Equal(t, PowerSequenceStateTypeScheduled, *rd.Alternatives[0].PowerSequence[0].State.State)
	assert.Equal(t, PowerSequenceStateTypeScheduled, *rd.Alternatives[0].PowerSequence[1].State.State)

	// Alt 2 sequences untouched
	assert.Equal(t, PowerSequenceStateTypeInactive, *rd.Alternatives[1].PowerSequence[0].State.State)
	assert.Equal(t, PowerSequenceStateTypeInactive, *rd.Alternatives[1].PowerSequence[1].State.State)

	// Original not modified
	assert.Equal(t, PowerSequenceStateTypeInactive, *existing.Alternatives[0].PowerSequence[0].State.State)
}

// TestExplicitRFE_PowerSequenceDescription_SequenceId verifies that
// PowerSequenceDescription.SequenceId restricts to specific sequences only.
func TestExplicitRFE_PowerSequenceDescription_SequenceId(t *testing.T) {
	existing := baseExplicitRFEData()

	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{State: util.Ptr(PowerSequenceStateTypeScheduled)},
			}},
		}},
	}

	sel := &SmartEnergyManagementPsDataSelectorsType{
		PowerSequenceDescription: &PowerSequenceDescriptionListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(10)),
		},
	}

	result, ok := existing.UpdateList(false, false, update, filterPartialWithSelectors(sel), nil, nil)
	assert.True(t, ok)
	rd := result.(*SmartEnergyManagementPsDataType)

	// SeqId 10 updated; 11, 20 and 21 untouched
	assert.Equal(t, PowerSequenceStateTypeScheduled, *rd.Alternatives[0].PowerSequence[0].State.State, "seqId=10")
	assert.Equal(t, PowerSequenceStateTypeInactive, *rd.Alternatives[0].PowerSequence[1].State.State, "seqId=11")
	assert.Equal(t, PowerSequenceStateTypeInactive, *rd.Alternatives[1].PowerSequence[0].State.State, "seqId=20")
	assert.Equal(t, PowerSequenceStateTypeInactive, *rd.Alternatives[1].PowerSequence[1].State.State, "seqId=21")
}

// TestExplicitRFE_PowerTimeSlotSchedule_SequenceAndSlot verifies that
// PowerTimeSlotSchedule.SequenceId + SlotNumber targets a specific slot.
func TestExplicitRFE_PowerTimeSlotSchedule_SequenceAndSlot(t *testing.T) {
	existing := baseExplicitRFEData()

	newDuration := DurationType("PT99M")
	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
					Schedule: &PowerTimeSlotScheduleDataType{
						DefaultDuration: util.Ptr(newDuration),
					},
				}},
			}},
		}},
	}

	sel := &SmartEnergyManagementPsDataSelectorsType{
		PowerTimeSlotSchedule: &PowerTimeSlotScheduleListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(10)),
			SlotNumber: util.Ptr(PowerTimeSlotNumberType(2)),
		},
	}

	result, ok := existing.UpdateList(false, false, update, filterPartialWithSelectors(sel), nil, nil)
	assert.True(t, ok)
	rd := result.(*SmartEnergyManagementPsDataType)

	seq10 := rd.Alternatives[0].PowerSequence[0]
	// Slot 1 of seq10 unchanged
	assert.Equal(t, DurationType("PT10M"), *seq10.PowerTimeSlot[0].Schedule.DefaultDuration, "slot 1 unchanged")
	// Slot 2 of seq10 updated
	assert.Equal(t, newDuration, *seq10.PowerTimeSlot[1].Schedule.DefaultDuration, "slot 2 updated")

	// Seq 11 slots untouched
	seq11 := rd.Alternatives[0].PowerSequence[1]
	assert.Equal(t, DurationType("PT10M"), *seq11.PowerTimeSlot[0].Schedule.DefaultDuration)
	assert.Equal(t, DurationType("PT20M"), *seq11.PowerTimeSlot[1].Schedule.DefaultDuration)
}

// TestExplicitRFE_PowerTimeSlotValue_ValueType verifies that
// PowerTimeSlotValue.ValueType targets a specific value within a slot.
func TestExplicitRFE_PowerTimeSlotValue_ValueType(t *testing.T) {
	existing := baseExplicitRFEData()

	// Update only the "power" value in seq 10, slot 1 — not powerMax
	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
					ValueList: &SmartEnergyManagementPsPowerTimeSlotValueListType{
						Value: []PowerTimeSlotValueDataType{{
							// No ValueType identifier — selector provides the target
							Value: &ScaledNumberType{Number: util.Ptr(NumberType(9999)), Scale: util.Ptr(ScaleType(0))},
						}},
					},
				}},
			}},
		}},
	}

	sel := &SmartEnergyManagementPsDataSelectorsType{
		PowerTimeSlotValue: &PowerTimeSlotValueListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(10)),
			SlotNumber: util.Ptr(PowerTimeSlotNumberType(1)),
			ValueType:  util.Ptr(PowerTimeSlotValueTypeTypePower),
		},
	}

	result, ok := existing.UpdateList(false, false, update, filterPartialWithSelectors(sel), nil, nil)
	assert.True(t, ok)
	rd := result.(*SmartEnergyManagementPsDataType)

	slot1 := rd.Alternatives[0].PowerSequence[0].PowerTimeSlot[0]
	// Power value updated
	assert.Equal(t, NumberType(9999), *slot1.ValueList.Value[0].Value.Number, "power value updated")
	// PowerMax value untouched
	assert.Equal(t, NumberType(2000), *slot1.ValueList.Value[1].Value.Number, "powerMax unchanged")
}

// TestExplicitRFE_CombinedSelectors_AltAndSeq verifies that AlternativesRelation
// combined with PowerSequenceDescription acts as AND (both must match).
func TestExplicitRFE_CombinedSelectors_AltAndSeq(t *testing.T) {
	existing := baseExplicitRFEData()

	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{State: util.Ptr(PowerSequenceStateTypeRunning)},
			}},
		}},
	}

	// Alt 1, seqId 10 only
	sel := &SmartEnergyManagementPsDataSelectorsType{
		AlternativesRelation: &PowerSequenceAlternativesRelationListDataSelectorsType{
			AlternativesId: util.Ptr(AlternativesIdType(1)),
		},
		PowerSequenceDescription: &PowerSequenceDescriptionListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(10)),
		},
	}

	result, ok := existing.UpdateList(false, false, update, filterPartialWithSelectors(sel), nil, nil)
	assert.True(t, ok)
	rd := result.(*SmartEnergyManagementPsDataType)

	assert.Equal(t, PowerSequenceStateTypeRunning, *rd.Alternatives[0].PowerSequence[0].State.State, "alt1 seqId=10 updated")
	assert.Equal(t, PowerSequenceStateTypeInactive, *rd.Alternatives[0].PowerSequence[1].State.State, "alt1 seqId=11 untouched")
	assert.Equal(t, PowerSequenceStateTypeInactive, *rd.Alternatives[1].PowerSequence[0].State.State, "alt2 untouched")
}

// TestExplicitRFE_NilSelectors_MatchAll verifies that a selector struct with
// all nil fields matches everything (same as no selector at all — but via explicit path).
func TestExplicitRFE_NilSelectors_MatchAll(t *testing.T) {
	existing := baseExplicitRFEData()

	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{State: util.Ptr(PowerSequenceStateTypeScheduled)},
			}},
		}},
	}

	// Empty selector — all fields nil — should match every alternative and sequence
	sel := &SmartEnergyManagementPsDataSelectorsType{}

	result, ok := existing.UpdateList(false, false, update, filterPartialWithSelectors(sel), nil, nil)
	assert.True(t, ok)
	rd := result.(*SmartEnergyManagementPsDataType)

	for i, alt := range rd.Alternatives {
		for j, seq := range alt.PowerSequence {
			assert.Equal(t, PowerSequenceStateTypeScheduled, *seq.State.State,
				"alt[%d].seq[%d] should be updated", i, j)
		}
	}
}

// TestExplicitRFE_PersistTrue_ModifiesOriginal verifies the persist flag is honoured.
func TestExplicitRFE_PersistTrue_ModifiesOriginal(t *testing.T) {
	existing := baseExplicitRFEData()

	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{State: util.Ptr(PowerSequenceStateTypeRunning)},
			}},
		}},
	}

	sel := &SmartEnergyManagementPsDataSelectorsType{
		PowerSequenceDescription: &PowerSequenceDescriptionListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(10)),
		},
	}

	_, ok := existing.UpdateList(false, true, update, filterPartialWithSelectors(sel), nil, nil)
	assert.True(t, ok)

	// persist=true: original is modified in place
	assert.Equal(t, PowerSequenceStateTypeRunning, *existing.Alternatives[0].PowerSequence[0].State.State)
	// Other sequence untouched
	assert.Equal(t, PowerSequenceStateTypeInactive, *existing.Alternatives[0].PowerSequence[1].State.State)
}

// TestExplicitRFE_SlotScheduleFields verifies that slot-level Schedule fields
// (DefaultDuration, SlotActivated, TimePeriod) are correctly merged.
func TestExplicitRFE_SlotScheduleFields(t *testing.T) {
	existing := baseExplicitRFEData()

	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
					Schedule: &PowerTimeSlotScheduleDataType{
						DefaultDuration: util.Ptr(DurationType("PT45M")),
						SlotActivated:   util.Ptr(true),
					},
				}},
			}},
		}},
	}

	sel := &SmartEnergyManagementPsDataSelectorsType{
		PowerTimeSlotSchedule: &PowerTimeSlotScheduleListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(10)),
			SlotNumber: util.Ptr(PowerTimeSlotNumberType(1)),
		},
	}

	result, ok := existing.UpdateList(false, false, update, filterPartialWithSelectors(sel), nil, nil)
	assert.True(t, ok)
	rd := result.(*SmartEnergyManagementPsDataType)

	slot := rd.Alternatives[0].PowerSequence[0].PowerTimeSlot[0]
	assert.Equal(t, DurationType("PT45M"), *slot.Schedule.DefaultDuration)
	assert.Equal(t, true, *slot.Schedule.SlotActivated)
	// SlotNumber key not overwritten
	assert.Equal(t, PowerTimeSlotNumberType(1), *slot.Schedule.SlotNumber)
}

// TestExplicitRFE_NoMatchingEntry_NoChange verifies that a selector that matches
// nothing leaves the data unchanged and still returns success.
func TestExplicitRFE_NoMatchingEntry_NoChange(t *testing.T) {
	existing := baseExplicitRFEData()
	original := existing.deepCopy()

	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				State: &PowerSequenceStateDataType{State: util.Ptr(PowerSequenceStateTypeRunning)},
			}},
		}},
	}

	sel := &SmartEnergyManagementPsDataSelectorsType{
		AlternativesRelation: &PowerSequenceAlternativesRelationListDataSelectorsType{
			AlternativesId: util.Ptr(AlternativesIdType(99)), // does not exist
		},
	}

	result, ok := existing.UpdateList(false, false, update, filterPartialWithSelectors(sel), nil, nil)
	assert.True(t, ok)
	rd := result.(*SmartEnergyManagementPsDataType)
	assert.Equal(t, original, rd)
}
