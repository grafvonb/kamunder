package resource

import (
	"context"

	rsvc "github.com/grafvonb/kamunder/internal/services/resource"
	"github.com/grafvonb/kamunder/kamunder/ferrors"
	"github.com/grafvonb/kamunder/kamunder/options"
)

type API interface {
	DeployProcessDefinition(ctx context.Context, tenantId string, units []DeploymentUnitData, opts ...options.FacadeOption) (ProcessDefinitionDeployment, error)
}

type client struct{ api rsvc.API }

func New(api rsvc.API) API { return &client{api: api} }

func (c *client) DeployProcessDefinition(ctx context.Context, tenantId string, units []DeploymentUnitData, opts ...options.FacadeOption) (ProcessDefinitionDeployment, error) {
	pdd, err := c.api.Deploy(ctx, tenantId, toDeploymentUnitDatas(units), options.MapFacadeOptionsToCallOptions(opts)...)
	if err != nil {
		return ProcessDefinitionDeployment{}, ferrors.FromDomain(err)
	}
	return fromProcessDefinitionDeployment(pdd), nil
}
