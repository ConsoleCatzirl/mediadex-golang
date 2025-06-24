package mlog

import (
	"log/slog"
	"os"
)

var level = new(slog.LevelVar)
var logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))

func Verbose() {
	level.Set(slog.LevelDebug)
}

func Error(msg string, err error, extra ...any) {
	args := []any{}
	args = append(args, "reason", err.Error())
	args = append(args, extra...)
	logger.Error(msg, args...)
}

func Warn(msg string, err error, extra ...any) {
	args := []any{}
	args = append(args, "reason", err.Error())
	args = append(args, extra...)
	logger.Warn(msg, args...)
}

func Info(msg string, extra ...any) {
	logger.Info(msg, extra...)
}

func Trace(fxn string, msg string, extra ...any) {
	args := []any{}
	args = append(args, "function", fxn)
	args = append(args, extra...)
	logger.Debug(msg, args...)
}
