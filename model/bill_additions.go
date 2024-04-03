package model

// BillListDataType

var _ Updater = (*BillListDataType)(nil)

func (r *BillListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []BillDataType
	if newList != nil {
		newData = newList.(*BillListDataType).BillData
	}

	data, success := UpdateList(remoteWrite, r.BillData, newData, filterPartial, filterDelete)

	if success {
		r.BillData = data
	}

	return success
}

// BillConstraintsListDataType

var _ Updater = (*BillConstraintsListDataType)(nil)

func (r *BillConstraintsListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []BillConstraintsDataType
	if newList != nil {
		newData = newList.(*BillConstraintsListDataType).BillConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.BillConstraintsData, newData, filterPartial, filterDelete)

	if success {
		r.BillConstraintsData = data
	}

	return success
}

// BillDescriptionListDataType

var _ Updater = (*BillDescriptionListDataType)(nil)

func (r *BillDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []BillDescriptionDataType
	if newList != nil {
		newData = newList.(*BillDescriptionListDataType).BillDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.BillDescriptionData, newData, filterPartial, filterDelete)

	if success {
		r.BillDescriptionData = data
	}

	return success
}
