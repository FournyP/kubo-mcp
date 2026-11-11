package settings

import (
	"time"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"

	common_settings "github.com/FournyP/kubo-mcp/src/kubo_mcp_common/settings"
)

const DefaultRequestTimeout = 5 * time.Minute

type DownloaderSettings struct {
	RequestTimeout       time.Duration
	AllowPrivateNetworks bool
}

func NewDownloaderSettings(logger *zap.Logger, k *koanf.Koanf) *DownloaderSettings {
	allowPrivateNetworks := common_settings.Bool(logger, k, "downloader.allow_private_networks", false)
	if allowPrivateNetworks {
		logger.Warn("downloader.allow_private_networks is enabled: add_file_from_url can reach internal services")
	}

	return &DownloaderSettings{
		RequestTimeout:       common_settings.Duration(logger, k, "downloader.request_timeout", DefaultRequestTimeout),
		AllowPrivateNetworks: allowPrivateNetworks,
	}
}
