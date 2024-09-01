package model

// LoadControlEventListDataType

var _ Updater = (*LoadControlEventListDataType)(nil)

func (r *LoadControlEventListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []LoadControlEventDataType
	if newList != nil {
		newData = newList.(*LoadControlEventListDataType).LoadControlEventData
	}

	data, success := UpdateList(remoteWrite, r.LoadControlEventData, newData, filterPartial, filterDelete)

	if success && persist {
		r.LoadControlEventData = data
	}

	return data, success
}

// LoadControlStateListDataType

var _ Updater = (*LoadControlStateListDataType)(nil)

func (r *LoadControlStateListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []LoadControlStateDataType
	if newList != nil {
		newData = newList.(*LoadControlStateListDataType).LoadControlStateData
	}

	data, success := UpdateList(remoteWrite, r.LoadControlStateData, newData, filterPartial, filterDelete)

	if success && persist {
		r.LoadControlStateData = data
	}

	return data, success
}

// LoadControlLimitListDataType

var _ Updater = (*LoadControlLimitListDataType)(nil)

func (r *LoadControlLimitListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []LoadControlLimitDataType
	if newList != nil {
		newData = newList.(*LoadControlLimitListDataType).LoadControlLimitData
	}

	data, success := UpdateList(remoteWrite, r.LoadControlLimitData, newData, filterPartial, filterDelete)

	if success && persist {
		r.LoadControlLimitData = data
	}

	return data, success
}

// LoadControlLimitConstraintsListDataType

var _ Updater = (*LoadControlLimitConstraintsListDataType)(nil)

func (r *LoadControlLimitConstraintsListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []LoadControlLimitConstraintsDataType
	if newList != nil {
		newData = newList.(*LoadControlLimitConstraintsListDataType).LoadControlLimitConstraintsData
	}

	data, success := UpdateList(remoteWrite, r.LoadControlLimitConstraintsData, newData, filterPartial, filterDelete)

	if success && persist {
		r.LoadControlLimitConstraintsData = data
	}

	return data, success
}

// LoadControlLimitDescriptionListDataType

var _ Updater = (*LoadControlLimitDescriptionListDataType)(nil)

func (r *LoadControlLimitDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []LoadControlLimitDescriptionDataType
	if newList != nil {
		newData = newList.(*LoadControlLimitDescriptionListDataType).LoadControlLimitDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.LoadControlLimitDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.LoadControlLimitDescriptionData = data
	}

	return data, success
}
