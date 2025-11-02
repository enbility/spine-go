package model

// IdentificationListDataType

var _ Updater = (*IdentificationListDataType)(nil)

func (r *IdentificationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []IdentificationDataType
	if newList != nil {
		newData = newList.(*IdentificationListDataType).IdentificationData
	}

	data, success := UpdateList(remoteWrite, r.IdentificationData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.IdentificationData = data
	}

	return persist, success
}

// SessionIdentificationListDataType

var _ Updater = (*SessionIdentificationListDataType)(nil)

func (r *SessionIdentificationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []SessionIdentificationDataType
	if newList != nil {
		newData = newList.(*SessionIdentificationListDataType).SessionIdentificationData
	}

	data, success := UpdateList(remoteWrite, r.SessionIdentificationData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.SessionIdentificationData = data
	}

	return persist, success
}

// SessionMeasurementRelationListDataType

var _ Updater = (*SessionMeasurementRelationListDataType)(nil)

func (r *SessionMeasurementRelationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []SessionMeasurementRelationDataType
	if newList != nil {
		newData = newList.(*SessionMeasurementRelationListDataType).SessionMeasurementRelationData
	}

	data, success := UpdateList(remoteWrite, r.SessionMeasurementRelationData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.SessionMeasurementRelationData = data
	}

	return persist, success
}

// IdentificationListDataType PartialReader implementation
var _ PartialReader = (*IdentificationListDataType)(nil)

func (r *IdentificationListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function with both selector and elements filtering
	filteredItems, success := partialListDataRead[IdentificationDataType, IdentificationListDataSelectorsType, IdentificationDataElementsType](r.IdentificationData, filter)
	if !success {
		return r, false
	}

	result := &IdentificationListDataType{
		IdentificationData: filteredItems,
	}
	return result, true
}

// SessionIdentificationListDataType PartialReader implementation
var _ PartialReader = (*SessionIdentificationListDataType)(nil)

func (r *SessionIdentificationListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function with both selector and elements filtering
	filteredItems, success := partialListDataRead[SessionIdentificationDataType, SessionIdentificationListDataSelectorsType, SessionIdentificationDataElementsType](r.SessionIdentificationData, filter)
	if !success {
		return r, false
	}

	result := &SessionIdentificationListDataType{
		SessionIdentificationData: filteredItems,
	}
	return result, true
}

// SessionMeasurementRelationListDataType PartialReader implementation
var _ PartialReader = (*SessionMeasurementRelationListDataType)(nil)

func (r *SessionMeasurementRelationListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles array matching for MeasurementId
	filteredItems, success := partialListDataRead[SessionMeasurementRelationDataType, SessionMeasurementRelationListDataSelectorsType, SessionMeasurementRelationDataElementsType](r.SessionMeasurementRelationData, filter)
	if !success {
		return r, false
	}

	result := &SessionMeasurementRelationListDataType{
		SessionMeasurementRelationData: filteredItems,
	}
	return result, true
}
