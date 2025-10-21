package v88

import (
	"context"
	"io"

	camundav88 "github.com/grafvonb/kamunder/internal/clients/camunda/v88/camunda"
)

type GenResourceClient interface {
	CreateDeploymentWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...camundav88.RequestEditorFn) (*camundav88.CreateDeploymentResponse, error)
}

var _ GenResourceClient = (*camundav88.ClientWithResponses)(nil)
