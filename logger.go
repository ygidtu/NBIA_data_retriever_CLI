package main

import (
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap/zapcore"
)

// newEncoderConfig create EncoderConfig for zap
func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// timeEncoder format logger time format
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// setLogger init the zap logger
func setLogger(debug bool, logfile string) {
	encoder := newEncoderConfig()
	level := zap.InfoLevel
	if debug {
		level = zap.DebugLevel
	}

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoder), zapcore.AddSync(os.Stdout), level)
	logger_ := zap.New(core, zap.AddCaller())
	if logfile != "" {
		_ = os.MkdirAll(filepath.Dir(logfile), os.ModePerm)
		f, err := os.OpenFile(logfile, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			logger.Warnf("failed save log to %s: %v", logfile, err)
		} else {
			core = zapcore.NewTee(
				zapcore.NewCore(zapcore.NewJSONEncoder(encoder), zapcore.AddSync(f), zap.DebugLevel),
				zapcore.NewCore(zapcore.NewConsoleEncoder(encoder), zapcore.AddSync(os.Stdout), level),
			)
		}
		logger_ = zap.New(core, zap.AddCaller())
	}

	defer logger_.Sync()
	logger = logger_.Sugar()
}
