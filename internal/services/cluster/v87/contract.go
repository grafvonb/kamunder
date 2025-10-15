package v87

import (
	"context"

	camundav87 "github.com/grafvonb/kamunder/internal/clients/camunda/v87/camunda"
)

type GenClusterClient interface {
	GetTopologyWithResponse(ctx context.Context, reqEditors ...camundav87.RequestEditorFn) (*camundav87.GetTopologyResponse, error)
}

var _ GenClusterClient = (*camundav87.ClientWithResponses)(nil)
