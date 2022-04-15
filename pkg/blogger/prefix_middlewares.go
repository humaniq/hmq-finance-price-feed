package blogger

import (
	"bytes"
	"context"
	"fmt"
	"time"
)

func CtxStringValues(keys ...string) func(next BufferedHandler) BufferedHandler {
	return func(next BufferedHandler) BufferedHandler {
		return BufferedHandlerFunc(func(ctx context.Context, level uint8, b bytes.Buffer, text string, args ...interface{}) {
			b.WriteString("[CTX:")
			for _, key := range keys {
				if value, ok := ctx.Value(key).(string); ok {
					b.WriteString(fmt.Sprintf("%s=%s,", key, value))
				}
			}
			b.WriteRune(']')
			next.HandleLogBuffer(ctx, level, b, text, args...)
		})
	}
}

func CurrentTimeFormat(format string) func(next BufferedHandler) BufferedHandler {
	return func(next BufferedHandler) BufferedHandler {
		return BufferedHandlerFunc(func(ctx context.Context, level uint8, b bytes.Buffer, text string, args ...interface{}) {
			b.WriteRune('[')
			b.WriteString(time.Now().Format(format))
			b.WriteRune(']')
			next.HandleLogBuffer(ctx, level, b, text, args...)
		})
	}
}

func LevelPrefix(next BufferedHandler) BufferedHandler {
	return BufferedHandlerFunc(func(ctx context.Context, level uint8, b bytes.Buffer, text string, args ...interface{}) {
		switch level {
		case Fatal:
			b.WriteString("[FATAL]")
		case Critical:
			b.WriteString("[CRITICAL]")
		case Error:
			b.WriteString("[ERROR]")
		case Info:
			b.WriteString("[INFO]")
		case Debug:
			b.WriteString("[DEBUG]")
		case Trace:
			b.WriteString("[TRACE]")
		case Unsafe:
			b.WriteString("!!!UNSAFE!!!")
		default:
			b.WriteString("???UNKNOWN???")
		}
		next.HandleLogBuffer(ctx, level, b, text, args...)
	})
}
