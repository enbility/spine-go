package model

// BillListDataType

var _ Updater = (*BillListDataType)(nil)

func (r *BillListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []BillDataType
	if newList != nil {
		newData = newList.(*BillListDataType).BillData
	}

	data, success := UpdateList(remoteWrite, r.BillData, newData, filterPartial, filterDelete)

	if success && persist {
		r.BillData = data
	}

	return data, success
}

// BillConstraintsListDataType

var _ Updater = (*BillConstraintsListDataType)(nil)

func (r *BillConstraintsListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []BillConstraintsDataType
	if newList != nil {
		newData = newList.(*BillConstraintsListDataType).BillConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.BillConstraintsData, newData, filterPartial, filterDelete)

	if success && persist {
		r.BillConstraintsData = data
	}

	return data, success
}

// BillDescriptionListDataType

var _ Updater = (*BillDescriptionListDataType)(nil)

func (r *BillDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []BillDescriptionDataType
	if newList != nil {
		newData = newList.(*BillDescriptionListDataType).BillDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.BillDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.BillDescriptionData = data
	}

	return data, success
}
