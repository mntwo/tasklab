package zaplog

import (
	"context"
	"io"

	"github.com/mntwo/tasklab/internal/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	l      *zap.SugaredLogger
	config *config
}

func NewZapLogger(opts ...Option) *ZapLogger {
	config := defaultConfig()

	// apply options
	for _, opt := range opts {
		opt.apply(config)
	}

	cores := make([]zapcore.Core, 0, len(config.coreConfigs))
	for _, coreConfig := range config.coreConfigs {
		cores = append(cores, zapcore.NewCore(coreConfig.Enc, coreConfig.Ws, coreConfig.Lvl))
	}

	logger := zap.New(zapcore.NewTee(cores[:]...), config.zapOpts...)

	return &ZapLogger{
		l:      logger.Sugar(),
		config: config,
	}
}

func (l *ZapLogger) CtxLog(ctx context.Context, level logger.Level, msg string, kvs ...interface{}) {
	slog := l.l
	if len(l.config.extraKeys) > 0 {
		for _, k := range l.config.extraKeys {
			if ctx.Value(k) == nil {
				continue
			}
			slog = slog.With(string(k), ctx.Value(k))
		}
	}
	switch level {
	case logger.LevelDebug, logger.LevelTrace:
		slog.Debugw(msg, kvs...)
	case logger.LevelInfo:
		slog.Infow(msg, kvs...)
	case logger.LevelNotice, logger.LevelWarn:
		slog.Warnw(msg, kvs...)
	case logger.LevelError:
		slog.Errorw(msg, kvs...)
	case logger.LevelFatal:
		slog.Fatalw(msg, kvs...)
	default:
		slog.Warnw(msg, kvs...)
	}
}

func (l *ZapLogger) CtxTrace(ctx context.Context, msg string, v ...interface{}) {
	l.CtxLog(ctx, logger.LevelTrace, msg, v...)
}

func (l *ZapLogger) CtxDebug(ctx context.Context, msg string, v ...interface{}) {
	l.CtxLog(ctx, logger.LevelDebug, msg, v...)
}

func (l *ZapLogger) CtxInfo(ctx context.Context, msg string, v ...interface{}) {
	l.CtxLog(ctx, logger.LevelInfo, msg, v...)
}

func (l *ZapLogger) CtxNotice(ctx context.Context, msg string, v ...interface{}) {
	l.CtxLog(ctx, logger.LevelWarn, msg, v...)
}

func (l *ZapLogger) CtxWarn(ctx context.Context, msg string, v ...interface{}) {
	l.CtxLog(ctx, logger.LevelWarn, msg, v...)
}

func (l *ZapLogger) CtxError(ctx context.Context, msg string, v ...interface{}) {
	l.CtxLog(ctx, logger.LevelError, msg, v...)
}

func (l *ZapLogger) CtxFatal(ctx context.Context, msg string, v ...interface{}) {
	l.CtxLog(ctx, logger.LevelFatal, msg, v...)
}

func (l *ZapLogger) SetLevel(level logger.Level) {
	var lvl zapcore.Level
	switch level {
	case logger.LevelTrace, logger.LevelDebug:
		lvl = zap.DebugLevel
	case logger.LevelInfo:
		lvl = zap.InfoLevel
	case logger.LevelWarn, logger.LevelNotice:
		lvl = zap.WarnLevel
	case logger.LevelError:
		lvl = zap.ErrorLevel
	case logger.LevelFatal:
		lvl = zap.FatalLevel
	default:
		lvl = zap.WarnLevel
	}

	l.config.coreConfigs[0].Lvl = lvl

	cores := make([]zapcore.Core, 0, len(l.config.coreConfigs))
	for _, coreConfig := range l.config.coreConfigs {
		cores = append(cores, zapcore.NewCore(coreConfig.Enc, coreConfig.Ws, coreConfig.Lvl))
	}

	logger := zap.New(
		zapcore.NewTee(cores[:]...),
		l.config.zapOpts...)

	l.l = logger.Sugar()
}

func (l *ZapLogger) SetOutput(writer io.Writer) {
	l.config.coreConfigs[0].Ws = zapcore.AddSync(writer)

	cores := make([]zapcore.Core, 0, len(l.config.coreConfigs))
	for _, coreConfig := range l.config.coreConfigs {
		cores = append(cores, zapcore.NewCore(coreConfig.Enc, coreConfig.Ws, coreConfig.Lvl))
	}

	logger := zap.New(
		zapcore.NewTee(cores[:]...),
		l.config.zapOpts...)

	l.l = logger.Sugar()
}

func (l *ZapLogger) Sync() {
	_ = l.l.Sync()
}
