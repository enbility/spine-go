package util

import (
	"encoding/json"
	"reflect"
)

// quick way to a struct into another
func DeepCopy[A any](source, dest A) {
	byt, _ := json.Marshal(source)
	_ = json.Unmarshal(byt, dest)
}

// checck if a value is nil for any type
func IsNil(data any) bool {
	if data == nil {
		return true
	}
	switch reflect.TypeOf(data).Kind() {
	case reflect.Ptr,
		reflect.Map,
		reflect.Array,
		reflect.Chan,
		reflect.Slice:
		return reflect.ValueOf(data).IsNil()
	}
	return false
}
