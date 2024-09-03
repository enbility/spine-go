package model

// DeviceConfigurationKeyValueListDataType

var _ Updater = (*DeviceConfigurationKeyValueListDataType)(nil)

func (r *DeviceConfigurationKeyValueListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []DeviceConfigurationKeyValueDataType
	if newList != nil {
		newData = newList.(*DeviceConfigurationKeyValueListDataType).DeviceConfigurationKeyValueData
	}

	data, success := UpdateList(remoteWrite, r.DeviceConfigurationKeyValueData, newData, filterPartial, filterDelete)

	if success && persist {
		r.DeviceConfigurationKeyValueData = data
	}

	return data, success
}

// DeviceConfigurationKeyValueDescriptionListDataType

var _ Updater = (*DeviceConfigurationKeyValueDescriptionListDataType)(nil)

func (r *DeviceConfigurationKeyValueDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []DeviceConfigurationKeyValueDescriptionDataType
	if newList != nil {
		newData = newList.(*DeviceConfigurationKeyValueDescriptionListDataType).DeviceConfigurationKeyValueDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.DeviceConfigurationKeyValueDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.DeviceConfigurationKeyValueDescriptionData = data
	}

	return data, success
}

// DeviceConfigurationKeyValueConstraintsListDataType

var _ Updater = (*DeviceConfigurationKeyValueConstraintsListDataType)(nil)

func (r *DeviceConfigurationKeyValueConstraintsListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []DeviceConfigurationKeyValueConstraintsDataType
	if newList != nil {
		newData = newList.(*DeviceConfigurationKeyValueConstraintsListDataType).DeviceConfigurationKeyValueConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.DeviceConfigurationKeyValueConstraintsData, newData, filterPartial, filterDelete)

	if success && persist {
		r.DeviceConfigurationKeyValueConstraintsData = data
	}

	return data, success
}
