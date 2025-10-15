package oauth2

import (
	"context"
	"io"

	auth "github.com/grafvonb/kamunder/internal/clients/auth/oauth2"
	"github.com/grafvonb/kamunder/internal/services/auth/core"
)

const formContentType = "application/x-www-form-urlencoded"

type GenAuthClient interface {
	RequestTokenWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...auth.RequestEditorFn) (*auth.RequestTokenResponse, error)
}

var _ core.Authenticator = (*Service)(nil)
