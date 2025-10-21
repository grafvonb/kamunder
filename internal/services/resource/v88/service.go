package v88

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/grafvonb/kamunder/config"
	camundav88 "github.com/grafvonb/kamunder/internal/clients/camunda/v88/camunda"
	d "github.com/grafvonb/kamunder/internal/domain"
	"github.com/grafvonb/kamunder/internal/services"
	"github.com/grafvonb/kamunder/internal/services/httpc"
)

type Service struct {
	c   GenResourceClient
	cfg *config.Config
	log *slog.Logger
}

type Option func(*Service)

func WithClient(c GenResourceClient) Option { return func(s *Service) { s.c = c } }

func New(cfg *config.Config, httpClient *http.Client, log *slog.Logger, opts ...Option) (*Service, error) {
	c, err := camundav88.NewClientWithResponses(
		cfg.APIs.Camunda.BaseURL,
		camundav88.WithHTTPClient(httpClient),
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

func (s *Service) Deploy(ctx context.Context, unit []byte, opts ...services.CallOption) (d.Deployment, error) {
	_ = services.ApplyCallOptions(opts)
	resp, err := s.c.CreateDeploymentWithBodyWithResponse(ctx, "application/xml", bytes.NewReader(unit))
	if err != nil {
		return d.Deployment{}, err
	}
	if err = httpc.HttpStatusErr(resp.HTTPResponse, resp.Body); err != nil {
		return d.Deployment{}, err
	}
	if resp.JSON200 == nil {
		return d.Deployment{}, fmt.Errorf("%w: 200 OK but empty payload; body=%s",
			d.ErrMalformedResponse, string(resp.Body))
	}
	return fromDeploymentResult(*resp.JSON200), nil
}
