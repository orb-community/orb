ARG PKTVISOR_TAG=latest-develop
FROM golang:1.17-alpine AS builder

WORKDIR /go/src/github.com/ns1labs/orb
COPY go.mod .
RUN go mod tidy
COPY . .
RUN apk update && apk add make build-base git unzip wget

RUN arch=$(arch | sed s/aarch64/arm64/ ) && \
    cd /tmp && \
    wget -O cloudprober.zip "https://github.com/cloudprober/cloudprober/releases/download/v0.11.9/cloudprober-v0.11.9-linux-${arch}.zip" && \
    unzip -j cloudprober.zip

RUN mkdir /tmp/build && CGO_ENABLED=1 make agent_bin && mv build/orb-agent /tmp/build/orb-agent

FROM ns1labs/pktvisor:${PKTVISOR_TAG}

RUN mkdir /etc/cloudprober
COPY --from=builder /tmp/build/orb-agent /usr/local/bin/orb-agent
COPY --from=builder /go/src/github.com/ns1labs/orb/agent/docker/agent.yaml /etc/orb/agent.yaml
COPY --from=builder /go/src/github.com/ns1labs/orb/agent/docker/orb-agent-entry.sh /usr/local/bin/orb-agent-entry.sh
COPY --from=builder /go/src/github.com/ns1labs/orb/agent/docker/run-agent.sh /run-agent.sh
COPY --from=builder /tmp/cloudprober /usr/local/sbin/cloudprober

RUN chmod a+x /run-agent.sh

ENTRYPOINT [ "/usr/local/bin/orb-agent-entry.sh" ]
