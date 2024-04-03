package model

// MeasurementListDataType

var _ Updater = (*MeasurementListDataType)(nil)

func (r *MeasurementListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []MeasurementDataType
	if newList != nil {
		newData = newList.(*MeasurementListDataType).MeasurementData
	}

	data, success := UpdateList(remoteWrite, r.MeasurementData, newData, filterPartial, filterDelete)

	if success {
		r.MeasurementData = data
	}

	return success
}

// MeasurementSeriesListDataType

var _ Updater = (*MeasurementSeriesListDataType)(nil)

func (r *MeasurementSeriesListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []MeasurementSeriesDataType
	if newList != nil {
		newData = newList.(*MeasurementSeriesListDataType).MeasurementSeriesData
	}

	data, success := UpdateList(remoteWrite, r.MeasurementSeriesData, newData, filterPartial, filterDelete)

	if success {
		r.MeasurementSeriesData = data
	}

	return success
}

// MeasurementConstraintsListDataType

var _ Updater = (*MeasurementConstraintsListDataType)(nil)

func (r *MeasurementConstraintsListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []MeasurementConstraintsDataType
	if newList != nil {
		newData = newList.(*MeasurementConstraintsListDataType).MeasurementConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.MeasurementConstraintsData, newData, filterPartial, filterDelete)

	if success {
		r.MeasurementConstraintsData = data
	}

	return success
}

// MeasurementDescriptionListDataType

var _ Updater = (*MeasurementDescriptionListDataType)(nil)

func (r *MeasurementDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []MeasurementDescriptionDataType
	if newList != nil {
		newData = newList.(*MeasurementDescriptionListDataType).MeasurementDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.MeasurementDescriptionData, newData, filterPartial, filterDelete)

	if success {
		r.MeasurementDescriptionData = data
	}

	return success
}

// MeasurementThresholdRelationListDataType

var _ Updater = (*MeasurementThresholdRelationListDataType)(nil)

func (r *MeasurementThresholdRelationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []MeasurementThresholdRelationDataType
	if newList != nil {
		newData = newList.(*MeasurementThresholdRelationListDataType).MeasurementThresholdRelationData
	}

	data, success := UpdateList(remoteWrite, r.MeasurementThresholdRelationData, newData, filterPartial, filterDelete)

	if success {
		r.MeasurementThresholdRelationData = data
	}

	return success
}
