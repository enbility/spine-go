package model

// ElectricalConnectionStateListDataType

var _ Updater = (*ElectricalConnectionStateListDataType)(nil)

func (r *ElectricalConnectionStateListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []ElectricalConnectionStateDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionStateListDataType).ElectricalConnectionStateData
	}

	data, success := UpdateList(remoteWrite, r.ElectricalConnectionStateData, newData, filterPartial, filterDelete)

	if success && persist {
		r.ElectricalConnectionStateData = data
	}

	return data, success
}

// ElectricalConnectionPermittedValueSetListDataType

var _ Updater = (*ElectricalConnectionPermittedValueSetListDataType)(nil)

func (r *ElectricalConnectionPermittedValueSetListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []ElectricalConnectionPermittedValueSetDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionPermittedValueSetListDataType).ElectricalConnectionPermittedValueSetData
	}

	data, success := UpdateList(remoteWrite, r.ElectricalConnectionPermittedValueSetData, newData, filterPartial, filterDelete)

	if success && persist {
		r.ElectricalConnectionPermittedValueSetData = data
	}

	return data, success
}

// ElectricalConnectionDescriptionListDataType

var _ Updater = (*ElectricalConnectionDescriptionListDataType)(nil)

func (r *ElectricalConnectionDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []ElectricalConnectionDescriptionDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionDescriptionListDataType).ElectricalConnectionDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.ElectricalConnectionDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.ElectricalConnectionDescriptionData = data
	}

	return data, success
}

// ElectricalConnectionCharacteristicListDataType

var _ Updater = (*ElectricalConnectionCharacteristicListDataType)(nil)

func (r *ElectricalConnectionCharacteristicListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []ElectricalConnectionCharacteristicDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionCharacteristicListDataType).ElectricalConnectionCharacteristicData
	}

	data, success := UpdateList(remoteWrite, r.ElectricalConnectionCharacteristicData, newData, filterPartial, filterDelete)

	if success && persist {
		r.ElectricalConnectionCharacteristicData = data
	}

	return data, success
}

// ElectricalConnectionParameterDescriptionListDataType

var _ Updater = (*ElectricalConnectionParameterDescriptionListDataType)(nil)

func (r *ElectricalConnectionParameterDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []ElectricalConnectionParameterDescriptionDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionParameterDescriptionListDataType).ElectricalConnectionParameterDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.ElectricalConnectionParameterDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.ElectricalConnectionParameterDescriptionData = data
	}

	return data, success
}
