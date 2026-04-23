package model

// AlarmListDataType

var _ Updater = (*AlarmListDataType)(nil)

func (r *AlarmListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []AlarmDataType
	if newList != nil {
		newData = newList.(*AlarmListDataType).AlarmData
	}

	data, success := UpdateList(remoteWrite, r.AlarmData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.AlarmData = data
	}

	return data, success
}
