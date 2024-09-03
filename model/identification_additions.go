package model

// IdentificationListDataType

var _ Updater = (*IdentificationListDataType)(nil)

func (r *IdentificationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []IdentificationDataType
	if newList != nil {
		newData = newList.(*IdentificationListDataType).IdentificationData
	}

	data, success := UpdateList(remoteWrite, r.IdentificationData, newData, filterPartial, filterDelete)

	if success && persist {
		r.IdentificationData = data
	}

	return persist, success
}

// SessionIdentificationListDataType

var _ Updater = (*SessionIdentificationListDataType)(nil)

func (r *SessionIdentificationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []SessionIdentificationDataType
	if newList != nil {
		newData = newList.(*SessionIdentificationListDataType).SessionIdentificationData
	}

	data, success := UpdateList(remoteWrite, r.SessionIdentificationData, newData, filterPartial, filterDelete)

	if success && persist {
		r.SessionIdentificationData = data
	}

	return persist, success
}

// SessionMeasurementRelationListDataType

var _ Updater = (*SessionMeasurementRelationListDataType)(nil)

func (r *SessionMeasurementRelationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []SessionMeasurementRelationDataType
	if newList != nil {
		newData = newList.(*SessionMeasurementRelationListDataType).SessionMeasurementRelationData
	}

	data, success := UpdateList(remoteWrite, r.SessionMeasurementRelationData, newData, filterPartial, filterDelete)

	if success && persist {
		r.SessionMeasurementRelationData = data
	}

	return persist, success
}
