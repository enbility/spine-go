package model

// AlarmListDataType

var _ Updater = (*AlarmListDataType)(nil)

func (r *AlarmListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []AlarmDataType
	if newList != nil {
		newData = newList.(*AlarmListDataType).AlarmListData
	}

	data, success := UpdateList(remoteWrite, r.AlarmListData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.AlarmListData = data
	}

	return data, success
}

// AlarmListDataType PartialReader implementation
var _ PartialReader = (*AlarmListDataType)(nil)

func (r *AlarmListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[AlarmDataType, AlarmListDataSelectorsType, AlarmDataElementsType](r.AlarmListData, filter)
	if !success {
		return r, false
	}

	result := &AlarmListDataType{
		AlarmListData: filteredItems,
	}
	return result, true
}
