package cmd

import (
	"fmt"

	"github.com/grafvonb/kamunder/kamunder/ferrors"
	"github.com/spf13/cobra"
)

var (
	flagWalkKey  string
	flagWalkMode string
)

var validWalkModes = map[string]bool{
	"parent":   true,
	"children": true,
	"family":   true,
}

var walkProcessInstanceCmd = &cobra.Command{
	Use:     "process-instance",
	Short:   "Traverse (walk) the parent/child graph of process instances.",
	Aliases: []string{"pi", "pis"},
	Run: func(cmd *cobra.Command, args []string) {
		cli, log, err := NewCli(cmd)
		if err != nil {
			ferrors.HandleAndExit(log, err)
		}
		if !validWalkModes[flagWalkMode] {
			ferrors.HandleAndExit(log, fmt.Errorf("%w: invalid value for --walk: %q (must be parent, children, or family)", ferrors.ErrBadRequest, flagWalkMode))
		}

		var path KeysPath
		var chain Chain

		switch flagWalkMode {
		case "parent":
			_, path, chain, err = cli.Ancestry(cmd.Context(), flagWalkKey, collectOptions()...)
			if err != nil {
				ferrors.HandleAndExit(log, err)
			}
		case "children":
			path, _, chain, err = cli.Descendants(cmd.Context(), flagWalkKey, collectOptions()...)
			if err != nil {
				ferrors.HandleAndExit(log, err)
			}
		case "family":
			path, _, chain, err = cli.Family(cmd.Context(), flagWalkKey, collectOptions()...)
			if err != nil {
				ferrors.HandleAndExit(log, err)
			}
		default:
			ferrors.HandleAndExit(log, fmt.Errorf("%w: invalid value for --walk: %q (must be parent, children, or family)", ferrors.ErrBadRequest, flagWalkMode))
		}
		if flagPDKeysOnly {
			cmd.Println(path.KeysOnly(chain))
			return
		}
		cmd.Println(path.StandardLine(chain))
	},
}

func init() {
	walkCmd.AddCommand(walkProcessInstanceCmd)

	fs := walkProcessInstanceCmd.Flags()
	fs.StringVar(&flagWalkKey, "key", "", "start walking from this process instance key")
	_ = walkProcessInstanceCmd.MarkFlagRequired("key")
	fs.StringVar(&flagWalkMode, "mode", "parent", "walk mode: parent, children, family")
	_ = walkProcessInstanceCmd.MarkFlagRequired("mode")

	// view options
	fs.BoolVarP(&flagPDKeysOnly, "keys-only", "", false, "only print the keys of the resources")
}
