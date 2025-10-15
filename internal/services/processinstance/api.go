package processinstance

import (
	"context"

	d "github.com/grafvonb/camunder/internal/domain"
	v87 "github.com/grafvonb/camunder/internal/services/processinstance/v87"
	v88 "github.com/grafvonb/camunder/internal/services/processinstance/v88"
)

type API interface {
	GetProcessInstanceByKey(ctx context.Context, key int64) (d.ProcessInstance, error)
	GetDirectChildrenOfProcessInstance(ctx context.Context, key int64) ([]d.ProcessInstance, error)
	FilterProcessInstanceWithOrphanParent(ctx context.Context, items []d.ProcessInstance) ([]d.ProcessInstance, error)
	SearchForProcessInstances(ctx context.Context, filter d.ProcessInstanceSearchFilterOpts, size int32) ([]d.ProcessInstance, error)
	CancelProcessInstance(ctx context.Context, key int64) (d.CancelResponse, error)
	DeleteProcessInstance(ctx context.Context, key int64) (d.ChangeStatus, error)
	DeleteProcessInstanceWithCancel(ctx context.Context, key int64) (d.ChangeStatus, error)
	WaitForProcessInstanceState(ctx context.Context, key int64, state d.State) error
}

var _ API = (*v87.Service)(nil)
var _ API = (*v88.Service)(nil)
