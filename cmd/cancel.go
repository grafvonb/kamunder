package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagCancelWait bool
)

var cancelCmd = &cobra.Command{
	Use:     "cancel",
	Short:   "Cancel resources",
	Aliases: []string{"c", "cn", "stop", "abort"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
	SuggestFor: []string{"cancle", "cancl"},
}

func init() {
	rootCmd.AddCommand(cancelCmd)

	cancelCmd.PersistentFlags().BoolVarP(&flagCancelWait, "wait", "w", false, "wait for the cancellation to be completed")

	addBackoffFlagsAndBindings(cancelCmd, viper.GetViper())
}
