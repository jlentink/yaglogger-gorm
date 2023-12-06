package yaglogger_gorm

import (
	"context"
	"errors"
	"fmt"
	log "github.com/jlentink/yaglogger"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

func NewLogger() LogWrapper {
	return LogWrapper{
		instance:                  log.GetInstance(),
		IgnoreRecordNotFoundError: true,
		SlowThreshold:             200 * time.Millisecond,
		SilenceQueries:            false,
	}
}

type LogWrapper struct {
	instance                  *log.Logger
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
	SilenceQueries            bool
}

func (l LogWrapper) LogMode(level logger.LogLevel) logger.Interface {
	switch level {
	case logger.Silent:
		l.instance.SetLevel(log.LevelFatal)
	case logger.Error:
		l.instance.SetLevel(log.LevelError)
	case logger.Info:
		l.instance.SetLevel(log.LevelInfo)
	default:
		l.instance.SetLevel(log.LevelInfo)
	}
	return l
}

func (l LogWrapper) Info(_ context.Context, msg string, args ...interface{}) {
	l.instance.Info(msg, args...)
}

func (l LogWrapper) Warn(_ context.Context, msg string, args ...interface{}) {
	l.instance.Warn(msg, args...)
}

func (l LogWrapper) Error(_ context.Context, msg string, args ...interface{}) {
	l.instance.Warn(msg, args...)
}

func (l LogWrapper) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.instance.Level <= log.LevelFatal {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.instance.Error("%s %s - [%.3fms] [rows:%v] %s", utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.instance.Error("%s %s - [%.3fms] [rows:%v] %s", utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.instance.Warn("%s - [warn] %s %d %s %s", utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.instance.Warn("%s - [warn] %s %d %s %s", utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	default:
		if !l.SilenceQueries {
			sql, rows := fc()
			if rows == -1 {
				l.instance.Debug("%s - [%.3fms] [rows:%v] %s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				l.instance.Debug("%s - [%.3fms] [rows:%v] %s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		}
	}
}
