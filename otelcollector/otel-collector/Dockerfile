FROM golang:1.19-alpine AS builder

ARG OTEL_COLLECTOR_TAG=v0.60.0

WORKDIR /go/src/
RUN apk update && apk add make build-base git ca-certificates
RUN git clone --depth 1 --branch $OTEL_COLLECTOR_TAG https://github.com/open-telemetry/opentelemetry-collector-contrib.git
RUN cd opentelemetry-collector-contrib && make install-tools && make otelcontribcol

FROM alpine:latest as prep
RUN apk --update add ca-certificates

FROM debian:bullseye-slim

COPY --from=prep /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/opentelemetry-collector-contrib/bin/otelcontribcol_linux_amd64 /otelcontribcol
EXPOSE 4317 55680 55679
ENTRYPOINT ["/otelcontribcol"]
CMD ["--config", "/etc/otel/config.yaml"]