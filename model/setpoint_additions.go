package model

// SetpointListDataType

var _ Updater = (*SetpointListDataType)(nil)

func (r *SetpointListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []SetpointDataType
	if newList != nil {
		newData = newList.(*SetpointListDataType).SetpointData
	}

	data, success := UpdateList(remoteWrite, r.SetpointData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.SetpointData = data
	}

	return data, success
}

// SetpointDescriptionListDataType

var _ Updater = (*SetpointDescriptionListDataType)(nil)

func (r *SetpointDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []SetpointDescriptionDataType
	if newList != nil {
		newData = newList.(*SetpointDescriptionListDataType).SetpointDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.SetpointDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.SetpointDescriptionData = data
	}

	return data, success
}

// SetpointListDataType PartialReader implementation
var _ PartialReader = (*SetpointListDataType)(nil)

func (r *SetpointListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[SetpointDataType, SetpointListDataSelectorsType, SetpointDataElementsType](r.SetpointData, filter)
	if !success {
		return r, false
	}

	result := &SetpointListDataType{
		SetpointData: filteredItems,
	}
	return result, true
}

// SetpointDescriptionListDataType PartialReader implementation
var _ PartialReader = (*SetpointDescriptionListDataType)(nil)

func (r *SetpointDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[SetpointDescriptionDataType, SetpointDescriptionListDataSelectorsType, SetpointDescriptionDataElementsType](r.SetpointDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &SetpointDescriptionListDataType{
		SetpointDescriptionData: filteredItems,
	}
	return result, true
}
