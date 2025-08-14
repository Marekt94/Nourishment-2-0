package logging

import (
	"io"
	"strings"

	"github.com/rs/zerolog"
)

type ZerologLogger struct {
	logger zerolog.Logger
}

// GinWriter wraps zerolog logger to provide io.Writer interface for Gin
type GinWriter struct {
	logger zerolog.Logger
}

func (w *GinWriter) Write(p []byte) (n int, err error) {
	// Remove trailing newline if present
	message := string(p)
	if len(message) > 0 && message[len(message)-1] == '\n' {
		message = message[:len(message)-1]
	}

	// Parse Gin log level from message and log at appropriate level
	if strings.Contains(message, "[GIN-debug]") {
		w.logger.Debug().Msg(strings.Replace(message, "[GIN-debug]", "[GIN]", 1))
	} else if strings.Contains(message, "[GIN-trace]") {
		w.logger.Trace().Msg(strings.Replace(message, "[GIN-trace]", "[GIN]", 1))
	} else if strings.Contains(message, "[GIN-warn]") {
		w.logger.Warn().Msg(strings.Replace(message, "[GIN-warn]", "[GIN]", 1))
	} else if strings.Contains(message, "[GIN-warning]") {
		w.logger.Warn().Msg(strings.Replace(message, "[GIN-warning]", "[GIN]", 1))
	} else if strings.Contains(message, "[GIN-error]") {
		w.logger.Error().Msg(strings.Replace(message, "[GIN-error]", "[GIN]", 1))
	} else if strings.Contains(message, "[GIN-fatal]") {
		w.logger.Fatal().Msg(strings.Replace(message, "[GIN-fatal]", "[GIN]", 1))
	} else if strings.Contains(message, "[GIN-panic]") {
		w.logger.Panic().Msg(strings.Replace(message, "[GIN-panic]", "[GIN]", 1))
	} else if strings.Contains(message, "[GIN]") {
		// Standard Gin logs (HTTP requests) - log as Info
		w.logger.Info().Msg(message)
	} else {
		// Default to Info level for other Gin messages
		w.logger.Info().Msg("[GIN] " + message)
	}

	return len(p), nil
}

func NewZerologLogger() *ZerologLogger {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        getLogWriter(),
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

func (l *ZerologLogger) Writer() io.Writer {
	return &GinWriter{logger: l.logger}
}
