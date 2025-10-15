package processdefinition

import (
	"context"

	d "github.com/grafvonb/kamunder/internal/domain"
	v87 "github.com/grafvonb/kamunder/internal/services/processdefinition/v87"
	v88 "github.com/grafvonb/kamunder/internal/services/processdefinition/v88"
)

type API interface {
	GetProcessDefinitionByKey(ctx context.Context, key int64) (d.ProcessDefinition, error)
	SearchProcessDefinitions(ctx context.Context, filter d.ProcessDefinitionSearchFilterOpts, size int32) ([]d.ProcessDefinition, error)
}

var _ API = (*v87.Service)(nil)
var _ API = (*v88.Service)(nil)
