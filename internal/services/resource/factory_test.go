package resource_test

import (
	"net/http"
	"testing"

	"log/slog"

	"github.com/grafvonb/kamunder/config"
	"github.com/grafvonb/kamunder/internal/services/resource"
	"github.com/grafvonb/kamunder/toolx"
	"github.com/stretchr/testify/require"
)

func testConfig() *config.Config {
	return &config.Config{
		APIs: config.APIs{},
	}
}

func TestFactory_V87(t *testing.T) {
	cfg := testConfig()
	cfg.APIs.Version = toolx.V87
	svc, err := resource.New(cfg, &http.Client{}, slog.Default())
	require.NoError(t, err)
	require.NotNil(t, svc)
}

func TestFactory_V88(t *testing.T) {
	cfg := testConfig()
	cfg.APIs.Version = toolx.V88
	svc, err := resource.New(cfg, &http.Client{}, slog.Default())
	require.NoError(t, err)
	require.NotNil(t, svc)
}

func TestFactory_Unknown(t *testing.T) {
	cfg := testConfig()
	cfg.APIs.Version = "v0"
	svc, err := resource.New(cfg, &http.Client{}, slog.Default())
	require.Error(t, err)
	require.Nil(t, svc)
	require.Contains(t, err.Error(), "unknown API version")
}
