package cmd

import (
	"fmt"

	"github.com/grafvonb/kamunder/kamunder"
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

		log.Debug("fetching cluster topology")
		topology, err := cli.GetClusterTopology(cmd.Context())
		if err != nil {
			log.Error(fmt.Sprintf("error fetching topology: %v", err))
			return
		}
		cmd.Println(ToJSONString(topology))
	},
}

func init() {
	getCmd.AddCommand(getClusterTopologyCmd)

	AddBackoffFlagsAndBindings(getCmd, viper.GetViper())
}
