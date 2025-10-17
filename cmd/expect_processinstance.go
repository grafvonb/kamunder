package cmd

import (
	"fmt"

	"github.com/grafvonb/kamunder/kamunder"
	"github.com/grafvonb/kamunder/kamunder/process"
	"github.com/grafvonb/kamunder/toolx/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagExpectPIKey   string
	flagExpectPIState string
)

var expectProcessInstanceCmd = &cobra.Command{
	Use:     "process-instance",
	Short:   "Expect a process instance to reach a certain state.",
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
		st, ok := process.ParseState(flagExpectPIState)
		if ok && st != process.StateAll {
			log.Info(fmt.Sprintf("waiting for process instance %s to reach state %s", flagExpectPIKey, st))
			err = cli.WaitForProcessInstanceState(cmd.Context(), flagExpectPIKey, st, collectOptions()...)
			if err != nil {
				log.Error(fmt.Sprintf("error waiting for a process instance %s to reach a %s state: %v", flagCancelPIKey, st, err))
				return
			}
			log.Info(fmt.Sprintf("process instance %s reached desired state %s", flagCancelPIKey, st))
		} else {
			log.Error(fmt.Sprintf("invalid process instance state: %s", flagPIState))
			return
		}
	},
}

func init() {
	expectCmd.AddCommand(expectProcessInstanceCmd)

	AddBackoffFlagsAndBindings(expectProcessInstanceCmd, viper.GetViper())

	expectProcessInstanceCmd.Flags().StringVarP(&flagExpectPIKey, "key", "k", "", "process instance key to expect a state for")
	_ = expectProcessInstanceCmd.MarkFlagRequired("key")
	expectProcessInstanceCmd.Flags().StringVarP(&flagExpectPIState, "state", "s", "", "state of a process instance: ACTIVE, COMPLETED, CANCELED or ABSENT")
	_ = expectProcessInstanceCmd.MarkFlagRequired("state")
}
