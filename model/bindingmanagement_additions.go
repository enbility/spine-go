package model

// BindingManagementEntryListDataType

var _ Updater = (*BindingManagementEntryListDataType)(nil)

func (r *BindingManagementEntryListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []BindingManagementEntryDataType
	if newList != nil {
		newData = newList.(*BindingManagementEntryListDataType).BindingManagementEntryData
	}

	data, success := UpdateList(remoteWrite, r.BindingManagementEntryData, newData, filterPartial, filterDelete)

	if success {
		r.BindingManagementEntryData = data
	}

	return success
}
