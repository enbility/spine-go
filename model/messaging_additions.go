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
