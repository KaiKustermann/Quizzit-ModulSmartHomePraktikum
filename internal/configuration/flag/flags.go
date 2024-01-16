// Package configflag provides all possible flags of the application.
// Also utilizes an easy access to the given flags at any time
// and means to patch the [QuizzitConfig] with flag values, if set.
package configflag

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

// AppFlags serves as container to hold all flags in one spot.
type AppFlags struct {
	ConfigFile    string
	QuestionsFile *string
	HttpPort      *int
	LogLevel      *log.Level
}

// flags is our local instance holding our settings
var flags = AppFlags{}

// InitFlags defines the application's flags, parses them and reads them into our AppFlags struct.
func InitFlags() {
	configFile := flag.String("config", "./config.yaml", "Relative path to the config file")
	questionsFile := flag.String("questions", "", "Relative path to the questions file. Leave empty for default")
	httpPort := flag.Int("http-port", 0, "Port for the HTTP Server. Put '0' for default")
	logLevel := flag.String("log-level", "", "Granularity of log output, see logrus.ParseLevel. Leave empty for default")
	flag.Parse()
	flags.ConfigFile = *configFile

	if *questionsFile != "" {
		flags.QuestionsFile = questionsFile
	}

	if *httpPort != 0 {
		flags.HttpPort = httpPort
	}

	if *logLevel != "" {
		lvl, err := log.ParseLevel(*logLevel)
		if err == nil {
			flags.LogLevel = &lvl
		} else {
			log.Warnf("Failed parsing Log Level '%s' %e", *logLevel, err)
		}
	}

}

// GetAppFlags returns the flags
func GetAppFlags() AppFlags {
	return flags
}
