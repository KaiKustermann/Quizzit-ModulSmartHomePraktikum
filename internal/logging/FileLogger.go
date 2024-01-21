package logging

import (
	"io"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LumberJackConfig struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

type FileLoggerHook struct {
	logWriter io.Writer
	formatter log.Formatter
	level     log.Level
}

func NewFileLoggerHook(level log.Level, formatter log.Formatter, config LumberJackConfig) (*FileLoggerHook, error) {

	hook := FileLoggerHook{
		level:     level,
		formatter: formatter,
	}
	hook.logWriter = &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	return &hook, nil
}

func (hook *FileLoggerHook) SetLevel(level log.Level) {
	hook.level = level
}

func (hook *FileLoggerHook) Levels() []log.Level {
	return log.AllLevels[:hook.level+1]
}

func (hook *FileLoggerHook) Fire(entry *log.Entry) (err error) {
	if entry.Level > hook.level {
		return nil
	}
	b, err := hook.formatter.Format(entry)
	if err != nil {
		return err
	}
	hook.logWriter.Write(b)
	return nil
}
