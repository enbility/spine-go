package model

// HvacSystemFunctionListDataType

var _ Updater = (*HvacSystemFunctionListDataType)(nil)

func (r *HvacSystemFunctionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []HvacSystemFunctionDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionListDataType).HvacSystemFunctionData
	}

	data, success := UpdateList(remoteWrite, r.HvacSystemFunctionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.HvacSystemFunctionData = data
	}

	return data, success
}

// HvacSystemFunctionOperationModeRelationListDataType

var _ Updater = (*HvacSystemFunctionOperationModeRelationListDataType)(nil)

func (r *HvacSystemFunctionOperationModeRelationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []HvacSystemFunctionOperationModeRelationDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionOperationModeRelationListDataType).HvacSystemFunctionOperationModeRelationData
	}

	data, success := UpdateList(remoteWrite, r.HvacSystemFunctionOperationModeRelationData, newData, filterPartial, filterDelete)

	if success && persist {
		r.HvacSystemFunctionOperationModeRelationData = data
	}

	return data, success
}

// HvacSystemFunctionSetpointRelationListDataType

var _ Updater = (*HvacSystemFunctionSetpointRelationListDataType)(nil)

func (r *HvacSystemFunctionSetpointRelationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []HvacSystemFunctionSetpointRelationDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionSetpointRelationListDataType).HvacSystemFunctionSetpointRelationData
	}

	data, success := UpdateList(remoteWrite, r.HvacSystemFunctionSetpointRelationData, newData, filterPartial, filterDelete)

	if success && persist {
		r.HvacSystemFunctionSetpointRelationData = data
	}

	return data, success
}

// HvacSystemFunctionPowerSequenceRelationListDataType

var _ Updater = (*HvacSystemFunctionPowerSequenceRelationListDataType)(nil)

func (r *HvacSystemFunctionPowerSequenceRelationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []HvacSystemFunctionPowerSequenceRelationDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionPowerSequenceRelationListDataType).HvacSystemFunctionPowerSequenceRelationData
	}

	data, success := UpdateList(remoteWrite, r.HvacSystemFunctionPowerSequenceRelationData, newData, filterPartial, filterDelete)

	if success && persist {
		r.HvacSystemFunctionPowerSequenceRelationData = data
	}

	return data, success
}

// HvacSystemFunctionDescriptionListDataType

var _ Updater = (*HvacSystemFunctionDescriptionListDataType)(nil)

func (r *HvacSystemFunctionDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []HvacSystemFunctionDescriptionDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionDescriptionListDataType).HvacSystemFunctionDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.HvacSystemFunctionDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.HvacSystemFunctionDescriptionData = data
	}

	return data, success
}

// HvacOperationModeDescriptionListDataType

var _ Updater = (*HvacOperationModeDescriptionListDataType)(nil)

func (r *HvacOperationModeDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []HvacOperationModeDescriptionDataType
	if newList != nil {
		newData = newList.(*HvacOperationModeDescriptionListDataType).HvacOperationModeDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.HvacOperationModeDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.HvacOperationModeDescriptionData = data
	}

	return data, success
}

// HvacOverrunListDataType

var _ Updater = (*HvacOverrunListDataType)(nil)

func (r *HvacOverrunListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []HvacOverrunDataType
	if newList != nil {
		newData = newList.(*HvacOverrunListDataType).HvacOverrunData
	}

	data, success := UpdateList(remoteWrite, r.HvacOverrunData, newData, filterPartial, filterDelete)

	if success && persist {
		r.HvacOverrunData = data
	}

	return data, success
}

// HvacOverrunDescriptionListDataType

var _ Updater = (*HvacOverrunDescriptionListDataType)(nil)

func (r *HvacOverrunDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []HvacOverrunDescriptionDataType
	if newList != nil {
		newData = newList.(*HvacOverrunDescriptionListDataType).HvacOverrunDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.HvacOverrunDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.HvacOverrunDescriptionData = data
	}

	return data, success
}
