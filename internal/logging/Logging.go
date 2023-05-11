package logging

import (
	"fmt"
	"path"
	"runtime"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

func SetUpLogging() {
	// TODO: Set Log level via process env or args!
	log.SetLevel(log.DebugLevel)
	var formatter = nested.Formatter{
		// HideKeys:        true,
		CallerFirst:     true,
		FieldsOrder:     []string{"time", "component", "category"},
		TimestampFormat: time.RFC3339,
		CustomCallerFormatter: func(f *runtime.Frame) string {
			filename := path.Base(f.File)
			return fmt.Sprintf(" %s:%d::%s()", filename, f.Line, f.Function)
		},
	}
	log.SetFormatter(&formatter)
	log.SetReportCaller(true)
}

// Enhance Log with metadata from envelope
func EnvelopeLog(envelope dto.WebsocketMessagePublish) *log.Entry {
	return log.WithFields(log.Fields{
		// "body":          envelope.Body,
		"correlationId": envelope.CorrelationId,
		"messageType":   envelope.MessageType,
	})
}
