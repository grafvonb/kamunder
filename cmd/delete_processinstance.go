package cmd

import (
	"fmt"

	"github.com/grafvonb/kamunder/kamunder"
	"github.com/grafvonb/kamunder/toolx/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagDeletePIKey      string
	flagDeleteWithCancel bool
)

var deleteProcessInstanceCmd = &cobra.Command{
	Use:     "process-instance",
	Short:   "Delete a process instance by its key.",
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

		_, err = cli.DeleteProcessInstance(cmd.Context(), flagDeletePIKey, collectOptions()...)
		if err != nil {
			log.Error(fmt.Sprintf("deleteling process instance: %v", err))
			return
		}
	},
}

func init() {
	deleteCmd.AddCommand(deleteProcessInstanceCmd)

	AddBackoffFlagsAndBindings(deleteProcessInstanceCmd, viper.GetViper())

	deleteProcessInstanceCmd.Flags().StringVarP(&flagDeletePIKey, "key", "k", "", "process instance key to delete")
	_ = deleteProcessInstanceCmd.MarkFlagRequired("key")
	deleteProcessInstanceCmd.Flags().BoolVar(&flagDeleteWithCancel, "with-cancel", false, "cancel the process instance before deleting it")
}
