package model

// TimeSeriesListDataType

var _ Updater = (*TimeSeriesListDataType)(nil)

func (r *TimeSeriesListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []TimeSeriesDataType
	if newList != nil {
		newData = newList.(*TimeSeriesListDataType).TimeSeriesData
	}

	r.TimeSeriesData = UpdateList(remoteWrite, r.TimeSeriesData, newData, filterPartial, filterDelete)
}

// TimeSeriesDescriptionListDataType

var _ Updater = (*TimeSeriesDescriptionListDataType)(nil)

func (r *TimeSeriesDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []TimeSeriesDescriptionDataType
	if newList != nil {
		newData = newList.(*TimeSeriesDescriptionListDataType).TimeSeriesDescriptionData
	}

	r.TimeSeriesDescriptionData = UpdateList(remoteWrite, r.TimeSeriesDescriptionData, newData, filterPartial, filterDelete)
}

// TimeSeriesConstraintsListDataType

var _ Updater = (*TimeSeriesConstraintsListDataType)(nil)

func (r *TimeSeriesConstraintsListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []TimeSeriesConstraintsDataType
	if newList != nil {
		newData = newList.(*TimeSeriesConstraintsListDataType).TimeSeriesConstraintsData
	}

	r.TimeSeriesConstraintsData = UpdateList(remoteWrite, r.TimeSeriesConstraintsData, newData, filterPartial, filterDelete)
}
