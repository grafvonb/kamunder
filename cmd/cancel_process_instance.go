package cmd

import (
	"fmt"

	"github.com/grafvonb/kamunder/internal/services/processinstance"
	"github.com/grafvonb/kamunder/toolx/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagCancelPIKey int64
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

		svc, err := processinstance.New(svcs.Config, svcs.HTTP.Client(), log)
		if err != nil {
			log.Error(fmt.Sprintf("creating process instance service: %v", err))
			return
		}
		_, err = svc.CancelProcessInstance(cmd.Context(), flagCancelPIKey)
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
}
