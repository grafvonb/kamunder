package process

import (
	"context"

	d "github.com/grafvonb/kamunder/internal/domain"
	pdsvc "github.com/grafvonb/kamunder/internal/services/processdefinition"
	pisvc "github.com/grafvonb/kamunder/internal/services/processinstance"
	"github.com/grafvonb/kamunder/kamunder/options"
	"github.com/grafvonb/kamunder/toolx"
)

type API interface {
	GetProcessDefinitionByKey(ctx context.Context, key int64) (ProcessDefinition, error)
	SearchProcessDefinitions(ctx context.Context, filter ProcessDefinitionSearchFilterOpts, size int32) (ProcessDefinitions, error)

	GetProcessInstanceByKey(ctx context.Context, key int64) (ProcessInstance, error)
	SearchForProcessInstances(ctx context.Context, filter ProcessInstanceSearchFilterOpts, size int32) (ProcessInstances, error)
	CancelProcessInstance(ctx context.Context, key int64, option ...options.FacadeOption) (CancelResponse, error)
	GetDirectChildrenOfProcessInstance(ctx context.Context, key int64) (ProcessInstances, error)
	FilterProcessInstanceWithOrphanParent(ctx context.Context, items []ProcessInstance) ([]ProcessInstance, error)
	DeleteProcessInstance(ctx context.Context, key int64) (ChangeStatus, error)
	DeleteProcessInstanceWithCancel(ctx context.Context, key int64) (ChangeStatus, error)
	WaitForProcessInstanceState(ctx context.Context, key int64, desiredState State) error
}

type client struct {
	pdApi pdsvc.API
	piApi pisvc.API
}

func New(pdApi pdsvc.API, piApi pisvc.API) API {
	return &client{
		pdApi: pdApi,
		piApi: piApi,
	}
}

func (c *client) GetProcessDefinitionByKey(ctx context.Context, key int64) (ProcessDefinition, error) {
	pd, err := c.pdApi.GetProcessDefinitionByKey(ctx, key)
	if err != nil {
		return ProcessDefinition{}, err
	}
	return fromDomainProcessDefinition(pd), nil
}

func (c *client) SearchProcessDefinitions(ctx context.Context, filter ProcessDefinitionSearchFilterOpts, size int32) (ProcessDefinitions, error) {
	pds, err := c.pdApi.SearchProcessDefinitions(ctx, toDomainProcessDefinitionFilter(filter), size)
	if err != nil {
		return ProcessDefinitions{}, err
	}
	return fromDomainProcessDefinitions(pds), nil
}

func (c *client) GetProcessInstanceByKey(ctx context.Context, key int64) (ProcessInstance, error) {
	pi, err := c.piApi.GetProcessInstanceByKey(ctx, key)
	if err != nil {
		return ProcessInstance{}, err
	}
	return fromDomainProcessInstance(pi), nil
}

func (c *client) SearchForProcessInstances(ctx context.Context, filter ProcessInstanceSearchFilterOpts, size int32) (ProcessInstances, error) {
	pis, err := c.piApi.SearchForProcessInstances(ctx, toDomainProcessInstanceFilter(filter), size)
	if err != nil {
		return ProcessInstances{}, err
	}
	return fromDomainProcessInstances(pis), nil
}

func (c *client) CancelProcessInstance(ctx context.Context, key int64, opts ...options.FacadeOption) (CancelResponse, error) {
	resp, err := c.piApi.CancelProcessInstance(ctx, key, options.MapFacadeOptionsToCallOptions(opts)...)
	if err != nil {
		return CancelResponse{}, err
	}
	return CancelResponse{StatusCode: resp.StatusCode, Status: resp.Status}, nil
}

func (c *client) GetDirectChildrenOfProcessInstance(ctx context.Context, key int64) (ProcessInstances, error) {
	children, err := c.piApi.GetDirectChildrenOfProcessInstance(ctx, key)
	if err != nil {
		return ProcessInstances{}, err
	}
	return fromDomainProcessInstances(children), nil
}

func (c *client) FilterProcessInstanceWithOrphanParent(ctx context.Context, items []ProcessInstance) ([]ProcessInstance, error) {
	in := toolx.MapSlice(items, toDomainProcessInstance)
	out, err := c.piApi.FilterProcessInstanceWithOrphanParent(ctx, in)
	if err != nil {
		return nil, err
	}
	return toolx.MapSlice(out, fromDomainProcessInstance), nil
}

func (c *client) DeleteProcessInstance(ctx context.Context, key int64) (ChangeStatus, error) {
	s, err := c.piApi.DeleteProcessInstance(ctx, key)
	if err != nil {
		return ChangeStatus{}, err
	}
	return ChangeStatus{Deleted: s.Deleted, Message: s.Message}, nil
}

func (c *client) DeleteProcessInstanceWithCancel(ctx context.Context, key int64) (ChangeStatus, error) {
	s, err := c.piApi.DeleteProcessInstanceWithCancel(ctx, key)
	if err != nil {
		return ChangeStatus{}, err
	}
	return ChangeStatus{Deleted: s.Deleted, Message: s.Message}, nil
}

func (c *client) WaitForProcessInstanceState(ctx context.Context, key int64, desiredState State) error {
	return c.piApi.WaitForProcessInstanceState(ctx, key, d.State(desiredState))
}
