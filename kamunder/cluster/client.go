package cluster

import (
	"context"

	csvc "github.com/grafvonb/kamunder/internal/services/cluster"
	"github.com/grafvonb/kamunder/kamunder/ferrors"
	"github.com/grafvonb/kamunder/kamunder/options"
)

type API interface {
	GetClusterTopology(ctx context.Context, opts ...options.FacadeOption) (Topology, error)
}

type client struct{ api csvc.API }

func New(api csvc.API) API { return &client{api: api} }

func (c *client) GetClusterTopology(ctx context.Context, opts ...options.FacadeOption) (Topology, error) {
	t, err := c.api.GetClusterTopology(ctx, options.MapFacadeOptionsToCallOptions(opts)...)
	if err != nil {
		return Topology{}, ferrors.FromDomain(err)
	}
	return fromDomainTopology(t), nil
}
