package model

// SubscriptionManagementEntryListDataType

var _ Updater = (*SubscriptionManagementEntryListDataType)(nil)

func (r *SubscriptionManagementEntryListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []SubscriptionManagementEntryDataType
	if newList != nil {
		newData = newList.(*SubscriptionManagementEntryListDataType).SubscriptionManagementEntryData
	}

	data, success := UpdateList(remoteWrite, r.SubscriptionManagementEntryData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.SubscriptionManagementEntryData = data
	}

	return data, success
}
