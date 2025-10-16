package v88

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/grafvonb/kamunder/config"
	camundav88 "github.com/grafvonb/kamunder/internal/clients/camunda/v88/camunda"
	operatev88 "github.com/grafvonb/kamunder/internal/clients/camunda/v88/operate"
	d "github.com/grafvonb/kamunder/internal/domain"
	"github.com/grafvonb/kamunder/internal/services/processinstance/state"
)

// nolint
type Service struct {
	cc  *camundav88.ClientWithResponses
	oc  *operatev88.ClientWithResponses
	cfg *config.Config
	log *slog.Logger
}

func (s Service) GetProcessInstanceByKey(ctx context.Context, key int64) (d.ProcessInstance, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetDirectChildrenOfProcessInstance(ctx context.Context, key int64) ([]d.ProcessInstance, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) FilterProcessInstanceWithOrphanParent(ctx context.Context, items []d.ProcessInstance) ([]d.ProcessInstance, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) SearchForProcessInstances(ctx context.Context, filter d.ProcessInstanceSearchFilterOpts, size int32) ([]d.ProcessInstance, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) CancelProcessInstance(ctx context.Context, key int64) (d.CancelResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) DeleteProcessInstance(ctx context.Context, key int64) (d.ChangeStatus, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) DeleteProcessInstanceWithCancel(ctx context.Context, key int64) (d.ChangeStatus, error) {
	//TODO implement me
	panic("implement me")
}

type Option func(*Service)

func New(cfg *config.Config, httpClient *http.Client, log *slog.Logger, opts ...Option) (*Service, error) {
	return &Service{
		cc:  &camundav88.ClientWithResponses{},
		oc:  &operatev88.ClientWithResponses{},
		cfg: &config.Config{},
		log: &slog.Logger{},
	}, nil
}

func (s *Service) WaitForProcessInstanceState(ctx context.Context, key int64, st d.State) error {
	return state.WaitForProcessInstanceState(ctx, s, s.cfg, s.log, key, st)
}
