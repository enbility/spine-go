package model

// StateInformationListDataType

var _ Updater = (*StateInformationListDataType)(nil)

func (r *StateInformationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []StateInformationDataType
	if newList != nil {
		newData = newList.(*StateInformationListDataType).StateInformationData
	}

	data, success := UpdateList(remoteWrite, r.StateInformationData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.StateInformationData = data
	}

	return data, success
}
