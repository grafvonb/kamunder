package config

import (
	"fmt"
	"strings"
)

type HTTP struct {
	Timeout string `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
}

func (h *HTTP) Validate() error {
	if strings.TrimSpace(h.Timeout) == "" {
		return fmt.Errorf("timeout must not be empty")
	}
	return nil
}
