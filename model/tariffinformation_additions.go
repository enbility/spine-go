package model

// TariffListDataType

var _ Updater = (*TariffListDataType)(nil)

func (r *TariffListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []TariffDataType
	if newList != nil {
		newData = newList.(*TariffListDataType).TariffData
	}

	data, success := UpdateList(remoteWrite, r.TariffData, newData, filterPartial, filterDelete)

	if success {
		r.TariffData = data
	}

	return success
}

// TariffTierRelationListDataType

var _ Updater = (*TariffTierRelationListDataType)(nil)

func (r *TariffTierRelationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []TariffTierRelationDataType
	if newList != nil {
		newData = newList.(*TariffTierRelationListDataType).TariffTierRelationData
	}

	data, success := UpdateList(remoteWrite, r.TariffTierRelationData, newData, filterPartial, filterDelete)

	if success {
		r.TariffTierRelationData = data
	}

	return success
}

// TariffBoundaryRelationListDataType

var _ Updater = (*TariffBoundaryRelationListDataType)(nil)

func (r *TariffBoundaryRelationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []TariffBoundaryRelationDataType
	if newList != nil {
		newData = newList.(*TariffBoundaryRelationListDataType).TariffBoundaryRelationData
	}

	data, success := UpdateList(remoteWrite, r.TariffBoundaryRelationData, newData, filterPartial, filterDelete)

	if success {
		r.TariffBoundaryRelationData = data
	}

	return success
}

// TariffDescriptionListDataType

var _ Updater = (*TariffDescriptionListDataType)(nil)

func (r *TariffDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []TariffDescriptionDataType
	if newList != nil {
		newData = newList.(*TariffDescriptionListDataType).TariffDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.TariffDescriptionData, newData, filterPartial, filterDelete)

	if success {
		r.TariffDescriptionData = data
	}

	return success
}

// TierBoundaryListDataType

var _ Updater = (*TierBoundaryListDataType)(nil)

func (r *TierBoundaryListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []TierBoundaryDataType
	if newList != nil {
		newData = newList.(*TierBoundaryListDataType).TierBoundaryData
	}

	data, success := UpdateList(remoteWrite, r.TierBoundaryData, newData, filterPartial, filterDelete)

	if success {
		r.TierBoundaryData = data
	}

	return success
}

// TierBoundaryDescriptionListDataType

var _ Updater = (*TierBoundaryDescriptionListDataType)(nil)

func (r *TierBoundaryDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []TierBoundaryDescriptionDataType
	if newList != nil {
		newData = newList.(*TierBoundaryDescriptionListDataType).TierBoundaryDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.TierBoundaryDescriptionData, newData, filterPartial, filterDelete)

	if success {
		r.TierBoundaryDescriptionData = data
	}

	return success
}

// CommodityListDataType

var _ Updater = (*CommodityListDataType)(nil)

func (r *CommodityListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []CommodityDataType
	if newList != nil {
		newData = newList.(*CommodityListDataType).CommodityData
	}

	data, success := UpdateList(remoteWrite, r.CommodityData, newData, filterPartial, filterDelete)

	if success {
		r.CommodityData = data
	}

	return success
}

// TierListDataType

var _ Updater = (*TierListDataType)(nil)

func (r *TierListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []TierDataType
	if newList != nil {
		newData = newList.(*TierListDataType).TierData
	}

	data, success := UpdateList(remoteWrite, r.TierData, newData, filterPartial, filterDelete)

	if success {
		r.TierData = data
	}

	return success
}

// TierIncentiveRelationListDataType

var _ Updater = (*TierIncentiveRelationListDataType)(nil)

func (r *TierIncentiveRelationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []TierIncentiveRelationDataType
	if newList != nil {
		newData = newList.(*TierIncentiveRelationListDataType).TierIncentiveRelationData
	}

	data, success := UpdateList(remoteWrite, r.TierIncentiveRelationData, newData, filterPartial, filterDelete)

	if success {
		r.TierIncentiveRelationData = data
	}

	return success
}

// TierDescriptionListDataType

var _ Updater = (*TierDescriptionListDataType)(nil)

func (r *TierDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []TierDescriptionDataType
	if newList != nil {
		newData = newList.(*TierDescriptionListDataType).TierDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.TierDescriptionData, newData, filterPartial, filterDelete)

	if success {
		r.TierDescriptionData = data
	}

	return success
}

// IncentiveListDataType

var _ Updater = (*IncentiveListDataType)(nil)

func (r *IncentiveListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []IncentiveDataType
	if newList != nil {
		newData = newList.(*IncentiveListDataType).IncentiveData
	}

	data, success := UpdateList(remoteWrite, r.IncentiveData, newData, filterPartial, filterDelete)

	if success {
		r.IncentiveData = data
	}

	return success
}

// IncentiveDescriptionListDataType

var _ Updater = (*IncentiveDescriptionListDataType)(nil)

func (r *IncentiveDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool {
	var newData []IncentiveDescriptionDataType
	if newList != nil {
		newData = newList.(*IncentiveDescriptionListDataType).IncentiveDescriptionData
	}

	data, success := UpdateList(remoteWrite, r.IncentiveDescriptionData, newData, filterPartial, filterDelete)

	if success {
		r.IncentiveDescriptionData = data
	}

	return success
}
