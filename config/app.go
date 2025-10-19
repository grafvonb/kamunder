package config

import "github.com/grafvonb/kamunder/internal/services/common"

type App struct {
	Tenant  string               `mapstructure:"tenant" json:"tenant" yaml:"tenant"`
	Backoff common.BackoffConfig `mapstructure:"backoff" json:"backoff" yaml:"backoff"`
}
