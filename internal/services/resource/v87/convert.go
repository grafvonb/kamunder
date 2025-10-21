package v87

import (
	camundav87 "github.com/grafvonb/kamunder/internal/clients/camunda/v87/camunda"
	d "github.com/grafvonb/kamunder/internal/domain"
	"github.com/grafvonb/kamunder/toolx"
)

func fromDeploymentResult(r camundav87.DeploymentResult) d.Deployment {
	return d.Deployment{
		TenantId: toolx.Deref(r.TenantId, ""),
		Units:    make([]d.DeploymentUnit, 0),
		Key:      "<unknown>",
	}
}
