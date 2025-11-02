package model

// HvacSystemFunctionListDataType

var _ Updater = (*HvacSystemFunctionListDataType)(nil)

func (r *HvacSystemFunctionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []HvacSystemFunctionDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionListDataType).HvacSystemFunctionData
	}

	data, success := UpdateList(remoteWrite, r.HvacSystemFunctionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.HvacSystemFunctionData = data
	}

	return data, success
}

// HvacSystemFunctionOperationModeRelationListDataType

var _ Updater = (*HvacSystemFunctionOperationModeRelationListDataType)(nil)

func (r *HvacSystemFunctionOperationModeRelationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []HvacSystemFunctionOperationModeRelationDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionOperationModeRelationListDataType).HvacSystemFunctionOperationModeRelationData
	}

	data, success := UpdateList(remoteWrite, r.HvacSystemFunctionOperationModeRelationData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.HvacSystemFunctionOperationModeRelationData = data
	}

	return data, success
}

// HvacSystemFunctionSetpointRelationListDataType

var _ Updater = (*HvacSystemFunctionSetpointRelationListDataType)(nil)

func (r *HvacSystemFunctionSetpointRelationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []HvacSystemFunctionSetpointRelationDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionSetpointRelationListDataType).HvacSystemFunctionSetpointRelationData
	}

	data, success := UpdateList(remoteWrite, r.HvacSystemFunctionSetpointRelationData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.HvacSystemFunctionSetpointRelationData = data
	}

	return data, success
}

// HvacSystemFunctionPowerSequenceRelationListDataType

var _ Updater = (*HvacSystemFunctionPowerSequenceRelationListDataType)(nil)

func (r *HvacSystemFunctionPowerSequenceRelationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []HvacSystemFunctionPowerSequenceRelationDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionPowerSequenceRelationListDataType).HvacSystemFunctionPowerSequenceRelationData
	}

	data, success := UpdateList(remoteWrite, r.HvacSystemFunctionPowerSequenceRelationData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.HvacSystemFunctionPowerSequenceRelationData = data
	}

	return data, success
}

// HvacSystemFunctionDescriptionListDataType

var _ Updater = (*HvacSystemFunctionDescriptionListDataType)(nil)

func (r *HvacSystemFunctionDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []HvacSystemFunctionDescriptionDataType
	if newList != nil {
		newData = newList.(*HvacSystemFunctionDescriptionListDataType).HvacSystemFunctionDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.HvacSystemFunctionDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.HvacSystemFunctionDescriptionData = data
	}

	return data, success
}

// HvacOperationModeDescriptionListDataType

var _ Updater = (*HvacOperationModeDescriptionListDataType)(nil)

func (r *HvacOperationModeDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []HvacOperationModeDescriptionDataType
	if newList != nil {
		newData = newList.(*HvacOperationModeDescriptionListDataType).HvacOperationModeDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.HvacOperationModeDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.HvacOperationModeDescriptionData = data
	}

	return data, success
}

// HvacOverrunListDataType

var _ Updater = (*HvacOverrunListDataType)(nil)

func (r *HvacOverrunListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []HvacOverrunDataType
	if newList != nil {
		newData = newList.(*HvacOverrunListDataType).HvacOverrunData
	}

	data, success := UpdateList(remoteWrite, r.HvacOverrunData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.HvacOverrunData = data
	}

	return data, success
}

// HvacOverrunDescriptionListDataType

var _ Updater = (*HvacOverrunDescriptionListDataType)(nil)

func (r *HvacOverrunDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []HvacOverrunDescriptionDataType
	if newList != nil {
		newData = newList.(*HvacOverrunDescriptionListDataType).HvacOverrunDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.HvacOverrunDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.HvacOverrunDescriptionData = data
	}

	return data, success
}

// HvacOverrunDescriptionListDataType PartialReader implementation
var _ PartialReader = (*HvacOverrunDescriptionListDataType)(nil)

func (r *HvacOverrunDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[HvacOverrunDescriptionDataType, HvacOverrunDescriptionListDataSelectorsType, HvacOverrunDescriptionDataElementsType](r.HvacOverrunDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &HvacOverrunDescriptionListDataType{
		HvacOverrunDescriptionData: filteredItems,
	}
	return result, true
}

// HvacSystemFunctionListDataType PartialReader implementation
var _ PartialReader = (*HvacSystemFunctionListDataType)(nil)

func (r *HvacSystemFunctionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[HvacSystemFunctionDataType, HvacSystemFunctionListDataSelectorsType, HvacSystemFunctionDataElementsType](r.HvacSystemFunctionData, filter)
	if !success {
		return r, false
	}

	result := &HvacSystemFunctionListDataType{
		HvacSystemFunctionData: filteredItems,
	}
	return result, true
}

// HvacSystemFunctionOperationModeRelationListDataType PartialReader implementation
var _ PartialReader = (*HvacSystemFunctionOperationModeRelationListDataType)(nil)

func (r *HvacSystemFunctionOperationModeRelationListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[HvacSystemFunctionOperationModeRelationDataType, HvacSystemFunctionOperationModeRelationListDataSelectorsType, HvacSystemFunctionOperationModeRelationDataElementsType](r.HvacSystemFunctionOperationModeRelationData, filter)
	if !success {
		return r, false
	}

	result := &HvacSystemFunctionOperationModeRelationListDataType{
		HvacSystemFunctionOperationModeRelationData: filteredItems,
	}
	return result, true
}

// HvacSystemFunctionSetpointRelationListDataType PartialReader implementation
var _ PartialReader = (*HvacSystemFunctionSetpointRelationListDataType)(nil)

func (r *HvacSystemFunctionSetpointRelationListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[HvacSystemFunctionSetpointRelationDataType, HvacSystemFunctionSetpointRelationListDataSelectorsType, HvacSystemFunctionSetpointRelationDataElementsType](r.HvacSystemFunctionSetpointRelationData, filter)
	if !success {
		return r, false
	}

	result := &HvacSystemFunctionSetpointRelationListDataType{
		HvacSystemFunctionSetpointRelationData: filteredItems,
	}
	return result, true
}

// HvacSystemFunctionPowerSequenceRelationListDataType PartialReader implementation
var _ PartialReader = (*HvacSystemFunctionPowerSequenceRelationListDataType)(nil)

func (r *HvacSystemFunctionPowerSequenceRelationListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[HvacSystemFunctionPowerSequenceRelationDataType, HvacSystemFunctionPowerSequenceRelationListDataSelectorsType, HvacSystemFunctionPowerSequenceRelationDataElementsType](r.HvacSystemFunctionPowerSequenceRelationData, filter)
	if !success {
		return r, false
	}

	result := &HvacSystemFunctionPowerSequenceRelationListDataType{
		HvacSystemFunctionPowerSequenceRelationData: filteredItems,
	}
	return result, true
}

// HvacSystemFunctionDescriptionListDataType PartialReader implementation
var _ PartialReader = (*HvacSystemFunctionDescriptionListDataType)(nil)

func (r *HvacSystemFunctionDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[HvacSystemFunctionDescriptionDataType, HvacSystemFunctionDescriptionListDataSelectorsType, HvacSystemFunctionDescriptionDataElementsType](r.HvacSystemFunctionDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &HvacSystemFunctionDescriptionListDataType{
		HvacSystemFunctionDescriptionData: filteredItems,
	}
	return result, true
}

// HvacOperationModeDescriptionListDataType PartialReader implementation
var _ PartialReader = (*HvacOperationModeDescriptionListDataType)(nil)

func (r *HvacOperationModeDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[HvacOperationModeDescriptionDataType, HvacOperationModeDescriptionListDataSelectorsType, HvacOperationModeDescriptionDataElementsType](r.HvacOperationModeDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &HvacOperationModeDescriptionListDataType{
		HvacOperationModeDescriptionData: filteredItems,
	}
	return result, true
}

// HvacOverrunListDataType PartialReader implementation
var _ PartialReader = (*HvacOverrunListDataType)(nil)

func (r *HvacOverrunListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[HvacOverrunDataType, HvacOverrunListDataSelectorsType, HvacOverrunDataElementsType](r.HvacOverrunData, filter)
	if !success {
		return r, false
	}

	result := &HvacOverrunListDataType{
		HvacOverrunData: filteredItems,
	}
	return result, true
}
