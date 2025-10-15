package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	config2 "github.com/grafvonb/kamunder/config"
	"github.com/grafvonb/kamunder/internal/services/auth"
	authcore "github.com/grafvonb/kamunder/internal/services/auth/core"
	"github.com/grafvonb/kamunder/internal/services/httpc"
	"github.com/grafvonb/kamunder/toolx"
	"github.com/grafvonb/kamunder/toolx/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagShowConfig bool // show effective config and exit
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kamunder",
	Short: "Kamunder is a CLI tool to interact with Camunda 8.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		v := viper.New()
		if err := initViper(v, cmd); err != nil {
			return err
		}
		// retrieve and validate config
		cfg, err := retrieveConfig(v)
		if err != nil {
			return err
		}
		cmd.SetContext(cfg.ToContext(cmd.Context()))

		if flagShowConfig {
			cfgpath := v.ConfigFileUsed()
			if cfgpath != "" {
				cmd.Println("config loaded:", cfgpath)
			}
			cmd.Println(cfg.String())
			os.Exit(0)
		}
		// Setup logger
		log := logging.New(logging.LoggerConfig{
			Level:      v.GetString("log.level"),
			Format:     v.GetString("log.format"),
			WithSource: v.GetBool("log.with_source"),
		})
		cmd.SetContext(logging.ToContext(cmd.Context(), log))

		if cmd.Name() == "help" || cmd.Name() == "version" || cmd.Name() == "completion" {
			return nil
		}
		if cmd.Flags().Changed("help") {
			return nil
		}

		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("validate config: %w", err)
		}

		httpSvc, err := httpc.New(cfg, log, httpc.WithCookieJar())
		if err != nil {
			return fmt.Errorf("http service: %w", err)
		}
		authenticator, err := auth.BuildAuthenticator(cfg, httpSvc.Client(), log)
		if err != nil {
			return fmt.Errorf("auth build: %w", err)
		}
		if err := authenticator.Init(cmd.Context()); err != nil {
			return fmt.Errorf("auth init: %w", err)
		}
		httpSvc.InstallAuthEditor(authenticator.Editor())

		ctx := httpSvc.ToContext(cmd.Context())
		ctx = authcore.ToContext(ctx, authenticator)
		cmd.SetContext(ctx)

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
		// return runUI(cmd, args)
	},
	SilenceUsage:  true,
	SilenceErrors: false,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	pf := rootCmd.PersistentFlags()

	pf.String("config", "", "path to config file")

	pf.String("log-level", "info", "log level (debug, info, warn, error)")
	pf.String("log-format", "plain", "log format (json, plain, text)")
	pf.Bool("log-with-source", false, "include source file and line number in logs")

	pf.String("tenant", "", "default tenant ID")

	pf.String("auth-token-url", "", "auth token URL")
	pf.String("auth-client-id", "", "auth client ID")
	pf.String("auth-client-secret", "", "auth client secret")
	pf.StringToString("auth-scopes", nil, "auth scopes as key=value (repeatable or comma-separated)")

	pf.String("http-timeout", "", "HTTP timeout (Go duration, e.g. 30s)")

	pf.StringP("camunda-apis-version", "a", string(toolx.Current), fmt.Sprintf("Camunda API version (supported: %v)", toolx.Supported()))
	pf.String("camunda-base-url", "", "Camunda API base URL")
	pf.String("operate-base-url", "", "Operate API base URL")
	pf.String("tasklist-base-url", "", "Tasklist API base URL")

	// TODO show-config flag should be in a "config" subcommand
	pf.BoolVar(&flagShowConfig, "show-config", false, "print effective config (secrets redacted)")

	// TODO add --dry-run flag to commands that perform actions
}

func initViper(v *viper.Viper, cmd *cobra.Command) error {
	// Resolve precedence: flags > env > config file > defaults
	fs := cmd.Flags()
	_ = v.BindPFlag("config", fs.Lookup("config"))

	_ = v.BindPFlag("log.level", fs.Lookup("log-level"))
	_ = v.BindPFlag("log.format", fs.Lookup("log-format"))
	_ = v.BindPFlag("log.with_source", fs.Lookup("log-with-source"))

	_ = v.BindPFlag("app.tenant", fs.Lookup("tenant"))
	_ = v.BindPFlag("auth.token_url", fs.Lookup("auth-token-url"))
	_ = v.BindPFlag("auth.client_id", fs.Lookup("auth-client-id"))
	_ = v.BindPFlag("auth.client_secret", fs.Lookup("auth-client-secret"))
	_ = v.BindPFlag("http.timeout", fs.Lookup("http-timeout"))

	_ = v.BindPFlag("apis.version", fs.Lookup("camunda-apis-version"))
	_ = v.BindPFlag("apis.camunda_api.base_url", fs.Lookup("camunda-base-url"))
	_ = v.BindPFlag("apis.operate_api.base_url", fs.Lookup("operate-base-url"))
	_ = v.BindPFlag("apis.tasklist_api.base_url", fs.Lookup("tasklist-base-url"))

	_ = v.BindPFlag("tmp.auth_scopes", fs.Lookup("auth-scopes"))

	// Force hardcoded keys
	v.Set("apis.camunda_api.key", config2.CamundaApiKeyConst)
	v.Set("apis.operate_api.key", config2.OperateApiKeyConst)
	v.Set("apis.tasklist_api.key", config2.TasklistApiKeyConst)

	// Defaults
	v.SetDefault("http.timeout", "30s")

	// Config file discovery
	if cfgFile := v.GetString("config"); cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")

		// Search config paths (in order):
		// Look in the current dir (./config.yaml)
		// Then $XDG_CONFIG_HOME/kamunder/config.yaml
		// Then $HOME/.config/kamunder/config.yaml
		// Finally fallback to $HOME/.kamunder/config.yaml
		v.AddConfigPath(".")
		if xdg, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok && xdg != "" {
			v.AddConfigPath(filepath.Join(xdg, "kamunder"))
		} else if home, err := os.UserHomeDir(); err == nil {
			v.AddConfigPath(filepath.Join(home, ".config", "kamunder"))
		}
		if home, err := os.UserHomeDir(); err == nil {
			v.AddConfigPath(filepath.Join(home, ".kamunder"))
		}
	}

	// ENV: CAMUNDER_AUTH_CLIENT_ID, etc.
	v.SetEnvPrefix("KAMUNDER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read config (ignore "not found")
	if err := v.ReadInConfig(); err != nil {
		var nf viper.ConfigFileNotFoundError
		if !errors.As(err, &nf) {
			return fmt.Errorf("read config: %w", err)
		}
	}
	return nil
}

func retrieveConfig(v *viper.Viper) (*config2.Config, error) {
	var cfg config2.Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	if tmpScopes := v.GetStringMapString("tmp.auth_scopes"); len(tmpScopes) > 0 {
		if cfg.Auth.OAuth2.Scopes == nil {
			cfg.Auth.OAuth2.Scopes = make(map[string]string, len(tmpScopes))
		}
		for k, scope := range tmpScopes {
			cfg.Auth.OAuth2.Scopes[strings.TrimSpace(k)] = strings.TrimSpace(scope)
		}
	}
	return &cfg, nil
}
