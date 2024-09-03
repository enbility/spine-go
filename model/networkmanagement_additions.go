package model

// NetworkManagementDeviceDescriptionListDataType

var _ Updater = (*NetworkManagementDeviceDescriptionListDataType)(nil)

func (r *NetworkManagementDeviceDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []NetworkManagementDeviceDescriptionDataType
	if newList != nil {
		newData = newList.(*NetworkManagementDeviceDescriptionListDataType).NetworkManagementDeviceDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.NetworkManagementDeviceDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.NetworkManagementDeviceDescriptionData = data
	}

	return data, success
}

// NetworkManagementEntityDescriptionListDataType

var _ Updater = (*NetworkManagementEntityDescriptionListDataType)(nil)

func (r *NetworkManagementEntityDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []NetworkManagementEntityDescriptionDataType
	if newList != nil {
		newData = newList.(*NetworkManagementEntityDescriptionListDataType).NetworkManagementEntityDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.NetworkManagementEntityDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.NetworkManagementEntityDescriptionData = data
	}

	return data, success
}

// NetworkManagementFeatureDescriptionListDataType

var _ Updater = (*NetworkManagementFeatureDescriptionListDataType)(nil)

func (r *NetworkManagementFeatureDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []NetworkManagementFeatureDescriptionDataType
	if newList != nil {
		newData = newList.(*NetworkManagementFeatureDescriptionListDataType).NetworkManagementFeatureDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.NetworkManagementFeatureDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.NetworkManagementFeatureDescriptionData = data
	}

	return data, success
}
