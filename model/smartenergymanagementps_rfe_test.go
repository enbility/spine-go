package model

// SmartEnergyManagementPs Restricted Function Exchange (RFE) Tests
//
// This file tests the remoteWrite=true execution path of UpdateList, which
// corresponds to a remote power-sequences client writing to the local server
// via a SPINE "write (restricted)" message (RFE per §5.3.4).
//
// Permission authority:
//   - Table 119 §4.3.24.1: only "powerSequences configuration" and
//     "powerSequences state change" use write(restricted) on smartEnergyManagementPsData
//   - Table 123 §4.3.24.4.4.1: permitted fields for configuration write
//   - Table 124 §4.3.24.4.5.1: permitted fields for state change write
//
// All other fields are server-owned output-only (Table 119 Out direction).
// Presence of ANY server-owned field in a remoteWrite=true payload MUST cause
// the entire UpdateList call to be rejected (success=false, data unchanged).
//

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Fixture helpers
// ---------------------------------------------------------------------------

// makeRFEBase returns a rich SmartEnergyManagementPsDataType with all field
// categories populated so blocked-field tests can verify data is unchanged.
// Layout:
//   - NodeScheduleInformation: all 5 fields set
//   - Alternatives[0]: AlternativesId=1, one sequence (SequenceId=1)
//   - Description with label and powerUnit
//   - State with all 7 fields set (server-owned values)
//   - Schedule with StartTime and EndTime
//   - ScheduleConstraints at sequence level
//   - SchedulePreference
//   - OperatingConstraintsInterrupt
//   - Slot 1 (constrained: MinDuration+MaxDuration+OptionalSlot): full schedule+constraints+valueList
//   - Slot 2 (implicit: no ScheduleConstraints): minimal schedule+valueList
//   - Alternatives[1]: AlternativesId=2, one sequence (SequenceId=2), minimal
func makeRFEBase() *SmartEnergyManagementPsDataType {
	slot1 := makeConstrainedDurationOptionalSlot(1)
	slot2 := makeImplicitSlot(2)

	return &SmartEnergyManagementPsDataType{
		NodeScheduleInformation: &PowerSequenceNodeScheduleInformationDataType{
			NodeRemoteControllable:           util.Ptr(true),
			SupportsSingleSlotSchedulingOnly: util.Ptr(false),
			AlternativesCount:                util.Ptr(uint(2)),
			TotalSequencesCountMax:           util.Ptr(uint(2)),
			SupportsReselection:              util.Ptr(true),
		},
		Alternatives: []SmartEnergyManagementPsAlternativesType{
			{
				Relation: &SmartEnergyManagementPsAlternativesRelationType{
					AlternativesId: util.Ptr(AlternativesIdType(1)),
				},
				PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
					Description: &PowerSequenceDescriptionDataType{
						SequenceId:  util.Ptr(PowerSequenceIdType(1)),
						Description: util.Ptr(DescriptionType("existing-desc")),
						PowerUnit:   util.Ptr(UnitOfMeasurementTypeW),
					},
					State: &PowerSequenceStateDataType{
						State:                      util.Ptr(PowerSequenceStateTypeInactive),
						ActiveSlotNumber:           util.Ptr(PowerTimeSlotNumberType(5)),
						ElapsedSlotTime:            util.Ptr(DurationType("PT5M")),
						RemainingSlotTime:          util.Ptr(DurationType("PT25M")),
						SequenceRemoteControllable: util.Ptr(true),
						ActiveRepetitionNumber:     util.Ptr(uint(0)),
						RemainingPauseTime:         util.Ptr(DurationType("PT0S")),
					},
					Schedule: &PowerSequenceScheduleDataType{
						StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT1H")),
						EndTime:   util.Ptr(AbsoluteOrRelativeTimeType("PT2H")),
					},
					ScheduleConstraints: &PowerSequenceScheduleConstraintsDataType{
						EarliestStartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT0S")),
						LatestEndTime:     util.Ptr(AbsoluteOrRelativeTimeType("PT3H")),
					},
					SchedulePreference: &PowerSequenceSchedulePreferenceDataType{
						Greenest: util.Ptr(true),
						Cheapest: util.Ptr(false),
					},
					OperatingConstraintsInterrupt: &OperatingConstraintsInterruptDataType{
						IsPausable: util.Ptr(true),
					},
					PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{slot1, slot2},
				}},
			},
			{
				Relation: &SmartEnergyManagementPsAlternativesRelationType{
					AlternativesId: util.Ptr(AlternativesIdType(2)),
				},
				PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
					Description: &PowerSequenceDescriptionDataType{
						SequenceId: util.Ptr(PowerSequenceIdType(2)),
						PowerUnit:  util.Ptr(UnitOfMeasurementTypeW),
					},
					State: &PowerSequenceStateDataType{
						State: util.Ptr(PowerSequenceStateTypeInactive),
					},
					Schedule: &PowerSequenceScheduleDataType{
						StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT2H")),
					},
				}},
			},
		},
	}
}

// makeConstrainedDurationOptionalSlot returns a slot with both duration and
// activation configurable (§4.3.24.3.2.7).
func makeConstrainedDurationOptionalSlot(slotNumber PowerTimeSlotNumberType) SmartEnergyManagementPsPowerTimeSlotType {
	return SmartEnergyManagementPsPowerTimeSlotType{
		Schedule: &PowerTimeSlotScheduleDataType{
			SlotNumber: util.Ptr(slotNumber),
			TimePeriod: &TimePeriodType{
				StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT0S")),
			},
			DefaultDuration:     util.Ptr(DurationType("PT20M")),
			DurationUncertainty: util.Ptr(DurationType("PT1M")),
			SlotActivated:       util.Ptr(true),
		},
		ScheduleConstraints: &PowerTimeSlotScheduleConstraintsDataType{
			MinDuration:  util.Ptr(DurationType("PT10M")),
			MaxDuration:  util.Ptr(DurationType("PT30M")),
			OptionalSlot: util.Ptr(true),
		},
		ValueList: &SmartEnergyManagementPsPowerTimeSlotValueListType{
			Value: []PowerTimeSlotValueDataType{{
				ValueType: util.Ptr(PowerTimeSlotValueTypeTypePower),
				Value:     &ScaledNumberType{Number: util.Ptr(NumberType(3000)), Scale: util.Ptr(ScaleType(0))},
			}},
		},
	}
}

// makeImplicitSlot returns a slot with no ScheduleConstraints (implicit schedule).
// Implicit slots are NOT configurable (§4.3.24.3.2.6.2).
func makeImplicitSlot(slotNumber PowerTimeSlotNumberType) SmartEnergyManagementPsPowerTimeSlotType {
	return SmartEnergyManagementPsPowerTimeSlotType{
		Schedule: &PowerTimeSlotScheduleDataType{
			SlotNumber:      util.Ptr(slotNumber),
			DefaultDuration: util.Ptr(DurationType("PT15M")),
		},
		// ScheduleConstraints intentionally nil → implicit
		ValueList: &SmartEnergyManagementPsPowerTimeSlotValueListType{
			Value: []PowerTimeSlotValueDataType{{
				ValueType: util.Ptr(PowerTimeSlotValueTypeTypePower),
				Value:     &ScaledNumberType{Number: util.Ptr(NumberType(1500)), Scale: util.Ptr(ScaleType(0))},
			}},
		},
	}
}

// rfeSeqUpdate builds a minimal remoteWrite=true payload targeting sequenceId=1
// in alternativesId=1 with the provided sequence content.
func rfeSeqUpdate(seq SmartEnergyManagementPsPowerSequenceType) *SmartEnergyManagementPsDataType {
	return &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(1)),
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{seq},
		}},
	}
}

// rfeSeqWithKey builds a rfeSeqUpdate payload where the sequence also carries
// the sequenceId key (needed for key-based matching).
func rfeSeqWithKey(seq SmartEnergyManagementPsPowerSequenceType) *SmartEnergyManagementPsDataType {
	seq.Description = &PowerSequenceDescriptionDataType{
		SequenceId: util.Ptr(PowerSequenceIdType(1)),
	}
	return rfeSeqUpdate(seq)
}

// ---------------------------------------------------------------------------
// Section 1+2: Full replace and filterDelete rejection (TC-3, TC-5)
// ---------------------------------------------------------------------------

// TC-3: remoteWrite=true + filterPartial=nil → RFE requires partial write; full
// replace must be rejected.
// Spec: §5.3.4 (RFE write must use filterPartial)
func TestRFE_RemoteWrite_FullReplace_Rejected(t *testing.T) {
	existing := makeRFEBase()
	originalState := *existing.Alternatives[0].PowerSequence[0].State.State

	update := makeRFEBase() // arbitrary full payload

	result, success := existing.UpdateList(true, true, update, nil, nil, nil)

	assert.False(t, success, "remoteWrite=true full replace (filterPartial=nil) must be rejected")
	assert.Equal(t, existing, result, "rejected call must return original")
	assert.Equal(t, originalState, *existing.Alternatives[0].PowerSequence[0].State.State,
		"existing data must be unchanged after rejection")
}

// TC-5: remoteWrite=true + filterDelete set → not supported.
func TestRFE_RemoteWrite_FilterDelete_Rejected(t *testing.T) {
	existing := makeRFEBase()

	deleteFilter := &FilterType{CmdControl: &CmdControlType{Delete: &ElementTagType{}}}
	result, success := existing.UpdateList(true, true, existing, NewFilterTypePartial(), deleteFilter, nil)

	assert.False(t, success, "remoteWrite=true with filterDelete must be rejected")
	assert.Equal(t, existing, result)
}

// ---------------------------------------------------------------------------
// Section 3: Remote write blocked fields (TC-6 – TC-31)
// All cases assert success=false and that existing data is unchanged.
// Spec: Table 119 (server Out only), Tables 123+124 (only permitted writes)
// ---------------------------------------------------------------------------

// TC-6, TC-7, TC-8, TC-93: nodeScheduleInformation fields
// Table 119: nodeScheduleInformation is server Out only; not present in Tables 123/124.
func TestRFE_NodeScheduleInformation_Rejected(t *testing.T) {
	tests := []struct {
		name   string
		update *SmartEnergyManagementPsDataType
	}{
		{
			name: "nodeRemoteControllable",
			update: &SmartEnergyManagementPsDataType{
				NodeScheduleInformation: &PowerSequenceNodeScheduleInformationDataType{
					NodeRemoteControllable: util.Ptr(false), // existing is true
				},
			},
		},
		{
			name: "alternativesCount",
			update: &SmartEnergyManagementPsDataType{
				NodeScheduleInformation: &PowerSequenceNodeScheduleInformationDataType{
					AlternativesCount: util.Ptr(uint(99)),
				},
			},
		},
		{
			name: "totalSequencesCountMax",
			update: &SmartEnergyManagementPsDataType{
				NodeScheduleInformation: &PowerSequenceNodeScheduleInformationDataType{
					TotalSequencesCountMax: util.Ptr(uint(99)),
				},
			},
		},
		{
			name: "supportsSingleSlotSchedulingOnly",
			update: &SmartEnergyManagementPsDataType{
				NodeScheduleInformation: &PowerSequenceNodeScheduleInformationDataType{
					SupportsSingleSlotSchedulingOnly: util.Ptr(true),
				},
			},
		},
		{
			name: "supportsReselection",
			update: &SmartEnergyManagementPsDataType{
				NodeScheduleInformation: &PowerSequenceNodeScheduleInformationDataType{
					SupportsReselection: util.Ptr(false), // existing is true
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			existing := makeRFEBase()
			origCount := *existing.NodeScheduleInformation.AlternativesCount

			result, success := existing.UpdateList(true, true, tt.update, NewFilterTypePartial(), nil, nil)

			assert.False(t, success, "nodeScheduleInformation write must be rejected")
			assert.Equal(t, existing, result)
			assert.Equal(t, origCount, *existing.NodeScheduleInformation.AlternativesCount,
				"nodeScheduleInformation must be unchanged")
		})
	}
}

// TC-9, TC-10: powerSequence.description fields
// Table 119: description is server Out only; not in Tables 123/124.
func TestRFE_DescriptionFields_Rejected(t *testing.T) {
	tests := []struct {
		name string
		desc *PowerSequenceDescriptionDataType
	}{
		{
			name: "description_label",
			desc: &PowerSequenceDescriptionDataType{
				SequenceId:  util.Ptr(PowerSequenceIdType(1)),
				Description: util.Ptr(DescriptionType("hacked-label")),
			},
		},
		{
			name: "powerUnit_change",
			desc: &PowerSequenceDescriptionDataType{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				PowerUnit:  util.Ptr(UnitOfMeasurementTypeW),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			existing := makeRFEBase()
			origDesc := *existing.Alternatives[0].PowerSequence[0].Description.Description

			update := rfeSeqUpdate(SmartEnergyManagementPsPowerSequenceType{
				Description: tt.desc,
			})
			result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

			assert.False(t, success, "description field write must be rejected")
			assert.Equal(t, existing, result)
			assert.Equal(t, origDesc, *existing.Alternatives[0].PowerSequence[0].Description.Description,
				"description must be unchanged")
		})
	}
}

// TC-11 – TC-14: powerSequence.state sub-fields (excluding state.state which IS allowed)
// Table 124 permits ONLY state.state; all other state sub-fields are server-assigned.
func TestRFE_StateSubfields_Rejected(t *testing.T) {
	tests := []struct {
		name  string
		state *PowerSequenceStateDataType
	}{
		{
			name:  "activeSlotNumber",
			state: &PowerSequenceStateDataType{ActiveSlotNumber: util.Ptr(PowerTimeSlotNumberType(99))},
		},
		{
			name:  "elapsedSlotTime",
			state: &PowerSequenceStateDataType{ElapsedSlotTime: util.Ptr(DurationType("PT99M"))},
		},
		{
			name:  "remainingSlotTime",
			state: &PowerSequenceStateDataType{RemainingSlotTime: util.Ptr(DurationType("PT99M"))},
		},
		{
			name:  "remainingPauseTime",
			state: &PowerSequenceStateDataType{RemainingPauseTime: util.Ptr(DurationType("PT99M"))},
		},
		{
			name:  "activeRepetitionNumber",
			state: &PowerSequenceStateDataType{ActiveRepetitionNumber: util.Ptr(uint(99))},
		},
		{
			name:  "sequenceRemoteControllable",
			state: &PowerSequenceStateDataType{SequenceRemoteControllable: util.Ptr(false)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			existing := makeRFEBase()
			origActive := *existing.Alternatives[0].PowerSequence[0].State.ActiveSlotNumber

			update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
				State: tt.state,
			})
			result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

			assert.False(t, success, "state sub-field (server-owned) write must be rejected")
			assert.Equal(t, existing, result)
			assert.Equal(t, origActive, *existing.Alternatives[0].PowerSequence[0].State.ActiveSlotNumber,
				"state must be unchanged")
		})
	}
}

// TC-15, TC-16: powerSequence.scheduleConstraints fields
// Table 119: sequence-level scheduleConstraints is server Out only; absent from Table 123.
func TestRFE_SequenceScheduleConstraints_Rejected(t *testing.T) {
	tests := []struct {
		name        string
		constraints *PowerSequenceScheduleConstraintsDataType
	}{
		{
			name: "earliestStartTime",
			constraints: &PowerSequenceScheduleConstraintsDataType{
				EarliestStartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT99H")),
			},
		},
		{
			name: "latestEndTime",
			constraints: &PowerSequenceScheduleConstraintsDataType{
				LatestEndTime: util.Ptr(AbsoluteOrRelativeTimeType("PT99H")),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			existing := makeRFEBase()
			origEarliestStart := *existing.Alternatives[0].PowerSequence[0].ScheduleConstraints.EarliestStartTime

			update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
				ScheduleConstraints: tt.constraints,
			})
			result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

			assert.False(t, success, "sequence-level scheduleConstraints write must be rejected")
			assert.Equal(t, existing, result)
			assert.Equal(t, origEarliestStart,
				*existing.Alternatives[0].PowerSequence[0].ScheduleConstraints.EarliestStartTime,
				"scheduleConstraints must be unchanged")
		})
	}
}

// TC-17, TC-18: powerSequence.schedulePreference fields
// Table 119: schedulePreference is server Out only; absent from Table 123.
func TestRFE_SchedulePreference_Rejected(t *testing.T) {
	tests := []struct {
		name string
		pref *PowerSequenceSchedulePreferenceDataType
	}{
		{
			name: "greenest",
			pref: &PowerSequenceSchedulePreferenceDataType{Greenest: util.Ptr(false)},
		},
		{
			name: "cheapest",
			pref: &PowerSequenceSchedulePreferenceDataType{Cheapest: util.Ptr(true)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			existing := makeRFEBase()
			origGreenest := *existing.Alternatives[0].PowerSequence[0].SchedulePreference.Greenest

			update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
				SchedulePreference: tt.pref,
			})
			result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

			assert.False(t, success, "schedulePreference write must be rejected")
			assert.Equal(t, existing, result)
			assert.Equal(t, origGreenest,
				*existing.Alternatives[0].PowerSequence[0].SchedulePreference.Greenest,
				"schedulePreference must be unchanged")
		})
	}
}

// TC-19: powerSequence.operatingConstraints* fields
// Table 119: operatingConstraints are server Out only; absent from Table 123.
func TestRFE_OperatingConstraints_Rejected(t *testing.T) {
	tests := []struct {
		name string
		seq  SmartEnergyManagementPsPowerSequenceType
	}{
		{
			name: "operatingConstraintsInterrupt",
			seq: SmartEnergyManagementPsPowerSequenceType{
				OperatingConstraintsInterrupt: &OperatingConstraintsInterruptDataType{
					IsPausable: util.Ptr(false),
				},
			},
		},
		{
			name: "operatingConstraintsDuration",
			seq: SmartEnergyManagementPsPowerSequenceType{
				OperatingConstraintsDuration: &OperatingConstraintsDurationDataType{
					ActiveDurationMin: util.Ptr(DurationType("PT5M")),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			existing := makeRFEBase()
			origPausable := *existing.Alternatives[0].PowerSequence[0].OperatingConstraintsInterrupt.IsPausable

			update := rfeSeqWithKey(tt.seq)
			result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

			assert.False(t, success, "operatingConstraints write must be rejected")
			assert.Equal(t, existing, result)
			assert.Equal(t, origPausable,
				*existing.Alternatives[0].PowerSequence[0].OperatingConstraintsInterrupt.IsPausable,
				"operatingConstraintsInterrupt must be unchanged")
		})
	}
}

// TC-20, TC-21: powerTimeSlot.schedule.timePeriod fields
// Table 123 permits only slotNumber, defaultDuration, slotActivated; timePeriod is absent.
func TestRFE_TimePeriod_Rejected(t *testing.T) {
	tests := []struct {
		name       string
		timePeriod *TimePeriodType
	}{
		{
			name:       "timePeriod_startTime",
			timePeriod: &TimePeriodType{StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT99H"))},
		},
		{
			name:       "timePeriod_endTime",
			timePeriod: &TimePeriodType{EndTime: util.Ptr(AbsoluteOrRelativeTimeType("PT99H"))},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			existing := makeRFEBase()
			origTimePeriod := existing.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].Schedule.TimePeriod.StartTime

			update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
				PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
					Schedule: &PowerTimeSlotScheduleDataType{
						SlotNumber: util.Ptr(PowerTimeSlotNumberType(1)),
						TimePeriod: tt.timePeriod,
					},
				}},
			})
			result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

			assert.False(t, success, "timePeriod write must be rejected")
			assert.Equal(t, existing, result)
			assert.Equal(t, origTimePeriod,
				existing.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].Schedule.TimePeriod.StartTime,
				"timePeriod must be unchanged")
		})
	}
}

// TC-22: powerTimeSlot.schedule.durationUncertainty
// Table 123 does not include durationUncertainty.
func TestRFE_DurationUncertainty_Rejected(t *testing.T) {
	existing := makeRFEBase()
	origUncertainty := *existing.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].Schedule.DurationUncertainty

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
			Schedule: &PowerTimeSlotScheduleDataType{
				SlotNumber:          util.Ptr(PowerTimeSlotNumberType(1)),
				DurationUncertainty: util.Ptr(DurationType("PT99M")),
			},
		}},
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.False(t, success, "durationUncertainty write must be rejected")
	assert.Equal(t, existing, result)
	assert.Equal(t, origUncertainty,
		*existing.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].Schedule.DurationUncertainty,
		"durationUncertainty must be unchanged")
}

// TC-23, TC-24, TC-25: powerTimeSlot.scheduleConstraints fields
// Table 119: slot-level scheduleConstraints is server Out only; absent from Table 123.
func TestRFE_SlotScheduleConstraints_Rejected(t *testing.T) {
	tests := []struct {
		name        string
		constraints *PowerTimeSlotScheduleConstraintsDataType
	}{
		{
			name: "minDuration",
			constraints: &PowerTimeSlotScheduleConstraintsDataType{
				MinDuration: util.Ptr(DurationType("PT1M")),
			},
		},
		{
			name: "maxDuration",
			constraints: &PowerTimeSlotScheduleConstraintsDataType{
				MaxDuration: util.Ptr(DurationType("PT99M")),
			},
		},
		{
			name: "optionalSlot",
			constraints: &PowerTimeSlotScheduleConstraintsDataType{
				OptionalSlot: util.Ptr(false), // existing is true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			existing := makeRFEBase()
			origMin := *existing.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ScheduleConstraints.MinDuration

			update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
				PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
					Schedule: &PowerTimeSlotScheduleDataType{
						SlotNumber: util.Ptr(PowerTimeSlotNumberType(1)),
					},
					ScheduleConstraints: tt.constraints,
				}},
			})
			result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

			assert.False(t, success, "slot-level scheduleConstraints write must be rejected")
			assert.Equal(t, existing, result)
			assert.Equal(t, origMin,
				*existing.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ScheduleConstraints.MinDuration,
				"slot scheduleConstraints must be unchanged")
		})
	}
}

// TC-26 – TC-29: powerTimeSlot.valueList
// Table 119: power values are server Out only. Table 123 contains no valueList entry at all.
func TestRFE_ValueList_Rejected(t *testing.T) {
	tests := []struct {
		name      string
		valueList *SmartEnergyManagementPsPowerTimeSlotValueListType
	}{
		{
			name: "entire_valueList",
			valueList: &SmartEnergyManagementPsPowerTimeSlotValueListType{
				Value: []PowerTimeSlotValueDataType{{
					ValueType: util.Ptr(PowerTimeSlotValueTypeTypePower),
					Value:     &ScaledNumberType{Number: util.Ptr(NumberType(9999))},
				}},
			},
		},
		{
			name: "value_Number",
			valueList: &SmartEnergyManagementPsPowerTimeSlotValueListType{
				Value: []PowerTimeSlotValueDataType{{
					ValueType: util.Ptr(PowerTimeSlotValueTypeTypePower),
					Value:     &ScaledNumberType{Number: util.Ptr(NumberType(9999))},
				}},
			},
		},
		{
			name: "value_Scale",
			valueList: &SmartEnergyManagementPsPowerTimeSlotValueListType{
				Value: []PowerTimeSlotValueDataType{{
					ValueType: util.Ptr(PowerTimeSlotValueTypeTypePower),
					Value:     &ScaledNumberType{Scale: util.Ptr(ScaleType(3))},
				}},
			},
		},
		{
			name: "adding_new_valueType",
			valueList: &SmartEnergyManagementPsPowerTimeSlotValueListType{
				Value: []PowerTimeSlotValueDataType{{
					ValueType: util.Ptr(PowerTimeSlotValueTypeTypeEnergy),
					Value:     &ScaledNumberType{Number: util.Ptr(NumberType(500))},
				}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			existing := makeRFEBase()
			origPowerValue := *existing.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value[0].Value.Number

			update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
				PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
					Schedule:  &PowerTimeSlotScheduleDataType{SlotNumber: util.Ptr(PowerTimeSlotNumberType(1))},
					ValueList: tt.valueList,
				}},
			})
			result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

			assert.False(t, success, "valueList write must be rejected")
			assert.Equal(t, existing, result)
			assert.Equal(t, origPowerValue,
				*existing.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value[0].Value.Number,
				"valueList must be unchanged")
		})
	}
}

// TC-30: Payload contains ONLY blocked fields → entire write rejected; existing data preserved.
func TestRFE_OnlyBlockedFields_Rejected(t *testing.T) {
	existing := makeRFEBase()
	origControllable := *existing.NodeScheduleInformation.NodeRemoteControllable

	update := &SmartEnergyManagementPsDataType{
		NodeScheduleInformation: &PowerSequenceNodeScheduleInformationDataType{
			NodeRemoteControllable: util.Ptr(false),
		},
	}
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.False(t, success)
	assert.Equal(t, existing, result)
	assert.Equal(t, origControllable, *existing.NodeScheduleInformation.NodeRemoteControllable)
}

// TC-31: Mixed payload: one allowed field (schedule.startTime) + one blocked field
// (nodeScheduleInformation) → entire write rejected; even the allowed field is NOT applied.
func TestRFE_MixedPayload_AllOrNothing_Rejected(t *testing.T) {
	existing := makeRFEBase()
	origStart := *existing.Alternatives[0].PowerSequence[0].Schedule.StartTime
	origControllable := *existing.NodeScheduleInformation.NodeRemoteControllable

	update := &SmartEnergyManagementPsDataType{
		NodeScheduleInformation: &PowerSequenceNodeScheduleInformationDataType{
			NodeRemoteControllable: util.Ptr(false), // blocked
		},
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(1)),
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				Description: &PowerSequenceDescriptionDataType{SequenceId: util.Ptr(PowerSequenceIdType(1))},
				Schedule: &PowerSequenceScheduleDataType{
					StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT99H")), // allowed
				},
			}},
		}},
	}
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.False(t, success, "mixed payload (blocked+allowed) must be entirely rejected")
	assert.Equal(t, existing, result)
	// The allowed field must NOT have been applied
	assert.Equal(t, origStart, *existing.Alternatives[0].PowerSequence[0].Schedule.StartTime,
		"startTime must NOT be applied when write is rejected")
	assert.Equal(t, origControllable, *existing.NodeScheduleInformation.NodeRemoteControllable,
		"nodeScheduleInformation must NOT be applied when write is rejected")
}

// TC-32: Payload contains ONLY allowed fields → no rejection (baseline).
// Confirms the rejection check is not over-broad.
func TestRFE_AllowedFieldsOnly_NoRejection(t *testing.T) {
	existing := makeRFEBase()

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		Schedule: &PowerSequenceScheduleDataType{
			StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT3H")),
		},
	})
	_, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.True(t, success, "payload with only allowed fields must not be rejected")
}

// ---------------------------------------------------------------------------
// Section 4: Remote write allowed fields (TC-33 – TC-38)
// ---------------------------------------------------------------------------

// TC-33: remoteWrite=true, write schedule.startTime → applied.
// Spec: Table 123 (Config A: startTime M)
func TestRFE_AllowedFields_ScheduleStartTime(t *testing.T) {
	existing := makeRFEBase()
	newStart := util.Ptr(AbsoluteOrRelativeTimeType("PT3H"))

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		Schedule: &PowerSequenceScheduleDataType{StartTime: newStart},
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.True(t, success)
	assert.NotNil(t, result)
	resultData := result.(*SmartEnergyManagementPsDataType)
	assert.Equal(t, *newStart, *resultData.Alternatives[0].PowerSequence[0].Schedule.StartTime,
		"startTime must be updated")
	// Other schedule field unchanged
	assert.Equal(t, AbsoluteOrRelativeTimeType("PT2H"),
		*resultData.Alternatives[0].PowerSequence[0].Schedule.EndTime,
		"endTime must remain unchanged")
}

// TC-34: remoteWrite=true, write schedule.endTime → applied.
// Spec: Table 123 (Config B: endTime M)
func TestRFE_AllowedFields_ScheduleEndTime(t *testing.T) {
	existing := makeRFEBase()
	newEnd := util.Ptr(AbsoluteOrRelativeTimeType("PT4H"))

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		Schedule: &PowerSequenceScheduleDataType{EndTime: newEnd},
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.True(t, success)
	resultData := result.(*SmartEnergyManagementPsDataType)
	assert.Equal(t, *newEnd, *resultData.Alternatives[0].PowerSequence[0].Schedule.EndTime,
		"endTime must be updated")
	assert.Equal(t, AbsoluteOrRelativeTimeType("PT1H"),
		*resultData.Alternatives[0].PowerSequence[0].Schedule.StartTime,
		"startTime must remain unchanged")
}

// TC-35: remoteWrite=true, write defaultDuration on a constrained slot
// (both minDuration and maxDuration present) → applied.
// Spec: Table 123 (Config C: defaultDuration O); §4.3.24.3.2.7
func TestRFE_AllowedFields_DefaultDuration_ConstrainedSlot(t *testing.T) {
	existing := makeRFEBase()
	newDuration := util.Ptr(DurationType("PT25M"))

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		Schedule: &PowerSequenceScheduleDataType{
			StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT1H")),
		},
		PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
			Schedule: &PowerTimeSlotScheduleDataType{
				SlotNumber:      util.Ptr(PowerTimeSlotNumberType(1)),
				DefaultDuration: newDuration,
			},
		}},
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.True(t, success)
	resultData := result.(*SmartEnergyManagementPsDataType)
	// Find slot 1
	var slot1 *SmartEnergyManagementPsPowerTimeSlotType
	for i := range resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot {
		if *resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[i].Schedule.SlotNumber == PowerTimeSlotNumberType(1) {
			slot1 = &resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[i]
			break
		}
	}
	assert.NotNil(t, slot1, "slot 1 must be found")
	if slot1 != nil {
		assert.Equal(t, *newDuration, *slot1.Schedule.DefaultDuration,
			"defaultDuration on constrained slot must be updated")
	}
}

// TC-36: remoteWrite=true, write slotActivated on a slot with optionalSlot=true → applied.
// Spec: Table 123 (Config C: slotActivated O); §4.3.24.3.2.7
func TestRFE_AllowedFields_SlotActivated_OptionalSlot(t *testing.T) {
	existing := makeRFEBase()

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		Schedule: &PowerSequenceScheduleDataType{
			StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT1H")),
		},
		PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
			Schedule: &PowerTimeSlotScheduleDataType{
				SlotNumber:    util.Ptr(PowerTimeSlotNumberType(1)),
				SlotActivated: util.Ptr(false), // was true
			},
		}},
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.True(t, success)
	resultData := result.(*SmartEnergyManagementPsDataType)
	var slot1 *SmartEnergyManagementPsPowerTimeSlotType
	for i := range resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot {
		if *resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[i].Schedule.SlotNumber == PowerTimeSlotNumberType(1) {
			slot1 = &resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[i]
			break
		}
	}
	assert.NotNil(t, slot1)
	if slot1 != nil {
		assert.Equal(t, false, *slot1.Schedule.SlotActivated,
			"slotActivated on optional slot must be updated")
	}
}

// TC-37: remoteWrite=true, write state.state → applied.
// Spec: Table 124 (state.state M)
func TestRFE_AllowedFields_StateState(t *testing.T) {
	existing := makeRFEBase()

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		State: &PowerSequenceStateDataType{
			State: util.Ptr(PowerSequenceStateTypeScheduled),
		},
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.True(t, success)
	resultData := result.(*SmartEnergyManagementPsDataType)
	assert.Equal(t, PowerSequenceStateTypeScheduled,
		*resultData.Alternatives[0].PowerSequence[0].State.State,
		"state.state must be updated")
	// Server-assigned fields must remain unchanged
	assert.Equal(t, PowerTimeSlotNumberType(5),
		*resultData.Alternatives[0].PowerSequence[0].State.ActiveSlotNumber,
		"activeSlotNumber must remain unchanged")
}

// TC-38: remoteWrite=true, Config C: schedule.startTime + powerTimeSlot.schedule → both applied.
// Spec: Table 123 (Config C: startTime M, powerTimeSlot M)
func TestRFE_ConfigC_StartTimeAndSlotData_Applied(t *testing.T) {
	existing := makeRFEBase()
	newStart := util.Ptr(AbsoluteOrRelativeTimeType("PT30M"))
	newDuration := util.Ptr(DurationType("PT25M"))

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		Schedule: &PowerSequenceScheduleDataType{StartTime: newStart},
		PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
			Schedule: &PowerTimeSlotScheduleDataType{
				SlotNumber:      util.Ptr(PowerTimeSlotNumberType(1)),
				DefaultDuration: newDuration,
			},
		}},
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.True(t, success)
	resultData := result.(*SmartEnergyManagementPsDataType)
	assert.Equal(t, *newStart,
		*resultData.Alternatives[0].PowerSequence[0].Schedule.StartTime)
	var slot1 *SmartEnergyManagementPsPowerTimeSlotType
	for i := range resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot {
		if *resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[i].Schedule.SlotNumber == PowerTimeSlotNumberType(1) {
			slot1 = &resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[i]
			break
		}
	}
	assert.NotNil(t, slot1)
	if slot1 != nil {
		assert.Equal(t, *newDuration, *slot1.Schedule.DefaultDuration)
	}
}

// ---------------------------------------------------------------------------
// Section 5: Config Option A / B / C NV (Not Valid) enforcement (TC-39 – TC-44)
// Spec: Table 123 §4.3.24.4.4.1
// ---------------------------------------------------------------------------

// TC-40: Config Option A + endTime also present → NV violation → rejected.
// Config A: endTime NV.
func TestRFE_ConfigOptionA_EndTime_IsNV_Rejected(t *testing.T) {
	existing := makeRFEBase()
	origStart := *existing.Alternatives[0].PowerSequence[0].Schedule.StartTime

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		// Both startTime and endTime present but no powerTimeSlot → ambiguous (Config A violated)
		Schedule: &PowerSequenceScheduleDataType{
			StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT30M")),
			EndTime:   util.Ptr(AbsoluteOrRelativeTimeType("PT60M")), // NV for Config A
		},
		// no powerTimeSlot → not Config C
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.False(t, success, "Config A with endTime (NV) must be rejected")
	assert.Equal(t, existing, result)
	assert.Equal(t, origStart, *existing.Alternatives[0].PowerSequence[0].Schedule.StartTime)
}

// TC-42: Config Option B + startTime also present → NV violation → rejected.
// Config B: startTime NV.
func TestRFE_ConfigOptionB_StartTime_IsNV_Rejected(t *testing.T) {
	existing := makeRFEBase()

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		Schedule: &PowerSequenceScheduleDataType{
			StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT0M")), // NV for Config B
			EndTime:   util.Ptr(AbsoluteOrRelativeTimeType("PT4H")),
		},
		// no powerTimeSlot → not Config C
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.False(t, success, "Config B with startTime (NV) must be rejected")
	assert.Equal(t, existing, result)
}

// TC-43: Config Option C: startTime + all configurable slot fields → all applied.
// Spec: Table 123 (Config C: M/O)
func TestRFE_ConfigOptionC_AllConfigurableFields_Applied(t *testing.T) {
	existing := makeRFEBase()

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		Schedule: &PowerSequenceScheduleDataType{
			StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT30M")),
		},
		PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
			Schedule: &PowerTimeSlotScheduleDataType{
				SlotNumber:      util.Ptr(PowerTimeSlotNumberType(1)),
				DefaultDuration: util.Ptr(DurationType("PT25M")),
				SlotActivated:   util.Ptr(false),
			},
		}},
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.True(t, success, "valid Config C must succeed")
	resultData := result.(*SmartEnergyManagementPsDataType)
	assert.Equal(t, AbsoluteOrRelativeTimeType("PT30M"),
		*resultData.Alternatives[0].PowerSequence[0].Schedule.StartTime)
	var slot1 *SmartEnergyManagementPsPowerTimeSlotType
	for i := range resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot {
		if *resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[i].Schedule.SlotNumber == PowerTimeSlotNumberType(1) {
			slot1 = &resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[i]
			break
		}
	}
	assert.NotNil(t, slot1)
	if slot1 != nil {
		assert.Equal(t, DurationType("PT25M"), *slot1.Schedule.DefaultDuration)
		assert.Equal(t, false, *slot1.Schedule.SlotActivated)
	}
}

// TC-44: Config Option C with endTime instead of startTime → applied.
// Spec: Table 123 (Config C: endTime O as alternative to startTime)
func TestRFE_ConfigOptionC_EndTimeAlternative_Applied(t *testing.T) {
	existing := makeRFEBase()

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		Schedule: &PowerSequenceScheduleDataType{
			EndTime: util.Ptr(AbsoluteOrRelativeTimeType("PT4H")),
		},
		PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
			Schedule: &PowerTimeSlotScheduleDataType{
				SlotNumber:      util.Ptr(PowerTimeSlotNumberType(1)),
				DefaultDuration: util.Ptr(DurationType("PT20M")),
			},
		}},
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.True(t, success, "Config C with endTime must succeed")
	resultData := result.(*SmartEnergyManagementPsDataType)
	assert.Equal(t, AbsoluteOrRelativeTimeType("PT4H"),
		*resultData.Alternatives[0].PowerSequence[0].Schedule.EndTime)
}

// ---------------------------------------------------------------------------
// Section 6: Implicit vs. constrained slot schedule (TC-45 – TC-52)
// Spec: §4.3.24.3.2.6.2, §4.3.24.3.2.7
// ---------------------------------------------------------------------------

// TC-45: Implicit slot (no ScheduleConstraints): remoteWrite=true, write defaultDuration → rejected.
// §4.3.24.3.2.6.2: implicit slots are NOT configurable.
func TestRFE_ImplicitSlot_DefaultDuration_Rejected(t *testing.T) {
	existing := makeRFEBase() // slot 2 is implicit
	origDuration := *existing.Alternatives[0].PowerSequence[0].PowerTimeSlot[1].Schedule.DefaultDuration

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		Schedule: &PowerSequenceScheduleDataType{
			StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT1H")),
		},
		PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
			Schedule: &PowerTimeSlotScheduleDataType{
				SlotNumber:      util.Ptr(PowerTimeSlotNumberType(2)), // implicit slot
				DefaultDuration: util.Ptr(DurationType("PT99M")),
			},
		}},
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.False(t, success, "defaultDuration write on implicit slot must be rejected")
	assert.Equal(t, existing, result)
	assert.Equal(t, origDuration,
		*existing.Alternatives[0].PowerSequence[0].PowerTimeSlot[1].Schedule.DefaultDuration,
		"implicit slot defaultDuration must be unchanged")
}

// TC-46: Implicit slot: remoteWrite=true, write slotActivated → rejected.
// §4.3.24.3.2.6.2: implicit slots are NOT configurable.
func TestRFE_ImplicitSlot_SlotActivated_Rejected(t *testing.T) {
	existing := makeRFEBase() // slot 2 is implicit

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		Schedule: &PowerSequenceScheduleDataType{
			StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT1H")),
		},
		PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
			Schedule: &PowerTimeSlotScheduleDataType{
				SlotNumber:    util.Ptr(PowerTimeSlotNumberType(2)), // implicit slot
				SlotActivated: util.Ptr(false),
			},
		}},
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.False(t, success, "slotActivated write on implicit slot must be rejected")
	assert.Equal(t, existing, result)
}

// TC-47: Constrained slot (both minDuration+maxDuration present): remoteWrite=true,
// write defaultDuration within range → applied.
// §4.3.24.3.2.7: duration configurable when min+max present.
func TestRFE_ConstrainedSlot_DefaultDuration_Applied(t *testing.T) {
	existing := makeRFEBase() // slot 1 is constrained (min=PT10M, max=PT30M)

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		Schedule: &PowerSequenceScheduleDataType{
			StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT1H")),
		},
		PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
			Schedule: &PowerTimeSlotScheduleDataType{
				SlotNumber:      util.Ptr(PowerTimeSlotNumberType(1)),
				DefaultDuration: util.Ptr(DurationType("PT25M")), // within [PT10M, PT30M]
			},
		}},
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.True(t, success, "defaultDuration within constraints must succeed")
	resultData := result.(*SmartEnergyManagementPsDataType)
	var slot1 *SmartEnergyManagementPsPowerTimeSlotType
	for i := range resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot {
		if *resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[i].Schedule.SlotNumber == PowerTimeSlotNumberType(1) {
			slot1 = &resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[i]
			break
		}
	}
	assert.NotNil(t, slot1)
	if slot1 != nil {
		assert.Equal(t, DurationType("PT25M"), *slot1.Schedule.DefaultDuration)
	}
}

// TC-48: Constrained slot (optionalSlot=true): remoteWrite=true, write slotActivated=false → applied.
// §4.3.24.3.2.7: activation configurable when optionalSlot=true.
func TestRFE_ConstrainedSlot_SlotActivated_Applied(t *testing.T) {
	existing := makeRFEBase() // slot 1 has OptionalSlot=true

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		Schedule: &PowerSequenceScheduleDataType{
			StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT1H")),
		},
		PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
			Schedule: &PowerTimeSlotScheduleDataType{
				SlotNumber:    util.Ptr(PowerTimeSlotNumberType(1)),
				SlotActivated: util.Ptr(false),
			},
		}},
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.True(t, success, "slotActivated on optional slot must succeed")
	resultData := result.(*SmartEnergyManagementPsDataType)
	var slot1 *SmartEnergyManagementPsPowerTimeSlotType
	for i := range resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot {
		if *resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[i].Schedule.SlotNumber == PowerTimeSlotNumberType(1) {
			slot1 = &resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[i]
			break
		}
	}
	assert.NotNil(t, slot1)
	if slot1 != nil {
		assert.Equal(t, false, *slot1.Schedule.SlotActivated)
	}
}

// TC-49: Constrained slot with only earliestStart/latestEnd (no min/max duration, no optionalSlot):
// remoteWrite=true, write defaultDuration → rejected (duration NOT configurable).
// §4.3.24.3.2.7: duration only configurable when BOTH minDuration and maxDuration present.
func TestRFE_ConstrainedSlot_NoDurationConstraints_DefaultDuration_Rejected(t *testing.T) {
	existing := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(1)),
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				Description: &PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(PowerSequenceIdType(1)),
				},
				Schedule: &PowerSequenceScheduleDataType{
					StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT1H")),
				},
				PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
					Schedule: &PowerTimeSlotScheduleDataType{
						SlotNumber:      util.Ptr(PowerTimeSlotNumberType(1)),
						DefaultDuration: util.Ptr(DurationType("PT20M")),
					},
					// Has scheduleConstraints but no min/max duration → not duration-configurable
					ScheduleConstraints: &PowerTimeSlotScheduleConstraintsDataType{
						EarliestStartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT0S")),
						LatestEndTime:     util.Ptr(AbsoluteOrRelativeTimeType("PT2H")),
						// MinDuration/MaxDuration intentionally absent
					},
				}},
			}},
		}},
	}
	origDuration := *existing.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].Schedule.DefaultDuration

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		Schedule: &PowerSequenceScheduleDataType{
			StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT1H")),
		},
		PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
			Schedule: &PowerTimeSlotScheduleDataType{
				SlotNumber:      util.Ptr(PowerTimeSlotNumberType(1)),
				DefaultDuration: util.Ptr(DurationType("PT15M")),
			},
		}},
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.False(t, success, "defaultDuration on non-duration-configurable slot must be rejected")
	assert.Equal(t, existing, result)
	assert.Equal(t, origDuration,
		*existing.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].Schedule.DefaultDuration)
}

// TC-50: Constrained slot without optionalSlot=true: remoteWrite=true, write slotActivated → rejected.
// §4.3.24.3.2.7: activation only configurable when optionalSlot=true present.
func TestRFE_ConstrainedSlot_NoOptionalSlot_SlotActivated_Rejected(t *testing.T) {
	existing := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(1)),
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				Description: &PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(PowerSequenceIdType(1)),
				},
				Schedule: &PowerSequenceScheduleDataType{
					StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT1H")),
				},
				PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
					Schedule: &PowerTimeSlotScheduleDataType{
						SlotNumber:    util.Ptr(PowerTimeSlotNumberType(1)),
						SlotActivated: util.Ptr(true),
					},
					ScheduleConstraints: &PowerTimeSlotScheduleConstraintsDataType{
						MinDuration: util.Ptr(DurationType("PT10M")),
						MaxDuration: util.Ptr(DurationType("PT30M")),
						// OptionalSlot intentionally absent → slotActivated NOT configurable
					},
				}},
			}},
		}},
	}

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		Schedule: &PowerSequenceScheduleDataType{
			StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT1H")),
		},
		PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
			Schedule: &PowerTimeSlotScheduleDataType{
				SlotNumber:    util.Ptr(PowerTimeSlotNumberType(1)),
				SlotActivated: util.Ptr(false),
			},
		}},
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.False(t, success, "slotActivated on non-optional slot must be rejected")
	assert.Equal(t, existing, result)
	assert.Equal(t, true, *existing.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].Schedule.SlotActivated)
}

// TC-51: Sequence has both implicit and constrained slots side by side.
// remoteWrite=true, write defaultDuration targeting both slots.
// Expected: constrained slot applied; implicit slot's write causes entire rejection.
// (Strict: any NV write in payload rejects the whole call.)
func TestRFE_MixedSlots_ImplicitAndConstrained_Rejected(t *testing.T) {
	existing := makeRFEBase() // slot1=constrained, slot2=implicit
	origSlot1Duration := *existing.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].Schedule.DefaultDuration
	origSlot2Duration := *existing.Alternatives[0].PowerSequence[0].PowerTimeSlot[1].Schedule.DefaultDuration

	// Both slots targeted; slot 2 (implicit) makes the whole payload invalid
	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		Schedule: &PowerSequenceScheduleDataType{
			StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT1H")),
		},
		PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{
			{
				Schedule: &PowerTimeSlotScheduleDataType{
					SlotNumber:      util.Ptr(PowerTimeSlotNumberType(1)),
					DefaultDuration: util.Ptr(DurationType("PT25M")), // constrained → valid
				},
			},
			{
				Schedule: &PowerTimeSlotScheduleDataType{
					SlotNumber:      util.Ptr(PowerTimeSlotNumberType(2)),
					DefaultDuration: util.Ptr(DurationType("PT20M")), // implicit → NV → rejects all
				},
			},
		},
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.False(t, success, "payload targeting implicit slot must be rejected entirely")
	assert.Equal(t, existing, result)
	assert.Equal(t, origSlot1Duration,
		*existing.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].Schedule.DefaultDuration,
		"constrained slot must NOT be updated (whole call rejected)")
	assert.Equal(t, origSlot2Duration,
		*existing.Alternatives[0].PowerSequence[0].PowerTimeSlot[1].Schedule.DefaultDuration,
		"implicit slot must remain unchanged")
}

// TC-52: Implicit slot + remoteWrite=false (server self-write): write defaultDuration → applied.
// Local writes have no RFE field restrictions.
func TestRFE_ImplicitSlot_LocalWrite_Applied(t *testing.T) {
	existing := makeRFEBase() // slot 2 is implicit

	// remoteWrite=false → no restrictions; server can set any field
	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(1)),
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				Description: &PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(PowerSequenceIdType(1)),
				},
				PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
					Schedule: &PowerTimeSlotScheduleDataType{
						SlotNumber:      util.Ptr(PowerTimeSlotNumberType(2)), // implicit slot
						DefaultDuration: util.Ptr(DurationType("PT99M")),
					},
				}},
			}},
		}},
	}
	result, success := existing.UpdateList(false, true, update, NewFilterTypePartial(), nil, nil)

	assert.True(t, success, "local write on implicit slot must succeed")
	resultData := result.(*SmartEnergyManagementPsDataType)
	var slot2 *SmartEnergyManagementPsPowerTimeSlotType
	for i := range resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot {
		if *resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[i].Schedule.SlotNumber == PowerTimeSlotNumberType(2) {
			slot2 = &resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[i]
			break
		}
	}
	assert.NotNil(t, slot2)
	if slot2 != nil {
		assert.Equal(t, DurationType("PT99M"), *slot2.Schedule.DefaultDuration,
			"local write must update implicit slot defaultDuration")
	}
}

// ---------------------------------------------------------------------------
// Section 7: Cross-cutting tests (TC-69, TC-73, TC-74, TC-82, TC-96, TC-99–TC-101, TC-104)
// ---------------------------------------------------------------------------

// TC-69: remoteWrite=true, sequenceId=1 selector, schedule.startTime allowed →
// only sequence 1 updated; sequence 2 unchanged.
func TestRFE_SequenceIdSelector_ScheduleStartTime(t *testing.T) {
	existing := makeRFEBase()
	origSeq2Start := *existing.Alternatives[1].PowerSequence[0].Schedule.StartTime

	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(1)),
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				Description: &PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(PowerSequenceIdType(1)),
				},
				Schedule: &PowerSequenceScheduleDataType{
					StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT30M")),
				},
			}},
		}},
	}
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.True(t, success)
	resultData := result.(*SmartEnergyManagementPsDataType)
	assert.Equal(t, AbsoluteOrRelativeTimeType("PT30M"),
		*resultData.Alternatives[0].PowerSequence[0].Schedule.StartTime,
		"sequence 1 startTime must be updated")
	assert.Equal(t, origSeq2Start,
		*resultData.Alternatives[1].PowerSequence[0].Schedule.StartTime,
		"sequence 2 startTime must be unchanged")
}

// TC-73: remoteWrite=true, slotNumber=1 selector, write defaultDuration on constrained slot →
// only slot 1 updated.
func TestRFE_SlotNumberSelector_DefaultDuration_ConstrainedSlot(t *testing.T) {
	existing := makeRFEBase() // slot 1 = constrained, slot 2 = implicit
	origSlot2Duration := *existing.Alternatives[0].PowerSequence[0].PowerTimeSlot[1].Schedule.DefaultDuration

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		Schedule: &PowerSequenceScheduleDataType{
			StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT1H")),
		},
		PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
			Schedule: &PowerTimeSlotScheduleDataType{
				SlotNumber:      util.Ptr(PowerTimeSlotNumberType(1)),
				DefaultDuration: util.Ptr(DurationType("PT25M")),
			},
		}},
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.True(t, success)
	resultData := result.(*SmartEnergyManagementPsDataType)
	var slot1, slot2 *SmartEnergyManagementPsPowerTimeSlotType
	for i := range resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot {
		n := *resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[i].Schedule.SlotNumber
		if n == 1 {
			slot1 = &resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[i]
		} else if n == 2 {
			slot2 = &resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[i]
		}
	}
	assert.NotNil(t, slot1)
	if slot1 != nil {
		assert.Equal(t, DurationType("PT25M"), *slot1.Schedule.DefaultDuration, "slot 1 must be updated")
	}
	assert.NotNil(t, slot2)
	if slot2 != nil {
		assert.Equal(t, origSlot2Duration, *slot2.Schedule.DefaultDuration, "slot 2 must be unchanged")
	}
}

// TC-74: remoteWrite=true, slotNumber=1 selector, write timePeriod.startTime → rejected.
// timePeriod is not in Table 123.
func TestRFE_SlotNumberSelector_TimePeriod_Rejected(t *testing.T) {
	existing := makeRFEBase()

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		Schedule: &PowerSequenceScheduleDataType{
			StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT1H")),
		},
		PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
			Schedule: &PowerTimeSlotScheduleDataType{
				SlotNumber: util.Ptr(PowerTimeSlotNumberType(1)),
				TimePeriod: &TimePeriodType{
					StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT99H")),
				},
			},
		}},
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.False(t, success, "timePeriod.startTime write must be rejected even with valid slotNumber")
	assert.Equal(t, existing, result)
}

// TC-82: remoteWrite=true + slotNumber selector + allowed defaultDuration + blocked valueList →
// entire write rejected (any blocked field rejects whole call).
func TestRFE_MixedSlotPayload_ValueListPlusAllowed_Rejected(t *testing.T) {
	existing := makeRFEBase()
	origDuration := *existing.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].Schedule.DefaultDuration

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		Schedule: &PowerSequenceScheduleDataType{
			StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT1H")),
		},
		PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
			Schedule: &PowerTimeSlotScheduleDataType{
				SlotNumber:      util.Ptr(PowerTimeSlotNumberType(1)),
				DefaultDuration: util.Ptr(DurationType("PT25M")), // allowed
			},
			ValueList: &SmartEnergyManagementPsPowerTimeSlotValueListType{ // blocked
				Value: []PowerTimeSlotValueDataType{{
					ValueType: util.Ptr(PowerTimeSlotValueTypeTypePower),
					Value:     &ScaledNumberType{Number: util.Ptr(NumberType(9999))},
				}},
			},
		}},
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.False(t, success, "valueList present → entire write rejected")
	assert.Equal(t, existing, result)
	assert.Equal(t, origDuration,
		*existing.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].Schedule.DefaultDuration,
		"defaultDuration must NOT have been applied (whole write rejected)")
}

// TC-96: remoteWrite=true, persist=true, allowed field (schedule.startTime) → applied; original modified.
func TestRFE_PersistTrue_AllowedField(t *testing.T) {
	existing := makeRFEBase()

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		Schedule: &PowerSequenceScheduleDataType{
			StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT3H")),
		},
	})
	_, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.True(t, success)
	// persist=true: original must be modified
	assert.Equal(t, AbsoluteOrRelativeTimeType("PT3H"),
		*existing.Alternatives[0].PowerSequence[0].Schedule.StartTime,
		"with persist=true, original must be modified in place")
}

// TC-99: remoteWrite=true, payload has ONLY slotNumber (the key field) with no data fields.
// slotNumber is M as an identifier (Table 123); key-only write is a valid no-op.
// success=true, no data changed.
func TestRFE_KeyOnly_NoDataFields_NoChange(t *testing.T) {
	existing := makeRFEBase()
	origDuration := *existing.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].Schedule.DefaultDuration

	update := rfeSeqWithKey(SmartEnergyManagementPsPowerSequenceType{
		Schedule: &PowerSequenceScheduleDataType{
			StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT1H")),
		},
		PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{{
			Schedule: &PowerTimeSlotScheduleDataType{
				SlotNumber: util.Ptr(PowerTimeSlotNumberType(1)),
				// No other fields: key-only, no data to apply
			},
		}},
	})
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	assert.True(t, success, "key-only payload (slotNumber only) is valid per Table 123")
	resultData := result.(*SmartEnergyManagementPsDataType)
	var slot1 *SmartEnergyManagementPsPowerTimeSlotType
	for i := range resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot {
		if *resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[i].Schedule.SlotNumber == PowerTimeSlotNumberType(1) {
			slot1 = &resultData.Alternatives[0].PowerSequence[0].PowerTimeSlot[i]
			break
		}
	}
	assert.NotNil(t, slot1)
	if slot1 != nil {
		assert.Equal(t, origDuration, *slot1.Schedule.DefaultDuration, "no data fields → no change")
	}
}

// TC-104: remoteWrite=true, alternativesId=1 selector + sequenceId from alternatives group 2
// → sequence not found in group 1 → no update.
func TestRFE_CrossGroup_SequenceIdMismatch_NoUpdate(t *testing.T) {
	existing := makeRFEBase()
	origSeq1Start := *existing.Alternatives[0].PowerSequence[0].Schedule.StartTime
	origSeq2Start := *existing.Alternatives[1].PowerSequence[0].Schedule.StartTime

	// Target alternativesId=1 but use sequenceId=2 (which belongs to alternativesId=2)
	update := &SmartEnergyManagementPsDataType{
		Alternatives: []SmartEnergyManagementPsAlternativesType{{
			Relation: &SmartEnergyManagementPsAlternativesRelationType{
				AlternativesId: util.Ptr(AlternativesIdType(1)), // group 1
			},
			PowerSequence: []SmartEnergyManagementPsPowerSequenceType{{
				Description: &PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(PowerSequenceIdType(2)), // lives in group 2, not group 1
				},
				Schedule: &PowerSequenceScheduleDataType{
					StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT99H")),
				},
			}},
		}},
	}
	result, success := existing.UpdateList(true, true, update, NewFilterTypePartial(), nil, nil)

	// Success=true (no error, just no match found) but no data changed
	assert.True(t, success, "no matching sequence is a success no-op, not an error")
	resultData := result.(*SmartEnergyManagementPsDataType)
	assert.Equal(t, origSeq1Start, *resultData.Alternatives[0].PowerSequence[0].Schedule.StartTime,
		"group 1 sequence must be unchanged")
	assert.Equal(t, origSeq2Start, *resultData.Alternatives[1].PowerSequence[0].Schedule.StartTime,
		"group 2 sequence must be unchanged")
}
