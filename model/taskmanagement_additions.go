package model

// TaskManagementJobListDataType

var _ Updater = (*TaskManagementJobListDataType)(nil)

func (r *TaskManagementJobListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []TaskManagementJobDataType
	if newList != nil {
		newData = newList.(*TaskManagementJobListDataType).TaskManagementJobData
	}

	data, success := UpdateList(remoteWrite, r.TaskManagementJobData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.TaskManagementJobData = data
	}

	return data, success
}

// TaskManagementJobRelationListDataType

var _ Updater = (*TaskManagementJobRelationListDataType)(nil)

func (r *TaskManagementJobRelationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []TaskManagementJobRelationDataType
	if newList != nil {
		newData = newList.(*TaskManagementJobRelationListDataType).TaskManagementJobRelationData
	}

	data, success := UpdateList(remoteWrite, r.TaskManagementJobRelationData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.TaskManagementJobRelationData = data
	}

	return data, success
}

// TaskManagementJobDescriptionListDataType

var _ Updater = (*TaskManagementJobDescriptionListDataType)(nil)

func (r *TaskManagementJobDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []TaskManagementJobDescriptionDataType
	if newList != nil {
		newData = newList.(*TaskManagementJobDescriptionListDataType).TaskManagementJobDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.TaskManagementJobDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.TaskManagementJobDescriptionData = data
	}

	return data, success
}

// TaskManagementJobListDataType PartialReader implementation
var _ PartialReader = (*TaskManagementJobListDataType)(nil)

func (r *TaskManagementJobListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[TaskManagementJobDataType, TaskManagementJobListDataSelectorsType, TaskManagementJobDataElementsType](r.TaskManagementJobData, filter)
	if !success {
		return r, false
	}

	result := &TaskManagementJobListDataType{
		TaskManagementJobData: filteredItems,
	}
	return result, true
}

// TaskManagementJobRelationListDataType PartialReader implementation
var _ PartialReader = (*TaskManagementJobRelationListDataType)(nil)

func (r *TaskManagementJobRelationListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[TaskManagementJobRelationDataType, TaskManagementJobRelationListDataSelectorsType, TaskManagementJobRelationDataElementsType](r.TaskManagementJobRelationData, filter)
	if !success {
		return r, false
	}

	result := &TaskManagementJobRelationListDataType{
		TaskManagementJobRelationData: filteredItems,
	}
	return result, true
}

// TaskManagementJobDescriptionListDataType PartialReader implementation
var _ PartialReader = (*TaskManagementJobDescriptionListDataType)(nil)

func (r *TaskManagementJobDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[TaskManagementJobDescriptionDataType, TaskManagementJobDescriptionListDataSelectorsType, TaskManagementJobDescriptionDataElementsType](r.TaskManagementJobDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &TaskManagementJobDescriptionListDataType{
		TaskManagementJobDescriptionData: filteredItems,
	}
	return result, true
}
