package logger

import (
	"bytes"
	"context"
	"os"

	"github.com/humaniq/hmq-finance-price-feed/pkg/blogger"
)

type StdLogger struct {
	level uint8
	out   blogger.BufferedHandler
	err   blogger.BufferedHandler
}

func NewStdLogger(level uint8) *StdLogger {
	handlerOut := blogger.LevelPrefix(
		blogger.CurrentTimeFormat("(2006-01-02)(15:04:05)")(
			blogger.CtxStringValues("requestUid", "uid", "wallet")(
				blogger.IOWriter(os.Stdout, true),
			),
		),
	)
	handlerErr := blogger.LevelPrefix(
		blogger.CurrentTimeFormat("2006-01-02T15:04:05Z07:00")(
			blogger.CtxStringValues("uid", "wallet")(
				blogger.IOWriter(os.Stderr, true),
			),
		),
	)
	return &StdLogger{
		level: level,
		out:   handlerOut,
		err:   handlerErr,
	}
}

func (l *StdLogger) Fatal(ctx context.Context, text string, args ...interface{}) {
	if l.level&blogger.Fatal > 0 {
		l.err.HandleLogBuffer(ctx, blogger.Fatal, bytes.Buffer{}, text, args...)
	}
}
func (l *StdLogger) Critical(ctx context.Context, text string, args ...interface{}) {
	if l.level&blogger.Critical > 0 {
		l.err.HandleLogBuffer(ctx, blogger.Critical, bytes.Buffer{}, text, args...)
	}
}
func (l *StdLogger) Error(ctx context.Context, text string, args ...interface{}) {
	if l.level&blogger.Error > 0 {
		l.err.HandleLogBuffer(ctx, blogger.Error, bytes.Buffer{}, text, args...)
	}
}
func (l *StdLogger) Warn(ctx context.Context, text string, args ...interface{}) {
	if l.level&blogger.Warn > 0 {
		l.out.HandleLogBuffer(ctx, blogger.Warn, bytes.Buffer{}, text, args...)
	}
}
func (l *StdLogger) Info(ctx context.Context, text string, args ...interface{}) {
	if l.level&blogger.Info > 0 {
		l.out.HandleLogBuffer(ctx, blogger.Info, bytes.Buffer{}, text, args...)
	}
}
func (l *StdLogger) Debug(ctx context.Context, text string, args ...interface{}) {
	if l.level&blogger.Debug > 0 {
		l.out.HandleLogBuffer(ctx, blogger.Debug, bytes.Buffer{}, text, args...)
	}
}
func (l *StdLogger) Trace(ctx context.Context, text string, args ...interface{}) {
	if l.level&blogger.Trace > 0 {
		l.out.HandleLogBuffer(ctx, blogger.Trace, bytes.Buffer{}, text, args...)
	}
}
func (l *StdLogger) Unsafe(ctx context.Context, text string, args ...interface{}) {
	if l.level&blogger.Unsafe > 0 {
		l.out.HandleLogBuffer(ctx, blogger.Unsafe, bytes.Buffer{}, text, args...)
	}
}
