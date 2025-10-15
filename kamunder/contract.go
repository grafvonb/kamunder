package kamunder

import (
	"context"

	"github.com/grafvonb/camunder/kamunder/cluster"
	"github.com/grafvonb/camunder/kamunder/process"
	"github.com/grafvonb/camunder/kamunder/task"
)

type API interface {
	Capabilities(ctx context.Context) (Capabilities, error)
	process.API
	task.API
	cluster.API
}

type Capabilities struct {
	APIVersion string
	Features   map[Feature]bool
}
type Feature string

func (c Capabilities) Has(f Feature) bool { return c.Features[f] }
