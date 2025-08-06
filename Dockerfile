FROM postgres:17.5-alpine3.22 AS postgres

COPY internal/db/migration/migration.sql /docker-entrypoint-initdb.d/migration.sql

FROM golang:1.24.5-alpine3.22 AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -C /app/cmd/server -o /app/cmd/capybelga

FROM busybox:musl AS busybox

FROM debian:bookworm-slim AS certs
RUN apt update && apt install -y ca-certificates

FROM scratch AS base

COPY --from=busybox /etc/passwd /etc/passwd
COPY --from=busybox /etc/group /etc/group
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /app/ssl/certs/

RUN --mount=from=busybox,dst=/usr/ ["busybox", "sh", "-c", "mkdir -p /app && chmod 777 /app"]
RUN --mount=from=busybox,dst=/usr/ ["busybox", "sh", "-c", "addgroup -S tux -g 1010 && adduser -S tux -u 1010 --ingroup tux --disabled-password"]

ENV HOME=/app
ENV USER=tux
ENV PATH=/usr/local/bin:/app
ENV SSL_CERT_DIR=/app/ssl/certs

FROM base AS release

USER tux

COPY --from=builder --chown=1010:1010 /app/cmd/capybelga /app/capybelga

EXPOSE 8080

ENTRYPOINT ["/app/capybelga"]
