FROM golang:1.19-alpine AS builder
ARG SVC
ARG GOARCH
ARG GOARM

WORKDIR /go/src/github.com/etaques/orb
COPY . .
RUN apk update \
    && apk add make
RUN make $SVC \
    && mv build/orb-$SVC /exe

FROM alpine:latest
ARG SVC

RUN if [[ "maestro" == "$SVC" ]]; then apk update && apk add --no-cache docker-cli bash curl && curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl && chmod +x ./kubectl && mv ./kubectl /usr/local/bin/kubectl; fi
# Certificates are needed so that mailing util can work.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /exe /
ENTRYPOINT ["/exe"]
