FROM python:3.12-alpine AS pybuilder

RUN apk add --no-cache \
    gcc \
    musl-dev \
    libffi-dev \
    python3-dev

WORKDIR /wheels

RUN pip wheel --no-cache-dir --wheel-dir=/wheels certbot-dns-duckdns==1.6

FROM certbot/certbot:v4.0.0 AS certbot-duckdns-base

COPY --from=pybuilder /wheels /wheels

RUN pip install --no-cache-dir --no-index --find-links=/wheels /wheels/*.whl \
    && rm -rf /wheels

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