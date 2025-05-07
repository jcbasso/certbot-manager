package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Viper instance (package level)
var v *viper.Viper

var (
	Defaults = Default{
		Staging:        true,
		NoEffEmail:     true,
		Cmd:            "certonly",
		ConfigFilePath: "./config.toml",
		CertbotPath:    "certbot",
		LogLevel:       "info",
	}
)

// Default holds default settings
type Default struct {
	Staging        bool
	NoEffEmail     bool
	Cmd            string
	ConfigFilePath string
	CertbotPath    string
	LogLevel       string
}

// Config holds the application configuration
type Config struct {
	Globals      Globals       `mapstructure:"globals"`
	Certificates []Certificate `mapstructure:"certificate"`
	CertbotPath  string
	LogLevel     string
}

type CommonConfigs struct {
	Cmd                 string `mapstructure:"cmd"`
	Email               string `mapstructure:"email"`
	WebrootPath         string `mapstructure:"webroot_path"`
	Staging             *bool  `mapstructure:"staging"`
	NoEffEmail          *bool  `mapstructure:"no_eff_email"`
	KeyType             string `mapstructure:"key_type"`
	InitialForceRenewal *bool  `mapstructure:"initial_force_renewal"`
	Args                string `mapstructure:"args"`
	Authenticator       string `mapstructure:"authenticator"`
	// Seconds to wait for DNS propagation (only used if authenticator is dns-*)
	DNSPropagationSeconds     *int   `mapstructure:"dns_propagation_seconds"`
	CloudflareCredentialsPath string `mapstructure:"cloudflare_credentials_path"`
	DuckDNSToken              string `mapstructure:"duckdns_token"`
}

// Globals holds global settings
type Globals struct {
	RenewalCron   string `mapstructure:"renewal_cron"`
	CommonConfigs `mapstructure:",squash"`
}

// Certificate represents a single certificate request
type Certificate struct {
	Domains       []string `mapstructure:"domains"`
	CommonConfigs `mapstructure:",squash"`
}

// Load initializes Viper and loads the configuration.
func Load() (*Config, error) {
	v = viper.New()

	pflag.StringP("config", "c", Defaults.ConfigFilePath, "Path to the configuration file (e.g., /app/config.toml)")
	pflag.String("certbot-path", Defaults.CertbotPath, "Path to the certbot executable")
	pflag.String("log-level", Defaults.LogLevel, "Logging level (debug, info, warn, error, fatal, panic)")
	help := pflag.BoolP("help", "h", false, "Show help message")

	pflag.Parse()

	if *help {
		fmt.Println("Certbot Manager: Manages Let's Encrypt certificates using webroot or DNS challenges.")
		fmt.Println("\nUsage:")
		pflag.PrintDefaults()
		fmt.Println("\nEnvironment Variables:")
		fmt.Println("  CERTBOT_MANAGER_* : Can override config values (e.g., CERTBOT_MANAGER_GLOBALS_EMAIL).")
		os.Exit(0)
	}

	// Args
	if err := v.BindPFlag("certbotPath", pflag.Lookup("certbot-path")); err != nil {
		log.Printf("Warning: could not bind certbot-path flag: %v", err) // Use standard log before logrus setup
	}
	if err := v.BindPFlag("logLevel", pflag.Lookup("log-level")); err != nil { // Bind log level flag
		log.Printf("Warning: could not bind log-level flag: %v", err)
	}

	// Defaults
	v.SetDefault("globals.staging", Defaults.Staging)
	v.SetDefault("globals.no_eff_email", Defaults.NoEffEmail)
	v.SetDefault("globals.cmd", Defaults.Cmd)

	// Env Vars
	v.SetEnvPrefix("CERTBOT_MANAGER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	var c Config
	bindEnvsRecursive("globals", reflect.ValueOf(&c.Globals), v)

	// Config File
	configFilePath, _ := pflag.CommandLine.GetString("config") // Use the parsed value

	if _, err := os.Stat(configFilePath); err == nil {
		// File exists
		v.SetConfigFile(configFilePath)
		//log.Printf("Attempting to load configuration from: %s", configFilePath)

		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file '%s': %w", configFilePath, err)
		}
		//log.Printf("Successfully loaded config file: %s", v.ConfigFileUsed())
	} else if os.IsNotExist(err) {
		// File specified by flag (or default) does not exist
		return nil, fmt.Errorf("config file '%s' not found", configFilePath)
	} else {
		// Different error stating the file (e.g., permission denied)
		return nil, fmt.Errorf("error checking config file '%s': %w", configFilePath, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	// Validations
	if cfg.Globals.RenewalCron == "" {
		return nil, fmt.Errorf("globals.RenewalCron is empty")
	}

	return &cfg, nil
}
