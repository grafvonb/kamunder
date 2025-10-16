package v88

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/grafvonb/kamunder/config"
	camundav88 "github.com/grafvonb/kamunder/internal/clients/camunda/v88/camunda"
	operatev88 "github.com/grafvonb/kamunder/internal/clients/camunda/v88/operate"
	d "github.com/grafvonb/kamunder/internal/domain"
	"github.com/grafvonb/kamunder/internal/services/httpc"
	"github.com/grafvonb/kamunder/internal/services/processinstance/state"
	"github.com/grafvonb/kamunder/toolx"
)

// nolint
type Service struct {
	cc  *camundav88.ClientWithResponses
	oc  *operatev88.ClientWithResponses
	cfg *config.Config
	log *slog.Logger
}

type Option func(*Service)

func New(cfg *config.Config, httpClient *http.Client, log *slog.Logger, opts ...Option) (*Service, error) {
	cc, err := camundav88.NewClientWithResponses(
		cfg.APIs.Camunda.BaseURL,
		camundav88.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}
	co, err := operatev88.NewClientWithResponses(
		cfg.APIs.Operate.BaseURL,
		operatev88.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}
	s := &Service{oc: co, cc: cc, cfg: cfg, log: log}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

func (s *Service) GetDirectChildrenOfProcessInstance(ctx context.Context, key int64) ([]d.ProcessInstance, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) FilterProcessInstanceWithOrphanParent(ctx context.Context, items []d.ProcessInstance) ([]d.ProcessInstance, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) SearchForProcessInstances(ctx context.Context, filter d.ProcessInstanceSearchFilterOpts, size int32) ([]d.ProcessInstance, error) {
	s.log.Debug(fmt.Sprintf("searching for process instances with filter: %+v", filter))
	st := operatev88.ProcessInstanceState(filter.State)
	f := operatev88.ProcessInstance{
		TenantId:          &s.cfg.App.Tenant,
		BpmnProcessId:     &filter.BpmnProcessId,
		ProcessVersion:    toolx.PtrIfNonZero(filter.ProcessVersion),
		ProcessVersionTag: &filter.ProcessVersionTag,
		State:             &st,
		ParentKey:         toolx.PtrIfNonZero(filter.ParentKey),
	}
	body := operatev88.SearchProcessInstancesJSONRequestBody{
		Filter: &f,
		Size:   &size,
	}
	resp, err := s.oc.SearchProcessInstancesWithResponse(ctx, body)
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
	return toolx.DerefSlicePtr(resp.JSON200.Items, fromProcessInstanceResponse), nil
}

func (s *Service) CancelProcessInstance(ctx context.Context, key int64) (d.CancelResponse, error) {
	s.log.Debug(fmt.Sprintf("checking if process instance with key %d is in allowable state to cancel", key))
	st, err := s.GetProcessInstanceStateByKey(ctx, key)
	if err != nil {
		return d.CancelResponse{}, err
	}
	if st.IsTerminal() {
		s.log.Info(fmt.Sprintf("process instance with key %d is already in state %s, no need to cancel", key, st))
		return d.CancelResponse{
			StatusCode: http.StatusOK,
			Status:     fmt.Sprintf("process instance with key %d is already in state %s, no need to cancel", key, st),
		}, nil
	}
	s.log.Debug(fmt.Sprintf("cancelling process instance with key %d", key))
	resp, err := s.cc.CancelProcessInstanceWithResponse(ctx, strconv.FormatInt(key, 10), camundav88.CancelProcessInstanceJSONRequestBody{})
	if err != nil {
		return d.CancelResponse{}, err
	}
	if err = httpc.HttpStatusErr(resp.HTTPResponse, resp.Body); err != nil {
		return d.CancelResponse{}, err
	}
	s.log.Info(fmt.Sprintf("process instance with key %d was successfully cancelled", key))
	return d.CancelResponse{
		StatusCode: resp.StatusCode(),
		Status:     resp.Status(),
	}, nil
}

func (s *Service) GetProcessInstanceStateByKey(ctx context.Context, key int64) (d.State, error) {
	s.log.Debug(fmt.Sprintf("checking state of process instance with key %d", key))
	pi, err := s.cc.GetProcessInstanceWithResponse(ctx, strconv.FormatInt(key, 10))
	if err != nil {
		return "", fmt.Errorf("fetching process instance with key %d: %w", key, err)
	}
	if err = httpc.HttpStatusErr(pi.HTTPResponse, pi.Body); err != nil {
		return "", fmt.Errorf("fetching process instance with key %d: %w", key, err)
	}
	if pi.JSON200 == nil {
		return "", fmt.Errorf("%w: 200 OK but empty payload; body=%s",
			d.ErrMalformedResponse, string(pi.Body))
	}
	st := d.State(pi.JSON200.State)
	s.log.Debug(fmt.Sprintf("process instance with key %d is in state %s", key, st))
	return st, nil
}

func (s *Service) DeleteProcessInstance(ctx context.Context, key int64) (d.ChangeStatus, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) DeleteProcessInstanceWithCancel(ctx context.Context, key int64) (d.ChangeStatus, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetProcessInstanceByKey(ctx context.Context, key int64) (d.ProcessInstance, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) WaitForProcessInstanceState(ctx context.Context, key int64, st d.State) error {
	return state.WaitForProcessInstanceState(ctx, s, s.cfg, s.log, key, st)
}
