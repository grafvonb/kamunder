package v88

import (
	operatev88 "github.com/grafvonb/kamunder/internal/clients/camunda/v88/operate"
	d "github.com/grafvonb/kamunder/internal/domain"
	"github.com/grafvonb/kamunder/toolx"
)

func fromProcessDefinitionResponse(r operatev88.ProcessDefinition) d.ProcessDefinition {
	return d.ProcessDefinition{
		BpmnProcessId: toolx.Deref(r.BpmnProcessId, ""),
		Key:           toolx.Int64PtrToString(r.Key),
		Name:          toolx.Deref(r.Name, ""),
		TenantId:      toolx.Deref(r.TenantId, ""),
		Version:       toolx.Deref(r.Version, int32(0)),
		VersionTag:    toolx.Deref(r.VersionTag, ""),
	}
}
