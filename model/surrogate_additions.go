package model

// SurrogateDescriptionListDataType

var _ Updater = (*SurrogateDescriptionListDataType)(nil)

func (r *SurrogateDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []SurrogateDescriptionDataType
	if newList != nil {
		newData = newList.(*SurrogateDescriptionListDataType).SurrogateDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.SurrogateDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.SurrogateDescriptionData = data
	}

	return data, success
}
