package model

// LoadControlEventListDataType

var _ Updater = (*LoadControlEventListDataType)(nil)

func (r *LoadControlEventListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []LoadControlEventDataType
	if newList != nil {
		newData = newList.(*LoadControlEventListDataType).LoadControlEventData
	}

	data, success := UpdateList(remoteWrite, r.LoadControlEventData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.LoadControlEventData = data
	}

	return data, success
}

// LoadControlStateListDataType

var _ Updater = (*LoadControlStateListDataType)(nil)

func (r *LoadControlStateListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []LoadControlStateDataType
	if newList != nil {
		newData = newList.(*LoadControlStateListDataType).LoadControlStateData
	}

	data, success := UpdateList(remoteWrite, r.LoadControlStateData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.LoadControlStateData = data
	}

	return data, success
}

// LoadControlLimitListDataType

var _ Updater = (*LoadControlLimitListDataType)(nil)

func (r *LoadControlLimitListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []LoadControlLimitDataType
	if newList != nil {
		newData = newList.(*LoadControlLimitListDataType).LoadControlLimitData
	}

	data, success := UpdateList(remoteWrite, r.LoadControlLimitData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.LoadControlLimitData = data
	}

	return data, success
}

// LoadControlLimitConstraintsListDataType

var _ Updater = (*LoadControlLimitConstraintsListDataType)(nil)

func (r *LoadControlLimitConstraintsListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []LoadControlLimitConstraintsDataType
	if newList != nil {
		newData = newList.(*LoadControlLimitConstraintsListDataType).LoadControlLimitConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.LoadControlLimitConstraintsData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.LoadControlLimitConstraintsData = data
	}

	return data, success
}

// LoadControlLimitDescriptionListDataType

var _ Updater = (*LoadControlLimitDescriptionListDataType)(nil)

func (r *LoadControlLimitDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []LoadControlLimitDescriptionDataType
	if newList != nil {
		newData = newList.(*LoadControlLimitDescriptionListDataType).LoadControlLimitDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.LoadControlLimitDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.LoadControlLimitDescriptionData = data
	}

	return data, success
}

// LoadControlLimitListDataType PartialReader implementation
var _ PartialReader = (*LoadControlLimitListDataType)(nil)

func (r *LoadControlLimitListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function with both selector and elements filtering
	filteredItems, success := partialListDataRead[LoadControlLimitDataType, LoadControlLimitListDataSelectorsType, LoadControlLimitDataElementsType](r.LoadControlLimitData, filter)
	if !success {
		return r, false
	}

	result := &LoadControlLimitListDataType{
		LoadControlLimitData: filteredItems,
	}
	return result, true
}

// LoadControlEventListDataType PartialReader implementation
var _ PartialReader = (*LoadControlEventListDataType)(nil)

func (r *LoadControlEventListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function with both selector and elements filtering
	filteredItems, success := partialListDataRead[LoadControlEventDataType, LoadControlEventListDataSelectorsType, LoadControlEventDataElementsType](r.LoadControlEventData, filter)
	if !success {
		return r, false
	}

	result := &LoadControlEventListDataType{
		LoadControlEventData: filteredItems,
	}
	return result, true
}

// LoadControlStateListDataType PartialReader implementation
var _ PartialReader = (*LoadControlStateListDataType)(nil)

func (r *LoadControlStateListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function with both selector and elements filtering
	filteredItems, success := partialListDataRead[LoadControlStateDataType, LoadControlStateListDataSelectorsType, LoadControlStateDataElementsType](r.LoadControlStateData, filter)
	if !success {
		return r, false
	}

	result := &LoadControlStateListDataType{
		LoadControlStateData: filteredItems,
	}
	return result, true
}

// LoadControlLimitConstraintsListDataType PartialReader implementation
var _ PartialReader = (*LoadControlLimitConstraintsListDataType)(nil)

func (r *LoadControlLimitConstraintsListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function with both selector and elements filtering
	filteredItems, success := partialListDataRead[LoadControlLimitConstraintsDataType, LoadControlLimitConstraintsListDataSelectorsType, LoadControlLimitConstraintsDataElementsType](r.LoadControlLimitConstraintsData, filter)
	if !success {
		return r, false
	}

	result := &LoadControlLimitConstraintsListDataType{
		LoadControlLimitConstraintsData: filteredItems,
	}
	return result, true
}

// LoadControlLimitDescriptionListDataType PartialReader implementation
var _ PartialReader = (*LoadControlLimitDescriptionListDataType)(nil)

func (r *LoadControlLimitDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function with both selector and elements filtering
	filteredItems, success := partialListDataRead[LoadControlLimitDescriptionDataType, LoadControlLimitDescriptionListDataSelectorsType, LoadControlLimitDescriptionDataElementsType](r.LoadControlLimitDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: filteredItems,
	}
	return result, true
}
