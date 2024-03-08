package model

// BindingManagementEntryListDataType

var _ Updater = (*BindingManagementEntryListDataType)(nil)

func (r *BindingManagementEntryListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []BindingManagementEntryDataType
	if newList != nil {
		newData = newList.(*BindingManagementEntryListDataType).BindingManagementEntryData
	}

	r.BindingManagementEntryData = UpdateList(remoteWrite, r.BindingManagementEntryData, newData, filterPartial, filterDelete)
}
