package cmd

import (
	"github.com/spf13/cobra"
)

var (
	flagOneLine bool
)

var getCmd = &cobra.Command{
	Use:     "get",
	Short:   "Get resources",
	Aliases: []string{"g", "read"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.PersistentFlags().BoolVar(&flagOneLine, "one-line", false, "output one line per item")
}

//nolint:unused
func printOutput[T any](x T, linetoprint func(i T) string) {
	if flagAsJson {
		getCmd.Print(ToJSONString(x))
		return
	}
	getCmd.Println(linetoprint(x))
}

//nolint:unused
func printOutputs[T any](xs []T, linetoprint func(i T) string) {
	if flagAsJson {
		getCmd.Print(ToJSONString(xs))
		return
	}
	for _, x := range xs {
		getCmd.Println(linetoprint(x))
	}
}
