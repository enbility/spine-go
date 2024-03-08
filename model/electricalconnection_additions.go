package model

// ElectricalConnectionStateListDataType

var _ Updater = (*ElectricalConnectionStateListDataType)(nil)

func (r *ElectricalConnectionStateListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []ElectricalConnectionStateDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionStateListDataType).ElectricalConnectionStateData
	}

	r.ElectricalConnectionStateData = UpdateList(remoteWrite, r.ElectricalConnectionStateData, newData, filterPartial, filterDelete)
}

// ElectricalConnectionPermittedValueSetListDataType

var _ Updater = (*ElectricalConnectionPermittedValueSetListDataType)(nil)

func (r *ElectricalConnectionPermittedValueSetListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []ElectricalConnectionPermittedValueSetDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionPermittedValueSetListDataType).ElectricalConnectionPermittedValueSetData
	}

	r.ElectricalConnectionPermittedValueSetData = UpdateList(remoteWrite, r.ElectricalConnectionPermittedValueSetData, newData, filterPartial, filterDelete)
}

// ElectricalConnectionDescriptionListDataType

var _ Updater = (*ElectricalConnectionDescriptionListDataType)(nil)

func (r *ElectricalConnectionDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []ElectricalConnectionDescriptionDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionDescriptionListDataType).ElectricalConnectionDescriptionData
	}

	r.ElectricalConnectionDescriptionData = UpdateList(remoteWrite, r.ElectricalConnectionDescriptionData, newData, filterPartial, filterDelete)
}

// ElectricalConnectionCharacteristicListDataType

var _ Updater = (*ElectricalConnectionCharacteristicListDataType)(nil)

func (r *ElectricalConnectionCharacteristicListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []ElectricalConnectionCharacteristicDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionCharacteristicListDataType).ElectricalConnectionCharacteristicListData
	}

	r.ElectricalConnectionCharacteristicListData = UpdateList(remoteWrite, r.ElectricalConnectionCharacteristicListData, newData, filterPartial, filterDelete)
}

// ElectricalConnectionParameterDescriptionListDataType

var _ Updater = (*ElectricalConnectionParameterDescriptionListDataType)(nil)

func (r *ElectricalConnectionParameterDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []ElectricalConnectionParameterDescriptionDataType
	if newList != nil {
		newData = newList.(*ElectricalConnectionParameterDescriptionListDataType).ElectricalConnectionParameterDescriptionData
	}

	r.ElectricalConnectionParameterDescriptionData = UpdateList(remoteWrite, r.ElectricalConnectionParameterDescriptionData, newData, filterPartial, filterDelete)
}
