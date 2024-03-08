package model

// MessagingListDataType

var _ Updater = (*MessagingListDataType)(nil)

func (r *MessagingListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []MessagingDataType
	if newList != nil {
		newData = newList.(*MessagingListDataType).MessagingData
	}

	r.MessagingData = UpdateList(remoteWrite, r.MessagingData, newData, filterPartial, filterDelete)
}
