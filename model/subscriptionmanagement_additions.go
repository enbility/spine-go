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

// SubscriptionManagementEntryListDataType PartialReader implementation
var _ PartialReader = (*SubscriptionManagementEntryListDataType)(nil)

func (r *SubscriptionManagementEntryListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[SubscriptionManagementEntryDataType, SubscriptionManagementEntryListDataSelectorsType, SubscriptionManagementEntryDataElementsType](r.SubscriptionManagementEntryData, filter)
	if !success {
		return r, false
	}

	result := &SubscriptionManagementEntryListDataType{
		SubscriptionManagementEntryData: filteredItems,
	}
	return result, true
}
