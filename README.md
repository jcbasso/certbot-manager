<div align="center">

# Certbot Manager

<!-- TODO: Replace with actual badge SVGs or URLs -->
[![go version](docs/assets/badge/go-version-badge.svg)](go.mod)
[![release](docs/assets/badge/release-badge.svg)](https://github.com/YOUR_USERNAME/YOUR_REPO/releases)
[![CC BY-NC-SA 4.0](docs/assets/badge/license.svg)](https://creativecommons.org/licenses/by-nc-sa/4.0/)

Go application to automatically obtain and renew multiple Let's Encrypt certificates using Certbot, configured via a
TOML file. Designed to run alongside a reverse proxy like Nginx.

<a href="#">
    <img src="docs/assets/logo.png" alt="Certbot Manager Logo" style="width: 50%; height: auto;">
</a>

</div>

---

## Features

* Automated certificate acquisition via Certbot.
* Automated certificate renewal via Certbot through Go cron scheduler.
* Declarative configuration using a TOML file (`config.toml`).
* Support for different Certbot authenticators (`webroot`, `dns-duckdns`).
* Customizable Certbot arguments per certificate.
* Leveled logging controllable via flags or environment variables.
* Designed for containerized environments (Docker).
* Open to extensibility for additional features, flags and authenticator plugins.

## Standalone Usage

> [!IMPORTANT]
> Since `certbot-manager` leverages Certbot, you must first install Certbot separately on your system if running
standalone.
> See [Certbot Installation](https://certbot.eff.org/instructions) for instructions.

1. Download the latest binary from the [Releases](https://github.com/YOUR_USERNAME/YOUR_REPO/releases) page.
2. Create your `config.toml` file (see [Configuration](#configuration)).
3. Ensure required environment variables (like `DUCKDNS_TOKEN`) are set if needed.
4. Run the binary:
   ```bash
   # Example: Use a config file in the current directory
   ./certbot-manager --config=./config.toml --log-level=debug

   # Example: Specify a different path
   ./certbot-manager -c /etc/certbot-manager/config.toml
   ```

## Docker Compose Usage

This application is primarily intended to be run as a Docker container using Docker Compose, alongside your web
server/proxy container. The Certbot dependency is included in the Docker image.

```yaml
services:
  certbot-manager:
    image: your-dockerhub-username/certbot-manager:latest
    container_name: certbot-manager
    # env_file:
    #   - .env
    environment:
    # CERTBOT_MANAGER_LOGLEVEL: "debug"
    # CERTBOT_MANAGER_GLOBALS_EMAIL: "prod-admin@example.com"
    # DUCKDNS_TOKEN: ${DUCKDNS_TOKEN} # Reads from .env file
    volumes:
      - ./config.toml:/app/config.toml:ro
      - letsencrypt_data:/etc/letsencrypt
      - acme_challenge_webroot:/var/www/acme-challenge
    restart: unless-stopped
    # command: ["--config", "/etc/custom/my-config.toml", "--log-level", "debug"]

volumes:
  letsencrypt_data:
  acme_challenge_webroot:
```

**Steps:**

1. Create your `config.toml` file in the same directory as your `docker-compose.yml` (or adjust the volume mount).
2. If using DNS authenticators requiring secrets (like DuckDNS), create a `.env` file in your project root containing
   the secrets (e.g., `DUCKDNS_TOKEN=your_secret_token`).
3. Ensure your `docker-compose.yml` correctly defines the `certbot-manager` service and shared volumes.
4. Run `docker compose up -d` to start the service.

## Configuration

Configuration is managed through a combination of a TOML file, command-line arguments, and environment variables,
following this precedence: **Flags > Environment Variables > Config File > Defaults**.

### Configuration TOML File (`config.toml`)

This file defines the certificates to manage and global/specific settings. The application looks for the file path
specified by the `--config` flag (default: `./config.toml` when run standalone, `/app/config.toml` is the typical mount
point in Docker).

**Structure:**

* `[globals]`: Default settings applied to all certificates unless overridden.
* `[[certificate]]`: Defines settings for a specific certificate request (can have multiple blocks for multiple
  certificates).

**Key `[globals]` Fields:**

| Key                     | Required | Description                                                                   | Example                      |
|-------------------------|----------|-------------------------------------------------------------------------------|------------------------------|
| `email`                 | Yes      | Default contact email for Let's Encrypt registration/recovery.                | `"admin@example.com"`        |
| `webroot_path`          | No       | Default path for `webroot` authenticator's ACME challenges.                   | `"/var/www/acme-challenge"`  |
| `staging`               | No       | Use Let's Encrypt staging server (default: `false`). Recommended for testing. | `true`                       |
| `key_type`              | No       | Preferred key type (`ecdsa` or `rsa`, default: `""` -> certbot default).      | `"ecdsa"`                    |
| `renewal_cron`          | No       | Cron expression for periodic renewal checks.                                  | `"0 0 3 * * *"` (3 AM daily)|
| `initial_force_renewal` | No       | Use `--force-renewal` on first run (default: `false`).                        | `true`                       |
| `no_eff_email`          | No       | Disable EFF mailing list signup (default: `true`).                            | `true`                       |

**Key `[[certificate]]` Fields:**

| Key                       | Required | Description                                                                                                                                        | Example                              |
|---------------------------|----------|----------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------|
| `domains`                 | Yes      | List of domain names for this certificate (SANs). First is primary name.                                                                           | `["example.com", "www.example.com"]` |
| `authenticator`           | No       | Certbot authenticator for this cert (default: `webroot`). See [Supported Authenticators](#supported-authenticators).                               | `"dns-duckdns"`                      |
| `args`                    | No       | **Raw string** of additional arguments passed *directly* to certbot for this certificate only. Very useful if an argument was yet not implemented. | `"--preferred-challenges http-01"`   |
| `dns_propagation_seconds` | No       | Wait time (seconds) for DNS challenges (default: `60`). Used by DNS authenticators.                                                                | `60`                                 |
| `email`                   | No       | Override global email.                                                                                                                             | `"specific@example.com"`             |
| `webroot_path`            | No       | Override global `webroot_path` (only used if authenticator is `webroot`).                                                                          | `"/var/www/specific-app"`            |
| `staging`                 | No       | Override global staging setting.                                                                                                                   | `true`                               |
| `key_type`                | No       | Override global key type.                                                                                                                          | `"rsa"`                              |
| `initial_force_renewal`   | No       | Override global initial force renewal setting.                                                                                                     | `true`                               |

> [!TIP]
> See the example [config.toml](./example-config.toml) for detailed structure.

### Command-Line Arguments

Arguments passed via the command line override environment variables and the configuration file.

```text
# See help message for current flags and defaults
./certbot-manager --help
```

| Flag             | Shorthand | Description                                             | Default                  |
|------------------|-----------|---------------------------------------------------------|--------------------------|
| `--config`       | `-c`      | Path to the TOML configuration file.                    | `/app/config.toml` (in Docker) / `./config.toml` (standalone default) |
| `--certbot-path` |           | Path to the `certbot` executable.                       | `certbot` (uses PATH)    |
| `--log-level`    |           | Logging level (debug, info, warn, error, fatal, panic). | `info`                   |
| `--help`         | `-h`      | Show this help message and exit.                        |                          |

### Environment Variables

Set these in your shell, via a `.env` file used by Docker Compose, or directly in the `environment:` section
of `docker-compose.yml`.

| Environment Variable | Required                     | Description                                                                                                                                                              |
|----------------------|------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `DUCKDNS_TOKEN`      | Yes (if using `dns-duckdns`) | API token for your DuckDNS account.                                                                                                                                      |
| `CERTBOT_MANAGER_*`  | No                           | Overrides corresponding config values (e.g., `CERTBOT_MANAGER_GLOBALS_EMAIL=...`, `CERTBOT_MANAGER_LOGLEVEL=debug`). Follow Viper's key structure (use `_` for nesting). |

## Supported Authenticators

The `authenticator` field in the `[[certificate]]` section determines how domain ownership is verified:

* `webroot` (Default): Uses HTTP-01 challenge. Requires web server configuration.
* `dns-duckdns`: Uses DNS-01 challenge
  via [infinityofspace/certbot_dns_duckdns](https://github.com/infinityofspace/certbot_dns_duckdns).
  Requires `DUCKDNS_TOKEN` env var.

<!-- TODO: Add more authenticators here as they are implemented -->

## Development

* Use specified [go version](#certbot-manager).
* Follow [git branching model & release specifications](docs/git-branching-model.md).

## License

This project is licensed under the [Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License](LICENSE).