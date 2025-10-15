package domain

type ProcessDefinition struct {
	BpmnProcessId string
	Key           int64
	Name          string
	TenantId      string
	Version       int32
	VersionTag    string
}

type ProcessDefinitionSearchFilterOpts struct {
	Key           int64
	BpmnProcessId string
	Version       int32
	VersionTag    string
}
