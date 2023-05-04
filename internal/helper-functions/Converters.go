package helpers

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// Takes 'input' and writes to 'output'.
//
// 'input' must be json-able.
func InterfaceToStruct(input any, output any) error {
	// TODO: Fix this bad workaround to create the needed DTO
	bytes, err := json.Marshal(input)
	if err != nil {
		log.Debug("Failed marshalling input to JSON")
		return err
	}
	decode_err := json.Unmarshal(bytes, &output)

	if decode_err != nil {
		log.Debug("Failed unmarshalling from JSON")
		return decode_err
	}
	return nil
}
