package cluster

import (
	"context"

	csvc "github.com/grafvonb/kamunder/internal/services/cluster"
)

type API interface {
	GetClusterTopology(ctx context.Context) (Topology, error)
}

type client struct{ api csvc.API }

func New(api csvc.API) API { return &client{api: api} }

func (c *client) GetClusterTopology(ctx context.Context) (Topology, error) {
	t, err := c.api.GetClusterTopology(ctx)
	if err != nil {
		return Topology{}, err
	}
	return fromDomainTopology(t), nil
}
