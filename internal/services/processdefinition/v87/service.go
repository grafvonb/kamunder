package v87

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/grafvonb/camunder/config"
	operatev87 "github.com/grafvonb/camunder/internal/clients/camunda/v87/operate"
	d "github.com/grafvonb/camunder/internal/domain"
	"github.com/grafvonb/camunder/internal/services/httpc"
	"github.com/grafvonb/camunder/toolx"
)

type Service struct {
	c   *operatev87.ClientWithResponses
	cfg *config.Config
	log *slog.Logger
}

type Option func(*Service)

func New(cfg *config.Config, httpClient *http.Client, log *slog.Logger, opts ...Option) (*Service, error) {
	c, err := operatev87.NewClientWithResponses(
		cfg.APIs.Operate.BaseURL,
		operatev87.WithHTTPClient(httpClient),
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

func (s *Service) GetProcessDefinitionByKey(ctx context.Context, key int64) (d.ProcessDefinition, error) {
	resp, err := s.c.GetProcessDefinitionByKeyWithResponse(ctx, key)
	if err != nil {
		return d.ProcessDefinition{}, err
	}
	if err = httpc.HttpStatusErr(resp.HTTPResponse, resp.Body); err != nil {
		return d.ProcessDefinition{}, err
	}
	if resp.JSON200 == nil {
		return d.ProcessDefinition{}, fmt.Errorf("%w: 200 OK but empty payload; body=%s",
			d.ErrMalformedResponse, string(resp.Body))
	}
	return fromProcessDefinitionResponse(*resp.JSON200), nil
}

func (s *Service) SearchProcessDefinitions(ctx context.Context, filter d.ProcessDefinitionSearchFilterOpts, size int32) ([]d.ProcessDefinition, error) {
	body := operatev87.QueryProcessDefinition{
		Filter: &operatev87.ProcessDefinition{
			BpmnProcessId: &filter.BpmnProcessId,
			Version:       toolx.PtrIfNonZero(filter.Version),
			VersionTag:    &filter.VersionTag,
		},
		Size: &size,
	}
	resp, err := s.c.SearchProcessDefinitionsWithResponse(ctx, body)
	if err != nil {
		return nil, err
	}
	if err = httpc.HttpStatusErr(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("%w: 200 OK but empty payload; body=%s",
			d.ErrMalformedResponse, string(resp.Body))
	}
	return toolx.DerefSlicePtr(resp.JSON200.Items, fromProcessDefinitionResponse), nil
}
