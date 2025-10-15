//go:build integration

package cookie_test

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	config2 "github.com/grafvonb/camunder/config"
	"github.com/grafvonb/camunder/internal/services/auth/cookie"
	"github.com/grafvonb/camunder/internal/testx"
	"github.com/stretchr/testify/require"
)

func TestCookie_Login_OK_IT(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode")
	}

	cfg := &config2.Config{
		Auth: config2.Auth{
			Mode: config2.ModeCookie,
			Cookie: config2.AuthCookieSession{
				BaseURL:  testx.RequireEnvWithPrefix(t, "COOKIE_BASE_URL"),
				Username: testx.RequireEnvWithPrefix(t, "COOKIE_USERNAME"),
				Password: testx.RequireEnvWithPrefix(t, "COOKIE_PASSWORD"),
			},
		},
	}

	httpClient := &http.Client{Timeout: 15 * time.Second}
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	svc, err := cookie.New(cfg, httpClient, log)
	require.NoError(t, err)
	require.NotNil(t, svc)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	t.Logf("trying to authenticate aginst %s with user %q", cfg.Auth.Cookie.BaseURL, cfg.Auth.Cookie.Username)
	err = svc.Init(ctx)
	require.NoError(t, err)
	require.True(t, svc.IsAuthenticated())
	t.Log("success: got authenticated")
}
