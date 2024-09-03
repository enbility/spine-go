package model

// MeasurementListDataType

var _ Updater = (*MeasurementListDataType)(nil)

func (r *MeasurementListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []MeasurementDataType
	if newList != nil {
		newData = newList.(*MeasurementListDataType).MeasurementData
	}

	data, success := UpdateList(remoteWrite, r.MeasurementData, newData, filterPartial, filterDelete)

	if success && persist {
		r.MeasurementData = data
	}

	return data, success
}

// MeasurementSeriesListDataType

var _ Updater = (*MeasurementSeriesListDataType)(nil)

func (r *MeasurementSeriesListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []MeasurementSeriesDataType
	if newList != nil {
		newData = newList.(*MeasurementSeriesListDataType).MeasurementSeriesData
	}

	data, success := UpdateList(remoteWrite, r.MeasurementSeriesData, newData, filterPartial, filterDelete)

	if success && persist {
		r.MeasurementSeriesData = data
	}

	return data, success
}

// MeasurementConstraintsListDataType

var _ Updater = (*MeasurementConstraintsListDataType)(nil)

func (r *MeasurementConstraintsListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []MeasurementConstraintsDataType
	if newList != nil {
		newData = newList.(*MeasurementConstraintsListDataType).MeasurementConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.MeasurementConstraintsData, newData, filterPartial, filterDelete)

	if success && persist {
		r.MeasurementConstraintsData = data
	}

	return data, success
}

// MeasurementDescriptionListDataType

var _ Updater = (*MeasurementDescriptionListDataType)(nil)

func (r *MeasurementDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []MeasurementDescriptionDataType
	if newList != nil {
		newData = newList.(*MeasurementDescriptionListDataType).MeasurementDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.MeasurementDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.MeasurementDescriptionData = data
	}

	return data, success
}

// MeasurementThresholdRelationListDataType

var _ Updater = (*MeasurementThresholdRelationListDataType)(nil)

func (r *MeasurementThresholdRelationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []MeasurementThresholdRelationDataType
	if newList != nil {
		newData = newList.(*MeasurementThresholdRelationListDataType).MeasurementThresholdRelationData
	}

	data, success := UpdateList(remoteWrite, r.MeasurementThresholdRelationData, newData, filterPartial, filterDelete)

	if success && persist {
		r.MeasurementThresholdRelationData = data
	}

	return data, success
}
