# Primary Key Tag Guidelines for spine-go

## Overview

The `primarykey` EEBus tag is used to distinguish primary identifiers from sub-identifiers in composite key data types. This improves the precision of list update filtering and aligns with SPINE specification terminology.

## When to Use primarykey Tag

### Rule 1: Composite Key Types
If a data type has **multiple fields** with `eebus:"key"` tags, the PRIMARY IDENTIFIER field must also have the `primarykey` tag:

```go
type ExampleDataType struct {
    PrimaryId *IdType `json:"primaryId,omitempty" eebus:"key,primarykey"`  // PRIMARY
    SubId     *IdType `json:"subId,omitempty" eebus:"key"`               // SUB
}
```

### Rule 2: Single Key Types
Single key types (only one field with `eebus:"key"`) do NOT need the `primarykey` tag:

```go
type SimpleDataType struct {
    Id    *IdType `json:"id,omitempty" eebus:"key"`  // No primarykey needed
    Value *int    `json:"value,omitempty"`
}
```

## How to Identify Primary vs Sub Identifiers

Consult the SPINE specification for each data type. The spec clearly states:
- **PRIMARY IDENTIFIER**: "SHALL be set as PRIMARY IDENTIFIER"
- **SUB IDENTIFIER**: "SHOULD be set" or "MAY be set"
- **FOREIGN IDENTIFIER**: References to other entities

### Examples from SPINE:

1. **MeasurementDataType**
   - `measurementId`: PRIMARY IDENTIFIER (mandatory)
   - `valueType`: SUB IDENTIFIER (SHOULD be set)

2. **SetpointDescriptionDataType**
   - `setpointId`: PRIMARY IDENTIFIER
   - `measurementId`: FOREIGN IDENTIFIER
   - `timeTableId`: FOREIGN IDENTIFIER

## Implementation Examples

### Correct Implementation
```go
// Composite keys with primarykey tag
type MeasurementDataType struct {
    MeasurementId *MeasurementIdType        `json:"measurementId,omitempty" eebus:"key,primarykey"`
    ValueType     *MeasurementValueTypeType `json:"valueType,omitempty" eebus:"key"`
    Value         *ScaledNumberType         `json:"value,omitempty"`
}

// Single key without primarykey tag
type BillDataType struct {
    BillId   *BillIdType   `json:"billId,omitempty" eebus:"key"`
    BillType *BillTypeType `json:"billType,omitempty"`
}
```

### Incorrect Implementation
```go
// WRONG: Composite keys without primarykey tag
type BadExampleType struct {
    Id1 *IdType `eebus:"key"`  // Which is primary?
    Id2 *IdType `eebus:"key"`  // Ambiguous!
}

// WRONG: Single key with unnecessary primarykey tag
type OverEngineeredType struct {
    Id *IdType `eebus:"key,primarykey"`  // Redundant for single keys
}
```

## Impact on Filtering

The `primarykey` tag affects how entries are filtered during list updates:

### With primarykey Tag (Improved Behavior)
```go
// Only this is filtered:
{MeasurementId: 1}                      // ✗ Only primary key

// These are NOT filtered:
{MeasurementId: 1, ValueType: "value"}  // ✓ Has sub-identifier
{MeasurementId: 1, Value: 100}          // ✓ Has data
```

### Without primarykey Tag (Old Behavior)
```go
// Both would be filtered:
{MeasurementId: 1}                      // ✗ Only keys
{MeasurementId: 1, ValueType: "value"}  // ✗ All keys, no data
```

## Checklist for New Data Types

When adding a new data type:

1. ☐ Count the fields with `eebus:"key"` tag
2. ☐ If count > 1, check SPINE spec for PRIMARY IDENTIFIER
3. ☐ Add `primarykey` to the PRIMARY IDENTIFIER field
4. ☐ Test that filtering works correctly
5. ☐ Document any FOREIGN IDENTIFIERs in comments

## Current Status

All composite key types in spine-go have been migrated:
- ✅ MeasurementDataType
- ✅ MeasurementSeriesDataType
- ✅ ElectricalConnectionPermittedValueSetDataType
- ✅ ElectricalConnectionParameterDescriptionDataType
- ✅ ElectricalConnectionCharacteristicDataType
- ✅ SetpointDescriptionDataType

## Testing

Always test composite key types with these scenarios:
1. Entry with only primary key → Should be filtered
2. Entry with primary + sub keys → Should NOT be filtered
3. Entry with primary key + data → Should NOT be filtered

```go
// Example test
func TestYourDataType_PrimaryKey(t *testing.T) {
    // Verify primarykey tag
    primaryKeys := fieldNamesWithEEBusTag(EEBusTagPrimaryKey, YourDataType{})
    assert.Equal(t, []string{"YourPrimaryId"}, primaryKeys)
    
    // Test filtering behavior
    assert.True(t, hasPrimaryKeyOnly(YourDataType{
        YourPrimaryId: util.Ptr(1),
    }))
    assert.False(t, hasPrimaryKeyOnly(YourDataType{
        YourPrimaryId: util.Ptr(1),
        YourSubId:     util.Ptr(2),
    }))
}
```