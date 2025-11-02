package model

// ElectricalConnectionStateListDataType

var _ Updater = (*ElectricalConnectionStateListDataType)(nil)

func (r *ElectricalConnectionStateListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []ElectricalConnectionStateDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionStateListDataType).ElectricalConnectionStateData
	}

	data, success := UpdateList(remoteWrite, r.ElectricalConnectionStateData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.ElectricalConnectionStateData = data
	}

	return data, success
}

// ElectricalConnectionPermittedValueSetListDataType

var _ Updater = (*ElectricalConnectionPermittedValueSetListDataType)(nil)

func (r *ElectricalConnectionPermittedValueSetListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []ElectricalConnectionPermittedValueSetDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionPermittedValueSetListDataType).ElectricalConnectionPermittedValueSetData
	}

	data, success := UpdateList(remoteWrite, r.ElectricalConnectionPermittedValueSetData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.ElectricalConnectionPermittedValueSetData = data
	}

	return data, success
}

// ElectricalConnectionDescriptionListDataType

var _ Updater = (*ElectricalConnectionDescriptionListDataType)(nil)

func (r *ElectricalConnectionDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []ElectricalConnectionDescriptionDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionDescriptionListDataType).ElectricalConnectionDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.ElectricalConnectionDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.ElectricalConnectionDescriptionData = data
	}

	return data, success
}

// ElectricalConnectionCharacteristicListDataType

var _ Updater = (*ElectricalConnectionCharacteristicListDataType)(nil)

func (r *ElectricalConnectionCharacteristicListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []ElectricalConnectionCharacteristicDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionCharacteristicListDataType).ElectricalConnectionCharacteristicData
	}

	data, success := UpdateList(remoteWrite, r.ElectricalConnectionCharacteristicData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.ElectricalConnectionCharacteristicData = data
	}

	return data, success
}

// ElectricalConnectionParameterDescriptionListDataType

var _ Updater = (*ElectricalConnectionParameterDescriptionListDataType)(nil)

func (r *ElectricalConnectionParameterDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []ElectricalConnectionParameterDescriptionDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionParameterDescriptionListDataType).ElectricalConnectionParameterDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.ElectricalConnectionParameterDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.ElectricalConnectionParameterDescriptionData = data
	}

	return data, success
}

// ElectricalConnectionPermittedValueSetListDataType PartialReader implementation
var _ PartialReader = (*ElectricalConnectionPermittedValueSetListDataType)(nil)

func (r *ElectricalConnectionPermittedValueSetListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - handles multiple selector fields automatically
	filteredItems, success := partialListDataRead[ElectricalConnectionPermittedValueSetDataType, ElectricalConnectionPermittedValueSetListDataSelectorsType, ElectricalConnectionPermittedValueSetDataElementsType](r.ElectricalConnectionPermittedValueSetData, filter)
	if !success {
		return r, false
	}

	result := &ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: filteredItems,
	}
	return result, true
}

// ElectricalConnectionDescriptionListDataType PartialReader implementation
var _ PartialReader = (*ElectricalConnectionDescriptionListDataType)(nil)

func (r *ElectricalConnectionDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function with both selector and elements filtering
	filteredItems, success := partialListDataRead[ElectricalConnectionDescriptionDataType, ElectricalConnectionDescriptionListDataSelectorsType, ElectricalConnectionDescriptionDataElementsType](r.ElectricalConnectionDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: filteredItems,
	}
	return result, true
}

// ElectricalConnectionStateListDataType PartialReader implementation
var _ PartialReader = (*ElectricalConnectionStateListDataType)(nil)

func (r *ElectricalConnectionStateListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function with both selector and elements filtering
	filteredItems, success := partialListDataRead[ElectricalConnectionStateDataType, ElectricalConnectionStateListDataSelectorsType, ElectricalConnectionStateDataElementsType](r.ElectricalConnectionStateData, filter)
	if !success {
		return r, false
	}

	result := &ElectricalConnectionStateListDataType{
		ElectricalConnectionStateData: filteredItems,
	}
	return result, true
}

// ElectricalConnectionParameterDescriptionListDataType PartialReader implementation
var _ PartialReader = (*ElectricalConnectionParameterDescriptionListDataType)(nil)

func (r *ElectricalConnectionParameterDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function with both selector and elements filtering
	filteredItems, success := partialListDataRead[ElectricalConnectionParameterDescriptionDataType, ElectricalConnectionParameterDescriptionListDataSelectorsType, ElectricalConnectionParameterDescriptionDataElementsType](r.ElectricalConnectionParameterDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: filteredItems,
	}
	return result, true
}

// ElectricalConnectionCharacteristicListDataType PartialReader implementation
var _ PartialReader = (*ElectricalConnectionCharacteristicListDataType)(nil)

func (r *ElectricalConnectionCharacteristicListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles all field matching
	filteredItems, success := partialListDataRead[ElectricalConnectionCharacteristicDataType, ElectricalConnectionCharacteristicListDataSelectorsType, ElectricalConnectionCharacteristicDataElementsType](r.ElectricalConnectionCharacteristicData, filter)
	if !success {
		return r, false
	}

	result := &ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: filteredItems,
	}
	return result, true
}
