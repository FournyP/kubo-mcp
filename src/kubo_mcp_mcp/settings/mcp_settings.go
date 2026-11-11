package settings

import (
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"

	common_settings "github.com/FournyP/kubo-mcp/src/kubo_mcp_common/settings"
)

type McpSettings struct {
	// EnabledTools is empty when every tool is exposed.
	EnabledTools []string
}

func NewMcpSettings(logger *zap.Logger, k *koanf.Koanf) *McpSettings {
	enabledTools := common_settings.StringSlice(k, "mcp.enabled_tools")
	if len(enabledTools) == 0 {
		logger.Info("mcp.enabled_tools is not set, exposing every tool")
	}

	return &McpSettings{EnabledTools: enabledTools}
}
