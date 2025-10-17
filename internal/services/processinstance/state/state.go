package state

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/grafvonb/kamunder/config"
	d "github.com/grafvonb/kamunder/internal/domain"
	"github.com/grafvonb/kamunder/internal/services"
)

type PIGetter interface {
	GetProcessInstanceByKey(ctx context.Context, key string, opts ...services.CallOption) (d.ProcessInstance, error)
}

// WaitForProcessInstanceState waits until the instance reaches the desired state.
// - Respects ctx cancellation/deadline; augments with cfg.Timeout if set
// - Returns nil on success or an error on failure/timeout.
func WaitForProcessInstanceState(ctx context.Context, s PIGetter, cfg *config.Config, log *slog.Logger, key string, desiredState d.State, opts ...services.CallOption) error {
	_ = services.ApplyCallOptions(opts)

	backoff := cfg.App.Backoff
	if backoff.Timeout > 0 {
		deadline := time.Now().Add(backoff.Timeout)
		if dl, ok := ctx.Deadline(); !ok || deadline.Before(dl) {
			var cancel context.CancelFunc
			ctx, cancel = context.WithDeadline(ctx, deadline)
			defer cancel()
		}
	}

	attempts := 0
	delay := backoff.InitialDelay

	for {
		if errInDelay := ctx.Err(); errInDelay != nil {
			return errInDelay
		}
		attempts++

		pi, errInDelay := s.GetProcessInstanceByKey(ctx, key)
		if errInDelay == nil {
			if pi.State.EqualsIgnoreCase(desiredState) {
				log.Debug(fmt.Sprintf("process instance %s reached desired state %q", key, desiredState))
				return nil
			}
			log.Debug(fmt.Sprintf("process instance %s currently in state %q; waiting...", key, pi.State))
		} else if errInDelay != nil {
			if strings.Contains(errInDelay.Error(), "status 404") {
				log.Debug(fmt.Sprintf("process instance %s is absent (not found); waiting...", key))
			} else {
				log.Error(fmt.Sprintf("fetching state for %q failed: %v (will retry)", key, errInDelay))
			}
		}
		if backoff.MaxRetries > 0 && attempts >= backoff.MaxRetries {
			return fmt.Errorf("exceeded max_retries (%d) waiting for state %q", backoff.MaxRetries, desiredState)
		}
		select {
		case <-time.After(delay):
			delay = backoff.NextDelay(delay)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
