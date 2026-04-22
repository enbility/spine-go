package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

// makeAllFieldsData returns a SmartEnergyManagementPsDataType with every leaf
// field populated so that deepCopy can be verified exhaustively.
//
// Two Alternatives, each with two PowerSequences, each with two PowerTimeSlots,
// each with two ValueList entries — enough depth and breadth to catch any
// missing copy at any nesting level.
func makeAllFieldsData() *SmartEnergyManagementPsDataType {
	return &SmartEnergyManagementPsDataType{
		NodeScheduleInformation: &PowerSequenceNodeScheduleInformationDataType{
			NodeRemoteControllable:           util.Ptr(true),
			SupportsSingleSlotSchedulingOnly: util.Ptr(false),
			AlternativesCount:                util.Ptr(uint(2)),
			TotalSequencesCountMax:           util.Ptr(uint(4)),
			SupportsReselection:              util.Ptr(true),
		},
		Alternatives: []SmartEnergyManagementPsAlternativesType{
			{
				Relation: &SmartEnergyManagementPsAlternativesRelationType{
					AlternativesId: util.Ptr(AlternativesIdType(1)),
					SequenceId:     []PowerSequenceIdType{10, 11},
				},
				PowerSequence: []SmartEnergyManagementPsPowerSequenceType{
					makeFullSequence(10, "seq-ten", true),
					makeFullSequence(11, "seq-eleven", false),
				},
			},
			{
				Relation: &SmartEnergyManagementPsAlternativesRelationType{
					AlternativesId: util.Ptr(AlternativesIdType(2)),
					SequenceId:     []PowerSequenceIdType{20},
				},
				PowerSequence: []SmartEnergyManagementPsPowerSequenceType{
					makeFullSequence(20, "seq-twenty", true),
				},
			},
		},
	}
}

// makeFullSequence returns a SmartEnergyManagementPsPowerSequenceType with every
// field and sub-field populated.  boolFlag drives boolean fields so that
// different sequences carry different values (aiding independence checks).
func makeFullSequence(id uint, label string, boolFlag bool) SmartEnergyManagementPsPowerSequenceType {
	sid := PowerSequenceIdType(id)
	return SmartEnergyManagementPsPowerSequenceType{
		Description: &PowerSequenceDescriptionDataType{
			SequenceId:              &sid,
			Description:             util.Ptr(DescriptionType(label)),
			PositiveEnergyDirection: util.Ptr(EnergyDirectionTypeConsume),
			PowerUnit:               util.Ptr(UnitOfMeasurementTypeW),
			EnergyUnit:              util.Ptr(UnitOfMeasurementTypeWh),
			ValueSource:             util.Ptr(MeasurementValueSourceTypeMeasuredValue),
			Scope:                   util.Ptr(PowerSequenceScopeTypeForecast),
			TaskIdentifier:          util.Ptr(uint(id + 100)),
		},
		State: &PowerSequenceStateDataType{
			State:                      util.Ptr(PowerSequenceStateTypeInactive),
			ActiveSlotNumber:           util.Ptr(PowerTimeSlotNumberType(id + 1)),
			ElapsedSlotTime:            util.Ptr(DurationType("PT1M")),
			RemainingSlotTime:          util.Ptr(DurationType("PT29M")),
			SequenceRemoteControllable: util.Ptr(boolFlag),
			ActiveRepetitionNumber:     util.Ptr(uint(0)),
			RemainingPauseTime:         util.Ptr(DurationType("PT0S")),
		},
		Schedule: &PowerSequenceScheduleDataType{
			StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT1H")),
			EndTime:   util.Ptr(AbsoluteOrRelativeTimeType("PT2H")),
		},
		ScheduleConstraints: &PowerSequenceScheduleConstraintsDataType{
			EarliestStartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT0S")),
			LatestStartTime:   util.Ptr(AbsoluteOrRelativeTimeType("PT30M")),
			EarliestEndTime:   util.Ptr(AbsoluteOrRelativeTimeType("PT1H")),
			LatestEndTime:     util.Ptr(AbsoluteOrRelativeTimeType("PT3H")),
			OptionalSequence:  util.Ptr(boolFlag),
		},
		SchedulePreference: &PowerSequenceSchedulePreferenceDataType{
			Greenest: util.Ptr(boolFlag),
			Cheapest: util.Ptr(!boolFlag),
		},
		OperatingConstraintsInterrupt: &OperatingConstraintsInterruptDataType{
			IsPausable:                  util.Ptr(boolFlag),
			IsStoppable:                 util.Ptr(!boolFlag),
			NotInterruptibleAtHighPower: util.Ptr(false),
			MaxCyclesPerDay:             util.Ptr(uint(3)),
		},
		OperatingConstraintsDuration: &OperatingConstraintsDurationDataType{
			ActiveDurationMin:    util.Ptr(DurationType("PT10M")),
			ActiveDurationMax:    util.Ptr(DurationType("PT120M")),
			PauseDurationMin:     util.Ptr(DurationType("PT5M")),
			PauseDurationMax:     util.Ptr(DurationType("PT60M")),
			ActiveDurationSumMin: util.Ptr(DurationType("PT30M")),
			ActiveDurationSumMax: util.Ptr(DurationType("PT480M")),
		},
		OperatingConstraintsResumeImplication: &OperatingConstraintsResumeImplicationDataType{
			ResumeEnergyEstimated: &ScaledNumberType{
				Number: util.Ptr(NumberType(500)),
				Scale:  util.Ptr(ScaleType(0)),
			},
			EnergyUnit: util.Ptr(UnitOfMeasurementTypeWh),
			ResumeCostEstimated: &ScaledNumberType{
				Number: util.Ptr(NumberType(25)),
				Scale:  util.Ptr(ScaleType(-2)),
			},
			Currency: util.Ptr(CurrencyTypeEur),
		},
		PowerTimeSlot: []SmartEnergyManagementPsPowerTimeSlotType{
			makeFullTimeSlot(1, boolFlag),
			makeFullTimeSlot(2, !boolFlag),
		},
	}
}

// makeFullTimeSlot returns a SmartEnergyManagementPsPowerTimeSlotType with every
// field and sub-field populated.
func makeFullTimeSlot(slotNum PowerTimeSlotNumberType, boolFlag bool) SmartEnergyManagementPsPowerTimeSlotType {
	return SmartEnergyManagementPsPowerTimeSlotType{
		Schedule: &PowerTimeSlotScheduleDataType{
			SlotNumber: util.Ptr(slotNum),
			TimePeriod: &TimePeriodType{
				StartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT0S")),
				EndTime:   util.Ptr(AbsoluteOrRelativeTimeType("PT20M")),
			},
			DefaultDuration:     util.Ptr(DurationType("PT20M")),
			DurationUncertainty: util.Ptr(DurationType("PT1M")),
			SlotActivated:       util.Ptr(boolFlag),
			Description:         util.Ptr(DescriptionType("slot-desc")),
		},
		ScheduleConstraints: &PowerTimeSlotScheduleConstraintsDataType{
			EarliestStartTime: util.Ptr(AbsoluteOrRelativeTimeType("PT0S")),
			LatestEndTime:     util.Ptr(AbsoluteOrRelativeTimeType("PT2H")),
			MinDuration:       util.Ptr(DurationType("PT10M")),
			MaxDuration:       util.Ptr(DurationType("PT30M")),
			OptionalSlot:      util.Ptr(boolFlag),
		},
		ValueList: &SmartEnergyManagementPsPowerTimeSlotValueListType{
			Value: []PowerTimeSlotValueDataType{
				{
					ValueType: util.Ptr(PowerTimeSlotValueTypeTypePower),
					Value:     &ScaledNumberType{Number: util.Ptr(NumberType(3000)), Scale: util.Ptr(ScaleType(0))},
				},
				{
					ValueType: util.Ptr(PowerTimeSlotValueTypeTypePowerMax),
					Value:     &ScaledNumberType{Number: util.Ptr(NumberType(3500)), Scale: util.Ptr(ScaleType(0))},
				},
			},
		},
	}
}

// ---------------------------------------------------------------------------

// TestDeepCopy_AllFields verifies that deepCopy produces a value-identical,
// fully-independent copy of a SmartEnergyManagementPsDataType.
func TestDeepCopy_AllFields(t *testing.T) {

	// --- sub-test 1: every field value matches the original ---
	t.Run("all_fields_present_and_equal", func(t *testing.T) {
		original := makeAllFieldsData()
		copied := original.deepCopy()

		// assert.Equal performs a recursive deep-equality check.
		// A single failure here means at least one field is missing or wrong.
		assert.Equal(t, original, copied, "deepCopy must produce a value-identical copy")
	})

	// --- sub-test 2: nil receiver returns nil ---
	t.Run("nil_returns_nil", func(t *testing.T) {
		var s *SmartEnergyManagementPsDataType
		assert.Nil(t, s.deepCopy())
	})

	// --- sub-test 3: top-level pointer fields are independent ---
	t.Run("pointer_independence_top_level", func(t *testing.T) {
		original := makeAllFieldsData()
		copied := original.deepCopy()

		// Mutate NodeScheduleInformation fields in the copy.
		*copied.NodeScheduleInformation.NodeRemoteControllable = false
		*copied.NodeScheduleInformation.AlternativesCount = 99
		*copied.NodeScheduleInformation.TotalSequencesCountMax = 99
		*copied.NodeScheduleInformation.SupportsReselection = false

		assert.Equal(t, true, *original.NodeScheduleInformation.NodeRemoteControllable,
			"NodeRemoteControllable: copy mutation must not affect original")
		assert.Equal(t, uint(2), *original.NodeScheduleInformation.AlternativesCount,
			"AlternativesCount: copy mutation must not affect original")
		assert.Equal(t, uint(4), *original.NodeScheduleInformation.TotalSequencesCountMax,
			"TotalSequencesCountMax: copy mutation must not affect original")
		assert.Equal(t, true, *original.NodeScheduleInformation.SupportsReselection,
			"SupportsReselection: copy mutation must not affect original")
	})

	// --- sub-test 4: Alternatives slice is independent ---
	t.Run("slice_independence_alternatives", func(t *testing.T) {
		original := makeAllFieldsData()
		originalLen := len(original.Alternatives)

		copied := original.deepCopy()
		copied.Alternatives = append(copied.Alternatives, SmartEnergyManagementPsAlternativesType{})

		assert.Equal(t, originalLen, len(original.Alternatives),
			"appending to copied.Alternatives must not affect original")
	})

	// --- sub-test 5: Relation and its SequenceId slice are independent ---
	t.Run("pointer_independence_relation", func(t *testing.T) {
		original := makeAllFieldsData()
		copied := original.deepCopy()

		*copied.Alternatives[0].Relation.AlternativesId = AlternativesIdType(99)
		copied.Alternatives[0].Relation.SequenceId = append(
			copied.Alternatives[0].Relation.SequenceId, PowerSequenceIdType(999))

		assert.Equal(t, AlternativesIdType(1),
			*original.Alternatives[0].Relation.AlternativesId,
			"AlternativesId: copy mutation must not affect original")
		assert.Equal(t, 2, len(original.Alternatives[0].Relation.SequenceId),
			"SequenceId slice: append to copy must not affect original length")
	})

	// --- sub-test 6: Description fields are independent ---
	t.Run("pointer_independence_description", func(t *testing.T) {
		original := makeAllFieldsData()
		copied := original.deepCopy()

		seq := &copied.Alternatives[0].PowerSequence[0]
		*seq.Description.Description = DescriptionType("mutated")
		*seq.Description.PowerUnit = UnitOfMeasurementTypeV
		*seq.Description.EnergyUnit = UnitOfMeasurementTypeVAh
		*seq.Description.TaskIdentifier = 9999

		orig := &original.Alternatives[0].PowerSequence[0]
		assert.Equal(t, DescriptionType("seq-ten"), *orig.Description.Description,
			"Description.Description: copy mutation must not affect original")
		assert.Equal(t, UnitOfMeasurementTypeW, *orig.Description.PowerUnit,
			"Description.PowerUnit: copy mutation must not affect original")
		assert.Equal(t, UnitOfMeasurementTypeWh, *orig.Description.EnergyUnit,
			"Description.EnergyUnit: copy mutation must not affect original")
		assert.Equal(t, uint(110), *orig.Description.TaskIdentifier,
			"Description.TaskIdentifier: copy mutation must not affect original")
	})

	// --- sub-test 7: State fields are independent ---
	t.Run("pointer_independence_state", func(t *testing.T) {
		original := makeAllFieldsData()
		copied := original.deepCopy()

		seq := &copied.Alternatives[0].PowerSequence[0]
		*seq.State.State = PowerSequenceStateTypeRunning
		*seq.State.ActiveSlotNumber = PowerTimeSlotNumberType(99)
		*seq.State.ElapsedSlotTime = DurationType("PT99M")
		*seq.State.RemainingSlotTime = DurationType("PT99M")
		*seq.State.ActiveRepetitionNumber = 5
		*seq.State.RemainingPauseTime = DurationType("PT99M")

		orig := &original.Alternatives[0].PowerSequence[0]
		assert.Equal(t, PowerSequenceStateTypeInactive, *orig.State.State,
			"State.State: copy mutation must not affect original")
		assert.Equal(t, PowerTimeSlotNumberType(11), *orig.State.ActiveSlotNumber,
			"State.ActiveSlotNumber: copy mutation must not affect original")
		assert.Equal(t, DurationType("PT1M"), *orig.State.ElapsedSlotTime,
			"State.ElapsedSlotTime: copy mutation must not affect original")
		assert.Equal(t, DurationType("PT29M"), *orig.State.RemainingSlotTime,
			"State.RemainingSlotTime: copy mutation must not affect original")
		assert.Equal(t, uint(0), *orig.State.ActiveRepetitionNumber,
			"State.ActiveRepetitionNumber: copy mutation must not affect original")
		assert.Equal(t, DurationType("PT0S"), *orig.State.RemainingPauseTime,
			"State.RemainingPauseTime: copy mutation must not affect original")
	})

	// --- sub-test 8: Schedule fields are independent ---
	t.Run("pointer_independence_schedule", func(t *testing.T) {
		original := makeAllFieldsData()
		copied := original.deepCopy()

		seq := &copied.Alternatives[0].PowerSequence[0]
		*seq.Schedule.StartTime = AbsoluteOrRelativeTimeType("PT99H")
		*seq.Schedule.EndTime = AbsoluteOrRelativeTimeType("PT99H")

		orig := &original.Alternatives[0].PowerSequence[0]
		assert.Equal(t, AbsoluteOrRelativeTimeType("PT1H"), *orig.Schedule.StartTime,
			"Schedule.StartTime: copy mutation must not affect original")
		assert.Equal(t, AbsoluteOrRelativeTimeType("PT2H"), *orig.Schedule.EndTime,
			"Schedule.EndTime: copy mutation must not affect original")
	})

	// --- sub-test 9: ScheduleConstraints fields are independent ---
	t.Run("pointer_independence_schedule_constraints", func(t *testing.T) {
		original := makeAllFieldsData()
		copied := original.deepCopy()

		seq := &copied.Alternatives[0].PowerSequence[0]
		*seq.ScheduleConstraints.EarliestStartTime = AbsoluteOrRelativeTimeType("PT99H")
		*seq.ScheduleConstraints.LatestStartTime = AbsoluteOrRelativeTimeType("PT99H")
		*seq.ScheduleConstraints.EarliestEndTime = AbsoluteOrRelativeTimeType("PT99H")
		*seq.ScheduleConstraints.LatestEndTime = AbsoluteOrRelativeTimeType("PT99H")
		*seq.ScheduleConstraints.OptionalSequence = false

		orig := &original.Alternatives[0].PowerSequence[0]
		assert.Equal(t, AbsoluteOrRelativeTimeType("PT0S"), *orig.ScheduleConstraints.EarliestStartTime,
			"ScheduleConstraints.EarliestStartTime: copy mutation must not affect original")
		assert.Equal(t, AbsoluteOrRelativeTimeType("PT30M"), *orig.ScheduleConstraints.LatestStartTime,
			"ScheduleConstraints.LatestStartTime: copy mutation must not affect original")
		assert.Equal(t, AbsoluteOrRelativeTimeType("PT1H"), *orig.ScheduleConstraints.EarliestEndTime,
			"ScheduleConstraints.EarliestEndTime: copy mutation must not affect original")
		assert.Equal(t, AbsoluteOrRelativeTimeType("PT3H"), *orig.ScheduleConstraints.LatestEndTime,
			"ScheduleConstraints.LatestEndTime: copy mutation must not affect original")
		assert.Equal(t, true, *orig.ScheduleConstraints.OptionalSequence,
			"ScheduleConstraints.OptionalSequence: copy mutation must not affect original")
	})

	// --- sub-test 10: SchedulePreference fields are independent ---
	t.Run("pointer_independence_schedule_preference", func(t *testing.T) {
		original := makeAllFieldsData()
		copied := original.deepCopy()

		seq := &copied.Alternatives[0].PowerSequence[0]
		*seq.SchedulePreference.Greenest = false
		*seq.SchedulePreference.Cheapest = true

		orig := &original.Alternatives[0].PowerSequence[0]
		assert.Equal(t, true, *orig.SchedulePreference.Greenest,
			"SchedulePreference.Greenest: copy mutation must not affect original")
		assert.Equal(t, false, *orig.SchedulePreference.Cheapest,
			"SchedulePreference.Cheapest: copy mutation must not affect original")
	})

	// --- sub-test 11: OperatingConstraintsInterrupt fields are independent ---
	t.Run("pointer_independence_operating_constraints_interrupt", func(t *testing.T) {
		original := makeAllFieldsData()
		copied := original.deepCopy()

		seq := &copied.Alternatives[0].PowerSequence[0]
		*seq.OperatingConstraintsInterrupt.IsPausable = false
		*seq.OperatingConstraintsInterrupt.IsStoppable = true
		*seq.OperatingConstraintsInterrupt.MaxCyclesPerDay = 99

		orig := &original.Alternatives[0].PowerSequence[0]
		assert.Equal(t, true, *orig.OperatingConstraintsInterrupt.IsPausable,
			"OperatingConstraintsInterrupt.IsPausable: copy mutation must not affect original")
		assert.Equal(t, false, *orig.OperatingConstraintsInterrupt.IsStoppable,
			"OperatingConstraintsInterrupt.IsStoppable: copy mutation must not affect original")
		assert.Equal(t, uint(3), *orig.OperatingConstraintsInterrupt.MaxCyclesPerDay,
			"OperatingConstraintsInterrupt.MaxCyclesPerDay: copy mutation must not affect original")
	})

	// --- sub-test 12: OperatingConstraintsDuration fields are independent ---
	t.Run("pointer_independence_operating_constraints_duration", func(t *testing.T) {
		original := makeAllFieldsData()
		copied := original.deepCopy()

		seq := &copied.Alternatives[0].PowerSequence[0]
		*seq.OperatingConstraintsDuration.ActiveDurationMin = DurationType("PT99M")
		*seq.OperatingConstraintsDuration.ActiveDurationMax = DurationType("PT99M")
		*seq.OperatingConstraintsDuration.PauseDurationMin = DurationType("PT99M")
		*seq.OperatingConstraintsDuration.PauseDurationMax = DurationType("PT99M")
		*seq.OperatingConstraintsDuration.ActiveDurationSumMin = DurationType("PT99M")
		*seq.OperatingConstraintsDuration.ActiveDurationSumMax = DurationType("PT99M")

		orig := &original.Alternatives[0].PowerSequence[0]
		assert.Equal(t, DurationType("PT10M"), *orig.OperatingConstraintsDuration.ActiveDurationMin,
			"OperatingConstraintsDuration.ActiveDurationMin: copy mutation must not affect original")
		assert.Equal(t, DurationType("PT120M"), *orig.OperatingConstraintsDuration.ActiveDurationMax,
			"OperatingConstraintsDuration.ActiveDurationMax: copy mutation must not affect original")
		assert.Equal(t, DurationType("PT5M"), *orig.OperatingConstraintsDuration.PauseDurationMin,
			"OperatingConstraintsDuration.PauseDurationMin: copy mutation must not affect original")
		assert.Equal(t, DurationType("PT60M"), *orig.OperatingConstraintsDuration.PauseDurationMax,
			"OperatingConstraintsDuration.PauseDurationMax: copy mutation must not affect original")
		assert.Equal(t, DurationType("PT30M"), *orig.OperatingConstraintsDuration.ActiveDurationSumMin,
			"OperatingConstraintsDuration.ActiveDurationSumMin: copy mutation must not affect original")
		assert.Equal(t, DurationType("PT480M"), *orig.OperatingConstraintsDuration.ActiveDurationSumMax,
			"OperatingConstraintsDuration.ActiveDurationSumMax: copy mutation must not affect original")
	})

	// --- sub-test 13: OperatingConstraintsResumeImplication fields are independent ---
	t.Run("pointer_independence_operating_constraints_resume", func(t *testing.T) {
		original := makeAllFieldsData()
		copied := original.deepCopy()

		seq := &copied.Alternatives[0].PowerSequence[0]
		*seq.OperatingConstraintsResumeImplication.ResumeEnergyEstimated.Number = NumberType(999)
		*seq.OperatingConstraintsResumeImplication.ResumeEnergyEstimated.Scale = ScaleType(3)
		*seq.OperatingConstraintsResumeImplication.ResumeCostEstimated.Number = NumberType(999)
		*seq.OperatingConstraintsResumeImplication.EnergyUnit = UnitOfMeasurementTypeVAh
		*seq.OperatingConstraintsResumeImplication.Currency = CurrencyTypeUsd

		orig := &original.Alternatives[0].PowerSequence[0]
		assert.Equal(t, NumberType(500), *orig.OperatingConstraintsResumeImplication.ResumeEnergyEstimated.Number,
			"ResumeEnergyEstimated.Number: copy mutation must not affect original")
		assert.Equal(t, ScaleType(0), *orig.OperatingConstraintsResumeImplication.ResumeEnergyEstimated.Scale,
			"ResumeEnergyEstimated.Scale: copy mutation must not affect original")
		assert.Equal(t, NumberType(25), *orig.OperatingConstraintsResumeImplication.ResumeCostEstimated.Number,
			"ResumeCostEstimated.Number: copy mutation must not affect original")
		assert.Equal(t, UnitOfMeasurementTypeWh, *orig.OperatingConstraintsResumeImplication.EnergyUnit,
			"ResumeImplication.EnergyUnit: copy mutation must not affect original")
		assert.Equal(t, CurrencyTypeEur, *orig.OperatingConstraintsResumeImplication.Currency,
			"ResumeImplication.Currency: copy mutation must not affect original")
	})

	// --- sub-test 14: PowerTimeSlot slice is independent ---
	t.Run("slice_independence_power_time_slot", func(t *testing.T) {
		original := makeAllFieldsData()
		originalSlotLen := len(original.Alternatives[0].PowerSequence[0].PowerTimeSlot)

		copied := original.deepCopy()
		copied.Alternatives[0].PowerSequence[0].PowerTimeSlot = append(
			copied.Alternatives[0].PowerSequence[0].PowerTimeSlot,
			SmartEnergyManagementPsPowerTimeSlotType{})

		assert.Equal(t, originalSlotLen,
			len(original.Alternatives[0].PowerSequence[0].PowerTimeSlot),
			"appending to copied PowerTimeSlot slice must not affect original")
	})

	// --- sub-test 15: PowerTimeSlot.Schedule fields are independent ---
	t.Run("pointer_independence_timeslot_schedule", func(t *testing.T) {
		original := makeAllFieldsData()
		copied := original.deepCopy()

		slot := &copied.Alternatives[0].PowerSequence[0].PowerTimeSlot[0]
		*slot.Schedule.SlotNumber = PowerTimeSlotNumberType(99)
		*slot.Schedule.TimePeriod.StartTime = AbsoluteOrRelativeTimeType("PT99H")
		*slot.Schedule.TimePeriod.EndTime = AbsoluteOrRelativeTimeType("PT99H")
		*slot.Schedule.DefaultDuration = DurationType("PT99M")
		*slot.Schedule.DurationUncertainty = DurationType("PT99M")
		*slot.Schedule.SlotActivated = false
		*slot.Schedule.Description = DescriptionType("mutated")

		origSlot := &original.Alternatives[0].PowerSequence[0].PowerTimeSlot[0]
		assert.Equal(t, PowerTimeSlotNumberType(1), *origSlot.Schedule.SlotNumber,
			"SlotNumber: copy mutation must not affect original")
		assert.Equal(t, AbsoluteOrRelativeTimeType("PT0S"), *origSlot.Schedule.TimePeriod.StartTime,
			"TimePeriod.StartTime: copy mutation must not affect original")
		assert.Equal(t, AbsoluteOrRelativeTimeType("PT20M"), *origSlot.Schedule.TimePeriod.EndTime,
			"TimePeriod.EndTime: copy mutation must not affect original")
		assert.Equal(t, DurationType("PT20M"), *origSlot.Schedule.DefaultDuration,
			"DefaultDuration: copy mutation must not affect original")
		assert.Equal(t, DurationType("PT1M"), *origSlot.Schedule.DurationUncertainty,
			"DurationUncertainty: copy mutation must not affect original")
		assert.Equal(t, true, *origSlot.Schedule.SlotActivated,
			"SlotActivated: copy mutation must not affect original")
		assert.Equal(t, DescriptionType("slot-desc"), *origSlot.Schedule.Description,
			"Schedule.Description: copy mutation must not affect original")
	})

	// --- sub-test 16: PowerTimeSlot.ScheduleConstraints fields are independent ---
	t.Run("pointer_independence_timeslot_schedule_constraints", func(t *testing.T) {
		original := makeAllFieldsData()
		copied := original.deepCopy()

		slot := &copied.Alternatives[0].PowerSequence[0].PowerTimeSlot[0]
		*slot.ScheduleConstraints.EarliestStartTime = AbsoluteOrRelativeTimeType("PT99H")
		*slot.ScheduleConstraints.LatestEndTime = AbsoluteOrRelativeTimeType("PT99H")
		*slot.ScheduleConstraints.MinDuration = DurationType("PT99M")
		*slot.ScheduleConstraints.MaxDuration = DurationType("PT99M")
		*slot.ScheduleConstraints.OptionalSlot = false

		origSlot := &original.Alternatives[0].PowerSequence[0].PowerTimeSlot[0]
		assert.Equal(t, AbsoluteOrRelativeTimeType("PT0S"), *origSlot.ScheduleConstraints.EarliestStartTime,
			"SlotScheduleConstraints.EarliestStartTime: copy mutation must not affect original")
		assert.Equal(t, AbsoluteOrRelativeTimeType("PT2H"), *origSlot.ScheduleConstraints.LatestEndTime,
			"SlotScheduleConstraints.LatestEndTime: copy mutation must not affect original")
		assert.Equal(t, DurationType("PT10M"), *origSlot.ScheduleConstraints.MinDuration,
			"SlotScheduleConstraints.MinDuration: copy mutation must not affect original")
		assert.Equal(t, DurationType("PT30M"), *origSlot.ScheduleConstraints.MaxDuration,
			"SlotScheduleConstraints.MaxDuration: copy mutation must not affect original")
		assert.Equal(t, true, *origSlot.ScheduleConstraints.OptionalSlot,
			"SlotScheduleConstraints.OptionalSlot: copy mutation must not affect original")
	})

	// --- sub-test 17: ValueList entries are independent ---
	t.Run("pointer_independence_value_list", func(t *testing.T) {
		original := makeAllFieldsData()
		origPowerNumber := *original.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value[0].Value.Number
		origMaxNumber := *original.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value[1].Value.Number

		copied := original.deepCopy()
		*copied.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value[0].Value.Number = NumberType(9999)
		*copied.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value[1].Value.Number = NumberType(9999)

		assert.Equal(t, origPowerNumber,
			*original.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value[0].Value.Number,
			"ValueList[0].Value.Number: copy mutation must not affect original")
		assert.Equal(t, origMaxNumber,
			*original.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value[1].Value.Number,
			"ValueList[1].Value.Number: copy mutation must not affect original")
	})

	// --- sub-test 18: PowerSequence slice is independent ---
	t.Run("slice_independence_power_sequence", func(t *testing.T) {
		original := makeAllFieldsData()
		originalSeqLen := len(original.Alternatives[0].PowerSequence)

		copied := original.deepCopy()
		copied.Alternatives[0].PowerSequence = append(
			copied.Alternatives[0].PowerSequence,
			SmartEnergyManagementPsPowerSequenceType{})

		assert.Equal(t, originalSeqLen, len(original.Alternatives[0].PowerSequence),
			"appending to copied PowerSequence slice must not affect original")
	})
}
