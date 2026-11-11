package configurations

import (
	"github.com/knadh/koanf/v2"
	"go.uber.org/dig"

	kubo_mcp_configurations "github.com/FournyP/kubo-mcp/src/kubo_mcp/configurations"
	kubo_mcp_host_settings "github.com/FournyP/kubo-mcp/src/kubo_mcp_host/settings"
	kubo_mcp_http_downloader_configurations "github.com/FournyP/kubo-mcp/src/kubo_mcp_http_downloader/configurations"
	kubo_mcp_http_kubo_configurations "github.com/FournyP/kubo-mcp/src/kubo_mcp_http_kubo/configurations"
	kubo_mcp_mcp_configurations "github.com/FournyP/kubo-mcp/src/kubo_mcp_mcp/configurations"
)

func ConfigureDI(k *koanf.Koanf) *dig.Container {
	container := dig.New()

	err := container.Provide(func() *koanf.Koanf { return k })
	if err != nil {
		panic(err)
	}

	err = container.Provide(kubo_mcp_host_settings.NewHostSettings)
	if err != nil {
		panic(err)
	}

	kubo_mcp_configurations.AddKuboMcpConfiguration(container)
	kubo_mcp_http_kubo_configurations.AddKuboMcpHttpKuboConfiguration(container)
	kubo_mcp_http_downloader_configurations.AddKuboMcpHttpDownloaderConfiguration(container)
	kubo_mcp_mcp_configurations.AddKuboMcpMCPConfiguration(container)

	return container
}
