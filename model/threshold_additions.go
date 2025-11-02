package model

// ThresholdListDataType

var _ Updater = (*ThresholdListDataType)(nil)

func (r *ThresholdListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []ThresholdDataType
	if newList != nil {
		newData = newList.(*ThresholdListDataType).ThresholdData
	}

	data, success := UpdateList(remoteWrite, r.ThresholdData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.ThresholdData = data
	}

	return data, success
}

// ThresholdConstraintsListDataType

var _ Updater = (*ThresholdConstraintsListDataType)(nil)

func (r *ThresholdConstraintsListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []ThresholdConstraintsDataType
	if newList != nil {
		newData = newList.(*ThresholdConstraintsListDataType).ThresholdConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.ThresholdConstraintsData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.ThresholdConstraintsData = data
	}

	return data, success
}

// ThresholdDescriptionListDataType

var _ Updater = (*ThresholdDescriptionListDataType)(nil)

func (r *ThresholdDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []ThresholdDescriptionDataType
	if newList != nil {
		newData = newList.(*ThresholdDescriptionListDataType).ThresholdDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.ThresholdDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.ThresholdDescriptionData = data
	}

	return data, success
}

// ThresholdListDataType PartialReader implementation
var _ PartialReader = (*ThresholdListDataType)(nil)

func (r *ThresholdListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[ThresholdDataType, ThresholdListDataSelectorsType, ThresholdDataElementsType](r.ThresholdData, filter)
	if !success {
		return r, false
	}

	result := &ThresholdListDataType{
		ThresholdData: filteredItems,
	}
	return result, true
}

// ThresholdConstraintsListDataType PartialReader implementation
var _ PartialReader = (*ThresholdConstraintsListDataType)(nil)

func (r *ThresholdConstraintsListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[ThresholdConstraintsDataType, ThresholdConstraintsListDataSelectorsType, ThresholdConstraintsDataElementsType](r.ThresholdConstraintsData, filter)
	if !success {
		return r, false
	}

	result := &ThresholdConstraintsListDataType{
		ThresholdConstraintsData: filteredItems,
	}
	return result, true
}

// ThresholdDescriptionListDataType PartialReader implementation
var _ PartialReader = (*ThresholdDescriptionListDataType)(nil)

func (r *ThresholdDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[ThresholdDescriptionDataType, ThresholdDescriptionListDataSelectorsType, ThresholdDescriptionDataElementsType](r.ThresholdDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &ThresholdDescriptionListDataType{
		ThresholdDescriptionData: filteredItems,
	}
	return result, true
}
