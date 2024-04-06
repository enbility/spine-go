package model

// AlarmListDataType

var _ Updater = (*AlarmListDataType)(nil)

func (r *AlarmListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []AlarmDataType
	if newList != nil {
		newData = newList.(*AlarmListDataType).AlarmListData
	}

	data, success := UpdateList(remoteWrite, r.AlarmListData, newData, filterPartial, filterDelete)

	if success {
		r.AlarmListData = data
	}

	return success
}
