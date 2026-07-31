package model

import (
	"github.com/enbility/spine-go/util"
)

// SmartEnergyManagementPsDataType

var _ Updater = (*SmartEnergyManagementPsDataType)(nil)

// UpdateList implements the Updater interface for SmartEnergyManagementPsDataType
// with proper key-based matching and SPINE-compliant update semantics.
//
// This implementation supports:
//   - Key-based matching for alternatives (alternativesId)
//   - Key-based matching for sequences (sequenceId)
//   - Composite key matching for time slot values (slotNumber + valueType)
//   - "Update all" semantics when keys are missing per SPINE spec
//   - Proper atomicity through persist flag handling
func (s *SmartEnergyManagementPsDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	// Type check
	newData, ok := newList.(*SmartEnergyManagementPsDataType)
	if !ok {
		return s, false
	}

	// Delete filters not supported for non-list types
	if filterDelete != nil {
		return s, false
	}

	// If no partial filter, do full replacement
	if filterPartial == nil {
		// If it is a remoteWrite, it shall be rejected
		if remoteWrite {
			return s, false
		}

		if persist {
			*s = *newData.deepCopy()
		}
		return newData, true
	}

	// remoteWrite=true: validate payload before any mutation (all-or-nothing per §5.3.4).
	if remoteWrite {
		if !s.validateRemoteWritePayload(newData) {
			return s, false
		}
	}

	// Explicit RFE: selectors in the filter narrow the target scope;
	// the payload contains no identifiers — just fields to apply.
	if filterData, err := filterPartial.Data(cmdFunction); err == nil && filterData.Selector != nil {
		if sel, ok := filterData.Selector.(*SmartEnergyManagementPsDataSelectorsType); ok {
			return s.applyWithSelectors(persist, newData, sel)
		}
	}

	// Implicit RFE: identifiers in the payload act as keys.
	var result *SmartEnergyManagementPsDataType
	if persist {
		result = s
	} else {
		result = s.deepCopy()
	}

	for _, newAlternative := range newData.Alternatives {
		s.mergeAlternative(result, &newAlternative)
	}

	if newData.NodeScheduleInformation != nil {
		if result.NodeScheduleInformation == nil {
			result.NodeScheduleInformation = &PowerSequenceNodeScheduleInformationDataType{}
		}
		s.mergeNodeScheduleInformation(result.NodeScheduleInformation, newData.NodeScheduleInformation)
	}

	return result, true
}

// applyWithSelectors implements explicit RFE (§5.3.4.2 with <SELECTORS>).
// Selectors narrow which entries are targeted; the payload has no identifiers.
func (s *SmartEnergyManagementPsDataType) applyWithSelectors(
	persist bool,
	newData *SmartEnergyManagementPsDataType,
	sel *SmartEnergyManagementPsDataSelectorsType,
) (any, bool) {
	var result *SmartEnergyManagementPsDataType
	if persist {
		result = s
	} else {
		result = s.deepCopy()
	}

	var seqPayload *SmartEnergyManagementPsPowerSequenceType
	var slotPayload *SmartEnergyManagementPsPowerTimeSlotType
	if len(newData.Alternatives) > 0 && len(newData.Alternatives[0].PowerSequence) > 0 {
		seqPayload = &newData.Alternatives[0].PowerSequence[0]
		if len(seqPayload.PowerTimeSlot) > 0 {
			slotPayload = &seqPayload.PowerTimeSlot[0]
		}
	}

	for i := range result.Alternatives {
		alt := &result.Alternatives[i]
		if !s.altMatchesSel(alt, sel.AlternativesRelation) {
			continue
		}
		for j := range alt.PowerSequence {
			seq := &alt.PowerSequence[j]
			if !s.seqMatchesSel(seq, sel) {
				continue
			}
			if seqPayload != nil {
				s.mergeSequenceTopLevelFields(seq, seqPayload)
			}
			if slotPayload == nil {
				continue
			}
			for k := range seq.PowerTimeSlot {
				slot := &seq.PowerTimeSlot[k]
				if !s.slotMatchesSel(slot, sel) {
					continue
				}
				s.mergeTimeSlotFieldsWithValueSel(slot, slotPayload, sel.PowerTimeSlotValue)
			}
		}
	}

	if newData.NodeScheduleInformation != nil {
		if result.NodeScheduleInformation == nil {
			result.NodeScheduleInformation = &PowerSequenceNodeScheduleInformationDataType{}
		}
		s.mergeNodeScheduleInformation(result.NodeScheduleInformation, newData.NodeScheduleInformation)
	}

	return result, true
}

// altMatchesSel returns true if alt matches the AlternativesRelation selector.
func (s *SmartEnergyManagementPsDataType) altMatchesSel(
	alt *SmartEnergyManagementPsAlternativesType,
	sel *PowerSequenceAlternativesRelationListDataSelectorsType,
) bool {
	if sel == nil {
		return true
	}
	if alt.Relation == nil {
		return false
	}
	rel := (*PowerSequenceAlternativesRelationDataType)(alt.Relation)
	if sel.AlternativesId != nil {
		if rel.AlternativesId == nil || *rel.AlternativesId != *sel.AlternativesId {
			return false
		}
	}
	return true
}

// seqMatchesSel returns true if seq matches sequence-level selectors.
// PowerSequenceDescription, PowerTimeSlotSchedule.SequenceId, and
// PowerTimeSlotValue.SequenceId all contribute AND conditions.
func (s *SmartEnergyManagementPsDataType) seqMatchesSel(
	seq *SmartEnergyManagementPsPowerSequenceType,
	sel *SmartEnergyManagementPsDataSelectorsType,
) bool {
	var seqId *PowerSequenceIdType
	if seq.Description != nil {
		seqId = seq.Description.SequenceId
	}
	if sel.PowerSequenceDescription != nil && sel.PowerSequenceDescription.SequenceId != nil {
		if seqId == nil || *seqId != *sel.PowerSequenceDescription.SequenceId {
			return false
		}
	}
	if sel.PowerTimeSlotSchedule != nil && sel.PowerTimeSlotSchedule.SequenceId != nil {
		if seqId == nil || *seqId != *sel.PowerTimeSlotSchedule.SequenceId {
			return false
		}
	}
	if sel.PowerTimeSlotValue != nil && sel.PowerTimeSlotValue.SequenceId != nil {
		if seqId == nil || *seqId != *sel.PowerTimeSlotValue.SequenceId {
			return false
		}
	}
	return true
}

// slotMatchesSel returns true if slot matches slot-level selectors.
func (s *SmartEnergyManagementPsDataType) slotMatchesSel(
	slot *SmartEnergyManagementPsPowerTimeSlotType,
	sel *SmartEnergyManagementPsDataSelectorsType,
) bool {
	var slotNum *PowerTimeSlotNumberType
	if slot.Schedule != nil {
		slotNum = slot.Schedule.SlotNumber
	}
	if sel.PowerTimeSlotSchedule != nil && sel.PowerTimeSlotSchedule.SlotNumber != nil {
		if slotNum == nil || *slotNum != *sel.PowerTimeSlotSchedule.SlotNumber {
			return false
		}
	}
	if sel.PowerTimeSlotValue != nil && sel.PowerTimeSlotValue.SlotNumber != nil {
		if slotNum == nil || *slotNum != *sel.PowerTimeSlotValue.SlotNumber {
			return false
		}
	}
	return true
}

// Key-based matching helper functions

// findAlternativeByKey finds an alternative in the result by alternativesId
func (s *SmartEnergyManagementPsDataType) findAlternativeByKey(result *SmartEnergyManagementPsDataType, alternativesId *AlternativesIdType) *SmartEnergyManagementPsAlternativesType {
	if alternativesId == nil {
		return nil
	}

	for i := range result.Alternatives {
		if result.Alternatives[i].Relation != nil &&
			result.Alternatives[i].Relation.AlternativesId != nil &&
			*result.Alternatives[i].Relation.AlternativesId == *alternativesId {
			return &result.Alternatives[i]
		}
	}
	return nil
}

// findSequenceByKey finds a power sequence within an alternative by sequenceId
func (s *SmartEnergyManagementPsDataType) findSequenceByKey(alternative *SmartEnergyManagementPsAlternativesType, sequenceId *PowerSequenceIdType) *SmartEnergyManagementPsPowerSequenceType {
	if sequenceId == nil {
		return nil
	}

	for i := range alternative.PowerSequence {
		if alternative.PowerSequence[i].Description != nil &&
			alternative.PowerSequence[i].Description.SequenceId != nil &&
			*alternative.PowerSequence[i].Description.SequenceId == *sequenceId {
			return &alternative.PowerSequence[i]
		}
	}
	return nil
}

// findTimeSlotValueByKey finds a time slot value by composite key (slotNumber + valueType)
func (s *SmartEnergyManagementPsDataType) findTimeSlotValueByKey(timeSlot *SmartEnergyManagementPsPowerTimeSlotType,
	slotNumber *PowerTimeSlotNumberType, valueType *PowerTimeSlotValueTypeType) *PowerTimeSlotValueDataType {

	if timeSlot.ValueList == nil || slotNumber == nil || valueType == nil {
		return nil
	}

	for i := range timeSlot.ValueList.Value {
		value := &timeSlot.ValueList.Value[i]
		if value.ValueType != nil && *value.ValueType == *valueType {
			return value
		}
	}
	return nil
}

// findTimeSlotByKey finds a time slot by slotNumber
func (s *SmartEnergyManagementPsDataType) findTimeSlotByKey(sequence *SmartEnergyManagementPsPowerSequenceType, slotNumber *PowerTimeSlotNumberType) *SmartEnergyManagementPsPowerTimeSlotType {
	if slotNumber == nil {
		return nil
	}

	for i := range sequence.PowerTimeSlot {
		if sequence.PowerTimeSlot[i].Schedule != nil &&
			sequence.PowerTimeSlot[i].Schedule.SlotNumber != nil &&
			*sequence.PowerTimeSlot[i].Schedule.SlotNumber == *slotNumber {
			return &sequence.PowerTimeSlot[i]
		}
	}
	return nil
}

// Merge helper functions

// mergeAlternative merges a new alternative into the result using key-based matching
func (s *SmartEnergyManagementPsDataType) mergeAlternative(result *SmartEnergyManagementPsDataType, newAlternative *SmartEnergyManagementPsAlternativesType) {
	// Get alternativesId for key-based matching
	var alternativesId *AlternativesIdType
	if newAlternative.Relation != nil {
		alternativesId = newAlternative.Relation.AlternativesId
	}

	// Handle missing key - apply "update all" semantics per SPINE spec
	if alternativesId == nil {
		// Check if this looks like a positional update attempt
		if s.looksLikePositionalUpdate(newAlternative) {
			// For positional-style updates with missing keys, only update first alternative
			// to maintain backward compatibility
			if len(result.Alternatives) > 0 {
				s.mergeAlternativeFields(&result.Alternatives[0], newAlternative)
			}
		} else {
			// Apply update to ALL alternatives (true "update all" semantics)
			for i := range result.Alternatives {
				s.mergeAlternativeFields(&result.Alternatives[i], newAlternative)
			}
		}
		return
	}

	// Key-based matching: find target alternative
	targetAlternative := s.findAlternativeByKey(result, alternativesId)
	if targetAlternative == nil {
		// Alternative not found - could create new one, but for OHPCF we expect it to exist
		return
	}

	// Merge fields from new alternative to target
	s.mergeAlternativeFields(targetAlternative, newAlternative)
}

// mergeAlternativeFields merges fields from newAlternative to target
func (s *SmartEnergyManagementPsDataType) mergeAlternativeFields(target, newAlternative *SmartEnergyManagementPsAlternativesType) {
	// Handle positional vs key-based sequence merging
	if s.looksLikePositionalUpdate(newAlternative) {
		// Positional-style update: only process sequences with content at valid positions
		for i, newSequence := range newAlternative.PowerSequence {
			if s.hasSequenceContent(&newSequence) {
				// Strict positional matching: only update if position exists
				if i < len(target.PowerSequence) {
					// Use existing sequence at this position
					s.mergeSequenceFields(&target.PowerSequence[i], &newSequence)
				}
				// If position doesn't exist, ignore this sequence (don't fall back to key-based)
			}
		}
	} else {
		// Semantic update: use key-based matching for all sequences
		for _, newSequence := range newAlternative.PowerSequence {
			s.mergePowerSequence(target, &newSequence)
		}
	}
}

// mergePowerSequence merges a power sequence using key-based matching
func (s *SmartEnergyManagementPsDataType) mergePowerSequence(alternative *SmartEnergyManagementPsAlternativesType, newSequence *SmartEnergyManagementPsPowerSequenceType) {
	// Get sequenceId for key-based matching
	var sequenceId *PowerSequenceIdType
	if newSequence.Description != nil {
		sequenceId = newSequence.Description.SequenceId
	}

	// Handle missing key
	if sequenceId == nil {
		// Check if this sequence has any meaningful content
		if s.hasSequenceContent(newSequence) {
			// Apply update to ALL sequences in this alternative (true "update all" semantics)
			for i := range alternative.PowerSequence {
				s.mergeSequenceFields(&alternative.PowerSequence[i], newSequence)
			}
		}
		// If no meaningful content, ignore this sequence (positional placeholder)
		return
	}

	// Key-based matching: find target sequence
	targetSequence := s.findSequenceByKey(alternative, sequenceId)
	if targetSequence == nil {
		// Sequence not found - could create new one, but for OHPCF we expect it to exist
		return
	}

	// Merge fields from new sequence to target
	s.mergeSequenceFields(targetSequence, newSequence)
}

// mergeSequenceFields merges fields from newSequence to target
func (s *SmartEnergyManagementPsDataType) mergeSequenceFields(target, newSequence *SmartEnergyManagementPsPowerSequenceType) {
	s.mergeSequenceTopLevelFields(target, newSequence)
	for _, newTimeSlot := range newSequence.PowerTimeSlot {
		s.mergePowerTimeSlot(target, &newTimeSlot)
	}
}

// mergeSequenceTopLevelFields merges only non-slot fields (used by explicit RFE path).
func (s *SmartEnergyManagementPsDataType) mergeSequenceTopLevelFields(target, newSequence *SmartEnergyManagementPsPowerSequenceType) {
	if newSequence.State != nil {
		if target.State == nil {
			target.State = &PowerSequenceStateDataType{}
		}
		s.mergeStateFields(target.State, newSequence.State)
	}
	if newSequence.Schedule != nil {
		if target.Schedule == nil {
			target.Schedule = &PowerSequenceScheduleDataType{}
		}
		s.mergeScheduleFields(target.Schedule, newSequence.Schedule)
	}
}

// mergeStateFields merges non-nil fields from newState to target
func (s *SmartEnergyManagementPsDataType) mergeStateFields(target, newState *PowerSequenceStateDataType) {
	if newState.State != nil {
		target.State = newState.State
	}
	if newState.ActiveSlotNumber != nil {
		target.ActiveSlotNumber = newState.ActiveSlotNumber
	}
	if newState.ElapsedSlotTime != nil {
		target.ElapsedSlotTime = newState.ElapsedSlotTime
	}
	if newState.RemainingSlotTime != nil {
		target.RemainingSlotTime = newState.RemainingSlotTime
	}
	if newState.SequenceRemoteControllable != nil {
		target.SequenceRemoteControllable = newState.SequenceRemoteControllable
	}
	if newState.ActiveRepetitionNumber != nil {
		target.ActiveRepetitionNumber = newState.ActiveRepetitionNumber
	}
	if newState.RemainingPauseTime != nil {
		target.RemainingPauseTime = newState.RemainingPauseTime
	}
}

// mergeScheduleFields merges non-nil fields from newSchedule to target
func (s *SmartEnergyManagementPsDataType) mergeScheduleFields(target, newSchedule *PowerSequenceScheduleDataType) {
	if newSchedule.StartTime != nil {
		target.StartTime = newSchedule.StartTime
	}
	if newSchedule.EndTime != nil {
		target.EndTime = newSchedule.EndTime
	}
}

// mergePowerTimeSlot merges a power time slot using composite key matching
func (s *SmartEnergyManagementPsDataType) mergePowerTimeSlot(sequence *SmartEnergyManagementPsPowerSequenceType, newTimeSlot *SmartEnergyManagementPsPowerTimeSlotType) {
	// Get slotNumber for key-based matching
	var slotNumber *PowerTimeSlotNumberType
	if newTimeSlot.Schedule != nil {
		slotNumber = newTimeSlot.Schedule.SlotNumber
	}

	// Handle missing key - apply "update all" semantics per SPINE spec
	if slotNumber == nil {
		// Apply update to ALL time slots in this sequence
		for i := range sequence.PowerTimeSlot {
			s.mergeTimeSlotFields(&sequence.PowerTimeSlot[i], newTimeSlot)
		}
		return
	}

	// Key-based matching: find target time slot
	targetTimeSlot := s.findTimeSlotByKey(sequence, slotNumber)
	if targetTimeSlot == nil {
		// Time slot not found - could create new one, but for OHPCF we expect it to exist
		return
	}

	// Merge fields from new time slot to target
	s.mergeTimeSlotFields(targetTimeSlot, newTimeSlot)
}

// mergeTimeSlotFields merges fields from newTimeSlot to target.
func (s *SmartEnergyManagementPsDataType) mergeTimeSlotFields(target, newTimeSlot *SmartEnergyManagementPsPowerTimeSlotType) {
	s.mergeTimeSlotFieldsWithValueSel(target, newTimeSlot, nil)
}

// mergeTimeSlotFieldsWithValueSel merges slot fields; when valueSel is set,
// only values whose ValueType matches the selector are updated.
func (s *SmartEnergyManagementPsDataType) mergeTimeSlotFieldsWithValueSel(
	target, newTimeSlot *SmartEnergyManagementPsPowerTimeSlotType,
	valueSel *PowerTimeSlotValueListDataSelectorsType,
) {
	if newTimeSlot.Schedule != nil {
		if target.Schedule == nil {
			target.Schedule = &PowerTimeSlotScheduleDataType{}
		}
		s.mergeSlotScheduleFields(target.Schedule, newTimeSlot.Schedule)
	}
	if newTimeSlot.ScheduleConstraints != nil {
		if target.ScheduleConstraints == nil {
			target.ScheduleConstraints = &PowerTimeSlotScheduleConstraintsDataType{}
		}
		s.mergeSlotScheduleConstraintsFields(target.ScheduleConstraints, newTimeSlot.ScheduleConstraints)
	}
	if newTimeSlot.ValueList == nil {
		return
	}
	if target.ValueList == nil {
		target.ValueList = &SmartEnergyManagementPsPowerTimeSlotValueListType{}
	}
	for i := range newTimeSlot.ValueList.Value {
		nv := &newTimeSlot.ValueList.Value[i]
		if valueSel != nil && valueSel.ValueType != nil {
			for v := range target.ValueList.Value {
				if target.ValueList.Value[v].ValueType != nil &&
					*target.ValueList.Value[v].ValueType == *valueSel.ValueType {
					s.mergeValueFields(&target.ValueList.Value[v], nv)
				}
			}
		} else {
			s.mergeTimeSlotValue(target, nv)
		}
	}
}

// mergeSlotScheduleFields merges non-nil, non-key fields from newSched to target.
func (s *SmartEnergyManagementPsDataType) mergeSlotScheduleFields(target, newSched *PowerTimeSlotScheduleDataType) {
	if newSched.TimePeriod != nil {
		if target.TimePeriod == nil {
			target.TimePeriod = &TimePeriodType{}
		}
		if newSched.TimePeriod.StartTime != nil {
			target.TimePeriod.StartTime = newSched.TimePeriod.StartTime
		}
		if newSched.TimePeriod.EndTime != nil {
			target.TimePeriod.EndTime = newSched.TimePeriod.EndTime
		}
	}
	if newSched.DefaultDuration != nil {
		target.DefaultDuration = newSched.DefaultDuration
	}
	if newSched.DurationUncertainty != nil {
		target.DurationUncertainty = newSched.DurationUncertainty
	}
	if newSched.SlotActivated != nil {
		target.SlotActivated = newSched.SlotActivated
	}
	if newSched.Description != nil {
		target.Description = newSched.Description
	}
}

// mergeSlotScheduleConstraintsFields merges non-nil, non-key fields.
func (s *SmartEnergyManagementPsDataType) mergeSlotScheduleConstraintsFields(target, newConstr *PowerTimeSlotScheduleConstraintsDataType) {
	if newConstr.EarliestStartTime != nil {
		target.EarliestStartTime = newConstr.EarliestStartTime
	}
	if newConstr.LatestEndTime != nil {
		target.LatestEndTime = newConstr.LatestEndTime
	}
	if newConstr.MinDuration != nil {
		target.MinDuration = newConstr.MinDuration
	}
	if newConstr.MaxDuration != nil {
		target.MaxDuration = newConstr.MaxDuration
	}
	if newConstr.OptionalSlot != nil {
		target.OptionalSlot = newConstr.OptionalSlot
	}
}

// mergeTimeSlotValue merges a time slot value using composite key matching (valueType)
func (s *SmartEnergyManagementPsDataType) mergeTimeSlotValue(timeSlot *SmartEnergyManagementPsPowerTimeSlotType, newValue *PowerTimeSlotValueDataType) {
	// Get slotNumber from timeSlot and valueType for composite key matching
	var slotNumber *PowerTimeSlotNumberType
	if timeSlot.Schedule != nil {
		slotNumber = timeSlot.Schedule.SlotNumber
	}

	var valueType *PowerTimeSlotValueTypeType = newValue.ValueType

	// Handle missing keys - apply "update all" semantics per SPINE spec
	if valueType == nil {
		// Apply update to ALL values in this time slot
		for i := range timeSlot.ValueList.Value {
			s.mergeValueFields(&timeSlot.ValueList.Value[i], newValue)
		}
		return
	}

	// Composite key matching: find target value by valueType
	targetValue := s.findTimeSlotValueByKey(timeSlot, slotNumber, valueType)
	if targetValue == nil {
		// Value not found - could create new one, but for OHPCF we expect it to exist
		return
	}

	// Merge fields from new value to target
	s.mergeValueFields(targetValue, newValue)
}

// mergeValueFields merges non-nil fields from newValue to target
func (s *SmartEnergyManagementPsDataType) mergeValueFields(target, newValue *PowerTimeSlotValueDataType) {
	if newValue.Value != nil {
		if target.Value == nil {
			target.Value = &ScaledNumberType{}
		}
		if newValue.Value.Number != nil {
			target.Value.Number = newValue.Value.Number
		}
		if newValue.Value.Scale != nil {
			target.Value.Scale = newValue.Value.Scale
		}
	}
}

// mergeNodeScheduleInformation merges non-nil fields from newInfo to target
func (s *SmartEnergyManagementPsDataType) mergeNodeScheduleInformation(target, newInfo *PowerSequenceNodeScheduleInformationDataType) {
	if newInfo.NodeRemoteControllable != nil {
		target.NodeRemoteControllable = newInfo.NodeRemoteControllable
	}
	if newInfo.SupportsSingleSlotSchedulingOnly != nil {
		target.SupportsSingleSlotSchedulingOnly = newInfo.SupportsSingleSlotSchedulingOnly
	}
	if newInfo.AlternativesCount != nil {
		target.AlternativesCount = newInfo.AlternativesCount
	}
	if newInfo.TotalSequencesCountMax != nil {
		target.TotalSequencesCountMax = newInfo.TotalSequencesCountMax
	}
	if newInfo.SupportsReselection != nil {
		target.SupportsReselection = newInfo.SupportsReselection
	}
}

// Helper functions to detect positional vs semantic updates

// looksLikePositionalUpdate detects if an alternative update is attempting positional access
func (s *SmartEnergyManagementPsDataType) looksLikePositionalUpdate(alternative *SmartEnergyManagementPsAlternativesType) bool {
	// If alternative has mixed empty/filled sequences, it's likely positional
	hasEmpty := false
	hasFilled := false

	for _, seq := range alternative.PowerSequence {
		if s.hasSequenceContent(&seq) {
			hasFilled = true
		} else {
			hasEmpty = true
		}
	}

	// Mixed content suggests positional array access
	return hasEmpty && hasFilled
}

// hasSequenceContent checks if a sequence contains meaningful update data
func (s *SmartEnergyManagementPsDataType) hasSequenceContent(sequence *SmartEnergyManagementPsPowerSequenceType) bool {
	return sequence.State != nil ||
		sequence.Schedule != nil ||
		len(sequence.PowerTimeSlot) > 0 ||
		sequence.ScheduleConstraints != nil ||
		sequence.SchedulePreference != nil ||
		sequence.OperatingConstraintsInterrupt != nil ||
		sequence.OperatingConstraintsDuration != nil ||
		sequence.OperatingConstraintsResumeImplication != nil
}

// validateRemoteWritePayload returns true when the payload contains only fields
// that are writable by a remote client per Tables 123/124 §4.3.24.4.4-4.5.
// Any blocked field causes the whole write to be rejected (all-or-nothing).
func (s *SmartEnergyManagementPsDataType) validateRemoteWritePayload(newData *SmartEnergyManagementPsDataType) bool {
	// NodeScheduleInformation is server Out only (Table 119)
	if newData.NodeScheduleInformation != nil {
		return false
	}
	for _, newAlt := range newData.Alternatives {
		for _, newSeq := range newAlt.PowerSequence {
			// Description: only SequenceId is a key; label and powerUnit are server-owned
			if newSeq.Description != nil {
				d := newSeq.Description
				if d.Description != nil || d.PowerUnit != nil {
					return false
				}
			}
			// State: only State.State is permitted (Table 124); all other sub-fields are server-assigned
			if newSeq.State != nil {
				st := newSeq.State
				if st.ActiveSlotNumber != nil || st.ElapsedSlotTime != nil ||
					st.RemainingSlotTime != nil || st.SequenceRemoteControllable != nil ||
					st.ActiveRepetitionNumber != nil || st.RemainingPauseTime != nil {
					return false
				}
			}
			// Sequence-level ScheduleConstraints, SchedulePreference, OperatingConstraints: server Out only
			if newSeq.ScheduleConstraints != nil || newSeq.SchedulePreference != nil ||
				newSeq.OperatingConstraintsInterrupt != nil || newSeq.OperatingConstraintsDuration != nil ||
				newSeq.OperatingConstraintsResumeImplication != nil {
				return false
			}
			// Config Option NV: startTime AND endTime both present with no slot data (Table 123)
			if newSeq.Schedule != nil {
				sched := newSeq.Schedule
				if sched.StartTime != nil && sched.EndTime != nil && len(newSeq.PowerTimeSlot) == 0 {
					return false
				}
			}
			// Slot-level checks
			for _, newSlot := range newSeq.PowerTimeSlot {
				// slot.schedule.timePeriod: not in Table 123
				if newSlot.Schedule != nil && newSlot.Schedule.TimePeriod != nil {
					return false
				}
				// slot.schedule.durationUncertainty: not in Table 123
				if newSlot.Schedule != nil && newSlot.Schedule.DurationUncertainty != nil {
					return false
				}
				// Slot-level ScheduleConstraints and ValueList: server Out only
				if newSlot.ScheduleConstraints != nil || newSlot.ValueList != nil {
					return false
				}
				// Configurability checks: implicit slots and slot-specific constraints
				hasDefaultDuration := newSlot.Schedule != nil && newSlot.Schedule.DefaultDuration != nil
				hasSlotActivated := newSlot.Schedule != nil && newSlot.Schedule.SlotActivated != nil
				if !hasDefaultDuration && !hasSlotActivated {
					continue
				}
				existingSlot := s.findSlotForValidation(&newAlt, &newSeq, newSlot.Schedule)
				if existingSlot == nil {
					continue // slot not found; allow (no constraints to check)
				}
				if existingSlot.ScheduleConstraints == nil {
					// Implicit slot — not configurable at all
					return false
				}
				constr := existingSlot.ScheduleConstraints
				if hasDefaultDuration {
					if constr.MinDuration == nil || constr.MaxDuration == nil {
						return false // slot not duration-configurable
					}
				}
				if hasSlotActivated {
					if constr.OptionalSlot == nil || !*constr.OptionalSlot {
						return false // slot not activation-configurable
					}
				}
			}
		}
	}
	return true
}

// findSlotForValidation resolves the existing slot that corresponds to a slot entry in
// the incoming payload so configurability constraints can be checked.
func (s *SmartEnergyManagementPsDataType) findSlotForValidation(
	newAlt *SmartEnergyManagementPsAlternativesType,
	newSeq *SmartEnergyManagementPsPowerSequenceType,
	newSlotSched *PowerTimeSlotScheduleDataType,
) *SmartEnergyManagementPsPowerTimeSlotType {
	if newSlotSched == nil || newSlotSched.SlotNumber == nil {
		return nil
	}
	var altId *AlternativesIdType
	if newAlt.Relation != nil {
		altId = newAlt.Relation.AlternativesId
	}
	existingAlt := s.findAlternativeByKey(s, altId)
	if existingAlt == nil {
		return nil
	}
	var seqId *PowerSequenceIdType
	if newSeq.Description != nil {
		seqId = newSeq.Description.SequenceId
	}
	existingSeq := s.findSequenceByKey(existingAlt, seqId)
	if existingSeq == nil {
		return nil
	}
	return s.findTimeSlotByKey(existingSeq, newSlotSched.SlotNumber)
}

// deepCopy returns a fully independent copy of s via JSON round-trip.
// This guarantees every field is copied regardless of future struct changes.
func (s *SmartEnergyManagementPsDataType) deepCopy() *SmartEnergyManagementPsDataType {
	if s == nil {
		return nil
	}
	result := &SmartEnergyManagementPsDataType{}
	util.DeepCopy(s, result)
	return result
}
