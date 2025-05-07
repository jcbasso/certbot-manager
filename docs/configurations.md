# Certbot Manager Configuration Details

This document provides detailed information about configuring Certbot Manager through its TOML file, command-line
arguments, and environment variables.

Configuration is managed through a combination of these methods, following this precedence: **Command-Line Flags >
Environment
Variables > Config File > Built-in Application Defaults**.

## Command-Line Arguments

Arguments passed via the command line have the highest precedence.

```text
# Run with --help to see the most up-to-date list
./certbot-manager --help
```

| Flag             | Shorthand | Description                                             | Default (Application Level) |
|------------------|-----------|---------------------------------------------------------|-----------------------------|
| `--config`       | `-c`      | Path to the TOML configuration file.                    | `./config.toml`             |
| `--certbot-path` |           | Path to the `certbot` executable.                       | `certbot` (uses PATH)       |
| `--log-level`    |           | Logging level (debug, info, warn, error, fatal, panic). | `info`                      |
| `--help`         | `-h`      | Show this help message and exit.                        |                             |

## Configuration TOML File (`config.toml`)

The primary configuration is done via a TOML file. It defines global default settings and settings for each individual
certificate to be managed. The application looks for the file path specified by the `--config` flag.

**Structure:**

* `[globals]`: Defines default settings that apply to all certificates.
* `[[certificate]]`: Defines settings for a specific certificate. Settings here override those in `[globals]`.

### Common Configuration Fields

These fields can be set both in the `[globals]` section (to act as defaults for all certificates) and in
each `[[certificate]]` section (to override the global setting for that specific certificate).

| Key                           | TOML Type | Required (Context)         | Description                                                                                                                        | Example                            | Default (App Level) |
|-------------------------------|-----------|----------------------------|------------------------------------------------------------------------------------------------------------------------------------|------------------------------------|---------------------|
| `cmd`                         | String    | No                         | Certbot subcommand to run (`certonly`, `run`, `enhance`, `none`).                                                                  | `"certonly"`                       | `"certonly"`        |
| `email`                       | String    | Yes (Overall)              | Contact email for Let's Encrypt. Must be set either in `[globals]` or in every `[[certificate]]`.                                  | `"admin@example.com"`              | None                |
| `webroot_path`                | String    | No (If not webroot)        | Path for `webroot` authenticator's ACME challenges. Required if `authenticator` is `webroot`.                                      | `"/var/www/acme-challenge"`        | None                |
| `staging`                     | Boolean   | No                         | Use Let's Encrypt staging server. Recommended for testing.                                                                         | `true`                             | `true`              |
| `no_eff_email`                | Boolean   | No                         | Disable EFF mailing list signup when registering.                                                                                  | `false`                            | `true`              |
| `key_type`                    | String    | No                         | Preferred key type (`ecdsa` or `rsa`). If empty, Certbot's default is used.                                                        | `"ecdsa"`                          | None                |
| `initial_force_renewal`       | Boolean   | No                         | Use `--force-renewal` on the first run for this certificate context.                                                               | `true`                             | None                |
| `args`                        | String    | No                         | **Raw string** of additional arguments passed *directly* to Certbot. Useful for flags not yet implemented directly.                | `"--preferred-challenges http-01"` | None                |
| `authenticator`               | String    | No                         | Certbot authenticator method. See [Supported Authenticators](#supported-authenticators) in the main README.                        | `"dns-duckdns"`                    | None                |
| `dns_propagation_seconds`     | Integer   | No (If not DNS)            | Wait time (seconds) for DNS challenges to propagate. Used by DNS authenticators.                                                   | `60`                               | None                |
| `duckdns_token`               | String    | No (If not DuckDNS)        | DuckDNS API token. Value here takes precedence for this specific certificate or global setting.                                    | `"123456-78910"`                   | None                |
| `cloudflare_credentials_path` | String    | No (If not Cloudflare DNS) | Cloudflare DNS credentials .ini path. See [dns-cloudflare documentation](https://certbot-dns-cloudflare.readthedocs.io/en/stable/) | `"cloudflare.ini"`                 | None                |

### `[globals]` Section Specific Fields

These fields are specific to the `[globals]` section and define application-wide behavior.

| Key            | TOML Type | Required | Description                                  | Example            | Default (App Level) |
|----------------|-----------|----------|----------------------------------------------|--------------------|---------------------|
| `renewal_cron` | String    | Yes      | Cron expression for periodic renewal checks. | `"0 0 0,12 * * *"` | None                |

### `[[certificate]]` Section Specific Fields

These fields are specific to each `[[certificate]]` block.

| Key       | TOML Type        | Required | Description                                                                                                         | Example                              |
|-----------|------------------|----------|---------------------------------------------------------------------------------------------------------------------|--------------------------------------|
| `domains` | Array of Strings | Yes      | List of domain names for this certificate (SANs). The first domain is the primary name for the certificate lineage. | `["example.com", "www.example.com"]` |

**Configuration Override Logic (within TOML):**

1. The application first looks for a "Common Configuration Field" setting within a specific `[[certificate]]` block.
2. If the field is not found in the `[[certificate]]` block, it then looks for that field in the `[globals]` block.
3. If the field is not found in `[globals]` either, the application's built-in "Default (App Level)" for that field will
   be used (as listed in the "Common Configuration Fields" table).

See the example [config.toml](../example.config.toml) in the project root for detailed structure and
comments. <!-- Adjust path as needed -->

## Environment Variables

Environment variables provide a way to configure `certbot-manager` dynamically, often useful for secrets or for
overriding settings without modifying the configuration file. **Values set via environment variables will override
corresponding values defined in the TOML configuration file.**

| Environment Variable        | Overrides                                                         | Description                                                                                                                                                                         |
|-----------------------------|-------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `CERTBOT_MANAGER_GLOBALS_*` | TOML key: `globals.<FIELD_NAME>` or `globals.<COMMON_FIELD_NAME>` | Overrides any field within the `[globals]` section of your `config.toml`. For example, to override `globals.renewal_cron`, use `CERTBOT_MANAGER_GLOBALS_RENEWALCRON="0 0 1 * * *"`. |

**How to Set Environment Variables:**

* **In your shell (for standalone usage):**
  ```bash
  export CERTBOT_MANAGER_GLOBALS_DUCKDNS_TOKEN="123456-78910"
  export CERTBOT_MANAGER_GLOBALS_EMAIL="override@example.com"
  ./certbot-manager --config=./config.toml
  ```
* **Using a `.env` file (often used with Docker Compose but can be sourced by shell scripts too):**
  Create a file named `.env` in your project directory:
  ```env
  CERTBOT_MANAGER_GLOBALS_DUCKDNS_TOKEN=123456-78910
  CERTBOT_MANAGER_GLOBALS_EMAIL=override@example.com
  ```
  If using Docker Compose, it will typically pick this up automatically. For standalone, you might source
  it: `source .env && ./certbot-manager ...`
* **Directly in Docker Compose (if using Docker):**
  ```yaml
  # docker-compose.yml
  services:
    certbot_manager:
      # ...
      environment:
        CERTBOT_MANAGER_GLOBALS_DUCKDNS_TOKEN: "123456-78910"
        CERTBOT_MANAGER_GLOBALS_EMAIL: "override@example.com"
  ```

**Pattern for `CERTBOT_MANAGER_GLOBALS_*` Variables:**

To override a key within the `[globals]` section of your `config.toml` file using an environment variable:

1. Start with the prefix `CERTBOT_MANAGER_GLOBALS_`.
2. Append the TOML key name (as defined for the "Common Configuration Fields" or "Globals Specific Fields",
   e.g., `email`, `webroot_path`, `renewal_cron`) in `UPPERCASE`.

**Examples for `[globals]` overrides:**

* To set `globals.email`:
  `CERTBOT_MANAGER_GLOBALS_EMAIL="admin@example.com"`
* To set `globals.staging` (boolean values are strings like `"true"` or `"false"`):
  `CERTBOT_MANAGER_GLOBALS_STAGING="true"`
* To set `globals.renewal_cron`:
  `CERTBOT_MANAGER_GLOBALS_RENEWALCRON="0 0 1 * * *"`

> **Note:**
> * For boolean environment variables, use string values like `"true"` or `"false"`.
> * Environment variables for specific certificate blocks (e.g., `CERTBOT_MANAGER_CERTIFICATES_0_EMAIL`) are not
    currently supported; use the TOML file for per-certificate overrides.