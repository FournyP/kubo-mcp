package services_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_downloader/errors"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_downloader/services"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_downloader/settings"
)

type WhenDownloadingARemoteFileTestingSuite struct {
	sut *services.RemoteFileDownloader

	serverUrl string
}

func WhenDownloadingARemoteFileBeforeEach(
	t *testing.T,
	allowPrivateNetworks bool,
	handler http.Handler,
) *WhenDownloadingARemoteFileTestingSuite {
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	downloaderSettings := &settings.DownloaderSettings{
		RequestTimeout:       5 * time.Second,
		AllowPrivateNetworks: allowPrivateNetworks,
	}

	return &WhenDownloadingARemoteFileTestingSuite{
		sut:       services.NewRemoteFileDownloader(downloaderSettings, services.NewDestinationGuard()),
		serverUrl: server.URL,
	}
}

func servePng(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "image/png")
	_, _ = writer.Write([]byte("hello"))
}

func TestWhenDownloadingARemoteFile(t *testing.T) {
	t.Parallel()

	t.Run("Given private networks are blocked", func(t *testing.T) {
		t.Parallel()

		// httptest listens on loopback, which the guard must refuse.
		t.Run("Should return a blocked destination error", func(t *testing.T) {
			t.Parallel()

			suite := WhenDownloadingARemoteFileBeforeEach(t, false, http.HandlerFunc(servePng))

			_, err := suite.sut.Download(context.Background(), suite.serverUrl+"/cat.png")

			var blockedDestinationError *errors.BlockedDestinationError
			require.ErrorAs(t, err, &blockedDestinationError)
			assert.Contains(t, blockedDestinationError.Address, "127.0.0.1")
		})
	})

	t.Run("Given the server answers with a non-OK status", func(t *testing.T) {
		t.Parallel()

		handler := http.NotFoundHandler()

		t.Run("Should return a status error", func(t *testing.T) {
			t.Parallel()

			suite := WhenDownloadingARemoteFileBeforeEach(t, true, handler)

			_, err := suite.sut.Download(context.Background(), suite.serverUrl+"/cat.png")

			assert.EqualError(t, err, "download failed with status 404")
		})
	})

	t.Run("Given too many redirects", func(t *testing.T) {
		t.Parallel()

		handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			http.Redirect(writer, request, "/loop", http.StatusFound)
		})

		t.Run("Should return a redirect error", func(t *testing.T) {
			t.Parallel()

			suite := WhenDownloadingARemoteFileBeforeEach(t, true, handler)

			_, err := suite.sut.Download(context.Background(), suite.serverUrl+"/loop")

			assert.ErrorContains(t, err, "stopped after 5 redirects")
		})
	})

	t.Run("Given a Content-Disposition file name", func(t *testing.T) {
		t.Parallel()

		handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-Disposition", `attachment; filename="dog.png"`)
			servePng(writer, request)
		})

		t.Run("Should use the Content-Disposition file name", func(t *testing.T) {
			t.Parallel()

			suite := WhenDownloadingARemoteFileBeforeEach(t, true, handler)

			remoteFile, err := suite.sut.Download(context.Background(), suite.serverUrl+"/download?id=1")

			require.NoError(t, err)
			defer remoteFile.Body.Close()
			assert.Equal(t, "dog.png", remoteFile.Name)
		})
	})

	t.Run("Given a redirect", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()
		mux.HandleFunc("/short", func(writer http.ResponseWriter, request *http.Request) {
			http.Redirect(writer, request, "/images/cat.png", http.StatusFound)
		})
		mux.HandleFunc("/images/cat.png", servePng)

		t.Run("Should name the file after the final URL", func(t *testing.T) {
			t.Parallel()

			suite := WhenDownloadingARemoteFileBeforeEach(t, true, mux)

			remoteFile, err := suite.sut.Download(context.Background(), suite.serverUrl+"/short")

			require.NoError(t, err)
			defer remoteFile.Body.Close()
			assert.Equal(t, "cat.png", remoteFile.Name)
		})
	})

	t.Run("Given a URL without a path", func(t *testing.T) {
		t.Parallel()

		t.Run("Should leave the file name empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenDownloadingARemoteFileBeforeEach(t, true, http.HandlerFunc(servePng))

			remoteFile, err := suite.sut.Download(context.Background(), suite.serverUrl)

			require.NoError(t, err)
			defer remoteFile.Body.Close()
			assert.Empty(t, remoteFile.Name)
		})
	})

	t.Run("Given the server answers with the file", func(t *testing.T) {
		t.Parallel()

		t.Run("Should populate every field", func(t *testing.T) {
			t.Parallel()

			suite := WhenDownloadingARemoteFileBeforeEach(t, true, http.HandlerFunc(servePng))

			remoteFile, err := suite.sut.Download(context.Background(), suite.serverUrl+"/images/cat.png")

			require.NoError(t, err)
			defer remoteFile.Body.Close()
			content, err := io.ReadAll(remoteFile.Body)
			require.NoError(t, err)
			assert.Equal(t, "cat.png", remoteFile.Name)
			assert.Equal(t, "image/png", remoteFile.ContentType)
			assert.Equal(t, int64(5), remoteFile.SizeBytes)
			assert.Equal(t, "hello", string(content))
		})
	})
}
