package configurations

import (
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/dig"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_mcp"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_mcp/middlewares"
)

const (
	HostName       = "kubo-mcp"
	HealthEndpoint = "/health"
	instructions   = "Stores files on IPFS through a Kubo node. Share the returned gatewayUrl or ipfsUri."
)

// Version is set with -ldflags at build time.
var Version = "dev"

func ConfigureServer(container *dig.Container) {
	err := container.Provide(newMCPServer)
	if err != nil {
		panic(err)
	}

	err = container.Provide(newHttpHandler)
	if err != nil {
		panic(err)
	}
}

func newMCPServer(toolRegistrar *kubo_mcp_mcp.ToolRegistrar) (*mcp.Server, error) {
	server := mcp.NewServer(
		&mcp.Implementation{Name: HostName, Version: Version},
		&mcp.ServerOptions{Instructions: instructions},
	)

	if err := toolRegistrar.Register(server); err != nil {
		return nil, err
	}

	return server, nil
}

type newHttpHandlerParams struct {
	dig.In

	BearerAuthMiddleware *middlewares.BearerAuthMiddleware
	MCPController        *kubo_mcp_mcp.MCPController
}

func newHttpHandler(params newHttpHandlerParams) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET "+HealthEndpoint, func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("ok"))
	})

	// Every route but /health goes through auth.
	authenticated := http.NewServeMux()
	params.MCPController.RegisterRoutes(authenticated)
	mux.Handle("/", params.BearerAuthMiddleware.Wrap(authenticated))

	return mux
}
