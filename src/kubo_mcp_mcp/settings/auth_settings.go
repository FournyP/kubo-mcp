package settings

import (
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"

	common_settings "github.com/FournyP/kubo-mcp/src/kubo_mcp_common/settings"
)

type AuthSettings struct {
	// ApiKeys is empty when /mcp is open.
	ApiKeys []string
}

func NewAuthSettings(logger *zap.Logger, k *koanf.Koanf) *AuthSettings {
	apiKeys := common_settings.StringSlice(k, "auth.api_keys")
	if len(apiKeys) == 0 {
		logger.Warn("auth.api_keys is not set: /mcp accepts unauthenticated requests")
	}

	return &AuthSettings{ApiKeys: apiKeys}
}
