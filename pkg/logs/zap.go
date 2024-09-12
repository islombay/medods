package logs

import (
	"os"

	"github.com/streamingfast/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/term"
)

func newZapLogger(namespace, level string) *zap.Logger {
	glLevel := parseLevel(level)

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= glLevel && lvl < zapcore.ErrorLevel
	})

	logStdErrorWriter := zapcore.Lock(os.Stderr)
	logStdInfoWriter := zapcore.Lock(os.Stdout)

	isTTY := term.IsTerminal(int(os.Stderr.Fd()))

	core := zapcore.NewTee(
		zapcore.NewCore(logging.NewEncoder(4, isTTY), logStdErrorWriter, highPriority),
		zapcore.NewCore(logging.NewEncoder(4, isTTY), logStdInfoWriter, lowPriority),
	)

	logger := zap.New(
		core,
		zap.AddCaller(), zap.AddCallerSkip(1),
	)

	logger = logger.Named(namespace)
	zap.RedirectStdLog(logger)

	return logger
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelDPanic:
		return zapcore.DPanicLevel
	case LevelPanic:
		return zapcore.PanicLevel
	case LevelFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
