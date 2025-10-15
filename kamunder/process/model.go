package process

import (
	"fmt"
)

type ProcessDefinition struct {
	BpmnProcessId string `json:"bpmnProcessId,omitempty"`
	Key           int64  `json:"key,omitempty"`
	Name          string `json:"name,omitempty"`
	TenantId      string `json:"tenantId,omitempty"`
	Version       int32  `json:"version,omitempty"`
	VersionTag    string `json:"versionTag,omitempty"`
}

type ProcessDefinitions struct {
	Total int32               `json:"total,omitempty"`
	Items []ProcessDefinition `json:"items,omitempty"`
}

type ProcessDefinitionSearchFilterOpts struct {
	Key           int64  `json:"key,omitempty"`
	BpmnProcessId string `json:"bpmnProcessId,omitempty"`
	Version       int32  `json:"version,omitempty"`
	VersionTag    string `json:"versionTag,omitempty"`
}

type ProcessInstance struct {
	BpmnProcessId             string `json:"bpmnProcessId,omitempty"`
	EndDate                   string `json:"endDate,omitempty"`
	Incident                  bool   `json:"incident,omitempty"`
	Key                       int64  `json:"key,omitempty"`
	ParentFlowNodeInstanceKey int64  `json:"parentFlowNodeInstanceKey,omitempty"`
	ParentKey                 int64  `json:"parentKey,omitempty"`
	ParentProcessInstanceKey  int64  `json:"parentProcessInstanceKey,omitempty"`
	ProcessDefinitionKey      int64  `json:"processDefinitionKey,omitempty"`
	ProcessVersion            int32  `json:"processVersion,omitempty"`
	ProcessVersionTag         string `json:"processVersionTag,omitempty"`
	StartDate                 string `json:"startDate,omitempty"`
	State                     State  `json:"state,omitempty"`
	TenantId                  string `json:"tenantId,omitempty"`
}

type ProcessInstances struct {
	Total int32             `json:"total,omitempty"`
	Items []ProcessInstance `json:"items,omitempty"`
}

type ProcessInstanceSearchFilterOpts struct {
	Key               int64
	BpmnProcessId     string
	ProcessVersion    int32
	ProcessVersionTag string
	State             State
	ParentKey         int64
}

type CancelResponse struct {
	StatusCode int
	Status     string
}

type ChangeStatus struct {
	Deleted int64
	Message string
}

func (c ChangeStatus) String() string {
	return fmt.Sprintf("deleted: %d, message: %s", c.Deleted, c.Message)
}
