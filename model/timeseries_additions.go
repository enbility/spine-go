package model

// TimeSeriesListDataType

var _ Updater = (*TimeSeriesListDataType)(nil)

func (r *TimeSeriesListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []TimeSeriesDataType
	if newList != nil {
		newData = newList.(*TimeSeriesListDataType).TimeSeriesData
	}

	data, success := UpdateList(remoteWrite, r.TimeSeriesData, newData, filterPartial, filterDelete)

	if success && persist {
		r.TimeSeriesData = data
	}

	return data, success
}

// TimeSeriesDescriptionListDataType

var _ Updater = (*TimeSeriesDescriptionListDataType)(nil)

func (r *TimeSeriesDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []TimeSeriesDescriptionDataType
	if newList != nil {
		newData = newList.(*TimeSeriesDescriptionListDataType).TimeSeriesDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.TimeSeriesDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.TimeSeriesDescriptionData = data
	}

	return data, success
}

// TimeSeriesConstraintsListDataType

var _ Updater = (*TimeSeriesConstraintsListDataType)(nil)

func (r *TimeSeriesConstraintsListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []TimeSeriesConstraintsDataType
	if newList != nil {
		newData = newList.(*TimeSeriesConstraintsListDataType).TimeSeriesConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.TimeSeriesConstraintsData, newData, filterPartial, filterDelete)

	if success && persist {
		r.TimeSeriesConstraintsData = data
	}

	return data, success
}
