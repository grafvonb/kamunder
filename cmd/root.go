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
	flagViewAsJson   bool
	flagViewKeysOnly bool
	flagQuiet        bool
)

var rootCmd = &cobra.Command{
	Use:   "kamunder",
	Short: "Kamunder is a CLI tool to interact with Camunda 8",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		v := viper.New()
		if err := initViper(v, cmd); err != nil {
			return err
		}
		if hasHelpFlag(cmd) {
			return nil
		}

		cfg, err := retrieveAndNormalizeConfig(v)
		if err != nil {
			return err
		}
		ctx := cfg.ToContext(cmd.Context())

		if flagQuiet {
			v.Set("log.level", "error")
		}
		log := logging.New(logging.LoggerConfig{
			Level:      v.GetString("log.level"),
			Format:     v.GetString("log.format"),
			WithSource: v.GetBool("log.with_source"),
		})
		ctx = logging.ToContext(ctx, log)

		if pathcfg := v.ConfigFileUsed(); pathcfg != "" {
			log.Debug("config loaded: " + pathcfg)
		} else {
			log.Debug("no config file loaded, using defaults and environment variables")
		}
		if isUtilityCommand(cmd) {
			cmd.SetContext(ctx)
			return nil
		}

		if err = cfg.Validate(); err != nil {
			return fmt.Errorf("validate config:\n%w", err)
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
	pf.BoolVarP(&flagQuiet, "quiet", "q", false, "suppress all output, except errors")
	pf.BoolVarP(&flagViewAsJson, "json", "j", false, "output as JSON (where applicable)")
	pf.BoolVar(&flagViewKeysOnly, "keys-only", false, "output as keys only (where applicable)")

	pf.String("config", "", "path to config file")

	pf.String("log-level", "info", "log level (debug, info, warn, error)")
	pf.String("log-format", "plain", "log format (json, plain, text)")
	pf.Bool("log-with-source", false, "include source file and line number in logs")

	pf.String("tenant", "", "default tenant ID")

	pf.String("auth-mode", "oauth2", "authentication mode (oauth2, cookie)")
	pf.String("auth-oauth2-client-id", "", "auth client ID")
	pf.String("auth-oauth2-client-secret", "", "auth client secret")
	pf.String("auth-oauth2-token-url", "", "auth token URL")
	pf.StringToString("auth-oauth2-scopes", nil, "auth scopes as key=value (repeatable or comma-separated)")
	pf.String("auth-cookie-base-url", "", "auth cookie base URL")
	pf.String("auth-cookie-username", "", "auth cookie username")
	pf.String("auth-cookie-password", "", "auth cookie password")

	pf.String("http-timeout", "", "HTTP timeout (Go duration, e.g. 30s)")

	pf.StringP("camunda-version", "a", string(toolx.CurrentCamundaVersion), fmt.Sprintf("Camunda version (%s) expected. Causes usage of specific API versions.", toolx.SupportedCamundaVersionsString()))
	pf.String("api-camunda-base-url", "", "Camunda API base URL")
	pf.String("api-operate-base-url", "", "Operate API base URL")
	pf.String("api-tasklist-base-url", "", "Tasklist API base URL")
}

func initViper(v *viper.Viper, cmd *cobra.Command) error {
	fs := cmd.Flags()

	_ = v.BindPFlag("config", fs.Lookup("config"))

	_ = v.BindPFlag("log.level", fs.Lookup("log-level"))
	_ = v.BindPFlag("log.format", fs.Lookup("log-format"))
	_ = v.BindPFlag("log.with_source", fs.Lookup("log-with-source"))

	_ = v.BindPFlag("app.tenant", fs.Lookup("tenant"))

	_ = v.BindPFlag("auth.mode", fs.Lookup("auth-mode"))
	_ = v.BindPFlag("auth.oauth2.client_id", fs.Lookup("auth-oauth2-client-id"))
	_ = v.BindPFlag("auth.oauth2.client_secret", fs.Lookup("auth-oauth2-client-secret"))
	_ = v.BindPFlag("auth.oauth2.token_url", fs.Lookup("auth-oauth2-token-url"))
	_ = v.BindPFlag("auth.oauth2.scopes", fs.Lookup("auth-oauth2-scopes"))
	_ = v.BindPFlag("auth.cookie.base_url", fs.Lookup("auth-cookie-base-url"))
	_ = v.BindPFlag("auth.cookie.username", fs.Lookup("auth-cookie-username"))
	_ = v.BindPFlag("auth.cookie.password", fs.Lookup("auth-cookie-password"))

	_ = v.BindPFlag("http.timeout", fs.Lookup("http-timeout"))

	_ = v.BindPFlag("apis.version", fs.Lookup("camunda-version"))
	_ = v.BindPFlag("apis.camunda_api.base_url", fs.Lookup("api-camunda-base-url"))
	_ = v.BindPFlag("apis.operate_api.base_url", fs.Lookup("api-operate-base-url"))
	_ = v.BindPFlag("apis.tasklist_api.base_url", fs.Lookup("api-tasklist-base-url"))

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

func retrieveAndNormalizeConfig(v *viper.Viper) (*config.Config, error) {
	var cfg config.Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	if err := cfg.Normalize(); err != nil {
		return nil, fmt.Errorf("normalize config: %w", err)
	}
	return &cfg, nil
}
