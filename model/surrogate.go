package model

type SurrogateIdType uint

type SurrogateScopeEnumType string

type SurrogateScopeType string

type SurrogateOtherSourceType string

type SurrogateSourcesType struct {
	SpineSource []EntityAddressType        `json:"spineSource,omitempty"`
	OtherSource []SurrogateOtherSourceType `json:"otherSource,omitempty"`
}

type SurrogateDescriptionDataType struct {
	SurrogateId      *SurrogateIdType      `json:"surrogateId,omitempty" eebus:"key"`
	SurrogateScope   *SurrogateScopeType   `json:"surrogateScope,omitempty"`
	SurrogateSources *SurrogateSourcesType `json:"surrogateSources,omitempty"`
	Label            *LabelType            `json:"label,omitempty"`
	Description      *DescriptionType      `json:"description,omitempty"`
}

type SurrogateSourcesElementsType struct {
	SpineSource *EntityAddressElementsType `json:"spineSource,omitempty"`
	OtherSource *ElementTagType            `json:"otherSource,omitempty"`
}

type SurrogateDescriptionDataElementsType struct {
	SurrogateId      *ElementTagType               `json:"surrogateId,omitempty"`
	SurrogateScope   *ElementTagType               `json:"surrogateScope,omitempty"`
	SurrogateSources *SurrogateSourcesElementsType `json:"surrogateSources,omitempty"`
	Label            *ElementTagType               `json:"label,omitempty"`
	Description      *ElementTagType               `json:"description,omitempty"`
}

type SurrogateDescriptionListDataType struct {
	SurrogateDescriptionData []SurrogateDescriptionDataType `json:"surrogateDescriptionData,omitempty"`
}

type SurrogateDescriptionListDataSelectorsSourcesType struct {
	SpineSource *EntityAddressType `json:"spineSource,omitempty"`
	OtherSource *string            `json:"otherSource,omitempty"`
}

type SurrogateDescriptionListDataSelectorsType struct {
	SurrogateId      *SurrogateIdType                                  `json:"surrogateId,omitempty"`
	SurrogateScope   *SurrogateScopeType                               `json:"surrogateScope,omitempty"`
	SurrogateSources *SurrogateDescriptionListDataSelectorsSourcesType `json:"surrogateSources,omitempty"`
}
