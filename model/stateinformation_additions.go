package model

// StateInformationListDataType

var _ Updater = (*StateInformationListDataType)(nil)

func (r *StateInformationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []StateInformationDataType
	if newList != nil {
		newData = newList.(*StateInformationListDataType).StateInformationData
	}

	data, success := UpdateList(remoteWrite, r.StateInformationData, newData, filterPartial, filterDelete)

	if success && persist {
		r.StateInformationData = data
	}

	return data, success
}
