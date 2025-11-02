package model

// TimeSeriesListDataType

var _ Updater = (*TimeSeriesListDataType)(nil)

func (r *TimeSeriesListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []TimeSeriesDataType
	if newList != nil {
		newData = newList.(*TimeSeriesListDataType).TimeSeriesData
	}

	data, success := UpdateList(remoteWrite, r.TimeSeriesData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.TimeSeriesData = data
	}

	return data, success
}

// TimeSeriesDescriptionListDataType

var _ Updater = (*TimeSeriesDescriptionListDataType)(nil)

func (r *TimeSeriesDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []TimeSeriesDescriptionDataType
	if newList != nil {
		newData = newList.(*TimeSeriesDescriptionListDataType).TimeSeriesDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.TimeSeriesDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.TimeSeriesDescriptionData = data
	}

	return data, success
}

// TimeSeriesConstraintsListDataType

var _ Updater = (*TimeSeriesConstraintsListDataType)(nil)

func (r *TimeSeriesConstraintsListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []TimeSeriesConstraintsDataType
	if newList != nil {
		newData = newList.(*TimeSeriesConstraintsListDataType).TimeSeriesConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.TimeSeriesConstraintsData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.TimeSeriesConstraintsData = data
	}

	return data, success
}

// TimeSeriesListDataType PartialReader implementation
var _ PartialReader = (*TimeSeriesListDataType)(nil)

func (r *TimeSeriesListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[TimeSeriesDataType, TimeSeriesListDataSelectorsType, TimeSeriesDataElementsType](r.TimeSeriesData, filter)
	if !success {
		return r, false
	}

	result := &TimeSeriesListDataType{
		TimeSeriesData: filteredItems,
	}
	return result, true
}

// TimeSeriesDescriptionListDataType PartialReader implementation
var _ PartialReader = (*TimeSeriesDescriptionListDataType)(nil)

func (r *TimeSeriesDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[TimeSeriesDescriptionDataType, TimeSeriesDescriptionListDataSelectorsType, TimeSeriesDescriptionDataElementsType](r.TimeSeriesDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &TimeSeriesDescriptionListDataType{
		TimeSeriesDescriptionData: filteredItems,
	}
	return result, true
}

// TimeSeriesConstraintsListDataType PartialReader implementation
var _ PartialReader = (*TimeSeriesConstraintsListDataType)(nil)

func (r *TimeSeriesConstraintsListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[TimeSeriesConstraintsDataType, TimeSeriesConstraintsListDataSelectorsType, TimeSeriesConstraintsDataElementsType](r.TimeSeriesConstraintsData, filter)
	if !success {
		return r, false
	}

	result := &TimeSeriesConstraintsListDataType{
		TimeSeriesConstraintsData: filteredItems,
	}
	return result, true
}
