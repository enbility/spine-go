package model

// TaskManagementJobListDataType

var _ Updater = (*TaskManagementJobListDataType)(nil)

func (r *TaskManagementJobListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []TaskManagementJobDataType
	if newList != nil {
		newData = newList.(*TaskManagementJobListDataType).TaskManagementJobData
	}

	data, success := UpdateList(remoteWrite, r.TaskManagementJobData, newData, filterPartial, filterDelete)

	if success {
		r.TaskManagementJobData = data
	}

	return success
}

// TaskManagementJobRelationListDataType

var _ Updater = (*TaskManagementJobRelationListDataType)(nil)

func (r *TaskManagementJobRelationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []TaskManagementJobRelationDataType
	if newList != nil {
		newData = newList.(*TaskManagementJobRelationListDataType).TaskManagementJobRelationData
	}

	data, success := UpdateList(remoteWrite, r.TaskManagementJobRelationData, newData, filterPartial, filterDelete)

	if success {
		r.TaskManagementJobRelationData = data
	}

	return success
}

// TaskManagementJobDescriptionListDataType

var _ Updater = (*TaskManagementJobDescriptionListDataType)(nil)

func (r *TaskManagementJobDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []TaskManagementJobDescriptionDataType
	if newList != nil {
		newData = newList.(*TaskManagementJobDescriptionListDataType).TaskManagementJobDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.TaskManagementJobDescriptionData, newData, filterPartial, filterDelete)

	if success {
		r.TaskManagementJobDescriptionData = data
	}

	return success
}
