package options

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

// AppFlags serves as container to hold all flags in one spot.
type AppFlags struct {
	ConfigFile    string
	LogLevel      *log.Level
	QuestionsFile *string
}

// Local instance holding our settings
var flags = AppFlags{}

// Define application's flags, parse them and read them into our options struct.
func InitFlags() {
	configFile := flag.String("config", "./config.yaml", "Relative path to the config file")
	questionsFile := flag.String("questions", "", "Relative path to the questions file")
	logLevel := flag.String("log-level", "", "Granularity of log output, see logrus.ParseLevel")
	flag.Parse()
	flags.ConfigFile = *configFile

	if *questionsFile != "" {
		flags.QuestionsFile = questionsFile
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

func GetAppFlags() AppFlags {
	return flags
}
