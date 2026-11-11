package kubo_mcp_mcp

import (
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_mcp/settings"
)

type ToolRegistrar struct {
	tools       []ToolInterface
	mcpSettings *settings.McpSettings
}

func NewToolRegistrar(tools []ToolInterface, mcpSettings *settings.McpSettings) *ToolRegistrar {
	return &ToolRegistrar{tools: tools, mcpSettings: mcpSettings}
}

// Register exposes every tool, or only mcp.enabled_tools when set; an unknown name fails startup.
func (r *ToolRegistrar) Register(server *mcp.Server) error {
	toolsByName := make(map[string]ToolInterface, len(r.tools))
	for _, tool := range r.tools {
		toolsByName[tool.Name()] = tool
	}

	if len(r.mcpSettings.EnabledTools) == 0 {
		for _, tool := range r.tools {
			tool.Register(server)
		}
		return nil
	}

	for _, name := range r.mcpSettings.EnabledTools {
		if _, exists := toolsByName[name]; !exists {
			return fmt.Errorf("mcp.enabled_tools: unknown tool %q", name)
		}
	}

	registered := make(map[string]bool, len(r.mcpSettings.EnabledTools))
	for _, name := range r.mcpSettings.EnabledTools {
		if !registered[name] {
			toolsByName[name].Register(server)
			registered[name] = true
		}
	}

	return nil
}
