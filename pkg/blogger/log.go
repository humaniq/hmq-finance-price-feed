package blogger

import (
	"bytes"
	"context"
	"io"
)

type Log struct {
	handler BufferedHandler
}

func NewLog(
	middlewares []func(next BufferedHandler) BufferedHandler,
	out BufferedHandler,
) *Log {
	handler := out
	for i := len(middlewares) - 1; i >= 0; i-- {
		fn := middlewares[i]
		handler = fn(handler)
	}
	return &Log{
		handler: handler,
	}
}
func (l *Log) Fatal(ctx context.Context, text string, args ...interface{}) {
	l.handler.HandleLogBuffer(ctx, LevelFatal, bytes.Buffer{}, text, args)
}
func (l *Log) Critical(ctx context.Context, text string, args ...interface{}) {
	l.handler.HandleLogBuffer(ctx, LevelFatal, bytes.Buffer{}, text, args)
}
func (l *Log) Error(ctx context.Context, text string, args ...interface{}) {
	l.handler.HandleLogBuffer(ctx, LevelFatal, bytes.Buffer{}, text, args)
}
func (l *Log) Warn(ctx context.Context, text string, args ...interface{}) {
	l.handler.HandleLogBuffer(ctx, LevelFatal, bytes.Buffer{}, text, args)
}
func (l *Log) Info(ctx context.Context, text string, args ...interface{}) {
	l.handler.HandleLogBuffer(ctx, LevelFatal, bytes.Buffer{}, text, args)
}
func (l *Log) Debug(ctx context.Context, text string, args ...interface{}) {
	l.handler.HandleLogBuffer(ctx, LevelFatal, bytes.Buffer{}, text, args)
}
func (l *Log) Trace(ctx context.Context, text string, args ...interface{}) {
	l.handler.HandleLogBuffer(ctx, LevelFatal, bytes.Buffer{}, text, args)
}
func (l *Log) Unsafe(ctx context.Context, text string, args ...interface{}) {
	l.handler.HandleLogBuffer(ctx, LevelFatal, bytes.Buffer{}, text, args)
}

type IOWriterRouter struct {
	out    BufferedHandler
	err    BufferedHandler
	unsafe BufferedHandler
}

func NewIOWriterRouter(out, err, unsafe io.Writer, forceNewLine bool) *IOWriterRouter {
	return &IOWriterRouter{
		out:    IOWriter(out, forceNewLine),
		err:    IOWriter(err, forceNewLine),
		unsafe: IOWriter(unsafe, forceNewLine),
	}
}
func (iowr *IOWriterRouter) HandleLogBuffer(ctx context.Context, level uint8, b bytes.Buffer, text string, args ...interface{}) {
	if level <= LevelError {
		iowr.err.HandleLogBuffer(ctx, level, b, text, args)
		return
	}
	if level < LevelUnsafe {
		iowr.out.HandleLogBuffer(ctx, level, b, text, args)
		return
	}
	iowr.unsafe.HandleLogBuffer(ctx, level, b, text, args)
}
