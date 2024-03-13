// Package asyncapiutils provides utility and mapping functions to work with the generated asyncapi DTOs
package asyncapiutils

import (
	"encoding/json"
	"errors"
	"regexp"
	"unicode"
	"unicode/utf8"

	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
)

// MarshalToLowerCamelCaseJSON works like json.marshal, but the JSON keys are lowerCamelCase
//
// The behavior is achieved by regex-replaceing the startLetter of all JSON keys.
// For this to work the '"key"' MUST be followed by ':' without whitespace in between!
//
// As seen in https://gist.github.com/piersy/b9934790a8892db1a603820c0c23e4a7
func MarshalToLowerCamelCaseJSON(data any) ([]byte, error) {
	keyMatchRegex := regexp.MustCompile(`\"(\w+)\":`)
	marshalled, err := json.Marshal(data)
	if err != nil {
		return marshalled, err
	}
	converted := keyMatchRegex.ReplaceAllFunc(
		marshalled,
		func(match []byte) []byte {
			// Empty keys are valid JSON, only lowercase if we do not have an
			// empty key.
			if len(match) > 2 {
				// Decode first rune after the double quotes
				r, width := utf8.DecodeRune(match[1:])
				r = unicode.ToLower(r)
				utf8.EncodeRune(match[1:width+1], r)
			}
			return match
		},
	)
	return converted, err
}

// ParseWebsocketMessage unmarshals the given bytes to a [WebsocketMessagePublish]
//
// Also validates 'MessageType' to be present
func ParseWebsocketMessage(payload []byte) (asyncapi.WebsocketMessagePublish, error) {
	var parsedPayload asyncapi.WebsocketMessagePublish
	err := json.Unmarshal(payload, &parsedPayload)
	if err != nil {
		log.Debug("Could not unmarshal JSON", err)
		return parsedPayload, err
	}
	if parsedPayload.MessageType == "" {
		err = errors.New("envelope message type is <empty>")
		return parsedPayload, err
	}
	return parsedPayload, nil
}
