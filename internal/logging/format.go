package logging

import (
	"fmt"
	"path"
	"runtime"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
)

var formatter = nested.Formatter{
	// HideKeys:        true,
	CallerFirst:     true,
	FieldsOrder:     []string{"time", "component", "category"},
	TimestampFormat: time.RFC3339,
	CustomCallerFormatter: func(f *runtime.Frame) string {
		filename := path.Base(f.File)
		return fmt.Sprintf(" %s:%d", filename, f.Line)
	},
}
