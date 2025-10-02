package logger

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger() {
	if err := os.MkdirAll("logs", 0o755); err != nil {
		log.Fatalf("Failed to create logs dir: %v", err)
	}

	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	fileWriter := zapcore.AddSync(file)
	consoleWriter := zapcore.AddSync(os.Stdout)

	encCfg := zap.NewProductionEncoderConfig()
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	level := zapcore.InfoLevel
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encCfg), fileWriter, level),
		zapcore.NewCore(zapcore.NewConsoleEncoder(encCfg), consoleWriter, level),
	)

	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}
