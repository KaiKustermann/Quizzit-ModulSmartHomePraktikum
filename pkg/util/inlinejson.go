package util

import (
	"bytes"
	"encoding/json"
)

// JsonString returns a JSON string representation of 'v'
func JsonString(v any) string {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "")
	enc.Encode(v)
	return buf.String()
}
