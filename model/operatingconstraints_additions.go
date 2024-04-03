package model

// OperatingConstraintsInterruptListDataType

var _ Updater = (*OperatingConstraintsInterruptListDataType)(nil)

func (r *OperatingConstraintsInterruptListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []OperatingConstraintsInterruptDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsInterruptListDataType).OperatingConstraintsInterruptData
	}

	data, success := UpdateList(remoteWrite, r.OperatingConstraintsInterruptData, newData, filterPartial, filterDelete)

	if success {
		r.OperatingConstraintsInterruptData = data
	}

	return success
}

// OperatingConstraintsDurationListDataType

var _ Updater = (*OperatingConstraintsDurationListDataType)(nil)

func (r *OperatingConstraintsDurationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []OperatingConstraintsDurationDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsDurationListDataType).OperatingConstraintsDurationData
	}

	data, success := UpdateList(remoteWrite, r.OperatingConstraintsDurationData, newData, filterPartial, filterDelete)

	if success {
		r.OperatingConstraintsDurationData = data
	}

	return success
}

// OperatingConstraintsPowerDescriptionListDataType

var _ Updater = (*OperatingConstraintsPowerDescriptionListDataType)(nil)

func (r *OperatingConstraintsPowerDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []OperatingConstraintsPowerDescriptionDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsPowerDescriptionListDataType).OperatingConstraintsPowerDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.OperatingConstraintsPowerDescriptionData, newData, filterPartial, filterDelete)

	if success {
		r.OperatingConstraintsPowerDescriptionData = data
	}

	return success
}

// OperatingConstraintsPowerRangeListDataType

var _ Updater = (*OperatingConstraintsPowerRangeListDataType)(nil)

func (r *OperatingConstraintsPowerRangeListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []OperatingConstraintsPowerRangeDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsPowerRangeListDataType).OperatingConstraintsPowerRangeData
	}

	data, success := UpdateList(remoteWrite, r.OperatingConstraintsPowerRangeData, newData, filterPartial, filterDelete)

	if success {
		r.OperatingConstraintsPowerRangeData = data
	}

	return success
}

// OperatingConstraintsPowerLevelListDataType

var _ Updater = (*OperatingConstraintsPowerLevelListDataType)(nil)

func (r *OperatingConstraintsPowerLevelListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []OperatingConstraintsPowerLevelDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsPowerLevelListDataType).OperatingConstraintsPowerLevelData
	}

	data, success := UpdateList(remoteWrite, r.OperatingConstraintsPowerLevelData, newData, filterPartial, filterDelete)

	if success {
		r.OperatingConstraintsPowerLevelData = data
	}

	return success
}

// OperatingConstraintsResumeImplicationListDataType

var _ Updater = (*OperatingConstraintsResumeImplicationListDataType)(nil)

func (r *OperatingConstraintsResumeImplicationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []OperatingConstraintsResumeImplicationDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsResumeImplicationListDataType).OperatingConstraintsResumeImplicationData
	}

	data, success := UpdateList(remoteWrite, r.OperatingConstraintsResumeImplicationData, newData, filterPartial, filterDelete)

	if success {
		r.OperatingConstraintsResumeImplicationData = data
	}

	return success
}
