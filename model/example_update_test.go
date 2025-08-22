package model_test

import (
	"fmt"
	"log"

	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// Example_updateList_measurementDataDuplicatePrevention demonstrates the core problem
// that the primary key filtering system solves: preventing duplicate entries when
// remote SPINE devices send structural messages followed by data messages.
func Example_updateList_measurementDataDuplicatePrevention() {
	// Initial state: Our local measurement data store is empty
	var existingData []model.MeasurementDataType

	// First message from remote device (e.g., during discovery/initialization)
	// This is a "structural" message that only contains identifiers
	structuralMessage := []model.MeasurementDataType{
		{MeasurementId: util.Ptr(model.MeasurementIdType(0))},
		{MeasurementId: util.Ptr(model.MeasurementIdType(4))},
		{MeasurementId: util.Ptr(model.MeasurementIdType(7))},
	}

	// Process the structural message
	result, success := model.UpdateList(false, existingData, structuralMessage, nil, nil, nil)
	if !success {
		log.Fatal("Update failed")
	}

	fmt.Printf("After structural message: %d entries\n", len(result))
	// Output shows 0 entries - all were filtered as primary-key-only

	// Second message from remote device with actual measurement data
	dataMessage := []model.MeasurementDataType{
		{
			MeasurementId: util.Ptr(model.MeasurementIdType(4)),
			ValueType:     util.Ptr(model.MeasurementValueTypeType("power")),
			Value:         model.NewScaledNumberType(1500.0),
			ValueSource:   util.Ptr(model.MeasurementValueSourceType("measuredValue")),
			ValueState:    util.Ptr(model.MeasurementValueStateType("normal")),
		},
	}

	// Process the data message
	result, success = model.UpdateList(false, result, dataMessage, nil, nil, nil)
	if !success {
		log.Fatal("Update failed")
	}

	fmt.Printf("After data message: %d entries\n", len(result))
	fmt.Printf("Entry details: MeasurementId=%d, ValueType=%s, Value=%.0f\n",
		*result[0].MeasurementId,
		*result[0].ValueType,
		result[0].Value.GetValue())

	// Output:
	// After structural message: 0 entries
	// After data message: 1 entries
	// Entry details: MeasurementId=4, ValueType=power, Value=1500
}

// Example_updateList_compositeKeyHandling shows how composite keys work with
// the primarykey tag to distinguish primary identifiers from sub-identifiers.
func Example_updateList_compositeKeyHandling() {
	// Existing measurement with both MeasurementId and ValueType
	existingData := []model.MeasurementDataType{
		{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			ValueType:     util.Ptr(model.MeasurementValueTypeType("power")),
			Value:         model.NewScaledNumberType(100.0),
		},
	}

	// Update attempts from remote device
	updates := []model.MeasurementDataType{
		// This entry only has primary key - will be filtered
		{MeasurementId: util.Ptr(model.MeasurementIdType(1))},
		// This entry has full composite key and data - will be processed
		{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			ValueType:     util.Ptr(model.MeasurementValueTypeType("power")),
			Value:         model.NewScaledNumberType(150.0),
		},
		// New entry with different ValueType - will be added
		{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			ValueType:     util.Ptr(model.MeasurementValueTypeType("voltage")),
			Value:         model.NewScaledNumberType(230.0),
		},
	}

	result, success := model.UpdateList(false, existingData, updates, nil, nil, nil)
	if !success {
		log.Fatal("Update failed")
	}

	fmt.Printf("Total entries: %d\n", len(result))
	for _, entry := range result {
		fmt.Printf("- MeasurementId=%d, ValueType=%s, Value=%.0f\n",
			*entry.MeasurementId,
			*entry.ValueType,
			entry.Value.GetValue())
	}

	// Output:
	// Total entries: 2
	// - MeasurementId=1, ValueType=power, Value=150
	// - MeasurementId=1, ValueType=voltage, Value=230
}

// Example_updateList_remoteWritePermissions demonstrates how the writecheck
// mechanism protects local data from unauthorized remote modifications.
func Example_updateList_remoteWritePermissions() {
	// LoadControl data with write permission control
	// Note: IsLimitChangeable is the writecheck field for LoadControlLimitDataType
	existingData := []model.LoadControlLimitDataType{
		{
			LimitId:           util.Ptr(model.LoadControlLimitIdType(0)),
			IsLimitChangeable: util.Ptr(false), // writecheck=false: denies remote writes
			Value:             model.NewScaledNumberType(16.0),
		},
		{
			LimitId:           util.Ptr(model.LoadControlLimitIdType(1)),
			IsLimitChangeable: util.Ptr(true), // writecheck=true: allows remote writes
			Value:             model.NewScaledNumberType(32.0),
		},
	}

	// Remote device attempts to update both limits
	remoteUpdates := []model.LoadControlLimitDataType{
		{
			LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
			Value:   model.NewScaledNumberType(20.0), // Will be rejected
		},
		{
			LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
			Value:   model.NewScaledNumberType(25.0), // Will be accepted
		},
	}

	// Process with remoteWrite=true to enforce permissions
	// Note: When any update fails due to permissions, success=false
	// but allowed updates are still applied
	result, success := model.UpdateList(true, existingData, remoteUpdates, nil, nil, nil)

	fmt.Printf("Overall success: %v (false because limit 0 was denied)\n", success)
	fmt.Printf("Limit 0 value: %.0f (unchanged)\n", result[0].Value.GetValue())
	fmt.Printf("Limit 1 value: %.0f (updated)\n", result[1].Value.GetValue())

	// Output:
	// Overall success: false (false because limit 0 was denied)
	// Limit 0 value: 16 (unchanged)
	// Limit 1 value: 25 (updated)
}

// Example_updateList_partialUpdateWithFilters shows how to use filters
// to selectively update specific entries in a list.
func Example_updateList_partialUpdateWithFilters() {
	// Multiple measurements in the system
	existingData := []model.MeasurementDataType{
		{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			ValueType:     util.Ptr(model.MeasurementValueTypeType("power")),
			Value:         model.NewScaledNumberType(100.0),
			ValueState:    util.Ptr(model.MeasurementValueStateType("normal")),
		},
		{
			MeasurementId: util.Ptr(model.MeasurementIdType(2)),
			ValueType:     util.Ptr(model.MeasurementValueTypeType("voltage")),
			Value:         model.NewScaledNumberType(230.0),
			ValueState:    util.Ptr(model.MeasurementValueStateType("normal")),
		},
	}

	// Create a filter to only update MeasurementId=1
	filterPartial := model.NewFilterTypePartial()
	filterPartial.MeasurementListDataSelectors = &model.MeasurementListDataSelectorsType{
		MeasurementId: util.Ptr(model.MeasurementIdType(1)),
	}

	// Update data - only affects filtered entry
	updateData := []model.MeasurementDataType{
		{
			Value:      model.NewScaledNumberType(150.0),
			ValueState: util.Ptr(model.MeasurementValueStateType("abnormal")),
		},
	}

	result, success := model.UpdateList(false, existingData, updateData, filterPartial, nil, nil)
	if !success {
		log.Fatal("Update failed")
	}

	for _, entry := range result {
		fmt.Printf("MeasurementId=%d: Value=%.0f, State=%s\n",
			*entry.MeasurementId,
			entry.Value.GetValue(),
			*entry.ValueState)
	}

	// Output:
	// MeasurementId=1: Value=150, State=abnormal
	// MeasurementId=2: Value=230, State=normal
}

// Example_updateList_broadcastUpdate demonstrates the "update all" semantics
// when incoming data lacks complete identifiers.
func Example_updateList_broadcastUpdate() {
	// Multiple LoadControl limits
	existingData := []model.LoadControlLimitDataType{
		{
			LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
			Value:   model.NewScaledNumberType(16.0),
		},
		{
			LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
			Value:   model.NewScaledNumberType(32.0),
		},
		{
			LimitId: util.Ptr(model.LoadControlLimitIdType(2)),
			Value:   model.NewScaledNumberType(25.0),
		},
	}

	// Update without identifiers - applies to all entries
	broadcastUpdate := []model.LoadControlLimitDataType{
		{
			// No LimitId specified - triggers "update all"
			IsLimitChangeable: util.Ptr(true),
			IsLimitActive:     util.Ptr(true),
		},
	}

	result, success := model.UpdateList(false, existingData, broadcastUpdate, nil, nil, nil)
	if !success {
		log.Fatal("Update failed")
	}

	fmt.Println("All limits updated with broadcast values:")
	for _, limit := range result {
		fmt.Printf("LimitId=%d: Changeable=%v, Active=%v\n",
			*limit.LimitId,
			*limit.IsLimitChangeable,
			*limit.IsLimitActive)
	}

	// Output:
	// All limits updated with broadcast values:
	// LimitId=0: Changeable=true, Active=true
	// LimitId=1: Changeable=true, Active=true
	// LimitId=2: Changeable=true, Active=true
}

// Example_updateList_errorHandling shows proper error handling patterns
// for update operations.
func Example_updateList_errorHandling() {
	existingData := []model.LoadControlLimitDataType{
		{
			LimitId:           util.Ptr(model.LoadControlLimitIdType(0)),
			IsLimitChangeable: util.Ptr(false), // writecheck field - denies remote writes
			Value:             model.NewScaledNumberType(16.0),
		},
	}

	// Remote update attempt
	remoteUpdate := []model.LoadControlLimitDataType{
		{
			LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
			Value:   model.NewScaledNumberType(20.0),
		},
	}

	// Attempt remote update
	result, success := model.UpdateList(true, existingData, remoteUpdate, nil, nil, nil)

	if !success {
		// In production, log the failure with context
		fmt.Println("Update failed: Remote write permission denied")
		fmt.Printf("Attempted to update LimitId=%d from %.0f to %.0f\n",
			*remoteUpdate[0].LimitId,
			existingData[0].Value.GetValue(),
			remoteUpdate[0].Value.GetValue())

		// Take appropriate action based on your use case:
		// - Send error response to remote device
		// - Log security event
		// - Retry with different permissions
		// - Alert system administrator
	}

	// Data remains unchanged
	fmt.Printf("Current value: %.0f\n", result[0].Value.GetValue())

	// Output:
	// Update failed: Remote write permission denied
	// Attempted to update LimitId=0 from 16 to 20
	// Current value: 16
}

// Example_updateList_deleteFilterUsage demonstrates how to use delete filters
// to remove specific entries or fields from the data.
func Example_updateList_deleteFilterUsage() {
	// Initial measurement data
	existingData := []model.MeasurementDataType{
		{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			ValueType:     util.Ptr(model.MeasurementValueTypeType("power")),
			Value:         model.NewScaledNumberType(100.0),
			ValueState:    util.Ptr(model.MeasurementValueStateType("normal")),
		},
		{
			MeasurementId: util.Ptr(model.MeasurementIdType(2)),
			ValueType:     util.Ptr(model.MeasurementValueTypeType("voltage")),
			Value:         model.NewScaledNumberType(230.0),
			ValueState:    util.Ptr(model.MeasurementValueStateType("normal")),
		},
	}

	// Create delete filter for MeasurementId=1
	filterDelete := &model.FilterType{
		CmdControl: &model.CmdControlType{Delete: &model.ElementTagType{}},
	}
	filterDelete.MeasurementListDataSelectors = &model.MeasurementListDataSelectorsType{
		MeasurementId: util.Ptr(model.MeasurementIdType(1)),
	}

	// Process deletion
	result, success := model.UpdateList(false, existingData, nil, nil, filterDelete, nil)
	if !success {
		log.Fatal("Delete operation failed")
	}

	fmt.Printf("Remaining entries: %d\n", len(result))
	for _, entry := range result {
		fmt.Printf("MeasurementId=%d, ValueType=%s\n",
			*entry.MeasurementId,
			*entry.ValueType)
	}

	// Output:
	// Remaining entries: 1
	// MeasurementId=2, ValueType=voltage
}

// Example_customUpdater demonstrates implementing the Updater interface
// for custom update logic.
type DeviceMeasurements struct {
	measurements []model.MeasurementDataType
	maxEntries   int
}

func (d *DeviceMeasurements) UpdateList(remoteWrite, persist bool, newList any,
	filterPartial, filterDelete *model.FilterType) (any, bool) {

	// Type assertion for incoming data
	newData, ok := newList.([]model.MeasurementDataType)
	if !ok {
		return d.measurements, false
	}

	// Apply size limit before processing
	if len(d.measurements)+len(newData) > d.maxEntries {
		// In production, implement proper handling:
		// - Remove oldest entries
		// - Reject update
		// - Send notification
		fmt.Printf("Warning: Update would exceed max entries (%d)\n", d.maxEntries)
	}

	// Delegate to standard UpdateList implementation
	result, success := model.UpdateList(remoteWrite, d.measurements, newData,
		filterPartial, filterDelete, nil)

	if success && persist {
		d.measurements = result
		// In production: persist to database/storage
	}

	return result, success
}

func Example_customUpdater() {
	// Create custom updater with constraints
	device := &DeviceMeasurements{
		measurements: []model.MeasurementDataType{},
		maxEntries:   100,
	}

	// Add some measurements
	newData := []model.MeasurementDataType{
		{
			MeasurementId: util.Ptr(model.MeasurementIdType(1)),
			ValueType:     util.Ptr(model.MeasurementValueTypeType("power")),
			Value:         model.NewScaledNumberType(1500.0),
		},
	}

	result, success := device.UpdateList(false, true, newData, nil, nil)
	if !success {
		log.Fatal("Update failed")
	}

	measurements := result.([]model.MeasurementDataType)
	fmt.Printf("Stored measurements: %d\n", len(measurements))

	// Output:
	// Stored measurements: 1
}
