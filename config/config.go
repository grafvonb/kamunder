package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

var (
	ErrNoBaseURL      = errors.New("no base_url provided in api configuration")
	ErrNoTokenURL     = errors.New("no token_url provided in auth configuration")
	ErrNoClientID     = errors.New("no client_id provided in auth configuration")
	ErrNoClientSecret = errors.New("no client_secret provided in auth configuration")

	ErrNoConfigInContext       = errors.New("no config in context")
	ErrInvalidServiceInContext = errors.New("invalid config in context")
)

type Config struct {
	Config string `mapstructure:"config" json:"config" yaml:"config"`

	App  App  `mapstructure:"app" json:"app" yaml:"app"`
	Auth Auth `mapstructure:"auth" json:"auth" yaml:"auth"`
	APIs APIs `mapstructure:"apis" json:"apis" yaml:"apis"`
	HTTP HTTP `mapstructure:"http" json:"http" yaml:"http"`
}

func (c *Config) Normalize() error {
	var errs []error
	if err := c.APIs.Normalize(); err != nil {
		errs = append(errs, fmt.Errorf("apis:\n%w", err))
	}
	return errors.Join(errs...)
}

func (c *Config) Validate() error {
	var errs []error
	if err := c.Auth.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("auth:\n%w", err))
	}
	if err := c.APIs.Validate(c.Auth.OAuth2.Scopes); err != nil {
		errs = append(errs, fmt.Errorf("apis:\n%w", err))
	}
	if err := c.HTTP.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("http:\n%w", err))
	}
	return errors.Join(errs...)
}

func (c *Config) ToSanitizedYAML() (string, error) {
	return c.toYaml(yamlExportOptions{
		template: false,
		sanitizeKeys: []string{
			"client_secret",
			"password",
			"token",
		},
	})
}

func (c *Config) ToTemplateYAML() (string, error) {
	return c.toYaml(yamlExportOptions{
		template:     true,
		sanitizeKeys: nil,
	})
}

type yamlExportOptions struct {
	template     bool     // if true: blank all leaf values
	sanitizeKeys []string // mask these keys with "***"
}

func (c *Config) toYaml(opts yamlExportOptions) (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return "", err
	}

	if len(opts.sanitizeKeys) > 0 {
		sanitize(m, opts.sanitizeKeys)
	}
	if opts.template {
		blankAllLeaves(m)
	}

	out, err := yaml.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func blankAllLeaves(m map[string]any) {
	for k, v := range m {
		switch x := v.(type) {
		case map[string]any:
			blankAllLeaves(x)
		case []any:
			m[k] = []any{}
		default:
			m[k] = ""
		}
	}
}

func sanitize(m map[string]any, sensitive []string) {
	for k, v := range m {
		if isSensitive(k, sensitive) {
			m[k] = "*****"
			continue
		}
		switch x := v.(type) {
		case map[string]any:
			sanitize(x, sensitive)
		case []any:
			for _, e := range x {
				if sub, ok := e.(map[string]any); ok {
					sanitize(sub, sensitive)
				}
			}
		}
	}
}

func isSensitive(k string, sensitive []string) bool {
	for _, s := range sensitive {
		if k == s {
			return true
		}
	}
	return false
}

type ctxConfigKey struct{}

func (c *Config) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxConfigKey{}, c)
}

func FromContext(ctx context.Context) (*Config, error) {
	v := ctx.Value(ctxConfigKey{})
	if v == nil {
		return nil, ErrNoConfigInContext
	}
	c, ok := v.(*Config)
	if !ok || c == nil {
		return nil, ErrInvalidServiceInContext
	}
	return c, nil
}
