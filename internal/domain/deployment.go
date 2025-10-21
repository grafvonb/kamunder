package domain

type Deployment struct {
	Key      string           `json:"key,omitempty"`
	Units    []DeploymentUnit `json:"units,omitempty"`
	TenantId string           `json:"tenantId,omitempty"`
}

type DeploymentUnit struct {
	ProcessDefinition ProcessDefinitionDeployment `json:"processDefinition,omitempty"`
}

type ProcessDefinitionDeployment struct {
	ProcessDefinitionId      string `json:"processDefinitionId,omitempty"`
	ProcessDefinitionKey     string `json:"processDefinitionKey,omitempty"`
	ProcessDefinitionVersion int32  `json:"processDefinitionVersion,omitempty"`
	ResourceName             string `json:"resourceName,omitempty"`
	TenantId                 string `json:"tenantId,omitempty"`
}
