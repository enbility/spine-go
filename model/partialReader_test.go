package model

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create pointers (commonly used in Go for optional fields)
func Ptr[T any](v T) *T {
	return &v
}

// TestGenericHelperFunctions tests the core generic helper functions
func TestMatchesSelector(t *testing.T) {
	type TestItem struct {
		ID          *int
		Name        *string
		DeviceAddr  *DeviceAddressType
		EntityAddr  *EntityAddressType
		FeatureAddr *FeatureAddressType
		Tags        []string
	}

	type TestSelector struct {
		ID          *int
		Name        *string
		DeviceAddr  *DeviceAddressType
		EntityAddr  *EntityAddressType
		FeatureAddr *FeatureAddressType
		Tags        *string
	}

	device1 := AddressDeviceType("device1")
	entity1 := AddressEntityType(1)
	feature1 := AddressFeatureType(100)

	item := TestItem{
		ID:   Ptr(1),
		Name: Ptr("test"),
		DeviceAddr: &DeviceAddressType{
			Device: &device1,
		},
		EntityAddr: &EntityAddressType{
			Device: &device1,
			Entity: []AddressEntityType{entity1},
		},
		FeatureAddr: &FeatureAddressType{
			Device:  &device1,
			Entity:  []AddressEntityType{entity1},
			Feature: &feature1,
		},
		Tags: []string{"tag1", "tag2"},
	}

	// Test 1: Nil selector should match all items
	assert.True(t, matchesSelector(item, nil))

	// Test 2: Empty selector (all fields nil) should match all items
	emptySelector := &TestSelector{}
	assert.True(t, matchesSelector(item, emptySelector))

	// Test 3: Matching selector should return true
	matchingSelector := &TestSelector{
		ID:   Ptr(1),
		Name: Ptr("test"),
	}
	assert.True(t, matchesSelector(item, matchingSelector))

	// Test 4: Non-matching selector should return false
	nonMatchingSelector := &TestSelector{
		ID: Ptr(2),
	}
	assert.False(t, matchesSelector(item, nonMatchingSelector))

	// Test 5: Address matching with device address
	deviceAddrSelector := &TestSelector{
		DeviceAddr: &DeviceAddressType{Device: &device1},
	}
	assert.True(t, matchesSelector(item, deviceAddrSelector))

	// Test 6: Non-matching device address
	device2 := AddressDeviceType("device2")
	nonMatchingDeviceSelector := &TestSelector{
		DeviceAddr: &DeviceAddressType{Device: &device2},
	}
	assert.False(t, matchesSelector(item, nonMatchingDeviceSelector))

	// Test 7: Array field matching
	arrayMatchSelector := &TestSelector{
		Tags: Ptr("tag1"), // Should find tag1 in the array
	}
	assert.True(t, matchesSelector(item, arrayMatchSelector))

	// Test 8: Non-matching array field
	arrayNonMatchSelector := &TestSelector{
		Tags: Ptr("tag3"), // tag3 is not in the array
	}
	assert.False(t, matchesSelector(item, arrayNonMatchSelector))
}

func TestFilterListDataBySelectors(t *testing.T) {
	// For this test, we'll use the existing ElectricalConnectionStateDataType
	data := []ElectricalConnectionStateDataType{
		{ElectricalConnectionId: Ptr(ElectricalConnectionIdType(1)), CurrentEnergyMode: Ptr(EnergyModeTypeConsume)},
		{ElectricalConnectionId: Ptr(ElectricalConnectionIdType(2)), CurrentEnergyMode: Ptr(EnergyModeTypeProduce)},
		{ElectricalConnectionId: Ptr(ElectricalConnectionIdType(3)), CurrentEnergyMode: Ptr(EnergyModeTypeConsume)},
	}

	// Test 1: Nil filter should return all data
	result, success := filterListDataBySelectors[ElectricalConnectionStateDataType, ElectricalConnectionStateListDataSelectorsType](data, nil)
	assert.True(t, success)
	assert.Len(t, result, 3)

	// Test 2: Filter with matching selector
	filter := &FilterType{
		ElectricalConnectionStateListDataSelectors: &ElectricalConnectionStateListDataSelectorsType{
			ElectricalConnectionId: Ptr(ElectricalConnectionIdType(2)),
		},
	}

	result, success = filterListDataBySelectors[ElectricalConnectionStateDataType, ElectricalConnectionStateListDataSelectorsType](data, filter)
	assert.True(t, success)
	assert.Len(t, result, 1)
	assert.Equal(t, ElectricalConnectionIdType(2), *result[0].ElectricalConnectionId)
}

func TestPartialListDataRead(t *testing.T) {
	// Use existing types from the model
	data := []MeasurementDataType{
		{
			MeasurementId: Ptr(MeasurementIdType(1)),
			Value:         &ScaledNumberType{Number: Ptr(NumberType(100)), Scale: Ptr(ScaleType(0))},
			Timestamp:     Ptr(AbsoluteOrRelativeTimeType("2023-01-01T00:00:00Z")),
		},
		{
			MeasurementId: Ptr(MeasurementIdType(2)),
			Value:         &ScaledNumberType{Number: Ptr(NumberType(200)), Scale: Ptr(ScaleType(0))},
			Timestamp:     Ptr(AbsoluteOrRelativeTimeType("2023-01-01T01:00:00Z")),
		},
	}

	// Test 1: Filter with selector and elements
	filter := &FilterType{
		MeasurementListDataSelectors: &MeasurementListDataSelectorsType{
			MeasurementId: Ptr(MeasurementIdType(1)),
		},
		MeasurementDataElements: &MeasurementDataElementsType{
			MeasurementId: Ptr(ElementTagType{}), // Only include MeasurementId
			Value:         Ptr(ElementTagType{}), // Only include Value
			// Timestamp is excluded
		},
	}

	result, success := partialListDataRead[MeasurementDataType, MeasurementListDataSelectorsType, MeasurementDataElementsType](data, filter)
	assert.True(t, success)
	assert.Len(t, result, 1)
	assert.Equal(t, MeasurementIdType(1), *result[0].MeasurementId)
	assert.Equal(t, NumberType(100), *result[0].Value.Number)
	// Timestamp should be filtered out (nil)
	assert.Nil(t, result[0].Timestamp)
}

func TestFilterElementsFromItems(t *testing.T) {
	// Use real measurement data types
	items := []MeasurementDataType{
		{
			MeasurementId: Ptr(MeasurementIdType(1)),
			Value:         &ScaledNumberType{Number: Ptr(NumberType(100))},
			Timestamp:     Ptr(AbsoluteOrRelativeTimeType("2023-01-01T00:00:00Z")),
		},
		{
			MeasurementId: Ptr(MeasurementIdType(2)),
			Value:         &ScaledNumberType{Number: Ptr(NumberType(200))},
			Timestamp:     Ptr(AbsoluteOrRelativeTimeType("2023-01-01T01:00:00Z")),
		},
	}

	// Test 1: Include only MeasurementId and Value
	elements := &MeasurementDataElementsType{
		MeasurementId: Ptr(ElementTagType{}),
		Value:         Ptr(ElementTagType{}),
	}
	result := filterElementsFromItems(items, elements)

	assert.Len(t, result, 2)
	assert.Equal(t, MeasurementIdType(1), *result[0].MeasurementId)
	assert.NotNil(t, result[0].Value)
	assert.Nil(t, result[0].Timestamp) // Should be filtered out

	// Test 2: Empty elements should return original items
	emptyElements := &MeasurementDataElementsType{}
	result = filterElementsFromItems(items, emptyElements)
	assert.Len(t, result, 2)
	assert.Equal(t, items, result) // Should be unchanged
}

func TestAddressMatchingFunctions(t *testing.T) {
	device1 := AddressDeviceType("device1")
	device2 := AddressDeviceType("device2")
	entity1 := AddressEntityType(1)
	entity2 := AddressEntityType(2)
	feature1 := AddressFeatureType(100)
	feature2 := AddressFeatureType(200)

	// Test deviceAddressesMatch
	t.Run("DeviceAddressesMatch", func(t *testing.T) {
		addr1 := &DeviceAddressType{Device: &device1}
		addr2 := &DeviceAddressType{Device: &device1}
		addr3 := &DeviceAddressType{Device: &device2}

		assert.True(t, deviceAddressesMatch(addr1, addr2))
		assert.False(t, deviceAddressesMatch(addr1, addr3))
		assert.True(t, deviceAddressesMatch(nil, nil))
		assert.False(t, deviceAddressesMatch(addr1, nil))
		assert.False(t, deviceAddressesMatch(nil, addr1))
	})

	// Test entityAddressesMatch
	t.Run("EntityAddressesMatch", func(t *testing.T) {
		addr1 := &EntityAddressType{Device: &device1, Entity: []AddressEntityType{entity1}}
		addr2 := &EntityAddressType{Device: &device1, Entity: []AddressEntityType{entity1}}
		addr3 := &EntityAddressType{Device: &device1, Entity: []AddressEntityType{entity2}}
		addr4 := &EntityAddressType{Device: &device2, Entity: []AddressEntityType{entity1}}

		assert.True(t, entityAddressesMatch(addr1, addr2))
		assert.False(t, entityAddressesMatch(addr1, addr3)) // Different entity
		assert.False(t, entityAddressesMatch(addr1, addr4)) // Different device
		assert.True(t, entityAddressesMatch(nil, nil))
		assert.False(t, entityAddressesMatch(addr1, nil))
	})

	// Test featureAddressesMatch
	t.Run("FeatureAddressesMatch", func(t *testing.T) {
		addr1 := &FeatureAddressType{Device: &device1, Entity: []AddressEntityType{entity1}, Feature: &feature1}
		addr2 := &FeatureAddressType{Device: &device1, Entity: []AddressEntityType{entity1}, Feature: &feature1}
		addr3 := &FeatureAddressType{Device: &device1, Entity: []AddressEntityType{entity1}, Feature: &feature2}

		assert.True(t, featureAddressesMatch(addr1, addr2))
		assert.True(t, addressesMatch(addr1, addr2)) // Test the alias function too
		assert.False(t, featureAddressesMatch(addr1, addr3))
		assert.False(t, addressesMatch(addr1, addr3))
	})
}

func TestHelperFunctions(t *testing.T) {
	// Test isFieldNilOrEmpty
	t.Run("IsFieldNilOrEmpty", func(t *testing.T) {
		var nilPtr *int
		validPtr := Ptr(42)
		emptySlice := []string{}
		nonEmptySlice := []string{"item"}
		emptyString := ""
		nonEmptyString := "test"

		assert.True(t, isFieldNilOrEmpty(reflect.ValueOf(nilPtr)))
		assert.False(t, isFieldNilOrEmpty(reflect.ValueOf(validPtr)))
		assert.True(t, isFieldNilOrEmpty(reflect.ValueOf(emptySlice)))
		assert.False(t, isFieldNilOrEmpty(reflect.ValueOf(nonEmptySlice)))
		assert.True(t, isFieldNilOrEmpty(reflect.ValueOf(emptyString)))
		assert.False(t, isFieldNilOrEmpty(reflect.ValueOf(nonEmptyString)))
	})

	// Test getNonNilElementNames
	t.Run("GetNonNilElementNames", func(t *testing.T) {
		type TestElements struct {
			Field1 *bool
			Field2 *bool
			Field3 *bool
		}

		elements := &TestElements{
			Field1: Ptr(true),
			Field3: Ptr(false), // false is still a valid non-nil value
		}

		names := getNonNilElementNames(elements)
		assert.Contains(t, names, "Field1")
		assert.Contains(t, names, "Field3")
		assert.NotContains(t, names, "Field2")
		assert.Len(t, names, 2)
	})

	// Test isElementFieldNil
	t.Run("IsElementFieldNil", func(t *testing.T) {
		var nilPtr *int
		validPtr := Ptr(42)

		assert.True(t, isElementFieldNil(nil))
		assert.True(t, isElementFieldNil(nilPtr))
		assert.False(t, isElementFieldNil(validPtr))
		assert.False(t, isElementFieldNil(42))
		assert.False(t, isElementFieldNil("test"))
	})
}

func TestAddressMatching(t *testing.T) {
	device1 := AddressDeviceType("device1")
	device2 := AddressDeviceType("device2")
	entity1 := AddressEntityType(1)
	entity2 := AddressEntityType(2)

	// Test device address matching
	addr1 := &DeviceAddressType{Device: &device1}
	addr2 := &DeviceAddressType{Device: &device1}
	addr3 := &DeviceAddressType{Device: &device2}

	assert.True(t, deviceAddressesMatch(addr1, addr2))
	assert.False(t, deviceAddressesMatch(addr1, addr3))
	assert.True(t, deviceAddressesMatch(nil, nil))
	assert.False(t, deviceAddressesMatch(addr1, nil))

	// Test entity address matching
	entityAddr1 := &EntityAddressType{Device: &device1, Entity: []AddressEntityType{entity1}}
	entityAddr2 := &EntityAddressType{Device: &device1, Entity: []AddressEntityType{entity1}}
	entityAddr3 := &EntityAddressType{Device: &device1, Entity: []AddressEntityType{entity2}}

	assert.True(t, entityAddressesMatch(entityAddr1, entityAddr2))
	assert.False(t, entityAddressesMatch(entityAddr1, entityAddr3))

	// Test feature address matching
	feature1 := AddressFeatureType(100)
	featureAddr1 := &FeatureAddressType{Device: &device1, Entity: []AddressEntityType{entity1}, Feature: &feature1}
	featureAddr2 := &FeatureAddressType{Device: &device1, Entity: []AddressEntityType{entity1}, Feature: &feature1}

	assert.True(t, featureAddressesMatch(featureAddr1, featureAddr2))
	assert.True(t, addressesMatch(featureAddr1, featureAddr2))
}

// Comprehensive tests for all data type implementations
func TestAllListDataTypeImplementations(t *testing.T) {
	device1 := AddressDeviceType("device1")

	tests := []struct {
		name          string
		dataType      PartialReader
		filter        *FilterType
		expectedCount int
		shouldSucceed bool
		description   string
	}{
		{
			name: "ElectricalConnectionStateListDataType_WithFilter",
			dataType: &ElectricalConnectionStateListDataType{
				ElectricalConnectionStateData: []ElectricalConnectionStateDataType{
					{ElectricalConnectionId: Ptr(ElectricalConnectionIdType(1))},
					{ElectricalConnectionId: Ptr(ElectricalConnectionIdType(2))},
				},
			},
			filter: &FilterType{
				ElectricalConnectionStateListDataSelectors: &ElectricalConnectionStateListDataSelectorsType{
					ElectricalConnectionId: Ptr(ElectricalConnectionIdType(1)),
				},
			},
			expectedCount: 1,
			shouldSucceed: true,
			description:   "Should filter by electrical connection ID",
		},
		{
			name: "MeasurementListDataType_WithFilter",
			dataType: &MeasurementListDataType{
				MeasurementData: []MeasurementDataType{
					{MeasurementId: Ptr(MeasurementIdType(1))},
					{MeasurementId: Ptr(MeasurementIdType(2))},
					{MeasurementId: Ptr(MeasurementIdType(3))},
				},
			},
			filter: &FilterType{
				MeasurementListDataSelectors: &MeasurementListDataSelectorsType{
					MeasurementId: Ptr(MeasurementIdType(2)),
				},
			},
			expectedCount: 1,
			shouldSucceed: true,
			description:   "Should filter by measurement ID",
		},
		{
			name: "SessionMeasurementRelationListDataType_WithArrayFilter",
			dataType: &SessionMeasurementRelationListDataType{
				SessionMeasurementRelationData: []SessionMeasurementRelationDataType{
					{
						SessionId:     Ptr(SessionIdType(1)),
						MeasurementId: []MeasurementIdType{MeasurementIdType(10), MeasurementIdType(20)},
					},
					{
						SessionId:     Ptr(SessionIdType(2)),
						MeasurementId: []MeasurementIdType{MeasurementIdType(30)},
					},
				},
			},
			filter: &FilterType{
				SessionMeasurementRelationListDataSelectors: &SessionMeasurementRelationListDataSelectorsType{
					MeasurementId: Ptr(MeasurementIdType(10)), // Should match first item's array
				},
			},
			expectedCount: 1,
			shouldSucceed: true,
			description:   "Should filter by measurement ID in array",
		},
		{
			name: "NetworkManagementDeviceDescriptionListDataType_WithDeviceAddress",
			dataType: &NetworkManagementDeviceDescriptionListDataType{
				NetworkManagementDeviceDescriptionData: []NetworkManagementDeviceDescriptionDataType{
					{DeviceAddress: &DeviceAddressType{Device: &device1}},
					{DeviceAddress: &DeviceAddressType{Device: Ptr(AddressDeviceType("device2"))}},
				},
			},
			filter: &FilterType{
				NetworkManagementDeviceDescriptionListDataSelectors: &NetworkManagementDeviceDescriptionListDataSelectorsType{
					DeviceAddress: &DeviceAddressType{Device: &device1},
				},
			},
			expectedCount: 1,
			shouldSucceed: true,
			description:   "Should filter by device address",
		},
		{
			name: "SpecificationVersionListDataType_AlwaysReturnsAll",
			dataType: &SpecificationVersionListDataType{
				SpecificationVersionData: []SpecificationVersionDataType{
					SpecificationVersionDataType("1.0.0"),
					SpecificationVersionDataType("2.0.0"),
				},
			},
			filter: &FilterType{
				// Empty selector type - should always return all
			},
			expectedCount: 2,
			shouldSucceed: true,
			description:   "Should always return all data (empty selector)",
		},
		{
			name: "ElectricalConnectionCharacteristicListDataType_WithMultipleFilters",
			dataType: &ElectricalConnectionCharacteristicListDataType{
				ElectricalConnectionCharacteristicData: []ElectricalConnectionCharacteristicDataType{
					{ElectricalConnectionId: Ptr(ElectricalConnectionIdType(1)), CharacteristicId: Ptr(ElectricalConnectionCharacteristicIdType(10))},
					{ElectricalConnectionId: Ptr(ElectricalConnectionIdType(1)), CharacteristicId: Ptr(ElectricalConnectionCharacteristicIdType(20))},
					{ElectricalConnectionId: Ptr(ElectricalConnectionIdType(2)), CharacteristicId: Ptr(ElectricalConnectionCharacteristicIdType(10))},
				},
			},
			filter: &FilterType{
				ElectricalConnectionCharacteristicListDataSelectors: &ElectricalConnectionCharacteristicListDataSelectorsType{
					ElectricalConnectionId: Ptr(ElectricalConnectionIdType(1)),
					CharacteristicId:       Ptr(ElectricalConnectionCharacteristicIdType(10)),
				},
			},
			expectedCount: 1,
			shouldSucceed: true,
			description:   "Should filter by both connection and characteristic ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test with nil filter first
			result, success := tt.dataType.ReadPartialData(nil)
			require.True(t, success, "ReadPartialData with nil filter should always succeed")
			require.NotNil(t, result, "Result should not be nil")

			// Test with the specific filter
			if tt.filter != nil {
				result, success := tt.dataType.ReadPartialData(tt.filter)
				if tt.shouldSucceed {
					require.True(t, success, "ReadPartialData should succeed: %s", tt.description)
					require.NotNil(t, result, "Result should not be nil")

					// Use reflection to check the length of the filtered data
					resultValue := reflect.ValueOf(result).Elem()
					for i := 0; i < resultValue.NumField(); i++ {
						field := resultValue.Field(i)
						if field.Kind() == reflect.Slice {
							assert.Equal(t, tt.expectedCount, field.Len(), "Filtered result count should match expected: %s", tt.description)
							break
						}
					}
				} else {
					require.False(t, success, "ReadPartialData should fail: %s", tt.description)
				}
			}
		})
	}
}

// Test edge cases and error conditions
func TestEdgeCasesAndErrorConditions(t *testing.T) {
	t.Run("EmptyData", func(t *testing.T) {
		emptyData := &MeasurementListDataType{
			MeasurementData: []MeasurementDataType{},
		}

		filter := &FilterType{
			MeasurementListDataSelectors: &MeasurementListDataSelectorsType{
				MeasurementId: Ptr(MeasurementIdType(1)),
			},
		}

		result, success := emptyData.ReadPartialData(filter)
		assert.True(t, success)
		resultData := result.(*MeasurementListDataType)
		assert.Len(t, resultData.MeasurementData, 0)
	})

	t.Run("InvalidFilterType", func(t *testing.T) {
		data := &MeasurementListDataType{
			MeasurementData: []MeasurementDataType{
				{MeasurementId: Ptr(MeasurementIdType(1))},
			},
		}

		// Create filter with wrong selector type (should still work due to generic handling)
		filter := &FilterType{
			ElectricalConnectionStateListDataSelectors: &ElectricalConnectionStateListDataSelectorsType{
				ElectricalConnectionId: Ptr(ElectricalConnectionIdType(1)),
			},
		}

		result, success := data.ReadPartialData(filter)
		assert.True(t, success) // Should succeed and return all data
		resultData := result.(*MeasurementListDataType)
		assert.Len(t, resultData.MeasurementData, 1) // All data returned
	})

	t.Run("ElementsFiltering", func(t *testing.T) {
		data := &MeasurementListDataType{
			MeasurementData: []MeasurementDataType{
				{
					MeasurementId: Ptr(MeasurementIdType(1)),
					Value:         &ScaledNumberType{Number: Ptr(NumberType(100))},
					Timestamp:     Ptr(AbsoluteOrRelativeTimeType("2023-01-01T00:00:00Z")),
				},
			},
		}

		// Filter to include only MeasurementId
		filter := &FilterType{
			MeasurementDataElements: &MeasurementDataElementsType{
				MeasurementId: Ptr(ElementTagType{}),
				// Value and Timestamp not included
			},
		}

		result, success := data.ReadPartialData(filter)
		assert.True(t, success)
		resultData := result.(*MeasurementListDataType)
		assert.Len(t, resultData.MeasurementData, 1)

		// Check that only MeasurementId is included
		item := resultData.MeasurementData[0]
		assert.NotNil(t, item.MeasurementId)
		assert.Nil(t, item.Value)     // Should be filtered out
		assert.Nil(t, item.Timestamp) // Should be filtered out
	})
}

// Test the new FieldMatcher interface and implementations
func TestFieldMatcherInterface(t *testing.T) {
	device1 := AddressDeviceType("device1")
	device2 := AddressDeviceType("device2")
	entity1 := AddressEntityType(1)
	feature1 := AddressFeatureType(100)

	// Test AddressMatcher
	t.Run("AddressMatcher", func(t *testing.T) {
		matcher := AddressMatcher{}

		// Test FeatureAddress matching
		itemAddr := &FeatureAddressType{
			Device:  &device1,
			Entity:  []AddressEntityType{entity1},
			Feature: &feature1,
		}

		// Partial selector - only device specified
		selectorAddr := &FeatureAddressType{
			Device: &device1,
		}

		itemField := reflect.ValueOf(itemAddr)
		selectorField := reflect.ValueOf(selectorAddr)

		assert.True(t, matcher.MatchesField(itemField, selectorField), "Should match when selector device matches")

		// Non-matching device
		nonMatchingSelector := &FeatureAddressType{
			Device: &device2,
		}
		nonMatchingSelectorField := reflect.ValueOf(nonMatchingSelector)
		assert.False(t, matcher.MatchesField(itemField, nonMatchingSelectorField), "Should not match when selector device differs")

		// Nil selector should match all
		nilSelectorField := reflect.ValueOf((*FeatureAddressType)(nil))
		assert.True(t, matcher.MatchesField(itemField, nilSelectorField), "Nil selector should match all")

		// Test DeviceAddress matching
		deviceItemAddr := &DeviceAddressType{Device: &device1}
		deviceSelectorAddr := &DeviceAddressType{Device: &device1}
		deviceItemField := reflect.ValueOf(deviceItemAddr)
		deviceSelectorField := reflect.ValueOf(deviceSelectorAddr)

		assert.True(t, matcher.MatchesField(deviceItemField, deviceSelectorField), "Device addresses should match")
	})

	// Test StructMatcher
	t.Run("StructMatcher", func(t *testing.T) {
		matcher := StructMatcher{}

		type TestStruct struct {
			ID   *int
			Name *string
		}

		item := TestStruct{ID: Ptr(1), Name: Ptr("test")}
		selector := TestStruct{ID: Ptr(1)}

		itemField := reflect.ValueOf(item)
		selectorField := reflect.ValueOf(selector)

		assert.True(t, matcher.MatchesField(itemField, selectorField), "Should match when struct fields match")

		// Non-matching struct
		nonMatchingSelector := TestStruct{ID: Ptr(2)}
		nonMatchingSelectorField := reflect.ValueOf(nonMatchingSelector)
		assert.False(t, matcher.MatchesField(itemField, nonMatchingSelectorField), "Should not match when struct fields differ")

		// Nil selector
		nilSelectorField := reflect.ValueOf((*TestStruct)(nil))
		assert.True(t, matcher.MatchesField(itemField, nilSelectorField), "Nil selector should match all")
	})

	// Test ArrayMatcher
	t.Run("ArrayMatcher", func(t *testing.T) {
		matcher := ArrayMatcher{}

		itemArray := []string{"tag1", "tag2", "tag3"}
		itemField := reflect.ValueOf(itemArray)

		// Test with single value that exists in array
		selectorValue := "tag2"
		selectorField := reflect.ValueOf(&selectorValue)

		assert.True(t, matcher.MatchesField(itemField, selectorField), "Should match when selector value exists in array")

		// Test with value that doesn't exist in array
		nonExistingValue := "tag4"
		nonExistingSelectorField := reflect.ValueOf(&nonExistingValue)
		assert.False(t, matcher.MatchesField(itemField, nonExistingSelectorField), "Should not match when selector value doesn't exist in array")

		// Test with nil selector
		nilSelectorField := reflect.ValueOf((*string)(nil))
		assert.True(t, matcher.MatchesField(itemField, nilSelectorField), "Nil selector should match all arrays")
	})

	// Test DefaultMatcher
	t.Run("DefaultMatcher", func(t *testing.T) {
		matcher := DefaultMatcher{}

		item := 42
		selector := 42

		itemField := reflect.ValueOf(item)
		selectorField := reflect.ValueOf(selector)

		assert.True(t, matcher.MatchesField(itemField, selectorField), "Should match equal values")

		// Non-matching values
		nonMatchingSelector := 43
		nonMatchingSelectorField := reflect.ValueOf(nonMatchingSelector)
		assert.False(t, matcher.MatchesField(itemField, nonMatchingSelectorField), "Should not match different values")
	})
}

// Test the new selector-style address matching functions
func TestSelectorStyleAddressMatching(t *testing.T) {
	device1 := AddressDeviceType("device1")
	device2 := AddressDeviceType("device2")
	entity1 := AddressEntityType(1)
	entity2 := AddressEntityType(2)
	feature1 := AddressFeatureType(100)
	feature2 := AddressFeatureType(200)

	t.Run("FeatureAddressMatchesSelector", func(t *testing.T) {
		itemAddr := &FeatureAddressType{
			Device:  &device1,
			Entity:  []AddressEntityType{entity1},
			Feature: &feature1,
		}

		// Test 1: Nil selector should match all
		assert.True(t, featureAddressMatchesSelector(itemAddr, nil))

		// Test 2: Empty selector should match all
		emptySelector := &FeatureAddressType{}
		assert.True(t, featureAddressMatchesSelector(itemAddr, emptySelector))

		// Test 3: Partial match - only device specified
		deviceOnlySelector := &FeatureAddressType{Device: &device1}
		assert.True(t, featureAddressMatchesSelector(itemAddr, deviceOnlySelector))

		// Test 4: Partial match - device and entity specified
		deviceEntitySelector := &FeatureAddressType{
			Device: &device1,
			Entity: []AddressEntityType{entity1},
		}
		assert.True(t, featureAddressMatchesSelector(itemAddr, deviceEntitySelector))

		// Test 5: Full match
		fullSelector := &FeatureAddressType{
			Device:  &device1,
			Entity:  []AddressEntityType{entity1},
			Feature: &feature1,
		}
		assert.True(t, featureAddressMatchesSelector(itemAddr, fullSelector))

		// Test 6: Non-matching device
		wrongDeviceSelector := &FeatureAddressType{Device: &device2}
		assert.False(t, featureAddressMatchesSelector(itemAddr, wrongDeviceSelector))

		// Test 7: Non-matching entity
		wrongEntitySelector := &FeatureAddressType{
			Device: &device1,
			Entity: []AddressEntityType{entity2},
		}
		assert.False(t, featureAddressMatchesSelector(itemAddr, wrongEntitySelector))

		// Test 8: Non-matching feature
		wrongFeatureSelector := &FeatureAddressType{
			Device:  &device1,
			Entity:  []AddressEntityType{entity1},
			Feature: &feature2,
		}
		assert.False(t, featureAddressMatchesSelector(itemAddr, wrongFeatureSelector))

		// Test 9: Nil item should not match non-nil selector
		assert.False(t, featureAddressMatchesSelector(nil, deviceOnlySelector))

		// Test 10: Nil item should match nil selector
		assert.True(t, featureAddressMatchesSelector(nil, nil))
	})

	t.Run("DeviceAddressMatchesSelector", func(t *testing.T) {
		itemAddr := &DeviceAddressType{Device: &device1}

		// Test 1: Nil selector should match all
		assert.True(t, deviceAddressMatchesSelector(itemAddr, nil))

		// Test 2: Empty selector should match all
		emptySelector := &DeviceAddressType{}
		assert.True(t, deviceAddressMatchesSelector(itemAddr, emptySelector))

		// Test 3: Matching device
		matchingSelector := &DeviceAddressType{Device: &device1}
		assert.True(t, deviceAddressMatchesSelector(itemAddr, matchingSelector))

		// Test 4: Non-matching device
		nonMatchingSelector := &DeviceAddressType{Device: &device2}
		assert.False(t, deviceAddressMatchesSelector(itemAddr, nonMatchingSelector))

		// Test 5: Nil item should not match non-nil selector
		assert.False(t, deviceAddressMatchesSelector(nil, matchingSelector))

		// Test 6: Nil item should match nil selector
		assert.True(t, deviceAddressMatchesSelector(nil, nil))
	})

	t.Run("EntityAddressMatchesSelector", func(t *testing.T) {
		itemAddr := &EntityAddressType{
			Device: &device1,
			Entity: []AddressEntityType{entity1},
		}

		// Test 1: Nil selector should match all
		assert.True(t, entityAddressMatchesSelector(itemAddr, nil))

		// Test 2: Empty selector should match all
		emptySelector := &EntityAddressType{}
		assert.True(t, entityAddressMatchesSelector(itemAddr, emptySelector))

		// Test 3: Device-only match
		deviceOnlySelector := &EntityAddressType{Device: &device1}
		assert.True(t, entityAddressMatchesSelector(itemAddr, deviceOnlySelector))

		// Test 4: Full match
		fullSelector := &EntityAddressType{
			Device: &device1,
			Entity: []AddressEntityType{entity1},
		}
		assert.True(t, entityAddressMatchesSelector(itemAddr, fullSelector))

		// Test 5: Non-matching device
		wrongDeviceSelector := &EntityAddressType{Device: &device2}
		assert.False(t, entityAddressMatchesSelector(itemAddr, wrongDeviceSelector))

		// Test 6: Non-matching entity
		wrongEntitySelector := &EntityAddressType{
			Device: &device1,
			Entity: []AddressEntityType{entity2},
		}
		assert.False(t, entityAddressMatchesSelector(itemAddr, wrongEntitySelector))

		// Test 7: Multi-entity address
		multiEntityAddr := &EntityAddressType{
			Device: &device1,
			Entity: []AddressEntityType{entity1, entity2},
		}
		multiEntitySelector := &EntityAddressType{
			Device: &device1,
			Entity: []AddressEntityType{entity1, entity2},
		}
		assert.True(t, entityAddressMatchesSelector(multiEntityAddr, multiEntitySelector))

		// Test 8: Different entity array length
		shortEntitySelector := &EntityAddressType{
			Device: &device1,
			Entity: []AddressEntityType{entity1},
		}
		assert.False(t, entityAddressMatchesSelector(multiEntityAddr, shortEntitySelector))
	})
}

// Test the getFieldMatcher function
func TestGetFieldMatcher(t *testing.T) {
	// Test array field detection
	arrayField := reflect.ValueOf([]string{"test"})
	singleField := reflect.ValueOf("test")
	matcher := getFieldMatcher(arrayField, singleField)
	assert.IsType(t, ArrayMatcher{}, matcher, "Should return ArrayMatcher for array fields")

	// Test address field detection
	device := AddressDeviceType("device1")
	deviceAddr := &DeviceAddressType{Device: &device}
	deviceField := reflect.ValueOf(deviceAddr)
	normalField := reflect.ValueOf("test")
	matcher = getFieldMatcher(deviceField, normalField)
	assert.IsType(t, AddressMatcher{}, matcher, "Should return AddressMatcher for address fields")

	// Test struct field detection
	type TestStruct struct {
		ID int
	}
	structValue := TestStruct{ID: 1}
	structField := reflect.ValueOf(structValue)
	matcher = getFieldMatcher(structField, structField)
	assert.IsType(t, StructMatcher{}, matcher, "Should return StructMatcher for struct fields")

	// Test default case
	intField := reflect.ValueOf(42)
	matcher = getFieldMatcher(intField, intField)
	assert.IsType(t, DefaultMatcher{}, matcher, "Should return DefaultMatcher for regular fields")
}

// Test type detection helper functions
func TestTypeDetectionHelpers(t *testing.T) {
	device := AddressDeviceType("device1")

	t.Run("IsAddressType", func(t *testing.T) {
		deviceAddr := &DeviceAddressType{Device: &device}
		entityAddr := &EntityAddressType{Device: &device}
		feature := AddressFeatureType(1)
		featureAddr := &FeatureAddressType{Device: &device, Feature: &feature}

		assert.True(t, isAddressType(reflect.ValueOf(deviceAddr)), "DeviceAddressType should be detected as address type")
		assert.True(t, isAddressType(reflect.ValueOf(entityAddr)), "EntityAddressType should be detected as address type")
		assert.True(t, isAddressType(reflect.ValueOf(featureAddr)), "FeatureAddressType should be detected as address type")

		// Non-address types
		assert.False(t, isAddressType(reflect.ValueOf("string")), "String should not be detected as address type")
		assert.False(t, isAddressType(reflect.ValueOf(42)), "Int should not be detected as address type")
		assert.False(t, isAddressType(reflect.ValueOf([]string{})), "Slice should not be detected as address type")
	})

	t.Run("IsStructType", func(t *testing.T) {
		type TestStruct struct {
			ID int
		}

		structValue := TestStruct{ID: 1}
		ptrToStruct := &structValue

		assert.True(t, isStructType(reflect.ValueOf(structValue)), "Struct should be detected as struct type")
		assert.True(t, isStructType(reflect.ValueOf(ptrToStruct)), "Pointer to struct should be detected as struct type")

		// Non-struct types
		assert.False(t, isStructType(reflect.ValueOf("string")), "String should not be detected as struct type")
		assert.False(t, isStructType(reflect.ValueOf(42)), "Int should not be detected as struct type")
		assert.False(t, isStructType(reflect.ValueOf([]string{})), "Slice should not be detected as struct type")

		// Nil pointer should not be detected as struct
		var nilPtr *TestStruct
		assert.False(t, isStructType(reflect.ValueOf(nilPtr)), "Nil pointer should not be detected as struct type")
	})
}

// Test edge cases for address matching
func TestAddressMatchingEdgeCases(t *testing.T) {
	device1 := AddressDeviceType("device1")
	entity1 := AddressEntityType(1)
	feature1 := AddressFeatureType(100)

	t.Run("FeatureAddressEdgeCases", func(t *testing.T) {
		// Test with nil device in item
		itemWithNilDevice := &FeatureAddressType{
			Entity:  []AddressEntityType{entity1},
			Feature: &feature1,
		}
		selectorWithDevice := &FeatureAddressType{Device: &device1}
		assert.False(t, featureAddressMatchesSelector(itemWithNilDevice, selectorWithDevice))

		// Test with nil entity in item but selector has entity
		itemWithNilEntity := &FeatureAddressType{
			Device:  &device1,
			Feature: &feature1,
		}
		selectorWithEntity := &FeatureAddressType{
			Device: &device1,
			Entity: []AddressEntityType{entity1},
		}
		assert.False(t, featureAddressMatchesSelector(itemWithNilEntity, selectorWithEntity))

		// Test with empty entity array in item
		itemWithEmptyEntity := &FeatureAddressType{
			Device:  &device1,
			Entity:  []AddressEntityType{},
			Feature: &feature1,
		}
		assert.False(t, featureAddressMatchesSelector(itemWithEmptyEntity, selectorWithEntity))

		// Test with nil feature in item but selector has feature
		itemWithNilFeature := &FeatureAddressType{
			Device: &device1,
			Entity: []AddressEntityType{entity1},
		}
		selectorWithFeature := &FeatureAddressType{
			Device:  &device1,
			Feature: &feature1,
		}
		assert.False(t, featureAddressMatchesSelector(itemWithNilFeature, selectorWithFeature))
	})

	t.Run("DeviceAddressEdgeCases", func(t *testing.T) {
		// Test with nil device in item
		itemWithNilDevice := &DeviceAddressType{}
		selectorWithDevice := &DeviceAddressType{Device: &device1}
		assert.False(t, deviceAddressMatchesSelector(itemWithNilDevice, selectorWithDevice))
	})

	t.Run("EntityAddressEdgeCases", func(t *testing.T) {
		// Test with nil device in item
		itemWithNilDevice := &EntityAddressType{
			Entity: []AddressEntityType{entity1},
		}
		selectorWithDevice := &EntityAddressType{Device: &device1}
		assert.False(t, entityAddressMatchesSelector(itemWithNilDevice, selectorWithDevice))

		// Test with nil entity in item but selector has entity
		itemWithNilEntity := &EntityAddressType{
			Device: &device1,
		}
		selectorWithEntity := &EntityAddressType{
			Device: &device1,
			Entity: []AddressEntityType{entity1},
		}
		assert.False(t, entityAddressMatchesSelector(itemWithNilEntity, selectorWithEntity))
	})
}
