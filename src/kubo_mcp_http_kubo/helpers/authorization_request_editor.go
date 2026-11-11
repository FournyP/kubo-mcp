package helpers

import (
	"context"
	"net/http"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_kubo/clients"
)

func NewAuthorizationRequestEditor(authorization string) clients.RequestEditorFn {
	return func(_ context.Context, request *http.Request) error {
		if authorization != "" {
			request.Header.Set("Authorization", authorization)
		}
		return nil
	}
}
