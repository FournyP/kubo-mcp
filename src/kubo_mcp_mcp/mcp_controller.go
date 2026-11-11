package kubo_mcp_mcp

import (
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/settings"
)

const (
	MCPEndpoint         = "/mcp"
	envelopeMarginBytes = int64(1 << 20)
)

type MCPController struct {
	httpHandler http.Handler
}

func NewMCPController(mcpServer *mcp.Server, ipfsSettings *settings.IpfsSettings) *MCPController {
	// Base64 inflates uploads by 4/3, plus room for the JSON-RPC envelope.
	maxRequestBodyBytes := ipfsSettings.MaxUploadSizeBytes/3*4 + envelopeMarginBytes

	streamableHandler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return mcpServer
	}, &mcp.StreamableHTTPOptions{
		Stateless:           true,
		JSONResponse:        true,
		MaxRequestBodyBytes: maxRequestBodyBytes,
	})

	return &MCPController{httpHandler: streamableHandler}
}

func (controller *MCPController) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle(MCPEndpoint, controller.httpHandler)
}
