#!/usr/bin/env bash

[[ "$VIZERD_ARGS" == "" ]] && VIZERD_ARGS="--full-api"

exec pktvisord $VIZERD_ARGS
