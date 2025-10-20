package cmd

import (
	"fmt"

	"github.com/grafvonb/kamunder/kamunder/ferrors"
	"github.com/grafvonb/kamunder/kamunder/process"
	"github.com/spf13/cobra"
)

const maxPISearchSize int32 = 1000

var (
	flagPIKey               string
	flagPIBpmnProcessID     string
	flagPIProcessVersion    int32
	flagPIProcessVersionTag string
	flagPIState             string
	flagPIParentKey         string
)

// command options
var (
	flagPIParentsOnly       bool
	flagPIChildrenOnly      bool
	flagPIOrphanParentsOnly bool
	flagPIIncidentsOnly     bool
	flagPINoIncidentsOnly   bool
)

var getProcessInstanceCmd = &cobra.Command{
	Use:     "process-instance",
	Short:   "Get process instances",
	Aliases: []string{"process-instances", "pi", "pis"},
	Run: func(cmd *cobra.Command, args []string) {
		cli, log, err := NewCli(cmd)
		if err != nil {
			ferrors.HandleAndExit(log, err)
		}

		if err != nil {
			ferrors.HandleAndExit(log, fmt.Errorf("error creating kamunder client: %w", err))
		}

		log.Debug(fmt.Sprintf("fetching process instances, render mode: %s", pickMode()))
		searchFilterOpts := populatePISearchFilterOpts()
		printFilter(cmd)
		if searchFilterOpts.Key != "" {
			log.Debug(fmt.Sprintf("searching by key: %s", searchFilterOpts.Key))
			pi, err := cli.GetProcessInstanceByKey(cmd.Context(), searchFilterOpts.Key)
			if err != nil {
				ferrors.HandleAndExit(log, fmt.Errorf("error fetching process instance by key %s: %w", searchFilterOpts.Key, err))
			}
			err = processInstanceView(cmd, pi)
			if err != nil {
				ferrors.HandleAndExit(log, fmt.Errorf("error rendering key-only view: %w", err))
			}
			log.Debug(fmt.Sprintf("searched by key, found process instance with key: %s", pi.Key))
		} else {
			log.Debug(fmt.Sprintf("searching by filter: %v", searchFilterOpts))
			pisr, err := cli.SearchForProcessInstances(cmd.Context(), searchFilterOpts, maxPISearchSize)
			if err != nil {
				ferrors.HandleAndExit(log, fmt.Errorf("error fetching process instances: %w", err))
			}
			if flagPIChildrenOnly && flagPIParentsOnly {
				ferrors.HandleAndExit(log, fmt.Errorf("%w: using both --children-only and --parents-only filters returns always no results", ferrors.ErrBadRequest))
			}
			if flagPIChildrenOnly {
				pisr = pisr.FilterChildrenOnly()
			}
			if flagPIParentsOnly {
				pisr = pisr.FilterParentsOnly()
			}
			if flagPIOrphanParentsOnly {
				pisr.Items, err = cli.FilterProcessInstanceWithOrphanParent(cmd.Context(), pisr.Items)
				if err != nil {
					ferrors.HandleAndExit(log, fmt.Errorf("error filtering orphan parents: %w", err))
				}
			}
			if flagPIIncidentsOnly {
				pisr = pisr.FilterByHavingIncidents(true)
			}
			if flagPINoIncidentsOnly {
				pisr = pisr.FilterByHavingIncidents(false)
			}
			err = listProcessInstancesView(cmd, pisr)
			if err != nil {
				ferrors.HandleAndExit(log, fmt.Errorf("error rendering items view: %w", err))
			}
			log.Debug(fmt.Sprintf("fetched process instances: %d", pisr.Total))
		}
	},
}

func init() {
	getCmd.AddCommand(getProcessInstanceCmd)

	fs := getProcessInstanceCmd.Flags()
	fs.StringVarP(&flagPIKey, "key", "k", "", "process instance key to fetch")
	fs.StringVarP(&flagPIBpmnProcessID, "bpmn-process-id", "b", "", "BPMN process ID to filter process instances")
	fs.Int32VarP(&flagPIProcessVersion, "process-version", "v", 0, "process definition version")
	fs.StringVar(&flagPIProcessVersionTag, "process-version-tag", "", "process definition version tag")

	// filtering options
	fs.StringVar(&flagPIParentKey, "parent-key", "", "parent process instance key to filter process instances")
	fs.StringVarP(&flagPIState, "state", "s", "all", "state to filter process instances: all, active, completed, canceled")
	fs.BoolVar(&flagPIParentsOnly, "parents-only", false, "show only parent process instances, meaning instances with no parent key set")
	fs.BoolVar(&flagPIChildrenOnly, "children-only", false, "show only child process instances, meaning instances that have a parent key set")
	fs.BoolVar(&flagPIOrphanParentsOnly, "orphan-parents-only", false, "show only child instances whose parent does not exist (return 404 on get by key)")
	fs.BoolVar(&flagPIIncidentsOnly, "incidents-only", false, "show only process instances that have incidents")
	fs.BoolVar(&flagPINoIncidentsOnly, "no-incidents-only", false, "show only process instances that have no incidents")
}

func populatePISearchFilterOpts() process.ProcessInstanceSearchFilterOpts {
	var filter process.ProcessInstanceSearchFilterOpts
	if flagPIKey != "" {
		filter.Key = flagPIKey
	}
	if flagPIParentKey != "" {
		filter.ParentKey = flagPIParentKey
	}
	if flagPIBpmnProcessID != "" {
		filter.BpmnProcessId = flagPIBpmnProcessID
	}
	if flagPIProcessVersion != 0 {
		filter.ProcessVersion = flagPIProcessVersion
	}
	if flagPIProcessVersionTag != "" {
		filter.ProcessVersionTag = flagPIProcessVersionTag
	}
	if flagPIState != "" && flagPIState != "all" {
		if state, ok := process.ParseState(flagPIState); ok {
			filter.State = state
		}
	}
	return filter
}
