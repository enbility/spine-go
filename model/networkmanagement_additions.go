package model

// NetworkManagementDeviceDescriptionListDataType

var _ Updater = (*NetworkManagementDeviceDescriptionListDataType)(nil)

func (r *NetworkManagementDeviceDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []NetworkManagementDeviceDescriptionDataType
	if newList != nil {
		newData = newList.(*NetworkManagementDeviceDescriptionListDataType).NetworkManagementDeviceDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.NetworkManagementDeviceDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.NetworkManagementDeviceDescriptionData = data
	}

	return data, success
}

// NetworkManagementEntityDescriptionListDataType

var _ Updater = (*NetworkManagementEntityDescriptionListDataType)(nil)

func (r *NetworkManagementEntityDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []NetworkManagementEntityDescriptionDataType
	if newList != nil {
		newData = newList.(*NetworkManagementEntityDescriptionListDataType).NetworkManagementEntityDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.NetworkManagementEntityDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.NetworkManagementEntityDescriptionData = data
	}

	return data, success
}

// NetworkManagementFeatureDescriptionListDataType

var _ Updater = (*NetworkManagementFeatureDescriptionListDataType)(nil)

func (r *NetworkManagementFeatureDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []NetworkManagementFeatureDescriptionDataType
	if newList != nil {
		newData = newList.(*NetworkManagementFeatureDescriptionListDataType).NetworkManagementFeatureDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.NetworkManagementFeatureDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.NetworkManagementFeatureDescriptionData = data
	}

	return data, success
}

// NetworkManagementDeviceDescriptionListDataType PartialReader implementation
var _ PartialReader = (*NetworkManagementDeviceDescriptionListDataType)(nil)

func (r *NetworkManagementDeviceDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles device address matching
	filteredItems, success := partialListDataRead[NetworkManagementDeviceDescriptionDataType, NetworkManagementDeviceDescriptionListDataSelectorsType, NetworkManagementDeviceDescriptionDataElementsType](r.NetworkManagementDeviceDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &NetworkManagementDeviceDescriptionListDataType{
		NetworkManagementDeviceDescriptionData: filteredItems,
	}
	return result, true
}

// NetworkManagementEntityDescriptionListDataType PartialReader implementation
var _ PartialReader = (*NetworkManagementEntityDescriptionListDataType)(nil)

func (r *NetworkManagementEntityDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles entity address matching
	filteredItems, success := partialListDataRead[NetworkManagementEntityDescriptionDataType, NetworkManagementEntityDescriptionListDataSelectorsType, NetworkManagementEntityDescriptionDataElementsType](r.NetworkManagementEntityDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &NetworkManagementEntityDescriptionListDataType{
		NetworkManagementEntityDescriptionData: filteredItems,
	}
	return result, true
}

// NetworkManagementFeatureDescriptionListDataType PartialReader implementation
var _ PartialReader = (*NetworkManagementFeatureDescriptionListDataType)(nil)

func (r *NetworkManagementFeatureDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles feature address matching
	filteredItems, success := partialListDataRead[NetworkManagementFeatureDescriptionDataType, NetworkManagementFeatureDescriptionListDataSelectorsType, NetworkManagementFeatureDescriptionDataElementsType](r.NetworkManagementFeatureDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &NetworkManagementFeatureDescriptionListDataType{
		NetworkManagementFeatureDescriptionData: filteredItems,
	}
	return result, true
}
