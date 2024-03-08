package model

// BillListDataType

var _ Updater = (*BillListDataType)(nil)

func (r *BillListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []BillDataType
	if newList != nil {
		newData = newList.(*BillListDataType).BillData
	}

	r.BillData = UpdateList(remoteWrite, r.BillData, newData, filterPartial, filterDelete)
}

// BillConstraintsListDataType

var _ Updater = (*BillConstraintsListDataType)(nil)

func (r *BillConstraintsListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []BillConstraintsDataType
	if newList != nil {
		newData = newList.(*BillConstraintsListDataType).BillConstraintsData
	}

	r.BillConstraintsData = UpdateList(remoteWrite, r.BillConstraintsData, newData, filterPartial, filterDelete)
}

// BillDescriptionListDataType

var _ Updater = (*BillDescriptionListDataType)(nil)

func (r *BillDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []BillDescriptionDataType
	if newList != nil {
		newData = newList.(*BillDescriptionListDataType).BillDescriptionData
	}

	r.BillDescriptionData = UpdateList(remoteWrite, r.BillDescriptionData, newData, filterPartial, filterDelete)
}
