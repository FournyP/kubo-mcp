package services

import (
	"context"
	"fmt"
	"mime"
	"net"
	"net/http"
	"path"
	"time"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/models"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_downloader/settings"
)

const (
	maxRedirects = 5
	userAgent    = "kubo-mcp"
)

type RemoteFileDownloader struct {
	httpClient *http.Client
}

func NewRemoteFileDownloader(
	downloaderSettings *settings.DownloaderSettings,
	destinationGuard *DestinationGuard,
) *RemoteFileDownloader {
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	if !downloaderSettings.AllowPrivateNetworks {
		dialer.Control = destinationGuard.Control
	}

	// No proxy: the guard must see the target's IP, not the proxy's.
	transport := &http.Transport{
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
	}

	return &RemoteFileDownloader{
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   downloaderSettings.RequestTimeout,
			CheckRedirect: func(request *http.Request, via []*http.Request) error {
				if len(via) >= maxRedirects {
					return fmt.Errorf("stopped after %d redirects", maxRedirects)
				}
				if request.URL.Scheme != "http" && request.URL.Scheme != "https" {
					return fmt.Errorf("redirect to unsupported scheme %q", request.URL.Scheme)
				}
				return nil
			},
		},
	}
}

func (d *RemoteFileDownloader) Download(ctx context.Context, rawUrl string) (*models.RemoteFile, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, rawUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build the download request: %w", err)
	}
	request.Header.Set("User-Agent", userAgent)

	response, err := d.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		return nil, fmt.Errorf("download failed with status %d", response.StatusCode)
	}

	return &models.RemoteFile{
		Name:        d.resolveFileName(response),
		ContentType: response.Header.Get("Content-Type"),
		SizeBytes:   response.ContentLength,
		Body:        response.Body,
	}, nil
}

func (d *RemoteFileDownloader) resolveFileName(response *http.Response) string {
	if _, params, err := mime.ParseMediaType(response.Header.Get("Content-Disposition")); err == nil {
		if fileName := path.Base(params["filename"]); params["filename"] != "" && fileName != "/" {
			return fileName
		}
	}

	fileName := path.Base(response.Request.URL.Path)
	if fileName == "/" || fileName == "." {
		return ""
	}
	return fileName
}
