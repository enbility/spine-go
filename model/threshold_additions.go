package model

// ThresholdListDataType

var _ Updater = (*ThresholdListDataType)(nil)

func (r *ThresholdListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []ThresholdDataType
	if newList != nil {
		newData = newList.(*ThresholdListDataType).ThresholdData
	}

	r.ThresholdData = UpdateList(remoteWrite, r.ThresholdData, newData, filterPartial, filterDelete)
}

// ThresholdConstraintsListDataType

var _ Updater = (*ThresholdConstraintsListDataType)(nil)

func (r *ThresholdConstraintsListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []ThresholdConstraintsDataType
	if newList != nil {
		newData = newList.(*ThresholdConstraintsListDataType).ThresholdConstraintsData
	}

	r.ThresholdConstraintsData = UpdateList(remoteWrite, r.ThresholdConstraintsData, newData, filterPartial, filterDelete)
}

// ThresholdDescriptionListDataType

var _ Updater = (*ThresholdDescriptionListDataType)(nil)

func (r *ThresholdDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []ThresholdDescriptionDataType
	if newList != nil {
		newData = newList.(*ThresholdDescriptionListDataType).ThresholdDescriptionData
	}

	r.ThresholdDescriptionData = UpdateList(remoteWrite, r.ThresholdDescriptionData, newData, filterPartial, filterDelete)
}
