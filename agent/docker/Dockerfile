ARG PKTVISOR_TAG=develop
FROM debian:bullseye-slim AS builder

ENV BUILD_DEPS "g++ cmake make git pkgconf jq python3-pip python3-setuptools ca-certificates golang"

RUN \
    apt-get update && \
    apt-get upgrade --yes --force-yes && \
    apt-get install --yes --force-yes --no-install-recommends ${BUILD_DEPS} && \
    pip3 install conan

WORKDIR /
RUN \
    git clone https://github.com/ns1labs/pktvisor.git /pktvisor-src && \
    cd /pktvisor-src && git checkout ${PKTVISOR_TAG}

WORKDIR /tmp/build
RUN \
    conan profile new --detect default && \
    conan profile update settings.compiler.libcxx=libstdc++11 default && \
    conan config set general.revisions_enabled=1

RUN \
    PKG_CONFIG_PATH=/local/lib/pkgconfig cmake -DCMAKE_BUILD_TYPE=Release /pktvisor-src && \
    make all test -j 4

WORKDIR /go/src/github.com/ns1labs/orb
COPY . .

RUN CGO_ENABLED=1 make agent_bin && mv build/orb-agent /tmp/build/orb-agent

FROM debian:bullseye-slim AS runtime

ENV RUNTIME_DEPS "curl ca-certificates"

RUN \
    apt-get update && \
    apt-get upgrade --yes --force-yes && \
    apt-get install --yes --force-yes --no-install-recommends ${RUNTIME_DEPS} && \
    rm -rf /var/lib/apt

COPY --from=builder /tmp/build/bin/pktvisord /usr/local/sbin/pktvisord
COPY --from=builder /tmp/build/orb-agent /usr/local/bin/orb-agent
COPY --from=builder /go/src/github.com/ns1labs/orb/agent/docker/agent.yaml /etc/orb/agent.yaml

ENTRYPOINT [ "/usr/local/bin/orb-agent" ]
