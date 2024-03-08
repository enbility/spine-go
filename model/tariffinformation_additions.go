package model

// TariffListDataType

var _ Updater = (*TariffListDataType)(nil)

func (r *TariffListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []TariffDataType
	if newList != nil {
		newData = newList.(*TariffListDataType).TariffData
	}

	r.TariffData = UpdateList(remoteWrite, r.TariffData, newData, filterPartial, filterDelete)
}

// TariffTierRelationListDataType

var _ Updater = (*TariffTierRelationListDataType)(nil)

func (r *TariffTierRelationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []TariffTierRelationDataType
	if newList != nil {
		newData = newList.(*TariffTierRelationListDataType).TariffTierRelationData
	}

	r.TariffTierRelationData = UpdateList(remoteWrite, r.TariffTierRelationData, newData, filterPartial, filterDelete)
}

// TariffBoundaryRelationListDataType

var _ Updater = (*TariffBoundaryRelationListDataType)(nil)

func (r *TariffBoundaryRelationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []TariffBoundaryRelationDataType
	if newList != nil {
		newData = newList.(*TariffBoundaryRelationListDataType).TariffBoundaryRelationData
	}

	r.TariffBoundaryRelationData = UpdateList(remoteWrite, r.TariffBoundaryRelationData, newData, filterPartial, filterDelete)
}

// TariffDescriptionListDataType

var _ Updater = (*TariffDescriptionListDataType)(nil)

func (r *TariffDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []TariffDescriptionDataType
	if newList != nil {
		newData = newList.(*TariffDescriptionListDataType).TariffDescriptionData
	}

	r.TariffDescriptionData = UpdateList(remoteWrite, r.TariffDescriptionData, newData, filterPartial, filterDelete)
}

// TierBoundaryListDataType

var _ Updater = (*TierBoundaryListDataType)(nil)

func (r *TierBoundaryListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []TierBoundaryDataType
	if newList != nil {
		newData = newList.(*TierBoundaryListDataType).TierBoundaryData
	}

	r.TierBoundaryData = UpdateList(remoteWrite, r.TierBoundaryData, newData, filterPartial, filterDelete)
}

// TierBoundaryDescriptionListDataType

var _ Updater = (*TierBoundaryDescriptionListDataType)(nil)

func (r *TierBoundaryDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []TierBoundaryDescriptionDataType
	if newList != nil {
		newData = newList.(*TierBoundaryDescriptionListDataType).TierBoundaryDescriptionData
	}

	r.TierBoundaryDescriptionData = UpdateList(remoteWrite, r.TierBoundaryDescriptionData, newData, filterPartial, filterDelete)
}

// CommodityListDataType

var _ Updater = (*CommodityListDataType)(nil)

func (r *CommodityListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []CommodityDataType
	if newList != nil {
		newData = newList.(*CommodityListDataType).CommodityData
	}

	r.CommodityData = UpdateList(remoteWrite, r.CommodityData, newData, filterPartial, filterDelete)
}

// TierListDataType

var _ Updater = (*TierListDataType)(nil)

func (r *TierListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []TierDataType
	if newList != nil {
		newData = newList.(*TierListDataType).TierData
	}

	r.TierData = UpdateList(remoteWrite, r.TierData, newData, filterPartial, filterDelete)
}

// TierIncentiveRelationListDataType

var _ Updater = (*TierIncentiveRelationListDataType)(nil)

func (r *TierIncentiveRelationListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []TierIncentiveRelationDataType
	if newList != nil {
		newData = newList.(*TierIncentiveRelationListDataType).TierIncentiveRelationData
	}

	r.TierIncentiveRelationData = UpdateList(remoteWrite, r.TierIncentiveRelationData, newData, filterPartial, filterDelete)
}

// TierDescriptionListDataType

var _ Updater = (*TierDescriptionListDataType)(nil)

func (r *TierDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []TierDescriptionDataType
	if newList != nil {
		newData = newList.(*TierDescriptionListDataType).TierDescriptionData
	}

	r.TierDescriptionData = UpdateList(remoteWrite, r.TierDescriptionData, newData, filterPartial, filterDelete)
}

// IncentiveListDataType

var _ Updater = (*IncentiveListDataType)(nil)

func (r *IncentiveListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []IncentiveDataType
	if newList != nil {
		newData = newList.(*IncentiveListDataType).IncentiveData
	}

	r.IncentiveData = UpdateList(remoteWrite, r.IncentiveData, newData, filterPartial, filterDelete)
}

// IncentiveDescriptionListDataType

var _ Updater = (*IncentiveDescriptionListDataType)(nil)

func (r *IncentiveDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []IncentiveDescriptionDataType
	if newList != nil {
		newData = newList.(*IncentiveDescriptionListDataType).IncentiveDescriptionData
	}

	r.IncentiveDescriptionData = UpdateList(remoteWrite, r.IncentiveDescriptionData, newData, filterPartial, filterDelete)
}
