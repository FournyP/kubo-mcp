package kubo_mcp_mcp

import "github.com/modelcontextprotocol/go-sdk/mcp"

type ToolInterface interface {
	Name() string
	Register(server *mcp.Server)
}
