package settings

import (
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const DefaultPort = 3000

type HostSettings struct {
	Port int
}

func NewHostSettings(logger *zap.Logger, k *koanf.Koanf) *HostSettings {
	port := k.Int("host.port")
	if port == 0 {
		port = DefaultPort
	}
	if port < 0 || port > 65535 {
		logger.Fatal("host.port must be a valid TCP port", zap.Int("value", port))
	}

	return &HostSettings{Port: port}
}
