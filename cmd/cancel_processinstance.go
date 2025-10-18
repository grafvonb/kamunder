package cmd

import (
	"fmt"

	"github.com/grafvonb/kamunder/kamunder/ferrors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagCancelPIKey        string
	flagCancelNoStateCheck bool
)

var cancelProcessInstanceCmd = &cobra.Command{
	Use:     "process-instance",
	Short:   "Cancel a process instance by its key.",
	Aliases: []string{"pi"},
	Run: func(cmd *cobra.Command, args []string) {
		cli, log, err := NewCli(cmd)
		if err != nil {
			ferrors.HandleAndExit(log, err)
		}

		_, err = cli.CancelProcessInstance(cmd.Context(), flagCancelPIKey, collectOptions()...)
		if err != nil {
			ferrors.HandleAndExit(log, fmt.Errorf("cancelling process instance: %w", err))
		}
	},
}

func init() {
	cancelCmd.AddCommand(cancelProcessInstanceCmd)

	AddBackoffFlagsAndBindings(cancelProcessInstanceCmd, viper.GetViper())

	cancelProcessInstanceCmd.Flags().StringVarP(&flagCancelPIKey, "key", "k", "", "process instance key to cancel")
	_ = cancelProcessInstanceCmd.MarkFlagRequired("key")
	cancelProcessInstanceCmd.Flags().BoolVar(&flagCancelNoStateCheck, "no-state-check", false, "skip checking the current state of the process instance before cancelling it")
}
