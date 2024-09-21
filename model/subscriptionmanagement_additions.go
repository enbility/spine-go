package model

// SubscriptionManagementEntryListDataType

var _ Updater = (*SubscriptionManagementEntryListDataType)(nil)

func (r *SubscriptionManagementEntryListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []SubscriptionManagementEntryDataType
	if newList != nil {
		newData = newList.(*SubscriptionManagementEntryListDataType).SubscriptionManagementEntryData
	}

	data, success := UpdateList(remoteWrite, r.SubscriptionManagementEntryData, newData, filterPartial, filterDelete)

	if success && persist {
		r.SubscriptionManagementEntryData = data
	}

	return data, success
}
