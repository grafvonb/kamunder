package v87

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/grafvonb/kamunder/config"
	camundav87 "github.com/grafvonb/kamunder/internal/clients/camunda/v87/camunda"
	d "github.com/grafvonb/kamunder/internal/domain"
	"github.com/grafvonb/kamunder/internal/services"
)

type Service struct {
	c   GenResourceClient
	cfg *config.Config
	log *slog.Logger
}

type Option func(*Service)

//nolint:unused
func WithClient(c GenResourceClient) Option { return func(s *Service) { s.c = c } }

func New(cfg *config.Config, httpClient *http.Client, log *slog.Logger, opts ...Option) (*Service, error) {
	c, err := camundav87.NewClientWithResponses(
		cfg.APIs.Camunda.BaseURL,
		camundav87.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}
	s := &Service{c: c, cfg: cfg, log: log}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

func (s *Service) Deploy(ctx context.Context, tenantId string, units []d.DeploymentUnitData, opts ...services.CallOption) (d.Deployment, error) {
	_ = services.ApplyCallOptions(opts)
	return d.Deployment{}, nil
}
