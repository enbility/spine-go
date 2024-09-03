package model

// ThresholdListDataType

var _ Updater = (*ThresholdListDataType)(nil)

func (r *ThresholdListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []ThresholdDataType
	if newList != nil {
		newData = newList.(*ThresholdListDataType).ThresholdData
	}

	data, success := UpdateList(remoteWrite, r.ThresholdData, newData, filterPartial, filterDelete)

	if success && persist {
		r.ThresholdData = data
	}

	return data, success
}

// ThresholdConstraintsListDataType

var _ Updater = (*ThresholdConstraintsListDataType)(nil)

func (r *ThresholdConstraintsListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []ThresholdConstraintsDataType
	if newList != nil {
		newData = newList.(*ThresholdConstraintsListDataType).ThresholdConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.ThresholdConstraintsData, newData, filterPartial, filterDelete)

	if success && persist {
		r.ThresholdConstraintsData = data
	}

	return data, success
}

// ThresholdDescriptionListDataType

var _ Updater = (*ThresholdDescriptionListDataType)(nil)

func (r *ThresholdDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []ThresholdDescriptionDataType
	if newList != nil {
		newData = newList.(*ThresholdDescriptionListDataType).ThresholdDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.ThresholdDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.ThresholdDescriptionData = data
	}

	return data, success
}
