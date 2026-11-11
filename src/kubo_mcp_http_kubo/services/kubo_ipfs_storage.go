package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"sort"
	"strconv"

	"github.com/samber/lo"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/models"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_kubo/clients"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_kubo/settings"
)

const streamErrorTrailer = "X-Stream-Error"

type KuboIpfsStorage struct {
	client       clients.ClientWithResponsesInterface
	kuboSettings *settings.KuboSettings
}

func NewKuboIpfsStorage(
	client clients.ClientWithResponsesInterface,
	kuboSettings *settings.KuboSettings,
) *KuboIpfsStorage {
	return &KuboIpfsStorage{client: client, kuboSettings: kuboSettings}
}

func (s *KuboIpfsStorage) Add(ctx context.Context, fileName string, content io.Reader) (*models.StoredFile, error) {
	ctx, cancel := context.WithTimeout(ctx, s.kuboSettings.UploadTimeout)
	defer cancel()

	// Streamed, so large files never sit in memory.
	pipeReader, pipeWriter := io.Pipe()
	multipartWriter := multipart.NewWriter(pipeWriter)
	copyDone := make(chan struct{})
	go func() {
		defer close(copyDone)
		part, err := multipartWriter.CreateFormFile("file", fileName)
		if err == nil {
			_, err = io.Copy(part, content)
		}
		if err == nil {
			err = multipartWriter.Close()
		}
		pipeWriter.CloseWithError(err)
	}()
	// Waits for the copy so content is untouched once Add returns.
	defer func() {
		pipeReader.Close()
		<-copyDone
	}()

	// CIDv1 (base32) for subdomain gateways.
	params := &clients.AddParams{CidVersion: lo.ToPtr(1), Pin: lo.ToPtr(true)}
	response, err := s.client.AddWithBodyWithResponse(ctx, params, multipartWriter.FormDataContentType(), pipeReader)
	if err != nil {
		return nil, fmt.Errorf("kubo add request failed: %w", err)
	}
	if response.JSON200 == nil {
		return nil, s.commandError("add", response.HTTPResponse, response.JSONDefault, response.Body)
	}

	sizeBytes, err := strconv.ParseInt(response.JSON200.Size, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the kubo add size %q: %w", response.JSON200.Size, err)
	}

	return &models.StoredFile{Cid: response.JSON200.Hash, Name: response.JSON200.Name, SizeBytes: sizeBytes}, nil
}

func (s *KuboIpfsStorage) Read(ctx context.Context, path string, maxBytes int64) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, s.kuboSettings.RequestTimeout)
	defer cancel()

	response, err := s.client.CatWithResponse(ctx, &clients.CatParams{Arg: path, Length: lo.ToPtr(maxBytes)})
	if err != nil {
		return nil, fmt.Errorf("kubo cat request failed: %w", err)
	}
	if response.StatusCode() != http.StatusOK {
		return nil, s.commandError("cat", response.HTTPResponse, response.JSONDefault, response.Body)
	}

	// Mid-stream failures arrive as a trailer after a 200.
	if streamError := response.HTTPResponse.Trailer.Get(streamErrorTrailer); streamError != "" {
		return nil, fmt.Errorf("kubo cat failed: %s", streamError)
	}

	return response.Body, nil
}

func (s *KuboIpfsStorage) Stat(ctx context.Context, path string) (*models.FileInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, s.kuboSettings.RequestTimeout)
	defer cancel()

	response, err := s.client.FilesStatWithResponse(ctx, &clients.FilesStatParams{Arg: "/ipfs/" + path})
	if err != nil {
		return nil, fmt.Errorf("kubo files/stat request failed: %w", err)
	}
	if response.JSON200 == nil {
		return nil, s.commandError("files/stat", response.HTTPResponse, response.JSONDefault, response.Body)
	}

	return &models.FileInfo{
		Cid:                 response.JSON200.Hash,
		Type:                response.JSON200.Type,
		SizeBytes:           response.JSON200.Size,
		CumulativeSizeBytes: response.JSON200.CumulativeSize,
		Blocks:              response.JSON200.Blocks,
	}, nil
}

func (s *KuboIpfsStorage) Pin(ctx context.Context, cid string) error {
	ctx, cancel := context.WithTimeout(ctx, s.kuboSettings.UploadTimeout)
	defer cancel()

	response, err := s.client.PinAddWithResponse(ctx, &clients.PinAddParams{Arg: cid})
	if err != nil {
		return fmt.Errorf("kubo pin/add request failed: %w", err)
	}
	if response.JSON200 == nil {
		return s.commandError("pin/add", response.HTTPResponse, response.JSONDefault, response.Body)
	}

	return nil
}

func (s *KuboIpfsStorage) Unpin(ctx context.Context, cid string) error {
	ctx, cancel := context.WithTimeout(ctx, s.kuboSettings.RequestTimeout)
	defer cancel()

	response, err := s.client.PinRmWithResponse(ctx, &clients.PinRmParams{Arg: cid})
	if err != nil {
		return fmt.Errorf("kubo pin/rm request failed: %w", err)
	}
	if response.JSON200 == nil {
		return s.commandError("pin/rm", response.HTTPResponse, response.JSONDefault, response.Body)
	}

	return nil
}

func (s *KuboIpfsStorage) ListPins(ctx context.Context) ([]models.Pin, error) {
	ctx, cancel := context.WithTimeout(ctx, s.kuboSettings.RequestTimeout)
	defer cancel()

	response, err := s.client.PinLsWithResponse(ctx, &clients.PinLsParams{Type: lo.ToPtr(clients.Recursive)})
	if err != nil {
		return nil, fmt.Errorf("kubo pin/ls request failed: %w", err)
	}
	if response.JSON200 == nil {
		return nil, s.commandError("pin/ls", response.HTTPResponse, response.JSONDefault, response.Body)
	}

	pins := make([]models.Pin, 0, len(response.JSON200.Keys))
	for cid, entry := range response.JSON200.Keys {
		pins = append(pins, models.Pin{Cid: cid, Type: entry.Type})
	}
	sort.Slice(pins, func(i, j int) bool { return pins[i].Cid < pins[j].Cid })

	return pins, nil
}

func (s *KuboIpfsStorage) commandError(
	command string,
	httpResponse *http.Response,
	kuboError *clients.KuboError,
	body []byte,
) error {
	statusCode := 0
	streamError := ""
	if httpResponse != nil {
		statusCode = httpResponse.StatusCode
		streamError = httpResponse.Trailer.Get(streamErrorTrailer)
	}

	switch {
	case kuboError != nil && kuboError.Message != "":
		return fmt.Errorf("kubo %s failed with status %d: %s", command, statusCode, kuboError.Message)
	case streamError != "":
		return fmt.Errorf("kubo %s failed with status %d: %s", command, statusCode, streamError)
	case len(body) > 0:
		return fmt.Errorf("kubo %s failed with status %d: %s", command, statusCode, lo.Substring(string(body), 0, 512))
	default:
		return fmt.Errorf("kubo %s failed with status %d", command, statusCode)
	}
}
