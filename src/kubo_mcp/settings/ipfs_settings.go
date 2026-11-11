package settings

import (
	"net/url"
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"

	common_settings "github.com/FournyP/kubo-mcp/src/kubo_mcp_common/settings"
)

const (
	DefaultMaxUploadSizeBytes = int64(100 << 20)
	DefaultMaxReadSizeBytes   = int64(5 << 20)
	DefaultPublicGatewayUrl   = "https://ipfs.io"
)

type IpfsSettings struct {
	MaxUploadSizeBytes int64
	MaxReadSizeBytes   int64
	PublicGatewayUrl   string
}

func NewIpfsSettings(logger *zap.Logger, k *koanf.Koanf) *IpfsSettings {
	maxUploadSizeBytes := common_settings.PositiveInt64(
		logger, k, "ipfs.max_upload_size_bytes", DefaultMaxUploadSizeBytes,
	)
	maxReadSizeBytes := common_settings.PositiveInt64(logger, k, "ipfs.max_read_size_bytes", DefaultMaxReadSizeBytes)

	publicGatewayUrl := strings.TrimSuffix(strings.TrimSpace(k.String("ipfs.public_gateway_url")), "/")
	if publicGatewayUrl == "" {
		publicGatewayUrl = DefaultPublicGatewayUrl
		logger.Info("ipfs.public_gateway_url is not set, using default value", zap.String("value", publicGatewayUrl))
	}

	parsedGatewayUrl, err := url.Parse(publicGatewayUrl)
	if err != nil || (parsedGatewayUrl.Scheme != "http" && parsedGatewayUrl.Scheme != "https") ||
		parsedGatewayUrl.Host == "" {
		logger.Fatal("ipfs.public_gateway_url must be an absolute http(s) URL", zap.String("value", publicGatewayUrl))
	}

	return &IpfsSettings{
		MaxUploadSizeBytes: maxUploadSizeBytes,
		MaxReadSizeBytes:   maxReadSizeBytes,
		PublicGatewayUrl:   publicGatewayUrl,
	}
}
