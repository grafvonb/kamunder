package options

import "github.com/grafvonb/kamunder/internal/services"

func WithNoStateCheck() FacadeOption { return func(c *FacadeCfg) { c.NoStateCheck = true } }
func WithCancel() FacadeOption       { return func(c *FacadeCfg) { c.WithCancel = true } }

type FacadeOption func(*FacadeCfg)

type FacadeCfg struct {
	NoStateCheck bool
	WithCancel   bool
}

func ApplyFacadeOptions(opts []FacadeOption) *FacadeCfg {
	c := &FacadeCfg{}
	for _, o := range opts {
		o(c)
	}
	return c
}

func MapFacadeOptionsToCallOptions(opts []FacadeOption) []services.CallOption {
	c := ApplyFacadeOptions(opts)
	var out []services.CallOption
	if c.NoStateCheck {
		out = append(out, services.WithNoStateCheck())
	}
	if c.WithCancel {
		out = append(out, services.WithCancel())
	}
	return out
}
