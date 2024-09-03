package model

// SpecificationVersionListDataType

var _ Updater = (*SpecificationVersionListDataType)(nil)

func (r *SpecificationVersionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []SpecificationVersionDataType
	if newList != nil {
		newData = newList.(*SpecificationVersionListDataType).SpecificationVersionData
	}

	data, success := UpdateList(remoteWrite, r.SpecificationVersionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.SpecificationVersionData = data
	}

	return data, success
}
