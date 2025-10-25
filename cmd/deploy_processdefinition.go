package cmd

import (
	"fmt"
	"os"

	"github.com/grafvonb/kamunder/config"
	"github.com/grafvonb/kamunder/kamunder/ferrors"
	"github.com/spf13/cobra"
)

var (
	flagDeployPDFiles   []string
	flagDeployPDWithRun bool
)

var deployProcessDefinitionCmd = &cobra.Command{
	Use:     "process-definition",
	Short:   "Deploy a process definition",
	Aliases: []string{"pd"},
	Run: func(cmd *cobra.Command, args []string) {
		cli, log, err := NewCli(cmd)
		if err != nil {
			ferrors.HandleAndExit(log, err)
		}
		cfg, err := config.FromContext(cmd.Context())
		if err != nil {
			ferrors.HandleAndExit(log, err)
		}
		if err := validateFiles(flagDeployPDFiles); err != nil {
			ferrors.HandleAndExit(log, fmt.Errorf("validating files: %w", err))
		}
		res, err := loadResources(flagDeployPDFiles, os.Stdin)
		if err != nil {
			ferrors.HandleAndExit(log, fmt.Errorf("collecting resources: %w", err))
		}
		log.Debug(fmt.Sprintf("deploying process definition(s) to tenant %s", cfg.App.Tenant))
		_, err = cli.DeployProcessDefinition(cmd.Context(), cfg.App.Tenant, res, collectOptions()...)
		if err != nil {
			ferrors.HandleAndExit(log, fmt.Errorf("deploying process definition: %w", err))
		}
		log.Info("process definition deployed successfully")
	},
}

func init() {
	deployCmd.AddCommand(deployProcessDefinitionCmd)
	deployProcessDefinitionCmd.Flags().StringSliceVarP(&flagDeployPDFiles, "files", "f", nil, "paths to BPMN/YAML files or '-' for stdin")
	_ = deployProcessDefinitionCmd.MarkFlagRequired("files")

	deployProcessDefinitionCmd.Flags().BoolVar(&flagDeployPDWithRun, "with-run", false, "start a process instance after deploy")
}
