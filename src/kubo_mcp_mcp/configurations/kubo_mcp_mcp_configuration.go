package configurations

import (
	"go.uber.org/dig"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_mcp"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_mcp/middlewares"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_mcp/settings"
)

type newToolRegistrarParams struct {
	dig.In

	Tools       []kubo_mcp_mcp.ToolInterface `group:"mcp_tools"`
	McpSettings *settings.McpSettings
}

func AddKuboMcpMCPConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(settings.NewMcpSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewAuthSettings)
	if err != nil {
		panic(err)
	}

	// Middlewares
	err = container.Provide(middlewares.NewBearerAuthMiddleware)
	if err != nil {
		panic(err)
	}

	// Tools
	err = container.Provide(
		kubo_mcp_mcp.NewAddFileTool,
		dig.As(new(kubo_mcp_mcp.ToolInterface)),
		dig.Group("mcp_tools"),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		kubo_mcp_mcp.NewAddFileFromUrlTool,
		dig.As(new(kubo_mcp_mcp.ToolInterface)),
		dig.Group("mcp_tools"),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		kubo_mcp_mcp.NewGetFileTool,
		dig.As(new(kubo_mcp_mcp.ToolInterface)),
		dig.Group("mcp_tools"),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		kubo_mcp_mcp.NewGetFileInfoTool,
		dig.As(new(kubo_mcp_mcp.ToolInterface)),
		dig.Group("mcp_tools"),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		kubo_mcp_mcp.NewPinFileTool,
		dig.As(new(kubo_mcp_mcp.ToolInterface)),
		dig.Group("mcp_tools"),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		kubo_mcp_mcp.NewUnpinFileTool,
		dig.As(new(kubo_mcp_mcp.ToolInterface)),
		dig.Group("mcp_tools"),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		kubo_mcp_mcp.NewListPinsTool,
		dig.As(new(kubo_mcp_mcp.ToolInterface)),
		dig.Group("mcp_tools"),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(params newToolRegistrarParams) *kubo_mcp_mcp.ToolRegistrar {
		return kubo_mcp_mcp.NewToolRegistrar(params.Tools, params.McpSettings)
	})
	if err != nil {
		panic(err)
	}

	// Controllers
	err = container.Provide(kubo_mcp_mcp.NewMCPController)
	if err != nil {
		panic(err)
	}
}
