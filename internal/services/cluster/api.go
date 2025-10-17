package cluster

import (
	"context"

	d "github.com/grafvonb/kamunder/internal/domain"
	"github.com/grafvonb/kamunder/internal/services"
	v87 "github.com/grafvonb/kamunder/internal/services/cluster/v87"
	v88 "github.com/grafvonb/kamunder/internal/services/cluster/v88"
)

type API interface {
	GetClusterTopology(ctx context.Context, opts ...services.CallOption) (d.Topology, error)
}

var _ API = (*v87.Service)(nil)
var _ API = (*v88.Service)(nil)
