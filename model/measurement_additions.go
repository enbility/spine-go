package model

// MeasurementListDataType

var _ Updater = (*MeasurementListDataType)(nil)

func (r *MeasurementListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []MeasurementDataType
	if newList != nil {
		newData = newList.(*MeasurementListDataType).MeasurementData
	}

	r.MeasurementData = UpdateList(remoteWrite, r.MeasurementData, newData, filterPartial, filterDelete)
}

// MeasurementSeriesListDataType

var _ Updater = (*MeasurementSeriesListDataType)(nil)

func (r *MeasurementSeriesListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []MeasurementSeriesDataType
	if newList != nil {
		newData = newList.(*MeasurementSeriesListDataType).MeasurementSeriesData
	}

	r.MeasurementSeriesData = UpdateList(remoteWrite, r.MeasurementSeriesData, newData, filterPartial, filterDelete)
}

// MeasurementConstraintsListDataType

var _ Updater = (*MeasurementConstraintsListDataType)(nil)

func (r *MeasurementConstraintsListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []MeasurementConstraintsDataType
	if newList != nil {
		newData = newList.(*MeasurementConstraintsListDataType).MeasurementConstraintsData
	}

	r.MeasurementConstraintsData = UpdateList(remoteWrite, r.MeasurementConstraintsData, newData, filterPartial, filterDelete)
}

// MeasurementDescriptionListDataType

var _ Updater = (*MeasurementDescriptionListDataType)(nil)

func (r *MeasurementDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []MeasurementDescriptionDataType
	if newList != nil {
		newData = newList.(*MeasurementDescriptionListDataType).MeasurementDescriptionData
	}

	r.MeasurementDescriptionData = UpdateList(remoteWrite, r.MeasurementDescriptionData, newData, filterPartial, filterDelete)
}

// MeasurementThresholdRelationListDataType

var _ Updater = (*MeasurementThresholdRelationListDataType)(nil)

func (r *MeasurementThresholdRelationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []MeasurementThresholdRelationDataType
	if newList != nil {
		newData = newList.(*MeasurementThresholdRelationListDataType).MeasurementThresholdRelationData
	}

	r.MeasurementThresholdRelationData = UpdateList(remoteWrite, r.MeasurementThresholdRelationData, newData, filterPartial, filterDelete)
}
