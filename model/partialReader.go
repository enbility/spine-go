package model

import "reflect"

// Updater interface for updating lists in data types
// Key Components:
// - matchesSelector(): Generic function using reflection to match selector criteria
// - filterListDataBySelectors(): Generic function for selector-based filtering
// - partialListDataRead(): Complete generic function with both selector and elements filtering
//
// Usage Pattern:
// Each list data type implementation follows this simple pattern:
//   filteredItems, success := partialListDataRead[DataType, SelectorType, ElementsType](data, filter)
//
// This approach automatically handles:
// - Reflection-based field matching in selectors
// - Special cases like address types and nested structs
// - Array/slice field matching (e.g., MeasurementId arrays)
// - Elements-based filtering

type PartialReader interface {
	// ReadPartialData reads partial data from the model.
	// The data is filtered by the given filter.
	// If the filter is nil, all data is returned.
	// Returns the filtered data and a boolean indicating success.
	ReadPartialData(filter *FilterType) (any, bool)
}

// Helper functions for address matching
func addressesMatch(addr1, addr2 *FeatureAddressType) bool {
	if addr1 == nil || addr2 == nil {
		return addr1 == addr2
	}

	// Compare device addresses
	if (addr1.Device == nil) != (addr2.Device == nil) {
		return false
	}
	if addr1.Device != nil && addr2.Device != nil && *addr1.Device != *addr2.Device {
		return false
	}

	// Compare entity addresses (Entity is a slice)
	if (addr1.Entity == nil) != (addr2.Entity == nil) {
		return false
	}
	if len(addr1.Entity) != len(addr2.Entity) {
		return false
	}
	for i, entity1 := range addr1.Entity {
		if entity1 != addr2.Entity[i] {
			return false
		}
	}

	// Compare feature addresses
	if (addr1.Feature == nil) != (addr2.Feature == nil) {
		return false
	}
	if addr1.Feature != nil && addr2.Feature != nil && *addr1.Feature != *addr2.Feature {
		return false
	}

	return true
}

func deviceAddressesMatch(addr1, addr2 *DeviceAddressType) bool {

	if addr1 == nil || addr2 == nil {
		return addr1 == addr2
	}

	return (addr1.Device == nil && addr2.Device == nil) ||
		(addr1.Device != nil && addr2.Device != nil && *addr1.Device == *addr2.Device)
}

func entityAddressesMatch(addr1, addr2 *EntityAddressType) bool {
	if addr1 == nil || addr2 == nil {
		return addr1 == addr2
	}

	// Compare device part
	if (addr1.Device == nil) != (addr2.Device == nil) {
		return false
	}
	if addr1.Device != nil && addr2.Device != nil && *addr1.Device != *addr2.Device {
		return false
	}

	// Compare entity part (Entity is a slice)
	if (addr1.Entity == nil) != (addr2.Entity == nil) {
		return false
	}
	if len(addr1.Entity) != len(addr2.Entity) {
		return false
	}
	for i, entity1 := range addr1.Entity {
		if entity1 != addr2.Entity[i] {
			return false
		}
	}

	return true
}

func featureAddressesMatch(addr1, addr2 *FeatureAddressType) bool {
	return addressesMatch(addr1, addr2)
}

// FieldMatcher interface for custom field comparison logic
type FieldMatcher interface {
	MatchesField(itemField, selectorField reflect.Value) bool
}

// AddressMatcher handles all address type comparisons
type AddressMatcher struct{}

func (am AddressMatcher) MatchesField(itemField, selectorField reflect.Value) bool {
	if selectorField.Kind() == reflect.Ptr && selectorField.IsNil() {
		return true
	}

	if itemField.Kind() == reflect.Ptr && itemField.IsNil() {
		return false
	}

	// Use selector-style address matching based on type
	if itemField.CanInterface() && selectorField.CanInterface() {
		if itemAddr, ok := itemField.Interface().(*FeatureAddressType); ok {
			if selectorAddr, ok := selectorField.Interface().(*FeatureAddressType); ok {
				return featureAddressMatchesSelector(itemAddr, selectorAddr)
			}
		}
		if itemAddr, ok := itemField.Interface().(*DeviceAddressType); ok {
			if selectorAddr, ok := selectorField.Interface().(*DeviceAddressType); ok {
				return deviceAddressMatchesSelector(itemAddr, selectorAddr)
			}
		}
		if itemAddr, ok := itemField.Interface().(*EntityAddressType); ok {
			if selectorAddr, ok := selectorField.Interface().(*EntityAddressType); ok {
				return entityAddressMatchesSelector(itemAddr, selectorAddr)
			}
		}
	}

	// Fall back to regular field comparison
	return compareFields(itemField, selectorField)
}

// Selector-style address matching functions (nil fields in selector are ignored)
func featureAddressMatchesSelector(itemAddr, selectorAddr *FeatureAddressType) bool {
	if selectorAddr == nil {
		return true
	}
	if itemAddr == nil {
		return false
	}

	// Compare device addresses - nil selector device matches all
	if selectorAddr.Device != nil {
		if itemAddr.Device == nil || *itemAddr.Device != *selectorAddr.Device {
			return false
		}
	}

	// Compare entity addresses - nil selector entity matches all
	if selectorAddr.Entity != nil {
		if itemAddr.Entity == nil {
			return false
		}
		if len(selectorAddr.Entity) != len(itemAddr.Entity) {
			return false
		}
		for i, selectorEntity := range selectorAddr.Entity {
			if selectorEntity != itemAddr.Entity[i] {
				return false
			}
		}
	}

	// Compare feature addresses - nil selector feature matches all
	if selectorAddr.Feature != nil {
		if itemAddr.Feature == nil || *itemAddr.Feature != *selectorAddr.Feature {
			return false
		}
	}

	return true
}

func deviceAddressMatchesSelector(itemAddr, selectorAddr *DeviceAddressType) bool {
	if selectorAddr == nil {
		return true
	}
	if itemAddr == nil {
		return false
	}

	// Compare device addresses - nil selector device matches all
	if selectorAddr.Device != nil {
		if itemAddr.Device == nil || *itemAddr.Device != *selectorAddr.Device {
			return false
		}
	}

	return true
}

func entityAddressMatchesSelector(itemAddr, selectorAddr *EntityAddressType) bool {
	if selectorAddr == nil {
		return true
	}
	if itemAddr == nil {
		return false
	}

	// Compare device addresses - nil selector device matches all
	if selectorAddr.Device != nil {
		if itemAddr.Device == nil || *itemAddr.Device != *selectorAddr.Device {
			return false
		}
	}

	// Compare entity addresses - nil selector entity matches all
	if selectorAddr.Entity != nil {
		if itemAddr.Entity == nil {
			return false
		}
		if len(selectorAddr.Entity) != len(itemAddr.Entity) {
			return false
		}
		for i, selectorEntity := range selectorAddr.Entity {
			if selectorEntity != itemAddr.Entity[i] {
				return false
			}
		}
	}

	return true
}

// StructMatcher handles nested struct comparisons
type StructMatcher struct{}

func (sm StructMatcher) MatchesField(itemField, selectorField reflect.Value) bool {
	if selectorField.Kind() == reflect.Ptr && selectorField.IsNil() {
		return true
	}

	if itemField.Kind() == reflect.Ptr && itemField.IsNil() {
		return selectorField.Kind() == reflect.Ptr && selectorField.IsNil()
	}

	// Recursively match using the generic selector matching
	return matchesSelector(itemField.Interface(), selectorField.Interface())
}

// ArrayMatcher handles array/slice comparisons
type ArrayMatcher struct{}

func (am ArrayMatcher) MatchesField(itemField, selectorField reflect.Value) bool {
	return matchesArrayField(itemField, selectorField)
}

// DefaultMatcher handles regular field comparisons
type DefaultMatcher struct{}

func (dm DefaultMatcher) MatchesField(itemField, selectorField reflect.Value) bool {
	return compareFields(itemField, selectorField)
}

// getFieldMatcher returns the appropriate matcher based on field types
func getFieldMatcher(itemField, selectorField reflect.Value) FieldMatcher {
	// Check if either field is an array/slice
	if (itemField.Kind() == reflect.Slice || itemField.Kind() == reflect.Array) ||
		(selectorField.Kind() == reflect.Slice || selectorField.Kind() == reflect.Array) {
		return ArrayMatcher{}
	}

	// Check if fields contain address types (by interface type)
	if isAddressType(itemField) || isAddressType(selectorField) {
		return AddressMatcher{}
	}

	// Check if fields are structs (for nested struct matching)
	if isStructType(itemField) && isStructType(selectorField) {
		return StructMatcher{}
	}

	// Default field comparison
	return DefaultMatcher{}
}

// isAddressType checks if a field contains an address type
func isAddressType(field reflect.Value) bool {
	if !field.CanInterface() {
		return false
	}

	switch field.Interface().(type) {
	case *FeatureAddressType, *DeviceAddressType, *EntityAddressType:
		return true
	}
	return false
}

// isStructType checks if a field is a struct or pointer to struct
func isStructType(field reflect.Value) bool {
	if field.Kind() == reflect.Ptr && !field.IsNil() {
		return field.Elem().Kind() == reflect.Struct
	}
	return field.Kind() == reflect.Struct
}

// Generic helper function using reflection to match selector criteria
// This eliminates the need for custom matcher functions by automatically comparing
// all non-nil fields in the selector with corresponding fields in the item
func matchesSelector(item any, selector any) bool {
	if selector == nil {
		return true
	}

	itemValue := reflect.ValueOf(item)
	selectorValue := reflect.ValueOf(selector)

	// Handle pointer types
	if selectorValue.Kind() == reflect.Ptr {
		if selectorValue.IsNil() {
			return true
		}
		selectorValue = selectorValue.Elem()
	}

	// Get the type for field iteration
	selectorType := selectorValue.Type()
	allFieldsNil := true

	// Iterate through all fields in the selector
	for i := 0; i < selectorValue.NumField(); i++ {
		selectorField := selectorValue.Field(i)
		selectorFieldType := selectorType.Field(i)

		// Skip unexported fields
		if !selectorField.CanInterface() {
			continue
		}

		// Check if the selector field is nil or zero value
		if isFieldNilOrEmpty(selectorField) {
			continue
		}

		allFieldsNil = false
		fieldName := selectorFieldType.Name

		// Get the corresponding field from the item
		itemField := getFieldByName(itemValue, fieldName)
		if !itemField.IsValid() {
			return false // Item doesn't have this field
		}

		// Use the appropriate matcher based on field types
		matcher := getFieldMatcher(itemField, selectorField)
		if !matcher.MatchesField(itemField, selectorField) {
			return false
		}
	}

	// If all selector fields are nil, include all items
	return allFieldsNil || true
}

// Helper function to check if a field is nil or empty
func isFieldNilOrEmpty(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.Ptr:
		return field.IsNil()
	case reflect.Slice, reflect.Array:
		return field.Len() == 0
	case reflect.String:
		return field.String() == ""
	case reflect.Interface:
		return field.IsNil()
	default:
		return field.IsZero()
	}
}

// Helper function to get a field by name from a struct, handling pointers
func getFieldByName(structValue reflect.Value, fieldName string) reflect.Value {
	if structValue.Kind() == reflect.Ptr {
		if structValue.IsNil() {
			return reflect.Value{}
		}
		structValue = structValue.Elem()
	}

	if structValue.Kind() != reflect.Struct {
		return reflect.Value{}
	}

	return structValue.FieldByName(fieldName)
}

// Helper function to match array/slice fields
func matchesArrayField(itemField, selectorField reflect.Value) bool {
	// Handle case where selector field is also a slice/array (like PowerSequenceDescriptionListDataSelectorsType.SequenceId)
	if selectorField.Kind() == reflect.Slice || selectorField.Kind() == reflect.Array {
		if selectorField.Len() == 0 {
			return true // Empty selector array matches all
		}

		// Check if the item field value exists in any of the selector array values
		if itemField.Kind() == reflect.Ptr && !itemField.IsNil() {
			itemValue := itemField.Elem()
			for i := 0; i < selectorField.Len(); i++ {
				selectorElement := selectorField.Index(i)
				if compareFields(itemValue, selectorElement) {
					return true
				}
			}
			return false
		} else if itemField.Kind() != reflect.Ptr {
			// Direct comparison with non-pointer item field
			for i := 0; i < selectorField.Len(); i++ {
				selectorElement := selectorField.Index(i)
				if compareFields(itemField, selectorElement) {
					return true
				}
			}
			return false
		}
		return false
	}

	// Original logic: If selector field is a pointer, check if selector value exists in item array
	if selectorField.Kind() == reflect.Ptr && !selectorField.IsNil() {
		// Check if the selector value exists in the item array
		for i := 0; i < itemField.Len(); i++ {
			itemElement := itemField.Index(i)
			if compareFields(itemElement, selectorField) {
				return true
			}
		}
		return false
	}

	return true // If selector is nil, match all
}

// Helper function to compare two fields
func compareFields(itemField, selectorField reflect.Value) bool {
	// Handle pointer fields
	if selectorField.Kind() == reflect.Ptr {
		if selectorField.IsNil() {
			return true // nil selector matches anything
		}

		if itemField.Kind() == reflect.Ptr {
			if itemField.IsNil() {
				return false // non-nil selector doesn't match nil item
			}
			return compareFields(itemField.Elem(), selectorField.Elem())
		} else {
			return compareFields(itemField, selectorField.Elem())
		}
	}

	if itemField.Kind() == reflect.Ptr {
		if itemField.IsNil() {
			return false // item is nil but selector is not
		}
		return compareFields(itemField.Elem(), selectorField)
	}

	// Direct value comparison
	if itemField.CanInterface() && selectorField.CanInterface() {
		return itemField.Interface() == selectorField.Interface()
	}

	return false
}

// Generic helper function using reflection for selector matching
func filterListDataBySelectors[T any, S any](data []T, filter *FilterType) ([]T, bool) {
	if filter == nil {
		return data, true
	}

	filterData, err := filter.Data(nil)
	if err != nil {
		return data, false
	}

	result := make([]T, 0)

	if filterData.Selector == nil {
		result = append(result, data...)
		return result, true
	}

	if selector, ok := filterData.Selector.(*S); ok {
		for _, item := range data {
			if matchesSelector(item, selector) {
				result = append(result, item)
			}
		}
	} else {
		// If selector type doesn't match, return all data
		result = append(result, data...)
	}

	return result, true
}

// Complete generic helper function with both selector and elements filtering
func partialListDataRead[T any, S any, E any](data []T, filter *FilterType) ([]T, bool) {
	if filter == nil {
		return data, true
	}

	// First apply selector filtering
	filteredItems, success := filterListDataBySelectors[T, S](data, filter)
	if !success {
		return data, false
	}

	// Then apply elements filtering if specified
	filterData, err := filter.Data(nil)
	if err != nil {
		return data, false
	}

	if filterData.Elements != nil {
		if elements, ok := filterData.Elements.(*E); ok {
			filteredItems = filterElementsFromItems(filteredItems, elements)
		}
	}

	return filteredItems, true
}

// Helper function to filter elements from items based on the Elements filter
// This function creates a copy of each item with only the fields specified in Elements
func filterElementsFromItems[T any, E any](items []T, elements E) []T {
	// Check if elements is nil by using reflection
	elementsValue := reflect.ValueOf(elements)
	if elementsValue.Kind() == reflect.Ptr && elementsValue.IsNil() {
		return items
	}

	elementsToInclude := getNonNilElementNames(elements)
	if len(elementsToInclude) == 0 {
		return items
	}

	result := make([]T, len(items))
	for i, item := range items {
		filteredItem := createFilteredItem(item, elementsToInclude)
		result[i] = filteredItem
	}

	return result
}

// Helper function to get the names of non-nil fields from an Elements struct
func getNonNilElementNames(element any) []string {
	var result []string

	v := reflect.ValueOf(element).Elem()
	t := reflect.TypeOf(element).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !isElementFieldNil(field.Interface()) {
			name := t.Field(i).Name
			result = append(result, name)
		}
	}

	return result
}

// Helper function to create a filtered copy of an item with only specified fields
func createFilteredItem[T any](item T, fieldsToInclude []string) T {
	// Create a new zero instance of the same type
	itemType := reflect.TypeOf(item)
	newItem := reflect.New(itemType).Elem()

	// Get the original item's value
	originalValue := reflect.ValueOf(item)

	// Copy only the specified fields
	for _, fieldName := range fieldsToInclude {
		originalField := originalValue.FieldByName(fieldName)
		newField := newItem.FieldByName(fieldName)

		if originalField.IsValid() && newField.IsValid() && newField.CanSet() {
			newField.Set(originalField)
		}
	}

	return newItem.Interface().(T)
}

// Helper function to check if a field value is nil (renamed to avoid redeclaration)
func isElementFieldNil(field interface{}) bool {
	if field == nil {
		return true
	}

	switch reflect.TypeOf(field).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(field).IsNil()
	}
	return false
}
