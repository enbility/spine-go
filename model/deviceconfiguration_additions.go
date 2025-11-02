package model

// DeviceConfigurationKeyValueListDataType

var _ Updater = (*DeviceConfigurationKeyValueListDataType)(nil)

func (r *DeviceConfigurationKeyValueListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []DeviceConfigurationKeyValueDataType
	if newList != nil {
		newData = newList.(*DeviceConfigurationKeyValueListDataType).DeviceConfigurationKeyValueData
	}

	data, success := UpdateList(remoteWrite, r.DeviceConfigurationKeyValueData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.DeviceConfigurationKeyValueData = data
	}

	return data, success
}

// DeviceConfigurationKeyValueDescriptionListDataType

var _ Updater = (*DeviceConfigurationKeyValueDescriptionListDataType)(nil)

func (r *DeviceConfigurationKeyValueDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []DeviceConfigurationKeyValueDescriptionDataType
	if newList != nil {
		newData = newList.(*DeviceConfigurationKeyValueDescriptionListDataType).DeviceConfigurationKeyValueDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.DeviceConfigurationKeyValueDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.DeviceConfigurationKeyValueDescriptionData = data
	}

	return data, success
}

// DeviceConfigurationKeyValueConstraintsListDataType

var _ Updater = (*DeviceConfigurationKeyValueConstraintsListDataType)(nil)

func (r *DeviceConfigurationKeyValueConstraintsListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []DeviceConfigurationKeyValueConstraintsDataType
	if newList != nil {
		newData = newList.(*DeviceConfigurationKeyValueConstraintsListDataType).DeviceConfigurationKeyValueConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.DeviceConfigurationKeyValueConstraintsData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.DeviceConfigurationKeyValueConstraintsData = data
	}

	return data, success
}

// DeviceConfigurationKeyValueListDataType PartialReader implementation
var _ PartialReader = (*DeviceConfigurationKeyValueListDataType)(nil)

func (r *DeviceConfigurationKeyValueListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function with both selector and elements filtering
	filteredItems, success := partialListDataRead[DeviceConfigurationKeyValueDataType, DeviceConfigurationKeyValueListDataSelectorsType, DeviceConfigurationKeyValueDataElementsType](r.DeviceConfigurationKeyValueData, filter)
	if !success {
		return r, false
	}

	result := &DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: filteredItems,
	}
	return result, true
}

// DeviceConfigurationKeyValueDescriptionListDataType PartialReader implementation
var _ PartialReader = (*DeviceConfigurationKeyValueDescriptionListDataType)(nil)

func (r *DeviceConfigurationKeyValueDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function with both selector and elements filtering
	filteredItems, success := partialListDataRead[DeviceConfigurationKeyValueDescriptionDataType, DeviceConfigurationKeyValueDescriptionListDataSelectorsType, DeviceConfigurationKeyValueDescriptionDataElementsType](r.DeviceConfigurationKeyValueDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: filteredItems,
	}
	return result, true
}

// DeviceConfigurationKeyValueConstraintsListDataType PartialReader implementation
var _ PartialReader = (*DeviceConfigurationKeyValueConstraintsListDataType)(nil)

func (r *DeviceConfigurationKeyValueConstraintsListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function with both selector and elements filtering
	filteredItems, success := partialListDataRead[DeviceConfigurationKeyValueConstraintsDataType, DeviceConfigurationKeyValueConstraintsListDataSelectorsType, DeviceConfigurationKeyValueConstraintsDataElementsType](r.DeviceConfigurationKeyValueConstraintsData, filter)
	if !success {
		return r, false
	}

	result := &DeviceConfigurationKeyValueConstraintsListDataType{
		DeviceConfigurationKeyValueConstraintsData: filteredItems,
	}
	return result, true
}
