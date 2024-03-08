package model

// SpecificationVersionListDataType

var _ Updater = (*SpecificationVersionListDataType)(nil)

func (r *SpecificationVersionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []SpecificationVersionDataType
	if newList != nil {
		newData = newList.(*SpecificationVersionListDataType).SpecificationVersionData
	}

	r.SpecificationVersionData = UpdateList(remoteWrite, r.SpecificationVersionData, newData, filterPartial, filterDelete)
}
