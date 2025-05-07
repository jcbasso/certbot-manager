FROM certbot/certbot:v4.0.0 AS certbot-duckdns-base

RUN pip install certbot_dns_cloudflare==4.0.0 certbot-dns-duckdns==1.6.0

FROM golang:1.24.1-alpine AS gobuilder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/certbot-manager ./cmd/certbot-manager

FROM certbot-duckdns-base AS final

COPY --from=gobuilder /app/certbot-manager /usr/local/bin/certbot-manager
RUN chmod +x /usr/local/bin/certbot-manager

WORKDIR /app

ENTRYPOINT ["/usr/local/bin/certbot-manager"]