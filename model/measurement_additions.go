package model

// MeasurementListDataType

var _ Updater = (*MeasurementListDataType)(nil)

func (r *MeasurementListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []MeasurementDataType
	if newList != nil {
		newData = newList.(*MeasurementListDataType).MeasurementData
	}

	data, success := UpdateList(remoteWrite, r.MeasurementData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.MeasurementData = data
	}

	return data, success
}

// MeasurementSeriesListDataType

var _ Updater = (*MeasurementSeriesListDataType)(nil)

func (r *MeasurementSeriesListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []MeasurementSeriesDataType
	if newList != nil {
		newData = newList.(*MeasurementSeriesListDataType).MeasurementSeriesData
	}

	data, success := UpdateList(remoteWrite, r.MeasurementSeriesData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.MeasurementSeriesData = data
	}

	return data, success
}

// MeasurementConstraintsListDataType

var _ Updater = (*MeasurementConstraintsListDataType)(nil)

func (r *MeasurementConstraintsListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []MeasurementConstraintsDataType
	if newList != nil {
		newData = newList.(*MeasurementConstraintsListDataType).MeasurementConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.MeasurementConstraintsData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.MeasurementConstraintsData = data
	}

	return data, success
}

// MeasurementDescriptionListDataType

var _ Updater = (*MeasurementDescriptionListDataType)(nil)

func (r *MeasurementDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []MeasurementDescriptionDataType
	if newList != nil {
		newData = newList.(*MeasurementDescriptionListDataType).MeasurementDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.MeasurementDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.MeasurementDescriptionData = data
	}

	return data, success
}

// MeasurementThresholdRelationListDataType

var _ Updater = (*MeasurementThresholdRelationListDataType)(nil)

func (r *MeasurementThresholdRelationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []MeasurementThresholdRelationDataType
	if newList != nil {
		newData = newList.(*MeasurementThresholdRelationListDataType).MeasurementThresholdRelationData
	}

	data, success := UpdateList(remoteWrite, r.MeasurementThresholdRelationData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.MeasurementThresholdRelationData = data
	}

	return data, success
}

// MeasurementListDataType PartialReader implementation
var _ PartialReader = (*MeasurementListDataType)(nil)

func (r *MeasurementListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function with both selector and elements filtering
	filteredItems, success := partialListDataRead[MeasurementDataType, MeasurementListDataSelectorsType, MeasurementDataElementsType](r.MeasurementData, filter)
	if !success {
		return r, false
	}

	result := &MeasurementListDataType{
		MeasurementData: filteredItems,
	}
	return result, true
}

// MeasurementConstraintsListDataType PartialReader implementation
var _ PartialReader = (*MeasurementConstraintsListDataType)(nil)

func (r *MeasurementConstraintsListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function with both selector and elements filtering
	filteredItems, success := partialListDataRead[MeasurementConstraintsDataType, MeasurementConstraintsListDataSelectorsType, MeasurementConstraintsDataElementsType](r.MeasurementConstraintsData, filter)
	if !success {
		return r, false
	}

	result := &MeasurementConstraintsListDataType{
		MeasurementConstraintsData: filteredItems,
	}
	return result, true
}

// MeasurementDescriptionListDataType PartialReader implementation
var _ PartialReader = (*MeasurementDescriptionListDataType)(nil)

func (r *MeasurementDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function with both selector and elements filtering
	filteredItems, success := partialListDataRead[MeasurementDescriptionDataType, MeasurementDescriptionListDataSelectorsType, MeasurementDescriptionDataElementsType](r.MeasurementDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &MeasurementDescriptionListDataType{
		MeasurementDescriptionData: filteredItems,
	}
	return result, true
}

// MeasurementSeriesListDataType PartialReader implementation
var _ PartialReader = (*MeasurementSeriesListDataType)(nil)

func (r *MeasurementSeriesListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[MeasurementSeriesDataType, MeasurementSeriesListDataSelectorsType, MeasurementSeriesDataElementsType](r.MeasurementSeriesData, filter)
	if !success {
		return r, false
	}

	result := &MeasurementSeriesListDataType{
		MeasurementSeriesData: filteredItems,
	}
	return result, true
}

// MeasurementThresholdRelationListDataType PartialReader implementation
var _ PartialReader = (*MeasurementThresholdRelationListDataType)(nil)

func (r *MeasurementThresholdRelationListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[MeasurementThresholdRelationDataType, MeasurementThresholdRelationListDataSelectorsType, MeasurementThresholdRelationDataElementsType](r.MeasurementThresholdRelationData, filter)
	if !success {
		return r, false
	}

	result := &MeasurementThresholdRelationListDataType{
		MeasurementThresholdRelationData: filteredItems,
	}
	return result, true
}
