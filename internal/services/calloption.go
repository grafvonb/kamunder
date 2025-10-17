package services

func WithNoStateCheck() CallOption { return func(c *CallCfg) { c.NoStateCheck = true } }
func WithCancel() CallOption       { return func(c *CallCfg) { c.WithCancel = true } }

type CallOption func(*CallCfg)

type CallCfg struct {
	NoStateCheck bool
	WithCancel   bool
}

func ApplyCallOptions(opts []CallOption) *CallCfg {
	c := &CallCfg{}
	for _, o := range opts {
		o(c)
	}
	return c
}
