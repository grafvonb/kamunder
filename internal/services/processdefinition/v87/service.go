package v87

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/grafvonb/kamunder/config"
	operatev87 "github.com/grafvonb/kamunder/internal/clients/camunda/v87/operate"
	d "github.com/grafvonb/kamunder/internal/domain"
	"github.com/grafvonb/kamunder/internal/services"
	"github.com/grafvonb/kamunder/internal/services/httpc"
	"github.com/grafvonb/kamunder/toolx"
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

func (s *Service) GetProcessDefinitionByKey(ctx context.Context, key string, opts ...services.CallOption) (d.ProcessDefinition, error) {
	_ = services.ApplyCallOptions(opts)
	oldKey, err := toolx.StringToInt64(key)
	if err != nil {
		return d.ProcessDefinition{}, fmt.Errorf("converting process definition key %q to int64: %w", key, err)
	}
	resp, err := s.c.GetProcessDefinitionByKeyWithResponse(ctx, oldKey)
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

func (s *Service) SearchProcessDefinitions(ctx context.Context, filter d.ProcessDefinitionSearchFilterOpts, size int32, opts ...services.CallOption) ([]d.ProcessDefinition, error) {
	_ = services.ApplyCallOptions(opts)
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
