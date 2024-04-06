package model

// LoadControlEventListDataType

var _ Updater = (*LoadControlEventListDataType)(nil)

func (r *LoadControlEventListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []LoadControlEventDataType
	if newList != nil {
		newData = newList.(*LoadControlEventListDataType).LoadControlEventData
	}

	data, success := UpdateList(remoteWrite, r.LoadControlEventData, newData, filterPartial, filterDelete)

	if success {
		r.LoadControlEventData = data
	}

	return success
}

// LoadControlStateListDataType

var _ Updater = (*LoadControlStateListDataType)(nil)

func (r *LoadControlStateListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []LoadControlStateDataType
	if newList != nil {
		newData = newList.(*LoadControlStateListDataType).LoadControlStateData
	}

	data, success := UpdateList(remoteWrite, r.LoadControlStateData, newData, filterPartial, filterDelete)

	if success {
		r.LoadControlStateData = data
	}

	return success
}

// LoadControlLimitListDataType

var _ Updater = (*LoadControlLimitListDataType)(nil)

func (r *LoadControlLimitListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []LoadControlLimitDataType
	if newList != nil {
		newData = newList.(*LoadControlLimitListDataType).LoadControlLimitData
	}

	data, success := UpdateList(remoteWrite, r.LoadControlLimitData, newData, filterPartial, filterDelete)

	if success {
		r.LoadControlLimitData = data
	}

	return success
}

// LoadControlLimitConstraintsListDataType

var _ Updater = (*LoadControlLimitConstraintsListDataType)(nil)

func (r *LoadControlLimitConstraintsListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []LoadControlLimitConstraintsDataType
	if newList != nil {
		newData = newList.(*LoadControlLimitConstraintsListDataType).LoadControlLimitConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.LoadControlLimitConstraintsData, newData, filterPartial, filterDelete)

	if success {
		r.LoadControlLimitConstraintsData = data
	}

	return success
}

// LoadControlLimitDescriptionListDataType

var _ Updater = (*LoadControlLimitDescriptionListDataType)(nil)

func (r *LoadControlLimitDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []LoadControlLimitDescriptionDataType
	if newList != nil {
		newData = newList.(*LoadControlLimitDescriptionListDataType).LoadControlLimitDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.LoadControlLimitDescriptionData, newData, filterPartial, filterDelete)

	if success {
		r.LoadControlLimitDescriptionData = data
	}

	return success
}
