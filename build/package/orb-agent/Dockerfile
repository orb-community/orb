#FROM consul:1.9.3 as consul
FROM timberio/vector:nightly-debian as vector
FROM ns1labs/orb-agent-slim:develop as orb-agent
FROM ns1labs/pktvisor-prometheus-export:latest as pktvisor-export
FROM ns1labs/pktvisor:develop

COPY --from=vector /usr/bin/vector /usr/local/bin/vector
COPY --from=orb-agent /usr/local/bin/orb-agent /usr/local/bin/orb-agent
COPY --from=pktvisor-export /usr/bin/pktvisor_exporter /usr/local/bin/pktvisor_prometheus
RUN chmod +x /usr/local/bin/pktvisor_prometheus

RUN  apt-get update \
  && apt-get install -y runit-init \
  && rm -rf /var/lib/apt \
  && mkdir -p /etc/runit/ \
  && mkdir -p /etc/vector/ \
  && mkdir /var/lib/vector \
  && rm -rf /etc/service/*

COPY files/vector.toml /etc/vector/vector.toml
COPY scripts/run-orb-agent.sh /etc/service/orb-agent/run
COPY scripts/run-vector.sh /etc/service/vector/run
COPY scripts/run-pktvisord.sh /etc/service/pktvisord/run
COPY scripts/entrypoint /usr/local/bin/entrypoint

ENTRYPOINT ["/usr/local/bin/entrypoint"]
