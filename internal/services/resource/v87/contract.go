package v87

import (
	"context"
	"io"

	camundav87 "github.com/grafvonb/kamunder/internal/clients/camunda/v87/camunda"
)

type GenResourceClient interface {
	PostDeploymentsWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...camundav87.RequestEditorFn) (*camundav87.PostDeploymentsResponse, error)
}

var _ GenResourceClient = (*camundav87.ClientWithResponses)(nil)
