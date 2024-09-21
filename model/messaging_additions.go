package model

// MessagingListDataType

var _ Updater = (*MessagingListDataType)(nil)

func (r *MessagingListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []MessagingDataType
	if newList != nil {
		newData = newList.(*MessagingListDataType).MessagingData
	}

	data, success := UpdateList(remoteWrite, r.MessagingData, newData, filterPartial, filterDelete)

	if success && persist {
		r.MessagingData = data
	}

	return data, success
}
