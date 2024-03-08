package model

// IdentificationListDataType

var _ Updater = (*IdentificationListDataType)(nil)

func (r *IdentificationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []IdentificationDataType
	if newList != nil {
		newData = newList.(*IdentificationListDataType).IdentificationData
	}

	r.IdentificationData = UpdateList(remoteWrite, r.IdentificationData, newData, filterPartial, filterDelete)
}

// SessionIdentificationListDataType

var _ Updater = (*SessionIdentificationListDataType)(nil)

func (r *SessionIdentificationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []SessionIdentificationDataType
	if newList != nil {
		newData = newList.(*SessionIdentificationListDataType).SessionIdentificationData
	}

	r.SessionIdentificationData = UpdateList(remoteWrite, r.SessionIdentificationData, newData, filterPartial, filterDelete)
}

// SessionMeasurementRelationListDataType

var _ Updater = (*SessionMeasurementRelationListDataType)(nil)

func (r *SessionMeasurementRelationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []SessionMeasurementRelationDataType
	if newList != nil {
		newData = newList.(*SessionMeasurementRelationListDataType).SessionMeasurementRelationData
	}

	r.SessionMeasurementRelationData = UpdateList(remoteWrite, r.SessionMeasurementRelationData, newData, filterPartial, filterDelete)
}
