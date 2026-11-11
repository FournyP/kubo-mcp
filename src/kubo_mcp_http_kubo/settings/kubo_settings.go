package settings

import (
	"net/url"
	"strings"
	"time"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"

	common_settings "github.com/FournyP/kubo-mcp/src/kubo_mcp_common/settings"
)

const (
	DefaultApiUrl         = "http://127.0.0.1:5001"
	DefaultRequestTimeout = 60 * time.Second
	DefaultUploadTimeout  = 10 * time.Minute
)

type KuboSettings struct {
	ApiUrl string
	// Authorization is sent verbatim as the Authorization header.
	Authorization  string
	RequestTimeout time.Duration
	UploadTimeout  time.Duration
}

func NewKuboSettings(logger *zap.Logger, k *koanf.Koanf) *KuboSettings {
	apiUrl := strings.TrimSuffix(strings.TrimSpace(k.String("kubo.api_url")), "/")
	if apiUrl == "" {
		apiUrl = DefaultApiUrl
		logger.Info("kubo.api_url is not set, using default value", zap.String("value", apiUrl))
	}

	parsedApiUrl, err := url.Parse(apiUrl)
	if err != nil || (parsedApiUrl.Scheme != "http" && parsedApiUrl.Scheme != "https") || parsedApiUrl.Host == "" {
		logger.Fatal("kubo.api_url must be an absolute http(s) URL", zap.String("value", apiUrl))
	}

	return &KuboSettings{
		ApiUrl:         apiUrl,
		Authorization:  strings.TrimSpace(k.String("kubo.authorization")),
		RequestTimeout: common_settings.Duration(logger, k, "kubo.request_timeout", DefaultRequestTimeout),
		UploadTimeout:  common_settings.Duration(logger, k, "kubo.upload_timeout", DefaultUploadTimeout),
	}
}
