package log

import (
	"context"
	"fmt"
	"os"

	"github.com/mntwo/tasklab/internal/config"
	"github.com/mntwo/tasklab/internal/logger"
	"github.com/mntwo/tasklab/internal/logger/zaplog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var defaultLog logger.Logger

func init() {
	defaultLog = New()
}

func New() *zaplog.ZapLogger {
	var coreConfigs []zaplog.CoreConfig
	if getType("std") == "std" {
		coreConfigs = []zaplog.CoreConfig{
			{
				Enc: zapcore.NewConsoleEncoder(humanEncoderConfig()),
				Ws:  zapcore.AddSync(os.Stdout),
				Lvl: zap.NewAtomicLevelAt(zapcore.DebugLevel),
			},
		}
	} else {
		coreConfigs = []zaplog.CoreConfig{
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer("all"),
				Lvl: zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
					return lev == zap.DebugLevel
				}),
			},
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer("debug"),
				Lvl: zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
					return lev == zap.DebugLevel
				}),
			},
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer("info"),
				Lvl: zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
					return lev == zap.InfoLevel
				}),
			},
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer("warn"),
				Lvl: zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
					return lev == zap.WarnLevel
				}),
			},
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer("error"),
				Lvl: zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
					return lev == zap.ErrorLevel
				}),
			},
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer("panic"),
				Lvl: zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
					return lev == zap.PanicLevel
				}),
			},
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer("fatal"),
				Lvl: zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
					return lev == zap.FatalLevel
				}),
			},
		}
	}
	return zaplog.NewZapLogger(
		zaplog.WithExtraKeys([]zaplog.ExtraKey{"trace", "request_id"}),
		zaplog.WithCores(coreConfigs...),
		zaplog.WithZapOptions(zap.AddCaller(), zap.AddCallerSkip(3), zap.Development()),
	)
}

func humanEncoderConfig() zapcore.EncoderConfig {
	cfg := defaultEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncodeDuration = zapcore.StringDurationEncoder
	cfg.EncodeCaller = zapcore.ShortCallerEncoder
	return cfg
}

func getWriteSyncer(suffix string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("logs/%s.log", suffix),
		MaxSize:    getMaxSize(10),
		MaxBackups: getMaxBackups(3),
		MaxAge:     getMaxAge(7),
		Compress:   getCompress(false),
		LocalTime:  getLocalTime(false),
	}
	return zapcore.AddSync(lumberJackLogger)
}

func defaultEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "name",
		TimeKey:        "ts",
		CallerKey:      "caller",
		FunctionKey:    "",
		StacktraceKey:  "stacktrace",
		LineEnding:     "\n",
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func getType(defaultType string) string {
	logCfg := config.GetLog()
	if logCfg == nil {
		return defaultType
	}
	if logCfg.Type != "std" && logCfg.Type != "file" {
		logCfg.Type = "file"
	}
	return logCfg.Type
}

func getMaxSize(defaultMaxSize int) int {
	logCfg := config.GetLog()
	if logCfg == nil || logCfg.MaxSize <= 0 {
		return defaultMaxSize
	}
	return logCfg.MaxSize
}

func getMaxBackups(defaultMaxBackups int) int {
	logCfg := config.GetLog()
	if logCfg == nil {
		return defaultMaxBackups
	}
	return logCfg.MaxBackups
}

func getMaxAge(defaultMaxAge int) int {
	logCfg := config.GetLog()
	if logCfg == nil {
		return defaultMaxAge
	}
	return logCfg.MaxAge
}

func getLocalTime(defaultLocalTime bool) bool {
	logCfg := config.GetLog()
	if logCfg == nil {
		return defaultLocalTime
	}
	return logCfg.LocalTime
}

func getCompress(defaultCompress bool) bool {
	logCfg := config.GetLog()
	if logCfg == nil {
		return defaultCompress
	}
	return logCfg.Compress
}

func Log(ctx context.Context, level logger.Level, msg string, fields ...zap.Field) {
	var args []interface{}
	for _, field := range fields {
		args = append(args, field)
	}
	defaultLog.CtxLog(ctx, level, msg, args...)
}

func Trace(ctx context.Context, msg string, fields ...zap.Field) {
	var args []interface{}
	for _, field := range fields {
		args = append(args, field)
	}
	defaultLog.CtxTrace(ctx, msg, args...)
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	var args []interface{}
	for _, field := range fields {
		args = append(args, field)
	}
	defaultLog.CtxDebug(ctx, msg, args...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	var args []interface{}
	for _, field := range fields {
		args = append(args, field)
	}
	defaultLog.CtxInfo(ctx, msg, args...)
}

func Notice(ctx context.Context, msg string, fields ...zap.Field) {
	var args []interface{}
	for _, field := range fields {
		args = append(args, field)
	}
	defaultLog.CtxNotice(ctx, msg, args...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	var args []interface{}
	for _, field := range fields {
		args = append(args, field)
	}
	defaultLog.CtxWarn(ctx, msg, args...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	var args []interface{}
	for _, field := range fields {
		args = append(args, field)
	}
	defaultLog.CtxError(ctx, msg, args...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	var args []interface{}
	for _, field := range fields {
		args = append(args, field)
	}
	defaultLog.CtxFatal(ctx, msg, args...)
}

func Sync() {
	defaultLog.Sync()
}
