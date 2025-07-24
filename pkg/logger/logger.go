package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap logger
type RootLogger struct {
	logger *zap.Logger
}

// New creates a new logger instance
func NewRootLogger(format string, logLevel string) *RootLogger {
	var config zap.Config
	logFormat := Format(format)

	switch logFormat {
	case Console:
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.Encoding = "console"
		config.EncoderConfig.ConsoleSeparator = " "
	case JSON:
		config = zap.NewProductionConfig()
	default:
		config = zap.NewProductionConfig() // Default to production config for unknown formats
	}

	// Set the log level from config
	config.Level = zap.NewAtomicLevelAt(GetLogLevel(logLevel))
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	return &RootLogger{
		logger: logger,
	}
}

// Debug logs a debug message
func (l *RootLogger) Debug(message string) {
	l.logger.Debug(message)
}

// Info logs an info message
func (l *RootLogger) Info(message string) {
	l.logger.Info(message)
}

// Warn logs a warning message
func (l *RootLogger) Warn(message string) {
	l.logger.Warn(message)
}

// Error logs an error message
func (l *RootLogger) Error(message string, err error) {
	l.logger.Error(message + " " + err.Error())
}

// Fatal logs a fatal message and exits
func (l *RootLogger) Fatal(message string, err error) {
	l.logger.Fatal(message + " " + err.Error())
}
