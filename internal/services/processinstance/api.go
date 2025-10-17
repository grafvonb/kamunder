package processinstance

import (
	"context"

	d "github.com/grafvonb/kamunder/internal/domain"
	"github.com/grafvonb/kamunder/internal/services"
	v87 "github.com/grafvonb/kamunder/internal/services/processinstance/v87"
	v88 "github.com/grafvonb/kamunder/internal/services/processinstance/v88"
)

type API interface {
	GetProcessInstanceByKey(ctx context.Context, key string, opts ...services.CallOption) (d.ProcessInstance, error)
	GetDirectChildrenOfProcessInstance(ctx context.Context, key string, opts ...services.CallOption) ([]d.ProcessInstance, error)
	FilterProcessInstanceWithOrphanParent(ctx context.Context, items []d.ProcessInstance, opts ...services.CallOption) ([]d.ProcessInstance, error)
	SearchForProcessInstances(ctx context.Context, filter d.ProcessInstanceSearchFilterOpts, size int32, opts ...services.CallOption) ([]d.ProcessInstance, error)
	CancelProcessInstance(ctx context.Context, key string, opts ...services.CallOption) (d.CancelResponse, error)
	DeleteProcessInstance(ctx context.Context, key string, opts ...services.CallOption) (d.ChangeStatus, error)
	GetProcessInstanceStateByKey(ctx context.Context, key string, opts ...services.CallOption) (d.State, error)
	WaitForProcessInstanceState(ctx context.Context, key string, state d.State, opts ...services.CallOption) error
}

var _ API = (*v87.Service)(nil)
var _ API = (*v88.Service)(nil)
