package model

// HvacSystemFunctionListDataType

var _ Updater = (*HvacSystemFunctionListDataType)(nil)

func (r *HvacSystemFunctionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []HvacSystemFunctionDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionListDataType).HvacSystemFunctionData
	}

	data, success := UpdateList(remoteWrite, r.HvacSystemFunctionData, newData, filterPartial, filterDelete)

	if success {
		r.HvacSystemFunctionData = data
	}

	return success
}

// HvacSystemFunctionOperationModeRelationListDataType

var _ Updater = (*HvacSystemFunctionOperationModeRelationListDataType)(nil)

func (r *HvacSystemFunctionOperationModeRelationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []HvacSystemFunctionOperationModeRelationDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionOperationModeRelationListDataType).HvacSystemFunctionOperationModeRelationData
	}

	data, success := UpdateList(remoteWrite, r.HvacSystemFunctionOperationModeRelationData, newData, filterPartial, filterDelete)

	if success {
		r.HvacSystemFunctionOperationModeRelationData = data
	}

	return success
}

// HvacSystemFunctionSetpointRelationListDataType

var _ Updater = (*HvacSystemFunctionSetpointRelationListDataType)(nil)

func (r *HvacSystemFunctionSetpointRelationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []HvacSystemFunctionSetpointRelationDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionSetpointRelationListDataType).HvacSystemFunctionSetpointRelationData
	}

	data, success := UpdateList(remoteWrite, r.HvacSystemFunctionSetpointRelationData, newData, filterPartial, filterDelete)

	if success {
		r.HvacSystemFunctionSetpointRelationData = data
	}

	return success
}

// HvacSystemFunctionPowerSequenceRelationListDataType

var _ Updater = (*HvacSystemFunctionPowerSequenceRelationListDataType)(nil)

func (r *HvacSystemFunctionPowerSequenceRelationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []HvacSystemFunctionPowerSequenceRelationDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionPowerSequenceRelationListDataType).HvacSystemFunctionPowerSequenceRelationData
	}

	data, success := UpdateList(remoteWrite, r.HvacSystemFunctionPowerSequenceRelationData, newData, filterPartial, filterDelete)

	if success {
		r.HvacSystemFunctionPowerSequenceRelationData = data
	}

	return success
}

// HvacSystemFunctionDescriptionListDataType

var _ Updater = (*HvacSystemFunctionDescriptionListDataType)(nil)

func (r *HvacSystemFunctionDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []HvacSystemFunctionDescriptionDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionDescriptionListDataType).HvacSystemFunctionDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.HvacSystemFunctionDescriptionData, newData, filterPartial, filterDelete)

	if success {
		r.HvacSystemFunctionDescriptionData = data
	}

	return success
}

// HvacOperationModeDescriptionListDataType

var _ Updater = (*HvacOperationModeDescriptionListDataType)(nil)

func (r *HvacOperationModeDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []HvacOperationModeDescriptionDataType
	if newList != nil {
		newData = newList.(*HvacOperationModeDescriptionListDataType).HvacOperationModeDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.HvacOperationModeDescriptionData, newData, filterPartial, filterDelete)

	if success {
		r.HvacOperationModeDescriptionData = data
	}

	return success
}

// HvacOverrunListDataType

var _ Updater = (*HvacOverrunListDataType)(nil)

func (r *HvacOverrunListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []HvacOverrunDataType
	if newList != nil {
		newData = newList.(*HvacOverrunListDataType).HvacOverrunData
	}

	data, success := UpdateList(remoteWrite, r.HvacOverrunData, newData, filterPartial, filterDelete)

	if success {
		r.HvacOverrunData = data
	}

	return success
}

// HvacOverrunDescriptionListDataType

var _ Updater = (*HvacOverrunDescriptionListDataType)(nil)

func (r *HvacOverrunDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []HvacOverrunDescriptionDataType
	if newList != nil {
		newData = newList.(*HvacOverrunDescriptionListDataType).HvacOverrunDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.HvacOverrunDescriptionData, newData, filterPartial, filterDelete)

	if success {
		r.HvacOverrunDescriptionData = data
	}

	return success
}
