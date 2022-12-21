FROM golang:1.19-alpine AS builder

WORKDIR /go/src/healthcheck
COPY go.mod .
RUN go mod tidy
COPY . .
RUN apk update && apk add make build-base git
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o healthcheck
RUN chmod a+x /go/src/healthcheck/run.sh

FROM etaques/opentelemetry-collector-contrib:0.60.10 as otel

FROM alpine:latest as prep
RUN apk --update add ca-certificates

FROM debian:bullseye-slim
COPY --from=prep /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/healthcheck/healthcheck /healthcheck
COPY --from=builder /go/src/healthcheck/run.sh /run.sh
COPY --from=otel /otelcontribcol /usr/local/bin/otelcontribcol

ENTRYPOINT [ "/run.sh" ]