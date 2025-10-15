package kamunder

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/grafvonb/kamunder/config"
	csvc "github.com/grafvonb/kamunder/internal/services/cluster"
	pdsvc "github.com/grafvonb/kamunder/internal/services/processdefinition"
	pisvc "github.com/grafvonb/kamunder/internal/services/processinstance"

	"github.com/grafvonb/kamunder/kamunder/cluster"
	"github.com/grafvonb/kamunder/kamunder/process"
	"github.com/grafvonb/kamunder/kamunder/task"
)

type Option func(*cfg)

func WithConfig(c *config.Config) Option   { return func(x *cfg) { x.cfg = c } }
func WithHTTPClient(h *http.Client) Option { return func(x *cfg) { x.http = h } }
func WithLogger(l *slog.Logger) Option     { return func(x *cfg) { x.log = l } }

func New(opts ...Option) (API, error) {
	c := cfg{
		http: &http.Client{Timeout: 30 * time.Second},
		log:  slog.Default(),
	}
	for _, o := range opts {
		o(&c)
	}
	if c.cfg == nil {
		c.cfg = &config.Config{}
	}
	if c.http == nil {
		c.http = &http.Client{Timeout: 30 * time.Second}
	}
	if c.log == nil {
		c.log = slog.Default()
	}

	// wire internals
	cAPI, err := csvc.New(c.cfg, c.http, c.log)
	if err != nil {
		return nil, err
	}
	pdAPI, err := pdsvc.New(c.cfg, c.http, c.log)
	if err != nil {
		return nil, err
	}
	piAPI, err := pisvc.New(c.cfg, c.http, c.log)
	if err != nil {
		return nil, err
	}

	return &client{
		ClusterAPI: cluster.New(cAPI),
		ProcessAPI: process.New(pdAPI, piAPI),
		TaskAPI:    task.New(pdAPI, piAPI),
		capsFunc: func(context.Context) (Capabilities, error) {
			return Capabilities{
				APIVersion: string(c.cfg.APIs.Version),
				Features:   map[Feature]bool{},
			}, nil
		},
	}, nil
}

type cfg struct {
	cfg  *config.Config
	http *http.Client
	log  *slog.Logger
}

type ClusterAPI = cluster.API
type ProcessAPI = process.API
type TaskAPI = task.API

var _ API = (*client)(nil)

type client struct {
	ClusterAPI
	ProcessAPI
	TaskAPI

	capsFunc func(context.Context) (Capabilities, error)
}

func (c *client) Capabilities(ctx context.Context) (Capabilities, error) { return c.capsFunc(ctx) }
