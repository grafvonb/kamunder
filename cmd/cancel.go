package cmd

import (
	"github.com/spf13/cobra"
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
}
