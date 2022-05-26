FROM golang:1.15.1-alpine3.12

LABEL author="Everton Taques<everton.taques@encora.com>"
LABEL maintainer="ns1labs"
LABEL version="1.0.0"

RUN mkdir -p /workspace

WORKDIR /workspace

COPY ./entrypoint.sh /entrypoint.sh

RUN apk add git && \
wget https://github.com/cloudposse/github-commenter/releases/download/0.7.0/github-commenter_linux_amd64 -O /github-commenter && \
apk add jq && \
chmod +x /github-commenter /entrypoint.sh

CMD ["/entrypoint.sh"]
