package model

// TimeTableListDataType

var _ Updater = (*TimeTableListDataType)(nil)

func (r *TimeTableListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []TimeTableDataType
	if newList != nil {
		newData = newList.(*TimeTableListDataType).TimeTableData
	}

	r.TimeTableData = UpdateList(remoteWrite, r.TimeTableData, newData, filterPartial, filterDelete)
}

// TimeTableConstraintsListDataType

var _ Updater = (*TimeTableConstraintsListDataType)(nil)

func (r *TimeTableConstraintsListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []TimeTableConstraintsDataType
	if newList != nil {
		newData = newList.(*TimeTableConstraintsListDataType).TimeTableConstraintsData
	}

	r.TimeTableConstraintsData = UpdateList(remoteWrite, r.TimeTableConstraintsData, newData, filterPartial, filterDelete)
}

// TimeTableDescriptionListDataType

var _ Updater = (*TimeTableDescriptionListDataType)(nil)

func (r *TimeTableDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []TimeTableDescriptionDataType
	if newList != nil {
		newData = newList.(*TimeTableDescriptionListDataType).TimeTableDescriptionData
	}

	r.TimeTableDescriptionData = UpdateList(remoteWrite, r.TimeTableDescriptionData, newData, filterPartial, filterDelete)
}
