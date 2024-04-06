package model

// NetworkManagementDeviceDescriptionListDataType

var _ Updater = (*NetworkManagementDeviceDescriptionListDataType)(nil)

func (r *NetworkManagementDeviceDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []NetworkManagementDeviceDescriptionDataType
	if newList != nil {
		newData = newList.(*NetworkManagementDeviceDescriptionListDataType).NetworkManagementDeviceDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.NetworkManagementDeviceDescriptionData, newData, filterPartial, filterDelete)

	if success {
		r.NetworkManagementDeviceDescriptionData = data
	}

	return success
}

// NetworkManagementEntityDescriptionListDataType

var _ Updater = (*NetworkManagementEntityDescriptionListDataType)(nil)

func (r *NetworkManagementEntityDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []NetworkManagementEntityDescriptionDataType
	if newList != nil {
		newData = newList.(*NetworkManagementEntityDescriptionListDataType).NetworkManagementEntityDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.NetworkManagementEntityDescriptionData, newData, filterPartial, filterDelete)

	if success {
		r.NetworkManagementEntityDescriptionData = data
	}

	return success
}

// NetworkManagementFeatureDescriptionListDataType

var _ Updater = (*NetworkManagementFeatureDescriptionListDataType)(nil)

func (r *NetworkManagementFeatureDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []NetworkManagementFeatureDescriptionDataType
	if newList != nil {
		newData = newList.(*NetworkManagementFeatureDescriptionListDataType).NetworkManagementFeatureDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.NetworkManagementFeatureDescriptionData, newData, filterPartial, filterDelete)

	if success {
		r.NetworkManagementFeatureDescriptionData = data
	}

	return success
}
