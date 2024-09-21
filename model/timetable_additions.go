package model

// TimeTableListDataType

var _ Updater = (*TimeTableListDataType)(nil)

func (r *TimeTableListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []TimeTableDataType
	if newList != nil {
		newData = newList.(*TimeTableListDataType).TimeTableData
	}

	data, success := UpdateList(remoteWrite, r.TimeTableData, newData, filterPartial, filterDelete)

	if success && persist {
		r.TimeTableData = data
	}

	return data, success
}

// TimeTableConstraintsListDataType

var _ Updater = (*TimeTableConstraintsListDataType)(nil)

func (r *TimeTableConstraintsListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []TimeTableConstraintsDataType
	if newList != nil {
		newData = newList.(*TimeTableConstraintsListDataType).TimeTableConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.TimeTableConstraintsData, newData, filterPartial, filterDelete)

	if success && persist {
		r.TimeTableConstraintsData = data
	}

	return data, success
}

// TimeTableDescriptionListDataType

var _ Updater = (*TimeTableDescriptionListDataType)(nil)

func (r *TimeTableDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []TimeTableDescriptionDataType
	if newList != nil {
		newData = newList.(*TimeTableDescriptionListDataType).TimeTableDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.TimeTableDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.TimeTableDescriptionData = data
	}

	return data, success
}
