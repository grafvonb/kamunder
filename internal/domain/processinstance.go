package domain

import (
	"fmt"
)

type ProcessInstance struct {
	BpmnProcessId             string
	EndDate                   string
	Incident                  bool
	Key                       int64
	ParentFlowNodeInstanceKey int64
	ParentKey                 int64
	ProcessDefinitionKey      int64
	ProcessVersion            int32
	ProcessVersionTag         string
	StartDate                 string
	State                     State
	TenantId                  string
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
