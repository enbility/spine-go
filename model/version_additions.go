package model

// SpecificationVersionListDataType

var _ Updater = (*SpecificationVersionListDataType)(nil)

func (r *SpecificationVersionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []SpecificationVersionDataType
	if newList != nil {
		newData = newList.(*SpecificationVersionListDataType).SpecificationVersionData
	}

	data, success := UpdateList(remoteWrite, r.SpecificationVersionData, newData, filterPartial, filterDelete)

	if success {
		r.SpecificationVersionData = data
	}

	return success
}
