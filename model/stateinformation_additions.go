package model

// StateInformationListDataType

var _ Updater = (*StateInformationListDataType)(nil)

func (r *StateInformationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []StateInformationDataType
	if newList != nil {
		newData = newList.(*StateInformationListDataType).StateInformationData
	}

	data, success := UpdateList(remoteWrite, r.StateInformationData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.StateInformationData = data
	}

	return data, success
}

// StateInformationListDataType PartialReader implementation
var _ PartialReader = (*StateInformationListDataType)(nil)

func (r *StateInformationListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[StateInformationDataType, StateInformationListDataSelectorsType, StateInformationDataElementsType](r.StateInformationData, filter)
	if !success {
		return r, false
	}

	result := &StateInformationListDataType{
		StateInformationData: filteredItems,
	}
	return result, true
}
