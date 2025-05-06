<div align="center">

# Certbot Manager

<!-- TODO: Replace with actual badge SVGs or URLs -->
[![go version](docs/assets/badge/go-version-badge.svg)](go.mod)
[![release](docs/assets/badge/release-badge.svg)](https://github.com/jcbasso/certbot-manager/releases)
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
> standalone.
> See [Certbot Installation](https://certbot.eff.org/instructions) for instructions.

1. Download the latest binary from the [Releases](https://github.com/jcbasso/certbot-manager/releases) page.
2. Create your `config.toml` file (see [Configuration](#configuration) section below and the
   detailed [Configuration Details](docs/configurations.md) document).
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
    image: ghcr.io/jcbasso/certbot-manager:latest
    container_name: certbot-manager
    env_file:
      - .env
    volumes:
      - ./config.toml:/app/config.toml:ro
      - letsencrypt_data:/etc/letsencrypt
      - acme_challenge_webroot:/var/www/acme-challenge
    restart: unless-stopped

volumes:
  letsencrypt_data:
  acme_challenge_webroot:
```

**Steps:**

1. Create your `config.toml` file in the same directory as your `docker-compose.yml` (or adjust the volume mount). See
   the [Configuration](#configuration) summary below and the detailed [Configuration Details](docs/configurations.md)
   document.
2. If using DNS authenticators requiring secrets (like DuckDNS), create a `.env` file in your project root containing
   the secrets (e.g., `DUCKDNS_TOKEN=your_secret_token`).
3. Ensure your `docker-compose.yml` correctly defines the `certbot-manager` service and shared volumes.
4. Run `docker compose up -d` to start the service (add `--build` if building locally).

## Configuration

Configuration is managed through a combination of a TOML file, command-line arguments, and environment variables,
following this precedence: **Flags > Environment Variables > Config File > Defaults**.

The primary configuration is done via a TOML file. This file allows you to define:

* **`[globals]`**: Default settings that apply to all certificates unless overridden. This includes settings like your
  default email address for Let's Encrypt, whether to use the staging environment, the Certbot command to run (
  e.g., `certonly`), the cron schedule for renewals, and more.
* **`[[certificate]]`**: One or more blocks, each defining a specific certificate to be managed. Here you list
  the `domains` for the certificate, and you can override any of the global settings (like email, staging,
  authenticator, etc.) or provide specific Certbot arguments for this certificate only.

> [!TIP]
> For a comprehensive list of all configuration options, their descriptions, defaults, and examples, please see the
**[Configuration Details](docs/configurations.md)** document.
>
> An example configuration file can also be found
> at [config.toml.example](./config.toml.example). <!-- TODO: Create this example file -->

Key aspects you will configure include:

* **Domains:** The list of domain names for each certificate.
* **Authenticator:** The method Certbot will use to prove you control the domains (e.g., `webroot` or `dns-duckdns`).
* **Certbot Command:** The main Certbot action to perform (e.g., `certonly` for just getting certs, or `run` if you want
  Certbot to also install them, though this is less common when `certbot-manager` is paired with a separate proxy).
* **Custom Arguments:** Pass any raw Certbot arguments for fine-grained control.

## Supported Authenticators

The `authenticator` field in the `[[certificate]]` section of your `config.toml` determines how domain ownership is
verified:

* `webroot` (Default): Uses the HTTP-01 challenge. Requires your web server (e.g., Nginx) to be configured to serve
  files from the `webroot_path` for the `/.well-known/acme-challenge/` URI.
* `dns-duckdns`: Uses the DNS-01 challenge via
  the [infinityofspace/certbot_dns_duckdns](https://github.com/infinityofspace/certbot_dns_duckdns) plugin. Requires
  the `DUCKDNS_TOKEN` environment variable to be set.

<!-- TODO: Add more authenticators here as they are implemented -->

## Development

* Requires Go version specified in [go.mod](go.mod).
* Follow [git branching model & release specifications](docs/git-branching-model.md).

## License

This project is licensed under
the [Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License](LICENSE).