package model

// MessagingListDataType

var _ Updater = (*MessagingListDataType)(nil)

func (r *MessagingListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []MessagingDataType
	if newList != nil {
		newData = newList.(*MessagingListDataType).MessagingData
	}

	data, success := UpdateList(remoteWrite, r.MessagingData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.MessagingData = data
	}

	return data, success
}

// MessagingListDataType PartialReader implementation
var _ PartialReader = (*MessagingListDataType)(nil)

func (r *MessagingListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[MessagingDataType, MessagingListDataSelectorsType, MessagingDataElementsType](r.MessagingData, filter)
	if !success {
		return r, false
	}

	result := &MessagingListDataType{
		MessagingData: filteredItems,
	}
	return result, true
}
