package logger

import (
	"context"
	"io"
)

type Logger interface {
	CtxLog(ctx context.Context, lvl Level, msg string, kvs ...interface{})
	CtxTrace(ctx context.Context, msg string, kvs ...interface{})
	CtxDebug(ctx context.Context, msg string, kvs ...interface{})
	CtxInfo(ctx context.Context, msg string, kvs ...interface{})
	CtxNotice(ctx context.Context, msg string, kvs ...interface{})
	CtxWarn(ctx context.Context, msg string, kvs ...interface{})
	CtxError(ctx context.Context, msg string, kvs ...interface{})
	CtxFatal(ctx context.Context, msg string, kvs ...interface{})
	SetOutput(io.Writer)
	SetLevel(Level)
	Sync()
}

type Level int

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelNotice
	LevelWarn
	LevelError
	LevelFatal
)
