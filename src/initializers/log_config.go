package initializers

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

func LogConfig(dirPath, logPath string) (*zap.Logger, error) {
	err := os.Mkdir(dirPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("error when creating log dir: %w", err)
	}

	var fullLogPath = logPath + time.Now().Format("02.01.2006_15:04:05") + ".log"

	zapConfig := zap.Config{
		Development:      true,
		Encoding:         "json",
		ErrorOutputPaths: []string{"stderr", fullLogPath},
		OutputPaths:      []string{"stdout", fullLogPath},
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
	}

	logger, err := zapConfig.Build(zap.AddStacktrace(zapcore.ErrorLevel))

	if err != nil {
		return nil, fmt.Errorf("error when creating logger from config: %w", err)
	}

	return logger, nil
}
