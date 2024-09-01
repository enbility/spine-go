package model

// SupplyConditionListDataType

var _ Updater = (*SupplyConditionListDataType)(nil)

func (r *SupplyConditionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []SupplyConditionDataType
	if newList != nil {
		newData = newList.(*SupplyConditionListDataType).SupplyConditionData
	}

	data, success := UpdateList(remoteWrite, r.SupplyConditionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.SupplyConditionData = data
	}

	return data, success
}

// SupplyConditionDescriptionListDataType

var _ Updater = (*SupplyConditionDescriptionListDataType)(nil)

func (r *SupplyConditionDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []SupplyConditionDescriptionDataType
	if newList != nil {
		newData = newList.(*SupplyConditionDescriptionListDataType).SupplyConditionDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.SupplyConditionDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.SupplyConditionDescriptionData = data
	}

	return data, success
}

// SupplyConditionThresholdRelationListDataType

var _ Updater = (*SupplyConditionThresholdRelationListDataType)(nil)

func (r *SupplyConditionThresholdRelationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []SupplyConditionThresholdRelationDataType
	if newList != nil {
		newData = newList.(*SupplyConditionThresholdRelationListDataType).SupplyConditionThresholdRelationData
	}

	data, success := UpdateList(remoteWrite, r.SupplyConditionThresholdRelationData, newData, filterPartial, filterDelete)

	if success && persist {
		r.SupplyConditionThresholdRelationData = data
	}

	return data, success
}
