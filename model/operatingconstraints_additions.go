package model

// OperatingConstraintsInterruptListDataType

var _ Updater = (*OperatingConstraintsInterruptListDataType)(nil)

func (r *OperatingConstraintsInterruptListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []OperatingConstraintsInterruptDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsInterruptListDataType).OperatingConstraintsInterruptData
	}

	r.OperatingConstraintsInterruptData = UpdateList(remoteWrite, r.OperatingConstraintsInterruptData, newData, filterPartial, filterDelete)
}

// OperatingConstraintsDurationListDataType

var _ Updater = (*OperatingConstraintsDurationListDataType)(nil)

func (r *OperatingConstraintsDurationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []OperatingConstraintsDurationDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsDurationListDataType).OperatingConstraintsDurationData
	}

	r.OperatingConstraintsDurationData = UpdateList(remoteWrite, r.OperatingConstraintsDurationData, newData, filterPartial, filterDelete)
}

// OperatingConstraintsPowerDescriptionListDataType

var _ Updater = (*OperatingConstraintsPowerDescriptionListDataType)(nil)

func (r *OperatingConstraintsPowerDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []OperatingConstraintsPowerDescriptionDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsPowerDescriptionListDataType).OperatingConstraintsPowerDescriptionData
	}

	r.OperatingConstraintsPowerDescriptionData = UpdateList(remoteWrite, r.OperatingConstraintsPowerDescriptionData, newData, filterPartial, filterDelete)
}

// OperatingConstraintsPowerRangeListDataType

var _ Updater = (*OperatingConstraintsPowerRangeListDataType)(nil)

func (r *OperatingConstraintsPowerRangeListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []OperatingConstraintsPowerRangeDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsPowerRangeListDataType).OperatingConstraintsPowerRangeData
	}

	r.OperatingConstraintsPowerRangeData = UpdateList(remoteWrite, r.OperatingConstraintsPowerRangeData, newData, filterPartial, filterDelete)
}

// OperatingConstraintsPowerLevelListDataType

var _ Updater = (*OperatingConstraintsPowerLevelListDataType)(nil)

func (r *OperatingConstraintsPowerLevelListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []OperatingConstraintsPowerLevelDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsPowerLevelListDataType).OperatingConstraintsPowerLevelData
	}

	r.OperatingConstraintsPowerLevelData = UpdateList(remoteWrite, r.OperatingConstraintsPowerLevelData, newData, filterPartial, filterDelete)
}

// OperatingConstraintsResumeImplicationListDataType

var _ Updater = (*OperatingConstraintsResumeImplicationListDataType)(nil)

func (r *OperatingConstraintsResumeImplicationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []OperatingConstraintsResumeImplicationDataType
	if newList != nil {
		newData = newList.(*OperatingConstraintsResumeImplicationListDataType).OperatingConstraintsResumeImplicationData
	}

	r.OperatingConstraintsResumeImplicationData = UpdateList(remoteWrite, r.OperatingConstraintsResumeImplicationData, newData, filterPartial, filterDelete)
}
