package dao

import (
	"errors"
	"github.com/rs/zerolog"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type gormZeroLogger struct {
	*zerolog.Logger
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
}

func (l *gormZeroLogger) LogMode(level logger.LogLevel) logger.Interface {
	var targetLevel zerolog.Level
	switch level {
	case logger.Silent:
		targetLevel = zerolog.NoLevel
	case logger.Error:
		targetLevel = zerolog.ErrorLevel
	case logger.Warn:
		targetLevel = zerolog.WarnLevel
	case logger.Info:
		targetLevel = zerolog.DebugLevel
	default:
		targetLevel = zerolog.DebugLevel
	}
	l2 := l.Logger.Level(targetLevel)
	return &gormZeroLogger{
		Logger: &l2,
	}
}

func (l *gormZeroLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.Logger.Info().Msgf(msg, args...)
}

func (l *gormZeroLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.Logger.Warn().Msgf(msg, args...)
}

func (l *gormZeroLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.Logger.Error().Msgf(msg, args...)
}

func (l *gormZeroLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := map[string]interface{}{
		"sql":      sql,
		"duration": elapsed,
	}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		l.Logger.Error().Err(err).Fields(fields).Msg("[GORM] query error")
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.Logger.Warn().Fields(fields).Msgf("[GORM] slow query")
		return
	}

	l.Logger.Debug().Fields(fields).Msgf("[GORM] query")
}
