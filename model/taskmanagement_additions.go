package model

// TaskManagementJobListDataType

var _ Updater = (*TaskManagementJobListDataType)(nil)

func (r *TaskManagementJobListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []TaskManagementJobDataType
	if newList != nil {
		newData = newList.(*TaskManagementJobListDataType).TaskManagementJobData
	}

	data, success := UpdateList(remoteWrite, r.TaskManagementJobData, newData, filterPartial, filterDelete)

	if success && persist {
		r.TaskManagementJobData = data
	}

	return data, success
}

// TaskManagementJobRelationListDataType

var _ Updater = (*TaskManagementJobRelationListDataType)(nil)

func (r *TaskManagementJobRelationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []TaskManagementJobRelationDataType
	if newList != nil {
		newData = newList.(*TaskManagementJobRelationListDataType).TaskManagementJobRelationData
	}

	data, success := UpdateList(remoteWrite, r.TaskManagementJobRelationData, newData, filterPartial, filterDelete)

	if success && persist {
		r.TaskManagementJobRelationData = data
	}

	return data, success
}

// TaskManagementJobDescriptionListDataType

var _ Updater = (*TaskManagementJobDescriptionListDataType)(nil)

func (r *TaskManagementJobDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []TaskManagementJobDescriptionDataType
	if newList != nil {
		newData = newList.(*TaskManagementJobDescriptionListDataType).TaskManagementJobDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.TaskManagementJobDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.TaskManagementJobDescriptionData = data
	}

	return data, success
}
