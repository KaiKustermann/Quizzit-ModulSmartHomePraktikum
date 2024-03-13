// Package asyncapiutils provides utility and mapping functions to work with the generated asyncapi DTOs
package asyncapiutils

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// InterfaceToStruct takes 'input' and writes to 'output'.
//
// This is achieved by marshalling and unmarshalling the input into output.
//
// Thus, 'input' AND 'output' must support JSON and use the same keys.
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
