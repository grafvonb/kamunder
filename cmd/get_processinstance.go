package cmd

import (
	"fmt"

	"github.com/grafvonb/kamunder/kamunder"
	"github.com/grafvonb/kamunder/kamunder/process"
	"github.com/grafvonb/kamunder/toolx/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const maxPISearchSize int32 = 1000

var (
	flagPIKey               int64
	flagPIBpmnProcessID     string
	flagPIProcessVersion    int32
	flagPIProcessVersionTag string
	flagPIState             string
	flagPIParentKey         int64
)

// command options
var (
	flagPIParentsOnly       bool
	flagPIChildrenOnly      bool
	flagPIOrphanParentsOnly bool
	flagPIIncidentsOnly     bool
	flagPINoIncidentsOnly   bool
)

// view options
var (
	flagPIKeysOnly bool
)

var getProcessInstanceCmd = &cobra.Command{
	Use:     "process-instance",
	Short:   "Get process instances.",
	Aliases: []string{"process-instances", "pi", "pis"},
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

		log.Debug("fetching process instances")
		searchFilterOpts := populatePISearchFilterOpts()
		printFilter(cmd)
		if searchFilterOpts.Key > 0 {
			log.Debug(fmt.Sprintf("searching by key: %d", searchFilterOpts.Key))
			pi, err := cli.GetProcessInstanceByKey(cmd.Context(), searchFilterOpts.Key)
			if err != nil {
				log.Error(fmt.Sprintf("error fetching process instance by key %d: %v", searchFilterOpts.Key, err))
				return
			}
			err = processInstanceView(cmd, pi)
			if err != nil {
				log.Error(fmt.Sprintf("error rendering key-only view: %v", err))
				return
			}
			log.Debug(fmt.Sprintf("searched by key, found process instance with key: %d", pi.Key))
		} else {
			log.Debug(fmt.Sprintf("searching by filter: %v", searchFilterOpts))
			pisr, err := cli.SearchForProcessInstances(cmd.Context(), searchFilterOpts, maxPISearchSize)
			if err != nil {
				log.Error(fmt.Sprintf("error fetching process instances: %v", err))
				return
			}
			if flagPIChildrenOnly && flagPIParentsOnly {
				log.Error("using both --children-only and --parents-only filters returns always no results")
				return
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
					log.Error(fmt.Sprintf("error filtering orphan parents: %v", err))
					return
				}
			}
			if flagPIIncidentsOnly {
				pisr = pisr.FilterByHavingIncidents(true)
			}
			if flagPINoIncidentsOnly {
				pisr = pisr.FilterByHavingIncidents(false)
			}
			if flagPIKeysOnly {
				err = listKeyOnlyProcessInstancesView(cmd, pisr)
				if err != nil {
					log.Error(fmt.Sprintf("error rendering keys-only view: %v", err))
				}
				return
			}
			err = listProcessInstancesView(cmd, pisr)
			if err != nil {
				log.Error(fmt.Sprintf("error rendering items view: %v", err))
			}
			log.Debug(fmt.Sprintf("fetched process instances: %d", pisr.Total))
		}
	},
}

func init() {
	getCmd.AddCommand(getProcessInstanceCmd)

	AddBackoffFlagsAndBindings(getProcessInstanceCmd, viper.GetViper())

	fs := getProcessInstanceCmd.Flags()
	fs.Int64VarP(&flagPIKey, "key", "k", 0, "resource key (e.g. process instance) to fetch")
	fs.StringVarP(&flagPIBpmnProcessID, "bpmn-process-id", "b", "", "BPMN process ID to filter process instances")
	fs.Int32VarP(&flagPIProcessVersion, "process-version", "v", 0, "process definition version")
	fs.StringVar(&flagPIProcessVersionTag, "process-version-tag", "", "process definition version tag")

	// filtering options
	fs.Int64Var(&flagPIParentKey, "parent-key", 0, "parent process instance key to filter process instances")
	fs.StringVarP(&flagPIState, "state", "s", "all", "state to filter process instances: all, active, completed, canceled")
	fs.BoolVar(&flagPIParentsOnly, "parents-only", false, "show only parent process instances, meaning instances with no parent key set")
	fs.BoolVar(&flagPIChildrenOnly, "children-only", false, "show only child process instances, meaning instances that have a parent key set")
	fs.BoolVar(&flagPIOrphanParentsOnly, "orphan-parents-only", false, "show only child instances whose parent does not exist (return 404 on get by key)")
	fs.BoolVar(&flagPIIncidentsOnly, "incidents-only", false, "show only process instances that have incidents")
	fs.BoolVar(&flagPINoIncidentsOnly, "no-incidents-only", false, "show only process instances that have no incidents")

	// view options
	fs.BoolVar(&flagPIKeysOnly, "keys-only", false, "show only keys in output")
}

func populatePISearchFilterOpts() process.ProcessInstanceSearchFilterOpts {
	var filter process.ProcessInstanceSearchFilterOpts
	if flagPIKey != 0 {
		filter.Key = flagPIKey
	}
	if flagPIParentKey != 0 {
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
		state, err := process.ParseState(flagPIState)
		if err == nil {
			filter.State = state
		}
	}
	return filter
}
