package v87

import (
	operatev87 "github.com/grafvonb/camunder/internal/clients/camunda/v87/operate"
	d "github.com/grafvonb/camunder/internal/domain"
	"github.com/grafvonb/camunder/toolx"
)

func fromProcessInstanceResponse(r operatev87.ProcessInstance) d.ProcessInstance {
	return d.ProcessInstance{
		BpmnProcessId:             toolx.Deref(r.BpmnProcessId, ""),
		EndDate:                   toolx.Deref(r.EndDate, ""),
		Incident:                  toolx.Deref(r.Incident, false),
		Key:                       toolx.Deref(r.Key, int64(0)),
		ParentFlowNodeInstanceKey: toolx.Deref(r.ParentFlowNodeInstanceKey, int64(0)),
		ParentKey:                 toolx.Deref(r.ParentKey, int64(0)),
		ProcessDefinitionKey:      toolx.Deref(r.ProcessDefinitionKey, int64(0)),
		ProcessVersion:            toolx.Deref(r.ProcessVersion, int32(0)),
		ProcessVersionTag:         toolx.Deref(r.ProcessVersionTag, ""),
		StartDate:                 toolx.Deref(r.StartDate, ""),
		State:                     d.State(*r.State),
		TenantId:                  toolx.Deref(r.TenantId, ""),
	}
}
