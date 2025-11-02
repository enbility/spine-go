package model

// OperatingConstraintsInterruptListDataType

var _ Updater = (*OperatingConstraintsInterruptListDataType)(nil)

func (r *OperatingConstraintsInterruptListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []OperatingConstraintsInterruptDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsInterruptListDataType).OperatingConstraintsInterruptData
	}

	data, success := UpdateList(remoteWrite, r.OperatingConstraintsInterruptData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.OperatingConstraintsInterruptData = data
	}

	return data, success
}

// OperatingConstraintsDurationListDataType

var _ Updater = (*OperatingConstraintsDurationListDataType)(nil)

func (r *OperatingConstraintsDurationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []OperatingConstraintsDurationDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsDurationListDataType).OperatingConstraintsDurationData
	}

	data, success := UpdateList(remoteWrite, r.OperatingConstraintsDurationData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.OperatingConstraintsDurationData = data
	}

	return data, success
}

// OperatingConstraintsPowerDescriptionListDataType

var _ Updater = (*OperatingConstraintsPowerDescriptionListDataType)(nil)

func (r *OperatingConstraintsPowerDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []OperatingConstraintsPowerDescriptionDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsPowerDescriptionListDataType).OperatingConstraintsPowerDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.OperatingConstraintsPowerDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.OperatingConstraintsPowerDescriptionData = data
	}

	return data, success
}

// OperatingConstraintsPowerRangeListDataType

var _ Updater = (*OperatingConstraintsPowerRangeListDataType)(nil)

func (r *OperatingConstraintsPowerRangeListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []OperatingConstraintsPowerRangeDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsPowerRangeListDataType).OperatingConstraintsPowerRangeData
	}

	data, success := UpdateList(remoteWrite, r.OperatingConstraintsPowerRangeData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.OperatingConstraintsPowerRangeData = data
	}

	return data, success
}

// OperatingConstraintsPowerLevelListDataType

var _ Updater = (*OperatingConstraintsPowerLevelListDataType)(nil)

func (r *OperatingConstraintsPowerLevelListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []OperatingConstraintsPowerLevelDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsPowerLevelListDataType).OperatingConstraintsPowerLevelData
	}

	data, success := UpdateList(remoteWrite, r.OperatingConstraintsPowerLevelData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.OperatingConstraintsPowerLevelData = data
	}

	return data, success
}

// OperatingConstraintsResumeImplicationListDataType

var _ Updater = (*OperatingConstraintsResumeImplicationListDataType)(nil)

func (r *OperatingConstraintsResumeImplicationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []OperatingConstraintsResumeImplicationDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsResumeImplicationListDataType).OperatingConstraintsResumeImplicationData
	}

	data, success := UpdateList(remoteWrite, r.OperatingConstraintsResumeImplicationData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.OperatingConstraintsResumeImplicationData = data
	}

	return data, success
}

// OperatingConstraintsInterruptListDataType PartialReader implementation
var _ PartialReader = (*OperatingConstraintsInterruptListDataType)(nil)

func (r *OperatingConstraintsInterruptListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[OperatingConstraintsInterruptDataType, OperatingConstraintsInterruptListDataSelectorsType, OperatingConstraintsInterruptDataElementsType](r.OperatingConstraintsInterruptData, filter)
	if !success {
		return r, false
	}

	result := &OperatingConstraintsInterruptListDataType{
		OperatingConstraintsInterruptData: filteredItems,
	}
	return result, true
}

// OperatingConstraintsDurationListDataType PartialReader implementation
var _ PartialReader = (*OperatingConstraintsDurationListDataType)(nil)

func (r *OperatingConstraintsDurationListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[OperatingConstraintsDurationDataType, OperatingConstraintsDurationListDataSelectorsType, OperatingConstraintsDurationDataElementsType](r.OperatingConstraintsDurationData, filter)
	if !success {
		return r, false
	}

	result := &OperatingConstraintsDurationListDataType{
		OperatingConstraintsDurationData: filteredItems,
	}
	return result, true
}

// OperatingConstraintsPowerDescriptionListDataType PartialReader implementation
var _ PartialReader = (*OperatingConstraintsPowerDescriptionListDataType)(nil)

func (r *OperatingConstraintsPowerDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[OperatingConstraintsPowerDescriptionDataType, OperatingConstraintsPowerDescriptionListDataSelectorsType, OperatingConstraintsPowerDescriptionDataElementsType](r.OperatingConstraintsPowerDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &OperatingConstraintsPowerDescriptionListDataType{
		OperatingConstraintsPowerDescriptionData: filteredItems,
	}
	return result, true
}

// OperatingConstraintsPowerRangeListDataType PartialReader implementation
var _ PartialReader = (*OperatingConstraintsPowerRangeListDataType)(nil)

func (r *OperatingConstraintsPowerRangeListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[OperatingConstraintsPowerRangeDataType, OperatingConstraintsPowerRangeListDataSelectorsType, OperatingConstraintsPowerRangeDataElementsType](r.OperatingConstraintsPowerRangeData, filter)
	if !success {
		return r, false
	}

	result := &OperatingConstraintsPowerRangeListDataType{
		OperatingConstraintsPowerRangeData: filteredItems,
	}
	return result, true
}

// OperatingConstraintsPowerLevelListDataType PartialReader implementation
var _ PartialReader = (*OperatingConstraintsPowerLevelListDataType)(nil)

func (r *OperatingConstraintsPowerLevelListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[OperatingConstraintsPowerLevelDataType, OperatingConstraintsPowerLevelListDataSelectorsType, OperatingConstraintsPowerLevelDataElementsType](r.OperatingConstraintsPowerLevelData, filter)
	if !success {
		return r, false
	}

	result := &OperatingConstraintsPowerLevelListDataType{
		OperatingConstraintsPowerLevelData: filteredItems,
	}
	return result, true
}

// OperatingConstraintsResumeImplicationListDataType PartialReader implementation
var _ PartialReader = (*OperatingConstraintsResumeImplicationListDataType)(nil)

func (r *OperatingConstraintsResumeImplicationListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[OperatingConstraintsResumeImplicationDataType, OperatingConstraintsResumeImplicationListDataSelectorsType, OperatingConstraintsResumeImplicationDataElementsType](r.OperatingConstraintsResumeImplicationData, filter)
	if !success {
		return r, false
	}

	result := &OperatingConstraintsResumeImplicationListDataType{
		OperatingConstraintsResumeImplicationData: filteredItems,
	}
	return result, true
}
