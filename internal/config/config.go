package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Viper instance (package level)
var v *viper.Viper

var (
	Defaults = Default{
		Cmd:                   "certonly",
		RenewalCron:           "0 0 0,12 * * *",
		WebrootPath:           "/var/www/acme-challenge",
		Staging:               false,
		NoEffEmail:            true,
		KeyType:               "",
		InitialForceRenewal:   false,
		Authenticator:         "webroot",
		DNSPropagationSeconds: 60,
		ConfigFilePath:        "./config.toml",
		LogLevel:              "info",
	}
)

// Default holds default settings
type Default struct {
	Cmd                   string
	RenewalCron           string
	WebrootPath           string
	Staging               bool
	NoEffEmail            bool
	KeyType               string
	InitialForceRenewal   bool
	Authenticator         string
	DNSPropagationSeconds int
	ConfigFilePath        string
	LogLevel              string
}

// Config holds the application configuration
type Config struct {
	Globals      Globals       `mapstructure:"globals"`
	Certificates []Certificate `mapstructure:"certificate"`
	CertbotPath  string
	LogLevel     string
}

// Globals holds global settings
type Globals struct {
	Cmd                 string `mapstructure:"cmd"`
	Email               string `mapstructure:"email"`
	WebrootPath         string `mapstructure:"webroot_path"`
	Staging             bool   `mapstructure:"staging"`
	NoEffEmail          bool   `mapstructure:"no_eff_email"`
	KeyType             string `mapstructure:"key_type"`
	InitialForceRenewal bool   `mapstructure:"initial_force_renewal"`
	RenewalCron         string `mapstructure:"renewal_cron"`
}

// Certificate represents a single certificate request
type Certificate struct {
	Cmd                 string   `mapstructure:"cmd"`
	Domains             []string `mapstructure:"domains"`
	Email               string   `mapstructure:"email"`
	WebrootPath         string   `mapstructure:"webroot_path"`
	Staging             *bool    `mapstructure:"staging"`
	KeyType             string   `mapstructure:"key_type"`
	InitialForceRenewal *bool    `mapstructure:"initial_force_renewal"`
	Args                string   `mapstructure:"args"`
	Authenticator       string   `mapstructure:"authenticator"`
	// Seconds to wait for DNS propagation (only used if authenticator is dns-*)
	DNSPropagationSeconds *int `mapstructure:"dns_propagation_seconds"`
}

// Load initializes Viper and loads the configuration.
func Load() (*Config, error) {
	v = viper.New()

	pflag.StringP("config", "c", Defaults.ConfigFilePath, "Path to the configuration file (e.g., /app/config.toml)")
	pflag.String("certbot-path", "certbot", "Path to the certbot executable")
	pflag.String("log-level", Defaults.LogLevel, "Logging level (debug, info, warn, error, fatal, panic)")
	help := pflag.BoolP("help", "h", false, "Show help message")

	pflag.Parse()

	if *help {
		fmt.Println("Certbot Manager: Manages Let's Encrypt certificates using webroot or DNS challenges.")
		fmt.Println("\nUsage:")
		pflag.PrintDefaults()
		fmt.Println("\nEnvironment Variables:")
		fmt.Println("  DUCKDNS_TOKEN: Required if using 'dns-duckdns' authenticator.")
		fmt.Println("  CERTBOT_MANAGER_* : Can override config values (e.g., CERTBOT_MANAGER_GLOBALS_EMAIL).")
		os.Exit(0)
	}

	// Args
	v.SetDefault("certbotPath", "certbot")
	v.SetDefault("logLevel", Defaults.LogLevel)

	if err := v.BindPFlag("certbotPath", pflag.Lookup("certbot-path")); err != nil {
		log.Printf("Warning: could not bind certbot-path flag: %v", err) // Use standard log before logrus setup
	}
	if err := v.BindPFlag("logLevel", pflag.Lookup("log-level")); err != nil { // Bind log level flag
		log.Printf("Warning: could not bind log-level flag: %v", err)
	}

	// Default
	v.SetDefault("globals.cmd", Defaults.Cmd)
	v.SetDefault("globals.webroot_path", Defaults.WebrootPath)
	v.SetDefault("globals.staging", Defaults.Staging)
	v.SetDefault("globals.no_eff_email", Defaults.NoEffEmail)
	v.SetDefault("globals.key_type", Defaults.KeyType)
	v.SetDefault("globals.initial_force_renewal", Defaults.InitialForceRenewal)
	v.SetDefault("globals.renewal_cron", Defaults.RenewalCron)

	// Env Vars
	v.SetEnvPrefix("CERTBOT_MANAGER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

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

	cfg.CertbotPath = v.GetString("certbotPath")

	if cfg.Globals.RenewalCron == "" {
		log.Printf("Warning: Globals.RenewalCron is empty, falling back to default: %s", Defaults.RenewalCron)
		cfg.Globals.RenewalCron = Defaults.RenewalCron
	}

	return &cfg, nil
}
