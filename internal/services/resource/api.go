package resource

import (
	"context"

	d "github.com/grafvonb/kamunder/internal/domain"
	"github.com/grafvonb/kamunder/internal/services"
	v87 "github.com/grafvonb/kamunder/internal/services/resource/v87"
	v88 "github.com/grafvonb/kamunder/internal/services/resource/v88"
)

type API interface {
	Deploy(ctx context.Context, tenantId string, units []d.DeploymentUnitData, opts ...services.CallOption) (d.Deployment, error)
}

var _ API = (*v87.Service)(nil)
var _ API = (*v88.Service)(nil)
