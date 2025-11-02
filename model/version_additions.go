package model

// SpecificationVersionListDataType

var _ Updater = (*SpecificationVersionListDataType)(nil)

func (r *SpecificationVersionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []SpecificationVersionDataType
	if newList != nil {
		newData = newList.(*SpecificationVersionListDataType).SpecificationVersionData
	}

	data, success := UpdateList(remoteWrite, r.SpecificationVersionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.SpecificationVersionData = data
	}

	return data, success
}

// SpecificationVersionListDataType PartialReader implementation
var _ PartialReader = (*SpecificationVersionListDataType)(nil)

func (r *SpecificationVersionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	// SpecificationVersionListDataSelectorsType is empty, so we always return all data
	return r, true
}
