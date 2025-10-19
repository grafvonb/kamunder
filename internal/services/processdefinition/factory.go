package processdefinition

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/grafvonb/kamunder/config"
	"github.com/grafvonb/kamunder/internal/services"
	v87 "github.com/grafvonb/kamunder/internal/services/processdefinition/v87"
	v88 "github.com/grafvonb/kamunder/internal/services/processdefinition/v88"
	"github.com/grafvonb/kamunder/toolx"
)

func New(cfg *config.Config, httpClient *http.Client, log *slog.Logger) (API, error) {
	v := cfg.APIs.Version
	switch v {
	case toolx.V88:
		return v88.New(cfg, httpClient, log)
	case toolx.V87:
		return v87.New(cfg, httpClient, log)
	default:
		return nil, fmt.Errorf("%w: %q (supported: %v)", services.ErrUnknownAPIVersion, v, toolx.SupportedCamundaVersionsString())
	}
}
