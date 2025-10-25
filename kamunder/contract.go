package kamunder

import (
	"context"

	"github.com/grafvonb/kamunder/kamunder/cluster"
	"github.com/grafvonb/kamunder/kamunder/process"
	"github.com/grafvonb/kamunder/kamunder/resource"
	"github.com/grafvonb/kamunder/kamunder/task"
)

type API interface {
	Capabilities(ctx context.Context) (Capabilities, error)
	process.API
	task.API
	cluster.API
	resource.API
}

type Capabilities struct {
	APIVersion string
	Features   map[Feature]bool
}
type Feature string

func (c Capabilities) Has(f Feature) bool { return c.Features[f] }
