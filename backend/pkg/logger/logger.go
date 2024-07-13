package logger

import (
	"context"
	"github.com/rs/zerolog"
	gormLogger "gorm.io/gorm/logger"
	"os"
	"strings"
	"time"
)

type Logger interface {
	Debug(ctx context.Context, message string, args ...any)
	Info(ctx context.Context, message string, args ...any)
	Warn(ctx context.Context, message string, args ...any)
	Error(ctx context.Context, message string, args ...any)
	Fatal(ctx context.Context, message string, args ...any)

	LogMode(level gormLogger.LogLevel) gormLogger.Interface
	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}

type logImpl struct {
	logger *zerolog.Logger
}

var _ Logger = (*logImpl)(nil)

type LogLevel string

const (
	ErrorLevel LogLevel = "error"
	WarnLevel  LogLevel = "warn"
	InfoLevel  LogLevel = "info"
	DebugLevel LogLevel = "debug"
)

func New(level string) {
	var l zerolog.Level

	switch LogLevel(strings.ToLower(level)) {
	case ErrorLevel:
		l = zerolog.ErrorLevel
	case WarnLevel:
		l = zerolog.WarnLevel
	case InfoLevel:
		l = zerolog.InfoLevel
	case DebugLevel:
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	skipFrameCount := 1
	logger := zerolog.
		New(os.Stdout).
		With().Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
		Logger()

	log = &logImpl{&logger}
}

func (l *logImpl) Debug(ctx context.Context, message string, args ...interface{}) {
	l.prepare(ctx).Debug().Msgf(message, args...)
}

func (l *logImpl) Info(ctx context.Context, message string, args ...interface{}) {
	l.prepare(ctx).Info().Msgf(message, args...)
}

func (l *logImpl) Warn(ctx context.Context, message string, args ...interface{}) {
	l.prepare(ctx).Warn().Msgf(message, args...)
}

func (l *logImpl) Error(ctx context.Context, message string, args ...interface{}) {
	l.prepare(ctx).Error().Msgf(message, args...)
}

func (l *logImpl) Fatal(ctx context.Context, message string, args ...interface{}) {
	l.prepare(ctx).Fatal().Msgf(message, args...)

	os.Exit(1)
}

func (l *logImpl) prepare(ctx context.Context) *zerolog.Logger {
	reqID, ok := ctx.Value("req-id").(string)
	if !ok || reqID == "" {
		reqID = "not-provided"
	}

	logger := l.logger.With().
		Str("req-id", reqID).
		Logger()

	return &logger
}

func (l *logImpl) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	var logLevel LogLevel

	switch level {
	case gormLogger.Silent:
		logLevel = DebugLevel
	case gormLogger.Error:
		logLevel = ErrorLevel
	case gormLogger.Warn:
		logLevel = WarnLevel
	case gormLogger.Info:
		logLevel = InfoLevel
	}

	New(string(logLevel))
	return log
}

func (l *logImpl) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	event := l.prepare(ctx).Info().
		Dur("elapsed", elapsed).
		Int64("rows", rows).
		Str("sql", sql)

	if err != nil {
		event = event.Err(err)
	}

	event.Msg("")
}

var log *logImpl

// Log Get logger
func Log() Logger {
	if log == nil {
		New("")
	}
	return log
}
