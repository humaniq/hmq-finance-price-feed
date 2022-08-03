package blogger

import (
	"bytes"
	"context"
)

func LogLevelFilter(logLevel uint8) func(next BufferedHandler) BufferedHandler {
	return func(next BufferedHandler) BufferedHandler {
		return BufferedHandlerFunc(func(ctx context.Context, level uint8, b bytes.Buffer, text string, args ...interface{}) {
			if logLevel&level > 0 {
				return
			}
			next.HandleLogBuffer(ctx, level, b, text, args)
		})
	}
}
