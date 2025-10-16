package cmd

import (
	"fmt"
	"strings"

	d "github.com/grafvonb/kamunder/internal/domain"
	"github.com/grafvonb/kamunder/internal/services/processinstance"
	"github.com/grafvonb/kamunder/internal/services/processinstance/state"
	"github.com/grafvonb/kamunder/toolx/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var supportedResourcesForExpect = ResourceTypes{
	"pi": "process-instance",
}

var (
	flagExpectKey int64
)

// expectCmd represents the cancel command
var expectCmd = &cobra.Command{
	Use:     "expect [resource name] [key]",
	Short:   "Expect a resource of a given type to change (e.g. its state) by its key. " + supportedResourcesForExpect.PrettyString(),
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"e", "exp", "await"},
	Run: func(cmd *cobra.Command, args []string) {
		log := logging.FromContext(cmd.Context())
		if err := requireAnyFlag(cmd, "state"); err != nil {
			log.Error(err.Error())
			return
		}
		rn := strings.ToLower(args[0])
		svcs, err := NewFromContext(cmd.Context())
		if err != nil {
			log.Error(fmt.Sprintf("Error initializing service from context: %v", err))
			return
		}

		switch rn {
		case "process-instance", "pi":
			svc, err := processinstance.New(svcs.Config, svcs.HTTP.Client(), log)
			if err != nil {
				log.Error(fmt.Sprintf("error creating process instance service: %v", err))
				return
			}

			st, ok := d.ParseState(flagPIState)
			if ok && st != d.StateAll {
				log.Info(fmt.Sprintf("waiting for process instance %d to reach state %q", flagExpectKey, st))
				err = state.WaitForProcessInstanceState(cmd.Context(), svc, svcs.Config, log, flagExpectKey, st)
				if err != nil {
					log.Error(fmt.Sprintf("error waiting for a process instance to reach a %q state: %v", st, err))
					return
				}
			}
		default:
			log.Error(fmt.Sprintf("unknown resource type %q", rn))
		}
	},
}

func init() {
	rootCmd.AddCommand(expectCmd)

	AddBackoffFlagsAndBindings(expectCmd, viper.GetViper())

	expectCmd.Flags().Int64VarP(&flagExpectKey, "key", "k", 0, "resource key (e.g. process instance)")
	_ = expectCmd.MarkFlagRequired("key")

	expectCmd.Flags().StringVarP(&flagPIState, "state", "s", "", "state of a process instance: active, completed, canceled or absent")
}
