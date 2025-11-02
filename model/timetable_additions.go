package model

// TimeTableListDataType

var _ Updater = (*TimeTableListDataType)(nil)

func (r *TimeTableListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []TimeTableDataType
	if newList != nil {
		newData = newList.(*TimeTableListDataType).TimeTableData
	}

	data, success := UpdateList(remoteWrite, r.TimeTableData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.TimeTableData = data
	}

	return data, success
}

// TimeTableConstraintsListDataType

var _ Updater = (*TimeTableConstraintsListDataType)(nil)

func (r *TimeTableConstraintsListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []TimeTableConstraintsDataType
	if newList != nil {
		newData = newList.(*TimeTableConstraintsListDataType).TimeTableConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.TimeTableConstraintsData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.TimeTableConstraintsData = data
	}

	return data, success
}

// TimeTableDescriptionListDataType

var _ Updater = (*TimeTableDescriptionListDataType)(nil)

func (r *TimeTableDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []TimeTableDescriptionDataType
	if newList != nil {
		newData = newList.(*TimeTableDescriptionListDataType).TimeTableDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.TimeTableDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.TimeTableDescriptionData = data
	}

	return data, success
}

// TimeTableListDataType PartialReader implementation
var _ PartialReader = (*TimeTableListDataType)(nil)

func (r *TimeTableListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[TimeTableDataType, TimeTableListDataSelectorsType, TimeTableDataElementsType](r.TimeTableData, filter)
	if !success {
		return r, false
	}

	result := &TimeTableListDataType{
		TimeTableData: filteredItems,
	}
	return result, true
}

// TimeTableConstraintsListDataType PartialReader implementation
var _ PartialReader = (*TimeTableConstraintsListDataType)(nil)

func (r *TimeTableConstraintsListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[TimeTableConstraintsDataType, TimeTableConstraintsListDataSelectorsType, TimeTableConstraintsDataElementsType](r.TimeTableConstraintsData, filter)
	if !success {
		return r, false
	}

	result := &TimeTableConstraintsListDataType{
		TimeTableConstraintsData: filteredItems,
	}
	return result, true
}

// TimeTableDescriptionListDataType PartialReader implementation
var _ PartialReader = (*TimeTableDescriptionListDataType)(nil)

func (r *TimeTableDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[TimeTableDescriptionDataType, TimeTableDescriptionListDataSelectorsType, TimeTableDescriptionDataElementsType](r.TimeTableDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &TimeTableDescriptionListDataType{
		TimeTableDescriptionData: filteredItems,
	}
	return result, true
}
