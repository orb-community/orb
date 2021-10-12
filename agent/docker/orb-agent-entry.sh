#!/usr/bin/env bash
#
# entry point for orb-agent
#

# -e ORB_CLOUD_ADDRESS=<prefill-host-name> \
# -e ORB_BACKENDS_PKTVISOR_IFACE=[ETH-INTERFACE] \

ORB_AGENT_BIN="/usr/local/bin/orb-agent"

exec $ORB_AGENT_BIN
