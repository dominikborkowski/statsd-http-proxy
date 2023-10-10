ARG VERSION
FROM golang:1.21-alpine as builder
LABEL org.opencontainers.image.authors="dominik.borkowski@gmail.com"
LABEL org.opencontainers.image.source="https://github.com/dominikborkowski/statsd-http-proxy"
LABEL org.opencontainers.image.description="HTTP proxy for StatsD"
COPY . /source

RUN apk add --no-cache --virtual .build-deps \
    ca-certificates \
    make \
    wget \
    git \
    curl \
    go \
    musl-dev && \
    # update certs
    update-ca-certificates

RUN cd /source && \
    make build

# golang:1.21-alpine is based on alpine 3.18
FROM alpine:3.18

# for heartbeat
RUN apk add --no-cache netcat-openbsd curl
COPY --from=builder /source/bin/statsd-http-proxy /usr/local/bin/

# start service
EXPOSE 8825
ENTRYPOINT ["/usr/local/bin/statsd-http-proxy", "--http-host="]
