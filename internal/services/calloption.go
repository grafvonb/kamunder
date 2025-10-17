package services

func WithNoStateCheck() CallOption { return func(c *CallCfg) { c.NoStateCheck = true } }
func WithForce() CallOption        { return func(c *CallCfg) { c.Force = true } }

type CallOption func(*CallCfg)

type CallCfg struct {
	NoStateCheck bool
	Force        bool
}

func ApplyCallOptions(opts []CallOption) *CallCfg {
	c := &CallCfg{}
	for _, o := range opts {
		o(c)
	}
	return c
}
