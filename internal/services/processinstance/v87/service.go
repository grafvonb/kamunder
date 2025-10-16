package v87

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/grafvonb/kamunder/config"
	camundav87 "github.com/grafvonb/kamunder/internal/clients/camunda/v87/camunda"
	operatev87 "github.com/grafvonb/kamunder/internal/clients/camunda/v87/operate"
	d "github.com/grafvonb/kamunder/internal/domain"
	"github.com/grafvonb/kamunder/internal/services/httpc"
	"github.com/grafvonb/kamunder/internal/services/processinstance/state"
	"github.com/grafvonb/kamunder/toolx"
)

const wrongStateMessage400 = "Process instances needs to be in one of the states [COMPLETED, CANCELED]"

type Service struct {
	cc  *camundav87.ClientWithResponses
	oc  *operatev87.ClientWithResponses
	cfg *config.Config
	log *slog.Logger
}

type Option func(*Service)

func New(cfg *config.Config, httpClient *http.Client, log *slog.Logger, opts ...Option) (*Service, error) {
	cc, err := camundav87.NewClientWithResponses(
		cfg.APIs.Camunda.BaseURL,
		camundav87.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}
	co, err := operatev87.NewClientWithResponses(
		cfg.APIs.Operate.BaseURL,
		operatev87.WithHTTPClient(httpClient),
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

func (s *Service) GetProcessInstanceByKey(ctx context.Context, key int64) (d.ProcessInstance, error) {
	resp, err := s.oc.GetProcessInstanceByKeyWithResponse(ctx, key)
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
	return fromProcessInstanceResponse(*resp.JSON200), nil
}

func (s *Service) GetDirectChildrenOfProcessInstance(ctx context.Context, key int64) ([]d.ProcessInstance, error) {
	filter := d.ProcessInstanceSearchFilterOpts{
		ParentKey: key,
	}
	resp, err := s.SearchForProcessInstances(ctx, filter, 1000)
	if err != nil {
		return nil, fmt.Errorf("searching for children of process instance with key %d: %w", key, err)
	}
	return resp, nil
}

func (s *Service) FilterProcessInstanceWithOrphanParent(ctx context.Context, items []d.ProcessInstance) ([]d.ProcessInstance, error) {
	if items == nil {
		return nil, nil
	}
	var result []d.ProcessInstance
	for _, it := range items {
		if it.ParentKey == 0 {
			continue
		}
		_, err := s.GetProcessInstanceByKey(ctx, it.ParentKey)
		if err != nil && strings.Contains(err.Error(), "status 404") {
			result = append(result, it)
		} else if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (s *Service) SearchForProcessInstances(ctx context.Context, filter d.ProcessInstanceSearchFilterOpts, size int32) ([]d.ProcessInstance, error) {
	s.log.Debug(fmt.Sprintf("searching for process instances with filter: %+v", filter))
	st := operatev87.ProcessInstanceState(filter.State)
	f := operatev87.ProcessInstance{
		TenantId:          &s.cfg.App.Tenant,
		BpmnProcessId:     &filter.BpmnProcessId,
		ProcessVersion:    toolx.PtrIfNonZero(filter.ProcessVersion),
		ProcessVersionTag: &filter.ProcessVersionTag,
		State:             &st,
		ParentKey:         toolx.PtrIfNonZero(filter.ParentKey),
	}
	body := operatev87.SearchProcessInstancesJSONRequestBody{
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
	resp, err := s.cc.PostProcessInstancesProcessInstanceKeyCancellationWithResponse(ctx, strconv.Itoa(int(key)),
		camundav87.PostProcessInstancesProcessInstanceKeyCancellationJSONRequestBody{})
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
	pi, err := s.oc.GetProcessInstanceByKeyWithResponse(ctx, key)
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
	st := d.State(*pi.JSON200.State)
	s.log.Debug(fmt.Sprintf("process instance with key %d is in state %s", key, st))
	return st, nil
}

func (s *Service) DeleteProcessInstance(ctx context.Context, key int64) (d.ChangeStatus, error) {
	s.log.Debug(fmt.Sprintf("trying to delete process instance with key %d...", key))
	resp, err := s.oc.DeleteProcessInstanceAndAllDependantDataByKeyWithResponse(ctx, key)
	if err != nil {
		return d.ChangeStatus{}, err
	}
	if resp.StatusCode() != http.StatusOK {
		return d.ChangeStatus{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	s.log.Info(fmt.Sprintf("process instance with key %d was successfully deleted", key))
	return d.ChangeStatus{
		Deleted: toolx.Deref(resp.JSON200.Deleted, 0),
		Message: toolx.Deref(resp.JSON200.Message, ""),
	}, nil
}

func (s *Service) DeleteProcessInstanceWithCancel(ctx context.Context, key int64) (d.ChangeStatus, error) {
	s.log.Debug(fmt.Sprintf("trying to delete process instance with key %d...", key))
	resp, err := s.oc.DeleteProcessInstanceAndAllDependantDataByKeyWithResponse(ctx, key)
	if resp.StatusCode() == http.StatusBadRequest &&
		resp.ApplicationproblemJSON400 != nil &&
		*resp.ApplicationproblemJSON400.Message == wrongStateMessage400 {
		s.log.Info(fmt.Sprintf("process instance with key %d not in state COMPLETED or CANCELED, cancelling it first...", key))
		_, err = s.CancelProcessInstance(ctx, key)
		if err != nil {
			return d.ChangeStatus{}, fmt.Errorf("error cancelling process instance with key %d: %w", key, err)
		}
		s.log.Info(fmt.Sprintf("waiting for process instance with key %d to be cancelled by workflow engine...", key))
		if err = state.WaitForProcessInstanceState(ctx, s, s.cfg, s.log, key, d.StateCanceled); err != nil {
			return d.ChangeStatus{}, fmt.Errorf("waiting for canceled state failed for %d: %w", key, err)
		}
		resp, err = s.oc.DeleteProcessInstanceAndAllDependantDataByKeyWithResponse(ctx, key)
	}
	if err != nil {
		return d.ChangeStatus{}, err
	}
	if resp.StatusCode() != http.StatusOK {
		return d.ChangeStatus{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	s.log.Info(fmt.Sprintf("process instance with key %d was successfully deleted", key))
	return d.ChangeStatus{
		Deleted: toolx.Deref(resp.JSON200.Deleted, 0),
		Message: toolx.Deref(resp.JSON200.Message, ""),
	}, nil
}

func (s *Service) WaitForProcessInstanceState(ctx context.Context, key int64, st d.State) error {
	return state.WaitForProcessInstanceState(ctx, s, s.cfg, s.log, key, st)
}
