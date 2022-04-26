package blogger

import (
	"bytes"
	"context"
)

type Logger interface {
	Fatal(ctx context.Context, text string, args ...interface{})
	Critical(ctx context.Context, text string, args ...interface{})
	Error(ctx context.Context, text string, args ...interface{})
	Warn(ctx context.Context, text string, args ...interface{})
	Info(ctx context.Context, text string, args ...interface{})
	Debug(ctx context.Context, text string, args ...interface{})
	Trace(ctx context.Context, text string, args ...interface{})
	Unsafe(ctx context.Context, text string, args ...interface{})
}

type BufferedHandler interface {
	HandleLogBuffer(ctx context.Context, level uint8, b bytes.Buffer, text string, args ...interface{})
}
type BufferedHandlerFunc func(ctx context.Context, level uint8, b bytes.Buffer, text string, args ...interface{})

func (bhf BufferedHandlerFunc) HandleLogBuffer(ctx context.Context, level uint8, b bytes.Buffer, text string, args ...interface{}) {
	bhf(ctx, level, b, text, args...)
}
