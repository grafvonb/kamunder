package domain

type ProcessDefinition struct {
	BpmnProcessId string
	Key           string
	Name          string
	TenantId      string
	Version       int32
	VersionTag    string
}

type ProcessDefinitionSearchFilterOpts struct {
	Key           string
	BpmnProcessId string
	Version       int32
	VersionTag    string
}
