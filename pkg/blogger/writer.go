package blogger

import (
	"bytes"
	"context"
	"fmt"
	"io"
)

func IOWriter(writer io.Writer, forceNewline bool) BufferedHandler {
	return BufferedHandlerFunc(func(ctx context.Context, level uint8, b bytes.Buffer, text string, args ...interface{}) {
		b.WriteRune(' ')
		b.WriteString(fmt.Sprintf(text, args...))
		if forceNewline {
			b.WriteRune('\n')
		}
		writer.Write(b.Bytes())
	})
}
