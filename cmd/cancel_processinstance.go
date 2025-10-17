package cmd

import (
	"fmt"

	"github.com/grafvonb/kamunder/kamunder"
	"github.com/grafvonb/kamunder/kamunder/options"
	"github.com/grafvonb/kamunder/toolx/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagCancelPIKey        int64
	flagCancelNoStateCheck bool
)

var cancelProcessInstanceCmd = &cobra.Command{
	Use:     "process-instance",
	Short:   "Cancel a process instance by its key.",
	Aliases: []string{"pi"},
	Run: func(cmd *cobra.Command, args []string) {
		log := logging.FromContext(cmd.Context())
		svcs, err := NewFromContext(cmd.Context())
		if err != nil {
			log.Error(fmt.Sprintf("%v", err))
			return
		}
		cli, err := kamunder.New(
			kamunder.WithConfig(svcs.Config),
			kamunder.WithHTTPClient(svcs.HTTP.Client()),
			kamunder.WithLogger(log),
		)
		if err != nil {
			log.Error(fmt.Sprintf("error creating kamunder client: %v", err))
			return
		}

		_, err = cli.CancelProcessInstance(cmd.Context(), flagCancelPIKey, collectOptions()...)
		if err != nil {
			log.Error(fmt.Sprintf("cancelling process instance: %v", err))
			return
		}
	},
}

func init() {
	cancelCmd.AddCommand(cancelProcessInstanceCmd)

	AddBackoffFlagsAndBindings(cancelProcessInstanceCmd, viper.GetViper())

	cancelProcessInstanceCmd.Flags().Int64VarP(&flagCancelPIKey, "key", "k", 0, "process instance key to cancel")
	_ = cancelProcessInstanceCmd.MarkFlagRequired("key")
	cancelProcessInstanceCmd.Flags().BoolVar(&flagCancelNoStateCheck, "no-state-check", false, "skip checking the current state of the process instance before cancelling it")
}

func collectOptions() []options.FacadeOption {
	var opts []options.FacadeOption
	if flagCancelNoStateCheck {
		opts = append(opts, options.WithNoStateCheck())
	}
	return opts
}
