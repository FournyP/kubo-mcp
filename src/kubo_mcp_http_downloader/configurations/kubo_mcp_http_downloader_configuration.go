package configurations

import (
	"go.uber.org/dig"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/services/interfaces"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_downloader/services"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_downloader/settings"
)

func AddKuboMcpHttpDownloaderConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(settings.NewDownloaderSettings)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(services.NewDestinationGuard)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewRemoteFileDownloader,
		dig.As(new(interfaces.RemoteFileDownloaderInterface)),
	)
	if err != nil {
		panic(err)
	}
}
