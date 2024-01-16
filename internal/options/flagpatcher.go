package options

import log "github.com/sirupsen/logrus"

func (conf *QuizzitConfig) patchwithFlags() {
	fl := GetAppFlags()
	conf.Log.patchWithLogFlag(fl.LogLevel)
}

func (conf *LogConfig) patchWithLogFlag(flag *log.Level) {
	if flag == nil {
		log.Debug("Log Level is nil, not patching")
		return
	}
	conf.Level = *flag
}
