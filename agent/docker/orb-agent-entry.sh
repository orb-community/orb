#!/usr/bin/env bash
#
# entry point for orb-agent
#

# -e ORB_CLOUD_ADDRESS=<prefill-host-name> \
# -e ORB_BACKENDS_PKTVISOR_IFACE=[ETH-INTERFACE] \

orb_agent_bin="${ORB_AGENT_BIN:-/usr/local/bin/orb-agent}"

if [[ -n "${ORB_CLOUD_ADDRESS}" ]]; then
  ORB_CLOUD_API_ADDRESS="https://${ORB_CLOUD_ADDRESS}"
  ORB_CLOUD_MQTT_ADDRESS="tls://${ORB_CLOUD_ADDRESS}:8883"
  export ORB_CLOUD_API_ADDRESS ORB_CLOUD_MQTT_ADDRESS
fi

exec "$orb_agent_bin" "$@"
