FROM golang:1.24.1-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/certbot-manager ./cmd/certbot-manager

FROM certbot/certbot:latest

COPY --from=builder /app/certbot-manager /usr/local/bin/certbot-manager
RUN chmod +x /usr/local/bin/certbot-manager

# Install the DuckDNS plugin using pip
RUN apk add --no-cache --virtual .build-deps \
        gcc \
        musl-dev \
        libffi-dev \
        python3-dev \
    && pip install --no-cache-dir certbot-dns-duckdns \
    && pip install certbot_dns_duckdns -U \
    && apk del .build-deps # Clean up build dependencies

# Create directory for webroot challenges (permissions handled by volume mount ideally)
# RUN mkdir -p /var/www/acme-challenge && chown nobody:nogroup /var/www/acme-challenge

WORKDIR /app

ENTRYPOINT ["/usr/local/bin/certbot-manager"]