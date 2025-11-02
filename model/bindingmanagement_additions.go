package model

// BindingManagementEntryListDataType

var _ Updater = (*BindingManagementEntryListDataType)(nil)

func (r *BindingManagementEntryListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []BindingManagementEntryDataType
	if newList != nil {
		newData = newList.(*BindingManagementEntryListDataType).BindingManagementEntryData
	}

	data, success := UpdateList(remoteWrite, r.BindingManagementEntryData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.BindingManagementEntryData = data
	}

	return data, success
}

// BindingManagementEntryListDataType PartialReader implementation
var _ PartialReader = (*BindingManagementEntryListDataType)(nil)

func (r *BindingManagementEntryListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[BindingManagementEntryDataType, BindingManagementEntryListDataSelectorsType, BindingManagementEntryDataElementsType](r.BindingManagementEntryData, filter)
	if !success {
		return r, false
	}

	result := &BindingManagementEntryListDataType{
		BindingManagementEntryData: filteredItems,
	}
	return result, true
}
