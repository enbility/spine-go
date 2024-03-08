package model

// PowerTimeSlotScheduleListDataType

var _ Updater = (*PowerTimeSlotScheduleListDataType)(nil)

func (r *PowerTimeSlotScheduleListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []PowerTimeSlotScheduleDataType
	if newList != nil {
		newData = newList.(*PowerTimeSlotScheduleListDataType).PowerTimeSlotScheduleData
	}

	r.PowerTimeSlotScheduleData = UpdateList(remoteWrite, r.PowerTimeSlotScheduleData, newData, filterPartial, filterDelete)
}

// PowerTimeSlotValueListDataType

var _ Updater = (*PowerTimeSlotValueListDataType)(nil)

func (r *PowerTimeSlotValueListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []PowerTimeSlotValueDataType
	if newList != nil {
		newData = newList.(*PowerTimeSlotValueListDataType).PowerTimeSlotValueData
	}

	r.PowerTimeSlotValueData = UpdateList(remoteWrite, r.PowerTimeSlotValueData, newData, filterPartial, filterDelete)
}

// PowerTimeSlotScheduleConstraintsListDataType

var _ Updater = (*PowerTimeSlotScheduleConstraintsListDataType)(nil)

func (r *PowerTimeSlotScheduleConstraintsListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []PowerTimeSlotScheduleConstraintsDataType
	if newList != nil {
		newData = newList.(*PowerTimeSlotScheduleConstraintsListDataType).PowerTimeSlotScheduleConstraintsData
	}

	r.PowerTimeSlotScheduleConstraintsData = UpdateList(remoteWrite, r.PowerTimeSlotScheduleConstraintsData, newData, filterPartial, filterDelete)
}

// PowerSequenceAlternativesRelationListDataType

var _ Updater = (*PowerSequenceAlternativesRelationListDataType)(nil)

func (r *PowerSequenceAlternativesRelationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []PowerSequenceAlternativesRelationDataType
	if newList != nil {
		newData = newList.(*PowerSequenceAlternativesRelationListDataType).PowerSequenceAlternativesRelationData
	}

	r.PowerSequenceAlternativesRelationData = UpdateList(remoteWrite, r.PowerSequenceAlternativesRelationData, newData, filterPartial, filterDelete)
}

// PowerSequenceDescriptionListDataType

var _ Updater = (*PowerSequenceDescriptionListDataType)(nil)

func (r *PowerSequenceDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []PowerSequenceDescriptionDataType
	if newList != nil {
		newData = newList.(*PowerSequenceDescriptionListDataType).PowerSequenceDescriptionData
	}

	r.PowerSequenceDescriptionData = UpdateList(remoteWrite, r.PowerSequenceDescriptionData, newData, filterPartial, filterDelete)
}

// PowerSequenceStateListDataType

var _ Updater = (*PowerSequenceStateListDataType)(nil)

func (r *PowerSequenceStateListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []PowerSequenceStateDataType
	if newList != nil {
		newData = newList.(*PowerSequenceStateListDataType).PowerSequenceStateData
	}

	r.PowerSequenceStateData = UpdateList(remoteWrite, r.PowerSequenceStateData, newData, filterPartial, filterDelete)
}

// PowerSequenceScheduleListDataType

var _ Updater = (*PowerSequenceScheduleListDataType)(nil)

func (r *PowerSequenceScheduleListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []PowerSequenceScheduleDataType
	if newList != nil {
		newData = newList.(*PowerSequenceScheduleListDataType).PowerSequenceScheduleData
	}

	r.PowerSequenceScheduleData = UpdateList(remoteWrite, r.PowerSequenceScheduleData, newData, filterPartial, filterDelete)
}

// PowerSequenceScheduleConstraintsListDataType

var _ Updater = (*PowerSequenceScheduleConstraintsListDataType)(nil)

func (r *PowerSequenceScheduleConstraintsListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []PowerSequenceScheduleConstraintsDataType
	if newList != nil {
		newData = newList.(*PowerSequenceScheduleConstraintsListDataType).PowerSequenceScheduleConstraintsData
	}

	r.PowerSequenceScheduleConstraintsData = UpdateList(remoteWrite, r.PowerSequenceScheduleConstraintsData, newData, filterPartial, filterDelete)
}

// PowerSequencePriceListDataType

var _ Updater = (*PowerSequencePriceListDataType)(nil)

func (r *PowerSequencePriceListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []PowerSequencePriceDataType
	if newList != nil {
		newData = newList.(*PowerSequencePriceListDataType).PowerSequencePriceData
	}

	r.PowerSequencePriceData = UpdateList(remoteWrite, r.PowerSequencePriceData, newData, filterPartial, filterDelete)
}

// PowerSequenceSchedulePreferenceListDataType

var _ Updater = (*PowerSequenceSchedulePreferenceListDataType)(nil)

func (r *PowerSequenceSchedulePreferenceListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []PowerSequenceSchedulePreferenceDataType
	if newList != nil {
		newData = newList.(*PowerSequenceSchedulePreferenceListDataType).PowerSequenceSchedulePreferenceData
	}

	r.PowerSequenceSchedulePreferenceData = UpdateList(remoteWrite, r.PowerSequenceSchedulePreferenceData, newData, filterPartial, filterDelete)
}
