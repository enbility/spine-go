package util

import (
	"encoding/json"
)

// quick way to a struct into another
func DeepCopy[A any](source, dest A) {
	byt, _ := json.Marshal(source)
	_ = json.Unmarshal(byt, dest)
}
