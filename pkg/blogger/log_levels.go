package blogger

import "strings"

const (
	LevelOff      uint8 = 0
	LevelFatal    uint8 = 1
	LevelCritical uint8 = 3
	LevelError    uint8 = 7
	LevelWarn     uint8 = 15
	LevelInfo     uint8 = 31
	LevelDebug    uint8 = 63
	LevelTrace    uint8 = 127
	LevelUnsafe   uint8 = 255

	LevelAll     uint8 = 63
	LevelDefault uint8 = 31
)

const (
	Fatal    uint8 = 1
	Critical uint8 = 2
	Error    uint8 = 4
	Warn     uint8 = 8
	Info     uint8 = 16
	Debug    uint8 = 32
	Trace    uint8 = 64
	Unsafe   uint8 = 128
)

func StringToLevel(value string) uint8 {
	switch strings.ToLower(value) {
	case "off":
		return LevelOff
	case "fatal":
		return LevelFatal
	case "critical":
		return LevelCritical
	case "error":
		return LevelError
	case "warn":
		return LevelWarn
	case "info":
		return LevelInfo
	case "debug":
		return LevelDebug
	case "trace":
		return LevelTrace
	case "unsafe":
		return LevelUnsafe
	default:
		return LevelDefault
	}
}
