ARG PKTVISOR_TAG=latest-develop
FROM ns1labs/pktvisor:${PKTVISOR_TAG} AS builder

WORKDIR /go/src/github.com/ns1labs/orb
COPY . .

ENV BUILD_DEPS "g++ make ca-certificates golang"

RUN apt-get update && \
    apt-get install --yes --force-yes --no-install-recommends ${BUILD_DEPS}
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -ldflags "-X 'github.com/ns1labs/orb/buildinfo.version=0.9.0-develop'" -o ../build/orb-agent ../cmd/agent/main.go
RUN mkdir /tmp/build && mv orb-agent /tmp/build/orb-agent

FROM ns1labs/pktvisor:${PKTVISOR_TAG}

COPY --from=builder /tmp/build/orb-agent /usr/local/bin/orb-agent
COPY --from=builder /go/src/github.com/ns1labs/orb/agent/docker/agent.yaml /etc/orb/agent.yaml
COPY --from=builder /go/src/github.com/ns1labs/orb/agent/docker/orb-agent-entry.sh /usr/local/bin/orb-agent-entry.sh

ENTRYPOINT [ "/usr/local/bin/orb-agent-entry.sh" ]
