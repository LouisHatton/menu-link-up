package log

import (
	"context"

	"github.com/LouisHatton/menu-link-up/internal/objectstore"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

func NewLogger() *Logger {

	config := zap.NewDevelopmentConfig()
	// config.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder // Uncomment to get full caller file
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	logger, err := config.Build(zap.AddCallerSkip(1))

	// logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return &Logger{
		logger: logger,
	}
}

func NewNop() *Logger {
	return &Logger{
		logger: zap.NewNop(),
	}
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	child := Logger{logger: l.logger.With(fields...)}
	return &child
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func Error(err error) zap.Field {
	return zap.Error(err)
}

func String(key string, val string) zap.Field {
	return zap.String(key, val)
}

func Context(ctx context.Context) zap.Field {
	return zap.Skip()
}

func FileId(fileId string) zap.Field {
	return zap.String("fileId", fileId)
}

func UserId(userId string) zap.Field {
	return zap.String("userId", userId)
}

func RequestedId(id string) zap.Field {
	return zap.String("requestedId", id)
}

func FileLocation(loc objectstore.FileLocation) zap.Field {
	return zap.Object("fileLocation", loc)
}
