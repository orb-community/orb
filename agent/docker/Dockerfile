ARG PKTVISOR_TAG=develop

FROM golang:1.16-alpine AS builder
ARG GOARCH
ARG GOARM

WORKDIR /go/src/github.com/ns1labs/orb
COPY . .
RUN apk update \
    && apk add make gcc musl-dev
RUN CGO_ENABLED=1 make agent_bin \
    && mv build/orb-agent /exe

FROM ns1labs/pktvisor:${PKTVISOR_TAG}

COPY --from=builder /exe /usr/local/bin/orb-agent
COPY --from=builder /go/src/github.com/ns1labs/orb/agent/docker/agent.yaml /etc/orb/agent.yaml

ENTRYPOINT /usr/local/bin/orb-agent
