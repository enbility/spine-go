package model

// HvacSystemFunctionListDataType

var _ Updater = (*HvacSystemFunctionListDataType)(nil)

func (r *HvacSystemFunctionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []HvacSystemFunctionDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionListDataType).HvacSystemFunctionData
	}

	r.HvacSystemFunctionData = UpdateList(remoteWrite, r.HvacSystemFunctionData, newData, filterPartial, filterDelete)
}

// HvacSystemFunctionOperationModeRelationListDataType

var _ Updater = (*HvacSystemFunctionOperationModeRelationListDataType)(nil)

func (r *HvacSystemFunctionOperationModeRelationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []HvacSystemFunctionOperationModeRelationDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionOperationModeRelationListDataType).HvacSystemFunctionOperationModeRelationData
	}

	r.HvacSystemFunctionOperationModeRelationData = UpdateList(remoteWrite, r.HvacSystemFunctionOperationModeRelationData, newData, filterPartial, filterDelete)
}

// HvacSystemFunctionSetpointRelationListDataType

var _ Updater = (*HvacSystemFunctionSetpointRelationListDataType)(nil)

func (r *HvacSystemFunctionSetpointRelationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []HvacSystemFunctionSetpointRelationDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionSetpointRelationListDataType).HvacSystemFunctionSetpointRelationData
	}

	r.HvacSystemFunctionSetpointRelationData = UpdateList(remoteWrite, r.HvacSystemFunctionSetpointRelationData, newData, filterPartial, filterDelete)
}

// HvacSystemFunctionPowerSequenceRelationListDataType

var _ Updater = (*HvacSystemFunctionPowerSequenceRelationListDataType)(nil)

func (r *HvacSystemFunctionPowerSequenceRelationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []HvacSystemFunctionPowerSequenceRelationDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionPowerSequenceRelationListDataType).HvacSystemFunctionPowerSequenceRelationData
	}

	r.HvacSystemFunctionPowerSequenceRelationData = UpdateList(remoteWrite, r.HvacSystemFunctionPowerSequenceRelationData, newData, filterPartial, filterDelete)
}

// HvacSystemFunctionDescriptionListDataType

var _ Updater = (*HvacSystemFunctionDescriptionListDataType)(nil)

func (r *HvacSystemFunctionDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []HvacSystemFunctionDescriptionDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionDescriptionListDataType).HvacSystemFunctionDescriptionData
	}

	r.HvacSystemFunctionDescriptionData = UpdateList(remoteWrite, r.HvacSystemFunctionDescriptionData, newData, filterPartial, filterDelete)
}

// HvacOperationModeDescriptionListDataType

var _ Updater = (*HvacOperationModeDescriptionListDataType)(nil)

func (r *HvacOperationModeDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []HvacOperationModeDescriptionDataType
	if newList != nil {
		newData = newList.(*HvacOperationModeDescriptionListDataType).HvacOperationModeDescriptionData
	}

	r.HvacOperationModeDescriptionData = UpdateList(remoteWrite, r.HvacOperationModeDescriptionData, newData, filterPartial, filterDelete)
}

// HvacOverrunListDataType

var _ Updater = (*HvacOverrunListDataType)(nil)

func (r *HvacOverrunListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []HvacOverrunDataType
	if newList != nil {
		newData = newList.(*HvacOverrunListDataType).HvacOverrunData
	}

	r.HvacOverrunData = UpdateList(remoteWrite, r.HvacOverrunData, newData, filterPartial, filterDelete)
}

// HvacOverrunDescriptionListDataType

var _ Updater = (*HvacOverrunDescriptionListDataType)(nil)

func (r *HvacOverrunDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []HvacOverrunDescriptionDataType
	if newList != nil {
		newData = newList.(*HvacOverrunDescriptionListDataType).HvacOverrunDescriptionData
	}

	r.HvacOverrunDescriptionData = UpdateList(remoteWrite, r.HvacOverrunDescriptionData, newData, filterPartial, filterDelete)
}
