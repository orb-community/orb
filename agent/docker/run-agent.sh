#!/bin/bash
# orb agent binary location. by default, matches orb-agent container (see Dockerfile)
orb_agent_bin="${ORB_AGENT_BIN:-/usr/local/bin/orb-agent}"
echo "Starting orb-agent : $orb_agent_bin with args $#"

if [ $# -eq 0 ]; then
  "$orb_agent_bin" run &
  echo $! > /var/run/orb-agent.pid
else
  "$orb_agent_bin" "$@" &
  echo $! > /var/run/orb-agent.pid
fi
