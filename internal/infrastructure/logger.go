package infrastructure

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// InitLogger initializes Zap logger in production mode
func InitLogger() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder // ← Human readable time

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.AddSync(os.Stdout),
		zap.InfoLevel,
	)

	Logger = zap.New(core)
}

// GetLogger returns the global logger
func GetLogger() *zap.Logger {
	return Logger
}

// Sync flushes logs before shutdown
func SyncLogger() {
	if Logger != nil {
		Logger.Sync()
	}
}
