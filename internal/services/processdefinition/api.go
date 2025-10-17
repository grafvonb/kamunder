package processdefinition

import (
	"context"

	d "github.com/grafvonb/kamunder/internal/domain"
	"github.com/grafvonb/kamunder/internal/services"
	v87 "github.com/grafvonb/kamunder/internal/services/processdefinition/v87"
	v88 "github.com/grafvonb/kamunder/internal/services/processdefinition/v88"
)

type API interface {
	GetProcessDefinitionByKey(ctx context.Context, key string, opts ...services.CallOption) (d.ProcessDefinition, error)
	SearchProcessDefinitions(ctx context.Context, filter d.ProcessDefinitionSearchFilterOpts, size int32, opts ...services.CallOption) ([]d.ProcessDefinition, error)
}

var _ API = (*v87.Service)(nil)
var _ API = (*v88.Service)(nil)
