package configurations

import (
	"go.uber.org/dig"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/settings"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases"
)

func AddKuboMcpConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(settings.NewIpfsSettings)
	if err != nil {
		panic(err)
	}

	// Use cases
	err = container.Provide(use_cases.NewAddFile, dig.As(new(use_cases.AddFileInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewAddFileFromUrl, dig.As(new(use_cases.AddFileFromUrlInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewGetFile, dig.As(new(use_cases.GetFileInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewGetFileInfo, dig.As(new(use_cases.GetFileInfoInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewPinFile, dig.As(new(use_cases.PinFileInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewUnpinFile, dig.As(new(use_cases.UnpinFileInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewListPins, dig.As(new(use_cases.ListPinsInterface)))
	if err != nil {
		panic(err)
	}
}
