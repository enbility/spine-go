package model

import (
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

type TestUpdateData struct {
	Id           *uint `eebus:"key"`
	IsChangeable *bool `eebus:"writecheck"`
	DataItem     *int
}

type TestUpdater struct {
	// updateSelectorHashKey *string
	// deleteSelectorHashKey *string
}

func TestUpdateList_NewItem(t *testing.T) {
	existingData := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(1))}}
	newData := []TestUpdateData{{Id: util.Ptr(uint(2)), DataItem: util.Ptr(int(2))}}

	expectedResult := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(1))}, {Id: util.Ptr(uint(2)), DataItem: util.Ptr(int(2))}}

	// Act
	result, boolV := UpdateList(false, existingData, newData, nil, nil, nil)

	assert.True(t, boolV)
	assert.Equal(t, expectedResult, result)

	expectedResult = []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(1))}}

	// Act
	result, boolV = UpdateList(true, existingData, newData, nil, nil, nil)

	assert.False(t, boolV)
	assert.Equal(t, expectedResult, result)
}

func TestUpdateList_ChangedItem(t *testing.T) {
	existingData := []TestUpdateData{{Id: util.Ptr(uint(1)), IsChangeable: util.Ptr(false), DataItem: util.Ptr(int(1))}}
	newData := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(2))}}

	expectedResult := []TestUpdateData{{Id: util.Ptr(uint(1)), IsChangeable: util.Ptr(false), DataItem: util.Ptr(int(2))}}

	// Act
	result, boolV := UpdateList(false, existingData, newData, nil, nil, nil)

	assert.True(t, boolV)
	assert.Equal(t, expectedResult, result)

	expectedResult = []TestUpdateData{{Id: util.Ptr(uint(1)), IsChangeable: util.Ptr(false), DataItem: util.Ptr(int(1))}}

	// Act
	result, boolV = UpdateList(true, existingData, newData, nil, nil, nil)

	assert.False(t, boolV)
	assert.Equal(t, expectedResult, result)
}

func TestUpdateList_NewAndChangedItem(t *testing.T) {
	existingData := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(1))}}
	newData := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(2))}, {Id: util.Ptr(uint(3)), DataItem: util.Ptr(int(3))}}

	expectedResult := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(2))}, {Id: util.Ptr(uint(3)), DataItem: util.Ptr(int(3))}}

	// Act
	result, boolV := UpdateList(false, existingData, newData, nil, nil, nil)

	assert.True(t, boolV)
	assert.Equal(t, expectedResult, result)

	expectedResult = []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(1))}}

	// Act
	result, boolV = UpdateList(true, existingData, newData, nil, nil, nil)

	assert.False(t, boolV)
	assert.Equal(t, expectedResult, result)
}

func TestUpdateList_ItemWithNoIdentifier(t *testing.T) {
	existingData := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(1))}, {Id: util.Ptr(uint(2)), DataItem: util.Ptr(int(2))}}
	newData := []TestUpdateData{{DataItem: util.Ptr(int(3))}}

	expectedResult := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(3))}, {Id: util.Ptr(uint(2)), DataItem: util.Ptr(int(3))}}

	// Act
	result, boolV := UpdateList(false, existingData, newData, nil, nil, nil)

	assert.True(t, boolV)
	assert.Equal(t, expectedResult, result)

	expectedResult = []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(3))}, {Id: util.Ptr(uint(2)), DataItem: util.Ptr(int(3))}}

	// Act
	result, boolV = UpdateList(true, existingData, newData, nil, nil, nil)

	assert.False(t, boolV)
	assert.Equal(t, expectedResult, result)
}

func TestUpdateList_FilterDelete(t *testing.T) {
	existingData := []LoadControlLimitDataType{
		{
			LimitId: util.Ptr(LoadControlLimitIdType(0)),
			Value:   NewScaledNumberType(0),
		},
		{
			LimitId: util.Ptr(LoadControlLimitIdType(1)),
			Value:   NewScaledNumberType(0),
		},
	}
	newData := []LoadControlLimitDataType{
		{
			LimitId: util.Ptr(LoadControlLimitIdType(1)),
			Value:   NewScaledNumberType(10),
		},
	}

	expectedResult := []LoadControlLimitDataType{
		{
			LimitId: util.Ptr(LoadControlLimitIdType(1)),
			Value:   NewScaledNumberType(10),
		},
	}

	filterDelete := &FilterType{CmdControl: &CmdControlType{Delete: &ElementTagType{}}}
	filterDelete.CmdControl.Delete = new(ElementTagType)
	filterDelete.LoadControlLimitListDataSelectors = &LoadControlLimitListDataSelectorsType{
		LimitId: util.Ptr(LoadControlLimitIdType(0)),
	}

	filterPartial := NewFilterTypePartial()
	filterPartial.LoadControlLimitListDataSelectors = &LoadControlLimitListDataSelectorsType{
		LimitId: util.Ptr(LoadControlLimitIdType(1)),
	}

	// Act
	result, boolV := UpdateList(false, existingData, newData, filterPartial, filterDelete, nil)

	assert.True(t, boolV)
	assert.Equal(t, expectedResult, result)

	newData = []LoadControlLimitDataType{
		{
			Value: NewScaledNumberType(10),
		},
	}

	expectedResult = []LoadControlLimitDataType{
		{
			LimitId: util.Ptr(LoadControlLimitIdType(0)),
			Value:   NewScaledNumberType(0),
		},
		{
			LimitId: util.Ptr(LoadControlLimitIdType(1)),
			Value:   NewScaledNumberType(0),
		},
	}

	// Act
	result, boolV = UpdateList(true, existingData, newData, filterPartial, filterDelete, nil)

	assert.False(t, boolV)
	assert.Equal(t, expectedResult, result)
}

func TestRemoveFieldFromType(t *testing.T) {
	items := &LoadControlLimitListDataType{
		LoadControlLimitData: []LoadControlLimitDataType{
			{
				LimitId: util.Ptr(LoadControlLimitIdType(1)),
				Value:   NewScaledNumberType(16.0),
			},
		},
	}

	elements := &LoadControlLimitDataElementsType{
		Value: &ScaledNumberElementsType{},
	}

	RemoveElementFromItem(&items.LoadControlLimitData[0], elements)

	var nilValue *ScaledNumberType

	assert.Equal(t, nilValue, items.LoadControlLimitData[0].Value)
}

// TestUpdateList_FilterDeleteDataError tests UpdateList behavior when filterDelete.Data() fails.
// This test verifies that UpdateList correctly ignores invalid filters for backward compatibility.
// Invalid filters (those without proper function fields) are skipped rather than causing failure.
func TestUpdateList_FilterDeleteDataError(t *testing.T) {
	existingData := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(1))}}
	newData := []TestUpdateData{{Id: util.Ptr(uint(2)), DataItem: util.Ptr(int(2))}}

	// Create a FilterType without proper EEBus tags that will cause Data() to return an error
	// This filter has no fields with the required "fct" (function) and "typ" (type) tags
	invalidFilterDelete := &FilterType{
		CmdControl: &CmdControlType{Delete: &ElementTagType{}},
		FilterId:   util.Ptr(FilterIdType(1)),
		// No selector fields with proper eebus tags - this will cause Data() to fail
	}

	// Act - UpdateList ignores invalid filters (backward compatibility)
	result, success := UpdateList(false, existingData, newData, nil, invalidFilterDelete, nil)

	// Assert - operation should succeed, ignoring the invalid filter
	assert.True(t, success)
	// Result should contain merged data since invalid filter is ignored
	expectedResult := []TestUpdateData{
		{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(1))},
		{Id: util.Ptr(uint(2)), DataItem: util.Ptr(int(2))},
	}
	assert.Equal(t, expectedResult, result)
}

// TestUpdateList_FilterPartialDataError tests UpdateList behavior when filterPartial.Data() fails.
// This test verifies that UpdateList correctly ignores invalid filters for backward compatibility.
// When partial filter is invalid, normal merge processing occurs instead.
func TestUpdateList_FilterPartialDataError(t *testing.T) {
	existingData := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(1))}}
	newData := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(2))}}

	// Create a FilterType without proper EEBus tags that will cause Data() to return an error
	// This filter has CmdControl.Partial set but no selector fields with proper eebus tags
	invalidFilterPartial := &FilterType{
		CmdControl: &CmdControlType{Partial: &ElementTagType{}},
		FilterId:   util.Ptr(FilterIdType(2)),
		// No selector fields with proper eebus tags - this will cause Data() to fail
	}

	// Act - UpdateList ignores invalid filters (backward compatibility)
	result, success := UpdateList(false, existingData, newData, invalidFilterPartial, nil, nil)

	// Assert - operation should succeed, ignoring the invalid filter
	assert.True(t, success)
	// Result should contain updated data since invalid partial filter is ignored
	expectedResult := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(2))}}
	assert.Equal(t, expectedResult, result)
}

// TestUpdateList_BothFiltersDataError tests UpdateList behavior when both filters fail.
// This test verifies that UpdateList correctly ignores all invalid filters for backward compatibility.
// When both filters are invalid, normal merge processing occurs.
func TestUpdateList_BothFiltersDataError(t *testing.T) {
	existingData := []TestUpdateData{{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(1))}}
	newData := []TestUpdateData{{Id: util.Ptr(uint(2)), DataItem: util.Ptr(int(2))}}

	// Create invalid filters that will cause Data() to return errors
	invalidFilterDelete := &FilterType{
		CmdControl: &CmdControlType{Delete: &ElementTagType{}},
		FilterId:   util.Ptr(FilterIdType(3)),
		// Missing required eebus tags
	}

	invalidFilterPartial := &FilterType{
		CmdControl: &CmdControlType{Partial: &ElementTagType{}},
		FilterId:   util.Ptr(FilterIdType(4)),
		// Missing required eebus tags
	}

	// Act - UpdateList ignores invalid filters (backward compatibility)
	result, success := UpdateList(false, existingData, newData, invalidFilterPartial, invalidFilterDelete, nil)

	// Assert - operation should succeed, ignoring both invalid filters
	assert.True(t, success)
	// Result should contain merged data since both invalid filters are ignored
	expectedResult := []TestUpdateData{
		{Id: util.Ptr(uint(1)), DataItem: util.Ptr(int(1))},
		{Id: util.Ptr(uint(2)), DataItem: util.Ptr(int(2))},
	}
	assert.Equal(t, expectedResult, result)
}

// TestUpdateList_ValidDeleteInvalidPartialFilters tests mixed scenario with one valid and one invalid filter.
// This test verifies that UpdateList processes valid filters and ignores invalid ones.
// The valid delete filter is processed while the invalid partial filter is skipped.
func TestUpdateList_ValidDeleteInvalidPartialFilters(t *testing.T) {
	existingData := []LoadControlLimitDataType{
		{
			LimitId: util.Ptr(LoadControlLimitIdType(1)),
			Value:   NewScaledNumberType(10),
		},
	}
	newData := []LoadControlLimitDataType{
		{
			LimitId: util.Ptr(LoadControlLimitIdType(1)),
			Value:   NewScaledNumberType(20),
		},
	}

	// Valid delete filter with proper eebus tags
	validFilterDelete := &FilterType{CmdControl: &CmdControlType{Delete: &ElementTagType{}}}
	validFilterDelete.LoadControlLimitListDataSelectors = &LoadControlLimitListDataSelectorsType{
		LimitId: util.Ptr(LoadControlLimitIdType(0)),
	}

	// Invalid partial filter without proper eebus tags
	invalidFilterPartial := &FilterType{
		CmdControl: &CmdControlType{Partial: &ElementTagType{}},
		FilterId:   util.Ptr(FilterIdType(5)),
		// Missing required selector with eebus tags
	}

	// Act - operation should succeed, ignoring invalid partial filter
	result, success := UpdateList(false, existingData, newData, invalidFilterPartial, validFilterDelete, nil)

	// Assert - operation should succeed, processing valid delete and ignoring invalid partial
	assert.True(t, success)
	// Result should contain updated data (delete ignored limitId(0), partial ignored due to invalid)
	expectedResult := []LoadControlLimitDataType{
		{
			LimitId: util.Ptr(LoadControlLimitIdType(1)),
			Value:   NewScaledNumberType(20),
		},
	}
	assert.Equal(t, expectedResult, result)
}

// TestUpdateList_InvalidDeleteValidPartialFilters tests mixed scenario with invalid delete and valid partial filter.
// This test verifies that UpdateList processes valid filters and ignores invalid ones.
// The valid partial filter is processed while the invalid delete filter is skipped.
func TestUpdateList_InvalidDeleteValidPartialFilters(t *testing.T) {
	existingData := []LoadControlLimitDataType{
		{
			LimitId: util.Ptr(LoadControlLimitIdType(1)),
			Value:   NewScaledNumberType(10),
		},
	}
	newData := []LoadControlLimitDataType{
		{
			Value: NewScaledNumberType(20),
		},
	}

	// Invalid delete filter without proper eebus tags
	invalidFilterDelete := &FilterType{
		CmdControl: &CmdControlType{Delete: &ElementTagType{}},
		FilterId:   util.Ptr(FilterIdType(6)),
		// Missing required selector with eebus tags
	}

	// Valid partial filter with proper eebus tags
	validFilterPartial := NewFilterTypePartial()
	validFilterPartial.LoadControlLimitListDataSelectors = &LoadControlLimitListDataSelectorsType{
		LimitId: util.Ptr(LoadControlLimitIdType(1)),
	}

	// Act - operation should succeed, ignoring invalid delete filter
	result, success := UpdateList(false, existingData, newData, validFilterPartial, invalidFilterDelete, nil)

	// Assert - operation should succeed, processing valid partial and ignoring invalid delete
	assert.True(t, success)
	// Result should contain data updated by partial filter (delete filter ignored)
	expectedResult := []LoadControlLimitDataType{
		{
			LimitId: util.Ptr(LoadControlLimitIdType(1)),
			Value:   NewScaledNumberType(20),
		},
	}
	assert.Equal(t, expectedResult, result)
}

// TODO: Fix, as these tests won't work right now as TestUpdater doesn't use FilterProvider and its data structure
/*
func TestUpdateList_UpdateSelector(t *testing.T) {
	existingData := []TestUpdateData{{Id: util.Ptr(1), DataItem: 1}, {Id: util.Ptr(2), DataItem: 2}}
	newData := []TestUpdateData{{DataItem: 3}}

	dataProvider := &TestUpdater{
		updateSelectorHashKey: util.Ptr("1"),
	}
	expectedResult := []TestUpdateData{{Id: util.Ptr(1), DataItem: 3}, {Id: util.Ptr(2), DataItem: 2}}

	// Act
	result := UpdateList[TestUpdateData](existingData, newData, dataProvider)

	assert.Equal(t, expectedResult, result)
}

func TestUpdateList_DeleteSelector(t *testing.T) {
	existingData := []TestUpdateData{{Id: util.Ptr(1), DataItem: 1}, {Id: util.Ptr(2), DataItem: 2}}
	newData := []TestUpdateData{{Id: util.Ptr(0), DataItem: 0}}

	dataProvider := &TestUpdater{
		deleteSelectorHashKey: util.Ptr("1"),
	}
	expectedResult := []TestUpdateData{{Id: util.Ptr(2), DataItem: 2}}

	// Act
	result := UpdateList[TestUpdateData](existingData, newData, dataProvider)

	assert.Equal(t, expectedResult, result)
}
*/
