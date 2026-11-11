package configurations

import (
	"github.com/knadh/koanf/v2"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ConfigureLogging(container *dig.Container, k *koanf.Koanf) {
	zapConfig := zap.NewProductionConfig()

	level, err := zapcore.ParseLevel(k.String("logging.level"))
	if err != nil {
		level = zapcore.InfoLevel
	}
	zapConfig.Level.SetLevel(level)

	logger, err := zapConfig.Build()
	if err != nil {
		panic(err)
	}

	if err := container.Provide(func() *zap.Logger { return logger }); err != nil {
		logger.Fatal("Failed to provide logger", zap.Error(err))
	}
}
