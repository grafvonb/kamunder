package cmd

import (
	"github.com/spf13/cobra"
)

var walkCmd = &cobra.Command{
	Use:     "walk",
	Short:   "Traverse (walk) the parent/child graph of resource type",
	Aliases: []string{"w", "traverse"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
	SuggestFor: []string{"walkk", "travers"},
}

func init() {
	rootCmd.AddCommand(walkCmd)
}
