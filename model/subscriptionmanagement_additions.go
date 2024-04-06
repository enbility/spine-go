package model

// SubscriptionManagementEntryListDataType

var _ Updater = (*SubscriptionManagementEntryListDataType)(nil)

func (r *SubscriptionManagementEntryListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []SubscriptionManagementEntryDataType
	if newList != nil {
		newData = newList.(*SubscriptionManagementEntryListDataType).SubscriptionManagementEntryData
	}

	data, success := UpdateList(remoteWrite, r.SubscriptionManagementEntryData, newData, filterPartial, filterDelete)

	if success {
		r.SubscriptionManagementEntryData = data
	}

	return success
}
