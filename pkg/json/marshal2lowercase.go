// Package jsonutil provides utility functions around JSON and JSONable structs
package jsonutil

import (
	"encoding/json"
	"regexp"
	"unicode"
	"unicode/utf8"
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
