package model

// TariffListDataType

var _ Updater = (*TariffListDataType)(nil)

func (r *TariffListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []TariffDataType
	if newList != nil {
		newData = newList.(*TariffListDataType).TariffData
	}

	data, success := UpdateList(remoteWrite, r.TariffData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.TariffData = data
	}

	return data, success
}

// TariffTierRelationListDataType

var _ Updater = (*TariffTierRelationListDataType)(nil)

func (r *TariffTierRelationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []TariffTierRelationDataType
	if newList != nil {
		newData = newList.(*TariffTierRelationListDataType).TariffTierRelationData
	}

	data, success := UpdateList(remoteWrite, r.TariffTierRelationData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.TariffTierRelationData = data
	}

	return data, success
}

// TariffBoundaryRelationListDataType

var _ Updater = (*TariffBoundaryRelationListDataType)(nil)

func (r *TariffBoundaryRelationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []TariffBoundaryRelationDataType
	if newList != nil {
		newData = newList.(*TariffBoundaryRelationListDataType).TariffBoundaryRelationData
	}

	data, success := UpdateList(remoteWrite, r.TariffBoundaryRelationData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.TariffBoundaryRelationData = data
	}

	return data, success
}

// TariffDescriptionListDataType

var _ Updater = (*TariffDescriptionListDataType)(nil)

func (r *TariffDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []TariffDescriptionDataType
	if newList != nil {
		newData = newList.(*TariffDescriptionListDataType).TariffDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.TariffDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.TariffDescriptionData = data
	}

	return data, success
}

// TierBoundaryListDataType

var _ Updater = (*TierBoundaryListDataType)(nil)

func (r *TierBoundaryListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []TierBoundaryDataType
	if newList != nil {
		newData = newList.(*TierBoundaryListDataType).TierBoundaryData
	}

	data, success := UpdateList(remoteWrite, r.TierBoundaryData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.TierBoundaryData = data
	}

	return data, success
}

// TierBoundaryDescriptionListDataType

var _ Updater = (*TierBoundaryDescriptionListDataType)(nil)

func (r *TierBoundaryDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []TierBoundaryDescriptionDataType
	if newList != nil {
		newData = newList.(*TierBoundaryDescriptionListDataType).TierBoundaryDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.TierBoundaryDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.TierBoundaryDescriptionData = data
	}

	return data, success
}

// CommodityListDataType

var _ Updater = (*CommodityListDataType)(nil)

func (r *CommodityListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []CommodityDataType
	if newList != nil {
		newData = newList.(*CommodityListDataType).CommodityData
	}

	data, success := UpdateList(remoteWrite, r.CommodityData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.CommodityData = data
	}

	return data, success
}

// TierListDataType

var _ Updater = (*TierListDataType)(nil)

func (r *TierListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []TierDataType
	if newList != nil {
		newData = newList.(*TierListDataType).TierData
	}

	data, success := UpdateList(remoteWrite, r.TierData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.TierData = data
	}

	return data, success
}

// TierIncentiveRelationListDataType

var _ Updater = (*TierIncentiveRelationListDataType)(nil)

func (r *TierIncentiveRelationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []TierIncentiveRelationDataType
	if newList != nil {
		newData = newList.(*TierIncentiveRelationListDataType).TierIncentiveRelationData
	}

	data, success := UpdateList(remoteWrite, r.TierIncentiveRelationData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.TierIncentiveRelationData = data
	}

	return data, success
}

// TierDescriptionListDataType

var _ Updater = (*TierDescriptionListDataType)(nil)

func (r *TierDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []TierDescriptionDataType
	if newList != nil {
		newData = newList.(*TierDescriptionListDataType).TierDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.TierDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.TierDescriptionData = data
	}

	return data, success
}

// IncentiveListDataType

var _ Updater = (*IncentiveListDataType)(nil)

func (r *IncentiveListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []IncentiveDataType
	if newList != nil {
		newData = newList.(*IncentiveListDataType).IncentiveData
	}

	data, success := UpdateList(remoteWrite, r.IncentiveData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.IncentiveData = data
	}

	return data, success
}

// IncentiveDescriptionListDataType

var _ Updater = (*IncentiveDescriptionListDataType)(nil)

func (r *IncentiveDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool) {
	var newData []IncentiveDescriptionDataType
	if newList != nil {
		newData = newList.(*IncentiveDescriptionListDataType).IncentiveDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.IncentiveDescriptionData, newData, filterPartial, filterDelete, cmdFunction)

	if success && persist {
		r.IncentiveDescriptionData = data
	}

	return data, success
}

// TariffListDataType PartialReader implementation
var _ PartialReader = (*TariffListDataType)(nil)

func (r *TariffListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[TariffDataType, TariffListDataSelectorsType, TariffDataElementsType](r.TariffData, filter)
	if !success {
		return r, false
	}

	result := &TariffListDataType{
		TariffData: filteredItems,
	}
	return result, true
}

// TariffTierRelationListDataType PartialReader implementation
var _ PartialReader = (*TariffTierRelationListDataType)(nil)

func (r *TariffTierRelationListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[TariffTierRelationDataType, TariffTierRelationListDataSelectorsType, TariffTierRelationDataElementsType](r.TariffTierRelationData, filter)
	if !success {
		return r, false
	}

	result := &TariffTierRelationListDataType{
		TariffTierRelationData: filteredItems,
	}
	return result, true
}

// TariffBoundaryRelationListDataType PartialReader implementation
var _ PartialReader = (*TariffBoundaryRelationListDataType)(nil)

func (r *TariffBoundaryRelationListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[TariffBoundaryRelationDataType, TariffBoundaryRelationListDataSelectorsType, TariffBoundaryRelationDataElementsType](r.TariffBoundaryRelationData, filter)
	if !success {
		return r, false
	}

	result := &TariffBoundaryRelationListDataType{
		TariffBoundaryRelationData: filteredItems,
	}
	return result, true
}

// TariffDescriptionListDataType PartialReader implementation
var _ PartialReader = (*TariffDescriptionListDataType)(nil)

func (r *TariffDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[TariffDescriptionDataType, TariffDescriptionListDataSelectorsType, TariffDescriptionDataElementsType](r.TariffDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &TariffDescriptionListDataType{
		TariffDescriptionData: filteredItems,
	}
	return result, true
}

// TierBoundaryListDataType PartialReader implementation
var _ PartialReader = (*TierBoundaryListDataType)(nil)

func (r *TierBoundaryListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[TierBoundaryDataType, TierBoundaryListDataSelectorsType, TierBoundaryDataElementsType](r.TierBoundaryData, filter)
	if !success {
		return r, false
	}

	result := &TierBoundaryListDataType{
		TierBoundaryData: filteredItems,
	}
	return result, true
}

// TierBoundaryDescriptionListDataType PartialReader implementation
var _ PartialReader = (*TierBoundaryDescriptionListDataType)(nil)

func (r *TierBoundaryDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[TierBoundaryDescriptionDataType, TierBoundaryDescriptionListDataSelectorsType, TierBoundaryDescriptionDataElementsType](r.TierBoundaryDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &TierBoundaryDescriptionListDataType{
		TierBoundaryDescriptionData: filteredItems,
	}
	return result, true
}

// CommodityListDataType PartialReader implementation
var _ PartialReader = (*CommodityListDataType)(nil)

func (r *CommodityListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[CommodityDataType, CommodityListDataSelectorsType, CommodityDataElementsType](r.CommodityData, filter)
	if !success {
		return r, false
	}

	result := &CommodityListDataType{
		CommodityData: filteredItems,
	}
	return result, true
}

// TierListDataType PartialReader implementation
var _ PartialReader = (*TierListDataType)(nil)

func (r *TierListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[TierDataType, TierListDataSelectorsType, TierDataElementsType](r.TierData, filter)
	if !success {
		return r, false
	}

	result := &TierListDataType{
		TierData: filteredItems,
	}
	return result, true
}

// TierIncentiveRelationListDataType PartialReader implementation
var _ PartialReader = (*TierIncentiveRelationListDataType)(nil)

func (r *TierIncentiveRelationListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[TierIncentiveRelationDataType, TierIncentiveRelationListDataSelectorsType, TierIncentiveRelationDataElementsType](r.TierIncentiveRelationData, filter)
	if !success {
		return r, false
	}

	result := &TierIncentiveRelationListDataType{
		TierIncentiveRelationData: filteredItems,
	}
	return result, true
}

// TierDescriptionListDataType PartialReader implementation
var _ PartialReader = (*TierDescriptionListDataType)(nil)

func (r *TierDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[TierDescriptionDataType, TierDescriptionListDataSelectorsType, TierDescriptionDataElementsType](r.TierDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &TierDescriptionListDataType{
		TierDescriptionData: filteredItems,
	}
	return result, true
}

// IncentiveListDataType PartialReader implementation
var _ PartialReader = (*IncentiveListDataType)(nil)

func (r *IncentiveListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[IncentiveDataType, IncentiveListDataSelectorsType, IncentiveDataElementsType](r.IncentiveData, filter)
	if !success {
		return r, false
	}

	result := &IncentiveListDataType{
		IncentiveData: filteredItems,
	}
	return result, true
}

// IncentiveDescriptionListDataType PartialReader implementation
var _ PartialReader = (*IncentiveDescriptionListDataType)(nil)

func (r *IncentiveDescriptionListDataType) ReadPartialData(filter *FilterType) (any, bool) {
	if filter == nil {
		return r, true
	}

	// Use complete generic function - automatically handles address matching
	filteredItems, success := partialListDataRead[IncentiveDescriptionDataType, IncentiveDescriptionListDataSelectorsType, IncentiveDescriptionDataElementsType](r.IncentiveDescriptionData, filter)
	if !success {
		return r, false
	}

	result := &IncentiveDescriptionListDataType{
		IncentiveDescriptionData: filteredItems,
	}
	return result, true
}
