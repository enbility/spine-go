package model

// PowerTimeSlotScheduleListDataType

var _ Updater = (*PowerTimeSlotScheduleListDataType)(nil)

func (r *PowerTimeSlotScheduleListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []PowerTimeSlotScheduleDataType
	if newList != nil {
		newData = newList.(*PowerTimeSlotScheduleListDataType).PowerTimeSlotScheduleData
	}

	data, success := UpdateList(remoteWrite, r.PowerTimeSlotScheduleData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.PowerTimeSlotScheduleData = data
	}

	return data, success
}

// PowerTimeSlotValueListDataType

var _ Updater = (*PowerTimeSlotValueListDataType)(nil)

func (r *PowerTimeSlotValueListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []PowerTimeSlotValueDataType
	if newList != nil {
		newData = newList.(*PowerTimeSlotValueListDataType).PowerTimeSlotValueData
	}

	data, success := UpdateList(remoteWrite, r.PowerTimeSlotValueData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.PowerTimeSlotValueData = data
	}

	return data, success
}

// PowerTimeSlotScheduleConstraintsListDataType

var _ Updater = (*PowerTimeSlotScheduleConstraintsListDataType)(nil)

func (r *PowerTimeSlotScheduleConstraintsListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []PowerTimeSlotScheduleConstraintsDataType
	if newList != nil {
		newData = newList.(*PowerTimeSlotScheduleConstraintsListDataType).PowerTimeSlotScheduleConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.PowerTimeSlotScheduleConstraintsData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.PowerTimeSlotScheduleConstraintsData = data
	}

	return data, success
}

// PowerSequenceAlternativesRelationListDataType

var _ Updater = (*PowerSequenceAlternativesRelationListDataType)(nil)

func (r *PowerSequenceAlternativesRelationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []PowerSequenceAlternativesRelationDataType
	if newList != nil {
		newData = newList.(*PowerSequenceAlternativesRelationListDataType).PowerSequenceAlternativesRelationData
	}

	data, success := UpdateList(remoteWrite, r.PowerSequenceAlternativesRelationData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.PowerSequenceAlternativesRelationData = data
	}

	return data, success
}

// PowerSequenceDescriptionListDataType

var _ Updater = (*PowerSequenceDescriptionListDataType)(nil)

func (r *PowerSequenceDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []PowerSequenceDescriptionDataType
	if newList != nil {
		newData = newList.(*PowerSequenceDescriptionListDataType).PowerSequenceDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.PowerSequenceDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.PowerSequenceDescriptionData = data
	}

	return data, success
}

// PowerSequenceStateListDataType

var _ Updater = (*PowerSequenceStateListDataType)(nil)

func (r *PowerSequenceStateListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []PowerSequenceStateDataType
	if newList != nil {
		newData = newList.(*PowerSequenceStateListDataType).PowerSequenceStateData
	}

	data, success := UpdateList(remoteWrite, r.PowerSequenceStateData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.PowerSequenceStateData = data
	}

	return data, success
}

// PowerSequenceScheduleListDataType

var _ Updater = (*PowerSequenceScheduleListDataType)(nil)

func (r *PowerSequenceScheduleListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []PowerSequenceScheduleDataType
	if newList != nil {
		newData = newList.(*PowerSequenceScheduleListDataType).PowerSequenceScheduleData
	}

	data, success := UpdateList(remoteWrite, r.PowerSequenceScheduleData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.PowerSequenceScheduleData = data
	}

	return data, success
}

// PowerSequenceScheduleConstraintsListDataType

var _ Updater = (*PowerSequenceScheduleConstraintsListDataType)(nil)

func (r *PowerSequenceScheduleConstraintsListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []PowerSequenceScheduleConstraintsDataType
	if newList != nil {
		newData = newList.(*PowerSequenceScheduleConstraintsListDataType).PowerSequenceScheduleConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.PowerSequenceScheduleConstraintsData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.PowerSequenceScheduleConstraintsData = data
	}

	return data, success
}

// PowerSequencePriceListDataType

var _ Updater = (*PowerSequencePriceListDataType)(nil)

func (r *PowerSequencePriceListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []PowerSequencePriceDataType
	if newList != nil {
		newData = newList.(*PowerSequencePriceListDataType).PowerSequencePriceData
	}

	data, success := UpdateList(remoteWrite, r.PowerSequencePriceData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.PowerSequencePriceData = data
	}

	return data, success
}

// PowerSequenceSchedulePreferenceListDataType

var _ Updater = (*PowerSequenceSchedulePreferenceListDataType)(nil)

func (r *PowerSequenceSchedulePreferenceListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []PowerSequenceSchedulePreferenceDataType
	if newList != nil {
		newData = newList.(*PowerSequenceSchedulePreferenceListDataType).PowerSequenceSchedulePreferenceData
	}

	data, success := UpdateList(remoteWrite, r.PowerSequenceSchedulePreferenceData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.PowerSequenceSchedulePreferenceData = data
	}

	return data, success
}

// PowerTimeSlotScheduleListDataType PartialReader implementation
var _ PartialReader = (*PowerTimeSlotScheduleListDataType)(nil)

func (r *PowerTimeSlotScheduleListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[PowerTimeSlotScheduleDataType, PowerTimeSlotScheduleListDataSelectorsType, PowerTimeSlotScheduleDataElementsType](r.PowerTimeSlotScheduleData, filter)
	if !success {
		return r, false
	}

	result := &PowerTimeSlotScheduleListDataType{
		PowerTimeSlotScheduleData: filteredItems,
	}
	return result, true
}

// PowerTimeSlotValueListDataType PartialReader implementation
var _ PartialReader = (*PowerTimeSlotValueListDataType)(nil)

func (r *PowerTimeSlotValueListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[PowerTimeSlotValueDataType, PowerTimeSlotValueListDataSelectorsType, PowerTimeSlotValueDataElementsType](r.PowerTimeSlotValueData, filter)
	if !success {
		return r, false
	}

	result := &PowerTimeSlotValueListDataType{
		PowerTimeSlotValueData: filteredItems,
	}
	return result, true
}

// PowerTimeSlotScheduleConstraintsListDataType PartialReader implementation
var _ PartialReader = (*PowerTimeSlotScheduleConstraintsListDataType)(nil)

func (r *PowerTimeSlotScheduleConstraintsListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[PowerTimeSlotScheduleConstraintsDataType, PowerTimeSlotScheduleConstraintsListDataSelectorsType, PowerTimeSlotScheduleConstraintsDataElementsType](r.PowerTimeSlotScheduleConstraintsData, filter)
	if !success {
		return r, false
	}

	result := &PowerTimeSlotScheduleConstraintsListDataType{
		PowerTimeSlotScheduleConstraintsData: filteredItems,
	}
	return result, true
}

// PowerSequenceAlternativesRelationListDataType PartialReader implementation
var _ PartialReader = (*PowerSequenceAlternativesRelationListDataType)(nil)

func (r *PowerSequenceAlternativesRelationListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[PowerSequenceAlternativesRelationDataType, PowerSequenceAlternativesRelationListDataSelectorsType, PowerSequenceAlternativesRelationDataElementsType](r.PowerSequenceAlternativesRelationData, filter)
	if !success {
		return r, false
	}

	result := &PowerSequenceAlternativesRelationListDataType{
		PowerSequenceAlternativesRelationData: filteredItems,
	}
	return result, true
}

// PowerSequenceDescriptionListDataType PartialReader implementation
var _ PartialReader = (*PowerSequenceDescriptionListDataType)(nil)

func (r *PowerSequenceDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[PowerSequenceDescriptionDataType, PowerSequenceDescriptionListDataSelectorsType, PowerSequenceDescriptionDataElementsType](r.PowerSequenceDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &PowerSequenceDescriptionListDataType{
		PowerSequenceDescriptionData: filteredItems,
	}
	return result, true
}

// PowerSequenceStateListDataType PartialReader implementation
var _ PartialReader = (*PowerSequenceStateListDataType)(nil)

func (r *PowerSequenceStateListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[PowerSequenceStateDataType, PowerSequenceStateListDataSelectorsType, PowerSequenceStateDataElementsType](r.PowerSequenceStateData, filter)
	if !success {
		return r, false
	}

	result := &PowerSequenceStateListDataType{
		PowerSequenceStateData: filteredItems,
	}
	return result, true
}

// PowerSequenceScheduleListDataType PartialReader implementation
var _ PartialReader = (*PowerSequenceScheduleListDataType)(nil)

func (r *PowerSequenceScheduleListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[PowerSequenceScheduleDataType, PowerSequenceScheduleListDataSelectorsType, PowerSequenceScheduleDataElementsType](r.PowerSequenceScheduleData, filter)
	if !success {
		return r, false
	}

	result := &PowerSequenceScheduleListDataType{
		PowerSequenceScheduleData: filteredItems,
	}
	return result, true
}

// PowerSequenceScheduleConstraintsListDataType PartialReader implementation
var _ PartialReader = (*PowerSequenceScheduleConstraintsListDataType)(nil)

func (r *PowerSequenceScheduleConstraintsListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[PowerSequenceScheduleConstraintsDataType, PowerSequenceScheduleConstraintsListDataSelectorsType, PowerSequenceScheduleConstraintsDataElementsType](r.PowerSequenceScheduleConstraintsData, filter)
	if !success {
		return r, false
	}

	result := &PowerSequenceScheduleConstraintsListDataType{
		PowerSequenceScheduleConstraintsData: filteredItems,
	}
	return result, true
}

// PowerSequencePriceListDataType PartialReader implementation
var _ PartialReader = (*PowerSequencePriceListDataType)(nil)

func (r *PowerSequencePriceListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[PowerSequencePriceDataType, PowerSequencePriceListDataSelectorsType, PowerSequencePriceDataElementsType](r.PowerSequencePriceData, filter)
	if !success {
		return r, false
	}

	result := &PowerSequencePriceListDataType{
		PowerSequencePriceData: filteredItems,
	}
	return result, true
}

// PowerSequenceSchedulePreferenceListDataType PartialReader implementation
var _ PartialReader = (*PowerSequenceSchedulePreferenceListDataType)(nil)

func (r *PowerSequenceSchedulePreferenceListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[PowerSequenceSchedulePreferenceDataType, PowerSequenceSchedulePreferenceListDataSelectorsType, PowerSequenceSchedulePreferenceDataElementsType](r.PowerSequenceSchedulePreferenceData, filter)
	if !success {
		return r, false
	}

	result := &PowerSequenceSchedulePreferenceListDataType{
		PowerSequenceSchedulePreferenceData: filteredItems,
	}
	return result, true
}
