package utils

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"time"

	"github.com/rs/zerolog"

	"gorm.io/gorm/logger"
)

type GormLogger struct {
}

func (l GormLogger) LogMode(logger.LogLevel) logger.Interface {
	// log mode is ignored. Gormlogger will use the global log level
	return l
}

func (l GormLogger) Error(ctx context.Context, msg string, opts ...interface{}) {
	log.Ctx(ctx).Error().Msg(fmt.Sprintf(msg, opts...))
}

func (l GormLogger) Warn(ctx context.Context, msg string, opts ...interface{}) {
	log.Ctx(ctx).Warn().Msg(fmt.Sprintf(msg, opts...))
}

func (l GormLogger) Info(ctx context.Context, msg string, opts ...interface{}) {
	log.Ctx(ctx).Info().Msg(fmt.Sprintf(msg, opts...))
}

func (l GormLogger) Trace(ctx context.Context, begin time.Time, f func() (string, int64), err error) {
	zl := log.Ctx(ctx)
	var event *zerolog.Event

	if err != nil {
		event = zl.Debug()
	} else {
		event = zl.Trace()
	}

	var durKey string

	switch zerolog.DurationFieldUnit {
	case time.Nanosecond:
		durKey = "elapsed_ns"
	case time.Microsecond:
		durKey = "elapsed_us"
	case time.Millisecond:
		durKey = "elapsed_ms"
	case time.Second:
		durKey = "elapsed"
	case time.Minute:
		durKey = "elapsed_min"
	case time.Hour:
		durKey = "elapsed_hr"
	default:
		zl.Error().Interface("zerolog.DurationFieldUnit", zerolog.DurationFieldUnit).Msg("gorm logger encountered a unknown value for DurationFieldUnit")
		durKey = "elapsed_"
	}

	event.Dur(durKey, time.Since(begin))

	sql, rows := f()
	if sql != "" {
		event.Str("sql", sql)
	}
	if rows > -1 {
		event.Int64("rows", rows)
	}

	event.Send()

	return
}
