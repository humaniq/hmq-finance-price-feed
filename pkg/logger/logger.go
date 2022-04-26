package logger

import (
	"context"

	"github.com/humaniq/hmq-finance-price-feed/pkg/blogger"
)

var instance blogger.Logger

func Fatal(ctx context.Context, text string, args ...interface{}) {
	get().Fatal(ctx, text, args...)
}
func Critical(ctx context.Context, text string, args ...interface{}) {
	get().Critical(ctx, text, args...)
}
func Error(ctx context.Context, text string, args ...interface{}) {
	get().Error(ctx, text, args...)
}
func Warn(ctx context.Context, text string, args ...interface{}) {
	get().Warn(ctx, text, args...)
}
func Info(ctx context.Context, text string, args ...interface{}) {
	get().Info(ctx, text, args...)
}
func Debug(ctx context.Context, text string, args ...interface{}) {
	get().Debug(ctx, text, args...)
}
func Trace(ctx context.Context, text string, args ...interface{}) {
	get().Trace(ctx, text, args...)
}
func Unsafe(ctx context.Context, text string, args ...interface{}) {
	get().Unsafe(ctx, text, args...)
}

func InitDefault(level uint8) {
	instance = NewStdLogger(level)
}
func Init(logger blogger.Logger) {
	instance = logger
}

func get() blogger.Logger {
	if instance == nil {
		instance = NewStdLogger(blogger.LevelDefault)
	}
	return instance
}
