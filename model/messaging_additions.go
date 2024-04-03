package model

// MessagingListDataType

var _ Updater = (*MessagingListDataType)(nil)

func (r *MessagingListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []MessagingDataType
	if newList != nil {
		newData = newList.(*MessagingListDataType).MessagingData
	}

	data, success := UpdateList(remoteWrite, r.MessagingData, newData, filterPartial, filterDelete)

	if success {
		r.MessagingData = data
	}

	return success
}
