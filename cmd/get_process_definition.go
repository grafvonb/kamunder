package cmd

import (
	"fmt"

	"github.com/grafvonb/kamunder/kamunder"
	"github.com/grafvonb/kamunder/kamunder/process"
	"github.com/grafvonb/kamunder/toolx/logging"
	"github.com/spf13/cobra"
)

const maxPDSearchSize int32 = 1000

var (
	flagPDKey               int64
	flagPDBpmnProcessID     string
	flagPDProcessVersion    int32
	flagPDProcessVersionTag string
)

// view options
var (
	flagPDKeysOnly bool
)

var getProcessDefinitionCmd = &cobra.Command{
	Use:     "process-definition",
	Short:   "Get deployed process definitions.",
	Aliases: []string{"pd"},
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

		log.Debug("fetching process definitions")
		searchFilterOpts := populatePDSearchFilterOpts()
		if searchFilterOpts.Key > 0 {
			log.Debug(fmt.Sprintf("searching by key: %d", searchFilterOpts.Key))
			pd, err := cli.GetProcessDefinitionByKey(cmd.Context(), searchFilterOpts.Key)
			if err != nil {
				log.Error(fmt.Sprintf("error fetching process definition by key %d: %v", searchFilterOpts.Key, err))
				return
			}
			err = processDefinitionView(cmd, pd)
			if err != nil {
				log.Error(fmt.Sprintf("error rendering key-only view: %v", err))
				return
			}
		} else {
			log.Debug(fmt.Sprintf("searching by filter: %v", searchFilterOpts))
			pds, err := cli.SearchProcessDefinitions(cmd.Context(), searchFilterOpts, maxPDSearchSize)
			if err != nil {
				log.Error(fmt.Sprintf("error fetching process definitions: %v", err))
				return
			}
			if flagPDKeysOnly {
				err = listKeyOnlyProcessDefinitionsView(cmd, pds)
				if err != nil {
					log.Error(fmt.Sprintf("error rendering keys-only view: %v", err))
				}
				return
			}
			err = listProcessDefinitionsView(cmd, pds)
			if err != nil {
				log.Error(fmt.Sprintf("error rendering items view: %v", err))
			}
		}

	},
}

func init() {
	getCmd.AddCommand(getProcessDefinitionCmd)

	fs := getProcessDefinitionCmd.Flags()
	fs.Int64VarP(&flagPDKey, "key", "k", 0, "resource key (e.g. process instance) to fetch")
	fs.StringVarP(&flagPDBpmnProcessID, "bpmn-process-id", "b", "", "BPMN process ID to filter process instances")
	fs.Int32VarP(&flagPDProcessVersion, "process-version", "v", 0, "process definition version")
	fs.StringVar(&flagPDProcessVersionTag, "process-version-tag", "", "process definition version tag")

	// view options
	fs.BoolVar(&flagPDKeysOnly, "keys-only", false, "show only keys in output")
	fs.BoolVar(&flagOneLine, "one-line", false, "output one line per item")
}

func populatePDSearchFilterOpts() process.ProcessDefinitionSearchFilterOpts {
	var filter process.ProcessDefinitionSearchFilterOpts
	if flagPDKey != 0 {
		filter.Key = flagPDKey
	}
	if flagPDBpmnProcessID != "" {
		filter.BpmnProcessId = flagPDBpmnProcessID
	}
	if flagPDProcessVersion != 0 {
		filter.Version = flagPDProcessVersion
	}
	if flagPDProcessVersionTag != "" {
		filter.VersionTag = flagPDProcessVersionTag
	}
	return filter
}
