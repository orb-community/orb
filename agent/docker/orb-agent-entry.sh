#!/usr/bin/env bash
#
# entry point for orb-agent
#

# orb agent binary location. by default, matches orb-agent container (see Dockerfile)
orb_agent_bin="${ORB_AGENT_BIN:-/usr/local/bin/orb-agent}"

# support generating API and MQTT addresses with one host name in ORB_CLOUD_ADDRESS
if [[ -n "${ORB_CLOUD_ADDRESS}" ]]; then
  ORB_CLOUD_API_ADDRESS="https://${ORB_CLOUD_ADDRESS}"
  ORB_CLOUD_MQTT_ADDRESS="tls://${ORB_CLOUD_ADDRESS}:8883"
  export ORB_CLOUD_API_ADDRESS ORB_CLOUD_MQTT_ADDRESS
fi

# support generating simple default pktvisor PCAP taps

tmpfile=$(mktemp /tmp/orb-agent-pktvisor-conf.XXXXXX)
trap 'rm -f "$tmpfile"' EXIT

# simplest: specify just interface, creates tap named "default_pcap"
# PKTVISOR_PCAP_IFACE_DEFAULT=en0
if [[ -n "${PKTVISOR_PCAP_IFACE_DEFAULT}" ]]; then
(
cat <<END
version: "1.0"

visor:
  taps:
    default_pcap:
      input_type: pcap
      config:
        iface: "$PKTVISOR_PCAP_IFACE_DEFAULT"
END
) >"$tmpfile"

  export ORB_BACKENDS_PKTVISOR_CONFIG_FILE="$tmpfile"
fi

# or specify pair of TAPNAME:IFACE
# TODO allow multiple, split on comma
# PKTVISOR_PCAP_IFACE_TAPS=default_pcap:en0

exec "$orb_agent_bin" "$@"
