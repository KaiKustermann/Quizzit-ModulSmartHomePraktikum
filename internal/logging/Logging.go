package logging

import (
	log "github.com/sirupsen/logrus"
)

// SetUpLogFormat sets a custom formatter for log output
func SetUpLogFormat() {
	log.SetFormatter(&formatter)
	log.SetReportCaller(true)
}

// CreateFileLoggingHook creates a logrus.hook to log to file
func CreateFileLoggingHook() (*FileLoggerHook, error) {
	return NewFileLoggerHook(log.TraceLevel, &formatter, LumberJackConfig{
		Filename:   "logs/quizzit.log",
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     90, //days
	})
}
