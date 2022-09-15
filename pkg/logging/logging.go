package logging

import (
	"os"

	"github.com/neermitt/opsos/pkg/globals"
	"github.com/neermitt/opsos/pkg/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
)

func init() {
	config := zap.NewDevelopmentConfig()
	logLevelText := os.Getenv(globals.LogLevelEnvVariable)
	if len(logLevelText) > 0 {
		logLevel, err := zapcore.ParseLevel(logLevelText)
		if err != nil {
			utils.PrintErrorToStdErrorAndExit(err)
		}
		config.Level.SetLevel(logLevel)
	} else {
		config.Level.SetLevel(zap.WarnLevel)
	}

	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	Logger, _ = config.Build()
}

func InitLogger(level zapcore.Level) {
	Logger = Logger.WithOptions(zap.IncreaseLevel(level))
}
