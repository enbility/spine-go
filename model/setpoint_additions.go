package model

// SetpointListDataType

var _ Updater = (*SetpointListDataType)(nil)

func (r *SetpointListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []SetpointDataType
	if newList != nil {
		newData = newList.(*SetpointListDataType).SetpointData
	}

	data, success := UpdateList(remoteWrite, r.SetpointData, newData, filterPartial, filterDelete)

	if success && persist {
		r.SetpointData = data
	}

	return data, success
}

// SetpointDescriptionListDataType

var _ Updater = (*SetpointDescriptionListDataType)(nil)

func (r *SetpointDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []SetpointDescriptionDataType
	if newList != nil {
		newData = newList.(*SetpointDescriptionListDataType).SetpointDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.SetpointDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.SetpointDescriptionData = data
	}

	return data, success
}
