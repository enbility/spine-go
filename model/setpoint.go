package model

type SetpointIdType uint

type SetpointTypeType string

const (
	SetpointTypeTypeValueAbsolute SetpointTypeType = "valueAbsolute"
	SetpointTypeTypeValueRelative SetpointTypeType = "valueRelative"
)

type SetpointDataType struct {
	SetpointId               *SetpointIdType   `json:"setpointId,omitempty" eebus:"key"`
	Value                    *ScaledNumberType `json:"value,omitempty"`
	ValueMin                 *ScaledNumberType `json:"valueMin,omitempty"`
	ValueMax                 *ScaledNumberType `json:"valueMax,omitempty"`
	ValueToleranceAbsolute   *ScaledNumberType `json:"valueToleranceAbsolute,omitempty"`
	ValueTolerancePercentage *ScaledNumberType `json:"valueTolerancePercentage,omitempty"`
	IsSetpointChangeable     *bool             `json:"isSetpointChangeable,omitempty" eebus:"writecheck"`
	IsSetpointActive         *bool             `json:"isSetpointActive,omitempty"`
	TimePeriod               *TimePeriodType   `json:"timePeriod,omitempty"`
}

type SetpointDataElementsType struct {
	SetpointId               *ElementTagType           `json:"setpointId,omitempty"`
	Value                    *ScaledNumberElementsType `json:"value,omitempty"`
	ValueMin                 *ScaledNumberElementsType `json:"valueMin,omitempty"`
	ValueMax                 *ScaledNumberElementsType `json:"valueMax,omitempty"`
	ValueToleranceAbsolute   *ScaledNumberElementsType `json:"valueToleranceAbsolute,omitempty"`
	ValueTolerancePercentage *ScaledNumberElementsType `json:"valueTolerancePercentage,omitempty"`
	IsSetpointChangeable     *ElementTagType           `json:"isSetpointChangeable,omitempty"`
	IsSetpointActive         *ElementTagType           `json:"isSetpointActive,omitempty"`
	TimePeriod               *TimePeriodElementsType   `json:"timePeriod,omitempty"`
}

type SetpointListDataType struct {
	SetpointData []SetpointDataType `json:"setpointData,omitempty"`
}

type SetpointListDataSelectorsType struct {
	SetpointId *SetpointIdType `json:"setpointId,omitempty"`
}

type SetpointConstraintsDataType struct {
	SetpointId       *SetpointIdType   `json:"setpointId,omitempty" eebus:"key"`
	SetpointRangeMin *ScaledNumberType `json:"setpointRangeMin,omitempty"`
	SetpointRangeMax *ScaledNumberType `json:"setpointRangeMax,omitempty"`
	SetpointStepSize *ScaledNumberType `json:"setpointStepSize,omitempty"`
}

type SetpointConstraintsDataElementsType struct {
	SetpointId       *ElementTagType           `json:"setpointId,omitempty"`
	SetpointRangeMin *ScaledNumberElementsType `json:"setpointRangeMin,omitempty"`
	SetpointRangeMax *ScaledNumberElementsType `json:"setpointRangeMax,omitempty"`
	SetpointStepSize *ScaledNumberElementsType `json:"setpointStepSize,omitempty"`
}

type SetpointConstraintsListDataType struct {
	SetpointConstraintsData []SetpointConstraintsDataType `json:"setpointConstraintsData,omitempty"`
}

type SetpointConstraintsListDataSelectorsType struct {
	SetpointId *SetpointIdType `json:"setpointId,omitempty"`
}

type SetpointDescriptionDataType struct {
	SetpointId    *SetpointIdType        `json:"setpointId,omitempty" eebus:"key"`
	MeasurementId *SetpointIdType        `json:"measurementId,omitempty" eebus:"key"`
	TimeTableId   *SetpointIdType        `json:"timeTableId,omitempty" eebus:"key"`
	SetpointType  *SetpointTypeType      `json:"setpointType,omitempty"`
	Unit          *UnitOfMeasurementType `json:"unit,omitempty"`
	ScopeType     *ScopeTypeType         `json:"scopeType,omitempty"`
	Label         *LabelType             `json:"label,omitempty"`
	Description   *DescriptionType       `json:"description,omitempty"`
}

type SetpointDescriptionDataElementsType struct {
	SetpointId    *ElementTagType `json:"setpointId,omitempty"`
	MeasurementId *ElementTagType `json:"measurementId,omitempty"`
	TimeTableId   *ElementTagType `json:"timeTableId,omitempty"`
	SetpointType  *ElementTagType `json:"setpointType,omitempty"`
	Unit          *ElementTagType `json:"unit,omitempty"`
	ScopeType     *ElementTagType `json:"scopeType,omitempty"`
	Label         *ElementTagType `json:"label,omitempty"`
	Description   *ElementTagType `json:"description,omitempty"`
}

type SetpointDescriptionListDataType struct {
	SetpointDescriptionData []SetpointDescriptionDataType `json:"setpointDescriptionData,omitempty"`
}

type SetpointDescriptionListDataSelectorsType struct {
	SetpointId    *SetpointIdType    `json:"setpointId,omitempty"`
	MeasurementId *MeasurementIdType `json:"measurementId,omitempty"`
	TimeTableId   *TimeTableIdType   `json:"timeTableId,omitempty"`
	SetpointType  *SetpointTypeType  `json:"setpointType,omitempty"`
	ScopeType     *ScopeTypeType     `json:"scopeType,omitempty"`
}
