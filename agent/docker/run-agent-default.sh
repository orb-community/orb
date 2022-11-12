#!/bin/bash
# orb agent binary location. by default, matches orb-agent container (see Dockerfile)
orb_agent_bin="${ORB_AGENT_BIN:-/usr/local/bin/orb-agent}"
#
# check if debug mode is enabled
DEBUG=''
if [[ "$2" == '-d' ]]; then
  DEBUG='-d'
fi
if [ $# -eq 0 ]; then
  "$orb_agent_bin" run -c /opt/orb/agent_default.yaml &
  echo $! > /var/run/orb-agent.pid
else
  "$orb_agent_bin" run $DEBUG -c /opt/orb/agent_default.yaml &
  echo $! > /var/run/orb-agent.pid
fi
