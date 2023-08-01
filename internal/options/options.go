package options

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type QuizzitOptions struct {
	HttpPort               string
	LogLevel               log.Level
	ScoredPointsToWin      int
	HybridDieSearchTimeout time.Duration
}

// Local instance holding our settings
var optionsInstance = QuizzitOptions{
	HttpPort:               getServerPort(),
	LogLevel:               getLogLevel(),
	ScoredPointsToWin:      5,
	HybridDieSearchTimeout: 30 * time.Second,
}

/*
Get Quizzit Options
*/
func GetQuizzitOptions() QuizzitOptions {
	return optionsInstance
}

/*
Take port from env 'HTTP_SERVER_PORT'
Default: 8080
*/
func getServerPort() string {
	port := os.Getenv("HTTP_SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

/*
Take log level from env 'LOG_LEVEL'
Default: Info
*/
func getLogLevel() log.Level {
	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = log.InfoLevel
	}
	return logLevel
}
