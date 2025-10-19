package cmd

import "github.com/spf13/cobra"

var utilityCommands = map[string]struct{}{
	"help":       {},
	"version":    {},
	"completion": {},
	"config":     {},
}

func isUtilityCommand(cmd *cobra.Command) bool {
	if cmd == nil {
		return false
	}
	_, ok := utilityCommands[cmd.Name()]
	return ok
}

func hasHelpFlag(cmd *cobra.Command) bool {
	if cmd == nil {
		return false
	}
	return cmd.Flags().Changed("help")
}
