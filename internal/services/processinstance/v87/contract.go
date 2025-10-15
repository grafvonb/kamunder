package v87

import (
	"context"

	operatev87 "github.com/grafvonb/camunder/internal/clients/camunda/v87/operate"
)

type GenClusterClient interface {
	GetProcessDefinitionByKeyWithResponse(ctx context.Context, key int64, reqEditors ...operatev87.RequestEditorFn) (*operatev87.GetProcessDefinitionByKeyResponse, error)
	SearchProcessDefinitionsWithResponse(ctx context.Context, body operatev87.SearchProcessDefinitionsJSONRequestBody, reqEditors ...operatev87.RequestEditorFn) (*operatev87.SearchProcessDefinitionsResponse, error)
}

var _ GenClusterClient = (*operatev87.ClientWithResponses)(nil)
