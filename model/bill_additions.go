package model

// BillListDataType

var _ Updater = (*BillListDataType)(nil)

func (r *BillListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []BillDataType
	if newList != nil {
		newData = newList.(*BillListDataType).BillData
	}

	data, success := UpdateList(remoteWrite, r.BillData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.BillData = data
	}

	return data, success
}

// BillConstraintsListDataType

var _ Updater = (*BillConstraintsListDataType)(nil)

func (r *BillConstraintsListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []BillConstraintsDataType
	if newList != nil {
		newData = newList.(*BillConstraintsListDataType).BillConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.BillConstraintsData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.BillConstraintsData = data
	}

	return data, success
}

// BillDescriptionListDataType

var _ Updater = (*BillDescriptionListDataType)(nil)

func (r *BillDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []BillDescriptionDataType
	if newList != nil {
		newData = newList.(*BillDescriptionListDataType).BillDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.BillDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.BillDescriptionData = data
	}

	return data, success
}

// BillListDataType PartialReader implementation
var _ PartialReader = (*BillListDataType)(nil)

func (r *BillListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[BillDataType, BillListDataSelectorsType, BillDataElementsType](r.BillData, filter)
	if !success {
		return r, false
	}

	result := &BillListDataType{
		BillData: filteredItems,
	}
	return result, true
}

// BillConstraintsListDataType PartialReader implementation
var _ PartialReader = (*BillConstraintsListDataType)(nil)

func (r *BillConstraintsListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[BillConstraintsDataType, BillConstraintsListDataSelectorsType, BillConstraintsDataElementsType](r.BillConstraintsData, filter)
	if !success {
		return r, false
	}

	result := &BillConstraintsListDataType{
		BillConstraintsData: filteredItems,
	}
	return result, true
}

// BillDescriptionListDataType PartialReader implementation
var _ PartialReader = (*BillDescriptionListDataType)(nil)

func (r *BillDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[BillDescriptionDataType, BillDescriptionListDataSelectorsType, BillDescriptionDataElementsType](r.BillDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &BillDescriptionListDataType{
		BillDescriptionData: filteredItems,
	}
	return result, true
}
