package model

// SetpointListDataType

var _ Updater = (*SetpointListDataType)(nil)

func (r *SetpointListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []SetpointDataType
	if newList != nil {
		newData = newList.(*SetpointListDataType).SetpointData
	}

	r.SetpointData = UpdateList(remoteWrite, r.SetpointData, newData, filterPartial, filterDelete)
}

// SetpointDescriptionListDataType

var _ Updater = (*SetpointDescriptionListDataType)(nil)

func (r *SetpointDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []SetpointDescriptionDataType
	if newList != nil {
		newData = newList.(*SetpointDescriptionListDataType).SetpointDescriptionData
	}

	r.SetpointDescriptionData = UpdateList(remoteWrite, r.SetpointDescriptionData, newData, filterPartial, filterDelete)
}
