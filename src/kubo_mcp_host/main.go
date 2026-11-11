package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/dig"
	"go.uber.org/zap"

	common_configurations "github.com/FournyP/kubo-mcp/src/kubo_mcp_common/configurations"
	common_helpers "github.com/FournyP/kubo-mcp/src/kubo_mcp_common/helpers"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_host/configurations"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_host/settings"
)

const (
	shutdownTimeout   = 10 * time.Second
	readHeaderTimeout = 10 * time.Second
)

func main() {
	k, err := common_helpers.LoadKoanfConfig("")
	if err != nil {
		panic(err)
	}

	container := configurations.ConfigureDI(k)
	common_configurations.ConfigureLogging(container, k)
	configurations.ConfigureServer(container)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err = container.Invoke(func(
		hostSettings *settings.HostSettings,
		httpHandler http.Handler,
		logger *zap.Logger,
	) error {
		return serveHttp(ctx, hostSettings, httpHandler, logger)
	})
	if err != nil && !errors.Is(err, context.Canceled) {
		// RootCause drops dig's wiring trace so the configuration error reads first.
		common_helpers.FatalOnInvokeError(container, "Server stopped", dig.RootCause(err))
	}
}

func serveHttp(
	ctx context.Context,
	hostSettings *settings.HostSettings,
	handler http.Handler,
	logger *zap.Logger,
) error {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", hostSettings.Port),
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	serveErr := make(chan error, 1)
	go func() {
		logger.Info("Serving MCP over HTTP",
			zap.Int("port", hostSettings.Port),
			zap.String("version", configurations.Version),
		)
		serveErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
	}

	logger.Info(fmt.Sprintf("Shutting down %s...", configurations.HostName))
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	return server.Shutdown(shutdownCtx)
}
