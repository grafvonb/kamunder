package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/grafvonb/kamunder/config"
	"github.com/grafvonb/kamunder/internal/services/auth"
	"github.com/grafvonb/kamunder/internal/services/auth/authenticator"
	"github.com/grafvonb/kamunder/internal/services/httpc"
	"github.com/grafvonb/kamunder/toolx"
	"github.com/grafvonb/kamunder/toolx/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagShowConfig bool
	//nolint:unused
	flagAsJson bool
)

var rootCmd = &cobra.Command{
	Use:   "kamunder",
	Short: "Kamunder is a CLI tool to interact with Camunda 8.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		v := viper.New()
		if err := initViper(v, cmd); err != nil {
			return err
		}

		if flagShowConfig {
			cfg, err := retrieveConfig(v, false)
			if err != nil {
				return err
			}
			if p := v.ConfigFileUsed(); p != "" {
				cmd.Println("config loaded:", p)
			}
			cmd.Println(cfg.String())
			os.Exit(0)
			return nil
		}

		if isUtilityCommand(cmd) || hasHelpFlag(cmd) {
			return nil
		}

		cfg, err := retrieveConfig(v, true)
		if err != nil {
			return err
		}
		ctx := cfg.ToContext(cmd.Context())

		log := logging.New(logging.LoggerConfig{
			Level:      v.GetString("log.level"),
			Format:     v.GetString("log.format"),
			WithSource: v.GetBool("log.with_source"),
		})
		ctx = logging.ToContext(ctx, log)

		if pathcfg := v.ConfigFileUsed(); pathcfg != "" {
			log.Debug("config loaded: " + pathcfg)
		}

		httpSvc, err := httpc.New(cfg, log, httpc.WithCookieJar())
		if err != nil {
			return fmt.Errorf("http service: %w", err)
		}
		ator, err := auth.BuildAuthenticator(cfg, httpSvc.Client(), log)
		if err != nil {
			return fmt.Errorf("auth build: %w", err)
		}
		if err := ator.Init(ctx); err != nil {
			return fmt.Errorf("auth init: %w", err)
		}
		httpSvc.InstallAuthEditor(ator.Editor())
		ctx = httpSvc.ToContext(ctx)

		ctx = authenticator.ToContext(ctx, ator)
		cmd.SetContext(ctx)

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
	SilenceUsage:  true,
	SilenceErrors: false,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
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

	pf.BoolVar(&flagShowConfig, "show-config", false, "print effective config (secrets redacted)")
}

func initViper(v *viper.Viper, cmd *cobra.Command) error {
	fs := cmd.Flags()

	_ = v.BindPFlag("config", fs.Lookup("config"))

	_ = v.BindPFlag("log.level", fs.Lookup("log-level"))
	_ = v.BindPFlag("log.format", fs.Lookup("log-format"))
	_ = v.BindPFlag("log.with_source", fs.Lookup("log-with-source"))

	_ = v.BindPFlag("app.tenant", fs.Lookup("tenant"))

	_ = v.BindPFlag("auth.token_url", fs.Lookup("auth-token-url"))
	_ = v.BindPFlag("auth.client_id", fs.Lookup("auth-client-id"))
	_ = v.BindPFlag("auth.client_secret", fs.Lookup("auth-client-secret"))
	_ = v.BindPFlag("tmp.auth_scopes", fs.Lookup("auth-scopes"))

	_ = v.BindPFlag("http.timeout", fs.Lookup("http-timeout"))

	_ = v.BindPFlag("apis.version", fs.Lookup("camunda-apis-version"))
	_ = v.BindPFlag("apis.camunda_api.base_url", fs.Lookup("camunda-base-url"))
	_ = v.BindPFlag("apis.operate_api.base_url", fs.Lookup("operate-base-url"))
	_ = v.BindPFlag("apis.tasklist_api.base_url", fs.Lookup("tasklist-base-url"))

	v.Set("apis.camunda_api.key", config.CamundaApiKeyConst)
	v.Set("apis.operate_api.key", config.OperateApiKeyConst)
	v.Set("apis.tasklist_api.key", config.TasklistApiKeyConst)

	v.SetDefault("http.timeout", "30s")

	v.SetEnvPrefix("KAMUNDER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Config file resolution and read
	if p := v.GetString("config"); p != "" {
		v.SetConfigFile(p)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("$XDG_CONFIG_HOME/kamunder")
		v.AddConfigPath("$HOME/.config/kamunder")
		v.AddConfigPath("$HOME/.kamunder")
		v.AddConfigPath("/etc/kamunder")
	}
	if err := v.ReadInConfig(); err != nil {
		var nf viper.ConfigFileNotFoundError
		if !errors.As(err, &nf) || v.GetString("config") != "" {
			return fmt.Errorf("read config file: %w", err)
		}
	}
	return nil
}

func retrieveConfig(v *viper.Viper, validate bool) (*config.Config, error) {
	var cfg config.Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if tmp := v.GetStringMapString("tmp.auth_scopes"); len(tmp) > 0 {
		if cfg.Auth.OAuth2.Scopes == nil {
			cfg.Auth.OAuth2.Scopes = make(map[string]string, len(tmp))
		}
		for k, scope := range tmp {
			k = strings.TrimSpace(k)
			scope = strings.TrimSpace(scope)
			if k == "" || scope == "" {
				continue
			}
			cfg.Auth.OAuth2.Scopes[k] = scope
		}
	}

	if validate {
		if err := cfg.Validate(); err != nil {
			return nil, fmt.Errorf("validate config: %w", err)
		}
	}

	return &cfg, nil
}
