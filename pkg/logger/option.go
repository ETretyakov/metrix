package logger

import (
	"time"

	"github.com/rs/zerolog"
)

const (
	DefaultCallerSkipFrameCount = 3
	DefaultTimeFieldFormat      = time.RFC3339Nano
	MessageFieldName            = "msg"
	TimestampFieldName          = "timestamp"
	StacktraceFieldName         = "stacktrace"
	VersionFieldName            = "version"
)

// DefaultOptions - the variable for default options.
var DefaultOptions = []Option{
	WithDefaults(),
	WithDefaultTimeFormat(),
}

// Option - the type function for logger options.
type Option func()

// WithDefaults - the method to add default options for a logger.
func WithDefaults() Option {
	return func() {
		zerolog.TimestampFieldName = TimestampFieldName
		zerolog.MessageFieldName = MessageFieldName
		zerolog.ErrorStackFieldName = StacktraceFieldName

		withStack = false
		withCaller = false
		withConsoleWriter = false
	}
}

// WithStack - the method to add stack option.
func WithStack() Option {
	return func() {
		withStack = true
	}
}

// WithCaller - the method to add caller option.
func WithCaller() Option {
	return func() {
		withCaller = true
	}
}

// WithConsoleWriter - the method to add console writer option.
func WithConsoleWriter() Option {
	return func() {
		withConsoleWriter = true
		bootstrapLogger = zerolog.New(zerolog.NewConsoleWriter(consoleWriterSetup(globalOutput)))
	}
}

// WithDefaultTimeFormat - the method to add default time format option.
func WithDefaultTimeFormat() Option {
	return func() {
		zerolog.TimeFieldFormat = DefaultTimeFieldFormat
	}
}

func applyOptions(opts ...Option) {
	for _, opt := range DefaultOptions {
		opt()
	}

	for _, opt := range opts {
		opt()
	}
}
