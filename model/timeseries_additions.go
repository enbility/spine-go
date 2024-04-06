package model

// TimeSeriesListDataType

var _ Updater = (*TimeSeriesListDataType)(nil)

func (r *TimeSeriesListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []TimeSeriesDataType
	if newList != nil {
		newData = newList.(*TimeSeriesListDataType).TimeSeriesData
	}

	data, success := UpdateList(remoteWrite, r.TimeSeriesData, newData, filterPartial, filterDelete)

	if success {
		r.TimeSeriesData = data
	}

	return success
}

// TimeSeriesDescriptionListDataType

var _ Updater = (*TimeSeriesDescriptionListDataType)(nil)

func (r *TimeSeriesDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []TimeSeriesDescriptionDataType
	if newList != nil {
		newData = newList.(*TimeSeriesDescriptionListDataType).TimeSeriesDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.TimeSeriesDescriptionData, newData, filterPartial, filterDelete)

	if success {
		r.TimeSeriesDescriptionData = data
	}

	return success
}

// TimeSeriesConstraintsListDataType

var _ Updater = (*TimeSeriesConstraintsListDataType)(nil)

func (r *TimeSeriesConstraintsListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []TimeSeriesConstraintsDataType
	if newList != nil {
		newData = newList.(*TimeSeriesConstraintsListDataType).TimeSeriesConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.TimeSeriesConstraintsData, newData, filterPartial, filterDelete)

	if success {
		r.TimeSeriesConstraintsData = data
	}

	return success
}
