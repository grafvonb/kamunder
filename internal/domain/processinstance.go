package domain

import (
	"fmt"
)

type ProcessInstance struct {
	BpmnProcessId             string
	EndDate                   string
	Incident                  bool
	Key                       string
	ParentFlowNodeInstanceKey string
	ParentKey                 string
	ProcessDefinitionKey      string
	ProcessVersion            int32
	ProcessVersionTag         string
	StartDate                 string
	State                     State
	TenantId                  string
}

type ProcessInstanceSearchFilterOpts struct {
	Key               string
	BpmnProcessId     string
	ProcessVersion    int32
	ProcessVersionTag string
	State             State
	ParentKey         string
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
