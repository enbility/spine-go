package model

// TariffListDataType

var _ Updater = (*TariffListDataType)(nil)

func (r *TariffListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []TariffDataType
	if newList != nil {
		newData = newList.(*TariffListDataType).TariffData
	}

	data, success := UpdateList(remoteWrite, r.TariffData, newData, filterPartial, filterDelete)

	if success && persist {
		r.TariffData = data
	}

	return data, success
}

// TariffTierRelationListDataType

var _ Updater = (*TariffTierRelationListDataType)(nil)

func (r *TariffTierRelationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []TariffTierRelationDataType
	if newList != nil {
		newData = newList.(*TariffTierRelationListDataType).TariffTierRelationData
	}

	data, success := UpdateList(remoteWrite, r.TariffTierRelationData, newData, filterPartial, filterDelete)

	if success && persist {
		r.TariffTierRelationData = data
	}

	return data, success
}

// TariffBoundaryRelationListDataType

var _ Updater = (*TariffBoundaryRelationListDataType)(nil)

func (r *TariffBoundaryRelationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []TariffBoundaryRelationDataType
	if newList != nil {
		newData = newList.(*TariffBoundaryRelationListDataType).TariffBoundaryRelationData
	}

	data, success := UpdateList(remoteWrite, r.TariffBoundaryRelationData, newData, filterPartial, filterDelete)

	if success && persist {
		r.TariffBoundaryRelationData = data
	}

	return data, success
}

// TariffDescriptionListDataType

var _ Updater = (*TariffDescriptionListDataType)(nil)

func (r *TariffDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []TariffDescriptionDataType
	if newList != nil {
		newData = newList.(*TariffDescriptionListDataType).TariffDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.TariffDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.TariffDescriptionData = data
	}

	return data, success
}

// TierBoundaryListDataType

var _ Updater = (*TierBoundaryListDataType)(nil)

func (r *TierBoundaryListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []TierBoundaryDataType
	if newList != nil {
		newData = newList.(*TierBoundaryListDataType).TierBoundaryData
	}

	data, success := UpdateList(remoteWrite, r.TierBoundaryData, newData, filterPartial, filterDelete)

	if success && persist {
		r.TierBoundaryData = data
	}

	return data, success
}

// TierBoundaryDescriptionListDataType

var _ Updater = (*TierBoundaryDescriptionListDataType)(nil)

func (r *TierBoundaryDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []TierBoundaryDescriptionDataType
	if newList != nil {
		newData = newList.(*TierBoundaryDescriptionListDataType).TierBoundaryDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.TierBoundaryDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.TierBoundaryDescriptionData = data
	}

	return data, success
}

// CommodityListDataType

var _ Updater = (*CommodityListDataType)(nil)

func (r *CommodityListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []CommodityDataType
	if newList != nil {
		newData = newList.(*CommodityListDataType).CommodityData
	}

	data, success := UpdateList(remoteWrite, r.CommodityData, newData, filterPartial, filterDelete)

	if success && persist {
		r.CommodityData = data
	}

	return data, success
}

// TierListDataType

var _ Updater = (*TierListDataType)(nil)

func (r *TierListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []TierDataType
	if newList != nil {
		newData = newList.(*TierListDataType).TierData
	}

	data, success := UpdateList(remoteWrite, r.TierData, newData, filterPartial, filterDelete)

	if success && persist {
		r.TierData = data
	}

	return data, success
}

// TierIncentiveRelationListDataType

var _ Updater = (*TierIncentiveRelationListDataType)(nil)

func (r *TierIncentiveRelationListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []TierIncentiveRelationDataType
	if newList != nil {
		newData = newList.(*TierIncentiveRelationListDataType).TierIncentiveRelationData
	}

	data, success := UpdateList(remoteWrite, r.TierIncentiveRelationData, newData, filterPartial, filterDelete)

	if success && persist {
		r.TierIncentiveRelationData = data
	}

	return data, success
}

// TierDescriptionListDataType

var _ Updater = (*TierDescriptionListDataType)(nil)

func (r *TierDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []TierDescriptionDataType
	if newList != nil {
		newData = newList.(*TierDescriptionListDataType).TierDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.TierDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.TierDescriptionData = data
	}

	return data, success
}

// IncentiveListDataType

var _ Updater = (*IncentiveListDataType)(nil)

func (r *IncentiveListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []IncentiveDataType
	if newList != nil {
		newData = newList.(*IncentiveListDataType).IncentiveData
	}

	data, success := UpdateList(remoteWrite, r.IncentiveData, newData, filterPartial, filterDelete)

	if success && persist {
		r.IncentiveData = data
	}

	return data, success
}

// IncentiveDescriptionListDataType

var _ Updater = (*IncentiveDescriptionListDataType)(nil)

func (r *IncentiveDescriptionListDataType) UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType) (any, bool) {
	var newData []IncentiveDescriptionDataType
	if newList != nil {
		newData = newList.(*IncentiveDescriptionListDataType).IncentiveDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.IncentiveDescriptionData, newData, filterPartial, filterDelete)

	if success && persist {
		r.IncentiveDescriptionData = data
	}

	return data, success
}
