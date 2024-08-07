package model

// StateInformationListDataType

var _ Updater = (*StateInformationListDataType)(nil)

func (r *StateInformationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []StateInformationDataType
	if newList != nil {
		newData = newList.(*StateInformationListDataType).StateInformationData
	}

	data, success := UpdateList(remoteWrite, r.StateInformationData, newData, filterPartial, filterDelete)

	if success {
		r.StateInformationData = data
	}

	return success
}
