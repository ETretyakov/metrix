// Module "logger" unifies approaches for logging.
package logger

import (
	"context"
	"io"
	"os"

	"github.com/rs/zerolog"
)

var (
	bootstrapLogger zerolog.Logger
	globalLogger    zerolog.Logger
	globalOutput    io.Writer

	withStack         bool
	withCaller        bool
	withConsoleWriter bool
)

func init() {
	bootstrapLogger = zerolog.Nop()
	globalLogger = bootstrapLogger
}

// InitDefault - the builder function to setup default logger.
func InitDefault(level string) {
	globalOutput = os.Stderr

	runLogFile, _ := os.OpenFile(
		"logs/stdout.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	multi := zerolog.MultiLevelWriter(globalOutput, runLogFile)
	bootstrapLogger = zerolog.New(multi)

	applyOptions()

	ctx := bootstrapLogger.With().Timestamp().Caller().Stack()

	globalLogger = ctx.Logger()

	GlobalLevelFromString(level)
}

// Init - the builder function to init logger with options.
func Init(w io.Writer, level string, opts ...Option) {
	if w == nil {
		w = os.Stderr
	}

	globalOutput = w

	bootstrapLogger = zerolog.New(globalOutput)

	applyOptions(opts...)

	ctx := bootstrapLogger.With().Timestamp()

	if withCaller {
		ctx = ctx.Caller()
	}

	if withStack {
		ctx = ctx.Stack()
	}

	globalLogger = ctx.Logger()

	GlobalLevelFromString(level)
}

func consoleWriterSetup(out io.Writer) func(cw *zerolog.ConsoleWriter) {
	return func(cw *zerolog.ConsoleWriter) {
		cw.Out = out

		cw.PartsOrder = []string{
			TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.CallerFieldName,
			MessageFieldName,
		}
	}
}

// Debug - the function for debug logs.
func Debug(ctx context.Context, msg string, kv ...any) {
	globalLogger.Debug().Fields(kv).Msg(msg)
}

// Info - the function for info logs.
func Info(ctx context.Context, msg string, kv ...any) {
	globalLogger.Info().Fields(kv).Msg(msg)
}

// Warn - the function for warn logs.
func Warn(ctx context.Context, msg string, kv ...any) {
	globalLogger.Warn().Fields(kv).Msg(msg)
}

// Error - the function for error logs.
func Error(ctx context.Context, msg string, err error, kv ...any) {
	globalLogger.Error().Fields(kv).Err(err).Msg(msg)
}

// Fatal - the function for fatal logs.
func Fatal(ctx context.Context, msg string, err error, kv ...any) {
	globalLogger.Fatal().Fields(kv).Err(err).Msg(msg)
}
