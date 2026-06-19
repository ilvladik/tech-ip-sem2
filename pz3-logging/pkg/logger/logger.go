package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if lvl := os.Getenv("LOG_LEVEL"); lvl != "" {
		if err := level.UnmarshalText([]byte(lvl)); err != nil {
			return nil, err
		}
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleEncoder := zapcore.NewJSONEncoder(encoderCfg)
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)

	cores := []zapcore.Core{consoleCore}

	logFile := os.Getenv("LOG_FILE")
	if logFile == "" {
		logFile = "logs/app.log"
	}
	if err := os.MkdirAll(filepath.Dir(logFile), 0o755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	fileEncoder := zapcore.NewJSONEncoder(encoderCfg)
	fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(file), level)
	cores = append(cores, fileCore)

	core := zapcore.NewTee(cores...)
	return zap.New(core, zap.AddCaller()), nil
}
