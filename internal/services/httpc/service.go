package httpc

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/grafvonb/kamunder/config"
	"github.com/grafvonb/kamunder/internal/services/auth/authenticator"
)

var (
	ErrNoHttpServiceInContext  = errors.New("no http service in context")
	ErrInvalidServiceInContext = errors.New("invalid http service in context")
)

type Service struct {
	c   *http.Client
	cfg *config.Config
	log *slog.Logger
}

type Option func(*Service)

func WithTimeout(d time.Duration) Option {
	return func(s *Service) { s.c.Timeout = d }
}

func WithTimeoutString(v string) Option {
	return func(s *Service) {
		if v == "" {
			return
		}
		if d, err := time.ParseDuration(v); err == nil {
			s.c.Timeout = d
		}
	}
}

// WithCookieJar Ensure cookie jar (needed for XSRF)
func WithCookieJar() Option {
	return func(s *Service) { _ = s.InstallCookieJar() }
}

// WithAuthEditor Install an auth editor transport now
func WithAuthEditor(ed authenticator.RequestEditor) Option {
	return func(s *Service) { s.InstallAuthEditor(ed) }
}

func New(cfg *config.Config, log *slog.Logger, opts ...Option) (*Service, error) {
	if cfg == nil {
		return nil, errors.New("cfg is nil")
	}
	d, err := time.ParseDuration(cfg.HTTP.Timeout)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{Timeout: d, Transport: &logTransport{log: log}}
	s := &Service{c: httpClient, cfg: cfg, log: log}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

func (s *Service) Client() *http.Client { return s.c }

func (s *Service) UseClient(c *http.Client) { s.c = c }

func (s *Service) InstallCookieJar() error {
	if s.c.Jar != nil {
		return nil
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	s.c.Jar = jar
	return nil
}

func (s *Service) InstallAuthEditor(ed authenticator.RequestEditor) {
	s.c.Transport = &authTransport{base: s.c.Transport, editor: ed}
}

type authTransport struct {
	base   http.RoundTripper
	editor authenticator.RequestEditor
}

func (t *authTransport) rt() http.RoundTripper {
	if t.base != nil {
		return t.base
	}
	return http.DefaultTransport
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.editor != nil {
		if err := t.editor(req.Context(), req); err != nil {
			return nil, err
		}
	}
	return t.rt().RoundTrip(req)
}

type logTransport struct {
	base http.RoundTripper
	log  *slog.Logger
}

func (t *logTransport) rt() http.RoundTripper {
	if t.base != nil {
		return t.base
	}
	if t.log == nil {
		t.log = slog.Default()
	}
	return http.DefaultTransport
}

func (t *logTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.log.Debug("calling: " + req.URL.String())
	return t.rt().RoundTrip(req)
}

type ctxKey struct{}

func (s *Service) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey{}, s)
}

func FromContext(ctx context.Context) (*Service, error) {
	v := ctx.Value(ctxKey{})
	if v == nil {
		return nil, ErrNoHttpServiceInContext
	}
	s, ok := v.(*Service)
	if !ok || s == nil {
		return nil, ErrInvalidServiceInContext
	}
	return s, nil
}

func MustClient(ctx context.Context) *http.Client {
	if s, err := FromContext(ctx); err == nil && s != nil {
		return s.c
	}
	return http.DefaultClient
}
