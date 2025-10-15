package auth

import (
	"fmt"
	"log/slog"
	"net/http"

	config2 "github.com/grafvonb/kamunder/config"
	"github.com/grafvonb/kamunder/internal/services/auth/cookie"
	"github.com/grafvonb/kamunder/internal/services/auth/core"
	"github.com/grafvonb/kamunder/internal/services/auth/oauth2"
)

func BuildAuthenticator(cfg *config2.Config, httpClient *http.Client, log *slog.Logger) (core.Authenticator, error) {
	switch cfg.Auth.Mode {
	case config2.ModeOAuth2, "":
		return oauth2.New(cfg, httpClient, log)
	case config2.ModeCookie:
		return cookie.New(cfg, httpClient, log)
	default:
		return nil, fmt.Errorf("unknown auth mode: %s", cfg.Auth.Mode)
	}
}
