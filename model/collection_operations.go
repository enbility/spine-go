package model

import (
	"fmt"
	"reflect"
	"slices"
)

// creates an hash key by using fields that have eebus tag "key"
func hashKey(data any) string {
	result := ""

	keys := fieldNamesWithEEBusTag(EEBusTagKey, data)

	if len(keys) == 0 {
		return result
	}

	v := reflect.ValueOf(data)

	for _, fieldName := range keys {
		f := v.FieldByName(fieldName)

		if f.IsNil() || !f.IsValid() {
			return result
		}

		switch f.Elem().Kind() {
		case reflect.String:
			value := f.Elem().String()

			if len(result) > 0 {
				result = fmt.Sprintf("%s|", result)
			}
			result = fmt.Sprintf("%s%s", result, value)

		case reflect.Uint:
			value := f.Elem().Uint()

			if len(result) > 0 {
				result = fmt.Sprintf("%s|", result)
			}
			result = fmt.Sprintf("%s%d", result, value)

		case reflect.Struct:
			value := f.Type()
			fmt.Println(value)

			if !f.CanInterface() {
				return result
			}

			c, ok := f.Interface().(UpdateHelper)
			if !ok {
				return result
			}

			if len(result) > 0 {
				result = fmt.Sprintf("%s|", result)
			}
			result = fmt.Sprintf("%s%s", result, c.String())

			return result
		default:
			return result
		}
	}

	return result
}

// check the eebus tag if is has a "writecheck" item
// and if so, if the value of that field is true
func writeAllowed(data any) bool {
	fields := fieldNamesWithEEBusTag(EEBusTagWriteCheck, data)
	// only one field in a struct may have this tag
	if len(fields) != 1 {
		return true
	}

	fieldName := fields[0]
	v := reflect.ValueOf(data)
	f := v.FieldByName(fieldName)

	if f.IsNil() || !f.IsValid() {
		return false
	}

	// if this is not a boolean, the tag is wrong which shouldn't happen
	// and we allow overwriting
	if f.Elem().Kind() != reflect.Bool {
		return true
	}

	value := f.Elem().Bool()
	return value
}

// update missing fields in destination with values from source
func updateFields[T any](remoteWrite bool, source T, destination *T) {
	if destination == nil {
		return
	}

	writeCheckFields := fieldNamesWithEEBusTag(EEBusTagWriteCheck, source)

	sV := reflect.ValueOf(source)
	sT := reflect.TypeOf(source)
	dV := reflect.ValueOf(destination).Elem()

	// if the fields don't match, don't do anything
	if sV.Kind() != reflect.Struct || sV.NumField() != dV.NumField() {
		return
	}

	for i := 0; i < sV.NumField(); i++ {
		value := sV.Field(i)
		fieldName := sT.Field(i).Name
		f := dV.FieldByName(fieldName)

		if !f.IsValid() ||
			!f.CanSet() {
			continue
		}

		// on local merge set all nil values
		// on remote writes only set nil values if it is not a "writecheck" tagged field
		if f.IsNil() ||
			(remoteWrite && len(writeCheckFields) > 0 && slices.Contains(writeCheckFields, fieldName)) {
			f.Set(value)
		}
	}
}

// Merges two slices into one. The item in the first slice will be replaced by the one in the second slice
// if the hash key is the same. Items in the second slice which are not in the first will be added.
//
// Parameter remoteWrite defines if this data came on from a remote service, as that is then to
// ignore the "writecheck" tagges fields and should only be allowed to write if the "writecheck" tagged field
// boolean is set to true
//
// returns:
//   - the new data set
//   - true if everything was successful, false if not
func Merge[T any](remoteWrite bool, s1 []T, s2 []T) ([]T, bool) {
	result := []T{}
	success := true

	m2 := ToMap(s2)

	// go through the first slice
	m1 := make(map[string]T, len(s1))
	for _, s1Item := range s1 {
		s1ItemHash := hashKey(s1Item)
		s2Item, exist := m2[s1ItemHash]
		writeAllowed := writeAllowed(s1Item)
		if !writeAllowed && remoteWrite {
			success = false
		}
		// if exists and overwriting is allowed
		if exist && (!remoteWrite || writeAllowed) {
			// add values from s1Item that don't exist in s2Item or shouldn't be
			// set in s2Item
			updateFields(remoteWrite, s1Item, &s2Item)

			// the item in the first slice will be replaced by the one of the second slice
			result = append(result, s2Item)
		} else {
			result = append(result, s1Item)
		}

		m1[s1ItemHash] = s1Item
	}

	// append items which were not in the first slice
	for _, s2Item := range s2 {
		s2ItemHash := hashKey(s2Item)
		_, exist := m1[s2ItemHash]
		if !exist && !remoteWrite {
			// only local updates can append data
			result = append(result, s2Item)
		}
	}

	return result, success
}

func ToMap[T any](s []T) map[string]T {
	result := make(map[string]T, len(s))
	for _, item := range s {
		result[hashKey(item)] = item
	}
	return result
}
