# Certbot Manager Configuration Details

This document provides detailed information about configuring Certbot Manager through its TOML file, command-line
arguments, and environment variables.

Configuration is managed through a combination of these methods, following this precedence: **Flags > Environment
Variables > Config File > Defaults**.

## Command-Line Arguments

Arguments passed via the command line override environment variables and the configuration file.

```text
# Run with --help to see the most up-to-date list
./certbot-manager --help
```

| Flag             | Shorthand | Description                                             | Default                                                               |
|------------------|-----------|---------------------------------------------------------|-----------------------------------------------------------------------|
| `--config`       | `-c`      | Path to the TOML configuration file.                    | `/app/config.toml` (in Docker) / `./config.toml` (standalone default) |
| `--certbot-path` |           | Path to the `certbot` executable.                       | `certbot` (uses PATH)                                                 |
| `--log-level`    |           | Logging level (debug, info, warn, error, fatal, panic). | `info`                                                                |
| `--help`         | `-h`      | Show this help message and exit.                        |                                                                       |

## Configuration TOML File (`config.toml`)

This file defines the certificates to manage and global/specific settings. The application looks for the file path
specified by the `--config` flag (default: `./config.toml` when run standalone, `/app/config.toml` is the typical mount
point in Docker).

**Structure:**

* `[globals]`: Default settings applied to all certificates unless overridden.
* `[[certificate]]`: Defines settings for a specific certificate request (can have multiple blocks for multiple
  certificates).

### `[globals]` Section Fields

These settings apply to all `[[certificate]]` blocks unless explicitly overridden within a specific block.

| Key                     | Required    | Description                                                                  | Example                      | Default                     |
|-------------------------|-------------|------------------------------------------------------------------------------|------------------------------|-----------------------------|
| `email`                 | Yes         | Default contact email for Let's Encrypt registration/recovery.               | `"admin@example.com"`        | (None - Must be set)        |
| `cmd`                   | No          | Default subcommand to run on certbot (`certonly`, `run`, `enhance`, `none`). | `"certonly"`                 | `"certonly"`                |
| `webroot_path`          | No          | Default path for `webroot` authenticator's ACME challenges.                  | `"/var/www/acme-challenge"`  | `"/var/www/acme-challenge"` |
| `staging`               | No          | Use Let's Encrypt staging server. Recommended for testing.                   | `true`                       | `false`                     |
| `key_type`              | No          | Preferred key type (`ecdsa` or `rsa`). If empty, Certbot's default is used.  | `"ecdsa"`                    | `""`                        |
| `renewal_cron`          | No          | Cron expression for periodic renewal checks.                                 | `"0 0 3 * * *"` (3 AM daily) | `"0 0 0,12 * * *"`          |
| `initial_force_renewal` | No          | Use `--force-renewal` on the first run for each certificate.                 | `true`                       | `false`                     |
| `no_eff_email`          | No          | Disable EFF mailing list signup when registering.                            | `false`                      | `true`                      |
 <!--                    | `agree_tos` | No                                                                           | Agree to Let's Encrypt ToS.  | `true`                      | `true` (enforced)         | -->

### `[[certificate]]` Section Fields

Each `[[certificate]]` block defines a separate request to Certbot.

| Key                       | Required | Description                                                                                                                                        | Example                              | Default                          |
|---------------------------|----------|----------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------|----------------------------------|
| `domains`                 | Yes      | List of domain names for this certificate (SANs). The first domain is the primary name for the certificate lineage.                                | `["example.com", "www.example.com"]` | (None - Must be set)             |
| `cmd`                     | No       | Override the global subcommand for this specific certificate (`certonly`, `run`, `enhance`, `none`).                                               | `"run"`                              | (Global `cmd`)                   |
| `args`                    | No       | **Raw string** of additional arguments passed *directly* to Certbot for this certificate only. Very useful if an argument was not yet implemented. | `"--preferred-challenges http-01"`   | `""`                             |
| `authenticator`           | No       | Certbot authenticator for this certificate. See [Supported Authenticators](#supported-authenticators) in the main README.                          | `"dns-duckdns"`                      | `"webroot"`                      |
| `dns_propagation_seconds` | No       | Wait time (seconds) for DNS challenges to propagate. Used by DNS authenticators.                                                                   | `120`                                | `60`                             |
| `email`                   | No       | Override global email for this specific certificate.                                                                                               | `"specific@example.com"`             | (Global `email`)                 |
| `webroot_path`            | No       | Override global `webroot_path` (only used if this certificate's authenticator is `webroot`).                                                       | `"/var/www/specific-app"`            | (Global `webroot_path`)          |
| `staging`                 | No       | Override global staging setting for this specific certificate.                                                                                     | `true`                               | (Global `staging`)               |
| `key_type`                | No       | Override global key type for this specific certificate.                                                                                            | `"rsa"`                              | (Global `key_type`)              |
| `initial_force_renewal`   | No       | Override global initial force renewal setting for this specific certificate.                                                                       | `true`                               | (Global `initial_force_renewal`) |

See the example [config.toml](../../config.toml.example) in the project root for detailed structure and
comments. <!-- Adjust path as needed -->

## Environment Variables

Set these in your shell, via a `.env` file used by Docker Compose, or directly in the `environment:` section
of `docker-compose.yml`. Environment variables override values from the TOML configuration file.

| Environment Variable                          | Corresponds to Config Key                 | Description                                                                                                       |
|-----------------------------------------------|-------------------------------------------|-------------------------------------------------------------------------------------------------------------------|
| `DUCKDNS_TOKEN`                               | (Specific to `dns-duckdns` authenticator) | API token for your DuckDNS account. Required if using the `dns-duckdns` authenticator for any certificate.        |
| `CERTBOT_MANAGER_GLOBALS_EMAIL`               | `globals.email`                           | Overrides the default email.                                                                                      |
| `CERTBOT_MANAGER_GLOBALS_CMD`                 | `globals.cmd`                             | Overrides the default Certbot command.                                                                            |
| `CERTBOT_MANAGER_GLOBALS_WEBROOTPATH`         | `globals.webroot_path`                    | Overrides the default webroot path.                                                                               |
| `CERTBOT_MANAGER_GLOBALS_STAGING`             | `globals.staging`                         | Overrides the default staging flag (e.g., `"true"` or `"false"`).                                                 |
| `CERTBOT_MANAGER_GLOBALS_KEYTYPE`             | `globals.key_type`                        | Overrides the default key type.                                                                                   |
| `CERTBOT_MANAGER_GLOBALS_RENEWALCRON`         | `globals.renewal_cron`                    | Overrides the cron expression for renewals.                                                                       |
| `CERTBOT_MANAGER_GLOBALS_INITIALFORCERENEWAL` | `globals.initial_force_renewal`           | Overrides the initial force renewal flag.                                                                         |
| `CERTBOT_MANAGER_GLOBALS_NOEFFEMAIL`          | `globals.no_eff_email`                    | Overrides the no EFF email flag.                                                                                  |
| `CERTBOT_MANAGER_CERTBOTPATH`                 | (Command line flag equivalent)            | Overrides the path to the Certbot executable.                                                                     |
| `CERTBOT_MANAGER_LOGLEVEL`                    | (Command line flag equivalent)            | Overrides the logging level.                                                                                      |
| `CERTBOT_MANAGER_CONFIG`                      | (Command line flag equivalent)            | Overrides the path to the config file (less common to set via env var if using the flag or default volume mount). |

> **Note:** For boolean environment variables, use string values like `"true"` or `"false"`. For nested TOML keys (
> like `globals.email`), use underscores in the environment variable name as shown (`CERTBOT_MANAGER_GLOBALS_EMAIL`).
