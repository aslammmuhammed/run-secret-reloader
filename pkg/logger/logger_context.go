package logger

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Constants
const TraceIDKey = "trace_id"

// LogLevel represents logging levels
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Interface defines the logging methods
type Interface interface {
	Debug(ctx context.Context, message string)
	Info(ctx context.Context, message string)
	Warn(ctx context.Context, message string)
	Error(ctx context.Context, message string)
	Fatal(ctx context.Context, message string)
}

// Logger wraps zap logger
type Logger struct {
	logger *zap.Logger
}

// Ensure Logger implements Interface
var _ Interface = (*Logger)(nil)

// New creates a new logger instance
func New(format string, logLevel string) *Logger {
	var config zap.Config

	switch format {
	case "console":
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.Encoding = "console"
		config.EncoderConfig.ConsoleSeparator = " "
	case "json":
		config = zap.NewProductionConfig()
	default:
		config.Encoding = "json" // Default to production config for unknown formats
	}

	// Set the log level from config
	config.Level = zap.NewAtomicLevelAt(GetLogLevel(logLevel))
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	return &Logger{
		logger: logger,
	}
}

// GetTraceIDFromContext extracts trace ID from context
func GetTraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	// Get trace ID from gin context
	if gc, ok := ctx.(*gin.Context); ok {
		if traceID := gc.GetString(TraceIDKey); traceID != "" {
			return traceID
		}
	}
	// Fallback to context value
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// log is a common logging function that handles different log levels
func (l *Logger) log(ctx context.Context, level LogLevel, message string) {
	traceID := GetTraceIDFromContext(ctx)
	var fields []zap.Field
	if traceID == "" {
		fields = []zap.Field{}
	} else {
		fields = []zap.Field{zap.String("trace_id", traceID)}
	}

	switch level {
	case DebugLevel:
		l.logger.Debug(message, fields...)
	case InfoLevel:
		l.logger.Info(message, fields...)
	case WarnLevel:
		l.logger.Warn(message, fields...)
	case ErrorLevel:
		l.logger.Error(message, fields...)
	case FatalLevel:
		l.logger.Fatal(message, fields...)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(ctx context.Context, message string) {
	l.log(ctx, DebugLevel, message)
}

// Info logs an info message
func (l *Logger) Info(ctx context.Context, message string) {
	l.log(ctx, InfoLevel, message)
}

// Warn logs a warning message
func (l *Logger) Warn(ctx context.Context, message string) {
	l.log(ctx, WarnLevel, message)
}

// Error logs an error message
func (l *Logger) Error(ctx context.Context, message string) {
	l.log(ctx, ErrorLevel, message)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(ctx context.Context, message string) {
	l.log(ctx, FatalLevel, message)
}
