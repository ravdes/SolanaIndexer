package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

var Logger *zap.SugaredLogger

func InitializeLogger() {
	var err error
	config := zap.NewDevelopmentConfig()
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("Jan 2 15:04:05.00")
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	encoderConfig.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(filepath.Base(caller.FullPath()))
	}

	config.EncoderConfig = encoderConfig

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	fileWriteSyncer := zapcore.AddSync(file)

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileEncoderConfig := encoderConfig
	fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	fileEncoder := zapcore.NewConsoleEncoder(fileEncoderConfig)

	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel)
	fileCore := zapcore.NewCore(fileEncoder, fileWriteSyncer, zap.DebugLevel)

	core := zapcore.NewTee(consoleCore, fileCore)

	zapLog := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	Logger = zapLog.Sugar()
}

func Info(msg string, keysAndValues ...any)  { Logger.Infow(msg, keysAndValues...) }
func Error(msg string, keysAndValues ...any) { Logger.Errorw(msg, keysAndValues...) }
func Debug(msg string, keysAndValues ...any) { Logger.Debugw(msg, keysAndValues...) }
func Warn(msg string, keysAndValues ...any)  { Logger.Warnw(msg, keysAndValues...) }
func Fatal(msg string, keysAndValues ...any) { Logger.Fatalw(msg, keysAndValues...) }
func Panic(msg string, keysAndValues ...any) { Logger.Panicw(msg, keysAndValues...) }

func Infof(template string, args ...any)  { Logger.Infof(template, args...) }
func Errorf(template string, args ...any) { Logger.Errorf(template, args...) }
func Debugf(template string, args ...any) { Logger.Debugf(template, args...) }
func Warnf(template string, args ...any)  { Logger.Warnf(template, args...) }
func Fatalf(template string, args ...any) { Logger.Fatalf(template, args...) }
func Panicf(template string, args ...any) { Logger.Panicf(template, args...) }
