package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zapcore.Field

var (
	Int    = zap.Int
	String = zap.String
	Error  = zap.Error
	Bool   = zap.Bool
	Any    = zap.Any
)

type LoggerInterface interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	DPanic(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
}

type logger struct {
	zap *zap.Logger
}

func NewLogger(namespace, level string) LoggerInterface {
	if level == "" {
		level = LevelInfo
	}

	l := logger{
		newZapLogger(namespace, level),
	}
	return &l
}

func (l *logger) Debug(msg string, fields ...Field) {
	l.zap.Debug(msg, fields...)
}

func (l *logger) Info(msg string, fields ...Field) {
	l.zap.Info(msg, fields...)
}

func (l *logger) Warn(msg string, fields ...Field) {
	l.zap.Warn(msg, fields...)
}

func (l *logger) Error(msg string, fields ...Field) {
	l.zap.Error(msg, fields...)
}

func (l *logger) DPanic(msg string, fields ...Field) {
	l.zap.DPanic(msg, fields...)
}

func (l *logger) Panic(msg string, fields ...Field) {
	l.zap.Panic(msg, fields...)
}

func (l *logger) Fatal(msg string, fields ...Field) {
	l.zap.Fatal(msg, fields...)
}

func GetNamed(l LoggerInterface, name string) LoggerInterface {
	switch v := l.(type) {
	case *logger:
		v.zap = v.zap.Named(name)
		return v
	default:
		l.Info("logger.GetNamed: invalid logger type")
		return l
	}
}

func WithFields(l LoggerInterface, fields ...Field) LoggerInterface {
	switch v := l.(type) {
	case *logger:
		return &logger{
			zap: v.zap.With(fields...),
		}
	default:
		l.Info("logger.WithFields: invalid logger type")
		return l
	}
}

func Cleanup(l LoggerInterface) error {
	switch v := l.(type) {
	case *logger:
		return v.zap.Sync()
	default:
		l.Info("logger.Cleanup: invalid logger type")
		return nil
	}
}
