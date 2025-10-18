package cmd

import "github.com/grafvonb/kamunder/kamunder/options"

func collectOptions() []options.FacadeOption {
	var opts []options.FacadeOption
	if flagCancelWait {
		opts = append(opts, options.WithWait())
	}
	if flagCancelNoStateCheck {
		opts = append(opts, options.WithNoStateCheck())
	}
	if flagDeleteWithCancel {
		opts = append(opts, options.WithCancel())
	}
	return opts
}
