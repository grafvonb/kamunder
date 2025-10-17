package v88

import (
	camundav88 "github.com/grafvonb/kamunder/internal/clients/camunda/v88/camunda"
	operatev88 "github.com/grafvonb/kamunder/internal/clients/camunda/v88/operate"
	d "github.com/grafvonb/kamunder/internal/domain"
	"github.com/grafvonb/kamunder/toolx"
)

//nolint:unused
func fromProcessInstanceResponse(r operatev88.ProcessInstance) d.ProcessInstance {
	return d.ProcessInstance{
		BpmnProcessId:             toolx.Deref(r.BpmnProcessId, ""),
		EndDate:                   toolx.Deref(r.EndDate, ""),
		Incident:                  toolx.Deref(r.Incident, false),
		Key:                       toolx.Int64PtrToString(r.Key),
		ParentFlowNodeInstanceKey: toolx.Int64PtrToString(r.ParentFlowNodeInstanceKey),
		ParentKey:                 toolx.Int64PtrToString(r.ParentKey),
		ProcessDefinitionKey:      toolx.Int64PtrToString(r.ProcessDefinitionKey),
		ProcessVersion:            toolx.Deref(r.ProcessVersion, int32(0)),
		ProcessVersionTag:         toolx.Deref(r.ProcessVersionTag, ""),
		StartDate:                 toolx.Deref(r.StartDate, ""),
		State:                     d.State(*r.State),
		TenantId:                  toolx.Deref(r.TenantId, ""),
	}
}

func fromProcessInstanceResult(r camundav88.ProcessInstanceResult) d.ProcessInstance {
	return d.ProcessInstance{}
}
