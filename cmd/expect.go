package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var expectCmd = &cobra.Command{
	Use:     "expect",
	Short:   "Expect resources to be in a certain state",
	Aliases: []string{"e", "exp"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
	SuggestFor: []string{"expecte", "exepct"},
}

func init() {
	rootCmd.AddCommand(expectCmd)

	addBackoffFlagsAndBindings(expectCmd, viper.GetViper())
}
