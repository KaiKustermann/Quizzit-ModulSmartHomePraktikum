// Package jsonutil provides utility functions around JSON and JSONable structs
package jsonutil

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// InterfaceToStruct takes 'input' and writes to 'output'.
//
// This is achieved by marshalling and unmarshalling the input into output.
//
// Thus, 'input' AND 'output' must support JSON and use the same keys.
//
// This is a workaround to transform an interface{} into a struct
func InterfaceToStruct[O any](input interface{}) (out O, err error) {
	bytes, err := json.Marshal(input)
	if err != nil {
		log.Debug("Failed marshalling input to JSON")
		return
	}
	err = json.Unmarshal(bytes, &out)
	if err != nil {
		log.Debug("Failed unmarshalling from JSON")
		return
	}
	return
}
