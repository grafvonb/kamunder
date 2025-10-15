package cmd

import (
	"fmt"
	"strings"

	"github.com/grafvonb/camunder/kamunder"
	"github.com/grafvonb/camunder/kamunder/process"
	"github.com/grafvonb/camunder/toolx/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const maxSearchSize int32 = 1000

var supportedResourcesForGet = ResourceTypes{
	"ct": "cluster-topology",
	"pd": "process-definition",
	"pi": "process-instance",
}

// filter options
var (
	flagKey               int64
	flagBpmnProcessID     string
	flagProcessVersion    int32
	flagProcessVersionTag string
	flagState             string
	flagParentKey         int64
)

// command options
var (
	flagParentsOnly       bool
	flagChildrenOnly      bool
	flagOrphanParentsOnly bool
	flagIncidentsOnly     bool
	flagNoIncidentsOnly   bool
)

// view options
var (
	flagKeysOnly bool
	flagOneLine  bool
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:     "get [resource type]",
	Short:   "List resources of a resource type. " + supportedResourcesForGet.PrettyString(),
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"g", "list", "ls", "g"},
	Run: func(cmd *cobra.Command, args []string) {
		log := logging.FromContext(cmd.Context())
		rn := strings.ToLower(args[0])
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
		switch rn {
		case "cluster-topology", "ct":
			log.Debug("fetching cluster topology")
			topology, err := cli.GetClusterTopology(cmd.Context())
			if err != nil {
				log.Error(fmt.Sprintf("error fetching topology: %v", err))
				return
			}
			cmd.Println(ToJSONString(topology))

		case "process-definition", "pd":
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
				pdsr, err := cli.SearchProcessDefinitions(cmd.Context(), searchFilterOpts, maxSearchSize)
				if err != nil {
					log.Error(fmt.Sprintf("error fetching process definitions: %v", err))
					return
				}
				if flagKeysOnly {
					err = listKeyOnlyProcessDefinitionsView(cmd, pdsr)
					if err != nil {
						log.Error(fmt.Sprintf("error rendering keys-only view: %v", err))
					}
					return
				}
				err = listProcessDefinitionsView(cmd, pdsr)
				if err != nil {
					log.Error(fmt.Sprintf("error rendering items view: %v", err))
				}
			}

		case "process-instance", "pi":
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
				pisr, err := cli.SearchForProcessInstances(cmd.Context(), searchFilterOpts, maxSearchSize)
				if err != nil {
					log.Error(fmt.Sprintf("error fetching process instances: %v", err))
					return
				}
				if flagChildrenOnly && flagParentsOnly {
					log.Error("using both --children-only and --parents-only filters returns always no results")
					return
				}
				if flagChildrenOnly {
					pisr = pisr.FilterChildrenOnly()
				}
				if flagParentsOnly {
					pisr = pisr.FilterParentsOnly()
				}
				if flagOrphanParentsOnly {
					pisr.Items, err = cli.FilterProcessInstanceWithOrphanParent(cmd.Context(), pisr.Items)
					if err != nil {
						log.Error(fmt.Sprintf("error filtering orphan parents: %v", err))
						return
					}
				}
				if flagIncidentsOnly {
					pisr = pisr.FilterByHavingIncidents(true)
				}
				if flagNoIncidentsOnly {
					pisr = pisr.FilterByHavingIncidents(false)
				}
				if flagKeysOnly {
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

		default:
			log.Error(fmt.Sprintf("unknown resource type: %s, supported: %s", rn, supportedResourcesForGet))
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	AddBackoffFlagsAndBindings(getCmd, viper.GetViper())

	fs := getCmd.Flags()
	fs.Int64VarP(&flagKey, "key", "k", 0, "resource key (e.g. process instance) to fetch")
	fs.StringVarP(&flagBpmnProcessID, "bpmn-process-id", "b", "", "BPMN process ID to filter process instances")
	fs.Int32VarP(&flagProcessVersion, "process-version", "v", 0, "process definition version")
	fs.StringVar(&flagProcessVersionTag, "process-version-tag", "", "process definition version tag")

	// filtering options
	fs.Int64Var(&flagParentKey, "parent-key", 0, "parent process instance key to filter process instances")
	fs.StringVarP(&flagState, "state", "s", "all", "state to filter process instances: all, active, completed, canceled")
	fs.BoolVar(&flagParentsOnly, "parents-only", false, "show only parent process instances, meaning instances with no parent key set")
	fs.BoolVar(&flagChildrenOnly, "children-only", false, "show only child process instances, meaning instances that have a parent key set")
	fs.BoolVar(&flagOrphanParentsOnly, "orphan-parents-only", false, "show only child instances whose parent does not exist (return 404 on get by key)")
	fs.BoolVar(&flagIncidentsOnly, "incidents-only", false, "show only process instances that have incidents")
	fs.BoolVar(&flagNoIncidentsOnly, "no-incidents-only", false, "show only process instances that have no incidents")

	// view options
	fs.BoolVar(&flagKeysOnly, "keys-only", false, "show only keys in output")
	fs.BoolVar(&flagOneLine, "one-line", false, "output one line per item")
}

func populatePISearchFilterOpts() process.ProcessInstanceSearchFilterOpts {
	var filter process.ProcessInstanceSearchFilterOpts
	if flagKey != 0 {
		filter.Key = flagKey
	}
	if flagParentKey != 0 {
		filter.ParentKey = flagParentKey
	}
	if flagBpmnProcessID != "" {
		filter.BpmnProcessId = flagBpmnProcessID
	}
	if flagProcessVersion != 0 {
		filter.ProcessVersion = flagProcessVersion
	}
	if flagProcessVersionTag != "" {
		filter.ProcessVersionTag = flagProcessVersionTag
	}
	if flagState != "" && flagState != "all" {
		state, err := process.ParseState(flagState)
		if err == nil {
			filter.State = state
		}
	}
	return filter
}

func populatePDSearchFilterOpts() process.ProcessDefinitionSearchFilterOpts {
	var filter process.ProcessDefinitionSearchFilterOpts
	if flagKey != 0 {
		filter.Key = flagKey
	}
	if flagBpmnProcessID != "" {
		filter.BpmnProcessId = flagBpmnProcessID
	}
	if flagProcessVersion != 0 {
		filter.Version = flagProcessVersion
	}
	if flagProcessVersionTag != "" {
		filter.VersionTag = flagProcessVersionTag
	}
	return filter
}
