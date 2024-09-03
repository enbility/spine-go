package model

// BindingManagementEntryListDataType

var _ Updater = (*BindingManagementEntryListDataType)(nil)

func (r *BindingManagementEntryListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []BindingManagementEntryDataType
	if newList != nil {
		newData = newList.(*BindingManagementEntryListDataType).BindingManagementEntryData
	}

	data, success := UpdateList(remoteWrite, r.BindingManagementEntryData, newData, filterPartial, filterDelete)

	if success && persist {
		r.BindingManagementEntryData = data
	}

	return data, success
}
