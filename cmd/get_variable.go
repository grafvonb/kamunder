package cmd

import (
	"fmt"

	"github.com/grafvonb/kamunder/kamunder/ferrors"
	"github.com/spf13/cobra"
)

var getVariableCmd = &cobra.Command{
	Use:     "variable",
	Short:   "Get a variable by its name from a process instance",
	Aliases: []string{"var"},
	Run: func(cmd *cobra.Command, args []string) {
		cli, log, err := NewCli(cmd)
		if err != nil {
			ferrors.HandleAndExit(log, err)
		}

		log.Debug("getting variable")
		topology, err := cli.GetClusterTopology(cmd.Context())
		if err != nil {
			ferrors.HandleAndExit(log, fmt.Errorf("error getting variable: %w", err))
		}
		cmd.Println(ToJSONString(topology))
	},
}

func init() {
	getCmd.AddCommand(getVariableCmd)
}
