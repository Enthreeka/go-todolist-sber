package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	sugarLogger *zap.SugaredLogger
}

func (l *Logger) Info(format string, v ...any) {
	l.sugarLogger.Infof(format, v...)
}

func (l *Logger) Error(format string, v ...any) {
	l.sugarLogger.Errorf(format, v...)
}

func (l *Logger) Fatal(format string, v ...any) {
	l.sugarLogger.Fatalf(format, v...)
}

func New() *Logger {

	config := zap.NewDevelopmentConfig()
	config.DisableStacktrace = true

	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	logger, _ := config.Build(zap.AddCallerSkip(1))
	sugarLogger := logger.Sugar()

	log := &Logger{
		sugarLogger: sugarLogger,
	}

	return log
}
