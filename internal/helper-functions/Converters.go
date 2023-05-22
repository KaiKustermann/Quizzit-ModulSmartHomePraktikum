package helpers

import (
	"encoding/json"
	"regexp"
	"unicode"
	"unicode/utf8"

	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
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

// Like json.marshal, but the JSON keys are lowerCamelCase
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

func QuestionToWebsocketMessageSubscribe(question dto.Question) dto.WebsocketMessageSubscribe {
	msg := dto.WebsocketMessageSubscribe{
		MessageType: "game/question/Question",
		Body:        question,
	}
	return msg
}

func CorrectnessFeedbackToWebsocketMessageSubscribe(correctnessFeedback dto.CorrectnessFeedback) dto.WebsocketMessageSubscribe {
	msg := dto.WebsocketMessageSubscribe{
		MessageType: "game/question/CorrectnessFeedback",
		Body:        correctnessFeedback,
	}
	return msg
}
