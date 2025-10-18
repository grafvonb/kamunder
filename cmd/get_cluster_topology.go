package cmd

import (
	"fmt"

	"github.com/grafvonb/kamunder/kamunder"
	"github.com/grafvonb/kamunder/kamunder/ferrors"
	"github.com/grafvonb/kamunder/toolx/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getClusterTopologyCmd = &cobra.Command{
	Use:     "cluster-topology",
	Short:   "Get the cluster topology of the connected Camunda 8 cluster.",
	Aliases: []string{"ct", "cluster-info", "ci"},
	Run: func(cmd *cobra.Command, args []string) {
		log := logging.FromContext(cmd.Context())
		svcs, err := NewFromContext(cmd.Context())
		if err != nil {
			ferrors.HandleAndExit(log, fmt.Errorf("error getting services from context: %w", err))
		}
		cli, err := kamunder.New(
			kamunder.WithConfig(svcs.Config),
			kamunder.WithHTTPClient(svcs.HTTP.Client()),
			kamunder.WithLogger(log),
		)
		if err != nil {
			ferrors.HandleAndExit(log, fmt.Errorf("error creating kamunder client: %w", err))
		}

		log.Debug("fetching cluster topology")
		topology, err := cli.GetClusterTopology(cmd.Context())
		if err != nil {
			ferrors.HandleAndExit(log, fmt.Errorf("error fetching topology: %w", err))
		}
		cmd.Println(ToJSONString(topology))
	},
}

func init() {
	getCmd.AddCommand(getClusterTopologyCmd)

	AddBackoffFlagsAndBindings(getCmd, viper.GetViper())
}
