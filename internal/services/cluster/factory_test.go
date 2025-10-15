package cluster

import (
	"net/http"
	"testing"

	"log/slog"

	config2 "github.com/grafvonb/kamunder/config"
	"github.com/grafvonb/kamunder/toolx"
	"github.com/stretchr/testify/require"
)

func testConfig() *config2.Config {
	return &config2.Config{
		APIs: config2.APIs{},
	}
}

func TestFactory_V87(t *testing.T) {
	cfg := testConfig()
	cfg.APIs.Version = toolx.V87
	svc, err := New(cfg, &http.Client{}, slog.Default())
	require.NoError(t, err)
	require.NotNil(t, svc)
}

func TestFactory_V88(t *testing.T) {
	cfg := testConfig()
	cfg.APIs.Version = toolx.V88
	svc, err := New(cfg, &http.Client{}, slog.Default())
	require.NoError(t, err)
	require.NotNil(t, svc)
}

func TestFactory_Unknown(t *testing.T) {
	cfg := testConfig()
	cfg.APIs.Version = "v0"
	svc, err := New(cfg, &http.Client{}, slog.Default())
	require.Error(t, err)
	require.Nil(t, svc)
	require.Contains(t, err.Error(), "unknown API version")
}
