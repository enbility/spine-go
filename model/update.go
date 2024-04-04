package model

import (
	"reflect"
	"sort"

	"github.com/enbility/spine-go/util"
)

type Updater interface {
	// returns true if no errors occured
	UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) bool
}

// Generates a new list of function items by applying the rules mentioned in the spec
// (EEBus_SPINE_TS_ProtocolSpecification.pdf; chapter "5.3.4 Restricted function exchange with cmdOptions").
// The given data provider is used the get the current items and the items and the filters in the payload.
//
// returns:
//   - the new data set
//   - true if everything was successful, false if not
func UpdateList[T any](remoteWrite bool, existingData []T, newData []T, filterPartial, filterDelete *FilterType) ([]T, bool) {
	success := true

	// process delete filter (with selectors and elements)
	if filterDelete != nil {
		if filterData, err := filterDelete.Data(); err == nil {
			updatedData, noErrors := deleteFilteredData(remoteWrite, existingData, filterData)
			if noErrors {
				existingData = updatedData
			} else {
				success = false
			}
		}
	}

	// process update filter (with selectors and elements)
	if filterPartial != nil {
		if filterData, err := filterPartial.Data(); err == nil {
			newData, noErrors := copyToSelectedData(remoteWrite, existingData, filterData, &newData[0])
			if !noErrors {
				success = false
			}
			return newData, success
		}
	}

	// check if items have no identifiers
	// Currently all fields marked as key are required
	// TODO: check how to handle if only one identifier is provided
	if len(newData) > 0 && !HasIdentifiers(newData[0]) {
		// no identifiers specified --> copy data to all existing items
		// (see EEBus_SPINE_TS_ProtocolSpecification.pdf, Table 7: Considered cmdOptions combinations for classifier "notify")
		newData, noErrors := copyToAllData(remoteWrite, existingData, &newData[0])
		if !noErrors {
			success = false
		}
		return newData, success
	}

	result, noErrors := Merge(remoteWrite, existingData, newData)
	if !noErrors {
		success = false
	}

	result = SortData(result)

	return result, success
}

// return a list of field names that have the eebus tag
func fieldNamesWithEEBusTag(tag EEBusTag, item any) []string {
	var result []string

	v := reflect.ValueOf(item)
	t := reflect.TypeOf(item)

	if v.Kind() != reflect.Struct {
		return result
	}

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() != reflect.Ptr {
			continue
		}

		sf := v.Type().Field(i)
		eebusTags := EEBusTags(sf)
		_, exists := eebusTags[tag]
		if !exists {
			continue
		}

		fieldName := t.Field(i).Name
		result = append(result, fieldName)
	}

	return result
}

func HasIdentifiers(data any) bool {
	keys := fieldNamesWithEEBusTag(EEBusTagKey, data)

	v := reflect.ValueOf(data)

	for _, fieldName := range keys {
		f := v.FieldByName(fieldName)

		if f.IsNil() || !f.IsValid() {
			return false
		}
	}

	return true
}

// sort slices by fields that have eebus tag "key"
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

// Copy data t elements matching the selected items
//
// Parameter remoteWrite defines if this data came on from a remote service, as that is then to
// ignore the "writecheck" tagges fields and should only be allowed to write if the "writecheck" tagged field
// boolean is set to true
//
// returns:
//   - the new data set
//   - true if everything was successful, false if not
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

// Copy data to all elements
//
// Parameter remoteWrite defines if this data came on from a remote service, as that is then to
// ignore the "writecheck" tagges fields and should only be allowed to write if the "writecheck" tagged field
// boolean is set to true
//
// returns:
//   - the new data set
//   - true if everything was successful, false if not
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

// Execute a partial delete filter
//
// Parameter remoteWrite defines if this data came on from a remote service, as that is then to
// ignore the "writecheck" tagges fields and should only be allowed to write if the "writecheck" tagged field
// boolean is set to true
//
// returns:
//   - the new data set
//   - true if everything was successful, false if not
func deleteFilteredData[T any](remoteWrite bool, existingData []T, filterData *FilterData) ([]T, bool) {
	success := true

	if filterData.Elements == nil && filterData.Selector == nil {
		return existingData, true
	}

	result := []T{}
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

func isFieldValueNil(field interface{}) bool {
	if field == nil {
		return true
	}

	switch reflect.TypeOf(field).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(field).IsNil()
	}
	return false
}

func nonNilElementNames(element any) []string {
	var result []string

	v := reflect.ValueOf(element).Elem()
	t := reflect.TypeOf(element).Elem()
	for i := 0; i < v.NumField(); i++ {
		isNil := isFieldValueNil(v.Field(i).Interface())
		if !isNil {
			name := t.Field(i).Name
			result = append(result, name)
		}
	}

	return result
}

func isStringValueInSlice(value string, list []string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

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

	for i := 0; i < sV.NumField(); i++ {
		value := sV.Field(i)
		if value.IsNil() {
			continue
		}

		fieldName := sT.Field(i).Name
		f := dV.FieldByName(fieldName)

		if !f.IsValid() {
			continue
		}
		if !f.CanSet() {
			continue
		}

		f.Set(value)
	}
}
