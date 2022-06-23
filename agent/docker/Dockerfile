ARG PKTVISOR_TAG=latest-develop
FROM golang:1.17-alpine AS builder

WORKDIR /go/src/github.com/ns1labs/orb
COPY . .
RUN apk update && apk add make build-base git
RUN mkdir /tmp/build && CGO_ENABLED=1 make agent_bin && mv build/orb-agent /tmp/build/orb-agent

FROM ns1labs/pktvisor:${PKTVISOR_TAG}

COPY --from=builder /tmp/build/orb-agent /usr/local/bin/orb-agent
COPY --from=builder /go/src/github.com/ns1labs/orb/agent/docker/agent.yaml /etc/orb/agent.yaml
COPY --from=builder /go/src/github.com/ns1labs/orb/agent/docker/orb-agent-entry.sh /usr/local/bin/orb-agent-entry.sh

ENTRYPOINT [ "/usr/local/bin/orb-agent-entry.sh" ]
