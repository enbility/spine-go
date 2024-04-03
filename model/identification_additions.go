package model

// IdentificationListDataType

var _ Updater = (*IdentificationListDataType)(nil)

func (r *IdentificationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []IdentificationDataType
	if newList != nil {
		newData = newList.(*IdentificationListDataType).IdentificationData
	}

	data, success := UpdateList(remoteWrite, r.IdentificationData, newData, filterPartial, filterDelete)

	if success {
		r.IdentificationData = data
	}

	return success
}

// SessionIdentificationListDataType

var _ Updater = (*SessionIdentificationListDataType)(nil)

func (r *SessionIdentificationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []SessionIdentificationDataType
	if newList != nil {
		newData = newList.(*SessionIdentificationListDataType).SessionIdentificationData
	}

	data, success := UpdateList(remoteWrite, r.SessionIdentificationData, newData, filterPartial, filterDelete)

	if success {
		r.SessionIdentificationData = data
	}

	return success
}

// SessionMeasurementRelationListDataType

var _ Updater = (*SessionMeasurementRelationListDataType)(nil)

func (r *SessionMeasurementRelationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []SessionMeasurementRelationDataType
	if newList != nil {
		newData = newList.(*SessionMeasurementRelationListDataType).SessionMeasurementRelationData
	}

	data, success := UpdateList(remoteWrite, r.SessionMeasurementRelationData, newData, filterPartial, filterDelete)

	if success {
		r.SessionMeasurementRelationData = data
	}

	return success
}
