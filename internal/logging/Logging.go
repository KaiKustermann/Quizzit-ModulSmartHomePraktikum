package logging

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

func SetUpLogFormat() {
	var formatter = nested.Formatter{
		// HideKeys:        true,
		CallerFirst:     true,
		FieldsOrder:     []string{"time", "component", "category"},
		TimestampFormat: time.RFC3339,
		CustomCallerFormatter: func(f *runtime.Frame) string {
			filename := path.Base(f.File)
			fun := strings.Replace(f.Function, "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server", "", 1)
			return fmt.Sprintf(" %s:%d::%s()", filename, f.Line, fun)
		},
	}
	log.SetFormatter(&formatter)
	log.SetReportCaller(true)
}
