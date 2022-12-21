#!/bin/bash
#
# entry point
#

stop1 () {
  printf "\rFinishing container.."
  exit 0
}

## Configuration ##
trap stop1 SIGINT
trap stop1 SIGTERM

# checking if log exist
if [ -d "/tmp/otel-collector.log" ]; then
    rm /tmp/otel-collector.log &
fi

# eternal loop
while true
do
  # pid file dont exist
  if [ ! -f "/var/run/collector.pid"  ]; then
    # running healthcheck to read logs to stdout
    ./healthcheck &
    # running in background
    nohup /usr/local/bin/otelcontribcol --config /etc/otelcol-contrib/config.yaml &>/tmp/otel-collector.log &
    # write pid on file
    echo $! > /var/run/collector.pid
    sleep 3
  else
    PID=$(cat /var/run/collector.pid)
    if [ ! -d "/proc/$PID" ]; then
       # if proccess is not running, stop container
       echo "$PID is not running"
       rm /var/run/collector.pid
       exit 1
    fi
    sleep 5
  fi
done
