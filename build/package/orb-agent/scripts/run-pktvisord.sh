#!/usr/bin/env bash

[[ "$VIZERD_ARGS" == "" ]] && VIZERD_ARGS="--admin-api"
/usr/local/bin/pktvisor_prometheus &
exec pktvisord $VIZERD_ARGS
