// Package model provides data structures and update mechanisms for the SPINE protocol.
//
// This file implements the core list update functionality used throughout the SPINE protocol
// to handle partial data updates, filtering, and merging operations while preventing duplicate
// entries and maintaining data consistency.
//
// # SPINE Protocol Context
//
// The SPINE (Smart Premises Interoperable Neutral-message Exchange) protocol requires
// sophisticated data update semantics to handle:
//   - Partial updates from remote devices
//   - Composite key structures with primary and sub-identifiers
//   - Anti-duplication measures for incomplete data
//   - Atomic operations with filtering support
//
// # EEBus Tag System
//
// Data structures use EEBus tags to define field behavior:
//   - `eebus:"key"` - Identifies key fields for uniqueness
//   - `eebus:"key,primarykey"` - Primary identifier in composite keys
//   - `eebus:"writecheck"` - Controls write permissions
//
// Example usage:
//
//	type MeasurementData struct {
//	    MeasurementId *uint   `eebus:"key,primarykey"`  // Primary identifier
//	    ValueType     *string `eebus:"key"`             // Sub-identifier
//	    Value         *int                              // Data field
//	}
//
// # Update Flow
//
// The update process follows this sequence:
//  1. Apply delete filters (removes matching entries)
//  2. Apply partial filters (updates specific fields)
//  3. Filter primary-key-only entries (prevents duplicates)
//  4. Merge remaining data with existing entries
//  5. Sort results by key fields
package model

import (
	"reflect"
	"slices"
	"sort"

	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/util"
)

// Updater defines the interface for data structures that can perform SPINE protocol list updates.
//
// This interface enables any data type to implement custom update logic while maintaining
// consistency with SPINE's Restricted Function Exchange (RFE) requirements.
//
// # Implementation Requirements
//
// Implementations must handle:
//   - Remote vs local write permission checking (remoteWrite parameter)
//   - Atomic operations with proper persistence control
//   - Filter-based partial updates and deletions
//   - Primary key-only entry filtering for duplicate prevention
//
// # Example Implementation
//
//	func (d *MyDataType) UpdateList(remoteWrite, persist bool, newList any,
//	                               filterPartial, filterDelete *FilterType) (any, bool) {
//	    if newData, ok := newList.([]MyDataType); ok {
//	        return UpdateList(remoteWrite, d.existingData, newData, filterPartial, filterDelete)
//	    }
//	    return nil, false
//	}
type Updater interface {
	// UpdateList performs data model specific list updates following SPINE protocol semantics.
	//
	// This method implements the core update logic for handling partial data updates,
	// filtering operations, and merge semantics as defined by the SPINE specification's
	// Restricted Function Exchange (RFE) requirements.
	//
	// # Parameters
	//
	//   - remoteWrite: true if data originates from a remote SPINE device.
	//     When true, write operations are only allowed if the target field's
	//     "writecheck" tagged boolean field is set to true. This enforces
	//     remote write permission semantics.
	//
	//   - persist: true if data should be persisted to storage.
	//     When false, creates temporary datasets for validation or preview
	//     operations without permanent storage.
	//
	//   - newList: the incoming data to be merged. Must be a slice of the
	//     appropriate data type matching the implementing structure.
	//
	//   - filterPartial: optional partial update filter. When provided,
	//     only updates fields in entries matching the filter selectors.
	//
	//   - filterDelete: optional deletion filter. When provided,
	//     removes entries or fields matching the filter criteria.
	//
	//   - cmdFunction is the command function for filter context
	//
	// # Returns
	//
	//   - any: the updated data set after applying all operations
	//   - bool: true if all operations completed successfully, false if any
	//          operation failed (e.g., write permission denied)
	//
	// # SPINE Protocol Compliance
	//
	// Implementations must follow SPINE Table 7 cmdOptions combinations
	// and handle atomic operations according to the protocol specification.
	UpdateList(remoteWrite, persist bool, newList any, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) (any, bool)
}

// UpdateList generates a new list by applying SPINE protocol update rules.
//
// This is the core generic function that implements SPINE's Restricted Function Exchange (RFE)
// semantics as defined in EEBus_SPINE_TS_ProtocolSpecification.pdf chapter 5.3.4.
// It handles partial updates, filtering, and anti-duplication measures critical for
// maintaining data consistency in multi-device SPINE networks.
//
// # Key Features
//
//   - Primary key-only entry filtering: Prevents duplicate entries from incomplete
//     remote data by filtering out entries containing only identifier fields
//   - Composite key support: Handles complex data structures with multiple identifiers
//   - Atomic operations: Ensures all-or-nothing update semantics
//   - Filter support: Implements partial updates and selective deletions
//   - Write permission enforcement: Respects "writecheck" tagged field permissions
//
// # Update Sequence
//
//  1. Delete filtering: Removes entries/fields matching delete filters
//  2. Partial filtering: Updates specific fields in matching entries
//  3. Primary key filtering: Removes entries with only key fields (anti-duplication)
//  4. Identifier handling: Processes entries without complete identifiers
//  5. Data merging: Combines new data with existing entries
//  6. Sorting: Orders results by key fields for consistency
//
// # Type Constraints
//
// Type T must be a struct with EEBus tags defining:
//   - Key fields: `eebus:"key"` for identification
//   - Primary keys: `eebus:"key,primarykey"` for composite key structures
//   - Write permissions: `eebus:"writecheck"` for remote write control
//
// # Example Usage
//
//	existing := []MeasurementDataType{
//	    {MeasurementId: util.Ptr(1), ValueType: util.Ptr("power"), Value: util.Ptr(100)},
//	}
//	new := []MeasurementDataType{
//	    {MeasurementId: util.Ptr(1)}, // Key-only - will be filtered
//	    {MeasurementId: util.Ptr(2), ValueType: util.Ptr("voltage"), Value: util.Ptr(220)},
//	}
//	result, success := UpdateList(false, existing, new, nil, nil)
//	// Result contains: entry 1 unchanged, entry 2 added
//
// For complete, runnable examples demonstrating all features, see example_update_test.go
//
// # Parameters
//
//   - remoteWrite: true if data comes from remote SPINE device (enables write permission checks)
//   - existingData: current data set to be updated
//   - newData: incoming data to merge
//   - filterPartial: optional filter for partial field updates
//   - filterDelete: optional filter for entry/field deletion
//   - cmdFunction is passed to filter.Data() for partial filters without selectors
//
// # Returns
//
//   - []T: updated and sorted data set
//   - bool: true if all operations succeeded, false if any failed
func UpdateList[T any](remoteWrite bool, existingData []T, newData []T, filterPartial, filterDelete *FilterType, cmdFunction *FunctionType) ([]T, bool) {
	success := true

	// STEP 1: Apply delete filters (Selective deletion)
	// Process delete operations first to remove entries or fields before merging.
	// This ensures deletions take precedence over updates in the operation sequence.
	if filterDelete != nil {
		if filterData, err := filterDelete.Data(cmdFunction); err == nil {
			updatedData, noErrors := deleteFilteredData(remoteWrite, existingData, filterData)
			if noErrors {
				existingData = updatedData
			} else {
				success = false
			}
		}
	}

	// STEP 2: Apply partial filters (Selective updates)
	// Process partial update operations to modify specific fields in matching entries.
	// When partial filters are used, skip normal merge processing and return early.
	if filterPartial != nil {
		if filterData, err := filterPartial.Data(cmdFunction); err == nil {
			// Only use selector-based copying if there are actual selectors
			// If there are no selectors, fall through to normal identifier-based merge
			if filterData.Selector != nil {
				newData, noErrors := copyToSelectedData(remoteWrite, existingData, filterData, &newData[0])
				if !noErrors {
					success = false
				}
				return newData, success
			}
		}
	}

	// STEP 3: Filter primary-key-only entries (Anti-duplication)
	// Remove entries that contain only key fields to prevent duplicate/incomplete records.
	// This is critical for SPINE protocol compliance as remote devices often send
	// "structure" messages with only key fields before sending actual data.
	originalCount := len(newData)
	newData = filterPrimaryKeyOnlyEntries(newData)
	if len(newData) == 0 {
		// All entries were filtered out, nothing to update
		if originalCount > 0 {
			logging.Log().Debugf("All %d incoming entries were key-only, no meaningful data to process", originalCount)
		}
		return existingData, success
	}

	// STEP 4: Handle incomplete identifiers (SPINE Table 7 semantics)
	// When entries lack complete key information, apply "update all" semantics
	// by copying the provided data to all existing entries. This follows SPINE
	// specification Table 7 for cmdOptions combinations with classifier "notify".
	// NOTE: SPINE spec is ambiguous about partial identifier handling in composite keys
	if len(newData) > 0 && !HasIdentifiers(newData[0]) {
		// No complete identifiers --> copy data to all existing items
		// This implements SPINE "broadcast update" semantics for incomplete keys
		newData, noErrors := copyToAllData(remoteWrite, existingData, &newData[0])
		if !noErrors {
			success = false
		}
		return newData, success
	}

	// STEP 5: Merge new data with existing entries
	// Combine the filtered new data with existing data using SPINE merge semantics.
	// This handles key matching, field copying, and maintains data consistency.
	result, noErrors := Merge(remoteWrite, existingData, newData)
	if !noErrors {
		success = false
	}

	// STEP 6: Sort results for consistent ordering
	// Ensure deterministic output by sorting entries based on their key fields.
	// This provides predictable results and easier debugging.
	result = SortData(result)

	return result, success
}

// fieldNamesWithEEBusTag extracts field names that contain a specific EEBus tag.
//
// This function uses reflection to inspect struct fields and identify those
// tagged with the specified EEBus tag. It's fundamental to the tag-based
// field processing system used throughout SPINE data operations.
//
// # Supported Tags
//
//   - EEBusTagKey: identifies key/identifier fields
//   - EEBusTagPrimaryKey: identifies primary keys in composite structures
//   - EEBusTagWriteCheck: identifies write permission control fields
//   - EEBusTagFunction: identifies function-specific fields
//   - EEBusTagType: identifies type-specific fields
//
// # Parameters
//
//   - tag: the EEBus tag to search for
//   - item: struct instance to inspect (must be a struct type)
//
// # Returns
//
//   - []string: slice of field names containing the specified tag
//     (empty slice if no matches or item is not a struct)
//
// # Example
//
//	type Data struct {
//	    ID    *uint   `eebus:"key,primarykey"`
//	    SubID *string `eebus:"key"`
//	    Value *int
//	}
//	keys := fieldNamesWithEEBusTag(EEBusTagKey, Data{})
//	// Returns: ["ID", "SubID"]
//	primary := fieldNamesWithEEBusTag(EEBusTagPrimaryKey, Data{})
//	// Returns: ["ID"]
func fieldNamesWithEEBusTag(tag EEBusTag, item any) []string {
	var result []string

	v := reflect.ValueOf(item)
	t := reflect.TypeOf(item)

	if v.Kind() != reflect.Struct {
		return result
	}

	// Iterate through all struct fields using reflection
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		// Only process pointer fields (SPINE protocol requirement)
		if f.Kind() != reflect.Ptr {
			continue
		}

		// Extract EEBus tags from field's struct definition
		sf := v.Type().Field(i)
		eebusTags := EEBusTags(sf)
		// Check if field contains the requested tag
		_, exists := eebusTags[tag]
		if !exists {
			continue
		}

		// Add matching field name to result
		fieldName := t.Field(i).Name
		result = append(result, fieldName)
	}

	return result
}

// HasIdentifiers checks if a struct instance has values for all of its key fields.
//
// This function verifies that all fields tagged with `eebus:"key"` contain
// non-nil values, ensuring the instance has complete identification information.
// This is critical for SPINE protocol operations that require full key specification.
//
// # SPINE Protocol Context
//
// The SPINE specification requires complete identifiers for most operations.
// Incomplete identifiers trigger special "update all" semantics where data
// is copied to all existing entries rather than merged with specific matches.
//
// # Parameters
//
//   - data: struct instance to check (must contain EEBus-tagged key fields)
//
// # Returns
//
//   - bool: true if all key fields have non-nil values, false otherwise
//     (returns true for structs with no key fields)
//
// # Example
//
//	type MeasurementData struct {
//	    MeasurementId *uint   `eebus:"key,primarykey"`
//	    ValueType     *string `eebus:"key"`
//	    Value         *int
//	}
//
//	complete := MeasurementData{
//	    MeasurementId: util.Ptr(uint(1)),
//	    ValueType:     util.Ptr("power"),
//	}
//	HasIdentifiers(complete) // Returns: true
//
//	incomplete := MeasurementData{
//	    MeasurementId: util.Ptr(uint(1)),
//	    // ValueType is nil
//	}
//	HasIdentifiers(incomplete) // Returns: false
func HasIdentifiers(data any) bool {
	keys := fieldNamesWithEEBusTag(EEBusTagKey, data)

	v := reflect.ValueOf(data)

	// Check each key field for non-nil values
	for _, fieldName := range keys {
		f := v.FieldByName(fieldName)

		// If any key field is nil or invalid, identifiers are incomplete
		if f.IsNil() || !f.IsValid() {
			return false
		}
	}

	return true
}

// hasPrimaryKeyOnly determines if an entry contains only primary key fields with no actual data.
//
// This is a critical anti-duplication function that identifies "structural" entries
// sent by remote SPINE devices that contain only identification information.
// Such entries are filtered out to prevent creation of duplicate or incomplete
// records in the local data store.
//
// # Primary Key Detection Strategy
//
// The function uses a hybrid approach to maintain backward compatibility:
//
//  1. For composite key types: Uses `eebus:"primarykey"` tags to distinguish
//     primary identifiers from sub-identifiers
//  2. For single key types: Falls back to simplified detection for types
//     that haven't been migrated to the new primarykey tag system
//
// # SPINE Protocol Context
//
// Remote devices often send "structure" messages containing only key fields
// to establish data schemas before sending actual data. These must be filtered
// to prevent:
//   - Duplicate entries with empty data
//   - Corruption of existing complete entries
//   - Protocol violations in multi-vendor scenarios
//
// # Parameters
//
//   - item: struct instance to analyze
//
// # Returns
//
//   - bool: true if entry contains only primary key data, false if it has
//     additional meaningful fields
//
// # Example
//
//	type MeasurementData struct {
//	    MeasurementId *uint   `eebus:"key,primarykey"`
//	    ValueType     *string `eebus:"key"`
//	    Value         *int
//	}
//
//	keyOnly := MeasurementData{MeasurementId: util.Ptr(uint(1))}
//	hasPrimaryKeyOnly(keyOnly) // Returns: true (should be filtered)
//
//	withData := MeasurementData{
//	    MeasurementId: util.Ptr(uint(1)),
//	    Value:         util.Ptr(100),
//	}
//	hasPrimaryKeyOnly(withData) // Returns: false (should be processed)
func hasPrimaryKeyOnly(item any) bool {
	primaryKeys := fieldNamesWithEEBusTag(EEBusTagPrimaryKey, item)
	if len(primaryKeys) == 0 {
		// No primarykey tags found - handle backward compatibility
		// This supports legacy data structures that haven't been migrated
		// to the new primarykey tag system
		keys := fieldNamesWithEEBusTag(EEBusTagKey, item)
		if len(keys) == 1 {
			// Single key type - use simplified legacy detection
			return hasOnlySingleKey(item, keys[0])
		}
		// Composite keys without primarykey tags are not filtered
		// (safer to process than risk data loss)
		return false
	}

	// Type has primarykey tags - use enhanced detection algorithm
	return hasPrimaryKeyOnlyNew(item, primaryKeys)
}

// hasOnlySingleKey checks if only the specified key field has a value in a single-key struct.
//
// This function provides backward compatibility for data types that use a single
// key field without primarykey tags. It ensures only the key field contains data
// and all other fields are at their zero values.
//
// # Backward Compatibility
//
// This function supports legacy data structures that haven't been migrated
// to the new primarykey tag system, maintaining existing behavior while
// allowing gradual migration to the enhanced composite key system.
//
// # Field Value Detection
//
// The function handles different field types appropriately:
//   - Pointers: checks for non-nil values
//   - Slices/Maps: checks for non-nil and non-empty
//   - Strings: checks for non-empty values
//   - Other types: checks for non-zero values
//
// # Parameters
//
//   - item: struct instance to analyze
//   - keyField: name of the single key field to check
//
// # Returns
//
//   - bool: true if only the key field has a value, false otherwise
//
// # Example
//
//	type SimpleData struct {
//	    ID    *uint   `eebus:"key"`
//	    Value *int
//	    Name  *string
//	}
//
//	keyOnly := SimpleData{ID: util.Ptr(uint(1))}
//	hasOnlySingleKey(keyOnly, "ID") // Returns: true
//
//	withData := SimpleData{ID: util.Ptr(uint(1)), Value: util.Ptr(42)}
//	hasOnlySingleKey(withData, "ID") // Returns: false
func hasOnlySingleKey(item any, keyField string) bool {
	v := reflect.ValueOf(item)
	t := reflect.TypeOf(item)

	if v.Kind() != reflect.Struct {
		return false
	}

	hasKey := false

	// Examine each field to determine if it has a meaningful value
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name

		// Determine if field contains data based on its type
		hasValue := false
		switch field.Kind() {
		case reflect.Ptr:
			// Pointer fields: check for non-nil
			hasValue = !field.IsNil()
		case reflect.Slice, reflect.Map:
			// Collection fields: check for non-nil and non-empty
			hasValue = !field.IsNil() && field.Len() > 0
		case reflect.String:
			// String fields: check for non-empty
			hasValue = field.String() != ""
		default:
			// Other types: check for non-zero values
			hasValue = !field.IsZero()
		}

		if hasValue {
			if fieldName == keyField {
				// Found the key field with a value
				hasKey = true
			} else {
				// Non-key field has value - not key-only
				return false
			}
		}
	}

	return hasKey
}

// hasPrimaryKeyOnlyNew checks if only primary key fields have values using the enhanced tag system.
//
// This function implements the new approach for composite key structures that use
// `eebus:"primarykey"` tags to distinguish primary identifiers from sub-identifiers.
// It provides more precise control over what constitutes "key-only" data in
// complex multi-field key scenarios.
//
// # Enhanced Primary Key Detection
//
// Unlike the legacy single-key approach, this function:
//   - Supports multiple primary key fields in composite structures
//   - Distinguishes primary keys from sub-identifiers
//   - Enables fine-grained filtering based on identifier hierarchy
//   - Provides better compatibility with complex SPINE data models
//
// # Algorithm
//
//  1. Iterate through all struct fields
//  2. Check if each field has a non-zero value
//  3. Classify fields as primary key or other data
//  4. Return true only if primary keys exist but no other data exists
//
// # Parameters
//
//   - item: struct instance to analyze
//   - primaryKeyFields: slice of field names tagged as primary keys
//
// # Returns
//
//   - bool: true if only primary key fields have values, false if any
//     non-primary-key fields contain data
//
// # Example
//
//	type CompositeData struct {
//	    DeviceID *uint   `eebus:"key,primarykey"`
//	    EntityID *uint   `eebus:"key,primarykey"`
//	    SubType  *string `eebus:"key"`
//	    Value    *int
//	}
//
//	primaryOnly := CompositeData{
//	    DeviceID: util.Ptr(uint(1)),
//	    EntityID: util.Ptr(uint(2)),
//	}
//	hasPrimaryKeyOnlyNew(primaryOnly, []string{"DeviceID", "EntityID"}) // Returns: true
//
//	withSubKey := CompositeData{
//	    DeviceID: util.Ptr(uint(1)),
//	    EntityID: util.Ptr(uint(2)),
//	    SubType:  util.Ptr("measurement"),
//	}
//	hasPrimaryKeyOnlyNew(withSubKey, []string{"DeviceID", "EntityID"}) // Returns: false
func hasPrimaryKeyOnlyNew(item any, primaryKeyFields []string) bool {
	v := reflect.ValueOf(item)
	t := reflect.TypeOf(item)

	if v.Kind() != reflect.Struct {
		return false
	}

	hasPrimaryKey := false
	hasOtherData := false

	// Analyze each field to categorize it as primary key or other data
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name
		// Check if this field is marked as a primary key
		isPrimaryKey := slices.Contains(primaryKeyFields, fieldName)

		// Determine if field contains meaningful data
		hasValue := false
		switch field.Kind() {
		case reflect.Ptr:
			// Pointer types: non-nil indicates value
			hasValue = !field.IsNil()
		case reflect.Slice, reflect.Map:
			// Collections: non-nil and non-empty indicates value
			hasValue = !field.IsNil() && field.Len() > 0
		case reflect.String:
			// Strings: non-empty indicates value
			hasValue = field.String() != ""
		default:
			// Other types: non-zero indicates value
			hasValue = !field.IsZero()
		}

		// Skip fields without values
		if !hasValue {
			continue
		}

		// Categorize fields with values
		if isPrimaryKey {
			hasPrimaryKey = true
		} else {
			hasOtherData = true
		}
	}

	return hasPrimaryKey && !hasOtherData
}

// filterPrimaryKeyOnlyEntries removes entries containing only key fields to prevent duplicates.
//
// This is the core anti-duplication mechanism that filters out "structural" entries
// commonly sent by remote SPINE devices. These entries contain only identification
// fields without meaningful data and would create duplicate or incomplete records
// if allowed to merge with existing data.
//
// # SPINE Protocol Context
//
// Remote devices often send messages in two phases:
//  1. Structure phase: entries with only key fields (filtered by this function)
//  2. Data phase: entries with keys + actual data (processed normally)
//
// This separation is common in SPINE implementations and filtering the structure
// phase prevents data corruption and duplicate entry creation.
//
// # Filtering Strategy
//
// The function:
//   - Identifies entries with only key/primary key fields
//   - Logs filtered entries for debugging
//   - Returns only entries containing meaningful data
//   - Maintains original order for non-filtered entries
//
// # Performance Considerations
//
// For large datasets, this function:
//   - Processes entries in single pass
//   - Only allocates new slice if filtering occurs
//   - Provides detailed logging for troubleshooting
//
// # Parameters
//
//   - data: slice of entries to filter
//
// # Returns
//
//   - []T: slice with key-only entries removed (nil if all entries filtered)
//
// # Example
//
//	input := []MeasurementData{
//	    {MeasurementId: util.Ptr(1)}, // Key-only - filtered
//	    {MeasurementId: util.Ptr(2), ValueType: util.Ptr("power"), Value: util.Ptr(100)}, // Data - kept
//	    {MeasurementId: util.Ptr(3)}, // Key-only - filtered
//	}
//	result := filterPrimaryKeyOnlyEntries(input)
//	// Returns: [{MeasurementId: 2, ValueType: "power", Value: 100}]
func filterPrimaryKeyOnlyEntries[T any](data []T) []T {
	if len(data) == 0 {
		return data
	}

	var result []T
	var filteredCount int

	// Process each entry to determine if it should be filtered
	for _, item := range data {
		if hasPrimaryKeyOnly(item) {
			// Entry contains only key data - filter it out
			filteredCount++
			// Provide detailed logging for debugging
			primaryKeys := fieldNamesWithEEBusTag(EEBusTagPrimaryKey, item)
			if len(primaryKeys) == 0 {
				// Legacy single key type
				keys := fieldNamesWithEEBusTag(EEBusTagKey, item)
				logging.Log().Debugf("Ignoring incoming %T with only key field %v (preventing duplicate entry): %+v",
					item, keys, item)
			} else {
				// Enhanced composite key type
				logging.Log().Debugf("Ignoring incoming %T with only primary key fields %v (preventing duplicate entry): %+v",
					item, primaryKeys, item)
			}
		} else {
			// Entry contains meaningful data - keep it
			result = append(result, item)
		}
	}

	if filteredCount > 0 {
		logging.Log().Debugf("Ignored %d incoming %T entries with only key fields to prevent duplicate/low-quality data",
			filteredCount, data)
	}

	return result
}

// SortData sorts slice entries by their EEBus key fields for consistent ordering.
//
// This function provides deterministic ordering of SPINE data by sorting entries
// based on their key fields (identified by `eebus:"key"` tags). Consistent
// ordering is important for reproducible results and easier debugging.
//
// # Sorting Algorithm
//
//   - Identifies all fields tagged with `eebus:"key"`
//   - Sorts entries by comparing key field values in order
//   - Only sorts entries with valid, non-nil uint pointer key fields
//   - Preserves original order for entries that cannot be compared
//
// # Key Field Requirements
//
// For sorting to work, key fields must be:
//   - Pointer types (*uint recommended)
//   - Non-nil values
//   - Comparable types (currently supports uint)
//
// # Performance
//
//   - Uses Go's standard sort.Slice for O(n log n) performance
//   - Handles edge cases gracefully (empty slices, missing keys)
//   - Early returns for unsortable data
//
// # Parameters
//
//   - data: slice of entries to sort
//
// # Returns
//
//   - []T: sorted slice (same slice, modified in place)
//
// # Example
//
//	data := []MeasurementData{
//	    {MeasurementId: util.Ptr(uint(3)), Value: util.Ptr(300)},
//	    {MeasurementId: util.Ptr(uint(1)), Value: util.Ptr(100)},
//	    {MeasurementId: util.Ptr(uint(2)), Value: util.Ptr(200)},
//	}
//	SortData(data)
//	// Result: entries ordered by MeasurementId: 1, 2, 3
func SortData[T any](data []T) []T {
	if len(data) == 0 {
		return data
	}

	keys := fieldNamesWithEEBusTag(EEBusTagKey, data[0])

	if len(keys) == 0 {
		return data
	}

	sort.Slice(data, func(i, j int) bool {
		item1 := data[i]
		item2 := data[j]

		item1V := reflect.ValueOf(item1)
		item2V := reflect.ValueOf(item2)

		// if the fields don't match, don't do anything
		if item1V.NumField() != item2V.NumField() {
			return false
		}

		for _, fieldName := range keys {
			f1 := item1V.FieldByName(fieldName)
			f2 := item2V.FieldByName(fieldName)
			if f1.Type().Kind() != reflect.Ptr || f2.Type().Kind() != reflect.Ptr {
				return false
			}

			if f1.IsNil() || f2.IsNil() || !f1.IsValid() || !f2.IsValid() {
				return false
			}

			if f1.Elem().Kind() != reflect.Uint || f2.Elem().Kind() != reflect.Uint {
				return false
			}

			value1 := f1.Elem().Uint()
			value2 := f2.Elem().Uint()

			if value1 != value2 {
				return value1 < value2
			}
		}

		return false
	})

	return data
}

// copyToSelectedData applies partial updates to entries matching filter selectors.
//
// This function implements SPINE's partial update semantics by finding entries
// that match the provided filter selectors and copying non-nil fields from
// the new data to those matching entries. This enables precise field-level
// updates without affecting other entries or fields.
//
// # Partial Update Process
//
//  1. Iterate through existing entries
//  2. Check each entry against filter selectors
//  3. For matching entries, copy non-nil fields from newData
//  4. Respect write permissions if remoteWrite is true
//
// # Write Permission Enforcement
//
// When remoteWrite is true, the function checks "writecheck" tagged fields
// to determine if remote modifications are allowed. Operations fail if
// write permissions are denied.
//
// # Parameters
//
//   - remoteWrite: true if data originates from remote SPINE device
//   - existingData: current data set to update
//   - filterData: filter containing selectors for matching entries
//   - newData: data to copy to matching entries
//
// # Returns
//
//   - []T: updated data set with selective modifications
//   - bool: true if all operations succeeded, false if write permissions denied
func copyToSelectedData[T any](remoteWrite bool, existingData []T, filterData *FilterData, newData *T) ([]T, bool) {
	if filterData.Selector == nil {
		return existingData, true
	}

	success := true

	for i := range existingData {
		if filterData.SelectorMatch(util.Ptr(existingData[i])) {
			writeAllowed := writeAllowed(existingData[i])
			if !writeAllowed && remoteWrite {
				success = false
				continue
			}

			CopyNonNilDataFromItemToItem(newData, &existingData[i])
			break
		}
	}
	return existingData, success
}

// copyToAllData applies updates to all existing entries (broadcast semantics).
//
// This function implements SPINE's "update all" semantics used when incoming
// data lacks complete key identifiers. It copies non-nil fields from the
// new data to every existing entry, effectively broadcasting the update.
//
// # Broadcast Update Semantics
//
// According to SPINE Table 7, when entries have incomplete identifiers,
// the update should be applied to all existing entries rather than
// creating new entries or failing the operation.
//
// # Use Cases
//
//   - Global configuration updates affecting all entries
//   - Status changes that apply to entire collections
//   - Broadcast notifications from remote devices
//
// # Write Permission Enforcement
//
// When remoteWrite is true, respects "writecheck" tagged field permissions.
// Individual entry updates may fail while others succeed.
//
// # Parameters
//
//   - remoteWrite: true if data originates from remote SPINE device
//   - existingData: current data set to update
//   - newData: data to copy to all existing entries
//
// # Returns
//
//   - []T: updated data set with broadcast modifications
//   - bool: true if all operations succeeded, false if any write permissions denied
func copyToAllData[T any](remoteWrite bool, existingData []T, newData *T) ([]T, bool) {
	success := true

	for i := range existingData {
		writeAllowed := writeAllowed(existingData[i])
		if !writeAllowed && remoteWrite {
			success = false
			continue
		}

		CopyNonNilDataFromItemToItem(newData, &existingData[i])
	}

	return existingData, success
}

// deleteFilteredData executes selective deletion operations based on filter criteria.
//
// This function implements SPINE's delete filter semantics, supporting both
// entry-level deletion (removing entire entries) and field-level deletion
// (removing specific fields from entries). The deletion strategy depends
// on the filter configuration.
//
// # Deletion Strategies
//
//  1. Selector + Elements: Remove specified fields from matching entries
//  2. Selector only: Remove entire entries that match selectors
//  3. Elements only: Remove specified fields from all entries
//
// # Filter Processing
//
// The function supports complex deletion patterns:
//   - Conditional deletion based on entry content
//   - Selective field removal preserving entry structure
//   - Bulk operations across multiple entries
//
// # Write Permission Enforcement
//
// When remoteWrite is true, deletion operations respect "writecheck"
// tagged field permissions. Unauthorized deletions are skipped.
//
// # Parameters
//
//   - remoteWrite: true if data originates from remote SPINE device
//   - existingData: current data set to process
//   - filterData: filter specifying deletion criteria
//
// # Returns
//
//   - []T: modified data set after deletions
//   - bool: true if all operations succeeded, false if write permissions denied
func deleteFilteredData[T any](remoteWrite bool, existingData []T, filterData *FilterData) ([]T, bool) {
	success := true

	if filterData.Elements == nil && filterData.Selector == nil {
		return existingData, true
	}

	var result []T
	for i := range existingData {
		writeAllowed := writeAllowed(existingData[i])
		if !writeAllowed && remoteWrite {
			success = false
			continue
		}

		if filterData.Selector != nil && filterData.Elements != nil {
			// selector and elements filter

			// remove the fields defined in element if the item matches
			if filterData.SelectorMatch(util.Ptr(existingData[i])) {
				RemoveElementFromItem(&existingData[i], filterData.Elements)
				result = append(result, existingData[i])
			} else {
				result = append(result, existingData[i])
			}
		} else if filterData.Selector != nil {
			// only selector filter

			// remove the whole item if the item matches
			if !filterData.SelectorMatch(util.Ptr(existingData[i])) {
				result = append(result, existingData[i])
			}
		} else {
			// only elements filter

			// remove the fields defined in element
			RemoveElementFromItem(&existingData[i], filterData.Elements)
			result = append(result, existingData[i])
		}
	}

	return result, success
}

// isFieldValueNil checks if a field contains a nil value using type-safe reflection.
//
// This utility function safely determines if a field value is nil, handling
// different types appropriately. It's used throughout the update system for
// nil-checking during field processing and value detection.
//
// # Supported Types
//
//   - Pointers: checks if pointer is nil
//   - Maps: checks if map is nil
//   - Arrays: checks if array is nil
//   - Channels: checks if channel is nil
//   - Slices: checks if slice is nil
//   - Other types: always returns false (cannot be nil)
//
// # Parameters
//
//   - field: the field value to check
//
// # Returns
//
//   - bool: true if field is nil, false otherwise
func isFieldValueNil(field interface{}) bool {
	if field == nil {
		return true
	}

	switch reflect.TypeOf(field).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(field).IsNil()
	default:
		return false
	}
}

// nonNilElementNames extracts field names from an element structure that contain non-nil values.
//
// This helper function is used in element-based deletion operations to identify
// which fields should be removed from target items. It examines an element
// template structure and returns the names of fields that have non-nil values.
//
// # Parameters
//
//   - element: pointer to element structure to examine
//
// # Returns
//
//   - []string: slice of field names with non-nil values
func nonNilElementNames(element any) []string {
	var result []string

	v := reflect.ValueOf(element).Elem()
	t := reflect.TypeOf(element).Elem()
	// Examine each field in the element structure
	for i := 0; i < v.NumField(); i++ {
		// Check if field contains a non-nil value
		isNil := isFieldValueNil(v.Field(i).Interface())
		if !isNil {
			// Non-nil field indicates it should be removed from target
			name := t.Field(i).Name
			result = append(result, name)
		}
	}

	return result
}

// isStringValueInSlice checks if a string value exists in a slice of strings.
//
// This utility function provides simple membership testing for string slices.
// It's used throughout the update system for field name matching and validation.
//
// # Parameters
//
//   - value: string value to search for
//   - list: slice of strings to search in
//
// # Returns
//
//   - bool: true if value is found in list, false otherwise
func isStringValueInSlice(value string, list []string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

// RemoveElementFromItem removes fields from an item based on a template element structure.
//
// This function implements SPINE's element-based deletion semantics by examining
// a template element structure and setting corresponding fields in the target
// item to their zero values. It's used for partial deletions in filter operations.
//
// # Element-Based Deletion
//
// The SPINE protocol supports selective field deletion using "element" structures
// that specify which fields to remove. Non-nil fields in the element template
// indicate which fields should be deleted from the target item.
//
// # Type Safety
//
//   - Uses reflection to match field names between element and item
//   - Verifies field count compatibility before processing
//   - Safely handles field access and modification
//
// # Parameters
//
//   - item: pointer to the target item to modify
//   - element: template structure indicating which fields to remove
//
// # Example
//
//	item := &MeasurementData{
//	    MeasurementId: util.Ptr(uint(1)),
//	    ValueType:     util.Ptr("power"),
//	    Value:         util.Ptr(100),
//	}
//	elements := &MeasurementDataElements{
//	    Value: &ScaledNumberElements{}, // Indicates Value field should be removed
//	}
//	RemoveElementFromItem(item, elements)
//	// Result: item.Value is now nil, other fields unchanged
func RemoveElementFromItem[T any, E any](item *T, element E) {
	fieldNamesToBeRemoved := nonNilElementNames(element)

	eV := reflect.ValueOf(element).Elem()
	eT := reflect.TypeOf(element).Elem()
	iV := reflect.ValueOf(item).Elem()

	// if the fields don't match, don't do anything
	if eV.NumField() != iV.NumField() {
		return
	}

	for i := 0; i < eV.NumField(); i++ {
		fieldName := eT.Field(i).Name
		if isStringValueInSlice(fieldName, fieldNamesToBeRemoved) {
			f := iV.FieldByName(fieldName)
			if !f.IsValid() {
				continue
			}
			if !f.CanSet() {
				continue
			}

			f.Set(reflect.Zero(f.Type()))
		}
	}
}

// CopyNonNilDataFromItemToItem copies non-nil fields from source to destination.
//
// This function implements SPINE's merge semantics by copying only fields that
// contain actual data (non-nil values) from the source to the destination.
// This preserves existing data in the destination while updating only the
// fields provided in the source.
//
// # Merge Semantics
//
//   - Only copies non-nil fields from source
//   - Preserves existing data in destination for fields not in source
//   - Handles type safety through reflection
//   - Supports all pointer-based field types
//
// # Field Processing
//
//   - Iterates through all fields in source struct
//   - Checks if source field is non-nil
//   - Copies non-nil fields to corresponding destination fields
//   - Skips nil fields to preserve destination data
//
// # Safety Checks
//
//   - Validates both source and destination are non-nil
//   - Ensures field count compatibility
//   - Verifies field accessibility and mutability
//
// # Parameters
//
//   - source: pointer to source item (provides new data)
//   - destination: pointer to destination item (receives updates)
//
// # Example
//
//	source := &MeasurementData{
//	    MeasurementId: util.Ptr(uint(1)),
//	    Value:         util.Ptr(200), // New value
//	    // ValueType is nil - will not overwrite destination
//	}
//	destination := &MeasurementData{
//	    MeasurementId: util.Ptr(uint(1)),
//	    ValueType:     util.Ptr("power"), // Preserved
//	    Value:         util.Ptr(100),     // Will be updated to 200
//	}
//	CopyNonNilDataFromItemToItem(source, destination)
//	// Result: destination.Value = 200, destination.ValueType = "power" (preserved)
func CopyNonNilDataFromItemToItem[T any](source *T, destination *T) {
	if source == nil || destination == nil {
		return
	}

	sV := reflect.ValueOf(source).Elem()
	sT := reflect.TypeOf(source).Elem()
	dV := reflect.ValueOf(destination).Elem()

	// if the fields don't match, don't do anything
	if sV.NumField() != dV.NumField() {
		return
	}

	// Copy each non-nil field from source to destination
	for i := 0; i < sV.NumField(); i++ {
		value := sV.Field(i)
		// Skip nil fields to preserve destination data
		if value.IsNil() {
			continue
		}

		// Find corresponding field in destination
		fieldName := sT.Field(i).Name
		f := dV.FieldByName(fieldName)

		// Validate field accessibility
		if !f.IsValid() {
			continue
		}
		if !f.CanSet() {
			continue
		}

		// Copy source field value to destination
		f.Set(value)
	}
}
