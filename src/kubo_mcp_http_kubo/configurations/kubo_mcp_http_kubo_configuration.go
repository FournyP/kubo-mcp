package configurations

import (
	"net/http"

	"go.uber.org/dig"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/services/interfaces"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_kubo/clients"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_kubo/helpers"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_kubo/services"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_kubo/settings"
)

func AddKuboMcpHttpKuboConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(settings.NewKuboSettings)
	if err != nil {
		panic(err)
	}

	// Clients: timeouts are per call, uploads need longer than lookups.
	err = container.Provide(func(kuboSettings *settings.KuboSettings) (clients.ClientWithResponsesInterface, error) {
		return clients.NewClientWithResponses(
			kuboSettings.ApiUrl,
			clients.WithHTTPClient(&http.Client{}),
			clients.WithRequestEditorFn(helpers.NewAuthorizationRequestEditor(kuboSettings.Authorization)),
		)
	})
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(services.NewKuboIpfsStorage, dig.As(new(interfaces.IpfsStorageInterface)))
	if err != nil {
		panic(err)
	}
}
