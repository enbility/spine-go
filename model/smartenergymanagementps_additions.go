package model

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
func (s *SmartEnergyManagementPsDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
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
		if persist {
			*s = *newData
		}
		return newData, true
	}

	// Create working copy for updates
	var result *SmartEnergyManagementPsDataType
	if persist {
		// Work directly on original data
		result = s
	} else {
		// Work on a deep copy to preserve original
		result = s.deepCopy()
	}

	// Merge alternatives using key-based matching
	for _, newAlternative := range newData.Alternatives {
		s.mergeAlternative(result, &newAlternative)
	}

	// Update top-level NodeScheduleInformation if provided
	if newData.NodeScheduleInformation != nil {
		if result.NodeScheduleInformation == nil {
			result.NodeScheduleInformation = &PowerSequenceNodeScheduleInformationDataType{}
		}
		s.mergeNodeScheduleInformation(result.NodeScheduleInformation, newData.NodeScheduleInformation)
	}

	return result, true
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
	// Merge state if provided
	if newSequence.State != nil {
		if target.State == nil {
			target.State = &PowerSequenceStateDataType{}
		}
		s.mergeStateFields(target.State, newSequence.State)
	}
	
	// Merge schedule if provided
	if newSequence.Schedule != nil {
		if target.Schedule == nil {
			target.Schedule = &PowerSequenceScheduleDataType{}
		}
		s.mergeScheduleFields(target.Schedule, newSequence.Schedule)
	}
	
	// Merge power time slots using composite key matching
	for _, newTimeSlot := range newSequence.PowerTimeSlot {
		s.mergePowerTimeSlot(target, &newTimeSlot)
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

// mergeTimeSlotFields merges fields from newTimeSlot to target
func (s *SmartEnergyManagementPsDataType) mergeTimeSlotFields(target, newTimeSlot *SmartEnergyManagementPsPowerTimeSlotType) {
	// Merge value list using composite key matching (slotNumber + valueType)
	if newTimeSlot.ValueList != nil {
		if target.ValueList == nil {
			target.ValueList = &SmartEnergyManagementPsPowerTimeSlotValueListType{}
		}
		
		for _, newValue := range newTimeSlot.ValueList.Value {
			s.mergeTimeSlotValue(target, &newValue)
		}
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

// deepCopy creates a deep copy of the SmartEnergyManagementPsDataType
func (s *SmartEnergyManagementPsDataType) deepCopy() *SmartEnergyManagementPsDataType {
	if s == nil {
		return nil
	}

	// Create a new instance
	result := &SmartEnergyManagementPsDataType{}

	// Copy NodeScheduleInformation
	if s.NodeScheduleInformation != nil {
		result.NodeScheduleInformation = &PowerSequenceNodeScheduleInformationDataType{}
		if s.NodeScheduleInformation.NodeRemoteControllable != nil {
			v := *s.NodeScheduleInformation.NodeRemoteControllable
			result.NodeScheduleInformation.NodeRemoteControllable = &v
		}
		if s.NodeScheduleInformation.SupportsSingleSlotSchedulingOnly != nil {
			v := *s.NodeScheduleInformation.SupportsSingleSlotSchedulingOnly
			result.NodeScheduleInformation.SupportsSingleSlotSchedulingOnly = &v
		}
		if s.NodeScheduleInformation.AlternativesCount != nil {
			v := *s.NodeScheduleInformation.AlternativesCount
			result.NodeScheduleInformation.AlternativesCount = &v
		}
		if s.NodeScheduleInformation.TotalSequencesCountMax != nil {
			v := *s.NodeScheduleInformation.TotalSequencesCountMax
			result.NodeScheduleInformation.TotalSequencesCountMax = &v
		}
		if s.NodeScheduleInformation.SupportsReselection != nil {
			v := *s.NodeScheduleInformation.SupportsReselection
			result.NodeScheduleInformation.SupportsReselection = &v
		}
	}

	// Copy Alternatives
	if len(s.Alternatives) > 0 {
		result.Alternatives = make([]SmartEnergyManagementPsAlternativesType, len(s.Alternatives))
		for i := range s.Alternatives {
			// Copy Relation
			if s.Alternatives[i].Relation != nil {
				result.Alternatives[i].Relation = &SmartEnergyManagementPsAlternativesRelationType{}
				if s.Alternatives[i].Relation.AlternativesId != nil {
					v := *s.Alternatives[i].Relation.AlternativesId
					result.Alternatives[i].Relation.AlternativesId = &v
				}
			}

			// Copy PowerSequence
			if len(s.Alternatives[i].PowerSequence) > 0 {
				result.Alternatives[i].PowerSequence = make([]SmartEnergyManagementPsPowerSequenceType, len(s.Alternatives[i].PowerSequence))
				for j := range s.Alternatives[i].PowerSequence {
					ps := &s.Alternatives[i].PowerSequence[j]
					rps := &result.Alternatives[i].PowerSequence[j]

					// Copy Description
					if ps.Description != nil {
						rps.Description = &PowerSequenceDescriptionDataType{}
						if ps.Description.SequenceId != nil {
							v := *ps.Description.SequenceId
							rps.Description.SequenceId = &v
						}
						if ps.Description.PowerUnit != nil {
							v := *ps.Description.PowerUnit
							rps.Description.PowerUnit = &v
						}
					}

					// Copy State
					if ps.State != nil {
						rps.State = &PowerSequenceStateDataType{}
						if ps.State.State != nil {
							v := *ps.State.State
							rps.State.State = &v
						}
						if ps.State.ActiveSlotNumber != nil {
							v := *ps.State.ActiveSlotNumber
							rps.State.ActiveSlotNumber = &v
						}
						if ps.State.ElapsedSlotTime != nil {
							v := *ps.State.ElapsedSlotTime
							rps.State.ElapsedSlotTime = &v
						}
						if ps.State.RemainingSlotTime != nil {
							v := *ps.State.RemainingSlotTime
							rps.State.RemainingSlotTime = &v
						}
						if ps.State.SequenceRemoteControllable != nil {
							v := *ps.State.SequenceRemoteControllable
							rps.State.SequenceRemoteControllable = &v
						}
						if ps.State.ActiveRepetitionNumber != nil {
							v := *ps.State.ActiveRepetitionNumber
							rps.State.ActiveRepetitionNumber = &v
						}
						if ps.State.RemainingPauseTime != nil {
							v := *ps.State.RemainingPauseTime
							rps.State.RemainingPauseTime = &v
						}
					}

					// Copy Schedule
					if ps.Schedule != nil {
						rps.Schedule = &PowerSequenceScheduleDataType{}
						if ps.Schedule.StartTime != nil {
							v := *ps.Schedule.StartTime
							rps.Schedule.StartTime = &v
						}
						if ps.Schedule.EndTime != nil {
							v := *ps.Schedule.EndTime
							rps.Schedule.EndTime = &v
						}
					}

					// Copy PowerTimeSlot
					if len(ps.PowerTimeSlot) > 0 {
						rps.PowerTimeSlot = make([]SmartEnergyManagementPsPowerTimeSlotType, len(ps.PowerTimeSlot))
						for k := range ps.PowerTimeSlot {
							pts := &ps.PowerTimeSlot[k]
							rpts := &rps.PowerTimeSlot[k]

							// Copy Schedule
							if pts.Schedule != nil {
								rpts.Schedule = &PowerTimeSlotScheduleDataType{}
								if pts.Schedule.SlotNumber != nil {
									v := *pts.Schedule.SlotNumber
									rpts.Schedule.SlotNumber = &v
								}
							}

							// Copy ValueList
							if pts.ValueList != nil {
								rpts.ValueList = &SmartEnergyManagementPsPowerTimeSlotValueListType{}
								if len(pts.ValueList.Value) > 0 {
									rpts.ValueList.Value = make([]PowerTimeSlotValueDataType, len(pts.ValueList.Value))
									for l := range pts.ValueList.Value {
										val := &pts.ValueList.Value[l]
										rval := &rpts.ValueList.Value[l]

										if val.ValueType != nil {
											v := *val.ValueType
											rval.ValueType = &v
										}
										if val.Value != nil {
											rval.Value = &ScaledNumberType{}
											if val.Value.Number != nil {
												v := *val.Value.Number
												rval.Value.Number = &v
											}
											if val.Value.Scale != nil {
												v := *val.Value.Scale
												rval.Value.Scale = &v
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return result
}
