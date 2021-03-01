FROM consul:1.9.3 as consul
FROM timberio/vector:0.11.X-distroless-static as vector
FROM ns1labs/pktvisor:develop

COPY --from=consul /usr/bin/telegraf /usr/local/bin/telegraf

RUN  apt-get update \
  && apt-get install -y runit-init dnsutils \
  && rm -rf /var/lib/apt \
  && mkdir -p /etc/runit/ \
  && mkdir -p /etc/telegraf/ \
  && rm -rf /etc/service/*

COPY misc/run-telegraf.sh /etc/service/telegraf/run
COPY misc/run-pktvisord.sh /etc/service/pktvisord/run
COPY misc/run-dig.sh /etc/service/dig/run
COPY misc/entrypoint /usr/local/bin/entrypoint

ENTRYPOINT /usr/local/bin/entrypoint
