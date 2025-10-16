package httpc

import (
	"log/slog"
	"net/http"

	"github.com/grafvonb/kamunder/internal/services/auth/authenticator"
)

type LogTransport struct {
	base http.RoundTripper

	Log *slog.Logger
}

func (t *LogTransport) rt() http.RoundTripper {
	if t.base != nil {
		return t.base
	}
	if t.Log == nil {
		t.Log = slog.Default()
	}
	return http.DefaultTransport
}

func (t *LogTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.Log.Debug("calling: " + req.URL.String())
	return t.rt().RoundTrip(req)
}

type AuthTransport struct {
	base http.RoundTripper

	Editor authenticator.RequestEditor
}

func (t *AuthTransport) rt() http.RoundTripper {
	if t.base != nil {
		return t.base
	}
	return http.DefaultTransport
}

func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Editor != nil {
		if err := t.Editor(req.Context(), req); err != nil {
			return nil, err
		}
	}
	return t.rt().RoundTrip(req)
}
