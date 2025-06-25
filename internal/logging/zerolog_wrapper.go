package logging

import (
	"os"

	"github.com/rs/zerolog"
)

type ZerologLogger struct {
	logger zerolog.Logger
}

func NewZerologLogger() *ZerologLogger {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006/01/02 15:04:05",
	}
	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()
	logger.Info().Msg("[ZERLOGGER] ZerologLogger initialized and connected!")
	return &ZerologLogger{
		logger: logger,
	}
}

func (l *ZerologLogger) Infof(format string, args ...interface{}) {
	l.logger.Info().Msgf(format, args...)
}

func (l *ZerologLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warn().Msgf(format, args...)
}

func (l *ZerologLogger) Tracef(format string, args ...interface{}) {
	l.logger.Trace().Msgf(format, args...)
}

func (l *ZerologLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debug().Msgf(format, args...)
}

func (l *ZerologLogger) Panicf(format string, args ...interface{}) {
	l.logger.Panic().Msgf(format, args...)
}
