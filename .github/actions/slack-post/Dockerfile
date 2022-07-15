FROM ubuntu:20.04

LABEL author="Everton Taques<everton.taques@encora.com>"
LABEL maintainer="ns1labs"
LABEL version="1.0.0"

RUN mkdir -p /workspace

WORKDIR /workspace

COPY ./entrypoint.sh /entrypoint.sh

RUN apt-get update \
    && apt-get install jq git curl -y

ENTRYPOINT ["/entrypoint.sh"]