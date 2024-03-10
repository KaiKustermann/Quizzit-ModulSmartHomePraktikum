// Package configflag provides all possible flags of the application.
// Also utilizes an easy access to the given flags at any time
// and means to patch the [QuizzitConfig] with flag values, if set.
package configflag

import (
	"flag"
	"time"

	log "github.com/sirupsen/logrus"
)

// flags is our local instance holding our settings
var flags = AppFlags{}

// InitFlags defines the application's flags, parses them and reads them into our AppFlags struct.
func InitFlags() {
	configPath := flag.String("config", "./config.yaml", "Relative path to the system config file")
	userConfigPath := flag.String("user-config", "./user-config.yaml", "Relative path to the user config file")
	catalogPath := flag.String("catalog", "", "Relative path to the catalog file. Leave empty for default")
	questionsPath := flag.String("questions", "", "Relative path to the questions file. Leave empty for default")
	httpPort := flag.Int("http-port", 0, "Port for the HTTP Server. Put '0' for default")
	dieEnabled := flag.String("die-enabled", "", "Disable any hybrid-die functionality. Use 'yes' and 'no'. Leave empty for default")
	// TODO: could maybe do duration as flag type
	dieSearchTimeout := flag.String("die-search-timeout", "", "Maximum time to wait for the hybrid die, see time.ParseDuration. Leave empty for default")
	logLevel := flag.String("log-level", "", "Granularity of log output, see logrus.ParseLevel. Leave empty for default")
	logFileLevel := flag.String("log-file-level", "", "Granularity of log output for logfile, see logrus.ParseLevel. Leave empty for default")
	flag.Parse()
	flags.ConfigPath = *configPath
	flags.UserConfigPath = *userConfigPath

	if *catalogPath != "" {
		flags.CatalogPath = catalogPath
	}

	if *questionsPath != "" {
		flags.QuestionsPath = questionsPath
	}

	if *httpPort != 0 {
		flags.HttpPort = httpPort
	}

	if *dieEnabled != "" {
		flags.DieEnabled = dieEnabled
	}

	if *dieSearchTimeout != "" {
		dur, err := time.ParseDuration(*dieSearchTimeout)
		if err == nil {
			flags.HybridDieSearchTimeout = &dur
		} else {
			log.Warnf("Failed parsing Hybrid Die Search Timeout '%s' %e", *dieSearchTimeout, err)
		}
	}

	if *logLevel != "" {
		lvl, err := log.ParseLevel(*logLevel)
		if err == nil {
			flags.LogLevel = &lvl
		} else {
			log.Warnf("Failed parsing Log Level '%s' %e", *logLevel, err)
		}
	}

	if *logFileLevel != "" {
		lvl, err := log.ParseLevel(*logFileLevel)
		if err == nil {
			flags.LogFileLevel = &lvl
		} else {
			log.Warnf("Failed parsing LogFile Level '%s' %e", *logLevel, err)
		}
	}

}

// GetAppFlags returns the flags
func GetAppFlags() AppFlags {
	return flags
}
