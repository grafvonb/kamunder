package cmd

import (
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete resources",
	Aliases: []string{"d", "del", "remove", "rm"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
	SuggestFor: []string{"deelte", "delet"},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
