package model

// SetpointListDataType

var _ Updater = (*SetpointListDataType)(nil)

func (r *SetpointListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []SetpointDataType
	if newList != nil {
		newData = newList.(*SetpointListDataType).SetpointData
	}

	data, success := UpdateList(remoteWrite, r.SetpointData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.SetpointData = data
	}

	return data, success
}

// SetpointConstraintsListDataType

var _ Updater = (*SetpointConstraintsListDataType)(nil)

func (r *SetpointConstraintsListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []SetpointConstraintsDataType
	if newList != nil {
		newData = newList.(*SetpointConstraintsListDataType).SetpointConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.SetpointConstraintsData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.SetpointConstraintsData = data
	}

	return data, success
}

// SetpointDescriptionListDataType

var _ Updater = (*SetpointDescriptionListDataType)(nil)

func (r *SetpointDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []SetpointDescriptionDataType
	if newList != nil {
		newData = newList.(*SetpointDescriptionListDataType).SetpointDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.SetpointDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.SetpointDescriptionData = data
	}

	return data, success
}
