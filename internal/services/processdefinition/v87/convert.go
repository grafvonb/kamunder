package v87

import (
	operatev87 "github.com/grafvonb/camunder/internal/clients/camunda/v87/operate"
	d "github.com/grafvonb/camunder/internal/domain"
	"github.com/grafvonb/camunder/toolx"
)

func fromProcessDefinitionResponse(r operatev87.ProcessDefinition) d.ProcessDefinition {
	return d.ProcessDefinition{
		BpmnProcessId: toolx.Deref(r.BpmnProcessId, ""),
		Key:           toolx.Deref(r.Key, int64(0)),
		Name:          toolx.Deref(r.Name, ""),
		TenantId:      toolx.Deref(r.TenantId, ""),
		Version:       toolx.Deref(r.Version, int32(0)),
		VersionTag:    toolx.Deref(r.VersionTag, ""),
	}
}
