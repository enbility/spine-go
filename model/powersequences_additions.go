package model

// PowerTimeSlotScheduleListDataType

var _ Updater = (*PowerTimeSlotScheduleListDataType)(nil)

func (r *PowerTimeSlotScheduleListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []PowerTimeSlotScheduleDataType
	if newList != nil {
		newData = newList.(*PowerTimeSlotScheduleListDataType).PowerTimeSlotScheduleData
	}

	data, success := UpdateList(remoteWrite, r.PowerTimeSlotScheduleData, newData, filterPartial, filterDelete)

	if success {
		r.PowerTimeSlotScheduleData = data
	}

	return success
}

// PowerTimeSlotValueListDataType

var _ Updater = (*PowerTimeSlotValueListDataType)(nil)

func (r *PowerTimeSlotValueListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []PowerTimeSlotValueDataType
	if newList != nil {
		newData = newList.(*PowerTimeSlotValueListDataType).PowerTimeSlotValueData
	}

	data, success := UpdateList(remoteWrite, r.PowerTimeSlotValueData, newData, filterPartial, filterDelete)

	if success {
		r.PowerTimeSlotValueData = data
	}

	return success
}

// PowerTimeSlotScheduleConstraintsListDataType

var _ Updater = (*PowerTimeSlotScheduleConstraintsListDataType)(nil)

func (r *PowerTimeSlotScheduleConstraintsListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []PowerTimeSlotScheduleConstraintsDataType
	if newList != nil {
		newData = newList.(*PowerTimeSlotScheduleConstraintsListDataType).PowerTimeSlotScheduleConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.PowerTimeSlotScheduleConstraintsData, newData, filterPartial, filterDelete)

	if success {
		r.PowerTimeSlotScheduleConstraintsData = data
	}

	return success
}

// PowerSequenceAlternativesRelationListDataType

var _ Updater = (*PowerSequenceAlternativesRelationListDataType)(nil)

func (r *PowerSequenceAlternativesRelationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []PowerSequenceAlternativesRelationDataType
	if newList != nil {
		newData = newList.(*PowerSequenceAlternativesRelationListDataType).PowerSequenceAlternativesRelationData
	}

	data, success := UpdateList(remoteWrite, r.PowerSequenceAlternativesRelationData, newData, filterPartial, filterDelete)

	if success {
		r.PowerSequenceAlternativesRelationData = data
	}

	return success
}

// PowerSequenceDescriptionListDataType

var _ Updater = (*PowerSequenceDescriptionListDataType)(nil)

func (r *PowerSequenceDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []PowerSequenceDescriptionDataType
	if newList != nil {
		newData = newList.(*PowerSequenceDescriptionListDataType).PowerSequenceDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.PowerSequenceDescriptionData, newData, filterPartial, filterDelete)

	if success {
		r.PowerSequenceDescriptionData = data
	}

	return success
}

// PowerSequenceStateListDataType

var _ Updater = (*PowerSequenceStateListDataType)(nil)

func (r *PowerSequenceStateListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []PowerSequenceStateDataType
	if newList != nil {
		newData = newList.(*PowerSequenceStateListDataType).PowerSequenceStateData
	}

	data, success := UpdateList(remoteWrite, r.PowerSequenceStateData, newData, filterPartial, filterDelete)

	if success {
		r.PowerSequenceStateData = data
	}

	return success
}

// PowerSequenceScheduleListDataType

var _ Updater = (*PowerSequenceScheduleListDataType)(nil)

func (r *PowerSequenceScheduleListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []PowerSequenceScheduleDataType
	if newList != nil {
		newData = newList.(*PowerSequenceScheduleListDataType).PowerSequenceScheduleData
	}

	data, success := UpdateList(remoteWrite, r.PowerSequenceScheduleData, newData, filterPartial, filterDelete)

	if success {
		r.PowerSequenceScheduleData = data
	}

	return success
}

// PowerSequenceScheduleConstraintsListDataType

var _ Updater = (*PowerSequenceScheduleConstraintsListDataType)(nil)

func (r *PowerSequenceScheduleConstraintsListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []PowerSequenceScheduleConstraintsDataType
	if newList != nil {
		newData = newList.(*PowerSequenceScheduleConstraintsListDataType).PowerSequenceScheduleConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.PowerSequenceScheduleConstraintsData, newData, filterPartial, filterDelete)

	if success {
		r.PowerSequenceScheduleConstraintsData = data
	}

	return success
}

// PowerSequencePriceListDataType

var _ Updater = (*PowerSequencePriceListDataType)(nil)

func (r *PowerSequencePriceListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []PowerSequencePriceDataType
	if newList != nil {
		newData = newList.(*PowerSequencePriceListDataType).PowerSequencePriceData
	}

	data, success := UpdateList(remoteWrite, r.PowerSequencePriceData, newData, filterPartial, filterDelete)

	if success {
		r.PowerSequencePriceData = data
	}

	return success
}

// PowerSequenceSchedulePreferenceListDataType

var _ Updater = (*PowerSequenceSchedulePreferenceListDataType)(nil)

func (r *PowerSequenceSchedulePreferenceListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []PowerSequenceSchedulePreferenceDataType
	if newList != nil {
		newData = newList.(*PowerSequenceSchedulePreferenceListDataType).PowerSequenceSchedulePreferenceData
	}

	data, success := UpdateList(remoteWrite, r.PowerSequenceSchedulePreferenceData, newData, filterPartial, filterDelete)

	if success {
		r.PowerSequenceSchedulePreferenceData = data
	}

	return success
}
