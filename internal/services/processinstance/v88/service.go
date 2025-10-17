package v88

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/grafvonb/kamunder/config"
	camundav88 "github.com/grafvonb/kamunder/internal/clients/camunda/v88/camunda"
	operatev88 "github.com/grafvonb/kamunder/internal/clients/camunda/v88/operate"
	d "github.com/grafvonb/kamunder/internal/domain"
	"github.com/grafvonb/kamunder/internal/services"
	"github.com/grafvonb/kamunder/internal/services/httpc"
	"github.com/grafvonb/kamunder/internal/services/processinstance/state"
	"github.com/grafvonb/kamunder/toolx"
)

const wrongStateMessage400 = "Process instances needs to be in one of the states [COMPLETED, CANCELED]"

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

func (s *Service) GetDirectChildrenOfProcessInstance(ctx context.Context, key string, opts ...services.CallOption) ([]d.ProcessInstance, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) FilterProcessInstanceWithOrphanParent(ctx context.Context, items []d.ProcessInstance, opts ...services.CallOption) ([]d.ProcessInstance, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) SearchForProcessInstances(ctx context.Context, filter d.ProcessInstanceSearchFilterOpts, size int32, opts ...services.CallOption) ([]d.ProcessInstance, error) {
	_ = services.ApplyCallOptions(opts)
	s.log.Debug(fmt.Sprintf("searching for process instances with filter: %+v", filter))
	st := operatev88.ProcessInstanceState(filter.State)
	pk, err := toolx.StringToInt64Ptr(filter.ParentKey)
	if err != nil {
		return nil, fmt.Errorf("parsing parent key %q to int64: %w", filter.ParentKey, err)
	}
	f := operatev88.ProcessInstance{
		TenantId:          &s.cfg.App.Tenant,
		BpmnProcessId:     &filter.BpmnProcessId,
		ProcessVersion:    toolx.PtrIfNonZero(filter.ProcessVersion),
		ProcessVersionTag: &filter.ProcessVersionTag,
		State:             &st,
		ParentKey:         pk,
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

func (s *Service) CancelProcessInstance(ctx context.Context, key string, opts ...services.CallOption) (d.CancelResponse, error) {
	cCfg := services.ApplyCallOptions(opts)
	if !cCfg.NoStateCheck {
		s.log.Debug(fmt.Sprintf("checking if process instance with key %s is in allowable state to cancel", key))
		st, err := s.GetProcessInstanceStateByKey(ctx, key)
		if err != nil {
			return d.CancelResponse{}, err
		}
		if st.IsTerminal() {
			s.log.Info(fmt.Sprintf("process instance with key %s is already in state %s, no need to cancel", key, st))
			return d.CancelResponse{
				StatusCode: http.StatusOK,
				Status:     fmt.Sprintf("process instance with key %s is already in state %s, no need to cancel", key, st),
			}, nil
		}
	} else {
		s.log.Debug("skipping process instance state check before cancellation as per call options")
	}
	s.log.Debug(fmt.Sprintf("cancelling process instance with key %s", key))
	resp, err := s.cc.CancelProcessInstanceWithResponse(ctx, key, camundav88.CancelProcessInstanceJSONRequestBody{})
	if err != nil {
		return d.CancelResponse{}, err
	}
	if err = httpc.HttpStatusErr(resp.HTTPResponse, resp.Body); err != nil {
		return d.CancelResponse{}, err
	}
	s.log.Info(fmt.Sprintf("process instance with key %s was successfully cancelled", key))
	return d.CancelResponse{
		StatusCode: resp.StatusCode(),
		Status:     resp.Status(),
	}, nil
}

func (s *Service) GetProcessInstanceStateByKey(ctx context.Context, key string, opts ...services.CallOption) (d.State, error) {
	_ = services.ApplyCallOptions(opts)
	s.log.Debug(fmt.Sprintf("checking state of process instance with key %s", key))
	pi, err := s.cc.GetProcessInstanceWithResponse(ctx, key)
	if err != nil {
		return "", fmt.Errorf("fetching process instance with key %s: %w", key, err)
	}
	if err = httpc.HttpStatusErr(pi.HTTPResponse, pi.Body); err != nil {
		return "", fmt.Errorf("fetching process instance with key %s: %w", key, err)
	}
	if pi.JSON200 == nil {
		return "", fmt.Errorf("%w: 200 OK but empty payload; body=%s",
			d.ErrMalformedResponse, string(pi.Body))
	}
	st := d.State(pi.JSON200.State)
	s.log.Debug(fmt.Sprintf("process instance with key %s is in state %s", key, st))
	return st, nil
}

func (s *Service) DeleteProcessInstance(ctx context.Context, key string, opts ...services.CallOption) (d.ChangeStatus, error) {
	cCfg := services.ApplyCallOptions(opts)
	oldKey, err := toolx.StringToInt64(key)
	if err != nil {
		return d.ChangeStatus{}, fmt.Errorf("parsing process instance key %q to int64: %w", key, err)
	}
	s.log.Debug(fmt.Sprintf("deleting process instance with key %d", oldKey))
	resp, err := s.oc.DeleteProcessInstanceAndAllDependantDataByKeyWithResponse(ctx, oldKey)
	if resp.StatusCode() == http.StatusBadRequest &&
		resp.ApplicationproblemJSON400 != nil &&
		*resp.ApplicationproblemJSON400.Message == wrongStateMessage400 {
		if cCfg.WithCancel {
			s.log.Info(fmt.Sprintf("process instance with key %s not in one of terminated states; cancelling it first", key))
			_, err = s.CancelProcessInstance(ctx, key)
			if err != nil {
				return d.ChangeStatus{}, fmt.Errorf("error cancelling process instance with key %s: %w", key, err)
			}
			s.log.Info(fmt.Sprintf("waiting for process instance with key %s to be cancelled by workflow engine...", key))
			if err = state.WaitForProcessInstanceState(ctx, s, s.cfg, s.log, key, d.StateCanceled, opts...); err != nil {
				return d.ChangeStatus{}, fmt.Errorf("waiting for canceled state failed for %s: %w", key, err)
			}
			s.log.Info(fmt.Sprintf("retrying deletion of process instance with key %d", oldKey))
			resp, err = s.oc.DeleteProcessInstanceAndAllDependantDataByKeyWithResponse(ctx, oldKey)
		}
	}
	if err != nil {
		return d.ChangeStatus{}, err
	}
	if err = httpc.HttpStatusErr(resp.HTTPResponse, resp.Body); err != nil {
		return d.ChangeStatus{}, err
	}
	s.log.Info(fmt.Sprintf("process instance with key %s was successfully deleted", key))
	return d.ChangeStatus{
		Deleted: toolx.Deref(resp.JSON200.Deleted, 0),
		Message: toolx.Deref(resp.JSON200.Message, ""),
	}, nil
}

func (s *Service) GetProcessInstanceByKey(ctx context.Context, key string, opts ...services.CallOption) (d.ProcessInstance, error) {
	_ = services.ApplyCallOptions(opts)

	resp, err := s.cc.GetProcessInstanceWithResponse(ctx, key)
	if err != nil {
		return d.ProcessInstance{}, err
	}
	if err = httpc.HttpStatusErr(resp.HTTPResponse, resp.Body); err != nil {
		return d.ProcessInstance{}, err
	}
	if resp.JSON200 == nil {
		return d.ProcessInstance{}, fmt.Errorf("%w: 200 OK but empty payload; body=%s",
			d.ErrMalformedResponse, string(resp.Body))
	}
	return fromProcessInstanceResult(*resp.JSON200), nil
}

func (s *Service) WaitForProcessInstanceState(ctx context.Context, key string, st d.State, opts ...services.CallOption) error {
	return state.WaitForProcessInstanceState(ctx, s, s.cfg, s.log, key, st, opts...)
}
