package model

import (
	"encoding/json"
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestElectricalConnectionStateListDataType_Update(t *testing.T) {
	sut := ElectricalConnectionStateListDataType{
		ElectricalConnectionStateData: []ElectricalConnectionStateDataType{
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				CurrentEnergyMode:      util.Ptr(EnergyModeTypeProduce),
			},
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(1)),
				CurrentEnergyMode:      util.Ptr(EnergyModeTypeProduce),
			},
		},
	}

	newData := ElectricalConnectionStateListDataType{
		ElectricalConnectionStateData: []ElectricalConnectionStateDataType{
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(1)),
				CurrentEnergyMode:      util.Ptr(EnergyModeTypeConsume),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.ElectricalConnectionStateData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, EnergyModeTypeProduce, *item1.CurrentEnergyMode)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ElectricalConnectionId))
	assert.Equal(t, EnergyModeTypeConsume, *item2.CurrentEnergyMode)
}

// verifies that a subset of existing items will be updated with identified new values
func TestElectricalConnectionPermittedValueSetListDataType_Update_Modify(t *testing.T) {
	sut := ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(0)),
				PermittedValueSet: []ScaledNumberSetType{
					{
						Range: []ScaledNumberRangeType{
							{
								Min: NewScaledNumberType(1),
							},
						},
					},
				},
			},
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(1)),
				PermittedValueSet: []ScaledNumberSetType{
					{
						Range: []ScaledNumberRangeType{
							{
								Min: NewScaledNumberType(6),
								Max: NewScaledNumberType(16),
							},
						},
					},
				},
			},
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(2)),
				PermittedValueSet: []ScaledNumberSetType{
					{
						Range: []ScaledNumberRangeType{
							{
								Min: NewScaledNumberType(6),
								Max: NewScaledNumberType(16),
							},
						},
					},
				},
			},
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(3)),
				PermittedValueSet: []ScaledNumberSetType{
					{
						Range: []ScaledNumberRangeType{
							{
								Min: NewScaledNumberType(6),
								Max: NewScaledNumberType(16),
							},
						},
					},
				},
			},
		},
	}

	newData := ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(1)),
				PermittedValueSet: []ScaledNumberSetType{
					{
						Range: []ScaledNumberRangeType{
							{
								Min: NewScaledNumberType(2),
								Max: NewScaledNumberType(16),
							},
						},
					},
				},
			},
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(2)),
				PermittedValueSet: []ScaledNumberSetType{
					{
						Range: []ScaledNumberRangeType{
							{
								Min: NewScaledNumberType(2),
								Max: NewScaledNumberType(16),
							},
						},
					},
				},
			},
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(3)),
				PermittedValueSet: []ScaledNumberSetType{
					{
						Range: []ScaledNumberRangeType{
							{
								Min: NewScaledNumberType(2),
								Max: NewScaledNumberType(16),
							},
						},
					},
				},
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.ElectricalConnectionPermittedValueSetData
	// check the non changing items
	assert.Equal(t, 4, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, 0, int(*item1.ParameterId))
	assert.Equal(t, 1, len(item1.PermittedValueSet))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 0, int(*item2.ElectricalConnectionId))
	assert.Equal(t, 1, int(*item2.ParameterId))
	assert.Equal(t, 1, len(item2.PermittedValueSet))
	valueSet := item2.PermittedValueSet[0]
	assert.Equal(t, 1, len(valueSet.Range))
	rangeSet := valueSet.Range[0]
	assert.Equal(t, 2.0, rangeSet.Min.GetValue())
	assert.Equal(t, 16.0, rangeSet.Max.GetValue())
}

// verifies that a subset of existing items will be updated with identified new values
func TestElectricalConnectionPermittedValueSetListDataType_Update_Modify_Selector(t *testing.T) {
	existingDataJson := `{
		"electricalConnectionPermittedValueSetData":[
			{
				"electricalConnectionId":0,
				"parameterId":0,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":1,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":1,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":2,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":3,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			}
		]
	}`

	var sut ElectricalConnectionPermittedValueSetListDataType
	err := json.Unmarshal([]byte(existingDataJson), &sut)
	if assert.Nil(t, err) == false {
		return
	}

	newDataJson := `{
		"electricalConnectionPermittedValueSetData":[
			{
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":2,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			}
		]
	}`

	var newData ElectricalConnectionPermittedValueSetListDataType
	err = json.Unmarshal([]byte(newDataJson), &newData)
	if assert.Nil(t, err) == false {
		return
	}

	partial := &FilterType{
		CmdControl: &CmdControlType{
			Partial: &ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetListDataSelectors: &ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ElectricalConnectionId: util.Ptr[ElectricalConnectionIdType](0),
			ParameterId:            util.Ptr[ElectricalConnectionParameterIdType](1),
		},
	}

	// Act
	_, success := sut.UpdateList(false, false, &newData, partial, nil, nil)
	assert.True(t, success)

	data := sut.ElectricalConnectionPermittedValueSetData
	// check the non changing items
	assert.Equal(t, 4, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, 0, int(*item1.ParameterId))
	assert.Equal(t, 1, len(item1.PermittedValueSet))
	item3 := data[2]
	assert.Equal(t, 0, int(*item3.ElectricalConnectionId))
	assert.Equal(t, 2, int(*item3.ParameterId))
	assert.Equal(t, 1, len(item3.PermittedValueSet))
	valueSet := item3.PermittedValueSet[0]
	assert.Equal(t, 1, len(valueSet.Range))
	rangeSet := valueSet.Range[0]
	assert.Equal(t, 6.0, rangeSet.Min.GetValue())
	assert.Equal(t, 16.0, rangeSet.Max.GetValue())

	// check properties of updated item
	item2 := sut.ElectricalConnectionPermittedValueSetData[1]
	assert.Equal(t, 0, int(*item2.ElectricalConnectionId))
	assert.Equal(t, 1, int(*item2.ParameterId))
	assert.Equal(t, 1, len(item2.PermittedValueSet))
	valueSet = item2.PermittedValueSet[0]
	assert.Equal(t, 1, len(valueSet.Range))
	rangeSet = valueSet.Range[0]
	assert.Equal(t, 2.0, rangeSet.Min.GetValue())
	assert.Equal(t, 16.0, rangeSet.Max.GetValue())
}

// verifies that a subset of existing items will be updated with identified new values
func TestElectricalConnectionPermittedValueSetListDataType_Update_Delete_Modify(t *testing.T) {
	existingDataJson := `{
		"electricalConnectionPermittedValueSetData":[
			{
				"electricalConnectionId":0,
				"parameterId":0,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":1,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":1,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":2,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":3,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			}
		]
	}`

	var sut ElectricalConnectionPermittedValueSetListDataType
	err := json.Unmarshal([]byte(existingDataJson), &sut)
	if assert.Nil(t, err) == false {
		return
	}

	newDataJson := `{
		"electricalConnectionPermittedValueSetData":[
			{
				"electricalConnectionId":0,
				"parameterId":1,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":2,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":2,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":2,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":3,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":2,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			}
		]
	}`

	var newData ElectricalConnectionPermittedValueSetListDataType
	err = json.Unmarshal([]byte(newDataJson), &newData)
	if assert.Nil(t, err) == false {
		return
	}

	deleteFilter := &FilterType{
		CmdControl: &CmdControlType{
			Delete: &ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetListDataSelectors: &ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ElectricalConnectionId: util.Ptr[ElectricalConnectionIdType](0),
			ParameterId:            util.Ptr[ElectricalConnectionParameterIdType](0),
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), deleteFilter, nil)
	assert.True(t, success)

	data := sut.ElectricalConnectionPermittedValueSetData
	// check the deleted item is gone
	assert.Equal(t, 3, len(data))
	// check properties of updated item
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, 1, int(*item1.ParameterId))
	assert.Equal(t, 1, len(item1.PermittedValueSet))
	valueSet := item1.PermittedValueSet[0]
	assert.Equal(t, 1, len(valueSet.Range))
	rangeSet := valueSet.Range[0]
	assert.Equal(t, 2.0, rangeSet.Min.GetValue())
	assert.Equal(t, 16.0, rangeSet.Max.GetValue())
}

// verifies that a subset of existing items will be updated with identified new values
func TestElectricalConnectionPermittedValueSetListDataType_Update_Delete(t *testing.T) {
	existingDataJson := `{
		"electricalConnectionPermittedValueSetData":[
			{
				"electricalConnectionId":0,
				"parameterId":0,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":1,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":1,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":2,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":3,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			}
		]
	}`

	var sut ElectricalConnectionPermittedValueSetListDataType
	err := json.Unmarshal([]byte(existingDataJson), &sut)
	if assert.Nil(t, err) == false {
		return
	}

	deleteFilter := &FilterType{
		CmdControl: &CmdControlType{
			Delete: &ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetListDataSelectors: &ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ElectricalConnectionId: util.Ptr[ElectricalConnectionIdType](0),
			ParameterId:            util.Ptr[ElectricalConnectionParameterIdType](0),
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, nil, nil, deleteFilter, nil)
	assert.True(t, success)

	data := sut.ElectricalConnectionPermittedValueSetData
	// check the deleted item is added again
	assert.Equal(t, 3, len(data))
	// check properties of remaining item
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, 1, int(*item1.ParameterId))
	assert.Equal(t, 1, len(item1.PermittedValueSet))
	valueSet := item1.PermittedValueSet[0]
	assert.Equal(t, 1, len(valueSet.Range))
	rangeSet := valueSet.Range[0]
	assert.Equal(t, 6.0, rangeSet.Min.GetValue())
	assert.Equal(t, 16.0, rangeSet.Max.GetValue())
}

// verifies that a subset of existing items will be updated with identified new values
func TestElectricalConnectionPermittedValueSetListDataType_Update_Delete_Element(t *testing.T) {
	existingDataJson := `{
		"electricalConnectionPermittedValueSetData":[
			{
				"electricalConnectionId":0,
				"parameterId":0,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":1,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":1,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":2,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":3,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			}
		]
	}`

	var sut ElectricalConnectionPermittedValueSetListDataType
	err := json.Unmarshal([]byte(existingDataJson), &sut)
	if assert.Nil(t, err) == false {
		return
	}

	deleteFilter := &FilterType{
		CmdControl: &CmdControlType{
			Delete: &ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetDataElements: &ElectricalConnectionPermittedValueSetDataElementsType{
			PermittedValueSet: &ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetListDataSelectors: &ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ElectricalConnectionId: util.Ptr[ElectricalConnectionIdType](0),
			ParameterId:            util.Ptr[ElectricalConnectionParameterIdType](0),
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, nil, nil, deleteFilter, nil)
	assert.True(t, success)

	data := sut.ElectricalConnectionPermittedValueSetData
	// check no items are deleted
	assert.Equal(t, 4, len(data))
	// check permitted value is removed from item with ID 0
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, 0, int(*item1.ParameterId))
	var nilValue []ScaledNumberSetType
	assert.Equal(t, nilValue, item1.PermittedValueSet)

	// check properties of remaining item
	item2 := data[1]
	assert.Equal(t, 0, int(*item2.ElectricalConnectionId))
	assert.Equal(t, 1, int(*item2.ParameterId))
	assert.Equal(t, 1, len(item2.PermittedValueSet))
	valueSet := item2.PermittedValueSet[0]
	assert.Equal(t, 1, len(valueSet.Range))
	rangeSet := valueSet.Range[0]
	assert.Equal(t, 6.0, rangeSet.Min.GetValue())
	assert.Equal(t, 16.0, rangeSet.Max.GetValue())
}

// verifies that a subset of existing items will be updated with identified new values
func TestElectricalConnectionPermittedValueSetListDataType_Update_Delete_OnlyElement(t *testing.T) {
	existingDataJson := `{
		"electricalConnectionPermittedValueSetData":[
			{
				"electricalConnectionId":0,
				"parameterId":0,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":1,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":1,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":2,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":3,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			}
		]
	}`

	var sut ElectricalConnectionPermittedValueSetListDataType
	err := json.Unmarshal([]byte(existingDataJson), &sut)
	if assert.Nil(t, err) == false {
		return
	}

	deleteFilter := &FilterType{
		CmdControl: &CmdControlType{
			Delete: &ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetDataElements: &ElectricalConnectionPermittedValueSetDataElementsType{
			PermittedValueSet: &ElementTagType{},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, nil, nil, deleteFilter, nil)
	assert.True(t, success)

	data := sut.ElectricalConnectionPermittedValueSetData
	// check no items are deleted
	assert.Equal(t, 4, len(data))
	// check permitted value is removed from item with ID 0
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, 0, int(*item1.ParameterId))
	var nilValue []ScaledNumberSetType
	assert.Equal(t, nilValue, item1.PermittedValueSet)

	// check properties
	item2 := data[1]
	assert.Equal(t, 0, int(*item2.ElectricalConnectionId))
	assert.Equal(t, 1, int(*item2.ParameterId))
	assert.Equal(t, nilValue, item2.PermittedValueSet)

	item3 := data[2]
	assert.Equal(t, 0, int(*item3.ElectricalConnectionId))
	assert.Equal(t, 2, int(*item3.ParameterId))
	assert.Equal(t, nilValue, item3.PermittedValueSet)

	item4 := data[3]
	assert.Equal(t, 0, int(*item4.ElectricalConnectionId))
	assert.Equal(t, 3, int(*item4.ParameterId))
	assert.Equal(t, nilValue, item4.PermittedValueSet)
}

// verifies that a subset of existing items will be updated with identified new values
func TestElectricalConnectionPermittedValueSetListDataType_Update_Delete_Add(t *testing.T) {
	existingDataJson := `{
		"electricalConnectionPermittedValueSetData":[
			{
				"electricalConnectionId":0,
				"parameterId":0,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":1,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":1,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":2,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":3,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			}
		]
	}`

	var sut ElectricalConnectionPermittedValueSetListDataType
	err := json.Unmarshal([]byte(existingDataJson), &sut)
	if assert.Nil(t, err) == false {
		return
	}

	newDataJson := `{
		"electricalConnectionPermittedValueSetData":[
			{
				"electricalConnectionId":0,
				"parameterId":0,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":1,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":1,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":2,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":2,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":2,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":3,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":2,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			}
		]
	}`

	var newData ElectricalConnectionPermittedValueSetListDataType
	err = json.Unmarshal([]byte(newDataJson), &newData)
	if assert.Nil(t, err) == false {
		return
	}

	deleteFilter := &FilterType{
		CmdControl: &CmdControlType{
			Delete: &ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetListDataSelectors: &ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ElectricalConnectionId: util.Ptr[ElectricalConnectionIdType](0),
			ParameterId:            util.Ptr[ElectricalConnectionParameterIdType](0),
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), deleteFilter, nil)
	assert.True(t, success)

	data := sut.ElectricalConnectionPermittedValueSetData
	// check the deleted item is added again
	assert.Equal(t, 4, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, 0, int(*item1.ParameterId))
	assert.Equal(t, 1, len(item1.PermittedValueSet))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 0, int(*item2.ElectricalConnectionId))
	assert.Equal(t, 1, int(*item2.ParameterId))
	assert.Equal(t, 1, len(item2.PermittedValueSet))
	valueSet := item2.PermittedValueSet[0]
	assert.Equal(t, 1, len(valueSet.Range))
	rangeSet := valueSet.Range[0]
	assert.Equal(t, 2.0, rangeSet.Min.GetValue())
	assert.Equal(t, 16.0, rangeSet.Max.GetValue())
}

// verifies that an item in the payload which is not in the existing data will be added
func TestElectricalConnectionPermittedValueSetListDataType_Update_NewItem(t *testing.T) {
	existingDataJson := `{
		"electricalConnectionPermittedValueSetData": [
		  {
			"electricalConnectionId": 1,
			"parameterId": 1,
			"permittedValueSet": [
			  {
				"range": [
				  {
					"min": { "number": 3, "scale": 0 },
					"max": { "number": 6, "scale": 0 }
				  }
				]
			  }
			]
		  }
		]
	}`

	var sut ElectricalConnectionPermittedValueSetListDataType
	err := json.Unmarshal([]byte(existingDataJson), &sut)
	if assert.Nil(t, err) == false {
		return
	}

	newDataJson := `{
		"electricalConnectionPermittedValueSetData": [
		  {
			"electricalConnectionId": 1,
			"parameterId": 2,
			"permittedValueSet": [
			  {
				"range": [
				  {
					"min": { "number": 9, "scale": 0 },
					"max": { "number": 19, "scale": 0 }
				  }
				]
			  },
			  {
				"range": [
				  {
					"min": { "number": 30, "scale": 0 },
					"max": { "number": 36, "scale": 0 }
				  }
				]
			  }
			]
		  }
		]
	}`

	var newData ElectricalConnectionPermittedValueSetListDataType
	err = json.Unmarshal([]byte(newDataJson), &newData)
	if assert.Nil(t, err) == false {
		return
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.ElectricalConnectionPermittedValueSetData
	// new item should be added
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 1, int(*item1.ElectricalConnectionId))
	assert.Equal(t, 1, int(*item1.ParameterId))
	assert.Equal(t, 1, len(item1.PermittedValueSet))
	// check properties of added item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ElectricalConnectionId))
	assert.Equal(t, 2, int(*item2.ParameterId))
	assert.Equal(t, 2, len(item2.PermittedValueSet))
}

// verifies that an item in the payload which has no identifiers will be copied to all existing data
// (see EEBus_SPINE_TS_ProtocolSpecification.pdf, Table 7: Considered cmdOptions combinations for classifier "notify")
func TestElectricalConnectionPermittedValueSetListDataType_UpdateWithoutIdenifiers(t *testing.T) {
	existingDataJson := `{
		"electricalConnectionPermittedValueSetData": [
		  {
			"electricalConnectionId": 1,
			"parameterId": 1,
			"permittedValueSet": [
			  {
				"range": [
				  {
					"min": { "number": 3, "scale": 0 },
					"max": { "number": 6, "scale": 0 }
				  }
				]
			  }
			]
		  },
		  {
			"electricalConnectionId": 1,
			"parameterId": 2,
			"permittedValueSet": [
			  {
				"range": [
				  {
					"min": { "number": 6, "scale": 0 },
					"max": { "number": 12, "scale": 0 }
				  }
				]
			  }
			]
		  }		]
	}`

	var sut ElectricalConnectionPermittedValueSetListDataType
	err := json.Unmarshal([]byte(existingDataJson), &sut)
	if assert.Nil(t, err) == false {
		return
	}

	// item with no identifiers
	newDataJson := `{
		"electricalConnectionPermittedValueSetData": [
		  {
			"permittedValueSet": [
			  {
				"range": [
				  {
					"min": { "number": 30, "scale": 0 },
					"max": { "number": 36, "scale": 0 }
				  }
				]
			  }
			]
		  }
		]
	}`

	var newData ElectricalConnectionPermittedValueSetListDataType
	err = json.Unmarshal([]byte(newDataJson), &newData)
	if assert.Nil(t, err) == false {
		return
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.ElectricalConnectionPermittedValueSetData
	// the new item should not be added
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 1, int(*item1.ElectricalConnectionId))
	assert.Equal(t, 1, int(*item1.ParameterId))
	assert.Equal(t, 1, len(item1.PermittedValueSet))
	valueSet := item1.PermittedValueSet[0]
	assert.Equal(t, 1, len(valueSet.Range))
	// the values of the item in the payload should be copied to the first item
	assert.Equal(t, 30, int(*valueSet.Range[0].Min.Number))
	assert.Equal(t, 0, int(*valueSet.Range[0].Min.Scale))
	assert.Equal(t, 36, int(*valueSet.Range[0].Max.Number))
	assert.Equal(t, 0, int(*valueSet.Range[0].Max.Scale))

	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ElectricalConnectionId))
	assert.Equal(t, 2, int(*item2.ParameterId))
	assert.Equal(t, 1, len(item2.PermittedValueSet))
	valueSet = item2.PermittedValueSet[0]
	assert.Equal(t, 1, len(valueSet.Range))
	// the values of the item in the payload should be also copied to the second item
	assert.Equal(t, 30, int(*valueSet.Range[0].Min.Number))
	assert.Equal(t, 0, int(*valueSet.Range[0].Min.Scale))
	assert.Equal(t, 36, int(*valueSet.Range[0].Max.Number))
	assert.Equal(t, 0, int(*valueSet.Range[0].Max.Scale))
}

func TestElectricalConnectionDescriptionListDataType_Update(t *testing.T) {
	sut := ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				PowerSupplyType:        util.Ptr(ElectricalConnectionVoltageTypeTypeAc),
			},
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(1)),
				PowerSupplyType:        util.Ptr(ElectricalConnectionVoltageTypeTypeAc),
			},
		},
	}

	newData := ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(1)),
				PowerSupplyType:        util.Ptr(ElectricalConnectionVoltageTypeTypeDc),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.ElectricalConnectionDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, ElectricalConnectionVoltageTypeTypeAc, *item1.PowerSupplyType)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ElectricalConnectionId))
	assert.Equal(t, ElectricalConnectionVoltageTypeTypeDc, *item2.PowerSupplyType)
}

func TestElectricalConnectionCharacteristicListDataType_Update(t *testing.T) {
	sut := ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []ElectricalConnectionCharacteristicDataType{
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(0)),
				CharacteristicId:       util.Ptr(ElectricalConnectionCharacteristicIdType(0)),
				CharacteristicType:     util.Ptr(ElectricalConnectionCharacteristicTypeTypeApparentPowerConsumptionNominalMax),
			},
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(1)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(0)),
				CharacteristicId:       util.Ptr(ElectricalConnectionCharacteristicIdType(1)),
				CharacteristicType:     util.Ptr(ElectricalConnectionCharacteristicTypeTypePowerConsumptionMax),
			},
		},
	}

	newData := ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []ElectricalConnectionCharacteristicDataType{
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(1)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(0)),
				CharacteristicId:       util.Ptr(ElectricalConnectionCharacteristicIdType(1)),
				CharacteristicType:     util.Ptr(ElectricalConnectionCharacteristicTypeTypeEnergyCapacityNominalMax),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.ElectricalConnectionCharacteristicData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, ElectricalConnectionCharacteristicTypeTypeApparentPowerConsumptionNominalMax, *item1.CharacteristicType)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ElectricalConnectionId))
	assert.Equal(t, ElectricalConnectionCharacteristicTypeTypeEnergyCapacityNominalMax, *item2.CharacteristicType)
}

func TestElectricalConnectionParameterDescriptionListDataType_Update(t *testing.T) {
	sut := ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(0)),
				VoltageType:            util.Ptr(ElectricalConnectionVoltageTypeTypeAc),
			},
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(1)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(0)),
				MeasurementId:          util.Ptr(MeasurementIdType(0)),
				VoltageType:            util.Ptr(ElectricalConnectionVoltageTypeTypeAc),
			},
		},
	}

	newData := ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(1)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(0)),
				VoltageType:            util.Ptr(ElectricalConnectionVoltageTypeTypeDc),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.ElectricalConnectionParameterDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, ElectricalConnectionVoltageTypeTypeAc, *item1.VoltageType)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ElectricalConnectionId))
	assert.Equal(t, ElectricalConnectionVoltageTypeTypeDc, *item2.VoltageType)
}

func TestElectricalConnectionStateListDataType_ReadPartialData(t *testing.T) {
	// Create test data
	testData := &ElectricalConnectionStateListDataType{
		ElectricalConnectionStateData: []ElectricalConnectionStateDataType{
			{
				ElectricalConnectionId: Ptr(ElectricalConnectionIdType(1)),
			},
			{
				ElectricalConnectionId: Ptr(ElectricalConnectionIdType(2)),
			},
		},
	}

	// Test 1: No filter - should return alludata
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*ElectricalConnectionStateListDataType)
	assert.Len(t, resultData.ElectricalConnectionStateData, 2)

	// Test 2: Filter with specific electrical connection ID
	filter := &FilterType{
		ElectricalConnectionStateListDataSelectors: &ElectricalConnectionStateListDataSelectorsType{
			ElectricalConnectionId: Ptr(ElectricalConnectionIdType(1)),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*ElectricalConnectionStateListDataType)
	assert.Len(t, resultData.ElectricalConnectionStateData, 1)
	assert.Equal(t, ElectricalConnectionIdType(1), *resultData.ElectricalConnectionStateData[0].ElectricalConnectionId)

	// Test 3: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestElectricalConnectionParameterDescriptionListDataType_ReadPartialData(t *testing.T) {
	connectionId := ElectricalConnectionIdType(1)
	parameterId := ElectricalConnectionParameterIdType(10)

	data := &ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []ElectricalConnectionParameterDescriptionDataType{
			{ElectricalConnectionId: &connectionId, ParameterId: &parameterId},
			{ElectricalConnectionId: Ptr(ElectricalConnectionIdType(2)), ParameterId: &parameterId},
		},
	}

	// Test with both connection and parameter ID filter
	filter := &FilterType{
		ElectricalConnectionParameterDescriptionListDataSelectors: &ElectricalConnectionParameterDescriptionListDataSelectorsType{
			ElectricalConnectionId: &connectionId,
			ParameterId:            &parameterId,
		},
	}

	result, ok := data.ReadPartialData(filter)
	assert.True(t, ok)

	resultData := result.(*ElectricalConnectionParameterDescriptionListDataType)
	assert.Len(t, resultData.ElectricalConnectionParameterDescriptionData, 1)
	assert.Equal(t, connectionId, *resultData.ElectricalConnectionParameterDescriptionData[0].ElectricalConnectionId)

	// Test 2: ReadPartialData with nil filter
	result, ok = data.ReadPartialData(nil)
	assert.True(t, ok)
	resultData = result.(*ElectricalConnectionParameterDescriptionListDataType)
	assert.Len(t, resultData.ElectricalConnectionParameterDescriptionData, 2)

	// Test 2: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success := data.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestElectricalConnectionPermittedValueSetListDataType_ReadPartialData(t *testing.T) {
	electricalConnectionId1 := ElectricalConnectionIdType(1)
	electricalConnectionId2 := ElectricalConnectionIdType(2)
	parameterId1 := ElectricalConnectionParameterIdType(10)
	parameterId2 := ElectricalConnectionParameterIdType(20)

	testData := &ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: &electricalConnectionId1,
				ParameterId:            &parameterId1,
				PermittedValueSet: []ScaledNumberSetType{
					{
						Value: []ScaledNumberType{
							{Number: util.Ptr(NumberType(100)), Scale: util.Ptr(ScaleType(0))},
							{Number: util.Ptr(NumberType(200)), Scale: util.Ptr(ScaleType(0))},
						},
					},
				},
			},
			{
				ElectricalConnectionId: &electricalConnectionId2,
				ParameterId:            &parameterId2,
				PermittedValueSet: []ScaledNumberSetType{
					{
						Value: []ScaledNumberType{
							{Number: util.Ptr(NumberType(300)), Scale: util.Ptr(ScaleType(0))},
						},
					},
				},
			},
		},
	}

	// Test 1: No filter - should return all data
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, testData, result)

	// Test 2: Filter by electrical connection ID
	filter := &FilterType{
		ElectricalConnectionPermittedValueSetListDataSelectors: &ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ElectricalConnectionId: &electricalConnectionId1,
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData := result.(*ElectricalConnectionPermittedValueSetListDataType)
	assert.Len(t, resultData.ElectricalConnectionPermittedValueSetData, 1)
	assert.Equal(t, electricalConnectionId1, *resultData.ElectricalConnectionPermittedValueSetData[0].ElectricalConnectionId)

	// Test 3: Filter by parameter ID
	filter = &FilterType{
		ElectricalConnectionPermittedValueSetListDataSelectors: &ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ParameterId: &parameterId2,
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*ElectricalConnectionPermittedValueSetListDataType)
	assert.Len(t, resultData.ElectricalConnectionPermittedValueSetData, 1)
	assert.Equal(t, parameterId2, *resultData.ElectricalConnectionPermittedValueSetData[0].ParameterId)

	// Test 4: Multiple criteria filter
	filter = &FilterType{
		ElectricalConnectionPermittedValueSetListDataSelectors: &ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ElectricalConnectionId: &electricalConnectionId1,
			ParameterId:            &parameterId1,
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*ElectricalConnectionPermittedValueSetListDataType)
	assert.Len(t, resultData.ElectricalConnectionPermittedValueSetData, 1)
	assert.Equal(t, electricalConnectionId1, *resultData.ElectricalConnectionPermittedValueSetData[0].ElectricalConnectionId)
	assert.Equal(t, parameterId1, *resultData.ElectricalConnectionPermittedValueSetData[0].ParameterId)

	// Test 5: No matching filter
	nonExistentId := ElectricalConnectionIdType(999)
	filter = &FilterType{
		ElectricalConnectionPermittedValueSetListDataSelectors: &ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ElectricalConnectionId: &nonExistentId,
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*ElectricalConnectionPermittedValueSetListDataType)
	assert.Len(t, resultData.ElectricalConnectionPermittedValueSetData, 0)

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestElectricalConnectionDescriptionListDataType_ReadPartialData(t *testing.T) {
	electricalConnectionId1 := ElectricalConnectionIdType(1)
	electricalConnectionId2 := ElectricalConnectionIdType(2)
	scopeType1 := ScopeTypeType("device")
	scopeType2 := ScopeTypeType("system")

	testData := &ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId: &electricalConnectionId1,
				PowerSupplyType:        util.Ptr(ElectricalConnectionVoltageTypeType("ac")),
				ScopeType:              &scopeType1,
				Label:                  util.Ptr(LabelType("AC Connection 1")),
				Description:            util.Ptr(DescriptionType("Primary AC connection")),
			},
			{
				ElectricalConnectionId: &electricalConnectionId2,
				PowerSupplyType:        util.Ptr(ElectricalConnectionVoltageTypeType("dc")),
				ScopeType:              &scopeType2,
				Label:                  util.Ptr(LabelType("DC Connection 1")),
				Description:            util.Ptr(DescriptionType("Primary DC connection")),
			},
		},
	}

	// Test 1: No filter - should return all data
	result, success := testData.ReadPartialData(nil)
	assert.True(t, success)
	assert.Equal(t, testData, result)

	// Test 2: Filter by electrical connection ID
	filter := &FilterType{
		ElectricalConnectionDescriptionListDataSelectors: &ElectricalConnectionDescriptionListDataSelectorsType{
			ElectricalConnectionId: &electricalConnectionId1,
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData := result.(*ElectricalConnectionDescriptionListDataType)
	assert.Len(t, resultData.ElectricalConnectionDescriptionData, 1)
	assert.Equal(t, electricalConnectionId1, *resultData.ElectricalConnectionDescriptionData[0].ElectricalConnectionId)
	assert.Equal(t, "ac", string(*resultData.ElectricalConnectionDescriptionData[0].PowerSupplyType))

	// Test 3: Filter by power supply type
	filter = &FilterType{
		ElectricalConnectionDescriptionListDataSelectors: &ElectricalConnectionDescriptionListDataSelectorsType{
			ScopeType: &scopeType2,
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*ElectricalConnectionDescriptionListDataType)
	assert.Len(t, resultData.ElectricalConnectionDescriptionData, 1)
	assert.Equal(t, "system", string(*resultData.ElectricalConnectionDescriptionData[0].ScopeType))
	assert.Equal(t, "DC Connection 1", string(*resultData.ElectricalConnectionDescriptionData[0].Label))

	// Test 4: Elements filtering test
	filter = &FilterType{
		ElectricalConnectionDescriptionDataElements: &ElectricalConnectionDescriptionDataElementsType{
			ElectricalConnectionId: &ElementTagType{},
			Label:                  &ElementTagType{},
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*ElectricalConnectionDescriptionListDataType)
	assert.Len(t, resultData.ElectricalConnectionDescriptionData, 2)
	// Check that only requested elements are included
	for _, item := range resultData.ElectricalConnectionDescriptionData {
		assert.NotNil(t, item.ElectricalConnectionId)
		assert.NotNil(t, item.Label)
		// PowerSupplyType and Description should be nil since not requested in elements
		assert.Nil(t, item.PowerSupplyType)
		assert.Nil(t, item.Description)
	}

	// Test 5: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestElectricalConnectionCharacteristicListDataType_ReadPartialData_Enhanced(t *testing.T) {
	electricalConnectionId1 := ElectricalConnectionIdType(1)
	electricalConnectionId2 := ElectricalConnectionIdType(2)
	characteristicId1 := ElectricalConnectionCharacteristicIdType(10)
	characteristicId2 := ElectricalConnectionCharacteristicIdType(20)

	testData := &ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []ElectricalConnectionCharacteristicDataType{
			{
				ElectricalConnectionId: &electricalConnectionId1,
				CharacteristicId:       &characteristicId1,
				CharacteristicContext:  util.Ptr(ElectricalConnectionCharacteristicContextType("entity")),
				CharacteristicType:     util.Ptr(ElectricalConnectionCharacteristicTypeType("contractualConsumptionNominalMax")),
				Value:                  &ScaledNumberType{Number: util.Ptr(NumberType(3000)), Scale: util.Ptr(ScaleType(0))},
			},
			{
				ElectricalConnectionId: &electricalConnectionId1,
				CharacteristicId:       &characteristicId2,
				CharacteristicContext:  util.Ptr(ElectricalConnectionCharacteristicContextType("entity")),
				CharacteristicType:     util.Ptr(ElectricalConnectionCharacteristicTypeType("contractualProductionNominalMax")),
				Value:                  &ScaledNumberType{Number: util.Ptr(NumberType(2000)), Scale: util.Ptr(ScaleType(0))},
			},
			{
				ElectricalConnectionId: &electricalConnectionId2,
				CharacteristicId:       &characteristicId1,
				CharacteristicContext:  util.Ptr(ElectricalConnectionCharacteristicContextType("entity")),
				CharacteristicType:     util.Ptr(ElectricalConnectionCharacteristicTypeType("contractualConsumptionNominalMax")),
				Value:                  &ScaledNumberType{Number: util.Ptr(NumberType(5000)), Scale: util.Ptr(ScaleType(0))},
			},
		},
	}

	// Test 1: Filter by electrical connection ID and characteristic ID
	filter := &FilterType{
		ElectricalConnectionCharacteristicListDataSelectors: &ElectricalConnectionCharacteristicListDataSelectorsType{
			ElectricalConnectionId: &electricalConnectionId1,
			CharacteristicId:       &characteristicId1,
		},
	}

	result, success := testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData := result.(*ElectricalConnectionCharacteristicListDataType)
	assert.Len(t, resultData.ElectricalConnectionCharacteristicData, 1)
	assert.Equal(t, electricalConnectionId1, *resultData.ElectricalConnectionCharacteristicData[0].ElectricalConnectionId)
	assert.Equal(t, characteristicId1, *resultData.ElectricalConnectionCharacteristicData[0].CharacteristicId)

	// Test 2: Filter by characteristic context
	filter = &FilterType{
		ElectricalConnectionCharacteristicListDataSelectors: &ElectricalConnectionCharacteristicListDataSelectorsType{
			CharacteristicContext: util.Ptr(ElectricalConnectionCharacteristicContextType("entity")),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*ElectricalConnectionCharacteristicListDataType)
	assert.Len(t, resultData.ElectricalConnectionCharacteristicData, 3) // All items have entity context

	// Test 3: Filter by characteristic type
	filter = &FilterType{
		ElectricalConnectionCharacteristicListDataSelectors: &ElectricalConnectionCharacteristicListDataSelectorsType{
			CharacteristicType: util.Ptr(ElectricalConnectionCharacteristicTypeType("contractualProductionNominalMax")),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*ElectricalConnectionCharacteristicListDataType)
	assert.Len(t, resultData.ElectricalConnectionCharacteristicData, 1) // Only one with this type
	assert.Equal(t, "contractualProductionNominalMax", string(*resultData.ElectricalConnectionCharacteristicData[0].CharacteristicType))

	// Test 4: Complex filtering with multiple criteria
	filter = &FilterType{
		ElectricalConnectionCharacteristicListDataSelectors: &ElectricalConnectionCharacteristicListDataSelectorsType{
			ElectricalConnectionId: &electricalConnectionId1,
			CharacteristicType:     util.Ptr(ElectricalConnectionCharacteristicTypeType("contractualConsumptionNominalMax")),
		},
	}

	result, success = testData.ReadPartialData(filter)
	assert.True(t, success)
	resultData = result.(*ElectricalConnectionCharacteristicListDataType)
	assert.Len(t, resultData.ElectricalConnectionCharacteristicData, 1)
	assert.Equal(t, electricalConnectionId1, *resultData.ElectricalConnectionCharacteristicData[0].ElectricalConnectionId)
	assert.Equal(t, "contractualConsumptionNominalMax", string(*resultData.ElectricalConnectionCharacteristicData[0].CharacteristicType))

	// Test 5: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestElectricalConnectionPermittedValueSetListDataType_ReadPartialData_UpdateList(t *testing.T) {
	electricalConnectionId1 := ElectricalConnectionIdType(1)
	electricalConnectionId2 := ElectricalConnectionIdType(2)
	parameterId1 := ElectricalConnectionParameterIdType(10)
	parameterId2 := ElectricalConnectionParameterIdType(20)

	testData := &ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: &electricalConnectionId1,
				ParameterId:            &parameterId1,
				PermittedValueSet: []ScaledNumberSetType{
					{
						Value: []ScaledNumberType{
							{Number: util.Ptr(NumberType(100)), Scale: util.Ptr(ScaleType(0))},
							{Number: util.Ptr(NumberType(200)), Scale: util.Ptr(ScaleType(0))},
						},
					},
				},
			},
		},
	}

	// Test UpdateList functionality
	newItem := ElectricalConnectionPermittedValueSetDataType{
		ElectricalConnectionId: &electricalConnectionId2,
		ParameterId:            &parameterId2,
		PermittedValueSet: []ScaledNumberSetType{
			{
				Value: []ScaledNumberType{
					{Number: util.Ptr(NumberType(300)), Scale: util.Ptr(ScaleType(0))},
				},
			},
		},
	}

	newData := &ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []ElectricalConnectionPermittedValueSetDataType{newItem},
	}

	updatedList, success := testData.UpdateList(false, true, newData, nil, nil, Ptr(FunctionTypeElectricalConnectionPermittedValueSetListData))
	assert.True(t, success)
	updatedListData := updatedList.([]ElectricalConnectionPermittedValueSetDataType)
	assert.Len(t, updatedListData, 2) // Should have both items now

	// Verify new item was added
	var foundNewItem *ElectricalConnectionPermittedValueSetDataType
	for _, item := range updatedListData {
		if *item.ElectricalConnectionId == electricalConnectionId2 && *item.ParameterId == parameterId2 {
			foundNewItem = &item
			break
		}
	}

	assert.NotNil(t, foundNewItem)
	assert.Equal(t, electricalConnectionId2, *foundNewItem.ElectricalConnectionId)
	assert.Equal(t, parameterId2, *foundNewItem.ParameterId)

	// Test 2: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestElectricalConnectionCharacteristicListDataType_ReadPartialData_UpdateList(t *testing.T) {
	electricalConnectionId1 := ElectricalConnectionIdType(1)
	characteristicId1 := ElectricalConnectionCharacteristicIdType(10)
	parameterId1 := ElectricalConnectionParameterIdType(5)

	testData := &ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []ElectricalConnectionCharacteristicDataType{
			{
				ElectricalConnectionId: &electricalConnectionId1,
				ParameterId:            &parameterId1,
				CharacteristicId:       &characteristicId1,
				CharacteristicContext:  util.Ptr(ElectricalConnectionCharacteristicContextType("entity")),
				CharacteristicType:     util.Ptr(ElectricalConnectionCharacteristicTypeType("contractualConsumptionNominalMax")),
				Value:                  &ScaledNumberType{Number: util.Ptr(NumberType(3000)), Scale: util.Ptr(ScaleType(0))},
			},
		},
	}

	// Test UpdateList - modify existing item
	updatedItem := ElectricalConnectionCharacteristicDataType{
		ElectricalConnectionId: &electricalConnectionId1,
		ParameterId:            &parameterId1,
		CharacteristicId:       &characteristicId1,
		CharacteristicContext:  util.Ptr(ElectricalConnectionCharacteristicContextType("entity")),
		CharacteristicType:     util.Ptr(ElectricalConnectionCharacteristicTypeType("contractualConsumptionNominalMax")),
		Value:                  &ScaledNumberType{Number: util.Ptr(NumberType(5000)), Scale: util.Ptr(ScaleType(0))},
	}

	newData := &ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []ElectricalConnectionCharacteristicDataType{updatedItem},
	}

	updatedList, success := testData.UpdateList(false, true, newData, nil, nil, Ptr(FunctionTypeElectricalConnectionCharacteristicListData))
	assert.True(t, success)
	updatedListData := updatedList.([]ElectricalConnectionCharacteristicDataType)
	assert.Len(t, updatedListData, 1)

	// Verify item was updated
	assert.Equal(t, float64(5000), updatedListData[0].Value.GetValue())

	// Test 2: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = testData.ReadPartialData(emptyFilter)
	assert.False(t, success)
}
