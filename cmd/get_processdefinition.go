package cmd

import (
	"fmt"

	"github.com/grafvonb/kamunder/kamunder/ferrors"
	"github.com/grafvonb/kamunder/kamunder/process"
	"github.com/spf13/cobra"
)

const maxPDSearchSize int32 = 1000

var (
	flagPDKey               string
	flagPDBpmnProcessID     string
	flagPDProcessVersion    int32
	flagPDProcessVersionTag string
)

var getProcessDefinitionCmd = &cobra.Command{
	Use:     "process-definition",
	Short:   "Get deployed process definitions",
	Aliases: []string{"processdefinition", "processdefinitions", "pd", "pds"},
	Run: func(cmd *cobra.Command, args []string) {
		cli, log, err := NewCli(cmd)
		if err != nil {
			ferrors.HandleAndExit(log, err)
		}

		log.Debug("fetching process definitions")
		searchFilterOpts := populatePDSearchFilterOpts()
		if searchFilterOpts.Key != "" {
			log.Debug(fmt.Sprintf("searching by key: %s", searchFilterOpts.Key))
			pd, err := cli.GetProcessDefinitionByKey(cmd.Context(), searchFilterOpts.Key)
			if err != nil {
				ferrors.HandleAndExit(log, fmt.Errorf("error fetching process definition by key %s: %w", searchFilterOpts.Key, err))
			}
			err = processDefinitionView(cmd, pd)
			if err != nil {
				ferrors.HandleAndExit(log, fmt.Errorf("error rendering key-only view: %w", err))
			}
		} else {
			log.Debug(fmt.Sprintf("searching by filter: %v", searchFilterOpts))
			pds, err := cli.SearchProcessDefinitions(cmd.Context(), searchFilterOpts, maxPDSearchSize)
			if err != nil {
				ferrors.HandleAndExit(log, fmt.Errorf("error fetching process definitions: %w", err))
			}
			err = listProcessDefinitionsView(cmd, pds)
			if err != nil {
				ferrors.HandleAndExit(log, fmt.Errorf("error rendering items view: %w", err))
			}
		}
	},
}

func init() {
	getCmd.AddCommand(getProcessDefinitionCmd)

	fs := getProcessDefinitionCmd.Flags()
	fs.StringVarP(&flagPDKey, "key", "k", "", "process definition key to fetch")
	fs.StringVarP(&flagPDBpmnProcessID, "bpmn-process-id", "b", "", "BPMN process ID to filter process instances")
	fs.Int32VarP(&flagPDProcessVersion, "process-version", "v", 0, "process definition version")
	fs.StringVar(&flagPDProcessVersionTag, "process-version-tag", "", "process definition version tag")
}

func populatePDSearchFilterOpts() process.ProcessDefinitionSearchFilterOpts {
	var filter process.ProcessDefinitionSearchFilterOpts
	if flagPDKey != "" {
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
