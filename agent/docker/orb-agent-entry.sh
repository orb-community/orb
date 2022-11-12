#!/usr/bin/env bash
#
# entry point for orb-agent
#

agentstop1 () {
  printf "\rFinishing container.."
  exit 0
}

agentstop2 () {
  if [ -f "/var/run/orb-agent.pid"  ]; then
    ID=$(cat /var/run/orb-agent.pid)
    kill -15 $ID
  fi
}

# check geodb folder and extract db
cd /geo-db/
if [ -f "asn.mmdb.gz" ]; then
  gzip -d asn.mmdb.gz
  gzip -d city.mmdb.gz
fi
cd /
#

# orb agent binary location. by default, matches orb-agent container (see Dockerfile)
orb_agent_bin="${ORB_AGENT_BIN:-/usr/local/bin/orb-agent}"

## Cloud API Configuration ##
# support generating API and MQTT addresses with one host name in ORB_CLOUD_ADDRESS
if [[ -n "${ORB_CLOUD_ADDRESS}" ]]; then
  ORB_CLOUD_API_ADDRESS="https://${ORB_CLOUD_ADDRESS}"
  ORB_CLOUD_MQTT_ADDRESS="tls://${ORB_CLOUD_ADDRESS}:8883"
  export ORB_CLOUD_API_ADDRESS ORB_CLOUD_MQTT_ADDRESS
fi

## Agent Configuration ##
# support generating simple default pktvisor PCAP taps
tmpfile=$(mktemp /tmp/orb-agent-pktvisor-conf.XXXXXX)
trap 'rm -f "$tmpfile"' EXIT
trap agentstop1 SIGINT
trap agentstop2 SIGTERM

#Add default configuration to a file
(
cat <<END
version: "1.0"

visor:
  taps:
END
) > "$tmpfile"

# Check NetFlow TAP parameters
if [ "${PKTVISOR_NETFLOW_BIND_ADDRESS}" = '' ]; then
  PKTVISOR_NETFLOW_BIND_ADDRESS='0.0.0.0'
fi
if [ "${PKTVISOR_NETFLOW_PORT_DEFAULT}" = '' ]; then
  PKTVISOR_NETFLOW_PORT_DEFAULT='2055'
fi
if [ "${PKTVISOR_NETFLOW}" = 'true' ]; then
(
cat <<END
    default_netflow:
      input_type: flow
      config:
        flow_type: netflow
        port: "$PKTVISOR_NETFLOW_PORT_DEFAULT"
        bind: "$PKTVISOR_NETFLOW_BIND_ADDRESS"
END
) >> "$tmpfile"
  export ORB_BACKENDS_PKTVISOR_CONFIG_FILE="$tmpfile"
fi

# Check SFlow tap Parameters
if [ "${PKTVISOR_SFLOW_BIND_ADDRESS}" = '' ]; then
  PKTVISOR_SFLOW_BIND_ADDRESS='0.0.0.0'
fi
if [ "${PKTVISOR_SFLOW_PORT_DEFAULT}" = '' ]; then
  PKTVISOR_SFLOW_PORT_DEFAULT='6343'
fi
if [ "${PKTVISOR_SFLOW}" = 'true' ]; then
(
cat <<END
    default_sflow:
      input_type: flow
      config:
        flow_type: sflow
        port: "$PKTVISOR_SFLOW_PORT_DEFAULT"
        bind: "$PKTVISOR_SFLOW_BIND_ADDRESS"
END
) >> "$tmpfile"
  export ORB_BACKENDS_PKTVISOR_CONFIG_FILE="$tmpfile"
fi

# Check DNS TAP Parameters
if [ "${PKTVISOR_DNSTAP_BIND_ADDRESS}" = '' ]; then
  PKTVISOR_DNSTAP_BIND_ADDRESS='0.0.0.0'
fi
if [ "${PKTVISOR_DNSTAP_PORT_DEFAULT}" = '' ]; then
  PKTVISOR_DNSTAP_PORT_DEFAULT='6000'
fi
if [ "${PKTVISOR_DNSTAP}" = 'true' ]; then
(
cat <<END
    default_dnstap:
      input_type: dnstap
      config:
        tcp: "${PKTVISOR_DNSTAP_BIND_ADDRESS}:${PKTVISOR_DNSTAP_PORT_DEFAULT}"

END
) >> "$tmpfile"
  export ORB_BACKENDS_PKTVISOR_CONFIG_FILE="$tmpfile"
fi

# special case: if the iface is "mock", then use "mock" pcap source
if [ "$PKTVISOR_PCAP_IFACE_DEFAULT" = 'mock' ]; then
  MAYBE_MOCK='pcap_source: mock'
fi
if [[ -n "${PKTVISOR_PCAP_IFACE_DEFAULT}" || "${PKTVISOR_DNSTAP}" != 'true' && "${PKTVISOR_SFLOW}" != 'true' && "${PKTVISOR_NETFLOW}" != 'true' ]]; then
  if [ "$PKTVISOR_PCAP_IFACE_DEFAULT" = '' ]; then
    PKTVISOR_PCAP_IFACE_DEFAULT='auto'
  fi
(
cat <<END
    default_pcap:
      input_type: pcap
      config:
        iface: "$PKTVISOR_PCAP_IFACE_DEFAULT"
        $MAYBE_MOCK
END
) >>"$tmpfile"
  export ORB_BACKENDS_PKTVISOR_CONFIG_FILE="$tmpfile"
fi

# Checking agent.yaml config file
CONFIG_FILE_EXISTS=false
if [ -f "/opt/orb/agent.yaml" ]
then
  echo "Contains agent.yaml config file"
  CONFIG_FILE_EXISTS=true
else
  echo "Not contains agent.yaml config file, setting default file"
  CONFIG_FILE_EXISTS=false
  cp -rf /etc/orb/agent_default.yaml /opt/orb/agent_default.yaml
fi

# or specify pair of TAPNAME:IFACE
# TODO allow multiple, split on comma
# PKTVISOR_PCAP_IFACE_TAPS=default_pcap:en0
# eternal loop
while true
do
  # pid file dont exist
  if [ ! -f "/var/run/orb-agent.pid"  ]; then
    # running orb-agent in background
    if [[ "$2" == '-c' || "$3" == '-c' ]]; then
      # if config file was passed, drop the built-in pktvisor configuration file
      ORB_BACKENDS_PKTVISOR_CONFIG_FILE=""
      nohup /run-agent.sh "$@" &
    else
      if [[ $CONFIG_FILE_EXISTS ]]; then
        # if config file is mounted, drop the built-in pktvisor configuration file
        ORB_BACKENDS_PKTVISOR_CONFIG_FILE=""
        nohup /run-agent.sh "$@" &
      else
        # if none config file is set, use the built-in pktvisor configuration file and agent_default.yaml
        # check if debug mode is enabled
        DEBUG=''
        if [[ "$2" == '-d' ]]; then
          DEBUG='-d'
        fi
        nohup /run-agent.sh run $DEBUG -c /opt/orb/agent_default.yaml &
      fi
    fi
    sleep 2
    if [ -d "/nohup.out" ]; then
       tail -f /nohup.out &
    fi
  else
    PID=$(cat /var/run/orb-agent.pid)
    if [ ! -d "/proc/$PID" ]; then
       # stop container
       echo "$PID is not running"
       rm /var/run/orb-agent.pid
       exit 1
    fi
    sleep 5
  fi
done
