package logger

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
)

// Level - type to declare logger options.
type Level = int

// Logger levels describes as int.
const (
	TraceLevel Level = iota - 1
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel

	DisabledLevel = 7
)

// GlobalLevel - the function to set up the global logging level.
func GlobalLevel(l Level) {
	if l != DisabledLevel {
		l = min(max(l, TraceLevel), ErrorLevel)
	}

	zerolog.SetGlobalLevel(zerolog.Level(l))
}

// GlobalLevelFromString - the function that sets up the global logging level using string.
func GlobalLevelFromString(str string) {
	level, _ := ParseLevel(str)

	GlobalLevel(level)
}

// ParseLevel - the function that parses the global logging level using string.
func ParseLevel(str string) (Level, error) {
	switch strings.ToLower(strings.TrimSpace(str)) {
	case zerolog.Level(TraceLevel).String():
		return TraceLevel, nil
	case zerolog.Level(DebugLevel).String():
		return DebugLevel, nil
	case zerolog.Level(InfoLevel).String():
		return InfoLevel, nil
	case zerolog.Level(WarnLevel).String():
		return WarnLevel, nil
	case zerolog.Level(ErrorLevel).String():
		return ErrorLevel, nil
	case zerolog.Level(FatalLevel).String():
		return FatalLevel, nil
	case zerolog.Level(DisabledLevel).String():
		return DisabledLevel, nil
	}

	return ErrorLevel, fmt.Errorf("unknown Level string: '%s', defaulting to ErrorLevel", str)
}
