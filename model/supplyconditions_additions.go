package model

// SupplyConditionListDataType

var _ Updater = (*SupplyConditionListDataType)(nil)

func (r *SupplyConditionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []SupplyConditionDataType
	if newList != nil {
		newData = newList.(*SupplyConditionListDataType).SupplyConditionData
	}

	data, success := UpdateList(remoteWrite, r.SupplyConditionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.SupplyConditionData = data
	}

	return data, success
}

// SupplyConditionDescriptionListDataType

var _ Updater = (*SupplyConditionDescriptionListDataType)(nil)

func (r *SupplyConditionDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []SupplyConditionDescriptionDataType
	if newList != nil {
		newData = newList.(*SupplyConditionDescriptionListDataType).SupplyConditionDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.SupplyConditionDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.SupplyConditionDescriptionData = data
	}

	return data, success
}

// SupplyConditionThresholdRelationListDataType

var _ Updater = (*SupplyConditionThresholdRelationListDataType)(nil)

func (r *SupplyConditionThresholdRelationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []SupplyConditionThresholdRelationDataType
	if newList != nil {
		newData = newList.(*SupplyConditionThresholdRelationListDataType).SupplyConditionThresholdRelationData
	}

	data, success := UpdateList(remoteWrite, r.SupplyConditionThresholdRelationData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.SupplyConditionThresholdRelationData = data
	}

	return data, success
}

// SupplyConditionListDataType PartialReader implementation
var _ PartialReader = (*SupplyConditionListDataType)(nil)

func (r *SupplyConditionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[SupplyConditionDataType, SupplyConditionListDataSelectorsType, SupplyConditionDataElementsType](r.SupplyConditionData, filter)
	if !success {
		return r, false
	}

	result := &SupplyConditionListDataType{
		SupplyConditionData: filteredItems,
	}
	return result, true
}

// SupplyConditionDescriptionListDataType PartialReader implementation
var _ PartialReader = (*SupplyConditionDescriptionListDataType)(nil)

func (r *SupplyConditionDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[SupplyConditionDescriptionDataType, SupplyConditionDescriptionListDataSelectorsType, SupplyConditionDescriptionDataElementsType](r.SupplyConditionDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &SupplyConditionDescriptionListDataType{
		SupplyConditionDescriptionData: filteredItems,
	}
	return result, true
}

// SupplyConditionThresholdRelationListDataType PartialReader implementation
var _ PartialReader = (*SupplyConditionThresholdRelationListDataType)(nil)

func (r *SupplyConditionThresholdRelationListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[SupplyConditionThresholdRelationDataType, SupplyConditionThresholdRelationListDataSelectorsType, SupplyConditionThresholdRelationDataElementsType](r.SupplyConditionThresholdRelationData, filter)
	if !success {
		return r, false
	}

	result := &SupplyConditionThresholdRelationListDataType{
		SupplyConditionThresholdRelationData: filteredItems,
	}
	return result, true
}
