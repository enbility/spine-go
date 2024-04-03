package model

// ElectricalConnectionStateListDataType

var _ Updater = (*ElectricalConnectionStateListDataType)(nil)

func (r *ElectricalConnectionStateListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []ElectricalConnectionStateDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionStateListDataType).ElectricalConnectionStateData
	}

	data, success := UpdateList(remoteWrite, r.ElectricalConnectionStateData, newData, filterPartial, filterDelete)

	if success {
		r.ElectricalConnectionStateData = data
	}

	return success
}

// ElectricalConnectionPermittedValueSetListDataType

var _ Updater = (*ElectricalConnectionPermittedValueSetListDataType)(nil)

func (r *ElectricalConnectionPermittedValueSetListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []ElectricalConnectionPermittedValueSetDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionPermittedValueSetListDataType).ElectricalConnectionPermittedValueSetData
	}

	data, success := UpdateList(remoteWrite, r.ElectricalConnectionPermittedValueSetData, newData, filterPartial, filterDelete)

	if success {
		r.ElectricalConnectionPermittedValueSetData = data
	}

	return success
}

// ElectricalConnectionDescriptionListDataType

var _ Updater = (*ElectricalConnectionDescriptionListDataType)(nil)

func (r *ElectricalConnectionDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []ElectricalConnectionDescriptionDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionDescriptionListDataType).ElectricalConnectionDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.ElectricalConnectionDescriptionData, newData, filterPartial, filterDelete)

	if success {
		r.ElectricalConnectionDescriptionData = data
	}

	return success
}

// ElectricalConnectionCharacteristicListDataType

var _ Updater = (*ElectricalConnectionCharacteristicListDataType)(nil)

func (r *ElectricalConnectionCharacteristicListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []ElectricalConnectionCharacteristicDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionCharacteristicListDataType).ElectricalConnectionCharacteristicData
	}

	data, success := UpdateList(remoteWrite, r.ElectricalConnectionCharacteristicData, newData, filterPartial, filterDelete)

	if success {
		r.ElectricalConnectionCharacteristicData = data
	}

	return success
}

// ElectricalConnectionParameterDescriptionListDataType

var _ Updater = (*ElectricalConnectionParameterDescriptionListDataType)(nil)

func (r *ElectricalConnectionParameterDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []ElectricalConnectionParameterDescriptionDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionParameterDescriptionListDataType).ElectricalConnectionParameterDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.ElectricalConnectionParameterDescriptionData, newData, filterPartial, filterDelete)

	if success {
		r.ElectricalConnectionParameterDescriptionData = data
	}

	return success
}
