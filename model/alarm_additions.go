package model

// AlarmListDataType

var _ Updater = (*AlarmListDataType)(nil)

func (r *AlarmListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []AlarmDataType
	if newList != nil {
		newData = newList.(*AlarmListDataType).AlarmListData
	}

	data, success := UpdateList(remoteWrite, r.AlarmListData, newData, filterPartial, filterDelete)

	if success && persist {
		r.AlarmListData = data
	}

	return data, success
}
